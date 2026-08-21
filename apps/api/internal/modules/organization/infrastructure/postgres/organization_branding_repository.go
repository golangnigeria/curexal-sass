package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrganizationBrandingRepository struct {
	server *server.Server
}

func NewOrganizationBrandingRepository(server *server.Server) *OrganizationBrandingRepository {
	return &OrganizationBrandingRepository{server: server}
}

func (r *OrganizationBrandingRepository) GetBranding(ctx context.Context, orgID uuid.UUID) (*domain.BrandingConfig, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, logo_url, COALESCE(primary_color, '#0F172A'), theme_branding,
		       custom_domain, COALESCE(custom_domain_status, 'PENDING'), version, updated_at, updated_by
		FROM organization.organizations
		WHERE id = $1
		LIMIT 1
	`

	var (
		b            domain.BrandingConfig
		updatedByStr *string
	)
	err := dbExec.QueryRow(ctx, stmt, orgID.String()).Scan(
		&b.OrganizationID, &b.LogoURL, &b.PrimaryColor, &b.ThemeBranding,
		&b.CustomDomain, &b.CustomDomainStatus, &b.Version, &b.UpdatedAt, &updatedByStr,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrganizationNotFound
		}
		return nil, fmt.Errorf("failed to query organization branding orgID=%s: %w", orgID, err)
	}

	if updatedByStr != nil && *updatedByStr != "" {
		if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
			b.UpdatedBy = &parsed
		}
	}

	return &b, nil
}

func (r *OrganizationBrandingRepository) UpdateBranding(ctx context.Context, orgID uuid.UUID, payload *domain.UpdateBrandingPayload, actorID uuid.UUID) (*domain.BrandingConfig, error) {
	dbExec := r.server.DB.Conn(ctx)

	themeJSON := []byte("{}")
	if len(payload.ThemeBranding) > 0 {
		themeJSON = payload.ThemeBranding
	}

	stmt := `
		UPDATE organization.organizations
		SET logo_url = COALESCE($1, logo_url),
		    primary_color = COALESCE($2, primary_color),
		    theme_branding = $3::jsonb,
		    custom_domain = $4,
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
		payload.LogoURL, payload.PrimaryColor, string(themeJSON), payload.CustomDomain, actorID.String(), orgID.String(), payload.Version,
	).Scan(&newVersion, &updatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "uk_org_custom_domain" {
			return nil, domain.ErrDuplicateCustomDomain
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOptimisticLockingConflict
		}
		return nil, fmt.Errorf("failed to update organization branding: %w", err)
	}

	return r.GetBranding(ctx, orgID)
}

func (r *OrganizationBrandingRepository) SaveNotificationConfig(ctx context.Context, config *domain.NotificationConfig, actorID uuid.UUID) (*domain.NotificationConfig, error) {
	dbExec := r.server.DB.Conn(ctx)
	if config.ID == uuid.Nil {
		config.ID = uuid.New()
	}

	metaJSON := []byte("{}")
	if len(config.ConfigMetadata) > 0 {
		metaJSON = config.ConfigMetadata
	}

	stmt := `
		INSERT INTO organization.notification_configs (
			id, organization_id, channel, provider, sender_email, sender_name,
			host, port, username, encrypted_password, encrypted_api_key, config_metadata, is_active, version, updated_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12::jsonb, COALESCE($13, TRUE), 1, $14)
		ON CONFLICT (organization_id, channel, provider)
		DO UPDATE SET
			sender_email = EXCLUDED.sender_email,
			sender_name = EXCLUDED.sender_name,
			host = EXCLUDED.host,
			port = EXCLUDED.port,
			username = EXCLUDED.username,
			encrypted_password = COALESCE(EXCLUDED.encrypted_password, organization.notification_configs.encrypted_password),
			encrypted_api_key = COALESCE(EXCLUDED.encrypted_api_key, organization.notification_configs.encrypted_api_key),
			config_metadata = EXCLUDED.config_metadata,
			is_active = EXCLUDED.is_active,
			version = organization.notification_configs.version + 1,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = EXCLUDED.updated_by
		RETURNING id, version, created_at, updated_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		config.ID, config.OrganizationID, config.Channel, config.Provider, config.SenderEmail, config.SenderName,
		config.Host, config.Port, config.Username, config.Password, config.APIKey, string(metaJSON), config.IsActive, actorID.String(),
	).Scan(&config.ID, &config.Version, &config.CreatedAt, &config.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to save notification config: %w", err)
	}

	return config, nil
}

func (r *OrganizationBrandingRepository) ListNotificationConfigs(ctx context.Context, orgID uuid.UUID) ([]domain.NotificationConfig, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, channel, provider, sender_email, sender_name,
		       host, port, username, encrypted_password, encrypted_api_key, config_metadata, is_active, version, created_at, updated_at
		FROM organization.notification_configs
		WHERE organization_id = $1
		ORDER BY channel ASC, provider ASC
	`

	rows, err := dbExec.Query(ctx, stmt, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification configs: %w", err)
	}
	defer rows.Close()

	var configs []domain.NotificationConfig
	for rows.Next() {
		var c domain.NotificationConfig
		err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.Channel, &c.Provider, &c.SenderEmail, &c.SenderName,
			&c.Host, &c.Port, &c.Username, &c.Password, &c.APIKey, &c.ConfigMetadata, &c.IsActive, &c.Version, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification config row: %w", err)
		}
		configs = append(configs, c)
	}

	return configs, nil
}

