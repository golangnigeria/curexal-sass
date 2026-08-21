import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiPatch } from "@/api/client";
import type {
  Organization,
  CreateOrganizationPayload,
  UpdateOrganizationPayload,
  OrganizationDocument,
  OrganizationSettings,
} from "@/api/contracts";

export function usePlatformOrganizations() {
  return useQuery({
    queryKey: ["platform", "organizations"],
    queryFn: () => apiGet<Organization[]>("/platform/organizations"),
  });
}

export function useOrganization(id: string) {
  return useQuery({
    queryKey: ["organizations", id],
    queryFn: () => apiGet<Organization>(`/organizations/${id}`),
    enabled: !!id,
  });
}

export function useOrganizationSettings(id: string) {
  return useQuery({
    queryKey: ["organizations", id, "settings"],
    queryFn: () => apiGet<OrganizationSettings>(`/organizations/${id}/settings`),
    enabled: !!id,
  });
}

export function useOrganizationDocuments(id: string) {
  return useQuery({
    queryKey: ["organizations", id, "documents"],
    queryFn: () => apiGet<OrganizationDocument[]>(`/organizations/${id}/documents`),
    enabled: !!id,
  });
}

export function useCreateOrganization() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateOrganizationPayload) =>
      apiPost<Organization>("/organizations", payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "organizations"] });
    },
  });
}

export function useUpdateOrganization(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: UpdateOrganizationPayload) =>
      apiPut<Organization>(`/organizations/${id}`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "organizations"] });
      queryClient.invalidateQueries({ queryKey: ["organizations", id] });
    },
  });
}

export function useApproveOrganization(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => apiPost<{ message: string }>(`/platform/organizations/${id}/approve`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "organizations"] });
      queryClient.invalidateQueries({ queryKey: ["organizations", id] });
    },
  });
}

export function useRejectOrganization(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (reason: string) =>
      apiPost<{ message: string }>(`/platform/organizations/${id}/reject`, { reason }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "organizations"] });
      queryClient.invalidateQueries({ queryKey: ["organizations", id] });
    },
  });
}

export function useReviewDocument(docId: string, orgId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { status: "approved" | "rejected"; rejectionReason?: string }) =>
      apiPatch<{ message: string }>(`/platform/documents/${docId}/review`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organizations", orgId, "documents"] });
    },
  });
}

export function useResendOwnerInvite(orgId: string) {
  return useMutation({
    mutationFn: () =>
      apiPost<{ message: string }>(`/organizations/${orgId}/resend-invite`),
  });
}
