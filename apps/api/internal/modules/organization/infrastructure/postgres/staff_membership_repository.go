package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type StaffMembershipRepository struct {
	server *server.Server
}

func NewStaffMembershipRepository(server *server.Server) *StaffMembershipRepository {
	return &StaffMembershipRepository{server: server}
}

func (r *StaffMembershipRepository) ListMembers(ctx context.Context, orgID uuid.UUID) ([]domain.StaffMemberDTO, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT m.id, m.organization_id, m.user_id::text, u.email, COALESCE(u.name, ''),
		       m.role, COALESCE(m.role_title, 'member'), CASE WHEN m.is_active THEN 'ACTIVE' ELSE 'INACTIVE' END, m.created_at
		FROM organization.organization_memberships m
		JOIN identity.users u ON u.id::text = m.user_id::text
		WHERE m.organization_id = $1
		ORDER BY m.created_at ASC
	`

	rows, err := dbExec.Query(ctx, stmt, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to query staff memberships: %w", err)
	}
	defer rows.Close()

	var members []domain.StaffMemberDTO
	for rows.Next() {
		var m domain.StaffMemberDTO
		err := rows.Scan(
			&m.MembershipID, &m.OrganizationID, &m.UserID, &m.Email, &m.FullName,
			&m.Role, &m.RoleTitle, &m.Status, &m.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan staff member row: %w", err)
		}
		members = append(members, m)
	}

	// Populate assigned branches & departments for each member
	for i := range members {
		branches, _ := r.listMemberBranches(ctx, members[i].MembershipID)
		depts, _ := r.listMemberDepartments(ctx, members[i].MembershipID)
		members[i].AssignedBranches = branches
		members[i].DepartmentAssignments = depts
	}

	return members, nil
}

func (r *StaffMembershipRepository) GetMemberByID(ctx context.Context, orgID, membershipID uuid.UUID) (*domain.StaffMemberDTO, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT m.id, m.organization_id, m.user_id::text, u.email, COALESCE(u.name, ''),
		       m.role, COALESCE(m.role_title, 'member'), CASE WHEN m.is_active THEN 'ACTIVE' ELSE 'INACTIVE' END, m.created_at
		FROM organization.organization_memberships m
		JOIN identity.users u ON u.id::text = m.user_id::text
		WHERE m.organization_id = $1 AND m.id = $2
		LIMIT 1
	`

	var m domain.StaffMemberDTO
	err := dbExec.QueryRow(ctx, stmt, orgID, membershipID).Scan(
		&m.MembershipID, &m.OrganizationID, &m.UserID, &m.Email, &m.FullName,
		&m.Role, &m.RoleTitle, &m.Status, &m.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrStaffMemberNotFound
		}
		return nil, fmt.Errorf("failed to query member id=%s: %w", membershipID, err)
	}

	m.AssignedBranches, _ = r.listMemberBranches(ctx, m.MembershipID)
	m.DepartmentAssignments, _ = r.listMemberDepartments(ctx, m.MembershipID)
	return &m, nil
}

