import React, { useState } from "react";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useOrgBranches } from "@/api/hooks/use-organization";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Building2,
  ChevronDown,
  Check,
  Plus,
  Star,
  Activity,
  Microscope,
  Stethoscope,
  Pill,
  Search,
  ArrowLeft,
} from "lucide-react";
import { Link, useLocation, useNavigate, useParams } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import { authClient } from "@/lib/auth-client";

function getBranchIcon(typeStr?: string) {
  const norm = (typeStr || "").toLowerCase();
  if (norm.includes("lab")) return Microscope;
  if (norm.includes("clinic")) return Stethoscope;
  if (norm.includes("pharmacy")) return Pill;
  if (norm.includes("hospital")) return Building2;
  return Activity;
}

function resolveDefaultWorkspaceModule(facilityType?: string): string {
  const norm = (facilityType || "").toLowerCase();
  if (norm.includes("clinic") || norm.includes("outpatient") || norm.includes("emr")) {
    return "clinical";
  }
  if (norm.includes("pharmacy")) {
    return "pharmacy";
  }
  if (norm.includes("radiology") || norm.includes("imaging")) {
    return "radiology";
  }
  return "clinical";
}

export function BranchSwitcher() {
  const queryClient = useQueryClient();
  const { data: bootstrap } = useBootstrap();
  const { data: branches } = useOrgBranches();
  const location = useLocation();
  const navigate = useNavigate();
  const { branchSlug } = useParams<{ branchSlug?: string }>();
  const [filterQuery, setFilterQuery] = useState("");

  const activeBranchSlug =
    branchSlug ||
    bootstrap?.branch?.slug ||
    bootstrap?.branch?.code?.toLowerCase() ||
    bootstrap?.workspace?.slug ||
    "main";
  const activeBranchName = bootstrap?.branch?.name || bootstrap?.workspace?.name || "Main Facility";
  const activeBranchType = bootstrap?.branch?.facilityType || bootstrap?.workspace?.facilityType || "Facility";
  const ActiveIcon = getBranchIcon(activeBranchType);

  const isOrgContext =
    location.pathname.startsWith("/organization") ||
    bootstrap?.contexts?.current === "organization";
  const isOrgAdminOrOwner =
    bootstrap?.platform?.isStaff === true ||
    bootstrap?.organization?.role === "owner" ||
    bootstrap?.organization?.role === "org_admin" ||
    bootstrap?.organization?.role === "org_regional_manager" ||
    bootstrap?.contexts?.available?.includes("organization");

  // Determine current workspace subpath (e.g., 'clinical', 'care-desk', 'reception', 'billing', 'dashboard')
  const pathParts = location.pathname.split("/").filter(Boolean);
  const currentWorkspacePath =
    pathParts.length > 1 && pathParts[0] === activeBranchSlug
      ? pathParts.slice(1).join("/")
      : pathParts.length > 0 && pathParts[0] !== "organization" && pathParts[0] !== "platform"
      ? pathParts[0]
      : "dashboard";

  const handleSelectBranch = async (b: any) => {
    const targetSlug = b.slug || b.code?.toLowerCase();
    const targetType = b.facilityTypeName || b.facilityTypeCode || b.facilityType || "Facility";
    if (targetSlug === activeBranchSlug && !isOrgContext) return;

    try {
      if (b.id) {
        await authClient.switchBranch(b.id);
      }
    } catch (err) {
      console.error("Failed to switch branch session:", err);
    }

    // Purge cached workspace data to prevent stale cross-facility data leakage
    queryClient.removeQueries({ queryKey: ["workspace"] });
    // Invalidate bootstrap context
    queryClient.invalidateQueries({ queryKey: ["bootstrap"] });

    // Shared subpaths that can carry over across branch types
    const sharedSubpaths = ["dashboard", "billing", "settings"];
    let nextPath = currentWorkspacePath;

    if (!sharedSubpaths.includes(currentWorkspacePath)) {
      nextPath = resolveDefaultWorkspaceModule(targetType);
    }

    navigate(`/${targetSlug}/${nextPath}`);
  };

  const branchList = (branches && branches.length > 0 ? branches : bootstrap?.availableBranches) || [];

  // For branch staff with only 1 accessible branch, show a static facility badge
  if (!isOrgAdminOrOwner && branchList.length <= 1) {
    return (
      <div className="flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-secondary/40 text-xs font-medium text-foreground border border-border/60">
        <ActiveIcon className="w-3.5 h-3.5 text-primary" />
        <span className="truncate max-w-[140px] font-semibold">{activeBranchName}</span>
      </div>
    );
  }

  const filteredBranches = branchList.filter((b: any) => {
    if (!filterQuery) return true;
    const q = filterQuery.toLowerCase();
    return (
      (b.name || "").toLowerCase().includes(q) ||
      (b.code || "").toLowerCase().includes(q) ||
      (b.slug || "").toLowerCase().includes(q)
    );
  });

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="outline"
          size="sm"
          className="h-8 px-2.5 gap-1.5 text-xs font-medium bg-card border-border hover:bg-secondary/60 transition-colors"
        >
          <ActiveIcon className="w-3.5 h-3.5 text-primary" />
          <span className="truncate max-w-[130px] font-semibold text-foreground">
            {activeBranchName}
          </span>
          <ChevronDown className="w-3 h-3 text-muted-foreground opacity-60 ml-0.5" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" className="w-72 text-xs bg-card border-border shadow-lg">
        <DropdownMenuLabel className="text-[11px] uppercase tracking-wider text-muted-foreground font-semibold py-1.5 flex items-center justify-between">
          <span>Branch Facilities Network</span>
          <span className="text-[10px] font-mono lowercase text-primary font-normal">{branchList.length} active</span>
        </DropdownMenuLabel>

        {branchList.length > 4 && (
          <div className="px-2 py-1.5">
            <div className="relative">
              <Search className="w-3 h-3 absolute left-2 top-2 text-muted-foreground" />
              <Input
                placeholder="Filter branches..."
                value={filterQuery}
                onChange={(e) => setFilterQuery(e.target.value)}
                className="h-7 text-xs pl-7 bg-secondary/30"
              />
            </div>
          </div>
        )}

        <DropdownMenuSeparator className="bg-border" />

        <div className="max-h-60 overflow-y-auto">
          {filteredBranches.length > 0 ? (
            filteredBranches.map((b: any) => {
              const bSlug = b.slug || b.code?.toLowerCase();
              const isSelected = !isOrgContext && (bSlug === activeBranchSlug || b.id === bootstrap?.branch?.id);
              const rawType = b.facilityTypeName || b.facilityTypeCode || b.facilityType || "Facility";
              const Icon = getBranchIcon(rawType);

              return (
                <DropdownMenuItem
                  key={b.id}
                  onClick={() => handleSelectBranch(b)}
                  className="flex items-center justify-between py-2 cursor-pointer focus:bg-secondary/60 gap-2"
                >
                  <div className="flex items-center gap-2 min-w-0">
                    <div className={`p-1.5 rounded-md shrink-0 ${isSelected ? "bg-primary text-white" : "bg-secondary text-muted-foreground"}`}>
                      <Icon className="w-3.5 h-3.5" />
                    </div>
                    <div className="space-y-0.5 min-w-0">
                      <div className="flex items-center gap-1">
                        <p className={`font-semibold truncate ${isSelected ? "text-primary" : "text-foreground"}`}>
                          {b.name}
                        </p>
                        {b.isHeadquarters && (
                          <Star className="w-2.5 h-2.5 fill-amber-500 text-amber-500 shrink-0" />
                        )}
                      </div>
                      <p className="text-[10px] text-muted-foreground font-mono truncate">
                        /{bSlug} • {rawType.replace(/_/g, " ")}
                      </p>
                    </div>
                  </div>
                  {isSelected && <Check className="w-4 h-4 text-primary shrink-0" />}
                </DropdownMenuItem>
              );
            })
          ) : (
            <div className="py-3 px-2 text-center text-[11px] text-muted-foreground">
              No matching branch facilities
            </div>
          )}
        </div>

        {isOrgAdminOrOwner && (
          <>
            <DropdownMenuSeparator className="bg-border" />
            <DropdownMenuItem asChild>
              <Link
                to="/organization/dashboard"
                className="flex items-center gap-2 py-1.5 text-foreground font-medium cursor-pointer"
              >
                <Building2 className="w-3.5 h-3.5 text-primary" />
                <span>Executive Organization HQ</span>
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link
                to="/organization/branches"
                className="flex items-center gap-2 py-1.5 text-primary font-medium cursor-pointer"
              >
                <Plus className="w-3.5 h-3.5" />
                <span>Manage Facilities Network</span>
              </Link>
            </DropdownMenuItem>
          </>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
