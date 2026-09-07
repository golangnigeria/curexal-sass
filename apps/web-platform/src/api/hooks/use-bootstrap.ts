import { useQuery } from "@tanstack/react-query";
import { apiGet } from "@/api/client";
import type { BootstrapContractResponse } from "@/api/contracts";

export function useBootstrap(branchSlug?: string) {
  let resolvedBranchSlug = branchSlug;
  if (!resolvedBranchSlug && typeof window !== "undefined" && window.location?.pathname) {
    const parts = window.location.pathname.split("/").filter(Boolean);
    if (
      parts.length > 0 &&
      parts[0] !== "login" &&
      parts[0] !== "platform" &&
      parts[0] !== "organization" &&
      parts[0] !== "workspace" &&
      parts[0] !== "api"
    ) {
      resolvedBranchSlug = parts[0];
    }
  }

  return useQuery({
    queryKey: ["bootstrap", resolvedBranchSlug || ""],
    queryFn: () => {
      const params = new URLSearchParams();
      if (resolvedBranchSlug) params.set("branch", resolvedBranchSlug);
      if (typeof window !== "undefined" && window.location.hostname) {
        params.set("host", window.location.hostname);
      }
      const qs = params.toString();
      return apiGet<BootstrapContractResponse>(`/bootstrap${qs ? `?${qs}` : ""}`);
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
    retry: 1,
  });
}
