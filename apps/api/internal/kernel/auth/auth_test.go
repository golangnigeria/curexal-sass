package auth_test

import (
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
