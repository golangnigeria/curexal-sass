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

type FacilityBranchRepository struct {
	server *server.Server
}

func NewFacilityBranchRepository(server *server.Server) *FacilityBranchRepository {
	return &FacilityBranchRepository{server: server}
}

func (r *FacilityBranchRepository) ListBranches(ctx context.Context, orgID uuid.UUID) ([]domain.FacilityBranch, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT b.id, b.organization_id, b.facility_type_id, ft.code AS ft_code, ft.name AS ft_name, ft.category AS ft_cat,
		       b.code, b.name, b.is_headquarters, b.email, b.phone, b.address, b.city, b.state, b.lga, COALESCE(b.country, 'Nigeria'),
		       b.operating_hours, b.status, b.version, b.created_at, b.updated_at, b.updated_by
		FROM organization.facility_branches b
		JOIN platform.facility_types ft ON ft.id = b.facility_type_id
		WHERE b.organization_id = $1
		ORDER BY b.is_headquarters DESC, b.name ASC
	`

	rows, err := dbExec.Query(ctx, stmt, orgID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query organization facility branches: %w", err)
	}
	defer rows.Close()

	var branches []domain.FacilityBranch
	for rows.Next() {
		var (
			b            domain.FacilityBranch
			updatedByStr *string
		)
		err := rows.Scan(
			&b.ID, &b.OrganizationID, &b.FacilityTypeID, &b.FacilityTypeCode, &b.FacilityTypeName, &b.FacilityTypeCategory,
			&b.Code, &b.Name, &b.IsHeadquarters, &b.Email, &b.Phone, &b.Address, &b.City, &b.State, &b.LGA, &b.Country,
			&b.OperatingHours, &b.Status, &b.Version, &b.CreatedAt, &b.UpdatedAt, &updatedByStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan facility branch row: %w", err)
		}
		if updatedByStr != nil && *updatedByStr != "" {
			if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
				b.UpdatedBy = &parsed
			}
		}
		branches = append(branches, b)
	}

	return branches, nil
}

func (r *FacilityBranchRepository) GetBranchByID(ctx context.Context, orgID, branchID uuid.UUID) (*domain.FacilityBranch, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT b.id, b.organization_id, b.facility_type_id, ft.code AS ft_code, ft.name AS ft_name, ft.category AS ft_cat,
		       b.code, b.name, b.is_headquarters, b.email, b.phone, b.address, b.city, b.state, b.lga, COALESCE(b.country, 'Nigeria'),
		       b.operating_hours, b.status, b.version, b.created_at, b.updated_at, b.updated_by
		FROM organization.facility_branches b
		JOIN platform.facility_types ft ON ft.id = b.facility_type_id
		WHERE b.organization_id = $1 AND b.id = $2
		LIMIT 1
	`

	var (
		b            domain.FacilityBranch
		updatedByStr *string
	)
	err := dbExec.QueryRow(ctx, stmt, orgID.String(), branchID.String()).Scan(
		&b.ID, &b.OrganizationID, &b.FacilityTypeID, &b.FacilityTypeCode, &b.FacilityTypeName, &b.FacilityTypeCategory,
		&b.Code, &b.Name, &b.IsHeadquarters, &b.Email, &b.Phone, &b.Address, &b.City, &b.State, &b.LGA, &b.Country,
		&b.OperatingHours, &b.Status, &b.Version, &b.CreatedAt, &b.UpdatedAt, &updatedByStr,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrFacilityBranchNotFound
		}
		return nil, fmt.Errorf("failed to query facility branch id=%s: %w", branchID, err)
	}

	if updatedByStr != nil && *updatedByStr != "" {
		if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
			b.UpdatedBy = &parsed
		}
	}

	return &b, nil
}

