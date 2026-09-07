import React from "react";
import { useBootstrap } from "./use-bootstrap";

interface BootstrapLoaderProps {
  children: React.ReactNode;
}

export const BootstrapLoader: React.FC<BootstrapLoaderProps> = ({ children }) => {
  const { isLoading } = useBootstrap();

  if (isLoading) {
    return (
      <div className="flex h-screen w-full items-center justify-center bg-background">
        <div className="flex flex-col items-center gap-3">
          <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          <p className="text-xs text-muted-foreground font-mono">Loading facility context & entitlements...</p>
        </div>
      </div>
    );
  }

  return <>{children}</>;
};
