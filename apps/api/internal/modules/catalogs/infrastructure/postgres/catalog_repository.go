package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/catalogs/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CatalogRepository struct {
	server *server.Server
}

func NewCatalogRepository(s *server.Server) *CatalogRepository {
	return &CatalogRepository{server: s}
}

func (r *CatalogRepository) tableForDomain(d domain.CatalogDomain) (string, error) {
	switch d {
	case domain.ClinicalDomain:
		return "platform.clinical_catalogs", nil
	case domain.LabDomain:
		return "platform.lab_catalogs", nil
	case domain.RadiologyDomain:
		return "platform.radiology_catalogs", nil
	case domain.PharmacyDomain:
		return "platform.pharmacy_catalogs", nil
	default:
		return "", domain.ErrInvalidCatalogDomain
	}
}

func (r *CatalogRepository) scanCatalogItem(catDomain domain.CatalogDomain, rows pgx.Rows) (*domain.CatalogItem, error) {
	var (
		item         domain.CatalogItem
		updatedByStr *string
	)
	item.Domain = catDomain

	if catDomain == domain.ClinicalDomain {
		var sysGroup string
		err := rows.Scan(
			&item.ID, &item.Category, &item.Code, &item.Name, &item.Description, &sysGroup,
			&item.IsActive, &item.Version, &item.CreatedAt, &item.UpdatedAt, &updatedByStr,
		)
		if err != nil {
			return nil, err
		}
		item.SystemGroup = sysGroup
	} else {
		var rawMeta []byte
		err := rows.Scan(
			&item.ID, &item.Category, &item.Code, &item.Name, &item.Description, &rawMeta,
			&item.IsActive, &item.Version, &item.CreatedAt, &item.UpdatedAt, &updatedByStr,
		)
		if err != nil {
			return nil, err
		}
		if len(rawMeta) > 0 {
			item.Metadata = rawMeta
			var m map[string]any
			if err := json.Unmarshal(rawMeta, &m); err == nil {
				if bp, ok := m["basePrice"]; ok {
					switch v := bp.(type) {
					case float64:
						item.BasePrice = v
					case int:
						item.BasePrice = float64(v)
					}
				}
			}
		}
	}

	if updatedByStr != nil && *updatedByStr != "" {
		if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
			item.UpdatedBy = &parsed
		}
	}

	return &item, nil
}

func (r *CatalogRepository) scanCatalogRow(catDomain domain.CatalogDomain, row pgx.Row) (*domain.CatalogItem, error) {
	var (
		item         domain.CatalogItem
		updatedByStr *string
	)
	item.Domain = catDomain

	if catDomain == domain.ClinicalDomain {
		var sysGroup string
		err := row.Scan(
			&item.ID, &item.Category, &item.Code, &item.Name, &item.Description, &sysGroup,
			&item.IsActive, &item.Version, &item.CreatedAt, &item.UpdatedAt, &updatedByStr,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, domain.ErrCatalogItemNotFound
			}
			return nil, err
		}
		item.SystemGroup = sysGroup
	} else {
		var rawMeta []byte
		err := row.Scan(
			&item.ID, &item.Category, &item.Code, &item.Name, &item.Description, &rawMeta,
			&item.IsActive, &item.Version, &item.CreatedAt, &item.UpdatedAt, &updatedByStr,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, domain.ErrCatalogItemNotFound
			}
			return nil, err
		}
		if len(rawMeta) > 0 {
			item.Metadata = rawMeta
			var m map[string]any
			if err := json.Unmarshal(rawMeta, &m); err == nil {
				if bp, ok := m["basePrice"]; ok {
					switch v := bp.(type) {
					case float64:
						item.BasePrice = v
					case int:
						item.BasePrice = float64(v)
					}
				}
			}
		}
	}

	if updatedByStr != nil && *updatedByStr != "" {
		if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
			item.UpdatedBy = &parsed
		}
	}

	return &item, nil
}

func (r *CatalogRepository) ListItems(ctx context.Context, catalogDomain domain.CatalogDomain, category string, activeOnly bool) ([]domain.CatalogItem, error) {
	table, err := r.tableForDomain(catalogDomain)
	if err != nil {
		return nil, err
	}

	dbExec := r.server.DB.Conn(ctx)
	var stmt string
	if catalogDomain == domain.ClinicalDomain {
		stmt = `
			SELECT id, category, code, name, COALESCE(description, ''), COALESCE(system_group, ''), is_active, version, created_at, updated_at, updated_by
			FROM platform.clinical_catalogs
			WHERE ($1 = '' OR category = $1)
			  AND ($2 = FALSE OR is_active = TRUE)
			ORDER BY code ASC
		`
	} else {
		stmt = fmt.Sprintf(`
			SELECT id, category, code, name, COALESCE(description, ''), metadata, is_active, version, created_at, updated_at, updated_by
			FROM %s
			WHERE ($1 = '' OR category = $1)
			  AND ($2 = FALSE OR is_active = TRUE)
			ORDER BY code ASC
		`, table)
	}

	rows, err := dbExec.Query(ctx, stmt, category, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to query master catalog items for domain=%s: %w", catalogDomain, err)
	}
	defer rows.Close()

	var items []domain.CatalogItem
	for rows.Next() {
		item, err := r.scanCatalogItem(catalogDomain, rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan catalog item: %w", err)
		}
		items = append(items, *item)
	}
	return items, nil
}

