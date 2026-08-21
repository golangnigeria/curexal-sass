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
	r.rolePermissions["super_admin"] = GetAllPermissions()
	r.rolePermissions["owner"] = []string{
		PermissionPasswordWrite, PermissionOrganizationRead, PermissionOrganizationSettingsWrite,
		PermissionOrganizationCreate, PermissionUsersRead, PermissionUsersWrite, PermissionAuditRead,
		PermissionOrganizationDocumentUpload, PermissionOrganizationDocumentRead,
	}
	r.rolePermissions["org_admin"] = []string{
		PermissionPasswordWrite, PermissionOrganizationRead, PermissionOrganizationSettingsWrite,
		PermissionUsersRead, PermissionUsersWrite, PermissionAuditRead,
		PermissionOrganizationDocumentUpload, PermissionOrganizationDocumentRead,
	}
	r.rolePermissions["org_regional_manager"] = []string{
		PermissionPasswordWrite, PermissionOrganizationRead, PermissionUsersRead, PermissionAuditRead,
		PermissionOrganizationDocumentUpload, PermissionOrganizationDocumentRead,
	}
	r.rolePermissions["org_quality_manager"] = []string{
		PermissionPasswordWrite, PermissionOrganizationRead, PermissionUsersRead, PermissionAuditRead,
		PermissionOrganizationDocumentUpload, PermissionOrganizationDocumentRead,
	}
	r.rolePermissions["org_finance_manager"] = []string{
		PermissionPasswordWrite, PermissionOrganizationRead, PermissionBillingRead, PermissionBillingWrite,
	}
	r.rolePermissions["org_hr_manager"] = []string{
		PermissionPasswordWrite, PermissionOrganizationRead, PermissionUsersRead, PermissionUsersWrite,
	}
	r.rolePermissions["branch_admin"] = []string{
		PermissionPasswordWrite, PermissionOrganizationRead, PermissionOrganizationSettingsWrite,
		PermissionUsersRead, PermissionUsersWrite, PermissionPatientView, PermissionPatientCreate, PermissionPatientUpdate, PermissionLabCreateOrder,
		PermissionOrganizationDocumentUpload, PermissionOrganizationDocumentRead,
	}
	r.rolePermissions["clinician"] = []string{
		PermissionPasswordWrite, PermissionPatientView, PermissionPatientCreate, PermissionPatientUpdate,
		"consultation:write", "prescription:write",
	}
	r.rolePermissions["technician"] = []string{
		PermissionPasswordWrite, PermissionPatientView, PermissionLabAccession, PermissionLabEnterResult, PermissionLabAuthorizeResult,
	}
	r.rolePermissions["customer_care"] = []string{
		PermissionPasswordWrite, PermissionPatientView, PermissionPatientCreate, PermissionPatientUpdate,
	}
	r.rolePermissions["cashier"] = []string{
		PermissionPasswordWrite, PermissionPatientView, PermissionBillingInvoice, PermissionBillingPayment,
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
