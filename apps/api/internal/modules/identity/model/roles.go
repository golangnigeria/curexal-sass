package model

// Role names — single source of truth for the application.
// These must match the seeded rows in migration 006.
const (
	// Global roles (apply to both HMS and LIMS contexts)
	RoleOwner   = "owner"
	RoleAdmin   = "admin"
	RolePatient = "patient"

	// HMS roles
	RoleDoctor       = "doctor"
	RoleNurse        = "nurse"
	RoleReceptionist = "receptionist"

	// LIMS roles
	RoleLabTech = "lab_tech"

	// Platform-Level roles
	RoleSuperAdmin        = "super_admin"
	RoleSupportAgent      = "super_support_agent"
	RoleSalesStaff        = "super_sales_staff"
	RoleComplianceOfficer = "super_compliance_officer"

	// Organization-Level roles
	RoleRegionalManager = "org_regional_manager"
	RoleQualityManager  = "org_quality_manager"
	RoleFinanceManager  = "org_finance_manager"
	RoleHRManager       = "org_hr_manager"

	// Branch-Level roles (No prefix)
	RoleBranchAdmin  = "branch_admin"
	RoleClinician    = "clinician"
	RoleTechnician   = "technician"
	RoleCustomerCare = "customer_care"
	RoleCashier      = "cashier"
	RoleAccountant   = "accountant"
)

// RoleSystem identifies which sub-application a role belongs to.
type RoleSystem string

const (
	SystemGlobal   RoleSystem = "global"
	SystemHMS      RoleSystem = "hms"
	SystemLIMS     RoleSystem = "lims"
	SystemPlatform RoleSystem = "platform"
)

// AllRoles is the ordered list of every seeded role name.
var AllRoles = []string{
	RoleOwner,
	RoleAdmin,
	RoleDoctor,
	RoleNurse,
	RoleReceptionist,
	RoleLabTech,
	RolePatient,
	RoleSuperAdmin,
	RoleSupportAgent,
	RoleSalesStaff,
	RoleComplianceOfficer,
	RoleRegionalManager,
	RoleQualityManager,
	RoleFinanceManager,
	RoleHRManager,
	RoleBranchAdmin,
	RoleClinician,
	RoleTechnician,
	RoleCustomerCare,
	RoleCashier,
	RoleAccountant,
}

// ClinicalRoles are roles that have patient-facing clinical access.
var ClinicalRoles = []string{
	RoleAdmin,
	RoleClinician,
}

// LabRoles are roles with LIMS access.
var LabRoles = []string{
	RoleAdmin,
	RoleTechnician,
}

// AdminRoles are roles with user/tenant management capability.
var AdminRoles = []string{
	RoleOwner,
	RoleAdmin,
	RoleRegionalManager,
	RoleHRManager,
	RoleBranchAdmin,
	RoleAccountant,
}

// IsValidRole returns true when the given role name exists in the seeded set.
func IsValidRole(role string) bool {
	for _, r := range AllRoles {
		if r == role {
			return true
		}
	}
	return false
}
