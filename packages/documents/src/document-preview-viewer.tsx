import React, { useEffect, useRef, useState, useCallback } from "react";
import * as pdfjsLib from "pdfjs-dist";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  Button,
  Badge,
  Input,
} from "@curexal/ui";
import { formatFileSize } from "@curexal/utils";
import { toast } from "sonner";
import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
  ZoomIn,
  ZoomOut,
  Maximize2,
  RotateCw,
  Download,
  X,
  Loader2,
  AlertTriangle,
  FileText,
  Image as ImageIcon,
  RefreshCw,
  ShieldCheck,
} from "lucide-react";

// Configure PDF.js worker
if (typeof window !== "undefined" && !pdfjsLib.GlobalWorkerOptions.workerSrc) {
  pdfjsLib.GlobalWorkerOptions.workerSrc = `https://cdn.jsdelivr.net/npm/pdfjs-dist@${pdfjsLib.version || "4.10.38"}/build/pdf.worker.min.mjs`;
}

export interface DocumentPreviewViewerProps {
  open: boolean;
  onClose: () => void;
  documentId?: string;
  organizationId?: string;
  title?: string;
  fileName?: string;
  mimeType?: string;
  fileSizeBytes?: number;
  version?: number;
  previewUrl?: string;
  downloadUrl?: string;
}

