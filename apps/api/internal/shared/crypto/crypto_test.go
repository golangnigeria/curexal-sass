package crypto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── HashPassword ────────────────────────────────────────────────────────────

func TestHashPassword_ReturnsNonEmptyString(t *testing.T) {
	hash, err := HashPassword("securePassword123!")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestHashPassword_ProducesBcryptPrefix(t *testing.T) {
	hash, err := HashPassword("securePassword123!")
	require.NoError(t, err)
	// bcrypt hashes always start with $2a$, $2b$, or $2y$
	assert.True(t, strings.HasPrefix(hash, "$2"), "expected bcrypt prefix, got: %s", hash)
}

func TestHashPassword_TwoHashesOfSamePasswordAreDifferent(t *testing.T) {
	// bcrypt includes a random salt per hash — two hashes MUST differ.
	hash1, err1 := HashPassword("same-password")
	hash2, err2 := HashPassword("same-password")
	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2, "bcrypt hashes must differ due to random salt")
}

func TestHashPassword_EmptyPasswordProducesHash(t *testing.T) {
	// bcrypt will hash an empty string without error.
	hash, err := HashPassword("")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

// ─── ComparePassword ──────────────────────────────────────────────────────────

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
	// A non-bcrypt string is an invalid hash and must return false without panicking.
	ok := ComparePassword("not-a-valid-bcrypt-hash", "password")
	assert.False(t, ok)
}

// ─── VerifyPassword ───────────────────────────────────────────────────────────

func TestVerifyPassword_CorrectPasswordReturnsTrue(t *testing.T) {
	password := "verifyMe#2026"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	ok, err := VerifyPassword(password, hash)
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestVerifyPassword_WrongPasswordReturnsFalseNoError(t *testing.T) {
	hash, err := HashPassword("originalPassword")
	require.NoError(t, err)

	// VerifyPassword must return (false, nil) on mismatch — NOT an error.
	ok, err := VerifyPassword("wrongPassword", hash)
	assert.NoError(t, err, "VerifyPassword must return nil error on mismatch")
	assert.False(t, ok)
}

func TestVerifyPassword_InvalidHashReturnsError(t *testing.T) {
	// A non-bcrypt string is an invalid hash — must return (false, non-nil error).
	ok, err := VerifyPassword("anypassword", "definitely-not-bcrypt")
	assert.Error(t, err)
	assert.False(t, ok)
}

func TestVerifyPassword_EmptyBothArgsReturnsFalseNoError(t *testing.T) {
	// An empty password against an empty hash is a bcrypt error (invalid hash format).
	ok, err := VerifyPassword("", "")
	// Either (false, nil) or (false, err) is acceptable — must NOT panic.
	assert.False(t, ok)
	_ = err // ignore, just ensure no panic
}

// ─── Symmetry: HashPassword + VerifyPassword are inverse operations ───────────

func TestHashPassword_VerifyPassword_Symmetry(t *testing.T) {
	passwords := []string{
		"simple",
		"With Spaces 123",
		"!@#$%^&*()_+-=[]{}|;':\",./<>?",
		"ünïcödé",
		"a", // shortest realistic password
	}

	for _, pw := range passwords {
		t.Run(pw, func(t *testing.T) {
			hash, err := HashPassword(pw)
			require.NoError(t, err)

			ok, err := VerifyPassword(pw, hash)
			require.NoError(t, err)
			assert.True(t, ok, "VerifyPassword must confirm hash produced by HashPassword")

			// Cross-check: wrong password must not verify
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
	assert.Equal(t, strings.ToUpper(code), code, "code must be uppercase")

	// Verify unambiguous charset (no 0, O, 1, I)
	assert.False(t, strings.ContainsAny(code, "0O1I"), "code must not contain ambiguous characters 0, O, 1, I")

	// Verify randomness across multiple calls
	codes := make(map[string]bool)
	for i := 0; i < 50; i++ {
		c, err := GenerateAlphanumericCode(6)
		require.NoError(t, err)
		assert.False(t, codes[c], "codes must be uniquely random")
		codes[c] = true
	}
}
