package authz

import "github.com/google/uuid"

// ScopeContext represents the discrete operational scope hierarchy.
type ScopeType string

const (
	ScopePlatform     ScopeType = "platform"
	ScopeOrganization ScopeType = "organization"
	ScopeBranch       ScopeType = "branch"
	ScopeWorkspace    ScopeType = "workspace"
	ScopePatient      ScopeType = "patient"
)

// ScopeContext holds the exact tenancy and domain boundaries of an execution.
type ScopeContext struct {
	Scope          ScopeType  `json:"scope"`
	OrganizationID *uuid.UUID `json:"organizationId,omitempty"`
	BranchID       *uuid.UUID `json:"branchId,omitempty"`
	WorkspaceID    *uuid.UUID `json:"workspaceId,omitempty"`
	PatientID      *uuid.UUID `json:"patientId,omitempty"`
	BranchSlug     string     `json:"branchSlug,omitempty"`
	BranchType     string     `json:"branchType,omitempty"`
	Product        string     `json:"product,omitempty"`
}

// Principal defines the security identity of the authenticated user.
type Principal struct {
	UserID          uuid.UUID      `json:"userId"`
	Email           string         `json:"email"`
	FullName        string         `json:"fullName"`
	IsPlatformStaff bool           `json:"isPlatformStaff"`
	IsPlatformAdmin bool           `json:"isPlatformAdmin"`
	Roles           []string       `json:"roles"`
	Permissions     []string       `json:"permissions"`
	Credentials     []string       `json:"credentials"` // e.g. "physician", "medical_laboratory_scientist", "pharmacist"
	Products        []string       `json:"products"`    // e.g. "core", "hms", "lis", "pharmacy", "billing"
	Capabilities    []string       `json:"capabilities"`
	Context         ScopeContext   `json:"context"`
}

// HasPermission checks if the principal possesses a specific RBAC permission.
func (p *Principal) HasPermission(perm string) bool {
	if p == nil {
		return false
	}
	if p.IsPlatformStaff && perm == "platform:admin" {
		return true
	}
	for _, up := range p.Permissions {
		if up == "*" || up == perm {
			return true
		}
	}
	return false
}

// HasCredential checks if the principal holds a verified professional credential.
func (p *Principal) HasCredential(credential string) bool {
	if p == nil {
		return false
	}
	for _, c := range p.Credentials {
		if c == credential {
			return true
		}
	}
	return false
}

// AuthzRequest represents an authorization evaluation request.
type AuthzRequest struct {
	Permission string
	Resource   string
	Action     string
	Scope      ScopeContext
	Credential string
}
