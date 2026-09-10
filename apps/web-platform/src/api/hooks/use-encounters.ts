import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  encounterService,
  type DispatchOrdersPayload,
  type EncounterFilter,
} from "../services/encounter.service";
import type {
  StartEncounterPayload,
  UpdateSOAPNotesPayload,
  AddDiagnosisPayload,
  CreatePrescriptionPayload,
} from "../contracts";

export const encounterKeys = {
  all: ["encounters"] as const,
  list: (filter?: EncounterFilter) => [...encounterKeys.all, "list", filter] as const,
  detail: (id: string) => [...encounterKeys.all, "detail", id] as const,
  diagnoses: (id: string) => [...encounterKeys.detail(id), "diagnoses"] as const,
  prescriptions: (id: string) => [...encounterKeys.detail(id), "prescriptions"] as const,
};

/**
 * Hook to list encounters
 */
export function useEncounters(filter?: EncounterFilter) {
  return useQuery({
    queryKey: encounterKeys.list(filter),
    queryFn: () => encounterService.listEncounters(filter),
  });
}

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
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: encounterKeys.all });
      queryClient.invalidateQueries({ queryKey: ["care-orchestration"] });
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
 * Hook to add ICD-10 diagnosis
 */
export function useAddDiagnosis() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      encounterId,
      payload,
    }: {
      encounterId: string;
      payload: AddDiagnosisPayload;
    }) => encounterService.addDiagnosis(encounterId, payload),
    onSuccess: (_, vars) => {
      queryClient.invalidateQueries({ queryKey: encounterKeys.detail(vars.encounterId) });
      queryClient.invalidateQueries({ queryKey: encounterKeys.diagnoses(vars.encounterId) });
    },
  });
}

/**
 * Hook to list diagnoses
 */
export function useEncounterDiagnoses(encounterId: string) {
  return useQuery({
    queryKey: encounterKeys.diagnoses(encounterId),
    queryFn: () => encounterService.listDiagnoses(encounterId),
    enabled: Boolean(encounterId),
  });
}

/**
 * Hook to create electronic prescription
 */
export function useCreatePrescription() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      encounterId,
      payload,
    }: {
      encounterId: string;
      payload: CreatePrescriptionPayload;
    }) => encounterService.createPrescription(encounterId, payload),
    onSuccess: (_, vars) => {
      queryClient.invalidateQueries({ queryKey: encounterKeys.detail(vars.encounterId) });
      queryClient.invalidateQueries({ queryKey: encounterKeys.prescriptions(vars.encounterId) });
      queryClient.invalidateQueries({ queryKey: ["care-orchestration"] });
    },
  });
}

/**
 * Hook to list prescriptions
 */
export function useEncounterPrescriptions(encounterId: string) {
  return useQuery({
    queryKey: encounterKeys.prescriptions(encounterId),
    queryFn: () => encounterService.listPrescriptions(encounterId),
    enabled: Boolean(encounterId),
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
 * Hook to complete encounter and trigger POS billing
 */
export function useCompleteEncounter() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (
      param: string | { encounterId: string; consultationFee?: number }
    ) => {
      const encounterId = typeof param === "string" ? param : param.encounterId;
      const consultationFee = typeof param === "string" ? undefined : param.consultationFee;
      return encounterService.completeEncounter(encounterId, consultationFee);
    },
    onSuccess: (_, vars) => {
      queryClient.invalidateQueries({ queryKey: encounterKeys.all });
      queryClient.invalidateQueries({ queryKey: ["care-orchestration"] });
      queryClient.invalidateQueries({ queryKey: ["billing"] });
    },
  });
}
