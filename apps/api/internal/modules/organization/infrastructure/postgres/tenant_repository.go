package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TenantRepository struct {
	server *server.Server
}

func NewTenantRepository(server *server.Server) *TenantRepository {
	return &TenantRepository{server: server}
}

func (r *TenantRepository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	dbExec := r.server.DB.Conn(ctx)
	var exists bool
	stmt := "SELECT EXISTS(SELECT 1 FROM workspace.workspaces WHERE slug = @slug)"
	err := dbExec.QueryRow(ctx, stmt, pgx.NamedArgs{"slug": slug}).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if slug exists: %w", err)
	}
	return exists, nil
}

func (r *TenantRepository) CreateTenant(ctx context.Context, userID string, name, slug, orgID, location, phone, address string, logoURL, currency *string, modules []string) (*domain.Tenant, error) {
	dbExec := r.server.DB.Conn(ctx)

	var orgExists bool
	err := dbExec.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM organization.organizations WHERE id = $1)", orgID).Scan(&orgExists)
	if err != nil {
		return nil, fmt.Errorf("failed to verify organization existence: %w", err)
	}
	if !orgExists {
		return nil, errs.ErrOrganizationNotFound
	}

	var orgOwnerID string
	err = dbExec.QueryRow(ctx, `
		SELECT user_id FROM organization.organization_memberships 
		WHERE organization_id = $1 AND role_title = 'owner'
		LIMIT 1
	`, orgID).Scan(&orgOwnerID)
	if err != nil {
		orgOwnerID = userID
	}

	type TenantSettings struct {
		Location string `json:"location,omitempty"`
		Phone    string `json:"phone,omitempty"`
		Address  string `json:"address,omitempty"`
	}

	settingsVal := TenantSettings{
		Location: location,
		Phone:    phone,
		Address:  address,
	}

	settingsJSON, errMarshal := json.Marshal(settingsVal)
	if errMarshal != nil {
		settingsJSON = []byte("{}")
	}

	tenantID := uuid.New().String()
	stmtInsertTenant := `
		INSERT INTO workspace.workspaces (id, name, slug, logo_url, settings, organization_id, currency)
		VALUES (@id, @name, @slug, @logo_url, @settings::jsonb, @organization_id, @currency)
		RETURNING id
	`
	curr := "NGN"
	if currency != nil && *currency != "" {
		curr = *currency
	}

	err = dbExec.QueryRow(ctx, stmtInsertTenant, pgx.NamedArgs{
		"id":              tenantID,
		"name":            name,
		"slug":            slug,
		"logo_url":        logoURL,
		"settings":        string(settingsJSON),
		"organization_id": orgID,
		"currency":        curr,
	}).Scan(&tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant record: %w", err)
	}

	// Insert owner membership
	membershipID := uuid.New().String()
	stmtInsertMembership := `
		INSERT INTO organization.organization_memberships (id, user_id, organization_id, role_title, is_active, joined_at)
		VALUES (@id, @user_id, @org_id, 'owner', TRUE, CURRENT_TIMESTAMP)
		ON CONFLICT (organization_id, user_id) DO NOTHING
	`
	_, err = dbExec.Exec(ctx, stmtInsertMembership, pgx.NamedArgs{
		"id":      membershipID,
		"user_id": orgOwnerID,
		"org_id":  orgID,
	})
	if err != nil {
		r.server.Logger.Error().Err(err).Str("user_id", orgOwnerID).Str("tenant_id", tenantID).Msg("failed to create tenant membership")
	}

	parsedTenantID, _ := uuid.Parse(tenantID)
	createdTenant := &domain.Tenant{
		ID:             parsedTenantID,
		Name:           name,
		Slug:           slug,
		LogoURL:        logoURL,
		Settings:       string(settingsJSON),
		OrganizationID: orgID,
		Currency:       curr,
		EnabledModules: modules,
	}

	return createdTenant, nil
}

