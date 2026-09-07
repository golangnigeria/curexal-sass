import React from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { authClient } from "@/lib/auth-client";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { RouteLoadingScreen } from "@/components/auth/platform-guard";

interface OrganizationGuardProps {
  children?: React.ReactNode;
}

export function OrganizationGuard({ children }: OrganizationGuardProps) {
  const location = useLocation();
  const { data: session, isPending: isSessionPending } = authClient.useSession();
  const { data: bootstrap, isPending: isBootstrapPending } = useBootstrap();

  // 1. Wait for bootstrap and session resolution
  if (isSessionPending || (session?.user && isBootstrapPending && !bootstrap && !session.bootstrap)) {
    return <RouteLoadingScreen message="Resolving organization executive context..." />;
  }

  // 2. Unauthenticated -> redirect to /login
  if (!session || !session.user) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  const effectiveBootstrap = bootstrap || session.bootstrap;
  const user = session.user;

  // 3. Organization Executive Authorization Resolution
  const isPlatformStaff =
    effectiveBootstrap?.platform?.isStaff === true ||
    user.isPlatformAdmin === true ||
    user.platformRole === "super_admin" ||
    user.role === "super_admin";

  const userOrgRole = (effectiveBootstrap?.organization?.role || user.organizationRole || user.role || "").toLowerCase();

  const isOrgAuthorized =
    isPlatformStaff ||
    userOrgRole === "owner" ||
    userOrgRole === "org_admin" ||
    userOrgRole === "org_regional_manager" ||
    userOrgRole === "org_quality_manager" ||
    userOrgRole === "org_finance_manager" ||
    userOrgRole === "org_hr_manager" ||
    userOrgRole === "admin" ||
    (effectiveBootstrap?.contexts?.current === "organization" && effectiveBootstrap?.contexts?.available?.includes("organization"));

  if (isOrgAuthorized) {
    return children ? <>{children}</> : <Outlet />;
  }

  // 4. Branch staff (e.g. Doctor, Nurse, MLS) -> redirect to branch clinical/operational workspace
  const isWorkspaceAuthorized =
    Boolean(effectiveBootstrap?.branch?.id) ||
    Boolean(effectiveBootstrap?.workspace?.id) ||
    Boolean(effectiveBootstrap?.availableBranches?.length) ||
    effectiveBootstrap?.contexts?.current === "workspace" ||
    Boolean(user.activeTenantId) ||
    Boolean(user.workspaceId) ||
    Boolean(effectiveBootstrap?.organization?.id);

  if (isWorkspaceAuthorized) {
    const activeBranchSlug =
      effectiveBootstrap?.branch?.slug ||
      effectiveBootstrap?.branch?.code?.toLowerCase() ||
      effectiveBootstrap?.availableBranches?.[0]?.slug ||
      effectiveBootstrap?.workspace?.slug ||
      "main";
    const userRoleStr = (user.role || (effectiveBootstrap?.identity as any)?.role || userOrgRole).toLowerCase();
    let mod = "dashboard";
    if (userRoleStr.includes("doc") || userRoleStr.includes("clin") || userRoleStr.includes("phys")) {
      mod = "clinical";
    } else if (userRoleStr.includes("nurse") || userRoleStr.includes("triage")) {
      mod = "care-desk";
    } else if (userRoleStr.includes("recept") || userRoleStr.includes("front")) {
      mod = "reception";
    } else if (userRoleStr.includes("cash") || userRoleStr.includes("acc") || userRoleStr.includes("bill")) {
      mod = "billing";
    } else {
      mod = "dashboard";
    }
    return <Navigate to={`/${activeBranchSlug}/${mod}`} replace />;
  }

  // 5. No usable context
  return <Navigate to="/login" replace />;
}
