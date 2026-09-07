import { useEffect, useState } from "react";
import {
  apiGet,
  apiPost,
  refreshCsrfToken,
  getCsrfToken,
  setClientActiveOrgId,
  setAuthToken,
} from "@curexal/api-client";
import {
  ROLES,
  type UserRoleResponse,
  type BootstrapContractResponse,
  type StaffSignInPayload,
  type BranchSelectionOption,
  type SelectBranchPayload,
  type SwitchBranchPayload,
} from "@curexal/contracts";

export interface StaffSessionData {
  user: UserRoleResponse;
  bootstrap?: BootstrapContractResponse;
}

export interface BranchSelectionData {
  requireBranchSelection: true;
  selectionToken: string;
  assignedBranches: BranchSelectionOption[];
  organization?: {
    id: string;
    name: string;
    slug: string;
  };
}

export type SignInResult =
  | (StaffSessionData & { requireBranchSelection?: false; accessToken?: string })
  | BranchSelectionData;

export function normalizeStaffUser(raw: any): UserRoleResponse | null {
  if (!raw) return null;
  const id = raw.identity?.user?.id || raw.id || raw.userId || raw.user?.id;
  if (!id) return null;

  const email = raw.identity?.user?.email || raw.email || raw.user?.email || "";
  const name =
    raw.identity?.user?.name ||
    raw.name ||
    raw.user?.user_metadata?.name ||
    raw.user?.user_metadata?.full_name ||
    email;
  const isPlatformAdmin = Boolean(
    raw.identity?.platform?.isPlatformAdmin ||
      raw.isPlatformAdmin ||
      raw.identity?.platform?.role === ROLES.SUPER_ADMIN ||
      raw.platformRole === ROLES.SUPER_ADMIN
  );
  const platformRole =
    raw.identity?.platform?.role ||
    raw.platformRole ||
    (isPlatformAdmin ? ROLES.SUPER_ADMIN : undefined);

  const orgId =
    raw.context?.activeOrganization?.id ||
    raw.identity?.organizations?.[0]?.id ||
    raw.organizationId ||
    raw.user?.user_metadata?.organization_id;
  const orgRole =
    raw.context?.activeOrganization?.role ||
    raw.identity?.organizations?.[0]?.role ||
    raw.organizationRole;

  const branchId =
    raw.context?.activeBranch?.id ||
    raw.activeBranchId ||
    raw.branchId;

  const workspaceRole =
    raw.context?.workspaceMembership?.role ||
    raw.context?.activeBranch?.role;

  const directRole =
    raw.identity?.user?.role ||
    raw.role;

  // Prioritize specific role (e.g. branch_admin, doctor, nurse, owner) over generic fallbacks
  const rolePool: string[] = [workspaceRole, directRole, orgRole, platformRole].filter(Boolean);
  const specificRole = rolePool.find((r) => r !== "member" && r !== "user");
  const effectiveRole = platformRole || specificRole || rolePool[0] || "";

  return {
    id,
    email,
    name,
    role: effectiveRole,
    platformRole: platformRole || undefined,
    organizationId: orgId,
    organizationRole: orgRole,
    branchId: branchId || undefined,
    activeBranchId: branchId || undefined,
    isPlatformAdmin,
    activeTenantId: raw.context?.activeTenant?.id || raw.activeTenantId || orgId,
    tenantSlug: raw.context?.activeTenant?.slug || raw.tenantSlug,
    availableTenants: raw.availableTenants || [],
    permissions: raw.permissions || [],
  };
}

let globalStaffSession: StaffSessionData | null = null;
let globalStaffIsPending = true;
const staffListeners = new Set<(session: { data: StaffSessionData | null; isPending: boolean }) => void>();

const notifyStaff = () => {
  staffListeners.forEach((l) => l({ data: globalStaffSession, isPending: globalStaffIsPending }));
};

let staffFetchPromise: Promise<StaffSessionData | null> | null = null;

export const fetchStaffSession = async (): Promise<StaffSessionData | null> => {
  if (staffFetchPromise) return staffFetchPromise;

  staffFetchPromise = (async () => {
    try {
      // 1. Fetch CSRF token for mutating requests
      await refreshCsrfToken();

      // 2. Fetch authenticated user profile from Go backend
      const rawUser = await apiGet<any>("/users/me").catch(() => null);
      const user = normalizeStaffUser(rawUser);

      if (user && user.id) {
        // 3. Fetch bootstrap contract from Go backend
        let bootstrap: BootstrapContractResponse | undefined;
        try {
          bootstrap = await apiGet<BootstrapContractResponse>("/bootstrap");
        } catch {
          // Bootstrap optional fallback
        }

        const activeOrgId = bootstrap?.organization?.id || user.organizationId || "";
        if (activeOrgId) {
          setClientActiveOrgId(activeOrgId);
        }

        globalStaffSession = { user, bootstrap };
      } else {
        globalStaffSession = null;
      }
    } catch {
      globalStaffSession = null;
    } finally {
      globalStaffIsPending = false;
      notifyStaff();
    }
    return globalStaffSession;
  })();

  return staffFetchPromise;
};

