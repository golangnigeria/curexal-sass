import { z } from "zod";

const envVarsSchema = z.object({
  VITE_API_URL: z.string().default("http://localhost:8080"),
  VITE_ENV: z.enum(["production", "development", "local"]).default("local"),
  VITE_APP_DOMAIN: z.string().default("localhost"),
  VITE_PORTAL_URL: z.string().default("http://localhost:5001"),
  VITE_PATIENT_PORTAL_URL: z.string().default("http://localhost:5004"),
  VITE_TENANT_URL: z.string().default("http://localhost:5003"),
  VITE_PUBLIC_URL: z.string().default("http://localhost:5005"),
});

const metaEnv = (import.meta as any).env || {};

const rawEnv = {
  VITE_API_URL: metaEnv.VITE_API_URL || "http://localhost:8080",
  VITE_ENV: metaEnv.VITE_ENV || "local",
  VITE_APP_DOMAIN: metaEnv.VITE_APP_DOMAIN || "localhost",
  VITE_PORTAL_URL: metaEnv.VITE_PORTAL_URL || "http://localhost:5001",
  VITE_PATIENT_PORTAL_URL: metaEnv.VITE_PATIENT_PORTAL_URL || "http://localhost:5004",
  VITE_TENANT_URL: metaEnv.VITE_TENANT_URL || "http://localhost:5003",
  VITE_PUBLIC_URL: metaEnv.VITE_PUBLIC_URL || "http://localhost:5005",
};

const parseResult = envVarsSchema.safeParse(rawEnv);

export const env = parseResult.success
  ? parseResult.data
  : {
      VITE_API_URL: "http://localhost:8080",
      VITE_ENV: "local" as const,
      VITE_APP_DOMAIN: "localhost",
      VITE_PORTAL_URL: "http://localhost:5001",
      VITE_PATIENT_PORTAL_URL: "http://localhost:5004",
      VITE_TENANT_URL: "http://localhost:5003",
      VITE_PUBLIC_URL: "http://localhost:5005",
    };
