package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/golangnigeria/curexal/internal/shared/mailer"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
)

type StaffMembershipService struct {
	staffRepo  domain.StaffMembershipRepository
	orgRepo    domain.OrganizationRepository
	branchRepo domain.FacilityBranchRepository
	auditRepo  auditDomain.AuditRepository
	mailer     *mailer.Mailer
}

func NewStaffMembershipService(
	staffRepo domain.StaffMembershipRepository,
	orgRepo domain.OrganizationRepository,
	branchRepo domain.FacilityBranchRepository,
	auditRepo auditDomain.AuditRepository,
	mailer *mailer.Mailer,
) *StaffMembershipService {
	return &StaffMembershipService{
		staffRepo:  staffRepo,
		orgRepo:    orgRepo,
		branchRepo: branchRepo,
		auditRepo:  auditRepo,
		mailer:     mailer,
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

func (s *StaffMembershipService) resolveActiveOrgUUID(ctx context.Context, principal *middleware.AuthenticatedPrincipal) (uuid.UUID, error) {
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

	if orgIDStr != "" {
		parsed, err := uuid.Parse(orgIDStr)
		if err == nil {
			return parsed, nil
		}
	}

	if s.orgRepo != nil && principal.UserID != "" {
		orgs, err := s.orgRepo.List(ctx, principal.UserID, s.isPlatformAdmin(principal))
		if err == nil && len(orgs) > 0 {
			return orgs[0].ID, nil
		}
	}

	return uuid.Nil, domain.ErrUnauthorizedTenantAccess
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
	return crypto.GenerateAlphanumericCode(6)
}

func (s *StaffMembershipService) ListMembers(ctx context.Context, principal *middleware.AuthenticatedPrincipal) ([]domain.StaffMemberDTO, error) {
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
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
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
	if err != nil {
		return nil, err
	}

	actorUUID := uuid.Nil
	if parsed, err := uuid.Parse(principal.UserID); err == nil {
		actorUUID = parsed
	}

	// 1. Evaluate subscription plan max_staff limit
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

	roleTitleVal := payload.Role
	if payload.RoleTitle != "" && payload.RoleTitle != "member" {
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

	// 5. Dispatch email invitation asynchronously via Resend
	if s.mailer != nil {
		inviteURL := fmt.Sprintf("http://localhost:5002/login?code=%s&email=%s", rawToken, payload.Email)
		go func() {
			ctxMail, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			_ = s.mailer.SendStaffInvitationEmail(ctxMail, payload.Email, org.Name, payload.Role, roleTitleVal, rawToken, inviteURL)
		}()
	}

	return &domain.CreateStaffInvitationResponse{
		Invitation: createdInv,
		RawToken:   rawToken,
	}, nil
}

func (s *StaffMembershipService) DirectCreateMember(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	payload *domain.DirectCreateStaffMemberPayload,
) (*domain.StaffMemberDTO, error) {
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
	if err != nil {
		return nil, err
	}

	actorUUID := uuid.Nil
	if parsed, err := uuid.Parse(principal.UserID); err == nil {
		actorUUID = parsed
	}

	// 1. Evaluate subscription plan max_staff limit
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

	// 2. Hash password
	passwordToUse := payload.Password
	if strings.TrimSpace(passwordToUse) == "" {
		passwordToUse = "password" // Default fallback password
	}
	hash, errHash := crypto.HashPassword(passwordToUse)
	if errHash != nil {
		return nil, fmt.Errorf("failed to hash staff password: %w", errHash)
	}

	// 3. Prepare branch IDs
	branchIDs := payload.BranchIDs
	if payload.FacilityBranchID != nil && *payload.FacilityBranchID != uuid.Nil {
		hasBranch := false
		for _, bID := range branchIDs {
			if bID == *payload.FacilityBranchID {
				hasBranch = true
				break
			}
		}
		if !hasBranch {
			branchIDs = append([]uuid.UUID{*payload.FacilityBranchID}, branchIDs...)
		}
	}

	roleVal := payload.Role
	if roleVal == "" {
		roleVal = "member"
	}
	roleTitleVal := payload.RoleTitle
	if roleTitleVal == "" {
		roleTitleVal = roleVal
	}

	member, errCreate := s.staffRepo.DirectCreateMember(
		ctx,
		orgUUID,
		payload.FullName,
		payload.Email,
		hash,
		roleVal,
		roleTitleVal,
		branchIDs,
		actorUUID,
	)
	if errCreate != nil {
		return nil, errCreate
	}

	// 4. Record audit log
	if s.auditRepo != nil {
		action := "STAFF_DIRECT_PROVISIONED"
		resType := "organization.organization_memberships"
		resID := member.MembershipID.String()
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

	return member, nil
}

func (s *StaffMembershipService) ListInvitations(ctx context.Context, principal *middleware.AuthenticatedPrincipal) ([]domain.StaffInvitation, error) {
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
	if err != nil {
		return nil, err
	}

	return s.staffRepo.ListInvitations(ctx, orgUUID)
}

func (s *StaffMembershipService) RevokeInvitation(ctx context.Context, principal *middleware.AuthenticatedPrincipal, inviteID uuid.UUID) error {
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
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
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
	if err != nil {
		return nil, err
	}

	actorUUID := uuid.Nil
	if parsed, err := uuid.Parse(principal.UserID); err == nil {
		actorUUID = parsed
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
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
	if err != nil {
		return nil, err
	}

	if !domain.IsValidDepartmentCode(deptCode) {
		return nil, domain.ErrInvalidDepartmentCode
	}

	actorUUID := uuid.Nil
	if parsed, err := uuid.Parse(principal.UserID); err == nil {
		actorUUID = parsed
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
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
	if err != nil {
		return nil, err
	}

	actorUUID := uuid.Nil
	if parsed, err := uuid.Parse(principal.UserID); err == nil {
		actorUUID = parsed
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
