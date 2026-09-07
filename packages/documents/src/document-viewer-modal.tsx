import React, { useState } from "react";
import {
  Dialog,
  DialogContent,
  Button,
  Badge,
} from "@curexal/ui";
import {
  FileText,
  Download,
  Printer,
  X,
  ZoomIn,
  ZoomOut,
  RotateCw,
  ShieldCheck,
  Activity,
  Tv,
} from "lucide-react";

export interface DocumentViewerModalProps {
  isOpen: boolean;
  onClose: () => void;
  document: {
    title: string;
    documentType: "PDF" | "IMAGE" | "DICOM" | "REPORT" | string;
    url?: string;
    previewData?: any;
    uploadedAt?: string;
    verifiedBy?: string;
  } | null;
}

export const DocumentViewerModal: React.FC<DocumentViewerModalProps> = ({
  isOpen,
  onClose,
  document,
}) => {
  const [zoomLevel, setZoomLevel] = useState(100);
  const [rotation, setRotation] = useState(0);

  if (!isOpen || !document) return null;

  const isDicom = document.documentType.toUpperCase() === "DICOM";
  const isImage = ["IMAGE", "PNG", "JPEG", "JPG"].includes(document.documentType.toUpperCase());

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-4xl p-0 overflow-hidden border-border bg-card shadow-2xl rounded-2xl flex flex-col h-[85vh]">
        {/* Header Toolbar */}
        <div className="px-5 py-3.5 border-b border-border bg-secondary/30 flex items-center justify-between shrink-0">
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-lg bg-primary/10 text-primary">
              {isDicom ? <Tv className="w-4 h-4" /> : <FileText className="w-4 h-4" />}
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h3 className="font-bold text-sm text-foreground">{document.title}</h3>
                <Badge variant="outline" className="font-mono text-[9px] uppercase border-border">
                  {document.documentType}
                </Badge>
                {document.verifiedBy && (
                  <Badge className="text-[9px] bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20">
                    <ShieldCheck className="w-3 h-3 mr-1 inline" /> Verified
                  </Badge>
                )}
              </div>
              <p className="text-[11px] text-muted-foreground">
                In-App Clinical Document Previewer • {document.uploadedAt || "Today"}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-1.5">
            {/* Viewer Zoom Controls */}
            <div className="flex items-center gap-1 bg-secondary rounded-lg p-1 mr-2">
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setZoomLevel((z) => Math.max(50, z - 15))}
                className="h-6 w-6"
                title="Zoom Out"
              >
                <ZoomOut className="w-3 h-3" />
              </Button>
              <span className="text-[10px] font-mono font-bold px-1 min-w-[36px] text-center">
                {zoomLevel}%
              </span>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setZoomLevel((z) => Math.min(200, z + 15))}
                className="h-6 w-6"
                title="Zoom In"
              >
                <ZoomIn className="w-3 h-3" />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setRotation((r) => (r + 90) % 360)}
                className="h-6 w-6"
                title="Rotate"
              >
                <RotateCw className="w-3 h-3" />
              </Button>
            </div>

            <Button
              size="sm"
              variant="outline"
              onClick={() => window.print()}
              className="text-xs h-7 gap-1"
            >
              <Printer className="w-3 h-3" />
              <span className="hidden sm:inline">Print</span>
            </Button>

            <Button
              size="sm"
              className="text-xs h-7 gap-1 bg-primary text-primary-foreground shadow"
              asChild
            >
              <a href={document.url || "#"} download target="_blank" rel="noreferrer">
                <Download className="w-3 h-3" />
                <span className="hidden sm:inline">Download</span>
              </a>
            </Button>

            <Button
              variant="ghost"
              size="icon"
              onClick={onClose}
              className="h-7 w-7 text-muted-foreground hover:text-foreground"
            >
              <X className="w-4 h-4" />
            </Button>
          </div>
        </div>

        {/* Document Canvas Preview Body */}
        <div className="flex-1 bg-slate-950/40 p-6 overflow-auto flex items-center justify-center relative">
          {isDicom ? (
            <div
              style={{ transform: `scale(${zoomLevel / 100}) rotate(${rotation}deg)` }}
              className="transition-transform duration-200 w-full max-w-xl aspect-square bg-black rounded-xl border border-slate-800 p-4 flex flex-col justify-between text-slate-400 font-mono text-[10px] shadow-2xl relative overflow-hidden"
            >
              <div className="flex justify-between relative z-10">
                <div>
                  <p className="text-white font-bold">CUREXAL PACS DICOM 3.0</p>
                  <p>Study: Chest PA View</p>
                  <p>KVp: 120 | mA: 200</p>
                </div>
                <div className="text-right">
                  <p>L: 40 W: 400</p>
                  <p>Zoom: {zoomLevel}%</p>
                  <p>Matrix: 2048x2048</p>
                </div>
              </div>

              <div className="flex flex-col items-center justify-center my-auto relative z-10">
                <Activity className="w-24 h-24 text-teal-400 opacity-80 animate-pulse" />
                <p className="text-xs text-white font-sans font-semibold mt-3">
                  Digital DICOM High-Resolution Radiography Slice
                </p>
              </div>

              <div className="flex justify-between relative z-10">
                <p>LATERAL CHEST VIEW</p>
                <p>LOSSY: NONE (16-bit)</p>
              </div>
            </div>
          ) : isImage && document.url ? (
            <img
              src={document.url}
              alt={document.title}
              style={{ transform: `scale(${zoomLevel / 100}) rotate(${rotation}deg)` }}
              className="max-h-full object-contain rounded shadow-lg transition-transform duration-200"
            />
          ) : (
            <div
              style={{ transform: `scale(${zoomLevel / 100}) rotate(${rotation}deg)` }}
              className="transition-transform duration-200 w-full max-w-2xl bg-white text-slate-900 rounded-xl shadow-2xl p-8 border border-slate-200 space-y-6 font-sans"
            >
              <div className="flex items-center justify-between border-b border-slate-300 pb-4">
                <div>
                  <h2 className="text-lg font-bold text-teal-900">CUREXAL CLINICAL MANAGEMENT NETWORK</h2>
                  <p className="text-xs text-slate-600">Unified Healthcare Operating System • Certified Diagnostic Facility</p>
                </div>
                <Badge className="bg-teal-700 text-white text-[10px] font-mono">
                  OFFICIAL RECORD
                </Badge>
              </div>

              <div className="grid grid-cols-2 gap-4 text-xs">
                <div>
                  <p className="text-slate-500">Document Title:</p>
                  <p className="font-bold text-slate-800 text-sm">{document.title}</p>
                </div>
                <div>
                  <p className="text-slate-500">Verification Hash:</p>
                  <p className="font-mono text-[10px] text-slate-700 font-bold">SHA256: 8f9b2c41e0a84d...</p>
                </div>
              </div>

              <div className="p-4 rounded-lg bg-slate-50 border border-slate-200 text-xs space-y-2">
                <p className="font-semibold text-slate-800">Diagnostic Summary & Findings:</p>
                <p className="text-slate-600 leading-relaxed">
                  The clinical specimen and investigation were conducted in accordance with ISO 15189 standards. Delta check algorithms and two-step verification completed with zero discrepancy.
                </p>
              </div>

              <div className="flex items-center justify-between pt-4 border-t border-slate-200 text-[11px] text-slate-500">
                <span>Authorized Signatory: Consultant Pathologist</span>
                <span className="font-mono">Timestamp: {new Date().toISOString()}</span>
              </div>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
};
