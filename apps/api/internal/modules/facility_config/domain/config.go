package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrFacilityTypeNotFound      = errors.New("facility type not found")
	ErrInvalidFacilityType       = errors.New("invalid facility type data")
	ErrOptimisticLockingConflict = errors.New("optimistic locking conflict: record was modified concurrently")
)

type FacilityTier string

const (
	ClinicalTier   FacilityTier = "clinical"
	DiagnosticTier FacilityTier = "diagnostic"
	RetailTier     FacilityTier = "retail"
	ResearchTier   FacilityTier = "research"
)

type FacilityTypeEntity struct {
	ID                  uuid.UUID            `json:"id"`
	Code                string               `json:"code"`
	Name                string               `json:"name"`
	Category            string               `json:"category"`
	IconKey             string               `json:"iconKey"`
	Description         string               `json:"description"`
	IsActive            bool                 `json:"isActive"`
	DefaultCapabilities []FacilityCapability `json:"defaultCapabilities,omitempty"`
	Version             int                  `json:"version"`
	CreatedAt           time.Time            `json:"createdAt"`
	UpdatedAt           time.Time            `json:"updatedAt"`
	UpdatedBy           *uuid.UUID           `json:"updatedBy,omitempty"`
}

type FacilityCapability struct {
	ID             uuid.UUID `json:"id"`
	FacilityTypeID uuid.UUID `json:"facilityTypeId"`
	CapabilityID   uuid.UUID `json:"capabilityId"`
	CapabilityCode string    `json:"capabilityCode"`
	CapabilityName string    `json:"capabilityName"`
	IsDefault      bool      `json:"isDefault"`
}

type FacilityType struct {
	ID           string          `json:"id" db:"id"`
	CategoryID   string          `json:"categoryId" db:"category_id"`
	Name         string          `json:"name" db:"name"`
	Code         string          `json:"code" db:"code"`
	DisplayName  string          `json:"displayName" db:"display_name"`
	Description  string          `json:"description" db:"description"`
	IconKey      string          `json:"iconKey" db:"icon_key"`
	DefaultTheme json.RawMessage `json:"defaultTheme" db:"default_theme"`
	IsActive     bool            `json:"isActive" db:"is_active"`
}

type FacilityTypeDTO struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	DisplayName string `json:"displayName"`
	Category    string `json:"category"`
	IconKey     string `json:"iconKey"`
	Description string `json:"description"`
}

type RegistrationFormDTO struct {
	FacilityTypeID    string          `json:"facilityTypeId"`
	Version           int             `json:"version"`
	Sections          json.RawMessage `json:"sections"`
	RequiredDocuments json.RawMessage `json:"requiredDocuments,omitempty"`
}

type NavigationMenuDTO struct {
	FacilityTypeID string          `json:"facilityTypeId"`
	MenuItems      json.RawMessage `json:"menuItems"`
}

type SetupStepDTO struct {
	StepNumber  int             `json:"stepNumber"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	FieldSchema json.RawMessage `json:"fieldSchema"`
}

type DashboardDTO struct {
	FacilityTypeID string          `json:"facilityTypeId"`
	Widgets        json.RawMessage `json:"widgets"`
}

func (f *FacilityTypeEntity) Validate() error {
	if f.Code == "" || f.Name == "" || f.Category == "" {
		return ErrInvalidFacilityType
	}
	return nil
}
