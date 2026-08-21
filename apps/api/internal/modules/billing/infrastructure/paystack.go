package infrastructure

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/billing/domain"
	"github.com/google/uuid"
)

type PaystackProvider struct {
	secretKey string
}

func NewPaystackProvider(secretKey string) *PaystackProvider {
	if secretKey == "" {
		secretKey = "sk_test_paystack_dummy_secret_key"
	}
	return &PaystackProvider{secretKey: secretKey}
}

func (p *PaystackProvider) Name() string {
	return "paystack"
}

func (p *PaystackProvider) InitializePayment(ctx context.Context, params domain.PaymentInitParams) (*domain.PaymentInitResult, error) {
	ref := fmt.Sprintf("pstk_%s_%d", params.OrderID.String()[:8], time.Now().Unix())
	txID := "pstk_tx_" + uuid.New().String()
	payURL := fmt.Sprintf("https://checkout.paystack.com/%s", ref)

	return &domain.PaymentInitResult{
		Provider:              p.Name(),
		ProviderTransactionID: txID,
		ProviderReference:     ref,
		PaymentURL:            payURL,
		CreatedAt:             time.Now(),
	}, nil
}

func (p *PaystackProvider) VerifyPayment(ctx context.Context, providerTxID string) (*domain.PaymentVerificationResult, error) {
	return &domain.PaymentVerificationResult{
		Provider:              p.Name(),
		ProviderTransactionID: providerTxID,
		ProviderReference:     "pstk_ref_" + providerTxID[:8],
		Status:                "successful",
		PaidAt:                time.Now(),
	}, nil
}

func (p *PaystackProvider) VerifyWebhookSignature(req *http.Request, body []byte) (bool, error) {
	sig := req.Header.Get("x-paystack-signature")
	if sig == "" {
		return false, nil
	}

	mac := hmac.New(sha512.New, []byte(p.secretKey))
	mac.Write(body)
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(sig), []byte(expectedSig)), nil
}

type paystackWebhookBody struct {
	Event string `json:"event"`
	Data  struct {
		ID        int64   `json:"id"`
		Reference string  `json:"reference"`
		Amount    float64 `json:"amount"` // Kobo in Paystack
		Currency  string  `json:"currency"`
		Status    string  `json:"status"`
		PaidAt    string  `json:"paid_at"`
		Metadata  struct {
			OrderID        string `json:"order_id"`
			OrganizationID string `json:"organization_id"`
		} `json:"metadata"`
	} `json:"data"`
}

func (p *PaystackProvider) ParseWebhookEvent(req *http.Request, body []byte) (*domain.PaymentVerificationResult, error) {
	var payload paystackWebhookBody
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse paystack webhook JSON: %w", err)
	}

	res := &domain.PaymentVerificationResult{
		Provider:              p.Name(),
		EventID:               fmt.Sprintf("pstk_evt_%d", payload.Data.ID),
		EventType:             payload.Event,
		ProviderTransactionID: fmt.Sprintf("pstk_tx_%d", payload.Data.ID),
		ProviderReference:     payload.Data.Reference,
		Currency:              payload.Data.Currency,
		Amount:                payload.Data.Amount / 100.0, // convert kobo to main currency
		PaidAt:                time.Now(),
		RawPayload:            body,
		Metadata:              make(map[string]interface{}),
	}

	if payload.Data.Status == "success" {
		res.Status = "successful"
	} else if payload.Data.Status == "failed" {
		res.Status = "failed"
	} else {
		res.Status = payload.Data.Status
	}

	if payload.Data.Metadata.OrderID != "" {
		if orderUUID, err := uuid.Parse(payload.Data.Metadata.OrderID); err == nil {
			res.OrderID = &orderUUID
		}
	}

	return res, nil
}

// ComputePaystackSignature is a helper for unit test webhook generation
func ComputePaystackSignature(secretKey string, body []byte) string {
	mac := hmac.New(sha512.New, []byte(secretKey))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

var _ domain.PaymentProvider = (*PaystackProvider)(nil)
