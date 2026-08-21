package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	billingApp "github.com/golangnigeria/curexal/internal/modules/billing/application"
	billingDomain "github.com/golangnigeria/curexal/internal/modules/billing/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PatientVisitHandler struct {
	server          *server.Server
	billingRegistry *billingApp.BillingPolicyRegistry
}

func NewPatientVisitHandler(s *server.Server) *PatientVisitHandler {
	return &PatientVisitHandler{
		server:          s,
		billingRegistry: billingApp.NewBillingPolicyRegistry(),
	}
}

func (h *PatientVisitHandler) RegisterPatientVisit(c echo.Context) error {
	var payload VisitRegistrationPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload mapping"})
	}

	if payload.PatientDemographics.FirstName == "" || payload.PatientDemographics.LastName == "" || payload.PatientDemographics.Phone == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing mandatory demographics"})
	}

	if len(payload.Investigations) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "investigation basket cannot be empty"})
	}

	ctx := c.Request().Context()
	for _, allocPayload := range payload.PaymentAllocations {
		policy, err := h.billingRegistry.Resolve(allocPayload.BillingPolicy)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		allocModel := &billingDomain.InvoicePaymentAllocation{
			GuarantorType:        allocPayload.GuarantorType,
			PayerPartnerID:       allocPayload.PayerPartnerID,
			AllocatedAmount:      allocPayload.AllocatedAmount,
			BillingPolicy:        allocPayload.BillingPolicy,
			PreAuthorizationCode: allocPayload.PreAuthorizationCode,
		}

		if err := policy.Validate(ctx, allocModel); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		}
	}

	for _, allocPayload := range payload.PaymentAllocations {
		policy, _ := h.billingRegistry.Resolve(allocPayload.BillingPolicy)
		allocModel := &billingDomain.InvoicePaymentAllocation{
			GuarantorType:        allocPayload.GuarantorType,
			PayerPartnerID:       allocPayload.PayerPartnerID,
			AllocatedAmount:      allocPayload.AllocatedAmount,
			BillingPolicy:        allocPayload.BillingPolicy,
			PreAuthorizationCode: allocPayload.PreAuthorizationCode,
		}
		_ = policy.Collect(ctx, allocModel)
	}

	db := h.server.DB
	var patientID uuid.UUID
	var labNumber string

	err := db.Tx(ctx, func(txCtx context.Context) error {
		tx := db.Conn(txCtx)

		patientID = uuid.New()
		genderVal := "UNSPECIFIED"
		if payload.PatientDemographics.Gender != "" {
			genderVal = payload.PatientDemographics.Gender
		}
		_, err := tx.Exec(txCtx, `
			INSERT INTO patients (id, first_name, last_name, dob, gender)
			VALUES ($1, $2, $3, '1990-01-01'::date, $4)
		`, patientID, payload.PatientDemographics.FirstName, payload.PatientDemographics.LastName, genderVal)
		if err != nil {
			return fmt.Errorf("failed to insert patient: %w", err)
		}

		emailVal := ""
		if payload.PatientDemographics.Email != "" {
			emailVal = payload.PatientDemographics.Email
		}
		_, err = tx.Exec(txCtx, `
			INSERT INTO patient_contacts (patient_id, phone, email)
			VALUES ($1, $2, $3)
		`, patientID, payload.PatientDemographics.Phone, emailVal)
		if err != nil {
			return fmt.Errorf("failed to insert patient contacts: %w", err)
		}

		visitID := uuid.New()
		year := time.Now().Format("06")
		month := time.Now().Format("01")
		day := time.Now().Format("02")
		labNumber = year + month + day + "-10294"

		_, err = tx.Exec(txCtx, `
			INSERT INTO patient_visits (id, patient_id, accession_no)
			VALUES ($1, $2, $3)
		`, visitID, patientID, labNumber)
		if err != nil {
			return fmt.Errorf("failed to insert visit: %w", err)
		}

		orderID := uuid.New()
		_, err = tx.Exec(txCtx, `
			INSERT INTO lab_orders (id, visit_id, status)
			VALUES ($1, $2, 'ORDERED')
		`, orderID, visitID)
		if err != nil {
			return fmt.Errorf("failed to insert lab order: %w", err)
		}

		for _, inv := range payload.Investigations {
			sampleID := uuid.New()
			barcode := fmt.Sprintf("BAR-%s-%s", labNumber, inv.InvestigationID)
			_, err = tx.Exec(txCtx, `
				INSERT INTO lab_samples (id, order_id, barcode, specimen_type, status)
				VALUES ($1, $2, $3, 'Blood', 'ORDERED')
			`, sampleID, orderID, barcode)
			if err != nil {
				return fmt.Errorf("failed to insert sample: %w", err)
			}

			_, err = tx.Exec(txCtx, `
				INSERT INTO lab_results (sample_id, parameter_name, value, unit, status)
				VALUES ($1, $2, NULL, NULL, 'ENTERED')
			`, sampleID, inv.InvestigationID)
			if err != nil {
				return fmt.Errorf("failed to insert result placeholder: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		h.server.Logger.Error().Err(err).Msg("Failed to register patient visit transactionally")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register patient visit: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":     "success",
		"message":    "visit registered successfully",
		"lab_number": labNumber,
		"patient_id": patientID.String(),
		"released":   true,
	})
}
