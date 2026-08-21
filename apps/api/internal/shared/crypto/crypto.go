package crypto

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using Bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// ComparePassword compares a hashed password with a plain text password.
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// VerifyPassword verifies a password against a bcrypt hash.
func VerifyPassword(password, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
