package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockBranchValidator struct {
	access bool
	err    error
}

func (m *mockBranchValidator) VerifyUserFacilityAccess(ctx context.Context, orgID, branchID uuid.UUID, userID string, isOrgAdmin bool) (bool, error) {
	return m.access, m.err
}

func TestRequireActiveBranch_MissingBranchContext_Returns428(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspace/clinical/encounters", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	orgID := uuid.New().String()
	principal := &platformAuth.AuthenticatedPrincipal{
		UserID:   uuid.New().String(),
		TenantID: orgID,
		Role:     "doctor",
	}
	c.Set(PrincipalKey, principal)

	handler := RequireActiveBranch(&mockBranchValidator{access: true}, nil)(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	err := handler(c)
	require.Error(t, err)

	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusPreconditionRequired, httpErr.Code)
}

func TestRequireActiveBranch_InvalidBranchFormat_Returns400(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspace/clinical/encounters", nil)
	req.Header.Set("X-Branch-ID", "not-a-valid-uuid")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	principal := &platformAuth.AuthenticatedPrincipal{
		UserID:   uuid.New().String(),
		TenantID: uuid.New().String(),
		Role:     "doctor",
	}
	c.Set(PrincipalKey, principal)

	handler := RequireActiveBranch(&mockBranchValidator{access: true}, nil)(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	err := handler(c)
	require.Error(t, err)

	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestRequireActiveBranch_PlatformAdminBypass(t *testing.T) {
	e := echo.New()
	branchID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspace/patients", nil)
	req.Header.Set("X-Branch-ID", branchID)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	principal := &platformAuth.AuthenticatedPrincipal{
		UserID:   uuid.New().String(),
		TenantID: uuid.New().String(),
		Platform: platformAuth.PlatformVector{
			IsSuperAdmin: true,
		},
	}
	c.Set(PrincipalKey, principal)

	// Even if validator returns false, platform admin must bypass
	handler := RequireActiveBranch(&mockBranchValidator{access: false}, nil)(func(c echo.Context) error {
		assert.Equal(t, branchID, c.Get("branch_id"))
		return c.String(http.StatusOK, "success")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireActiveBranch_OrgAdmin_AccessGranted(t *testing.T) {
	e := echo.New()
	branchID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/workspace/patients", nil)
	req.Header.Set("X-Branch-ID", branchID)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	principal := &platformAuth.AuthenticatedPrincipal{
		UserID:   uuid.New().String(),
		TenantID: uuid.New().String(),
		Role:     "owner",
	}
	c.Set(PrincipalKey, principal)

	handler := RequireActiveBranch(&mockBranchValidator{access: true}, nil)(func(c echo.Context) error {
		assert.Equal(t, branchID, c.Get("branch_id"))
		return c.String(http.StatusOK, "success")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireActiveBranch_StaffAssigned_AccessGranted(t *testing.T) {
	e := echo.New()
	branchID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspace/clinical/encounters", nil)
	req.Header.Set("X-Branch-ID", branchID)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	principal := &platformAuth.AuthenticatedPrincipal{
		UserID:   uuid.New().String(),
		TenantID: uuid.New().String(),
		Role:     "doctor",
	}
	c.Set(PrincipalKey, principal)

	handler := RequireActiveBranch(&mockBranchValidator{access: true}, nil)(func(c echo.Context) error {
		assert.Equal(t, branchID, c.Get("branch_id"))
		return c.String(http.StatusOK, "success")
	})

	err := handler(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireActiveBranch_StaffUnassigned_Returns403_AndLogsViolation(t *testing.T) {
	e := echo.New()
	branchID := uuid.New().String()
	orgID := uuid.New().String()
	userID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/workspace/clinical/encounters", nil)
	req.Header.Set("X-Branch-ID", branchID)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	principal := &platformAuth.AuthenticatedPrincipal{
		UserID:   userID,
		TenantID: orgID,
		Role:     "doctor",
	}
	c.Set(PrincipalKey, principal)

	auditLogger := &mockAuditLogger{}
	handler := RequireActiveBranch(&mockBranchValidator{access: false}, auditLogger)(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	err := handler(c)
	require.Error(t, err)

	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpErr.Code)

	// Wait briefly for asynchronous goroutine audit log
	time.Sleep(50 * time.Millisecond)
	events := auditLogger.getEvents()
	require.Len(t, events, 1)
	assert.Equal(t, userID, events[0].UserID)
	assert.Equal(t, orgID, events[0].TenantID)
	assert.Equal(t, branchID, events[0].BranchID)
	assert.Equal(t, "workspace:branch:assignment", events[0].RequiredPerm)
}
