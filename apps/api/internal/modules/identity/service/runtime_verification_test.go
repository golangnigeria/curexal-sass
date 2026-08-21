package service_test

import (
	"context"
	"testing"

	"github.com/golangnigeria/curexal/internal/modules/identity/repository"
	"github.com/golangnigeria/curexal/internal/modules/identity/service"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuntimeVerification_PasswordLifecycle(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil || cfg.Database.DSN() == "" {
		t.Skip("Skipping DB integration test; Database DSN is not set or config load failed")
	}

	ctx := context.Background()
	logger := zerolog.Nop()
	srv, err := server.New(cfg, &logger, nil)
	if err != nil {
		t.Skipf("Skipping integration test; failed to connect to database: %v", err)
	}

	credRepo := repository.NewCredentialRepository(srv)
	authSvc := service.NewAuthService(srv)

	testEmail := "runtime_audit_" + uuid.New().String()[:8] + "@curexal.internal"
	oldPassword := "OldP@ssword123!"
	newPassword := "NewP@ssword456!"

	oldHash, err := crypto.HashPassword(oldPassword)
	require.NoError(t, err)

	// Create test user in identity.users
	userID := uuid.New().String()
	_, err = srv.DB.Pool.Exec(ctx, `
		INSERT INTO identity.users (id, name, email, email_verified, is_platform_admin)
		VALUES ($1, 'Runtime Verification User', $2, TRUE, FALSE)
	`, userID, testEmail)
	require.NoError(t, err)

	defer func() {
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.password_histories WHERE user_id = $1`, userID)
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.credentials WHERE user_id = $1`, userID)
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.users WHERE id = $1`, userID)
	}()

	// Insert active credential record in identity.credentials
	accID := uuid.New().String()
	_, err = srv.DB.Pool.Exec(ctx, `
		INSERT INTO identity.credentials (id, account_id, auth_provider, user_id, password_hash)
		VALUES ($1, $2, 'credential', $3, $4)
	`, accID, testEmail, userID, oldHash)
	require.NoError(t, err)

	// Step 1: Inspect DB before password change
	credBefore, err := credRepo.GetByEmail(ctx, testEmail)
	require.NoError(t, err)
	assert.Equal(t, userID, credBefore.UserID)
	okOld, _ := crypto.VerifyPassword(oldPassword, credBefore.PasswordHash)
	assert.True(t, okOld, "Initial hash must match old password")

	// Step 2 & 3: Change Password and verify DB immediately
	newHash, err := crypto.HashPassword(newPassword)
	require.NoError(t, err)

	auditRecord := &repository.PasswordHistoryRecord{
		ChangedBy:    userID,
		ChangeReason: "PASSWORD_CHANGE",
		IPAddress:    "127.0.0.1",
		UserAgent:    "RuntimeVerification/1.0",
	}

	err = credRepo.UpdatePasswordHashWithAudit(ctx, userID, newHash, auditRecord)
	require.NoError(t, err)

	// Step 3: Verify single active credential contains the new hash
	var activeHash string
	err = srv.DB.Pool.QueryRow(ctx, `
		SELECT password_hash 
		FROM identity.credentials 
		WHERE user_id = $1 AND auth_provider = 'credential'
	`, userID).Scan(&activeHash)
	require.NoError(t, err)

	assert.NotEqual(t, oldHash, activeHash, "Stored active hash MUST differ from initial old hash")
	okNewHash, _ := crypto.VerifyPassword(newPassword, activeHash)
	assert.True(t, okNewHash, "Active credential hash MUST verify new password")

	// Verify history record contains old password hash
	history, err := credRepo.GetPasswordHistory(ctx, userID, 5)
	require.NoError(t, err)
	assert.Len(t, history, 1, "Password history MUST record the previous password hash")
	okHistoryOld, _ := crypto.VerifyPassword(oldPassword, history[0])
	assert.True(t, okHistoryOld, "History MUST contain the old password hash")

	// Step 4 — Login Tests
	// Test A: Old password login MUST FAIL
	_, errOldLogin := authSvc.SignInCredentials(ctx, testEmail, oldPassword, "127.0.0.1", "RuntimeVerification/1.0")
	assert.Error(t, errOldLogin, "Login with old password MUST return an authentication error")

	// Test B: New password login MUST SUCCEED
	userNewLogin, errNewLogin := authSvc.SignInCredentials(ctx, testEmail, newPassword, "127.0.0.1", "RuntimeVerification/1.0")
	assert.NoError(t, errNewLogin, "Login with new password MUST succeed")
	assert.Equal(t, userID, userNewLogin.ID)

	// Test C: Attempt old password 3 consecutive times & check failed login counter
	for i := 0; i < 3; i++ {
		_, _ = authSvc.SignInCredentials(ctx, testEmail, oldPassword, "127.0.0.1", "RuntimeVerification/1.0")
	}
	credAfterFailed, err := credRepo.GetByEmail(ctx, testEmail)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, credAfterFailed.FailedLoginAttempts, 3, "Failed login counter MUST increment on failed attempts")

	// Test D: New password login clears failed login counter
	userPostFailedLogin, errPostFailed := authSvc.SignInCredentials(ctx, testEmail, newPassword, "127.0.0.1", "RuntimeVerification/1.0")
	assert.NoError(t, errPostFailed, "Login with new password post failed attempts MUST succeed")
	assert.Equal(t, userID, userPostFailedLogin.ID)

	credReset, err := credRepo.GetByEmail(ctx, testEmail)
	require.NoError(t, err)
	assert.Equal(t, 0, credReset.FailedLoginAttempts, "Successful login MUST reset failed login counter to 0")

	// Step 5: Trace credential query result
	credAfter, err := credRepo.GetByEmail(ctx, testEmail)
	require.NoError(t, err)
	okNew, _ := crypto.VerifyPassword(newPassword, credAfter.PasswordHash)
	assert.True(t, okNew, "GetByEmail MUST return the updated new password hash")

	// Step 6: Verify zero duplicate credential records exist
	var duplicateCount int
	err = srv.DB.Pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM (
			SELECT auth_provider
			FROM identity.credentials
			WHERE user_id = $1
			GROUP BY auth_provider
			HAVING COUNT(*) > 1
		) dupes
	`, userID).Scan(&duplicateCount)
	require.NoError(t, err)
	assert.Equal(t, 0, duplicateCount, "Must have ZERO duplicate credential records")
}

