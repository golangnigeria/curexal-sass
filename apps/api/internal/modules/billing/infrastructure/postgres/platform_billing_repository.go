package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/billing/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PlatformBillingRepository struct {
	server *server.Server
}

func NewPlatformBillingRepository(server *server.Server) *PlatformBillingRepository {
	return &PlatformBillingRepository{server: server}
}

func (r *PlatformBillingRepository) ListPricingRules(ctx context.Context) ([]domain.PricingRule, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, target_type, target_code, currency, monthly_price, annual_price, vat_percentage, is_active, version, created_at, updated_at, updated_by
		FROM platform.pricing_rules
		ORDER BY target_type ASC, target_code ASC, currency ASC
	`
	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query platform pricing rules: %w", err)
	}
	defer rows.Close()

	var list []domain.PricingRule
	for rows.Next() {
		var (
			p            domain.PricingRule
			updatedByStr *string
		)
		err := rows.Scan(
			&p.ID, &p.TargetType, &p.TargetCode, &p.Currency, &p.MonthlyPrice, &p.AnnualPrice, &p.VATPercentage,
			&p.IsActive, &p.Version, &p.CreatedAt, &p.UpdatedAt, &updatedByStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pricing rule row: %w", err)
		}
		if updatedByStr != nil && *updatedByStr != "" {
			if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
				p.UpdatedBy = &parsed
			}
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PlatformBillingRepository) UpdatePricingRule(ctx context.Context, rule *domain.PricingRule, updatedBy uuid.UUID) (*domain.PricingRule, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE platform.pricing_rules
		SET monthly_price = $1,
		    annual_price = $2,
		    vat_percentage = $3,
		    is_active = $4,
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $5
		WHERE id = $6 AND version = $7
		RETURNING version, updated_at
	`
	var (
		newVersion int
		updatedAt  time.Time
	)
	err := dbExec.QueryRow(ctx, stmt,
		rule.MonthlyPrice, rule.AnnualPrice, rule.VATPercentage, rule.IsActive, updatedBy.String(), rule.ID.String(), rule.Version,
	).Scan(&newVersion, &updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOptimisticLockingConflict
		}
		return nil, fmt.Errorf("failed to update pricing rule: %w", err)
	}

	rule.Version = newVersion
	rule.UpdatedAt = updatedAt
	rule.UpdatedBy = &updatedBy

	return rule, nil
}

func (r *PlatformBillingRepository) ListPaymentGateways(ctx context.Context) ([]domain.PaymentGatewayConfig, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, provider_code, name, is_enabled, priority, supported_currencies, encrypted_secret_key, public_key, webhook_secret, version, created_at, updated_at, updated_by
		FROM platform.payment_gateways
		ORDER BY priority ASC
	`
	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query platform payment gateways: %w", err)
	}
	defer rows.Close()

	var list []domain.PaymentGatewayConfig
	for rows.Next() {
		var (
			g                     domain.PaymentGatewayConfig
			rawSupportedCurrencies []byte
			updatedByStr          *string
		)
		err := rows.Scan(
			&g.ID, &g.ProviderCode, &g.Name, &g.IsEnabled, &g.Priority, &rawSupportedCurrencies,
			&g.EncryptedSecretKey, &g.PublicKey, &g.WebhookSecret, &g.Version, &g.CreatedAt, &g.UpdatedAt, &updatedByStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment gateway row: %w", err)
		}
		if len(rawSupportedCurrencies) > 0 {
			_ = json.Unmarshal(rawSupportedCurrencies, &g.SupportedCurrencies)
		}
		if updatedByStr != nil && *updatedByStr != "" {
			if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
				g.UpdatedBy = &parsed
			}
		}
		list = append(list, g)
	}
	return list, nil
}

func (r *PlatformBillingRepository) GetPaymentGatewayByProvider(ctx context.Context, providerCode string) (*domain.PaymentGatewayConfig, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, provider_code, name, is_enabled, priority, supported_currencies, encrypted_secret_key, public_key, webhook_secret, version, created_at, updated_at, updated_by
		FROM platform.payment_gateways
		WHERE provider_code = $1
		LIMIT 1
	`
	var (
		g                     domain.PaymentGatewayConfig
		rawSupportedCurrencies []byte
		updatedByStr          *string
	)
	err := dbExec.QueryRow(ctx, stmt, providerCode).Scan(
		&g.ID, &g.ProviderCode, &g.Name, &g.IsEnabled, &g.Priority, &rawSupportedCurrencies,
		&g.EncryptedSecretKey, &g.PublicKey, &g.WebhookSecret, &g.Version, &g.CreatedAt, &g.UpdatedAt, &updatedByStr,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPaymentGatewayNotFound
		}
		return nil, fmt.Errorf("failed to query payment gateway provider=%s: %w", providerCode, err)
	}
	if len(rawSupportedCurrencies) > 0 {
		_ = json.Unmarshal(rawSupportedCurrencies, &g.SupportedCurrencies)
	}
	if updatedByStr != nil && *updatedByStr != "" {
		if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
			g.UpdatedBy = &parsed
		}
	}
	return &g, nil
}

func (r *PlatformBillingRepository) UpdatePaymentGateway(ctx context.Context, gateway *domain.PaymentGatewayConfig, updatedBy uuid.UUID) (*domain.PaymentGatewayConfig, error) {
	dbExec := r.server.DB.Conn(ctx)

	currenciesJSON, err := json.Marshal(gateway.SupportedCurrencies)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gateway supported currencies: %w", err)
	}

	stmt := `
		UPDATE platform.payment_gateways
		SET name = $1,
		    is_enabled = $2,
		    priority = $3,
		    supported_currencies = $4,
		    encrypted_secret_key = $5,
		    public_key = $6,
		    webhook_secret = $7,
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $8
		WHERE provider_code = $9 AND version = $10
		RETURNING id, version, updated_at
	`
	var (
		id        uuid.UUID
		newVersion int
		updatedAt time.Time
	)
	err = dbExec.QueryRow(ctx, stmt,
		gateway.Name, gateway.IsEnabled, gateway.Priority, currenciesJSON,
		gateway.EncryptedSecretKey, gateway.PublicKey, gateway.WebhookSecret,
		updatedBy.String(), gateway.ProviderCode, gateway.Version,
	).Scan(&id, &newVersion, &updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOptimisticLockingConflict
		}
		return nil, fmt.Errorf("failed to update payment gateway: %w", err)
	}

	gateway.ID = id
	gateway.Version = newVersion
	gateway.UpdatedAt = updatedAt
	gateway.UpdatedBy = &updatedBy

	return gateway, nil
}
