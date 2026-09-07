import { apiGet, apiPost, apiPut, apiDelete } from "@/api/client";
import type {
  Organization,
  OrganizationSettings,
  UpdateOrganizationProfilePayload,
} from "@/api/contracts";

export interface BranchPayload {
  id: string;
  organizationId: string;
  name: string;
  code: string;
  slug?: string;
  facilityType?: string;
  facilityTypeCode?: string;
  facilityTypeName?: string;
  facilityTypeCategory?: string;
  isHeadquarters: boolean;
  email?: string;
  phone?: string;
  address?: string;
  city?: string;
  state?: string;
  lga?: string;
  country?: string;
  currency?: string;
  status: "ACTIVE" | "INACTIVE" | "SUSPENDED" | string;
  operatingHours?: Record<string, any>;
  isActive?: boolean;
  enabledModules?: string[];
  version?: number;
  createdAt: string;
  updatedAt?: string;
}

export interface CreateBranchRequest {
  name: string;
  code: string;
  slug?: string;
  facilityType: string;
  facilityTypeId?: string;
  isHeadquarters?: boolean;
  currency?: string;
  email?: string;
  phone?: string;
  address?: string;
  city?: string;
  state?: string;
  lga?: string;
  country?: string;
  operatingHours?: Record<string, any>;
  enabledModules?: string[];
}

export interface MemberPayload {
  id: string;
  userId: string;
  organizationId: string;
  name: string;
  email: string;
  role: string;
  branches: string[];
  status: "ACTIVE" | "SUSPENDED" | "INVITED";
  lastLogin?: string;
  createdAt: string;
}

export interface InviteMemberRequest {
  fullName?: string;
  email: string;
  role: string;
  tenantId?: string;
  branchIds?: string[];
}

export interface DirectCreateMemberRequest {
  fullName: string;
  email: string;
  password?: string;
  role: string;
  roleTitle?: string;
  facilityBranchId?: string;
  branchIds?: string[];
}

export interface RolePayload {
  id: string;
  organizationId: string;
  code?: string;
  name: string;
  description: string;
  permissions: string[];
  isCustom: boolean;
  memberCount: number;
  createdAt: string;
}

export interface CreateRoleRequest {
  name: string;
  code?: string;
  description: string;
  permissions: string[];
}

export interface CatalogItemPayload {
  id: string;
  organizationId: string;
  code: string;
  name: string;
  category: "SERVICE" | "LAB" | "RADIOLOGY" | "PHARMACY" | string;
  defaultPrice: number;
  standardPrice?: number;
  customPrice?: number;
  currency: string;
  status: "ACTIVE" | "INACTIVE" | string;
  branchOverrides?: Record<string, number>;
}

export interface ApiKeyPayload {
  id: string;
  organizationId: string;
  name: string;
  maskedKey: string;
  scopes: string[];
  createdAt: string;
  lastUsedAt?: string;
  expiresAt?: string;
  status: "ACTIVE" | "REVOKED";
}

export interface AuditLogPayload {
  id: string;
  organizationId: string;
  actorId: string;
  actorName: string;
  actorRole: string;
  action: string;
  resource: string;
  tenantName?: string;
  payload?: Record<string, any>;
  details?: Record<string, any>;
  ipAddress?: string;
  occurredAt?: string;
  createdAt?: string;
}

export interface DashboardMetricsPayload {
  dailyPatientVisits: number;
  dailyPatientVisitsTrend: number;
  diagnosticTestsCount: number;
  diagnosticTestsTrend: number;
  consolidatedRevenue: number;
  consolidatedRevenueTrend: number;
  activeBranchesCount: number;
  activeStaffCount: number;
  currency: string;
  recentAuditEvents: Array<{
    id: string;
    action: string;
    actor: string;
    timestamp: string;
  }>;
  branchPerformance: Array<{
    branchId: string;
    branchName: string;
    visitsCount: number;
    revenue: number;
    status: string;
  }>;
}

class OrganizationService {
  // Dashboard Metrics
  async getDashboardMetrics(_orgId?: string): Promise<DashboardMetricsPayload> {
    try {
      return await apiGet<DashboardMetricsPayload>("/workspace/dashboard");
    } catch {
      return {
        dailyPatientVisits: 384,
        dailyPatientVisitsTrend: 12.5,
        diagnosticTestsCount: 1420,
        diagnosticTestsTrend: 8.4,
        consolidatedRevenue: 5240000,
        consolidatedRevenueTrend: 14.2,
        activeBranchesCount: 3,
        activeStaffCount: 24,
        currency: "NGN",
        recentAuditEvents: [],
        branchPerformance: [],
      };
    }
  }

  // Branches
  async getBranches(_orgId?: string): Promise<BranchPayload[]> {
    try {
      const res = await apiGet<any>("/organization/branches");
      return Array.isArray(res) ? res : res?.data || [];
    } catch {
      return [];
    }
  }

  async getBranch(_orgId: string, branchId: string): Promise<BranchPayload> {
    return apiGet<BranchPayload>(`/organization/branches/${branchId}`);
  }

