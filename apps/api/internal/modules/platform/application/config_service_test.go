package application_test

import (
	"context"
	"testing"
	"time"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/platform/application"
	"github.com/golangnigeria/curexal/internal/modules/platform/domain"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConfigRepository struct {
	mock.Mock
}

func (m *MockConfigRepository) GetGeneralSettings(ctx context.Context) (*domain.PlatformGeneralSettings, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PlatformGeneralSettings), args.Error(1)
}

func (m *MockConfigRepository) UpdateGeneralSettings(ctx context.Context, settings *domain.PlatformGeneralSettings, updatedBy uuid.UUID) (*domain.PlatformGeneralSettings, error) {
	args := m.Called(ctx, settings, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PlatformGeneralSettings), args.Error(1)
}

func (m *MockConfigRepository) GetSecurityPolicy(ctx context.Context) (*domain.IdentitySecurityPolicy, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IdentitySecurityPolicy), args.Error(1)
}

func (m *MockConfigRepository) UpdateSecurityPolicy(ctx context.Context, policy *domain.IdentitySecurityPolicy, updatedBy uuid.UUID) (*domain.IdentitySecurityPolicy, error) {
	args := m.Called(ctx, policy, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IdentitySecurityPolicy), args.Error(1)
}

type MockAuditRepository struct {
	mock.Mock
}

func (m *MockAuditRepository) Create(ctx context.Context, payload *auditDomain.CreateAuditLogPayload) (*auditDomain.AuditLog, error) {
	args := m.Called(ctx, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditDomain.AuditLog), args.Error(1)
}

func (m *MockAuditRepository) GetByID(ctx context.Context, id uuid.UUID) (*auditDomain.AuditLog, error) {
	args := m.Called(ctx, id)
	return nil, args.Error(1)
}

func (m *MockAuditRepository) ListTenantLogs(ctx context.Context, tenantID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]auditDomain.AuditLog, error) {
	return nil, nil
}

func (m *MockAuditRepository) ListPlatformLogs(ctx context.Context, orgID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]auditDomain.AuditLog, error) {
	return nil, nil
}

func (m *MockAuditRepository) GetStats(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID) (*auditDomain.AdminStats, error) {
	return nil, nil
}

func (m *MockAuditRepository) ListAll(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID, limit int, offset int) ([]auditDomain.AuditLog, error) {
	return nil, nil
}

func TestPlatformConfigService_UpdateGeneralSettings_Success(t *testing.T) {
	mockRepo := new(MockConfigRepository)
	mockAudit := new(MockAuditRepository)
	service := application.NewPlatformConfigService(mockRepo, mockAudit)

	adminID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: adminID.String(),
		Role:   "platform_admin",
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: true,
		},
		Identity: platformAuth.IdentityVector{
			FullName: "Platform Admin",
		},
	}

	configID := uuid.New()
	existingSettings := &domain.PlatformGeneralSettings{
		ID:                  configID,
		PlatformName:        "Curexal Old",
		SupportEmail:        "old@curexal.com",
		DefaultCountry:      "NG",
		DefaultCurrency:     "NGN",
		SupportedCountries:  []string{"NG", "GH"},
		SupportedCurrencies: []string{"NGN", "USD"},
		DefaultTimezone:     "Africa/Lagos",
		DefaultLocale:       "en",
		DateFormat:          "YYYY-MM-DD",
		TimeFormat:          "HH:mm",
		NumberFormat:        "standard",
		MeasurementUnits:    "metric",
		Status:              "ACTIVE",
		Version:             1,
	}

	inputSettings := &domain.PlatformGeneralSettings{
		PlatformName: "Curexal Health OS",
		SupportEmail: "support@curexal.com",
	}

	updatedSettings := *existingSettings
	updatedSettings.PlatformName = "Curexal Health OS"
	updatedSettings.SupportEmail = "support@curexal.com"
	updatedSettings.Version = 2
	updatedSettings.UpdatedAt = time.Now()

	mockRepo.On("GetGeneralSettings", mock.Anything).Return(existingSettings, nil)
	mockRepo.On("UpdateGeneralSettings", mock.Anything, mock.Anything, adminID).Return(&updatedSettings, nil)
	mockAudit.On("Create", mock.Anything, mock.Anything).Return(&auditDomain.AuditLog{}, nil)

	res, err := service.UpdateGeneralSettings(context.Background(), principal, inputSettings)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Curexal Health OS", res.PlatformName)
	assert.Equal(t, "support@curexal.com", res.SupportEmail)
	assert.Equal(t, 2, res.Version)
	mockRepo.AssertExpectations(t)
	mockAudit.AssertExpectations(t)
}


