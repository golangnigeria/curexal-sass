// ==========================================
// 1. System Roles & Context Scope Constants
// ==========================================

export const ROLES = {
  // Platform Control Plane
  SUPER_ADMIN: "super_admin",
  SUPER_SUPPORT_AGENT: "super_support_agent",
  SUPER_SALES_STAFF: "super_sales_staff",
  PLATFORM_ADMIN: "platform_admin",
  PLATFORM_STAFF: "platform_staff",

  // Canonical Clinic Roles (Day 1 Section 2.2 A)
  OWNER: "owner",
  ORG_ADMIN: "org_admin",
  DOCTOR: "doctor",
  NURSE: "nurse",
  RECEPTIONIST: "receptionist",
  CASHIER: "cashier",

  // Organization Control Plane (Executive HQ / Legacy)
  ORG_REGIONAL_MANAGER: "org_regional_manager",
  ORG_QUALITY_MANAGER: "org_quality_manager",
  ORG_FINANCE_MANAGER: "org_finance_manager",
  ORG_HR_MANAGER: "org_hr_manager",

  // Legacy Workspace Operational Aliases
  BRANCH_ADMIN: "branch_admin",
  CLINICIAN: "clinician",
  TECHNICIAN: "technician",
  CUSTOMER_CARE: "customer_care",
  MEMBER: "member",
} as const;

export type RoleType = (typeof ROLES)[keyof typeof ROLES];

export const CANONICAL_CLINIC_ROLES = [
  ROLES.OWNER,
  ROLES.ORG_ADMIN,
  ROLES.DOCTOR,
  ROLES.NURSE,
  ROLES.RECEPTIONIST,
  ROLES.CASHIER,
] as const;

export type CanonicalClinicRoleType = (typeof CANONICAL_CLINIC_ROLES)[number];

// 34 Granular Permission Registry (Day 1 Section 2.2 B)
export const PERMISSIONS = {
  // Clinic Workspace Governance & Identity (6)
  ORG_VIEW: "organization:view",
  ORG_MANAGE: "organization:manage",
  ORG_BRANCH_MANAGE: "organization:branch:manage",
  USERS_READ: "users:read",
  USERS_WRITE: "users:write",
  AUDIT_READ: "audit:read",

  // Patient & MPI (3)
  WORKSPACE_PATIENT_CREATE: "workspace:patient:create",
  WORKSPACE_PATIENT_READ: "workspace:patient:read",
  WORKSPACE_PATIENT_UPDATE: "workspace:patient:update",

  // Appointments & Operations (3)
  WORKSPACE_APPOINTMENT_READ: "workspace:appointment:read",
  WORKSPACE_APPOINTMENT_WRITE: "workspace:appointment:write",
  WORKSPACE_QUEUE_MANAGE: "workspace:queue:manage",

  // Clinical & Triage (9)
  WORKSPACE_TRIAGE_CREATE: "workspace:triage:create",
  WORKSPACE_TRIAGE_READ: "workspace:triage:read",
  WORKSPACE_CLINICAL_READ: "workspace:clinical:read",
  WORKSPACE_CLINICAL_WRITE: "workspace:clinical:write",
  WORKSPACE_CLINICAL_SIGN: "workspace:clinical:sign",
  WORKSPACE_PRESCRIPTION_WRITE: "workspace:prescription:write",
  WORKSPACE_PRESCRIPTION_READ: "workspace:prescription:read",
  WORKSPACE_DIAGNOSTIC_ORDER: "workspace:diagnostic:order",
  WORKSPACE_DIAGNOSTIC_READ: "workspace:diagnostic:read",

  // Documents (2)
  WORKSPACE_DOCUMENT_UPLOAD: "workspace:document:upload",
  WORKSPACE_DOCUMENT_READ: "workspace:document:read",

  // Financial, POS & Billing (5)
  WORKSPACE_BILLING_READ: "workspace:billing:read",
  WORKSPACE_BILLING_CHARGE: "workspace:billing:charge",
  WORKSPACE_BILLING_REFUND: "workspace:billing:refund",
  WORKSPACE_POS_SETTLE: "workspace:pos:settle",
  WORKSPACE_POS_SHIFT_CLOSE: "workspace:pos:shift_close",

  // Platform Super-Admin (7)
  PLATFORM_ADMIN: "platform:admin",
  PLATFORM_VIEW: "platform:view",
  PLATFORM_MANAGE: "platform:manage",
  PLATFORM_IMPERSONATE: "platform:impersonate",
  PLATFORM_CATALOGS_MANAGE: "platform:catalogs:manage",
  PLATFORM_PRICING_MANAGE: "platform:pricing:manage",
  DEMO_MANAGE: "demo:manage",
} as const;

