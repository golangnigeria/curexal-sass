import { describe, it, expect } from "bun:test";

interface UserContext {
  id?: string;
  role?: string;
  platformRole?: string;
  isPlatformAdmin?: boolean;
  organizationId?: string;
  activeTenantId?: string;
  workspaceId?: string;
}

interface BootstrapContext {
  platform?: { isStaff?: boolean; role?: string };
  organization?: { id?: string; name?: string; role?: string };
  workspace?: { id?: string; name?: string };
  contexts?: { current?: string; available?: string[] };
}

interface SessionContext {
  user?: UserContext;
  bootstrap?: BootstrapContext;
}

/**
 * Canonical Destination and Authorization Resolver
 * Shared authorization logic used by RootRedirect, Login, and Route Guards
 */
export function resolveDestination(session: SessionContext | null, returnTo?: string | null): string {
  if (!session || !session.user) {
    return "/login";
  }

  const effectiveBootstrap = session.bootstrap;
  const user = session.user;

  const isPlatformStaff = Boolean(
    effectiveBootstrap?.platform?.isStaff === true ||
    user?.isPlatformAdmin === true ||
    user?.platformRole === "super_admin" ||
    user?.role === "super_admin"
  );

  const isOrgAuthorized = Boolean(
    isPlatformStaff ||
    Boolean(effectiveBootstrap?.organization?.id) ||
    effectiveBootstrap?.contexts?.current === "organization" ||
    user?.role === "owner" ||
    user?.role === "org_admin" ||
    user?.role === "org_regional_manager" ||
    Boolean(user?.organizationId)
  );

  const isWorkspaceAuthorized = Boolean(
    isOrgAuthorized ||
    Boolean(effectiveBootstrap?.workspace?.id) ||
    effectiveBootstrap?.contexts?.current === "workspace" ||
    Boolean(user?.activeTenantId) ||
    Boolean(user?.workspaceId)
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
  if (isWorkspaceAuthorized) return "/workspace/dashboard";

  return "/login";
}

/**
 * Route Gate Evaluator
 */
export function evaluateRouteAccess(
  route: string,
  session: SessionContext | null
): { allowed: boolean; redirectTo?: string } {
  if (!session || !session.user) {
    return { allowed: false, redirectTo: "/login" };
  }

  const effectiveBootstrap = session.bootstrap;
  const user = session.user;

  const isPlatformStaff = Boolean(
    effectiveBootstrap?.platform?.isStaff === true ||
    user?.isPlatformAdmin === true ||
    user?.platformRole === "super_admin" ||
    user?.role === "super_admin"
  );

  const isOrgAuthorized = Boolean(
    isPlatformStaff ||
    Boolean(effectiveBootstrap?.organization?.id) ||
    effectiveBootstrap?.contexts?.current === "organization" ||
    user?.role === "owner" ||
    user?.role === "org_admin" ||
    user?.role === "org_regional_manager" ||
    Boolean(user?.organizationId)
  );

  const isWorkspaceAuthorized = Boolean(
    isOrgAuthorized ||
    Boolean(effectiveBootstrap?.workspace?.id) ||
    effectiveBootstrap?.contexts?.current === "workspace" ||
    Boolean(user?.activeTenantId) ||
    Boolean(user?.workspaceId)
  );

  if (route.startsWith("/platform/")) {
    if (isPlatformStaff) return { allowed: true };
    if (isOrgAuthorized) return { allowed: false, redirectTo: "/organization/dashboard" };
    if (isWorkspaceAuthorized) return { allowed: false, redirectTo: "/workspace/dashboard" };
    return { allowed: false, redirectTo: "/login" };
  }

  if (route.startsWith("/organization/")) {
    if (isOrgAuthorized) return { allowed: true };
    if (isWorkspaceAuthorized) return { allowed: false, redirectTo: "/workspace/dashboard" };
    return { allowed: false, redirectTo: "/login" };
  }

  if (route.startsWith("/workspace/")) {
    if (isWorkspaceAuthorized) return { allowed: true };
    return { allowed: false, redirectTo: "/login" };
  }

  return { allowed: true };
}

describe("Role-Based Route Protection & Context Redirection Test Suite", () => {
  // Scenario A: Platform Admin
  it("Scenario A: Platform admin resolves to /platform/dashboard", () => {
    const session: SessionContext = {
      user: { id: "usr_1", isPlatformAdmin: true, role: "super_admin" },
      bootstrap: { platform: { isStaff: true, role: "super_admin" } },
    };
    expect(resolveDestination(session)).toBe("/platform/dashboard");
  });

  // Scenario B: Organization Owner
  it("Scenario B: Organization owner resolves to /organization/dashboard", () => {
    const session: SessionContext = {
      user: { id: "usr_2", isPlatformAdmin: false, role: "owner" },
      bootstrap: {
        platform: { isStaff: false, role: "" },
        organization: { id: "org_1", name: "Alpha Health", role: "owner" },
        contexts: { current: "organization" },
      },
    };
    expect(resolveDestination(session)).toBe("/organization/dashboard");
  });

  // Scenario C: Organization Member
  it("Scenario C: Organization member resolves to /organization/dashboard", () => {
    const session: SessionContext = {
      user: { id: "usr_3", isPlatformAdmin: false, role: "org_admin" },
      bootstrap: {
        platform: { isStaff: false, role: "" },
        organization: { id: "org_1", name: "Alpha Health", role: "org_admin" },
        contexts: { current: "organization" },
      },
    };
    expect(resolveDestination(session)).toBe("/organization/dashboard");
  });

  // Scenario D: Workspace User
  it("Scenario D: Workspace-only user resolves to /workspace/dashboard", () => {
    const session: SessionContext = {
      user: { id: "usr_4", isPlatformAdmin: false, role: "lab_technician", activeTenantId: "tenant_1" },
      bootstrap: {
        platform: { isStaff: false, role: "" },
        workspace: { id: "tenant_1", name: "Main Lab" },
        contexts: { current: "workspace" },
      },
    };
    expect(resolveDestination(session)).toBe("/workspace/dashboard");
  });

  // Scenario E: Unauthenticated
  it("Scenario E: Unauthenticated session resolves to /login", () => {
    expect(resolveDestination(null)).toBe("/login");
  });

  // Scenario G: Organization user enters /platform/dashboard
  it("Scenario G: Organization owner visiting /platform/dashboard is redirected to /organization/dashboard", () => {
    const session: SessionContext = {
      user: { id: "usr_2", isPlatformAdmin: false, role: "owner" },
      bootstrap: {
        platform: { isStaff: false, role: "" },
        organization: { id: "org_1", name: "Alpha Health", role: "owner" },
      },
    };
    const check = evaluateRouteAccess("/platform/dashboard", session);
    expect(check.allowed).toBe(false);
    expect(check.redirectTo).toBe("/organization/dashboard");
  });

  // Scenario H: Workspace user enters /platform/dashboard
  it("Scenario H: Workspace user visiting /platform/dashboard is redirected to /workspace/dashboard", () => {
    const session: SessionContext = {
      user: { id: "usr_4", isPlatformAdmin: false, role: "doctor", workspaceId: "ws_1" },
      bootstrap: {
        platform: { isStaff: false, role: "" },
        workspace: { id: "ws_1", name: "Clinical Unit" },
      },
    };
    const check = evaluateRouteAccess("/platform/dashboard", session);
    expect(check.allowed).toBe(false);
    expect(check.redirectTo).toBe("/workspace/dashboard");
  });

  // Scenario I: Organization user enters /organization/members
  it("Scenario I: Organization user visiting /organization/members is allowed", () => {
    const session: SessionContext = {
      user: { id: "usr_2", isPlatformAdmin: false, role: "owner" },
      bootstrap: {
        platform: { isStaff: false, role: "" },
        organization: { id: "org_1", name: "Alpha Health" },
      },
    };
    const check = evaluateRouteAccess("/organization/members", session);
    expect(check.allowed).toBe(true);
  });

  // Scenario J: Workspace-only user enters /organization/members
  it("Scenario J: Workspace-only user visiting /organization/members is redirected to /workspace/dashboard", () => {
    const session: SessionContext = {
      user: { id: "usr_4", isPlatformAdmin: false, role: "pharmacist", workspaceId: "ws_1" },
      bootstrap: {
        platform: { isStaff: false, role: "" },
        workspace: { id: "ws_1", name: "Pharmacy Branch" },
      },
    };
    const check = evaluateRouteAccess("/organization/members", session);
    expect(check.allowed).toBe(false);
    expect(check.redirectTo).toBe("/workspace/dashboard");
  });

  // Scenario K: Platform admin enters /platform/dashboard
  it("Scenario K: Platform admin visiting /platform/dashboard is allowed", () => {
    const session: SessionContext = {
      user: { id: "usr_1", isPlatformAdmin: true, role: "super_admin" },
      bootstrap: { platform: { isStaff: true, role: "super_admin" } },
    };
    const check = evaluateRouteAccess("/platform/dashboard", session);
    expect(check.allowed).toBe(true);
  });

  // Scenario L: Return URL attack prevention
  it("Scenario L: Return URL attack (Org user trying returnTo=/platform/dashboard) redirected to /organization/dashboard", () => {
    const session: SessionContext = {
      user: { id: "usr_2", isPlatformAdmin: false, role: "owner" },
      bootstrap: {
        platform: { isStaff: false, role: "" },
        organization: { id: "org_1", name: "Alpha Health" },
      },
    };
    const destination = resolveDestination(session, "/platform/dashboard");
    expect(destination).toBe("/organization/dashboard");
  });
});
