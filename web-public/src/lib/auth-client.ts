import { useEffect, useState } from "react";
import { getApiUrl } from "../api";

interface SessionUser {
  id: string;
  name: string;
  email: string;
  image: string | null;
  role: string;
  tenantId?: string;
  tenantSlug?: string;
  availableTenants?: Array<{ id: string; name: string; slug: string }>;
}

interface SessionData {
  user: SessionUser;
}

let globalSession: SessionData | null = null;
let globalIsPending = true;
const listeners = new Set<(session: { data: SessionData | null; isPending: boolean }) => void>();

const notify = () => {
  listeners.forEach((l) => l({ data: globalSession, isPending: globalIsPending }));
};

let fetchPromise: Promise<SessionData | null> | null = null;



const fetchSession = async (): Promise<SessionData | null> => {
  if (fetchPromise) return fetchPromise;

  fetchPromise = (async () => {
    try {
      const res = await fetch(getApiUrl("/users/me"), {
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      });
      if (res.ok) {
        const body = await res.json();
        if (body && body.id) {
          globalSession = {
            user: {
              id: body.id,
              name: body.name,
              email: body.email,
              image: body.image || null,
              role: body.role || "member",
              tenantId: body.activeTenantId || undefined,
              tenantSlug: body.tenantSlug || undefined,
              availableTenants: body.availableTenants || [],
            },
          };

          // Fetch and store the CSRF token in memory
          try {
            const csrfRes = await fetch(getApiUrl("/auth/csrf"), {
              credentials: "include",
            });
            if (csrfRes.ok) {
              const csrfBody = await csrfRes.json();
              if (csrfBody && csrfBody.csrfToken) {
                const apiModule = await import("../api");
                apiModule.setCsrfToken(csrfBody.csrfToken);
              }
            }
          } catch (csrfErr) {
            console.error("Failed to fetch CSRF token:", csrfErr);
          }
        } else {
          globalSession = null;
        }
      } else {
        globalSession = null;
      }
    } catch (err) {
      console.error("Failed to fetch session:", err);
      globalSession = null;
    } finally {
      globalIsPending = false;
      notify();
    }
    return globalSession;
  })();
  return fetchPromise;
};

// Start fetching session on load
fetchSession();

export const authClient = {
  useSession: () => {
    const [state, setState] = useState({ data: globalSession, isPending: globalIsPending });

    useEffect(() => {
      const listener = (s: { data: SessionData | null; isPending: boolean }) => setState(s);
      listeners.add(listener);

      // If the fetch already completed, make sure we have the latest state
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

  signOut: async () => {
    try {
      await fetch(getApiUrl("/auth/sign-out"), {
        method: "POST",
        credentials: "include",
      });
    } catch (err) {
      console.error("Sign out request failed:", err);
    }
    globalSession = null;
    globalIsPending = false;
    notify();
    window.location.href = "/";
  },

  // Helper to manually trigger refresh after successful Sign In / Sign Up
  refreshSession: async () => {
    fetchPromise = null;
    await fetchSession();
  }
};
