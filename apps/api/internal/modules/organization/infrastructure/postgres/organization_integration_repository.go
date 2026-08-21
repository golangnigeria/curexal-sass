package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrganizationIntegrationRepository struct {
	server *server.Server
}

func NewOrganizationIntegrationRepository(server *server.Server) *OrganizationIntegrationRepository {
	return &OrganizationIntegrationRepository{server: server}
}

func (r *OrganizationIntegrationRepository) CreateAPIKey(ctx context.Context, apiKey *domain.APIKey, keyHash string, actorID uuid.UUID) (*domain.APIKey, error) {
	dbExec := r.server.DB.Conn(ctx)
	if apiKey.ID == uuid.Nil {
		apiKey.ID = uuid.New()
	}

	scopesJSON := []byte("[]")
	if len(apiKey.Scopes) > 0 {
		scopesJSON = apiKey.Scopes
	}

	ipJSON := []byte("[]")
	if len(apiKey.IPWhitelist) > 0 {
		ipJSON = apiKey.IPWhitelist
	}

	rpm := 60
	if apiKey.RateLimitRPM > 0 {
		rpm = apiKey.RateLimitRPM
	}

	stmt := `
		INSERT INTO organization.api_keys (
			id, organization_id, name, key_prefix, key_hash, scopes, ip_whitelist, rate_limit_rpm, expires_at, is_active, version, created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7::jsonb, $8, $9, TRUE, 1, $10)
		RETURNING version, created_at, updated_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		apiKey.ID, apiKey.OrganizationID, apiKey.Name, apiKey.KeyPrefix, keyHash,
		string(scopesJSON), string(ipJSON), rpm, apiKey.ExpiresAt, actorID.String(),
	).Scan(&apiKey.Version, &apiKey.CreatedAt, &apiKey.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	apiKey.IsActive = true
	apiKey.RateLimitRPM = rpm
	apiKey.CreatedBy = &actorID
	return apiKey, nil
}

func (r *OrganizationIntegrationRepository) GetAPIKeyByID(ctx context.Context, orgID, keyID uuid.UUID) (*domain.APIKey, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, name, key_prefix, key_hash, scopes, ip_whitelist,
		       rate_limit_rpm, expires_at, last_used_at, is_active, version, created_at, updated_at, created_by
		FROM organization.api_keys
		WHERE organization_id = $1 AND id = $2
		LIMIT 1
	`

	var (
		k            domain.APIKey
		createdByStr *string
	)
	err := dbExec.QueryRow(ctx, stmt, orgID, keyID).Scan(
		&k.ID, &k.OrganizationID, &k.Name, &k.KeyPrefix, &k.KeyHash, &k.Scopes, &k.IPWhitelist,
		&k.RateLimitRPM, &k.ExpiresAt, &k.LastUsedAt, &k.IsActive, &k.Version, &k.CreatedAt, &k.UpdatedAt, &createdByStr,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAPIKeyNotFound
		}
		return nil, fmt.Errorf("failed to query API key: %w", err)
	}

	if createdByStr != nil && *createdByStr != "" {
		if parsed, pErr := uuid.Parse(*createdByStr); pErr == nil {
			k.CreatedBy = &parsed
		}
	}

	return &k, nil
}

