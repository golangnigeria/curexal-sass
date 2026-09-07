import type { BootstrapContractResponse, NavigationItemPayload } from "@curexal/contracts";

export interface WorkspaceNavItem {
  id: string;
  title: string;
  path: string;
  icon: string;
  moduleCode: string;
  requiredCapability?: string;
  badge?: string;
  description?: string;
}

/**
 * @deprecated Database-driven navigation is the single source of truth.
 * Kept for backwards compatibility with legacy imports.
 */
export const WORKSPACE_MODULE_REGISTRY: WorkspaceNavItem[] = [];

/**
 * Resolves workspace navigation items directly from authoritative backend bootstrap navigation.
 * Pure presentation transform with ZERO role hardcoding or frontend authorization guessing.
 */
export function resolveEntitledWorkspaces(
  bootstrap: BootstrapContractResponse | null | undefined,
  _user?: any
): WorkspaceNavItem[] {
  const navItems = bootstrap?.navigation || bootstrap?.structuredNavigation?.primary || [];
  return navItems.map((item: NavigationItemPayload) => {
    const segments = item.path.split("/").filter(Boolean);
    const subpath = segments.length > 1 ? segments[1] : segments[0] || "dashboard";
    return {
      id: item.id,
      title: item.title,
      path: subpath,
      icon: item.icon,
      moduleCode: subpath,
    };
  });
}
