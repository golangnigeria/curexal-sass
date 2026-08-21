package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/labstack/echo/v4"
)

func TestAuthorizationRegressionMatrix(t *testing.T) {
	resolver := NewMemoryPermissionResolver()

	tests := []struct {
		name           string
		role           string
		isPlatformAdmin bool
		permission     string
		expectedHas    bool
	}{
		// 1. Platform Super Admin (Bypass)
		{"Platform Admin - organization:read", "super_admin", true, PermissionOrganizationRead, true},
		{"Platform Admin - users:write", "super_admin", true, PermissionUsersWrite, true},
		{"Platform Admin - arbitrary permission", "super_admin", true, "custom:anything", true},

		// 2. Organization Owner
		{"Org Owner - organization:read", "owner", false, PermissionOrganizationRead, true},
		{"Org Owner - organization:settings:write", "owner", false, PermissionOrganizationSettingsWrite, true},
		{"Org Owner - organization:create", "owner", false, PermissionOrganizationCreate, true},
		{"Org Owner - users:read", "owner", false, PermissionUsersRead, true},
		{"Org Owner - users:write", "owner", false, PermissionUsersWrite, true},
		{"Org Owner - audit:read", "owner", false, PermissionAuditRead, true},

		// 3. Organization Admin
		{"Org Admin - organization:read", "org_admin", false, PermissionOrganizationRead, true},
		{"Org Admin - users:read", "org_admin", false, PermissionUsersRead, true},
		{"Org Admin - organization:create (denied)", "org_admin", false, PermissionOrganizationCreate, false},

		// 4. Branch Admin
		{"Branch Admin - organization:read", "branch_admin", false, PermissionOrganizationRead, true},
		{"Branch Admin - patient:view", "branch_admin", false, PermissionPatientView, true},
		{"Branch Admin - audit:read (denied)", "branch_admin", false, PermissionAuditRead, false},

		// 5. Staff Members (Clinician / Technician / Member)
		{"Clinician - patient:view", "clinician", false, PermissionPatientView, true},
		{"Clinician - organization:read (denied)", "clinician", false, PermissionOrganizationRead, false},
		{"Technician - laboratory:accession", "technician", false, PermissionLabAccession, true},
		{"Technician - organization:read (denied)", "technician", false, PermissionOrganizationRead, false},
		{"Member - identity:password:write", "member", false, PermissionPasswordWrite, true},
		{"Member - organization:read (denied)", "member", false, PermissionOrganizationRead, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			principal := &AuthenticatedPrincipal{
				UserID: "usr_test_123",
				Role:   tt.role,
				Platform: PlatformVector{
					IsPlatformAdmin: tt.isPlatformAdmin,
					IsPlatformStaff: tt.isPlatformAdmin,
					PlatformRole:    tt.role,
				},
				Organization: OrganizationVector{
					OrganizationRole: tt.role,
				},
			}

			perms, err := resolver.ResolvePermissions(context.Background(), principal)
			if err != nil {
				t.Fatalf("unexpected error resolving permissions: %v", err)
			}
			principal.Permissions = perms

			hasPerm, err := resolver.HasPermission(context.Background(), principal, tt.permission)
			if err != nil {
				t.Fatalf("unexpected error in HasPermission: %v", err)
			}

			if hasPerm != tt.expectedHas {
				t.Errorf("expected HasPermission(%s) = %v, got %v for role %s", tt.permission, tt.expectedHas, hasPerm, tt.role)
			}

			if principal.HasPermission(tt.permission) != tt.expectedHas {
				t.Errorf("expected principal.HasPermission(%s) = %v, got %v for role %s", tt.permission, tt.expectedHas, principal.HasPermission(tt.permission), tt.role)
			}
		})
	}
}

