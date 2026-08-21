package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCatalogItemNotFound       = errors.New("organization catalog item not found")
	ErrDuplicateCatalogCode      = errors.New("a catalog item with this code already exists in your organization domain")
	ErrInvalidCatalogDomain      = errors.New("invalid catalog domain type: must be CLINICAL, LABORATORY, RADIOLOGY, PHARMACY, or PROCEDURE")
	ErrBranchPriceNotFound       = errors.New("branch price override record not found")
	ErrInsuranceProviderNotFound = errors.New("insurance provider record not found")
	ErrDuplicateInsuranceCode    = errors.New("an insurance provider with this code already exists in your organization")
)

var AllowedCatalogDomainTypes = map[string]bool{
	"CLINICAL":   true,
	"LABORATORY": true,
	"RADIOLOGY":  true,
	"PHARMACY":   true,
	"PROCEDURE":  true,
}

func IsValidCatalogDomainType(domainType string) bool {
	return AllowedCatalogDomainTypes[strings.ToUpper(strings.TrimSpace(domainType))]
}

type OrganizationCatalogItem struct {
	ID              uuid.UUID  `json:"id"`
	OrganizationID  uuid.UUID  `json:"organizationId"`
	MasterCatalogID *uuid.UUID `json:"masterCatalogId,omitempty"`
	DomainType      string     `json:"domainType"`
	Code            string     `json:"code"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	BasePrice       float64    `json:"basePrice"`
	Currency        string     `json:"currency"`
	IsActive        bool       `json:"isActive"`
	Version         int        `json:"version"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	UpdatedBy       *uuid.UUID `json:"updatedBy,omitempty"`
}

type BranchPriceOverride struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   uuid.UUID  `json:"organizationId"`
	FacilityBranchID uuid.UUID  `json:"facilityBranchId"`
	CatalogItemID    uuid.UUID  `json:"catalogItemId"`
	OverridePrice    float64    `json:"overridePrice"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	UpdatedBy        *uuid.UUID `json:"updatedBy,omitempty"`
}

type InsuranceProvider struct {
	ID                 uuid.UUID  `json:"id"`
	OrganizationID     uuid.UUID  `json:"organizationId"`
	Name               string     `json:"name"`
	Code               string     `json:"code"`
	CoveragePercentage float64    `json:"coveragePercentage"`
	IsActive           bool       `json:"isActive"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	UpdatedBy          *uuid.UUID `json:"updatedBy,omitempty"`
}

type CreateCatalogItemPayload struct {
	MasterCatalogID *uuid.UUID `json:"masterCatalogId"`
	DomainType      string     `json:"domainType"`
	Code            string     `json:"code"`
	Name            string     `json:"name"`
	Description     *string    `json:"description"`
	BasePrice       float64    `json:"basePrice"`
	Currency        *string    `json:"currency"`
}

type UpdateCatalogItemPayload struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	BasePrice   *float64 `json:"basePrice"`
	Currency    *string  `json:"currency"`
	IsActive    *bool    `json:"isActive"`
	Version     int      `json:"version"`
}

type SetBranchPricePayload struct {
	FacilityBranchID uuid.UUID `json:"facilityBranchId"`
	OverridePrice    float64   `json:"overridePrice"`
}

type CreateInsuranceProviderPayload struct {
	Name               string   `json:"name"`
	Code               string   `json:"code"`
	CoveragePercentage *float64 `json:"coveragePercentage"`
}
