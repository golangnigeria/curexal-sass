package postgres

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/billing/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type POSRepository struct {
	server *server.Server
}

func NewPOSRepository(s *server.Server) *POSRepository {
	return &POSRepository{server: s}
}

// ListInvoices retrieves invoices for a tenant with optional status and patient filters
func (r *POSRepository) ListInvoices(
	ctx context.Context,
	tenantID string,
	status *string,
	patientID *string,
) ([]domain.Invoice, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT i.id, i.tenant_id, i.patient_id, i.encounter_id, i.invoice_number,
		       i.subtotal, i.discount_amount, i.tax_amount, i.total_amount,
		       i.amount_paid, i.balance_due, i.status, i.created_at, i.updated_at,
		       COALESCE(p.first_name || ' ' || p.last_name, 'Patient') AS patient_name,
		       COALESCE(p.mrn, '') AS mrn
		FROM billing.invoices i
		LEFT JOIN patient.patients p ON p.id = i.patient_id
		WHERE i.tenant_id = $1
		  AND ($2::text IS NULL OR i.status = $2)
		  AND ($3::uuid IS NULL OR i.patient_id = $3::uuid)
		ORDER BY i.created_at DESC
		LIMIT 100
	`
	rows, err := db.Query(ctx, query, tenantID, status, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []domain.Invoice
	for rows.Next() {
		var inv domain.Invoice
		err := rows.Scan(
			&inv.ID, &inv.TenantID, &inv.PatientID, &inv.EncounterID, &inv.InvoiceNumber,
			&inv.Subtotal, &inv.DiscountAmount, &inv.TaxAmount, &inv.TotalAmount,
			&inv.AmountPaid, &inv.BalanceDue, &inv.Status, &inv.CreatedAt, &inv.UpdatedAt,
			&inv.PatientName, &inv.MRN,
		)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, inv)
	}

	return invoices, nil
}

// GetInvoiceByID retrieves full invoice details with line items and payments
func (r *POSRepository) GetInvoiceByID(ctx context.Context, tenantID, invoiceID string) (*domain.Invoice, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT i.id, i.tenant_id, i.patient_id, i.encounter_id, i.invoice_number,
		       i.subtotal, i.discount_amount, i.tax_amount, i.total_amount,
		       i.amount_paid, i.balance_due, i.status, i.created_at, i.updated_at,
		       COALESCE(p.first_name || ' ' || p.last_name, 'Patient') AS patient_name,
		       COALESCE(p.mrn, '') AS mrn
		FROM billing.invoices i
		LEFT JOIN patient.patients p ON p.id = i.patient_id
		WHERE i.tenant_id = $1 AND i.id = $2
	`
	var inv domain.Invoice
	err := db.QueryRow(ctx, query, tenantID, invoiceID).Scan(
		&inv.ID, &inv.TenantID, &inv.PatientID, &inv.EncounterID, &inv.InvoiceNumber,
		&inv.Subtotal, &inv.DiscountAmount, &inv.TaxAmount, &inv.TotalAmount,
		&inv.AmountPaid, &inv.BalanceDue, &inv.Status, &inv.CreatedAt, &inv.UpdatedAt,
		&inv.PatientName, &inv.MRN,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// 1. Load Line Items
	itemQuery := `
		SELECT id, invoice_id, description, category, unit_price, quantity, total_price, created_at
		FROM billing.invoice_items
		WHERE invoice_id = $1
		ORDER BY created_at ASC
	`
	iRows, err := db.Query(ctx, itemQuery, invoiceID)
	if err == nil {
		for iRows.Next() {
			var itm domain.InvoiceItem
			_ = iRows.Scan(
				&itm.ID, &itm.InvoiceID, &itm.Description, &itm.Category,
				&itm.UnitPrice, &itm.Quantity, &itm.TotalPrice, &itm.CreatedAt,
			)
			inv.Items = append(inv.Items, itm)
		}
		iRows.Close()
	}

	// 2. Load Payments
	payQuery := `
		SELECT p.id, p.tenant_id, p.invoice_id, p.receipt_number, p.amount,
		       p.tender_type, p.tender_breakdown, p.cashier_id, p.status, p.paid_at,
		       COALESCE(u.display_name, 'Cashier') AS cashier_name
		FROM billing.payments p
		LEFT JOIN identity.users u ON u.id = p.cashier_id
		WHERE p.invoice_id = $1
		ORDER BY p.paid_at DESC
	`
	pRows, err := db.Query(ctx, payQuery, invoiceID)
	if err == nil {
		for pRows.Next() {
			var pay domain.Payment
			var breakdownRaw []byte
			_ = pRows.Scan(
				&pay.ID, &pay.TenantID, &pay.InvoiceID, &pay.ReceiptNumber, &pay.Amount,
				&pay.TenderType, &breakdownRaw, &pay.CashierID, &pay.Status, &pay.PaidAt,
				&pay.CashierName,
			)
			if len(breakdownRaw) > 0 {
				_ = json.Unmarshal(breakdownRaw, &pay.TenderBreakdown)
			}
			inv.Payments = append(inv.Payments, pay)
		}
		pRows.Close()
	}

	return &inv, nil
}

// ProcessPayment processes POS payment with split tender support and generates receipt
func (r *POSRepository) ProcessPayment(
	ctx context.Context,
	tenantID, invoiceID string,
	cashierID *string,
	payload domain.ProcessPaymentPayload,
) (*domain.PaymentReceipt, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}

	if payload.Amount <= 0 {
		return nil, errors.New("payment amount must be greater than zero")
	}

	var receipt domain.PaymentReceipt
	var encID *string
	var patientID string

	err := r.server.DB.RunInTx(ctx, func(txCtx context.Context) error {
		db := r.server.DB.Conn(txCtx)

		// 1. Fetch and lock invoice
		query := `
			SELECT i.id, i.patient_id, i.encounter_id, i.invoice_number,
			       i.total_amount, i.amount_paid, i.balance_due, i.status,
			       COALESCE(p.first_name || ' ' || p.last_name, 'Patient') AS patient_name,
			       COALESCE(p.mrn, '') AS mrn
			FROM billing.invoices i
			LEFT JOIN patient.patients p ON p.id = i.patient_id
			WHERE i.tenant_id = $1 AND i.id = $2
			FOR UPDATE OF i
		`
		var (
			invID, invNum, invStatus, patName, mrn string
			totalAmt, amtPaid, balDue               float64
		)
		err := db.QueryRow(txCtx, query, tenantID, invoiceID).Scan(
			&invID, &patientID, &encID, &invNum,
			&totalAmt, &amtPaid, &balDue, &invStatus,
			&patName, &mrn,
		)
		if err != nil {
			return fmt.Errorf("invoice not found: %w", err)
		}

		if invStatus == "PAID" || balDue <= 0.001 {
			return errors.New("invoice is already fully paid")
		}

		if payload.Amount > balDue+0.01 {
			return fmt.Errorf("payment amount (%.2f) exceeds balance due (%.2f)", payload.Amount, balDue)
		}

		// 2. Generate unique Receipt Number
		n, _ := rand.Int(rand.Reader, big.NewInt(90000))
		receiptNum := fmt.Sprintf("RCP-%d-%05d", time.Now().Year(), n.Int64()+10000)

		breakdownJSON, _ := json.Marshal(payload.TenderBreakdown)
		if len(payload.TenderBreakdown) == 0 {
			breakdownJSON = []byte("{}")
		}

		paymentID := uuid.New().String()
		payQuery := `
			INSERT INTO billing.payments (
				id, tenant_id, invoice_id, receipt_number, amount,
				tender_type, tender_breakdown, cashier_id, status, paid_at
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, 'SETTLED', NOW()
			)
		`
		_, err = db.Exec(
			txCtx, payQuery,
			paymentID, tenantID, invoiceID, receiptNum, payload.Amount,
			payload.TenderType, breakdownJSON, cashierID,
		)
		if err != nil {
			return fmt.Errorf("failed to record payment: %w", err)
		}

		// 3. Update invoice balance
		newAmtPaid := amtPaid + payload.Amount
		newBalDue := totalAmt - newAmtPaid
		newStatus := "PARTIALLY_PAID"
		if newBalDue <= 0.009 {
			newBalDue = 0.00
			newStatus = "PAID"
		}

		updQuery := `
			UPDATE billing.invoices
			SET amount_paid = $1, balance_due = $2, status = $3, updated_at = NOW()
			WHERE id = $4
		`
		_, err = db.Exec(txCtx, updQuery, newAmtPaid, newBalDue, newStatus, invoiceID)
		if err != nil {
			return fmt.Errorf("failed to update invoice balances: %w", err)
		}

		// 4. If fully paid and encounter linked, conclude Care Journey SETTLEMENT milestone
		if newStatus == "PAID" {
			msQuery := `
				UPDATE orchestration.care_journey_milestones
				SET status = 'COMPLETED', completed_at = NOW(), description = 'Invoice settled at cashier POS'
				WHERE tenant_id = $1 AND patient_id = $2 AND stage_code = 'SETTLEMENT'
			`
			_, _ = db.Exec(txCtx, msQuery, tenantID, patientID)
		}

		// 5. Construct receipt
		receipt = domain.PaymentReceipt{
			ReceiptNumber:   receiptNum,
			InvoiceNumber:   invNum,
			PatientName:     patName,
			MRN:             mrn,
			AmountPaid:      payload.Amount,
			PreviousBalance: balDue,
			NewBalance:      newBalDue,
			TenderType:      payload.TenderType,
			TenderBreakdown: payload.TenderBreakdown,
			Status:          "SETTLED",
			PaidAt:          time.Now().UTC().Format(time.RFC3339),
			CashierName:     "Cashier",
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &receipt, nil
}

// GetPaymentByReceipt retrieves payment receipt details
func (r *POSRepository) GetPaymentByReceipt(ctx context.Context, tenantID, receiptNumber string) (*domain.PaymentReceipt, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT p.receipt_number, i.invoice_number,
		       COALESCE(pt.first_name || ' ' || pt.last_name, 'Patient') AS patient_name,
		       COALESCE(pt.mrn, '') AS mrn,
		       p.amount, i.balance_due, p.tender_type, p.tender_breakdown,
		       p.status, p.paid_at,
		       COALESCE(u.display_name, 'Cashier') AS cashier_name
		FROM billing.payments p
		JOIN billing.invoices i ON i.id = p.invoice_id
		JOIN patient.patients pt ON pt.id = i.patient_id
		LEFT JOIN identity.users u ON u.id = p.cashier_id
		WHERE p.tenant_id = $1 AND p.receipt_number = $2
	`
	var (
		rec             domain.PaymentReceipt
		breakdownRaw    []byte
		paidAt          time.Time
		previousBalance float64
	)

	err := db.QueryRow(ctx, query, tenantID, receiptNumber).Scan(
		&rec.ReceiptNumber, &rec.InvoiceNumber, &rec.PatientName, &rec.MRN,
		&rec.AmountPaid, &rec.NewBalance, &rec.TenderType, &breakdownRaw,
		&rec.Status, &paidAt, &rec.CashierName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if len(breakdownRaw) > 0 {
		_ = json.Unmarshal(breakdownRaw, &rec.TenderBreakdown)
	}
	rec.PaidAt = paidAt.UTC().Format(time.RFC3339)
	rec.PreviousBalance = previousBalance

	return &rec, nil
}
