import React, { useState, useEffect } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import {
  Search,
  Activity,
  Shield,
  Layers,
  LogOut,
  Bell,
  HelpCircle,
  PanelLeftOpen,
  ChevronRight,
  Command,
} from "lucide-react";
import { authClient } from "@/lib/auth-client";
import { useDiagnostics } from "@/api/hooks/use-diagnostics";
import { useSidebar } from "@/components/layout/sidebar-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { BranchSwitcher } from "@/components/design-system/branch-switcher";

// Map paths to clean breadcrumb titles
const routeTitleMap: Record<string, string> = {
  // Platform Console routes
  "/platform/dashboard": "Platform Dashboard",
  "/platform/organizations": "Organizations Directory",
  "/platform/users": "User Directory",
  "/platform/marketplace": "Capability Marketplace",
  "/platform/pricing": "Pricing & Gateways",
  "/platform/facility-types": "Facility Types",
  "/platform/catalogs": "Master Catalogs",
  "/platform/audit": "Audit Trail & Telemetry",
  "/platform/diagnostics": "Diagnostics & Launch Gate",
  "/platform/demo-requests": "Inbound Demo Requests",
  "/platform/settings": "Console Settings",

  // Organization HQ routes
  "/organization/dashboard": "Executive HQ Dashboard",
  "/organization/branches": "Branch Facilities Network",
  "/organization/members": "Staff Roster & Access Control",
  "/organization/roles": "Roles & RBAC Permissions",
  "/organization/catalogs": "Catalogs & Custom Pricing",
  "/organization/billing": "Corporate Subscription & Billing",
  "/organization/branding": "Branding & Customization",
  "/organization/notifications": "Notification Settings",
  "/organization/integrations": "Developer APIs & Webhooks",
  "/organization/audit": "Corporate Audit Ledger",
  "/organization/settings": "Organization Profile & Settings",
};

