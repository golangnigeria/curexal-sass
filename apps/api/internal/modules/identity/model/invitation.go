package model

import (
	"time"

	"github.com/google/uuid"
)

// Invitation represents a pending invite to join a tenant workspace.
type Invitation struct {
	ID         uuid.UUID  `json:"id"         db:"id"`
	TenantID   uuid.UUID  `json:"tenantId"   db:"tenant_id"`
	Email      string     `json:"email"      db:"email"`
	RoleID     uuid.UUID  `json:"roleId"     db:"role_id"`
	Token      string     `json:"-"          db:"token"`
	InvitedBy  string     `json:"invitedBy"  db:"invited_by"`
	Status     string     `json:"status"     db:"status"`
	ExpiresAt  time.Time  `json:"expiresAt"  db:"expires_at"`
	AcceptedAt *time.Time `json:"acceptedAt" db:"accepted_at"`
	CreatedAt  time.Time  `json:"createdAt"  db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt"  db:"updated_at"`
}

// InviteMemberPayload is the request body for POST /organizations/:id/invite.
type InviteMemberPayload struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// AcceptInvitePayload is the request body for POST /auth/accept-invite.
type AcceptInvitePayload struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Invitation status constants.
const (
	InvitationStatusPending  = "pending"
	InvitationStatusAccepted = "accepted"
	InvitationStatusExpired  = "expired"
	InvitationStatusRevoked  = "revoked"
)
