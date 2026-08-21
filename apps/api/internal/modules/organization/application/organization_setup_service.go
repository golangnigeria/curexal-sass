package application

import (
	"context"
	"fmt"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
)

type OrganizationSetupService struct {
	orgRepo   domain.OrganizationRepository
	auditRepo auditDomain.AuditRepository
}

func NewOrganizationSetupService(orgRepo domain.OrganizationRepository, auditRepo auditDomain.AuditRepository) *OrganizationSetupService {
	return &OrganizationSetupService{
		orgRepo:   orgRepo,
		auditRepo: auditRepo,
	}
}

func (s *OrganizationSetupService) IsPlatformAdmin(principal *middleware.AuthenticatedPrincipal) bool {
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

func (s *OrganizationSetupService) resolveActiveOrgUUID(principal *middleware.AuthenticatedPrincipal) (uuid.UUID, error) {
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

func (s *OrganizationSetupService) GetProfile(ctx context.Context, principal *middleware.AuthenticatedPrincipal) (*domain.Organization, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	org, err := s.orgRepo.GetByID(ctx, orgUUID)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (s *OrganizationSetupService) UpdateProfile(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	payload *domain.UpdateOrganizationProfilePayload,
) (*domain.Organization, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	existing, errGet := s.orgRepo.GetByID(ctx, orgUUID)
	if errGet != nil {
		return nil, errGet
	}

	updated, errUp := s.orgRepo.UpdateProfile(ctx, orgUUID, payload, actorUUID)
	if errUp != nil {
		return nil, errUp
	}

	// Auto-advance setup state from PENDING_REGISTRATION -> PROFILE_COMPLETED
	if existing.SetupState == domain.SetupStatePendingRegistration {
		adv, errAdv := s.orgRepo.UpdateSetupState(ctx, orgUUID, domain.SetupStateProfileCompleted, 2, actorUUID)
		if errAdv == nil {
			updated = adv
		}
	}

	if s.auditRepo != nil {
		action := "ORGANIZATION_PROFILE_UPDATED"
		resType := "organization.organizations"
		resID := updated.ID.String()
		eventCat := "ORGANIZATION_GOVERNANCE"
		severity := "HIGH"
		status := "SUCCESS"

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    s.IsPlatformAdmin(principal),
			TenantID:      &resID,
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

func (s *OrganizationSetupService) SubmitForReview(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
) (*domain.Organization, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	existing, errGet := s.orgRepo.GetByID(ctx, orgUUID)
	if errGet != nil {
		return nil, errGet
	}

	if !domain.IsValidSetupTransition(existing.SetupState, domain.SetupStateUnderReview) {
		return nil, domain.ErrInvalidSetupStateTransition
	}

	updated, errUp := s.orgRepo.UpdateSetupState(ctx, orgUUID, domain.SetupStateUnderReview, 4, actorUUID)
	if errUp != nil {
		return nil, errUp
	}

	if s.auditRepo != nil {
		action := "ORGANIZATION_SUBMITTED_FOR_REVIEW"
		resType := "organization.organizations"
		resID := updated.ID.String()
		eventCat := "ORGANIZATION_GOVERNANCE"
		severity := "HIGH"
		status := "SUCCESS"

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    s.IsPlatformAdmin(principal),
			TenantID:      &resID,
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

func (s *OrganizationSetupService) VerifyOrganization(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	orgID uuid.UUID,
	approve bool,
	reason string,
) (*domain.Organization, error) {
	if !s.IsPlatformAdmin(principal) {
		return nil, fmt.Errorf("unauthorized: operation requires platform admin privileges")
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	targetState := domain.SetupStateVerified
	step := 5
	if !approve {
		targetState = domain.SetupStateRejected
		step = 3
	}

	updated, errUp := s.orgRepo.UpdateSetupState(ctx, orgID, targetState, step, actorUUID)
	if errUp != nil {
		return nil, errUp
	}

	if s.auditRepo != nil {
		action := "ORGANIZATION_VERIFIED"
		if !approve {
			action = "ORGANIZATION_REJECTED"
		}
		resType := "organization.organizations"
		resID := updated.ID.String()
		eventCat := "PLATFORM_VERIFICATION"
		severity := "CRITICAL"
		status := "SUCCESS"

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    true,
			TenantID:      &resID,
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
