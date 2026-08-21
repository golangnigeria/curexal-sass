package auth

type IdentityVector struct {
	UserID          string `json:"userId"`
	Email           string `json:"email"`
	FullName        string `json:"fullName"`
	AvatarURL       string `json:"avatarUrl,omitempty"`
	IsEmailVerified bool   `json:"isEmailVerified"`
	IsPhoneVerified bool   `json:"isPhoneVerified"`
}

type PlatformVector struct {
	IsPlatformStaff      bool   `json:"isPlatformStaff"`
	IsPlatformAdmin      bool   `json:"isPlatformAdmin"`
	IsSuperAdmin         bool   `json:"isSuperAdmin"`
	PlatformRole         string `json:"platformRole,omitempty"`
	ImpersonationAllowed bool   `json:"impersonationAllowed"`
}

type OrganizationVector struct {
	ActiveOrganizationID string `json:"activeOrganizationId,omitempty"`
	OrganizationName     string `json:"organizationName,omitempty"`
	OrganizationRole     string `json:"organizationRole,omitempty"`
}

type WorkspaceVector struct {
	ActiveWorkspaceID string `json:"activeWorkspaceId,omitempty"`
	WorkspaceName     string `json:"workspaceName,omitempty"`
	WorkspaceRole     string `json:"workspaceRole,omitempty"`
	Department        string `json:"department,omitempty"`
}

type ActiveContextVector struct {
	Type      string `json:"type"`
	ContextID string `json:"contextId,omitempty"`
}

type UserPreferencesVector struct {
	Theme              string `json:"theme"`
	Language           string `json:"language"`
	Timezone           string `json:"timezone"`
	DateFormat         string `json:"dateFormat"`
	NumberFormat       string `json:"numberFormat"`
	DefaultLandingPage string `json:"defaultLandingPage"`
}

type SecurityVector struct {
	Is2FAEnabled       bool   `json:"is2FAEnabled"`
	SessionID          string `json:"sessionId,omitempty"`
	LastPasswordChange string `json:"lastPasswordChange,omitempty"`
}

// AuthenticatedPrincipal is the immutable platform context populated per request by middleware.
type AuthenticatedPrincipal struct {
	UserID         string                `json:"userId"`
	SessionID      string                `json:"sessionId,omitempty"`
	TenantID       string                `json:"tenantId,omitempty"`
	OrganizationID string                `json:"organizationId,omitempty"`
	Role           string                `json:"role,omitempty"`
	Identity       IdentityVector        `json:"identity"`
	Platform       PlatformVector        `json:"platform"`
	Organization   OrganizationVector    `json:"organization"`
	Workspace      WorkspaceVector       `json:"workspace"`
	ActiveContext  ActiveContextVector   `json:"activeContext"`
	Permissions    []string              `json:"permissions"`
	Preferences    UserPreferencesVector `json:"preferences"`
	Security       SecurityVector        `json:"security"`
}

// HasPermission reports whether the principal has the target permission, granting automatic access to super admins, platform staff, or wildcard permission holders.
func (p *AuthenticatedPrincipal) HasPermission(permission string) bool {
	if p == nil {
		return false
	}
	if p.Platform.IsSuperAdmin || p.Platform.IsPlatformStaff || p.Platform.IsPlatformAdmin {
		return true
	}
	for _, userPerm := range p.Permissions {
		if userPerm == "*" || userPerm == permission {
			return true
		}
	}
	return false
}
