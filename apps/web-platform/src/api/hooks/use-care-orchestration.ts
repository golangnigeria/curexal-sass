import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  orchestrationService,
  type CreateCareRequestPayload,
} from "../services/orchestration.service";
import type { CareRequestFilter } from "../contracts";

export const careKeys = {
  all: ["care-orchestration"] as const,
  requests: (filter?: CareRequestFilter) => [...careKeys.all, "requests", filter] as const,
  request: (id: string) => [...careKeys.all, "request", id] as const,
  myRequests: () => [...careKeys.all, "my-requests"] as const,
  myJourney: () => [...careKeys.all, "my-journey"] as const,
};

/**
 * Hook to query care requests for clinical / care agent workspaces
 */
export function useCareRequests(filter?: CareRequestFilter) {
  return useQuery({
    queryKey: careKeys.requests(filter),
    queryFn: () => orchestrationService.listCareRequests(filter),
  });
}

/**
 * Hook to retrieve a single care request
 */
export function useCareRequest(requestId: string) {
  return useQuery({
    queryKey: careKeys.request(requestId),
    queryFn: () => orchestrationService.getCareRequestById(requestId),
    enabled: Boolean(requestId),
  });
}

/**
 * Hook to create a new Care Request
 */
export function useCreateCareRequest() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateCareRequestPayload) =>
      orchestrationService.createCareRequest(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: careKeys.all });
    },
  });
}

/**
 * Hook to fetch care requests for the active patient portal user
 */
export function useMyCareRequests() {
  return useQuery({
    queryKey: careKeys.myRequests(),
    queryFn: () => orchestrationService.getMyCareRequests(),
  });
}

/**
 * Hook to fetch the live Curexal Care Journey tracker for the active patient portal user
 */
export function useMyCareJourney() {
  return useQuery({
    queryKey: careKeys.myJourney(),
    queryFn: () => orchestrationService.getMyCareJourney(),
  });
}

/**
 * Hook to submit clinical triage vitals
 */
export function useSubmitTriage() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ requestId, payload }: { requestId: string; payload: any }) =>
      orchestrationService.submitTriage(requestId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: careKeys.all });
    },
  });
}

/**
 * Hook to get provider matching candidates
 */
export function useMatchProviders(requestId: string) {
  return useQuery({
    queryKey: [...careKeys.all, "match-providers", requestId],
    queryFn: () => orchestrationService.matchProviders(requestId),
    enabled: Boolean(requestId),
  });
}

/**
 * Hook to assign provider
 */
export function useAssignProvider() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ requestId, providerId }: { requestId: string; providerId: string }) =>
      orchestrationService.assignProvider(requestId, providerId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: careKeys.all });
    },
  });
}

