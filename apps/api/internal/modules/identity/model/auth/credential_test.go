package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCredential_IsLocked(t *testing.T) {
	now := time.Now()
	future := now.Add(1 * time.Minute)
	past := now.Add(-10 * time.Minute)

	t.Run("nil credential returns false", func(t *testing.T) {
		var c *Credential
		assert.False(t, c.IsLocked())
	})

	t.Run("locked until future is locked", func(t *testing.T) {
		c := &Credential{
			Status:      CredentialStatusLocked,
			LockedUntil: &future,
		}
		assert.True(t, c.IsLocked())
	})

	t.Run("locked until past is unlocked even if status is LOCKED", func(t *testing.T) {
		c := &Credential{
			Status:      CredentialStatusLocked,
			LockedUntil: &past,
		}
		assert.False(t, c.IsLocked())
	})

	t.Run("nil locked until with active status is unlocked", func(t *testing.T) {
		c := &Credential{
			Status:      CredentialStatusActive,
			LockedUntil: nil,
		}
		assert.False(t, c.IsLocked())
	})

	t.Run("nil locked until with locked status is locked (admin lockout)", func(t *testing.T) {
		c := &Credential{
			Status:      CredentialStatusLocked,
			LockedUntil: nil,
		}
		assert.True(t, c.IsLocked())
	})
}
