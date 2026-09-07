package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/audit/application"
	"github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/google/uuid"
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

func TestAuditService_LogEvent_HIPAA_Fields(t *testing.T) {
	mockRepo := new(MockAuditRepo)
	svc := application.NewAuditApplicationService(nil, mockRepo)

	ctx := context.Background()
	patientIDStr := uuid.New().String()
	patientUUID, _ := uuid.Parse(patientIDStr)
	branchIDStr := uuid.New().String()
	branchUUID, _ := uuid.Parse(branchIDStr)
	actorEmail := "doctor@curexal.com"
	logID := uuid.New()

	payload := &domain.CreateAuditLogPayload{
		Action:           "patient.chart.viewed",
		PatientID:        &patientIDStr,
		FacilityBranchID: &branchIDStr,
		ActorEmail:       &actorEmail,
		IsBreakGlass:     true,
		Severity:         "WARN",
		Status:           "SUCCESS",
	}

	expectedLog := &domain.AuditLog{
		ID:               logID,
		OccurredAt:       time.Now().UTC(),
		Action:           "patient.chart.viewed",
		PatientID:        &patientUUID,
		FacilityBranchID: &branchUUID,
		ActorEmail:       &actorEmail,
		IsBreakGlass:     true,
		Severity:         "WARN",
		Status:           "SUCCESS",
	}

	mockRepo.On("Create", ctx, payload).Return(expectedLog, nil)

	result, err := svc.LogEvent(ctx, payload)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, logID, result.ID)
	assert.True(t, result.IsBreakGlass)
	assert.Equal(t, &patientUUID, result.PatientID)
	assert.Equal(t, &branchUUID, result.FacilityBranchID)
	assert.Equal(t, &actorEmail, result.ActorEmail)

	mockRepo.AssertExpectations(t)
}

func TestAuditService_ListPatientDisclosures_HIPAA_Accounting(t *testing.T) {
	mockRepo := new(MockAuditRepo)
	svc := application.NewAuditApplicationService(nil, mockRepo)

	ctx := context.Background()
	patientUUID := uuid.New()
	limit := 50
	offset := 0

	expectedLogs := []domain.AuditLog{
		{
			ID:         uuid.New(),
			Action:     "patient.chart.viewed",
			PatientID:  &patientUUID,
			OccurredAt: time.Now().UTC(),
		},
		{
			ID:         uuid.New(),
			Action:     "laboratory.result.viewed",
			PatientID:  &patientUUID,
			OccurredAt: time.Now().UTC().Add(-10 * time.Minute),
		},
	}

	mockRepo.On("ListPatientDisclosures", ctx, patientUUID, limit, offset).Return(expectedLogs, nil)

	results, err := svc.ListPatientDisclosures(ctx, patientUUID, limit, offset)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "patient.chart.viewed", results[0].Action)
	assert.Equal(t, "laboratory.result.viewed", results[1].Action)

	mockRepo.AssertExpectations(t)
}