func (r *StaffMembershipRepository) listMemberBranches(ctx context.Context, membershipID uuid.UUID) ([]domain.FacilityBranch, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT b.id, b.organization_id, COALESCE(b.facility_type_id, '00000000-0000-0000-0000-000000000000'::uuid), COALESCE(ft.code, ''), COALESCE(ft.name, ''), COALESCE(ft.category, ''),
		       b.code, b.name, b.is_headquarters, b.email, b.phone, b.address, b.city, b.state, b.lga, COALESCE(b.country, 'Nigeria'),
		       b.operating_hours, b.status, b.version, b.created_at, b.updated_at
		FROM organization.membership_branches mb
		JOIN organization.facility_branches b ON b.id = mb.facility_branch_id
		LEFT JOIN platform.facility_types ft ON ft.id = b.facility_type_id
		WHERE mb.membership_id = $1
	`
	rows, err := dbExec.Query(ctx, stmt, membershipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.FacilityBranch
	for rows.Next() {
		var b domain.FacilityBranch
		_ = rows.Scan(
			&b.ID, &b.OrganizationID, &b.FacilityTypeID, &b.FacilityTypeCode, &b.FacilityTypeName, &b.FacilityTypeCategory,
			&b.Code, &b.Name, &b.IsHeadquarters, &b.Email, &b.Phone, &b.Address, &b.City, &b.State, &b.LGA, &b.Country,
			&b.OperatingHours, &b.Status, &b.Version, &b.CreatedAt, &b.UpdatedAt,
		)
		list = append(list, b)
	}
	return list, nil
}

func (r *StaffMembershipRepository) listMemberDepartments(ctx context.Context, membershipID uuid.UUID) ([]domain.DepartmentalMembership, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, membership_id, facility_branch_id, department_code, created_at, created_by
		FROM organization.departmental_memberships
		WHERE membership_id = $1
	`
	rows, err := dbExec.Query(ctx, stmt, membershipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.DepartmentalMembership
	for rows.Next() {
		var d domain.DepartmentalMembership
		var cByStr *string
		_ = rows.Scan(&d.ID, &d.MembershipID, &d.FacilityBranchID, &d.DepartmentCode, &d.CreatedAt, &cByStr)
		if cByStr != nil {
			if parsed, pErr := uuid.Parse(*cByStr); pErr == nil {
				d.CreatedBy = &parsed
			}
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *StaffMembershipRepository) CountActiveMembers(ctx context.Context, orgID uuid.UUID) (int, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `SELECT COUNT(*) FROM organization.organization_memberships WHERE organization_id = $1 AND is_active = TRUE`
	var count int
	err := dbExec.QueryRow(ctx, stmt, orgID.String()).Scan(&count)
	return count, err
}

func (r *StaffMembershipRepository) CreateInvitation(ctx context.Context, invite *domain.StaffInvitation) (*domain.StaffInvitation, error) {
	dbExec := r.server.DB.Conn(ctx)
	if invite.ID == uuid.Nil {
		invite.ID = uuid.New()
	}

	var branchIDStr *string
	if invite.FacilityBranchID != nil {
		str := invite.FacilityBranchID.String()
		branchIDStr = &str
	}

	stmt := `
		INSERT INTO organization.staff_invitations (
			id, organization_id, facility_branch_id, email, role, role_title, invite_token_hash, status, expires_at, invited_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'PENDING', $8, $9)
		ON CONFLICT (organization_id, email) DO UPDATE SET
			facility_branch_id = EXCLUDED.facility_branch_id,
			role = EXCLUDED.role,
			role_title = EXCLUDED.role_title,
			invite_token_hash = EXCLUDED.invite_token_hash,
			status = 'PENDING',
			expires_at = EXCLUDED.expires_at,
			invited_by = EXCLUDED.invited_by,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		invite.ID, invite.OrganizationID, branchIDStr, invite.Email, invite.Role, invite.RoleTitle, invite.InviteTokenHash, invite.ExpiresAt, invite.InvitedBy,
	).Scan(&invite.ID, &invite.CreatedAt, &invite.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create staff invitation: %w", err)
	}

	invite.Status = "PENDING"
	return invite, nil
}

func (r *StaffMembershipRepository) ListInvitations(ctx context.Context, orgID uuid.UUID) ([]domain.StaffInvitation, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, facility_branch_id, email, role, role_title, invite_token_hash, status, expires_at, invited_by, created_at, updated_at
		FROM organization.staff_invitations
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := dbExec.Query(ctx, stmt, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list staff invitations: %w", err)
	}
	defer rows.Close()

	var list []domain.StaffInvitation
	for rows.Next() {
		var inv domain.StaffInvitation
		var branchStr *string
		err := rows.Scan(
			&inv.ID, &inv.OrganizationID, &branchStr, &inv.Email, &inv.Role, &inv.RoleTitle,
			&inv.InviteTokenHash, &inv.Status, &inv.ExpiresAt, &inv.InvitedBy, &inv.CreatedAt, &inv.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan staff invitation row: %w", err)
		}
		if branchStr != nil && *branchStr != "" {
			if parsed, pErr := uuid.Parse(*branchStr); pErr == nil {
				inv.FacilityBranchID = &parsed
			}
		}
		list = append(list, inv)
	}
	return list, nil
}

func (r *StaffMembershipRepository) GetInvitationByHash(ctx context.Context, hash string) (*domain.StaffInvitation, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, facility_branch_id, email, role, role_title, invite_token_hash, status, expires_at, invited_by, created_at, updated_at
		FROM organization.staff_invitations
		WHERE invite_token_hash = $1
		LIMIT 1
	`

	var inv domain.StaffInvitation
	var branchStr *string
	err := dbExec.QueryRow(ctx, stmt, hash).Scan(
		&inv.ID, &inv.OrganizationID, &branchStr, &inv.Email, &inv.Role, &inv.RoleTitle,
		&inv.InviteTokenHash, &inv.Status, &inv.ExpiresAt, &inv.InvitedBy, &inv.CreatedAt, &inv.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrStaffInvitationNotFound
		}
		return nil, fmt.Errorf("failed to query invitation by hash: %w", err)
	}

	if branchStr != nil && *branchStr != "" {
		if parsed, pErr := uuid.Parse(*branchStr); pErr == nil {
			inv.FacilityBranchID = &parsed
		}
	}

	return &inv, nil
}

func (r *StaffMembershipRepository) RevokeInvitation(ctx context.Context, orgID, inviteID uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `UPDATE organization.staff_invitations SET status = 'REVOKED', updated_at = CURRENT_TIMESTAMP WHERE organization_id = $1 AND id = $2`
	res, err := dbExec.Exec(ctx, stmt, orgID, inviteID)
	if err != nil {
		return fmt.Errorf("failed to revoke invitation id=%s: %w", inviteID, err)
	}
	if res.RowsAffected() == 0 {
		return domain.ErrStaffInvitationNotFound
	}
	return nil
}

func (r *StaffMembershipRepository) AcceptInvitation(ctx context.Context, inviteID uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `UPDATE organization.staff_invitations SET status = 'ACCEPTED', updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := dbExec.Exec(ctx, stmt, inviteID)
	return err
}

func (r *StaffMembershipRepository) CheckPendingInviteExists(ctx context.Context, orgID uuid.UUID, email string) (bool, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `SELECT EXISTS(SELECT 1 FROM organization.staff_invitations WHERE organization_id = $1 AND LOWER(email) = LOWER($2) AND status = 'PENDING' AND expires_at > CURRENT_TIMESTAMP)`
	var exists bool
	err := dbExec.QueryRow(ctx, stmt, orgID, email).Scan(&exists)
	return exists, err
}

func (r *StaffMembershipRepository) AssignBranch(ctx context.Context, membershipID, branchID uuid.UUID, actorID uuid.UUID) (*domain.MembershipBranch, error) {
	dbExec := r.server.DB.Conn(ctx)
	b := &domain.MembershipBranch{
		ID:               uuid.New(),
		MembershipID:     membershipID,
		FacilityBranchID: branchID,
		CreatedBy:        &actorID,
	}

	stmt := `
		INSERT INTO organization.membership_branches (id, membership_id, facility_branch_id, created_by)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (membership_id, facility_branch_id) DO NOTHING
		RETURNING created_at
	`

	err := dbExec.QueryRow(ctx, stmt, b.ID, b.MembershipID, b.FacilityBranchID, actorID.String()).Scan(&b.CreatedAt)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to assign branch: %w", err)
	}

	return b, nil
}

func (r *StaffMembershipRepository) RemoveBranchAssignment(ctx context.Context, membershipID, branchID uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `DELETE FROM organization.membership_branches WHERE membership_id = $1 AND facility_branch_id = $2`
	_, err := dbExec.Exec(ctx, stmt, membershipID, branchID)
	return err
}

func (r *StaffMembershipRepository) AssignDepartment(ctx context.Context, membershipID, branchID uuid.UUID, deptCode string, actorID uuid.UUID) (*domain.DepartmentalMembership, error) {
	dbExec := r.server.DB.Conn(ctx)
	d := &domain.DepartmentalMembership{
		ID:               uuid.New(),
		MembershipID:     membershipID,
		FacilityBranchID: branchID,
		DepartmentCode:   deptCode,
		CreatedBy:        &actorID,
	}

	stmt := `
		INSERT INTO organization.departmental_memberships (id, membership_id, facility_branch_id, department_code, created_by)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (membership_id, facility_branch_id, department_code) DO NOTHING
		RETURNING created_at
	`

	err := dbExec.QueryRow(ctx, stmt, d.ID, d.MembershipID, d.FacilityBranchID, d.DepartmentCode, actorID.String()).Scan(&d.CreatedAt)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to assign department: %w", err)
	}

	return d, nil
}

func (r *StaffMembershipRepository) RemoveDepartmentAssignment(ctx context.Context, membershipID, branchID uuid.UUID, deptCode string) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `DELETE FROM organization.departmental_memberships WHERE membership_id = $1 AND facility_branch_id = $2 AND department_code = $3`
	_, err := dbExec.Exec(ctx, stmt, membershipID, branchID, deptCode)
	return err
}

func (r *StaffMembershipRepository) UpdateMemberRole(ctx context.Context, orgID, membershipID uuid.UUID, role, roleTitle string, actorID uuid.UUID) (*domain.StaffMemberDTO, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE organization.organization_memberships
		SET role = $1, role_title = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND organization_id = $4
	`
	res, err := dbExec.Exec(ctx, stmt, role, roleTitle, membershipID.String(), orgID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to update member role: %w", err)
	}
	if res.RowsAffected() == 0 {
		return nil, domain.ErrStaffMemberNotFound
	}

	return r.GetMemberByID(ctx, orgID, membershipID)
}

func (r *StaffMembershipRepository) DirectCreateMember(
	ctx context.Context,
	orgID uuid.UUID,
	fullName, email, passwordHash, role, roleTitle string,
	branchIDs []uuid.UUID,
	actorID uuid.UUID,
) (*domain.StaffMemberDTO, error) {
	dbExec := r.server.DB.Conn(ctx)
	cleanEmail := strings.ToLower(strings.TrimSpace(email))

	// 1. Check if user already exists in identity.users
	var userID string
	var userExists bool
	err := dbExec.QueryRow(ctx, `SELECT id FROM identity.users WHERE email = $1`, cleanEmail).Scan(&userID)
	if err == nil && userID != "" {
		userExists = true
	}

	if !userExists {
		newUserID := uuid.New().String()
		err = dbExec.QueryRow(ctx, `
			INSERT INTO identity.users (id, name, email, email_verified, is_platform_admin)
			VALUES ($1, $2, $3, TRUE, FALSE)
			RETURNING id
		`, newUserID, fullName, cleanEmail).Scan(&userID)
		if err != nil {
			return nil, fmt.Errorf("failed to create user record: %w", err)
		}

		// Insert credential record in identity.credentials
		credID := uuid.New().String()
		_, err = dbExec.Exec(ctx, `
			INSERT INTO identity.credentials (id, account_id, auth_provider, user_id, password_hash)
			VALUES ($1, $2, 'credential', $3, $4)
			ON CONFLICT (account_id) DO UPDATE SET password_hash = EXCLUDED.password_hash
		`, credID, cleanEmail, userID, passwordHash)
		if err != nil {
			return nil, fmt.Errorf("failed to create credential record: %w", err)
		}
	} else if passwordHash != "" {
		// Update credential password if provided
		_, _ = dbExec.Exec(ctx, `
			INSERT INTO identity.credentials (id, account_id, auth_provider, user_id, password_hash, created_at, updated_at)
			VALUES ($1, $2, 'credential', $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			ON CONFLICT (user_id, auth_provider) DO UPDATE SET password_hash = EXCLUDED.password_hash, updated_at = CURRENT_TIMESTAMP
		`, uuid.New().String(), cleanEmail, userID, passwordHash)
	}

	// 2. Insert or update organization membership
	var memID uuid.UUID
	memUUID := uuid.New()
	err = dbExec.QueryRow(ctx, `
		INSERT INTO organization.organization_memberships (id, organization_id, user_id, role, role_title, is_active, updated_at)
		VALUES ($1, $2, $3, $4, $5, TRUE, CURRENT_TIMESTAMP)
		ON CONFLICT (organization_id, user_id) DO UPDATE SET 
			role = EXCLUDED.role, 
			role_title = EXCLUDED.role_title, 
			is_active = TRUE,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`, memUUID, orgID.String(), userID, role, roleTitle).Scan(&memID)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization membership: %w", err)
	}

	// 3. Assign branches in organization.membership_branches
	for _, bID := range branchIDs {
		if bID != uuid.Nil {
			var actorIDParam *string
			if actorID != uuid.Nil {
				actStr := actorID.String()
				actorIDParam = &actStr
			}
			_, _ = dbExec.Exec(ctx, `
				INSERT INTO organization.membership_branches (id, membership_id, facility_branch_id, created_by)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (membership_id, facility_branch_id) DO NOTHING
			`, uuid.New(), memID, bID, actorIDParam)
		}
	}

	return r.GetMemberByID(ctx, orgID, memID)
}