// Initial session check
fetchStaffSession();

export const staffAuthClient = {
  useSession: () => {
    const [state, setState] = useState({ data: globalStaffSession, isPending: globalStaffIsPending });

    useEffect(() => {
      const listener = (s: { data: StaffSessionData | null; isPending: boolean }) => setState(s);
      staffListeners.add(listener);

      if (!globalStaffIsPending && (state.data !== globalStaffSession || state.isPending !== globalStaffIsPending)) {
        setState({ data: globalStaffSession, isPending: globalStaffIsPending });
      }

      return () => {
        staffListeners.delete(listener);
      };
    }, [state.data, state.isPending]);

    return state;
  },

  getSession: async () => {
    if (globalStaffIsPending) {
      return { data: await fetchStaffSession() };
    }
    return { data: globalStaffSession };
  },

  signIn: async (payload: StaffSignInPayload) => {
    await refreshCsrfToken();
    const res = await apiPost<any>("/auth/sign-in", payload);
    const data = res?.data || res;
    if (data?.targetUrl || data?.status === "redirect_required") {
      return data;
    }
    if (data?.requireBranchSelection || data?.status === "branch_selection_required") {
      return {
        ...data,
        requireBranchSelection: true,
        selectionToken: data.selectionToken,
        assignedBranches: data.assignedBranches || [],
        organization: data.organization,
      };
    }
    if (
      data?.status === "unassigned_facility_branch" ||
      data?.status === "organization_access_required" ||
      data?.status === "organization_selection_required"
    ) {
      return data;
    }
    const token = data?.accessToken || res?.accessToken;
    if (token) {
      setAuthToken(token);
    }
    staffFetchPromise = null;
    const session = await fetchStaffSession();
    return { session, ...data };
  },

  selectBranch: async (payload: SelectBranchPayload) => {
    await refreshCsrfToken();
    const res = await apiPost<any>("/auth/select-branch", payload);
    const data = res?.data || res;
    if (data?.targetUrl || data?.status === "redirect_required") {
      return data;
    }
    const token = data?.accessToken || res?.accessToken;
    if (token) {
      setAuthToken(token);
    }
    staffFetchPromise = null;
    const session = await fetchStaffSession();
    return { session, ...data };
  },

  switchBranch: async (branchId: string) => {
    await refreshCsrfToken();
    const res = await apiPost<any>("/auth/switch-branch", { branchId });
    const data = res?.data || res;
    const token = data?.accessToken || res?.accessToken;
    if (token) {
      setAuthToken(token);
    }
    staffFetchPromise = null;
    const session = await fetchStaffSession();
    return { session, ...data };
  },

  signUp: async (payload: { email: string; password: string; name?: string }) => {
    await refreshCsrfToken();
    const res = await apiPost<any>("/auth/sign-up", payload);
    staffFetchPromise = null;
    return res;
  },

  signOut: async (redirectPath: string = "/login") => {
    try {
      await apiPost("/auth/logout");
    } catch {
      try {
        await apiPost("/auth/sign-out");
      } catch {
        // Ignore
      }
    }

    setAuthToken("");
    globalStaffSession = null;
    globalStaffIsPending = false;
    notifyStaff();
    if (typeof window !== "undefined") {
      window.location.href = redirectPath;
    }
  },

  switchContext: async (payload: { targetContext: string; targetId?: string }) => {
    await apiPost("/context/switch", payload);
    staffFetchPromise = null;
    return await fetchStaffSession();
  },

  forgotPassword: async (email: string) => {
    await refreshCsrfToken();
    return await apiPost<any>("/auth/forgot-password", { email });
  },

  setPassword: async (payload: { email?: string; code?: string; token?: string; password: string }) => {
    await refreshCsrfToken();
    return await apiPost<any>("/auth/set-password", payload);
  },

  setSessionToken: (token: string) => {
    setAuthToken(token);
    staffFetchPromise = null;
    return fetchStaffSession();
  },

  getCsrfToken: () => getCsrfToken(),
  refreshSession: async () => {
    try {
      await refreshCsrfToken();
      const res = await apiPost<any>("/auth/refresh");
      const data = res?.data || res;
      const token = data?.accessToken || res?.accessToken;
      if (token) {
        setAuthToken(token);
      }
    } catch {
      // Refresh might fail if no valid session/refresh cookie
    }
    staffFetchPromise = null;
    return await fetchStaffSession();
  },
};

// Aliased export for compatibility
export const authClient = staffAuthClient;
