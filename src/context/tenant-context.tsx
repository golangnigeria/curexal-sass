import React, { createContext, useContext } from "react";

/**
 * Portal Tenant Context
 *
 * The portal (`app.curexal.com`) is the organization control plane.
 * Unlike the tenant workspace, it does NOT require a tenant to be resolved
 * from the subdomain before rendering. The tenant info is fetched lazily
 * from the user's session (via their tenantId).
 *
 * This context provides a simple nullable Tenant interface that pages
 * can optionally consume.
 */

interface TenantBranding {
  primaryColor: string;
  secondaryColor: string;
  fontFamily: string;
  orgType?: "hospital" | "laboratory" | "pharmacy";
}

export interface Tenant {
  id: string;
  name: string;
  slug: string;
  logoUrl?: string;
  branding?: TenantBranding;
  themeBranding?: string;
}

const TenantContext = createContext<Tenant | null>(null);

/**
 * Portal TenantProvider — passthrough wrapper.
 * Does NOT block rendering. Individual pages/components fetch tenant
 * data as needed from the session or API.
 */
export const TenantProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <TenantContext.Provider value={null}>
      {children}
    </TenantContext.Provider>
  );
};

export const useTenant = () => useContext(TenantContext);
