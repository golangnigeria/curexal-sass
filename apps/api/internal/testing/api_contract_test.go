// Package testing contains Phase 3 REST API Contract Tests for Curexal.
//
// These tests verify HTTP status codes, request body validation, and JSON response
// schemas for critical endpoints (/api/v1/users/me, /api/v1/auth/sign-in, /api/v1/catalogs).
//
// Tests use net/http/httptest with Echo directly and do NOT require a live database.
package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	catalogsDomain "github.com/golangnigeria/curexal/internal/modules/catalogs/domain"
	identityHandler "github.com/golangnigeria/curexal/internal/modules/identity/handler"
	identityRepo "github.com/golangnigeria/curexal/internal/modules/identity/repository"
	identityService "github.com/golangnigeria/curexal/internal/modules/identity/service"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createContractTestEcho sets up a minimal Echo instance with standard middleware.
func createContractTestEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	return e
}

// createContractTestServer creates a dummy server instance for handler initialization.
func createContractTestServer() *server.Server {
	logger := zerolog.Nop()
	return &server.Server{
		Config: &config.Config{
			Auth: config.AuthConfig{
				SecretKey: "test-secret-key-32-bytes-long!!",
			},
		},
		Logger: &logger,
	}
}

// ─── 1. /api/v1/users/me Contract Tests ───────────────────────────────────────

func TestAPISpec_UsersMe_Unauthenticated_Returns401(t *testing.T) {
	e := createContractTestEcho()
	s := createContractTestServer()

	userRoleHandler := identityHandler.NewUserRoleHandler(s, nil)
	e.GET("/api/v1/users/me", userRoleHandler.GetMe)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var errResp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err, "Response body must be valid JSON")
	assert.Contains(t, errResp, "message")
}

func TestAPISpec_UsersMe_RequireAuth_MiddlewareContract(t *testing.T) {
	e := createContractTestEcho()

	// RequireAuth middleware attached to protected endpoint
	protectedHandler := middleware.RequireAuth()(func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	e.GET("/api/v1/users/me", protectedHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	// Missing authorization context must yield 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

// ─── 2. /api/v1/auth/sign-in Contract Tests ───────────────────────────────────

func TestAPISpec_AuthSignIn_InvalidJSON_Returns400(t *testing.T) {
	e := createContractTestEcho()
	s := createContractTestServer()

	authService := identityService.NewAuthService(s)
	userRepo := identityRepo.NewUserRepository(s)
	authHandler := identityHandler.NewAuthHandler(s, authService, userRepo, nil, nil)
	e.POST("/api/v1/auth/sign-in", authHandler.SignIn)

	// Malformed JSON payload
	body := strings.NewReader(`{ email: "invalid-json", password: `)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err, "Response body must be valid JSON")
	assert.Contains(t, errResp, "message")
	assert.Contains(t, errResp["message"].(string), "Invalid request payload")
}

func TestAPISpec_AuthSignIn_MissingFields_ReturnsError(t *testing.T) {
	e := createContractTestEcho()
	s := createContractTestServer()

	authService := identityService.NewAuthService(s)
	userRepo := identityRepo.NewUserRepository(s)
	authHandler := identityHandler.NewAuthHandler(s, authService, userRepo, nil, nil)
	e.POST("/api/v1/auth/sign-in", authHandler.SignIn)

	// Array JSON instead of expected object payload
	body := strings.NewReader(`[1, 2, 3]`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ─── 3. /api/v1/platform/catalogs Contract Tests ──────────────────────────────
 
func TestAPISpec_MasterCatalogs_ResponseStructureContract(t *testing.T) {
	// Verify catalog response struct serialization contract
	sampleItem := catalogsDomain.CatalogItem{
		Domain:      catalogsDomain.LabDomain,
		Code:        "FBC-001",
		Name:        "Full Blood Count",
		Category:    "Haematology",
		SystemGroup: "Whole Blood",
		BasePrice:   5000.0,
		IsActive:    true,
	}

	jsonBytes, err := json.Marshal(sampleItem)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	require.NoError(t, err)

	// Assert core catalog fields exist in contract JSON response
	assert.Equal(t, "FBC-001", parsed["code"])
	assert.Equal(t, "Full Blood Count", parsed["name"])
	assert.Equal(t, "Haematology", parsed["category"])
	assert.Equal(t, "Whole Blood", parsed["systemGroup"])
	assert.Equal(t, 5000.0, parsed["basePrice"])
	assert.Equal(t, true, parsed["isActive"])
}
