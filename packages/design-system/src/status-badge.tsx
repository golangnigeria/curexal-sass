import React from "react";
import { cn } from "@curexal/utils";
import { resolveStatusMeta } from "./healthcare-status";

export interface StatusBadgeProps {
  status: string;
  label?: string;
  size?: "sm" | "md";
  showDot?: boolean;
  className?: string;
}

export const StatusBadge: React.FC<StatusBadgeProps> = ({
  status,
  label,
  size = "sm",
  showDot = true,
  className,
}) => {
  const meta = resolveStatusMeta(status);
  const displayLabel = label || meta.label;

  const variantStyles: Record<string, string> = {
    routine: "bg-teal-500/10 text-teal-700 dark:text-teal-300 border-teal-500/20",
    urgent: "bg-amber-500/10 text-amber-700 dark:text-amber-300 border-amber-500/20",
    emergency: "bg-rose-500/10 text-rose-700 dark:text-rose-300 border-rose-500/30 font-bold",
    critical: "bg-rose-950/40 text-rose-400 border-rose-600/40 font-bold animate-pulse",
    pending: "bg-slate-500/10 text-slate-700 dark:text-slate-300 border-slate-500/20",
    in_progress: "bg-sky-500/10 text-sky-700 dark:text-sky-300 border-sky-500/20",
    completed: "bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border-emerald-500/20 font-medium",
    cancelled: "bg-slate-500/10 text-slate-500 dark:text-slate-400 border-slate-500/20 line-through",
    rejected: "bg-rose-500/10 text-rose-700 dark:text-rose-400 border-rose-500/20",
  };

  const styleClass = variantStyles[meta.variant] || variantStyles.pending;

  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full border transition-colors",
        size === "sm" ? "px-2 py-0.5 text-[10px] font-mono" : "px-2.5 py-1 text-xs",
        styleClass,
        className
      )}
    >
      {showDot && (
        <span
          className={cn(
            "rounded-full shrink-0",
            size === "sm" ? "w-1.5 h-1.5" : "w-2 h-2",
            meta.dotColorClass || "bg-current"
          )}
        />
      )}
      <span className="truncate">{displayLabel}</span>
    </span>
  );
};