func (r *CatalogRepository) SearchItems(ctx context.Context, catalogDomain domain.CatalogDomain, query string) ([]domain.CatalogItem, error) {
	table, err := r.tableForDomain(catalogDomain)
	if err != nil {
		return nil, err
	}

	dbExec := r.server.DB.Conn(ctx)
	searchPattern := "%" + query + "%"
	var stmt string
	if catalogDomain == domain.ClinicalDomain {
		stmt = `
			SELECT id, category, code, name, COALESCE(description, ''), COALESCE(system_group, ''), is_active, version, created_at, updated_at, updated_by
			FROM platform.clinical_catalogs
			WHERE is_active = TRUE
			  AND ($1 = '' OR code ILIKE $1 OR name ILIKE $1 OR description ILIKE $1 OR system_group ILIKE $1)
			ORDER BY code ASC
			LIMIT 100
		`
	} else {
		stmt = fmt.Sprintf(`
			SELECT id, category, code, name, COALESCE(description, ''), metadata, is_active, version, created_at, updated_at, updated_by
			FROM %s
			WHERE is_active = TRUE
			  AND ($1 = '' OR code ILIKE $1 OR name ILIKE $1 OR description ILIKE $1)
			ORDER BY code ASC
			LIMIT 100
		`, table)
	}

	rows, err := dbExec.Query(ctx, stmt, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search catalog items for domain=%s: %w", catalogDomain, err)
	}
	defer rows.Close()

	var items []domain.CatalogItem
	for rows.Next() {
		item, err := r.scanCatalogItem(catalogDomain, rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search catalog item: %w", err)
		}
		items = append(items, *item)
	}
	return items, nil
}

func (r *CatalogRepository) GetItemByCode(ctx context.Context, catalogDomain domain.CatalogDomain, code string) (*domain.CatalogItem, error) {
	table, err := r.tableForDomain(catalogDomain)
	if err != nil {
		return nil, err
	}

	dbExec := r.server.DB.Conn(ctx)
	var stmt string
	if catalogDomain == domain.ClinicalDomain {
		stmt = `
			SELECT id, category, code, name, COALESCE(description, ''), COALESCE(system_group, ''), is_active, version, created_at, updated_at, updated_by
			FROM platform.clinical_catalogs
			WHERE code = $1
			LIMIT 1
		`
	} else {
		stmt = fmt.Sprintf(`
			SELECT id, category, code, name, COALESCE(description, ''), metadata, is_active, version, created_at, updated_at, updated_by
			FROM %s
			WHERE code = $1
			LIMIT 1
		`, table)
	}

	row := dbExec.QueryRow(ctx, stmt, code)
	return r.scanCatalogRow(catalogDomain, row)
}

func (r *CatalogRepository) GetItemByID(ctx context.Context, catalogDomain domain.CatalogDomain, id uuid.UUID) (*domain.CatalogItem, error) {
	table, err := r.tableForDomain(catalogDomain)
	if err != nil {
		return nil, err
	}

	dbExec := r.server.DB.Conn(ctx)
	var stmt string
	if catalogDomain == domain.ClinicalDomain {
		stmt = `
			SELECT id, category, code, name, COALESCE(description, ''), COALESCE(system_group, ''), is_active, version, created_at, updated_at, updated_by
			FROM platform.clinical_catalogs
			WHERE id = $1
			LIMIT 1
		`
	} else {
		stmt = fmt.Sprintf(`
			SELECT id, category, code, name, COALESCE(description, ''), metadata, is_active, version, created_at, updated_at, updated_by
			FROM %s
			WHERE id = $1
			LIMIT 1
		`, table)
	}

	row := dbExec.QueryRow(ctx, stmt, id)
	return r.scanCatalogRow(catalogDomain, row)
}

