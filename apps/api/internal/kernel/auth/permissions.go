package auth

// Centralized permission constants to eliminate magic string typos across backend.
const (
	// Identity & Auth
	PermissionPasswordWrite = "identity:password:write"

	// Platform & Diagnostics
	PermissionPlatformAdmin  = "platform:admin"
	PermissionPlatformHealth = "platform:health"

	// Organization Governance & Document Verification
	PermissionOrganizationRead           = "organization:read"
	PermissionOrganizationSettingsWrite  = "organization:settings:write"
	PermissionOrganizationCreate         = "organization:create"
	PermissionOrganizationUpdate         = "organization:update"
	PermissionOrganizationDelete         = "organization:delete"
	PermissionOrganizationDocumentUpload = "organization:document:upload"
	PermissionOrganizationDocumentRead   = "organization:document:read"
	PermissionOrganizationDocumentReview = "organization:document:review"
	PermissionOrganizationDocumentApprove = "organization:document:approve"
	PermissionOrganizationDocumentReject = "organization:document:reject"
	PermissionOrganizationVerify          = "organization:verify"

	// User Roster & Staff
	PermissionUsersRead  = "users:read"
	PermissionUsersWrite = "users:write"

	// Audit & Compliance
	PermissionAuditRead = "audit:read"

	// Patient & Clinical
	PermissionPatientView   = "patient:view"
	PermissionPatientCreate = "patient:create"
	PermissionPatientUpdate = "patient:update"

	// Laboratory & LIMS
	PermissionLabCreateOrder     = "laboratory:create_order"
	PermissionLabAccession       = "laboratory:accession"
	PermissionLabEnterResult     = "laboratory:enter_result"
	PermissionLabAuthorizeResult = "laboratory:authorize_result"

	// Billing & Finance
	PermissionBillingRead    = "billing:read"
	PermissionBillingWrite   = "billing:write"
	PermissionBillingInvoice = "billing:invoice"
	PermissionBillingPayment = "billing:payment"

	// Demo Management
	PermissionDemoRead = "demo:read"

	// Wildcard All (deprecated - do not return in user permission arrays)
	PermissionWildcard = "*"
)

// GetAllPermissions returns all explicit concrete permission codes across the system.
func GetAllPermissions() []string {
	return []string{
		PermissionPasswordWrite,
		PermissionPlatformAdmin,
		PermissionPlatformHealth,
		"platform:view",
		"platform:manage",
		"platform:impersonate",
		PermissionOrganizationRead,
		PermissionOrganizationSettingsWrite,
		PermissionOrganizationCreate,
		PermissionOrganizationUpdate,
		PermissionOrganizationDelete,
		PermissionOrganizationDocumentUpload,
		PermissionOrganizationDocumentRead,
		PermissionOrganizationDocumentReview,
		PermissionOrganizationDocumentApprove,
		PermissionOrganizationDocumentReject,
		PermissionOrganizationVerify,
		"organization:view",
		"organization:manage",
		PermissionUsersRead,
		PermissionUsersWrite,
		PermissionAuditRead,
		PermissionPatientView,
		PermissionPatientCreate,
		PermissionPatientUpdate,
		"workspace:patient:read",
		"workspace:patient:create",
		PermissionLabCreateOrder,
		PermissionLabAccession,
		PermissionLabEnterResult,
		PermissionLabAuthorizeResult,
		"workspace:sample:receive",
		"workspace:worksheet:update",
		"workspace:result:authorize",
		PermissionBillingRead,
		PermissionBillingWrite,
		PermissionBillingInvoice,
		PermissionBillingPayment,
		"workspace:billing:create",
		PermissionDemoRead,
		"workspace:clinical:read",
		"workspace:settings:manage",
		"consultation:write",
		"prescription:write",
		"appointments:write",
		"billing:refund",
		"settings:read",
		"settings:write",
		"support:impersonate",
		"orgs:read",
		"orgs:write",
	}
}

