import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost } from "@/api/client";
import type {
  DiagnosticsMetricsResponse,
  LaunchGateAudit,
  SystemHealthMetric,
} from "@/api/contracts";

export function useDiagnostics() {
  return useQuery({
    queryKey: ["platform", "diagnostics"],
    queryFn: () => apiGet<DiagnosticsMetricsResponse>("/platform/diagnostics"),
    refetchInterval: 30000, // Real-time refresh every 30 seconds
  });
}

export function useLaunchGateStatus() {
  return useQuery({
    queryKey: ["platform", "launch-gate", "status"],
    queryFn: () => apiGet<LaunchGateAudit>("/platform/launch-gate/status"),
  });
}

export function useHealthMetrics() {
  return useQuery({
    queryKey: ["platform", "health", "metrics"],
    queryFn: () => apiGet<SystemHealthMetric[]>("/platform/health/metrics"),
    refetchInterval: 60000,
  });
}

export function useVerifyLaunchGate() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => apiPost<LaunchGateAudit>("/platform/launch-gate/verify"),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "launch-gate"] });
      queryClient.invalidateQueries({ queryKey: ["platform", "diagnostics"] });
    },
  });
}
