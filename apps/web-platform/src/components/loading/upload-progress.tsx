import React from "react";
import { cn } from "@/lib/utils";
import { FileUp, CheckCircle2, AlertCircle } from "lucide-react";
import { CurexalLogoSymbol } from "@/components/brand/curexal-logo";

interface UploadProgressProps {
  progress: number; // 0 to 100
  fileName?: string;
  fileSize?: string;
  status?: "uploading" | "encrypting" | "success" | "error";
  errorMessage?: string;
  className?: string;
}

/**
 * UploadProgress: Document and asset upload progress bar with Curexal encryption states
 */
export function UploadProgress({
  progress,
  fileName = "document.pdf",
  fileSize,
  status = "uploading",
  errorMessage,
  className,
}: UploadProgressProps) {
  const clampedProgress = Math.min(100, Math.max(0, progress));

  return (
    <div
      className={cn(
        "rounded-xl border border-border/80 bg-card p-4 shadow-sm transition-all duration-300 space-y-3",
        className
      )}
    >
      <div className="flex items-center justify-between gap-3">
        <div className="flex items-center gap-3 min-w-0">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary shrink-0">
            {status === "success" ? (
              <CheckCircle2 className="h-5 w-5 text-emerald-500" />
            ) : status === "error" ? (
              <AlertCircle className="h-5 w-5 text-destructive" />
            ) : (
              <FileUp className="h-5 w-5 animate-pulse" />
            )}
          </div>
          <div className="min-w-0 flex-1">
            <p className="text-sm font-semibold text-foreground truncate">{fileName}</p>
            <div className="flex items-center gap-2 text-xs text-muted-foreground mt-0.5">
              {fileSize && <span>{fileSize}</span>}
              {fileSize && <span>•</span>}
              <span className="capitalize">
                {status === "uploading" && "Uploading to S3..."}
                {status === "encrypting" && "Applying AES-GCM Enc..."}
                {status === "success" && "Uploaded & Verified"}
                {status === "error" && (errorMessage || "Upload failed")}
              </span>
            </div>
          </div>
        </div>

        {/* Status Percentage Badge */}
        <div className="flex items-center gap-2 shrink-0">
          {status === "uploading" && (
            <CurexalLogoSymbol className="w-4 h-4 opacity-70 animate-spin [animation-duration:3s]" />
          )}
          <span className="text-sm font-bold tabular-nums text-foreground">
            {clampedProgress}%
          </span>
        </div>
      </div>

      {/* Progress Track */}
      <div className="h-2 w-full overflow-hidden rounded-full bg-muted/60">
        <div
          className={cn(
            "h-full rounded-full transition-all duration-300 ease-out",
            status === "error"
              ? "bg-destructive"
              : status === "success"
              ? "bg-emerald-500"
              : "bg-gradient-to-r from-emerald-400 via-primary to-blue-500"
          )}
          style={{ width: `${clampedProgress}%` }}
        />
      </div>
    </div>
  );
}
