import { useBootstrap } from "./use-bootstrap";

export function useCapabilities() {
  const { data: bootstrap, isLoading } = useBootstrap();

  const capabilities = bootstrap?.capabilities || [];
  const modules = bootstrap?.modules || [];
  const plan = bootstrap?.subscription?.plan || "smart";
  const limits = bootstrap?.limits || { maxBranches: 1, maxMembers: 5, storageGb: 10 };

  const hasCapability = (capabilityCode: string): boolean => {
    if (capabilities.includes("*")) return true;
    return capabilities.includes(capabilityCode);
  };

  const isModuleEnabled = (moduleCode: string): boolean => {
    const mod = modules.find((m) => m.code === moduleCode);
    return mod ? mod.enabled : false;
  };

  const isModuleLicensed = (moduleCode: string): boolean => {
    const mod = modules.find((m) => m.code === moduleCode);
    return mod ? mod.licensed : false;
  };

  const isUpgradeAvailable = (moduleCode: string): boolean => {
    const mod = modules.find((m) => m.code === moduleCode);
    return mod ? mod.upgradeAvailable : false;
  };

  return {
    capabilities,
    modules,
    plan,
    limits,
    isLoading,
    hasCapability,
    isModuleEnabled,
    isModuleLicensed,
    isUpgradeAvailable,
  };
}

export function useHasCapability(capabilityCode: string): boolean {
  const { hasCapability } = useCapabilities();
  return hasCapability(capabilityCode);
}
