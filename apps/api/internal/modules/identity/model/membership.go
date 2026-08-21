package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Membership represents a user's membership in a tenant with a specific role.
type Membership struct {
	ID        uuid.UUID  `json:"id"        db:"id"`
	UserID    string     `json:"userId"    db:"user_id"`
	TenantID  uuid.UUID  `json:"tenantId"  db:"tenant_id"`
	RoleID    uuid.UUID  `json:"roleId"    db:"role_id"`
	IsActive  bool       `json:"isActive"  db:"is_active"`
	InvitedBy *string    `json:"invitedBy" db:"invited_by"`
	JoinedAt  *time.Time `json:"joinedAt"  db:"joined_at"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
}

// MembershipWithDetails is returned from API responses and includes
// the role name, tenant name and user details for convenience.
type MembershipWithDetails struct {
	ID         string     `json:"id"         db:"id"`
	UserID     string     `json:"userId"     db:"user_id"`
	UserName   string     `json:"userName"   db:"user_name"`
	UserEmail  string     `json:"userEmail"  db:"user_email"`
	Name       string     `json:"name,omitempty"       db:"name"`
	Email      string     `json:"email,omitempty"      db:"email"`
	Role       string     `json:"role,omitempty"       db:"role"`
	TenantID   string     `json:"tenantId"   db:"tenant_id"`
	TenantName string     `json:"tenantName" db:"tenant_name"`
	RoleID     string     `json:"roleId"     db:"role_id"`
	RoleName   string     `json:"roleName"   db:"role_name"`
	RoleSystem string     `json:"roleSystem" db:"role_system"`
	IsActive   bool       `json:"isActive"   db:"is_active"`
	JoinedAt   *time.Time `json:"joinedAt"   db:"joined_at"`
	CreatedAt  time.Time  `json:"createdAt"  db:"created_at"`
}

// Role represents a named role in the system (seeded or tenant-custom).
type Role struct {
	ID          string     `json:"id"          db:"id"`
	Name        string     `json:"name"        db:"name"`
	Scope       string     `json:"scope"       db:"scope"`
	Description *string    `json:"description" db:"description"`
	System      string     `json:"system"      db:"system"`
	TenantID    *uuid.UUID `json:"tenantId"    db:"tenant_id"`
	IsSystem    bool       `json:"isSystem"    db:"is_system"`
	CreatedAt   time.Time  `json:"createdAt"   db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt"   db:"updated_at"`
}

// RoleWithPermissions is the Role struct extended with a list of permission names.
type RoleWithPermissions struct {
	Role
	Permissions []string `json:"permissions" db:"permissions"`
}

// Permission represents a fine-grained action on a resource.
type Permission struct {
	ID          uuid.UUID `json:"id"          db:"id"`
	Name        string    `json:"name"        db:"name"`
	Description *string   `json:"description" db:"description"`
	Resource    string    `json:"resource"    db:"resource"`
	Action      string    `json:"action"      db:"action"`
	System      string    `json:"system"      db:"system"`
	CreatedAt   time.Time `json:"createdAt"   db:"created_at"`
}

// User is the application representation of a user row.
type User struct {
	ID                  string     `json:"id"                  db:"id"`
	Name                string     `json:"name"                db:"name"`
	Email               string     `json:"email"               db:"email"`
	EmailVerified       bool       `json:"emailVerified"       db:"email_verified"`
	Image               *string    `json:"image"               db:"image"`
	IsPlatformAdmin     bool       `json:"isPlatformAdmin"     db:"is_platform_admin"`
	PlatformRole        *string    `json:"platformRole"        db:"platform_role"`
	FailedLoginAttempts int        `json:"failedLoginAttempts" db:"failed_login_attempts"`
	LockedUntil         *time.Time `json:"lockedUntil"         db:"locked_until"`
	CreatedAt           time.Time  `json:"createdAt"           db:"created_at"`
	UpdatedAt           time.Time  `json:"updatedAt"           db:"updated_at"`
}

type ActiveTenantContext struct {
	TenantID string  `json:"tenantId"`
	BranchID *string `json:"branchId,omitempty"`
}

func (a *ActiveTenantContext) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *ActiveTenantContext) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		s, ok := value.(string)
		if !ok {
			return errors.New("type assertion to []byte or string failed")
		}
		b = []byte(s)
	}
	return json.Unmarshal(b, a)
}

type TenantSelectorItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Session represents a user login session (with revocation support).
type Session struct {
	ID                  string               `json:"id"                  db:"id"`
	UserID              string               `json:"userId"              db:"user_id"`
	Token               string               `json:"token"               db:"token"`
	ExpiresAt           time.Time            `json:"expiresAt"           db:"expires_at"`
	IPAddress           *string              `json:"ipAddress"           db:"ip_address"`
	UserAgent           *string              `json:"userAgent"           db:"user_agent"`
	RevokedAt           *time.Time           `json:"revokedAt"           db:"revoked_at"`
	MfaVerified         bool                 `json:"mfaVerified"         db:"mfa_verified"`
	ActiveTenantContext *ActiveTenantContext `json:"activeTenantContext" db:"active_tenant_context"`
	CreatedAt           time.Time            `json:"createdAt"           db:"created_at"`
	UpdatedAt           time.Time            `json:"updatedAt"           db:"updated_at"`
}

// Me is the response payload for GET /api/v1/me — combines the user profile
// with the current tenant context and membership details.
type Me struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	Email            string               `json:"email"`
	EmailVerified    bool                 `json:"emailVerified"`
	Image            *string              `json:"image"`
	IsPlatformAdmin  bool                 `json:"isPlatformAdmin"`
	PlatformRole     *string              `json:"platformRole"`
	ActiveTenantID   string               `json:"activeTenantId"`
	TenantName       string               `json:"tenantName"`
	TenantSlug       string               `json:"tenantSlug"`
	MembershipID     string               `json:"membershipId"`
	Role             string               `json:"role"`
	Permissions      []string             `json:"permissions"`
	AvailableTenants []TenantSelectorItem `json:"availableTenants"`
	HasPortalAccess  bool                 `json:"hasPortalAccess"`
	Patient          *PatientContext      `json:"patient,omitempty"`
}
