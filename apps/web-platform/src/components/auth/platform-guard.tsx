import React from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { authClient } from "@/lib/auth-client";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { BootstrapLoader } from "@/components/loading/app-loader";

import { isPlatformHost, getCanonicalPlatformUrl, getCanonicalOrgUrl } from "@/lib/url-builder";

export { BootstrapLoader as RouteLoadingScreen };

interface PlatformGuardProps {
  children?: React.ReactNode;
}

export function PlatformGuard({ children }: PlatformGuardProps) {
  const location = useLocation();
  const { data: session, isPending: isSessionPending } = authClient.useSession();
  const { data: bootstrap, isPending: isBootstrapPending } = useBootstrap();

  // 1. Wait for bootstrap and session resolution before making any routing decisions
  if (isSessionPending || (session?.user && isBootstrapPending && !bootstrap && !session.bootstrap)) {
    return <BootstrapLoader message="Authorizing platform console session..." />;
  }

  // 2. Unauthenticated -> redirect to /login
  if (!session || !session.user) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  const effectiveBootstrap = bootstrap || session.bootstrap;
  const user = session.user;

  // 3. Platform Authorization Resolution
  const isPlatformAuthorized =
    effectiveBootstrap?.platform?.isStaff === true ||
    user.isPlatformAdmin === true ||
    user.platformRole === "super_admin" ||
    user.role === "super_admin";

  if (isPlatformAuthorized) {
    // Automatically redirect bare localhost / tenant domains to app.localhost:5002 / app.curexal.space
    if (typeof window !== "undefined" && !isPlatformHost()) {
      const canonicalPlatformUrl = getCanonicalPlatformUrl(location.pathname + location.search);
      if (canonicalPlatformUrl !== location.pathname + location.search) {
        window.location.replace(canonicalPlatformUrl);
        return <BootstrapLoader message="Connecting to Platform Control Center (app.localhost)..." />;
      }
    }
    return children ? <>{children}</> : <Outlet />;
  }

  // 4. Authenticated Non-Platform User -> redirect to highest authorized context
  const isOrgAuthorized =
    Boolean(effectiveBootstrap?.organization?.id) ||
    effectiveBootstrap?.contexts?.current === "organization" ||
    user.role === "owner" ||
    user.role === "org_admin" ||
    user.role === "org_regional_manager" ||
    Boolean(user.organizationId);

  if (isOrgAuthorized) {
    return <Navigate to="/organization/dashboard" replace />;
  }

  const isWorkspaceAuthorized =
    Boolean(effectiveBootstrap?.workspace?.id) ||
    effectiveBootstrap?.contexts?.current === "workspace" ||
    Boolean(user.activeTenantId) ||
    Boolean(user.workspaceId);

  if (isWorkspaceAuthorized) {
    const activeBranchSlug =
      effectiveBootstrap?.branch?.slug ||
      effectiveBootstrap?.branch?.code?.toLowerCase() ||
      effectiveBootstrap?.availableBranches?.[0]?.slug ||
      effectiveBootstrap?.workspace?.slug ||
      "main";
    return <Navigate to={`/${activeBranchSlug}/dashboard`} replace />;
  }

  // 5. Fallback if no usable context is authorized
  return <Navigate to="/login" replace />;
}
