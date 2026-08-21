package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	modelUser "github.com/golangnigeria/curexal/internal/modules/identity/model/user"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	server *server.Server
}

func NewUserRepository(s *server.Server) *UserRepository {
	return &UserRepository{server: s}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	db := r.server.DB.Conn(ctx)
	u := &model.User{}
	err := db.QueryRow(ctx, `
		SELECT id, name, email, email_verified, avatar_url, is_platform_admin, platform_role, COALESCE(failed_login_attempts, 0), locked_until 
		FROM identity.users 
		WHERE email = @email
	`, pgx.NamedArgs{"email": email}).Scan(&u.ID, &u.Name, &u.Email, &u.EmailVerified, &u.Image, &u.IsPlatformAdmin, &u.PlatformRole, &u.FailedLoginAttempts, &u.LockedUntil)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	db := r.server.DB.Conn(ctx)
	u := &model.User{}
	err := db.QueryRow(ctx, `
		SELECT id, name, email, email_verified, avatar_url, is_platform_admin, platform_role, COALESCE(failed_login_attempts, 0), locked_until 
		FROM identity.users 
		WHERE id = @id
	`, pgx.NamedArgs{"id": id}).Scan(&u.ID, &u.Name, &u.Email, &u.EmailVerified, &u.Image, &u.IsPlatformAdmin, &u.PlatformRole, &u.FailedLoginAttempts, &u.LockedUntil)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, u *model.User, passwordHash string) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `
		INSERT INTO identity.users (id, name, email, email_verified, avatar_url, is_platform_admin, platform_role)
		VALUES (@id, @name, @email, @emailVerified, @image, @isPlatformAdmin, @platformRole)
		ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email
	`, pgx.NamedArgs{
		"id":              u.ID,
		"name":            u.Name,
		"email":           u.Email,
		"emailVerified":   u.EmailVerified,
		"image":           u.Image,
		"isPlatformAdmin": u.IsPlatformAdmin,
		"platformRole":    u.PlatformRole,
	})
	if err != nil {
		return err
	}

	if passwordHash != "" {
		accountID := uuid.New().String()
		_, errAccount := db.Exec(ctx, `
			INSERT INTO identity.credentials (id, account_id, auth_provider, user_id, password_hash)
			VALUES (@id, @accountId, 'credential', @userId, @password)
		`, pgx.NamedArgs{
			"id":        accountID,
			"accountId": strings.ToLower(strings.TrimSpace(u.Email)),
			"userId":    u.ID,
			"password":  passwordHash,
		})
		if errAccount != nil {
			return errAccount
		}
	}
	return nil
}

type VerificationTokenRecord struct {
	ID        string
	UserID    string
	Token     string
	TokenType string
	Metadata  map[string]interface{}
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (r *UserRepository) CreateVerificationTokenWithMetadata(ctx context.Context, userID, token, tokenType string, metadata map[string]interface{}, expiresAt time.Time) error {
	db := r.server.DB.Conn(ctx)
	metaJSON, err := json.Marshal(metadata)
	if err != nil {
		metaJSON = []byte("{}")
	}
	_, err = db.Exec(ctx, `
		INSERT INTO identity.verification_tokens (id, user_id, token, token_type, metadata, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, uuid.New().String(), userID, token, tokenType, metaJSON, expiresAt)
	return err
}

func (r *UserRepository) CreateVerificationToken(ctx context.Context, token, email, tokenType string, expiresAt time.Time) error {
	db := r.server.DB.Conn(ctx)
	var userID string
	err := db.QueryRow(ctx, `SELECT id::text FROM identity.users WHERE LOWER(email) = LOWER($1)`, strings.TrimSpace(email)).Scan(&userID)
	if err != nil {
		return err
	}
	return r.CreateVerificationTokenWithMetadata(ctx, userID, token, tokenType, map[string]interface{}{"email": email}, expiresAt)
}

func (r *UserRepository) CreateSession(ctx context.Context, sess *model.Session) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `
		INSERT INTO identity.sessions (id, user_id, token, expires_at)
		VALUES (@id, @userID, @token, @expiresAt)
	`, pgx.NamedArgs{
		"id":        sess.ID,
		"userID":    sess.UserID,
		"token":     sess.Token,
		"expiresAt": sess.ExpiresAt,
	})
	return err
}

func (r *UserRepository) GetSessionByToken(ctx context.Context, token string) (*model.Session, error) {
	db := r.server.DB.Conn(ctx)
	sess := &model.Session{}
	err := db.QueryRow(ctx, `SELECT id, user_id::text, token, expires_at FROM identity.sessions WHERE token = @token`, pgx.NamedArgs{"token": token}).
		Scan(&sess.ID, &sess.UserID, &sess.Token, &sess.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (r *UserRepository) GetVerificationTokenRecord(ctx context.Context, token string) (*VerificationTokenRecord, error) {
	db := r.server.DB.Conn(ctx)
	rec := &VerificationTokenRecord{}
	var metaBytes []byte
	cleanToken := strings.TrimSpace(token)
	err := db.QueryRow(ctx, `
		SELECT id, user_id::text, token, token_type, COALESCE(metadata, '{}'::jsonb), expires_at, created_at
		FROM identity.verification_tokens
		WHERE UPPER(token) = UPPER($1) OR token = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, cleanToken).Scan(&rec.ID, &rec.UserID, &rec.Token, &rec.TokenType, &metaBytes, &rec.ExpiresAt, &rec.CreatedAt)
	if err != nil {
		return nil, err
	}
	if len(metaBytes) > 0 {
		_ = json.Unmarshal(metaBytes, &rec.Metadata)
	}
	return rec, nil
}

func (r *UserRepository) GetVerificationToken(ctx context.Context, token string) (string, string, time.Time, error) {
	rec, err := r.GetVerificationTokenRecord(ctx, token)
	if err != nil {
		return "", "", time.Time{}, err
	}
	if targetEmail, ok := rec.Metadata["new_email"].(string); ok && targetEmail != "" {
		return fmt.Sprintf("%s:%s", rec.UserID, targetEmail), rec.TokenType, rec.ExpiresAt, nil
	}
	db := r.server.DB.Conn(ctx)
	var email string
	_ = db.QueryRow(ctx, `SELECT email FROM identity.users WHERE id::text = $1`, rec.UserID).Scan(&email)
	return email, rec.TokenType, rec.ExpiresAt, nil
}

func (r *UserRepository) DeleteVerificationToken(ctx context.Context, token string) error {
	db := r.server.DB.Conn(ctx)
	cleanToken := strings.TrimSpace(token)
	_, err := db.Exec(ctx, `DELETE FROM identity.verification_tokens WHERE UPPER(token) = UPPER($1) OR token = $1`, cleanToken)
	return err
}

func (r *UserRepository) UpdateEmailVerified(ctx context.Context, email string, verified bool) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `UPDATE identity.users SET email_verified = @verified WHERE email = @email`, pgx.NamedArgs{"email": email, "verified": verified})
	return err
}

