package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidCiphertext = errors.New("invalid ciphertext or corrupted payload")
	ErrInvalidMasterKey  = errors.New("invalid master encryption key")
)

// DeriveKey ensures the master key string produces a 32-byte key suitable for AES-256-GCM.
func DeriveKey(masterKey string) []byte {
	if masterKey == "" {
		masterKey = "CUREXAL_DEFAULT_PLATFORM_AES_MASTER_KEY_2026"
	}
	hash := sha256.Sum256([]byte(masterKey))
	return hash[:]
}

// EncryptAEAD encrypts plaintext using AES-256-GCM AEAD encryption.
func EncryptAEAD(plaintext string, masterKey string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	key := DeriveKey(masterKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM block: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}

// DecryptAEAD decrypts ciphertext using AES-256-GCM AEAD decryption.
func DecryptAEAD(ciphertextHex string, masterKey string) (string, error) {
	if ciphertextHex == "" {
		return "", nil
	}

	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", ErrInvalidCiphertext
	}

	key := DeriveKey(masterKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM block: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce, encryptedMessage := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		return "", ErrInvalidCiphertext
	}

	return string(plaintext), nil
}

// RedactSecret returns a standard masked string for secrets exposed in API responses.
func RedactSecret(secret string) string {
	if secret == "" {
		return ""
	}
	return "••••••••"
}
