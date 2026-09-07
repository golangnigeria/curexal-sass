package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockValidator struct {
	active bool
	err    error
}

func (m *mockValidator) IsMembershipActive(ctx context.Context, userID, tenantID string) (bool, error) {
	return m.active, m.err
}

type mockAuditLogger struct {
	mu     sync.Mutex
	events []SecurityViolationEvent
}

func (m *mockAuditLogger) LogSecurityViolation(ctx context.Context, event SecurityViolationEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
	return nil
}

func (m *mockAuditLogger) getEvents() []SecurityViolationEvent {
	m.mu.Lock()
	defer m.mu.Unlock()
	copied := make([]SecurityViolationEvent, len(m.events))
	copy(copied, m.events)
	return copied
}

func TestRBACEnforcer_PlatformAdminBypass(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspace/patients", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	principal := &platformAuth.AuthenticatedPrincipal{
		UserID:   "admin-123",
		TenantID: "org-1",
		Platform: platformAuth.PlatformVector{
			IsSuperAdmin: true,
		},
	}
	c.Set(PrincipalKey, principal)

	enforcer := NewRBACEnforcer(nil, nil, nil)
	handler := enforcer.RequirePermission(platformAuth.PermWorkspacePatientRead)(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "org-1", c.Get("tenant_id"))
}

func TestRBACEnforcer_DirectPermissionGranted(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/clinical/encounters", nil)
	req.Header.Set("X-Branch-ID", "branch-south")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	principal := &platformAuth.AuthenticatedPrincipal{
		UserID:         "doc-1",
		TenantID:       "clinic-hq",
		ActiveBranchID: "branch-south",
		Permissions:    []string{platformAuth.PermWorkspaceClinicalWrite},
	}
	c.Set(PrincipalKey, principal)

	validator := &mockValidator{active: true}
	audit := &mockAuditLogger{}
	enforcer := NewRBACEnforcer(nil, validator, audit)

	handler := enforcer.RequirePermission(platformAuth.PermWorkspaceClinicalWrite)(func(c echo.Context) error {
		return c.String(http.StatusOK, "written")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "clinic-hq", c.Get("tenant_id"))
	assert.Equal(t, "branch-south", c.Get("branch_id"))
	assert.Empty(t, audit.getEvents())
}

func TestRBACEnforcer_MissingPermissionDenied(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/refunds", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	principal := &platformAuth.AuthenticatedPrincipal{
		UserID:      "nurse-1",
		TenantID:    "clinic-hq",
		Role:        "nurse",
		Permissions: []string{platformAuth.PermWorkspaceTriageCreate},
	}
	c.Set(PrincipalKey, principal)

	validator := &mockValidator{active: true}
	audit := &mockAuditLogger{}
	enforcer := NewRBACEnforcer(platformAuth.NewMemoryPermissionResolver(), validator, audit)

	handler := enforcer.RequirePermission(platformAuth.PermWorkspaceBillingRefund)(func(c echo.Context) error {
		return c.String(http.StatusOK, "refunded")
	})

	err := handler(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpErr.Code)
}

func TestRBACEnforcer_InactiveMembershipBlocked(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	principal := &platformAuth.AuthenticatedPrincipal{
		UserID:      "revoked-doctor",
		TenantID:    "clinic-hq",
		Role:        "doctor",
		Permissions: []string{platformAuth.PermWorkspacePatientRead},
	}
	c.Set(PrincipalKey, principal)

	// Membership is inactive/suspended
	validator := &mockValidator{active: false}
	audit := &mockAuditLogger{}
	enforcer := NewRBACEnforcer(nil, validator, audit)

	handler := enforcer.RequirePermission(platformAuth.PermWorkspacePatientRead)(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	err := handler(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpErr.Code)
	assert.Contains(t, httpErr.Message, "inactive or suspended")
}

func TestRBACEnforcer_CanonicalRoleResolver(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/prescriptions", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// User role is canonical 'doctor', with no pre-computed permissions slice in token
	principal := &platformAuth.AuthenticatedPrincipal{
		UserID:   "doc-smith",
		TenantID: "clinic-hq",
		Role:     "doctor",
	}
	c.Set(PrincipalKey, principal)

	validator := &mockValidator{active: true}
	resolver := platformAuth.NewMemoryPermissionResolver()
	enforcer := NewRBACEnforcer(resolver, validator, nil)

	handler := enforcer.RequirePermission(platformAuth.PermWorkspacePrescriptionWrite)(func(c echo.Context) error {
		return c.String(http.StatusOK, "prescribed")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