export type PermissionType = (typeof PERMISSIONS)[keyof typeof PERMISSIONS];

export const CONTEXT_SCOPES = {
  PLATFORM: "platform",
  ORGANIZATION: "organization",
  WORKSPACE: "workspace",
} as const;

export type ContextScopeType = (typeof CONTEXT_SCOPES)[keyof typeof CONTEXT_SCOPES];

// ==========================================
// 2. Identity & Session Contracts
// ==========================================

export interface IdentityPayload {
  id: string;
  email: string;
  displayName: string;
  avatar?: string;
  locale: string;
  timezone: string;
}

export interface PlatformPayload {
  isStaff: boolean;
  role: string;
}

export interface OrganizationPayload {
  id: string;
  name: string;
  slug?: string;
  logo?: string;
  role?: string;
  subscription: string;
  status?: string;
  setupState?: string;
  setup_state?: string;
  hostname?: string;
  isCustomDomain?: boolean;
}

export interface BranchPayload {
  id: string;
  name: string;
  slug: string;
  code: string;
  facilityType: string;
  isHeadquarters: boolean;
  city?: string;
  state?: string;
  operatingHours?: Record<string, any>;
}

export interface BranchSummaryPayload {
  id: string;
  name: string;
  slug: string;
  code: string;
  facilityType: string;
  isHeadquarters: boolean;
}

export interface WorkspacePayload {
  id: string;
  name: string;
  facilityType: string;
  slug: string;
  timezone: string;
  currency: string;
}

export interface SubscriptionPayload {
  plan: string;
  status: string;
  limits: Record<string, number>;
}

export interface ModuleCapabilityPayload {
  code: string;
  enabled: boolean;
  licensed: boolean;
  visible: boolean;
  upgradeAvailable: boolean;
  actions: string[];
}

export type NavigationStatus = "active" | "pending" | "disabled";

export interface NavigationItem {
  id: string;
  key: string;
  contextScope: string;
  moduleCode?: string;
  title: string;
  description?: string;
  icon: string;
  path: string;
  order: number;
  parentId?: string;
  requiredPermission?: string;
  requiredCapability?: string;
  status: NavigationStatus;
  isVisible: boolean;
  isActive: boolean;
  badgeKey?: string;
  badgeCount?: number;
  children?: NavigationItem[];
  createdAt?: string;
}

export interface NavigationContext {
  type: "platform" | "organization" | "workspace" | "patient" | string;
  organizationId?: string;
  organizationSlug?: string;
  branchId?: string;
  branchSlug?: string;
  branchType?: string;
  role?: string;
  capabilities?: string[];
}

export interface NavigationQueryParams {
  scope?: "platform" | "organization" | "workspace" | "patient" | string;
  branch?: string;
  host?: string;
}

export interface NavigationResponse {
  context: NavigationContext;
  items: NavigationItem[];
}

export interface NavigationItemPayload {
  id: string;
  title: string;
  icon: string;
  path: string;
  order?: number;
  children?: NavigationItemPayload[];
}

export interface BreadcrumbPayload {
  title: string;
  path: string;
}

export interface StructuredNavigationPayload {
  primary: NavigationItemPayload[];
  secondary?: NavigationItemPayload[];
  topbar?: NavigationItemPayload[];
  quickActions?: NavigationItemPayload[];
  breadcrumbs?: BreadcrumbPayload[];
}

export interface DashboardSectionPayload {
  id: string;
  title: string;
  type: string;
  items?: string[];
}

export interface DashboardPayload {
  widgets: string[];
  sections?: DashboardSectionPayload[];
}

export interface ContextsPayload {
  current: string;
  available: string[];
  default: string;
}

