package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	UserRoleKey       = platformAuth.UserRoleKey
	UserIDKey         = platformAuth.UserIDKey
	SessionIDKey      = platformAuth.SessionIDKey
	TenantSlugKey     = "tenant_slug"
	TenantNameKey     = "tenant_name"
	TenantIDKey       = platformAuth.TenantIDKey
	PatientContextKey = "patient_context"
	PrincipalKey      = platformAuth.PrincipalKey
)

type IdentityVector = platformAuth.IdentityVector
type PlatformVector = platformAuth.PlatformVector
type OrganizationVector = platformAuth.OrganizationVector
type WorkspaceVector = platformAuth.WorkspaceVector
type ActiveContextVector = platformAuth.ActiveContextVector
type UserPreferencesVector = platformAuth.UserPreferencesVector
type SecurityVector = platformAuth.SecurityVector
type AuthenticatedPrincipal = platformAuth.AuthenticatedPrincipal

type RequestContext struct {
	UserID          string
	TenantID        string
	Role            string
	IsPlatformAdmin bool
}

// Authenticate is a middleware that extracts JWT/session identity once per request and populates Context.
func Authenticate(cfg *config.Config) echo.MiddlewareFunc {
	return platformAuth.Authenticate(cfg)
}

// AuthenticateWithVerifier is a middleware that extracts identity and verifies tenant membership against the database.
func AuthenticateWithVerifier(cfg *config.Config, verifier platformAuth.TenantMembershipVerifier) echo.MiddlewareFunc {
	return platformAuth.AuthenticateWithVerifier(cfg, verifier)
}

func resolvePrincipal(c echo.Context, secretKey string) *AuthenticatedPrincipal {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			SecretKey:     secretKey,
			JWTCookieName: "jwt",
		},
	}
	return platformAuth.ResolvePrincipal(c, cfg)
}

func resolveRequestTenantID(c echo.Context) string {
	if c == nil {
		return ""
	}

	// 1. Direct Tenant / Organization Headers
	if tid := c.Request().Header.Get("X-Tenant-ID"); tid != "" {
		return tid
	}
	if orgID := c.Request().Header.Get("X-Organization-ID"); orgID != "" {
		return orgID
	}
	if activeTid := c.Request().Header.Get("X-Active-Tenant-ID"); activeTid != "" {
		return activeTid
	}

	// 2. Domain / Subdomain Resolved Organization ID
	if val := c.Get(platformAuth.ResolvedOrgIDKey); val != nil {
		if s, ok := val.(string); ok && s != "" {
			return s
		}
	}
	if resolvedOrgID := platformAuth.GetResolvedOrgID(c); resolvedOrgID != "" {
		return resolvedOrgID
	}

	// 3. Query Parameter Overrides
	if qOrgID := c.QueryParam("organization_id"); qOrgID != "" {
		return qOrgID
	}
	if qTenantID := c.QueryParam("tenant_id"); qTenantID != "" {
		return qTenantID
	}
	if qOrg := c.QueryParam("org_id"); qOrg != "" {
		return qOrg
	}

	// 4. Session / Active Org Cookie Fallbacks
	if cookie, err := c.Cookie("active_org_id"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	if cookie, err := c.Cookie("active_organization_id"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	if cookie, err := c.Cookie("tenant_id"); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

func GetPrincipal(c echo.Context) *AuthenticatedPrincipal {
	if val := c.Get(PrincipalKey); val != nil {
		if p, ok := val.(*AuthenticatedPrincipal); ok {
			return p
		}
	}
	return nil
}

func GetPrincipalFromContext(ctx context.Context) *AuthenticatedPrincipal {
	if c, ok := ctx.(echo.Context); ok {
		return GetPrincipal(c)
	}
	if p, ok := ctx.Value(PrincipalKey).(*AuthenticatedPrincipal); ok {
		return p
	}
	return nil
}

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

func RequirePermission(perm string, resolvers ...platformAuth.PermissionResolver) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			p := GetPrincipal(c)
			if p == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
			}
			if p.Platform.IsSuperAdmin || p.Platform.IsPlatformStaff || p.Platform.IsPlatformAdmin {
				return next(c)
			}

			if p.HasPermission(perm) {
				return next(c)
			}

			ctx := c.Request().Context()
			var resolver platformAuth.PermissionResolver
			if len(resolvers) > 0 && resolvers[0] != nil {
				resolver = resolvers[0]
			} else {
				resolver = platformAuth.NewMemoryPermissionResolver()
			}

			hasPerm, err := resolver.HasPermission(ctx, p, perm)
			if err == nil && hasPerm {
				return next(c)
			}

			c.Logger().Warn(fmt.Sprintf("RequirePermission DENIED: user_id='%s' role='%s' org_role='%s' required_perm='%s' uri='%s'", p.UserID, p.Role, p.Organization.OrganizationRole, perm, c.Request().RequestURI))

			return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to perform this action")
		}
	}
}

func RequirePlatformStaff() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			p := GetPrincipal(c)
			if p == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
			}
			if p.Platform.IsSuperAdmin || p.Platform.IsPlatformStaff || p.Platform.IsPlatformAdmin || IsPlatformStaff(c) {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden, "Platform staff privilege required")
		}
	}
}

type CapabilityChecker interface {
	HasCapability(ctx context.Context, orgID uuid.UUID, capabilityCode string) (bool, error)
}

