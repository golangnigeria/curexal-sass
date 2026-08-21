import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut } from "@/api/client";
import type {
  FacilityTypeEntity,
  RegistrationFormDTO,
  NavigationMenuDTO,
  SetupStepDTO,
  DashboardDTO,
} from "@/api/contracts";

export function useFacilityTypes() {
  return useQuery({
    queryKey: ["platform", "facility-types"],
    queryFn: () => apiGet<FacilityTypeEntity[]>("/platform/facility-types"),
  });
}

export function useFacilityTypeRegistrationForm(typeId: string) {
  return useQuery({
    queryKey: ["platform", "facility-types", typeId, "form"],
    queryFn: () =>
      apiGet<RegistrationFormDTO>(`/platform/facility-types/${typeId}/form`),
    enabled: !!typeId,
  });
}

export function useFacilityTypeNavigation(typeId: string) {
  return useQuery({
    queryKey: ["platform", "facility-types", typeId, "navigation"],
    queryFn: () =>
      apiGet<NavigationMenuDTO>(`/platform/facility-types/${typeId}/navigation`),
    enabled: !!typeId,
  });
}

export function useFacilityTypeSetupSteps(typeId: string) {
  return useQuery({
    queryKey: ["platform", "facility-types", typeId, "setup-steps"],
    queryFn: () =>
      apiGet<SetupStepDTO[]>(`/platform/facility-types/${typeId}/setup-steps`),
    enabled: !!typeId,
  });
}

export function useFacilityTypeDashboard(typeId: string) {
  return useQuery({
    queryKey: ["platform", "facility-types", typeId, "dashboard"],
    queryFn: () =>
      apiGet<DashboardDTO>(`/platform/facility-types/${typeId}/dashboard`),
    enabled: !!typeId,
  });
}

export function useCreateFacilityType() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: FacilityTypeEntity) =>
      apiPost<FacilityTypeEntity>("/platform/facility-types", payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "facility-types"] });
    },
  });
}

export function useUpdateFacilityType(typeId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: FacilityTypeEntity) =>
      apiPut<FacilityTypeEntity>(`/platform/facility-types/${typeId}`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "facility-types"] });
    },
  });
}