export interface BrandingPayload {
  logoUrl?: string;
  faviconUrl?: string;
  primaryColor: string;
  secondaryColor?: string;
  accentColor?: string;
  fontFamily?: string;
  borderRadius?: string;
  customDomain?: string;
  hideCurexalBadge?: boolean;
  themeBranding?: Record<string, any>;
}

export interface PreferencesPayload {
  theme: string;
  language: string;
  timezone: string;
}

export interface MetadataPayload {
  version: string;
  generatedAt: string;
  ttl: number;
  etag: string;
}

export interface BootstrapContractResponse {
  identity: IdentityPayload;
  platform: PlatformPayload;
  organization: OrganizationPayload;
  branch?: BranchPayload;
  workspace: WorkspacePayload;
  availableBranches?: BranchSummaryPayload[];
  subscription: SubscriptionPayload;
  modules: ModuleCapabilityPayload[];
  capabilities: string[];
  permissions: string[];
  navigation: NavigationItemPayload[];
  structuredNavigation: StructuredNavigationPayload;
  dashboard: DashboardPayload;
  contexts: ContextsPayload;
  branding: BrandingPayload;
  preferences: PreferencesPayload;
  featureFlags: Record<string, boolean>;
  limits: Record<string, number>;
  metadata: MetadataPayload;
}

export interface UserRoleResponse {
  id: string;
  email: string;
  name: string;
  phone?: string;
  role: string;
  platformRole?: string;
  organizationId?: string;
  organizationRole?: string;
  workspaceId?: string;
  branchId?: string;
  activeBranchId?: string;
  isPlatformAdmin: boolean;
  activeTenantId?: string;
  tenantSlug?: string;
  availableTenants?: Array<{ id: string; name: string; slug: string }>;
  permissions?: string[];
  status?: string;
  createdAt?: string;
}

export interface StaffSignInPayload {
  email: string;
  password: string;
  organizationSlug?: string;
  branchId?: string;
  branchCode?: string;
}

export interface BranchSelectionOption {
  id: string;
  organizationId: string;
  name: string;
  code: string;
  slug: string;
  facilityType?: string;
  isHeadquarters?: boolean;
  isActive: boolean;
}

export interface AuthSuccessResponse {
  accessToken: string;
  tokenType: string;
  expiresIn: number;
  user: UserRoleResponse;
}

export interface BranchSelectionResponse {
  requireBranchSelection: true;
  selectionToken: string;
  assignedBranches: BranchSelectionOption[];
  organization?: {
    id: string;
    name: string;
    slug: string;
  };
}

export interface SelectBranchPayload {
  selectionToken: string;
  branchId: string;
}

export interface SwitchBranchPayload {
  branchId: string;
}

export interface UnassignedBranchError {
  code: "UNASSIGNED_FACILITY_BRANCH";
  message: string;
  organization?: {
    id?: string;
    name?: string;
    slug?: string;
  };
}


export interface DirectoryUser {
  id: string;
  userId: string;
  userName?: string;
  name?: string;
  userEmail?: string;
  email?: string;
  tenantId?: string;
  tenantName?: string;
  organizationId?: string;
  organizationName?: string;
  roleId?: string;
  roleName?: string;
  role?: string;
  roleSystem?: string;
  isActive: boolean;
  joinedAt?: string;
  createdAt: string;
}

// ==========================================
// 3. Diagnostics & Telemetry Contracts
// ==========================================

export interface DatabaseStatus {
  status: string;
  openConnections: number;
  inUse: number;
  idle: number;
}

export interface TimeSeriesPoint {
  month: string;
  count: number;
}

export interface CapabilityDistPoint {
  code: string;
  name: string;
  count: number;
}

export interface TelemetryMetricsData {
  totalOrganizations: number;
  totalWorkspaces: number;
  totalUsers: number;
  organizationsGrowth: TimeSeriesPoint[];
  capabilityDistribution: CapabilityDistPoint[];
}

export interface DiagnosticsMetricsResponse {
  status: string;
  uptimeSeconds: number;
  database: DatabaseStatus;
  metrics: TelemetryMetricsData;
}

export interface LaunchGateCheckResult {
  checkName: string;
  status: "PASSED" | "FAILED";
  details: string;
}

export interface LaunchGateAudit {
  id: string;
  gateName: string;
  status: "PASSED" | "FAILED";
  checkResults: LaunchGateCheckResult[] | any;
  evaluatedAt: string;
  evaluatedBy?: string;
}