func (r *OrganizationBrandingRepository) GetNotificationConfig(ctx context.Context, orgID uuid.UUID, channel, provider string) (*domain.NotificationConfig, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, channel, provider, sender_email, sender_name,
		       host, port, username, encrypted_password, encrypted_api_key, config_metadata, is_active, version, created_at, updated_at
		FROM organization.notification_configs
		WHERE organization_id = $1 AND channel = $2 AND provider = $3
		LIMIT 1
	`

	var c domain.NotificationConfig
	err := dbExec.QueryRow(ctx, stmt, orgID, channel, provider).Scan(
		&c.ID, &c.OrganizationID, &c.Channel, &c.Provider, &c.SenderEmail, &c.SenderName,
		&c.Host, &c.Port, &c.Username, &c.Password, &c.APIKey, &c.ConfigMetadata, &c.IsActive, &c.Version, &c.CreatedAt, &c.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotificationConfigNotFound
		}
		return nil, fmt.Errorf("failed to query notification config: %w", err)
	}

	return &c, nil
}

func (r *OrganizationBrandingRepository) SaveNotificationTemplate(ctx context.Context, template *domain.NotificationTemplate, actorID uuid.UUID) (*domain.NotificationTemplate, error) {
	dbExec := r.server.DB.Conn(ctx)
	if template.ID == uuid.Nil {
		template.ID = uuid.New()
	}

	varsJSON := []byte("[]")
	if len(template.AllowedVariables) > 0 {
		varsJSON = template.AllowedVariables
	}

	stmt := `
		INSERT INTO organization.notification_templates (
			id, organization_id, template_key, channel, subject, body_template, allowed_variables, is_active, version, updated_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, COALESCE($8, TRUE), 1, $9)
		ON CONFLICT (organization_id, template_key, channel)
		DO UPDATE SET
			subject = EXCLUDED.subject,
			body_template = EXCLUDED.body_template,
			allowed_variables = EXCLUDED.allowed_variables,
			is_active = EXCLUDED.is_active,
			version = organization.notification_templates.version + 1,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = EXCLUDED.updated_by
		RETURNING id, version, created_at, updated_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		template.ID, template.OrganizationID, template.TemplateKey, template.Channel, template.Subject,
		template.BodyTemplate, string(varsJSON), template.IsActive, actorID.String(),
	).Scan(&template.ID, &template.Version, &template.CreatedAt, &template.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to save notification template: %w", err)
	}

	return template, nil
}

func (r *OrganizationBrandingRepository) ListNotificationTemplates(ctx context.Context, orgID uuid.UUID) ([]domain.NotificationTemplate, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, template_key, channel, subject, body_template, allowed_variables, is_active, version, created_at, updated_at
		FROM organization.notification_templates
		WHERE organization_id = $1
		ORDER BY template_key ASC, channel ASC
	`

	rows, err := dbExec.Query(ctx, stmt, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification templates: %w", err)
	}
	defer rows.Close()

	var list []domain.NotificationTemplate
	for rows.Next() {
		var t domain.NotificationTemplate
		err := rows.Scan(
			&t.ID, &t.OrganizationID, &t.TemplateKey, &t.Channel, &t.Subject,
			&t.BodyTemplate, &t.AllowedVariables, &t.IsActive, &t.Version, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification template row: %w", err)
		}
		list = append(list, t)
	}

	return list, nil
}

func (r *OrganizationBrandingRepository) GetNotificationTemplate(ctx context.Context, orgID uuid.UUID, templateKey, channel string) (*domain.NotificationTemplate, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, template_key, channel, subject, body_template, allowed_variables, is_active, version, created_at, updated_at
		FROM organization.notification_templates
		WHERE organization_id = $1 AND template_key = $2 AND channel = $3
		LIMIT 1
	`

	var t domain.NotificationTemplate
	err := dbExec.QueryRow(ctx, stmt, orgID, templateKey, channel).Scan(
		&t.ID, &t.OrganizationID, &t.TemplateKey, &t.Channel, &t.Subject,
		&t.BodyTemplate, &t.AllowedVariables, &t.IsActive, &t.Version, &t.CreatedAt, &t.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotificationTemplateNotFound
		}
		return nil, fmt.Errorf("failed to query notification template: %w", err)
	}

	return &t, nil
}

