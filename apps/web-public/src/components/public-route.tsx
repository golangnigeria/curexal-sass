import { authClient } from "@/lib/auth-client";
import { type ReactNode } from "react";
import { Navigate, useLocation } from "react-router-dom";
import { Activity } from "lucide-react";

interface PublicRouteProps {
  children: ReactNode;
  /** If true, redirect signed-in users to the dashboard (use for auth pages) */
  authOnly?: boolean;
}

export function PublicRoute({ children, authOnly = false }: PublicRouteProps) {
  const { data: session, isPending } = authClient.useSession();
  const location = useLocation();

  // Only block auth-specific pages (sign-in, sign-up) for signed-in users
  const isAuthPage = location.pathname.startsWith("/auth");

  if (isPending) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-background text-foreground gap-4">
        <div className="w-10 h-10 rounded-xl bg-primary/20 flex items-center justify-center animate-glow-pulse">
          <Activity className="h-5 w-5 text-primary" />
        </div>
      </div>
    );
  }

  if (session && isAuthPage) {
    return <Navigate to="/dashboard" replace />;
  }

  if (!session && isAuthPage) {
    return <>{children}</>;
  }

  return <>{children}</>;
}

