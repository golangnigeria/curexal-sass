import React from "react";
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
import { Building2, ChevronDown, Check, Plus } from "lucide-react";
import { Link } from "react-router-dom";

export function BranchSwitcher() {
  const { data: bootstrap } = useBootstrap();
  const { data: branches } = useOrgBranches();

  const activeBranchId = bootstrap?.workspace?.id;
  const activeBranchName = bootstrap?.workspace?.name || "Main Facility";
  const isOrgContext = bootstrap?.contexts?.current === "organization";

  if (!isOrgContext && (!branches || branches.length <= 1)) {
    return (
      <div className="flex items-center gap-1.5 px-2 py-1 rounded-md bg-secondary/30 text-xs font-medium text-muted-foreground border border-border/50">
        <Building2 className="w-3.5 h-3.5 text-primary" />
        <span className="truncate max-w-[140px]">{activeBranchName}</span>
      </div>
    );
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="outline"
          size="sm"
          className="h-8 px-2.5 gap-1.5 text-xs font-medium bg-card border-border hover:bg-secondary/60 transition-colors"
        >
          <Building2 className="w-3.5 h-3.5 text-primary" />
          <span className="truncate max-w-[130px] font-semibold text-foreground">
            {activeBranchName}
          </span>
          <ChevronDown className="w-3 h-3 text-muted-foreground opacity-60 ml-0.5" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" className="w-64 text-xs bg-card border-border shadow-lg">
        <DropdownMenuLabel className="text-[11px] uppercase tracking-wider text-muted-foreground font-semibold py-1.5">
          Branch Facilities
        </DropdownMenuLabel>
        <DropdownMenuSeparator className="bg-border" />
        {branches && branches.length > 0 ? (
          branches.map((b) => {
            const isSelected = b.id === activeBranchId;
            return (
              <DropdownMenuItem
                key={b.id}
                className="flex items-center justify-between py-2 cursor-pointer focus:bg-secondary/60"
              >
                <div className="space-y-0.5">
                  <p className={`font-semibold ${isSelected ? "text-primary" : "text-foreground"}`}>
                    {b.name}
                  </p>
                  <p className="text-[10px] text-muted-foreground font-mono">
                    {b.code} • {b.facilityType}
                  </p>
                </div>
                {isSelected && <Check className="w-4 h-4 text-primary shrink-0" />}
              </DropdownMenuItem>
            );
          })
        ) : (
          <div className="py-3 px-2 text-center text-[11px] text-muted-foreground">
            No additional branches configured
          </div>
        )}
        <DropdownMenuSeparator className="bg-border" />
        <DropdownMenuItem asChild>
          <Link
            to="/organization/branches"
            className="flex items-center gap-2 py-1.5 text-primary font-medium cursor-pointer"
          >
            <Plus className="w-3.5 h-3.5" />
            Manage Branch Facilities
          </Link>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
