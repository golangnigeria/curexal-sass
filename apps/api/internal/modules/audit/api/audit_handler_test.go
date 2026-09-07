package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/audit/api"
	"github.com/golangnigeria/curexal/internal/modules/audit/application"
	"github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockAuditRepo struct {
	mock.Mock
}

func (m *MockAuditRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuditLog), args.Error(1)
}

func (m *MockAuditRepo) Create(ctx context.Context, payload *domain.CreateAuditLogPayload) (*domain.AuditLog, error) {
	args := m.Called(ctx, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuditLog), args.Error(1)
}

func (m *MockAuditRepo) ListTenantLogs(ctx context.Context, tenantID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]domain.AuditLog, error) {
	args := m.Called(ctx, tenantID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.AuditLog), args.Error(1)
}

func (m *MockAuditRepo) ListPlatformLogs(ctx context.Context, orgID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]domain.AuditLog, error) {
	args := m.Called(ctx, orgID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.AuditLog), args.Error(1)
}

func (m *MockAuditRepo) GetStats(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID) (*domain.AdminStats, error) {
	args := m.Called(ctx, tenantID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AdminStats), args.Error(1)
}

func (m *MockAuditRepo) ListAll(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID, limit int, offset int) ([]domain.AuditLog, error) {
	args := m.Called(ctx, tenantID, orgID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.AuditLog), args.Error(1)
}

func (m *MockAuditRepo) ListPatientDisclosures(ctx context.Context, patientID uuid.UUID, limit, offset int) ([]domain.AuditLog, error) {
	args := m.Called(ctx, patientID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.AuditLog), args.Error(1)
}

func TestGetPatientDisclosures_InvalidUUID_Returns400(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/audit-logs/patient/invalid-uuid/disclosures", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("patientId")
	c.SetParamValues("invalid-uuid")

	mockRepo := new(MockAuditRepo)
	appSvc := application.NewAuditApplicationService(nil, mockRepo)
	h := api.NewAuditHandler(nil, appSvc)

	err := h.GetPatientDisclosures(c)
	require.Error(t, err)

	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestGetPatientDisclosures_Success(t *testing.T) {
	e := echo.New()
	patientID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/audit-logs/patient/"+patientID.String()+"/disclosures", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("patientId")
	c.SetParamValues(patientID.String())

	mockRepo := new(MockAuditRepo)
	expectedLogs := []domain.AuditLog{
		{
			ID:         uuid.New(),
			Action:     "patient.chart.viewed",
			PatientID:  &patientID,
			OccurredAt: time.Now().UTC(),
		},
	}
	mockRepo.On("ListPatientDisclosures", mock.Anything, patientID, 50, 0).Return(expectedLogs, nil)

	appSvc := application.NewAuditApplicationService(nil, mockRepo)
	// mock server with logger
	logger := zerolog.Nop()
	_ = logger
	h := api.NewAuditHandler(nil, appSvc)

	err := h.GetPatientDisclosures(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "patient.chart.viewed")
}
