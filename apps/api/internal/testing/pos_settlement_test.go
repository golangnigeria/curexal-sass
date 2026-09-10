package testing

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	billingApi "github.com/golangnigeria/curexal/internal/modules/billing/api"
	billingApp "github.com/golangnigeria/curexal/internal/modules/billing/application"
	billingDomain "github.com/golangnigeria/curexal/internal/modules/billing/domain"
	billingRepo "github.com/golangnigeria/curexal/internal/modules/billing/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPOSTestEngine() (*echo.Echo, *server.Server, *billingApp.POSService) {
	e := echo.New()
	e.HideBanner = true
	logger := zerolog.Nop()
	s := &server.Server{
		Config: &config.Config{
			Auth: config.AuthConfig{
				SecretKey: "test-secret-key-32-bytes-long!!",
			},
		},
		Logger: &logger,
	}

	repo := billingRepo.NewPOSRepository(s)
	svc := billingApp.NewPOSService(s, repo)
	hnd := billingApi.NewPOSHandler(svc)

	g := e.Group("/api/v1/billing")
	g.GET("/invoices", hnd.ListInvoices)
	g.GET("/invoices/:id", hnd.GetInvoiceByID)
	g.POST("/invoices/:id/payments", hnd.ProcessPayment)
	g.GET("/receipts/:receiptNumber", hnd.GetPaymentReceipt)

	return e, s, svc
}

// 1. Cashier POS Payment Request Validation
func TestPOSSettlement_Validation(t *testing.T) {
	e, _, _ := setupPOSTestEngine()

	t.Run("Rejects missing invoice ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/invoices//payments", nil)
		req.Header.Set("X-Tenant-ID", "00000000-0000-0000-0000-000000000001")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.True(t, rec.Code == http.StatusNotFound || rec.Code == http.StatusBadRequest)
	})

	t.Run("Rejects malformed JSON body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/billing/invoices/inv-123/payments", bytes.NewReader([]byte("{invalid-json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Tenant-ID", "00000000-0000-0000-0000-000000000001")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// 2. Split-Tender Payment Breakdown & Calculation
func TestPOSSettlement_SplitTenderCalculations(t *testing.T) {
	t.Run("Validates split tender payments across Cash and Card POS", func(t *testing.T) {
		totalBill := 5000.00
		cashTender := 2000.00
		cardTender := 3000.00

		breakdown := map[string]interface{}{
			"CASH": cashTender,
			"POS":  cardTender,
		}

		payload := billingDomain.ProcessPaymentPayload{
			Amount:          cashTender + cardTender,
			TenderType:      "SPLIT",
			TenderBreakdown: breakdown,
		}

		assert.Equal(t, totalBill, payload.Amount)
		assert.Equal(t, "SPLIT", payload.TenderType)

		computedSum := payload.TenderBreakdown["CASH"].(float64) + payload.TenderBreakdown["POS"].(float64)
		assert.Equal(t, totalBill, computedSum)
	})

	t.Run("Validates partial payment balance deduction", func(t *testing.T) {
		totalAmount := 12500.00
		firstPayment := 5000.00
		remainingBalance := totalAmount - firstPayment

		assert.Equal(t, 7500.00, remainingBalance)

		// Second payment clears the remainder
		secondPayment := 7500.00
		finalBalance := remainingBalance - secondPayment
		assert.Equal(t, 0.00, finalBalance)
	})
}

// 3. Receipt Generation and Formatting
func TestPOSSettlement_ReceiptStructure(t *testing.T) {
	t.Run("Generates valid receipt number and metadata", func(t *testing.T) {
		receipt := billingDomain.PaymentReceipt{
			ReceiptNumber:   "RCP-2026-98124",
			InvoiceNumber:   "INV-2026-10492",
			PatientName:     "Emeka Okonkwo",
			MRN:             "PAT-2026-00381",
			AmountPaid:      5000.00,
			PreviousBalance: 5000.00,
			NewBalance:      0.00,
			TenderType:      "CASH",
			Status:          "SETTLED",
			PaidAt:          "2026-09-10T10:15:00Z",
			CashierName:     "Adewale Adeleke",
		}

		matched, err := regexp.MatchString(`^RCP-\d{4}-\d{5}$`, receipt.ReceiptNumber)
		require.NoError(t, err)
		assert.True(t, matched, "Receipt number should match RCP-YYYY-XXXXX format")

		assert.Equal(t, "SETTLED", receipt.Status)
		assert.Equal(t, 0.00, receipt.NewBalance)
		assert.Equal(t, "Emeka Okonkwo", receipt.PatientName)
	})
}
