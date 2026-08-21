package application

import (
	"context"

	"github.com/golangnigeria/curexal/internal/modules/clinical/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
)

type ClinicalApplicationService struct {
	server      *server.Server
	catalogRepo domain.CatalogRepository
}

func NewClinicalApplicationService(server *server.Server, catalogRepo domain.CatalogRepository) *ClinicalApplicationService {
	return &ClinicalApplicationService{
		server:      server,
		catalogRepo: catalogRepo,
	}
}

func (s *ClinicalApplicationService) ListCatalog(ctx context.Context) ([]domain.CatalogItem, error) {
	return s.catalogRepo.ListCatalog(ctx)
}

func (s *ClinicalApplicationService) GetCatalogItemByID(ctx context.Context, id uuid.UUID) (*domain.CatalogItem, error) {
	return s.catalogRepo.GetCatalogItemByID(ctx, id)
}

func (s *ClinicalApplicationService) CreateCatalogItem(ctx context.Context, payload *domain.CreateCatalogItemPayload) (*domain.CatalogItem, error) {
	return s.catalogRepo.CreateCatalogItem(ctx, payload)
}

func (s *ClinicalApplicationService) UpdateCatalogMetadata(ctx context.Context, id uuid.UUID, payload *domain.CreateCatalogItemPayload) (*domain.CatalogItem, error) {
	return s.catalogRepo.UpdateCatalogMetadata(ctx, id, payload)
}

func (s *ClinicalApplicationService) UpdateCatalogPricing(ctx context.Context, id uuid.UUID, payload *domain.UpdatePricingPayload) (*domain.CatalogItem, error) {
	return s.catalogRepo.UpdateCatalogPricing(ctx, id, payload)
}

func (s *ClinicalApplicationService) ImportCatalog(ctx context.Context, items []domain.CreateCatalogItemPayload) (int, error) {
	return s.catalogRepo.ImportCatalog(ctx, items)
}