func (r *OrganizationBrandingRepository) CreateUserNotification(ctx context.Context, notif *domain.UserNotification) (*domain.UserNotification, error) {
	dbExec := r.server.DB.Conn(ctx)
	if notif.ID == uuid.Nil {
		notif.ID = uuid.New()
	}

	dataJSON := []byte("{}")
	if len(notif.Data) > 0 {
		dataJSON = notif.Data
	}

	stmt := `
		INSERT INTO organization.user_notifications (id, organization_id, user_id, notification_type, title, body, data, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8)
		RETURNING created_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		notif.ID, notif.OrganizationID, notif.UserID, notif.NotificationType, notif.Title, notif.Body, string(dataJSON), notif.ExpiresAt,
	).Scan(&notif.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user in-app notification: %w", err)
	}

	return notif, nil
}

func (r *OrganizationBrandingRepository) ListUserNotifications(ctx context.Context, orgID, userID uuid.UUID, limit int) ([]domain.UserNotification, error) {
	dbExec := r.server.DB.Conn(ctx)
	if limit <= 0 {
		limit = 50
	}

	stmt := `
		SELECT id, organization_id, user_id, notification_type, title, body, data, read_at, created_at, expires_at
		FROM organization.user_notifications
		WHERE organization_id = $1 AND user_id = $2
		ORDER BY created_at DESC
		LIMIT $3
	`

	rows, err := dbExec.Query(ctx, stmt, orgID, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list user notifications: %w", err)
	}
	defer rows.Close()

	var list []domain.UserNotification
	for rows.Next() {
		var n domain.UserNotification
		err := rows.Scan(
			&n.ID, &n.OrganizationID, &n.UserID, &n.NotificationType, &n.Title, &n.Body,
			&n.Data, &n.ReadAt, &n.CreatedAt, &n.ExpiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user notification row: %w", err)
		}
		list = append(list, n)
	}

	return list, nil
}

func (r *OrganizationBrandingRepository) MarkNotificationRead(ctx context.Context, orgID, userID, notifID uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `UPDATE organization.user_notifications SET read_at = CURRENT_TIMESTAMP WHERE id = $1 AND organization_id = $2 AND user_id = $3`
	res, err := dbExec.Exec(ctx, stmt, notifID, orgID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark notification read: %w", err)
	}
	if res.RowsAffected() == 0 {
		return domain.ErrUserNotificationNotFound
	}
	return nil
}

func (r *OrganizationBrandingRepository) MarkAllNotificationsRead(ctx context.Context, orgID, userID uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `UPDATE organization.user_notifications SET read_at = CURRENT_TIMESTAMP WHERE organization_id = $1 AND user_id = $2 AND read_at IS NULL`
	_, err := dbExec.Exec(ctx, stmt, orgID, userID)
	return err
}

func (r *OrganizationBrandingRepository) CreateNotificationDelivery(ctx context.Context, delivery *domain.NotificationDelivery) (*domain.NotificationDelivery, error) {
	dbExec := r.server.DB.Conn(ctx)
	if delivery.ID == uuid.Nil {
		delivery.ID = uuid.New()
	}

	stmt := `
		INSERT INTO organization.notification_deliveries (
			id, organization_id, notification_id, channel, provider, recipient, template_key,
			status, provider_message_id, attempt_count, last_error, queued_at, sent_at, delivered_at, failed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING created_at, updated_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		delivery.ID, delivery.OrganizationID, delivery.NotificationID, delivery.Channel, delivery.Provider,
		delivery.Recipient, delivery.TemplateKey, delivery.Status, delivery.ProviderMessageID, delivery.AttemptCount,
		delivery.LastError, delivery.QueuedAt, delivery.SentAt, delivery.DeliveredAt, delivery.FailedAt,
	).Scan(&delivery.CreatedAt, &delivery.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create notification delivery record: %w", err)
	}

	return delivery, nil
}

func (r *OrganizationBrandingRepository) ListNotificationDeliveries(ctx context.Context, orgID uuid.UUID, limit int) ([]domain.NotificationDelivery, error) {
	dbExec := r.server.DB.Conn(ctx)
	if limit <= 0 {
		limit = 50
	}

	stmt := `
		SELECT id, organization_id, notification_id, channel, provider, recipient, template_key,
		       status, provider_message_id, attempt_count, last_error, queued_at, sent_at, delivered_at, failed_at, created_at, updated_at
		FROM organization.notification_deliveries
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := dbExec.Query(ctx, stmt, orgID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification deliveries: %w", err)
	}
	defer rows.Close()

	var list []domain.NotificationDelivery
	for rows.Next() {
		var d domain.NotificationDelivery
		err := rows.Scan(
			&d.ID, &d.OrganizationID, &d.NotificationID, &d.Channel, &d.Provider, &d.Recipient, &d.TemplateKey,
			&d.Status, &d.ProviderMessageID, &d.AttemptCount, &d.LastError, &d.QueuedAt, &d.SentAt, &d.DeliveredAt, &d.FailedAt, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification delivery row: %w", err)
		}
		list = append(list, d)
	}

	return list, nil
}