func RequireCapability(capability string, checker CapabilityChecker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			p := GetPrincipal(c)
			if p == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
			}

			// Super Admin & Platform Admins bypass capability restrictions
			if p.Platform.IsSuperAdmin || p.Platform.IsPlatformStaff || p.Platform.IsPlatformAdmin {
				return next(c)
			}

			orgIDStr := p.Organization.ActiveOrganizationID
			if orgIDStr == "" {
				orgIDStr = c.Param("id")
			}
			if orgIDStr == "" {
				orgIDStr = c.Param("organizationId")
			}
			if orgIDStr == "" {
				orgIDStr = p.TenantID
			}

			if orgIDStr == "" {
				return echo.NewHTTPError(http.StatusForbidden, "Organization context required for capability verification")
			}

			orgUUID, errParse := uuid.Parse(orgIDStr)
			if errParse != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID format")
			}

			if checker != nil {
				hasCap, errCheck := checker.HasCapability(c.Request().Context(), orgUUID, capability)
				if errCheck == nil && hasCap {
					return next(c)
				}
			}

			c.Logger().Warn(fmt.Sprintf("RequireCapability DENIED: org_id='%s' user_id='%s' required_capability='%s' uri='%s'", orgIDStr, p.UserID, capability, c.Request().RequestURI))
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Organization does not have the required entitlement capability '%s'", capability))
		}
	}
}

func RequireTenant() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tenantID := GetActiveTenantID(c)
			if tenantID == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Active workspace context (tenant_id) is required")
			}
			return next(c)
		}
	}
}

func GetActiveTenantID(c echo.Context) string {
	if p := GetPrincipal(c); p != nil && p.TenantID != "" {
		return p.TenantID
	}
	if val := c.Get("tenant_id"); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return c.Request().Header.Get("X-Tenant-ID")
}

func GetOrganizationID(c echo.Context) string {
	if p := GetPrincipal(c); p != nil {
		if p.Organization.ActiveOrganizationID != "" {
			return p.Organization.ActiveOrganizationID
		}
		if p.OrganizationID != "" {
			return p.OrganizationID
		}
	}
	if val := c.Get("organization_id"); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return c.Request().Header.Get("X-Organization-ID")
}

func GetUserID(c echo.Context) string {
	if p := GetPrincipal(c); p != nil && p.UserID != "" {
		return p.UserID
	}
	if val := c.Get(UserIDKey); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return c.Request().Header.Get("X-User-ID")
}

func GetUserRole(c echo.Context) string {
	if val := c.Get(UserRoleKey); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return c.Request().Header.Get("X-User-Role")
}

func GetRequestID(c echo.Context) string {
	if val := c.Get("request_id"); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	if reqID := c.Response().Header().Get(echo.HeaderXRequestID); reqID != "" {
		return reqID
	}
	return c.Request().Header.Get("X-Request-ID")
}

func GetSessionID(c echo.Context) string {
	if val := c.Get("session_id"); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return c.Request().Header.Get("X-Session-ID")
}

func GetMembershipID(c echo.Context) string {
	if val := c.Get("membership_id"); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return c.Request().Header.Get("X-Membership-ID")
}

func GetPermissions(c echo.Context) []string {
	if val := c.Get("permissions"); val != nil {
		if perms, ok := val.([]string); ok {
			return perms
		}
	}
	return nil
}

func GetRequestContext(c echo.Context) *RequestContext {
	role := GetUserRole(c)
	return &RequestContext{
		UserID:          GetUserID(c),
		TenantID:        GetActiveTenantID(c),
		Role:            role,
		IsPlatformAdmin: IsPlatformStaff(c),
	}
}

func GetSubdomainFromHeaders(c echo.Context, baseDomain string) string {
	hosts := []string{
		c.Request().Header.Get("X-Forwarded-Host"),
		c.Request().Host,
	}

	if origin := c.Request().Header.Get("Origin"); origin != "" {
		if idx := strings.Index(origin, "://"); idx != -1 {
			origin = origin[idx+3:]
		}
		hosts = append(hosts, origin)
	}
	if referer := c.Request().Header.Get("Referer"); referer != "" {
		if idx := strings.Index(referer, "://"); idx != -1 {
			referer = referer[idx+3:]
		}
		if idx := strings.Index(referer, "/"); idx != -1 {
			referer = referer[:idx]
		}
		hosts = append(hosts, referer)
	}

	for _, host := range hosts {
		if host == "" {
			continue
		}
		if idx := strings.Index(host, ":"); idx != -1 {
			host = host[:idx]
		}

		if strings.HasSuffix(host, ".localhost") {
			parts := strings.Split(host, ".")
			if len(parts) >= 2 && parts[0] != "admin" && parts[0] != "www" && parts[0] != "app" {
				return parts[0]
			}
		}

		if baseDomain != "" && strings.HasSuffix(host, "."+baseDomain) {
			sub := strings.TrimSuffix(host, "."+baseDomain)
			if sub != "" && sub != "admin" && sub != "www" && sub != "app" {
				return sub
			}
		}
	}
	return ""
}

func IsPlatformStaff(c echo.Context) bool {
	if p := GetPrincipal(c); p != nil {
		if p.Platform.IsPlatformStaff {
			return true
		}
	}
	role := GetUserRole(c)
	return role == "super_admin" || role == "platform_staff" || role == "super_support_agent" || role == "super_sales_staff"
}

func GetPatientContext(c echo.Context) interface{} {
	return c.Get(PatientContextKey)
}

