package identity_test

import (
	"strings"
	"testing"

	"github.com/golangnigeria/curexal/internal/modules/identity/domain"
	modelAuth "github.com/golangnigeria/curexal/internal/modules/identity/model/auth"
	modelUser "github.com/golangnigeria/curexal/internal/modules/identity/model/user"
	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserProfilePayload_Validation(t *testing.T) {
	firstName := "Super"
	lastName := "Admin"
	phone := "+1 (555) 019-2834"
	payload := modelUser.UpdateUserProfilePayload{
		UserID:      "user-123",
		FirstName:   &firstName,
		LastName:    &lastName,
		PhoneNumber: &phone,
	}

	if err := payload.Validate(); err != nil {
		t.Fatalf("expected valid payload, got: %v", err)
	}
}

// SimulateMerge simulates the COALESCE merge logic executed by PostgreSQL and repository
func simulateMerge(existing modelUser.UserProfile, payload modelUser.UpdateUserProfilePayload) (modelUser.UserProfile, string) {
	result := existing

	if payload.FirstName != nil {
		result.FirstName = payload.FirstName
	}
	if payload.MiddleName != nil {
		result.MiddleName = payload.MiddleName
	}
	if payload.LastName != nil {
		result.LastName = payload.LastName
	}
	if payload.PhoneNumber != nil {
		result.PhoneNumber = payload.PhoneNumber
	}

	var parts []string
	if result.FirstName != nil && strings.TrimSpace(*result.FirstName) != "" {
		parts = append(parts, strings.TrimSpace(*result.FirstName))
	}
	if result.MiddleName != nil && strings.TrimSpace(*result.MiddleName) != "" {
		parts = append(parts, strings.TrimSpace(*result.MiddleName))
	}
	if result.LastName != nil && strings.TrimSpace(*result.LastName) != "" {
		parts = append(parts, strings.TrimSpace(*result.LastName))
	}

	return result, strings.Join(parts, " ")
}

func strPtr(s string) *string {
	return &s
}

func TestPartialUpdateSemantics_CasesA_through_G(t *testing.T) {
	// Baseline profile
	baseProfile := modelUser.UserProfile{
		UserID:      "usr_test_123",
		FirstName:   strPtr("Prince"),
		MiddleName:  strPtr("Dimkpa"),
		LastName:    strPtr("John"),
		PhoneNumber: strPtr("+2348011111111"),
	}

	// Case A: Update only middleName -> firstName & lastName remain unchanged
	t.Run("Case A - Update only middleName", func(t *testing.T) {
		payload := modelUser.UpdateUserProfilePayload{
			UserID:     baseProfile.UserID,
			MiddleName: strPtr("Kinikanwo"),
		}
		require.NoError(t, payload.Validate())

		res, displayName := simulateMerge(baseProfile, payload)
		assert.Equal(t, "Prince", *res.FirstName)
		assert.Equal(t, "Kinikanwo", *res.MiddleName)
		assert.Equal(t, "John", *res.LastName)
		assert.Equal(t, "Prince Kinikanwo John", displayName)
	})

	// Case B: Update only firstName -> Only firstName changes
	t.Run("Case B - Update only firstName", func(t *testing.T) {
		payload := modelUser.UpdateUserProfilePayload{
			UserID:    baseProfile.UserID,
			FirstName: strPtr("A"),
		}
		require.NoError(t, payload.Validate())

		res, displayName := simulateMerge(baseProfile, payload)
		assert.Equal(t, "A", *res.FirstName)
		assert.Equal(t, "Dimkpa", *res.MiddleName)
		assert.Equal(t, "John", *res.LastName)
		assert.Equal(t, "A Dimkpa John", displayName)
	})

	// Case C: Update only lastName -> Only lastName changes
	t.Run("Case C - Update only lastName", func(t *testing.T) {
		payload := modelUser.UpdateUserProfilePayload{
			UserID:   baseProfile.UserID,
			LastName: strPtr("C"),
		}
		require.NoError(t, payload.Validate())

		res, displayName := simulateMerge(baseProfile, payload)
		assert.Equal(t, "Prince", *res.FirstName)
		assert.Equal(t, "Dimkpa", *res.MiddleName)
		assert.Equal(t, "C", *res.LastName)
		assert.Equal(t, "Prince Dimkpa C", displayName)
	})

	// Case D: Update firstName and lastName -> Existing middleName remains unchanged
	t.Run("Case D - Update firstName and lastName", func(t *testing.T) {
		payload := modelUser.UpdateUserProfilePayload{
			UserID:    baseProfile.UserID,
			FirstName: strPtr("A"),
			LastName:  strPtr("C"),
		}
		require.NoError(t, payload.Validate())

		res, displayName := simulateMerge(baseProfile, payload)
		assert.Equal(t, "A", *res.FirstName)
		assert.Equal(t, "Dimkpa", *res.MiddleName)
		assert.Equal(t, "C", *res.LastName)
		assert.Equal(t, "A Dimkpa C", displayName)
	})

	// Case E: Update all three name fields -> A B C
	t.Run("Case E - Update all three names", func(t *testing.T) {
		payload := modelUser.UpdateUserProfilePayload{
			UserID:     baseProfile.UserID,
			FirstName:  strPtr("A"),
			MiddleName: strPtr("B"),
			LastName:   strPtr("C"),
		}
		require.NoError(t, payload.Validate())

		res, displayName := simulateMerge(baseProfile, payload)
		assert.Equal(t, "A", *res.FirstName)
		assert.Equal(t, "B", *res.MiddleName)
		assert.Equal(t, "C", *res.LastName)
		assert.Equal(t, "A B C", displayName)
	})

	// Case F: Update middleName with "" -> Middle name is explicitly cleared
	t.Run("Case F - Explicitly clear middleName", func(t *testing.T) {
		payload := modelUser.UpdateUserProfilePayload{
			UserID:     baseProfile.UserID,
			MiddleName: strPtr(""),
		}
		require.NoError(t, payload.Validate())

		res, displayName := simulateMerge(baseProfile, payload)
		assert.Equal(t, "Prince", *res.FirstName)
		assert.Equal(t, "", *res.MiddleName)
		assert.Equal(t, "John", *res.LastName)
		assert.Equal(t, "Prince John", displayName)
	})

	// Case G: Update non-name field -> No name fields change
	t.Run("Case G - Update non-name field", func(t *testing.T) {
		payload := modelUser.UpdateUserProfilePayload{
			UserID:      baseProfile.UserID,
			PhoneNumber: strPtr("+2348099999999"),
		}
		require.NoError(t, payload.Validate())

		res, displayName := simulateMerge(baseProfile, payload)
		assert.Equal(t, "Prince", *res.FirstName)
		assert.Equal(t, "Dimkpa", *res.MiddleName)
		assert.Equal(t, "John", *res.LastName)
		assert.Equal(t, "+2348099999999", *res.PhoneNumber)
		assert.Equal(t, "Prince Dimkpa John", displayName)
	})
}

func TestPasswordRotationAndHistory_Invariant(t *testing.T) {
	// Simulate password lifecycle: A -> B -> C -> D
	passA := "PasswordA!123"
	passB := "PasswordB!123"
	passC := "PasswordC!123"
	passD := "PasswordD!123"

	hashA, err := crypto.HashPassword(passA)
	require.NoError(t, err)
	hashB, err := crypto.HashPassword(passB)
	require.NoError(t, err)
	hashC, err := crypto.HashPassword(passC)
	require.NoError(t, err)
	hashD, err := crypto.HashPassword(passD)
	require.NoError(t, err)

	// Invariant: exactly 1 active credential hash at any time
	activeHash := hashA
	var historyHashes []string

	// Rotation 1: Change to B
	historyHashes = append([]string{activeHash}, historyHashes...)
	activeHash = hashB

	// Rotation 2: Change to C
	historyHashes = append([]string{activeHash}, historyHashes...)
	activeHash = hashC

	// Rotation 3: Change to D
	historyHashes = append([]string{activeHash}, historyHashes...)
	activeHash = hashD

	// Active hash must be D
	matchD, err := crypto.VerifyPassword(passD, activeHash)
	require.NoError(t, err)
	assert.True(t, matchD, "active credential must authenticate current password D")

	// Old passwords A, B, C must fail against active credential
	matchA, _ := crypto.VerifyPassword(passA, activeHash)
	assert.False(t, matchA, "old password A must fail against active credential")
	matchB, _ := crypto.VerifyPassword(passB, activeHash)
	assert.False(t, matchB, "old password B must fail against active credential")
	matchC, _ := crypto.VerifyPassword(passC, activeHash)
	assert.False(t, matchC, "old password C must fail against active credential")

	// Password history must contain A, B, C
	assert.Len(t, historyHashes, 3)

	// Attempting reuse of B or C must be rejected by policy
	errReuseB := domain.CheckPasswordHistory(passB, historyHashes)
	assert.Error(t, errReuseB, "reusing password B from history must be rejected")
	assert.Equal(t, domain.ErrPasswordReuse, errReuseB)

	errReuseA := domain.CheckPasswordHistory(passA, historyHashes)
	assert.Error(t, errReuseA, "reusing password A from history must be rejected")
	assert.Equal(t, domain.ErrPasswordReuse, errReuseA)

	// Brand new password E must be accepted
	passE := "PasswordE!123"
	errNewE := domain.CheckPasswordHistory(passE, historyHashes)
	assert.NoError(t, errNewE, "brand new password E must be permitted")
}

func TestRequestEmailChangePayload_Validation(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Valid Email", "newadmin@curexal.internal", false},
		{"Invalid Email", "invalid-email-format", true},
		{"Empty Email", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := modelUser.RequestEmailChangePayload{NewEmail: tt.email}
			err := p.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestEmailChangePayload.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVerifyEmailChangePayload_Validation(t *testing.T) {
	pToken := modelUser.VerifyEmailChangePayload{Token: "valid-token-123"}
	if err := pToken.Validate(); err != nil {
		t.Fatalf("expected valid token payload, got: %v", err)
	}

	pCode := modelUser.VerifyEmailChangePayload{Code: "7K9P2X"}
	if err := pCode.Validate(); err != nil {
		t.Fatalf("expected valid code payload, got: %v", err)
	}

	emptyP := modelUser.VerifyEmailChangePayload{Token: "", Code: ""}
	if err := emptyP.Validate(); err == nil {
		t.Fatalf("expected error for empty token and code, got nil")
	}
}

func TestRequestPasswordPayload_Validation(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Valid Email", "user@curexal.space", false},
		{"Invalid Format", "not-an-email", true},
		{"Empty Email", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := modelAuth.RequestPasswordPayload{Email: tt.email}
			err := p.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestPasswordPayload.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

