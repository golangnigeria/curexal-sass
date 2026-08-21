package domain

import (
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/golangnigeria/curexal/internal/shared/crypto"
)

var (
	ErrPasswordTooShort      = errors.New("password does not meet the minimum length requirement")
	ErrPasswordNoUppercase   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLowercase   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoDigit       = errors.New("password must contain at least one numeric digit")
	ErrPasswordNoSpecialChar = errors.New("password must contain at least one special character")
	ErrPasswordSameAsCurrent = errors.New("new password must be different from current password")
	ErrPasswordReuse         = errors.New("cannot reuse a recent password; please select a password you have not used recently")
)

// SecurityPolicy defines configurable security policy parameters for an organization or platform deployment.
type SecurityPolicy struct {
	MinLengthPlatform   int           `json:"minLengthPlatform"`
	MinLengthGeneral    int           `json:"minLengthGeneral"`
	PasswordHistoryDepth int          `json:"passwordHistoryDepth"`
	LockoutAttempts     int           `json:"lockoutAttempts"`
	LockoutDuration     time.Duration `json:"lockoutDuration"`
	PasswordExpiryDays  int           `json:"passwordExpiryDays"`
}

// DefaultSecurityPolicy returns standard enterprise security policy configuration.
func DefaultSecurityPolicy() *SecurityPolicy {
	return &SecurityPolicy{
		MinLengthPlatform:    12,
		MinLengthGeneral:     8,
		PasswordHistoryDepth: 5,
		LockoutAttempts:      5,
		LockoutDuration:      15 * time.Minute,
		PasswordExpiryDays:   90,
	}
}

// PasswordPolicyEvaluator handles domain logic for password validation.
type PasswordPolicyEvaluator struct {
	policy *SecurityPolicy
}

func NewPasswordPolicyEvaluator(p *SecurityPolicy) *PasswordPolicyEvaluator {
	if p == nil {
		p = DefaultSecurityPolicy()
	}
	return &PasswordPolicyEvaluator{policy: p}
}

// Validate evaluates password against complexity rules and minimum length.
func (p *PasswordPolicyEvaluator) Validate(password string, isPlatformUser bool) error {
	minLen := p.policy.MinLengthGeneral
	if isPlatformUser {
		minLen = p.policy.MinLengthPlatform
	}

	if len(password) < minLen {
		if isPlatformUser {
			return fmt.Errorf("%w: platform staff passwords require at least %d characters", ErrPasswordTooShort, minLen)
		}
		return fmt.Errorf("%w: passwords require at least %d characters", ErrPasswordTooShort, minLen)
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUppercase
	}
	if !hasLower {
		return ErrPasswordNoLowercase
	}
	if !hasDigit {
		return ErrPasswordNoDigit
	}
	if !hasSpecial {
		return ErrPasswordNoSpecialChar
	}

	return nil
}

// CheckPasswordHistory verifies that new password does not match any hash in the user's password history up to history depth.
func (p *PasswordPolicyEvaluator) CheckPasswordHistory(newPassword string, recentHashes []string) error {
	depth := p.policy.PasswordHistoryDepth
	if depth <= 0 {
		return nil
	}

	limit := len(recentHashes)
	if limit > depth {
		limit = depth
	}

	for i := 0; i < limit; i++ {
		hash := recentHashes[i]
		if hash == "" {
			continue
		}
		match, err := crypto.VerifyPassword(newPassword, hash)
		if err == nil && match {
			return ErrPasswordReuse
		}
	}
	return nil
}

// ValidatePassword is a package-level helper for quick evaluation.
func ValidatePassword(password string, isPlatformUser bool) error {
	evaluator := NewPasswordPolicyEvaluator(nil)
	return evaluator.Validate(password, isPlatformUser)
}

// CheckPasswordHistory is a package-level helper for history verification.
func CheckPasswordHistory(newPassword string, recentHashes []string) error {
	evaluator := NewPasswordPolicyEvaluator(nil)
	return evaluator.CheckPasswordHistory(newPassword, recentHashes)
}
