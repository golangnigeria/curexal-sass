package auth

import (
	"time"
)

type CredentialStatus string

const (
	CredentialStatusActive           CredentialStatus = "ACTIVE"
	CredentialStatusLocked           CredentialStatus = "LOCKED"
	CredentialStatusDisabled         CredentialStatus = "DISABLED"
	CredentialStatusPasswordExpired  CredentialStatus = "PASSWORD_EXPIRED"
	CredentialStatusPasswordResetReq CredentialStatus = "PASSWORD_RESET_REQUIRED"
	CredentialStatusInvited          CredentialStatus = "INVITED"
	CredentialStatusSuspended        CredentialStatus = "SUSPENDED"
)

// Credential represents the enterprise credential aggregate root owning authentication status and state machine logic.
type Credential struct {
	UserID                string           `json:"userId"                db:"user_id"`
	Email                 string           `json:"email"                 db:"email"`
	PasswordHash          string           `json:"-"                     db:"password_hash"`
	Status                CredentialStatus `json:"status"                db:"credential_status"`
	FailedLoginAttempts   int              `json:"failedLoginAttempts"   db:"failed_login_attempts"`
	LockedUntil           *time.Time       `json:"lockedUntil,omitempty" db:"locked_until"`
	PasswordChangedAt     *time.Time       `json:"passwordChangedAt,omitempty" db:"password_changed_at"`
	PasswordExpiresAt     *time.Time       `json:"passwordExpiresAt,omitempty" db:"password_expires_at"`
	LastSuccessfulLoginAt *time.Time       `json:"lastSuccessfulLoginAt,omitempty" db:"last_successful_login_at"`
	LastFailedLoginAt     *time.Time       `json:"lastFailedLoginAt,omitempty" db:"last_failed_login_at"`
}

// IsLocked reports whether the credential account is locked out.
func (c *Credential) IsLocked() bool {
	if c == nil {
		return false
	}
	if c.LockedUntil != nil {
		return c.LockedUntil.After(time.Now())
	}
	return c.Status == CredentialStatusLocked
}

// IsExpired reports whether the password has expired.
func (c *Credential) IsExpired() bool {
	if c == nil {
		return false
	}
	if c.Status == CredentialStatusPasswordExpired {
		return true
	}
	return c.PasswordExpiresAt != nil && c.PasswordExpiresAt.Before(time.Now())
}

// RecordSuccessfulLogin updates aggregate state upon valid authentication.
func (c *Credential) RecordSuccessfulLogin() {
	now := time.Now()
	c.FailedLoginAttempts = 0
	c.LockedUntil = nil
	c.LastSuccessfulLoginAt = &now
	if c.Status == CredentialStatusLocked {
		c.Status = CredentialStatusActive
	}
}

// RecordFailedLogin increments failed attempts and transitions status to LOCKED if threshold reached.
func (c *Credential) RecordFailedLogin(maxAttempts int, lockDuration time.Duration) {
	now := time.Now()
	c.FailedLoginAttempts++
	c.LastFailedLoginAt = &now
	if c.FailedLoginAttempts >= maxAttempts {
		lockedTime := now.Add(lockDuration)
		c.LockedUntil = &lockedTime
		c.Status = CredentialStatusLocked
	}
}

// ChangePassword updates the credential hash, resets failed attempts, and sets status to ACTIVE.
func (c *Credential) ChangePassword(newHash string) {
	now := time.Now()
	c.PasswordHash = newHash
	c.PasswordChangedAt = &now
	c.FailedLoginAttempts = 0
	c.LockedUntil = nil
	c.Status = CredentialStatusActive
}
