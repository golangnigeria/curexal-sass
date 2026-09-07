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

// HostAndSessionTenantResolver implements organization-centric branch workspace resolution.
// The Organization is resolved from the Host/domain by DomainResolverMiddleware.
// The operational Branch context is resolved from the X-Branch-Slug header, ?branch query param, or session tenant ID.
func HostAndSessionTenantResolver() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			branchSlug := req.Header.Get("X-Branch-Slug")
			if branchSlug == "" {
				branchSlug = c.QueryParam("branch")
			}

			// 2. Fallback to Authenticated Session Principal Tenant ID / Slug
			if branchSlug == "" {
				principal := GetPrincipal(c)
				if principal != nil {
					if principal.Workspace.ActiveWorkspaceID != "" {
						branchSlug = principal.Workspace.ActiveWorkspaceID
					} else if principal.TenantID != "" {
						branchSlug = principal.TenantID
					}
				}
			}

			// 3. Default fallback for initial/unauthenticated baseline
			if branchSlug == "" {
				branchSlug = "main-facility"
			}

			// Sanitize branch slug to prevent SQL/schema tampering
			branchSlug = strings.ToLower(strings.TrimSpace(branchSlug))
			if !validSlugRegex.MatchString(branchSlug) {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid branch identifier or slug")
			}

			tenantSchema := fmt.Sprintf("tenant_%s", strings.ReplaceAll(branchSlug, "-", "_"))

			c.Set(TenantSlugKey, branchSlug)
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
