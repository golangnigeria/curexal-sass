package application

import (
	"context"
	"fmt"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
)

type FacilityBranchService struct {
	branchRepo domain.FacilityBranchRepository
	orgRepo    domain.OrganizationRepository
	auditRepo  auditDomain.AuditRepository
}

func NewFacilityBranchService(
	branchRepo domain.FacilityBranchRepository,
	orgRepo domain.OrganizationRepository,
	auditRepo auditDomain.AuditRepository,
) *FacilityBranchService {
	return &FacilityBranchService{
		branchRepo: branchRepo,
		orgRepo:    orgRepo,
		auditRepo:  auditRepo,
	}
}

func (s *FacilityBranchService) isPlatformAdmin(principal *middleware.AuthenticatedPrincipal) bool {
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

func (s *FacilityBranchService) resolveActiveOrgUUID(principal *middleware.AuthenticatedPrincipal) (uuid.UUID, error) {
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

func (s *FacilityBranchService) getMaxBranchesForPlan(plan string) int {
	switch plan {
	case "smart":
		return 3
	case "optimize":
		return 10
	case "pro":
		return 25
	case "enterprise":
		return 1000 // Multi-facility enterprise tier
	default:
		return 1
	}
}

func (s *FacilityBranchService) ListBranches(ctx context.Context, principal *middleware.AuthenticatedPrincipal) ([]domain.FacilityBranch, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	return s.branchRepo.ListBranches(ctx, orgUUID)
}

func (s *FacilityBranchService) GetBranchByID(ctx context.Context, principal *middleware.AuthenticatedPrincipal, branchID uuid.UUID) (*domain.FacilityBranch, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	return s.branchRepo.GetBranchByID(ctx, orgUUID, branchID)
}

func (s *FacilityBranchService) CreateBranch(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	payload *domain.CreateFacilityBranchPayload,
) (*domain.FacilityBranch, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	// 1. Verify Facility Type Exists AND is ACTIVE
	active, errActive := s.branchRepo.CheckFacilityTypeActive(ctx, payload.FacilityTypeID)
	if errActive != nil || !active {
		return nil, domain.ErrInactiveFacilityType
	}

	// 2. Evaluate Max Branches Entitlement Limit
	org, errOrg := s.orgRepo.GetByID(ctx, orgUUID)
	if errOrg != nil {
		return nil, errOrg
	}
	currentBranchCount, errCount := s.branchRepo.CountActiveBranches(ctx, orgUUID)
	if errCount == nil {
		maxBranches := s.getMaxBranchesForPlan(org.Plan)
		if currentBranchCount >= maxBranches {
			return nil, domain.ErrMaxBranchesExceeded
		}
	}

	// 3. Single HQ Constraint check
	if payload.IsHeadquarters {
		hasHQ, errHQ := s.branchRepo.HasActiveHeadquarters(ctx, orgUUID)
		if errHQ == nil && hasHQ {
			return nil, domain.ErrHeadquartersConflict
		}
	}

	branchEntity := &domain.FacilityBranch{
		OrganizationID: orgUUID,
		FacilityTypeID: payload.FacilityTypeID,
		Code:           payload.Code,
		Name:           payload.Name,
		IsHeadquarters: payload.IsHeadquarters,
		Email:          payload.Email,
		Phone:          payload.Phone,
		Address:        payload.Address,
		City:           payload.City,
		State:          payload.State,
		LGA:            payload.LGA,
		OperatingHours: payload.OperatingHours,
	}

	if payload.Country != nil && *payload.Country != "" {
		branchEntity.Country = *payload.Country
	} else {
		branchEntity.Country = "Nigeria"
	}

	if errVal := branchEntity.Validate(); errVal != nil {
		return nil, errVal
	}

	created, errCreate := s.branchRepo.CreateBranch(ctx, branchEntity, actorUUID)
	if errCreate != nil {
		return nil, errCreate
	}

	if s.auditRepo != nil {
		action := "BRANCH_CREATED"
		resType := "organization.facility_branches"
		resID := created.ID.String()
		eventCat := "ORGANIZATION_OPERATIONS"
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

	return created, nil
}

func (s *FacilityBranchService) UpdateBranch(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	branchID uuid.UUID,
	payload *domain.UpdateFacilityBranchPayload,
) (*domain.FacilityBranch, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	existing, errGet := s.branchRepo.GetBranchByID(ctx, orgUUID, branchID)
	if errGet != nil {
		return nil, errGet
	}

	if payload.Name != nil {
		existing.Name = *payload.Name
	}
	if payload.IsHeadquarters != nil {
		existing.IsHeadquarters = *payload.IsHeadquarters
	}
	if payload.Email != nil {
		existing.Email = payload.Email
	}
	if payload.Phone != nil {
		existing.Phone = payload.Phone
	}
	if payload.Address != nil {
		existing.Address = payload.Address
	}
	if payload.City != nil {
		existing.City = payload.City
	}
	if payload.State != nil {
		existing.State = payload.State
	}
	if payload.LGA != nil {
		existing.LGA = payload.LGA
	}
	if len(payload.OperatingHours) > 0 {
		existing.OperatingHours = payload.OperatingHours
	}
	if payload.Status != nil {
		existing.Status = *payload.Status
	}
	existing.Version = payload.Version

	updated, errUp := s.branchRepo.UpdateBranch(ctx, existing, actorUUID)
	if errUp != nil {
		return nil, errUp
	}

	if s.auditRepo != nil {
		action := "BRANCH_UPDATED"
		resType := "organization.facility_branches"
		resID := updated.ID.String()
		eventCat := "ORGANIZATION_OPERATIONS"
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

func (s *FacilityBranchService) DeactivateBranch(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	branchID uuid.UUID,
) error {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	errDeact := s.branchRepo.DeactivateBranch(ctx, orgUUID, branchID, actorUUID)
	if errDeact != nil {
		return errDeact
	}

	if s.auditRepo != nil {
		action := "BRANCH_DEACTIVATED"
		resType := "organization.facility_branches"
		resID := branchID.String()
		eventCat := "ORGANIZATION_OPERATIONS"
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