func TestRuntimeVerification_PasswordRequestLifecycle(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil || cfg.Database.DSN() == "" {
		t.Skip("Skipping DB integration test; Database DSN is not set or config load failed")
	}

	ctx := context.Background()
	logger := zerolog.Nop()
	srv, err := server.New(cfg, &logger, nil)
	if err != nil {
		t.Skipf("Skipping integration test; failed to connect to database: %v", err)
	}

	credRepo := repository.NewCredentialRepository(srv)
	authSvc := service.NewAuthService(srv)

	testEmail := "pw_request_" + uuid.New().String()[:8] + "@curexal.internal"
	initialPassword := "InitP@ssword123!"

	initHash, err := crypto.HashPassword(initialPassword)
	require.NoError(t, err)

	// 1. Create verified test user
	userID := uuid.New().String()
	_, err = srv.DB.Pool.Exec(ctx, `
		INSERT INTO identity.users (id, name, email, email_verified, is_platform_admin)
		VALUES ($1, 'Password Request Test User', $2, TRUE, FALSE)
	`, userID, testEmail)
	require.NoError(t, err)

	defer func() {
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.password_requests WHERE user_id = $1`, userID)
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.password_histories WHERE user_id = $1`, userID)
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.credentials WHERE user_id = $1`, userID)
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.sessions WHERE user_id = $1`, userID)
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.users WHERE id = $1`, userID)
	}()

	// 2. Insert initial credentials
	accID := uuid.New().String()
	_, err = srv.DB.Pool.Exec(ctx, `
		INSERT INTO identity.credentials (id, account_id, auth_provider, user_id, password_hash)
		VALUES ($1, $2, 'credential', $3, $4)
	`, accID, testEmail, userID, initHash)
	require.NoError(t, err)

	// 3. Create active session to test session revocation
	sessID := "sess_" + uuid.New().String()[:16]
	_, err = srv.DB.Pool.Exec(ctx, `
		INSERT INTO identity.sessions (id, user_id, token, expires_at)
		VALUES ($1, $2, $3, NOW() + INTERVAL '1 hour')
	`, sessID, userID, "token_"+sessID)
	require.NoError(t, err)

	// Verify session exists before request
	var sessCountBefore int
	_ = srv.DB.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM identity.sessions WHERE user_id = $1`, userID).Scan(&sessCountBefore)
	assert.Equal(t, 1, sessCountBefore, "Active session must exist prior to password request")

	// 4. Execute Password Request for valid verified user
	err = authSvc.RequestPassword(ctx, testEmail, "https://curexal.space", "127.0.0.1", "TestAgent/1.0")
	require.NoError(t, err, "RequestPassword must succeed without error")

	// 5. Verify identity.password_requests record
	var reqCount int
	var reqStatus string
	err = srv.DB.Pool.QueryRow(ctx, `
		SELECT COUNT(*), COALESCE(status, '')
		FROM identity.password_requests
		WHERE user_id = $1 AND email = $2
		GROUP BY status
	`, userID, testEmail).Scan(&reqCount, &reqStatus)
	require.NoError(t, err)
	assert.Equal(t, 1, reqCount, "Password request MUST be recorded in identity.password_requests")
	assert.Equal(t, "DELIVERED", reqStatus)

	// 6. Verify password was changed and old password is now invalid
	credAfter, err := credRepo.GetByEmail(ctx, testEmail)
	require.NoError(t, err)
	assert.NotEqual(t, initHash, credAfter.PasswordHash, "Active password hash MUST be updated")

	_, errOldLogin := authSvc.SignInCredentials(ctx, testEmail, initialPassword, "127.0.0.1", "TestAgent/1.0")
	assert.Error(t, errOldLogin, "Old password MUST no longer work for sign-in")

	// 7. Verify all prior active sessions were revoked
	var sessCountAfter int
	_ = srv.DB.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM identity.sessions WHERE user_id = $1`, userID).Scan(&sessCountAfter)
	assert.Equal(t, 0, sessCountAfter, "All active sessions MUST be revoked upon password request")

	// 8. Verify anti-enumeration: Request password for unregistered email returns nil without inserting records
	errAnon := authSvc.RequestPassword(ctx, "nonexistent_anon_probe@curexal.internal", "", "127.0.0.1", "TestAgent/1.0")
	assert.NoError(t, errAnon, "RequestPassword for non-existent email MUST return nil (anti-enumeration)")

	var anonReqCount int
	_ = srv.DB.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM identity.password_requests WHERE email = $1`, "nonexistent_anon_probe@curexal.internal").Scan(&anonReqCount)
	assert.Equal(t, 0, anonReqCount, "No record should be inserted for non-existent email")

	// 9. Verify unverified email protection
	unverifiedEmail := "unverified_" + uuid.New().String()[:8] + "@curexal.internal"
	unverifiedID := uuid.New().String()
	_, err = srv.DB.Pool.Exec(ctx, `
		INSERT INTO identity.users (id, name, email, email_verified, is_platform_admin)
		VALUES ($1, 'Unverified User', $2, FALSE, FALSE)
	`, unverifiedID, unverifiedEmail)
	require.NoError(t, err)

	defer func() {
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.users WHERE id = $1`, unverifiedID)
	}()

	errUnverified := authSvc.RequestPassword(ctx, unverifiedEmail, "", "127.0.0.1", "TestAgent/1.0")
	assert.NoError(t, errUnverified, "RequestPassword for unverified user MUST return nil safely")

	var unverifiedReqCount int
	_ = srv.DB.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM identity.password_requests WHERE user_id = $1`, unverifiedID).Scan(&unverifiedReqCount)
	assert.Equal(t, 0, unverifiedReqCount, "Unverified user MUST NOT have password requests processed")
}

