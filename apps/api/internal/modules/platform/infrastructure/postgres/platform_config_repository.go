package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/platform/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PlatformConfigRepository struct {
	server *server.Server
}

func NewPlatformConfigRepository(server *server.Server) *PlatformConfigRepository {
	return &PlatformConfigRepository{server: server}
}

func (r *PlatformConfigRepository) GetGeneralSettings(ctx context.Context) (*domain.PlatformGeneralSettings, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, platform_name, platform_logo_url, platform_favicon_url, platform_description,
		       support_email, support_phone, support_website, default_country, default_currency,
		       supported_countries, supported_currencies, default_timezone, default_locale,
		       date_format, time_format, number_format, measurement_units, maintenance_mode,
		       announcement_banner, status, version, created_at, updated_at, updated_by
		FROM platform.general_settings
		ORDER BY created_at ASC LIMIT 1
	`
	row := dbExec.QueryRow(ctx, stmt)

	var (
		s                      domain.PlatformGeneralSettings
		rawSupportedCountries  []byte
		rawSupportedCurrencies []byte
		updatedByStr           *string
	)

	err := row.Scan(
		&s.ID, &s.PlatformName, &s.PlatformLogoURL, &s.PlatformFaviconURL, &s.PlatformDescription,
		&s.SupportEmail, &s.SupportPhone, &s.SupportWebsite, &s.DefaultCountry, &s.DefaultCurrency,
		&rawSupportedCountries, &rawSupportedCurrencies, &s.DefaultTimezone, &s.DefaultLocale,
		&s.DateFormat, &s.TimeFormat, &s.NumberFormat, &s.MeasurementUnits, &s.MaintenanceMode,
		&s.AnnouncementBanner, &s.Status, &s.Version, &s.CreatedAt, &s.UpdatedAt, &updatedByStr,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("platform general settings row not initialized: %w", err)
		}
		return nil, fmt.Errorf("failed to scan platform general settings: %w", err)
	}

	if len(rawSupportedCountries) > 0 {
		_ = json.Unmarshal(rawSupportedCountries, &s.SupportedCountries)
	}
	if len(rawSupportedCurrencies) > 0 {
		_ = json.Unmarshal(rawSupportedCurrencies, &s.SupportedCurrencies)
	}
	if updatedByStr != nil && *updatedByStr != "" {
		if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
			s.UpdatedBy = &parsed
		}
	}

	return &s, nil
}

func (r *PlatformConfigRepository) UpdateGeneralSettings(ctx context.Context, settings *domain.PlatformGeneralSettings, updatedBy uuid.UUID) (*domain.PlatformGeneralSettings, error) {
	dbExec := r.server.DB.Conn(ctx)

	countriesJSON, err := json.Marshal(settings.SupportedCountries)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal supported countries: %w", err)
	}
	currenciesJSON, err := json.Marshal(settings.SupportedCurrencies)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal supported currencies: %w", err)
	}

	stmt := `
		UPDATE platform.general_settings
		SET platform_name = $1,
		    platform_logo_url = $2,
		    platform_favicon_url = $3,
		    platform_description = $4,
		    support_email = $5,
		    support_phone = $6,
		    support_website = $7,
		    default_country = $8,
		    default_currency = $9,
		    supported_countries = $10,
		    supported_currencies = $11,
		    default_timezone = $12,
		    default_locale = $13,
		    date_format = $14,
		    time_format = $15,
		    number_format = $16,
		    measurement_units = $17,
		    maintenance_mode = $18,
		    announcement_banner = $19,
		    status = $20,
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $21
		WHERE id = $22 AND version = $23
		RETURNING version, updated_at
	`

	var (
		newVersion int
		updatedAt  time.Time
	)
	err = dbExec.QueryRow(ctx, stmt,
		settings.PlatformName, settings.PlatformLogoURL, settings.PlatformFaviconURL, settings.PlatformDescription,
		settings.SupportEmail, settings.SupportPhone, settings.SupportWebsite, settings.DefaultCountry, settings.DefaultCurrency,
		countriesJSON, currenciesJSON, settings.DefaultTimezone, settings.DefaultLocale,
		settings.DateFormat, settings.TimeFormat, settings.NumberFormat, settings.MeasurementUnits, settings.MaintenanceMode,
		settings.AnnouncementBanner, settings.Status, updatedBy.String(), settings.ID.String(), settings.Version,
	).Scan(&newVersion, &updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOptimisticLockingConflict
		}
		return nil, fmt.Errorf("failed to update platform general settings: %w", err)
	}

	settings.Version = newVersion
	settings.UpdatedAt = updatedAt
	settings.UpdatedBy = &updatedBy

	return settings, nil
}

func (r *PlatformConfigRepository) GetSecurityPolicy(ctx context.Context) (*domain.IdentitySecurityPolicy, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, min_password_length, password_require_uppercase, password_require_number,
		       password_require_symbol, password_expiration_days, email_verification_required,
		       mfa_policy, max_failed_login_attempts, account_lockout_duration_minutes,
		       session_max_duration_hours, refresh_token_duration_days, max_active_sessions,
		       version, created_at, updated_at, updated_by
		FROM platform.identity_security_policies
		ORDER BY created_at ASC LIMIT 1
	`
	row := dbExec.QueryRow(ctx, stmt)

	var (
		p            domain.IdentitySecurityPolicy
		updatedByStr *string
	)

	err := row.Scan(
		&p.ID, &p.MinPasswordLength, &p.PasswordRequireUppercase, &p.PasswordRequireNumber,
		&p.PasswordRequireSymbol, &p.PasswordExpirationDays, &p.EmailVerificationRequired,
		&p.MFAPolicy, &p.MaxFailedLoginAttempts, &p.AccountLockoutDurationMinutes,
		&p.SessionMaxDurationHours, &p.RefreshTokenDurationDays, &p.MaxActiveSessions,
		&p.Version, &p.CreatedAt, &p.UpdatedAt, &updatedByStr,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("platform identity security policy row not initialized: %w", err)
		}
		return nil, fmt.Errorf("failed to scan platform identity security policy: %w", err)
	}

	if updatedByStr != nil && *updatedByStr != "" {
		if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
			p.UpdatedBy = &parsed
		}
	}

	return &p, nil
}

