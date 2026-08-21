package application_test

import (
	"context"
	"testing"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/organization/application"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBrandingRepo struct {
	mock.Mock
}

func (m *MockBrandingRepo) GetBranding(ctx context.Context, orgID uuid.UUID) (*domain.BrandingConfig, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BrandingConfig), args.Error(1)
}

func (m *MockBrandingRepo) UpdateBranding(ctx context.Context, orgID uuid.UUID, payload *domain.UpdateBrandingPayload, actorID uuid.UUID) (*domain.BrandingConfig, error) {
	args := m.Called(ctx, orgID, payload, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BrandingConfig), args.Error(1)
}

func (m *MockBrandingRepo) SaveNotificationConfig(ctx context.Context, config *domain.NotificationConfig, actorID uuid.UUID) (*domain.NotificationConfig, error) {
	args := m.Called(ctx, config, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NotificationConfig), args.Error(1)
}

func (m *MockBrandingRepo) ListNotificationConfigs(ctx context.Context, orgID uuid.UUID) ([]domain.NotificationConfig, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.NotificationConfig), args.Error(1)
}

func (m *MockBrandingRepo) GetNotificationConfig(ctx context.Context, orgID uuid.UUID, channel, provider string) (*domain.NotificationConfig, error) {
	args := m.Called(ctx, orgID, channel, provider)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NotificationConfig), args.Error(1)
}

func (m *MockBrandingRepo) SaveNotificationTemplate(ctx context.Context, template *domain.NotificationTemplate, actorID uuid.UUID) (*domain.NotificationTemplate, error) {
	args := m.Called(ctx, template, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NotificationTemplate), args.Error(1)
}

func (m *MockBrandingRepo) ListNotificationTemplates(ctx context.Context, orgID uuid.UUID) ([]domain.NotificationTemplate, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.NotificationTemplate), args.Error(1)
}

func (m *MockBrandingRepo) GetNotificationTemplate(ctx context.Context, orgID uuid.UUID, templateKey, channel string) (*domain.NotificationTemplate, error) {
	args := m.Called(ctx, orgID, templateKey, channel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NotificationTemplate), args.Error(1)
}

func (m *MockBrandingRepo) CreateUserNotification(ctx context.Context, notif *domain.UserNotification) (*domain.UserNotification, error) {
	args := m.Called(ctx, notif)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserNotification), args.Error(1)
}

func (m *MockBrandingRepo) ListUserNotifications(ctx context.Context, orgID, userID uuid.UUID, limit int) ([]domain.UserNotification, error) {
	args := m.Called(ctx, orgID, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserNotification), args.Error(1)
}

func (m *MockBrandingRepo) MarkNotificationRead(ctx context.Context, orgID, userID, notifID uuid.UUID) error {
	args := m.Called(ctx, orgID, userID, notifID)
	return args.Error(0)
}

func (m *MockBrandingRepo) MarkAllNotificationsRead(ctx context.Context, orgID, userID uuid.UUID) error {
	args := m.Called(ctx, orgID, userID)
	return args.Error(0)
}

func (m *MockBrandingRepo) CreateNotificationDelivery(ctx context.Context, delivery *domain.NotificationDelivery) (*domain.NotificationDelivery, error) {
	args := m.Called(ctx, delivery)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NotificationDelivery), args.Error(1)
}

func (m *MockBrandingRepo) ListNotificationDeliveries(ctx context.Context, orgID uuid.UUID, limit int) ([]domain.NotificationDelivery, error) {
	args := m.Called(ctx, orgID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.NotificationDelivery), args.Error(1)
}

func TestOrganizationBrandingService_SaveNotificationConfig_EncryptsSecrets_RedactsResponse(t *testing.T) {
	mockBrandingRepo := new(MockBrandingRepo)
	mockAuditRepo := new(MockAuditRepo)
	masterKey := "TEST_MASTER_ENCRYPTION_KEY_2026"
	svc := application.NewOrganizationBrandingService(mockBrandingRepo, nil, mockAuditRepo, masterKey)

	orgID := uuid.New()
	actorID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	rawAPIKey := "re_1234567890_secret_resend_api_key"
	payload := &domain.SaveNotificationConfigPayload{
		Channel:  "EMAIL",
		Provider: "RESEND",
		APIKey:   &rawAPIKey,
	}

	// Verify that the config passed to repository contains encrypted AEAD payload (NOT raw API key)
	mockBrandingRepo.On("SaveNotificationConfig", mock.Anything, mock.MatchedBy(func(cfg *domain.NotificationConfig) bool {
		if cfg.APIKey == nil || *cfg.APIKey == rawAPIKey {
			return false // API key MUST NOT be stored in plaintext
		}
		// Decrypt to verify exact AEAD roundtrip
		dec, err := crypto.DecryptAEAD(*cfg.APIKey, masterKey)
		return err == nil && dec == rawAPIKey
	}), actorID).Return(&domain.NotificationConfig{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Channel:        "EMAIL",
		Provider:       "RESEND",
		APIKey:         payload.APIKey, // Repo returns saved pointer
		IsActive:       true,
	}, nil)

	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "NOTIFICATION_CONFIG_SAVED"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.SaveNotificationConfig(context.Background(), principal, payload)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.APIKey)
	assert.Equal(t, "••••••••", *res.APIKey) // Secret API key MUST be redacted in service response envelope

	mockBrandingRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestOrganizationBrandingService_UpdateBranding_Success(t *testing.T) {
	mockBrandingRepo := new(MockBrandingRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewOrganizationBrandingService(mockBrandingRepo, nil, mockAuditRepo, "")

	orgID := uuid.New()
	actorID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	color := "#10B981"
	domainStr := "portal.everight.curexal.health"
	payload := &domain.UpdateBrandingPayload{
		PrimaryColor: &color,
		CustomDomain: &domainStr,
		Version:      1,
	}

	updatedBranding := &domain.BrandingConfig{
		OrganizationID:     orgID,
		PrimaryColor:       "#10B981",
		CustomDomain:       &domainStr,
		CustomDomainStatus: "PENDING",
		Version:            2,
	}

	mockBrandingRepo.On("UpdateBranding", mock.Anything, orgID, payload, actorID).Return(updatedBranding, nil)
	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "BRANDING_UPDATED"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.UpdateBranding(context.Background(), principal, payload)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "#10B981", res.PrimaryColor)
	assert.Equal(t, "PENDING", res.CustomDomainStatus)

	mockBrandingRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestOrganizationBrandingService_ListUserNotifications_Success(t *testing.T) {
	mockBrandingRepo := new(MockBrandingRepo)
	svc := application.NewOrganizationBrandingService(mockBrandingRepo, nil, nil, "")

	orgID := uuid.New()
	userID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: userID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	expectedNotifs := []domain.UserNotification{
		{
			ID:               uuid.New(),
			OrganizationID:   orgID,
			UserID:           userID,
			NotificationType: "LAB_RESULT_READY",
			Title:            "Lab Result Ready",
			Body:             "Your Full Blood Count result is now available.",
		},
	}

	mockBrandingRepo.On("ListUserNotifications", mock.Anything, orgID, userID, 50).Return(expectedNotifs, nil)

	res, err := svc.ListUserNotifications(context.Background(), principal, 50)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "LAB_RESULT_READY", res[0].NotificationType)

	mockBrandingRepo.AssertExpectations(t)
}
