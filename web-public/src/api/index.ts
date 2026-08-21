import { env } from "@/config/env";
import { apiContract } from "./contracts";
import { initClient } from "@ts-rest/core";
import axios, {
  type Method,
  type AxiosError,
  isAxiosError,
  type AxiosResponse,
} from "axios";

type Headers = Awaited<
  ReturnType<NonNullable<Parameters<typeof initClient>[1]["api"]>>
>["headers"];

export type TApiClient = ReturnType<typeof useApiClient>;

export const getApiUrl = (path: string = ""): string => {
  const baseUrl = env.VITE_API_URL;
  return `${baseUrl}/api/v1${path}`;
};

export let csrfToken = "";
export const setCsrfToken = (val: string) => {
  csrfToken = val;
};

export const useApiClient = ({ isBlob = false }: { isBlob?: boolean } = {}) => {
  return initClient(apiContract, {
    baseUrl: "",
    baseHeaders: {
      "Content-Type": "application/json",
    },
    api: async ({ path, method, headers, body }) => {
      try {
        const reqHeaders = { ...headers } as any;
        if (csrfToken && ["post", "put", "patch", "delete"].includes(method.toLowerCase())) {
          reqHeaders["X-CSRF-Token"] = csrfToken;
        }

        const result = await axios.request({
          method: method as Method,
          url: getApiUrl(path),
          headers: reqHeaders,
          data: body,
          withCredentials: true,
          ...(isBlob ? { responseType: "blob" } : {}),
        });
        return {
          status: result.status,
          body: result.data,
          headers: result.headers as unknown as Headers,
        };
      } catch (e: any) {
        if (isAxiosError(e)) {
          const error = e as AxiosError;
          const response = error.response as AxiosResponse;
          return {
            status: response?.status || 500,
            body: response?.data || { message: "Internal server error" },
            headers: (response?.headers as unknown as Headers) || {},
          };
        }
        throw e;
      }
    },
  });
};
