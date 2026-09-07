package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
		       b.code, COALESCE(b.slug, b.code), b.name, b.is_headquarters, b.email, b.phone, b.address, b.city, b.state, b.lga, COALESCE(b.country, 'Nigeria'),
		       b.operating_hours, COALESCE(b.theme_branding, '{}'::jsonb), b.status, b.version, b.created_at, b.updated_at, b.updated_by
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
			&b.Code, &b.Slug, &b.Name, &b.IsHeadquarters, &b.Email, &b.Phone, &b.Address, &b.City, &b.State, &b.LGA, &b.Country,
			&b.OperatingHours, &b.ThemeBranding, &b.Status, &b.Version, &b.CreatedAt, &b.UpdatedAt, &updatedByStr,
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
		       b.code, COALESCE(b.slug, b.code), b.name, b.is_headquarters, b.email, b.phone, b.address, b.city, b.state, b.lga, COALESCE(b.country, 'Nigeria'),
		       b.operating_hours, COALESCE(b.theme_branding, '{}'::jsonb), b.status, b.version, b.created_at, b.updated_at, b.updated_by
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
		&b.Code, &b.Slug, &b.Name, &b.IsHeadquarters, &b.Email, &b.Phone, &b.Address, &b.City, &b.State, &b.LGA, &b.Country,
		&b.OperatingHours, &b.ThemeBranding, &b.Status, &b.Version, &b.CreatedAt, &b.UpdatedAt, &updatedByStr,
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
		       b.code, COALESCE(b.slug, b.code), b.name, b.is_headquarters, b.email, b.phone, b.address, b.city, b.state, b.lga, COALESCE(b.country, 'Nigeria'),
		       b.operating_hours, COALESCE(b.theme_branding, '{}'::jsonb), b.status, b.version, b.created_at, b.updated_at, b.updated_by
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
		&b.Code, &b.Slug, &b.Name, &b.IsHeadquarters, &b.Email, &b.Phone, &b.Address, &b.City, &b.State, &b.LGA, &b.Country,
		&b.OperatingHours, &b.ThemeBranding, &b.Status, &b.Version, &b.CreatedAt, &b.UpdatedAt, &updatedByStr,
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

