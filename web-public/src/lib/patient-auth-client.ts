import { useEffect, useState } from "react";
import { getApiUrl } from "../api";

export interface PatientContext {
  profileId: string;
  phone: string | null;
  dateOfBirth: string | null;
  gender: string | null;
  bloodGroup: string | null;
  genotype: string | null;
  city: string | null;
  state: string | null;
  country: string;
}

export interface PatientSessionUser {
  id: string;
  name: string;
  email: string;
  emailVerified: boolean;
  image: string | null;
  role?: string;
  patient: PatientContext | null;
}

export interface PatientSessionData {
  user: PatientSessionUser;
  patient: PatientContext | null;
}

let globalPatientSession: PatientSessionData | null = null;
let globalIsPending = true;
const listeners = new Set<(session: { data: PatientSessionData | null; isPending: boolean }) => void>();

const notify = () => {
  listeners.forEach((l) => l({ data: globalPatientSession, isPending: globalIsPending }));
};

let fetchPromise: Promise<PatientSessionData | null> | null = null;

const fetchPatientSession = async (): Promise<PatientSessionData | null> => {
  if (fetchPromise) return fetchPromise;

  fetchPromise = (async () => {
    try {
      const res = await fetch(getApiUrl("/users/me"), {
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      });
      if (res.ok) {
        const body = await res.json();
        const userId = body?.id || body?.identity?.user?.id;
        if (body && userId) {
          globalPatientSession = {
            user: {
              id: userId,
              name: body.name || body.identity?.user?.name || "",
              email: body.email || body.identity?.user?.email || "",
              emailVerified: body.emailVerified ?? body.identity?.user?.emailVerified ?? false,
              image: body.image || body.identity?.user?.avatarUrl || null,
              role: body.role || "member",
              patient: body.patient || null,
            },
            patient: body.patient || null,
          };
        } else {
          globalPatientSession = null;
        }
      } else {
        globalPatientSession = null;
      }
    } catch (err) {
      console.error("Failed to fetch patient session:", err);
      globalPatientSession = null;
    } finally {
      globalIsPending = false;
      notify();
    }
    return globalPatientSession;
  })();
  return fetchPromise;
};

// Start fetching session on load
fetchPatientSession();

export const patientAuthClient = {
  usePatientSession: () => {
    const [state, setState] = useState({ data: globalPatientSession, isPending: globalIsPending });

    useEffect(() => {
      const listener = (s: { data: PatientSessionData | null; isPending: boolean }) => setState(s);
      listeners.add(listener);

      if (!globalIsPending && (state.data !== globalPatientSession || state.isPending !== globalIsPending)) {
        setState({ data: globalPatientSession, isPending: globalIsPending });
      }

      return () => {
        listeners.delete(listener);
      };
    }, [state.data, state.isPending]);

    return state;
  },

  getSession: async () => {
    if (globalIsPending) {
      return { data: await fetchPatientSession() };
    }
    return { data: globalPatientSession };
  },

  signIn: async (email: string, password: string): Promise<PatientSessionData> => {
    const res = await fetch(getApiUrl("/auth/sign-in"), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ email, password }),
    });

    if (!res.ok) {
      const errData = await res.json().catch(() => ({}));
      throw new Error(errData.message || errData.error || "Sign in failed");
    }

    const body = await res.json();
    globalPatientSession = {
      user: {
        id: body.id,
        name: body.name,
        email: body.email,
        emailVerified: body.emailVerified,
        image: body.image || null,
        role: body.role || "member",
        patient: body.patient || null,
      },
      patient: body.patient || null,
    };
    globalIsPending = false;
    notify();
    return globalPatientSession;
  },

  register: async (payload: { name: string; email: string; password: string; phone?: string }): Promise<{ message: string }> => {
    const res = await fetch(getApiUrl("/auth/register"), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify(payload),
    });

    if (!res.ok) {
      const errData = await res.json().catch(() => ({}));
      throw new Error(errData.message || errData.error || "Registration failed");
    }

    return await res.json();
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
    globalPatientSession = null;
    globalIsPending = false;
    notify();
    window.location.href = "/login";
  },

  refreshSession: async () => {
    fetchPromise = null;
    await fetchPatientSession();
  },
};
