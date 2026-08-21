package domain

import (
	"time"

	"github.com/google/uuid"
)

type CatalogItem struct {
	ID                   uuid.UUID   `json:"id"`
	Name                 string      `json:"name"`
	Code                 string      `json:"code"`
	BasePrice            float64     `json:"basePrice"`
	Type                 string      `json:"type"` // "test" or "profile"
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
	ChildServiceIDs      []uuid.UUID `json:"childServiceIds,omitempty"`
	CreatedAt            time.Time   `json:"createdAt"`
	UpdatedAt            time.Time   `json:"updatedAt"`
}
