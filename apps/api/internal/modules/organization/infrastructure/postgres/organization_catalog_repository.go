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

type OrganizationCatalogRepository struct {
	server *server.Server
}

func NewOrganizationCatalogRepository(server *server.Server) *OrganizationCatalogRepository {
	return &OrganizationCatalogRepository{server: server}
}

func (r *OrganizationCatalogRepository) ListCatalogItems(ctx context.Context, orgID uuid.UUID, domainType string) ([]domain.OrganizationCatalogItem, error) {
	dbExec := r.server.DB.Conn(ctx)

	stmt := `
		SELECT id, organization_id, master_catalog_id, domain_type, code, name, description,
		       base_price, COALESCE(currency, 'NGN'), is_active, version, created_at, updated_at, updated_by
		FROM organization.catalog_items
		WHERE organization_id = $1
	`
	args := []any{orgID.String()}

	if domainType != "" {
		stmt += " AND domain_type = $2"
		args = append(args, domainType)
	}

	stmt += " ORDER BY domain_type ASC, name ASC"

	rows, err := dbExec.Query(ctx, stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query organization catalog items: %w", err)
	}
	defer rows.Close()

	var items []domain.OrganizationCatalogItem
	for rows.Next() {
		var (
			item          domain.OrganizationCatalogItem
			masterIDStr   *string
			updatedByStr  *string
		)
		err := rows.Scan(
			&item.ID, &item.OrganizationID, &masterIDStr, &item.DomainType, &item.Code, &item.Name, &item.Description,
			&item.BasePrice, &item.Currency, &item.IsActive, &item.Version, &item.CreatedAt, &item.UpdatedAt, &updatedByStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan catalog item row: %w", err)
		}
		if masterIDStr != nil && *masterIDStr != "" {
			if parsed, pErr := uuid.Parse(*masterIDStr); pErr == nil {
				item.MasterCatalogID = &parsed
			}
		}
		if updatedByStr != nil && *updatedByStr != "" {
			if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
				item.UpdatedBy = &parsed
			}
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *OrganizationCatalogRepository) GetCatalogItemByID(ctx context.Context, orgID, itemID uuid.UUID) (*domain.OrganizationCatalogItem, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, master_catalog_id, domain_type, code, name, description,
		       base_price, COALESCE(currency, 'NGN'), is_active, version, created_at, updated_at, updated_by
		FROM organization.catalog_items
		WHERE organization_id = $1 AND id = $2
		LIMIT 1
	`

	var (
		item          domain.OrganizationCatalogItem
		masterIDStr   *string
		updatedByStr  *string
	)
	err := dbExec.QueryRow(ctx, stmt, orgID.String(), itemID.String()).Scan(
		&item.ID, &item.OrganizationID, &masterIDStr, &item.DomainType, &item.Code, &item.Name, &item.Description,
		&item.BasePrice, &item.Currency, &item.IsActive, &item.Version, &item.CreatedAt, &item.UpdatedAt, &updatedByStr,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCatalogItemNotFound
		}
		return nil, fmt.Errorf("failed to query catalog item id=%s: %w", itemID, err)
	}

	if masterIDStr != nil && *masterIDStr != "" {
		if parsed, pErr := uuid.Parse(*masterIDStr); pErr == nil {
			item.MasterCatalogID = &parsed
		}
	}
	if updatedByStr != nil && *updatedByStr != "" {
		if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
			item.UpdatedBy = &parsed
		}
	}

	return &item, nil
}

func (r *OrganizationCatalogRepository) CreateCatalogItem(ctx context.Context, item *domain.OrganizationCatalogItem, actorID uuid.UUID) (*domain.OrganizationCatalogItem, error) {
	dbExec := r.server.DB.Conn(ctx)
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}

	curr := "NGN"
	if item.Currency != "" {
		curr = item.Currency
	}

	var masterStr *string
	if item.MasterCatalogID != nil {
		str := item.MasterCatalogID.String()
		masterStr = &str
	}

	stmt := `
		INSERT INTO organization.catalog_items (
			id, organization_id, master_catalog_id, domain_type, code, name, description,
			base_price, currency, is_active, version, updated_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, TRUE, 1, $10)
		RETURNING version, created_at, updated_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		item.ID, item.OrganizationID, masterStr, item.DomainType, item.Code, item.Name, item.Description,
		item.BasePrice, curr, actorID.String(),
	).Scan(&item.Version, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "uk_org_catalog_code" {
			return nil, domain.ErrDuplicateCatalogCode
		}
		return nil, fmt.Errorf("failed to create catalog item: %w", err)
	}

	item.UpdatedBy = &actorID
	item.IsActive = true
	return item, nil
}

func (r *OrganizationCatalogRepository) UpdateCatalogItem(ctx context.Context, item *domain.OrganizationCatalogItem, actorID uuid.UUID) (*domain.OrganizationCatalogItem, error) {
	dbExec := r.server.DB.Conn(ctx)

	stmt := `
		UPDATE organization.catalog_items
		SET name = $1,
		    description = $2,
		    base_price = $3,
		    currency = COALESCE($4, currency),
		    is_active = COALESCE($5, is_active),
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $6
		WHERE id = $7 AND organization_id = $8 AND version = $9
		RETURNING version, updated_at
	`

	var (
		newVersion int
		updatedAt  time.Time
	)
	err := dbExec.QueryRow(ctx, stmt,
		item.Name, item.Description, item.BasePrice, item.Currency, item.IsActive, actorID.String(), item.ID, item.OrganizationID, item.Version,
	).Scan(&newVersion, &updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCatalogItemNotFound
		}
		return nil, fmt.Errorf("failed to update catalog item id=%s: %w", item.ID, err)
	}

	item.Version = newVersion
	item.UpdatedAt = updatedAt
	item.UpdatedBy = &actorID
	return item, nil
}

