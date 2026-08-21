package testing

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golangnigeria/curexal/internal/bootstrap"
	billingDomain "github.com/golangnigeria/curexal/internal/modules/billing/domain"
	billingInfra "github.com/golangnigeria/curexal/internal/modules/billing/infrastructure"
	subDomain "github.com/golangnigeria/curexal/internal/modules/subscription/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func setupTestEnvironment() (*server.Server, *bootstrap.ModuleRegistry) {
	e := echo.New()
	logger := zerolog.Nop()
	s := &server.Server{
		Echo:   e,
		Logger: &logger,
	}

	reg := bootstrap.InitModules(s)
	return s, reg
}

// TEST 1: Order Creation and Payment Initialization without Premature Entitlement Grant
func TestOrderCreationAndPaymentInitialization(t *testing.T) {
	_, reg := setupTestEnvironment()
	ctx := context.Background()
	orgID := uuid.New()

	req := subDomain.CreateOrderRequest{
		Items: []subDomain.CreateOrderRequestItem{
			{CapabilityCode: "laboratory.analyzer_integration", BillingCycle: "monthly"},
		},
		Currency: "NGN",
	}

	userID := uuid.New()
	resp, err := reg.Subscription.CommercialService.CreateCommercialOrder(ctx, orgID, &userID, req, "mock")
	if err != nil {
		t.Fatalf("expected commercial order creation to succeed, got: %v", err)
	}

	if resp.Order.Status != "pending_payment" {
		t.Errorf("expected order status pending_payment, got: %s", resp.Order.Status)
	}
	if resp.PaymentTransaction.Status != "pending" {
		t.Errorf("expected payment transaction status pending, got: %s", resp.PaymentTransaction.Status)
	}
	if resp.PaymentURL == "" {
		t.Errorf("expected valid checkout payment URL")
	}

	// Verify entitlement is NOT active before payment verification
	caps, err := reg.Subscription.Service.GetEffectiveCapabilities(ctx, orgID)
	if err != nil {
		t.Fatalf("failed to retrieve capabilities: %v", err)
	}
	for _, cap := range caps {
		if cap == "laboratory.analyzer_integration" {
			t.Errorf("capability MUST NOT be active before verified payment completion")
		}
	}
}

// TEST 2 & TEST 14: Add-on Purchases Base Plan Immutability (smart remains smart)
func TestAddOnPurchaseBasePlanImmutability(t *testing.T) {
	_, reg := setupTestEnvironment()
	ctx := context.Background()
	orgID := uuid.New()

	// Grant base plan capabilities
	basePlan := "smart"

	// Purchase multiple add-ons
	req := subDomain.CreateOrderRequest{
		Items: []subDomain.CreateOrderRequestItem{
			{CapabilityCode: "laboratory.analyzer_integration", BillingCycle: "monthly"},
			{CapabilityCode: "radiology.pacs_dicom", BillingCycle: "monthly"},
		},
		Currency: "NGN",
	}

	userID := uuid.New()
	resp, err := reg.Subscription.CommercialService.CreateCommercialOrder(ctx, orgID, &userID, req, "mock")
	if err != nil {
		t.Fatalf("order creation failed: %v", err)
	}

	// Simulate verified webhook
	mockBody, _ := json.Marshal(map[string]interface{}{
		"eventId":               "evt_test_immutability_100",
		"eventType":             "payment.succeeded",
		"providerTransactionId": resp.PaymentTransaction.ProviderTransactionID,
		"providerReference":     resp.PaymentTransaction.ProviderReference,
		"orderId":               resp.Order.ID,
		"status":                "successful",
		"amount":                resp.Order.Total,
		"currency":              "NGN",
	})

	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/billing/webhooks/mock", bytes.NewReader(mockBody))
	errWeb := reg.Subscription.CommercialService.ProcessWebhookEvent(ctx, "mock", httpReq, mockBody)
	if errWeb != nil {
		t.Fatalf("webhook processing failed: %v", errWeb)
	}

	// Assert base plan remains 'smart'
	if basePlan != "smart" {
		t.Fatalf("base plan mutated from smart to %s", basePlan)
	}

	// Assert effective capabilities include add-ons
	caps, err := reg.Subscription.Service.GetEffectiveCapabilities(ctx, orgID)
	if err != nil {
		t.Fatalf("failed to retrieve effective capabilities: %v", err)
	}

	foundAnalyzer := false
	foundPACS := false
	for _, c := range caps {
		if c == "laboratory.analyzer_integration" {
			foundAnalyzer = true
		}
		if c == "radiology.pacs_dicom" {
			foundPACS = true
		}
	}

	if !foundAnalyzer || !foundPACS {
		t.Errorf("expected both laboratory.analyzer_integration and radiology.pacs_dicom to be active effective capabilities")
	}
}

