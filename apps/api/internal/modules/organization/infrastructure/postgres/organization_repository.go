package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OrganizationRepository struct {
	server *server.Server
}

func NewOrganizationRepository(server *server.Server) *OrganizationRepository {
	return &OrganizationRepository{server: server}
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, name, slug, plan, status, logo_url, custom_domain,
		       registration_number, license_number, tax_id, email, phone, address, city, state, lga, country,
		       setup_state, setup_step, completed_at, settings, version, created_at, updated_at, updated_by
		FROM organization.organizations
		WHERE id = $1
	`

	var (
		org          domain.Organization
		updatedByStr *string
		countryStr   *string
	)
	err := dbExec.QueryRow(ctx, stmt, id.String()).Scan(
		&org.ID, &org.Name, &org.Slug, &org.Plan, &org.Status, &org.LogoURL, &org.CustomDomain,
		&org.RegistrationNumber, &org.LicenseNumber, &org.TaxID, &org.Email, &org.Phone, &org.Address, &org.City, &org.State, &org.LGA, &countryStr,
		&org.SetupState, &org.SetupStep, &org.CompletedAt, &org.Settings, &org.Version, &org.CreatedAt, &org.UpdatedAt, &updatedByStr,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrganizationNotFound
		}
		return nil, fmt.Errorf("failed to execute get organization query for id=%s: %w", id, err)
	}

	if countryStr != nil {
		org.Country = *countryStr
	} else {
		org.Country = "Nigeria"
	}
	if updatedByStr != nil && *updatedByStr != "" {
		if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
			org.UpdatedBy = &parsed
		}
	}

	return &org, nil
}

func (r *OrganizationRepository) Create(ctx context.Context, name, slug, plan string) (*domain.Organization, error) {
	dbExec := r.server.DB.Conn(ctx)

	slugToUse := strings.ToLower(strings.TrimSpace(slug))
	if slugToUse == "" {
		slugToUse = strings.ToLower(strings.TrimSpace(name))
		slugToUse = strings.ReplaceAll(slugToUse, " ", "-")
	}

	planToUse := "smart"
	if plan != "" {
		planToUse = plan
	}

	orgID := uuid.New()
	stmt := `
		INSERT INTO organization.organizations (id, name, slug, plan, status, setup_state, setup_step, version)
		VALUES ($1, $2, $3, $4, 'pending_verification', 'PENDING_REGISTRATION', 1, 1)
		RETURNING created_at, updated_at
	`

	var createdAt, updatedAt time.Time
	err := dbExec.QueryRow(ctx, stmt, orgID, name, slugToUse, planToUse).Scan(&createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create organization query for name=%s: %w", name, err)
	}

	return r.GetByID(ctx, orgID)
}

func (r *OrganizationRepository) Update(ctx context.Context, id uuid.UUID, name, slug, plan, customDomain *string, settings map[string]any) (*domain.Organization, error) {
	dbExec := r.server.DB.Conn(ctx)

	query := "UPDATE organization.organizations SET "
	args := pgx.NamedArgs{"id": id.String()}
	updates := []string{}

	if name != nil {
		updates = append(updates, "name = @name")
		args["name"] = *name
	}
	if slug != nil {
		cleanSlug := strings.ToLower(strings.TrimSpace(*slug))
		updates = append(updates, "slug = @slug")
		args["slug"] = cleanSlug
	}
	if plan != nil {
		updates = append(updates, "plan = @plan")
		args["plan"] = *plan
	}
	if customDomain != nil {
		updates = append(updates, "custom_domain = @custom_domain")
		args["custom_domain"] = customDomain
	}
	if settings != nil {
		updates = append(updates, "settings = @settings")
		args["settings"] = settings
	}

	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}

	updates = append(updates, "version = version + 1", "updated_at = CURRENT_TIMESTAMP")
	query += strings.Join(updates, ", ") + " WHERE id = @id"

	_, err := dbExec.Exec(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to update organization id=%s: %w", id, err)
	}

	return r.GetByID(ctx, id)
}

func (r *OrganizationRepository) UpdateProfile(ctx context.Context, orgID uuid.UUID, payload *domain.UpdateOrganizationProfilePayload, actorID uuid.UUID) (*domain.Organization, error) {
	dbExec := r.server.DB.Conn(ctx)

	var updatedByVal *string
	if actorID != uuid.Nil {
		s := actorID.String()
		updatedByVal = &s
	}

	regNumber := payload.RegistrationNumber
	if regNumber == nil && payload.LegalName != nil && *payload.LegalName != "" {
		regNumber = payload.LegalName
	}

	jsonSettings := make(map[string]any)
	if payload.Currency != nil && *payload.Currency != "" {
		jsonSettings["currency"] = *payload.Currency
	}
	if payload.Timezone != nil && *payload.Timezone != "" {
		jsonSettings["timezone"] = *payload.Timezone
	}
	var settingsJSON *string
	if len(jsonSettings) > 0 {
		b, _ := json.Marshal(jsonSettings)
		s := string(b)
		settingsJSON = &s
	}

	stmt := `
		UPDATE organization.organizations
		SET name = COALESCE($1, name),
		    registration_number = COALESCE($2, registration_number),
		    license_number = COALESCE($3, license_number),
		    tax_id = COALESCE($4, tax_id),
		    email = COALESCE($5, email),
		    phone = COALESCE($6, phone),
		    address = COALESCE($7, address),
		    city = COALESCE($8, city),
		    state = COALESCE($9, state),
		    lga = COALESCE($10, lga),
		    country = COALESCE($11, country),
		    logo_url = COALESCE($12, logo_url),
		    custom_domain = COALESCE($13, custom_domain),
		    settings = CASE WHEN $14::jsonb IS NOT NULL THEN COALESCE(settings, '{}'::jsonb) || $14::jsonb ELSE settings END,
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = COALESCE($15, updated_by)
		WHERE id = $16 AND ($17 = 0 OR version = $17)
		RETURNING version
	`

	var newVersion int
	err := dbExec.QueryRow(ctx, stmt,
		payload.Name, regNumber, payload.LicenseNumber, payload.TaxID,
		payload.Email, payload.Phone, payload.Address, payload.City, payload.State, payload.LGA, payload.Country,
		payload.LogoURL, payload.CustomDomain, settingsJSON, updatedByVal, orgID.String(), payload.Version,
	).Scan(&newVersion)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOptimisticLockingConflict
		}
		return nil, fmt.Errorf("failed to update organization profile id=%s: %w", orgID, err)
	}

	return r.GetByID(ctx, orgID)
}

func (r *OrganizationRepository) UpdateSetupState(ctx context.Context, orgID uuid.UUID, newState domain.SetupState, step int, actorID uuid.UUID) (*domain.Organization, error) {
	dbExec := r.server.DB.Conn(ctx)

	var completedAt *time.Time
	if newState == domain.SetupStateVerified {
		now := time.Now()
		completedAt = &now
	}

	stmt := `
		UPDATE organization.organizations
		SET setup_state = $1,
		    setup_step = $2,
		    completed_at = COALESCE($3, completed_at),
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $4
		WHERE id = $5
		RETURNING version
	`

	var newVersion int
	err := dbExec.QueryRow(ctx, stmt, newState, step, completedAt, actorID.String(), orgID.String()).Scan(&newVersion)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrganizationNotFound
		}
		return nil, fmt.Errorf("failed to update organization setup state id=%s: %w", orgID, err)
	}

	return r.GetByID(ctx, orgID)
}

func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `DELETE FROM organization.organizations WHERE id = $1`
	_, err := dbExec.Exec(ctx, stmt, id.String())
	return err
}

func (r *OrganizationRepository) GetSettings(ctx context.Context, orgID uuid.UUID) (*domain.OrganizationSettings, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, logo_url, COALESCE(theme_branding::text, '{}'), custom_domain,
		       email, phone, registration_number, tax_id, address, settings
		FROM organization.organizations
		WHERE id = $1
	`

	var (
		idStr            string
		logoURL          *string
		themeBrandingStr string
		customDomain     *string
		email            *string
		phone            *string
		regNum           *string
		taxID            *string
		address          *string
		settingsMap      map[string]any
	)

	err := dbExec.QueryRow(ctx, stmt, orgID.String()).Scan(
		&idStr, &logoURL, &themeBrandingStr, &customDomain,
		&email, &phone, &regNum, &taxID, &address, &settingsMap,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrganizationNotFound
		}
		return nil, fmt.Errorf("failed to get organization settings id=%s: %w", orgID, err)
	}

	var businessType, timezone, currency, language *string
	if settingsMap != nil {
		if bt, ok := settingsMap["businessType"].(string); ok && bt != "" {
			businessType = &bt
		}
		if tz, ok := settingsMap["timezone"].(string); ok && tz != "" {
			timezone = &tz
		}
		if curr, ok := settingsMap["currency"].(string); ok && curr != "" {
			currency = &curr
		}
		if lang, ok := settingsMap["language"].(string); ok && lang != "" {
			language = &lang
		}
	}

	return &domain.OrganizationSettings{
		ID:             idStr,
		OrganizationID: orgID.String(),
		LogoURL:        logoURL,
		ThemeBranding:  themeBrandingStr,
		CustomDomain:   customDomain,
		SupportEmail:   email,
		SupportPhone:   phone,
		CACNumber:      regNum,
		TINNumber:      taxID,
		TaxNumber:      taxID,
		BusinessType:   businessType,
		Address:        address,
		Timezone:       timezone,
		Currency:       currency,
		Language:       language,
	}, nil
}