func (r *UserRepository) CreatePasswordResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `
		INSERT INTO identity.password_reset_tokens (id, user_id, token_hash, expires_at)
		VALUES (@id, @userID, @tokenHash, @expiresAt)
	`, pgx.NamedArgs{
		"id":        uuid.New().String(),
		"userID":    userID,
		"tokenHash": tokenHash,
		"expiresAt": expiresAt,
	})
	return err
}

func (r *UserRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (string, time.Time, *time.Time, error) {
	db := r.server.DB.Conn(ctx)
	var userID string
	var expiresAt time.Time
	var usedAt *time.Time
	err := db.QueryRow(ctx, `SELECT user_id::text, expires_at, used_at FROM identity.password_reset_tokens WHERE token_hash = @tokenHash`, pgx.NamedArgs{"tokenHash": tokenHash}).
		Scan(&userID, &expiresAt, &usedAt)
	return userID, expiresAt, usedAt, err
}

func (r *UserRepository) MarkPasswordResetTokenUsed(ctx context.Context, tokenHash string) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `UPDATE identity.password_reset_tokens SET used_at = NOW() WHERE token_hash = @tokenHash`, pgx.NamedArgs{"tokenHash": tokenHash})
	return err
}

func (r *UserRepository) GetPasswordSetupToken(ctx context.Context, tokenHash string) (string, *string, time.Time, *time.Time, error) {
	db := r.server.DB.Conn(ctx)
	var (
		userID    string
		orgID     *string
		expiresAt time.Time
		usedAt    *time.Time
	)
	err := db.QueryRow(ctx, `
		SELECT user_id::text, organization_id::text, expires_at, used_at 
		FROM identity.password_setup_tokens 
		WHERE token_hash = @tokenHash
	`, pgx.NamedArgs{"tokenHash": tokenHash}).Scan(&userID, &orgID, &expiresAt, &usedAt)
	return userID, orgID, expiresAt, usedAt, err
}

func (r *UserRepository) MarkPasswordSetupTokenUsed(ctx context.Context, tokenHash string) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `UPDATE identity.password_setup_tokens SET used_at = NOW() WHERE token_hash = @tokenHash`, pgx.NamedArgs{"tokenHash": tokenHash})
	return err
}

func (r *UserRepository) RecordPasswordRequest(ctx context.Context, userID, email, status, ipAddress, userAgent string) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `
		INSERT INTO identity.password_requests (id, user_id, email, status, ip_address, user_agent, created_at)
		VALUES (@id, @userID, @email, @status, @ipAddress, @userAgent, NOW())
	`, pgx.NamedArgs{
		"id":        uuid.New().String(),
		"userID":    userID,
		"email":     email,
		"status":    status,
		"ipAddress": ipAddress,
		"userAgent": userAgent,
	})
	return err
}

