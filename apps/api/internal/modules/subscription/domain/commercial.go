package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CommercialOrder struct {
	ID             uuid.UUID             `json:"id" db:"id"`
	OrderNumber    string                `json:"orderNumber" db:"order_number"`
	OrganizationID uuid.UUID             `json:"organizationId" db:"organization_id"`
	Status         string                `json:"status" db:"status"` // pending_payment, paid, cancelled, expired, refunded
	Currency       string                `json:"currency" db:"currency"`
	Subtotal       float64               `json:"subtotal" db:"subtotal"`
	Tax            float64               `json:"tax" db:"tax"`
	Discount       float64               `json:"discount" db:"discount"`
	Total          float64               `json:"total" db:"total"`
	CreatedBy      *uuid.UUID            `json:"createdBy,omitempty" db:"created_by"`
	PaidAt         *time.Time            `json:"paidAt,omitempty" db:"paid_at"`
	CancelledAt    *time.Time            `json:"cancelledAt,omitempty" db:"cancelled_at"`
	RefundedAt     *time.Time            `json:"refundedAt,omitempty" db:"refunded_at"`
	CreatedAt      time.Time             `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time             `json:"updatedAt" db:"updated_at"`
	Items          []CommercialOrderItem `json:"items,omitempty"`
}

type CommercialOrderItem struct {
	ID             uuid.UUID `json:"id" db:"id"`
	OrderID        uuid.UUID `json:"orderId" db:"order_id"`
	CapabilityID   uuid.UUID `json:"capabilityId" db:"capability_id"`
	CapabilityCode string    `json:"capabilityCode" db:"capability_code"`
	BillingCycle   string    `json:"billingCycle" db:"billing_cycle"` // monthly, annual
	UnitPrice      float64   `json:"unitPrice" db:"unit_price"`
	Amount         float64   `json:"amount" db:"amount"`
	Currency       string    `json:"currency" db:"currency"`
	Metadata       string    `json:"metadata,omitempty" db:"metadata"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
}

type PaymentTransaction struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	OrderID               uuid.UUID  `json:"orderId" db:"order_id"`
	OrganizationID        uuid.UUID  `json:"organizationId" db:"organization_id"`
	Provider              string     `json:"provider" db:"provider"`
	ProviderTransactionID string     `json:"providerTransactionId" db:"provider_transaction_id"`
	ProviderReference     string     `json:"providerReference" db:"provider_reference"`
	PaymentURL            string     `json:"paymentUrl,omitempty" db:"payment_url"`
	Status                string     `json:"status" db:"status"` // pending, successful, failed, refunded
	Amount                float64    `json:"amount" db:"amount"`
	Currency              string     `json:"currency" db:"currency"`
	Metadata              string     `json:"metadata,omitempty" db:"metadata"`
	PaidAt                *time.Time `json:"paidAt,omitempty" db:"paid_at"`
	CreatedAt             time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time  `json:"updatedAt" db:"updated_at"`
}

type WebhookEventLog struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Provider    string     `json:"provider" db:"provider"`
	EventID     string     `json:"eventId" db:"event_id"`
	EventType   string     `json:"eventType" db:"event_type"`
	Payload     string     `json:"payload" db:"payload"`
	Processed   bool       `json:"processed" db:"processed"`
	ProcessedAt *time.Time `json:"processedAt,omitempty" db:"processed_at"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
}

type CommercialInvoice struct {
	ID             uuid.UUID `json:"id" db:"id"`
	InvoiceNumber  string    `json:"invoiceNumber" db:"invoice_number"`
	OrganizationID uuid.UUID `json:"organizationId" db:"organization_id"`
	OrderID        uuid.UUID `json:"orderId" db:"order_id"`
	Subtotal       float64   `json:"subtotal" db:"subtotal"`
	Tax            float64   `json:"tax" db:"tax"`
	Total          float64   `json:"total" db:"total"`
	Currency       string    `json:"currency" db:"currency"`
	Status         string    `json:"status" db:"status"` // issued, paid, cancelled, refunded
	IssuedAt       time.Time `json:"issuedAt" db:"issued_at"`
	DueAt          time.Time `json:"dueAt" db:"due_at"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
}

type CommercialReceipt struct {
	ID                uuid.UUID `json:"id" db:"id"`
	ReceiptNumber     string    `json:"receiptNumber" db:"receipt_number"`
	PaymentID         uuid.UUID `json:"paymentId" db:"payment_id"`
	OrderID           uuid.UUID `json:"orderId" db:"order_id"`
	OrganizationID    uuid.UUID `json:"organizationId" db:"organization_id"`
	Amount            float64   `json:"amount" db:"amount"`
	Currency          string    `json:"currency" db:"currency"`
	ProviderReference string    `json:"providerReference" db:"provider_reference"`
	PaidAt            time.Time `json:"paidAt" db:"paid_at"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
}

type CreateOrderRequestItem struct {
	CapabilityCode string `json:"capabilityCode"`
	BillingCycle   string `json:"billingCycle"` // monthly, annual
}

type CreateOrderRequest struct {
	Items    []CreateOrderRequestItem `json:"items"`
	Currency string                   `json:"currency"` // NGN, USD
}

type CreateOrderResponse struct {
	Order              CommercialOrder    `json:"order"`
	PaymentTransaction PaymentTransaction `json:"paymentTransaction"`
	PaymentURL         string             `json:"paymentUrl"`
}

type CommercialRepository interface {
	CreateOrder(ctx context.Context, order *CommercialOrder, items []CommercialOrderItem) error
	GetOrder(ctx context.Context, orderID uuid.UUID) (*CommercialOrder, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*CommercialOrder, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string, paidAt *time.Time) error

	CreatePaymentTransaction(ctx context.Context, tx *PaymentTransaction) error
	GetPaymentTransaction(ctx context.Context, txID uuid.UUID) (*PaymentTransaction, error)
	GetPaymentTransactionByProviderID(ctx context.Context, provider, providerTxID string) (*PaymentTransaction, error)
	UpdatePaymentTransactionStatus(ctx context.Context, txID uuid.UUID, status string, paidAt *time.Time) error

	RecordWebhookEvent(ctx context.Context, event *WebhookEventLog) error
	IsWebhookProcessed(ctx context.Context, provider, eventID string) (bool, error)
	MarkWebhookProcessed(ctx context.Context, provider, eventID string) error

	CreateInvoice(ctx context.Context, invoice *CommercialInvoice) error
	GetInvoicesByOrg(ctx context.Context, orgID uuid.UUID) ([]CommercialInvoice, error)

	CreateReceipt(ctx context.Context, receipt *CommercialReceipt) error
	GetReceiptsByOrg(ctx context.Context, orgID uuid.UUID) ([]CommercialReceipt, error)
}
