package application

import (
	"context"
	"fmt"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/catalogs/domain"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
)

type MasterCatalogService struct {
	repo      domain.CatalogRepository
	auditRepo auditDomain.AuditRepository
}

func NewMasterCatalogService(repo domain.CatalogRepository, auditRepo auditDomain.AuditRepository) *MasterCatalogService {
	return &MasterCatalogService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (s *MasterCatalogService) IsPlatformAdmin(principal *middleware.AuthenticatedPrincipal) bool {
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

func (s *MasterCatalogService) ListItems(ctx context.Context, catalogDomain domain.CatalogDomain, category string, activeOnly bool) ([]domain.CatalogItem, error) {
	return s.repo.ListItems(ctx, catalogDomain, category, activeOnly)
}

func (s *MasterCatalogService) SearchItems(ctx context.Context, catalogDomain domain.CatalogDomain, query string) ([]domain.CatalogItem, error) {
	return s.repo.SearchItems(ctx, catalogDomain, query)
}

func (s *MasterCatalogService) GetItemByCode(ctx context.Context, catalogDomain domain.CatalogDomain, code string) (*domain.CatalogItem, error) {
	return s.repo.GetItemByCode(ctx, catalogDomain, code)
}

func (s *MasterCatalogService) GetItemByID(ctx context.Context, catalogDomain domain.CatalogDomain, id uuid.UUID) (*domain.CatalogItem, error) {
	return s.repo.GetItemByID(ctx, catalogDomain, id)
}

func (s *MasterCatalogService) CreateItem(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	item *domain.CatalogItem,
) (*domain.CatalogItem, error) {
	if !s.IsPlatformAdmin(principal) {
		return nil, domain.ErrUnauthorizedPlatformAdmin
	}
	if err := item.Validate(); err != nil {
		return nil, err
	}

	actorUUID, err := uuid.Parse(principal.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", err)
	}

	created, err := s.repo.CreateItem(ctx, item, actorUUID)
	if err != nil {
		return nil, err
	}

	if s.auditRepo != nil {
		action := "MASTER_CATALOG_ITEM_CREATED"
		resType := fmt.Sprintf("platform.%s_catalogs", item.Domain)
		resID := created.ID.String()
		eventCat := "ADMINISTRATIVE"
		severity := "HIGH"
		status := "SUCCESS"

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    true,
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

func (s *MasterCatalogService) UpdateItem(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	itemID uuid.UUID,
	catDomain domain.CatalogDomain,
	payload *domain.UpdateCatalogItemPayload,
) (*domain.CatalogItem, error) {
	if !s.IsPlatformAdmin(principal) {
		return nil, domain.ErrUnauthorizedPlatformAdmin
	}

	existing, err := s.repo.GetItemByID(ctx, catDomain, itemID)
	if err != nil {
		return nil, err
	}

	if payload.Category != nil {
		existing.Category = *payload.Category
	}
	if payload.Code != nil {
		existing.Code = *payload.Code
	}
	if payload.Name != nil {
		existing.Name = *payload.Name
	}
	if payload.Description != nil {
		existing.Description = *payload.Description
	}
	if payload.SystemGroup != nil {
		existing.SystemGroup = *payload.SystemGroup
	}
	if payload.BasePrice != nil {
		existing.BasePrice = *payload.BasePrice
	}
	if len(payload.Metadata) > 0 {
		existing.Metadata = payload.Metadata
	}
	if payload.IsActive != nil {
		existing.IsActive = *payload.IsActive
	}
	if payload.Version > 0 {
		existing.Version = payload.Version
	}

	if err := existing.Validate(); err != nil {
		return nil, err
	}

	actorUUID, err := uuid.Parse(principal.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", err)
	}

	updated, err := s.repo.UpdateItem(ctx, existing, actorUUID)
	if err != nil {
		return nil, err
	}

	if s.auditRepo != nil {
		action := "MASTER_CATALOG_ITEM_UPDATED"
		resType := fmt.Sprintf("platform.%s_catalogs", existing.Domain)
		resID := updated.ID.String()
		eventCat := "ADMINISTRATIVE"
		severity := "HIGH"
		status := "SUCCESS"

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    true,
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

// Baseline DTO compatibility methods
func (s *MasterCatalogService) GetSpecimenTypes(ctx context.Context) ([]domain.SpecimenType, error) {
	return s.repo.GetSpecimenTypes(ctx)
}

func (s *MasterCatalogService) GetTestCategories(ctx context.Context) ([]domain.TestCategory, error) {
	return s.repo.GetTestCategories(ctx)
}

func (s *MasterCatalogService) GetUnitsOfMeasure(ctx context.Context) ([]domain.UnitOfMeasure, error) {
	return s.repo.GetUnitsOfMeasure(ctx)
}

func (s *MasterCatalogService) GetSpecialties(ctx context.Context) ([]domain.MedicalSpecialty, error) {
	return s.repo.GetSpecialties(ctx)
}

func (s *MasterCatalogService) SearchICD10(ctx context.Context, query string) ([]domain.ICD10Code, error) {
	return s.repo.SearchICD10(ctx, query)
}
