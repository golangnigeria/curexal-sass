import type { ReactNode } from "react";
import { FolderOpen, ServerOff, SearchX } from "lucide-react";
import { Button }  from "./button";

interface EmptyStateProps {
  title: string;
  description: string;
  icon?: "folder" | "search" | "server";
  actionLabel?: string;
  onAction?: () => void;
  badge?: string;
}

export function EmptyState({
  title,
  description,
  icon = "folder",
  actionLabel,
  onAction,
  badge,
}: EmptyStateProps) {
  const IconComponent =
    icon === "search" ? SearchX : icon === "server" ? ServerOff : FolderOpen;

  return (
    <div className="flex flex-col items-center justify-center text-center p-8 md:p-12 rounded-2xl bg-slate-50/50 dark:bg-slate-900/50 border border-slate-200/80 dark:border-slate-800 my-4 max-w-2xl mx-auto">
      {badge && (
        <span className="inline-flex items-center gap-1.5 px-3 py-1 text-xs font-semibold rounded-full bg-teal-50 text-[#0F766E] border border-teal-200/60 dark:bg-teal-950/40 dark:text-teal-400 dark:border-teal-800/60 mb-4">
          {badge}
        </span>
      )}
      <div className="w-14 h-14 rounded-2xl bg-teal-50 dark:bg-teal-950/60 flex items-center justify-center text-[#0F766E] dark:text-teal-400 border border-teal-100 dark:border-teal-900 mb-4 shadow-sm">
        <IconComponent className="w-7 h-7" />
      </div>
      <h3 className="text-xl font-bold text-slate-900 dark:text-white tracking-tight mb-2">
        {title}
      </h3>
      <p className="text-sm text-slate-600 dark:text-slate-400 max-w-md mb-6 leading-relaxed">
        {description}
      </p>
      {actionLabel && onAction && (
        <Button
          onClick={onAction}
          className="bg-[#0F766E] hover:bg-[#115E59] text-white font-medium shadow-sm transition-all rounded-xl"
        >
          {actionLabel}
        </Button>
      )}
    </div>
  );
}
