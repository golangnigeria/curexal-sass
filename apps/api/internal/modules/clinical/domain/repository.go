package domain

import (
	"context"

	"github.com/google/uuid"
)

type CreateCatalogItemPayload struct {
	Name                 string      `json:"name"`
	Code                 string      `json:"code"`
	BasePrice            float64     `json:"basePrice"`
	Type                 string      `json:"type"`
	UrgencyPrice         float64     `json:"urgencyPrice"`
	CommissionValue      float64     `json:"commissionValue"`
	CommissionPercentage float64     `json:"commissionPercentage"`
	DiscountAmount       float64     `json:"discountAmount"`
	DiscountPercentage   float64     `json:"discountPercentage"`
	DisplayName          *string     `json:"displayName"`
	ShortName            *string     `json:"shortName"`
	RecoveryTime         *string     `json:"recoveryTime"`
	DepartmentID         *uuid.UUID  `json:"departmentId"`
	TestGroup            *string     `json:"testGroup"`
	TestCategory         *string     `json:"testCategory"`
	TatHours             int         `json:"tatHours"`
	ChildServiceIDs      []uuid.UUID `json:"childServiceIds"`
}

type UpdatePricingPayload struct {
	BasePrice            *float64 `json:"basePrice"`
	UrgencyPrice         *float64 `json:"urgencyPrice"`
	CommissionValue      *float64 `json:"commissionValue"`
	CommissionPercentage *float64 `json:"commissionPercentage"`
	DiscountAmount       *float64 `json:"discountAmount"`
	DiscountPercentage   *float64 `json:"discountPercentage"`
}

type CatalogRepository interface {
	ListCatalog(ctx context.Context) ([]CatalogItem, error)
	GetCatalogItemByID(ctx context.Context, id uuid.UUID) (*CatalogItem, error)
	CreateCatalogItem(ctx context.Context, payload *CreateCatalogItemPayload) (*CatalogItem, error)
	UpdateCatalogMetadata(ctx context.Context, id uuid.UUID, payload *CreateCatalogItemPayload) (*CatalogItem, error)
	UpdateCatalogPricing(ctx context.Context, id uuid.UUID, payload *UpdatePricingPayload) (*CatalogItem, error)
	ImportCatalog(ctx context.Context, items []CreateCatalogItemPayload) (int, error)
}
