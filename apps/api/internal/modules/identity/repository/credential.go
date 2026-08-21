package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/identity/model/auth"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
)

type PasswordHistoryRecord struct {
	ID           string    `db:"id"`
	CredentialID string    `db:"credential_id"`
	UserID       string    `db:"user_id"`
	PasswordHash string    `db:"password_hash"`
	ChangedBy    string    `db:"changed_by"`
	ChangeReason string    `db:"change_reason"`
	IPAddress    string    `db:"ip_address"`
	UserAgent    string    `db:"user_agent"`
	CreatedAt    time.Time `db:"created_at"`
}

type CredentialRepository struct {
	server *server.Server
}

func NewCredentialRepository(s *server.Server) *CredentialRepository {
	return &CredentialRepository{server: s}
}

// GetByEmail retrieves credential aggregate by user email address.
func (r *CredentialRepository) GetByEmail(ctx context.Context, email string) (*auth.Credential, error) {
	db := r.server.DB.Conn(ctx)
	cred := &auth.Credential{}
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	var statusStr string
	query := `
		SELECT u.id, u.email, COALESCE(c.password_hash, ''), COALESCE(u.credential_status, 'ACTIVE'), COALESCE(u.failed_login_attempts, 0), u.locked_until, u.password_changed_at, u.password_expires_at, u.last_successful_login_at, u.last_failed_login_at
		FROM identity.users u
		LEFT JOIN identity.credentials c ON c.user_id = u.id AND c.auth_provider = 'credential'
		WHERE LOWER(u.email) = $1 AND u.deleted_at IS NULL
	`
	err := db.QueryRow(ctx, query, cleanEmail).Scan(
		&cred.UserID, &cred.Email, &cred.PasswordHash, &statusStr, &cred.FailedLoginAttempts, &cred.LockedUntil, &cred.PasswordChangedAt, &cred.PasswordExpiresAt, &cred.LastSuccessfulLoginAt, &cred.LastFailedLoginAt,
	)
	if err != nil {
		return nil, err
	}
	cred.Status = auth.CredentialStatus(statusStr)
	return cred, nil
}

// GetByUserID retrieves credential aggregate by unique user ID.
func (r *CredentialRepository) GetByUserID(ctx context.Context, userID string) (*auth.Credential, error) {
	db := r.server.DB.Conn(ctx)
	cred := &auth.Credential{}
	var statusStr string
	query := `
		SELECT u.id, u.email, COALESCE(c.password_hash, ''), COALESCE(u.credential_status, 'ACTIVE'), COALESCE(u.failed_login_attempts, 0), u.locked_until, u.password_changed_at, u.password_expires_at, u.last_successful_login_at, u.last_failed_login_at
		FROM identity.users u
		LEFT JOIN identity.credentials c ON c.user_id = u.id AND c.auth_provider = 'credential'
		WHERE u.id::text = $1 AND u.deleted_at IS NULL
	`
	err := db.QueryRow(ctx, query, userID).Scan(
		&cred.UserID, &cred.Email, &cred.PasswordHash, &statusStr, &cred.FailedLoginAttempts, &cred.LockedUntil, &cred.PasswordChangedAt, &cred.PasswordExpiresAt, &cred.LastSuccessfulLoginAt, &cred.LastFailedLoginAt,
	)
	if err != nil {
		return nil, err
	}
	cred.Status = auth.CredentialStatus(statusStr)
	return cred, nil
}

