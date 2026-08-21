import React from "react";
import { Link, useLocation } from "react-router-dom";
import {
  Activity,
  Building2,
  Users,
  Store,
  Layers,
  Cpu,
  Settings,
  ShieldCheck,
  Shield,
  CreditCard,
  BookOpen,
  History,
  Inbox,
  LayoutDashboard,
  Palette,
  Bell,
  FileCheck,
  LucideIcon,
  ChevronRight,
  ChevronLeft,
  PanelLeftClose,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useSidebar } from "@/components/layout/sidebar-context";
import { CurexalLogoSymbol } from "@/components/brand/curexal-logo";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

// Map backend icon string names to Lucide icons
const iconMap: Record<string, LucideIcon> = {
  Activity,
  Building2,
  Users,
  Store,
  Layers,
  Cpu,
  Settings,
  ShieldCheck,
  Shield,
  CreditCard,
  BookOpen,
  History,
  Inbox,
  LayoutDashboard,
  Palette,
  Bell,
  FileCheck,
};

// Default canonical platform console navigation
const defaultPlatformNavigation = [
  { id: "nav_plat_dashboard", title: "Platform Dashboard", icon: "LayoutDashboard", path: "/platform/dashboard", order: 1 },
  { id: "nav_plat_orgs", title: "Organizations", icon: "Building2", path: "/platform/organizations", order: 2 },
  { id: "nav_plat_users", title: "User Directory", icon: "Users", path: "/platform/users", order: 3 },
  { id: "nav_plat_marketplace", title: "B2B Marketplace", icon: "Store", path: "/platform/marketplace", order: 4 },
  { id: "nav_plat_pricing", title: "Pricing & Billing", icon: "CreditCard", path: "/platform/pricing", order: 5 },
  { id: "nav_plat_facility_types", title: "Facility Types", icon: "Layers", path: "/platform/facility-types", order: 6 },
  { id: "nav_plat_catalogs", title: "Master Catalogs", icon: "BookOpen", path: "/platform/catalogs", order: 7 },
  { id: "nav_plat_audit", title: "Audit Trail", icon: "History", path: "/platform/audit", order: 8 },
  { id: "nav_plat_diag", title: "Diagnostics & Gate", icon: "Cpu", path: "/platform/diagnostics", order: 9 },
  { id: "nav_plat_demo", title: "Demo Requests", icon: "Inbox", path: "/platform/demo-requests", order: 10 },
  { id: "nav_plat_settings", title: "Console Settings", icon: "Settings", path: "/platform/settings", order: 11 },
];

function getRolePresentationBadge(role?: string, context?: string): string {
  if (!role && context === "organization") return "Executive HQ";
  switch (role?.toLowerCase()) {
    case "owner":
      return "Organization Owner";
    case "org_admin":
      return "Org Administrator";
    case "org_regional_manager":
      return "Regional Manager";
    case "super_admin":
      return "Super Admin";
    case "platform_admin":
      return "Platform Admin";
    case "platform_staff":
      return "Platform Staff";
    default:
      return context === "organization" ? "Organization HQ" : "Console";
  }
}

import { useBrandTheme } from "@/lib/theme/brand-theme-provider";

