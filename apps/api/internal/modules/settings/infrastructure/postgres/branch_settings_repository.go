package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golangnigeria/curexal/internal/modules/settings/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type BranchSettingsRepository struct {
	server *server.Server
}

func NewBranchSettingsRepository(server *server.Server) *BranchSettingsRepository {
	return &BranchSettingsRepository{server: server}
}

func (r *BranchSettingsRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.BranchSettings, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `SELECT 
		id, tenant_id, general_config, financial_config, inventory_config, 
		integrations_config, notifications_config, document_header_config, 
		patient_config, lims_config, consultation_config, staff_config, 
		created_at, updated_at 
	FROM branch_settings 
	WHERE tenant_id = @tenant_id`

	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"tenant_id": tenantID})
	if err != nil {
		return nil, fmt.Errorf("failed to query branch settings: %w", err)
	}

	settings, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.BranchSettings])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return r.CreateDefault(ctx, tenantID)
		}
		return nil, fmt.Errorf("failed to collect branch settings row: %w", err)
	}

	return &settings, nil
}

func (r *BranchSettingsRepository) CreateDefault(ctx context.Context, tenantID uuid.UUID) (*domain.BranchSettings, error) {
	dbExec := r.server.DB.Conn(ctx)

	stmt := `INSERT INTO branch_settings (
		tenant_id, general_config, financial_config, inventory_config, 
		integrations_config, notifications_config, document_header_config, 
		patient_config, lims_config, consultation_config, staff_config
	) VALUES (
		@tenant_id, '{}'::jsonb, '{}'::jsonb, '{}'::jsonb, '{}'::jsonb, '{}'::jsonb, 
		'{}'::jsonb, '{}'::jsonb, '{}'::jsonb, '{}'::jsonb, '{}'::jsonb
	) RETURNING 
		id, tenant_id, general_config, financial_config, inventory_config, 
		integrations_config, notifications_config, document_header_config, 
		patient_config, lims_config, consultation_config, staff_config, 
		created_at, updated_at`

	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{
		"tenant_id": tenantID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to provision default branch settings: %w", err)
	}

	settings, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.BranchSettings])
	if err != nil {
		return nil, fmt.Errorf("failed to collect default branch settings row: %w", err)
	}

	return &settings, nil
}

func (r *BranchSettingsRepository) UpsertSection(ctx context.Context, tenantID uuid.UUID, section string, payload map[string]any) (*domain.BranchSettings, error) {
	dbExec := r.server.DB.Conn(ctx)

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal section settings payload: %w", err)
	}

	columnName := section + "_config"
	stmt := fmt.Sprintf(`
		INSERT INTO branch_settings (tenant_id, %s)
		VALUES (@tenant_id, @payload::jsonb)
		ON CONFLICT (tenant_id) DO UPDATE
		SET %s = @payload::jsonb, updated_at = NOW()
		RETURNING 
			id, tenant_id, general_config, financial_config, inventory_config, 
			integrations_config, notifications_config, document_header_config, 
			patient_config, lims_config, consultation_config, staff_config, 
			created_at, updated_at
	`, columnName, columnName)

	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{
		"tenant_id": tenantID,
		"payload":   string(jsonBytes),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upsert branch settings section %s: %w", section, err)
	}

	settings, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.BranchSettings])
	if err != nil {
		return nil, fmt.Errorf("failed to collect updated branch settings row: %w", err)
	}

	return &settings, nil
}

func (r *BranchSettingsRepository) ResetToDefaults(ctx context.Context, tenantID uuid.UUID, section *string) (*domain.BranchSettings, error) {
	dbExec := r.server.DB.Conn(ctx)

	if section != nil && *section != "" {
		columnName := *section + "_config"
		stmt := fmt.Sprintf(`
			UPDATE branch_settings
			SET %s = '{}'::jsonb, updated_at = NOW()
			WHERE tenant_id = @tenant_id
			RETURNING 
				id, tenant_id, general_config, financial_config, inventory_config, 
				integrations_config, notifications_config, document_header_config, 
				patient_config, lims_config, consultation_config, staff_config, 
				created_at, updated_at
		`, columnName)
		rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"tenant_id": tenantID})
		if err != nil {
			return nil, fmt.Errorf("failed to reset section settings: %w", err)
		}
		settings, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.BranchSettings])
		if err != nil {
			return nil, err
		}
		return &settings, nil
	}

	stmt := `
		UPDATE branch_settings
		SET general_config = '{}'::jsonb, financial_config = '{}'::jsonb, inventory_config = '{}'::jsonb,
		    integrations_config = '{}'::jsonb, notifications_config = '{}'::jsonb, document_header_config = '{}'::jsonb,
		    patient_config = '{}'::jsonb, lims_config = '{}'::jsonb, consultation_config = '{}'::jsonb, staff_config = '{}'::jsonb,
		    updated_at = NOW()
		WHERE tenant_id = @tenant_id
		RETURNING 
			id, tenant_id, general_config, financial_config, inventory_config, 
			integrations_config, notifications_config, document_header_config, 
			patient_config, lims_config, consultation_config, staff_config, 
			created_at, updated_at
	`
	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"tenant_id": tenantID})
	if err != nil {
		return nil, fmt.Errorf("failed to reset all settings: %w", err)
	}
	settings, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.BranchSettings])
	if err != nil {
		return nil, err
	}
	return &settings, nil
}
