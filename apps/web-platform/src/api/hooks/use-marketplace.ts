import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiDelete } from "@/api/client";
import type {
  CapabilityCatalogItem,
  OrganizationEntitlement,
  EntitlementTrace,
} from "@/api/contracts";

export function useCapabilityCatalog(orgId?: string) {
  return useQuery({
    queryKey: ["marketplace", "capabilities", orgId || "all"],
    queryFn: () =>
      apiGet<CapabilityCatalogItem[]>(
        orgId
          ? `/organizations/${orgId}/marketplace/catalog`
          : "/marketplace/capabilities"
      ),
  });
}

export function useOrganizationCapabilities(orgId: string) {
  return useQuery({
    queryKey: ["organizations", orgId, "capabilities"],
    queryFn: () =>
      apiGet<{ organizationId: string; capabilities: string[] }>(
        `/organizations/${orgId}/capabilities`
      ),
    enabled: !!orgId,
  });
}

export function useOrganizationEntitlements(orgId: string) {
  return useQuery({
    queryKey: ["organizations", orgId, "entitlements"],
    queryFn: () =>
      apiGet<OrganizationEntitlement[]>(`/organizations/${orgId}/entitlements`),
    enabled: !!orgId,
  });
}

export function useEntitlementTrace(orgId: string, capabilityCode: string) {
  return useQuery({
    queryKey: ["organizations", orgId, "capabilities", capabilityCode, "trace"],
    queryFn: () =>
      apiGet<EntitlementTrace>(
        `/organizations/${orgId}/capabilities/trace?capability=${capabilityCode}`
      ),
    enabled: !!orgId && !!capabilityCode,
  });
}

export function useGrantCapability(orgId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { capabilityCode: string; source?: string; expiresAt?: string }) =>
      apiPost<{ message: string }>(`/platform/organizations/${orgId}/capabilities`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organizations", orgId] });
      queryClient.invalidateQueries({ queryKey: ["marketplace"] });
    },
  });
}

export function useRevokeCapability(orgId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (capabilityCode: string) =>
      apiDelete<{ message: string }>(
        `/platform/organizations/${orgId}/capabilities/${capabilityCode}`
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organizations", orgId] });
      queryClient.invalidateQueries({ queryKey: ["marketplace"] });
    },
  });
}

export function useStartTrialCapability(orgId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ capabilityCode, durationDays }: { capabilityCode: string; durationDays: number }) =>
      apiPost<{ message: string }>(
        `/platform/organizations/${orgId}/capabilities/${capabilityCode}/trial`,
        { durationDays }
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organizations", orgId] });
      queryClient.invalidateQueries({ queryKey: ["marketplace"] });
    },
  });
}
