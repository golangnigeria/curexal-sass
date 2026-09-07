package application

import (
	"context"
	"errors"

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

func (s *FacilityBranchService) resolveActiveOrgUUID(ctx context.Context, principal *middleware.AuthenticatedPrincipal) (uuid.UUID, error) {
	if principal == nil {
		return uuid.Nil, domain.ErrUnauthorizedTenantAccess
	}

	isPlatform := s.isPlatformAdmin(principal)

	// 1. Try parsing explicit requested organization ID
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
			if isPlatform {
				return parsed, nil
			}
			// Authoritative membership check - NEVER trust client supplied ID without DB verification
			if s.orgRepo != nil && principal.UserID != "" {
				isMember, errCheck := s.orgRepo.VerifyMembership(ctx, parsed, principal.UserID)
				if errCheck == nil && isMember {
					return parsed, nil
				}
			}
			return uuid.Nil, domain.ErrUnauthorizedTenantAccess
		}
	}

	// 2. Fallback to resolving the user's primary active organization from memberships
	if s.orgRepo != nil && principal.UserID != "" {
		orgs, err := s.orgRepo.List(ctx, principal.UserID, isPlatform)
		if err == nil && len(orgs) > 0 {
			return orgs[0].ID, nil
		}
	}

	return uuid.Nil, domain.ErrUnauthorizedTenantAccess
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
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
	if err != nil {
		return nil, err
	}

	return s.branchRepo.ListBranches(ctx, orgUUID)
}

func (s *FacilityBranchService) GetBranchByID(ctx context.Context, principal *middleware.AuthenticatedPrincipal, branchID uuid.UUID) (*domain.FacilityBranch, error) {
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
	if err != nil {
		return nil, err
	}

	isOrgAdmin := s.isPlatformAdmin(principal) || principal.HasPermission("organization:branch:manage") || principal.HasPermission("organization:manage") || principal.Role == "owner" || principal.Role == "org_admin"

	// Authoritative facility access verification
	hasAccess, errAccess := s.branchRepo.VerifyUserFacilityAccess(ctx, orgUUID, branchID, principal.UserID, isOrgAdmin)
	if errAccess != nil || !hasAccess {
		return nil, domain.ErrUnauthorizedTenantAccess
	}

	return s.branchRepo.GetBranchByID(ctx, orgUUID, branchID)
}

func (s *FacilityBranchService) CreateBranch(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	payload *domain.CreateFacilityBranchPayload,
) (*domain.FacilityBranch, error) {
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
	if err != nil {
		return nil, err
	}

	actorUUID := uuid.Nil
	if parsed, err := uuid.Parse(principal.UserID); err == nil {
		actorUUID = parsed
	}

	// 1. Resolve Facility Type ID if passed as Code/Slug or UUID Nil
	if payload.FacilityTypeID == uuid.Nil {
		codeToLookup := payload.FacilityTypeCode
		if codeToLookup == "" {
			codeToLookup = payload.FacilityType
		}
		if codeToLookup == "" {
			codeToLookup = "diagnostic_center"
		}
		ft, errFT := s.branchRepo.GetFacilityTypeByCode(ctx, codeToLookup)
		if errFT != nil || ft == nil {
			return nil, domain.ErrInvalidFacilityType
		}
		payload.FacilityTypeID = ft.ID
	}

	// 2. Verify Facility Type Exists AND is ACTIVE
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
		Slug:           payload.Slug,
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
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
	if err != nil {
		return nil, err
	}

	actorUUID := uuid.Nil
	if parsed, err := uuid.Parse(principal.UserID); err == nil {
		actorUUID = parsed
	}

	existing, errGet := s.branchRepo.GetBranchByID(ctx, orgUUID, branchID)
	if errGet != nil {
		return nil, errGet
	}

	if payload.Name != nil {
		existing.Name = *payload.Name
	}
	if payload.Slug != nil && *payload.Slug != "" {
		existing.Slug = *payload.Slug
		if errVal := existing.Validate(); errVal != nil {
			return nil, errVal
		}
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
	if payload.Version > 0 {
		existing.Version = payload.Version
	}

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
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
	if err != nil {
		return err
	}

	actorUUID := uuid.Nil
	if parsed, err := uuid.Parse(principal.UserID); err == nil {
		actorUUID = parsed
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

func (s *FacilityBranchService) SetHeadquarters(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	branchID uuid.UUID,
) (*domain.FacilityBranch, error) {
	orgUUID, err := s.resolveActiveOrgUUID(ctx, principal)
	if err != nil {
		return nil, err
	}

	actorUUID := uuid.Nil
	if parsed, err := uuid.Parse(principal.UserID); err == nil {
		actorUUID = parsed
	}

	// Verify branch exists in this organization
	branch, errGet := s.branchRepo.GetBranchByID(ctx, orgUUID, branchID)
	if errGet != nil {
		return nil, errGet
	}

	if branch.Status != "ACTIVE" {
		return nil, errors.New("cannot set inactive or suspended branch as headquarters")
	}

	errSet := s.branchRepo.SetHeadquarters(ctx, orgUUID, branchID, actorUUID)
	if errSet != nil {
		return nil, errSet
	}

	if s.auditRepo != nil {
		action := "HEADQUARTERS_CHANGED"
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

	return s.branchRepo.GetBranchByID(ctx, orgUUID, branchID)
}
