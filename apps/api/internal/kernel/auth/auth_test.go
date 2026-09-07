package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/labstack/echo/v4"
)

func TestParseSameSite(t *testing.T) {
	tests := []struct {
		input    string
		expected http.SameSite
	}{
		{"Lax", http.SameSiteLaxMode},
		{"lax", http.SameSiteLaxMode},
		{"Strict", http.SameSiteStrictMode},
		{"None", http.SameSiteNoneMode},
		{"Default", http.SameSiteDefaultMode},
		{"unknown", http.SameSiteDefaultMode},
	}

	for _, tt := range tests {
		got := platformAuth.ParseSameSite(tt.input)
		if got != tt.expected {
			t.Errorf("ParseSameSite(%q) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

func TestSetAndClearSessionCookies(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			SecretKey:         "test-secret-key-12345",
			JWTCookieName:     "jwt",
			RefreshCookieName: "refresh_token",
			CookieDomain:      "localhost",
			CookiePath:        "/",
			CookieSecure:      false,
			CookieHTTPOnly:    true,
			CookieSameSite:    "Default",
			JWTExpiration:     15 * time.Minute,
			RefreshExpiration: 30 * 24 * time.Hour,
		},
	}

	// 1. Set Session Cookies
	platformAuth.SetSessionCookies(c, cfg, "access-token-val", "refresh-token-val")
	cookies := rec.Result().Cookies()

	if len(cookies) != 2 {
		t.Fatalf("expected 2 cookies, got %d", len(cookies))
	}

	var jwtCookie, refreshCookie *http.Cookie
	for _, ck := range cookies {
		if ck.Name == "jwt" {
			jwtCookie = ck
		} else if ck.Name == "refresh_token" {
			refreshCookie = ck
		}
	}

	if jwtCookie == nil || jwtCookie.Value != "access-token-val" {
		t.Errorf("unexpected jwt cookie: %v", jwtCookie)
	}
	if jwtCookie.Domain != "" {
		t.Errorf("expected empty cookie domain for local dev, got '%s'", jwtCookie.Domain)
	}
	if refreshCookie == nil || refreshCookie.Value != "refresh-token-val" {
		t.Errorf("unexpected refresh cookie: %v", refreshCookie)
	}

	// 2. Clear Session Cookies
	recClear := httptest.NewRecorder()
	cClear := e.NewContext(req, recClear)
	platformAuth.ClearSessionCookies(cClear, cfg)

	clearCookies := recClear.Result().Cookies()
	if len(clearCookies) != 2 {
		t.Fatalf("expected 2 clear cookies, got %d", len(clearCookies))
	}
	for _, ck := range clearCookies {
		if ck.MaxAge != -1 {
			t.Errorf("expected cookie %s MaxAge -1, got %d", ck.Name, ck.MaxAge)
		}
	}
}

func TestJWTTokenGenerationAndResolutionOrder(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{
		Primary: config.Primary{
			Env: "test",
		},
		Auth: config.AuthConfig{
			SecretKey:         "test-secret-key-12345",
			JWTCookieName:     "jwt",
			RefreshCookieName: "refresh_token",
			CookieDomain:      "localhost",
			CookiePath:        "/",
			CookieSecure:      false,
			CookieHTTPOnly:    true,
			CookieSameSite:    "Default",
			JWTExpiration:     15 * time.Minute,
			RefreshExpiration: 30 * 24 * time.Hour,
			AllowTestHeaders:  true,
		},
	}

	role := "super_admin"
	jwtToken, err := platformAuth.GenerateAccessJWT(cfg, "usr_100", "sess_200", &role, true)
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	// Test 1: Resolution via Cookie
	reqCookie := httptest.NewRequest(http.MethodGet, "/", nil)
	reqCookie.AddCookie(&http.Cookie{Name: "jwt", Value: jwtToken})
	recCookie := httptest.NewRecorder()
	cCookie := e.NewContext(reqCookie, recCookie)

	pCookie := platformAuth.ResolvePrincipal(cCookie, cfg)
	if pCookie == nil || pCookie.UserID != "usr_100" || pCookie.SessionID != "sess_200" {
		t.Errorf("failed cookie resolution: got %v", pCookie)
	}

	// Test 2: Resolution via Bearer Header
	reqBearer := httptest.NewRequest(http.MethodGet, "/", nil)
	reqBearer.Header.Set("Authorization", "Bearer "+jwtToken)
	recBearer := httptest.NewRecorder()
	cBearer := e.NewContext(reqBearer, recBearer)

	pBearer := platformAuth.ResolvePrincipal(cBearer, cfg)
	if pBearer == nil || pBearer.UserID != "usr_100" {
		t.Errorf("failed bearer resolution: got %v", pBearer)
	}

	// Test 3: Resolution via Fallback Test Header
	reqHeader := httptest.NewRequest(http.MethodGet, "/", nil)
	reqHeader.Header.Set("X-User-ID", "usr_test_300")
	recHeader := httptest.NewRecorder()
	cHeader := e.NewContext(reqHeader, recHeader)

	pHeader := platformAuth.ResolvePrincipal(cHeader, cfg)
	if pHeader == nil || pHeader.UserID != "usr_test_300" {
		t.Errorf("failed fallback header resolution: got %v", pHeader)
	}
}

func TestSecurity_XUserRoleHeaderInjectionMustNeverElevatePrivilege(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{
		Auth: config.AuthConfig{
			SecretKey:         "test-secret-key-12345",
			JWTCookieName:     "jwt",
			RefreshCookieName: "refresh_token",
			CookieDomain:      "localhost",
			CookiePath:        "/",
			JWTExpiration:     15 * time.Minute,
			RefreshExpiration: 30 * 24 * time.Hour,
		},
	}

	// Standard member token without platform admin claims
	memberRole := "staff_nurse"
	jwtToken, err := platformAuth.GenerateAccessJWT(cfg, "usr_nurse_1", "sess_nurse_1", &memberRole, false)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Adversary attempts privilege escalation by setting X-User-Role: super_admin
	req := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-User-Role", "super_admin")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	principal := platformAuth.ResolvePrincipal(c, cfg)
	if principal == nil {
		t.Fatalf("expected principal to resolve from valid JWT")
	}

	// ASSERTION: Role MUST remain staff_nurse, and IsSuperAdmin / IsPlatformAdmin MUST be FALSE
	if principal.Platform.IsSuperAdmin {
		t.Errorf("SECURITY VULNERABILITY: principal was granted IsSuperAdmin via X-User-Role header injection!")
	}
	if principal.Platform.IsPlatformAdmin {
		t.Errorf("SECURITY VULNERABILITY: principal was granted IsPlatformAdmin via X-User-Role header injection!")
	}
	if principal.Role == "super_admin" {
		t.Errorf("SECURITY VULNERABILITY: principal role was overwritten with super_admin from client header!")
	}
	if principal.Role != "staff_nurse" {
		t.Errorf("expected role 'staff_nurse', got '%s'", principal.Role)
	}
}

type mockTenantVerifier struct {
	allowedTenantID string
}

func (m *mockTenantVerifier) VerifyMembership(ctx context.Context, userID, tenantID string) (bool, string, error) {
	if tenantID == m.allowedTenantID {
		return true, "doctor", nil
	}
	return false, "", nil
}

func TestSecurity_UnverifiedTenantHeaderIDORProtection(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{
		Auth: config.AuthConfig{
			SecretKey:     "test-secret-key-12345",
			JWTCookieName: "jwt",
			JWTExpiration: 15 * time.Minute,
		},
	}

	doctorRole := "doctor"
	jwtToken, err := platformAuth.GenerateAccessJWT(cfg, "usr_doc_1", "sess_doc_1", &doctorRole, false)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	verifier := &mockTenantVerifier{allowedTenantID: "tenant-legit-clinic-123"}

	// Scenario 1: Accessing authorized tenant
	reqAuthorized := httptest.NewRequest(http.MethodGet, "/", nil)
	reqAuthorized.Header.Set("Authorization", "Bearer "+jwtToken)
	reqAuthorized.Header.Set("X-Tenant-ID", "tenant-legit-clinic-123")
	recAuth := httptest.NewRecorder()
	cAuth := e.NewContext(reqAuthorized, recAuth)

	pAuth := platformAuth.ResolvePrincipalWithVerifier(cAuth, cfg, verifier)
	if pAuth == nil || pAuth.TenantID != "tenant-legit-clinic-123" {
		t.Errorf("expected verified tenant ID 'tenant-legit-clinic-123', got '%s'", pAuth.TenantID)
	}

	// Scenario 2: Attempting unauthorized cross-tenant IDOR access (Tenant B)
	reqAttacker := httptest.NewRequest(http.MethodGet, "/", nil)
	reqAttacker.Header.Set("Authorization", "Bearer "+jwtToken)
	reqAttacker.Header.Set("X-Tenant-ID", "tenant-victim-clinic-999")
	recAttacker := httptest.NewRecorder()
	cAttacker := e.NewContext(reqAttacker, recAttacker)

	pAttacker := platformAuth.ResolvePrincipalWithVerifier(cAttacker, cfg, verifier)
	if pAttacker == nil {
		t.Fatalf("expected principal to resolve")
	}
	if pAttacker.TenantID == "tenant-victim-clinic-999" {
		t.Errorf("SECURITY VULNERABILITY (IDOR): principal granted access to unverified tenant 'tenant-victim-clinic-999'!")
	}
	if pAttacker.TenantID != "" {
		t.Errorf("expected tenant ID to be stripped for unverified membership, got '%s'", pAttacker.TenantID)
	}
}