func (r *FacilityBranchRepository) GetBranchByCode(ctx context.Context, orgID uuid.UUID, code string) (*domain.FacilityBranch, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT b.id, b.organization_id, b.facility_type_id, ft.code AS ft_code, ft.name AS ft_name, ft.category AS ft_cat,
		       b.code, b.name, b.is_headquarters, b.email, b.phone, b.address, b.city, b.state, b.lga, COALESCE(b.country, 'Nigeria'),
		       b.operating_hours, b.status, b.version, b.created_at, b.updated_at, b.updated_by
		FROM organization.facility_branches b
		JOIN platform.facility_types ft ON ft.id = b.facility_type_id
		WHERE b.organization_id = $1 AND b.code = $2
		LIMIT 1
	`

	var (
		b            domain.FacilityBranch
		updatedByStr *string
	)
	err := dbExec.QueryRow(ctx, stmt, orgID.String(), code).Scan(
		&b.ID, &b.OrganizationID, &b.FacilityTypeID, &b.FacilityTypeCode, &b.FacilityTypeName, &b.FacilityTypeCategory,
		&b.Code, &b.Name, &b.IsHeadquarters, &b.Email, &b.Phone, &b.Address, &b.City, &b.State, &b.LGA, &b.Country,
		&b.OperatingHours, &b.Status, &b.Version, &b.CreatedAt, &b.UpdatedAt, &updatedByStr,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrFacilityBranchNotFound
		}
		return nil, fmt.Errorf("failed to query facility branch code=%s: %w", code, err)
	}

	if updatedByStr != nil && *updatedByStr != "" {
		if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
			b.UpdatedBy = &parsed
		}
	}

	return &b, nil
}

func (r *FacilityBranchRepository) CreateBranch(ctx context.Context, branch *domain.FacilityBranch, actorID uuid.UUID) (*domain.FacilityBranch, error) {
	dbExec := r.server.DB.Conn(ctx)
	if branch.ID == uuid.Nil {
		branch.ID = uuid.New()
	}

	countryVal := "Nigeria"
	if branch.Country != "" {
		countryVal = branch.Country
	}

	hoursJSON := []byte("{}")
	if len(branch.OperatingHours) > 0 {
		hoursJSON = branch.OperatingHours
	}

	stmt := `
		INSERT INTO organization.facility_branches (
			id, organization_id, facility_type_id, code, name, is_headquarters,
			email, phone, address, city, state, lga, country, operating_hours, status, version, updated_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14::jsonb, 'ACTIVE', 1, $15)
		RETURNING version, created_at, updated_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		branch.ID, branch.OrganizationID, branch.FacilityTypeID, branch.Code, branch.Name, branch.IsHeadquarters,
		branch.Email, branch.Phone, branch.Address, branch.City, branch.State, branch.LGA, countryVal, string(hoursJSON), actorID.String(),
	).Scan(&branch.Version, &branch.CreatedAt, &branch.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "uk_facility_single_headquarters" {
				return nil, domain.ErrHeadquartersConflict
			}
			if pgErr.ConstraintName == "uk_facility_org_code" {
				return nil, domain.ErrDuplicateBranchCode
			}
		}
		return nil, fmt.Errorf("failed to create facility branch: %w", err)
	}

	branch.UpdatedBy = &actorID
	return r.GetBranchByID(ctx, branch.OrganizationID, branch.ID)
}

func (r *FacilityBranchRepository) UpdateBranch(ctx context.Context, branch *domain.FacilityBranch, actorID uuid.UUID) (*domain.FacilityBranch, error) {
	dbExec := r.server.DB.Conn(ctx)

	hoursJSON := []byte("{}")
	if len(branch.OperatingHours) > 0 {
		hoursJSON = branch.OperatingHours
	}

	stmt := `
		UPDATE organization.facility_branches
		SET name = $1,
		    is_headquarters = $2,
		    email = $3,
		    phone = $4,
		    address = $5,
		    city = $6,
		    state = $7,
		    lga = $8,
		    operating_hours = $9::jsonb,
		    status = COALESCE($10, status),
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $11
		WHERE id = $12 AND organization_id = $13 AND version = $14
		RETURNING version, updated_at
	`

	var (
		newVersion int
		updatedAt  time.Time
	)
	err := dbExec.QueryRow(ctx, stmt,
		branch.Name, branch.IsHeadquarters, branch.Email, branch.Phone, branch.Address, branch.City, branch.State, branch.LGA,
		string(hoursJSON), branch.Status, actorID.String(), branch.ID, branch.OrganizationID, branch.Version,
	).Scan(&newVersion, &updatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "uk_facility_single_headquarters" {
				return nil, domain.ErrHeadquartersConflict
			}
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOptimisticLockingConflict
		}
		return nil, fmt.Errorf("failed to update facility branch id=%s: %w", branch.ID, err)
	}

	return r.GetBranchByID(ctx, branch.OrganizationID, branch.ID)
}

func (r *FacilityBranchRepository) DeactivateBranch(ctx context.Context, orgID, branchID uuid.UUID, actorID uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE organization.facility_branches
		SET status = 'INACTIVE',
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $1
		WHERE id = $2 AND organization_id = $3
	`
	res, err := dbExec.Exec(ctx, stmt, actorID.String(), branchID, orgID)
	if err != nil {
		return fmt.Errorf("failed to deactivate facility branch id=%s: %w", branchID, err)
	}
	if res.RowsAffected() == 0 {
		return domain.ErrFacilityBranchNotFound
	}
	return nil
}

func (r *FacilityBranchRepository) CountActiveBranches(ctx context.Context, orgID uuid.UUID) (int, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `SELECT COUNT(*) FROM organization.facility_branches WHERE organization_id = $1 AND status = 'ACTIVE'`
	var count int
	err := dbExec.QueryRow(ctx, stmt, orgID).Scan(&count)
	return count, err
}

func (r *FacilityBranchRepository) HasActiveHeadquarters(ctx context.Context, orgID uuid.UUID) (bool, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `SELECT EXISTS(SELECT 1 FROM organization.facility_branches WHERE organization_id = $1 AND is_headquarters = TRUE AND status = 'ACTIVE')`
	var exists bool
	err := dbExec.QueryRow(ctx, stmt, orgID).Scan(&exists)
	return exists, err
}

func (r *FacilityBranchRepository) CheckFacilityTypeActive(ctx context.Context, facilityTypeID uuid.UUID) (bool, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `SELECT EXISTS(SELECT 1 FROM platform.facility_types WHERE id = $1 AND is_active = TRUE)`
	var active bool
	err := dbExec.QueryRow(ctx, stmt, facilityTypeID).Scan(&active)
	return active, err
}
