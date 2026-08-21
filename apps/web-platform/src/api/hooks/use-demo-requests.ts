import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPut } from "@/api/client";
import type { DemoRequest } from "@/api/contracts";

export function useDemoRequests() {
  return useQuery({
    queryKey: ["demo-requests"],
    queryFn: () => apiGet<DemoRequest[]>("/demo-requests"),
  });
}

export function useUpdateDemoRequest(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { status?: string; notes?: string }) =>
      apiPut<DemoRequest>(`/demo-requests/${id}`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["demo-requests"] });
    },
  });
}
