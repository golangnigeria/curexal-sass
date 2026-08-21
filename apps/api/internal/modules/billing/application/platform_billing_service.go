package application

import (
	"context"
	"fmt"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/billing/domain"
	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
)

type PlatformBillingRepo interface {
	ListPricingRules(ctx context.Context) ([]domain.PricingRule, error)
	UpdatePricingRule(ctx context.Context, rule *domain.PricingRule, updatedBy uuid.UUID) (*domain.PricingRule, error)
	ListPaymentGateways(ctx context.Context) ([]domain.PaymentGatewayConfig, error)
	GetPaymentGatewayByProvider(ctx context.Context, providerCode string) (*domain.PaymentGatewayConfig, error)
	UpdatePaymentGateway(ctx context.Context, gateway *domain.PaymentGatewayConfig, updatedBy uuid.UUID) (*domain.PaymentGatewayConfig, error)
}

type PlatformPricingService struct {
	repo      PlatformBillingRepo
	auditRepo auditDomain.AuditRepository
}

func NewPlatformPricingService(repo PlatformBillingRepo, auditRepo auditDomain.AuditRepository) *PlatformPricingService {
	return &PlatformPricingService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func isPlatformAdmin(principal *middleware.AuthenticatedPrincipal) bool {
	if principal == nil {
		return false
	}
	if principal.Platform.IsPlatformAdmin || principal.Platform.IsSuperAdmin || principal.Platform.IsPlatformStaff {
		return true
	}
	if principal.Role == "super_admin" || principal.Role == "platform_admin" || principal.Role == "platform_staff" {
		return true
	}
	return false
}

func (s *PlatformPricingService) ListPricingRules(ctx context.Context) ([]domain.PricingRule, error) {
	return s.repo.ListPricingRules(ctx)
}

func (s *PlatformPricingService) UpdatePricingRule(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	rule *domain.PricingRule,
) (*domain.PricingRule, error) {
	if !isPlatformAdmin(principal) {
		return nil, domain.ErrUnauthorizedPlatformAdmin
	}
	if err := rule.Validate(); err != nil {
		return nil, err
	}

	actorUUID, err := uuid.Parse(principal.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", err)
	}

	updated, err := s.repo.UpdatePricingRule(ctx, rule, actorUUID)
	if err != nil {
		return nil, err
	}

	if s.auditRepo != nil {
		action := "PRICING_RULE_UPDATED"
		resType := "platform.pricing_rules"
		resID := updated.ID.String()
		eventCat := "COMMERCIAL"
		severity := "HIGH"
		status := "SUCCESS"

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    true,
			ActorID:       &principal.UserID,
			ActorName:     &principal.Identity.FullName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	return updated, nil
}

type PaymentGatewayVaultService struct {
	repo      PlatformBillingRepo
	auditRepo auditDomain.AuditRepository
	masterKey string
}

func NewPaymentGatewayVaultService(repo PlatformBillingRepo, auditRepo auditDomain.AuditRepository, masterKey string) *PaymentGatewayVaultService {
	if masterKey == "" {
		masterKey = "CUREXAL_DEFAULT_PLATFORM_AES_MASTER_KEY_2026"
	}
	return &PaymentGatewayVaultService{
		repo:      repo,
		auditRepo: auditRepo,
		masterKey: masterKey,
	}
}

func (s *PaymentGatewayVaultService) ListGateways(ctx context.Context) ([]domain.PaymentGatewayConfig, error) {
	gateways, err := s.repo.ListPaymentGateways(ctx)
	if err != nil {
		return nil, err
	}

	for i := range gateways {
		gateways[i].RedactedSecretKey = crypto.RedactSecret("MASKED_KEY")
	}

	return gateways, nil
}

func (s *PaymentGatewayVaultService) GetGateway(ctx context.Context, providerCode string) (*domain.PaymentGatewayConfig, error) {
	gateway, err := s.repo.GetPaymentGatewayByProvider(ctx, providerCode)
	if err != nil {
		return nil, err
	}

	gateway.RedactedSecretKey = crypto.RedactSecret("MASKED_KEY")
	return gateway, nil
}

type UpdateGatewayPayload struct {
	Name                string   `json:"name"`
	IsEnabled           bool     `json:"isEnabled"`
	Priority            int      `json:"priority"`
	SupportedCurrencies []string `json:"supportedCurrencies"`
	PlaintextSecretKey  string   `json:"secretKey"` // Plaintext input from admin for update
	PublicKey           *string  `json:"publicKey"`
	WebhookSecret       *string  `json:"webhookSecret"`
	Version             int      `json:"version"`
}

func (s *PaymentGatewayVaultService) UpdateGateway(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	providerCode string,
	payload *UpdateGatewayPayload,
) (*domain.PaymentGatewayConfig, error) {
	if !isPlatformAdmin(principal) {
		return nil, domain.ErrUnauthorizedPlatformAdmin
	}

	existing, err := s.repo.GetPaymentGatewayByProvider(ctx, providerCode)
	if err != nil {
		return nil, err
	}

	ciphertext := existing.EncryptedSecretKey
	if payload.PlaintextSecretKey != "" {
		encrypted, errEnc := crypto.EncryptAEAD(payload.PlaintextSecretKey, s.masterKey)
		if errEnc != nil {
			return nil, fmt.Errorf("failed to encrypt gateway secret key: %w", errEnc)
		}
		ciphertext = encrypted
	}

	existing.Name = payload.Name
	existing.IsEnabled = payload.IsEnabled
	existing.Priority = payload.Priority
	if len(payload.SupportedCurrencies) > 0 {
		existing.SupportedCurrencies = payload.SupportedCurrencies
	}
	existing.EncryptedSecretKey = ciphertext
	if payload.PublicKey != nil {
		existing.PublicKey = payload.PublicKey
	}
	if payload.WebhookSecret != nil {
		existing.WebhookSecret = payload.WebhookSecret
	}
	existing.Version = payload.Version

	if err := existing.Validate(); err != nil {
		return nil, err
	}

	actorUUID, err := uuid.Parse(principal.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", err)
	}

	updated, err := s.repo.UpdatePaymentGateway(ctx, existing, actorUUID)
	if err != nil {
		return nil, err
	}

	if s.auditRepo != nil {
		action := "PAYMENT_GATEWAY_UPDATED"
		resType := "platform.payment_gateways"
		resID := updated.ProviderCode
		eventCat := "COMMERCIAL_SECURITY"
		severity := "CRITICAL"
		status := "SUCCESS"

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    true,
			ActorID:       &principal.UserID,
			ActorName:     &principal.Identity.FullName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	updated.RedactedSecretKey = crypto.RedactSecret("MASKED_KEY")
	return updated, nil
}
