package auth

// Centralized permission constants to eliminate magic string typos across backend.
const (
	// Identity & Auth
	PermissionPasswordWrite = "identity:password:write"

	// Canonical Clinic Workspace Governance
	PermOrganizationView         = "organization:view"
	PermOrganizationManage       = "organization:manage"
	PermOrganizationBranchManage = "organization:branch:manage"
	PermUsersRead                = "users:read"
	PermUsersWrite               = "users:write"
	PermAuditRead                = "audit:read"

	// Canonical Patient & MPI
	PermWorkspacePatientCreate = "workspace:patient:create"
	PermWorkspacePatientRead   = "workspace:patient:read"
	PermWorkspacePatientUpdate = "workspace:patient:update"

	// Canonical Appointments & Operations
	PermWorkspaceAppointmentRead  = "workspace:appointment:read"
	PermWorkspaceAppointmentWrite = "workspace:appointment:write"
	PermWorkspaceQueueManage      = "workspace:queue:manage"

	// Canonical Clinical & Triage
	PermWorkspaceTriageCreate      = "workspace:triage:create"
	PermWorkspaceTriageRead        = "workspace:triage:read"
	PermWorkspaceClinicalRead      = "workspace:clinical:read"
	PermWorkspaceClinicalWrite     = "workspace:clinical:write"
	PermWorkspaceClinicalSign      = "workspace:clinical:sign"
	PermWorkspacePrescriptionWrite = "workspace:prescription:write"
	PermWorkspacePrescriptionRead  = "workspace:prescription:read"
	PermWorkspaceDiagnosticOrder   = "workspace:diagnostic:order"
	PermWorkspaceDiagnosticRead    = "workspace:diagnostic:read"

	// Canonical Documents
	PermWorkspaceDocumentUpload = "workspace:document:upload"
	PermWorkspaceDocumentRead   = "workspace:document:read"

	// Canonical Financial & Billing
	PermWorkspaceBillingRead   = "workspace:billing:read"
	PermWorkspaceBillingCharge = "workspace:billing:charge"
	PermWorkspaceBillingRefund = "workspace:billing:refund"
	PermWorkspacePosSettle     = "workspace:pos:settle"
	PermWorkspacePosShiftClose = "workspace:pos:shift_close"

	// Canonical Platform Super-Admin
	PermPlatformAdmin          = "platform:admin"
	PermPlatformView           = "platform:view"
	PermPlatformManage         = "platform:manage"
	PermPlatformImpersonate    = "platform:impersonate"
	PermPlatformCatalogsManage = "platform:catalogs:manage"
	PermPlatformPricingManage  = "platform:pricing:manage"
	PermDemoManage             = "demo:manage"

	// Platform & Diagnostics (Aliases)
	PermissionPlatformAdmin  = PermPlatformAdmin
	PermissionPlatformHealth = "platform:health"

	// Organization Governance Aliases
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

	// User Roster & Staff Aliases
	PermissionUsersRead  = PermUsersRead
	PermissionUsersWrite = PermUsersWrite

	// Audit & Compliance
	PermissionAuditRead = PermAuditRead

	// Patient & Clinical Aliases
	PermissionPatientView   = PermWorkspacePatientRead
	PermissionPatientCreate = PermWorkspacePatientCreate
	PermissionPatientUpdate = PermWorkspacePatientUpdate

	// Laboratory & LIMS
	PermissionLabCreateOrder     = "laboratory:create_order"
	PermissionLabAccession       = "laboratory:accession"
	PermissionLabEnterResult     = "laboratory:enter_result"
	PermissionLabAuthorizeResult = "laboratory:authorize_result"

	// Billing & Finance Aliases
	PermissionBillingRead    = PermWorkspaceBillingRead
	PermissionBillingWrite   = "billing:write"
	PermissionBillingInvoice = PermWorkspaceBillingCharge
	PermissionBillingPayment = PermWorkspacePosSettle

	// Demo Management
	PermissionDemoRead = PermDemoManage

	// Wildcard All (deprecated - do not return in user permission arrays)
	PermissionWildcard = "*"
)

// GetAllPermissions returns all explicit concrete permission codes across the system.
func GetAllPermissions() []string {
	return []string{
		// 27 Canonical Clinic Workspace Permissions
		PermOrganizationView,
		PermOrganizationManage,
		PermOrganizationBranchManage,
		PermUsersRead,
		PermUsersWrite,
		PermAuditRead,
		PermWorkspacePatientCreate,
		PermWorkspacePatientRead,
		PermWorkspacePatientUpdate,
		PermWorkspaceAppointmentRead,
		PermWorkspaceAppointmentWrite,
		PermWorkspaceQueueManage,
		PermWorkspaceTriageCreate,
		PermWorkspaceTriageRead,
		PermWorkspaceClinicalRead,
		PermWorkspaceClinicalWrite,
		PermWorkspaceClinicalSign,
		PermWorkspacePrescriptionWrite,
		PermWorkspacePrescriptionRead,
		PermWorkspaceDiagnosticOrder,
		PermWorkspaceDiagnosticRead,
		PermWorkspaceDocumentUpload,
		PermWorkspaceDocumentRead,
		PermWorkspaceBillingRead,
		PermWorkspaceBillingCharge,
		PermWorkspaceBillingRefund,
		PermWorkspacePosSettle,
		PermWorkspacePosShiftClose,

		// 7 Canonical Platform Permissions
		PermPlatformAdmin,
		PermPlatformView,
		PermPlatformManage,
		PermPlatformImpersonate,
		PermPlatformCatalogsManage,
		PermPlatformPricingManage,
		PermDemoManage,

		// Legacy / Module compatibility aliases
		PermissionPasswordWrite,
		PermissionPlatformHealth,
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
		PermissionLabCreateOrder,
		PermissionLabAccession,
		PermissionLabEnterResult,
		PermissionLabAuthorizeResult,
		PermissionBillingWrite,
	}
}
