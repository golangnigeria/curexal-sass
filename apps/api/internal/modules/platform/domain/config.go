package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrOptimisticLockingConflict = errors.New("optimistic concurrency check failed: record has been modified by another request")
	ErrUnauthorizedPlatformAdmin = errors.New("unauthorized: operation requires platform administrator privileges")
	ErrInvalidPlatformConfig     = errors.New("invalid platform general settings configuration")
	ErrInvalidSecurityPolicy     = errors.New("invalid platform security policy configuration")
)

// PlatformGeneralSettings represents the global platform configuration entity.
type PlatformGeneralSettings struct {
	ID                  uuid.UUID `json:"id"`
	PlatformName        string    `json:"platformName"`
	PlatformLogoURL     *string   `json:"platformLogoUrl,omitempty"`
	PlatformFaviconURL  *string   `json:"platformFaviconUrl,omitempty"`
	PlatformDescription *string   `json:"platformDescription,omitempty"`
	SupportEmail        string    `json:"supportEmail"`
	SupportPhone        string    `json:"supportPhone"`
	SupportWebsite      *string   `json:"supportWebsite,omitempty"`
	DefaultCountry      string    `json:"defaultCountry"`
	DefaultCurrency     string    `json:"defaultCurrency"`
	SupportedCountries  []string  `json:"supportedCountries"`
	SupportedCurrencies []string  `json:"supportedCurrencies"`
	DefaultTimezone     string    `json:"defaultTimezone"`
	DefaultLocale       string    `json:"defaultLocale"`
	DateFormat          string    `json:"dateFormat"`
	TimeFormat          string    `json:"timeFormat"`
	NumberFormat        string    `json:"numberFormat"`
	MeasurementUnits    string    `json:"measurementUnits"`
	MaintenanceMode     bool      `json:"maintenanceMode"`
	AnnouncementBanner  *string   `json:"announcementBanner,omitempty"`
	Status              string    `json:"status"`
	Version             int       `json:"version"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	UpdatedBy           *uuid.UUID`json:"updatedBy,omitempty"`
}

// IdentitySecurityPolicy represents global security thresholds and identity governance parameters.
type IdentitySecurityPolicy struct {
	ID                            uuid.UUID `json:"id"`
	MinPasswordLength             int       `json:"minPasswordLength"`
	PasswordRequireUppercase      bool      `json:"passwordRequireUppercase"`
	PasswordRequireNumber         bool      `json:"passwordRequireNumber"`
	PasswordRequireSymbol         bool      `json:"passwordRequireSymbol"`
	PasswordExpirationDays        int       `json:"passwordExpirationDays"`
	EmailVerificationRequired     bool      `json:"emailVerificationRequired"`
	MFAPolicy                     string    `json:"mfaPolicy"`
	MaxFailedLoginAttempts        int       `json:"maxFailedLoginAttempts"`
	AccountLockoutDurationMinutes int       `json:"accountLockoutDurationMinutes"`
	SessionMaxDurationHours       int       `json:"sessionMaxDurationHours"`
	RefreshTokenDurationDays      int       `json:"refreshTokenDurationDays"`
	MaxActiveSessions             int       `json:"maxActiveSessions"`
	Version                       int       `json:"version"`
	CreatedAt                     time.Time `json:"createdAt"`
	UpdatedAt                     time.Time `json:"updatedAt"`
	UpdatedBy                     *uuid.UUID`json:"updatedBy,omitempty"`
}

func (s *PlatformGeneralSettings) Validate() error {
	if s.PlatformName == "" || s.SupportEmail == "" || s.DefaultCountry == "" || s.DefaultCurrency == "" {
		return ErrInvalidPlatformConfig
	}
	return nil
}

func (p *IdentitySecurityPolicy) Validate() error {
	if p.MinPasswordLength < 8 {
		return ErrInvalidSecurityPolicy
	}
	if p.MaxFailedLoginAttempts < 1 || p.AccountLockoutDurationMinutes < 1 {
		return ErrInvalidSecurityPolicy
	}
	return nil
}
