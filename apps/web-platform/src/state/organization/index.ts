import { useBootstrap } from "@/api/hooks/use-bootstrap";

export function useOrganizationCapabilities() {
  const { data: bootstrap } = useBootstrap();

  const capabilities = new Set(bootstrap?.capabilities || []);

  return {
    capabilities,
    hasCapability: (code: string) => capabilities.has(code),
    planTier: (bootstrap?.organization as any)?.planTier || "Starter",
  };
}
