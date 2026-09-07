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
  Stethoscope,
  Microscope,
  Pill,
  CalendarPlus,
  LucideIcon,
  ChevronRight,
  ChevronLeft,
  PanelLeftClose,
  ArrowLeft,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useNavigation } from "@/api/hooks/use-navigation";
import { useSidebar } from "@/components/layout/sidebar-context";
import { CurexalLogoSymbol } from "@/components/brand/curexal-logo";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { ROLES, CONTEXT_SCOPES } from "@/api/contracts";
import type { NavigationItem } from "@curexal/contracts";
import { useBrandTheme } from "@/lib/theme/brand-theme-provider";

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
  Stethoscope,
  Microscope,
  Pill,
  CalendarPlus,
};

function getFacilityTypeIcon(typeStr?: string) {
  const norm = (typeStr || "").toLowerCase();
  if (norm.includes("lab")) return Microscope;
  if (norm.includes("clinic")) return Stethoscope;
  if (norm.includes("pharmacy")) return Pill;
  if (norm.includes("hospital")) return Building2;
  return Activity;
}

function getRolePresentationBadge(role?: string, context?: string): string {
  if (!role) return context === CONTEXT_SCOPES.PLATFORM ? "Platform Staff" : "Staff Member";
  switch (role.toLowerCase()) {
    case ROLES.SUPER_ADMIN:
      return "Super Admin";
    case ROLES.PLATFORM_ADMIN:
      return "Platform Admin";
    case ROLES.PLATFORM_STAFF:
      return "Platform Ops";
    case ROLES.OWNER:
      return "Organization Owner";
    case ROLES.ORG_ADMIN:
      return "HQ Administrator";
    case ROLES.ORG_REGIONAL_MANAGER:
      return "Regional Director";
    case ROLES.ORG_QUALITY_MANAGER:
      return "Quality Director";
    case ROLES.ORG_FINANCE_MANAGER:
      return "Finance Director";
    case ROLES.ORG_HR_MANAGER:
      return "HR Director";
    case ROLES.BRANCH_ADMIN:
      return "Branch Director";
    case ROLES.CLINICIAN:
      return "Lead Clinician";
    case ROLES.TECHNICIAN:
      return "Medical Lab Scientist";
    case ROLES.CASHIER:
      return "Billing Specialist";
    case ROLES.CUSTOMER_CARE:
      return "Care Coordinator";
    default:
      return context === CONTEXT_SCOPES.ORGANIZATION ? "Organization HQ" : "Console";
  }
}

