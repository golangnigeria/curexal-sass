package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PatientReferral struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	VisitID           uuid.UUID  `json:"visitId" db:"visit_id"`
	DoctorPartnerID   *uuid.UUID `json:"doctorPartnerId" db:"doctor_partner_id"`
	FacilityPartnerID *uuid.UUID `json:"facilityPartnerId" db:"facility_partner_id"`
	CampaignID        *uuid.UUID `json:"campaignId" db:"campaign_id"`
	MarketingRepID    *uuid.UUID `json:"marketingRepId" db:"marketing_rep_id"`
	CreatedBy         *uuid.UUID `json:"createdBy" db:"created_by"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
}

type InvoicePaymentAllocation struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	InvoiceID            uuid.UUID  `json:"invoiceId" db:"invoice_id"`
	GuarantorType        string     `json:"guarantorType" db:"guarantor_type"`
	PayerPartnerID       *uuid.UUID `json:"payerPartnerId" db:"payer_partner_id"`
	AllocatedAmount      float64    `json:"allocatedAmount" db:"allocated_amount"`
	BillingPolicy        string     `json:"billingPolicy" db:"billing_policy"`
	PreAuthorizationCode string     `json:"preAuthorizationCode" db:"pre_authorization_code"`
	Status               string     `json:"status" db:"status"`
	StatementItemID      *uuid.UUID `json:"statementItemId" db:"statement_item_id"`
	CreatedAt            time.Time  `json:"createdAt" db:"created_at"`
}

type FinancialLedgerEntry struct {
	ID            uuid.UUID `json:"id" db:"id"`
	TenantID      uuid.UUID `json:"tenantId" db:"tenant_id"`
	InvoiceID     uuid.UUID `json:"invoiceId" db:"invoice_id"`
	EntryType     string    `json:"entryType" db:"entry_type"`
	Amount        float64   `json:"amount" db:"amount"`
	PaymentMethod string    `json:"paymentMethod" db:"payment_method"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

type StatementOfAccount struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	CorporatePartnerID uuid.UUID `json:"corporatePartnerId" db:"corporate_partner_id"`
	BillingPeriodStart time.Time `json:"billingPeriodStart" db:"billing_period_start"`
	BillingPeriodEnd   time.Time `json:"billingPeriodEnd" db:"billing_period_end"`
	TotalInvoiced      float64   `json:"totalInvoiced" db:"total_invoiced"`
	Status             string    `json:"status" db:"status"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
}

type StatementInvoiceItem struct {
	ID              uuid.UUID `json:"id" db:"id"`
	StatementID     uuid.UUID `json:"statementId" db:"statement_id"`
	InvoiceID       uuid.UUID `json:"invoiceId" db:"invoice_id"`
	AllocatedAmount float64   `json:"allocatedAmount" db:"allocated_amount"`
	DisputeStatus   string    `json:"disputeStatus" db:"dispute_status"`
	Notes           string    `json:"notes" db:"notes"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}

type BillingPolicy interface {
	Validate(ctx context.Context, allocation *InvoicePaymentAllocation) error
	Collect(ctx context.Context, allocation *InvoicePaymentAllocation) error
	Reverse(ctx context.Context, allocation *InvoicePaymentAllocation) error
}

type CashBillingPolicy struct{}

func (p *CashBillingPolicy) Validate(ctx context.Context, allocation *InvoicePaymentAllocation) error {
	if allocation.AllocatedAmount <= 0 {
		return errors.New("cash billing amount must be greater than zero")
	}
	return nil
}

func (p *CashBillingPolicy) Collect(ctx context.Context, allocation *InvoicePaymentAllocation) error {
	allocation.Status = "SETTLED"
	return nil
}

func (p *CashBillingPolicy) Reverse(ctx context.Context, allocation *InvoicePaymentAllocation) error {
	allocation.Status = "UNPAID"
	return nil
}

type PrepaidBillingPolicy struct {
	WalletBalanceFunc func(walletID uuid.UUID) (float64, error)
	DebitWalletFunc   func(walletID uuid.UUID, amount float64) error
}

func (p *PrepaidBillingPolicy) Validate(ctx context.Context, allocation *InvoicePaymentAllocation) error {
	if allocation.PayerPartnerID == nil {
		return errors.New("prepaid billing requires a valid payer partner ID")
	}
	if allocation.AllocatedAmount <= 0 {
		return errors.New("allocated prepaid amount must be greater than zero")
	}
	return nil
}

func (p *PrepaidBillingPolicy) Collect(ctx context.Context, allocation *InvoicePaymentAllocation) error {
	allocation.Status = "SETTLED"
	return nil
}

func (p *PrepaidBillingPolicy) Reverse(ctx context.Context, allocation *InvoicePaymentAllocation) error {
	allocation.Status = "UNPAID"
	return nil
}

type CreditBillingPolicy struct {
	CreditLimitFunc func(partnerID uuid.UUID) (float64, error)
	ExposureFunc    func(partnerID uuid.UUID) (float64, error)
}

func (p *CreditBillingPolicy) Validate(ctx context.Context, allocation *InvoicePaymentAllocation) error {
	if allocation.PayerPartnerID == nil {
		return errors.New("credit billing requires a valid corporate/HMO partner ID")
	}
	if allocation.AllocatedAmount <= 0 {
		return errors.New("allocated credit amount must be greater than zero")
	}
	if allocation.PreAuthorizationCode == "" {
		return errors.New("credit billing requires a valid pre-authorization approval code")
	}

	if p.CreditLimitFunc != nil && p.ExposureFunc != nil {
		limit, err := p.CreditLimitFunc(*allocation.PayerPartnerID)
		if err != nil {
			return fmt.Errorf("retrieving credit limit: %w", err)
		}
		exposure, err := p.ExposureFunc(*allocation.PayerPartnerID)
		if err != nil {
			return fmt.Errorf("retrieving client exposure: %w", err)
		}

		if exposure+allocation.AllocatedAmount > limit {
			return fmt.Errorf("transaction declined: credit limit exceeded (Limit: ₦%.2f, Current Exposure: ₦%.2f, Attempted: ₦%.2f)", limit, exposure, allocation.AllocatedAmount)
		}
	}

	return nil
}

func (p *CreditBillingPolicy) Collect(ctx context.Context, allocation *InvoicePaymentAllocation) error {
	allocation.Status = "PENDING_APPROVAL"
	return nil
}

func (p *CreditBillingPolicy) Reverse(ctx context.Context, allocation *InvoicePaymentAllocation) error {
	allocation.Status = "UNPAID"
	return nil
}