func (r *UserRepository) RevokeAllUserSessions(ctx context.Context, userID string) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `DELETE FROM identity.sessions WHERE user_id::text = @userID`, pgx.NamedArgs{"userID": userID})
	return err
}


func (r *UserRepository) ListUsersByTenant(ctx context.Context, tenantID string, allBranches bool) ([]model.MembershipWithDetails, error) {
	dbExec := r.server.DB.Conn(ctx)
	var rows pgx.Rows
	var err error
	args := pgx.NamedArgs{"tenantID": tenantID}
	if allBranches {
		stmt := `
			SELECT m.id, m.user_id, u.name, u.email, COALESCE(m.organization_id::text, ''), COALESCE(t.name, ''), m.role_title, m.role_title, m.is_active, m.joined_at, m.created_at
			FROM organization.organization_memberships m
			JOIN identity.users u ON u.id = m.user_id
			LEFT JOIN workspace.workspaces t ON t.organization_id = m.organization_id
			WHERE m.organization_id = @tenantID
			ORDER BY m.created_at DESC
		`
		rows, err = dbExec.Query(ctx, stmt, args)
	} else {
		stmt := `
			SELECT m.id, m.user_id, u.name, u.email, COALESCE(m.organization_id::text, ''), COALESCE(t.name, ''), m.role_title, m.role_title, m.is_active, m.joined_at, m.created_at
			FROM organization.organization_memberships m
			JOIN identity.users u ON u.id = m.user_id
			LEFT JOIN workspace.workspaces t ON t.organization_id = m.organization_id
			WHERE m.organization_id = @tenantID
			ORDER BY m.created_at DESC
		`
		rows, err = dbExec.Query(ctx, stmt, args)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query users by tenant: %w", err)
	}
	defer rows.Close()

	var items []model.MembershipWithDetails
	for rows.Next() {
		var item model.MembershipWithDetails
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.UserName, &item.UserEmail,
			&item.TenantID, &item.TenantName, &item.RoleID, &item.RoleName,
			&item.IsActive, &item.JoinedAt, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan membership item: %w", err)
		}
		item.RoleSystem = "workspace"
		item.Name = item.UserName
		item.Email = item.UserEmail
		item.Role = item.RoleName
		items = append(items, item)
	}
	if items == nil {
		items = []model.MembershipWithDetails{}
	}
	return items, nil
}

func (r *UserRepository) ListUsersByOrganization(ctx context.Context, orgID string) ([]model.MembershipWithDetails, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT DISTINCT ON (m.id) m.id, m.user_id, u.name, u.email, COALESCE(m.organization_id::text, ''), COALESCE(o.name, ''), m.role_title, m.role_title, m.is_active, m.joined_at, m.created_at
		FROM organization.organization_memberships m
		JOIN identity.users u ON u.id = m.user_id
		LEFT JOIN organization.organizations o ON o.id = m.organization_id
		WHERE m.organization_id = @orgID
		ORDER BY m.id, m.created_at DESC
	`
	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"orgID": orgID})
	if err != nil {
		return nil, fmt.Errorf("failed to query users by organization: %w", err)
	}
	defer rows.Close()

	var items []model.MembershipWithDetails
	for rows.Next() {
		var item model.MembershipWithDetails
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.UserName, &item.UserEmail,
			&item.TenantID, &item.TenantName, &item.RoleID, &item.RoleName,
			&item.IsActive, &item.JoinedAt, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan organization membership item: %w", err)
		}
		item.RoleSystem = "workspace"
		item.Name = item.UserName
		item.Email = item.UserEmail
		item.Role = item.RoleName
		items = append(items, item)
	}
	if items == nil {
		items = []model.MembershipWithDetails{}
	}
	return items, nil
}

func (r *UserRepository) ListAllUsers(ctx context.Context) ([]model.MembershipWithDetails, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT 
			COALESCE(m.id::text, u.id::text), 
			u.id, 
			u.name, 
			u.email, 
			COALESCE(m.organization_id::text, ''), 
			COALESCE(o.name, 'Global Platform'), 
			COALESCE(m.role_title, u.platform_role, CASE WHEN u.is_platform_admin THEN 'super_admin' ELSE 'member' END), 
			COALESCE(m.role_title, u.platform_role, CASE WHEN u.is_platform_admin THEN 'super_admin' ELSE 'member' END), 
			COALESCE(m.is_active, true), 
			COALESCE(m.joined_at, u.created_at), 
			u.created_at
		FROM identity.users u
		LEFT JOIN organization.organization_memberships m ON m.user_id = u.id
		LEFT JOIN organization.organizations o ON o.id = m.organization_id
		ORDER BY u.created_at DESC
	`
	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query all users: %w", err)
	}
	defer rows.Close()

	var items []model.MembershipWithDetails
	for rows.Next() {
		var item model.MembershipWithDetails
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.UserName, &item.UserEmail,
			&item.TenantID, &item.TenantName, &item.RoleID, &item.RoleName,
			&item.IsActive, &item.JoinedAt, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan membership item: %w", err)
		}
		item.RoleSystem = "platform"
		item.Name = item.UserName
		item.Email = item.UserEmail
		item.Role = item.RoleName
		items = append(items, item)
	}
	if items == nil {
		items = []model.MembershipWithDetails{}
	}
	return items, nil
}