func (r *CatalogRepository) CreateItem(ctx context.Context, item *domain.CatalogItem, updatedBy uuid.UUID) (*domain.CatalogItem, error) {
	table, err := r.tableForDomain(item.Domain)
	if err != nil {
		return nil, err
	}
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}

	dbExec := r.server.DB.Conn(ctx)

	var (
		version   int
		createdAt time.Time
		updatedAt time.Time
	)

	if item.Domain == domain.ClinicalDomain {
		stmt := `
			INSERT INTO platform.clinical_catalogs (id, category, code, name, description, system_group, is_active, version, updated_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 1, $8)
			RETURNING version, created_at, updated_at
		`
		err = dbExec.QueryRow(ctx, stmt,
			item.ID, item.Category, item.Code, item.Name, item.Description, item.SystemGroup, item.IsActive, updatedBy.String(),
		).Scan(&version, &createdAt, &updatedAt)
	} else {
		var metaMap map[string]any
		if len(item.Metadata) > 0 {
			_ = json.Unmarshal(item.Metadata, &metaMap)
		}
		if metaMap == nil {
			metaMap = make(map[string]any)
		}
		if item.BasePrice > 0 {
			metaMap["basePrice"] = item.BasePrice
		}
		if metaBytes, err := json.Marshal(metaMap); err == nil {
			item.Metadata = metaBytes
		}

		stmt := fmt.Sprintf(`
			INSERT INTO %s (id, category, code, name, description, metadata, is_active, version, updated_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 1, $8)
			RETURNING version, created_at, updated_at
		`, table)
		err = dbExec.QueryRow(ctx, stmt,
			item.ID, item.Category, item.Code, item.Name, item.Description, item.Metadata, item.IsActive, updatedBy.String(),
		).Scan(&version, &createdAt, &updatedAt)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to insert catalog item into %s: %w", table, err)
	}

	item.Version = version
	item.CreatedAt = createdAt
	item.UpdatedAt = updatedAt
	item.UpdatedBy = &updatedBy

	return item, nil
}

func (r *CatalogRepository) UpdateItem(ctx context.Context, item *domain.CatalogItem, updatedBy uuid.UUID) (*domain.CatalogItem, error) {
	table, err := r.tableForDomain(item.Domain)
	if err != nil {
		return nil, err
	}

	dbExec := r.server.DB.Conn(ctx)

	var (
		newVersion int
		updatedAt  time.Time
	)

	if item.Domain == domain.ClinicalDomain {
		stmt := `
			UPDATE platform.clinical_catalogs
			SET category = $1,
			    code = $2,
			    name = $3,
			    description = $4,
			    system_group = $5,
			    is_active = $6,
			    version = version + 1,
			    updated_at = CURRENT_TIMESTAMP,
			    updated_by = $7
			WHERE id = $8 AND ($9 = 0 OR version = $9)
			RETURNING version, updated_at
		`
		err = dbExec.QueryRow(ctx, stmt,
			item.Category, item.Code, item.Name, item.Description, item.SystemGroup, item.IsActive, updatedBy.String(), item.ID, item.Version,
		).Scan(&newVersion, &updatedAt)
	} else {
		var metaMap map[string]any
		if len(item.Metadata) > 0 {
			_ = json.Unmarshal(item.Metadata, &metaMap)
		}
		if metaMap == nil {
			metaMap = make(map[string]any)
		}
		if item.BasePrice > 0 {
			metaMap["basePrice"] = item.BasePrice
		}
		if metaBytes, err := json.Marshal(metaMap); err == nil {
			item.Metadata = metaBytes
		}

		stmt := fmt.Sprintf(`
			UPDATE %s
			SET category = $1,
			    code = $2,
			    name = $3,
			    description = $4,
			    metadata = $5,
			    is_active = $6,
			    version = version + 1,
			    updated_at = CURRENT_TIMESTAMP,
			    updated_by = $7
			WHERE id = $8 AND ($9 = 0 OR version = $9)
			RETURNING version, updated_at
		`, table)
		err = dbExec.QueryRow(ctx, stmt,
			item.Category, item.Code, item.Name, item.Description, item.Metadata, item.IsActive, updatedBy.String(), item.ID, item.Version,
		).Scan(&newVersion, &updatedAt)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOptimisticLockingConflict
		}
		return nil, fmt.Errorf("failed to update catalog item in %s: %w", table, err)
	}

	item.Version = newVersion
	item.UpdatedAt = updatedAt
	item.UpdatedBy = &updatedBy

	return item, nil
}

// ─── Legacy DTO compatibility methods ──────────────────────────────────────────

