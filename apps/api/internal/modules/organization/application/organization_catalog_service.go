package application

import (
	"context"
	"fmt"
	"strings"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
)

type OrganizationCatalogService struct {
	catalogRepo domain.OrganizationCatalogRepository
	branchRepo  domain.FacilityBranchRepository
	orgRepo     domain.OrganizationRepository
	auditRepo   auditDomain.AuditRepository
}

func NewOrganizationCatalogService(
	catalogRepo domain.OrganizationCatalogRepository,
	branchRepo domain.FacilityBranchRepository,
	orgRepo domain.OrganizationRepository,
	auditRepo auditDomain.AuditRepository,
) *OrganizationCatalogService {
	return &OrganizationCatalogService{
		catalogRepo: catalogRepo,
		branchRepo:  branchRepo,
		orgRepo:     orgRepo,
		auditRepo:   auditRepo,
	}
}

func (s *OrganizationCatalogService) isPlatformAdmin(principal *middleware.AuthenticatedPrincipal) bool {
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

func (s *OrganizationCatalogService) resolveActiveOrgUUID(principal *middleware.AuthenticatedPrincipal) (uuid.UUID, error) {
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

func (s *OrganizationCatalogService) ListCatalogItems(ctx context.Context, principal *middleware.AuthenticatedPrincipal, domainType string) ([]domain.OrganizationCatalogItem, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	return s.catalogRepo.ListCatalogItems(ctx, orgUUID, strings.ToUpper(strings.TrimSpace(domainType)))
}

func (s *OrganizationCatalogService) GetCatalogItemByID(ctx context.Context, principal *middleware.AuthenticatedPrincipal, itemID uuid.UUID) (*domain.OrganizationCatalogItem, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	return s.catalogRepo.GetCatalogItemByID(ctx, orgUUID, itemID)
}

func (s *OrganizationCatalogService) CreateCatalogItem(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	payload *domain.CreateCatalogItemPayload,
) (*domain.OrganizationCatalogItem, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	domainTypeUpper := strings.ToUpper(strings.TrimSpace(payload.DomainType))
	if !domain.IsValidCatalogDomainType(domainTypeUpper) {
		return nil, domain.ErrInvalidCatalogDomain
	}

	currVal := "NGN"
	if payload.Currency != nil && *payload.Currency != "" {
		currVal = *payload.Currency
	}

	itemEntity := &domain.OrganizationCatalogItem{
		OrganizationID:  orgUUID,
		MasterCatalogID: payload.MasterCatalogID,
		DomainType:      domainTypeUpper,
		Code:            strings.ToLower(strings.TrimSpace(payload.Code)),
		Name:            payload.Name,
		Description:     payload.Description,
		BasePrice:       payload.BasePrice,
		Currency:        currVal,
	}

	created, errCreate := s.catalogRepo.CreateCatalogItem(ctx, itemEntity, actorUUID)
	if errCreate != nil {
		return nil, errCreate
	}

	if s.auditRepo != nil {
		action := "CATALOG_ITEM_CREATED"
		resType := "organization.catalog_items"
		resID := created.ID.String()
		eventCat := "ORGANIZATION_COMMERCIAL_GOVERNANCE"
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

func (s *OrganizationCatalogService) UpdateCatalogItem(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	itemID uuid.UUID,
	payload *domain.UpdateCatalogItemPayload,
) (*domain.OrganizationCatalogItem, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	existing, errGet := s.catalogRepo.GetCatalogItemByID(ctx, orgUUID, itemID)
	if errGet != nil {
		return nil, errGet
	}

	if payload.Name != nil {
		existing.Name = *payload.Name
	}
	if payload.Description != nil {
		existing.Description = payload.Description
	}
	if payload.BasePrice != nil {
		existing.BasePrice = *payload.BasePrice
	}
	if payload.Currency != nil {
		existing.Currency = *payload.Currency
	}
	if payload.IsActive != nil {
		existing.IsActive = *payload.IsActive
	}
	existing.Version = payload.Version

	updated, errUp := s.catalogRepo.UpdateCatalogItem(ctx, existing, actorUUID)
	if errUp != nil {
		return nil, errUp
	}

	if s.auditRepo != nil {
		action := "CATALOG_ITEM_UPDATED"
		resType := "organization.catalog_items"
		resID := updated.ID.String()
		eventCat := "ORGANIZATION_COMMERCIAL_GOVERNANCE"
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

func (s *OrganizationCatalogService) SetBranchPriceOverride(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	itemID uuid.UUID,
	payload *domain.SetBranchPricePayload,
) (*domain.BranchPriceOverride, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	// Verify branch belongs to org
	_, errBranch := s.branchRepo.GetBranchByID(ctx, orgUUID, payload.FacilityBranchID)
	if errBranch != nil {
		return nil, errBranch
	}

	// Verify item belongs to org
	_, errItem := s.catalogRepo.GetCatalogItemByID(ctx, orgUUID, itemID)
	if errItem != nil {
		return nil, errItem
	}

	overrideEntity := &domain.BranchPriceOverride{
		OrganizationID:   orgUUID,
		FacilityBranchID: payload.FacilityBranchID,
		CatalogItemID:    itemID,
		OverridePrice:    payload.OverridePrice,
	}

	setOverride, errSet := s.catalogRepo.SetBranchPriceOverride(ctx, overrideEntity, actorUUID)
	if errSet != nil {
		return nil, errSet
	}

	if s.auditRepo != nil {
		action := "BRANCH_PRICE_OVERRIDDEN"
		resType := "organization.branch_price_overrides"
		resID := setOverride.ID.String()
		eventCat := "ORGANIZATION_COMMERCIAL_GOVERNANCE"
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

	return setOverride, nil
}

func (s *OrganizationCatalogService) CreateInsuranceProvider(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	payload *domain.CreateInsuranceProviderPayload,
) (*domain.InsuranceProvider, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	covVal := 100.00
	if payload.CoveragePercentage != nil {
		covVal = *payload.CoveragePercentage
	}

	providerEntity := &domain.InsuranceProvider{
		OrganizationID:     orgUUID,
		Name:               payload.Name,
		Code:               strings.ToLower(strings.TrimSpace(payload.Code)),
		CoveragePercentage: covVal,
	}

	created, errCreate := s.catalogRepo.CreateInsuranceProvider(ctx, providerEntity, actorUUID)
	if errCreate != nil {
		return nil, errCreate
	}

	if s.auditRepo != nil {
		action := "INSURANCE_PROVIDER_ADDED"
		resType := "organization.insurance_providers"
		resID := created.ID.String()
		eventCat := "ORGANIZATION_COMMERCIAL_GOVERNANCE"
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

func (s *OrganizationCatalogService) ListInsuranceProviders(ctx context.Context, principal *middleware.AuthenticatedPrincipal) ([]domain.InsuranceProvider, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	return s.catalogRepo.ListInsuranceProviders(ctx, orgUUID)
}