func TestPlatformConfigService_UpdateGeneralSettings_Unauthorized(t *testing.T) {
	mockRepo := new(MockConfigRepository)
	mockAudit := new(MockAuditRepository)
	service := application.NewPlatformConfigService(mockRepo, mockAudit)

	userID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: userID.String(),
		Role:   "regular_user",
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: false,
		},
	}

	inputSettings := &domain.PlatformGeneralSettings{
		PlatformName:    "Curexal Health",
		SupportEmail:    "user@curexal.com",
		DefaultCountry:  "NG",
		DefaultCurrency: "NGN",
	}

	res, err := service.UpdateGeneralSettings(context.Background(), principal, inputSettings)

	assert.ErrorIs(t, err, domain.ErrUnauthorizedPlatformAdmin)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "UpdateGeneralSettings", mock.Anything, mock.Anything, mock.Anything)
}

func TestPlatformConfigService_UpdateSecurityPolicy_InvalidMinPassword(t *testing.T) {
	mockRepo := new(MockConfigRepository)
	mockAudit := new(MockAuditRepository)
	service := application.NewPlatformConfigService(mockRepo, mockAudit)

	adminID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: adminID.String(),
		Role:   "platform_admin",
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: true,
		},
	}

	existingPolicy := &domain.IdentitySecurityPolicy{
		ID:                            uuid.New(),
		MinPasswordLength:             8,
		MaxFailedLoginAttempts:        5,
		AccountLockoutDurationMinutes: 30,
		Version:                       1,
	}
	mockRepo.On("GetSecurityPolicy", mock.Anything).Return(existingPolicy, nil)

	policy := &domain.IdentitySecurityPolicy{
		MinPasswordLength:             6, // Invalid: must be >= 8
		MaxFailedLoginAttempts:        5,
		AccountLockoutDurationMinutes: 30,
	}

	res, err := service.UpdateSecurityPolicy(context.Background(), principal, policy)

	assert.ErrorIs(t, err, domain.ErrInvalidSecurityPolicy)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "UpdateSecurityPolicy", mock.Anything, mock.Anything, mock.Anything)
}

func TestPlatformConfigService_UpdateSecurityPolicy_Success(t *testing.T) {
	mockRepo := new(MockConfigRepository)
	mockAudit := new(MockAuditRepository)
	service := application.NewPlatformConfigService(mockRepo, mockAudit)

	adminID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: adminID.String(),
		Role:   "platform_admin",
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: true,
		},
		Identity: platformAuth.IdentityVector{
			FullName: "Platform Admin",
		},
	}

	policyID := uuid.New()
	existingPolicy := &domain.IdentitySecurityPolicy{
		ID:                            policyID,
		MinPasswordLength:             8,
		MaxFailedLoginAttempts:        5,
		AccountLockoutDurationMinutes: 30,
		Version:                       1,
	}

	inputPolicy := &domain.IdentitySecurityPolicy{
		MinPasswordLength: 12,
	}

	updatedPolicy := *existingPolicy
	updatedPolicy.MinPasswordLength = 12
	updatedPolicy.Version = 2

	mockRepo.On("GetSecurityPolicy", mock.Anything).Return(existingPolicy, nil)
	mockRepo.On("UpdateSecurityPolicy", mock.Anything, mock.Anything, adminID).Return(&updatedPolicy, nil)
	mockAudit.On("Create", mock.Anything, mock.Anything).Return(&auditDomain.AuditLog{}, nil)

	res, err := service.UpdateSecurityPolicy(context.Background(), principal, inputPolicy)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 12, res.MinPasswordLength)
	assert.Equal(t, 2, res.Version)
	mockRepo.AssertExpectations(t)
	mockAudit.AssertExpectations(t)
}