export function PlatformSidebar() {
  const location = useLocation();
  const { data: bootstrap } = useBootstrap();
  const { isCollapsed, toggleCollapse } = useSidebar();
  const { logoUrl } = useBrandTheme();

  const isPlatformContext = location.pathname.startsWith("/platform");
  const isOrgContext = !isPlatformContext && location.pathname.startsWith("/organization");
  const isWorkspaceContext = !isPlatformContext && !location.pathname.startsWith("/organization") && location.pathname !== "/login";

  const pathParts = location.pathname.split("/").filter(Boolean);
  const activeBranchSlug = pathParts[0] || bootstrap?.branch?.slug || bootstrap?.workspace?.slug || "main";

  const activeScope = isPlatformContext ? "platform" : isOrgContext ? "organization" : "workspace";
  const { data: navData } = useNavigation({
    scope: activeScope,
    branchSlug: isWorkspaceContext ? activeBranchSlug : undefined,
  });

  const orgName = bootstrap?.organization?.name || "Curexal";
  const activeBranch =
    bootstrap?.availableBranches?.find(
      (b: any) => b.slug === activeBranchSlug || b.code?.toLowerCase() === activeBranchSlug.toLowerCase()
    ) ||
    bootstrap?.branch ||
    bootstrap?.workspace;
  const facilityName = activeBranch?.name || activeBranchSlug.toUpperCase() + " Facility";
  const facilityType = activeBranch?.facilityType || "Diagnostic Facility";
  const FacilityIcon = getFacilityTypeIcon(facilityType);

  const isOrgAdminOrOwner =
    bootstrap?.platform?.isStaff === true ||
    bootstrap?.organization?.role === "owner" ||
    bootstrap?.organization?.role === "org_admin" ||
    bootstrap?.organization?.role === "org_regional_manager" ||
    bootstrap?.contexts?.available?.includes("organization") ||
    (bootstrap?.identity as any)?.role === "owner" ||
    (bootstrap?.identity as any)?.role === "org_admin";

  const activeRole = isPlatformContext
    ? (bootstrap?.platform?.role || (bootstrap?.platform?.isStaff ? "platform_staff" : undefined) || "super_admin")
    : (bootstrap?.organization?.role || bootstrap?.platform?.role);

  const roleBadge = isWorkspaceContext
    ? facilityType.replace(/_/g, " ")
    : getRolePresentationBadge(
        activeRole,
        isPlatformContext ? CONTEXT_SCOPES.PLATFORM : bootstrap?.contexts?.current
      );

  const homePath = isPlatformContext
    ? "/platform/dashboard"
    : isOrgContext
    ? "/organization/dashboard"
    : `/${activeBranchSlug}/dashboard`;

  // 100% Database-driven navigation items derived from backend Navigation API or Bootstrap contract
  const backendNavigation: NavigationItem[] = React.useMemo(() => {
    let rawItems: NavigationItem[] = [];

    // 1. Primary: Use dedicated /api/v1/navigation response items
    if (navData?.items && navData.items.length > 0) {
      rawItems = navData.items;
    } else if (bootstrap?.structuredNavigation?.primary?.length) {
      // 2. Secondary: Fallback to bootstrap contract navigation payload
      rawItems = bootstrap.structuredNavigation.primary.map((item): NavigationItem => {
        let p = item.path;
        if (isWorkspaceContext && p.includes("/:branch")) {
          p = p.replace("/:branch", `/${activeBranchSlug}`);
        }
        return {
          id: item.id,
          key: item.id,
          contextScope: isPlatformContext ? "platform" : isOrgContext ? "organization" : "workspace",
          title: item.title,
          icon: item.icon,
          path: p,
          order: item.order || 0,
          status: "active",
          isVisible: true,
          isActive: true,
          badgeCount: undefined,
          children: item.children as any,
        };
      });
    } else if (bootstrap?.navigation?.length) {
      rawItems = bootstrap.navigation.map((item): NavigationItem => {
        let p = item.path;
        if (isWorkspaceContext && p.includes("/:branch")) {
          p = p.replace("/:branch", `/${activeBranchSlug}`);
        }
        return {
          id: item.id,
          key: item.id,
          contextScope: isPlatformContext ? "platform" : isOrgContext ? "organization" : "workspace",
          title: item.title,
          icon: item.icon,
          path: p,
          order: item.order || 0,
          status: "active",
          isVisible: true,
          isActive: true,
          badgeCount: undefined,
          children: item.children as any,
        };
      });
    }

    // Defensive facility-type pruning to prevent cross-facility leakage during initial loading/transition
    if (isWorkspaceContext) {
      const normFT = (facilityType || "").toLowerCase();
      return rawItems.filter((item) => {
        const p = item.path.toLowerCase();
        // Outpatient Clinic: strictly prohibit LIS, RIS/PACS, and Inpatient HMS
        if (normFT.includes("clinic") || normFT.includes("outpatient")) {
          if (p.includes("/laboratory") || p.includes("/radiology") || p.includes("/hospital")) {
            return false;
          }
        } else if (normFT.includes("lab") && !normFT.includes("hospital")) {
          if (p.includes("/clinical") || p.includes("/radiology") || p.includes("/hospital") || p.includes("/pharmacy")) {
            return false;
          }
        } else if (normFT.includes("pharmacy") && !normFT.includes("hospital")) {
          if (p.includes("/clinical") || p.includes("/laboratory") || p.includes("/radiology") || p.includes("/hospital")) {
            return false;
          }
        } else if (normFT.includes("radiology") && !normFT.includes("hospital")) {
          if (p.includes("/clinical") || p.includes("/laboratory") || p.includes("/pharmacy") || p.includes("/hospital")) {
            return false;
          }
        }
        return true;
      });
    }

    return rawItems;
  }, [navData, bootstrap, isPlatformContext, isOrgContext, isWorkspaceContext, activeBranchSlug, facilityType]);

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
            className="flex items-center gap-3 overflow-hidden group min-w-0"
          >
            <div className="relative flex items-center justify-center p-1 rounded-xl bg-slate-900 shadow-md shrink-0 group-hover:scale-105 transition-transform">
              {isWorkspaceContext ? (
                <div className="w-7 h-7 flex items-center justify-center rounded-lg bg-primary/20 text-primary">
                  <FacilityIcon className="w-4 h-4" />
                </div>
              ) : isOrgContext && logoUrl ? (
                <img src={logoUrl} alt={orgName} className="w-7 h-7 object-contain rounded-lg" />
              ) : (
                <CurexalLogoSymbol className="w-7 h-7" />
              )}
            </div>
            {!isCollapsed && (
              <div className="flex flex-col truncate min-w-0">
                <span className="font-extrabold tracking-tight text-foreground text-sm leading-tight truncate">
                  {isWorkspaceContext ? facilityName : isOrgContext ? orgName : "CUREXAL"}
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

        {/* Quick Return to Org HQ button inside sidebar for Org Admins in Workspace */}
        {isWorkspaceContext && isOrgAdminOrOwner && !isCollapsed && (
          <div className="px-3 pt-3">
            <Button
              asChild
              variant="outline"
              size="sm"
              className="w-full justify-start text-xs h-8 gap-2 border-border/80 bg-secondary/40 hover:bg-secondary text-foreground"
            >
              <Link to="/organization/dashboard">
                <ArrowLeft className="w-3.5 h-3.5 text-primary" />
                <span>Return to Executive HQ</span>
              </Link>
            </Button>
          </div>
        )}

        {/* Navigation Links strictly rendered from database SSOT */}
        <div className="flex-1 overflow-y-auto px-3 py-4 space-y-1.5 scrollbar-thin">
          {!isCollapsed && (
            <div className="px-3 pb-2 text-[10px] font-bold uppercase tracking-wider text-muted-foreground/70">
              {isWorkspaceContext
                ? "Clinical & Facility Operations"
                : isOrgContext
                ? "Executive Management & Operations"
                : "Operations & Management"}
            </div>
          )}

          {backendNavigation.map((item) => {
            const IconComponent = iconMap[item.icon] || Activity;
            const targetPath = item.path;
            const isRootDashboard =
              targetPath === "/platform/dashboard" ||
              targetPath === "/organization/dashboard" ||
              targetPath === `/${activeBranchSlug}/dashboard`;
            const isActive =
              location.pathname === targetPath ||
              (!isRootDashboard && location.pathname.startsWith(targetPath));

            const isPending = item.status === "pending";

            const linkElement = (
              <Link
                key={item.id || item.path}
                to={isPending ? "#" : targetPath}
                onClick={isPending ? (e) => e.preventDefault() : undefined}
                className={cn(
                  "group relative flex items-center rounded-lg text-sm font-medium transition-all duration-200",
                  isCollapsed
                    ? "justify-center h-11 w-11 mx-auto p-0"
                    : "gap-3 px-3 py-2.5",
                  isActive
                    ? "bg-primary text-primary-foreground shadow-sm font-semibold"
                    : isPending
                    ? "opacity-60 cursor-not-allowed text-muted-foreground"
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
                    {isPending ? (
                      <Badge variant="outline" className="text-[9px] px-1 py-0 border-amber-500/50 text-amber-600 dark:text-amber-400">
                        Pending
                      </Badge>
                    ) : item.badgeCount !== undefined && item.badgeCount > 0 ? (
                      <Badge variant="secondary" className="text-[10px] px-1.5 py-0 bg-primary/20 text-primary">
                        {item.badgeCount}
                      </Badge>
                    ) : isActive ? (
                      <ChevronRight className="h-3.5 w-3.5 opacity-80" />
                    ) : null}
                  </>
                )}
              </Link>
            );

            if (isCollapsed) {
              return (
                <Tooltip key={item.id || item.path}>
                  <TooltipTrigger asChild>{linkElement}</TooltipTrigger>
                  <TooltipContent side="right" className="text-xs font-medium">
                    {item.title} {isPending && "(Pending)"}
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
                <ChevronRight className="h-4 w-4" />
              </Button>
              <span className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" title="Live Cluster" />
            </div>
          ) : (
            <>
              <div className="flex flex-col">
                <span className="font-mono text-[11px] text-muted-foreground">
                  {bootstrap?.metadata?.version || "Cluster v1.0.0"}
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
