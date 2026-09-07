package auth

import (
	"strings"

	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/labstack/echo/v4"
)

// ResolvePrincipal resolves identity from incoming requests following the standard resolution pipeline:
// Cookie -> Authorization Bearer -> X-Access-Token -> Verified JWT Claims
func ResolvePrincipal(c echo.Context, cfg *config.Config) *AuthenticatedPrincipal {
	return ResolvePrincipalWithProvider(c, cfg, nil)
}

// ResolvePrincipalWithVerifier resolves principal identity and validates tenant membership using TenantMembershipVerifier.
func ResolvePrincipalWithVerifier(c echo.Context, cfg *config.Config, verifier TenantMembershipVerifier) *AuthenticatedPrincipal {
	return ResolvePrincipalWithAll(c, cfg, nil, verifier)
}

// ResolvePrincipalWithProvider resolves principal identity using an optional IdentityProvider implementation.
func ResolvePrincipalWithProvider(c echo.Context, cfg *config.Config, provider IdentityProvider) *AuthenticatedPrincipal {
	return ResolvePrincipalWithAll(c, cfg, provider, nil)
}

// ResolvePrincipalWithAll resolves principal identity using provider and verifier.
func ResolvePrincipalWithAll(c echo.Context, cfg *config.Config, provider IdentityProvider, verifier TenantMembershipVerifier) *AuthenticatedPrincipal {
	// Stage 1: Try IdentityProvider session verification
	if provider != nil {
		var sessionToken string
		if customSession := c.Request().Header.Get("X-Session-Token"); customSession != "" {
			sessionToken = customSession
		}

		if sessionToken != "" {
			sess, err := provider.Authenticate(c.Request().Context(), sessionToken)
			if err == nil && sess != nil && sess.Active {
				tenantID := resolveRequestTenantID(c)
				role := "member"

				// Server-side membership verification
				if verifier != nil && tenantID != "" {
					valid, memberRole, errVer := verifier.VerifyMembership(c.Request().Context(), sess.IdentityID, tenantID)
					if errVer == nil && valid {
						if memberRole != "" {
							role = memberRole
						}
					} else {
						// Not a verified member of the requested tenant
						tenantID = ""
					}
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

	// If token found, parse and cryptographically verify JWT claims
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

			tenantID := resolveRequestTenantID(c)

			// Security: Role is derived STRICTLY from verified cryptographic claims, NEVER from client headers
			role := "member"
			if platformRole != "" {
				role = platformRole
			} else if orgRole != "" {
				role = orgRole
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

			// Security: Super admin privilege requires verified claim or server platform role, NEVER client header
			isSuperAdmin := claims.IsPlatformAdmin || platformRole == "super_admin"

			// Server-side tenant membership verification (unless user is platform super admin)
			if verifier != nil && tenantID != "" && !isSuperAdmin {
				valid, memberRole, errVer := verifier.VerifyMembership(c.Request().Context(), claims.Subject, tenantID)
				if errVer == nil && valid {
					if memberRole != "" && orgRole == "" {
						role = memberRole
						orgRole = memberRole
					}
				} else {
					// User is not an active verified member of the requested tenant
					tenantID = ""
				}
			}

			effectiveOrgRole := orgRole
			if effectiveOrgRole == "" {
				effectiveOrgRole = role
			}

			if tenantID == "" && claims.OrganizationID != "" {
				tenantID = claims.OrganizationID
			}
			branchID := claims.ActiveBranchID

			p := &AuthenticatedPrincipal{
				UserID:         claims.Subject,
				SessionID:      claims.SessionID,
				TenantID:       tenantID,
				OrganizationID: tenantID,
				ActiveBranchID: branchID,
				Role:           role,
				Identity: IdentityVector{
					UserID: claims.Subject,
					Email:  claims.Email,
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

	// 3. Fallback to X-User-ID header (STRICTLY for internal test environments when explicitly enabled; disabled in production)
	if cfg.Auth.AllowTestHeaders && cfg.Primary.Env != "production" {
		if internalUserID := c.Request().Header.Get("X-User-ID"); internalUserID != "" {
			tenantID := resolveRequestTenantID(c)
			role := c.Request().Header.Get("X-User-Role")
			if role == "" {
				role = "member"
			}
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
	if resolvedOrgID := GetResolvedOrgID(c); resolvedOrgID != "" {
		return resolvedOrgID
	}
	if val := c.Get(ResolvedOrgIDKey); val != nil {
		if s, ok := val.(string); ok && s != "" {
			return s
		}
	}

	// 3. Query Parameter Overrides (e.g. ?organization_id=... or ?tenant_id=...)
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
