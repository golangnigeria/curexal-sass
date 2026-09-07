import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import {
  Dialog,
  DialogContent,
} from "@/components/ui/dialog";
import {
  Search,
  Stethoscope,
  Microscope,
  Activity,
  Pill,
  Building2,
  CreditCard,
  Users,
  LayoutDashboard,
  Sparkles,
  Command,
  ArrowRight,
  UserPlus,
  Barcode,
  Layers,
  LucideIcon,
} from "lucide-react";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { useNavigation } from "@/api/hooks/use-navigation";

const iconMap: Record<string, LucideIcon> = {
  LayoutDashboard,
  Users,
  Activity,
  Stethoscope,
  Microscope,
  Pill,
  Layers,
  Building2,
  CreditCard,
};

interface CommandPaletteProps {
  isOpen: boolean;
  onClose: () => void;
}

export const CommandPalette: React.FC<CommandPaletteProps> = ({
  isOpen,
  onClose,
}) => {
  const navigate = useNavigate();
  const { data: bootstrap } = useBootstrap();
  const [query, setQuery] = useState("");

  const activeBranchSlug =
    bootstrap?.branch?.slug ||
    bootstrap?.branch?.code?.toLowerCase() ||
    bootstrap?.workspace?.slug ||
    "main";

  const activeScope = bootstrap?.contexts?.current || "workspace";
  const { data: navData } = useNavigation({
    scope: activeScope as any,
    branchSlug: activeScope === "workspace" ? activeBranchSlug : undefined,
  });

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "k") {
        e.preventDefault();
        if (isOpen) {
          onClose();
        }
      }
      if (e.key === "Escape" && isOpen) {
        onClose();
      }
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [isOpen, onClose]);

  const handleSelectRoute = (path: string) => {
    navigate(path);
    onClose();
  };

  const activeNavItems = navData?.items || [];
  const filteredNavItems = activeNavItems.filter(
    (item) =>
      item.title.toLowerCase().includes(query.toLowerCase()) ||
      item.key.toLowerCase().includes(query.toLowerCase()) ||
      item.description?.toLowerCase().includes(query.toLowerCase())
  );

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-xl p-0 overflow-hidden border-border bg-card shadow-2xl rounded-2xl">
        {/* Search Header */}
        <div className="flex items-center px-4 border-b border-border bg-secondary/20">
          <Search className="w-4 h-4 text-muted-foreground mr-2 shrink-0" />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Type a command or search workspaces (e.g. LIS, Clinical, POS)..."
            className="w-full py-3.5 bg-transparent text-sm text-foreground placeholder-muted-foreground focus:outline-none font-sans"
            autoFocus
          />
          <kbd className="hidden sm:inline-flex items-center gap-1 px-1.5 py-0.5 text-[10px] font-mono text-muted-foreground bg-secondary rounded border border-border">
            ESC
          </kbd>
        </div>

        {/* Command Body */}
        <div className="max-h-[60vh] overflow-y-auto p-3 space-y-4 text-xs">
          {/* Quick Actions */}
          <div className="space-y-1">
            <p className="px-2 text-[10px] font-mono font-bold uppercase tracking-wider text-muted-foreground">
              Quick Actions
            </p>
            <div className="space-y-1">
              <button
                type="button"
                onClick={() => handleSelectRoute(`/${activeBranchSlug}/reception`)}
                className="w-full flex items-center justify-between p-2 rounded-lg hover:bg-secondary text-foreground text-left transition-colors group cursor-pointer"
              >
                <div className="flex items-center gap-2.5">
                  <div className="p-1.5 rounded-md bg-teal-500/10 text-teal-600 dark:text-teal-400">
                    <UserPlus className="w-4 h-4" />
                  </div>
                  <span className="font-medium">New Patient Intake (MPI)</span>
                </div>
                <ArrowRight className="w-3.5 h-3.5 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity" />
              </button>

              <button
                type="button"
                onClick={() => handleSelectRoute(`/${activeBranchSlug}/care-desk`)}
                className="w-full flex items-center justify-between p-2 rounded-lg hover:bg-secondary text-foreground text-left transition-colors group cursor-pointer"
              >
                <div className="flex items-center gap-2.5">
                  <div className="p-1.5 rounded-md bg-teal-500/10 text-teal-600 dark:text-teal-400">
                    <Activity className="w-4 h-4" />
                  </div>
                  <span className="font-medium">Nursing Triage & Vitals (Care Desk)</span>
                </div>
                <ArrowRight className="w-3.5 h-3.5 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity" />
              </button>

              <button
                type="button"
                onClick={() => handleSelectRoute(`/${activeBranchSlug}/clinical`)}
                className="w-full flex items-center justify-between p-2 rounded-lg hover:bg-secondary text-foreground text-left transition-colors group cursor-pointer"
              >
                <div className="flex items-center gap-2.5">
                  <div className="p-1.5 rounded-md bg-indigo-500/10 text-indigo-600 dark:text-indigo-400">
                    <Stethoscope className="w-4 h-4" />
                  </div>
                  <span className="font-medium">Doctor Consultation Queue (EMR)</span>
                </div>
                <ArrowRight className="w-3.5 h-3.5 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity" />
              </button>

              <button
                type="button"
                onClick={() => handleSelectRoute(`/${activeBranchSlug}/billing`)}
                className="w-full flex items-center justify-between p-2 rounded-lg hover:bg-secondary text-foreground text-left transition-colors group cursor-pointer"
              >
                <div className="flex items-center gap-2.5">
                  <div className="p-1.5 rounded-md bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
                    <CreditCard className="w-4 h-4" />
                  </div>
                  <span className="font-medium">New Cashier Invoice (POS)</span>
                </div>
                <ArrowRight className="w-3.5 h-3.5 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity" />
              </button>
            </div>
          </div>

          {/* Database-driven Workspaces & Navigation Items */}
          <div className="space-y-1">
            <p className="px-2 text-[10px] font-mono font-bold uppercase tracking-wider text-muted-foreground">
              Workspaces & Modules
            </p>
            <div className="space-y-1">
              {filteredNavItems.map((item) => {
                const IconComponent = iconMap[item.icon] || LayoutDashboard;
                return (
                  <button
                    key={item.id}
                    type="button"
                    onClick={() => handleSelectRoute(item.path)}
                    className="w-full flex items-center justify-between p-2 rounded-lg hover:bg-secondary text-foreground text-left transition-colors group cursor-pointer"
                  >
                    <div className="flex items-center gap-2.5">
                      <div className="p-1.5 rounded-md bg-secondary text-muted-foreground group-hover:text-primary transition-colors">
                        <IconComponent className="w-4 h-4" />
                      </div>
                      <div>
                        <p className="font-medium text-foreground">{item.title}</p>
                        {item.description && (
                          <p className="text-[11px] text-muted-foreground">{item.description}</p>
                        )}
                      </div>
                    </div>
                    <span className="text-[10px] font-mono text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity">
                      Jump →
                    </span>
                  </button>
                );
              })}
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="px-4 py-2 bg-secondary/30 border-t border-border flex items-center justify-between text-[11px] text-muted-foreground">
          <span>Active Facility: <strong className="text-foreground">/{activeBranchSlug}</strong></span>
          <span className="flex items-center gap-1 font-mono text-[10px]">
            <span>Navigate with</span>
            <kbd className="px-1 py-0.5 bg-card border rounded text-[9px]">↑</kbd>
            <kbd className="px-1 py-0.5 bg-card border rounded text-[9px]">↓</kbd>
            <kbd className="px-1 py-0.5 bg-card border rounded text-[9px]">↵</kbd>
          </span>
        </div>
      </DialogContent>
    </Dialog>
  );
};
