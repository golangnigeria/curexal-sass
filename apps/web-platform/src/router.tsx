import React, { Suspense, lazy } from "react";
import { createBrowserRouter, Navigate, useParams, useLocation } from "react-router-dom";
import { ProtectedRoute } from "@/components/auth/protected-route";
import { PlatformGuard, RouteLoadingScreen } from "@/components/auth/platform-guard";
import { OrganizationGuard } from "@/components/auth/organization-guard";
import { WorkspaceGuard } from "@/components/auth/workspace-guard";
import { PlatformLayout } from "@/components/layout/platform-layout";
import { RouteErrorBoundary } from "@/components/ui/error-boundary";
import { authClient } from "@curexal/auth";
import { useBootstrap } from "@/api/hooks/use-bootstrap";
import { ROLES } from "@curexal/contracts";
import {
  isPlatformHost,
  isPatientHost,
  getCanonicalPlatformUrl,
  getCanonicalOrgUrl,
  getCanonicalPatientUrl,
} from "@curexal/tenant-context/url-builder";

// --- Lazy-Loaded Route Chunks ---
const withSuspense = (Component: React.ComponentType<any>, message: string = "Loading interface...") => {
  return (props: any) => (
    <Suspense fallback={<RouteLoadingScreen message={message} />}>
      <Component {...props} />
    </Suspense>
  );
};

// Auth & Public
const LoginPage = withSuspense(lazy(() => import("@/pages/auth/login")), "Connecting to secure login...");
const ExchangePage = withSuspense(lazy(() => import("@/pages/auth/exchange")), "Exchanging workspace credentials...");

// Platform Console Domain
const DashboardPage = withSuspense(lazy(() => import("@/pages/platform/dashboard")), "Loading platform dashboard...");
const OrganizationsPage = withSuspense(lazy(() => import("@/pages/platform/organizations")), "Loading organizations directory...");
const OrganizationDetailPage = withSuspense(lazy(() => import("@/pages/platform/organizations/org-detail")), "Loading organization profile...");
const UsersDirectoryPage = withSuspense(lazy(() => import("@/pages/platform/users")), "Loading users directory...");
const MarketplacePage = withSuspense(lazy(() => import("@/pages/platform/marketplace")), "Loading capability catalog...");
const PricingPage = withSuspense(lazy(() => import("@/pages/platform/pricing")), "Loading pricing rules...");
const FacilityTypesPage = withSuspense(lazy(() => import("@/pages/platform/facility-types")), "Loading facility types...");
const MasterCatalogsPage = withSuspense(lazy(() => import("@/pages/platform/catalogs")), "Loading master catalogs...");
const AuditLogsPage = withSuspense(lazy(() => import("@/pages/platform/audit")), "Loading security audit log...");
const DiagnosticsPage = withSuspense(lazy(() => import("@/pages/platform/diagnostics")), "Loading system diagnostics...");
const DemoRequestsPage = withSuspense(lazy(() => import("@/pages/platform/demo-requests")), "Loading demo requests...");
const SettingsPage = withSuspense(lazy(() => import("@/pages/platform/settings")), "Loading platform settings...");

// Organization HQ Domain
const OrganizationDashboardPage = withSuspense(lazy(() => import("@/pages/organization/dashboard")), "Loading executive headquarters...");
const OrganizationBranchesPage = withSuspense(lazy(() => import("@/pages/organization/branches")), "Loading branches directory...");
const OrganizationMembersPage = withSuspense(lazy(() => import("@/pages/organization/members")), "Loading staff directory...");
const OrganizationRolesPage = withSuspense(lazy(() => import("@/pages/organization/roles")), "Loading RBAC role matrix...");
const OrganizationCatalogsPage = withSuspense(lazy(() => import("@/pages/organization/catalogs")), "Loading clinical catalogs...");
const OrganizationBillingPage = withSuspense(lazy(() => import("@/pages/organization/billing")), "Loading billing & subscriptions...");
const OrganizationBrandingPage = withSuspense(lazy(() => import("@/pages/organization/branding")), "Loading white-label branding...");
const OrganizationNotificationsPage = withSuspense(lazy(() => import("@/pages/organization/notifications")), "Loading notification channels...");
const OrganizationIntegrationsPage = withSuspense(lazy(() => import("@/pages/organization/integrations")), "Loading API integrations...");
const OrganizationAuditPage = withSuspense(lazy(() => import("@/pages/organization/audit")), "Loading compliance audit logs...");
const OrganizationSettingsPage = withSuspense(lazy(() => import("@/pages/organization/settings")), "Loading organization settings...");
const OrganizationCompliancePage = withSuspense(lazy(() => import("@/pages/organization/compliance")), "Loading regulatory compliance vault...");

