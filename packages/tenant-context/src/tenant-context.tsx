import React, { createContext, useContext } from "react";
import { useLocation, useParams } from "react-router-dom";
import { authClient, fetchStaffSession } from "@curexal/auth";
import { resolveEntitledWorkspaces, WORKSPACE_MODULE_REGISTRY } from "@curexal/permissions";
import { isPlatformHost } from "./url-builder";
import type { OrganizationPayload, BranchPayload, BranchSummaryPayload } from "@curexal/contracts";

export interface TenantContextType {
  organization?: OrganizationPayload;
  activeBranch?: BranchPayload;
  availableBranches?: BranchSummaryPayload[];
  isPlatform: boolean;
  isLoading: boolean;
  refetch: () => Promise<any>;
}

const TenantContext = createContext<TenantContextType | null>(null);

export const TenantProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { data: session, isPending } = authClient.useSession();
  const bootstrap = session?.bootstrap;

  const value: TenantContextType = {
    organization: bootstrap?.organization,
    activeBranch: bootstrap?.branch,
    availableBranches: bootstrap?.availableBranches || [],
    isPlatform: isPlatformHost(),
    isLoading: isPending,
    refetch: () => fetchStaffSession(),
  };

  return <TenantContext.Provider value={value}>{children}</TenantContext.Provider>;
};

export function useTenant() {
  return useContext(TenantContext);
}

export function useActiveTenant() {
  const { data: session, isPending } = authClient.useSession();
  const bootstrap = session?.bootstrap;

  const activeBranch = bootstrap?.branch;
  const availableBranches = bootstrap?.availableBranches || [];
  const organization = bootstrap?.organization;
  const isPlatform = isPlatformHost();

  return {
    organization,
    activeBranch,
    availableBranches,
    isPlatform,
    isLoading: isPending,
    refetch: () => fetchStaffSession(),
  };
}

export function useOrganizationCapabilities() {
  const { data: session } = authClient.useSession();
  const bootstrap = session?.bootstrap;

  const capabilities = new Set(bootstrap?.capabilities || []);

  return {
    capabilities,
    hasCapability: (code: string) => capabilities.has(code),
    planTier: (bootstrap?.organization as any)?.planTier || "Starter",
  };
}

export function useActiveWorkspace() {
  const location = useLocation();
  const params = useParams<{ branchSlug?: string }>();
  const { data: session } = authClient.useSession();
  const bootstrap = session?.bootstrap;

  const pathParts = location.pathname.split("/").filter(Boolean);
  const branchSlug = params.branchSlug || (pathParts.length > 0 && pathParts[0] !== "organization" && pathParts[0] !== "platform" ? pathParts[0] : "main");
  const activeModuleCode = pathParts[1] || "dashboard";

  const entitledWorkspaces = resolveEntitledWorkspaces(bootstrap, session?.user);
  const activeWorkspace = entitledWorkspaces.find((w) => w.path === activeModuleCode) || null;

  return {
    branchSlug,
    activeModuleCode,
    activeWorkspace,
    entitledWorkspaces,
    allModules: WORKSPACE_MODULE_REGISTRY,
  };
}
