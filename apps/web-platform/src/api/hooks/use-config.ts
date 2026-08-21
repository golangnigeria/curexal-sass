import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPut, apiPatch } from "@/api/client";
import type { PlatformGeneralSettings, IdentitySecurityPolicy } from "@/api/contracts";

export function usePlatformConfig() {
  return useQuery({
    queryKey: ["platform", "config"],
    queryFn: () => apiGet<PlatformGeneralSettings>("/platform/config"),
  });
}

export function useUpdatePlatformConfig() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: PlatformGeneralSettings) =>
      apiPut<PlatformGeneralSettings>("/platform/config", payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "config"] });
    },
  });
}

export function useSecurityPolicy() {
  return useQuery({
    queryKey: ["platform", "security-policy"],
    queryFn: () => apiGet<IdentitySecurityPolicy>("/platform/security-policy"),
  });
}

export function useUpdateSecurityPolicy() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: IdentitySecurityPolicy) =>
      apiPut<IdentitySecurityPolicy>("/platform/security-policy", payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "security-policy"] });
    },
  });
}
