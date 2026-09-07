import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { patientService } from "../services/patient.service";
import type {
  DuplicateEvaluationRequest,
  PatientListFilter,
  RegisterCanonicalPatientPayload,
  SendPortalOTPPayload,
  SetPortalPINPayload,
} from "../contracts";

export const patientKeys = {
  all: ["patients"] as const,
  lists: () => [...patientKeys.all, "list"] as const,
  list: (filter?: PatientListFilter) => [...patientKeys.lists(), filter] as const,
  details: () => [...patientKeys.all, "detail"] as const,
  detail: (id: string) => [...patientKeys.details(), id] as const,
};

/**
 * Hook to query paginated patient directory with search & filters
 */
export function usePatients(filter?: PatientListFilter) {
  return useQuery({
    queryKey: patientKeys.list(filter),
    queryFn: () => patientService.listPatients(filter),
  });
}

/**
 * Hook to retrieve a single patient's profile
 */
export function usePatient(patientId: string) {
  return useQuery({
    queryKey: patientKeys.detail(patientId),
    queryFn: () => patientService.getPatientById(patientId),
    enabled: Boolean(patientId),
  });
}

/**
 * Hook to evaluate duplicates via Master Patient Index
 */
export function useResolveDuplicates() {
  return useMutation({
    mutationFn: (payload: DuplicateEvaluationRequest) =>
      patientService.resolveDuplicates(payload),
  });
}

/**
 * Hook to register a canonical patient record
 */
export function useRegisterPatient() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: RegisterCanonicalPatientPayload) =>
      patientService.registerPatient(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: patientKeys.lists() });
    },
  });
}

/**
 * Hook to send passwordless login OTP to patient portal account
 */
export function useSendPortalOTP() {
  return useMutation({
    mutationFn: (payload: SendPortalOTPPayload) =>
      patientService.sendPortalOTP(payload),
  });
}

/**
 * Hook to verify passwordless login OTP and establish patient session
 */
export function useVerifyPortalOTP() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { identifier: string; code: string }) =>
      patientService.verifyPortalOTP(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: patientKeys.all });
    },
  });
}

/**
 * Hook to set quick login PIN for patient portal account
 */
export function useSetPortalPIN() {
  return useMutation({
    mutationFn: ({ patientId, payload }: { patientId: string; payload: SetPortalPINPayload }) =>
      patientService.setPortalPIN(patientId, payload),
  });
}
