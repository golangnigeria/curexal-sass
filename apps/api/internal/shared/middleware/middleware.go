package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

func resolvePrincipal(c echo.Context, secretKey string) *AuthenticatedPrincipal {
	var tokenStr string

	// 1. Try extracting token from JWT Cookie
	if cookie, err := c.Cookie("jwt"); err == nil && cookie.Value != "" {
		tokenStr = cookie.Value
	}

	// 2. Try extracting token from Authorization header or X-Access-Token header
	if tokenStr == "" {
		authHeader := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else if customToken := c.Request().Header.Get("X-Access-Token"); customToken != "" {
			tokenStr = customToken
		}
	}

	if tokenStr != "" && secretKey != "" {
		token, errParse := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if errParse == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				sub, _ := claims.GetSubject()
				sid, _ := claims["sid"].(string)

				var platformRole string
				if pr, ok := claims["platform_role"].(string); ok && pr != "" {
					platformRole = pr
				}
				isPlatformAdmin, _ := claims["is_platform_admin"].(bool)

				tenantID := c.Request().Header.Get("X-Tenant-ID")
				role := c.Request().Header.Get("X-User-Role")
				if role == "" {
					role = platformRole
				}

				if sub != "" {
					isStaff := isPlatformAdmin || platformRole == "super_admin" || platformRole == "platform_staff" || platformRole == "super_support_agent" || platformRole == "super_sales_staff"
					return &AuthenticatedPrincipal{
						UserID:    sub,
						SessionID: sid,
						TenantID:  tenantID,
						Role:      role,
						Identity: IdentityVector{
							UserID: sub,
						},
						Platform: PlatformVector{
							IsPlatformStaff: isStaff,
							PlatformRole:    platformRole,
						},
						Organization: OrganizationVector{
							ActiveOrganizationID: tenantID,
						},
						Workspace: WorkspaceVector{
							ActiveWorkspaceID: tenantID,
							WorkspaceRole:     role,
						},
						ActiveContext: ActiveContextVector{
							Type:      "platform",
							ContextID: tenantID,
						},
						Preferences: UserPreferencesVector{
							Theme:              "dark",
							Language:           "en",
							Timezone:           "Africa/Lagos",
							DateFormat:         "YYYY-MM-DD",
							NumberFormat:       "en-NG",
							DefaultLandingPage: "/dashboard",
						},
						Security: SecurityVector{
							SessionID: sid,
						},
					}
				}
			}
		}
	}

	// 3. Fallback to X-User-ID header (for internal service calls/tests)
	if internalUserID := c.Request().Header.Get("X-User-ID"); internalUserID != "" {
		tenantID := c.Request().Header.Get("X-Tenant-ID")
		role := c.Request().Header.Get("X-User-Role")
		return &AuthenticatedPrincipal{
			UserID:   internalUserID,
			TenantID: tenantID,
			Role:     role,
			Identity: IdentityVector{
				UserID: internalUserID,
			},
			Platform: PlatformVector{
				IsPlatformStaff: role == "super_admin" || role == "platform_staff" || role == "super_support_agent" || role == "super_sales_staff",
				PlatformRole:    role,
			},
			Organization: OrganizationVector{
				ActiveOrganizationID: tenantID,
			},
			Workspace: WorkspaceVector{
				ActiveWorkspaceID: tenantID,
				WorkspaceRole:     role,
			},
			ActiveContext: ActiveContextVector{
				Type:      "platform",
				ContextID: tenantID,
			},
			Preferences: UserPreferencesVector{
				Theme:              "dark",
				Language:           "en",
				Timezone:           "Africa/Lagos",
				DateFormat:         "YYYY-MM-DD",
				NumberFormat:       "en-NG",
				DefaultLandingPage: "/dashboard",
			},
		}
	}

	return nil
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

