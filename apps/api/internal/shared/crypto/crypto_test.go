package crypto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── HashPassword (Argon2id) ──────────────────────────────────────────────────

func TestHashPassword_ReturnsNonEmptyString(t *testing.T) {
	hash, err := HashPassword("securePassword123!")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestHashPassword_ProducesArgon2idPrefix(t *testing.T) {
	hash, err := HashPassword("securePassword123!")
	require.NoError(t, err)
	expectedPrefix := "$argon2id$v=19$m=65536,t=3,p=2$"
	assert.True(t, strings.HasPrefix(hash, expectedPrefix), "expected Argon2id prefix %s, got: %s", expectedPrefix, hash)
}

func TestHashPassword_TwoHashesOfSamePasswordAreDifferent(t *testing.T) {
	hash1, err1 := HashPassword("same-password")
	hash2, err2 := HashPassword("same-password")
	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2, "Argon2id hashes must differ due to random salt")
}

func TestHashPassword_EmptyPasswordProducesHash(t *testing.T) {
	hash, err := HashPassword("")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

// ─── ComparePassword & VerifyPassword ─────────────────────────────────────────

func TestComparePassword_CorrectPasswordReturnsTrue(t *testing.T) {
	password := "myS3cret!Pass"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	ok := ComparePassword(hash, password)
	assert.True(t, ok)
}

func TestComparePassword_WrongPasswordReturnsFalse(t *testing.T) {
	hash, err := HashPassword("correctPassword")
	require.NoError(t, err)

	ok := ComparePassword(hash, "wrongPassword")
	assert.False(t, ok)
}

func TestComparePassword_EmptyPasswordAgainstHashReturnsFalse(t *testing.T) {
	hash, err := HashPassword("actualPassword")
	require.NoError(t, err)

	ok := ComparePassword(hash, "")
	assert.False(t, ok)
}

func TestComparePassword_InvalidHashReturnsFalse(t *testing.T) {
	ok := ComparePassword("not-a-valid-hash", "password")
	assert.False(t, ok)
}

// ─── Backwards Compatibility with Bcrypt ──────────────────────────────────────

func TestVerifyPassword_LegacyBcryptHashWorks(t *testing.T) {
	password := "password"
	standardBcryptHash := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"
	ok, err := VerifyPassword(password, standardBcryptHash)
	require.NoError(t, err)
	assert.True(t, ok)

	// Wrong password for Bcrypt
	okWrong, errWrong := VerifyPassword("wrongPassword", standardBcryptHash)
	assert.NoError(t, errWrong)
	assert.False(t, okWrong)
}

func TestNeedsRehash(t *testing.T) {
	standardBcryptHash := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"
	assert.True(t, NeedsRehash(standardBcryptHash), "Bcrypt hash must need rehash")

	freshArgon2Hash, err := HashPassword("testPass123!")
	require.NoError(t, err)
	assert.False(t, NeedsRehash(freshArgon2Hash), "Fresh Argon2id hash with standard parameters must NOT need rehash")

	olderArgon2Hash := "$argon2id$v=19$m=32768,t=2,p=1$c2FsdHNhbHQ$aGFzaGhhc2g"
	assert.True(t, NeedsRehash(olderArgon2Hash), "Argon2id hash with lower parameters must need rehash")
}

func TestVerifyPassword_WrongPasswordReturnsFalseNoError(t *testing.T) {
	hash, err := HashPassword("originalPassword")
	require.NoError(t, err)

	ok, err := VerifyPassword("wrongPassword", hash)
	assert.NoError(t, err, "VerifyPassword must return nil error on mismatch")
	assert.False(t, ok)
}

func TestVerifyPassword_InvalidHashReturnsError(t *testing.T) {
	ok, err := VerifyPassword("anypassword", "definitely-not-valid-hash")
	assert.Error(t, err)
	assert.False(t, ok)
}

// ─── Symmetry: HashPassword + VerifyPassword are inverse operations ───────────

func TestHashPassword_VerifyPassword_Symmetry(t *testing.T) {
	passwords := []string{
		"simple",
		"With Spaces 123",
		"!@#$%^&*()_+-=[]{}|;':\",./<>?",
		"ünïcödé",
		"a",
	}

	for _, pw := range passwords {
		t.Run(pw, func(t *testing.T) {
			hash, err := HashPassword(pw)
			require.NoError(t, err)

			ok, err := VerifyPassword(pw, hash)
			require.NoError(t, err)
			assert.True(t, ok, "VerifyPassword must confirm hash produced by HashPassword")

			ok2, err2 := VerifyPassword(pw+"wrong", hash)
			assert.NoError(t, err2)
			assert.False(t, ok2)
		})
	}
}

// ─── GenerateAlphanumericCode ─────────────────────────────────────────────────

func TestGenerateAlphanumericCode(t *testing.T) {
	code, err := GenerateAlphanumericCode(6)
	require.NoError(t, err)
	assert.Len(t, code, 6)

	// Verify unambiguous charset (no 0, O, 1, I)
	for _, ch := range code {
		assert.NotContains(t, []rune{'0', 'O', '1', 'I'}, ch)
	}
}
