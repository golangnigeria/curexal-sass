package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/billing/domain"
	"github.com/google/uuid"
)

type MockPaymentProvider struct {
	mu                 sync.Mutex
	InvalidSignature   bool
	SimulateFailure    bool
	RecordedInitParams []domain.PaymentInitParams
	CustomVerifyResult *domain.PaymentVerificationResult
}

func NewMockPaymentProvider() *MockPaymentProvider {
	return &MockPaymentProvider{
		RecordedInitParams: make([]domain.PaymentInitParams, 0),
	}
}

func (m *MockPaymentProvider) Name() string {
	return "mock"
}

func (m *MockPaymentProvider) InitializePayment(ctx context.Context, params domain.PaymentInitParams) (*domain.PaymentInitResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.SimulateFailure {
		return nil, fmt.Errorf("mock payment provider init simulated error")
	}

	m.RecordedInitParams = append(m.RecordedInitParams, params)
	txID := "mock_tx_" + uuid.New().String()
	ref := "mock_ref_" + params.OrderID.String()[:8]

	return &domain.PaymentInitResult{
		Provider:              m.Name(),
		ProviderTransactionID: txID,
		ProviderReference:     ref,
		PaymentURL:            fmt.Sprintf("https://mock-checkout.curexal.com/pay/%s", txID),
		CreatedAt:             time.Now(),
	}, nil
}

func (m *MockPaymentProvider) VerifyPayment(ctx context.Context, providerTxID string) (*domain.PaymentVerificationResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.CustomVerifyResult != nil {
		return m.CustomVerifyResult, nil
	}

	return &domain.PaymentVerificationResult{
		Provider:              m.Name(),
		ProviderTransactionID: providerTxID,
		ProviderReference:     "mock_ref_" + providerTxID[:8],
		Status:                "successful",
		PaidAt:                time.Now(),
	}, nil
}

func (m *MockPaymentProvider) VerifyWebhookSignature(req *http.Request, body []byte) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.InvalidSignature {
		return false, nil
	}
	return true, nil
}

type mockWebhookPayload struct {
	EventID               string     `json:"eventId"`
	EventType             string     `json:"eventType"`
	ProviderTransactionID string     `json:"providerTransactionId"`
	ProviderReference     string     `json:"providerReference"`
	OrderID               *uuid.UUID `json:"orderId"`
	Status                string     `json:"status"`
	Amount                float64    `json:"amount"`
	Currency              string     `json:"currency"`
}

func (m *MockPaymentProvider) ParseWebhookEvent(req *http.Request, body []byte) (*domain.PaymentVerificationResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.CustomVerifyResult != nil {
		return m.CustomVerifyResult, nil
	}

	var payload mockWebhookPayload
	if err := json.Unmarshal(body, &payload); err == nil && payload.ProviderTransactionID != "" {
		return &domain.PaymentVerificationResult{
			Provider:              m.Name(),
			EventID:               payload.EventID,
			EventType:             payload.EventType,
			ProviderTransactionID: payload.ProviderTransactionID,
			ProviderReference:     payload.ProviderReference,
			OrderID:               payload.OrderID,
			Status:                payload.Status,
			Amount:                payload.Amount,
			Currency:              payload.Currency,
			PaidAt:                time.Now(),
			RawPayload:            body,
		}, nil
	}

	return &domain.PaymentVerificationResult{
		Provider:              m.Name(),
		EventID:               "mock_evt_1001",
		EventType:             "payment.succeeded",
		ProviderTransactionID: "mock_tx_1001",
		ProviderReference:     "mock_ref_1001",
		Status:                "successful",
		Amount:                35000.00,
		Currency:              "NGN",
		PaidAt:                time.Now(),
		RawPayload:            body,
	}, nil
}

var _ domain.PaymentProvider = (*MockPaymentProvider)(nil)