func (r *OrganizationIntegrationRepository) GetAPIKeyByHash(ctx context.Context, keyHash string) (*domain.APIKey, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, name, key_prefix, key_hash, scopes, ip_whitelist,
		       rate_limit_rpm, expires_at, last_used_at, is_active, version, created_at, updated_at, created_by
		FROM organization.api_keys
		WHERE key_hash = $1
		LIMIT 1
	`

	var (
		k            domain.APIKey
		createdByStr *string
	)
	err := dbExec.QueryRow(ctx, stmt, keyHash).Scan(
		&k.ID, &k.OrganizationID, &k.Name, &k.KeyPrefix, &k.KeyHash, &k.Scopes, &k.IPWhitelist,
		&k.RateLimitRPM, &k.ExpiresAt, &k.LastUsedAt, &k.IsActive, &k.Version, &k.CreatedAt, &k.UpdatedAt, &createdByStr,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAPIKeyNotFound
		}
		return nil, fmt.Errorf("failed to query API key by hash: %w", err)
	}

	if createdByStr != nil && *createdByStr != "" {
		if parsed, pErr := uuid.Parse(*createdByStr); pErr == nil {
			k.CreatedBy = &parsed
		}
	}

	return &k, nil
}

func (r *OrganizationIntegrationRepository) ListAPIKeys(ctx context.Context, orgID uuid.UUID) ([]domain.APIKey, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, name, key_prefix, key_hash, scopes, ip_whitelist,
		       rate_limit_rpm, expires_at, last_used_at, is_active, version, created_at, updated_at, created_by
		FROM organization.api_keys
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := dbExec.Query(ctx, stmt, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}
	defer rows.Close()

	var list []domain.APIKey
	for rows.Next() {
		var (
			k            domain.APIKey
			createdByStr *string
		)
		err := rows.Scan(
			&k.ID, &k.OrganizationID, &k.Name, &k.KeyPrefix, &k.KeyHash, &k.Scopes, &k.IPWhitelist,
			&k.RateLimitRPM, &k.ExpiresAt, &k.LastUsedAt, &k.IsActive, &k.Version, &k.CreatedAt, &k.UpdatedAt, &createdByStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key row: %w", err)
		}
		if createdByStr != nil && *createdByStr != "" {
			if parsed, pErr := uuid.Parse(*createdByStr); pErr == nil {
				k.CreatedBy = &parsed
			}
		}
		list = append(list, k)
	}

	return list, nil
}

func (r *OrganizationIntegrationRepository) RevokeAPIKey(ctx context.Context, orgID, keyID uuid.UUID, actorID uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `UPDATE organization.api_keys SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP WHERE organization_id = $1 AND id = $2`
	res, err := dbExec.Exec(ctx, stmt, orgID, keyID)
	if err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}
	if res.RowsAffected() == 0 {
		return domain.ErrAPIKeyNotFound
	}
	return nil
}

func (r *OrganizationIntegrationRepository) CreateWebhookSubscription(ctx context.Context, sub *domain.WebhookSubscription, actorID uuid.UUID) (*domain.WebhookSubscription, error) {
	dbExec := r.server.DB.Conn(ctx)
	if sub.ID == uuid.Nil {
		sub.ID = uuid.New()
	}

	eventsJSON := []byte("[]")
	if len(sub.EventTypes) > 0 {
		eventsJSON = sub.EventTypes
	}

	stmt := `
		INSERT INTO organization.webhook_subscriptions (
			id, organization_id, name, target_url, event_types, encrypted_secret, is_active, version, created_by
		)
		VALUES ($1, $2, $3, $4, $5::jsonb, $6, TRUE, 1, $7)
		RETURNING version, created_at, updated_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		sub.ID, sub.OrganizationID, sub.Name, sub.TargetURL, string(eventsJSON), sub.SigningSecret, actorID.String(),
	).Scan(&sub.Version, &sub.CreatedAt, &sub.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "uk_org_webhook_target" {
			return nil, domain.ErrDuplicateWebhookTarget
		}
		return nil, fmt.Errorf("failed to create webhook subscription: %w", err)
	}

	sub.IsActive = true
	sub.CreatedBy = &actorID
	return sub, nil
}

