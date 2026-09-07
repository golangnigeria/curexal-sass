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
  Building2,
  Stethoscope,
  Microscope,
  Pill,
  ArrowLeft,
  Sparkles,
} from "lucide-react";
import { authClient } from "@/lib/auth-client";
import { useDiagnostics } from "@/api/hooks/use-diagnostics";
import { useSidebar } from "@/components/layout/sidebar-context";
import { useNavigation } from "@/api/hooks/use-navigation";
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
import { CommandPalette } from "@/features/search/command-palette";

function getFacilityTypeIcon(typeStr?: string) {
  const norm = (typeStr || "").toLowerCase();
  if (norm.includes("lab")) return Microscope;
  if (norm.includes("clinic")) return Stethoscope;
  if (norm.includes("pharmacy")) return Pill;
  if (norm.includes("hospital")) return Building2;
  return Activity;
}

export function PlatformTopbar() {
  const location = useLocation();
  const navigate = useNavigate();
  const { data: session } = authClient.useSession();
  const { isCollapsed, toggleCollapse } = useSidebar();

  const [searchQuery, setSearchQuery] = useState("");
  const [isCmdOpen, setIsCmdOpen] = useState(false);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        setIsCmdOpen((prev) => !prev);
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  const isPlatformContext = location.pathname.startsWith("/platform");
  const isOrgContext = !isPlatformContext && location.pathname.startsWith("/organization");
  const isWorkspaceContext = !isPlatformContext && !location.pathname.startsWith("/organization") && location.pathname !== "/login";

  const user = session?.user;
  const org = session?.bootstrap?.organization;
  const pathParts = location.pathname.split("/").filter(Boolean);
  const activeBranchSlug = pathParts[0] || session?.bootstrap?.branch?.slug || session?.bootstrap?.workspace?.slug || "main";

  const activeScope = isPlatformContext ? "platform" : isOrgContext ? "organization" : "workspace";
  const { data: navData } = useNavigation({
    scope: activeScope,
    branchSlug: isWorkspaceContext ? activeBranchSlug : undefined,
  });

  const activeBranch =
    session?.bootstrap?.availableBranches?.find(
      (b: any) => b.slug === activeBranchSlug || b.code?.toLowerCase() === activeBranchSlug.toLowerCase()
    ) ||
    session?.bootstrap?.branch ||
    session?.bootstrap?.workspace;
  const facilityType = activeBranch?.facilityType || "Facility";
  const FacilityIcon = getFacilityTypeIcon(facilityType);

  const isOrgAdminOrOwner =
    session?.bootstrap?.platform?.isStaff === true ||
    user?.isPlatformAdmin === true ||
    session?.bootstrap?.organization?.role === "owner" ||
    session?.bootstrap?.organization?.role === "org_admin" ||
    session?.bootstrap?.organization?.role === "org_regional_manager" ||
    user?.role === "owner" ||
    user?.role === "org_admin" ||
    user?.role === "org_regional_manager";

  const userInitials = user?.name
    ? user.name
        .split(" ")
        .map((n) => n[0])
        .join("")
        .slice(0, 2)
        .toUpperCase()
    : "CU";

function formatRoleBadgeTitle(rawRole?: string | null): string {
  if (!rawRole) return "Healthcare Specialist";
  const r = rawRole.toLowerCase().trim();
  switch (r) {
    case "branch_admin":
      return "Branch Administrator";
    case "org_admin":
      return "Organization Administrator";
    case "owner":
      return "Organization Owner";
    case "doctor":
    case "medical_doctor":
      return "Medical Doctor";
    case "nurse":
    case "registered_nurse":
      return "Registered Nurse";
    case "lab_technician":
    case "laboratory_scientist":
    case "lab_scientist":
      return "Laboratory Scientist";
    case "pharmacist":
      return "Pharmacist";
    case "radiologist":
      return "Radiologist";
    case "receptionist":
    case "front_desk":
      return "Front Desk Specialist";
    case "cashier":
    case "billing":
      return "Billing & Cashier Specialist";
    case "org_regional_manager":
      return "Regional Manager";
    case "org_quality_manager":
      return "Quality Assurance Manager";
    case "org_finance_manager":
      return "Finance Manager";
    case "org_hr_manager":
      return "HR Manager";
    case "super_admin":
      return "Super Admin";
    case "platform_admin":
      return "Platform Admin";
    case "platform_staff":
      return "Platform Staff";
    default:
      if (r === "user" || r === "member") return "Staff Member";
      return r
        .split(/[_\s]+/)
        .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
        .join(" ");
  }
}

  const rawRoleCandidate =
    (isPlatformContext
      ? session?.bootstrap?.platform?.role || user?.platformRole
      : isOrgContext
      ? session?.bootstrap?.organization?.role || user?.role
      : (activeBranch as any)?.role || session?.bootstrap?.organization?.role || user?.role) ||
    session?.bootstrap?.organization?.role ||
    user?.role;

  const displayRole = formatRoleBadgeTitle(rawRoleCandidate);

  // Derive current page title & breadcrumbs directly from database navigation models
  const matchingNavItem = navData?.items?.find(
    (item) => item.path === location.pathname || (location.pathname.startsWith(item.path) && item.path !== "/" && !item.path.endsWith("/dashboard"))
  );

  let currentTitle = matchingNavItem?.title;
  if (!currentTitle) {
    if (isWorkspaceContext && pathParts.length >= 2) {
      currentTitle = pathParts[1].charAt(0).toUpperCase() + pathParts[1].slice(1) + " Workspace";
    } else if (location.pathname.startsWith("/platform/organizations/")) {
      currentTitle = "Organization Details";
    } else if (isOrgContext) {
      currentTitle = "Organization Portal";
    } else {
      currentTitle = "Overview";
    }
  }

  const homePath = isPlatformContext
    ? "/platform/dashboard"
    : isOrgContext
    ? "/organization/dashboard"
    : `/${pathParts[0] || "main"}/dashboard`;
  const homeLabel = isPlatformContext
    ? "Platform"
    : isOrgContext
    ? org?.name || "Organization"
    : activeBranch?.name || "Facility";

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchQuery.trim()) return;

    const q = searchQuery.toLowerCase().trim();
    // Search within live backend navigation items
    const matched = navData?.items?.find(
      (item) =>
        item.title.toLowerCase().includes(q) ||
        item.key.toLowerCase().includes(q) ||
        item.description?.toLowerCase().includes(q)
    );

    if (matched) {
      navigate(matched.path);
      setSearchQuery("");
      return;
    }

    if (isWorkspaceContext) {
      const bSlug = pathParts[0] || "main";
      navigate(`/${bSlug}/dashboard`);
    }
    setSearchQuery("");
  };

  return (
    <header className="sticky top-0 z-20 flex h-16 w-full items-center justify-between border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/80 px-4 sm:px-6">
      {/* Left side: Breadcrumb & Context Navigation */}
      <div className="flex items-center gap-3">
        {isCollapsed && (
          <Button
            variant="ghost"
            size="icon"
            onClick={toggleCollapse}
            className="h-8 w-8 text-muted-foreground hover:text-foreground"
            title="Expand Sidebar"
          >
            <PanelLeftOpen className="h-4 w-4" />
          </Button>
        )}

        <nav aria-label="Breadcrumb" className="flex items-center space-x-1.5 text-xs text-muted-foreground">
          <Link
            to={homePath}
            className="flex items-center gap-1.5 font-medium hover:text-foreground transition-colors"
          >
            {isWorkspaceContext && <FacilityIcon className="w-3.5 h-3.5 text-primary" />}
            <span>{homeLabel}</span>
          </Link>
          <ChevronRight className="h-3.5 w-3.5 text-muted-foreground/60" />
          <span className="font-semibold text-foreground tracking-tight">{currentTitle}</span>
        </nav>
      </div>

      {/* Right side: Global Search, Branch Switcher & Controls */}
      <div className="flex items-center gap-3">
        {/* Branch / Facility Switcher strictly for Organization Context */}
        {!isPlatformContext && <BranchSwitcher />}

        {/* Global Quick Search */}
        <form onSubmit={handleSearchSubmit} className="relative hidden md:block">
          <Search className="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-muted-foreground" />
          <Input
            type="search"
            placeholder="Search console (⌘K)..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onClick={() => setIsCmdOpen(true)}
            className="w-48 lg:w-64 pl-8 pr-12 h-8 text-xs bg-secondary/50 border-border focus:bg-background transition-all"
          />
          <kbd className="absolute right-2 top-2 pointer-events-none hidden h-4 select-none items-center gap-1 rounded border border-border bg-muted px-1.5 font-mono text-[9px] font-medium text-muted-foreground sm:flex">
            ⌘K
          </kbd>
        </form>

        {/* Command Palette Trigger Button (Mobile) */}
        <Button
          variant="ghost"
          size="icon"
          onClick={() => setIsCmdOpen(true)}
          className="h-8 w-8 md:hidden text-muted-foreground hover:text-foreground"
          title="Command Palette"
        >
          <Command className="h-4 w-4" />
        </Button>

        {/* System Diagnostics Indicator (Platform Only) */}
        {isPlatformContext && (
          <Button
            asChild
            variant="ghost"
            size="sm"
            className="h-8 text-xs gap-1.5 text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 hover:bg-emerald-500/20 border border-emerald-500/30"
          >
            <Link to="/platform/diagnostics">
              <Sparkles className="w-3.5 h-3.5 animate-spin text-emerald-500" style={{ animationDuration: "3s" }} />
              <span className="font-semibold hidden sm:inline">Launch Gate Ready</span>
            </Link>
          </Button>
        )}

        {/* User Profile & Workspace Dropdown */}
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

      {/* Global Command Center Modal */}
      <CommandPalette isOpen={isCmdOpen} onClose={() => setIsCmdOpen(false)} />
    </header>
  );
}
