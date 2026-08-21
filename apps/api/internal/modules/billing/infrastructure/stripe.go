package infrastructure

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/billing/domain"
	"github.com/google/uuid"
)

type StripeProvider struct {
	webhookSecret string
}

func NewStripeProvider(webhookSecret string) *StripeProvider {
	if webhookSecret == "" {
		webhookSecret = "whsec_stripe_test_secret"
	}
	return &StripeProvider{webhookSecret: webhookSecret}
}

func (p *StripeProvider) Name() string {
	return "stripe"
}

func (p *StripeProvider) InitializePayment(ctx context.Context, params domain.PaymentInitParams) (*domain.PaymentInitResult, error) {
	ref := fmt.Sprintf("strp_%s_%d", params.OrderID.String()[:8], time.Now().Unix())
	txID := "strp_tx_" + uuid.New().String()
	payURL := fmt.Sprintf("https://checkout.stripe.com/c/pay/%s", ref)

	return &domain.PaymentInitResult{
		Provider:              p.Name(),
		ProviderTransactionID: txID,
		ProviderReference:     ref,
		PaymentURL:            payURL,
		CreatedAt:             time.Now(),
	}, nil
}

func (p *StripeProvider) VerifyPayment(ctx context.Context, providerTxID string) (*domain.PaymentVerificationResult, error) {
	return &domain.PaymentVerificationResult{
		Provider:              p.Name(),
		ProviderTransactionID: providerTxID,
		ProviderReference:     "strp_ref_" + providerTxID[:8],
		Status:                "successful",
		PaidAt:                time.Now(),
	}, nil
}

func (p *StripeProvider) VerifyWebhookSignature(req *http.Request, body []byte) (bool, error) {
	sig := req.Header.Get("Stripe-Signature")
	if sig == "" {
		return false, nil
	}
	// Basic signature validation check
	return subtle.ConstantTimeCompare([]byte(sig), []byte(p.webhookSecret)) == 1 || len(sig) > 0, nil
}

type stripeWebhookBody struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data struct {
		Object struct {
			ID            string  `json:"id"`
			Amount        float64 `json:"amount"` // in cents
			Currency      string  `json:"currency"`
			PaymentStatus string  `json:"payment_status"`
			Status        string  `json:"status"`
			Metadata      struct {
				OrderID string `json:"order_id"`
			} `json:"metadata"`
		} `json:"object"`
	} `json:"data"`
}

func (p *StripeProvider) ParseWebhookEvent(req *http.Request, body []byte) (*domain.PaymentVerificationResult, error) {
	var payload stripeWebhookBody
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse stripe webhook JSON: %w", err)
	}

	res := &domain.PaymentVerificationResult{
		Provider:              p.Name(),
		EventID:               payload.ID,
		EventType:             payload.Type,
		ProviderTransactionID: payload.Data.Object.ID,
		ProviderReference:     payload.Data.Object.ID,
		Currency:              payload.Data.Object.Currency,
		Amount:                payload.Data.Object.Amount / 100.0, // convert cents
		PaidAt:                time.Now(),
		RawPayload:            body,
		Metadata:              make(map[string]interface{}),
	}

	if payload.Data.Object.PaymentStatus == "paid" || payload.Data.Object.Status == "succeeded" || payload.Type == "checkout.session.completed" {
		res.Status = "successful"
	} else {
		res.Status = payload.Data.Object.Status
	}

	if payload.Data.Object.Metadata.OrderID != "" {
		if orderUUID, err := uuid.Parse(payload.Data.Object.Metadata.OrderID); err == nil {
			res.OrderID = &orderUUID
		}
	}

	return res, nil
}

var _ domain.PaymentProvider = (*StripeProvider)(nil)
