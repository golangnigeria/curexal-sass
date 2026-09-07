import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import path from "path";
import { defineConfig, loadEnv } from "vite";

export default defineConfig(({ mode }) => {
  const envDir = path.resolve(__dirname, "../../");
  const env = loadEnv(mode, envDir, "");
  return {
    envDir,
    plugins: [react(), tailwindcss()],
    server: {
      port: env.VITE_PORTAL_PORT ? parseInt(env.VITE_PORTAL_PORT, 10) : 5002,
      host: env.VITE_HOST || "::",
      allowedHosts: true,
      proxy: {
        "/api": {
          target: env.VITE_BACKEND_URL || "http://localhost:8080",
          changeOrigin: true,
          xfwd: true,
          configure: (proxy) => {
            proxy.on("proxyReq", (proxyReq, req) => {
              if (req.headers.host) {
                proxyReq.setHeader("X-Forwarded-Host", req.headers.host);
              }
            });
          },
        },
      },
    },
    build: {
      chunkSizeWarningLimit: 1200,
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (id.includes("pdfjs-dist")) {
              return "vendor-pdf";
            }
            if (id.includes("recharts")) {
              return "vendor-charts";
            }
            if (id.includes("framer-motion")) {
              return "vendor-motion";
            }
            if (id.includes("@radix-ui")) {
              return "vendor-radix";
            }
            if (id.includes("@tanstack/react-query")) {
              return "vendor-query";
            }
            if (id.includes("node_modules/react/") || id.includes("node_modules/react-dom/") || id.includes("node_modules/react-router")) {
              return "vendor-react";
            }
          },
        },
      },
    },
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
        "@curexal/contracts": path.resolve(__dirname, "../../packages/contracts/src"),
        "@curexal/utils": path.resolve(__dirname, "../../packages/utils/src"),
        "@curexal/ui": path.resolve(__dirname, "../../packages/ui/src"),
        "@curexal/design-system": path.resolve(__dirname, "../../packages/design-system/src"),
        "@curexal/api-client": path.resolve(__dirname, "../../packages/api-client/src"),
        "@curexal/auth": path.resolve(__dirname, "../../packages/auth/src"),
        "@curexal/permissions": path.resolve(__dirname, "../../packages/permissions/src"),
        "@curexal/tenant-context": path.resolve(__dirname, "../../packages/tenant-context/src"),
        "@curexal/documents": path.resolve(__dirname, "../../packages/documents/src"),
        "@curexal/patient": path.resolve(__dirname, "../../packages/patient/src"),
        "@curexal/validation": path.resolve(__dirname, "../../packages/validation/src"),
      },
    },
  };
});
