import React from "react";
import { Navigate, Outlet, useLocation, useParams } from "react-router-dom";
import { authClient } from "@/lib/auth-client";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { RouteLoadingScreen } from "@/components/auth/platform-guard";

interface WorkspaceGuardProps {
  children?: React.ReactNode;
}

export function WorkspaceGuard({ children }: WorkspaceGuardProps) {
  const location = useLocation();
  const { branchSlug } = useParams<{ branchSlug: string }>();
  const { data: session, isPending: isSessionPending } = authClient.useSession();
  const { data: bootstrap, isPending: isBootstrapPending } = useBootstrap(branchSlug);

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

  // 3. Organization Executive Resolution
  const isPlatformStaff =
    effectiveBootstrap?.platform?.isStaff === true ||
    user.isPlatformAdmin === true ||
    user.platformRole === "super_admin" ||
    user.role === "super_admin";

  const userOrgRole = (effectiveBootstrap?.organization?.role || user.organizationRole || user.role || "").toLowerCase();
  const isOrgExecutive =
    isPlatformStaff ||
    userOrgRole === "owner" ||
    userOrgRole === "org_admin" ||
    userOrgRole === "org_regional_manager" ||
    userOrgRole === "org_quality_manager" ||
    userOrgRole === "org_finance_manager" ||
    userOrgRole === "org_hr_manager" ||
    userOrgRole === "admin";

  // 4. Branch Verification for Branch-only Staff
  const availableBranches = effectiveBootstrap?.availableBranches || [];
  const activeBranchSlug =
    effectiveBootstrap?.branch?.slug ||
    effectiveBootstrap?.branch?.code?.toLowerCase() ||
    availableBranches[0]?.slug ||
    availableBranches[0]?.code?.toLowerCase() ||
    effectiveBootstrap?.workspace?.slug ||
    "main";

  if (!isOrgExecutive && branchSlug && availableBranches.length > 0) {
    const isAssigned = availableBranches.some(
      (b: any) =>
        b.slug?.toLowerCase() === branchSlug.toLowerCase() ||
        b.code?.toLowerCase() === branchSlug.toLowerCase() ||
        b.id === branchSlug
    );
    if (!isAssigned) {
      // Unauthorized cross-branch attempt -> redirect safely to assigned branch
      return <Navigate to={`/${activeBranchSlug}/dashboard`} replace />;
    }
  }

  const isWorkspaceAuthorized =
    isOrgExecutive ||
    Boolean(effectiveBootstrap?.workspace?.id) ||
    Boolean(effectiveBootstrap?.branch?.id) ||
    availableBranches.length > 0 ||
    effectiveBootstrap?.contexts?.current === "workspace" ||
    Boolean(user.activeTenantId) ||
    Boolean(user.workspaceId);

  if (isWorkspaceAuthorized) {
    return children ? <>{children}</> : <Outlet />;
  }

  // 5. No usable context
  return <Navigate to="/login" replace />;
}
