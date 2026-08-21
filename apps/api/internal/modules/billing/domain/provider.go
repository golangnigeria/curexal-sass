package domain

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProviderNotFound        = errors.New("payment provider not found")
	ErrInvalidWebhookSignature = errors.New("invalid webhook signature")
	ErrWebhookAlreadyProcessed = errors.New("webhook event already processed")
	ErrPaymentAmountMismatch   = errors.New("payment transaction amount mismatch")
	ErrPaymentCurrencyMismatch = errors.New("payment transaction currency mismatch")
)

type PaymentInitParams struct {
	OrderID        uuid.UUID              `json:"orderId"`
	OrganizationID uuid.UUID              `json:"organizationId"`
	CustomerEmail  string                 `json:"customerEmail"`
	CustomerName   string                 `json:"customerName"`
	Amount         float64                `json:"amount"`
	Currency       string                 `json:"currency"`
	CallbackURL    string                 `json:"callbackUrl"`
	Metadata       map[string]interface{} `json:"metadata"`
}

type PaymentInitResult struct {
	Provider              string    `json:"provider"`
	ProviderTransactionID string    `json:"providerTransactionId"`
	ProviderReference     string    `json:"providerReference"`
	PaymentURL            string    `json:"paymentUrl"`
	CreatedAt             time.Time `json:"createdAt"`
}

type PaymentVerificationResult struct {
	Provider              string                 `json:"provider"`
	EventID               string                 `json:"eventId"`
	EventType             string                 `json:"eventType"`
	ProviderTransactionID string                 `json:"providerTransactionId"`
	ProviderReference     string                 `json:"providerReference"`
	OrderID               *uuid.UUID             `json:"orderId,omitempty"`
	Status                string                 `json:"status"` // successful, failed, pending, refunded
	Amount                float64                `json:"amount"`
	Currency              string                 `json:"currency"`
	PaidAt                time.Time              `json:"paidAt"`
	RawPayload            []byte                 `json:"-"`
	Metadata              map[string]interface{} `json:"metadata,omitempty"`
}

type PaymentProvider interface {
	Name() string
	InitializePayment(ctx context.Context, params PaymentInitParams) (*PaymentInitResult, error)
	VerifyPayment(ctx context.Context, providerTxID string) (*PaymentVerificationResult, error)
	VerifyWebhookSignature(req *http.Request, body []byte) (bool, error)
	ParseWebhookEvent(req *http.Request, body []byte) (*PaymentVerificationResult, error)
}

type PaymentProviderRegistry struct {
	mu        sync.RWMutex
	providers map[string]PaymentProvider
}

func NewPaymentProviderRegistry() *PaymentProviderRegistry {
	return &PaymentProviderRegistry{
		providers: make(map[string]PaymentProvider),
	}
}

func (r *PaymentProviderRegistry) Register(p PaymentProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p.Name()] = p
}

func (r *PaymentProviderRegistry) Get(name string) (PaymentProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, exists := r.providers[name]
	if !exists {
		return nil, ErrProviderNotFound
	}
	return p, nil
}
