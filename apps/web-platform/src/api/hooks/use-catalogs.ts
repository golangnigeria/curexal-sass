import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiGet, apiPost, apiPut } from "@/api/client";
import type { CatalogDomain, CatalogItem, ICD10Code } from "@/api/contracts";

export function useCatalogItems(domain: CatalogDomain, category?: string, activeOnly?: boolean) {
  const queryParams = new URLSearchParams();
  if (category) queryParams.append("category", category);
  if (activeOnly) queryParams.append("active_only", "true");
  const queryStr = queryParams.toString();

  return useQuery({
    queryKey: ["platform", "catalogs", domain, category || "all", activeOnly ? "active" : "all"],
    queryFn: () =>
      apiGet<CatalogItem[]>(
        `/platform/catalogs/${domain}${queryStr ? `?${queryStr}` : ""}`
      ),
    enabled: !!domain,
  });
}

export function useSearchCatalogItems(domain: CatalogDomain, query: string) {
  return useQuery({
    queryKey: ["platform", "catalogs", domain, "search", query],
    queryFn: () =>
      apiGet<CatalogItem[]>(`/platform/catalogs/${domain}/search?q=${encodeURIComponent(query)}`),
    enabled: !!domain && !!query && query.length >= 2,
  });
}

export function useCreateCatalogItem(domain: CatalogDomain) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CatalogItem) =>
      apiPost<CatalogItem>(`/platform/catalogs/${domain}`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "catalogs", domain] });
    },
  });
}

export function useUpdateCatalogItem(domain: CatalogDomain, itemId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CatalogItem) =>
      apiPut<CatalogItem>(`/platform/catalogs/${domain}/${itemId}`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platform", "catalogs", domain] });
    },
  });
}

export function useSearchICD10(query: string) {
  return useQuery({
    queryKey: ["catalogs", "icd10", query],
    queryFn: () =>
      apiGet<ICD10Code[]>(`/catalogs/icd10?q=${encodeURIComponent(query)}`),
    enabled: !!query && query.length >= 2,
  });
}
