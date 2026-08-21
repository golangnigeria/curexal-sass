import React from "react";
import { Badge } from "@/components/ui/badge";

export type StatusType =
  | "active"
  | "pending"
  | "completed"
  | "cancelled"
  | "normal"
  | "critical"
  | "high"
  | "low"
  | "stat"
  | "routine"
  | "paid"
  | "overdue";

interface StatusPillProps {
  status: StatusType | string;
  label?: string;
  className?: string;
}

const statusStyles: Record<string, { bg: string; text: string; dot: string }> = {
  active: { bg: "bg-emerald-500/10 border-emerald-500/30", text: "text-emerald-700 dark:text-emerald-400", dot: "bg-emerald-500" },
  completed: { bg: "bg-emerald-500/10 border-emerald-500/30", text: "text-emerald-700 dark:text-emerald-400", dot: "bg-emerald-500" },
  paid: { bg: "bg-emerald-500/10 border-emerald-500/30", text: "text-emerald-700 dark:text-emerald-400", dot: "bg-emerald-500" },
  normal: { bg: "bg-emerald-500/10 border-emerald-500/30", text: "text-emerald-700 dark:text-emerald-400", dot: "bg-emerald-500" },
  
  pending: { bg: "bg-amber-500/10 border-amber-500/30", text: "text-amber-700 dark:text-amber-400", dot: "bg-amber-500" },
  routine: { bg: "bg-sky-500/10 border-sky-500/30", text: "text-sky-700 dark:text-sky-400", dot: "bg-sky-500" },
  
  critical: { bg: "bg-rose-500/10 border-rose-500/30", text: "text-rose-700 dark:text-rose-400", dot: "bg-rose-500" },
  high: { bg: "bg-rose-500/10 border-rose-500/30", text: "text-rose-700 dark:text-rose-400", dot: "bg-rose-500" },
  stat: { bg: "bg-rose-500/10 border-rose-500/30", text: "text-rose-700 dark:text-rose-400", dot: "bg-rose-500" },
  overdue: { bg: "bg-rose-500/10 border-rose-500/30", text: "text-rose-700 dark:text-rose-400", dot: "bg-rose-500" },
  cancelled: { bg: "bg-slate-500/10 border-slate-500/30", text: "text-slate-700 dark:text-slate-400", dot: "bg-slate-500" },
};

export function StatusPill({ status, label, className }: StatusPillProps) {
  const normalized = status.toLowerCase();
  const config = statusStyles[normalized] || {
    bg: "bg-secondary/50 border-border",
    text: "text-muted-foreground",
    dot: "bg-muted-foreground",
  };

  return (
    <Badge
      variant="outline"
      className={`inline-flex items-center gap-1.5 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider ${config.bg} ${config.text} ${className || ""}`}
    >
      <span className={`w-1.5 h-1.5 rounded-full ${config.dot}`} />
      {label || status.replace("_", " ")}
    </Badge>
  );
}
