package testing

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	billingDomain "github.com/golangnigeria/curexal/internal/modules/billing/domain"
	billingInfra "github.com/golangnigeria/curexal/internal/modules/billing/infrastructure"
	subDomain "github.com/golangnigeria/curexal/internal/modules/subscription/domain"
	"github.com/google/uuid"
)

// PRODUCTION ACCEPTANCE AUDIT TEST SUITE (MILESTONE 4)

func TestMilestone4_FullE2ECommercialLifecycle(t *testing.T) {
	_, reg := setupTestEnvironment()
	ctx := context.Background()
	orgID := uuid.New()

	// 1. Check initial state: organization.plan is "smart"
	plan, err := reg.Subscription.Service.GetOrganizationPlan(ctx, orgID)
	if err != nil {
		t.Fatalf("failed to retrieve plan: %v", err)
	}
	if plan != "smart" {
		t.Fatalf("expected initial plan to be smart, got: %s", plan)
	}

	// 2. Create order for add-on capability (laboratory.analyzer_integration)
	req := subDomain.CreateOrderRequest{
		Items: []subDomain.CreateOrderRequestItem{
			{CapabilityCode: "laboratory.analyzer_integration", BillingCycle: "monthly"},
		},
		Currency: "NGN",
	}
	userID := uuid.New()
	orderResp, err := reg.Subscription.CommercialService.CreateCommercialOrder(ctx, orgID, &userID, req, "mock")
	if err != nil {
		t.Fatalf("failed to create commercial order: %v", err)
	}

	// 3. Confirm BEFORE payment: entitlement is NOT effective
	capsBefore, _ := reg.Subscription.Service.GetEffectiveCapabilities(ctx, orgID)
	for _, c := range capsBefore {
		if c == "laboratory.analyzer_integration" {
			t.Fatalf("CRITICAL SECURITY VIOLATION: capability active BEFORE payment confirmation")
		}
	}

	// 4. Process verified webhook
	mockBody, _ := json.Marshal(map[string]interface{}{
		"eventId":               "evt_e2e_acceptance_101",
		"eventType":             "payment.succeeded",
		"providerTransactionId": orderResp.PaymentTransaction.ProviderTransactionID,
		"providerReference":     orderResp.PaymentTransaction.ProviderReference,
		"orderId":               orderResp.Order.ID,
		"status":                "successful",
		"amount":                orderResp.Order.Total,
		"currency":              "NGN",
	})
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/billing/webhooks/mock", bytes.NewReader(mockBody))
	errWeb := reg.Subscription.CommercialService.ProcessWebhookEvent(ctx, "mock", httpReq, mockBody)
	if errWeb != nil {
		t.Fatalf("webhook processing failed: %v", errWeb)
	}

	// 5. Confirm AFTER payment: entitlement is active with source_type = add_on
	entitlements, err := reg.Subscription.Service.GetOrganizationEntitlements(ctx, orgID)
	if err != nil {
		t.Fatalf("failed to query organization entitlements: %v", err)
	}
	foundActive := false
	for _, e := range entitlements {
		if e.CapabilityCode == "laboratory.analyzer_integration" && e.Status == "active" {
			foundActive = true
			break
		}
	}
	if !foundActive {
		t.Fatalf("expected entitlement laboratory.analyzer_integration to be active after payment")
	}

	// 6. Confirm plan is STILL "smart" (Base Plan Immutability)
	planAfter, _ := reg.Subscription.Service.GetOrganizationPlan(ctx, orgID)
	if planAfter != "smart" {
		t.Fatalf("CRITICAL ARCHITECTURAL VIOLATION: plan mutated from smart to %s", planAfter)
	}

	// 7. Revoke entitlement
	_ = reg.Subscription.Service.RevokeCapabilityAddOn(ctx, "", orgID, "laboratory.analyzer_integration")

	// 8. Confirm AFTER revocation: entitlement is NO LONGER effective
	capsAfterRevoke, _ := reg.Subscription.Service.GetEffectiveCapabilities(ctx, orgID)
	for _, c := range capsAfterRevoke {
		if c == "laboratory.analyzer_integration" {
			t.Fatalf("revoked entitlement must not contribute to effective capabilities")
		}
	}
}

