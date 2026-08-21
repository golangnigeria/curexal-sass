package auth

import (
	"net/http"
	"time"

	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/labstack/echo/v4"
)

// SetSessionCookies sets session and refresh token cookies according to AuthConfig.
func SetSessionCookies(c echo.Context, cfg *config.Config, accessToken, refreshToken string) {
	sameSite := ParseSameSite(cfg.Auth.CookieSameSite)

	domain := cfg.Auth.CookieDomain
	if domain == "localhost" || domain == "127.0.0.1" {
		domain = ""
	}

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

// ClearSessionCookies clears session cookies according to AuthConfig.
func ClearSessionCookies(c echo.Context, cfg *config.Config) {
	sameSite := ParseSameSite(cfg.Auth.CookieSameSite)

	domain := cfg.Auth.CookieDomain
	if domain == "localhost" || domain == "127.0.0.1" {
		domain = ""
	}

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