func (r *FacilityBranchRepository) GetBranchBySlug(ctx context.Context, orgID uuid.UUID, slug string) (*domain.FacilityBranch, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT b.id, b.organization_id, b.facility_type_id, ft.code AS ft_code, ft.name AS ft_name, ft.category AS ft_cat,
		       b.code, COALESCE(b.slug, b.code), b.name, b.is_headquarters, b.email, b.phone, b.address, b.city, b.state, b.lga, COALESCE(b.country, 'Nigeria'),
		       b.operating_hours, COALESCE(b.theme_branding, '{}'::jsonb), b.status, b.version, b.created_at, b.updated_at, b.updated_by
		FROM organization.facility_branches b
		JOIN platform.facility_types ft ON ft.id = b.facility_type_id
		WHERE b.organization_id = $1 AND (b.slug = $2 OR b.code = $2)
		LIMIT 1
	`

	var (
		b            domain.FacilityBranch
		updatedByStr *string
	)
	err := dbExec.QueryRow(ctx, stmt, orgID.String(), slug).Scan(
		&b.ID, &b.OrganizationID, &b.FacilityTypeID, &b.FacilityTypeCode, &b.FacilityTypeName, &b.FacilityTypeCategory,
		&b.Code, &b.Slug, &b.Name, &b.IsHeadquarters, &b.Email, &b.Phone, &b.Address, &b.City, &b.State, &b.LGA, &b.Country,
		&b.OperatingHours, &b.ThemeBranding, &b.Status, &b.Version, &b.CreatedAt, &b.UpdatedAt, &updatedByStr,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrFacilityBranchNotFound
		}
		return nil, fmt.Errorf("failed to query facility branch slug=%s: %w", slug, err)
	}

	if updatedByStr != nil && *updatedByStr != "" {
		if parsed, pErr := uuid.Parse(*updatedByStr); pErr == nil {
			b.UpdatedBy = &parsed
		}
	}

	return &b, nil
}

func getMaxBranchesForPlanStatic(plan string) int {
	switch strings.ToLower(plan) {
	case "smart":
		return 3
	case "optimize":
		return 10
	case "pro":
		return 25
	case "enterprise":
		return 1000
	default:
		return 1
	}
}

func (r *FacilityBranchRepository) CreateBranch(ctx context.Context, branch *domain.FacilityBranch, actorID uuid.UUID) (*domain.FacilityBranch, error) {
	if branch.ID == uuid.Nil {
		branch.ID = uuid.New()
	}
	if branch.Slug == "" {
		branch.Slug = branch.Code
	}

	countryVal := "Nigeria"
	if branch.Country != "" {
		countryVal = branch.Country
	}

	hoursJSON := []byte("{}")
	if len(branch.OperatingHours) > 0 {
		hoursJSON = branch.OperatingHours
	}

	tx, errTx := r.server.DB.Pool.Begin(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("failed to begin tx for create branch: %w", errTx)
	}
	defer tx.Rollback(ctx)

	// 1. Row-level lock on organization to serialize quota evaluation and prevent race conditions
	var orgPlan string
	errLock := tx.QueryRow(ctx, `SELECT plan FROM organization.organizations WHERE id = $1 FOR UPDATE`, branch.OrganizationID.String()).Scan(&orgPlan)
	if errLock != nil {
		return nil, fmt.Errorf("failed to lock organization row for quota check: %w", errLock)
	}

	// 2. Count active branches inside the serialized transaction
	var activeCount int
	errCount := tx.QueryRow(ctx, `SELECT COUNT(*) FROM organization.facility_branches WHERE organization_id = $1 AND status = 'ACTIVE'`, branch.OrganizationID.String()).Scan(&activeCount)
	if errCount != nil {
		return nil, fmt.Errorf("failed to count active branches: %w", errCount)
	}

	maxBranches := getMaxBranchesForPlanStatic(orgPlan)
	if activeCount >= maxBranches {
		return nil, domain.ErrMaxBranchesExceeded
	}

	stmt := `
		INSERT INTO organization.facility_branches (
			id, organization_id, facility_type_id, code, slug, name, is_headquarters,
			email, phone, address, city, state, lga, country, operating_hours, status, version, updated_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15::jsonb, 'ACTIVE', 1, $16)
		RETURNING version, created_at, updated_at
	`

	err := tx.QueryRow(ctx, stmt,
		branch.ID, branch.OrganizationID, branch.FacilityTypeID, branch.Code, branch.Slug, branch.Name, branch.IsHeadquarters,
		branch.Email, branch.Phone, branch.Address, branch.City, branch.State, branch.LGA, countryVal, string(hoursJSON), actorID.String(),
	).Scan(&branch.Version, &branch.CreatedAt, &branch.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "uk_facility_org_code" || pgErr.ConstraintName == "uk_facility_branches_org_slug" {
				return nil, domain.ErrDuplicateBranchCode
			}
		}
		return nil, fmt.Errorf("failed to create facility branch: %w", err)
	}

	if errCommit := tx.Commit(ctx); errCommit != nil {
		return nil, fmt.Errorf("failed to commit create branch tx: %w", errCommit)
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

	themeJSON := []byte("{}")
	if len(branch.ThemeBranding) > 0 {
		themeJSON = branch.ThemeBranding
	}

	stmt := `
		UPDATE organization.facility_branches
		SET name = $1,
		    slug = COALESCE(NULLIF($2, ''), slug),
		    is_headquarters = $3,
		    email = $4,
		    phone = $5,
		    address = $6,
		    city = $7,
		    state = $8,
		    lga = $9,
		    operating_hours = $10::jsonb,
		    theme_branding = $11::jsonb,
		    status = COALESCE($12, status),
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $13
		WHERE id = $14 AND organization_id = $15 AND version = $16
		RETURNING version, updated_at
	`

	var (
		newVersion int
		updatedAt  time.Time
	)
	err := dbExec.QueryRow(ctx, stmt,
		branch.Name, branch.Slug, branch.IsHeadquarters, branch.Email, branch.Phone, branch.Address, branch.City, branch.State, branch.LGA,
		string(hoursJSON), string(themeJSON), branch.Status, actorID.String(), branch.ID, branch.OrganizationID, branch.Version,
	).Scan(&newVersion, &updatedAt)

	if err != nil {

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

func (r *FacilityBranchRepository) SetHeadquarters(ctx context.Context, orgID, branchID, actorID uuid.UUID) error {
	// In a transaction, unset previous HQ and set new HQ
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for set headquarters: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Unset all active HQ in this org
	unsetStmt := `
		UPDATE organization.facility_branches
		SET is_headquarters = FALSE,
		    updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND is_headquarters = TRUE
	`
	if _, err := tx.Exec(ctx, unsetStmt, orgID.String()); err != nil {
		return fmt.Errorf("failed to unset previous headquarters: %w", err)
	}

	// 2. Set target branch as HQ
	setStmt := `
		UPDATE organization.facility_branches
		SET is_headquarters = TRUE,
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP,
		    updated_by = $1
		WHERE id = $2 AND organization_id = $3
	`
	res, err := tx.Exec(ctx, setStmt, actorID.String(), branchID.String(), orgID.String())
	if err != nil {
		return fmt.Errorf("failed to set new headquarters id=%s: %w", branchID, err)
	}
	if res.RowsAffected() == 0 {
		return domain.ErrFacilityBranchNotFound
	}

	return tx.Commit(ctx)
}

func (r *FacilityBranchRepository) GetFacilityTypeByCode(ctx context.Context, code string) (*domain.FacilityType, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, code, name, category, COALESCE(icon_key, 'Building2'), description, is_active, version, created_at, updated_at
		FROM platform.facility_types
		WHERE (code = $1 OR id::text = $1 OR LOWER(name) = LOWER($1)) AND is_active = TRUE
		LIMIT 1
	`

	var ft domain.FacilityType
	err := dbExec.QueryRow(ctx, stmt, code).Scan(
		&ft.ID, &ft.Code, &ft.Name, &ft.Category, &ft.IconKey, &ft.Description,
		&ft.IsActive, &ft.Version, &ft.CreatedAt, &ft.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInvalidFacilityType
		}
		return nil, fmt.Errorf("failed to get facility type by code=%s: %w", code, err)
	}

	return &ft, nil
}

