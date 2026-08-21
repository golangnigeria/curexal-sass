import React from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { authClient } from "@/lib/auth-client";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { RouteLoadingScreen } from "@/components/auth/platform-guard";

interface WorkspaceGuardProps {
  children?: React.ReactNode;
}

export function WorkspaceGuard({ children }: WorkspaceGuardProps) {
  const location = useLocation();
  const { data: session, isPending: isSessionPending } = authClient.useSession();
  const { data: bootstrap, isPending: isBootstrapPending } = useBootstrap();

  // 1. Wait for bootstrap and session resolution
  if (isSessionPending || (session?.user && isBootstrapPending && !bootstrap && !session.bootstrap)) {
    return <RouteLoadingScreen message="Connecting to clinical workspace..." />;
  }

  // 2. Unauthenticated -> redirect to /login
  if (!session || !session.user) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  const effectiveBootstrap = bootstrap || session.bootstrap;
  const user = session.user;

  // 3. Workspace Authorization Resolution
  const isPlatformStaff =
    effectiveBootstrap?.platform?.isStaff === true ||
    user.isPlatformAdmin === true ||
    user.platformRole === "super_admin" ||
    user.role === "super_admin";

  const isOrgAuthorized =
    isPlatformStaff ||
    Boolean(effectiveBootstrap?.organization?.id) ||
    effectiveBootstrap?.contexts?.current === "organization" ||
    user.role === "owner" ||
    user.role === "org_admin" ||
    user.role === "org_regional_manager" ||
    Boolean(user.organizationId);

  const isWorkspaceAuthorized =
    isOrgAuthorized ||
    Boolean(effectiveBootstrap?.workspace?.id) ||
    effectiveBootstrap?.contexts?.current === "workspace" ||
    Boolean(user.activeTenantId) ||
    Boolean(user.workspaceId);

  if (isWorkspaceAuthorized) {
    return children ? <>{children}</> : <Outlet />;
  }

  // 4. No usable context
  return <Navigate to="/login" replace />;
}
