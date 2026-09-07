package auth

import (
	"context"
)

// PermissionResolver abstracts resolution of permissions from underlying providers (Casbin, DB, Redis cache).
type PermissionResolver interface {
	ResolvePermissions(ctx context.Context, principal *AuthenticatedPrincipal) ([]string, error)
	HasPermission(ctx context.Context, principal *AuthenticatedPrincipal, permission string) (bool, error)
}

// MemoryPermissionResolver implements PermissionResolver using a configurable role-permission mapping.
type MemoryPermissionResolver struct {
	rolePermissions map[string][]string
}

// NewMemoryPermissionResolver returns a initialized MemoryPermissionResolver with default role mappings.
func NewMemoryPermissionResolver() *MemoryPermissionResolver {
	r := &MemoryPermissionResolver{
		rolePermissions: make(map[string][]string),
	}
	r.bootstrapDefaults()
	return r
}

func (r *MemoryPermissionResolver) bootstrapDefaults() {
	// Canonical Platform Roles
	r.rolePermissions["super_admin"] = GetAllPermissions()
	r.rolePermissions["super_support_agent"] = []string{
		PermPlatformView,
		PermPlatformImpersonate,
	}
	r.rolePermissions["super_sales_staff"] = []string{
		PermPlatformView,
		PermDemoManage,
	}

	// Canonical Clinic Roles (Day 1 Section 2.2)
	// 1. OWNER: Full authority across all 27 clinic workspace permissions
	r.rolePermissions["owner"] = []string{
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

		// Legacy aliases
		PermissionPasswordWrite,
		PermissionOrganizationRead,
		PermissionOrganizationSettingsWrite,
		PermissionOrganizationCreate,
		PermissionOrganizationDocumentUpload,
		PermissionOrganizationDocumentRead,
	}

	// 2. ORG_ADMIN: Practice Manager / Clinic Administrator
	r.rolePermissions["org_admin"] = []string{
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
		PermWorkspaceBillingRead,
		PermWorkspaceBillingCharge,
		PermWorkspaceBillingRefund,

		// Legacy aliases
		PermissionPasswordWrite,
		PermissionOrganizationRead,
		PermissionOrganizationSettingsWrite,
		PermissionOrganizationDocumentUpload,
		PermissionOrganizationDocumentRead,
	}

	// 3. DOCTOR: Attending Medical Doctor / Specialist
	r.rolePermissions["doctor"] = []string{
		PermUsersRead,
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

		// Legacy aliases
		PermissionPatientView,
		PermissionPatientCreate,
		PermissionPatientUpdate,
		"consultation:write",
		"prescription:write",
	}

	// 4. NURSE: Triage Nurse / Clinical Assistant
	r.rolePermissions["nurse"] = []string{
		PermWorkspacePatientCreate,
		PermWorkspacePatientRead,
		PermWorkspacePatientUpdate,
		PermWorkspaceAppointmentRead,
		PermWorkspaceAppointmentWrite,
		PermWorkspaceQueueManage,
		PermWorkspaceTriageCreate,
		PermWorkspaceTriageRead,
		PermWorkspaceClinicalRead,
		PermWorkspacePrescriptionRead,
		PermWorkspaceDiagnosticRead,
		PermWorkspaceDocumentUpload,
		PermWorkspaceDocumentRead,

		// Legacy aliases
		PermissionPatientView,
		PermissionPatientCreate,
		PermissionPatientUpdate,
	}

	// 5. RECEPTIONIST: Front Desk Officer
	r.rolePermissions["receptionist"] = []string{
		PermWorkspacePatientCreate,
		PermWorkspacePatientRead,
		PermWorkspacePatientUpdate,
		PermWorkspaceAppointmentRead,
		PermWorkspaceAppointmentWrite,
		PermWorkspaceQueueManage,
		PermWorkspaceDocumentUpload,
		PermWorkspaceDocumentRead,

		// Legacy aliases
		PermissionPatientView,
		PermissionPatientCreate,
		PermissionPatientUpdate,
	}

	// 6. CASHIER: Billing & Accounts Officer
	r.rolePermissions["cashier"] = []string{
		PermWorkspacePatientRead,
		PermWorkspacePrescriptionRead,
		PermWorkspaceDiagnosticRead,
		PermWorkspaceBillingRead,
		PermWorkspaceBillingCharge,
		PermWorkspacePosSettle,
		PermWorkspacePosShiftClose,

		// Legacy aliases
		PermissionPatientView,
		PermissionBillingInvoice,
		PermissionBillingPayment,
	}

	// Legacy / Compatibility role mappings
	r.rolePermissions["branch_admin"] = []string{
		PermissionPasswordWrite,
		PermissionOrganizationRead,
		PermissionOrganizationSettingsWrite,
		PermissionUsersRead,
		PermissionUsersWrite,
		PermissionPatientView,
		PermissionPatientCreate,
		PermissionPatientUpdate,
		PermWorkspacePatientRead,
		PermWorkspacePatientCreate,
		PermWorkspacePatientUpdate,
		PermissionLabCreateOrder,
		PermissionOrganizationDocumentUpload,
		PermissionOrganizationDocumentRead,
	}
	r.rolePermissions["clinician"] = r.rolePermissions["doctor"]
	r.rolePermissions["customer_care"] = r.rolePermissions["receptionist"]
	r.rolePermissions["technician"] = []string{
		PermWorkspacePatientRead,
		PermWorkspaceDiagnosticRead,
		PermissionPatientView,
		PermissionLabAccession,
		PermissionLabEnterResult,
		PermissionLabAuthorizeResult,
	}
	r.rolePermissions["member"] = []string{
		PermissionPasswordWrite,
	}
}

func (r *MemoryPermissionResolver) SetRolePermissions(role string, perms []string) {
	r.rolePermissions[role] = perms
}

func (r *MemoryPermissionResolver) ResolvePermissions(ctx context.Context, principal *AuthenticatedPrincipal) ([]string, error) {
	if principal == nil {
		return nil, nil
	}
	if principal.Platform.IsSuperAdmin || principal.Platform.IsPlatformStaff || principal.Platform.IsPlatformAdmin {
		return GetAllPermissions(), nil
	}

	role := principal.Role
	if role == "" {
		role = principal.Organization.OrganizationRole
	}
	if role == "" {
		role = principal.Platform.PlatformRole
	}
	if perms, exists := r.rolePermissions[role]; exists {
		return perms, nil
	}
	return []string{}, nil
}

func (r *MemoryPermissionResolver) HasPermission(ctx context.Context, principal *AuthenticatedPrincipal, permission string) (bool, error) {
	if principal == nil {
		return false, nil
	}
	if principal.Platform.IsSuperAdmin || principal.Platform.IsPlatformStaff || principal.Platform.IsPlatformAdmin {
		return true, nil
	}

	if principal.HasPermission(permission) {
		return true, nil
	}

	perms, err := r.ResolvePermissions(ctx, principal)
	if err != nil {
		return false, err
	}

	for _, p := range perms {
		if p == "*" || p == permission {
			return true, nil
		}
	}
	return false, nil
}