// Workspace Operational Domain (Healthcare Products)
const WorkspaceDashboardPage = withSuspense(lazy(() => import("@/pages/workspace/dashboard")), "Loading workspace overview...");
const WorkspaceClinicalPage = withSuspense(lazy(() => import("@/pages/workspace/clinical")), "Loading outpatient clinic EMR...");
const WorkspaceBillingPage = withSuspense(lazy(() => import("@/pages/workspace/billing")), "Loading cashier billing POS...");
const ReceptionWorkspacePage = withSuspense(lazy(() => import("@/pages/workspace/reception").then(m => ({ default: m.ReceptionWorkspacePage }))), "Loading patient intake MPI...");
const CareAgentDeskWorkspacePage = withSuspense(lazy(() => import("@/pages/workspace/care-desk").then(m => ({ default: m.CareAgentDeskWorkspacePage }))), "Loading triage desk...");
const ClinicalEncounterRoomPage = withSuspense(lazy(() => import("@/pages/workspace/clinical/encounter-room").then(m => ({ default: m.ClinicalEncounterRoomPage }))), "Connecting clinical consultation room...");

function resolveDefaultWorkspaceModule(_facilityType?: string): string {
  return "clinical";
}

function BranchDefaultRedirect() {
  const { branchSlug } = useParams<{ branchSlug: string }>();
  const { data: bootstrap } = useBootstrap();

  const targetBranch = bootstrap?.availableBranches?.find(
    (b: any) => b.slug === branchSlug || b.code?.toLowerCase() === branchSlug?.toLowerCase()
  );
  const facilityType = targetBranch?.facilityType || bootstrap?.branch?.facilityType || bootstrap?.workspace?.facilityType;
  const defaultModule = resolveDefaultWorkspaceModule(facilityType);

  return <Navigate to={`/${branchSlug}/${defaultModule}`} replace />;
}

function LegacyWorkspaceRedirect() {
  const location = useLocation();
  const { data: bootstrap } = useBootstrap();
  const branchSlug =
    bootstrap?.branch?.slug ||
    bootstrap?.branch?.code?.toLowerCase() ||
    bootstrap?.availableBranches?.[0]?.slug ||
    bootstrap?.availableBranches?.[0]?.code?.toLowerCase() ||
    bootstrap?.workspace?.slug ||
    "main";

  const subpath = location.pathname.replace(/^\/workspace\/?/, "") || "";
  let targetModule = subpath;
  if (!targetModule) {
    targetModule = resolveDefaultWorkspaceModule(
      bootstrap?.branch?.facilityType ||
      bootstrap?.availableBranches?.[0]?.facilityType ||
      bootstrap?.workspace?.facilityType
    );
  }

  return <Navigate to={`/${branchSlug}/${targetModule}`} replace />;
}

