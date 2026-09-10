import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { operationsService } from "../services/operations.service";
import type {
  AppointmentFilter,
  CreateAppointmentPayload,
} from "../contracts";

export const appointmentKeys = {
  all: ["appointments"] as const,
  list: (filter?: AppointmentFilter) => [...appointmentKeys.all, "list", filter] as const,
  detail: (id: string) => [...appointmentKeys.all, "detail", id] as const,
  providers: (status?: string) => ["providers", "profiles", status] as const,
};

export function useAppointments(filter?: AppointmentFilter) {
  return useQuery({
    queryKey: appointmentKeys.list(filter),
    queryFn: () => operationsService.listAppointments(filter),
  });
}

export function useAppointment(id: string) {
  return useQuery({
    queryKey: appointmentKeys.detail(id),
    queryFn: () => operationsService.getAppointmentById(id),
    enabled: Boolean(id),
  });
}

export function useCreateAppointment() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateAppointmentPayload) =>
      operationsService.createAppointment(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: appointmentKeys.all });
    },
  });
}

export function useProviderProfiles(status?: string, specialty?: string) {
  return useQuery({
    queryKey: appointmentKeys.providers(status),
    queryFn: () => operationsService.listProviderProfiles(status, specialty),
  });
}

export function useUpdateProviderStatus() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status, reason }: { id: string; status: string; reason?: string }) =>
      operationsService.updateProviderStatus(id, status, reason),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["providers"] });
    },
  });
}
