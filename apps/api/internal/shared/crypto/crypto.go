package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

// Argon2id configuration complying with HIPAA, OWASP, and Day 1 Production Specification:
// Memory: 64 MB (65536 KiB), Iterations: 3, Parallelism: 2 threads, Salt: 16 bytes, Key: 32 bytes.
const (
	argon2Memory      = 64 * 1024 // 64 MB
	argon2Iterations  = 3
	argon2Parallelism = 2
	argon2SaltLen     = 16
	argon2KeyLen      = 32
)

// HashPassword hashes a plain text password using Argon2id.
// Output format: $argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
func HashPassword(password string) (string, error) {
	salt := make([]byte, argon2SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate cryptographic salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, argon2Iterations, argon2Memory, argon2Parallelism, argon2KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argon2Memory, argon2Iterations, argon2Parallelism, b64Salt, b64Hash)

	return encoded, nil
}

// VerifyPassword verifies a plain text password against either an Argon2id or Bcrypt hash.
func VerifyPassword(password, encodedHash string) (bool, error) {
	if strings.HasPrefix(encodedHash, "$argon2id$") {
		return verifyArgon2id(password, encodedHash)
	}
	if strings.HasPrefix(encodedHash, "$2a$") || strings.HasPrefix(encodedHash, "$2b$") || strings.HasPrefix(encodedHash, "$2y$") {
		err := bcrypt.CompareHashAndPassword([]byte(encodedHash), []byte(password))
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
	return false, errors.New("unrecognized password hash algorithm")
}

// ComparePassword compares a hashed password with a plain text password.
func ComparePassword(hashedPassword, password string) bool {
	ok, err := VerifyPassword(password, hashedPassword)
	return err == nil && ok
}

// NeedsRehash returns true if the stored hash is not using the latest Argon2id parameters.
func NeedsRehash(encodedHash string) bool {
	if !strings.HasPrefix(encodedHash, "$argon2id$") {
		return true
	}
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return true
	}
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return true
	}
	var memory, iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return true
	}
	return memory != argon2Memory || iterations != argon2Iterations || parallelism != argon2Parallelism
}

func verifyArgon2id(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid argon2id hash format")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return false, errors.New("incompatible argon2 version")
	}

	var memory, iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return false, errors.New("invalid argon2id parameters")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, errors.New("invalid salt encoding")
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, errors.New("invalid hash encoding")
	}

	comparisonHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))

	if subtle.ConstantTimeCompare(hash, comparisonHash) == 1 {
		return true, nil
	}
	return false, nil
}

// GenerateAlphanumericCode generates a cryptographically secure uppercase alphanumeric verification code.
// Uses a 32-character unambiguous charset (excluding 0, O, 1, I).
func GenerateAlphanumericCode(length int) (string, error) {
	if length <= 0 {
		length = 6
	}
	const charset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate secure random code: %w", err)
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}
