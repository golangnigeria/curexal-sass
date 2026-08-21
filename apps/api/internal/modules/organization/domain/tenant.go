package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Slug           string    `json:"slug" db:"slug"`
	LogoURL        *string   `json:"logoUrl,omitempty" db:"logo_url"`
	Settings       string    `json:"settings" db:"settings"`
	OrganizationID string    `json:"organizationId" db:"organization_id"`
	Currency       string    `json:"currency" db:"currency"`
	EnabledModules []string  `json:"enabledModules,omitempty" db:"-"`
}

type TenantDomain struct {
	ID                 uuid.UUID  `json:"id"                 db:"id"`
	TenantID           uuid.UUID  `json:"tenantId"           db:"tenant_id"`
	Domain             string     `json:"domain"             db:"domain"`
	VerificationType   string     `json:"verificationType"   db:"verification_type"`
	VerificationTarget string     `json:"verificationTarget" db:"verification_target"`
	IsVerified         bool       `json:"isVerified"         db:"is_verified"`
	SSLStatus          string     `json:"sslStatus"          db:"ssl_status"`
	VerifiedAt         *time.Time `json:"verifiedAt"         db:"verified_at"`
}
