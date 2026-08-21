import React from "react";
import { createBrowserRouter, Navigate } from "react-router-dom";
import { ProtectedRoute } from "@/components/auth/protected-route";
import { PlatformGuard, RouteLoadingScreen } from "@/components/auth/platform-guard";
import { OrganizationGuard } from "@/components/auth/organization-guard";
import { WorkspaceGuard } from "@/components/auth/workspace-guard";
import { PlatformLayout } from "@/components/layout/platform-layout";
import { RouteErrorBoundary } from "@/components/ui/error-boundary";
import { authClient } from "@/lib/auth-client";
import { useBootstrap } from "@/api/hooks/use-bootstrap";

import LoginPage from "@/pages/auth/login";
import DashboardPage from "@/pages/platform/dashboard";
import OrganizationsPage from "@/pages/platform/organizations";
import OrganizationDetailPage from "@/pages/platform/organizations/org-detail";
import UsersDirectoryPage from "@/pages/platform/users";
import MarketplacePage from "@/pages/platform/marketplace";
import PricingPage from "@/pages/platform/pricing";
import FacilityTypesPage from "@/pages/platform/facility-types";
import MasterCatalogsPage from "@/pages/platform/catalogs";
import AuditLogsPage from "@/pages/platform/audit";
import DiagnosticsPage from "@/pages/platform/diagnostics";
import DemoRequestsPage from "@/pages/platform/demo-requests";
import SettingsPage from "@/pages/platform/settings";
import OrganizationDashboardPage from "@/pages/organization/dashboard";
import OrganizationBranchesPage from "@/pages/organization/branches";
import OrganizationMembersPage from "@/pages/organization/members";
import OrganizationRolesPage from "@/pages/organization/roles";
import OrganizationCatalogsPage from "@/pages/organization/catalogs";
import OrganizationBillingPage from "@/pages/organization/billing";
import OrganizationBrandingPage from "@/pages/organization/branding";
import OrganizationNotificationsPage from "@/pages/organization/notifications";
import OrganizationIntegrationsPage from "@/pages/organization/integrations";
import OrganizationAuditPage from "@/pages/organization/audit";
import OrganizationSettingsPage from "@/pages/organization/settings";

import WorkspaceDashboardPage from "@/pages/workspace/dashboard";
import WorkspaceLaboratoryPage from "@/pages/workspace/laboratory";
import WorkspaceClinicalPage from "@/pages/workspace/clinical";
import WorkspaceHospitalPage from "@/pages/workspace/hospital";
import WorkspaceRadiologyPage from "@/pages/workspace/radiology";
import WorkspacePharmacyPage from "@/pages/workspace/pharmacy";
import WorkspaceBillingPage from "@/pages/workspace/billing";

function RootRedirect() {
  const { data: session, isPending: isSessionPending } = authClient.useSession();
  const { data: bootstrap, isPending: isBootstrapPending } = useBootstrap();

  // Wait for session and bootstrap resolution to prevent redirect race
  if (isSessionPending || (session?.user && isBootstrapPending && !bootstrap && !session.bootstrap)) {
    return <RouteLoadingScreen message="Resolving destination..." />;
  }

  if (!session || !session.user) {
    return <Navigate to="/login" replace />;
  }

  const effectiveBootstrap = bootstrap || session.bootstrap;
  const user = session.user;

  // 1. Platform Staff -> /platform/dashboard
  if (
    effectiveBootstrap?.platform?.isStaff === true ||
    user.isPlatformAdmin === true ||
    user.platformRole === "super_admin" ||
    user.role === "super_admin"
  ) {
    return <Navigate to="/platform/dashboard" replace />;
  }

  // 2. Organization Authorized -> /organization/dashboard
  if (
    Boolean(effectiveBootstrap?.organization?.id) ||
    effectiveBootstrap?.contexts?.current === "organization" ||
    user.role === "owner" ||
    user.role === "org_admin" ||
    user.role === "org_regional_manager" ||
    Boolean(user.organizationId)
  ) {
    return <Navigate to="/organization/dashboard" replace />;
  }

  // 3. Workspace Authorized -> /workspace/dashboard
  if (
    Boolean(effectiveBootstrap?.workspace?.id) ||
    effectiveBootstrap?.contexts?.current === "workspace" ||
    Boolean(user.activeTenantId) ||
    Boolean(user.workspaceId)
  ) {
    return <Navigate to="/workspace/dashboard" replace />;
  }

  // 4. Fallback -> /login
  return <Navigate to="/login" replace />;
}

