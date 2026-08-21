import { useQuery } from "@tanstack/react-query";
import { apiGet } from "@/api/client";
import type { AuditLog, ListAuditLogsPayload, AuditAdminStats } from "@/api/contracts";

export function usePlatformAuditLogs(filters: ListAuditLogsPayload = {}) {
  const queryParams = new URLSearchParams();
  if (filters.limit) queryParams.append("limit", filters.limit.toString());
  if (filters.offset) queryParams.append("offset", filters.offset.toString());
  if (filters.category) queryParams.append("category", filters.category);
  if (filters.severity) queryParams.append("severity", filters.severity);
  if (filters.status) queryParams.append("status", filters.status);
  if (filters.actorId) queryParams.append("actorId", filters.actorId);
  if (filters.action) queryParams.append("action", filters.action);
  if (filters.resourceType) queryParams.append("resourceType", filters.resourceType);
  if (filters.resourceId) queryParams.append("resourceId", filters.resourceId);
  if (filters.organizationId) queryParams.append("organizationId", filters.organizationId);
  if (filters.startDate) queryParams.append("startDate", filters.startDate);
  if (filters.endDate) queryParams.append("endDate", filters.endDate);
  if (filters.search) queryParams.append("search", filters.search);

  const queryStr = queryParams.toString();

  return useQuery({
    queryKey: ["audit-logs", "platform", filters],
    queryFn: () =>
      apiGet<AuditLog[]>(`/audit-logs/platform${queryStr ? `?${queryStr}` : ""}`),
  });
}

export function useAuditStats() {
  return useQuery({
    queryKey: ["audit-logs", "stats"],
    queryFn: () => apiGet<AuditAdminStats>("/audit-logs/stats"),
  });
}
