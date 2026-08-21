package model

// Permission names — single source of truth for the application.
// Format: "resource:action"
// These must match the seeded rows in migration 006.
const (
	// ── HMS: Patients ──────────────────────────────────
	PermPatientsRead   = "patients:read"
	PermPatientsWrite  = "patients:write"
	PermPatientsDelete = "patients:delete"

	// ── HMS: Appointments ──────────────────────────────
	PermAppointmentsRead   = "appointments:read"
	PermAppointmentsWrite  = "appointments:write"
	PermAppointmentsDelete = "appointments:delete"

	// ── HMS: Prescriptions ─────────────────────────────
	PermPrescriptionsRead  = "prescriptions:read"
	PermPrescriptionsWrite = "prescriptions:write"

	// ── LIMS: Lab Results ──────────────────────────────
	PermLabResultsRead   = "lab_results:read"
	PermLabResultsWrite  = "lab_results:write"
	PermLabResultsDelete = "lab_results:delete"

	// ── LIMS: Lab Tests ────────────────────────────────
	PermLabTestsRead  = "lab_tests:read"
	PermLabTestsWrite = "lab_tests:write"

	// ── Global: Billing ────────────────────────────────
	PermBillingRead  = "billing:read"
	PermBillingWrite = "billing:write"

	// ── Global: User Management ────────────────────────
	PermUsersRead   = "users:read"
	PermUsersWrite  = "users:write"
	PermUsersDelete = "users:delete"

	// ── Global: Roles & Tenant Management ─────────────
	PermRolesManage  = "roles:manage"
	PermTenantManage = "tenant:manage"

	// ── Global: Own Profile ────────────────────────────
	PermProfileRead  = "profile:read"
	PermProfileWrite = "profile:write"
)

// AllPermissions is the full set of seeded permission names.
var AllPermissions = []string{
	PermPatientsRead, PermPatientsWrite, PermPatientsDelete,
	PermAppointmentsRead, PermAppointmentsWrite, PermAppointmentsDelete,
	PermPrescriptionsRead, PermPrescriptionsWrite,
	PermLabResultsRead, PermLabResultsWrite, PermLabResultsDelete,
	PermLabTestsRead, PermLabTestsWrite,
	PermBillingRead, PermBillingWrite,
	PermUsersRead, PermUsersWrite, PermUsersDelete,
	PermRolesManage, PermTenantManage,
	PermProfileRead, PermProfileWrite,
}

// IsValidPermission returns true when the given permission name exists in the seeded set.
func IsValidPermission(perm string) bool {
	for _, p := range AllPermissions {
		if p == perm {
			return true
		}
	}
	return false
}
