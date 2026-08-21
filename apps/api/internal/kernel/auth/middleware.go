package auth

import (
	"context"
	"net/http"

	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/labstack/echo/v4"
)

const (
	PrincipalKey = "principal"
	UserIDKey    = "user_id"
	SessionIDKey = "session_id"
	TenantIDKey  = "tenant_id"
	UserRoleKey  = "user_role"
)

// Authenticate is the platform Echo middleware that extracts identity using ResolvePrincipal and sets context.
func Authenticate(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			principal := ResolvePrincipal(c, cfg)
			if principal != nil {
				c.Set(PrincipalKey, principal)
				c.Set(UserIDKey, principal.UserID)
				if principal.SessionID != "" {
					c.Set(SessionIDKey, principal.SessionID)
				}
				if principal.TenantID != "" {
					c.Set(TenantIDKey, principal.TenantID)
				}
				effectiveRole := principal.Role
				if effectiveRole == "" {
					effectiveRole = principal.Organization.OrganizationRole
				}
				if effectiveRole != "" {
					c.Set(UserRoleKey, effectiveRole)
				}
				if principal.Organization.ActiveOrganizationID != "" {
					c.Set("organization_id", principal.Organization.ActiveOrganizationID)
				}
			}
			return next(c)
		}
	}
}

// RequireAuth middleware ensures an authenticated principal is present on protected routes.
func RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if GetUserID(c) == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
			}
			return next(c)
		}
	}
}

// GetPrincipal retrieves the AuthenticatedPrincipal from the Echo context.
func GetPrincipal(c echo.Context) *AuthenticatedPrincipal {
	if val := c.Get(PrincipalKey); val != nil {
		if p, ok := val.(*AuthenticatedPrincipal); ok {
			return p
		}
	}
	return nil
}

// GetUserID retrieves the authenticated user ID from context.
func GetUserID(c echo.Context) string {
	if val, ok := c.Get(UserIDKey).(string); ok {
		return val
	}
	if p := GetPrincipal(c); p != nil {
		return p.UserID
	}
	return ""
}

// GetSessionID retrieves the authenticated session ID from context.
func GetSessionID(c echo.Context) string {
	if val, ok := c.Get(SessionIDKey).(string); ok {
		return val
	}
	if p := GetPrincipal(c); p != nil {
		return p.SessionID
	}
	return ""
}

// GetPrincipalFromContext retrieves AuthenticatedPrincipal from standard Go context.
func GetPrincipalFromContext(ctx context.Context) *AuthenticatedPrincipal {
	if c, ok := ctx.(echo.Context); ok {
		return GetPrincipal(c)
	}
	if p, ok := ctx.Value(PrincipalKey).(*AuthenticatedPrincipal); ok {
		return p
	}
	return nil
}
