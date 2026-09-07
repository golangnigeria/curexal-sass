package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/labstack/echo/v4"
)

// resolveEffectiveCookieDomain dynamically determines the cookie domain from the request host.
// In the organization-centric model, host-only cookies (Domain = "") provide the highest
// security and isolation, preventing session bleeding across tenant subdomains and custom domains.
func resolveEffectiveCookieDomain(c echo.Context, cfg *config.Config) string {
	cleanHost := ExtractCleanHost(c.Request())

	if cleanHost == "" || cleanHost == "localhost" || cleanHost == "127.0.0.1" || strings.HasSuffix(cleanHost, ".localhost") {
		return ""
	}

	// For explicit configured cookie domain override (if not wildcard localhost)
	if cfg.Auth.CookieDomain != "" && !strings.Contains(cfg.Auth.CookieDomain, "localhost") && !strings.Contains(cfg.Auth.CookieDomain, ".curexal.space") {
		return cfg.Auth.CookieDomain
	}

	// Default to host-only cookie for strict tenant domain isolation
	return ""
}

// SetSessionCookies sets session and refresh token cookies according to AuthConfig.
func SetSessionCookies(c echo.Context, cfg *config.Config, accessToken, refreshToken string) {
	sameSite := ParseSameSite(cfg.Auth.CookieSameSite)
	domain := resolveEffectiveCookieDomain(c, cfg)

	jwtCookie := &http.Cookie{
		Name:     cfg.Auth.JWTCookieName,
		Value:    accessToken,
		Path:     cfg.Auth.CookiePath,
		Expires:  time.Now().Add(cfg.Auth.JWTExpiration),
		HttpOnly: cfg.Auth.CookieHTTPOnly,
		Secure:   cfg.Auth.CookieSecure,
		SameSite: sameSite,
		Domain:   domain,
	}
	c.SetCookie(jwtCookie)

	refreshCookie := &http.Cookie{
		Name:     cfg.Auth.RefreshCookieName,
		Value:    refreshToken,
		Path:     cfg.Auth.CookiePath,
		Expires:  time.Now().Add(cfg.Auth.RefreshExpiration),
		HttpOnly: cfg.Auth.CookieHTTPOnly,
		Secure:   cfg.Auth.CookieSecure,
		SameSite: sameSite,
		Domain:   domain,
	}
	c.SetCookie(refreshCookie)
}

// SetAccessJWTCookie updates only the access JWT cookie (e.g. after branch switch).
func SetAccessJWTCookie(c echo.Context, cfg *config.Config, accessToken string) {
	sameSite := ParseSameSite(cfg.Auth.CookieSameSite)
	domain := resolveEffectiveCookieDomain(c, cfg)

	jwtCookie := &http.Cookie{
		Name:     cfg.Auth.JWTCookieName,
		Value:    accessToken,
		Path:     cfg.Auth.CookiePath,
		Expires:  time.Now().Add(cfg.Auth.JWTExpiration),
		HttpOnly: cfg.Auth.CookieHTTPOnly,
		Secure:   cfg.Auth.CookieSecure,
		SameSite: sameSite,
		Domain:   domain,
	}
	c.SetCookie(jwtCookie)
}

// ClearSessionCookies clears session cookies according to AuthConfig.
func ClearSessionCookies(c echo.Context, cfg *config.Config) {
	sameSite := ParseSameSite(cfg.Auth.CookieSameSite)
	domain := resolveEffectiveCookieDomain(c, cfg)

	names := []string{cfg.Auth.JWTCookieName, cfg.Auth.RefreshCookieName}
	for _, name := range names {
		cookie := &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     cfg.Auth.CookiePath,
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: cfg.Auth.CookieHTTPOnly,
			Secure:   cfg.Auth.CookieSecure,
			SameSite: sameSite,
			Domain:   domain,
		}
		c.SetCookie(cookie)
	}
}
