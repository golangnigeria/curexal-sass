package crypto_test

import (
	"testing"

	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/stretchr/testify/assert"
)

func TestAEAD_EncryptAndDecrypt_Success(t *testing.T) {
	masterKey := "super_secret_platform_key_123"
	secretPayload := "sk_live_paystack_secret_key_998877"

	ciphertext, err := crypto.EncryptAEAD(secretPayload, masterKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, ciphertext)
	assert.NotEqual(t, secretPayload, ciphertext)

	decrypted, err := crypto.DecryptAEAD(ciphertext, masterKey)
	assert.NoError(t, err)
	assert.Equal(t, secretPayload, decrypted)
}

func TestAEAD_Decrypt_CorruptedCiphertext(t *testing.T) {
	masterKey := "super_secret_platform_key_123"
	corrupted := "deadbeefbadc0de"

	decrypted, err := crypto.DecryptAEAD(corrupted, masterKey)
	assert.Error(t, err)
	assert.Empty(t, decrypted)
}

func TestRedactSecret(t *testing.T) {
	secret := "sk_live_stripe_secret_key_abc"
	redacted := crypto.RedactSecret(secret)
	assert.Equal(t, "••••••••", redacted)
}
