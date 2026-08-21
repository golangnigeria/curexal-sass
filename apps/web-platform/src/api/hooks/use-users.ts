import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut } from "@/api/client";
import type { DirectoryUser } from "@/api/contracts";

export function usePlatformUsers() {
  return useQuery({
    queryKey: ["platform", "users"],
    queryFn: () => apiGet<DirectoryUser[]>("/users"),
  });
}