export function PlatformSidebar() {
  const location = useLocation();
  const { data: bootstrap } = useBootstrap();
  const { isCollapsed, toggleCollapse } = useSidebar();
  const { logoUrl } = useBrandTheme();

  const isOrgContext = bootstrap?.contexts?.current === "organization";
  const orgName = bootstrap?.organization?.name || "Curexal";
  const roleBadge = getRolePresentationBadge(
    bootstrap?.organization?.role || bootstrap?.platform?.role,
    bootstrap?.contexts?.current
  );
  const homePath = isOrgContext ? "/organization/dashboard" : "/platform/dashboard";

  // Unified navigation from backend bootstrap with fallback and route deduplication
  const backendNavigation = React.useMemo(() => {
    const rawItems =
      bootstrap?.structuredNavigation?.primary?.length
        ? bootstrap.structuredNavigation.primary
        : bootstrap?.navigation?.length
        ? bootstrap.navigation
        : defaultPlatformNavigation;

    const seenPaths = new Set<string>();
    return rawItems.filter((item) => {
      if (seenPaths.has(item.path)) return false;
      seenPaths.add(item.path);
      return true;
    });
  }, [bootstrap]);

  return (
    <TooltipProvider delayDuration={150}> 
      <aside
        className={cn(
          "fixed left-0 top-0 z-30 flex h-screen flex-col border-r border-sidebar-border bg-sidebar text-sidebar-foreground transition-all duration-300 ease-in-out",
          isCollapsed ? "w-20" : "w-64"
        )}
      >
        {/* Brand Header */}
        <div
          className={cn(
            "flex h-16 items-center border-b border-sidebar-border px-4 transition-all duration-300",
            isCollapsed ? "justify-center px-2" : "justify-between px-5"
          )}
        >
          <Link
            to={homePath}
            className="flex items-center gap-3 overflow-hidden group"
          >
            <div className="relative flex items-center justify-center p-1 rounded-xl bg-slate-900 shadow-md shrink-0 group-hover:scale-105 transition-transform">
              {isOrgContext && logoUrl ? (
                <img src={logoUrl} alt={orgName} className="w-7 h-7 object-contain rounded-lg" />
              ) : (
                <CurexalLogoSymbol className="w-7 h-7" />
              )}
            </div>
            {!isCollapsed && (
              <div className="flex flex-col truncate">
                <span className="font-extrabold tracking-tight text-foreground text-sm leading-tight truncate">
                  {isOrgContext ? orgName : "CUREXAL"}
                </span>
                <span className="text-[10px] tracking-wider font-semibold text-primary uppercase truncate">
                  {roleBadge}
                </span>
              </div>
            )}
          </Link>

          {!isCollapsed && (
            <Button
              variant="ghost"
              size="icon"
              onClick={toggleCollapse}
              className="h-8 w-8 text-muted-foreground hover:text-foreground hover:bg-sidebar-accent shrink-0"
              title="Collapse Sidebar"
            >
              <PanelLeftClose className="h-4 w-4" />
            </Button>
          )}
        </div>

        {/* Navigation Links strictly from backend */}
        <div className="flex-1 overflow-y-auto px-3 py-4 space-y-1.5 scrollbar-thin">
          {!isCollapsed && (
            <div className="px-3 pb-2 text-[10px] font-bold uppercase tracking-wider text-muted-foreground/70">
              {isOrgContext ? "Executive Management & Operations" : "Operations & Management"}
            </div>
          )}

          {backendNavigation.map((item) => {
            const IconComponent = iconMap[item.icon] || Activity;
            const targetPath = item.path;
            const isRootDashboard = targetPath === "/platform/dashboard" || targetPath === "/organization/dashboard";
            const isActive =
              location.pathname === targetPath ||
              (!isRootDashboard && location.pathname.startsWith(targetPath));

            const linkElement = (
              <Link
                key={item.id || item.path}
                to={targetPath}
                className={cn(
                  "group relative flex items-center rounded-lg text-sm font-medium transition-all duration-200",
                  isCollapsed
                    ? "justify-center h-11 w-11 mx-auto p-0"
                    : "gap-3 px-3 py-2.5",
                  isActive
                    ? "bg-primary text-primary-foreground shadow-sm font-semibold"
                    : "text-sidebar-foreground/80 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
                )}
              >
                <IconComponent
                  className={cn(
                    "shrink-0 transition-transform group-hover:scale-110",
                    isCollapsed ? "h-5 w-5" : "h-4 w-4",
                    isActive ? "text-primary-foreground" : "text-muted-foreground group-hover:text-foreground"
                  )}
                />
                {!isCollapsed && (
                  <>
                    <span className="flex-1 truncate">{item.title}</span>
                    {isActive && <ChevronRight className="h-3.5 w-3.5 opacity-80" />}
                  </>
                )}
              </Link>
            );

            if (isCollapsed) {
              return (
                <Tooltip key={item.id || item.path}>
                  <TooltipTrigger asChild>{linkElement}</TooltipTrigger>
                  <TooltipContent side="right" className="text-xs font-medium">
                    {item.title}
                  </TooltipContent>
                </Tooltip>
              );
            }

            return linkElement;
          })}
        </div>

        {/* Footer System State */}
        <div className="border-t border-sidebar-border p-3 bg-sidebar-accent/30 flex items-center justify-between">
          {isCollapsed ? (
            <div className="w-full flex flex-col items-center gap-2">
              <Button
                variant="ghost"
                size="icon"
                onClick={toggleCollapse}
                className="h-8 w-8 text-muted-foreground hover:text-foreground"
                title="Expand Sidebar"
              >
              </Button>
              <span className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" title="Live Cluster" />
            </div>
          ) : (
            <>
              <div className="flex flex-col">
                <span className="font-mono text-[11px] text-muted-foreground">
                  {bootstrap?.metadata?.version || "Cluster"}
                </span>
                <span className="flex items-center gap-1.5 font-medium text-[11px] text-emerald-600 dark:text-emerald-400">
                  <span className="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse" />
                  Live Cluster
                </span>
              </div>
              <Button
                variant="ghost"
                size="icon"
                onClick={toggleCollapse}
                className="h-7 w-7 text-muted-foreground hover:text-foreground hover:bg-sidebar-accent"
                title="Collapse Sidebar"
              >
                <ChevronLeft className="h-4 w-4" />
              </Button>
            </>
          )}
        </div>
      </aside>
    </TooltipProvider>
  );
}
