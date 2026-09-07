import axios, { type AxiosRequestConfig, isAxiosError } from "axios";
import { ApiError } from "./errors";

let customBaseUrl: string | null = null;
let inMemoryCsrfToken = "";
let inMemoryActiveOrgId = "";
let inMemoryAuthToken = "";

export const setApiBaseUrl = (url: string) => {
  customBaseUrl = url;
};

export const setAuthToken = (token: string) => {
  inMemoryAuthToken = token;
};

export const getAuthToken = (): string => inMemoryAuthToken;

export const getApiUrl = (path: string = ""): string => {
  const g = typeof globalThis !== "undefined" ? (globalThis as any) : null;
  const envBase = (typeof import.meta !== "undefined" && (import.meta as any)?.env?.VITE_API_URL) ||
    (g?.process?.env?.VITE_API_URL) ||
    (typeof window !== "undefined" && (window as any).__VITE_API_URL__) ||
    "";

  const baseUrl = customBaseUrl || envBase || "";
  const cleanPath = path.startsWith("/") ? path : `/${path}`;

  if (baseUrl) {
    return `${baseUrl}/api/v1${cleanPath}`;
  }
  return `/api/v1${cleanPath}`;
};

export const setCsrfToken = (token: string) => {
  inMemoryCsrfToken = token;
};

export const getCsrfToken = (): string => inMemoryCsrfToken;

export const setClientActiveOrgId = (orgId: string) => {
  inMemoryActiveOrgId = orgId;
};

export const getClientActiveOrgId = (): string => inMemoryActiveOrgId;

export const apiClient = axios.create({
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  // Prepend API v1 URL if relative and not already prefixed
  if (config.url && !config.url.startsWith("http")) {
    config.url = getApiUrl(config.url);
  }

  // Attach Authorization Bearer token if present
  if (inMemoryAuthToken && !config.headers["Authorization"]) {
    config.headers["Authorization"] = `Bearer ${inMemoryAuthToken}`;
  }

  // Attach CSRF token on mutating HTTP methods
  const method = config.method?.toLowerCase();
  if (inMemoryCsrfToken && method && ["post", "put", "patch", "delete"].includes(method)) {
    config.headers["X-CSRF-Token"] = inMemoryCsrfToken;
  }

  // Attach active organization headers if resolved in client memory
  if (inMemoryActiveOrgId) {
    config.headers["X-Organization-ID"] = inMemoryActiveOrgId;
    config.headers["X-Tenant-ID"] = inMemoryActiveOrgId;
  }

  // Attach active branch slug if in branch workspace route
  if (typeof window !== "undefined" && window.location?.pathname) {
    const parts = window.location.pathname.split("/").filter(Boolean);
    if (
      parts.length > 0 &&
      parts[0] !== "login" &&
      parts[0] !== "platform" &&
      parts[0] !== "organization" &&
      parts[0] !== "workspace" &&
      parts[0] !== "api" &&
      parts[0] !== "portal"
    ) {
      config.headers["X-Branch-Slug"] = parts[0];
    }
  }

  return config;
});

// Fetch and store CSRF token helper
export async function refreshCsrfToken(): Promise<string> {
  try {
    const res = await axios.get<{ csrfToken: string }>(getApiUrl("/auth/csrf"), {
      withCredentials: true,
    });
    if (res.data?.csrfToken) {
      setCsrfToken(res.data.csrfToken);
      return res.data.csrfToken;
    }
  } catch (err) {
    // CSRF fetch warning
  }
  return "";
}

function handleAxiosError(err: unknown): never {
  if (isAxiosError(err)) {
    const status = err.response?.status || 500;
    const body = err.response?.data;
    const message = body?.message || body?.error || err.message || "API request failed";
    const code = body?.code;
    const requestId = err.response?.headers?.["x-request-id"];
    throw new ApiError(message, status, { code, details: body, requestId });
  }
  if (err instanceof Error) {
    throw new ApiError(err.message, 500);
  }
  throw new ApiError("Unknown API Error", 500);
}

// Generic API fetchers with automatic response unwrapping and error normalization
export async function apiGet<T>(path: string, config?: AxiosRequestConfig): Promise<T> {
  try {
    const res = await apiClient.get<T>(path, config);
    if (res.data && typeof res.data === "object" && "data" in res.data && !Array.isArray(res.data)) {
      return (res.data as any).data as T;
    }
    return res.data;
  } catch (err) {
    return handleAxiosError(err);
  }
}

export async function apiPost<T, B = any>(path: string, body?: B, config?: AxiosRequestConfig): Promise<T> {
  try {
    const res = await apiClient.post<T>(path, body, config);
    if (res.data && typeof res.data === "object" && "data" in res.data && !Array.isArray(res.data)) {
      return (res.data as any).data as T;
    }
    return res.data;
  } catch (err) {
    return handleAxiosError(err);
  }
}

export async function apiPut<T, B = any>(path: string, body?: B, config?: AxiosRequestConfig): Promise<T> {
  try {
    const res = await apiClient.put<T>(path, body, config);
    if (res.data && typeof res.data === "object" && "data" in res.data && !Array.isArray(res.data)) {
      return (res.data as any).data as T;
    }
    return res.data;
  } catch (err) {
    return handleAxiosError(err);
  }
}

export async function apiPatch<T, B = any>(path: string, body?: B, config?: AxiosRequestConfig): Promise<T> {
  try {
    const res = await apiClient.patch<T>(path, body, config);
    if (res.data && typeof res.data === "object" && "data" in res.data && !Array.isArray(res.data)) {
      return (res.data as any).data as T;
    }
    return res.data;
  } catch (err) {
    return handleAxiosError(err);
  }
}

export async function apiDelete<T>(path: string, config?: AxiosRequestConfig): Promise<T> {
  try {
    const res = await apiClient.delete<T>(path, config);
    if (res.data && typeof res.data === "object" && "data" in res.data && !Array.isArray(res.data)) {
      return (res.data as any).data as T;
    }
    return res.data;
  } catch (err) {
    return handleAxiosError(err);
  }
}