func (r *TenantRepository) GetTenantByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT t.id, t.name, COALESCE(o.slug, t.slug) AS slug, t.logo_url, t.settings, t.organization_id, t.currency, t.created_at, t.updated_at 
		FROM workspace.workspaces t
		LEFT JOIN organization.organizations o ON o.id = t.organization_id
		WHERE t.id = @id
	`

	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"id": id.String()})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get tenant query for id=%s: %w", id, err)
	}

	t, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Tenant])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:tenant for id=%s: %w", id, err)
	}

	return &t, nil
}

func (r *TenantRepository) UpdateTenant(ctx context.Context, id uuid.UUID, name, slug, logoURL, currency *string, settingsJSON *string) (*domain.Tenant, error) {
	dbExec := r.server.DB.Conn(ctx)
	setClauses := []string{}
	args := pgx.NamedArgs{
		"id": id.String(),
	}

	if name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *name
	}
	if slug != nil {
		setClauses = append(setClauses, "slug = @slug")
		args["slug"] = strings.ToLower(strings.TrimSpace(*slug))
	}
	if logoURL != nil {
		setClauses = append(setClauses, "logo_url = @logo_url")
		args["logo_url"] = *logoURL
	}
	if currency != nil {
		setClauses = append(setClauses, "currency = @currency")
		args["currency"] = *currency
	}
	if settingsJSON != nil {
		setClauses = append(setClauses, "settings = @settings")
		args["settings"] = *settingsJSON
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no fields to update")
	}

	stmt := fmt.Sprintf(`
		WITH updated AS (
			UPDATE workspace.workspaces SET %s WHERE id = @id RETURNING id, name, slug, logo_url, settings, organization_id, currency, created_at, updated_at
		)
		SELECT u.id, u.name, COALESCE(o.slug, u.slug) AS slug, u.logo_url, u.settings, u.organization_id, u.currency, u.created_at, u.updated_at 
		FROM updated u
		LEFT JOIN organization.organizations o ON o.id = u.organization_id
	`, strings.Join(setClauses, ", "))

	rows, err := dbExec.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update tenant query for id=%s: %w", id, err)
	}

	t, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Tenant])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:tenant for id=%s: %w", id, err)
	}

	return &t, nil
}

func (r *TenantRepository) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := "DELETE FROM workspace.workspaces WHERE id = @id"

	result, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{"id": id.String()})
	if err != nil {
		return fmt.Errorf("failed to execute delete tenant query for id=%s: %w", id, err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("tenant not found")
	}

	return nil
}

func (r *TenantRepository) ListTenants(ctx context.Context) ([]domain.Tenant, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT t.id, t.name, COALESCE(o.slug, t.slug) AS slug, t.logo_url, t.settings, t.organization_id, t.currency, t.created_at, t.updated_at 
		FROM workspace.workspaces t
		LEFT JOIN organization.organizations o ON o.id = t.organization_id
		ORDER BY t.created_at DESC
	`

	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute list tenants query: %w", err)
	}
	defer rows.Close()

	tenantsList, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Tenant])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row list from table:tenant: %w", err)
	}

	return tenantsList, nil
}

func (r *TenantRepository) CountActiveMembers(ctx context.Context, tenantID uuid.UUID) (int, error) {
	dbExec := r.server.DB.Conn(ctx)
	var count int
	stmt := `SELECT COUNT(*) FROM organization.organization_memberships WHERE organization_id = @tenant_id AND is_active = TRUE`
	err := dbExec.QueryRow(ctx, stmt, pgx.NamedArgs{"tenant_id": tenantID.String()}).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active members for tenant: %w", err)
	}
	return count, nil
}

func (r *TenantRepository) ListBranchesByOrgID(ctx context.Context, orgID uuid.UUID) ([]domain.Tenant, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT t.id, t.name, t.slug, t.logo_url, COALESCE(t.settings::text, '{}') AS settings, COALESCE(t.organization_id::text, '') AS organization_id, COALESCE(t.currency, 'NGN') AS currency
		FROM workspace.workspaces t
		WHERE t.organization_id = @org_id
		ORDER BY t.created_at DESC
	`

	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"org_id": orgID.String()})
	if err != nil {
		return nil, fmt.Errorf("failed to list branches by organization id=%s: %w", orgID, err)
	}
	defer rows.Close()

	branches, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Tenant])
	if err != nil {
		return nil, fmt.Errorf("failed to collect branch list: %w", err)
	}

	return branches, nil
}

func (r *TenantRepository) ListAllTenantsSelector(ctx context.Context) ([]model.TenantSelectorItem, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := "SELECT id, name, slug FROM workspace.workspaces ORDER BY name ASC"
	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query all tenants selector: %w", err)
	}
	defer rows.Close()

	var items []model.TenantSelectorItem
	for rows.Next() {
		var item model.TenantSelectorItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Slug); err != nil {
			return nil, fmt.Errorf("failed to scan tenant selector item: %w", err)
		}
		items = append(items, item)
	}
	if items == nil {
		items = []model.TenantSelectorItem{}
	}
	return items, nil
}

func (r *TenantRepository) GetTenantBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT t.id, t.name, t.slug, t.logo_url, t.settings, t.organization_id, t.currency
		FROM workspace.workspaces t
		WHERE t.slug = $1 OR t.id = $1
		LIMIT 1
	`
	t := &domain.Tenant{}
	err := dbExec.QueryRow(ctx, stmt, slug).Scan(
		&t.ID, &t.Name, &t.Slug, &t.LogoURL, &t.Settings, &t.OrganizationID, &t.Currency,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}
