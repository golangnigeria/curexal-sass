import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  encounterService,
  type DispatchOrdersPayload,
  type StartEncounterPayload,
  type UpdateSOAPNotesPayload,
} from "../services/encounter.service";

export const encounterKeys = {
  all: ["encounters"] as const,
  detail: (id: string) => [...encounterKeys.all, "detail", id] as const,
};

/**
 * Hook to retrieve single active encounter
 */
export function useEncounter(encounterId: string) {
  return useQuery({
    queryKey: encounterKeys.detail(encounterId),
    queryFn: () => encounterService.getEncounterById(encounterId),
    enabled: Boolean(encounterId),
  });
}

/**
 * Hook to start encounter
 */
export function useStartEncounter() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: StartEncounterPayload) => encounterService.startEncounter(payload),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: encounterKeys.all });
    },
  });
}

/**
 * Hook to save physician SOAP notes
 */
export function useSaveSOAPNotes() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      encounterId,
      payload,
    }: {
      encounterId: string;
      payload: UpdateSOAPNotesPayload;
    }) => encounterService.saveSOAPNotes(encounterId, payload),
    onSuccess: (_, vars) => {
      queryClient.invalidateQueries({ queryKey: encounterKeys.detail(vars.encounterId) });
    },
  });
}

/**
 * Hook to dispatch clinical lab, radiology, and prescription orders
 */
export function useDispatchOrders() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      encounterId,
      payload,
    }: {
      encounterId: string;
      payload: DispatchOrdersPayload;
    }) => encounterService.dispatchOrders(encounterId, payload),
    onSuccess: (_, vars) => {
      queryClient.invalidateQueries({ queryKey: encounterKeys.detail(vars.encounterId) });
      queryClient.invalidateQueries({ queryKey: ["care-orchestration"] });
    },
  });
}

/**
 * Hook to complete encounter and discharge patient
 */
export function useCompleteEncounter() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (encounterId: string) => encounterService.completeEncounter(encounterId),
    onSuccess: (_, encounterId) => {
      queryClient.invalidateQueries({ queryKey: encounterKeys.detail(encounterId) });
      queryClient.invalidateQueries({ queryKey: ["care-orchestration"] });
    },
  });
}