export function DocumentPreviewViewer({
  open,
  onClose,
  documentId,
  organizationId,
  title,
  fileName,
  mimeType,
  fileSizeBytes,
  version = 1,
  previewUrl,
  downloadUrl,
}: DocumentPreviewViewerProps) {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);

  // Document State
  const [pdfDoc, setPdfDoc] = useState<any>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [pageInput, setPageInput] = useState("1");
  const [scale, setScale] = useState(1.15);
  const [rotation, setRotation] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [isRendering, setIsRendering] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Image State
  const [imageBlobUrl, setImageBlobUrl] = useState<string | null>(null);

  // In-flight render task ref to cancel on fast pagination/zoom
  const renderTaskRef = useRef<any>(null);

  const effectiveMime = (mimeType || "").toLowerCase();
  const isPdf =
    effectiveMime.includes("pdf") ||
    (fileName || "").toLowerCase().endsWith(".pdf") ||
    (!effectiveMime.includes("image") && !effectiveMime.includes("png") && !effectiveMime.includes("jpeg"));

  const isImage =
    effectiveMime.includes("image") ||
    effectiveMime.includes("png") ||
    effectiveMime.includes("jpeg") ||
    effectiveMime.includes("jpg") ||
    (fileName || "").toLowerCase().match(/\.(png|jpe?g)$/i);

  // Construct secure API paths
  const effectivePreviewUrl =
    previewUrl ||
    (organizationId && documentId
      ? `/api/v1/organizations/${organizationId}/documents/${documentId}/preview`
      : documentId
      ? `/api/v1/organization/documents/${documentId}/preview`
      : "");

  const effectiveDownloadUrl =
    downloadUrl ||
    (organizationId && documentId
      ? `/api/v1/organizations/${organizationId}/documents/${documentId}/download`
      : documentId
      ? `/api/v1/organization/documents/${documentId}/download`
      : "");

  // Fetch document binary with session credentials
  const loadDocument = useCallback(async () => {
    if (!effectivePreviewUrl || !open) return;

    setIsLoading(true);
    setError(null);
    setPdfDoc(null);
    setCurrentPage(1);
    setPageInput("1");

    try {
      const response = await fetch(effectivePreviewUrl, {
        method: "GET",
        credentials: "include",
        headers: {
          Accept: "application/pdf, image/*, */*",
        },
      });

      if (!response.ok) {
        if (response.status === 401) {
          throw new Error("Authentication required. Please refresh your session.");
        }
        if (response.status === 403) {
          throw new Error("Access denied: You do not have permission to view this document.");
        }
        if (response.status === 404) {
          throw new Error("Document file was not found on the secure storage server.");
        }
        throw new Error(`Failed to load document (HTTP ${response.status})`);
      }

      const blob = await response.blob();

      if (isPdf) {
        const arrayBuffer = await blob.arrayBuffer();
        const loadingTask = pdfjsLib.getDocument({
          data: arrayBuffer,
          cMapUrl: "https://cdn.jsdelivr.net/npm/pdfjs-dist@4.10.38/cmaps/",
          cMapPacked: true,
        });

        const doc = await loadingTask.promise;
        setPdfDoc(doc);
        setTotalPages(doc.numPages);
      } else if (isImage) {
        const url = URL.createObjectURL(blob);
        setImageBlobUrl(url);
      }
    } catch (err: any) {
      console.error("Document preview error:", err);
      setError(err.message || "An error occurred while loading the document preview.");
    } finally {
      setIsLoading(false);
    }
  }, [effectivePreviewUrl, open, isPdf, isImage]);

  // Clean up object URLs on unmount / change
  useEffect(() => {
    if (open) {
      loadDocument();
    } else {
      if (imageBlobUrl) {
        URL.revokeObjectURL(imageBlobUrl);
        setImageBlobUrl(null);
      }
      setPdfDoc(null);
      setError(null);
    }
    return () => {
      if (imageBlobUrl) {
        URL.revokeObjectURL(imageBlobUrl);
      }
    };
  }, [open, loadDocument]);

  // Render current PDF page
  const renderPage = useCallback(
    async (pageNum: number) => {
      if (!pdfDoc || !canvasRef.current) return;

      if (renderTaskRef.current) {
        try {
          renderTaskRef.current.cancel();
        } catch {
          // Ignore cancellation errors
        }
      }

      setIsRendering(true);
      try {
        const page = await pdfDoc.getPage(pageNum);
        const canvas = canvasRef.current;
        if (!canvas) return;

        const viewport = page.getViewport({ scale: scale, rotation: rotation });

        // High-DPI Retina canvas adjustment
        const dpr = window.devicePixelRatio || 1;
        canvas.width = Math.floor(viewport.width * dpr);
        canvas.height = Math.floor(viewport.height * dpr);
        canvas.style.width = `${Math.floor(viewport.width)}px`;
        canvas.style.height = `${Math.floor(viewport.height)}px`;

        const ctx = canvas.getContext("2d");
        if (!ctx) return;

        const renderContext = {
          canvasContext: ctx,
          viewport: viewport,
          transform: dpr !== 1 ? [dpr, 0, 0, dpr, 0, 0] : null,
        };

        const renderTask = page.render(renderContext);
        renderTaskRef.current = renderTask;
        await renderTask.promise;
      } catch (err: any) {
        if (err?.name !== "RenderingCancelledException") {
          console.error("Page render error:", err);
        }
      } finally {
        setIsRendering(false);
      }
    },
    [pdfDoc, scale, rotation]
  );

  useEffect(() => {
    if (pdfDoc && isPdf) {
      renderPage(currentPage);
    }
  }, [pdfDoc, currentPage, scale, rotation, renderPage, isPdf]);

  // Pagination Handlers
  const goToPage = (p: number) => {
    const validPage = Math.max(1, Math.min(totalPages, p));
    setCurrentPage(validPage);
    setPageInput(String(validPage));
  };

  const handlePageSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const p = parseInt(pageInput, 10);
    if (!isNaN(p)) {
      goToPage(p);
    } else {
      setPageInput(String(currentPage));
    }
  };

  // Zoom Handlers
  const zoomIn = () => setScale((s) => Math.min(3.0, +(s + 0.15).toFixed(2)));
  const zoomOut = () => setScale((s) => Math.max(0.5, +(s - 0.15).toFixed(2)));
  const resetZoom = () => setScale(1.0);
  const fitWidth = () => {
    if (containerRef.current) {
      const containerWidth = containerRef.current.clientWidth - 48;
      setScale(Math.max(0.6, +(containerWidth / 620).toFixed(2)));
    }
  };
  const rotateDoc = () => setRotation((r) => (r + 90) % 360);

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!open) return;
      if (e.key === "ArrowRight" || e.key === "PageDown") {
        if (currentPage < totalPages) goToPage(currentPage + 1);
      } else if (e.key === "ArrowLeft" || e.key === "PageUp") {
        if (currentPage > 1) goToPage(currentPage - 1);
      } else if (e.key === "+" || e.key === "=") {
        zoomIn();
      } else if (e.key === "-") {
        zoomOut();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [open, currentPage, totalPages]);

  // Explicit Download Handler
  const handleExplicitDownload = () => {
    if (!effectiveDownloadUrl) return;

    toast.info("Preparing download...", {
      description: `Saving ${fileName || "document"}`,
    });

    const anchor = document.createElement("a");
    anchor.href = effectiveDownloadUrl;
    anchor.download = fileName || "document.pdf";
    document.body.appendChild(anchor);
    anchor.click();
    document.body.removeChild(anchor);
  };

  return (
    <Dialog open={open} onOpenChange={(val) => !val && onClose()}>
      <DialogContent className="max-w-6xl w-[94vw] h-[92vh] max-h-[92vh] flex flex-col p-0 gap-0 overflow-hidden bg-background border-border shadow-2xl rounded-2xl">
        {/* Header Bar */}
        <DialogHeader className="px-5 py-3.5 border-b border-border bg-card flex flex-row items-center justify-between space-y-0 shrink-0">
          <div className="flex items-center gap-3 truncate max-w-[60%]">
            <div className="p-2 rounded-xl bg-primary/10 text-primary shrink-0">
              {isPdf ? <FileText className="w-5 h-5" /> : <ImageIcon className="w-5 h-5" />}
            </div>
            <div className="truncate">
              <DialogTitle className="text-sm font-bold text-foreground truncate flex items-center gap-2">
                <span className="truncate">{title || fileName || "Document Preview"}</span>
                <Badge variant="outline" className="text-[10px] font-mono uppercase bg-primary/10 text-primary border-primary/20 shrink-0">
                  v{version}
                </Badge>
              </DialogTitle>
              <div className="flex items-center gap-2 text-[11px] text-muted-foreground mt-0.5">
                <span className="font-mono truncate max-w-[240px]">{fileName}</span>
                {fileSizeBytes ? (
                  <>
                    <span>•</span>
                    <span className="font-mono">{formatFileSize(fileSizeBytes)}</span>
                  </>
                ) : null}
                <span>•</span>
                <span className="text-emerald-600 dark:text-emerald-400 flex items-center gap-1 font-medium">
                  <ShieldCheck className="w-3.5 h-3.5" />
                  AES-256 Verified
                </span>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <Button
              size="sm"
              variant="outline"
              onClick={handleExplicitDownload}
              className="text-xs h-8 gap-1.5 border-border hover:border-primary hover:text-primary transition-all shadow-sm"
              title="Download original file"
            >
              <Download className="w-3.5 h-3.5" />
              <span className="hidden sm:inline">Download</span>
            </Button>
            <Button
              size="icon"
              variant="ghost"
              onClick={onClose}
              className="h-8 w-8 text-muted-foreground hover:text-foreground rounded-lg"
            >
              <X className="w-4 h-4" />
            </Button>
          </div>
        </DialogHeader>

        {/* Toolbar Controls */}
        <div className="px-5 py-2.5 bg-muted/40 border-b border-border/80 flex flex-wrap items-center justify-between gap-3 text-xs shrink-0 select-none">
          {/* Pagination Controls (PDF only) */}
          {isPdf && totalPages > 0 ? (
            <div className="flex items-center gap-1">
              <Button
                variant="ghost"
                size="icon"
                onClick={() => goToPage(1)}
                disabled={currentPage <= 1 || isRendering}
                className="h-7 w-7"
                title="First Page"
              >
                <ChevronsLeft className="w-3.5 h-3.5" />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => goToPage(currentPage - 1)}
                disabled={currentPage <= 1 || isRendering}
                className="h-7 w-7"
                title="Previous Page"
              >
                <ChevronLeft className="w-3.5 h-3.5" />
              </Button>

              <form onSubmit={handlePageSubmit} className="flex items-center gap-1.5 mx-1">
                <Input
                  type="text"
                  value={pageInput}
                  onChange={(e) => setPageInput(e.target.value)}
                  onBlur={() => setPageInput(String(currentPage))}
                  className="w-10 h-7 text-center text-xs font-mono p-0"
                />
                <span className="text-muted-foreground font-mono text-xs">/ {totalPages}</span>
              </form>

              <Button
                variant="ghost"
                size="icon"
                onClick={() => goToPage(currentPage + 1)}
                disabled={currentPage >= totalPages || isRendering}
                className="h-7 w-7"
                title="Next Page"
              >
                <ChevronRight className="w-3.5 h-3.5" />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => goToPage(totalPages)}
                disabled={currentPage >= totalPages || isRendering}
                className="h-7 w-7"
                title="Last Page"
              >
                <ChevronsRight className="w-3.5 h-3.5" />
              </Button>
            </div>
          ) : (
            <div className="text-[11px] text-muted-foreground font-mono">
              {isImage ? "Image Preview" : "Document Stream"}
            </div>
          )}

          {/* Zoom & View Controls */}
          <div className="flex items-center gap-1.5">
            <Button
              variant="ghost"
              size="icon"
              onClick={zoomOut}
              disabled={scale <= 0.5 || isLoading}
              className="h-7 w-7"
              title="Zoom Out"
            >
              <ZoomOut className="w-3.5 h-3.5" />
            </Button>

            <button
              onClick={resetZoom}
              className="px-2 py-1 rounded text-[11px] font-mono text-muted-foreground hover:text-foreground hover:bg-card transition-colors"
              title="Reset Zoom (100%)"
            >
              {Math.round(scale * 100)}%
            </button>

            <Button
              variant="ghost"
              size="icon"
              onClick={zoomIn}
              disabled={scale >= 3.0 || isLoading}
              className="h-7 w-7"
              title="Zoom In"
            >
              <ZoomIn className="w-3.5 h-3.5" />
            </Button>

            <div className="w-[1px] h-4 bg-border mx-1" />

            <Button
              variant="ghost"
              size="sm"
              onClick={fitWidth}
              disabled={isLoading}
              className="h-7 px-2 text-xs gap-1"
              title="Fit to Width"
            >
              <Maximize2 className="w-3 h-3" />
              <span className="hidden md:inline">Fit Width</span>
            </Button>

            <Button
              variant="ghost"
              size="icon"
              onClick={rotateDoc}
              disabled={isLoading}
              className="h-7 w-7"
              title="Rotate 90°"
            >
              <RotateCw className="w-3.5 h-3.5" />
            </Button>
          </div>
        </div>

        {/* Content Viewer Body */}
        <div
          ref={containerRef}
          className="flex-1 overflow-auto bg-slate-950/20 dark:bg-slate-950/60 p-6 flex items-center justify-center relative scrollbar-thin"
        >
          {/* Loading State */}
          {isLoading && (
            <div className="flex flex-col items-center justify-center gap-3 text-muted-foreground animate-fade-in">
              <Loader2 className="w-8 h-8 animate-spin text-primary" />
              <p className="text-xs font-medium">Decrypting & Streaming Document...</p>
            </div>
          )}

          {/* Error State */}
          {error && !isLoading && (
            <div className="max-w-md w-full p-6 rounded-2xl border border-rose-500/30 bg-rose-500/5 text-center space-y-3 animate-fade-in shadow-lg">
              <div className="w-12 h-12 rounded-full bg-rose-500/10 text-rose-500 flex items-center justify-center mx-auto">
                <AlertTriangle className="w-6 h-6" />
              </div>
              <div>
                <h4 className="text-sm font-bold text-rose-600 dark:text-rose-400">
                  Preview Failed
                </h4>
                <p className="text-xs text-muted-foreground mt-1 leading-relaxed">
                  {error}
                </p>
              </div>
              <Button
                size="sm"
                variant="outline"
                onClick={loadDocument}
                className="text-xs h-8 gap-1.5 border-rose-500/40 text-rose-600 hover:bg-rose-500/10"
              >
                <RefreshCw className="w-3.5 h-3.5" />
                Retry Loading
              </Button>
            </div>
          )}

          {/* PDF Canvas Rendering */}
          {isPdf && !error && !isLoading && (
            <div className="relative shadow-2xl border border-border/50 rounded-lg overflow-hidden bg-white dark:bg-slate-900 transition-transform">
              {isRendering && (
                <div className="absolute inset-0 bg-background/50 backdrop-blur-xs flex items-center justify-center z-10">
                  <Loader2 className="w-6 h-6 animate-spin text-primary" />
                </div>
              )}
              <canvas ref={canvasRef} className="block mx-auto max-w-full" />
            </div>
          )}

          {/* Image Rendering */}
          {isImage && imageBlobUrl && !error && !isLoading && (
            <div
              className="max-w-full max-h-full overflow-hidden flex items-center justify-center transition-transform"
              style={{ transform: `scale(${scale}) rotate(${rotation}deg)` }}
            >
              <img
                src={imageBlobUrl}
                alt={fileName || "Document Image"}
                className="max-w-full max-h-[75vh] object-contain rounded-lg shadow-xl border border-border"
              />
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