func TestMilestone4_ConcurrentWebhookIdempotency(t *testing.T) {
	_, reg := setupTestEnvironment()
	ctx := context.Background()
	orgID := uuid.New()

	req := subDomain.CreateOrderRequest{
		Items: []subDomain.CreateOrderRequestItem{
			{CapabilityCode: "radiology.pacs_dicom", BillingCycle: "monthly"},
		},
		Currency: "NGN",
	}
	userID := uuid.New()
	orderResp, err := reg.Subscription.CommercialService.CreateCommercialOrder(ctx, orgID, &userID, req, "mock")
	if err != nil {
		t.Fatalf("order creation failed: %v", err)
	}

	mockBody, _ := json.Marshal(map[string]interface{}{
		"eventId":               "evt_concurrent_test_500",
		"eventType":             "payment.succeeded",
		"providerTransactionId": orderResp.PaymentTransaction.ProviderTransactionID,
		"providerReference":     orderResp.PaymentTransaction.ProviderReference,
		"orderId":               orderResp.Order.ID,
		"status":                "successful",
		"amount":                orderResp.Order.Total,
		"currency":              "NGN",
	})

	var wg sync.WaitGroup
	errCh := make(chan error, 10)
	workers := 10

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/billing/webhooks/mock", bytes.NewReader(mockBody))
			err := reg.Subscription.CommercialService.ProcessWebhookEvent(ctx, "mock", httpReq, mockBody)
			if err != nil && err != billingDomain.ErrWebhookAlreadyProcessed {
				errCh <- err
			}
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("concurrent webhook processing emitted error: %v", err)
	}

	// Assert exactly ONE entitlement active
	ents, _ := reg.Subscription.Service.GetOrganizationEntitlements(ctx, orgID)
	activeCount := 0
	for _, e := range ents {
		if e.CapabilityCode == "radiology.pacs_dicom" && e.Status == "active" {
			activeCount++
		}
	}
	if activeCount != 1 {
		t.Fatalf("expected exactly 1 active entitlement from concurrent webhooks, got: %d", activeCount)
	}
}

func TestMilestone4_TenantIsolationProtection(t *testing.T) {
	_, reg := setupTestEnvironment()
	ctx := context.Background()

	orgA := uuid.New()
	orgB := uuid.New()

	req := subDomain.CreateOrderRequest{
		Items: []subDomain.CreateOrderRequestItem{
			{CapabilityCode: "laboratory.analyzer_integration", BillingCycle: "monthly"},
		},
		Currency: "NGN",
	}

	userID := uuid.New()
	orderA, err := reg.Subscription.CommercialService.CreateCommercialOrder(ctx, orgA, &userID, req, "mock")
	if err != nil {
		t.Fatalf("failed to create order for Org A: %v", err)
	}

	// Attempt refund on Org A's order by Org B
	errRefund := reg.Subscription.CommercialService.ProcessRefund(ctx, orgB, orderA.Order.ID)
	if errRefund == nil {
		t.Fatalf("CRITICAL SECURITY VIOLATION: Org B successfully refunded Org A's order")
	}
}

func TestMilestone4_PaystackCryptographicSignature(t *testing.T) {
	paystack := billingInfra.NewPaystackProvider("sk_test_curexal_secret_999")

	body := []byte(`{"event":"charge.success","data":{"id":123,"reference":"ref_123"}}`)
	validSig := billingInfra.ComputePaystackSignature("sk_test_curexal_secret_999", body)

	// Valid signature test
	reqValid := httptest.NewRequest("POST", "/api/v1/billing/webhooks/paystack", bytes.NewReader(body))
	reqValid.Header.Set("x-paystack-signature", validSig)
	ok, err := paystack.VerifyWebhookSignature(reqValid, body)
	if err != nil || !ok {
		t.Fatalf("valid paystack signature verification failed")
	}

	// Invalid signature test
	reqInvalid := httptest.NewRequest("POST", "/api/v1/billing/webhooks/paystack", bytes.NewReader(body))
	reqInvalid.Header.Set("x-paystack-signature", "invalid_signature_hash")
	okInvalid, _ := paystack.VerifyWebhookSignature(reqInvalid, body)
	if okInvalid {
		t.Fatalf("CRITICAL SECURITY VIOLATION: invalid paystack signature accepted")
	}
}
