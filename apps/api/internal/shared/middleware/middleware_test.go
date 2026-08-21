package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecretKey = "test-secret-key-for-unit-tests-only"

func newTestConfig(secret string) *config.Config {
	return &config.Config{
		Auth: config.AuthConfig{
			SecretKey:         secret,
			JWTCookieName:     "jwt",
			RefreshCookieName: "refresh_token",
			CookieDomain:      "localhost",
			CookiePath:        "/",
			CookieSecure:      false,
			CookieHTTPOnly:    true,
			CookieSameSite:    "Default",
			JWTExpiration:     15 * time.Minute,
			RefreshExpiration: 30 * 24 * time.Hour,
			AllowTestHeaders:  true,
		},
	}
}

// ─── helpers ──────────────────────────────────────────────────────────────────

// buildSignedJWT mints a HS256 JWT identical to what GenerateToken produces:
// claims: sub = userID, sid = sessionID, exp = now + 15 min.
func buildSignedJWT(t *testing.T, userID, sessionID, secret string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"sub": userID,
		"sid": sessionID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString([]byte(secret))
	require.NoError(t, err)
	return signed
}

// buildExpiredJWT mints a JWT that is already expired.
func buildExpiredJWT(t *testing.T, userID, secret string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(-5 * time.Minute).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString([]byte(secret))
	require.NoError(t, err)
	return signed
}

// newEchoContext returns a fresh Echo context backed by an httptest recorder.
func newEchoContext(t *testing.T, req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	t.Helper()
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

// ─── GetUserID ────────────────────────────────────────────────────────────────

func TestGetUserID_FromPrincipalOnContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, _ := newEchoContext(t, req)
	c.Set(PrincipalKey, &AuthenticatedPrincipal{UserID: "usr_abc"})

	assert.Equal(t, "usr_abc", GetUserID(c))
}

func TestGetUserID_FromContextKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, _ := newEchoContext(t, req)
	c.Set(UserIDKey, "usr_fromCtx")

	assert.Equal(t, "usr_fromCtx", GetUserID(c))
}

func TestGetUserID_FromXUserIDHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-ID", "usr_header")
	c, _ := newEchoContext(t, req)

	assert.Equal(t, "usr_header", GetUserID(c))
}

func TestGetUserID_NoSourceReturnsEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, _ := newEchoContext(t, req)

	assert.Equal(t, "", GetUserID(c))
}

// ─── GetActiveTenantID ────────────────────────────────────────────────────────

func TestGetActiveTenantID_FromPrincipal(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, _ := newEchoContext(t, req)
	c.Set(PrincipalKey, &AuthenticatedPrincipal{TenantID: "tenant-123"})

	assert.Equal(t, "tenant-123", GetActiveTenantID(c))
}

func TestGetActiveTenantID_FromContextKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, _ := newEchoContext(t, req)
	c.Set("tenant_id", "tenant-ctx")

	assert.Equal(t, "tenant-ctx", GetActiveTenantID(c))
}

func TestGetActiveTenantID_FromXTenantIDHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-ID", "tenant-hdr")
	c, _ := newEchoContext(t, req)

	assert.Equal(t, "tenant-hdr", GetActiveTenantID(c))
}

// ─── GetUserRole ──────────────────────────────────────────────────────────────

func TestGetUserRole_FromContextKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, _ := newEchoContext(t, req)
	c.Set(UserRoleKey, "branch_admin")

	assert.Equal(t, "branch_admin", GetUserRole(c))
}

func TestGetUserRole_FromXUserRoleHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-Role", "technician")
	c, _ := newEchoContext(t, req)

	assert.Equal(t, "technician", GetUserRole(c))
}

// ─── IsPlatformStaff ──────────────────────────────────────────────────────────

func TestIsPlatformStaff_SuperAdminReturnsTrueViaPrincipal(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, _ := newEchoContext(t, req)
	c.Set(PrincipalKey, &AuthenticatedPrincipal{
		Platform: PlatformVector{IsPlatformStaff: true, PlatformRole: "super_admin"},
	})

	assert.True(t, IsPlatformStaff(c))
}

func TestIsPlatformStaff_BranchAdminReturnsFalse(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, _ := newEchoContext(t, req)
	// A tenant branch admin role — must NOT pass the platform staff check.
	c.Set(PrincipalKey, &AuthenticatedPrincipal{
		Role: "admin", // generic tenant admin string
		Platform: PlatformVector{
			IsPlatformStaff: false, // correctly set during resolvePrincipal
			PlatformRole:    "admin",
		},
	})

	assert.False(t, IsPlatformStaff(c))
}

