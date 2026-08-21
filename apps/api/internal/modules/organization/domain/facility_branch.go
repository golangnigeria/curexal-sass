package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrFacilityBranchNotFound  = errors.New("facility branch not found")
	ErrDuplicateBranchCode     = errors.New("a branch with this code already exists in your organization")
	ErrInvalidFacilityType     = errors.New("invalid facility type specification")
	ErrInactiveFacilityType    = errors.New("selected platform facility type is currently inactive")
	ErrHeadquartersConflict    = errors.New("organization already has an active headquarters facility")
	ErrMaxBranchesExceeded     = errors.New("organization branch limit reached for active subscription plan")
	ErrInvalidFacilityBranch   = errors.New("invalid facility branch parameters: name, code, and facility type are required")
)

type FacilityBranch struct {
	ID                   uuid.UUID       `json:"id"`
	OrganizationID       uuid.UUID       `json:"organizationId"`
	FacilityTypeID       uuid.UUID       `json:"facilityTypeId"`
	FacilityTypeCode     string          `json:"facilityTypeCode,omitempty"`
	FacilityTypeName     string          `json:"facilityTypeName,omitempty"`
	FacilityTypeCategory string          `json:"facilityTypeCategory,omitempty"`
	Code                 string          `json:"code"`
	Name                 string          `json:"name"`
	IsHeadquarters       bool            `json:"isHeadquarters"`
	Email                *string         `json:"email,omitempty"`
	Phone                *string         `json:"phone,omitempty"`
	Address              *string         `json:"address,omitempty"`
	City                 *string         `json:"city,omitempty"`
	State                *string         `json:"state,omitempty"`
	LGA                  *string         `json:"lga,omitempty"`
	Country              string          `json:"country"`
	OperatingHours       json.RawMessage `json:"operatingHours,omitempty"`
	Status               string          `json:"status"` // 'ACTIVE', 'INACTIVE', 'SUSPENDED'
	Version              int             `json:"version"`
	CreatedAt            time.Time       `json:"createdAt"`
	UpdatedAt            time.Time       `json:"updatedAt"`
	UpdatedBy            *uuid.UUID      `json:"updatedBy,omitempty"`
}

func (b *FacilityBranch) Validate() error {
	if b.Name == "" || b.Code == "" || b.FacilityTypeID == uuid.Nil {
		return ErrInvalidFacilityBranch
	}
	b.Code = strings.ToLower(strings.TrimSpace(b.Code))
	return nil
}

type CreateFacilityBranchPayload struct {
	FacilityTypeID uuid.UUID       `json:"facilityTypeId"`
	Code           string          `json:"code"`
	Name           string          `json:"name"`
	IsHeadquarters bool            `json:"isHeadquarters"`
	Email          *string         `json:"email"`
	Phone          *string         `json:"phone"`
	Address        *string         `json:"address"`
	City           *string         `json:"city"`
	State          *string         `json:"state"`
	LGA            *string         `json:"lga"`
	Country        *string         `json:"country"`
	OperatingHours json.RawMessage `json:"operatingHours"`
}

type UpdateFacilityBranchPayload struct {
	Name           *string         `json:"name"`
	IsHeadquarters *bool           `json:"isHeadquarters"`
	Email          *string         `json:"email"`
	Phone          *string         `json:"phone"`
	Address        *string         `json:"address"`
	City           *string         `json:"city"`
	State          *string         `json:"state"`
	LGA            *string         `json:"lga"`
	OperatingHours json.RawMessage `json:"operatingHours"`
	Status         *string         `json:"status"`
	Version        int             `json:"version"`
}
