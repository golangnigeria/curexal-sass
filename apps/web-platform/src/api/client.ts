import axios, { type AxiosRequestConfig } from "axios";
import { env } from "@/config/env";

export const getApiUrl = (path: string = ""): string => {
  const baseUrl = env.VITE_API_URL;
  return `${baseUrl}/api/v1${path.startsWith("/") ? path : `/${path}`}`;
};

let inMemoryCsrfToken = "";

export const setCsrfToken = (token: string) => {
  inMemoryCsrfToken = token;
};

export const getCsrfToken = (): string => inMemoryCsrfToken;

export const apiClient = axios.create({
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  // Prepend API v1 URL if relative
  if (config.url && !config.url.startsWith("http")) {
    config.url = getApiUrl(config.url);
  }

  // Attach CSRF token on mutating HTTP methods
  const method = config.method?.toLowerCase();
  if (inMemoryCsrfToken && method && ["post", "put", "patch", "delete"].includes(method)) {
    config.headers["X-CSRF-Token"] = inMemoryCsrfToken;
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
    console.warn("Could not fetch CSRF token:", err);
  }
  return "";
}

// Generic API fetchers
export async function apiGet<T>(path: string, config?: AxiosRequestConfig): Promise<T> {
  const res = await apiClient.get<T>(path, config);
  // Support both direct return and response.Success wrapper ({ data: T })
  if (res.data && typeof res.data === "object" && "data" in res.data && !Array.isArray(res.data)) {
    return (res.data as any).data as T;
  }
  return res.data;
}

export async function apiPost<T, B = any>(path: string, body?: B, config?: AxiosRequestConfig): Promise<T> {
  const res = await apiClient.post<T>(path, body, config);
  if (res.data && typeof res.data === "object" && "data" in res.data && !Array.isArray(res.data)) {
    return (res.data as any).data as T;
  }
  return res.data;
}

export async function apiPut<T, B = any>(path: string, body?: B, config?: AxiosRequestConfig): Promise<T> {
  const res = await apiClient.put<T>(path, body, config);
  if (res.data && typeof res.data === "object" && "data" in res.data && !Array.isArray(res.data)) {
    return (res.data as any).data as T;
  }
  return res.data;
}

export async function apiPatch<T, B = any>(path: string, body?: B, config?: AxiosRequestConfig): Promise<T> {
  const res = await apiClient.patch<T>(path, body, config);
  if (res.data && typeof res.data === "object" && "data" in res.data && !Array.isArray(res.data)) {
    return (res.data as any).data as T;
  }
  return res.data;
}

export async function apiDelete<T>(path: string, config?: AxiosRequestConfig): Promise<T> {
  const res = await apiClient.delete<T>(path, config);
  if (res.data && typeof res.data === "object" && "data" in res.data && !Array.isArray(res.data)) {
    return (res.data as any).data as T;
  }
  return res.data;
}
