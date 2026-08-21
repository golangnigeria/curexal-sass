import React, { useEffect, useState } from "react";
import { useNavigation } from "react-router-dom";
import { cn } from "@/lib/utils";

interface TopLoadingBarProps {
  isLoading?: boolean;
  className?: string;
}

/**
 * TopLoadingBar: Sleek YouTube / NProgress style top loading bar
 * Automatically tracks route navigation state or manual query loading
 */
export function TopLoadingBar({ isLoading: manualLoading, className }: TopLoadingBarProps) {
  let isNavigating = false;
  try {
    const navigation = useNavigation();
    isNavigating = navigation.state === "loading";
  } catch {
    // Ignore outside router context
  }

  const isLoading = manualLoading ?? isNavigating;
  const [progress, setProgress] = useState(0);
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    let interval: any;
    if (isLoading) {
      setVisible(true);
      setProgress(15);
      interval = setInterval(() => {
        setProgress((prev) => {
          if (prev >= 85) return prev + 1.5;
          if (prev >= 60) return prev + 3;
          return prev + 12;
        });
      }, 150);
    } else {
      if (visible) {
        setProgress(100);
        const timeout = setTimeout(() => {
          setVisible(false);
          setProgress(0);
        }, 300);
        return () => clearTimeout(timeout);
      }
    }
    return () => clearInterval(interval);
  }, [isLoading, visible]);

  if (!visible && progress === 0) return null;

  return (
    <div
      className={cn(
        "fixed top-0 left-0 right-0 z-50 h-[3px] pointer-events-none overflow-hidden",
        className
      )}
    >
      <div
        className="h-full bg-gradient-to-r from-emerald-400 via-primary to-cyan-400 shadow-[0_0_12px_rgba(16,185,129,0.8)] transition-all duration-200 ease-out"
        style={{
          width: `${progress}%`,
          opacity: visible ? 1 : 0,
        }}
      />
    </div>
  );
}
