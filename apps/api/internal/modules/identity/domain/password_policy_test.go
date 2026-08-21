package domain_test

import (
	"testing"

	"github.com/golangnigeria/curexal/internal/modules/identity/domain"
	"github.com/golangnigeria/curexal/internal/shared/crypto"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		isPlatformUser bool
		wantErr        bool
	}{
		{
			name:           "valid general user password",
			password:       "SecureP@ss1",
			isPlatformUser: false,
			wantErr:        false,
		},
		{
			name:           "general user password too short",
			password:       "P@ss1",
			isPlatformUser: false,
			wantErr:        true,
		},
		{
			name:           "platform user password valid (12+ chars)",
			password:       "SuperSecur3P@ssword",
			isPlatformUser: true,
			wantErr:        false,
		},
		{
			name:           "platform user password too short (10 chars)",
			password:       "Secur3P@ss1",
			isPlatformUser: true,
			wantErr:        true,
		},
		{
			name:           "missing uppercase",
			password:       "securep@ss123",
			isPlatformUser: false,
			wantErr:        true,
		},
		{
			name:           "missing lowercase",
			password:       "SECUREP@SS123",
			isPlatformUser: false,
			wantErr:        true,
		},
		{
			name:           "missing digit",
			password:       "SecureP@ssword",
			isPlatformUser: false,
			wantErr:        true,
		},
		{
			name:           "missing special char",
			password:       "SecurePassword123",
			isPlatformUser: false,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidatePassword(tt.password, tt.isPlatformUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%q, %v) error = %v, wantErr %v", tt.password, tt.isPlatformUser, err, tt.wantErr)
			}
		})
	}
}

func TestCheckPasswordHistory(t *testing.T) {
	oldPassword := "OldP@ssword123"
	oldHash, err := crypto.HashPassword(oldPassword)
	if err != nil {
		t.Fatalf("failed to hash old password: %v", err)
	}

	hashes := []string{oldHash}

	t.Run("rejects password in history", func(t *testing.T) {
		err := domain.CheckPasswordHistory("OldP@ssword123", hashes)
		if err == nil {
			t.Errorf("expected ErrPasswordReuse, got nil")
		}
	})

	t.Run("allows new password not in history", func(t *testing.T) {
		err := domain.CheckPasswordHistory("BrandNewP@ss123", hashes)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
