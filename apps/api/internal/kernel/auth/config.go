package auth

import (
	"net/http"
	"strings"
)

// ParseSameSite converts a string configuration ("Lax", "Strict", "None", "Default") into an http.SameSite enum.
func ParseSameSite(mode string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	case "default":
		return http.SameSiteDefaultMode
	default:
		return http.SameSiteDefaultMode
	}
}
