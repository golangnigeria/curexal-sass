package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrStaffMemberNotFound     = errors.New("staff membership record not found")
	ErrStaffInvitationNotFound = errors.New("staff invitation record not found")
	ErrStaffInvitationExpired  = errors.New("staff invitation has expired")
	ErrDuplicateStaffInvite    = errors.New("a pending invitation already exists for this email in your organization")
	ErrMaxStaffExceeded        = errors.New("organization staff limit reached for active subscription plan")
	ErrInvalidDepartmentCode   = errors.New("invalid department code: must be one of laboratory, clinical, pharmacy, radiology, billing, reception, administration")
)

var AllowedDepartmentCodes = map[string]bool{
	"laboratory":     true,
	"clinical":       true,
	"pharmacy":       true,
	"radiology":      true,
	"billing":        true,
	"reception":      true,
	"administration": true,
}

func IsValidDepartmentCode(code string) bool {
	return AllowedDepartmentCodes[strings.ToLower(strings.TrimSpace(code))]
}

func HashInviteToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}

type StaffMemberDTO struct {
	MembershipID          uuid.UUID                `json:"membershipId"`
	OrganizationID        uuid.UUID                `json:"organizationId"`
	UserID                string                   `json:"userId"`
	Email                 string                   `json:"email"`
	FullName              string                   `json:"fullName"`
	Role                  string                   `json:"role"`
	RoleTitle             string                   `json:"roleTitle"`
	Status                string                   `json:"status"`
	AssignedBranches      []FacilityBranch         `json:"assignedBranches,omitempty"`
	DepartmentAssignments []DepartmentalMembership `json:"departmentAssignments,omitempty"`
	CreatedAt             time.Time                `json:"createdAt"`
}

type MembershipBranch struct {
	ID               uuid.UUID  `json:"id"`
	MembershipID     uuid.UUID  `json:"membershipId"`
	FacilityBranchID uuid.UUID  `json:"facilityBranchId"`
	CreatedAt        time.Time  `json:"createdAt"`
	CreatedBy        *uuid.UUID `json:"createdBy,omitempty"`
}

type DepartmentalMembership struct {
	ID               uuid.UUID  `json:"id"`
	MembershipID     uuid.UUID  `json:"membershipId"`
	FacilityBranchID uuid.UUID  `json:"facilityBranchId"`
	DepartmentCode   string     `json:"departmentCode"`
	CreatedAt        time.Time  `json:"createdAt"`
	CreatedBy        *uuid.UUID `json:"createdBy,omitempty"`
}

type StaffInvitation struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   uuid.UUID  `json:"organizationId"`
	FacilityBranchID *uuid.UUID `json:"facilityBranchId,omitempty"`
	Email            string     `json:"email"`
	Role             string     `json:"role"`
	RoleTitle        string     `json:"roleTitle"`
	InviteTokenHash  string     `json:"-"` // Omit hash from JSON serialization for security
	Status           string     `json:"status"`
	ExpiresAt        time.Time  `json:"expiresAt"`
	InvitedBy        uuid.UUID  `json:"invitedBy"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

type CreateStaffInvitationPayload struct {
	Email            string     `json:"email"`
	Role             string     `json:"role"`
	RoleTitle        string     `json:"roleTitle"`
	FacilityBranchID *uuid.UUID `json:"facilityBranchId"`
}

type CreateStaffInvitationResponse struct {
	Invitation *StaffInvitation `json:"invitation"`
	RawToken   string           `json:"rawToken"` // Exposed ONCE upon creation for delivery
}