func (r *UserRepository) ListRoles(ctx context.Context, scopeFilter ...string) ([]model.RoleWithPermissions, error) {
	dbExec := r.server.DB.Conn(ctx)

	stmt := `
		SELECT r.id, r.code, COALESCE(r.context_scope, 'workspace'), COALESCE(r.description, ''), NULL::text, r.created_at,
		       COALESCE(ARRAY_AGG(p.code) FILTER (WHERE p.code IS NOT NULL), '{}') as permissions
		FROM "authorization".roles r
		LEFT JOIN "authorization".role_permissions rp ON rp.role_id = r.id
		LEFT JOIN "authorization".permissions p ON p.id = rp.permission_id
		GROUP BY r.id, r.code, r.context_scope, r.description, r.created_at
		ORDER BY r.code ASC
	`

	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()

	var roles []model.RoleWithPermissions
	for rows.Next() {
		var item model.RoleWithPermissions
		var tenantID *string
		if err := rows.Scan(
			&item.ID, &item.Name, &item.Scope, &item.Description, &tenantID,
			&item.CreatedAt, &item.Permissions,
		); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		if tenantID != nil {
			uID, _ := uuid.Parse(*tenantID)
			item.TenantID = &uID
		}
		item.IsSystem = tenantID == nil
		item.System = item.Scope
		if item.System == "" {
			if tenantID == nil {
				item.System = "platform"
			} else {
				item.System = "workspace"
			}
		}
		item.UpdatedAt = item.CreatedAt
		roles = append(roles, item)
	}
	if roles == nil {
		roles = []model.RoleWithPermissions{}
	}
	return roles, nil
}

func (r *UserRepository) ListAvailableTenants(ctx context.Context, userID string) ([]model.TenantSelectorItem, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT DISTINCT o.id, o.name, o.slug
		FROM organization.organization_memberships m
		JOIN organization.organizations o ON o.id = m.organization_id
		WHERE m.user_id::text = @userID AND m.is_active = TRUE
		ORDER BY o.name ASC
	`
	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"userID": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to query available tenants: %w", err)
	}
	defer rows.Close()

	var items []model.TenantSelectorItem
	for rows.Next() {
		var item model.TenantSelectorItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Slug); err != nil {
			return nil, fmt.Errorf("failed to scan tenant item: %w", err)
		}
		items = append(items, item)
	}
	if items == nil {
		items = []model.TenantSelectorItem{}
	}
	return items, nil
}

func (r *UserRepository) GetActiveMembership(ctx context.Context, userID, tenantID string) (membershipID string, roleName string, err error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT m.id, m.role_title
		FROM organization.organization_memberships m
		WHERE m.user_id::text = @userID AND m.organization_id::text = @tenantID AND m.is_active = TRUE
	`
	err = dbExec.QueryRow(ctx, stmt, pgx.NamedArgs{"userID": userID, "tenantID": tenantID}).Scan(&membershipID, &roleName)
	if err != nil {
		return "", "", err
	}
	return membershipID, roleName, nil
}

func (r *UserRepository) ListPermissions(ctx context.Context) ([]string, error) {
	dbExec := r.server.DB.Conn(ctx)
	rows, err := dbExec.Query(ctx, `SELECT code FROM "authorization".permissions ORDER BY code ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query permissions: %w", err)
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err == nil {
			permissions = append(permissions, code)
		}
	}
	if permissions == nil {
		permissions = []string{}
	}
	return permissions, nil
}

func (r *UserRepository) ListPermissionsByRole(ctx context.Context, roleName, tenantID string) ([]string, error) {
	if roleName == "owner" || roleName == "super_admin" || roleName == "org_admin" {
		return r.ListPermissions(ctx)
	}
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT DISTINCT p.code
		FROM "authorization".roles r
		JOIN "authorization".role_permissions rp ON rp.role_id = r.id
		JOIN "authorization".permissions p ON p.id = rp.permission_id
		WHERE r.code = @roleName OR r.name = @roleName
	`
	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"roleName": roleName, "tenantID": tenantID})
	if err != nil {
		r.server.Logger.Warn().Err(err).Str("role", roleName).Msg("role permissions missing or query failed - failing closed with empty permissions")
		return []string{}, nil
	}
	defer rows.Close()

	var permissions []string = []string{}
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err == nil && code != "" {
			permissions = append(permissions, code)
		}
	}
	return permissions, nil
}

