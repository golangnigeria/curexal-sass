package testing

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golangnigeria/curexal/internal/bootstrap"
	"github.com/golangnigeria/curexal/internal/modules/identity/repository"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_BrowserPasswordChangeWorkflow tests the complete HTTP API lifecycle:
// Sign-In -> Me Profile -> Change Password (HTTP 200) -> Sign-In with Old Password (401) -> Sign-In with New Password (200).
func TestE2E_BrowserPasswordChangeWorkflow(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil || cfg.Database.DSN() == "" {
		t.Skip("Skipping HTTP E2E integration test; Database DSN is not configured")
	}

	logger := zerolog.Nop()
	srv, err := server.New(cfg, &logger, nil)
	if err != nil {
		t.Skipf("Skipping HTTP E2E integration test; DB connection failed: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	srv.Echo = e

	// Register all real modules & routes
	registeredModules := bootstrap.InitModules(srv)
	registeredModules.RegisterRoutes(srv)

	ctx := context.Background()
	credRepo := repository.NewCredentialRepository(srv)

	testEmail := "e2e_browser_" + uuid.New().String()[:8] + "@curexal.internal"
	oldPassword := "OldP@ssword123!"
	newPassword := "NewP@ssword456!"

	oldHash, err := crypto.HashPassword(oldPassword)
	require.NoError(t, err)

	// Step 1: Create user in database
	userID := uuid.New().String()
	err = srv.DB.Pool.QueryRow(ctx, `
		INSERT INTO identity.users (id, name, email, email_verified, is_platform_admin)
		VALUES ($1, 'E2E Browser Test User', $2, TRUE, FALSE)
		RETURNING id
	`, userID, testEmail).Scan(&userID)
	require.NoError(t, err)

	defer func() {
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.sessions WHERE user_id = $1`, userID)
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.credentials WHERE user_id = $1`, userID)
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.password_histories WHERE user_id = $1`, userID)
		_, _ = srv.DB.Pool.Exec(ctx, `DELETE FROM identity.users WHERE id = $1`, userID)
	}()

	accID := uuid.New().String()
	_, err = srv.DB.Pool.Exec(ctx, `
		INSERT INTO identity.credentials (id, user_id, auth_provider, password_hash)
		VALUES ($1, $2, 'credential', $3)
	`, accID, userID, oldHash)
	require.NoError(t, err)

	// Step 2: Login via HTTP POST /api/v1/auth/sign-in using OLD password
	loginBody, _ := json.Marshal(map[string]string{
		"email":    testEmail,
		"password": oldPassword,
	})
	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", bytes.NewReader(loginBody))
	reqLogin.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recLogin := httptest.NewRecorder()
	e.ServeHTTP(recLogin, reqLogin)

	assert.Equal(t, http.StatusOK, recLogin.Code, "Initial login with old password MUST return 200 OK")
	cookies := recLogin.Result().Cookies()
	require.NotEmpty(t, cookies, "Initial login MUST issue session cookies")

	// Extract jwt cookie
	var jwtCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "jwt" || c.Name == cfg.Auth.JWTCookieName {
			jwtCookie = c
			break
		}
	}

	// Step 3: Call GET /api/v1/users/me with cookies
	reqMe := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	if jwtCookie != nil {
		reqMe.AddCookie(jwtCookie)
	}
	recMe := httptest.NewRecorder()
	e.ServeHTTP(recMe, reqMe)
	assert.Equal(t, http.StatusOK, recMe.Code, "GET /users/me MUST succeed with valid session cookie")

	// Step 4: Change Password via HTTP PUT /api/v1/users/me/password
	changePwBody, _ := json.Marshal(map[string]string{
		"currentPassword": oldPassword,
		"newPassword":     newPassword,
	})
	reqChangePw := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/password", bytes.NewReader(changePwBody))
	reqChangePw.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if jwtCookie != nil {
		reqChangePw.AddCookie(jwtCookie)
	}
	recChangePw := httptest.NewRecorder()
	e.ServeHTTP(recChangePw, reqChangePw)

	assert.Equal(t, http.StatusOK, recChangePw.Code, "Password change HTTP endpoint MUST return 200 OK")

	// Step 5: Verify database immediately post-HTTP change
	var updatedHash string
	_ = srv.DB.Pool.QueryRow(ctx, `SELECT password_hash FROM identity.credentials WHERE user_id = $1 AND auth_provider = 'credential'`, userID).Scan(&updatedHash)
	validNew, _ := crypto.VerifyPassword(newPassword, updatedHash)
	assert.True(t, validNew, "identity.credentials.password_hash MUST match new password after HTTP update")

	// Step 6: Attempt Login via HTTP POST /api/v1/auth/sign-in using OLD password
	oldLoginAttemptBody, _ := json.Marshal(map[string]string{
		"email":    testEmail,
		"password": oldPassword,
	})
	reqOldAttempt := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", bytes.NewReader(oldLoginAttemptBody))
	reqOldAttempt.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recOldAttempt := httptest.NewRecorder()
	e.ServeHTTP(recOldAttempt, reqOldAttempt)

	assert.Equal(t, http.StatusUnauthorized, recOldAttempt.Code, "Login with OLD password MUST return 401 Unauthorized after HTTP password change")

	// Step 7: Attempt Login via HTTP POST /api/v1/auth/sign-in using NEW password
	newLoginAttemptBody, _ := json.Marshal(map[string]string{
		"email":    testEmail,
		"password": newPassword,
	})
	reqNewAttempt := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", bytes.NewReader(newLoginAttemptBody))
	reqNewAttempt.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recNewAttempt := httptest.NewRecorder()
	e.ServeHTTP(recNewAttempt, reqNewAttempt)

	assert.Equal(t, http.StatusOK, recNewAttempt.Code, "Login with NEW password MUST return 200 OK after HTTP password change")
	newCookies := recNewAttempt.Result().Cookies()
	require.NotEmpty(t, newCookies, "New login MUST issue fresh session cookies")

	// Step 8: Verify GetByEmail returns updated hash
	credFinal, err := credRepo.GetByEmail(ctx, testEmail)
	require.NoError(t, err)
	okNew, _ := crypto.VerifyPassword(newPassword, credFinal.PasswordHash)
	assert.True(t, okNew, "GetByEmail MUST verify against new password")
}
