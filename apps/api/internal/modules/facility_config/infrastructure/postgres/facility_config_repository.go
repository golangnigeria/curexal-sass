package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/facility_config/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FacilityConfigRepository struct {
	server *server.Server
}

func NewFacilityConfigRepository(server *server.Server) *FacilityConfigRepository {
	return &FacilityConfigRepository{server: server}
}

func defaultFacilityTypes() []domain.FacilityTypeDTO {
	return []domain.FacilityTypeDTO{
		{
			ID:          "fac_hospital",
			Code:        "hospital",
			DisplayName: "Hospital",
			Category:    "clinical",
			IconKey:     "Building2",
			Description: "Full service inpatient and outpatient tertiary hospital",
		},
		{
			ID:          "fac_clinic",
			Code:        "clinic",
			DisplayName: "Outpatient Clinic",
			Category:    "clinical",
			IconKey:     "Stethoscope",
			Description: "Primary healthcare and outpatient clinical practice",
		},
		{
			ID:          "fac_diagnostic",
			Code:        "diagnostic_center",
			DisplayName: "Diagnostic Center",
			Category:    "diagnostic",
			IconKey:     "Activity",
			Description: "Diagnostic pathology and imaging center",
		},
		{
			ID:          "fac_pharmacy",
			Code:        "pharmacy",
			DisplayName: "Community Pharmacy",
			Category:    "retail",
			IconKey:     "Pill",
			Description: "Retail and clinical prescription pharmacy",
		},
		{
			ID:          "fac_laboratory",
			Code:        "laboratory",
			DisplayName: "Medical Laboratory",
			Category:    "diagnostic",
			IconKey:     "Microscope",
			Description: "Pathology, clinical biochemistry and microbiology laboratory",
		},
	}
}

func (r *FacilityConfigRepository) GetActiveFacilityTypes(ctx context.Context) ([]domain.FacilityTypeDTO, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT 
			id::text, 
			code, 
			name AS display_name, 
			category, 
			COALESCE(icon_key, 'Building2') AS icon_key, 
			COALESCE(description, '') AS description
		FROM platform.facility_types
		WHERE is_active = TRUE
		ORDER BY name ASC
	`
	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return defaultFacilityTypes(), nil
	}
	defer rows.Close()

	var types []domain.FacilityTypeDTO
	for rows.Next() {
		var t domain.FacilityTypeDTO
		var desc string
		err := rows.Scan(&t.ID, &t.Code, &t.DisplayName, &t.Category, &t.IconKey, &desc)
		if err != nil {
			return nil, fmt.Errorf("failed to scan facility type row: %w", err)
		}
		t.Description = desc
		types = append(types, t)
	}
	if len(types) == 0 {
		return defaultFacilityTypes(), nil
	}
	return types, nil
}

func (r *FacilityConfigRepository) ListFacilityTypeEntities(ctx context.Context) ([]domain.FacilityTypeEntity, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, code, name, category, icon_key, description, is_active, version, created_at, updated_at, updated_by
		FROM platform.facility_types
		ORDER BY name ASC
	`
	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query platform facility types: %w", err)
	}
	defer rows.Close()

	var list []domain.FacilityTypeEntity
	for rows.Next() {
		var (
			ft           domain.FacilityTypeEntity
			updatedByStr *string
			iconKey      *string
			desc         *string
		)
		err := rows.Scan(
			&ft.ID, &ft.Code, &ft.Name, &ft.Category, &iconKey, &desc,
			&ft.IsActive, &ft.Version, &ft.CreatedAt, &ft.UpdatedAt, &updatedByStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan facility type entity: %w", err)
		}
		if iconKey != nil {
			ft.IconKey = *iconKey
		}
		if desc != nil {
			ft.Description = *desc
		}
		if updatedByStr != nil && *updatedByStr != "" {
			if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
				ft.UpdatedBy = &parsed
			}
		}

		caps, _ := r.getCapabilitiesForFacilityType(ctx, ft.ID)
		ft.DefaultCapabilities = caps

		list = append(list, ft)
	}

	return list, nil
}