func (r *UserRepository) CheckMembershipExists(ctx context.Context, userID, tenantID string) (bool, error) {
	dbExec := r.server.DB.Conn(ctx)
	var exists bool
	stmt := `SELECT EXISTS(
		SELECT 1 FROM organization.organization_memberships 
		WHERE user_id = @userID AND (organization_id = @tenantID OR id::text = @tenantID) AND is_active = TRUE
	)`
	err := dbExec.QueryRow(ctx, stmt, pgx.NamedArgs{"userID": userID, "tenantID": tenantID}).Scan(&exists)
	return exists, err
}

func (r *UserRepository) UpdateActiveTenantContext(ctx context.Context, sessionID, atc string) error {
	dbExec := r.server.DB.Conn(ctx)
	_, err := dbExec.Exec(ctx, `UPDATE identity.sessions SET updated_at = CURRENT_TIMESTAMP WHERE id = @sessionID`, pgx.NamedArgs{"sessionID": sessionID})
	return err
}

func (r *UserRepository) UpdateSessionActiveTenant(ctx context.Context, sessionID string, atc *model.ActiveTenantContext) error {
	dbExec := r.server.DB.Conn(ctx)
	_, err := dbExec.Exec(ctx, `UPDATE identity.sessions SET updated_at = CURRENT_TIMESTAMP WHERE id = @sessionID`, pgx.NamedArgs{"sessionID": sessionID})
	return err
}

func (r *UserRepository) GetRoleIDByName(ctx context.Context, roleName, tenantID string) (string, error) {
	dbExec := r.server.DB.Conn(ctx)
	var roleID string
	stmt := `SELECT id::text FROM "authorization".roles WHERE name = @roleName OR code = @roleName LIMIT 1`
	err := dbExec.QueryRow(ctx, stmt, pgx.NamedArgs{"roleName": roleName}).Scan(&roleID)
	return roleID, err
}

func (r *UserRepository) GetRoleDetailsByName(ctx context.Context, roleName, tenantID string) (string, string, bool, error) {
	dbExec := r.server.DB.Conn(ctx)
	var roleID, name string
	stmt := `
		SELECT id::text, name, is_system
		FROM "authorization".roles
		WHERE name = @roleName OR code = @roleName
		LIMIT 1
	`
	var isSystem bool
	err := dbExec.QueryRow(ctx, stmt, pgx.NamedArgs{"roleName": roleName}).Scan(&roleID, &name, &isSystem)
	if err != nil {
		return "", "", false, err
	}
	return roleID, name, isSystem, nil
}

