package auth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSupabaseIdentityProvider_AuthenticateJWT(t *testing.T) {
	jwtSecret := "super-secret-supabase-jwt-token-key-for-testing-12345"
	provider := NewSupabaseIdentityProvider(SupabaseConfig{
		JWTSecret: jwtSecret,
	})

	assert.Equal(t, "supabase", provider.Name())

	t.Run("Valid JWT Token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, SupabaseClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "usr_supa_12345",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			Email: "doctor@curexal.com",
			UserMetadata: map[string]interface{}{
				"name":       "Dr. Sarah Connor",
				"avatar_url": "https://cdn.curexal.space/avatars/sarah.png",
			},
		})

		signedToken, err := token.SignedString([]byte(jwtSecret))
		require.NoError(t, err)

		sess, err := provider.Authenticate(context.Background(), signedToken)
		require.NoError(t, err)
		require.NotNil(t, sess)

		assert.Equal(t, "usr_supa_12345", sess.IdentityID)
		assert.Equal(t, "doctor@curexal.com", sess.Traits.Email)
		assert.Equal(t, "Dr. Sarah Connor", sess.Traits.Name)
		assert.Equal(t, "https://cdn.curexal.space/avatars/sarah.png", sess.Traits.AvatarURL)
		assert.True(t, sess.Active)
	})

	t.Run("Invalid JWT Signature", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, SupabaseClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "usr_supa_12345",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
			Email: "bad@curexal.com",
		})

		signedToken, err := token.SignedString([]byte("wrong-secret-key"))
		require.NoError(t, err)

		sess, err := provider.Authenticate(context.Background(), signedToken)
		assert.Error(t, err)
		assert.Nil(t, sess)
	})

	t.Run("Empty Token Returns ErrSessionNotFound", func(t *testing.T) {
		sess, err := provider.Authenticate(context.Background(), "")
		assert.ErrorIs(t, err, ErrSessionNotFound)
		assert.Nil(t, sess)
	})
}
