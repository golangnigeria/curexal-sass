import React from "react";
import { AlertTriangle, RefreshCw, Home } from "lucide-react";
import { Button } from "./button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "./card";

export interface RouteErrorBoundaryProps {
  error?: any;
  homePath?: string;
  homeLabel?: string;
}

export function RouteErrorBoundary({
  error: propError,
  homePath = "/",
  homeLabel = "Home",
}: RouteErrorBoundaryProps) {
  let errorMessage = "An unexpected error occurred.";
  let errorStatus = "Application Error";

  const error = propError || (typeof window !== "undefined" ? (window as any).__LAST_ERROR__ : null);

  if (error && typeof error === "object") {
    if ("status" in error) {
      errorStatus = `${error.status} ${error.statusText || ""}`;
      errorMessage = error.data?.message || error.statusText || errorMessage;
    } else if (error instanceof Error) {
      errorMessage = error.message;
    }
  } else if (typeof error === "string") {
    errorMessage = error;
  }

  return (
    <div className="min-h-[60vh] w-full flex items-center justify-center p-6 bg-background">
      <Card className="max-w-md w-full border-destructive/20 shadow-card card-enterprise">
        <CardHeader className="text-center pb-3">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-2xl bg-destructive/10 text-destructive mb-2">
            <AlertTriangle className="h-6 w-6" />
          </div>
          <CardTitle className="text-lg font-bold text-foreground">
            {errorStatus}
          </CardTitle>
          <CardDescription className="text-xs">
            {errorMessage}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-3 pt-2">
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => typeof window !== "undefined" && window.location.reload()}
              className="w-full text-xs h-8 gap-1.5"
            >
              <RefreshCw className="h-3.5 w-3.5" /> Reload Page
            </Button>
            <Button
              asChild
              size="sm"
              className="w-full bg-primary text-primary-foreground text-xs h-8 gap-1.5"
            >
              <a href={homePath}>
                <Home className="h-3.5 w-3.5" /> {homeLabel}
              </a>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