func (r *FacilityConfigRepository) getCapabilitiesForFacilityType(ctx context.Context, facilityTypeID uuid.UUID) ([]domain.FacilityCapability, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT fc.id, fc.facility_type_id, fc.capability_id, c.code, c.name, fc.is_default
		FROM platform.facility_capabilities fc
		JOIN subscription.capabilities c ON c.id = fc.capability_id
		WHERE fc.facility_type_id = $1
		ORDER BY c.code ASC
	`
	rows, err := dbExec.Query(ctx, stmt, facilityTypeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.FacilityCapability
	for rows.Next() {
		var fc domain.FacilityCapability
		if err := rows.Scan(&fc.ID, &fc.FacilityTypeID, &fc.CapabilityID, &fc.CapabilityCode, &fc.CapabilityName, &fc.IsDefault); err == nil {
			list = append(list, fc)
		}
	}
	return list, nil
}

func (r *FacilityConfigRepository) CreateFacilityType(ctx context.Context, ft *domain.FacilityTypeEntity, updatedBy uuid.UUID) (*domain.FacilityTypeEntity, error) {
	dbExec := r.server.DB.Conn(ctx)
	if ft.ID == uuid.Nil {
		ft.ID = uuid.New()
	}

	stmt := `
		INSERT INTO platform.facility_types (id, code, name, category, icon_key, description, is_active, version, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 1, $8)
		RETURNING version, created_at, updated_at
	`
	var (
		version   int
		createdAt time.Time
		updatedAt time.Time
	)
	err := dbExec.QueryRow(ctx, stmt,
		ft.ID, ft.Code, ft.Name, ft.Category, ft.IconKey, ft.Description, ft.IsActive, updatedBy.String(),
	).Scan(&version, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert platform facility type: %w", err)
	}

	ft.Version = version
	ft.CreatedAt = createdAt
	ft.UpdatedAt = updatedAt
	ft.UpdatedBy = &updatedBy

	return ft, nil
}

func (r *FacilityConfigRepository) UpdateFacilityType(ctx context.Context, ft *domain.FacilityTypeEntity, updatedBy uuid.UUID) (*domain.FacilityTypeEntity, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE platform.facility_types
		SET name = $1,
		    category = $2,
		    icon_key = $3,
		    description = $4,
		    is_active = $5,
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $6
		WHERE id = $7 AND version = $8
		RETURNING version, updated_at
	`
	var (
		newVersion int
		updatedAt  time.Time
	)
	err := dbExec.QueryRow(ctx, stmt,
		ft.Name, ft.Category, ft.IconKey, ft.Description, ft.IsActive, updatedBy.String(), ft.ID, ft.Version,
	).Scan(&newVersion, &updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOptimisticLockingConflict
		}
		return nil, fmt.Errorf("failed to update platform facility type: %w", err)
	}

	ft.Version = newVersion
	ft.UpdatedAt = updatedAt
	ft.UpdatedBy = &updatedBy

	return ft, nil
}

func (r *FacilityConfigRepository) GetRegistrationForm(ctx context.Context, typeID string) (*domain.RegistrationFormDTO, error) {
	sectionsJSON := json.RawMessage(`[
		{"id": "sec_org_info", "title": "Facility Information", "fields": ["facility_name", "registration_number", "license_number"]},
		{"id": "sec_contact", "title": "Contact Details", "fields": ["email", "phone", "address", "city", "state"]}
	]`)
	docsJSON := json.RawMessage(`[
		{"id": "doc_cac", "title": "CAC Certificate", "required": true},
		{"id": "doc_license", "title": "Operational Practice License", "required": true}
	]`)
	return &domain.RegistrationFormDTO{
		FacilityTypeID:    typeID,
		Version:           1,
		Sections:          sectionsJSON,
		RequiredDocuments: docsJSON,
	}, nil
}

func (r *FacilityConfigRepository) GetNavigationMenu(ctx context.Context, typeID string) (*domain.NavigationMenuDTO, error) {
	itemsJSON := json.RawMessage(`[
		{"id": "nav_dashboard", "title": "Dashboard", "path": "/workspace/dashboard", "icon": "LayoutDashboard"},
		{"id": "nav_reception", "title": "Reception", "path": "/workspace/patients", "icon": "UserPlus"}
	]`)
	return &domain.NavigationMenuDTO{
		FacilityTypeID: typeID,
		MenuItems:      itemsJSON,
	}, nil
}

func (r *FacilityConfigRepository) GetSetupSteps(ctx context.Context, typeID string) ([]domain.SetupStepDTO, error) {
	return []domain.SetupStepDTO{
		{StepNumber: 1, Title: "Basic Profile", Description: "Set up facility name and contact details"},
		{StepNumber: 2, Title: "Departments", Description: "Configure operational units"},
		{StepNumber: 3, Title: "Verification", Description: "Submit regulatory licenses"},
	}, nil
}

func (r *FacilityConfigRepository) GetDashboard(ctx context.Context, typeID string) (*domain.DashboardDTO, error) {
	widgetsJSON := json.RawMessage(`[
		{"id": "wgt_visits", "title": "Today Visits", "type": "counter"},
		{"id": "wgt_revenue", "title": "Daily Revenue", "type": "chart"}
	]`)
	return &domain.DashboardDTO{
		FacilityTypeID: typeID,
		Widgets:        widgetsJSON,
	}, nil
}

func (r *FacilityConfigRepository) GetTenantOverrides(ctx context.Context, tenantID uuid.UUID, branchID *uuid.UUID) (map[string]any, error) {
	return map[string]any{}, nil
}