func (r *PlatformConfigRepository) UpdateSecurityPolicy(ctx context.Context, policy *domain.IdentitySecurityPolicy, updatedBy uuid.UUID) (*domain.IdentitySecurityPolicy, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE platform.identity_security_policies
		SET min_password_length = $1,
		    password_require_uppercase = $2,
		    password_require_number = $3,
		    password_require_symbol = $4,
		    password_expiration_days = $5,
		    email_verification_required = $6,
		    mfa_policy = $7,
		    max_failed_login_attempts = $8,
		    account_lockout_duration_minutes = $9,
		    session_max_duration_hours = $10,
		    refresh_token_duration_days = $11,
		    max_active_sessions = $12,
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $13
		WHERE id = $14 AND version = $15
		RETURNING version, updated_at
	`

	var (
		newVersion int
		updatedAt  time.Time
	)

	err := dbExec.QueryRow(ctx, stmt,
		policy.MinPasswordLength, policy.PasswordRequireUppercase, policy.PasswordRequireNumber,
		policy.PasswordRequireSymbol, policy.PasswordExpirationDays, policy.EmailVerificationRequired,
		policy.MFAPolicy, policy.MaxFailedLoginAttempts, policy.AccountLockoutDurationMinutes,
		policy.SessionMaxDurationHours, policy.RefreshTokenDurationDays, policy.MaxActiveSessions,
		updatedBy.String(), policy.ID.String(), policy.Version,
	).Scan(&newVersion, &updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOptimisticLockingConflict
		}
		return nil, fmt.Errorf("failed to update platform identity security policy: %w", err)
	}

	policy.Version = newVersion
	policy.UpdatedAt = updatedAt
	policy.UpdatedBy = &updatedBy

	return policy, nil
}
