import React from "react";
import { Outlet } from "react-router-dom";
import { PlatformSidebar } from "./platform-sidebar";
import { PlatformTopbar } from "./platform-topbar";
import { SidebarProvider, useSidebar } from "./sidebar-context";
import { cn } from "@/lib/utils";

import { TopLoadingBar } from "@/components/loading";

function PlatformLayoutContent() {
  const { isCollapsed } = useSidebar();

  return (
    <div className="min-h-screen bg-background text-foreground flex">
      {/* Global Top Route Loading Bar */}
      <TopLoadingBar />

      {/* Fixed Collapsible Sidebar */}
      <PlatformSidebar />

      {/* Main Content Area dynamically adjusting margin */}
      <div
        className={cn(
          "flex-1 flex flex-col min-h-screen transition-all duration-300 ease-in-out",
          isCollapsed ? "pl-20" : "pl-64"
        )}
      >
        <PlatformTopbar />
        <main className="flex-1 p-8 max-w-7xl w-full mx-auto animate-fade-in">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

import { BrandThemeProvider } from "@/lib/theme/brand-theme-provider";

export function PlatformLayout() {
  return (
    <BrandThemeProvider>
      <SidebarProvider>
        <PlatformLayoutContent />
      </SidebarProvider>
    </BrandThemeProvider>
  );
}
