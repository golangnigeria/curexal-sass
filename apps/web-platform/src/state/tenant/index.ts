import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { getOrgSlugFromHostname, isPlatformHost } from "@/lib/url-builder";

export function useActiveTenant() {
  const { data: bootstrap, isLoading, refetch } = useBootstrap();

  const activeBranch = bootstrap?.branch || bootstrap?.workspace;
  const availableBranches = bootstrap?.availableBranches || [];
  const organization = bootstrap?.organization;
  const isPlatform = isPlatformHost();

  return {
    organization,
    activeBranch,
    availableBranches,
    isPlatform,
    isLoading,
    refetch,
  };
}
