import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useBootstrap } from "./use-bootstrap";
import {
  organizationService,
  CreateBranchRequest,
  InviteMemberRequest,
  CreateRoleRequest,
} from "../services/organization.service";

export function useOrgDashboardMetrics() {
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useQuery({
    queryKey: ["organization", orgId, "dashboard"],
    queryFn: () => organizationService.getDashboardMetrics(orgId!),
    enabled: Boolean(orgId),
    staleTime: 1000 * 60, // 1 minute
  });
}

export function useOrgBranches() {
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useQuery({
    queryKey: ["organization", orgId, "branches"],
    queryFn: () => organizationService.getBranches(orgId!),
    enabled: Boolean(orgId),
    staleTime: 1000 * 30,
  });
}

export function useCreateBranch() {
  const queryClient = useQueryClient();
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useMutation({
    mutationFn: (data: CreateBranchRequest) => organizationService.createBranch(orgId!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organization", orgId, "branches"] });
      queryClient.invalidateQueries({ queryKey: ["organization", orgId, "dashboard"] });
      queryClient.invalidateQueries({ queryKey: ["bootstrap"] });
    },
  });
}

export function useOrgMembers() {
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useQuery({
    queryKey: ["organization", orgId, "members"],
    queryFn: () => organizationService.getMembers(orgId!),
    enabled: Boolean(orgId),
    staleTime: 1000 * 30,
  });
}

export function useInviteMember() {
  const queryClient = useQueryClient();
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useMutation({
    mutationFn: (data: InviteMemberRequest) => organizationService.inviteMember(orgId!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organization", orgId, "members"] });
      queryClient.invalidateQueries({ queryKey: ["organization", orgId, "dashboard"] });
    },
  });
}

export function useOrgRoles() {
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useQuery({
    queryKey: ["organization", orgId, "roles"],
    queryFn: () => organizationService.getRoles(orgId!),
    enabled: Boolean(orgId),
    staleTime: 1000 * 60,
  });
}

export function useCreateRole() {
  const queryClient = useQueryClient();
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useMutation({
    mutationFn: (data: CreateRoleRequest) => organizationService.createRole(orgId!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organization", orgId, "roles"] });
    },
  });
}

export function useOrgCatalogs() {
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useQuery({
    queryKey: ["organization", orgId, "catalogs"],
    queryFn: () => organizationService.getCatalogs(orgId!),
    enabled: Boolean(orgId),
    staleTime: 1000 * 60,
  });
}

export function useUpdateCatalogPrice() {
  const queryClient = useQueryClient();
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useMutation({
    mutationFn: ({ itemId, customPrice }: { itemId: string; customPrice: number }) =>
      organizationService.updateCatalogPrice(orgId!, itemId, customPrice),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organization", orgId, "catalogs"] });
    },
  });
}

export function useOrgIntegrations() {
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useQuery({
    queryKey: ["organization", orgId, "integrations"],
    queryFn: () => organizationService.getApiKeys(orgId!),
    enabled: Boolean(orgId),
    staleTime: 1000 * 60,
  });
}

export function useCreateApiKey() {
  const queryClient = useQueryClient();
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useMutation({
    mutationFn: ({ name, scopes }: { name: string; scopes: string[] }) =>
      organizationService.createApiKey(orgId!, name, scopes),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organization", orgId, "integrations"] });
    },
  });
}

export function useOrgAuditLogs() {
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useQuery({
    queryKey: ["organization", orgId, "audit"],
    queryFn: () => organizationService.getAuditLogs(orgId!),
    enabled: Boolean(orgId),
    staleTime: 1000 * 30,
  });
}

export function useSubscribeCapability() {
  const queryClient = useQueryClient();
  const { data: bootstrap } = useBootstrap();
  const orgId = bootstrap?.organization?.id;

  return useMutation({
    mutationFn: ({ capabilityCode, currency }: { capabilityCode: string; currency?: string }) =>
      organizationService.subscribeCapability(orgId!, capabilityCode, currency),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bootstrap"] });
      queryClient.invalidateQueries({ queryKey: ["organization", orgId, "dashboard"] });
    },
  });
}
