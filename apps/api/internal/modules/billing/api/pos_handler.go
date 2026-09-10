package api

import (
	"net/http"
	"strings"

	"github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/modules/billing/application"
	"github.com/golangnigeria/curexal/internal/modules/billing/domain"
	"github.com/labstack/echo/v4"
)

type POSHandler struct {
	service *application.POSService
}

func NewPOSHandler(s *application.POSService) *POSHandler {
	return &POSHandler{service: s}
}

func resolveTenantID(c echo.Context) string {
	principal := auth.GetPrincipal(c)
	if principal != nil && principal.TenantID != "" {
		return principal.TenantID
	}
	tid := c.Request().Header.Get("X-Tenant-ID")
	if tid == "" {
		tid = c.Request().Header.Get("X-Branch-ID")
	}
	if tid == "" {
		tid = c.QueryParam("organization_id")
	}
	return tid
}

// ListInvoices lists clinic patient invoices
func (h *POSHandler) ListInvoices(c echo.Context) error {
	tenantID := resolveTenantID(c)

	var statusPtr *string
	if s := c.QueryParam("status"); s != "" {
		statusPtr = &s
	}

	var patientIDPtr *string
	if pt := c.QueryParam("patient_id"); pt != "" {
		patientIDPtr = &pt
	}

	invoices, err := h.service.ListInvoices(c.Request().Context(), tenantID, statusPtr, patientIDPtr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    invoices,
		"count":   len(invoices),
	})
}

// GetInvoiceByID returns a single invoice with line items and previous payments
func (h *POSHandler) GetInvoiceByID(c echo.Context) error {
	tenantID := resolveTenantID(c)
	invoiceID := c.Param("id")
	if strings.TrimSpace(invoiceID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invoice ID is required")
	}

	inv, err := h.service.GetInvoiceByID(c.Request().Context(), tenantID, invoiceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if inv == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invoice not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    inv,
	})
}

// ProcessPayment processes cashier POS payment and returns receipt
func (h *POSHandler) ProcessPayment(c echo.Context) error {
	tenantID := resolveTenantID(c)
	invoiceID := c.Param("id")
	if strings.TrimSpace(invoiceID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invoice ID is required")
	}

	var payload domain.ProcessPaymentPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payment payload")
	}

	var cashierID *string
	principal := auth.GetPrincipal(c)
	if principal != nil && principal.UserID != "" {
		cashierID = &principal.UserID
	}

	receipt, err := h.service.ProcessPayment(c.Request().Context(), tenantID, invoiceID, cashierID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Payment recorded successfully",
		"data":    receipt,
	})
}

// GetPaymentReceipt retrieves a receipt by receipt number
func (h *POSHandler) GetPaymentReceipt(c echo.Context) error {
	tenantID := resolveTenantID(c)
	receiptNum := c.Param("receiptNumber")
	if strings.TrimSpace(receiptNum) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Receipt number is required")
	}

	receipt, err := h.service.GetPaymentReceipt(c.Request().Context(), tenantID, receiptNum)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if receipt == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Receipt not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    receipt,
	})
}