func (r *FacilityBranchRepository) ListFacilityTypes(ctx context.Context) ([]domain.FacilityType, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, code, name, category, COALESCE(icon_key, 'Building2'), description, is_active, version, created_at, updated_at
		FROM platform.facility_types
		WHERE is_active = TRUE
		ORDER BY name ASC
	`

	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query facility types: %w", err)
	}
	defer rows.Close()

	var types []domain.FacilityType
	for rows.Next() {
		var ft domain.FacilityType
		if err := rows.Scan(
			&ft.ID, &ft.Code, &ft.Name, &ft.Category, &ft.IconKey, &ft.Description,
			&ft.IsActive, &ft.Version, &ft.CreatedAt, &ft.UpdatedAt,
		); err == nil {
			types = append(types, ft)
		}
	}

	return types, nil
}

func (r *FacilityBranchRepository) VerifyUserFacilityAccess(ctx context.Context, orgID, branchID uuid.UUID, userID string, isOrgAdmin bool) (bool, error) {
	dbExec := r.server.DB.Conn(ctx)

	// 1. Verify branch exists, belongs to target organization, and is ACTIVE
	branchCheckStmt := `
		SELECT EXISTS (
			SELECT 1 FROM organization.facility_branches
			WHERE id = $1 AND organization_id = $2 AND status = 'ACTIVE'
		)
	`
	var branchActive bool
	if err := dbExec.QueryRow(ctx, branchCheckStmt, branchID.String(), orgID.String()).Scan(&branchActive); err != nil {
		return false, fmt.Errorf("failed to verify facility branch active status: %w", err)
	}
	if !branchActive {
		return false, nil
	}

	// 2. Organization governance administrators (owner / org_admin) have access to all branches within their organization
	if isOrgAdmin {
		return true, nil
	}

	// 3. For clinical and operational staff, verify explicit branch assignment in organization.membership_branches
	assignmentStmt := `
		SELECT EXISTS (
			SELECT 1
			FROM organization.membership_branches mb
			JOIN organization.organization_memberships om ON om.id = mb.membership_id
			WHERE om.user_id = $1::uuid
			  AND om.organization_id = $2
			  AND mb.facility_branch_id = $3
			  AND om.is_active = TRUE
		)
	`
	var isAssigned bool
	if err := dbExec.QueryRow(ctx, assignmentStmt, userID, orgID.String(), branchID.String()).Scan(&isAssigned); err != nil {
		return false, fmt.Errorf("failed to verify staff facility branch assignment: %w", err)
	}

	return isAssigned, nil
}

func (r *FacilityBranchRepository) IsUserAssignedToBranch(ctx context.Context, userID, orgID, branchID uuid.UUID) (bool, error) {
	return r.VerifyUserFacilityAccess(ctx, orgID, branchID, userID.String(), false)
}

