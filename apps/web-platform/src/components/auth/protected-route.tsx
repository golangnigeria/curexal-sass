import React, { useEffect } from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { authClient } from "@/lib/auth-client";
import { AppLoader } from "@/components/loading/app-loader";

interface ProtectedRouteProps {
  requiredPermission?: string;
}

export function ProtectedRoute({ requiredPermission }: ProtectedRouteProps) {
  const location = useLocation();
  const { data: session, isPending } = authClient.useSession();

  if (isPending) {
    return <AppLoader message="Initializing Curexal Console..." />;
  }

  if (!session || !session.user) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // Check if authorized for platform or organization console context
  const currentContext = session.bootstrap?.contexts?.current;
  const isAuthorizedContext =
    currentContext === "platform" ||
    currentContext === "organization" ||
    session.user.isPlatformAdmin ||
    session.user.role === "owner" ||
    session.user.role === "org_admin" ||
    session.bootstrap?.organization?.id !== undefined;

  if (!isAuthorizedContext) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center p-6 text-center bg-background">
        <div className="max-w-md space-y-4 rounded-xl border border-destructive/20 bg-card p-8 shadow-card">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10 text-destructive">
            <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <h2 className="text-xl font-bold">Access Restricted</h2>
          <p className="text-sm text-muted-foreground">
            Your account ({session.user.email}) does not hold active administrator or organization privileges.
          </p>
          <button
            onClick={() => authClient.signOut()}
            className="w-full rounded-lg bg-primary py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
          >
            Sign In with Different Account
          </button>
        </div>
      </div>
    );
  }

  // If specific permission required, verify against user's permissions
  if (
    requiredPermission &&
    session.user.permissions &&
    !session.user.permissions.includes(requiredPermission) &&
    !session.user.isPlatformAdmin &&
    session.user.platformRole !== "super_admin"
  ) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center p-6 text-center bg-background">
        <div className="max-w-md space-y-4 rounded-xl border border-border bg-card p-8 shadow-card">
          <h2 className="text-lg font-bold">Permission Required</h2>
          <p className="text-xs text-muted-foreground">
            You require the <code className="font-mono text-primary">{requiredPermission}</code> permission to access this resource.
          </p>
        </div>
      </div>
    );
  }

  return <Outlet />;
}
