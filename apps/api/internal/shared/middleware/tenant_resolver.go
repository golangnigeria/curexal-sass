package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	TenantSchemaKey = "tenant_schema"
)

var validSlugRegex = regexp.MustCompile(`^[a-z0-9_-]+$`)

// HostAndSessionTenantResolver implements ADR 030: Host-Driven & Session-Based Tenant Resolution.
// It derives the tenant context strictly from the HTTP Host header or session principal.
// Spoofable client headers (like X-Tenant-Slug) are forbidden/ignored.
func HostAndSessionTenantResolver() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			host := req.Host

			// Strip port if present
			if colonIdx := strings.Index(host, ":"); colonIdx != -1 {
				host = host[:colonIdx]
			}

			tenantSlug := ""

			// 1. Resolve Tenant Slug from Subdomain (e.g. main-lab.curexal.com or main-lab.localhost)
			parts := strings.Split(host, ".")
			if len(parts) >= 3 || (len(parts) == 2 && parts[1] == "localhost") {
				candidate := strings.ToLower(parts[0])
				if candidate != "app" && candidate != "admin" && candidate != "hq" && candidate != "org font" && candidate != "api" && candidate != "www" {
					tenantSlug = candidate
				}
			}

			// 2. Fallback to Authenticated Session Principal Tenant ID / Slug
			if tenantSlug == "" {
				principal := GetPrincipal(c)
				if principal != nil && principal.TenantID != "" {
					tenantSlug = principal.TenantID
				}
			}

			// Default fallback for dev environment if unauthenticated
			if tenantSlug == "" {
				tenantSlug = "main-facility"
			}

			// Sanitize tenant slug to prevent SQL/schema tampering
			tenantSlug = strings.ToLower(strings.TrimSpace(tenantSlug))
			if !validSlugRegex.MatchString(tenantSlug) {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid tenant host or identifier")
			}

			tenantSchema := fmt.Sprintf("tenant_%s", strings.ReplaceAll(tenantSlug, "-", "_"))

			c.Set(TenantSlugKey, tenantSlug)
			c.Set(TenantSchemaKey, tenantSchema)

			return next(c)
		}
	}
}

// GetTenantSchema retrieves the resolved tenant schema name from Echo context.
func GetTenantSchema(c echo.Context) string {
	if val, ok := c.Get(TenantSchemaKey).(string); ok {
		return val
	}
	return "public"
}
