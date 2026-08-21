import React from "react";
import { createBrowserRouter, Navigate } from "react-router-dom";
import { ProtectedRoute } from "@/components/auth/protected-route";
import { PlatformLayout } from "@/components/layout/platform-layout";
import { RouteErrorBoundary } from "@/components/ui/error-boundary";
import { authClient } from "@/lib/auth-client";

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
  const { data: session, isPending } = authClient.useSession();
  if (isPending) return null;
  if (session?.bootstrap?.contexts?.current === "organization") {
    return <Navigate to="/organization/dashboard" replace />;
  }
  if (session?.bootstrap?.contexts?.current === "workspace") {
    return <Navigate to="/workspace/dashboard" replace />;
  }
  return <Navigate to="/platform/dashboard" replace />;
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
          // Platform Console Routes
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

          // Organization HQ Dedicated Routes
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

          // Specialized Facility Blueprint Workspace Routes
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
  {
    path: "*",
    element: <RootRedirect />,
  },
]);