func (r *OrganizationRepository) UpdateSettings(ctx context.Context, orgID uuid.UUID, logoURL, themeBranding, customDomain, supportEmail, supportPhone, cacNumber, tinNumber, taxNumber, businessType, address, timezone, currency, language *string) (*domain.OrganizationSettings, error) {
	dbExec := r.server.DB.Conn(ctx)

	query := "UPDATE organization.organizations SET "
	args := pgx.NamedArgs{"id": orgID.String()}
	updates := []string{}

	if logoURL != nil {
		updates = append(updates, "logo_url = @logo_url")
		args["logo_url"] = logoURL
	}
	if themeBranding != nil {
		updates = append(updates, "theme_branding = @theme_branding::jsonb")
		args["theme_branding"] = *themeBranding
	}
	if customDomain != nil {
		updates = append(updates, "custom_domain = @custom_domain")
		args["custom_domain"] = customDomain
	}
	if supportEmail != nil {
		updates = append(updates, "email = @email")
		args["email"] = supportEmail
	}
	if supportPhone != nil {
		updates = append(updates, "phone = @phone")
		args["phone"] = supportPhone
	}
	if cacNumber != nil {
		updates = append(updates, "registration_number = @registration_number")
		args["registration_number"] = cacNumber
	}
	if tinNumber != nil {
		updates = append(updates, "tax_id = @tax_id")
		args["tax_id"] = tinNumber
	} else if taxNumber != nil {
		updates = append(updates, "tax_id = @tax_id")
		args["tax_id"] = taxNumber
	}
	if address != nil {
		updates = append(updates, "address = @address")
		args["address"] = address
	}

	// Merge extra settings into jsonb
	jsonSettings := make(map[string]any)
	if businessType != nil {
		jsonSettings["businessType"] = *businessType
	}
	if timezone != nil {
		jsonSettings["timezone"] = *timezone
	}
	if currency != nil {
		jsonSettings["currency"] = *currency
	}
	if language != nil {
		jsonSettings["language"] = *language
	}

	if len(jsonSettings) > 0 {
		settingsBytes, _ := json.Marshal(jsonSettings)
		updates = append(updates, "settings = COALESCE(settings, '{}'::jsonb) || @extra_settings::jsonb")
		args["extra_settings"] = string(settingsBytes)
	}

	if len(updates) > 0 {
		updates = append(updates, "version = version + 1", "updated_at = CURRENT_TIMESTAMP")
		query += strings.Join(updates, ", ") + " WHERE id = @id"

		_, err := dbExec.Exec(ctx, query, args)
		if err != nil {
			return nil, fmt.Errorf("failed to update organization settings id=%s: %w", orgID, err)
		}
	}

	return r.GetSettings(ctx, orgID)
}

