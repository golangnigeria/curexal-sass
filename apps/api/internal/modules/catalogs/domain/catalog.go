package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCatalogItemNotFound       = errors.New("master reference catalog item not found")
	ErrInvalidCatalogDomain      = errors.New("invalid master catalog domain (expected 'clinical', 'lab', 'radiology', or 'pharmacy')")
	ErrInvalidCatalogItem        = errors.New("invalid catalog item parameters: code, name, and category are required")
	ErrOptimisticLockingConflict = errors.New("optimistic locking conflict: record was modified concurrently")
	ErrUnauthorizedPlatformAdmin = errors.New("unauthorized: administrative operation requires platform admin privileges")
)

type CatalogDomain string

const (
	ClinicalDomain  CatalogDomain = "clinical"
	LabDomain       CatalogDomain = "lab"
	RadiologyDomain CatalogDomain = "radiology"
	PharmacyDomain  CatalogDomain = "pharmacy"
)

type CatalogItem struct {
	ID          uuid.UUID       `json:"id"`
	Domain      CatalogDomain   `json:"domain"`
	Category    string          `json:"category"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	SystemGroup string          `json:"systemGroup,omitempty"`
	BasePrice   float64         `json:"basePrice"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	IsActive    bool            `json:"isActive"`
	Version     int             `json:"version"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	UpdatedBy   *uuid.UUID      `json:"updatedBy,omitempty"`
}

type UpdateCatalogItemPayload struct {
	Category    *string          `json:"category,omitempty"`
	Code        *string          `json:"code,omitempty"`
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	SystemGroup *string          `json:"systemGroup,omitempty"`
	BasePrice   *float64         `json:"basePrice,omitempty"`
	Metadata    json.RawMessage  `json:"metadata,omitempty"`
	IsActive    *bool            `json:"isActive,omitempty"`
	Version     int              `json:"version,omitempty"`
}

func (c *CatalogItem) Validate() error {
	if c.Code == "" || c.Name == "" || c.Category == "" {
		return ErrInvalidCatalogItem
	}
	if c.Domain != ClinicalDomain && c.Domain != LabDomain && c.Domain != RadiologyDomain && c.Domain != PharmacyDomain {
		return ErrInvalidCatalogDomain
	}
	return nil
}

// Baseline DTO compatibility models
type SpecimenType struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Container string `json:"container"`
	IsActive  bool   `json:"isActive"`
}

type TestCategory struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

type UnitOfMeasure struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	IsActive bool   `json:"isActive"`
}

type MedicalSpecialty struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ICD10Code struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type CatalogDataResponse struct {
	SpecimenTypes  []SpecimenType     `json:"specimenTypes"`
	TestCategories []TestCategory     `json:"testCategories"`
	UnitsOfMeasure []UnitOfMeasure    `json:"unitsOfMeasure"`
	Specialties    []MedicalSpecialty `json:"specialties"`
}
