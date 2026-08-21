package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockIdentityProvider_Lifecycle(t *testing.T) {
	ctx := context.Background()
	provider := NewMockIdentityProvider()

	assert.Equal(t, "mock", provider.Name())

	// 1. Create Identity
	params := &CreateIdentityParams{
		Traits: IdentityTraits{
			Email: "test.user@curexal.internal",
			Name:  "Test User",
		},
		Password: "SecretPassword123!",
	}
	identityID, err := provider.CreateIdentity(ctx, params)
	require.NoError(t, err)
	assert.NotEmpty(t, identityID)

	// 2. Fetch Identity Traits
	traits, err := provider.GetIdentity(ctx, identityID)
	require.NoError(t, err)
	assert.Equal(t, "test.user@curexal.internal", traits.Email)
	assert.Equal(t, "Test User", traits.Name)

	// 3. Add & Authenticate Session
	sessionToken := "sess_token_123"
	sess := &IdentitySession{
		ID:         "sess_123",
		IdentityID: identityID,
		Active:     true,
		Traits:     *traits,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}
	provider.AddSession(sessionToken, sess)

	authSess, err := provider.Authenticate(ctx, sessionToken)
	require.NoError(t, err)
	assert.True(t, authSess.Active)
	assert.Equal(t, identityID, authSess.IdentityID)

	// 4. Update Identity
	updateParams := &UpdateIdentityParams{
		Traits: IdentityTraits{
			Email: "updated.user@curexal.internal",
			Name:  "Updated Name",
		},
	}
	err = provider.UpdateIdentity(ctx, identityID, updateParams)
	require.NoError(t, err)

	updatedTraits, err := provider.GetIdentity(ctx, identityID)
	require.NoError(t, err)
	assert.Equal(t, "updated.user@curexal.internal", updatedTraits.Email)

	// 5. Delete Identity
	err = provider.DeleteIdentity(ctx, identityID)
	require.NoError(t, err)

	_, errDeleted := provider.GetIdentity(ctx, identityID)
	assert.ErrorIs(t, errDeleted, ErrIdentityNotFound)
}
