package application_test

import (
	"context"
	"testing"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/platform/application"
	"github.com/golangnigeria/curexal/internal/modules/platform/domain"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLaunchGateRepo struct {
	mock.Mock
}

func (m *MockLaunchGateRepo) SaveLaunchGateAudit(ctx context.Context, audit *domain.LaunchGateAudit, actorID *uuid.UUID) (*domain.LaunchGateAudit, error) {
	args := m.Called(ctx, audit, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LaunchGateAudit), args.Error(1)
}

func (m *MockLaunchGateRepo) GetLatestLaunchGateAudit(ctx context.Context, gateName string) (*domain.LaunchGateAudit, error) {
	args := m.Called(ctx, gateName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LaunchGateAudit), args.Error(1)
}

func (m *MockLaunchGateRepo) SaveSystemHealthMetric(ctx context.Context, metric *domain.SystemHealthMetric) (*domain.SystemHealthMetric, error) {
	args := m.Called(ctx, metric)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SystemHealthMetric), args.Error(1)
}

func (m *MockLaunchGateRepo) ListSystemHealthMetrics(ctx context.Context) ([]domain.SystemHealthMetric, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.SystemHealthMetric), args.Error(1)
}

type MockAuditRepo struct {
	mock.Mock
}

func (m *MockAuditRepo) GetByID(ctx context.Context, id uuid.UUID) (*auditDomain.AuditLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditDomain.AuditLog), args.Error(1)
}

func (m *MockAuditRepo) Create(ctx context.Context, payload *auditDomain.CreateAuditLogPayload) (*auditDomain.AuditLog, error) {
	args := m.Called(ctx, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditDomain.AuditLog), args.Error(1)
}

func (m *MockAuditRepo) ListTenantLogs(ctx context.Context, tenantID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]auditDomain.AuditLog, error) {
	args := m.Called(ctx, tenantID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]auditDomain.AuditLog), args.Error(1)
}

func (m *MockAuditRepo) ListPlatformLogs(ctx context.Context, orgID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]auditDomain.AuditLog, error) {
	args := m.Called(ctx, orgID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]auditDomain.AuditLog), args.Error(1)
}

func (m *MockAuditRepo) GetStats(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID) (*auditDomain.AdminStats, error) {
	args := m.Called(ctx, tenantID, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditDomain.AdminStats), args.Error(1)
}

func (m *MockAuditRepo) ListAll(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID, limit int, offset int) ([]auditDomain.AuditLog, error) {
	args := m.Called(ctx, tenantID, orgID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]auditDomain.AuditLog), args.Error(1)
}

func TestLaunchGateService_VerifyProductionReadiness_Evaluates10Checks_Success(t *testing.T) {
	mockLaunchGateRepo := new(MockLaunchGateRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewLaunchGateService(mockLaunchGateRepo, mockAuditRepo)

	actorID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Role:   "super_admin",
		Platform: platformAuth.PlatformVector{
			IsSuperAdmin:    true,
			IsPlatformAdmin: true,
		},
	}

	mockLaunchGateRepo.On("SaveLaunchGateAudit", mock.Anything, mock.MatchedBy(func(a *domain.LaunchGateAudit) bool {
		return a.GateName == "PHASE_10_FINAL_PRODUCTION_GATE" && a.Status == "PASSED"
	}), &actorID).Return(&domain.LaunchGateAudit{
		ID:       uuid.New(),
		GateName: "PHASE_10_FINAL_PRODUCTION_GATE",
		Status:   "PASSED",
	}, nil)

	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "LAUNCH_GATE_VERIFIED"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.VerifyProductionReadiness(context.Background(), principal)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "PASSED", res.Status)

	mockLaunchGateRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestLaunchGateService_VerifyProductionReadiness_NonAdmin_Fails(t *testing.T) {
	mockLaunchGateRepo := new(MockLaunchGateRepo)
	svc := application.NewLaunchGateService(mockLaunchGateRepo, nil)

	principal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Role:   "clinician",
	}

	res, err := svc.VerifyProductionReadiness(context.Background(), principal)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "only platform administrators")

	mockLaunchGateRepo.AssertNotCalled(t, "SaveLaunchGateAudit", mock.Anything, mock.Anything, mock.Anything)
}