func TestIsPlatformStaff_PlatformRolesReturnTrue(t *testing.T) {
	for _, role := range []string{"super_admin", "platform_staff", "super_support_agent", "super_sales_staff"} {
		t.Run(role, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			c, _ := newEchoContext(t, req)
			c.Set(UserRoleKey, role)

			assert.True(t, IsPlatformStaff(c))
		})
	}
}

func TestIsPlatformStaff_NonPlatformRolesReturnFalse(t *testing.T) {
	for _, role := range []string{"owner", "branch_admin", "technician", "cashier", "customer_care", "clinician", "admin"} {
		t.Run(role, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			c, _ := newEchoContext(t, req)
			c.Set(UserRoleKey, role)

			assert.False(t, IsPlatformStaff(c))
		})
	}
}

// ─── RequireAuth middleware ────────────────────────────────────────────────────

func TestRequireAuth_WithAuthenticatedUser_CallsNext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(PrincipalKey, &AuthenticatedPrincipal{UserID: "usr_abc"})

	called := false
	handler := RequireAuth()(func(c echo.Context) error {
		called = true
		return c.NoContent(http.StatusOK)
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.True(t, called, "next handler must be called for authenticated user")
}

func TestRequireAuth_WithoutUser_Returns401(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No principal set — unauthenticated request.

	handler := RequireAuth()(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	err := handler(c)
	require.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, he.Code)
}

// ─── RequirePermission middleware ──────────────────────────────────────────────

func TestRequirePermission_PlatformStaffBypassesPermissionCheck(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, httptest.NewRecorder())
	c.Set(PrincipalKey, &AuthenticatedPrincipal{
		UserID:      "usr_staff",
		Permissions: []string{}, // no explicit permissions
		Platform:    PlatformVector{IsPlatformStaff: true},
	})

	called := false
	handler := RequirePermission("catalog:write")(func(c echo.Context) error {
		called = true
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.True(t, called, "platform staff must bypass permission check")
}

func TestRequirePermission_WildcardPermissionGrantsAccess(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, httptest.NewRecorder())
	c.Set(PrincipalKey, &AuthenticatedPrincipal{
		UserID:      "usr_super",
		Permissions: []string{"*"},
		Platform:    PlatformVector{IsPlatformStaff: false},
	})

	called := false
	handler := RequirePermission("any:specific:permission")(func(c echo.Context) error {
		called = true
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestRequirePermission_ExactPermissionMatch_GrantsAccess(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, httptest.NewRecorder())
	c.Set(PrincipalKey, &AuthenticatedPrincipal{
		UserID:      "usr_tech",
		Permissions: []string{"laboratory:enter_result", "laboratory:accession"},
		Platform:    PlatformVector{IsPlatformStaff: false},
	})

	called := false
	handler := RequirePermission("laboratory:enter_result")(func(c echo.Context) error {
		called = true
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestRequirePermission_MissingPermission_Returns403(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, httptest.NewRecorder())
	c.Set(PrincipalKey, &AuthenticatedPrincipal{
		UserID:      "usr_cashier",
		Permissions: []string{"billing:read"},
		Platform:    PlatformVector{IsPlatformStaff: false},
	})

	handler := RequirePermission("catalog:write")(func(c echo.Context) error {
		return nil
	})

	err := handler(c)
	require.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusForbidden, he.Code)
}

func TestRequirePermission_NilPrincipal_Returns401(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, httptest.NewRecorder())
	// No principal set at all.

	handler := RequirePermission("catalog:write")(func(c echo.Context) error {
		return nil
	})

	err := handler(c)
	require.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, he.Code)
}

// ─── RequireTenant middleware ─────────────────────────────────────────────────

func TestRequireTenant_WithTenantID_CallsNext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, httptest.NewRecorder())
	c.Set(PrincipalKey, &AuthenticatedPrincipal{TenantID: "tenant-xyz"})

	called := false
	handler := RequireTenant()(func(c echo.Context) error {
		called = true
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestRequireTenant_WithoutTenantID_Returns400(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, httptest.NewRecorder())
	// No tenant context set.

	handler := RequireTenant()(func(c echo.Context) error {
		return nil
	})

	err := handler(c)
	require.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, he.Code)
}

// ─── Authenticate middleware + JWT parsing ─────────────────────────────────────

func TestAuthenticate_ValidBearerToken_PopulatesPrincipal(t *testing.T) {
	userID := "usr_01test"
	sessionID := "sess_01test"
	token := buildSignedJWT(t, userID, sessionID, testSecretKey)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c := e.NewContext(req, httptest.NewRecorder())

	called := false
	handler := Authenticate(newTestConfig(testSecretKey))(func(c echo.Context) error {
		called = true
		p := GetPrincipal(c)
		require.NotNil(t, p)
		assert.Equal(t, userID, p.UserID)
		assert.Equal(t, sessionID, p.SessionID)
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestAuthenticate_ValidXAccessToken_PopulatesPrincipal(t *testing.T) {
	userID := "usr_xaccess"
	token := buildSignedJWT(t, userID, "sess_xaccess", testSecretKey)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Access-Token", token)
	c := e.NewContext(req, httptest.NewRecorder())

	called := false
	handler := Authenticate(newTestConfig(testSecretKey))(func(c echo.Context) error {
		called = true
		p := GetPrincipal(c)
		require.NotNil(t, p)
		assert.Equal(t, userID, p.UserID)
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestAuthenticate_XUserIDFallback_PopulatesPrincipal(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-ID", "usr_internal")
	req.Header.Set("X-Tenant-ID", "tenant-internal")
	req.Header.Set("X-User-Role", "super_admin")
	c := e.NewContext(req, httptest.NewRecorder())

	called := false
	handler := Authenticate(newTestConfig(testSecretKey))(func(c echo.Context) error {
		called = true
		p := GetPrincipal(c)
		require.NotNil(t, p)
		assert.Equal(t, "usr_internal", p.UserID)
		assert.Equal(t, "tenant-internal", p.TenantID)
		assert.True(t, p.Platform.IsPlatformStaff, "super_admin via X-User-ID must set IsPlatformStaff=true")
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestAuthenticate_ExpiredToken_PrincipalIsNil(t *testing.T) {
	expiredToken := buildExpiredJWT(t, "usr_expired", testSecretKey)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	c := e.NewContext(req, httptest.NewRecorder())

	called := false
	handler := Authenticate(newTestConfig(testSecretKey))(func(c echo.Context) error {
		called = true
		// Authenticate middleware must still call next — it does NOT block on expiry.
		// But the principal must be nil since the token is invalid.
		p := GetPrincipal(c)
		assert.Nil(t, p)
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestAuthenticate_WrongSecretKey_PrincipalIsNil(t *testing.T) {
	token := buildSignedJWT(t, "usr_abc", "sess_abc", "correct-secret")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c := e.NewContext(req, httptest.NewRecorder())

	handler := Authenticate(newTestConfig("wrong-secret"))(func(c echo.Context) error {
		p := GetPrincipal(c)
		assert.Nil(t, p, "token signed with different key must not produce a principal")
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
}

func TestAuthenticate_NoBearerPrefix_IsIgnored(t *testing.T) {
	token := buildSignedJWT(t, "usr_abc", "sess_abc", testSecretKey)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// Missing "Bearer " prefix — resolvePrincipal must NOT parse this as JWT.
	req.Header.Set("Authorization", token)
	c := e.NewContext(req, httptest.NewRecorder())

	handler := Authenticate(newTestConfig(testSecretKey))(func(c echo.Context) error {
		p := GetPrincipal(c)
		assert.Nil(t, p)
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
}

// ─── Token extraction precedence ─────────────────────────────────────────────

func TestAuthenticate_CookieTakesPrecedenceOverHeader(t *testing.T) {
	cookieUserID := "usr_from_cookie"
	headerUserID := "usr_from_header"

	cookieToken := buildSignedJWT(t, cookieUserID, "sess_cookie", testSecretKey)
	headerToken := buildSignedJWT(t, headerUserID, "sess_header", testSecretKey)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "jwt", Value: cookieToken})
	req.Header.Set("Authorization", "Bearer "+headerToken)
	c := e.NewContext(req, httptest.NewRecorder())

	handler := Authenticate(newTestConfig(testSecretKey))(func(c echo.Context) error {
		p := GetPrincipal(c)
		require.NotNil(t, p)
		// Cookie must win per the documented precedence order in resolvePrincipal.
		assert.Equal(t, cookieUserID, p.UserID, "cookie token must take precedence over Authorization header")
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
}

// ─── GetPrincipalFromContext — stdlib context.Context path ────────────────────

func TestGetPrincipalFromContext_WithStdlibContextValue_ReturnsPrincipal(t *testing.T) {
	p := &AuthenticatedPrincipal{UserID: "usr_stdlib"}
	// Inject into a stdlib context using the same PrincipalKey the middleware writes.
	ctx := context.WithValue(context.Background(), PrincipalKey, p)

	got := GetPrincipalFromContext(ctx)
	require.NotNil(t, got)
	assert.Equal(t, "usr_stdlib", got.UserID)
}

func TestGetPrincipalFromContext_EmptyContext_ReturnsNil(t *testing.T) {
	got := GetPrincipalFromContext(context.Background())
	assert.Nil(t, got)
}

// ─── GetSessionID / GetMembershipID ──────────────────────────────────────────

func TestGetSessionID_FromContextKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, _ := newEchoContext(t, req)
	c.Set("session_id", "sess_abc")
	assert.Equal(t, "sess_abc", GetSessionID(c))
}

func TestGetSessionID_FromXSessionIDHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Session-ID", "sess_hdr")
	c, _ := newEchoContext(t, req)
	assert.Equal(t, "sess_hdr", GetSessionID(c))
}

func TestGetMembershipID_FromContextKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, _ := newEchoContext(t, req)
	c.Set("membership_id", "mem_abc")
	assert.Equal(t, "mem_abc", GetMembershipID(c))
}

func TestGetMembershipID_FromXMembershipIDHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Membership-ID", "mem_hdr")
	c, _ := newEchoContext(t, req)
	assert.Equal(t, "mem_hdr", GetMembershipID(c))
}

// ─── GetSubdomainFromHeaders ──────────────────────────────────────────────────

func TestGetSubdomainFromHeaders_ValidSubdomain(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "acme.curexal.com"
	c, _ := newEchoContext(t, req)

	sub := GetSubdomainFromHeaders(c, "curexal.com")
	assert.Equal(t, "acme", sub)
}

func TestGetSubdomainFromHeaders_NoSubdomain(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "curexal.com"
	c, _ := newEchoContext(t, req)

	sub := GetSubdomainFromHeaders(c, "curexal.com")
	assert.Equal(t, "", sub)
}

func TestGetSubdomainFromHeaders_HostWithPort_StripsPort(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "clinic1.curexal.com:443"
	c, _ := newEchoContext(t, req)

	sub := GetSubdomainFromHeaders(c, "curexal.com")
	assert.Equal(t, "clinic1", sub)
}

func TestGetSubdomainFromHeaders_LocalhostNoBaseDomain_ReturnsEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "localhost:5002"
	c, _ := newEchoContext(t, req)

	sub := GetSubdomainFromHeaders(c, "")
	assert.Equal(t, "", sub)
}

// ─── GetPermissions ───────────────────────────────────────────────────────────

func TestGetPermissions_FromContextKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, _ := newEchoContext(t, req)
	c.Set("permissions", []string{"lab:read", "lab:write"})

	perms := GetPermissions(c)
	assert.Equal(t, []string{"lab:read", "lab:write"}, perms)
}

func TestGetPermissions_ReturnsNilWhenNotSet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, _ := newEchoContext(t, req)

	perms := GetPermissions(c)
	assert.Nil(t, perms)
}

// ─── Critical security invariant: branch 'admin' ≠ platform staff ─────────────

func TestAuthenticate_BranchAdminViaXUserID_IsPlatformStaffFalse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-ID", "usr_branch_admin")
	req.Header.Set("X-User-Role", "admin") // generic tenant "admin" — NOT a platform role
	c := e.NewContext(req, httptest.NewRecorder())

	handler := Authenticate(newTestConfig(testSecretKey))(func(c echo.Context) error {
		p := GetPrincipal(c)
		require.NotNil(t, p)
		// CRITICAL: branch "admin" must NOT elevate IsPlatformStaff.
		assert.False(t, p.Platform.IsPlatformStaff,
			"branch role 'admin' must NOT set IsPlatformStaff=true")
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
}

// ─── Default preferences are set on JWT path ──────────────────────────────────

func TestAuthenticate_TokenPath_SetsDefaultPreferences(t *testing.T) {
	token := buildSignedJWT(t, "usr_prefs", "sess_prefs", testSecretKey)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c := e.NewContext(req, httptest.NewRecorder())

	handler := Authenticate(newTestConfig(testSecretKey))(func(c echo.Context) error {
		p := GetPrincipal(c)
		require.NotNil(t, p)
		assert.Equal(t, "dark", p.Preferences.Theme)
		assert.Equal(t, "en", p.Preferences.Language)
		assert.Equal(t, "Africa/Lagos", p.Preferences.Timezone)
		assert.Equal(t, "YYYY-MM-DD", p.Preferences.DateFormat)
		assert.Equal(t, "en-NG", p.Preferences.NumberFormat)
		assert.Equal(t, "/dashboard", p.Preferences.DefaultLandingPage)
		return nil
	})

	err := handler(c)
	assert.NoError(t, err)
}

// suppress unused import warning — strings is used in the IsPlatformStaff role tests via HasPrefix.
var _ = strings.HasPrefix
