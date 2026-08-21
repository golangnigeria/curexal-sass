package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
)

type StaffMembershipService struct {
	staffRepo  domain.StaffMembershipRepository
	orgRepo    domain.OrganizationRepository
	branchRepo domain.FacilityBranchRepository
	auditRepo  auditDomain.AuditRepository
}

func NewStaffMembershipService(
	staffRepo domain.StaffMembershipRepository,
	orgRepo domain.OrganizationRepository,
	branchRepo domain.FacilityBranchRepository,
	auditRepo auditDomain.AuditRepository,
) *StaffMembershipService {
	return &StaffMembershipService{
		staffRepo:  staffRepo,
		orgRepo:    orgRepo,
		branchRepo: branchRepo,
		auditRepo:  auditRepo,
	}
}

func (s *StaffMembershipService) isPlatformAdmin(principal *middleware.AuthenticatedPrincipal) bool {
	if principal == nil {
		return false
	}
	if principal.Platform.IsPlatformAdmin || principal.Platform.IsSuperAdmin || principal.Platform.IsPlatformStaff {
		return true
	}
	if principal.Role == "super_admin" || principal.Role == "platform_admin" || principal.Role == "platform_staff" {
		return true
	}
	return false
}

func (s *StaffMembershipService) resolveActiveOrgUUID(principal *middleware.AuthenticatedPrincipal) (uuid.UUID, error) {
	if principal == nil {
		return uuid.Nil, domain.ErrUnauthorizedTenantAccess
	}

	orgIDStr := principal.Organization.ActiveOrganizationID
	if orgIDStr == "" {
		orgIDStr = principal.OrganizationID
	}
	if orgIDStr == "" {
		orgIDStr = principal.TenantID
	}

	if orgIDStr == "" {
		return uuid.Nil, domain.ErrUnauthorizedTenantAccess
	}

	parsed, err := uuid.Parse(orgIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid active organization ID: %w", err)
	}

	return parsed, nil
}

func (s *StaffMembershipService) getMaxStaffForPlan(plan string) int {
	switch plan {
	case "smart":
		return 10
	case "optimize":
		return 50
	case "pro":
		return 200
	case "enterprise":
		return 10000
	default:
		return 5
	}
}

func (s *StaffMembershipService) generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate cryptographically secure token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func (s *StaffMembershipService) ListMembers(ctx context.Context, principal *middleware.AuthenticatedPrincipal) ([]domain.StaffMemberDTO, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	return s.staffRepo.ListMembers(ctx, orgUUID)
}

