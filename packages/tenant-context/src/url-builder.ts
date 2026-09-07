/**
 * Canonical URL builder utilities for Curexal organization-centric routing.
 */

export const RESERVED_BRANCH_SLUGS = new Set([
  "organization",
  "platform",
  "api",
  "auth",
  "login",
  "register",
  "settings",
  "account",
  "profile",
  "notifications",
  "dashboard",
  "workspace",
  "admin",
  "public",
  "portal",
]);

/**
 * Builds an organization HQ path or full URL.
 */
export function buildOrgHqPath(subpath: string = "dashboard"): string {
  const cleanSubpath = subpath.startsWith("/") ? subpath.slice(1) : subpath;
  return `/organization/${cleanSubpath}`;
}

/**
 * Builds a branch-scoped workspace path (e.g. /owerri/laboratory).
 */
export function buildBranchWorkspacePath(
  branchSlug: string,
  workspaceType: string = "laboratory"
): string {
  const cleanBranch = (branchSlug || "main").toLowerCase().trim();
  const cleanWorkspace = (workspaceType || "laboratory").toLowerCase().trim();
  return `/${cleanBranch}/${cleanWorkspace}`;
}

/**
 * Checks if the current browser window is running on a custom domain vs a Curexal platform domain.
 */
export function isCurrentCustomDomain(): boolean {
  if (typeof window === "undefined") return false;
  const host = window.location.hostname.toLowerCase();
  return (
    !host.endsWith(".curexal.space") &&
    !host.endsWith(".curexal.internal") &&
    !host.endsWith(".localhost") &&
    host !== "localhost" &&
    host !== "127.0.0.1" &&
    host !== "app.curexal.space"
  );
}

/**
 * Checks if the current browser window is on the canonical Platform Control Plane host.
 */
export function isPlatformHost(): boolean {
  if (typeof window === "undefined") return false;
  const host = window.location.hostname.toLowerCase();
  return host === "app.localhost" || host === "app.curexal.space" || host === "app.curexal.internal";
}

/**
 * Checks if the current browser window is on the Patient Portal subdomain/host (patient.localhost or patient.curexal.space).
 */
export function isPatientHost(): boolean {
  if (typeof window === "undefined") return false;
  const host = window.location.hostname.toLowerCase();
  return (
    host === "patient.localhost" ||
    host === "patient.curexal.space" ||
    host === "patient.curexal.internal" ||
    host.startsWith("patient.")
  );
}

/**
 * Extracts organization subdomain/slug from the active window hostname.
 */
export function getOrgSlugFromHostname(): string | null {
  if (typeof window === "undefined") return null;
  const host = window.location.hostname.toLowerCase();
  if (isPlatformHost() || isPatientHost()) return null;
  if (host.endsWith(".localhost")) {
    const sub = host.replace(".localhost", "");
    return sub && sub !== "app" && sub !== "patient" ? sub : null;
  }
  if (host.endsWith(".curexal.space")) {
    const sub = host.replace(".curexal.space", "");
    return sub && sub !== "app" && sub !== "api" && sub !== "patient" ? sub : null;
  }
  return null;
}

/**
 * Returns the absolute canonical URL for the Patient Portal (patient.localhost:5002 or patient.curexal.space).
 */
export function getCanonicalPatientUrl(path: string = "/dashboard"): string {
  if (typeof window === "undefined") return path;
  const cleanPath = path.startsWith("/") ? path : `/${path}`;
  const host = window.location.hostname.toLowerCase();

  if (isPatientHost()) {
    return cleanPath;
  }

  if (host.endsWith("localhost") || host === "127.0.0.1") {
    return `http://patient.localhost:5003${cleanPath}`;
  }

  if (host.endsWith(".curexal.space") || host === "curexal.space") {
    return `https://patient.curexal.space${cleanPath}`;
  }

  return cleanPath;
}

/**
 * Returns the absolute canonical URL for the Platform Console (app.localhost:5002 or app.curexal.space).
 */
export function getCanonicalPlatformUrl(path: string = "/platform/dashboard"): string {
  if (typeof window === "undefined") return path;
  const cleanPath = path.startsWith("/") ? path : `/${path}`;
  const host = window.location.hostname.toLowerCase();
  const port = window.location.port;

  if (isPlatformHost()) {
    return cleanPath;
  }

  if (host.endsWith("localhost") || host === "127.0.0.1") {
    const targetPort = port ? `:${port}` : ":5002";
    return `http://app.localhost${targetPort}${cleanPath}`;
  }

  if (host.endsWith(".curexal.space") || host === "curexal.space") {
    return `https://app.curexal.space${cleanPath}`;
  }

  return cleanPath;
}

/**
 * Returns the absolute canonical URL for an organization domain.
 */
export function getCanonicalOrgUrl(orgSlug: string, path: string = "/organization/dashboard"): string {
  if (typeof window === "undefined" || !orgSlug) return path;
  const cleanPath = path.startsWith("/") ? path : `/${path}`;
  const host = window.location.hostname.toLowerCase();
  const port = window.location.port;

  const currentOrgSlug = getOrgSlugFromHostname();
  if (currentOrgSlug === orgSlug) {
    return cleanPath;
  }

  if (host.endsWith("localhost") || host === "127.0.0.1") {
    const targetPort = port ? `:${port}` : ":5002";
    return `http://${orgSlug}.localhost${targetPort}${cleanPath}`;
  }

  if (host.endsWith(".curexal.space") || host === "curexal.space" || host === "app.curexal.space") {
    return `https://${orgSlug}.curexal.space${cleanPath}`;
  }

  return cleanPath;
}