func (r *OrganizationIntegrationRepository) GetWebhookSubscriptionByID(ctx context.Context, orgID, subID uuid.UUID) (*domain.WebhookSubscription, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, name, target_url, event_types, encrypted_secret,
		       is_active, version, created_at, updated_at, created_by
		FROM organization.webhook_subscriptions
		WHERE organization_id = $1 AND id = $2
		LIMIT 1
	`

	var (
		s            domain.WebhookSubscription
		createdByStr *string
	)
	err := dbExec.QueryRow(ctx, stmt, orgID, subID).Scan(
		&s.ID, &s.OrganizationID, &s.Name, &s.TargetURL, &s.EventTypes, &s.SigningSecret,
		&s.IsActive, &s.Version, &s.CreatedAt, &s.UpdatedAt, &createdByStr,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrWebhookSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to query webhook subscription: %w", err)
	}

	if createdByStr != nil && *createdByStr != "" {
		if parsed, pErr := uuid.Parse(*createdByStr); pErr == nil {
			s.CreatedBy = &parsed
		}
	}

	return &s, nil
}

func (r *OrganizationIntegrationRepository) ListWebhookSubscriptions(ctx context.Context, orgID uuid.UUID) ([]domain.WebhookSubscription, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, name, target_url, event_types, encrypted_secret,
		       is_active, version, created_at, updated_at, created_by
		FROM organization.webhook_subscriptions
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := dbExec.Query(ctx, stmt, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhook subscriptions: %w", err)
	}
	defer rows.Close()

	var list []domain.WebhookSubscription
	for rows.Next() {
		var (
			s            domain.WebhookSubscription
			createdByStr *string
		)
		err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.Name, &s.TargetURL, &s.EventTypes, &s.SigningSecret,
			&s.IsActive, &s.Version, &s.CreatedAt, &s.UpdatedAt, &createdByStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan webhook subscription row: %w", err)
		}
		if createdByStr != nil && *createdByStr != "" {
			if parsed, pErr := uuid.Parse(*createdByStr); pErr == nil {
				s.CreatedBy = &parsed
			}
		}
		list = append(list, s)
	}

	return list, nil
}

func (r *OrganizationIntegrationRepository) DeleteWebhookSubscription(ctx context.Context, orgID, subID uuid.UUID, actorID uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `DELETE FROM organization.webhook_subscriptions WHERE organization_id = $1 AND id = $2`
	res, err := dbExec.Exec(ctx, stmt, orgID, subID)
	if err != nil {
		return fmt.Errorf("failed to delete webhook subscription: %w", err)
	}
	if res.RowsAffected() == 0 {
		return domain.ErrWebhookSubscriptionNotFound
	}
	return nil
}

func (r *OrganizationIntegrationRepository) CreateWebhookDelivery(ctx context.Context, delivery *domain.WebhookDelivery) (*domain.WebhookDelivery, error) {
	dbExec := r.server.DB.Conn(ctx)
	if delivery.ID == uuid.Nil {
		delivery.ID = uuid.New()
	}

	payloadJSON := []byte("{}")
	if len(delivery.Payload) > 0 {
		payloadJSON = delivery.Payload
	}

	stmt := `
		INSERT INTO organization.webhook_deliveries (
			id, organization_id, subscription_id, event_type, payload, response_status, response_body, attempt_count, last_error, status
		)
		VALUES ($1, $2, $3, $4, $5::jsonb, $6, $7, $8, $9, $10)
		RETURNING created_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		delivery.ID, delivery.OrganizationID, delivery.SubscriptionID, delivery.EventType, string(payloadJSON),
		delivery.ResponseStatus, delivery.ResponseBody, delivery.AttemptCount, delivery.LastError, delivery.Status,
	).Scan(&delivery.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create webhook delivery record: %w", err)
	}

	return delivery, nil
}

func (r *OrganizationIntegrationRepository) ListWebhookDeliveries(ctx context.Context, orgID uuid.UUID, limit int) ([]domain.WebhookDelivery, error) {
	dbExec := r.server.DB.Conn(ctx)
	if limit <= 0 {
		limit = 50
	}

	stmt := `
		SELECT id, organization_id, subscription_id, event_type, payload, response_status, response_body, attempt_count, last_error, status, created_at
		FROM organization.webhook_deliveries
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := dbExec.Query(ctx, stmt, orgID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhook deliveries: %w", err)
	}
	defer rows.Close()

	var list []domain.WebhookDelivery
	for rows.Next() {
		var d domain.WebhookDelivery
		err := rows.Scan(
			&d.ID, &d.OrganizationID, &d.SubscriptionID, &d.EventType, &d.Payload,
			&d.ResponseStatus, &d.ResponseBody, &d.AttemptCount, &d.LastError, &d.Status, &d.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan webhook delivery row: %w", err)
		}
		list = append(list, d)
	}

	return list, nil
}
