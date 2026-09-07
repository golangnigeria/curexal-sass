package authz

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RequirePermission creates an Echo middleware checking if the authenticated principal has a required permission.
func RequirePermission(authzService *Service, permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract principal from context
			val := c.Get("principal")
			if val == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication context")
			}

			principal, ok := val.(*Principal)
			if !ok || principal == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid principal format")
			}

			if !authzService.Can(c.Request().Context(), principal, permission, principal.Context) {
				return echo.NewHTTPError(http.StatusForbidden, "insufficient privileges for this action")
			}

			return next(c)
		}
	}
}
