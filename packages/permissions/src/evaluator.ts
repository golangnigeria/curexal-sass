import {
  ROLES,
  CONTEXT_SCOPES,
  type UserRoleResponse,
  type BootstrapContractResponse,
} from "@curexal/contracts";

export interface SessionAuthContext {
  user?: Partial<UserRoleResponse>;
  bootstrap?: Partial<BootstrapContractResponse>;
}

/**
 * Check if the user has a specific permission string.
 */
export function can(
  permission: string,
  user: Partial<UserRoleResponse> | null | undefined
): boolean {
  if (!user) return false;
  if (user.isPlatformAdmin || user.platformRole === ROLES.SUPER_ADMIN || user.role === ROLES.SUPER_ADMIN) {
    return true;
  }
  return Boolean(user.permissions && user.permissions.includes(permission));
}

/**
 * Check if the user has a specific role.
 */
export function hasRole(
  role: string,
  user: Partial<UserRoleResponse> | null | undefined
): boolean {
  if (!user) return false;
  const target = role.toLowerCase();
  return (
    (user.role || "").toLowerCase() === target ||
    (user.platformRole || "").toLowerCase() === target ||
    (user.organizationRole || "").toLowerCase() === target
  );
}

/**
 * Check if an organization has a specific licensed/enabled capability.
 */
export function hasCapability(
  capabilityCode: string,
  bootstrap: Partial<BootstrapContractResponse> | null | undefined
): boolean {
  if (!bootstrap) return false;
  const capabilities = new Set(bootstrap.capabilities || []);
  return capabilities.has(capabilityCode);
}

/**
 * Canonical Destination and Authorization Resolver
 * Shared authorization logic used by RootRedirect, Login, and Route Guards
 */
export function resolveDestination(
  session: SessionAuthContext | null,
  returnTo?: string | null
): string {
  if (!session || !session.user) {
    return "/login";
  }

  const effectiveBootstrap = session.bootstrap;
  const user = session.user;

  const isPlatformStaff = Boolean(
    effectiveBootstrap?.platform?.isStaff === true ||
      user?.isPlatformAdmin === true ||
      user?.platformRole === ROLES.SUPER_ADMIN ||
      user?.role === ROLES.SUPER_ADMIN
  );

  const userOrgRole = (effectiveBootstrap?.organization?.role || user?.organizationRole || user?.role || "").toLowerCase();
  const isOrgAuthorized = Boolean(
    isPlatformStaff ||
      userOrgRole === "owner" ||
      userOrgRole === "org_admin" ||
      userOrgRole === "org_regional_manager" ||
      userOrgRole === "org_quality_manager" ||
      userOrgRole === "org_finance_manager" ||
      userOrgRole === "org_hr_manager" ||
      userOrgRole === "admin" ||
      (effectiveBootstrap?.contexts?.current === CONTEXT_SCOPES.ORGANIZATION &&
        effectiveBootstrap?.contexts?.available?.includes(CONTEXT_SCOPES.ORGANIZATION))
  );

  const isWorkspaceAuthorized = Boolean(
    Boolean(effectiveBootstrap?.workspace?.id) ||
      Boolean(effectiveBootstrap?.branch?.id) ||
      effectiveBootstrap?.contexts?.current === CONTEXT_SCOPES.WORKSPACE ||
      Boolean(user?.activeTenantId) ||
      Boolean(user?.workspaceId) ||
      !isOrgAuthorized
  );

  // Validate returnTo if present (prevent open redirect & privilege escalation)
  if (returnTo && returnTo.startsWith("/") && !returnTo.startsWith("//")) {
    if (returnTo.startsWith("/platform/")) {
      if (isPlatformStaff) return returnTo;
    } else if (returnTo.startsWith("/organization/")) {
      if (isOrgAuthorized) return returnTo;
    } else if (returnTo.startsWith("/workspace/")) {
      if (isWorkspaceAuthorized) return returnTo;
    }
  }

  // Canonical destination resolution
  if (isPlatformStaff) return "/platform/dashboard";
  if (isOrgAuthorized) return "/organization/dashboard";
  if (isWorkspaceAuthorized) {
    const bSlug =
      effectiveBootstrap?.branch?.slug ||
      effectiveBootstrap?.branch?.code?.toLowerCase() ||
      effectiveBootstrap?.workspace?.slug ||
      "main";
    const userRoleStr = (user?.role || userOrgRole).toLowerCase();
    const mod =
      userRoleStr.includes("doc") || userRoleStr.includes("clin")
        ? "clinical"
        : userRoleStr.includes("nurse") || userRoleStr.includes("triage")
        ? "triage"
        : userRoleStr.includes("recept") || userRoleStr.includes("front")
        ? "appointments"
        : userRoleStr.includes("lab") || userRoleStr.includes("scien")
        ? "laboratory"
        : userRoleStr.includes("pharm")
        ? "pharmacy"
        : userRoleStr.includes("rad")
        ? "radiology"
        : userRoleStr.includes("cash") || userRoleStr.includes("bill")
        ? "billing"
        : "dashboard";
    return `/${bSlug}/${mod}`;
  }

  return "/login";
}

/**
 * Route Gate Evaluator
 */
export function evaluateRouteAccess(
  route: string,
  session: SessionAuthContext | null
): { allowed: boolean; redirectTo?: string } {
  if (!session || !session.user) {
    return { allowed: false, redirectTo: "/login" };
  }

  const effectiveBootstrap = session.bootstrap;
  const user = session.user;

  const isPlatformStaff = Boolean(
    effectiveBootstrap?.platform?.isStaff === true ||
      user?.isPlatformAdmin === true ||
      user?.platformRole === ROLES.SUPER_ADMIN ||
      user?.role === ROLES.SUPER_ADMIN
  );

  const isOrgAuthorized = Boolean(
    isPlatformStaff ||
      Boolean(effectiveBootstrap?.organization?.id) ||
      effectiveBootstrap?.contexts?.current === CONTEXT_SCOPES.ORGANIZATION ||
      user?.role === ROLES.OWNER ||
      user?.role === ROLES.ORG_ADMIN ||
      user?.role === ROLES.ORG_REGIONAL_MANAGER ||
      Boolean(user?.organizationId)
  );

  const isWorkspaceAuthorized = Boolean(
    isOrgAuthorized ||
      Boolean(effectiveBootstrap?.workspace?.id) ||
      effectiveBootstrap?.contexts?.current === CONTEXT_SCOPES.WORKSPACE ||
      Boolean(user?.activeTenantId) ||
      Boolean(user?.workspaceId)
  );

  if (route.startsWith("/platform")) {
    if (isPlatformStaff) return { allowed: true };
    if (isOrgAuthorized) return { allowed: false, redirectTo: "/organization/dashboard" };
    if (isWorkspaceAuthorized) {
      const bSlug = effectiveBootstrap?.branch?.slug || "main";
      return { allowed: false, redirectTo: `/${bSlug}/dashboard` };
    }
    return { allowed: false, redirectTo: "/login" };
  }

  if (route.startsWith("/organization")) {
    if (isOrgAuthorized) return { allowed: true };
    if (isWorkspaceAuthorized) {
      const bSlug = effectiveBootstrap?.branch?.slug || "main";
      return { allowed: false, redirectTo: `/${bSlug}/dashboard` };
    }
    return { allowed: false, redirectTo: "/login" };
  }

  if (route.startsWith("/workspace") || route.includes("/")) {
    if (isWorkspaceAuthorized) return { allowed: true };
    return { allowed: false, redirectTo: "/login" };
  }

  return { allowed: true };
}
