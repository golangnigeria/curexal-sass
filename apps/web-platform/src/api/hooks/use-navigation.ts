import { useQuery } from "@tanstack/react-query";
import { navigationService } from "@curexal/api-client";
import type { NavigationResponse, NavigationQueryParams } from "@curexal/contracts";

export interface UseNavigationOptions {
  scope?: "platform" | "organization" | "workspace" | "patient";
  branchSlug?: string;
}

export function useNavigation(scopeOrBranch?: string | UseNavigationOptions) {
  let resolvedScope: string | undefined;
  let resolvedBranchSlug: string | undefined;

  if (typeof scopeOrBranch === "object" && scopeOrBranch !== null) {
    resolvedScope = scopeOrBranch.scope;
    resolvedBranchSlug = scopeOrBranch.branchSlug;
  } else if (typeof scopeOrBranch === "string") {
    if (
      scopeOrBranch === "platform" ||
      scopeOrBranch === "organization" ||
      scopeOrBranch === "workspace" ||
      scopeOrBranch === "patient"
    ) {
      resolvedScope = scopeOrBranch;
    } else {
      resolvedBranchSlug = scopeOrBranch;
    }
  }

  // Infer scope and branch from location if not explicitly provided
  if (typeof window !== "undefined" && window.location?.pathname) {
    const pathname = window.location.pathname;
    const parts = pathname.split("/").filter(Boolean);

    if (!resolvedScope) {
      if (pathname.startsWith("/platform")) {
        resolvedScope = "platform";
      } else if (pathname.startsWith("/organization")) {
        resolvedScope = "organization";
      } else if (parts.length > 0 && !["login", "portal", "api"].includes(parts[0])) {
        resolvedScope = "workspace";
        if (!resolvedBranchSlug) {
          resolvedBranchSlug = parts[0];
        }
      }
    }
  }

  // Sanitize branch slug so route keywords are never sent as branch
  if (
    resolvedBranchSlug === "platform" ||
    resolvedBranchSlug === "organization" ||
    resolvedBranchSlug === "workspace" ||
    resolvedBranchSlug === "login" ||
    resolvedBranchSlug === "portal"
  ) {
    resolvedBranchSlug = undefined;
  }

  return useQuery<NavigationResponse>({
    queryKey: ["navigation", resolvedScope || "auto", resolvedBranchSlug || "default"],
    queryFn: () => {
      const host = typeof window !== "undefined" ? window.location.hostname : undefined;
      const params: NavigationQueryParams = {
        scope: resolvedScope,
        branch: resolvedBranchSlug,
        host,
      };
      return navigationService.getNavigation(params);
    },
    staleTime: 1000 * 60 * 5, // 5 minutes cache
    retry: 2,
  });
}
