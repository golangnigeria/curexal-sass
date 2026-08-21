import React from "react";
import { CurexalLogoSymbol } from "@/components/brand/curexal-logo";
import { cn } from "@/lib/utils";

interface AppLoaderProps {
  message?: string;
  subMessage?: string;
  className?: string;
}

/**
 * AppLoader: Fullscreen branded application loader
 * Used for initial app mount, critical system boot, and full reload sequences
 */
export function AppLoader({
  message = "Loading Curexal Platform...",
  subMessage = "Initializing secure clinical workspace & encryption keys",
  className,
}: AppLoaderProps) {
  return (
    <div
      className={cn(
        "fixed inset-0 z-50 flex flex-col items-center justify-center bg-background/95 backdrop-blur-md transition-all",
        className
      )}
    >
      {/* Ambient Radial Glow */}
      <div className="absolute h-72 w-72 rounded-full bg-primary/15 blur-3xl -z-10 animate-pulse" />

      {/* Central Animated Logo Container */}
      <div className="relative flex items-center justify-center mb-6">
        {/* Outer Pulsing Glow Ring */}
        <div className="absolute -inset-3 rounded-2xl bg-gradient-to-tr from-emerald-500/30 via-primary/25 to-blue-600/30 blur-md animate-pulse" />

        {/* Orbit Ring */}
        <div className="absolute -inset-2.5 rounded-2xl border border-primary/30 animate-spin [animation-duration:8s]" />

        {/* Logo Shield Frame */}
        <div className="relative flex h-20 w-20 items-center justify-center rounded-2xl bg-slate-950 border border-border/80 shadow-2xl p-2.5">
          <CurexalLogoSymbol className="w-14 h-14 animate-bounce [animation-duration:2.5s]" />
        </div>
      </div>

      {/* Branded Text & Status */}
      <div className="flex flex-col items-center gap-1.5 text-center px-4 max-w-sm">
        <div className="flex items-center gap-2">
          <span className="font-extrabold tracking-tight text-lg text-foreground">CUREXAL</span>
          <span className="text-[11px] font-semibold uppercase tracking-wider text-emerald-400 bg-emerald-950/60 border border-emerald-800/50 px-2 py-0.5 rounded-full">
            Health OS
          </span>
        </div>
        <p className="text-sm font-medium text-foreground/90 animate-pulse mt-1">{message}</p>
        {subMessage && <p className="text-xs text-muted-foreground">{subMessage}</p>}
      </div>

      {/* Sleek Gradient Progress Indicator */}
      <div className="w-48 h-1 bg-muted rounded-full overflow-hidden mt-6">
        <div className="h-full bg-gradient-to-r from-emerald-400 via-primary to-blue-500 rounded-full animate-indeterminate" />
      </div>
    </div>
  );
}

/**
 * BootstrapLoader: Context and tenant resolution loader
 * Used while resolving user contexts, organization profiles, and capabilities
 */
export function BootstrapLoader({
  message = "Resolving platform authorization...",
  className,
}: {
  message?: string;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "flex h-screen w-screen flex-col items-center justify-center bg-background text-center px-4",
        className
      )}
    >
      <div className="relative flex items-center justify-center mb-5">
        <div className="absolute -inset-2 rounded-xl bg-primary/20 blur-sm animate-pulse" />
        <div className="relative flex h-14 w-14 items-center justify-center rounded-xl bg-slate-950 border border-border p-2 shadow-lg">
          <CurexalLogoSymbol className="w-10 h-10 animate-pulse" />
        </div>
      </div>
      <p className="text-sm font-semibold text-foreground tracking-tight">{message}</p>
      <p className="text-xs text-muted-foreground mt-1 animate-pulse">Syncing role contracts & entitlements</p>
    </div>
  );
}

/**
 * RouteLoader: Page / route transition loader
 * Used during route changes, lazy component loading, and view transitions
 */
export function RouteLoader({
  message = "Loading view...",
  className,
}: {
  message?: string;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "flex h-[60vh] w-full flex-col items-center justify-center gap-3 text-center p-6",
        className
      )}
    >
      <div className="relative flex items-center justify-center">
        <div className="h-10 w-10 animate-spin rounded-full border-2 border-primary/20 border-t-primary" />
        <div className="absolute inset-0 flex items-center justify-center">
          <CurexalLogoSymbol className="w-4 h-4 opacity-70" />
        </div>
      </div>
      <span className="text-xs font-medium text-muted-foreground animate-pulse">{message}</span>
    </div>
  );
}
