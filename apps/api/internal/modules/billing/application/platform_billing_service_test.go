package application_test

import (
	"context"
	"testing"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/billing/application"
	"github.com/golangnigeria/curexal/internal/modules/billing/domain"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBillingRepo struct {
	mock.Mock
}

func (m *MockBillingRepo) ListPricingRules(ctx context.Context) ([]domain.PricingRule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PricingRule), args.Error(1)
}

func (m *MockBillingRepo) UpdatePricingRule(ctx context.Context, rule *domain.PricingRule, updatedBy uuid.UUID) (*domain.PricingRule, error) {
	args := m.Called(ctx, rule, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PricingRule), args.Error(1)
}

func (m *MockBillingRepo) ListPaymentGateways(ctx context.Context) ([]domain.PaymentGatewayConfig, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PaymentGatewayConfig), args.Error(1)
}

func (m *MockBillingRepo) GetPaymentGatewayByProvider(ctx context.Context, providerCode string) (*domain.PaymentGatewayConfig, error) {
	args := m.Called(ctx, providerCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PaymentGatewayConfig), args.Error(1)
}

func (m *MockBillingRepo) UpdatePaymentGateway(ctx context.Context, gateway *domain.PaymentGatewayConfig, updatedBy uuid.UUID) (*domain.PaymentGatewayConfig, error) {
	args := m.Called(ctx, gateway, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PaymentGatewayConfig), args.Error(1)
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

func TestPaymentGatewayVaultService_UpdateGateway_EncryptsSecretAndRedacts(t *testing.T) {
	mockRepo := new(MockBillingRepo)
	mockAudit := new(MockAuditRepo)
	masterKey := "test_master_encryption_key_2026"
	svc := application.NewPaymentGatewayVaultService(mockRepo, mockAudit, masterKey)

	adminID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: adminID.String(),
		Role:   "platform_admin",
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: true,
		},
		Identity: platformAuth.IdentityVector{
			FullName: "Super Admin",
		},
	}

	existingGateway := &domain.PaymentGatewayConfig{
		ID:                  uuid.New(),
		ProviderCode:        "paystack",
		Name:                "Paystack",
		IsEnabled:           true,
		Priority:            1,
		SupportedCurrencies: []string{"NGN", "USD"},
		EncryptedSecretKey:  "old_ciphertext",
		Version:             1,
	}

	mockRepo.On("GetPaymentGatewayByProvider", mock.Anything, "paystack").Return(existingGateway, nil)

	plaintextSecret := "sk_live_paystack_super_secret_9988"

	mockRepo.On("UpdatePaymentGateway", mock.Anything, mock.MatchedBy(func(g *domain.PaymentGatewayConfig) bool {
		// Verify ciphertext is NOT plaintext and can be decrypted back with masterKey
		decrypted, err := crypto.DecryptAEAD(g.EncryptedSecretKey, masterKey)
		return err == nil && decrypted == plaintextSecret
	}), adminID).Return(existingGateway, nil)

	var loggedPayload *auditDomain.CreateAuditLogPayload
	mockAudit.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		loggedPayload = p
		return p.Action == "PAYMENT_GATEWAY_UPDATED"
	})).Return(&auditDomain.AuditLog{}, nil)

	updatePayload := &application.UpdateGatewayPayload{
		Name:                "Paystack Updated",
		IsEnabled:           true,
		Priority:            1,
		SupportedCurrencies: []string{"NGN", "USD"},
		PlaintextSecretKey:  plaintextSecret,
		Version:             1,
	}

	res, err := svc.UpdateGateway(context.Background(), principal, "paystack", updatePayload)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "••••••••", res.RedactedSecretKey)
	assert.NotEqual(t, plaintextSecret, res.EncryptedSecretKey)

	// Verify Audit Log payload NEVER contained secret key
	assert.NotNil(t, loggedPayload)
	assert.Nil(t, loggedPayload.BeforeState)
	assert.Nil(t, loggedPayload.AfterState)

	mockRepo.AssertExpectations(t)
	mockAudit.AssertExpectations(t)
}

func TestPaymentGatewayVaultService_ListGateways_RedactsSecrets(t *testing.T) {
	mockRepo := new(MockBillingRepo)
	svc := application.NewPaymentGatewayVaultService(mockRepo, nil, "key")

	mockRepo.On("ListPaymentGateways", mock.Anything).Return([]domain.PaymentGatewayConfig{
		{ProviderCode: "paystack", EncryptedSecretKey: "ae_secret1"},
		{ProviderCode: "stripe", EncryptedSecretKey: "ae_secret2"},
	}, nil)

	res, err := svc.ListGateways(context.Background())

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "••••••••", res[0].RedactedSecretKey)
	assert.Equal(t, "••••••••", res[1].RedactedSecretKey)
}

func TestPlatformPricingService_UpdatePricingRule_Unauthorized(t *testing.T) {
	mockRepo := new(MockBillingRepo)
	svc := application.NewPlatformPricingService(mockRepo, nil)

	principal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Role:   "organization_user",
	}

	rule := &domain.PricingRule{
		TargetType:   "plan",
		TargetCode:   "smart",
		Currency:     "NGN",
		MonthlyPrice: 50000.00,
	}

	res, err := svc.UpdatePricingRule(context.Background(), principal, rule)

	assert.ErrorIs(t, err, domain.ErrUnauthorizedPlatformAdmin)
	assert.Nil(t, res)
}
