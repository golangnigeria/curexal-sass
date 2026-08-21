package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPricingRuleNotFound       = errors.New("pricing rule not found")
	ErrInvalidPricingRule        = errors.New("invalid pricing rule configuration")
	ErrPaymentGatewayNotFound    = errors.New("payment gateway configuration not found")
	ErrInvalidGatewayConfig      = errors.New("invalid payment gateway configuration")
	ErrOptimisticLockingConflict = errors.New("optimistic concurrency conflict: record has been modified by another request")
	ErrUnauthorizedPlatformAdmin = errors.New("unauthorized: operation requires platform administrator privileges")
)

type PricingRule struct {
	ID            uuid.UUID  `json:"id"`
	TargetType    string     `json:"targetType"` // 'plan' or 'capability'
	TargetCode    string     `json:"targetCode"`
	Currency      string     `json:"currency"`
	MonthlyPrice  float64    `json:"monthlyPrice"`
	AnnualPrice   float64    `json:"annualPrice"`
	VATPercentage float64    `json:"vatPercentage"`
	IsActive      bool       `json:"isActive"`
	Version       int        `json:"version"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	UpdatedBy     *uuid.UUID `json:"updatedBy,omitempty"`
}

type PaymentGatewayConfig struct {
	ID                  uuid.UUID  `json:"id"`
	ProviderCode        string     `json:"providerCode"` // 'paystack', 'flutterwave', 'stripe', 'mock'
	Name                string     `json:"name"`
	IsEnabled           bool       `json:"isEnabled"`
	Priority            int        `json:"priority"`
	SupportedCurrencies []string   `json:"supportedCurrencies"`
	EncryptedSecretKey  string     `json:"-"`                   // Ciphertext only; omitted from JSON
	RedactedSecretKey   string     `json:"secretKey,omitempty"` // Exposed in DTO as '••••••••'
	PublicKey           *string    `json:"publicKey,omitempty"`
	WebhookSecret       *string    `json:"webhookSecret,omitempty"`
	Version             int        `json:"version"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
	UpdatedBy           *uuid.UUID `json:"updatedBy,omitempty"`
}

func (p *PricingRule) Validate() error {
	if p.TargetType == "" || p.TargetCode == "" || p.Currency == "" {
		return ErrInvalidPricingRule
	}
	if p.MonthlyPrice < 0 || p.AnnualPrice < 0 || p.VATPercentage < 0 {
		return ErrInvalidPricingRule
	}
	return nil
}

func (g *PaymentGatewayConfig) Validate() error {
	if g.ProviderCode == "" || g.Name == "" {
		return ErrInvalidGatewayConfig
	}
	return nil
}