export interface SystemHealthMetric {
  id: string;
  componentName: string;
  status: "HEALTHY" | "DEGRADED" | "UNHEALTHY";
  metrics: Record<string, any>;
  checkedAt: string;
}

// ==========================================
// 4. Platform Config & Security Policies
// ==========================================

export interface PlatformGeneralSettings {
  id?: string;
  platformName: string;
  platformLogoUrl?: string;
  platformFaviconUrl?: string;
  platformDescription?: string;
  supportEmail: string;
  supportPhone: string;
  supportWebsite?: string;
  defaultCountry: string;
  defaultCurrency: string;
  supportedCountries: string[];
  supportedCurrencies: string[];
  defaultTimezone: string;
  defaultLocale: string;
  dateFormat: string;
  timeFormat: string;
  numberFormat: string;
  measurementUnits: string;
  maintenanceMode: boolean;
  announcementBanner?: string;
  status?: string;
  version: number;
  createdAt?: string;
  updatedAt?: string;
  updatedBy?: string;
}

export interface IdentitySecurityPolicy {
  id?: string;
  minPasswordLength: number;
  passwordRequireUppercase: boolean;
  passwordRequireNumber: boolean;
  passwordRequireSymbol: boolean;
  passwordExpirationDays: number;
  emailVerificationRequired: boolean;
  mfaPolicy: string;
  maxFailedLoginAttempts: number;
  accountLockoutDurationMinutes: number;
  sessionMaxDurationHours: number;
  refreshTokenDurationDays: number;
  maxActiveSessions: number;
  version: number;
  createdAt?: string;
  updatedAt?: string;
  updatedBy?: string;
}

// ==========================================
// 5. Organizations & Documents
// ==========================================

export interface OrganizationSettings {
  id?: string;
  organizationId: string;
  logoUrl?: string;
  themeBranding?: string;
  customDomain?: string;
  supportEmail?: string;
  supportPhone?: string;
  cacNumber?: string;
  tinNumber?: string;
  taxNumber?: string;
  businessType?: string;
  address?: string;
  timezone?: string;
  currency?: string;
  language?: string;
}

export interface Organization {
  id: string;
  name: string;
  slug: string;
  plan: "smart" | "optimize" | "pro" | "enterprise" | string;
  logoUrl?: string;
  customDomain?: string;
  status: "active" | "inactive" | "pending_verification" | "suspended" | string;
  registrationNumber?: string;
  legalName?: string;
  licenseNumber?: string;
  taxId?: string;
  email?: string;
  phone?: string;
  address?: string;
  city?: string;
  state?: string;
  lga?: string;
  country?: string;
  setupState?: string;
  setupStep?: number;
  completedAt?: string;
  ownerId?: string;
  settings?: OrganizationSettings | Record<string, any>;
  currency?: string;
  timezone?: string;
  version?: number;
  memberCount?: number;
  branchCount?: number;
  createdAt: string;
  updatedAt: string;
}

export interface UpdateOrganizationProfilePayload {
  name?: string;
  legalName?: string;
  registrationNumber?: string;
  licenseNumber?: string;
  taxId?: string;
  email?: string;
  phone?: string;
  address?: string;
  city?: string;
  state?: string;
  lga?: string;
  country?: string;
  logoUrl?: string;
  customDomain?: string;
  currency?: string;
  timezone?: string;
  version?: number;
}

export interface CreateOrganizationPayload {
  name: string;
  slug: string;
  plan?: string;
  email?: string;
  phone?: string;
  address?: string;
  city?: string;
  state?: string;
  country?: string;
  owner?: {
    email: string;
    name?: string;
  };
  ownerEmail?: string;
  ownerName?: string;
}

export interface UpdateOrganizationPayload {
  name?: string;
  slug?: string;
  plan?: string;
  customDomain?: string;
  settings?: Record<string, any>;
}