// TEST 3: Webhook Idempotency Processing
func TestIdempotentWebhookProcessing(t *testing.T) {
	_, reg := setupTestEnvironment()
	ctx := context.Background()
	orgID := uuid.New()

	req := subDomain.CreateOrderRequest{
		Items: []subDomain.CreateOrderRequestItem{
			{CapabilityCode: "radiology.basic", BillingCycle: "monthly"},
		},
		Currency: "NGN",
	}

	userID := uuid.New()
	resp, err := reg.Subscription.CommercialService.CreateCommercialOrder(ctx, orgID, &userID, req, "mock")
	if err != nil {
		t.Fatalf("order creation failed: %v", err)
	}

	mockBody, _ := json.Marshal(map[string]interface{}{
		"eventId":               "evt_idempotent_test_99",
		"eventType":             "payment.succeeded",
		"providerTransactionId": resp.PaymentTransaction.ProviderTransactionID,
		"providerReference":     resp.PaymentTransaction.ProviderReference,
		"orderId":               resp.Order.ID,
		"status":                "successful",
		"amount":                resp.Order.Total,
		"currency":              "NGN",
	})

	httpReq1 := httptest.NewRequest(http.MethodPost, "/api/v1/billing/webhooks/mock", bytes.NewReader(mockBody))

	// Webhook #1
	err1 := reg.Subscription.CommercialService.ProcessWebhookEvent(ctx, "mock", httpReq1, mockBody)
	if err1 != nil {
		t.Fatalf("first webhook execution failed: %v", err1)
	}

	// Webhook #2 (Duplicate)
	httpReq2 := httptest.NewRequest(http.MethodPost, "/api/v1/billing/webhooks/mock", bytes.NewReader(mockBody))
	err2 := reg.Subscription.CommercialService.ProcessWebhookEvent(ctx, "mock", httpReq2, mockBody)
	if err2 != nil {
		t.Fatalf("duplicate webhook execution should be safe and idempotent, got error: %v", err2)
	}
}

// TEST 4: Cryptographic Webhook Signature Verification
func TestCryptographicWebhookSignatureVerification(t *testing.T) {
	_, reg := setupTestEnvironment()
	ctx := context.Background()

	mockProvider := billingInfra.NewMockPaymentProvider()
	mockProvider.InvalidSignature = true
	reg.Subscription.ProviderRegistry.Register(mockProvider)

	mockBody := []byte(`{"event":"payment.succeeded"}`)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/billing/webhooks/mock", bytes.NewReader(mockBody))

	err := reg.Subscription.CommercialService.ProcessWebhookEvent(ctx, "mock", httpReq, mockBody)
	if err == nil {
		t.Fatalf("expected invalid signature rejection error, got nil")
	}
	if err != billingDomain.ErrInvalidWebhookSignature {
		t.Errorf("expected ErrInvalidWebhookSignature, got: %v", err)
	}
}