export function PlatformTopbar() {
  const location = useLocation();
  const navigate = useNavigate();
  const { data: session } = authClient.useSession();
  const { isCollapsed, toggleCollapse } = useSidebar();

  const [searchQuery, setSearchQuery] = useState("");

  const isOrgContext = session?.bootstrap?.contexts?.current === "organization";
  const user = session?.user;
  const userInitials = user?.name
    ? user.name
        .split(" ")
        .map((n) => n[0])
        .join("")
        .slice(0, 2)
        .toUpperCase()
    : "CU";

  const displayRole =
    session?.bootstrap?.organization?.role === "owner" || user?.role === "owner"
      ? "Organization Owner"
      : session?.bootstrap?.organization?.role === "org_admin" || user?.role === "org_admin"
      ? "Org Administrator"
      : user?.platformRole || user?.role || (isOrgContext ? "Organization Member" : "Platform Staff");

  // Derive current page title
  const currentTitle =
    routeTitleMap[location.pathname] ||
    (location.pathname.startsWith("/platform/organizations/")
      ? "Organization Details"
      : isOrgContext
      ? "Organization Portal"
      : "Platform Console");

  const homePath = isOrgContext ? "/organization/dashboard" : "/platform/dashboard";
  const homeLabel = isOrgContext ? "Organization" : "Platform";

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchQuery.trim()) return;

    const q = searchQuery.toLowerCase().trim();
    if (isOrgContext) {
      if (q.includes("branch") || q.includes("facility")) {
        navigate("/organization/branches");
      } else if (q.includes("staff") || q.includes("member") || q.includes("user")) {
        navigate("/organization/members");
      } else if (q.includes("cat") || q.includes("price") || q.includes("test")) {
        navigate("/organization/catalogs");
      } else if (q.includes("bill") || q.includes("sub") || q.includes("plan")) {
        navigate("/organization/billing");
      } else if (q.includes("brand") || q.includes("logo") || q.includes("theme")) {
        navigate("/organization/branding");
      } else if (q.includes("api") || q.includes("key") || q.includes("hook")) {
        navigate("/organization/integrations");
      } else if (q.includes("audit") || q.includes("log")) {
        navigate("/organization/audit");
      } else {
        navigate("/organization/dashboard");
      }
    } else {
      if (q.includes("org") || q.includes("hospital") || q.includes("clinic") || q.includes("lab")) {
        navigate(`/platform/organizations?q=${encodeURIComponent(q)}`);
      } else if (q.includes("market") || q.includes("cap") || q.includes("addon")) {
        navigate("/platform/marketplace");
      } else if (q.includes("price") || q.includes("bill") || q.includes("pay") || q.includes("gate")) {
        navigate("/platform/pricing");
      } else if (q.includes("audit") || q.includes("log") || q.includes("event")) {
        navigate("/platform/audit");
      } else if (q.includes("diag") || q.includes("health") || q.includes("metric") || q.includes("gate")) {
        navigate("/platform/diagnostics");
      } else if (q.includes("cat") || q.includes("icd") || q.includes("test")) {
        navigate("/platform/catalogs");
      } else if (q.includes("demo") || q.includes("lead")) {
        navigate("/platform/demo-requests");
      } else if (q.includes("set") || q.includes("sec") || q.includes("policy")) {
        navigate("/platform/settings");
      } else {
        navigate(`/platform/organizations?q=${encodeURIComponent(q)}`);
      }
    }
  };

  return (
    <header className="sticky top-0 z-20 flex h-14 w-full items-center justify-between border-b border-border bg-card/95 px-6 backdrop-blur-sm">
      {/* Left: Sidebar Toggle & Page Breadcrumbs */}
      <div className="flex items-center gap-3">
        {isCollapsed && (
          <>
            <Button
              variant="ghost"
              size="icon"
              onClick={toggleCollapse}
              className="h-8 w-8 text-muted-foreground hover:text-foreground hover:bg-secondary/60"
              title="Expand Sidebar"
            >
              <PanelLeftOpen className="h-4 w-4" />
            </Button>
            <div className="h-4 w-[1px] bg-border hidden sm:block" />
          </>
        )}

        <BranchSwitcher />

        <div className="h-4 w-[1px] bg-border hidden sm:block" />

        <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
          <Link to={homePath} className="hover:text-foreground transition-colors font-medium">
            {homeLabel}
          </Link>
          <ChevronRight className="h-3.5 w-3.5 opacity-60" />
          <span className="font-semibold text-foreground truncate max-w-[200px] md:max-w-none">
            {currentTitle}
          </span>
        </div>
      </div>

      {/* Center: Global Search Bar */}
      <div className="hidden md:flex flex-1 max-w-md mx-6">
        <form onSubmit={handleSearchSubmit} className="relative w-full">
          <Search className="absolute left-3 top-2.5 h-3.5 w-3.5 text-muted-foreground" />
          <Input
            type="search"
            placeholder={
              isOrgContext
                ? "Search branches, staff, catalogs, audit ledger..."
                : "Search organizations, capabilities, catalogs, audit logs..."
            }
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-9 pr-12 text-xs h-8 bg-secondary/40 border-border/80 focus-visible:bg-secondary/80 transition-colors rounded-lg"
          />
          <div className="absolute right-2.5 top-2 hidden lg:flex items-center gap-0.5 text-[10px] font-mono text-muted-foreground bg-muted/60 px-1.5 py-0.5 rounded border border-border/60">
            <span>⌘</span>
            <span>K</span>
          </div>
        </form>
      </div>

      {/* Right: Telemetry Health, Notifications & Admin Menu */}
      <div className="flex items-center gap-3">
        {/* Discrete Live Cluster Dot */}
        <Link
          to={isOrgContext ? "/organization/audit" : "/platform/diagnostics"}
          className="hidden sm:flex items-center gap-2 rounded-md border border-border/60 bg-secondary/30 px-2.5 py-1 text-xs text-muted-foreground hover:text-foreground hover:border-border transition-colors"
          title="Cluster Health & Status"
        >
          <span className="relative flex h-2 w-2">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75" />
            <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500" />
          </span>
          <span className="font-mono text-[11px]">Live Cluster</span>
        </Link>

        {/* Audit Log Quick Link */}
        <Button
          asChild
          variant="ghost"
          size="icon"
          className="h-8 w-8 text-muted-foreground hover:text-foreground hover:bg-secondary/60 hidden sm:flex"
        >
          <Link to={isOrgContext ? "/organization/audit" : "/platform/audit"} title="Audit Trail">
            <Bell className="h-4 w-4" />
          </Link>
        </Button>

        {/* User Profile Menu */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              className="h-8 flex items-center gap-2 rounded-lg pl-1.5 pr-2.5 hover:bg-secondary/60 border border-transparent hover:border-border transition-all"
            >
              <Avatar className="h-6 w-6 border border-border">
                <AvatarFallback className="bg-primary text-primary-foreground text-[10px] font-bold">
                  {userInitials}
                </AvatarFallback>
              </Avatar>
              <span className="text-xs font-medium text-foreground truncate max-w-[120px] hidden sm:inline">
                {user?.name || "User"}
              </span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56 text-xs bg-card border-border shadow-lg">
            <DropdownMenuLabel className="font-normal py-2">
              <div className="flex flex-col space-y-1">
                <p className="font-semibold text-foreground text-xs">{user?.name}</p>
                <p className="text-[11px] text-muted-foreground truncate">{user?.email}</p>
                <Badge
                  variant="outline"
                  className="w-fit text-[9px] uppercase tracking-wider font-mono px-1.5 py-0 mt-1 border-primary/40 text-primary bg-primary/10"
                >
                  {displayRole}
                </Badge>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator className="bg-border" />
            <DropdownMenuItem asChild>
              <Link to={isOrgContext ? "/organization/settings" : "/platform/settings"} className="flex items-center gap-2 py-1.5 cursor-pointer">
                <Shield className="h-3.5 w-3.5 text-muted-foreground" />
                Settings & Profile
              </Link>
            </DropdownMenuItem>
            {isOrgContext ? (
              <DropdownMenuItem asChild>
                <Link to="/organization/billing" className="flex items-center gap-2 py-1.5 cursor-pointer">
                  <Layers className="h-3.5 w-3.5 text-muted-foreground" />
                  Corporate Subscription
                </Link>
              </DropdownMenuItem>
            ) : (
              <DropdownMenuItem asChild>
                <Link to="/platform/diagnostics" className="flex items-center gap-2 py-1.5 cursor-pointer">
                  <Activity className="h-3.5 w-3.5 text-muted-foreground" />
                  Launch Gate & Diagnostics
                </Link>
              </DropdownMenuItem>
            )}
            <DropdownMenuSeparator className="bg-border" />
            <DropdownMenuItem
              onClick={() => authClient.signOut()}
              className="text-destructive focus:text-destructive focus:bg-destructive/10 py-1.5 cursor-pointer flex items-center gap-2"
            >
              <LogOut className="h-3.5 w-3.5" />
              Sign Out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
