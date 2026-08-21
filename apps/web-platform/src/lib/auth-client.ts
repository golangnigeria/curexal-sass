import { useEffect, useState } from "react";
import { apiGet, apiPost, refreshCsrfToken, setCsrfToken } from "@/api/client";
import type { UserRoleResponse, BootstrapContractResponse } from "@/api/contracts";

interface SessionData {
  user: UserRoleResponse;
  bootstrap?: BootstrapContractResponse;
}

function normalizeUser(raw: any): UserRoleResponse | null {
  if (!raw) return null;
  const id = raw.identity?.user?.id || raw.id || raw.userId;
  if (!id) return null;

  const email = raw.identity?.user?.email || raw.email || "";
  const name = raw.identity?.user?.name || raw.name || email;
  const isPlatformAdmin = Boolean(
    raw.identity?.platform?.isPlatformAdmin ??
      raw.isPlatformAdmin ??
      raw.identity?.platform?.role === "super_admin" ??
      raw.platformRole === "super_admin"
  );
  const platformRole =
    raw.identity?.platform?.role ||
    raw.platformRole ||
    (isPlatformAdmin ? "super_admin" : (raw.role || ""));

  return {
    id,
    email,
    name,
    role: platformRole || raw.role || "",
    platformRole: platformRole || undefined,
    isPlatformAdmin,
    activeTenantId: raw.context?.activeTenant?.id || raw.activeTenantId,
    tenantSlug: raw.context?.activeTenant?.slug || raw.tenantSlug,
    availableTenants: raw.availableTenants || [],
    permissions: raw.permissions || [],
  };
}

let globalSession: SessionData | null = null;
let globalIsPending = true;
const listeners = new Set<(session: { data: SessionData | null; isPending: boolean }) => void>();

const notify = () => {
  listeners.forEach((l) => l({ data: globalSession, isPending: globalIsPending }));
};

let fetchPromise: Promise<SessionData | null> | null = null;

export const fetchSession = async (): Promise<SessionData | null> => {
  if (fetchPromise) return fetchPromise;

  fetchPromise = (async () => {
    try {
      // 1. Fetch CSRF token first
      await refreshCsrfToken();

      // 2. Fetch authenticated user profile
      const rawUser = await apiGet<any>("/users/me");
      const user = normalizeUser(rawUser);

      if (user && user.id) {
        // 3. Fetch bootstrap contract
        let bootstrap: BootstrapContractResponse | undefined;
        try {
          bootstrap = await apiGet<BootstrapContractResponse>("/bootstrap");
        } catch (bErr) {
          console.warn("Bootstrap fetch warning:", bErr);
        }

        globalSession = { user, bootstrap };
      } else {
        globalSession = null;
      }
    } catch (err) {
      globalSession = null;
    } finally {
      globalIsPending = false;
      notify();
    }
    return globalSession;
  })();

  return fetchPromise;
};

// Initial trigger
fetchSession();

export const authClient = {
  useSession: () => {
    const [state, setState] = useState({ data: globalSession, isPending: globalIsPending });

    useEffect(() => {
      const listener = (s: { data: SessionData | null; isPending: boolean }) => setState(s);
      listeners.add(listener);

      if (!globalIsPending && (state.data !== globalSession || state.isPending !== globalIsPending)) {
        setState({ data: globalSession, isPending: globalIsPending });
      }

      return () => {
        listeners.delete(listener);
      };
    }, [state.data, state.isPending]);

    return state;
  },

  getSession: async () => {
    if (globalIsPending) {
      return { data: await fetchSession() };
    }
    return { data: globalSession };
  },

  signIn: async (payload: { email: string; password: string }) => {
    await refreshCsrfToken();
    const res = await apiPost<any>("/auth/sign-in", payload);
    fetchPromise = null;
    const session = await fetchSession();
    return session || res;
  },

  signOut: async () => {
    try {
      await apiPost("/auth/sign-out");
    } catch (err) {
      console.error("Sign out failed:", err);
    }
    globalSession = null;
    globalIsPending = false;
    notify();
    window.location.href = "/login";
  },

  switchContext: async (payload: { targetContext: string; targetId?: string }) => {
    await apiPost("/context/switch", payload);
    fetchPromise = null;
    return await fetchSession();
  },

  forgotPassword: async (email: string) => {
    await refreshCsrfToken();
    return await apiPost<any>("/auth/forgot-password", { email });
  },

  setPassword: async (payload: { email?: string; code?: string; token?: string; password: string }) => {
    await refreshCsrfToken();
    return await apiPost<any>("/auth/set-password", payload);
  },

  refreshSession: async () => {
    fetchPromise = null;
    return await fetchSession();
  },
};