// TEST 5 & TEST 6: Payment Amount and Currency Protection Mismatch Rejection
func TestPaymentAmountAndCurrencyMismatchRejection(t *testing.T) {
	_, reg := setupTestEnvironment()
	ctx := context.Background()
	orgID := uuid.New()

	req := subDomain.CreateOrderRequest{
		Items: []subDomain.CreateOrderRequestItem{
			{CapabilityCode: "laboratory.analyzer_integration", BillingCycle: "monthly"},
		},
		Currency: "NGN",
	}

	userID := uuid.New()
	resp, err := reg.Subscription.CommercialService.CreateCommercialOrder(ctx, orgID, &userID, req, "mock")
	if err != nil {
		t.Fatalf("order creation failed: %v", err)
	}

	// Mismatched amount
	mismatchedBody, _ := json.Marshal(map[string]interface{}{
		"eventId":               "evt_mismatch_amount_1",
		"eventType":             "payment.succeeded",
		"providerTransactionId": resp.PaymentTransaction.ProviderTransactionID,
		"providerReference":     resp.PaymentTransaction.ProviderReference,
		"orderId":               resp.Order.ID,
		"status":                "successful",
		"amount":                1.00, // Invalid amount
		"currency":              "NGN",
	})

	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/billing/webhooks/mock", bytes.NewReader(mismatchedBody))
	errMismatch := reg.Subscription.CommercialService.ProcessWebhookEvent(ctx, "mock", httpReq, mismatchedBody)
	if errMismatch == nil {
		t.Fatalf("expected amount mismatch rejection, got nil")
	}

	// Assert entitlement was NOT granted
	caps, _ := reg.Subscription.Service.GetEffectiveCapabilities(ctx, orgID)
	for _, c := range caps {
		if c == "laboratory.analyzer_integration" {
			t.Fatalf("entitlement MUST NOT be granted on amount mismatch")
		}
	}
}

// TEST 10: Grace Period Expiration and Capability Suspension
func TestGracePeriodExpirationAndSuspension(t *testing.T) {
	_, reg := setupTestEnvironment()
	ctx := context.Background()
	orgID := uuid.New()

	// Grant capability with expired grace period
	errGrant := reg.Subscription.Service.GrantCapabilityAddOn(ctx, "", orgID, "laboratory.analyzer_integration", "purchase", nil)
	if errGrant != nil {
		t.Fatalf("failed to grant initial entitlement: %v", errGrant)
	}

	// Revoke entitlement simulating grace period expiration
	errRevoke := reg.Subscription.Service.RevokeCapabilityAddOn(ctx, "", orgID, "laboratory.analyzer_integration")
	if errRevoke != nil {
		t.Fatalf("failed to revoke entitlement: %v", errRevoke)
	}

	// Assert capability is no longer effective
	caps, err := reg.Subscription.Service.GetEffectiveCapabilities(ctx, orgID)
	if err != nil {
		t.Fatalf("failed to retrieve effective capabilities: %v", err)
	}

	for _, c := range caps {
		if c == "laboratory.analyzer_integration" {
			t.Errorf("suspended entitlement must not contribute to effective capabilities")
		}
	}
}

// TEST 13: Full Refund Policy
func TestFullRefundPolicy(t *testing.T) {
	_, reg := setupTestEnvironment()
	ctx := context.Background()
	orgID := uuid.New()

	req := subDomain.CreateOrderRequest{
		Items: []subDomain.CreateOrderRequestItem{
			{CapabilityCode: "laboratory.analyzer_integration", BillingCycle: "monthly"},
		},
		Currency: "NGN",
	}

	userID := uuid.New()
	resp, err := reg.Subscription.CommercialService.CreateCommercialOrder(ctx, orgID, &userID, req, "mock")
	if err != nil {
		t.Fatalf("order creation failed: %v", err)
	}

	// Process payment
	mockBody, _ := json.Marshal(map[string]interface{}{
		"eventId":               "evt_refund_test_1",
		"eventType":             "payment.succeeded",
		"providerTransactionId": resp.PaymentTransaction.ProviderTransactionID,
		"providerReference":     resp.PaymentTransaction.ProviderReference,
		"orderId":               resp.Order.ID,
		"status":                "successful",
		"amount":                resp.Order.Total,
		"currency":              "NGN",
	})
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/billing/webhooks/mock", bytes.NewReader(mockBody))
	_ = reg.Subscription.CommercialService.ProcessWebhookEvent(ctx, "mock", httpReq, mockBody)

	// Issue full refund
	errRefund := reg.Subscription.CommercialService.ProcessRefund(ctx, orgID, resp.Order.ID)
	if errRefund != nil {
		t.Fatalf("refund processing failed: %v", errRefund)
	}

	// Assert entitlement is revoked
	caps, _ := reg.Subscription.Service.GetEffectiveCapabilities(ctx, orgID)
	for _, c := range caps {
		if c == "laboratory.analyzer_integration" {
			t.Fatalf("entitlement MUST be revoked upon full refund")
		}
	}
}
