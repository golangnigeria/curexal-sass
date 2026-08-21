package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/platform/api"
	"github.com/golangnigeria/curexal/internal/modules/platform/application"
	"github.com/golangnigeria/curexal/internal/modules/platform/domain"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConfigRepo struct {
	mock.Mock
}

func (m *MockConfigRepo) GetGeneralSettings(ctx context.Context) (*domain.PlatformGeneralSettings, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PlatformGeneralSettings), args.Error(1)
}

func (m *MockConfigRepo) UpdateGeneralSettings(ctx context.Context, settings *domain.PlatformGeneralSettings, updatedBy uuid.UUID) (*domain.PlatformGeneralSettings, error) {
	args := m.Called(ctx, settings, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PlatformGeneralSettings), args.Error(1)
}

func (m *MockConfigRepo) GetSecurityPolicy(ctx context.Context) (*domain.IdentitySecurityPolicy, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IdentitySecurityPolicy), args.Error(1)
}

func (m *MockConfigRepo) UpdateSecurityPolicy(ctx context.Context, policy *domain.IdentitySecurityPolicy, updatedBy uuid.UUID) (*domain.IdentitySecurityPolicy, error) {
	args := m.Called(ctx, policy, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IdentitySecurityPolicy), args.Error(1)
}

type MockAuditRepo struct {
	mock.Mock
}

func (m *MockAuditRepo) Create(ctx context.Context, payload *auditDomain.CreateAuditLogPayload) (*auditDomain.AuditLog, error) {
	args := m.Called(ctx, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditDomain.AuditLog), args.Error(1)
}

func (m *MockAuditRepo) GetByID(ctx context.Context, id uuid.UUID) (*auditDomain.AuditLog, error) {
	return nil, nil
}

func (m *MockAuditRepo) ListTenantLogs(ctx context.Context, tenantID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]auditDomain.AuditLog, error) {
	return nil, nil
}

func (m *MockAuditRepo) ListPlatformLogs(ctx context.Context, orgID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]auditDomain.AuditLog, error) {
	return nil, nil
}

func (m *MockAuditRepo) GetStats(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID) (*auditDomain.AdminStats, error) {
	return nil, nil
}

func (m *MockAuditRepo) ListAll(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID, limit int, offset int) ([]auditDomain.AuditLog, error) {
	return nil, nil
}

func TestPlatformConfigHandler_GetPlatformConfig(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/platform/config", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := new(MockConfigRepo)
	mockAudit := new(MockAuditRepo)

	mockRepo.On("GetGeneralSettings", mock.Anything).Return(&domain.PlatformGeneralSettings{
		PlatformName:    "Curexal",
		DefaultCountry:  "NG",
		DefaultCurrency: "NGN",
	}, nil)

	svc := application.NewPlatformConfigService(mockRepo, mockAudit)
	handler := api.NewPlatformConfigHandler(svc)

	err := handler.GetPlatformConfig(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Curexal")
}

func TestPlatformConfigHandler_UpdatePlatformConfig_ForbiddenForNonAdmin(t *testing.T) {
	e := echo.New()
	jsonBody := `{"platformName":"Curexal Health","supportEmail":"admin@curexal.com","defaultCountry":"NG","defaultCurrency":"NGN"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/platform/config", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	nonAdminPrincipal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Role:   "organization_user",
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: false,
		},
	}
	c.Set(middleware.PrincipalKey, nonAdminPrincipal)

	mockRepo := new(MockConfigRepo)
	mockAudit := new(MockAuditRepo)
	svc := application.NewPlatformConfigService(mockRepo, mockAudit)
	handler := api.NewPlatformConfigHandler(svc)

	err := handler.UpdatePlatformConfig(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestPlatformConfigHandler_UpdatePlatformConfig_PartialSuccess(t *testing.T) {
	e := echo.New()
	jsonBody := `{"platformName":"Curexal Health OS","supportEmail":"support@curexal.com"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/platform/config", strings.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminPrincipal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Role:   "platform_admin",
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: true,
		},
		Identity: platformAuth.IdentityVector{
			FullName: "Admin User",
		},
	}
	c.Set(middleware.PrincipalKey, adminPrincipal)

	mockRepo := new(MockConfigRepo)
	mockAudit := new(MockAuditRepo)

	existing := &domain.PlatformGeneralSettings{
		ID:              uuid.New(),
		PlatformName:    "Curexal",
		SupportEmail:    "support@curexal.space",
		DefaultCountry:  "NG",
		DefaultCurrency: "NGN",
		Version:         1,
	}
	updated := *existing
	updated.PlatformName = "Curexal Health OS"
	updated.SupportEmail = "support@curexal.com"
	updated.Version = 2

	mockRepo.On("GetGeneralSettings", mock.Anything).Return(existing, nil)
	mockRepo.On("UpdateGeneralSettings", mock.Anything, mock.Anything, mock.Anything).Return(&updated, nil)
	mockAudit.On("Create", mock.Anything, mock.Anything).Return(&auditDomain.AuditLog{}, nil)

	svc := application.NewPlatformConfigService(mockRepo, mockAudit)
	handler := api.NewPlatformConfigHandler(svc)

	err := handler.UpdatePlatformConfig(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Curexal Health OS")
	assert.Contains(t, rec.Body.String(), "support@curexal.com")
}

