package application_test

import (
	"context"
	"encoding/json"
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

type MockIntegrationRepo struct {
	mock.Mock
}

func (m *MockIntegrationRepo) CreateAPIKey(ctx context.Context, apiKey *domain.APIKey, keyHash string, actorID uuid.UUID) (*domain.APIKey, error) {
	args := m.Called(ctx, apiKey, keyHash, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.APIKey), args.Error(1)
}

func (m *MockIntegrationRepo) GetAPIKeyByID(ctx context.Context, orgID, keyID uuid.UUID) (*domain.APIKey, error) {
	args := m.Called(ctx, orgID, keyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.APIKey), args.Error(1)
}

func (m *MockIntegrationRepo) GetAPIKeyByHash(ctx context.Context, keyHash string) (*domain.APIKey, error) {
	args := m.Called(ctx, keyHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.APIKey), args.Error(1)
}

func (m *MockIntegrationRepo) ListAPIKeys(ctx context.Context, orgID uuid.UUID) ([]domain.APIKey, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.APIKey), args.Error(1)
}

func (m *MockIntegrationRepo) RevokeAPIKey(ctx context.Context, orgID, keyID uuid.UUID, actorID uuid.UUID) error {
	args := m.Called(ctx, orgID, keyID, actorID)
	return args.Error(0)
}

func (m *MockIntegrationRepo) CreateWebhookSubscription(ctx context.Context, sub *domain.WebhookSubscription, actorID uuid.UUID) (*domain.WebhookSubscription, error) {
	args := m.Called(ctx, sub, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookSubscription), args.Error(1)
}

func (m *MockIntegrationRepo) GetWebhookSubscriptionByID(ctx context.Context, orgID, subID uuid.UUID) (*domain.WebhookSubscription, error) {
	args := m.Called(ctx, orgID, subID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookSubscription), args.Error(1)
}

func (m *MockIntegrationRepo) ListWebhookSubscriptions(ctx context.Context, orgID uuid.UUID) ([]domain.WebhookSubscription, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.WebhookSubscription), args.Error(1)
}

func (m *MockIntegrationRepo) DeleteWebhookSubscription(ctx context.Context, orgID, subID uuid.UUID, actorID uuid.UUID) error {
	args := m.Called(ctx, orgID, subID, actorID)
	return args.Error(0)
}

func (m *MockIntegrationRepo) CreateWebhookDelivery(ctx context.Context, delivery *domain.WebhookDelivery) (*domain.WebhookDelivery, error) {
	args := m.Called(ctx, delivery)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookDelivery), args.Error(1)
}

func (m *MockIntegrationRepo) ListWebhookDeliveries(ctx context.Context, orgID uuid.UUID, limit int) ([]domain.WebhookDelivery, error) {
	args := m.Called(ctx, orgID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.WebhookDelivery), args.Error(1)
}

