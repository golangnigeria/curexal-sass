import React, { createContext, useContext, useEffect } from "react";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { BrandingPayload } from "@/api/contracts";

interface BrandThemeContextType {
  branding?: BrandingPayload;
  logoUrl?: string;
  orgName?: string;
  primaryColor: string;
  isCustomBranded: boolean;
}

const BrandThemeContext = createContext<BrandThemeContextType>({
  primaryColor: "#0284c7",
  isCustomBranded: false,
});

export function useBrandTheme() {
  return useContext(BrandThemeContext);
}

// Convert Hex color to contrast foreground
function getAccessibleTextColor(hexColor: string): string {
  let c = hexColor.substring(1);
  if (c.length === 3) {
    c = c.split("").map((x) => x + x).join("");
  }
  const num = parseInt(c, 16);
  if (isNaN(num)) return "#ffffff";
  const r = (num >> 16) & 255;
  const g = (num >> 8) & 255;
  const b = num & 255;
  // Perceptive luminance formula
  const yiq = (r * 299 + g * 587 + b * 114) / 1000;
  return yiq >= 140 ? "#0f172a" : "#ffffff";
}

// Convert Hex to HSL string representation for Tailwind compat
function hexToHslString(hex: string): string {
  let c = hex.replace("#", "");
  if (c.length === 3) c = c.split("").map((x) => x + x).join("");
  const num = parseInt(c, 16);
  if (isNaN(num)) return "204 90% 40%";
  const r = ((num >> 16) & 255) / 255;
  const g = ((num >> 8) & 255) / 255;
  const b = (num & 255) / 255;

  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  let h = 0;
  let s = 0;
  const l = (max + min) / 2;

  if (max !== min) {
    const d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
    switch (max) {
      case r:
        h = (g - b) / d + (g < b ? 6 : 0);
        break;
      case g:
        h = (b - r) / d + 2;
        break;
      case b:
        h = (r - g) / d + 4;
        break;
    }
    h /= 6;
  }
  return `${Math.round(h * 360)} ${Math.round(s * 100)}% ${Math.round(l * 100)}%`;
}

export function BrandThemeProvider({ children }: { children: React.ReactNode }) {
  const { data: bootstrap } = useBootstrap();
  const branding = bootstrap?.branding;
  const orgName = bootstrap?.organization?.name;
  const isCustomBranded = Boolean(branding?.primaryColor || branding?.logoUrl);

  const primaryColor = branding?.primaryColor || "#0284c7";
  const logoUrl = branding?.logoUrl || (branding?.themeBranding?.logoUrl as string);

  useEffect(() => {
    if (!branding) return;

    const root = document.documentElement;

    // 1. Inject Dynamic Primary & Semantic Color Tokens
    if (primaryColor) {
      root.style.setProperty("--brand-primary", primaryColor);
      root.style.setProperty("--primary", hexToHslString(primaryColor));

      const fgColor = getAccessibleTextColor(primaryColor);
      root.style.setProperty("--primary-foreground", hexToHslString(fgColor));
      root.style.setProperty("--sidebar-primary", hexToHslString(primaryColor));
      root.style.setProperty("--sidebar-primary-foreground", hexToHslString(fgColor));
    }

    // 2. Inject Secondary / Accent
    if (branding.secondaryColor) {
      root.style.setProperty("--brand-secondary", branding.secondaryColor);
    }
    if (branding.accentColor) {
      root.style.setProperty("--brand-accent", branding.accentColor);
    }

    // 3. Inject Typography & Border Radius
    if (branding.fontFamily) {
      root.style.setProperty("--font-sans", `'${branding.fontFamily}', -apple-system, BlinkMacSystemFont, sans-serif`);
    }
    if (branding.borderRadius) {
      root.style.setProperty("--radius", branding.borderRadius);
    }

    // 4. Update Favicon if provided
    if (branding.faviconUrl) {
      let link = document.querySelector("link[rel~='icon']") as HTMLLinkElement;
      if (!link) {
        link = document.createElement("link");
        link.rel = "icon";
        document.head.appendChild(link);
      }
      link.href = branding.faviconUrl;
    }

    // 5. Update Document Title if in Organization Context
    if (bootstrap?.contexts?.current === "organization" && orgName) {
      document.title = `${orgName} | Healthcare Console`;
    }
  }, [branding, bootstrap, orgName, primaryColor]);

  return (
    <BrandThemeContext.Provider
      value={{
        branding,
        logoUrl,
        orgName,
        primaryColor,
        isCustomBranded,
      }}
    >
      {children}
    </BrandThemeContext.Provider>
  );
}