func (s *StaffMembershipService) CreateInvitation(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	payload *domain.CreateStaffInvitationPayload,
) (*domain.CreateStaffInvitationResponse, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	// 1. Check existing pending invite
	pending, errCheck := s.staffRepo.CheckPendingInviteExists(ctx, orgUUID, payload.Email)
	if errCheck == nil && pending {
		return nil, domain.ErrDuplicateStaffInvite
	}

	// 2. Evaluate subscription plan max_staff limit
	org, errOrg := s.orgRepo.GetByID(ctx, orgUUID)
	if errOrg != nil {
		return nil, errOrg
	}
	currentMemberCount, errCount := s.staffRepo.CountActiveMembers(ctx, orgUUID)
	if errCount == nil {
		maxStaff := s.getMaxStaffForPlan(org.Plan)
		if currentMemberCount >= maxStaff {
			return nil, domain.ErrMaxStaffExceeded
		}
	}

	// 3. Generate cryptographically secure token and SHA-256 hash
	rawToken, errGen := s.generateSecureToken()
	if errGen != nil {
		return nil, errGen
	}
	tokenHash := domain.HashInviteToken(rawToken)

	roleTitleVal := "member"
	if payload.RoleTitle != "" {
		roleTitleVal = payload.RoleTitle
	}

	invEntity := &domain.StaffInvitation{
		OrganizationID:   orgUUID,
		FacilityBranchID: payload.FacilityBranchID,
		Email:            payload.Email,
		Role:             payload.Role,
		RoleTitle:        roleTitleVal,
		InviteTokenHash:  tokenHash,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
		InvitedBy:        actorUUID,
	}

	createdInv, errCreate := s.staffRepo.CreateInvitation(ctx, invEntity)
	if errCreate != nil {
		return nil, errCreate
	}

	if s.auditRepo != nil {
		action := "STAFF_INVITED"
		resType := "organization.staff_invitations"
		resID := createdInv.ID.String()
		eventCat := "ORGANIZATION_STAFF_GOVERNANCE"
		severity := "HIGH"
		status := "SUCCESS"
		orgIDStr := orgUUID.String()

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    s.isPlatformAdmin(principal),
			TenantID:      &orgIDStr,
			ActorID:       &principal.UserID,
			ActorName:     &principal.Identity.FullName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	return &domain.CreateStaffInvitationResponse{
		Invitation: createdInv,
		RawToken:   rawToken,
	}, nil
}

func (s *StaffMembershipService) ListInvitations(ctx context.Context, principal *middleware.AuthenticatedPrincipal) ([]domain.StaffInvitation, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	return s.staffRepo.ListInvitations(ctx, orgUUID)
}

func (s *StaffMembershipService) RevokeInvitation(ctx context.Context, principal *middleware.AuthenticatedPrincipal, inviteID uuid.UUID) error {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return err
	}

	errRev := s.staffRepo.RevokeInvitation(ctx, orgUUID, inviteID)
	if errRev != nil {
		return errRev
	}

	if s.auditRepo != nil {
		action := "STAFF_INVITATION_REVOKED"
		resType := "organization.staff_invitations"
		resID := inviteID.String()
		eventCat := "ORGANIZATION_STAFF_GOVERNANCE"
		severity := "HIGH"
		status := "SUCCESS"
		orgIDStr := orgUUID.String()

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    s.isPlatformAdmin(principal),
			TenantID:      &orgIDStr,
			ActorID:       &principal.UserID,
			ActorName:     &principal.Identity.FullName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	return nil
}

func (s *StaffMembershipService) AssignBranch(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	membershipID, branchID uuid.UUID,
) (*domain.MembershipBranch, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	// Verify branch belongs to active org
	_, errBranch := s.branchRepo.GetBranchByID(ctx, orgUUID, branchID)
	if errBranch != nil {
		return nil, errBranch
	}

	// Verify member belongs to active org
	_, errMem := s.staffRepo.GetMemberByID(ctx, orgUUID, membershipID)
	if errMem != nil {
		return nil, errMem
	}

	assignment, errAssign := s.staffRepo.AssignBranch(ctx, membershipID, branchID, actorUUID)
	if errAssign != nil {
		return nil, errAssign
	}

	if s.auditRepo != nil {
		action := "STAFF_BRANCH_ASSIGNED"
		resType := "organization.membership_branches"
		resID := assignment.ID.String()
		eventCat := "ORGANIZATION_STAFF_GOVERNANCE"
		severity := "MEDIUM"
		status := "SUCCESS"
		orgIDStr := orgUUID.String()

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    s.isPlatformAdmin(principal),
			TenantID:      &orgIDStr,
			ActorID:       &principal.UserID,
			ActorName:     &principal.Identity.FullName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	return assignment, nil
}

func (s *StaffMembershipService) AssignDepartment(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	membershipID, branchID uuid.UUID,
	deptCode string,
) (*domain.DepartmentalMembership, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	if !domain.IsValidDepartmentCode(deptCode) {
		return nil, domain.ErrInvalidDepartmentCode
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	// Verify branch belongs to active org
	_, errBranch := s.branchRepo.GetBranchByID(ctx, orgUUID, branchID)
	if errBranch != nil {
		return nil, errBranch
	}

	// Verify member belongs to active org
	_, errMem := s.staffRepo.GetMemberByID(ctx, orgUUID, membershipID)
	if errMem != nil {
		return nil, errMem
	}

	assignment, errAssign := s.staffRepo.AssignDepartment(ctx, membershipID, branchID, deptCode, actorUUID)
	if errAssign != nil {
		return nil, errAssign
	}

	if s.auditRepo != nil {
		action := "STAFF_DEPARTMENT_ASSIGNED"
		resType := "organization.departmental_memberships"
		resID := assignment.ID.String()
		eventCat := "ORGANIZATION_STAFF_GOVERNANCE"
		severity := "MEDIUM"
		status := "SUCCESS"
		orgIDStr := orgUUID.String()

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    s.isPlatformAdmin(principal),
			TenantID:      &orgIDStr,
			ActorID:       &principal.UserID,
			ActorName:     &principal.Identity.FullName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	return assignment, nil
}

func (s *StaffMembershipService) UpdateMemberRole(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	membershipID uuid.UUID,
	newRole, newRoleTitle string,
) (*domain.StaffMemberDTO, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	updated, errUp := s.staffRepo.UpdateMemberRole(ctx, orgUUID, membershipID, newRole, newRoleTitle, actorUUID)
	if errUp != nil {
		return nil, errUp
	}

	if s.auditRepo != nil {
		action := "STAFF_ROLE_UPDATED"
		resType := "organization.organization_memberships"
		resID := membershipID.String()
		eventCat := "ORGANIZATION_STAFF_GOVERNANCE"
		severity := "HIGH"
		status := "SUCCESS"
		orgIDStr := orgUUID.String()

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    s.isPlatformAdmin(principal),
			TenantID:      &orgIDStr,
			ActorID:       &principal.UserID,
			ActorName:     &principal.Identity.FullName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	return updated, nil
}
