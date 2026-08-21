package model

type PermissionScope string

const (
	// Patient Scopes
	PermPatientView   PermissionScope = "patient:view"
	PermPatientCreate PermissionScope = "patient:create"
	PermPatientUpdate PermissionScope = "patient:update"
	PermPatientDelete PermissionScope = "patient:delete"

	// Laboratory LIMS Scopes
	PermLabOrderCreate PermissionScope = "laboratory:create_order"
	PermLabAccession   PermissionScope = "laboratory:accession"
	PermLabResultEnter PermissionScope = "laboratory:enter_result"
	PermLabResultAuth  PermissionScope = "laboratory:authorize_result"

	// Clinic HMS Scopes
	PermConsultationWrite PermissionScope = "consultation:write"
	PermPrescriptionWrite PermissionScope = "prescription:write"

	// Billing Scopes
	PermBillingInvoice PermissionScope = "billing:invoice"
	PermBillingRefund  PermissionScope = "billing:refund"

	// Administration Scopes
	PermOrgManage      PermissionScope = "organization:manage"
	PermSettingsUpdate PermissionScope = "settings:update"
	PermUsersRead      PermissionScope = "users:read"
	PermUsersWrite     PermissionScope = "users:write"
)

type EnforceRequest struct {
	Subject  string `json:"subject"`
	Tenant   string `json:"tenant"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type EnforceResponse struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

type UserPermissionsResponse struct {
	Subject     string   `json:"subject"`
	Tenant      string   `json:"tenant"`
	Permissions []string `json:"permissions"`
}