export interface OrganizationDocument {
  id: string;
  organizationId: string;
  documentType: string;
  originalFilename?: string;
  fileName: string;
  storageKey?: string;
  mimeType: string;
  fileSize?: number;
  fileSizeBytes?: number;
  checksumSha256?: string;
  uploadedBy?: string;
  uploadedAt?: string;
  status: "pending" | "approved" | "rejected" | "expired" | string;
  version?: number;
  reviewedBy?: string;
  reviewedAt?: string;
  rejectionReason?: string;
  presignedUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export interface DocumentWithPresignedURL {
  document: OrganizationDocument;
  presignedUrl: string;
}

// ==========================================
// 6. Capabilities & Marketplace
// ==========================================

export interface CapabilityPrice {
  id?: string;
  capabilityId?: string;
  monthlyPrice: number;
  annualPrice: number;
  currency: string;
}

export interface CapabilityEntity {
  id?: string;
  code: string;
  name: string;
  category: string;
  description?: string;
  isAddOn?: boolean;
}

export interface CapabilityCatalogItem {
  id?: string;
  code?: string;
  name?: string;
  category?: string;
  description?: string;
  basePrice?: number;
  monthlyPrice?: number;
  annualPrice?: number;
  currency?: string;
  isAddOn?: boolean;
  requiredCapabilities?: string[];
  capability?: CapabilityEntity;
  prices?: CapabilityPrice[];
  dependencies?: string[];
  alreadyIncluded?: boolean;
  isEffective?: boolean;
}

export interface OrganizationEntitlement {
  id: string;
  organizationId: string;
  capabilityCode: string;
  source: "plan" | "addon" | "override" | "trial" | string;
  status: "active" | "expired" | "revoked";
  expiresAt?: string;
  createdAt: string;
}

export interface EntitlementTrace {
  capabilityCode: string;
  isEffective: boolean;
  grantedBy: string;
  sourceType: string;
  expiresAt?: string;
  ruleHierarchy: Array<{
    source: string;
    allowed: boolean;
    reason: string;
  }>;
}

// ==========================================
// 7. Pricing Rules & Payment Gateways
// ==========================================

export interface PricingRule {
  id?: string;
  targetType: "plan" | "capability";
  targetCode: string;
  currency: string;
  monthlyPrice: number;
  annualPrice: number;
  vatPercentage: number;
  isActive: boolean;
  version: number;
  createdAt?: string;
  updatedAt?: string;
  updatedBy?: string;
}

export interface PaymentGatewayConfig {
  id?: string;
  providerCode: "paystack" | "flutterwave" | "stripe" | "mock" | string;
  name: string;
  isEnabled: boolean;
  priority: number;
  supportedCurrencies: string[];
  secretKey?: string;
  publicKey?: string;
  webhookSecret?: string;
  version: number;
  createdAt?: string;
  updatedAt?: string;
  updatedBy?: string;
}

export interface UpdateGatewayPayload {
  name?: string;
  isEnabled?: boolean;
  priority?: number;
  supportedCurrencies?: string[];
  secretKey?: string;
  publicKey?: string;
  webhookSecret?: string;
  version: number;
}

// ==========================================
// 8. Facility Types
// ==========================================

export interface FacilityCapability {
  id?: string;
  facilityTypeId?: string;
  capabilityId?: string;
  capabilityCode: string;
  capabilityName: string;
  isDefault: boolean;
}

export interface FacilityTypeEntity {
  id?: string;
  code: string;
  name: string;
  category: "clinical" | "diagnostic" | "retail" | "research" | string;
  iconKey: string;
  description: string;
  isActive: boolean;
  defaultCapabilities?: FacilityCapability[];
  version: number;
  createdAt?: string;
  updatedAt?: string;
  updatedBy?: string;
}

export interface RegistrationFormDTO {
  facilityTypeId: string;
  version: number;
  sections: any;
  requiredDocuments?: any;
}

export interface NavigationMenuDTO {
  facilityTypeId: string;
  menuItems: any;
}

export interface SetupStepDTO {
  stepNumber: number;
  title: string;
  description: string;
  fieldSchema: any;
}

export interface DashboardDTO {
  facilityTypeId: string;
  widgets: any;
}

// ==========================================
// 9. Master Reference Catalogs
// ==========================================

export type CatalogDomain = "clinical" | "lab" | "radiology" | "pharmacy";

export interface CatalogItem {
  id?: string;
  domain: CatalogDomain;
  category: string;
  code: string;
  name: string;
  description?: string;
  systemGroup?: string;
  basePrice?: number;
  metadata?: any;
  isActive: boolean;
  version: number;
  createdAt?: string;
  updatedAt?: string;
  updatedBy?: string;
}

export interface ICD10Code {
  code: string;
  description: string;
  category: string;
}

// ==========================================
// 10. Audit Logs
// ==========================================

export interface AuditLog {
  id: string;
  organizationId?: string;
  facilityBranchId?: string;
  patientId?: string;
  tenantId?: string;
  actorId?: string;
  actorEmail?: string;
  actorRole?: string;
  ipAddress?: string;
  userAgent?: string;
  category: string;
  severity: "info" | "warn" | "error" | "critical" | string;
  action: string;
  resourceType: string;
  resourceId?: string;
  status: "success" | "failure" | string;
  isBreakGlass?: boolean;
  prevRecordHash?: string;
  recordHash?: string;
  details?: Record<string, any>;
  createdAt: string;
}

export interface ListAuditLogsPayload {
  limit?: number;
  offset?: number;
  category?: string;
  severity?: string;
  status?: string;
  actorId?: string;
  action?: string;
  resourceType?: string;
  resourceId?: string;
  organizationId?: string;
  startDate?: string;
  endDate?: string;
  search?: string;
}

export interface AuditAdminStats {
  totalEvents: number;
  errorEvents: number;
  warningEvents: number;
  activeActors: number;
  recentSecurityIncidents: number;
}

// ==========================================
// 11. Demo Requests
// ==========================================

export interface DemoRequest {
  id: string;
  laboratoryName: string;
  contactName: string;
  email: string;
  phone?: string;
  dailyVolume?: string;
  notes?: string;
  status: "pending" | "scheduled" | "completed" | string;
  createdAt: string;
  updatedAt: string;
}

// ==========================================
// 12. Patient Domain & Care Orchestration
// ==========================================

export type MatchConfidence = "NONE" | "LOW" | "PROBABLE_DUPLICATE" | "EXACT_MATCH";

export interface PatientContact {
  id?: string;
  patientId?: string;
  system: "PHONE" | "EMAIL" | "WHATSAPP" | string;
  value: string;
  useType: "MOBILE" | "HOME" | "WORK" | "EMERGENCY" | string;
  isPrimary: boolean;
  verifiedAt?: string;
  createdAt?: string;
}

export interface PortalAccount {
  id: string;
  patientId: string;
  tenantId: string;
  identifier: string;
  status: "INVITED" | "PENDING_VERIFICATION" | "ACTIVE" | "LOCKED" | "SUSPENDED" | string;
  mfaEnabled: boolean;
  lastLoginAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Consent {
  id: string;
  patientId: string;
  consentType: "TELEHEALTH" | "DATA_SHARING" | "RESEARCH" | "PROXY_ACCESS" | string;
  status: "ACTIVE" | "REVOKED" | "EXPIRED" | string;
  grantedBy: string;
  grantedAt: string;
  expiresAt?: string;
  metadata?: Record<string, any>;
}

export interface CanonicalPatient {
  id: string;
  userId?: string;
  tenantId: string;
  mrn: string;
  firstName: string;
  middleName?: string;
  lastName: string;
  gender: "MALE" | "FEMALE" | "OTHER" | string;
  dateOfBirth: string;
  bloodGroup?: string;
  genotype?: string;
  nin?: string;
  status: "DISCOVERED" | "IDENTIFIED" | "REGISTERED" | "SUSPENDED" | "DEACTIVATED" | string;
  registrationChannel: "RECEPTION" | "PORTAL" | "REFERRAL" | "API" | string;
  metadata?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
  contacts?: PatientContact[];
  portal?: PortalAccount;
}

export interface DuplicateEvaluationRequest {
  firstName: string;
  middleName?: string;
  lastName: string;
  dateOfBirth?: string;
  phone?: string;
  email?: string;
  nin?: string;
}

export interface DuplicateMatchCandidate {
  patientId: string;
  mrn: string;
  firstName: string;
  lastName: string;
  dateOfBirth: string;
  gender: string;
  matchedSignals: string[];
  confidenceScore: number;
  confidenceLevel: MatchConfidence;
}

export interface DuplicateEvaluationResponse {
  matchStatus: MatchConfidence;
  candidates: DuplicateMatchCandidate[];
}

export interface RegisterCanonicalPatientPayload {
  firstName: string;
  middleName?: string;
  lastName: string;
  gender: "MALE" | "FEMALE" | "OTHER" | string;
  dateOfBirth: string; // YYYY-MM-DD
  phone: string;
  email?: string;
  bloodGroup?: string;
  genotype?: string;
  nin?: string;
  address?: string;
  registrationChannel?: string;
  forceRegistration?: boolean;
}

export interface PatientListFilter {
  query?: string;
  status?: string;
  gender?: string;
  limit?: number;
  offset?: number;
}

export interface PatientListResponse {
  items: CanonicalPatient[];
  total: number;
  limit: number;
  offset: number;
}

export interface SendPortalOTPPayload {
  identifier: string;
}

export interface VerifyPortalOTPPayload {
  identifier: string;
  code: string;
}

export interface SetPortalPINPayload {
  pin: string;
}

export interface PortalAuthResponse {
  token?: string;
  accessToken?: string;
  refreshToken?: string;
  patientId?: string;
  mrn?: string;
  firstName?: string;
  lastName?: string;
  status?: string;
  hasPin?: boolean;
  expiresAt?: string;
  patient?: CanonicalPatient;
}

export interface CareRequest {
  id: string;
  tenantId: string;
  patientId: string;
  patientName?: string;
  mrn?: string;
  requestNumber: string;
  serviceType: "GENERAL_CONSULTATION" | "SPECIALIST" | "LAB_TEST" | "REFILL" | "TELEHEALTH" | string;
  preferredMode: "IN_PERSON" | "VIDEO" | "AUDIO" | "ASYNC_CHAT" | string;
  urgency: "ROUTINE" | "URGENT" | "EMERGENCY" | string;
  status: "SUBMITTED" | "TRIAGED" | "ASSIGNED_AGENT" | "MATCHED" | "IN_PROGRESS" | "COMPLETED" | "CANCELLED" | string;
  chiefComplaint?: string;
  symptomsJson?: any[];
  preferredTimeWindow?: any;
  assignedCareAgentId?: string;
  matchedProviderId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CareRequestFilter {
  status?: string;
  urgency?: string;
  serviceType?: string;
  limit?: number;
  offset?: number;
}

export interface CareRequestListResponse {
  items: CareRequest[];
  total: number;
  limit: number;
  offset: number;
}

export interface TriageAssessment {
  id: string;
  careRequestId: string;
  patientId: string;
  assessorId?: string;
  acuityLevel: "GREEN" | "YELLOW" | "RED";
  systolicBp?: number;
  diastolicBp?: number;
  pulseRate?: number;
  temperature?: number;
  spo2?: number;
  respiratoryRate?: number;
  painScore?: number;
  triageNotes?: string;
  createdAt: string;
}

export interface SubmitTriagePayload {
  systolicBp?: number;
  diastolicBp?: number;
  pulseRate?: number;
  temperature?: number;
  spo2?: number;
  respiratoryRate?: number;
  painScore?: number;
  triageNotes?: string;
  acuityOverride?: "GREEN" | "YELLOW" | "RED" | string;
}

export interface MatchedProviderCandidate {
  providerId: string;
  userId: string;
  providerName: string;
  specialtyCode: string;
  status: "ON_DUTY" | "ON_BREAK" | "OFF_DUTY" | "BUSY" | string;
  currentActiveQueue: number;
  maxActiveQueue: number;
  matchScore: number;
  matchingReasons: string[];
}

export interface ProviderMatchingResponse {
  careRequestId: string;
  requiredMode: string;
  specialty: string;
  candidates: MatchedProviderCandidate[];
}

export interface CareJourneyMilestone {
  id: string;
  tenantId: string;
  patientId: string;
  encounterId?: string;
  stageCode: "INTAKE" | "TRIAGE" | "CONSULTATION" | "LAB_WORKLIST" | "RADIOLOGY_STUDY" | "PHARMACY_DISPENSE" | "SETTLEMENT" | "FOLLOW_UP" | string;
  title: string;
  description?: string;
  status: "PENDING" | "IN_PROGRESS" | "COMPLETED" | "SKIPPED" | "BLOCKED" | string;
  blockingReason?: string;
  completedAt?: string;
  createdAt: string;
}

// ==========================================
// 13. Healthcare Products & Operational Workspaces
// ==========================================

// --- Laboratory (LIS / LIMS) ---
export interface LabResultEntry {
  parameter: string;
  value: string | number;
  unit: string;
  referenceRange?: string;
  refRange?: string;
  flag?: "NORMAL" | "HIGH" | "LOW" | "PANIC" | string;
  deltaCheckWarning?: string;
}

export interface LabSpecimen {
  id: string;
  barcode: string;
  patientId?: string;
  patientName?: string;
  testCode: string;
  testName: string;
  specimenType: string;
  status: "ordered" | "collected" | "accessioned" | "in_analysis" | "results_pending" | "authorized" | "released" | "rejected";
  priority: "routine" | "urgent" | "stat";
  collectedAt?: string;
  receivedAt?: string;
  authorizedAt?: string;
  authorizedBy?: string;
  results?: LabResultEntry[];
}

export interface LabSamplePayload {
  id: string;
  orderNumber?: string;
  sampleId?: string;
  barcode?: string;
  patientName: string;
  patientMrn: string;
  testName: string;
  testCode?: string;
  specimenType: string;
  status: "ordered" | "collected" | "accessioned" | "in_analysis" | "results_pending" | "authorized" | "released";
  priority: "routine" | "stat" | "urgent";
  receivedAt: string;
  collectedAt?: string;
  results?: Array<{
    parameter: string;
    value: string;
    unit: string;
    referenceRange?: string;
    refRange?: string;
    flag?: "NORMAL" | "HIGH" | "LOW" | "PANIC";
  }>;
}

// --- Outpatient Clinic & EMR / HMS ---
export interface ConsultationQueuePayload {
  id: string;
  queueNumber: number;
  patientName: string;
  patientMrn: string;
  triageCategory: "emergency" | "urgent" | "routine";
  vitals?: { bp: string; pulse: string; temp: string; weight: string; spo2: string };
  status: "waiting" | "in_consultation" | "completed";
  waitingSince: string;
}

export interface SoapEncounterNote {
  patientId: string;
  encounterId?: string;
  subjective: string;
  objective: string;
  assessment: string;
  plan: string;
  signedBy?: string;
  signedAt?: string;
}

// --- Inpatient & Hospital Management ---
export interface HospitalBedPayload {
  id: string;
  wardName: string;
  bedNumber: string;
  status: "occupied" | "available" | "cleaning" | "maintenance";
  patientName?: string;
  patientMrn?: string;
  admissionDate?: string;
  attendingPhysician?: string;
}

// --- Radiology (RIS & PACS) ---
export interface RadiologyScanPayload {
  id: string;
  accessionNumber: string;
  patientName: string;
  patientMrn: string;
  modality: "X-RAY" | "ULTRASOUND" | "CT" | "MRI" | "MAMMOGRAPHY" | string;
  procedureName: string;
  status: "scheduled" | "in_progress" | "images_acquired" | "reported" | "verified";
  scheduledAt: string;
  radiologistName?: string;
  pacsUrl?: string;
}

// --- Pharmacy & Dispensary ---
export interface PharmacyPrescriptionItem {
  medicationName?: string;
  drugName?: string;
  dosage: string;
  quantity: number;
  batchNumber?: string;
  expiryDate?: string;
}

export interface PharmacyPrescriptionPayload {
  id: string;
  rxNumber?: string;
  prescriptionNumber?: string;
  patientName: string;
  patientMrn?: string;
  prescribingDoctor: string;
  status: "pending_verification" | "pending_dispense" | "dispensed" | "cancelled";
  items: PharmacyPrescriptionItem[];
  createdAt: string;
}

// --- Cashier POS & Billing ---
export interface InvoiceItem {
  id: string;
  description: string;
  quantity: number;
  unitPrice: number;
  total: number;
}

export interface InvoicePayload {
  id: string;
  invoiceNumber: string;
  patientName: string;
  patientMrn: string;
  items: InvoiceItem[];
  subtotal: number;
  taxAmount: number;
  totalAmount: number;
  amountPaid: number;
  status: "draft" | "pending" | "partially_paid" | "paid" | "voided";
  createdAt: string;
}