func (r *OrganizationRepository) List(ctx context.Context, userID string, isPlatformAdmin bool) ([]domain.Organization, error) {
	dbExec := r.server.DB.Conn(ctx)

	var (
		rows pgx.Rows
		err  error
	)

	if isPlatformAdmin {
		stmt := `SELECT id FROM organization.organizations ORDER BY name ASC`
		rows, err = dbExec.Query(ctx, stmt)
	} else {
		stmt := `
			SELECT DISTINCT o.id
			FROM organization.organizations o
			JOIN organization.organization_memberships m ON m.organization_id = o.id
			WHERE (m.user_id = $1 OR m.user_id = $2) AND m.is_active = TRUE
			ORDER BY o.id ASC
		`
		rows, err = dbExec.Query(ctx, stmt, userID, userID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	defer rows.Close()

	var orgs []domain.Organization
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err == nil {
			if o, errGet := r.GetByID(ctx, id); errGet == nil {
				orgs = append(orgs, *o)
			}
		}
	}
	return orgs, nil
}

func (r *OrganizationRepository) VerifyMembership(ctx context.Context, orgID uuid.UUID, userID string) (bool, error) {
	if userID == "" || orgID == uuid.Nil {
		return false, nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT EXISTS (
			SELECT 1 FROM organization.organization_memberships
			WHERE organization_id = $1 AND (user_id = $2 OR user_id = $3) AND is_active = TRUE
		)
	`
	var exists bool
	err := dbExec.QueryRow(ctx, stmt, orgID.String(), userID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to verify organization membership for user=%s org=%s: %w", userID, orgID, err)
	}
	return exists, nil
}
