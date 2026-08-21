import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { LucideIcon, TrendingUp, TrendingDown } from "lucide-react";

interface StatCardProps {
  title: string;
  value: string | number;
  icon: LucideIcon;
  iconColorClass?: string;
  trendPercentage?: number;
  trendLabel?: string;
  isLoading?: boolean;
}

export function StatCard({
  title,
  value,
  icon: Icon,
  iconColorClass = "text-primary bg-primary/10",
  trendPercentage,
  trendLabel = "vs last period",
  isLoading,
}: StatCardProps) {
  const isPositive = (trendPercentage ?? 0) >= 0;

  return (
    <Card className="border-border shadow-sm hover:shadow-md transition-shadow">
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
          {title}
        </CardTitle>
        <div className={`p-2 rounded-lg ${iconColorClass}`}>
          <Icon className="w-4 h-4" />
        </div>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold text-foreground">
          {isLoading ? <span className="animate-pulse">...</span> : value}
        </div>
        {trendPercentage !== undefined && (
          <p
            className={`text-[11px] flex items-center gap-1 font-medium mt-1 ${
              isPositive ? "text-emerald-600 dark:text-emerald-400" : "text-rose-600 dark:text-rose-400"
            }`}
          >
            {isPositive ? <TrendingUp className="w-3 h-3" /> : <TrendingDown className="w-3 h-3" />}
            {isPositive ? `+${trendPercentage}%` : `${trendPercentage}%`} {trendLabel}
          </p>
        )}
      </CardContent>
    </Card>
  );
}