function RootRedirect() {
  // 0. Patient Portal Subdomain (patient.localhost or patient.curexal.space)
  if (isPatientHost()) {
    if (typeof window !== "undefined") {
      const canonicalPatientUrl = getCanonicalPatientUrl("/dashboard");
      window.location.replace(canonicalPatientUrl);
      return <RouteLoadingScreen message="Connecting to Patient Health Portal..." />;
    }
    return <Navigate to="/login" replace />;
  }

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

  // 1. Platform Staff -> /platform/dashboard on canonical app.localhost:5002 / app.curexal.space
  if (
    effectiveBootstrap?.platform?.isStaff === true ||
    user.isPlatformAdmin === true ||
    user.platformRole === ROLES.SUPER_ADMIN ||
    user.role === ROLES.SUPER_ADMIN
  ) {
    if (typeof window !== "undefined" && !isPlatformHost()) {
      const canonicalPlatformUrl = getCanonicalPlatformUrl("/platform/dashboard");
      if (canonicalPlatformUrl !== "/platform/dashboard") {
        window.location.replace(canonicalPlatformUrl);
        return <RouteLoadingScreen message="Connecting to Platform Console..." />;
      }
    }
    return <Navigate to="/platform/dashboard" replace />;
  }

  // 2. Organization HQ Authorized -> /organization/dashboard on canonical org domain
  const isOrgExecutive =
    effectiveBootstrap?.contexts?.current === "organization" ||
    effectiveBootstrap?.organization?.role === "owner" ||
    effectiveBootstrap?.organization?.role === "org_admin" ||
    effectiveBootstrap?.organization?.role === "org_regional_manager" ||
    user.role === ROLES.OWNER ||
    user.role === ROLES.ORG_ADMIN ||
    user.role === ROLES.ORG_REGIONAL_MANAGER ||
    (user as any).membershipRole === "owner" ||
    (user as any).membershipRole === "org_admin" ||
    (user as any).organizationRole === "owner" ||
    (user as any).organizationRole === "org_admin";

  if (isOrgExecutive) {
    const orgSlug = effectiveBootstrap?.organization?.slug;
    if (orgSlug && typeof window !== "undefined" && isPlatformHost()) {
      const canonicalOrgUrl = getCanonicalOrgUrl(orgSlug, "/organization/dashboard");
      window.location.replace(canonicalOrgUrl);
      return <RouteLoadingScreen message="Connecting to Organization Portal..." />;
    }
    return <Navigate to="/organization/dashboard" replace />;
  }

  // 3. Branch / Workspace Scope -> Dynamic facility-type default module
  const activeBranchSlug =
    effectiveBootstrap?.branch?.slug ||
    effectiveBootstrap?.branch?.code?.toLowerCase() ||
    effectiveBootstrap?.availableBranches?.[0]?.slug ||
    effectiveBootstrap?.availableBranches?.[0]?.code?.toLowerCase() ||
    effectiveBootstrap?.workspace?.slug ||
    "ho-01";
  const defaultModule = resolveDefaultWorkspaceModule(
    effectiveBootstrap?.branch?.facilityType || effectiveBootstrap?.workspace?.facilityType
  );
  return <Navigate to={`/${activeBranchSlug}/${defaultModule}`} replace />;
}

export const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
    errorElement: <RouteErrorBoundary />,
  },
  {
    path: "/auth/exchange",
    element: <ExchangePage />,
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
              {
                path: "organization/compliance",
                element: <OrganizationCompliancePage />,
              },
            ],
          },

          // Legacy /workspace routes (Automatic backward-compatible redirect to branch path)
          {
            path: "workspace/*",
            element: <LegacyWorkspaceRedirect />,
          },

          // Specialized Branch Workspace Routes: /:branchSlug/* (Protected by WorkspaceGuard)
          {
            path: ":branchSlug",
            element: <WorkspaceGuard />,
            children: [
              {
                index: true,
                element: <BranchDefaultRedirect />,
              },
              {
                path: "dashboard",
                element: <WorkspaceDashboardPage />,
              },
              {
                path: "reception",
                element: <ReceptionWorkspacePage />,
              },
              {
                path: "care-desk",
                element: <CareAgentDeskWorkspacePage />,
              },
              {
                path: "patients",
                element: <ReceptionWorkspacePage />,
              },
              {
                path: "clinical",
                element: <WorkspaceClinicalPage />,
              },
              {
                path: "clinical/encounters/:encounterId",
                element: <ClinicalEncounterRoomPage />,
              },
              {
                path: "billing",
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