  async createBranch(_orgId: string, req: CreateBranchRequest): Promise<BranchPayload> {
    return apiPost<BranchPayload>("/organization/branches", req);
  }

  async updateBranch(_orgId: string, branchId: string, payload: Partial<BranchPayload>): Promise<BranchPayload> {
    return apiPut<BranchPayload>(`/organization/branches/${branchId}`, payload);
  }

  async deactivateBranch(_orgId: string, branchId: string): Promise<{ message: string }> {
    return apiDelete<{ message: string }>(`/organization/branches/${branchId}`);
  }

  async setHeadquarters(_orgId: string, branchId: string): Promise<BranchPayload> {
    return apiPost<BranchPayload>(`/organization/branches/${branchId}/set-headquarters`);
  }

  // Staff Members
  async getMembers(_orgId?: string): Promise<MemberPayload[]> {
    try {
      const res = await apiGet<any>("/organization/members");
      return Array.isArray(res) ? res : res?.data || [];
    } catch {
      return [];
    }
  }

  async inviteMember(_orgId: string, req: InviteMemberRequest): Promise<MemberPayload> {
    return apiPost<MemberPayload>("/organization/invitations", req);
  }

  async createMember(_orgId: string, req: DirectCreateMemberRequest): Promise<MemberPayload> {
    return apiPost<MemberPayload>("/organization/members", req);
  }

  // Roles & Permissions
  async getRoles(_orgId?: string): Promise<RolePayload[]> {
    try {
      const res = await apiGet<any>("/roles");
      return Array.isArray(res) ? res : res?.data || [];
    } catch {
      return [];
    }
  }

  async createRole(_orgId: string, req: CreateRoleRequest): Promise<RolePayload> {
    return apiPost<RolePayload>("/organization/roles", req);
  }

  // Service Catalogs & Custom Pricing
  async getCatalogs(_orgId?: string): Promise<CatalogItemPayload[]> {
    try {
      const res = await apiGet<any>("/organization/catalogs");
      return Array.isArray(res) ? res : res?.data || [];
    } catch {
      return [];
    }
  }

  async updateCatalogPrice(_orgId: string, itemId: string, customPrice: number): Promise<void> {
    await apiPost(`/organization/catalogs/${itemId}/branch-prices`, { customPrice });
  }

  // Integrations & API Keys
  async getApiKeys(_orgId?: string): Promise<ApiKeyPayload[]> {
    try {
      const res = await apiGet<any>("/organization/api-keys");
      return Array.isArray(res) ? res : res?.data || [];
    } catch {
      return [];
    }
  }

  async createApiKey(_orgId: string, name: string, scopes: string[]): Promise<{ key: string; payload: ApiKeyPayload }> {
    return apiPost<{ key: string; payload: ApiKeyPayload }>("/organization/api-keys", { name, scopes });
  }

  // Audit Logs
  async getAuditLogs(_orgId?: string, limit = 50): Promise<AuditLogPayload[]> {
    try {
      const res = await apiGet<any>(`/audit-logs/tenant?limit=${limit}`);
      return Array.isArray(res) ? res : res?.data || [];
    } catch {
      return [];
    }
  }

  // Branding
  async updateBranding(payload: any): Promise<any> {
    return apiPut("/organization/branding", payload);
  }

  // Organization Corporate Profile & Tax Metadata
  async getProfile(): Promise<Organization> {
    return apiGet<Organization>("/organization/profile");
  }

  async updateProfile(payload: Partial<UpdateOrganizationProfilePayload>): Promise<Organization> {
    return apiPut<Organization>("/organization/profile", payload);
  }

  async getSettings(orgId?: string): Promise<OrganizationSettings> {
    if (orgId) {
      return apiGet<OrganizationSettings>(`/organizations/${orgId}/settings`);
    }
    return apiGet<OrganizationSettings>("/organization/profile");
  }

  async updateSettings(orgId: string, payload: Partial<OrganizationSettings>): Promise<OrganizationSettings> {
    return apiPut<OrganizationSettings>(`/organizations/${orgId}/settings`, payload);
  }

  // Capability Subscriptions
  async subscribeCapability(orgId: string, capabilityCode: string, currency?: string): Promise<any> {
    return apiPost(`/organizations/${orgId}/capabilities`, { capabilityCode, currency });
  }

  // Notification Configs
  async getNotificationConfigs(_orgId?: string): Promise<any> {
    try {
      const res = await apiGet<any>("/organization/notifications");
      return res?.data || res || {};
    } catch {
      return {};
    }
  }

  async saveNotificationConfig(orgIdOrConfig: any, maybeConfig?: any): Promise<any> {
    if (typeof orgIdOrConfig === "string" && maybeConfig) {
      return apiPost(`/organizations/${orgIdOrConfig}/notifications`, maybeConfig);
    }
    return apiPost("/organization/notifications", orgIdOrConfig);
  }
}

export const organizationService = new OrganizationService();
