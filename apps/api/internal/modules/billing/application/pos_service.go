package application

import (
	"context"
	"errors"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/billing/domain"
	"github.com/golangnigeria/curexal/internal/modules/billing/infrastructure/postgres"
)

type POSService struct {
	server *server.Server
	repo   *postgres.POSRepository
}

func NewPOSService(s *server.Server, repo *postgres.POSRepository) *POSService {
	return &POSService{
		server: s,
		repo:   repo,
	}
}

// ListInvoices retrieves clinic invoices
func (s *POSService) ListInvoices(
	ctx context.Context,
	tenantID string,
	status *string,
	patientID *string,
) ([]domain.Invoice, error) {
	return s.repo.ListInvoices(ctx, tenantID, status, patientID)
}

// GetInvoiceByID returns detailed invoice
func (s *POSService) GetInvoiceByID(ctx context.Context, tenantID, invoiceID string) (*domain.Invoice, error) {
	if invoiceID == "" {
		return nil, errors.New("invoice ID is required")
	}
	return s.repo.GetInvoiceByID(ctx, tenantID, invoiceID)
}

// ProcessPayment clears or partially pays an invoice
func (s *POSService) ProcessPayment(
	ctx context.Context,
	tenantID, invoiceID string,
	cashierID *string,
	payload domain.ProcessPaymentPayload,
) (*domain.PaymentReceipt, error) {
	if invoiceID == "" {
		return nil, errors.New("invoice ID is required")
	}
	return s.repo.ProcessPayment(ctx, tenantID, invoiceID, cashierID, payload)
}

// GetPaymentReceipt fetches receipt
func (s *POSService) GetPaymentReceipt(ctx context.Context, tenantID, receiptNumber string) (*domain.PaymentReceipt, error) {
	if receiptNumber == "" {
		return nil, errors.New("receipt number is required")
	}
	return s.repo.GetPaymentByReceipt(ctx, tenantID, receiptNumber)
}
