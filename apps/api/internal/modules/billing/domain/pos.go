package domain

import "time"

// Invoice represents a patient bill for clinic encounters, medications, or services
type Invoice struct {
	ID             string        `json:"id" db:"id"`
	TenantID       string        `json:"tenantId" db:"tenant_id"`
	PatientID      string        `json:"patientId" db:"patient_id"`
	EncounterID    *string       `json:"encounterId,omitempty" db:"encounter_id"`
	InvoiceNumber  string        `json:"invoiceNumber" db:"invoice_number"`
	Subtotal       float64       `json:"subtotal" db:"subtotal"`
	DiscountAmount float64       `json:"discountAmount" db:"discount_amount"`
	TaxAmount      float64       `json:"taxAmount" db:"tax_amount"`
	TotalAmount    float64       `json:"totalAmount" db:"total_amount"`
	AmountPaid     float64       `json:"amountPaid" db:"amount_paid"`
	BalanceDue     float64       `json:"balanceDue" db:"balance_due"`
	Status         string        `json:"status" db:"status"` // UNPAID, PARTIALLY_PAID, PAID, VOID
	CreatedAt      time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time     `json:"updatedAt" db:"updated_at"`

	// Enriched fields for cashier UI
	PatientName *string       `json:"patientName,omitempty"`
	MRN         *string       `json:"mrn,omitempty"`
	Items       []InvoiceItem `json:"items,omitempty"`
	Payments    []Payment     `json:"payments,omitempty"`
}

// InvoiceItem represents a line item charge on an invoice
type InvoiceItem struct {
	ID          string    `json:"id" db:"id"`
	InvoiceID   string    `json:"invoiceId" db:"invoice_id"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"` // CONSULTATION, PHARMACY, LAB, PROCEDURE, OTHER
	UnitPrice   float64   `json:"unitPrice" db:"unit_price"`
	Quantity    int       `json:"quantity" db:"quantity"`
	TotalPrice  float64   `json:"totalPrice" db:"total_price"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

// Payment represents a monetary settlement against an invoice
type Payment struct {
	ID              string                 `json:"id" db:"id"`
	TenantID        string                 `json:"tenantId" db:"tenant_id"`
	InvoiceID       string                 `json:"invoiceId" db:"invoice_id"`
	ReceiptNumber   string                 `json:"receiptNumber" db:"receipt_number"`
	Amount          float64                `json:"amount" db:"amount"`
	TenderType      string                 `json:"tenderType" db:"tender_type"` // CASH, POS, BANK_TRANSFER, SPLIT
	TenderBreakdown map[string]interface{} `json:"tenderBreakdown" db:"tender_breakdown"`
	CashierID       *string                `json:"cashierId,omitempty" db:"cashier_id"`
	CashierName     *string                `json:"cashierName,omitempty"`
	Status          string                 `json:"status" db:"status"` // SETTLED, FAILED, REFUNDED
	PaidAt          time.Time              `json:"paidAt" db:"paid_at"`
}

// ProcessPaymentPayload payload submitted by cashier at POS terminal
type ProcessPaymentPayload struct {
	Amount          float64                `json:"amount" validate:"required,gt=0"`
	TenderType      string                 `json:"tenderType" validate:"required"` // CASH, POS, BANK_TRANSFER, SPLIT
	TenderBreakdown map[string]interface{} `json:"tenderBreakdown"`
	Notes           *string                `json:"notes"`
}

// PaymentReceipt generated upon successful settlement
type PaymentReceipt struct {
	ReceiptNumber   string                 `json:"receiptNumber"`
	InvoiceNumber   string                 `json:"invoiceNumber"`
	PatientName     string                 `json:"patientName"`
	MRN             string                 `json:"mrn"`
	AmountPaid      float64                `json:"amountPaid"`
	PreviousBalance float64                `json:"previousBalance"`
	NewBalance      float64                `json:"newBalance"`
	TenderType      string                 `json:"tenderType"`
	TenderBreakdown map[string]interface{} `json:"tenderBreakdown,omitempty"`
	Status          string                 `json:"status"`
	PaidAt          string                 `json:"paidAt"`
	CashierName     string                 `json:"cashierName"`
}
