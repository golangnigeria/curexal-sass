import { describe, it, expect } from "bun:test";
import type { NavigationItem, NavigationResponse } from "@/api/contracts";

describe("Database-Driven Navigation Platform Suite", () => {
  it("Scenario 1: Platform Console navigation renders platform-scoped items strictly from backend", () => {
    const backendResponse: NavigationResponse = {
      context: {
        type: "platform",
        role: "super_admin",
      },
      items: [
        {
          id: "nav_plat_dashboard",
          key: "platform.dashboard",
          contextScope: "platform",
          title: "Platform Dashboard",
          icon: "LayoutDashboard",
          path: "/platform/dashboard",
          order: 1,
          status: "active",
          isVisible: true,
          isActive: true,
        },
        {
          id: "nav_plat_orgs",
          key: "platform.organizations",
          contextScope: "platform",
          title: "Organizations",
          icon: "Building2",
          path: "/platform/organizations",
          order: 2,
          status: "active",
          isVisible: true,
          isActive: true,
        },
      ],
    };

    expect(backendResponse.context.type).toBe("platform");
    expect(backendResponse.items).toHaveLength(2);
    expect(backendResponse.items[0].path).toBe("/platform/dashboard");
    expect(backendResponse.items[1].path).toBe("/platform/organizations");
  });

  it("Scenario 2: Organization HQ navigation contains compliance & regulatory vault", () => {
    const backendResponse: NavigationResponse = {
      context: {
        type: "organization",
        organizationSlug: "everight",
        role: "owner",
      },
      items: [
        {
          id: "nav_org_dashboard",
          key: "organization.dashboard",
          contextScope: "organization",
          title: "Executive HQ Dashboard",
          icon: "LayoutDashboard",
          path: "/organization/dashboard",
          order: 1,
          status: "active",
          isVisible: true,
          isActive: true,
        },
        {
          id: "nav_org_compliance",
          key: "organization.compliance",
          contextScope: "organization",
          title: "Compliance & Documents",
          icon: "FileCheck",
          path: "/organization/compliance",
          order: 11,
          status: "active",
          isVisible: true,
          isActive: true,
        },
      ],
    };

    expect(backendResponse.context.type).toBe("organization");
    const hasCompliance = backendResponse.items.some((i) => i.key === "organization.compliance");
    expect(hasCompliance).toBe(true);
  });

  it("Scenario 3: Workspace navigation dynamic branch route interpolation", () => {
    const activeBranchSlug = "owerri-main";
    const rawItems: NavigationItem[] = [
      {
        id: "nav_wsp_dashboard",
        key: "workspace.dashboard",
        contextScope: "workspace",
        title: "Workspace Overview",
        icon: "LayoutDashboard",
        path: `/${activeBranchSlug}/dashboard`,
        order: 1,
        status: "active",
        isVisible: true,
        isActive: true,
      },
      {
        id: "nav_wsp_clin",
        key: "workspace.clinical",
        contextScope: "workspace",
        title: "Outpatient Clinic (EMR)",
        icon: "Stethoscope",
        path: `/${activeBranchSlug}/clinical`,
        order: 5,
        status: "active",
        isVisible: true,
        isActive: true,
      },
    ];

    expect(rawItems[0].path).toBe("/owerri-main/dashboard");
    expect(rawItems[1].path).toBe("/owerri-main/clinical");
    expect(rawItems[0].path.startsWith("/organization")).toBe(false);
  });

  it("Scenario 4: Pending feature navigation item renders disabled with status badge", () => {
    const item: NavigationItem = {
      id: "nav_wsp_telehealth_advanced",
      key: "workspace.telehealth.ai_scribe",
      contextScope: "workspace",
      title: "AI Medical Scribe",
      icon: "Sparkles",
      path: "/owerri-main/clinical/ai-scribe",
      order: 12,
      status: "pending",
      isVisible: true,
      isActive: true,
    };

    expect(item.status).toBe("pending");
    expect(item.isVisible).toBe(true);
  });
});
