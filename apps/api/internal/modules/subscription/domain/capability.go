package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Capability struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Code        string    `json:"code" db:"code"`
	Module      string    `json:"module" db:"module"`
	TierLevel   string    `json:"tierLevel" db:"tier_level"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	IsActive    bool      `json:"isActive" db:"is_active"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type OrganizationEntitlement struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	OrganizationID uuid.UUID  `json:"organizationId" db:"organization_id"`
	CapabilityID   uuid.UUID  `json:"capabilityId" db:"capability_id"`
	CapabilityCode string     `json:"capabilityCode" db:"capability_code"`
	Source         string     `json:"source" db:"source"` // plan, purchase, promotion, trial, platform, enterprise_contract
	Status         string     `json:"status" db:"status"` // active, suspended, revoked, expired
	StartsAt       time.Time  `json:"startsAt" db:"starts_at"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty" db:"expires_at"`
	GrantedBy      *uuid.UUID `json:"grantedBy,omitempty" db:"granted_by"`
	Metadata       string     `json:"metadata" db:"metadata"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

type EntitlementTraceStep struct {
	Step    string `json:"step"`
	Source  string `json:"source,omitempty"`
	Status  string `json:"status"`
	Details string `json:"details"`
}

type EntitlementTrace struct {
	OrganizationID      uuid.UUID              `json:"organizationId"`
	BasePlan            string                 `json:"basePlan"`
	RequestedCapability string                 `json:"requestedCapability"`
	Allowed             bool                   `json:"allowed"`
	Steps               []EntitlementTraceStep `json:"steps"`
}

type CapabilityPrice struct {
	ID            uuid.UUID `json:"id" db:"id"`
	CapabilityID  uuid.UUID `json:"capabilityId" db:"capability_id"`
	Currency      string    `json:"currency" db:"currency"`
	BillingPeriod string    `json:"billingPeriod" db:"billing_period"`
	Price         float64   `json:"price" db:"price"`
	IsActive      bool      `json:"isActive" db:"is_active"`
}

type CapabilitySubscription struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	OrganizationID     uuid.UUID  `json:"organizationId" db:"organization_id"`
	CapabilityID       uuid.UUID  `json:"capabilityId" db:"capability_id"`
	Status             string     `json:"status" db:"status"` // pending_payment, active, past_due, cancelled
	BillingCycle       string     `json:"billingCycle" db:"billing_cycle"`
	Price              float64    `json:"price" db:"price"`
	Currency           string     `json:"currency" db:"currency"`
	StartedAt          time.Time  `json:"startedAt" db:"started_at"`
	CurrentPeriodStart time.Time  `json:"currentPeriodStart" db:"current_period_start"`
	CurrentPeriodEnd   time.Time  `json:"currentPeriodEnd" db:"current_period_end"`
	CancelledAt        *time.Time `json:"cancelledAt,omitempty" db:"cancelled_at"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time  `json:"updatedAt" db:"updated_at"`
}

type CapabilityCatalogItem struct {
	Capability      Capability        `json:"capability"`
	Prices          []CapabilityPrice `json:"prices"`
	Dependencies    []string          `json:"dependencies"`
	AlreadyIncluded bool              `json:"alreadyIncluded"`
	IsEffective     bool              `json:"isEffective"`
}

type EntitlementRepository interface {
	GetPlanBaseCapabilities(ctx context.Context, planCode string) ([]string, error)
	GetOrganizationAddOnCapabilities(ctx context.Context, orgID uuid.UUID) ([]string, error)
	GetCapabilityDependencies(ctx context.Context, capabilityCodes []string) ([]string, error)
	GetOrganizationEntitlements(ctx context.Context, orgID uuid.UUID) ([]OrganizationEntitlement, error)
	GrantOrganizationEntitlement(ctx context.Context, entitlement *OrganizationEntitlement) error
	RevokeOrganizationEntitlement(ctx context.Context, orgID uuid.UUID, capabilityCode string) error
	GetCapabilityByCode(ctx context.Context, code string) (*Capability, error)
	GetAllCapabilities(ctx context.Context) ([]Capability, error)
	GetCapabilityPrices(ctx context.Context, capabilityID uuid.UUID) ([]CapabilityPrice, error)
	CreateCapabilitySubscription(ctx context.Context, sub *CapabilitySubscription) error
	GetCapabilitySubscription(ctx context.Context, subID uuid.UUID) (*CapabilitySubscription, error)
	UpdateCapabilitySubscriptionStatus(ctx context.Context, subID uuid.UUID, status string) error
}