func (r *OrganizationCatalogRepository) SetBranchPriceOverride(ctx context.Context, override *domain.BranchPriceOverride, actorID uuid.UUID) (*domain.BranchPriceOverride, error) {
	dbExec := r.server.DB.Conn(ctx)
	if override.ID == uuid.Nil {
		override.ID = uuid.New()
	}

	stmt := `
		INSERT INTO organization.branch_price_overrides (
			id, organization_id, facility_branch_id, catalog_item_id, override_price, updated_by
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (facility_branch_id, catalog_item_id)
		DO UPDATE SET override_price = EXCLUDED.override_price, updated_at = CURRENT_TIMESTAMP, updated_by = EXCLUDED.updated_by
		RETURNING created_at, updated_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		override.ID, override.OrganizationID, override.FacilityBranchID, override.CatalogItemID, override.OverridePrice, actorID.String(),
	).Scan(&override.CreatedAt, &override.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to set branch price override: %w", err)
	}

	override.UpdatedBy = &actorID
	return override, nil
}

func (r *OrganizationCatalogRepository) GetBranchPriceOverride(ctx context.Context, orgID, branchID, itemID uuid.UUID) (*domain.BranchPriceOverride, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, facility_branch_id, catalog_item_id, override_price, created_at, updated_at, updated_by
		FROM organization.branch_price_overrides
		WHERE organization_id = $1 AND facility_branch_id = $2 AND catalog_item_id = $3
		LIMIT 1
	`

	var (
		override     domain.BranchPriceOverride
		updatedByStr *string
	)
	err := dbExec.QueryRow(ctx, stmt, orgID, branchID, itemID).Scan(
		&override.ID, &override.OrganizationID, &override.FacilityBranchID, &override.CatalogItemID,
		&override.OverridePrice, &override.CreatedAt, &override.UpdatedAt, &updatedByStr,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBranchPriceNotFound
		}
		return nil, fmt.Errorf("failed to query branch price override: %w", err)
	}

	if updatedByStr != nil && *updatedByStr != "" {
		if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
			override.UpdatedBy = &parsed
		}
	}

	return &override, nil
}

func (r *OrganizationCatalogRepository) ListBranchPriceOverrides(ctx context.Context, orgID, branchID uuid.UUID) ([]domain.BranchPriceOverride, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, facility_branch_id, catalog_item_id, override_price, created_at, updated_at, updated_by
		FROM organization.branch_price_overrides
		WHERE organization_id = $1 AND facility_branch_id = $2
	`

	rows, err := dbExec.Query(ctx, stmt, orgID, branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to query branch price overrides: %w", err)
	}
	defer rows.Close()

	var overrides []domain.BranchPriceOverride
	for rows.Next() {
		var (
			o            domain.BranchPriceOverride
			updatedByStr *string
		)
		err := rows.Scan(
			&o.ID, &o.OrganizationID, &o.FacilityBranchID, &o.CatalogItemID,
			&o.OverridePrice, &o.CreatedAt, &o.UpdatedAt, &updatedByStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan branch price override row: %w", err)
		}
		if updatedByStr != nil && *updatedByStr != "" {
			if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
				o.UpdatedBy = &parsed
			}
		}
		overrides = append(overrides, o)
	}

	return overrides, nil
}

func (r *OrganizationCatalogRepository) CreateInsuranceProvider(ctx context.Context, provider *domain.InsuranceProvider, actorID uuid.UUID) (*domain.InsuranceProvider, error) {
	dbExec := r.server.DB.Conn(ctx)
	if provider.ID == uuid.Nil {
		provider.ID = uuid.New()
	}

	coverageVal := 100.00
	if provider.CoveragePercentage > 0 {
		coverageVal = provider.CoveragePercentage
	}

	stmt := `
		INSERT INTO organization.insurance_providers (
			id, organization_id, name, code, coverage_percentage, is_active, updated_by
		)
		VALUES ($1, $2, $3, $4, $5, TRUE, $6)
		RETURNING created_at, updated_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		provider.ID, provider.OrganizationID, provider.Name, provider.Code, coverageVal, actorID.String(),
	).Scan(&provider.CreatedAt, &provider.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "uk_org_insurance_code" {
			return nil, domain.ErrDuplicateInsuranceCode
		}
		return nil, fmt.Errorf("failed to create insurance provider: %w", err)
	}

	provider.IsActive = true
	provider.CoveragePercentage = coverageVal
	provider.UpdatedBy = &actorID
	return provider, nil
}

func (r *OrganizationCatalogRepository) ListInsuranceProviders(ctx context.Context, orgID uuid.UUID) ([]domain.InsuranceProvider, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, name, code, coverage_percentage, is_active, created_at, updated_at, updated_by
		FROM organization.insurance_providers
		WHERE organization_id = $1
		ORDER BY name ASC
	`

	rows, err := dbExec.Query(ctx, stmt, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list insurance providers: %w", err)
	}
	defer rows.Close()

	var providers []domain.InsuranceProvider
	for rows.Next() {
		var (
			p            domain.InsuranceProvider
			updatedByStr *string
		)
		err := rows.Scan(
			&p.ID, &p.OrganizationID, &p.Name, &p.Code, &p.CoveragePercentage, &p.IsActive,
			&p.CreatedAt, &p.UpdatedAt, &updatedByStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan insurance provider row: %w", err)
		}
		if updatedByStr != nil && *updatedByStr != "" {
			if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
				p.UpdatedBy = &parsed
			}
		}
		providers = append(providers, p)
	}

	return providers, nil
}
