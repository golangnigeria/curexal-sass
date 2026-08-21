import React from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

/**
 * PageSkeleton: Full page structure skeleton
 * Used as a placeholder for full dashboard layouts, analytics screens, and configuration pages
 */
export function PageSkeleton({
  hasMetrics = true,
  cardCount = 4,
  className,
}: {
  hasMetrics?: boolean;
  cardCount?: number;
  className?: string;
}) {
  return (
    <div className={cn("space-y-6 p-6 max-w-7xl mx-auto w-full animate-fade-in", className)}>
      {/* Header Skeleton */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 border-b border-border/40 pb-5">
        <div className="space-y-2">
          <Skeleton className="h-8 w-56 rounded-lg" />
          <Skeleton className="h-4 w-80 rounded-md" />
        </div>
        <div className="flex items-center gap-3">
          <Skeleton className="h-9 w-24 rounded-lg" />
          <Skeleton className="h-9 w-32 rounded-lg" />
        </div>
      </div>

      {/* Metric Cards Skeleton */}
      {hasMetrics && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {Array.from({ length: cardCount }).map((_, i) => (
            <CardSkeleton key={i} />
          ))}
        </div>
      )}

      {/* Content Area Skeleton */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-4">
          <TableSkeleton rows={5} columns={4} />
        </div>
        <div className="space-y-4">
          <div className="rounded-xl border border-border/60 bg-card p-5 space-y-4">
            <Skeleton className="h-5 w-36 rounded-md" />
            <Skeleton className="h-32 w-full rounded-lg" />
            <div className="space-y-2 pt-2">
              <Skeleton className="h-4 w-full rounded-md" />
              <Skeleton className="h-4 w-4/5 rounded-md" />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

/**
 * TableSkeleton: Tabular data loading placeholder
 * Used for user tables, organization lists, audit logs, and catalog grids
 */
export function TableSkeleton({
  rows = 5,
  columns = 5,
  showHeader = true,
  className,
}: {
  rows?: number;
  columns?: number;
  showHeader?: boolean;
  className?: string;
}) {
  return (
    <div className={cn("rounded-xl border border-border/60 bg-card overflow-hidden shadow-sm", className)}>
      {showHeader && (
        <div className="flex items-center justify-between p-4 border-b border-border/40 bg-muted/20">
          <div className="flex items-center gap-3">
            <Skeleton className="h-5 w-32 rounded-md" />
            <Skeleton className="h-4 w-12 rounded-full" />
          </div>
          <div className="flex items-center gap-2">
            <Skeleton className="h-8 w-48 rounded-lg" />
            <Skeleton className="h-8 w-20 rounded-lg" />
          </div>
        </div>
      )}

      {/* Table Header Row */}
      <div className="grid grid-cols-12 gap-4 px-6 py-3 border-b border-border/40 bg-muted/40 text-xs font-semibold">
        {Array.from({ length: columns }).map((_, i) => (
          <div key={i} className={cn("flex items-center", i === 0 ? "col-span-3" : "col-span-2")}>
            <Skeleton className="h-3.5 w-20 rounded-md" />
          </div>
        ))}
      </div>

      {/* Table Rows */}
      <div className="divide-y divide-border/30">
        {Array.from({ length: rows }).map((_, r) => (
          <div key={r} className="grid grid-cols-12 gap-4 px-6 py-4 items-center hover:bg-muted/10 transition-colors">
            {Array.from({ length: columns }).map((_, c) => (
              <div key={c} className={cn("flex items-center gap-3", c === 0 ? "col-span-3" : "col-span-2")}>
                {c === 0 && <Skeleton className="h-8 w-8 rounded-full shrink-0" />}
                <Skeleton className={cn("h-4 rounded-md", c === 0 ? "w-28" : "w-20")} />
              </div>
            ))}
          </div>
        ))}
      </div>
    </div>
  );
}

/**
 * CardSkeleton: Metric / KPI Card placeholder
 */
export function CardSkeleton({ className }: { className?: string }) {
  return (
    <div className={cn("rounded-xl border border-border/60 bg-card p-5 space-y-3 shadow-sm", className)}>
      <div className="flex items-center justify-between">
        <Skeleton className="h-4 w-24 rounded-md" />
        <Skeleton className="h-8 w-8 rounded-lg" />
      </div>
      <Skeleton className="h-7 w-20 rounded-md" />
      <Skeleton className="h-3.5 w-36 rounded-md" />
    </div>
  );
}
