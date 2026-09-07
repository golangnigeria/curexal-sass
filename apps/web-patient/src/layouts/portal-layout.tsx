import React from "react";
import { Link, Outlet, useLocation, useNavigate } from "react-router-dom";
import {
  LayoutDashboard,
  CalendarPlus,
  LogOut,
  Activity,
  LucideIcon,
} from "lucide-react";
import { CurexalLogoSymbol, useBrandTheme } from "@curexal/design-system";
import { useQuery } from "@tanstack/react-query";
import { navigationService } from "@curexal/api-client";
import type { NavigationResponse } from "@curexal/contracts";

const iconMap: Record<string, LucideIcon> = {
  LayoutDashboard,
  CalendarPlus,
  Activity,
};

export const PatientPortalLayout: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const { orgName, logoUrl } = useBrandTheme();

  const organizationDisplay = orgName || "Curexal Health Network";

  const handleLogout = () => {
    localStorage.removeItem("curexal_portal_token");
    localStorage.removeItem("curexal_portal_patient");
    navigate("/login");
  };

  const { data: navData, isLoading: isNavLoading } = useQuery<NavigationResponse>({
    queryKey: ["patient-navigation"],
    queryFn: () =>
      navigationService.getNavigation({
        scope: "patient",
        host: typeof window !== "undefined" ? window.location.hostname : "patient.localhost",
      }),
    staleTime: 1000 * 60 * 5,
  });

  const navItems = navData?.items || [];

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-[#0B1120] text-slate-900 dark:text-white flex flex-col font-sans selection:bg-teal-500/20 relative">
      {/* Background Radial Glow */}
      <div className="fixed inset-0 bg-[radial-gradient(ellipse_80%_60%_at_50%_-20%,rgba(15,118,110,0.08),rgba(255,255,255,0))] dark:bg-[radial-gradient(ellipse_80%_60%_at_50%_-20%,rgba(20,184,166,0.12),rgba(11,17,32,0))] pointer-events-none" />

      {/* Top Navbar */}
      <header className="sticky top-0 z-40 w-full bg-white/90 dark:bg-[#0B1120]/90 backdrop-blur-xl border-b border-slate-200/80 dark:border-slate-800 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-6">
            <Link to="/dashboard" className="flex items-center gap-3 group">
              <div className="p-1 rounded-xl bg-teal-50 dark:bg-teal-950/60 border border-teal-200/60 dark:border-teal-800/60 shadow-sm group-hover:scale-105 transition-transform flex items-center justify-center">
                {logoUrl ? (
                  <img src={logoUrl} alt={organizationDisplay} className="w-7 h-7 object-contain rounded" />
                ) : (
                  <CurexalLogoSymbol className="w-7 h-7" />
                )}
              </div>
              <div className="flex flex-col">
                <span className="font-extrabold tracking-tight text-slate-900 dark:text-white text-base">
                  {organizationDisplay}
                </span>
                <span className="text-[10px] font-bold text-teal-600 dark:text-teal-400 uppercase tracking-wider">
                  Patient Health Workspace
                </span>
              </div>
            </Link>

            {/* Navigation Tabs (DB-Driven) */}
            <nav className="hidden md:flex items-center gap-1.5 ml-4">
              {navItems.map((item) => {
                const Icon = iconMap[item.icon] || Activity;
                const isActive = location.pathname === item.path;
                return (
                  <Link
                    key={item.id || item.path}
                    to={item.path}
                    className={`flex items-center gap-2 px-3.5 py-1.5 rounded-xl text-xs font-semibold transition-all ${
                      isActive
                        ? "bg-teal-500/10 text-teal-700 dark:text-teal-300 border border-teal-500/30 shadow-sm"
                        : "text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-slate-800/60"
                    }`}
                  >
                    <Icon className="w-3.5 h-3.5" />
                    <span>{item.title}</span>
                  </Link>
                );
              })}
            </nav>
          </div>

          {/* Right Header Actions */}
          <div className="flex items-center gap-3">
            <div className="hidden sm:flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 text-xs font-semibold">
              <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
              <span>Care Network Online</span>
            </div>

            <button
              onClick={handleLogout}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-semibold text-slate-500 hover:text-rose-600 hover:bg-rose-500/10 border border-slate-200 dark:border-slate-800 transition-colors cursor-pointer"
              title="Sign Out"
            >
              <LogOut className="w-3.5 h-3.5" />
              <span className="hidden sm:inline">Sign Out</span>
            </button>
          </div>
        </div>
      </header>

      {/* Main Content Body */}
      <main className="flex-1 max-w-7xl w-full mx-auto p-4 sm:p-6 lg:p-8 z-10">
        <Outlet />
      </main>

      {/* Public Footer */}
      <footer className="w-full border-t border-slate-200/80 dark:border-slate-800 py-6 text-center text-xs text-slate-500 dark:text-slate-400 bg-white/50 dark:bg-[#0B1120]/50 backdrop-blur-md z-10">
        <div className="max-w-7xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-3">
          <div className="flex items-center gap-2">
            <CurexalLogoSymbol className="w-4 h-4" />
            <span className="font-bold text-slate-900 dark:text-white">&copy; {new Date().getFullYear()} {organizationDisplay}</span>
          </div>
          <span className="text-[11px]">Enterprise Encrypted • HIPAA & NDPR Verified</span>
        </div>
      </footer>
    </div>
  );
};
