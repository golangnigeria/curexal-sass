package model

import (
	"time"
)

// ──────────────────────────────────────────────────────────────────────────────
// Enterprise Session Vector
// The canonical response for any authenticated user — staff, patient, or admin.
// ──────────────────────────────────────────────────────────────────────────────

// UserContext — pure identity. No healthcare or org-specific fields.
type UserContext struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	EmailVerified bool    `json:"emailVerified"`
	Image         *string `json:"image"`
}

// PlatformContext — control-plane authority.
type PlatformContext struct {
	IsPlatformAdmin bool    `json:"isPlatformAdmin"`
	Role            *string `json:"role"`
	HasPortalAccess bool    `json:"hasPortalAccess"`
}

// OrganizationContext — corporate governance scope. nil = no org membership.
type OrganizationContext struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
	Slug *string `json:"slug"`
	Role *string `json:"role"`
}

// WorkspaceContext — active branch execution scope. nil = no workspace.
type WorkspaceContext struct {
	ID           *string `json:"id"`
	Name         *string `json:"name"`
	Slug         *string `json:"slug"`
	MembershipID *string `json:"membershipId"`
}

// PatientContext — patient healthcare persona. nil = user has no patient profile.
type PatientContext struct {
	ProfileID   string  `json:"profileId"`
	Phone       *string `json:"phone"`
	DateOfBirth *string `json:"dateOfBirth"`
	Gender      *string `json:"gender"`
	BloodGroup  *string `json:"bloodGroup"`
	Genotype    *string `json:"genotype"`
	City        *string `json:"city"`
	State       *string `json:"state"`
	Country     string  `json:"country"`
}

// ActiveContext — tracks which persona/UI is currently active for this session.
// persona values: "patient" | "staff" | "provider" | "platform"
type ActiveContext struct {
	Persona string `json:"persona"`
}

// EnterpriseSession is the canonical authenticated session response.
// Returned by POST /auth/sign-in, POST /auth/verify-otp, GET /users/me.
//
// Vector fields are nil when the user does not have that persona:
//   - Organization == nil → user is not an org member
//   - Workspace    == nil → user has no active branch workspace
//   - Patient      == nil → user has no patient profile
type EnterpriseSession struct {
	User             UserContext           `json:"user"`
	Platform         PlatformContext       `json:"platform"`
	Organization     *OrganizationContext  `json:"organization"`
	Workspace        *WorkspaceContext     `json:"workspace"`
	Patient          *PatientContext       `json:"patient"`
	Permissions      []string             `json:"permissions"`
	AvailableTenants []TenantSelectorItem  `json:"availableTenants"`
	ActiveContext    ActiveContext         `json:"activeContext"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Legacy support types (kept for internal use during migration)
// ──────────────────────────────────────────────────────────────────────────────

// UserBaseline represents the immutable user identity credentials.
type UserBaseline struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	FirstName     *string   `json:"firstName,omitempty"`
	LastName      *string   `json:"lastName,omitempty"`
	MiddleName    *string   `json:"middleName,omitempty"`
	EmailVerified bool      `json:"emailVerified"`
	AvatarURL     *string   `json:"avatarUrl"`
	CreatedAt     time.Time `json:"createdAt"`
}

// PlatformCapability represents Platform Control Plane authority.
type PlatformCapability struct {
	IsPlatformAdmin bool    `json:"isPlatformAdmin"`
	Role            *string `json:"role"`
}

// OrganizationSummary represents an organization membership entry.
type OrganizationSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Role string `json:"role"`
}

// IdentityVector contains immutable user identity attributes.
type IdentityVector struct {
	User          UserBaseline          `json:"user"`
	Platform      PlatformCapability    `json:"platform"`
	Organizations []OrganizationSummary `json:"organizations"`
}

// ActiveOrganizationContext represents active corporate governance context.
type ActiveOrganizationContext struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
	Slug *string `json:"slug"`
	Role *string `json:"role"`
}

// ActiveTenantContextSummary represents active execution tenant details.
type ActiveTenantContextSummary struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
	Slug *string `json:"slug"`
	Type string  `json:"type"` // "organization" | "branch"
}

// WorkspaceMembershipContext represents active branch membership details.
type WorkspaceMembershipContext struct {
	MembershipID *string `json:"membershipId"`
	Role         *string `json:"role"`
}

// ContextVector contains mutable active execution scope details.
type ContextVector struct {
	ActiveOrganization  ActiveOrganizationContext   `json:"activeOrganization"`
	ActiveTenant        ActiveTenantContextSummary  `json:"activeTenant"`
	ActiveBranch        ActiveTenantContextSummary  `json:"activeBranch"`
	WorkspaceMembership WorkspaceMembershipContext  `json:"workspaceMembership"`
}

// PrincipalMetadata contains session metadata details.
type PrincipalMetadata struct {
	SessionID string    `json:"sessionId"`
	IssuedAt  time.Time `json:"issuedAt"`
}

// AuthenticatedPrincipal is the canonical session payload for Curexal.
// Serves as the single response contract for both POST /auth/sign-in and GET /users/me.
type AuthenticatedPrincipal struct {
	Identity         IdentityVector       `json:"identity"`
	Context          ContextVector        `json:"context"`
	Permissions      []string             `json:"permissions"`
	Metadata         PrincipalMetadata    `json:"metadata"`
	AvailableTenants []TenantSelectorItem `json:"availableTenants"`
}

// CasbinSubject represents structured subject data for domain authorization.
type CasbinSubject struct {
	UserID              string
	PlatformRole        string
	OrganizationRoles   []string
	WorkspaceMembership string
	TenantID            string
}
