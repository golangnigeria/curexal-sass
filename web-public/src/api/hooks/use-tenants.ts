import { useQuery } from "@tanstack/react-query";
import { useApiClient } from "..";

export const useTenantDetails = (id: string | undefined) => {
  const api = useApiClient() as any;

  return useQuery({
    queryKey: ["tenant", id],
    queryFn: async () => {
      if (!id || !api?.Tenant?.getTenant) return null;
      const res = await api.Tenant.getTenant({
        params: { id },
      });
      if (res?.status !== 200) {
        throw new Error((res?.body as any)?.message || "Failed to fetch tenant details");
      }
      return res.body;
    },
    enabled: !!id,
  });
};

export const useActiveTenant = () => {
  const api = useApiClient() as any;

  return useQuery({
    queryKey: ["tenant", "active"],
    queryFn: async () => {
      if (!api?.Tenant?.getActiveTenant) return null;
      const res = await api.Tenant.getActiveTenant();
      if (res?.status !== 200) {
        const err: any = new Error((res?.body as any)?.message || "Failed to fetch active tenant");
        err.status = res?.status;
        throw err;
      }
      return res.body;
    },
    retry: (failureCount, error: any) => {
      if (error?.status === 401 || error?.status === 403 || error?.status === 404) return false;
      return failureCount < 2;
    },
  });
};

export const useTenants = () => {
  return { tenants: [], loading: false, error: null, refetch: () => {} };
};
