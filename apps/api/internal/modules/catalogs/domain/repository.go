package domain

import (
	"context"

	"github.com/google/uuid"
)

type CatalogRepository interface {
	ListItems(ctx context.Context, domain CatalogDomain, category string, activeOnly bool) ([]CatalogItem, error)
	SearchItems(ctx context.Context, domain CatalogDomain, query string) ([]CatalogItem, error)
	GetItemByCode(ctx context.Context, domain CatalogDomain, code string) (*CatalogItem, error)
	GetItemByID(ctx context.Context, domain CatalogDomain, id uuid.UUID) (*CatalogItem, error)
	CreateItem(ctx context.Context, item *CatalogItem, updatedBy uuid.UUID) (*CatalogItem, error)
	UpdateItem(ctx context.Context, item *CatalogItem, updatedBy uuid.UUID) (*CatalogItem, error)

	// Baseline legacy compatibility methods
	GetSpecimenTypes(ctx context.Context) ([]SpecimenType, error)
	GetTestCategories(ctx context.Context) ([]TestCategory, error)
	GetUnitsOfMeasure(ctx context.Context) ([]UnitOfMeasure, error)
	GetSpecialties(ctx context.Context) ([]MedicalSpecialty, error)
	SearchICD10(ctx context.Context, query string) ([]ICD10Code, error)
}