func (r *UserRepository) AddOrReactivateMembership(ctx context.Context, userID, tenantID, roleID string) error {
	dbExec := r.server.DB.Conn(ctx)
	id := uuid.New().String()
	stmt := `
		INSERT INTO organization.organization_memberships (id, user_id, organization_id, role_title, is_active, joined_at)
		VALUES (@id, @userID, @tenantID, @roleID, TRUE, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, organization_id) DO UPDATE SET
			role_title = EXCLUDED.role_title,
			is_active = TRUE,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{"id": id, "userID": userID, "tenantID": tenantID, "roleID": roleID})
	return err
}

func (r *UserRepository) UpdateMembershipRole(ctx context.Context, membershipID, roleID, tenantID string) (bool, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE organization.organization_memberships
		SET role_title = @roleID, updated_at = CURRENT_TIMESTAMP
		WHERE id = @membershipID AND organization_id = @tenantID AND is_active = TRUE
	`
	tag, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{"membershipID": membershipID, "roleID": roleID, "tenantID": tenantID})
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (r *UserRepository) DeleteMembership(ctx context.Context, membershipID, tenantID string) (bool, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE organization.organization_memberships
		SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP
		WHERE id = @membershipID AND organization_id = @tenantID
	`
	tag, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{"membershipID": membershipID, "tenantID": tenantID})
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (r *UserRepository) DeactivateMembership(ctx context.Context, membershipID, tenantID string) (bool, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE organization.organization_memberships
		SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP
		WHERE id = @membershipID AND (organization_id = @tenantID OR id::text = @membershipID) AND is_active = TRUE
	`
	res, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{"membershipID": membershipID, "tenantID": tenantID})
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func (r *UserRepository) ActivateMembership(ctx context.Context, membershipID, tenantID string) (bool, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE organization.organization_memberships
		SET is_active = TRUE, updated_at = CURRENT_TIMESTAMP
		WHERE id = @membershipID AND (organization_id = @tenantID OR id::text = @membershipID) AND is_active = FALSE
	`
	res, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{"membershipID": membershipID, "tenantID": tenantID})
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func (r *UserRepository) GetProfile(ctx context.Context, userID string) (*modelUser.UserProfile, error) {
	dbExec := r.server.DB.Conn(ctx)
	u := &modelUser.UserProfile{UserID: userID}

	var email, userImage *string
	var firstName, lastName, middleName, profilePhone, bio, profileAvatar *string
	var gender, nationality, primaryPhone, timezone, languageCode, emergencyContact *string
	var dateOfBirth *time.Time

	err := dbExec.QueryRow(ctx, `
		SELECT 
			u.email, 
			u.avatar_url,
			p.first_name, 
			p.last_name, 
			p.middle_name, 
			p.phone_number, 
			p.bio, 
			p.avatar_url, 
			p.gender, 
			p.date_of_birth, 
			p.nationality, 
			p.timezone, 
			p.emergency_contact
		FROM identity.users u
		LEFT JOIN identity.user_profiles p ON p.user_id = u.id
		WHERE u.id::text = @userID
	`, pgx.NamedArgs{"userID": userID}).Scan(
		&email, &userImage,
		&firstName, &lastName, &middleName, &profilePhone, &bio, &profileAvatar,
		&gender, &dateOfBirth, &nationality, &timezone, &emergencyContact,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query user profile: %w", err)
	}

	// 1. Authoritative name fields directly from user_profiles (Single Source of Truth)
	u.FirstName = firstName
	u.LastName = lastName
	u.MiddleName = middleName

	// 2. Resolve avatar
	if profileAvatar != nil && *profileAvatar != "" {
		u.AvatarURL = profileAvatar
	} else if userImage != nil {
		u.AvatarURL = userImage
	}

	// 3. Resolve phone
	if profilePhone != nil && *profilePhone != "" {
		u.PhoneNumber = profilePhone
	}
	if primaryPhone != nil && *primaryPhone != "" {
		u.PrimaryPhone = primaryPhone
	} else if u.PhoneNumber != nil {
		u.PrimaryPhone = u.PhoneNumber
	}

	// 4. Resolve remaining profile fields
	u.Bio = bio
	u.Gender = gender
	u.DateOfBirth = dateOfBirth
	u.Nationality = nationality
	u.Timezone = timezone
	u.LanguageCode = languageCode
	u.EmergencyContact = emergencyContact

	return u, nil
}

func (r *UserRepository) UpsertProfile(ctx context.Context, payload *modelUser.UpdateUserProfilePayload) (*modelUser.UserProfile, error) {
	var resultProfile *modelUser.UserProfile

	err := r.server.DB.RunInTx(ctx, func(txCtx context.Context) error {
		dbExec := r.server.DB.Conn(txCtx)

		var parsedDOB *time.Time
		if payload.DateOfBirth != nil && *payload.DateOfBirth != "" {
			if t, err := time.Parse("2006-01-02", *payload.DateOfBirth); err == nil {
				parsedDOB = &t
			} else if t, err := time.Parse(time.RFC3339, *payload.DateOfBirth); err == nil {
				parsedDOB = &t
			}
		}

		// 1. Insert or update identity.user_profiles (Single Source of Truth)
		_, err := dbExec.Exec(txCtx, `
			INSERT INTO identity.user_profiles (
				user_id, first_name, last_name, middle_name, phone_number, bio, avatar_url, 
				gender, date_of_birth, nationality, timezone, emergency_contact, updated_at
			)
			VALUES (
				@userID, @firstName, @lastName, @middleName, @phone, @bio, @avatarUrl, 
				@gender, @dob, @nationality, @timezone, @emergencyContact, CURRENT_TIMESTAMP
			)
			ON CONFLICT (user_id) DO UPDATE SET
				first_name = COALESCE(@firstName, user_profiles.first_name),
				last_name = COALESCE(@lastName, user_profiles.last_name),
				middle_name = COALESCE(@middleName, user_profiles.middle_name),
				phone_number = COALESCE(@phone, user_profiles.phone_number),
				bio = COALESCE(@bio, user_profiles.bio),
				avatar_url = COALESCE(@avatarUrl, user_profiles.avatar_url),
				gender = COALESCE(@gender, user_profiles.gender),
				date_of_birth = COALESCE(@dob, user_profiles.date_of_birth),
				nationality = COALESCE(@nationality, user_profiles.nationality),
				timezone = COALESCE(@timezone, user_profiles.timezone),
				emergency_contact = COALESCE(@emergencyContact, user_profiles.emergency_contact),
				updated_at = CURRENT_TIMESTAMP
		`, pgx.NamedArgs{
			"userID":           payload.UserID,
			"firstName":        payload.FirstName,
			"lastName":         payload.LastName,
			"middleName":       payload.MiddleName,
			"phone":            payload.PhoneNumber,
			"bio":              payload.Bio,
			"avatarUrl":        payload.AvatarURL,
			"gender":           payload.Gender,
			"dob":              parsedDOB,
			"nationality":      payload.Nationality,
			"timezone":         payload.Timezone,
			"emergencyContact": payload.EmergencyContact,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert profile: %w", err)
		}

		// 2. Fetch updated profile directly from user_profiles to compute derived display name
		p, err := r.GetProfile(txCtx, payload.UserID)
		if err != nil {
			return fmt.Errorf("failed to reload profile: %w", err)
		}

		// 3. Synchronize identity.users.name and avatar_url transactionally
		var nameParts []string
		if p.FirstName != nil && strings.TrimSpace(*p.FirstName) != "" {
			nameParts = append(nameParts, strings.TrimSpace(*p.FirstName))
		}
		if p.MiddleName != nil && strings.TrimSpace(*p.MiddleName) != "" {
			nameParts = append(nameParts, strings.TrimSpace(*p.MiddleName))
		}
		if p.LastName != nil && strings.TrimSpace(*p.LastName) != "" {
			nameParts = append(nameParts, strings.TrimSpace(*p.LastName))
		}

		derivedName := strings.Join(nameParts, " ")
		if derivedName != "" {
			_, err = dbExec.Exec(txCtx, `
				UPDATE identity.users 
				SET name = @name, updated_at = CURRENT_TIMESTAMP 
				WHERE id::text = @userID
			`, pgx.NamedArgs{"name": derivedName, "userID": payload.UserID})
			if err != nil {
				return fmt.Errorf("failed to sync derived name: %w", err)
			}
		}

		if payload.AvatarURL != nil {
			_, err = dbExec.Exec(txCtx, `
				UPDATE identity.users 
				SET avatar_url = @avatarUrl, updated_at = CURRENT_TIMESTAMP 
				WHERE id::text = @userID
			`, pgx.NamedArgs{"avatarUrl": *payload.AvatarURL, "userID": payload.UserID})
			if err != nil {
				return fmt.Errorf("failed to sync avatar: %w", err)
			}
		}

		resultProfile = p
		return nil
	})

	if err != nil {
		return nil, err
	}
	return resultProfile, nil
}

func (r *UserRepository) UpdateUserEmail(ctx context.Context, userID, newEmail string) error {
	dbExec := r.server.DB.Conn(ctx)
	_, err := dbExec.Exec(ctx, `
		UPDATE identity.users
		SET email = @email, email_verified = TRUE, updated_at = CURRENT_TIMESTAMP
		WHERE id::text = @userID
	`, pgx.NamedArgs{"email": newEmail, "userID": userID})
	if err != nil {
		return fmt.Errorf("failed to update user email: %w", err)
	}

	// Also update credential account if exists
	_, _ = dbExec.Exec(ctx, `
		UPDATE identity.credentials
		SET account_id = @email, updated_at = CURRENT_TIMESTAMP
		WHERE user_id::text = @userID AND auth_provider = 'credential'
	`, pgx.NamedArgs{"email": newEmail, "userID": userID})

	return nil
}

func (r *UserRepository) GetProfessionalProfiles(ctx context.Context, userID string) ([]modelUser.ProfessionalProfile, error) {
	return []modelUser.ProfessionalProfile{}, nil
}

func (r *UserRepository) CreateProfessionalProfile(ctx context.Context, userID string, payload *modelUser.CreateProfessionalProfilePayload) (*modelUser.ProfessionalProfile, error) {
	return &modelUser.ProfessionalProfile{
		ID:        uuid.New().String(),
		UserID:    userID,
		Specialty: payload.Specialty,
	}, nil
}

func (r *UserRepository) RevokeSession(ctx context.Context, sessionID string) error {
	dbExec := r.server.DB.Conn(ctx)
	_, err := dbExec.Exec(ctx, `DELETE FROM identity.sessions WHERE id = @sessionID`, pgx.NamedArgs{"sessionID": sessionID})
	return err
}

func (r *UserRepository) IsBranchOnlyUser(ctx context.Context, userID string) (bool, error) {
	dbExec := r.server.DB.Conn(ctx)
	var isBranchOnly bool
	stmt := `
		SELECT EXISTS(
			SELECT 1 FROM organization.organization_memberships m
			WHERE m.user_id::text = @userID 
			  AND m.is_active = TRUE
			  AND m.role_title NOT IN ('owner', 'org_admin', 'super_admin')
		)
	`
	err := dbExec.QueryRow(ctx, stmt, pgx.NamedArgs{"userID": userID}).Scan(&isBranchOnly)
	return isBranchOnly, err
}

func (r *UserRepository) GetUserTenantSlugFallback(ctx context.Context, userID string) (string, error) {
	dbExec := r.server.DB.Conn(ctx)
	var slug string
	stmt := `
		SELECT w.slug
		FROM organization.organization_memberships m
		JOIN workspace.workspaces w ON w.organization_id = m.organization_id
		WHERE m.user_id::text = @userID AND m.is_active = TRUE
		ORDER BY m.created_at ASC
		LIMIT 1
	`
	err := dbExec.QueryRow(ctx, stmt, pgx.NamedArgs{"userID": userID}).Scan(&slug)
	if err != nil {
		return "", err
	}
	return slug, nil
}

func (r *UserRepository) CheckUserWorkspaceAccess(ctx context.Context, userID string, subdomain string) (bool, error) {
	dbExec := r.server.DB.Conn(ctx)
	var hasAccess bool
	stmt := `
		SELECT EXISTS(
			SELECT 1 FROM organization.organization_memberships m
			JOIN workspace.workspaces w ON w.organization_id = m.organization_id
			WHERE m.user_id::text = @userID AND (w.slug = @subdomain OR w.id::text = @subdomain) AND m.is_active = TRUE
		)
	`
	err := dbExec.QueryRow(ctx, stmt, pgx.NamedArgs{"userID": userID, "subdomain": subdomain}).Scan(&hasAccess)
	return hasAccess, err
}

func (r *UserRepository) GetSessionByID(ctx context.Context, sessionID string) (*model.Session, error) {
	dbExec := r.server.DB.Conn(ctx)
	sess := &model.Session{}
	stmt := `SELECT id, user_id::text, token, expires_at FROM identity.sessions WHERE id = @sessionID`
	err := dbExec.QueryRow(ctx, stmt, pgx.NamedArgs{"sessionID": sessionID}).Scan(
		&sess.ID, &sess.UserID, &sess.Token, &sess.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (r *UserRepository) GetPermissionOverrides(ctx context.Context, userID, tenantID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func (r *UserRepository) CreatePermissionOverride(ctx context.Context, userID, tenantID, permission, overrideType, createdBy string, expiresAt *time.Time) (string, error) {
	return uuid.New().String(), nil
}

func (r *UserRepository) DeletePermissionOverride(ctx context.Context, overrideID string) error {
	return nil
}

func (r *UserRepository) GetSignatures(ctx context.Context, userID, tenantID string) ([]modelUser.ProfessionalSignature, error) {
	return []modelUser.ProfessionalSignature{}, nil
}

func (r *UserRepository) CreateSignature(ctx context.Context, userID, tenantID string, payload *modelUser.CreateProfessionalSignaturePayload) (*modelUser.ProfessionalSignature, error) {
	return &modelUser.ProfessionalSignature{
		ID:       uuid.New().String(),
		UserID:   userID,
		TenantID: tenantID,
	}, nil
}

func (r *UserRepository) DeactivateSignature(ctx context.Context, sigID string) (bool, error) {
	return true, nil
}

func (r *UserRepository) GetTenantEmployment(ctx context.Context, userID, tenantID string) (*modelUser.TenantStaff, error) {
	return &modelUser.TenantStaff{
		ID:       uuid.New().String(),
		UserID:   userID,
		TenantID: tenantID,
	}, nil
}

func (r *UserRepository) UpdateTenantEmployment(ctx context.Context, userID, tenantID string, payload *modelUser.UpdateTenantStaffPayload) (*modelUser.TenantStaff, error) {
	return &modelUser.TenantStaff{
		ID:       uuid.New().String(),
		UserID:   userID,
		TenantID: tenantID,
	}, nil
}

func (r *UserRepository) GetCompetencies(ctx context.Context, staffID string) ([]modelUser.ProfessionalCompetency, error) {
	return []modelUser.ProfessionalCompetency{}, nil
}

func (r *UserRepository) CreateCompetency(ctx context.Context, staffID string, payload *modelUser.CreateProfessionalCompetencyPayload) (*modelUser.ProfessionalCompetency, error) {
	return &modelUser.ProfessionalCompetency{
		ID:      uuid.New().String(),
		StaffID: staffID,
	}, nil
}

func (r *UserRepository) GetPrivileges(ctx context.Context, staffID string) ([]modelUser.ClinicalPrivilege, error) {
	return []modelUser.ClinicalPrivilege{}, nil
}

func (r *UserRepository) CreatePrivilege(ctx context.Context, staffID string, payload *modelUser.CreateClinicalPrivilegePayload, approvedBy string) (*modelUser.ClinicalPrivilege, error) {
	return &modelUser.ClinicalPrivilege{
		ID:      uuid.New().String(),
		StaffID: staffID,
	}, nil
}

func (r *UserRepository) CreateVerification(ctx context.Context, profileID, actorID string, payload *modelUser.CreateProfessionalVerificationPayload) (*modelUser.ProfessionalVerification, error) {
	return &modelUser.ProfessionalVerification{
		ID:                    uuid.New().String(),
		ProfessionalProfileID: profileID,
	}, nil
}