// GetPasswordHistory retrieves up to limit recent password hashes for a user from identity.password_histories.
func (r *CredentialRepository) GetPasswordHistory(ctx context.Context, userID string, limit int) ([]string, error) {
	db := r.server.DB.Conn(ctx)
	if limit <= 0 {
		limit = 5
	}
	query := `
		SELECT password_hash
		FROM identity.password_histories
		WHERE user_id::text = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hashes []string
	for rows.Next() {
		var h string
		if err := rows.Scan(&h); err == nil && h != "" {
			hashes = append(hashes, h)
		}
	}
	return hashes, nil
}

// Save persists the complete state of a Credential aggregate root to the database.
func (r *CredentialRepository) Save(ctx context.Context, cred *auth.Credential) error {
	db := r.server.DB.Conn(ctx)
	query := `
		UPDATE identity.users
		SET credential_status = $1,
		    failed_login_attempts = $2,
		    locked_until = $3,
		    password_changed_at = $4,
		    password_expires_at = $5,
		    last_successful_login_at = $6,
		    last_failed_login_at = $7,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id::text = $8
	`
	statusStr := string(cred.Status)
	if statusStr == "" {
		statusStr = string(auth.CredentialStatusActive)
	}
	_, err := db.Exec(ctx, query,
		statusStr, cred.FailedLoginAttempts, cred.LockedUntil, cred.PasswordChangedAt, cred.PasswordExpiresAt, cred.LastSuccessfulLoginAt, cred.LastFailedLoginAt, cred.UserID,
	)
	if err != nil {
		return err
	}

	if cred.PasswordHash != "" {
		cleanEmail := strings.ToLower(strings.TrimSpace(cred.Email))
		res, err := db.Exec(ctx, `
			UPDATE identity.credentials
			SET password_hash = $1, account_id = COALESCE(NULLIF($2, ''), account_id), updated_at = CURRENT_TIMESTAMP
			WHERE user_id::text = $3 AND auth_provider = 'credential'
		`, cred.PasswordHash, cleanEmail, cred.UserID)
		if err != nil {
			return err
		}
		if res.RowsAffected() == 0 {
			_, err = db.Exec(ctx, `
				INSERT INTO identity.credentials (id, account_id, auth_provider, user_id, password_hash, created_at, updated_at)
				VALUES ($1, $2, 'credential', $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			`, uuid.New().String(), cleanEmail, cred.UserID, cred.PasswordHash)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdatePasswordHashWithAudit performs an atomic credential update and records a password history entry inside a single transaction.
func (r *CredentialRepository) UpdatePasswordHashWithAudit(ctx context.Context, userID string, newHash string, record *PasswordHistoryRecord) error {
	return r.server.DB.RunInTx(ctx, func(txCtx context.Context) error {
		db := r.server.DB.Conn(txCtx)
		now := time.Now()

		// 1. Fetch current active password hash to archive into password_histories
		var currentHash string
		_ = db.QueryRow(txCtx, `
			SELECT password_hash 
			FROM identity.credentials 
			WHERE user_id::text = $1 AND auth_provider = 'credential'
		`, userID).Scan(&currentHash)

		if currentHash != "" {
			changeReason := "PASSWORD_CHANGE"
			var ipAddr, userAgent *string
			if record != nil {
				if record.ChangeReason != "" {
					changeReason = record.ChangeReason
				}
				if record.IPAddress != "" {
					ipAddr = &record.IPAddress
				}
				if record.UserAgent != "" {
					userAgent = &record.UserAgent
				}
			}
			_, err := db.Exec(txCtx, `
				INSERT INTO identity.password_histories (user_id, password_hash, change_reason, ip_address, user_agent, created_at)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, userID, currentHash, changeReason, ipAddr, userAgent, now)
			if err != nil {
				return fmt.Errorf("failed to archive password history: %w", err)
			}
		}

		// 2. Update existing active credential row in-place (or insert initial credential if absent)
		var userEmail string
		_ = db.QueryRow(txCtx, `SELECT email FROM identity.users WHERE id::text = $1`, userID).Scan(&userEmail)
		cleanEmail := strings.ToLower(strings.TrimSpace(userEmail))

		res, err := db.Exec(txCtx, `
			UPDATE identity.credentials
			SET password_hash = $1,
			    account_id = COALESCE(NULLIF($2, ''), account_id),
			    updated_at = $3
			WHERE user_id::text = $4 AND auth_provider = 'credential'
		`, newHash, cleanEmail, now, userID)
		if err != nil {
			return fmt.Errorf("failed to update active credential: %w", err)
		}

		if res.RowsAffected() == 0 {
			credID := uuid.New().String()
			_, err = db.Exec(txCtx, `
				INSERT INTO identity.credentials (id, account_id, auth_provider, user_id, password_hash, created_at, updated_at)
				VALUES ($1, $2, 'credential', $3, $4, $5, $5)
			`, credID, cleanEmail, userID, newHash, now)
			if err != nil {
				return fmt.Errorf("failed to insert initial credential: %w", err)
			}
		}

		// 3. Update identity.users credential lifecycle fields
		_, err = db.Exec(txCtx, `
			UPDATE identity.users
			SET credential_status = 'ACTIVE',
			    failed_login_attempts = 0,
			    locked_until = NULL,
			    password_changed_at = $1,
			    updated_at = $1
			WHERE id::text = $2
		`, now, userID)
		if err != nil {
			return fmt.Errorf("failed to update user password lifecycle: %w", err)
		}

		return nil
	})
}

// UpdatePasswordHash updates password hash with default audit parameters.
func (r *CredentialRepository) UpdatePasswordHash(ctx context.Context, userID string, newHash string) error {
	return r.UpdatePasswordHashWithAudit(ctx, userID, newHash, &PasswordHistoryRecord{
		ChangeReason: "PASSWORD_CHANGE",
	})
}

// IncrementFailedLogin increments failed login counter and triggers a lockout if attempts reach 5.
func (r *CredentialRepository) IncrementFailedLogin(ctx context.Context, userID string) (int, *time.Time, error) {
	db := r.server.DB.Conn(ctx)
	var attempts int
	var lockedUntil *time.Time

	lockoutDuration := 1 * time.Minute
	if r.server != nil && r.server.Config != nil && r.server.Config.Auth.LockoutDuration > 0 {
		lockoutDuration = r.server.Config.Auth.LockoutDuration
	}
	lockoutSeconds := int(lockoutDuration.Seconds())
	if lockoutSeconds <= 0 {
		lockoutSeconds = 60
	}

	err := db.QueryRow(ctx, `
		UPDATE identity.users
		SET failed_login_attempts = COALESCE(failed_login_attempts, 0) + 1,
		    last_failed_login_at = CURRENT_TIMESTAMP,
		    credential_status = CASE WHEN COALESCE(failed_login_attempts, 0) + 1 >= 5 THEN 'LOCKED' ELSE credential_status END,
		    locked_until = CASE 
		        WHEN COALESCE(failed_login_attempts, 0) + 1 >= 5 THEN CURRENT_TIMESTAMP + ($2 * INTERVAL '1 second')
		        ELSE locked_until 
		    END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING failed_login_attempts, locked_until
	`, userID, lockoutSeconds).Scan(&attempts, &lockedUntil)
	if err != nil {
		return 0, nil, err
	}
	return attempts, lockedUntil, nil
}

// ResetFailedLogin clears failed login count and lifts account lockout.
func (r *CredentialRepository) ResetFailedLogin(ctx context.Context, userID string) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `
		UPDATE identity.users
		SET failed_login_attempts = 0,
		    locked_until = NULL,
		    credential_status = 'ACTIVE',
		    last_successful_login_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, userID)
	return err
}
