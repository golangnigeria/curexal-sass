import React from "react";
import { cn } from "@/lib/utils";
import { CurexalLogoSymbol } from "@/components/brand/curexal-logo";

interface ButtonSpinnerProps {
  size?: "sm" | "md" | "lg";
  withLogo?: boolean;
  className?: string;
}

/**
 * ButtonSpinner: Interactive button loading spinner with optional mini Curexal brand node
 */
export function ButtonSpinner({
  size = "sm",
  withLogo = false,
  className,
}: ButtonSpinnerProps) {
  const sizeMap = {
    sm: "h-4 w-4 border-2",
    md: "h-5 w-5 border-2",
    lg: "h-6 w-6 border-[2.5px]",
  };

  const logoSizeMap = {
    sm: "w-2.5 h-2.5",
    md: "w-3 h-3",
    lg: "w-3.5 h-3.5",
  };

  return (
    <span className={cn("inline-flex items-center justify-center relative shrink-0", className)}>
      <span
        className={cn(
          "animate-spin rounded-full border-current border-t-transparent",
          sizeMap[size]
        )}
      />
      {withLogo && (
        <span className="absolute inset-0 flex items-center justify-center">
          <CurexalLogoSymbol className={cn("opacity-80", logoSizeMap[size])} />
        </span>
      )}
    </span>
  );
}
