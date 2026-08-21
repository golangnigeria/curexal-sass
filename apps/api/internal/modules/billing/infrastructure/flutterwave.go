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

type FlutterwaveProvider struct {
	secretHash string
}

func NewFlutterwaveProvider(secretHash string) *FlutterwaveProvider {
	if secretHash == "" {
		secretHash = "flw_test_secret_hash_value"
	}
	return &FlutterwaveProvider{secretHash: secretHash}
}

func (p *FlutterwaveProvider) Name() string {
	return "flutterwave"
}

func (p *FlutterwaveProvider) InitializePayment(ctx context.Context, params domain.PaymentInitParams) (*domain.PaymentInitResult, error) {
	ref := fmt.Sprintf("flw_%s_%d", params.OrderID.String()[:8], time.Now().Unix())
	txID := "flw_tx_" + uuid.New().String()
	payURL := fmt.Sprintf("https://checkout.flutterwave.com/v3/hosted/pay/%s", ref)

	return &domain.PaymentInitResult{
		Provider:              p.Name(),
		ProviderTransactionID: txID,
		ProviderReference:     ref,
		PaymentURL:            payURL,
		CreatedAt:             time.Now(),
	}, nil
}

func (p *FlutterwaveProvider) VerifyPayment(ctx context.Context, providerTxID string) (*domain.PaymentVerificationResult, error) {
	return &domain.PaymentVerificationResult{
		Provider:              p.Name(),
		ProviderTransactionID: providerTxID,
		ProviderReference:     "flw_ref_" + providerTxID[:8],
		Status:                "successful",
		PaidAt:                time.Now(),
	}, nil
}

func (p *FlutterwaveProvider) VerifyWebhookSignature(req *http.Request, body []byte) (bool, error) {
	hash := req.Header.Get("verif-hash")
	if hash == "" {
		return false, nil
	}

	return subtle.ConstantTimeCompare([]byte(hash), []byte(p.secretHash)) == 1, nil
}

type flutterwaveWebhookBody struct {
	Event string `json:"event"`
	Data  struct {
		ID        int64   `json:"id"`
		TxRef     string  `json:"tx_ref"`
		FlwRef    string  `json:"flw_ref"`
		Amount    float64 `json:"amount"`
		Currency  string  `json:"currency"`
		Status    string  `json:"status"`
		CreatedAt string  `json:"created_at"`
		Meta      struct {
			OrderID string `json:"order_id"`
		} `json:"meta"`
	} `json:"data"`
}

func (p *FlutterwaveProvider) ParseWebhookEvent(req *http.Request, body []byte) (*domain.PaymentVerificationResult, error) {
	var payload flutterwaveWebhookBody
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse flutterwave webhook JSON: %w", err)
	}

	res := &domain.PaymentVerificationResult{
		Provider:              p.Name(),
		EventID:               fmt.Sprintf("flw_evt_%d", payload.Data.ID),
		EventType:             payload.Event,
		ProviderTransactionID: fmt.Sprintf("flw_tx_%d", payload.Data.ID),
		ProviderReference:     payload.Data.TxRef,
		Currency:              payload.Data.Currency,
		Amount:                payload.Data.Amount,
		PaidAt:                time.Now(),
		RawPayload:            body,
		Metadata:              make(map[string]interface{}),
	}

	if payload.Data.Status == "successful" {
		res.Status = "successful"
	} else if payload.Data.Status == "failed" {
		res.Status = "failed"
	} else {
		res.Status = payload.Data.Status
	}

	if payload.Data.Meta.OrderID != "" {
		if orderUUID, err := uuid.Parse(payload.Data.Meta.OrderID); err == nil {
			res.OrderID = &orderUUID
		}
	}

	return res, nil
}

var _ domain.PaymentProvider = (*FlutterwaveProvider)(nil)
