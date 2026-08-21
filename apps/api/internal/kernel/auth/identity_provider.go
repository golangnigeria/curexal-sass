package auth

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSessionNotFound  = errors.New("identity session not found")
	ErrIdentityNotFound = errors.New("identity not found")
)

// IdentityTraits defines standard identity properties stored with the authentication provider.
// Strict Architectural Requirement: Zero business domain data shall ever be stored in traits.
type IdentityTraits struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// IdentitySession represents an authenticated session issued by the identity provider.
type IdentitySession struct {
	ID              string         `json:"id"`
	IdentityID      string         `json:"identity_id"`
	Active          bool           `json:"active"`
	Traits          IdentityTraits `json:"traits"`
	ExpiresAt       time.Time      `json:"expires_at"`
	AuthenticatedAt time.Time      `json:"authenticated_at"`
}

// CreateIdentityParams contains attributes required to provision a new authentication identity.
type CreateIdentityParams struct {
	Traits         IdentityTraits `json:"traits"`
	Password       string         `json:"password,omitempty"`
	HashedPassword string         `json:"hashed_password,omitempty"` // For migration of existing Bcrypt hashes
}

// UpdateIdentityParams contains attributes to update an existing identity.
type UpdateIdentityParams struct {
	Traits IdentityTraits `json:"traits"`
}

// IdentityProvider is the interface implemented by external authentication providers (e.g. ORY Kratos, MockProvider).
// This abstraction guarantees 100% vendor independence and prevents ORY SDK calls from scattering throughout business code.
type IdentityProvider interface {
	// Name returns the provider's unique identifier (e.g., "kratos", "mock").
	Name() string

	// Authenticate validates a session token or cookie value and returns the active IdentitySession.
	Authenticate(ctx context.Context, sessionToken string) (*IdentitySession, error)

	// GetSession retrieves an active session by session ID.
	GetSession(ctx context.Context, sessionID string) (*IdentitySession, error)

	// GetIdentity fetches the traits of an identity by identity ID.
	GetIdentity(ctx context.Context, identityID string) (*IdentityTraits, error)

	// CreateIdentity provisions a new identity in the authentication provider.
	CreateIdentity(ctx context.Context, params *CreateIdentityParams) (string, error)

	// UpdateIdentity modifies traits for an existing identity.
	UpdateIdentity(ctx context.Context, identityID string, params *UpdateIdentityParams) error

	// DeleteIdentity removes an identity from the authentication provider.
	DeleteIdentity(ctx context.Context, identityID string) error
}