export const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
    errorElement: <RouteErrorBoundary />,
  },
  {
    path: "/",
    element: <ProtectedRoute />,
    errorElement: <RouteErrorBoundary />,
    children: [
      {
        element: <PlatformLayout />,
        errorElement: <RouteErrorBoundary />,
        children: [
          {
            index: true,
            element: <RootRedirect />,
          },
          // Platform Console Routes (Protected by PlatformGuard)
          {
            element: <PlatformGuard />,
            children: [
              {
                path: "platform/dashboard",
                element: <DashboardPage />,
              },
              {
                path: "platform/organizations",
                element: <OrganizationsPage />,
              },
              {
                path: "platform/organizations/:id",
                element: <OrganizationDetailPage />,
              },
              {
                path: "platform/users",
                element: <UsersDirectoryPage />,
              },
              {
                path: "platform/marketplace",
                element: <MarketplacePage />,
              },
              {
                path: "platform/pricing",
                element: <PricingPage />,
              },
              {
                path: "platform/facility-types",
                element: <FacilityTypesPage />,
              },
              {
                path: "platform/catalogs",
                element: <MasterCatalogsPage />,
              },
              {
                path: "platform/audit",
                element: <AuditLogsPage />,
              },
              {
                path: "platform/diagnostics",
                element: <DiagnosticsPage />,
              },
              {
                path: "platform/demo-requests",
                element: <DemoRequestsPage />,
              },
              {
                path: "platform/settings",
                element: <SettingsPage />,
              },
            ],
          },

          // Organization HQ Dedicated Routes (Protected by OrganizationGuard)
          {
            element: <OrganizationGuard />,
            children: [
              {
                path: "organization/dashboard",
                element: <OrganizationDashboardPage />,
              },
              {
                path: "organization/branches",
                element: <OrganizationBranchesPage />,
              },
              {
                path: "organization/members",
                element: <OrganizationMembersPage />,
              },
              {
                path: "organization/roles",
                element: <OrganizationRolesPage />,
              },
              {
                path: "organization/catalogs",
                element: <OrganizationCatalogsPage />,
              },
              {
                path: "organization/billing",
                element: <OrganizationBillingPage />,
              },
              {
                path: "organization/branding",
                element: <OrganizationBrandingPage />,
              },
              {
                path: "organization/notifications",
                element: <OrganizationNotificationsPage />,
              },
              {
                path: "organization/integrations",
                element: <OrganizationIntegrationsPage />,
              },
              {
                path: "organization/audit",
                element: <OrganizationAuditPage />,
              },
              {
                path: "organization/settings",
                element: <OrganizationSettingsPage />,
              },
            ],
          },

          // Specialized Facility Blueprint Workspace Routes (Protected by WorkspaceGuard)
          {
            element: <WorkspaceGuard />,
            children: [
              {
                path: "workspace/dashboard",
                element: <WorkspaceDashboardPage />,
              },
              {
                path: "workspace/laboratory",
                element: <WorkspaceLaboratoryPage />,
              },
              {
                path: "workspace/clinical",
                element: <WorkspaceClinicalPage />,
              },
              {
                path: "workspace/hospital",
                element: <WorkspaceHospitalPage />,
              },
              {
                path: "workspace/radiology",
                element: <WorkspaceRadiologyPage />,
              },
              {
                path: "workspace/pharmacy",
                element: <WorkspacePharmacyPage />,
              },
              {
                path: "workspace/billing",
                element: <WorkspaceBillingPage />,
              },
            ],
          },
        ],
      },
    ],
  },
  {
    path: "*",
    element: <RootRedirect />,
  },
]);