func TestOrganizationIntegrationService_CreateAPIKey_GeneratesKey_HashesToken_ReturnsRawOnce(t *testing.T) {
	mockIntegrationRepo := new(MockIntegrationRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewOrganizationIntegrationService(mockIntegrationRepo, nil, mockAuditRepo, "")

	orgID := uuid.New()
	actorID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	scopesJSON, _ := json.Marshal([]string{"patients:read", "appointments:write"})
	payload := &domain.CreateAPIKeyPayload{
		Name:   "EMR System API Key",
		Scopes: scopesJSON,
	}

	mockIntegrationRepo.On("CreateAPIKey", mock.Anything, mock.MatchedBy(func(k *domain.APIKey) bool {
		return k.Name == "EMR System API Key" && k.KeyPrefix != ""
	}), mock.MatchedBy(func(hash string) bool {
		return len(hash) == 64 // SHA-256 hex string length
	}), actorID).Return(&domain.APIKey{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           "EMR System API Key",
		KeyPrefix:      "cx_live_1234",
		IsActive:       true,
	}, nil)

	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "API_KEY_CREATED"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.CreateAPIKey(context.Background(), principal, payload)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.RawKey, "cx_live_") // Raw key returned ONCE
	assert.Equal(t, "EMR System API Key", res.APIKey.Name)

	mockIntegrationRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestOrganizationIntegrationService_CreateWebhookSubscription_SSRFBoundary_RejectsForbiddenURLs(t *testing.T) {
	mockIntegrationRepo := new(MockIntegrationRepo)
	svc := application.NewOrganizationIntegrationService(mockIntegrationRepo, nil, nil, "")

	orgID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	forbiddenURLs := []string{
		"http://127.0.0.1:8080/webhook",
		"http://localhost:3000/webhook",
		"http://169.254.169.254/latest/meta-data",
		"http://10.0.0.15/receiver",
		"http://192.168.1.50/receiver",
	}

	for _, u := range forbiddenURLs {
		payload := &domain.CreateWebhookSubscriptionPayload{
			Name:      "SSRF Attack Target",
			TargetURL: u,
		}
		res, err := svc.CreateWebhookSubscription(context.Background(), principal, payload)
		assert.ErrorIs(t, err, domain.ErrSSRFURLForbidden, "URL should be forbidden: "+u)
		assert.Nil(t, res)
	}

	mockIntegrationRepo.AssertNotCalled(t, "CreateWebhookSubscription", mock.Anything, mock.Anything, mock.Anything)
}

func TestOrganizationIntegrationService_CreateWebhookSubscription_ValidURL_EncryptsSigningSecret(t *testing.T) {
	mockIntegrationRepo := new(MockIntegrationRepo)
	mockAuditRepo := new(MockAuditRepo)
	masterKey := "TEST_MASTER_ENCRYPTION_KEY_2026"
	svc := application.NewOrganizationIntegrationService(mockIntegrationRepo, nil, mockAuditRepo, masterKey)

	orgID := uuid.New()
	actorID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	rawSecret := "whsec_custom_secret_key_12345"
	payload := &domain.CreateWebhookSubscriptionPayload{
		Name:          "Valid Patient Webhook",
		TargetURL:     "https://api.everight-health.com/v1/curexal-webhook",
		SigningSecret: &rawSecret,
	}

	mockIntegrationRepo.On("CreateWebhookSubscription", mock.Anything, mock.MatchedBy(func(sub *domain.WebhookSubscription) bool {
		if sub.SigningSecret == nil || *sub.SigningSecret == rawSecret {
			return false // Signing secret MUST NOT be stored in plaintext
		}
		dec, err := crypto.DecryptAEAD(*sub.SigningSecret, masterKey)
		return err == nil && dec == rawSecret
	}), actorID).Return(&domain.WebhookSubscription{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           "Valid Patient Webhook",
		TargetURL:      payload.TargetURL,
		SigningSecret:  payload.SigningSecret,
		IsActive:       true,
	}, nil)

	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "WEBHOOK_SUBSCRIPTION_CREATED"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.CreateWebhookSubscription(context.Background(), principal, payload)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.SigningSecret)
	assert.Equal(t, "••••••••", *res.SigningSecret) // Secret MUST be redacted in service response envelope

	mockIntegrationRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestOrganizationIntegrationService_ComputeHMACSignature_ValidatesSymmetry(t *testing.T) {
	secret := "whsec_super_secret_signing_key"
	timestamp := "1770912000"
	payload := `{"event":"PATIENT_REGISTERED","patient_id":"123"}`

	sig1 := application.ComputeHMACSignature(secret, timestamp, payload)
	sig2 := application.ComputeHMACSignature(secret, timestamp, payload)

	assert.NotEmpty(t, sig1)
	assert.Equal(t, sig1, sig2) // Deterministic HMAC-SHA256 calculation

	// Tampered payload produces different signature
	tamperedSig := application.ComputeHMACSignature(secret, timestamp, `{"event":"PATIENT_REGISTERED","patient_id":"999"}`)
	assert.NotEqual(t, sig1, tamperedSig)
}
