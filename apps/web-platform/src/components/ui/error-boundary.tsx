import React from "react";
import { useRouteError, isRouteErrorResponse, Link } from "react-router-dom";
import { AlertTriangle, RefreshCw, Home } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export function RouteErrorBoundary() {
  const error = useRouteError();

  let errorMessage = "An unexpected error occurred.";
  let errorStatus = "Application Error";

  if (isRouteErrorResponse(error)) {
    errorStatus = `${error.status} ${error.statusText}`;
    errorMessage = error.data?.message || error.statusText;
  } else if (error instanceof Error) {
    errorMessage = error.message;
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
              onClick={() => window.location.reload()}
              className="w-full text-xs h-8 gap-1.5"
            >
              <RefreshCw className="h-3.5 w-3.5" /> Reload Page
            </Button>
            <Button
              asChild
              size="sm"
              className="w-full bg-primary text-primary-foreground text-xs h-8 gap-1.5"
            >
              <Link to="/platform/dashboard">
                <Home className="h-3.5 w-3.5" /> Control Center
              </Link>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
