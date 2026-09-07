import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut, apiPatch } from "@/api/client";
import type {
  Organization,
  CreateOrganizationPayload,
  UpdateOrganizationPayload,
  UpdateOrganizationProfilePayload,
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

export function useOrganizationProfile() {
  return useQuery({
    queryKey: ["organization", "profile"],
    queryFn: () => apiGet<Organization>("/organization/profile"),
  });
}

export function useUpdateOrganizationProfile() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Partial<UpdateOrganizationProfilePayload>) =>
      apiPut<Organization>("/organization/profile", payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organization", "profile"] });
      queryClient.invalidateQueries({ queryKey: ["organizations"] });
      queryClient.invalidateQueries({ queryKey: ["platform", "organizations"] });
      queryClient.invalidateQueries({ queryKey: ["bootstrap"] });
    },
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
    queryFn: async () => {
      const raw = await apiGet<any[]>(`/organizations/${id}/documents`);
      if (!Array.isArray(raw)) return [];
      return raw.map((item: any) => {
        const doc = item.document || item;
        const resolvedFilename =
          doc.originalFilename ||
          doc.original_filename ||
          doc.filename ||
          doc.fileName ||
          "";
        const resolvedSize =
          doc.fileSizeBytes ??
          doc.file_size_bytes ??
          doc.fileSize ??
          0;
        return {
          id: doc.id,
          organizationId: doc.organizationId || doc.organization_id,
          documentType: doc.documentType || doc.document_type || "",
          originalFilename: resolvedFilename,
          fileName: resolvedFilename,
          storageKey: doc.storageKey || doc.storage_key || "",
          mimeType: doc.mimeType || doc.mime_type || "",
          fileSizeBytes: Number(resolvedSize),
          fileSize: Number(resolvedSize),
          checksumSha256: doc.checksumSha256 || doc.checksum_sha256 || "",
          uploadedBy: doc.uploadedBy || doc.uploaded_by,
          uploadedAt: doc.uploadedAt || doc.uploaded_at || doc.createdAt || doc.created_at,
          status: (doc.status || "pending").toLowerCase(),
          version: doc.version || 1,
          reviewedBy: doc.reviewedBy || doc.reviewed_by,
          reviewedAt: doc.reviewedAt || doc.reviewed_at,
          rejectionReason: doc.rejectionReason || doc.rejection_reason,
          presignedUrl: item.presignedUrl || item.presigned_url || doc.presignedUrl || "",
          createdAt: doc.createdAt || doc.created_at,
          updatedAt: doc.updatedAt || doc.updated_at,
        } as OrganizationDocument;
      });
    },
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

export function useUpdateOrganizationSettings(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Partial<OrganizationSettings>) =>
      apiPut<OrganizationSettings>(`/organizations/${id}/settings`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organizations", id, "settings"] });
      queryClient.invalidateQueries({ queryKey: ["organizations", id] });
      queryClient.invalidateQueries({ queryKey: ["organization", "profile"] });
      queryClient.invalidateQueries({ queryKey: ["platform", "organizations"] });
      queryClient.invalidateQueries({ queryKey: ["bootstrap"] });
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
      queryClient.invalidateQueries({ queryKey: ["organization", "profile"] });
      queryClient.invalidateQueries({ queryKey: ["organizations", id, "documents"] });
      queryClient.invalidateQueries({ queryKey: ["bootstrap"] });
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
      queryClient.invalidateQueries({ queryKey: ["organization", "profile"] });
      queryClient.invalidateQueries({ queryKey: ["organizations", id, "documents"] });
      queryClient.invalidateQueries({ queryKey: ["bootstrap"] });
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
      queryClient.invalidateQueries({ queryKey: ["organizations", orgId] });
      queryClient.invalidateQueries({ queryKey: ["organization", "profile"] });
      queryClient.invalidateQueries({ queryKey: ["bootstrap"] });
      queryClient.invalidateQueries({ queryKey: ["platform", "organizations"] });
    },
  });
}

export function useResendOwnerInvite(orgId: string) {
  return useMutation({
    mutationFn: () =>
      apiPost<{ message: string }>(`/organizations/${orgId}/resend-invite`),
  });
}