func TestRequirePermissionMiddlewareMatrix(t *testing.T) {
	resolver := NewMemoryPermissionResolver()

	requirePermMiddleware := func(targetPerm string) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				p := GetPrincipal(c)
				if p == nil {
					return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
				}
				if p.HasPermission(targetPerm) {
					return next(c)
				}
				has, _ := resolver.HasPermission(c.Request().Context(), p, targetPerm)
				if has {
					return next(c)
				}
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to perform this action")
			}
		}
	}

	tests := []struct {
		name         string
		principal    *AuthenticatedPrincipal
		permission   string
		expectedCode int
	}{
		{
			name:         "Anonymous user returns 401 Unauthorized",
			principal:    nil,
			permission:   PermissionOrganizationRead,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Organization Owner accessing GET /organizations returns 200 OK",
			principal: &AuthenticatedPrincipal{
				UserID: "usr_owner_1",
				Role:   "owner",
			},
			permission:   PermissionOrganizationRead,
			expectedCode: http.StatusOK,
		},
		{
			name: "Clinician accessing GET /organizations returns 403 Forbidden",
			principal: &AuthenticatedPrincipal{
				UserID: "usr_doc_1",
				Role:   "clinician",
			},
			permission:   PermissionOrganizationRead,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "Platform Super Admin accessing GET /organizations returns 200 OK via bypass",
			principal: &AuthenticatedPrincipal{
				UserID: "usr_admin_1",
				Platform: PlatformVector{
					IsPlatformAdmin: true,
				},
			},
			permission:   PermissionOrganizationRead,
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/organizations", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.principal != nil {
				perms, _ := resolver.ResolvePermissions(context.Background(), tt.principal)
				tt.principal.Permissions = perms
				c.Set(PrincipalKey, tt.principal)
			}

			handler := requirePermMiddleware(tt.permission)(func(c echo.Context) error {
				return c.String(http.StatusOK, "OK")
			})

			err := handler(c)
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					if he.Code != tt.expectedCode {
						t.Errorf("expected HTTP status %d, got %d", tt.expectedCode, he.Code)
					}
					return
				}
				t.Fatalf("unexpected error: %v", err)
			}

			if rec.Code != tt.expectedCode {
				t.Errorf("expected HTTP status %d, got %d", tt.expectedCode, rec.Code)
			}
		})
	}
}

func TestEmpiricalJWTAndPrincipalResolution(t *testing.T) {
	cfg := &config.Config{}
	cfg.Auth.SecretKey = "curexal_empirical_verification_secret_key_123456"
	cfg.Auth.JWTCookieName = "curexal_jwt"
	cfg.Auth.JWTExpiration = 15 * time.Minute

	orgRole := "owner"
	platformRole := "member"
	tokenStr, err := GenerateAccessJWT(cfg, "usr_owner_999", "sess_empirical_123", &platformRole, false, &orgRole)
	if err != nil {
		t.Fatalf("failed to generate access JWT: %v", err)
	}

	claims, err := ParseAccessJWT(cfg, tokenStr)
	if err != nil {
		t.Fatalf("failed to parse access JWT: %v", err)
	}

	if claims.Subject != "usr_owner_999" {
		t.Errorf("expected Subject 'usr_owner_999', got '%s'", claims.Subject)
	}
	if claims.OrganizationRole == nil || *claims.OrganizationRole != "owner" {
		t.Errorf("expected OrganizationRole 'owner', got '%v'", claims.OrganizationRole)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	req.Header.Set("X-Tenant-ID", "tenant_branch_alpha")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	principal := ResolvePrincipal(c, cfg)
	if principal == nil {
		t.Fatalf("ResolvePrincipal returned nil for valid JWT token")
	}

	if principal.UserID != "usr_owner_999" {
		t.Errorf("expected principal.UserID 'usr_owner_999', got '%s'", principal.UserID)
	}
	if principal.Organization.OrganizationRole != "owner" {
		t.Errorf("expected OrganizationRole 'owner', got '%s'", principal.Organization.OrganizationRole)
	}

	// Verify permissions
	hasRead := false
	for _, p := range principal.Permissions {
		if p == PermissionUsersRead {
			hasRead = true
			break
		}
	}
	if !hasRead {
		t.Errorf("expected Organization Owner to have 'users:read' permission, got permissions: %v", principal.Permissions)
	}
}

type mockRoleFinder struct {
	data map[string][]string
}

func (m *mockRoleFinder) ListPermissionsByRole(ctx context.Context, roleName, tenantID string) ([]string, error) {
	if perms, ok := m.data[roleName]; ok {
		return perms, nil
	}
	return []string{}, nil
}

func TestEmpiricalDatabaseRolePermissionsResolution(t *testing.T) {
	finder := &mockRoleFinder{
		data: map[string][]string{
			"owner": {"organization:read", "organization:settings:write", "users:read", "users:write", "audit:read"},
		},
	}

	provider := NewDatabaseRolePermissionProvider(finder)
	resolver := NewProviderPermissionResolver(provider)

	principal := &AuthenticatedPrincipal{
		UserID: "usr_db_test_100",
		Role:   "owner",
		Organization: OrganizationVector{
			OrganizationRole: "owner",
		},
	}

	perms, err := resolver.ResolvePermissions(context.Background(), principal)
	if err != nil {
		t.Fatalf("failed to resolve database permissions: %v", err)
	}

	expected := map[string]bool{
		"organization:read":           true,
		"organization:settings:write": true,
		"users:read":                  true,
		"users:write":                 true,
		"audit:read":                  true,
	}

	if len(perms) != len(expected) {
		t.Errorf("expected %d permissions from database provider, got %d (%v)", len(expected), len(perms), perms)
	}

	for _, p := range perms {
		if !expected[p] {
			t.Errorf("unexpected permission resolved from database: %s", p)
		}
	}
}
