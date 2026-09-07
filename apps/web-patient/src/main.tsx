import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider } from "@curexal/design-system";
import { BrandThemeProvider } from "@curexal/design-system";
import { Toaster } from "@curexal/ui";
import { router } from "./router";
import "./index.css";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5,
      retry: 1,
    },
  },
});

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <ThemeProvider attribute="class" defaultTheme="dark" enableSystem>
        <BrandThemeProvider>
          <RouterProvider router={router} />
          <Toaster position="top-right" richColors />
        </BrandThemeProvider>
      </ThemeProvider>
    </QueryClientProvider>
  </React.StrictMode>
);
