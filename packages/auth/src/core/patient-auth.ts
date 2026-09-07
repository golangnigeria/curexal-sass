import { useEffect, useState } from "react";
import { apiGet, apiPost, refreshCsrfToken } from "@curexal/api-client";
import type { CanonicalPatient } from "@curexal/contracts";

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
let globalPatientIsPending = true;
const patientListeners = new Set<(session: { data: PatientSessionData | null; isPending: boolean }) => void>();

const notifyPatient = () => {
  patientListeners.forEach((l) => l({ data: globalPatientSession, isPending: globalPatientIsPending }));
};

let patientFetchPromise: Promise<PatientSessionData | null> | null = null;

export const fetchPatientSession = async (): Promise<PatientSessionData | null> => {
  if (patientFetchPromise) return patientFetchPromise;

  patientFetchPromise = (async () => {
    try {
      await refreshCsrfToken();
      const body = await apiGet<any>("/users/me");
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
    } catch {
      globalPatientSession = null;
    } finally {
      globalPatientIsPending = false;
      notifyPatient();
    }
    return globalPatientSession;
  })();

  return patientFetchPromise;
};

// Initial trigger
fetchPatientSession();

export const patientAuthClient = {
  usePatientSession: () => {
    const [state, setState] = useState({ data: globalPatientSession, isPending: globalPatientIsPending });

    useEffect(() => {
      const listener = (s: { data: PatientSessionData | null; isPending: boolean }) => setState(s);
      patientListeners.add(listener);

      if (!globalPatientIsPending && (state.data !== globalPatientSession || state.isPending !== globalPatientIsPending)) {
        setState({ data: globalPatientSession, isPending: globalPatientIsPending });
      }

      return () => {
        patientListeners.delete(listener);
      };
    }, [state.data, state.isPending]);

    return state;
  },

  getSession: async () => {
    if (globalPatientIsPending) {
      return { data: await fetchPatientSession() };
    }
    return { data: globalPatientSession };
  },

  signIn: async (payload: { email: string; password: string }): Promise<PatientSessionData> => {
    await refreshCsrfToken();
    const body = await apiPost<any>("/auth/sign-in", payload);
    patientFetchPromise = null;
    const session = await fetchPatientSession();
    return session || {
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
  },

  register: async (payload: { name: string; email: string; password: string; phone?: string }): Promise<{ message: string }> => {
    await refreshCsrfToken();
    return await apiPost<{ message: string }>("/auth/register", payload);
  },

  signOut: async (redirectPath: string = "/portal/login") => {
    try {
      await apiPost("/auth/sign-out");
    } catch (err) {
      console.error("Patient sign out failed:", err);
    }
    globalPatientSession = null;
    globalPatientIsPending = false;
    notifyPatient();
    if (typeof window !== "undefined") {
      window.location.href = redirectPath;
    }
  },

  refreshSession: async () => {
    patientFetchPromise = null;
    await fetchPatientSession();
  },
};
