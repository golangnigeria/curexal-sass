package auth

import (
	"context"
	"errors"
	"sync"
	"time"
)

// MockIdentityProvider is an in-memory implementation of IdentityProvider for fast, hermetic unit tests.
type MockIdentityProvider struct {
	mu         sync.RWMutex
	Sessions   map[string]*IdentitySession
	Identities map[string]*IdentityTraits
}

func NewMockIdentityProvider() *MockIdentityProvider {
	return &MockIdentityProvider{
		Sessions:   make(map[string]*IdentitySession),
		Identities: make(map[string]*IdentityTraits),
	}
}

func (m *MockIdentityProvider) Name() string {
	return "mock"
}

func (m *MockIdentityProvider) Authenticate(ctx context.Context, sessionToken string) (*IdentitySession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sess, ok := m.Sessions[sessionToken]
	if !ok || !sess.Active {
		return nil, ErrSessionNotFound
	}
	return sess, nil
}

func (m *MockIdentityProvider) GetSession(ctx context.Context, sessionID string) (*IdentitySession, error) {
	return m.Authenticate(ctx, sessionID)
}

func (m *MockIdentityProvider) GetIdentity(ctx context.Context, identityID string) (*IdentityTraits, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	traits, ok := m.Identities[identityID]
	if !ok {
		return nil, ErrIdentityNotFound
	}
	return traits, nil
}

func (m *MockIdentityProvider) CreateIdentity(ctx context.Context, params *CreateIdentityParams) (string, error) {
	if params == nil {
		return "", errors.New("params cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	identityID := "id_mock_" + params.Traits.Email
	m.Identities[identityID] = &params.Traits
	return identityID, nil
}

func (m *MockIdentityProvider) UpdateIdentity(ctx context.Context, identityID string, params *UpdateIdentityParams) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Identities[identityID]; !ok {
		return ErrIdentityNotFound
	}
	m.Identities[identityID] = &params.Traits
	return nil
}

func (m *MockIdentityProvider) DeleteIdentity(ctx context.Context, identityID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.Identities, identityID)
	return nil
}

func (m *MockIdentityProvider) AddSession(token string, sess *IdentitySession) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if sess.AuthenticatedAt.IsZero() {
		sess.AuthenticatedAt = time.Now()
	}
	m.Sessions[token] = sess
}
