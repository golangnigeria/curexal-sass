import React from "react";
import { CurexalLogoSymbol } from "@/components/brand/curexal-logo";
import { cn } from "@/lib/utils";

interface BrandPulseProps {
  size?: "sm" | "md" | "lg";
  text?: string;
  className?: string;
}

/**
 * BrandPulse: Modern medical radar / pulse loading effect with Curexal brand core
 * Ideal for card refreshes, background synchronization indicators, and modal confirmations
 */
export function BrandPulse({ size = "md", text, className }: BrandPulseProps) {
  const sizeMap = {
    sm: { container: "h-10 w-10", logo: "w-5 h-5", ripple: "h-12 w-12" },
    md: { container: "h-16 w-16", logo: "w-8 h-8", ripple: "h-20 w-20" },
    lg: { container: "h-24 w-24", logo: "w-12 h-12", ripple: "h-32 w-32" },
  };

  const selected = sizeMap[size];

  return (
    <div className={cn("flex flex-col items-center justify-center gap-3 p-4 select-none", className)}>
      <div className="relative flex items-center justify-center">
        {/* Radar Expanding Rings */}
        <div
          className={cn(
            "absolute rounded-full border-2 border-emerald-500/40 animate-ping opacity-75",
            selected.ripple
          )}
        />
        <div
          className={cn(
            "absolute rounded-full bg-emerald-500/10 animate-pulse [animation-duration:2s]",
            selected.ripple
          )}
        />

        {/* Shield Frame with Glowing Logo */}
        <div
          className={cn(
            "relative flex items-center justify-center rounded-2xl bg-slate-950 border border-emerald-500/30 shadow-lg shadow-emerald-900/20",
            selected.container
          )}
        >
          <CurexalLogoSymbol className={cn("animate-pulse [animation-duration:1.5s]", selected.logo)} />
        </div>
      </div>

      {text && (
        <span className="text-xs font-semibold tracking-tight text-muted-foreground animate-pulse">
          {text}
        </span>
      )}
    </div>
  );
}