func (r *CatalogRepository) GetSpecimenTypes(ctx context.Context) ([]domain.SpecimenType, error) {
	items, err := r.ListItems(ctx, domain.LabDomain, "specimen", true)
	if err != nil || len(items) == 0 {
		return []domain.SpecimenType{
			{ID: "spec-01", Code: "EDTA", Name: "Whole Blood (EDTA)", Container: "Purple Cap EDTA Tube", IsActive: true},
			{ID: "spec-02", Code: "SERUM", Name: "Serum (Clot Activator)", Container: "Red / Yellow SST Tube", IsActive: true},
			{ID: "spec-03", Code: "PLASMA", Name: "Plasma (Sodium Heparin)", Container: "Green Cap Tube", IsActive: true},
			{ID: "spec-04", Code: "URINE", Name: "Urine (Midstream)", Container: "Sterile Urine Cup", IsActive: true},
		}, nil
	}

	var res []domain.SpecimenType
	for _, it := range items {
		container := "Standard Container"
		if len(it.Metadata) > 0 {
			var m map[string]string
			if _ = json.Unmarshal(it.Metadata, &m); m["container"] != "" {
				container = m["container"]
			}
		}
		res = append(res, domain.SpecimenType{
			ID:        it.ID.String(),
			Code:      it.Code,
			Name:      it.Name,
			Container: container,
			IsActive:  it.IsActive,
		})
	}
	return res, nil
}

func (r *CatalogRepository) GetTestCategories(ctx context.Context) ([]domain.TestCategory, error) {
	items, err := r.ListItems(ctx, domain.LabDomain, "test_category", true)
	if err != nil || len(items) == 0 {
		return []domain.TestCategory{
			{ID: "cat-01", Code: "HEM", Name: "Hematology", Description: "Blood cell counts and coagulation testing", IsActive: true},
			{ID: "cat-02", Code: "CHM", Name: "Clinical Chemistry", Description: "Metabolic panels, enzymes, and electrolytes", IsActive: true},
		}, nil
	}

	var res []domain.TestCategory
	for _, it := range items {
		res = append(res, domain.TestCategory{
			ID:          it.ID.String(),
			Code:        it.Code,
			Name:        it.Name,
			Description: it.Description,
			IsActive:    it.IsActive,
		})
	}
	return res, nil
}

func (r *CatalogRepository) GetUnitsOfMeasure(ctx context.Context) ([]domain.UnitOfMeasure, error) {
	items, err := r.ListItems(ctx, domain.LabDomain, "uom", true)
	if err != nil || len(items) == 0 {
		return []domain.UnitOfMeasure{
			{ID: "uom-01", Code: "MG_DL", Name: "Milligrams per Deciliter", Symbol: "mg/dL", IsActive: true},
			{ID: "uom-02", Code: "MMOL_L", Name: "Millimoles per Liter", Symbol: "mmol/L", IsActive: true},
		}, nil
	}

	var res []domain.UnitOfMeasure
	for _, it := range items {
		sym := it.Description
		res = append(res, domain.UnitOfMeasure{
			ID:       it.ID.String(),
			Code:     it.Code,
			Name:     it.Name,
			Symbol:   sym,
			IsActive: it.IsActive,
		})
	}
	return res, nil
}

func (r *CatalogRepository) GetSpecialties(ctx context.Context) ([]domain.MedicalSpecialty, error) {
	items, err := r.ListItems(ctx, domain.ClinicalDomain, "specialty", true)
	if err != nil || len(items) == 0 {
		return []domain.MedicalSpecialty{
			{ID: "spec-med-01", Code: "GEN", Name: "General Practice / Internal Medicine"},
			{ID: "spec-med-02", Code: "PATH", Name: "Pathology & Laboratory Medicine"},
			{ID: "spec-med-03", Code: "RAD", Name: "Radiology & Diagnostic Imaging"},
		}, nil
	}

	var res []domain.MedicalSpecialty
	for _, it := range items {
		res = append(res, domain.MedicalSpecialty{
			ID:   it.ID.String(),
			Code: it.Code,
			Name: it.Name,
		})
	}
	return res, nil
}

func (r *CatalogRepository) SearchICD10(ctx context.Context, query string) ([]domain.ICD10Code, error) {
	items, err := r.SearchItems(ctx, domain.ClinicalDomain, query)
	if err != nil || len(items) == 0 {
		return []domain.ICD10Code{
			{Code: "E11.9", Description: "Type 2 diabetes mellitus without complications", Category: "Endocrine"},
			{Code: "I10", Description: "Essential (primary) hypertension", Category: "Cardiovascular"},
			{Code: "B54", Description: "Unspecified malaria", Category: "Infectious"},
		}, nil
	}

	var res []domain.ICD10Code
	for _, it := range items {
		cat := it.SystemGroup
		if cat == "" {
			cat = "General"
		}
		res = append(res, domain.ICD10Code{
			Code:        it.Code,
			Description: it.Name,
			Category:    cat,
		})
	}
	return res, nil
}
