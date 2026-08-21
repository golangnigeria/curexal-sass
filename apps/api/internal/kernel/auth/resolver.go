package auth

import (
	"strings"

	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/labstack/echo/v4"
)

// ResolvePrincipal resolves identity from incoming requests following the standard resolution pipeline:
// Cookie -> Authorization Bearer -> X-Access-Token -> X-User-ID
func ResolvePrincipal(c echo.Context, cfg *config.Config) *AuthenticatedPrincipal {
	return ResolvePrincipalWithProvider(c, cfg, nil)
}

// ResolvePrincipalWithProvider resolves principal identity using an optional IdentityProvider implementation.
func ResolvePrincipalWithProvider(c echo.Context, cfg *config.Config, provider IdentityProvider) *AuthenticatedPrincipal {
	// Stage 1: Try IdentityProvider session verification (MockProvider)
	if provider != nil {
		var sessionToken string
		if customSession := c.Request().Header.Get("X-Session-Token"); customSession != "" {
			sessionToken = customSession
		}

		if sessionToken != "" {
			sess, err := provider.Authenticate(c.Request().Context(), sessionToken)
			if err == nil && sess != nil && sess.Active {
				tenantID := c.Request().Header.Get("X-Tenant-ID")
				role := c.Request().Header.Get("X-User-Role")
				if role == "" {
					role = "member"
				}
				p := &AuthenticatedPrincipal{
					UserID:    sess.IdentityID,
					SessionID: sess.ID,
					TenantID:  tenantID,
					Role:      role,
					Identity: IdentityVector{
						UserID: sess.IdentityID,
					},
					Platform: PlatformVector{
						IsPlatformStaff: false,
						IsPlatformAdmin: false,
						IsSuperAdmin:    false,
						PlatformRole:    "",
					},
					Organization: OrganizationVector{
						ActiveOrganizationID: tenantID,
						OrganizationRole:     role,
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
						SessionID: sess.ID,
					},
				}
				resolver := NewProviderPermissionResolver()
				perms, _ := resolver.ResolvePermissions(c.Request().Context(), p)
				p.Permissions = perms
				return p
			}
		}
	}

	var tokenStr string

	// 2. Try extracting token from JWT Cookie
	if cookie, err := c.Cookie(cfg.Auth.JWTCookieName); err == nil && cookie.Value != "" {
		tokenStr = cookie.Value
	}

	// 3. Try extracting token from Authorization header (Bearer token) or X-Access-Token header
	if tokenStr == "" {
		authHeader := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else if customToken := c.Request().Header.Get("X-Access-Token"); customToken != "" {
			tokenStr = customToken
		}
	}

	// If token found, parse JWT claims
	if tokenStr != "" {
		claims, err := ParseAccessJWT(cfg, tokenStr)
		if err == nil && claims != nil && claims.Subject != "" {
			var platformRole string
			if claims.PlatformRole != nil {
				platformRole = *claims.PlatformRole
			}
			var orgRole string
			if claims.OrganizationRole != nil {
				orgRole = *claims.OrganizationRole
			}
			tenantID := c.Request().Header.Get("X-Tenant-ID")
			role := c.Request().Header.Get("X-User-Role")
			if role == "" {
				if platformRole != "" {
					role = platformRole
				} else if orgRole != "" {
					role = orgRole
				}
			}

			isStaff := claims.IsPlatformAdmin
			if !isStaff && platformRole != "" {
				for _, staffRole := range cfg.Auth.PlatformStaffRoles {
					if platformRole == staffRole {
						isStaff = true
						break
					}
				}
			}

			isSuperAdmin := claims.IsPlatformAdmin || platformRole == "super_admin" || role == "super_admin"
			effectiveOrgRole := orgRole
			if effectiveOrgRole == "" {
				effectiveOrgRole = role
			}
			p := &AuthenticatedPrincipal{
				UserID:    claims.Subject,
				SessionID: claims.SessionID,
				TenantID:  tenantID,
				Role:      role,
				Identity: IdentityVector{
					UserID: claims.Subject,
				},
				Platform: PlatformVector{
					IsPlatformStaff: isStaff,
					IsPlatformAdmin: claims.IsPlatformAdmin,
					IsSuperAdmin:    isSuperAdmin,
					PlatformRole:    platformRole,
				},
				Organization: OrganizationVector{
					ActiveOrganizationID: tenantID,
					OrganizationRole:     effectiveOrgRole,
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
					SessionID: claims.SessionID,
				},
			}
			resolver := NewProviderPermissionResolver()
			perms, _ := resolver.ResolvePermissions(c.Request().Context(), p)
			p.Permissions = perms
			return p
		}
	}

	// 3. Fallback to X-User-ID header (for internal service calls/testing if enabled)
	if cfg.Auth.AllowTestHeaders {
		if internalUserID := c.Request().Header.Get("X-User-ID"); internalUserID != "" {
			tenantID := c.Request().Header.Get("X-Tenant-ID")
			role := c.Request().Header.Get("X-User-Role")
			isStaff := role == "super_admin" || role == "platform_staff" || role == "super_support_agent" || role == "super_sales_staff"
			isSuperAdmin := role == "super_admin"
			p := &AuthenticatedPrincipal{
				UserID:   internalUserID,
				TenantID: tenantID,
				Role:     role,
				Identity: IdentityVector{
					UserID: internalUserID,
				},
				Platform: PlatformVector{
					IsPlatformStaff: isStaff,
					IsPlatformAdmin: isSuperAdmin,
					IsSuperAdmin:    isSuperAdmin,
					PlatformRole:    role,
				},
				Organization: OrganizationVector{
					ActiveOrganizationID: tenantID,
					OrganizationRole:     role,
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
					SessionID: "test-session",
				},
			}
			resolver := NewProviderPermissionResolver()
			perms, _ := resolver.ResolvePermissions(c.Request().Context(), p)
			p.Permissions = perms
			return p
		}
	}

	return nil
}
