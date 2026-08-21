import { useQuery } from "@tanstack/react-query";
import { apiGet } from "@/api/client";
import type { BootstrapContractResponse } from "@/api/contracts";

export function useBootstrap() {
  return useQuery({
    queryKey: ["bootstrap"],
    queryFn: () => apiGet<BootstrapContractResponse>("/bootstrap"),
    staleTime: 1000 * 60 * 5, // 5 minutes
    retry: 1,
  });
}
