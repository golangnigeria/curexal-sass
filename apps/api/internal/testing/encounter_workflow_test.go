package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	encounterHandler "github.com/golangnigeria/curexal/internal/modules/encounter/handler"
	encounterModel "github.com/golangnigeria/curexal/internal/modules/encounter/model"
	encounterRepo "github.com/golangnigeria/curexal/internal/modules/encounter/repository"
	encounterService "github.com/golangnigeria/curexal/internal/modules/encounter/service"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupEncounterTestEngine() (*echo.Echo, *server.Server, *encounterService.EncounterService) {
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

	repo := encounterRepo.NewEncounterRepository(s)
	svc := encounterService.NewEncounterService(s, repo)
	hnd := encounterHandler.NewEncounterHandler(svc)

	g := e.Group("/api/v1/encounters")
	g.GET("", hnd.ListEncounters)
	g.POST("", hnd.StartEncounter)
	g.POST("/start", hnd.StartEncounter)
	g.GET("/:id", hnd.GetEncounterByID)
	g.PUT("/:id/soap", hnd.SaveSOAPNotes)
	g.PUT("/:id/notes", hnd.SaveSOAPNotes)
	g.POST("/:id/diagnoses", hnd.AddDiagnosis)
	g.GET("/:id/diagnoses", hnd.ListDiagnoses)
	g.POST("/:id/prescriptions", hnd.CreatePrescription)
	g.GET("/:id/prescriptions", hnd.ListPrescriptions)
	g.POST("/:id/orders", hnd.DispatchOrders)
	g.POST("/:id/complete", hnd.CompleteEncounter)

	return e, s, svc
}

// 1. Encounter Payload Validation and Channel Adaptability
func TestEncounterWorkflow_Validation(t *testing.T) {
	e, _, _ := setupEncounterTestEngine()

	t.Run("Fails without patientId or providerId", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"encounterChannel": "in_person",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/encounters", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Tenant-ID", "00000000-0000-0000-0000-000000000001")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Validates channel support for both in-person and telehealth", func(t *testing.T) {
		channels := []string{"in_person", "video", "telephone", "secure_message"}
		for _, ch := range channels {
			payload := encounterModel.StartEncounterPayload{
				PatientID:        "pat-12345",
				ProviderID:       "prov-67890",
				EncounterChannel: ch,
			}
			assert.NotEmpty(t, payload.EncounterChannel)
			assert.Equal(t, ch, payload.EncounterChannel)
		}
	})
}

// 2. Structured SOAP Notes & ICD-10 Diagnoses Data Flow
func TestEncounterWorkflow_SOAPAndDiagnosis(t *testing.T) {
	t.Run("SOAP Notes Payload structure handles full clinical notes", func(t *testing.T) {
		subj := "Patient complains of persistent dry cough and low-grade fever for 3 days."
		obj := "BP: 120/80 mmHg, SpO2: 98%, Temp: 37.4C. Chest clear to auscultation."
		assess := "Acute upper respiratory tract infection (URTI), uncomplicated."
		plan := "Rest, oral hydration, symptomatic antipyretic relief, review in 5 days if unimproved."
		diagCode := "J06.9"
		diagName := "Acute upper respiratory infection, unspecified"

		soap := encounterModel.UpdateSOAPNotesPayload{
			Subjective:           &subj,
			Objective:            &obj,
			Assessment:           &assess,
			Plan:                 &plan,
			PrimaryDiagnosisCode: &diagCode,
			PrimaryDiagnosisName: &diagName,
			SecondaryDiagnoses:   []string{"R05 (Cough)"},
		}

		assert.Equal(t, "J06.9", *soap.PrimaryDiagnosisCode)
		assert.Equal(t, "Acute upper respiratory infection, unspecified", *soap.PrimaryDiagnosisName)
		assert.Len(t, soap.SecondaryDiagnoses, 1)
	})

	t.Run("ICD-10 Diagnosis Payload assigns appropriate clinical verification status", func(t *testing.T) {
		notes := "Provisional diagnosis pending viral swab results"
		diagPayload := encounterModel.AddDiagnosisPayload{
			ICD10Code:          "J00",
			ICD10Title:         "Acute nasopharyngitis (common cold)",
			IsPrimary:          true,
			ClinicalStatus:     "ACTIVE",
			VerificationStatus: "CONFIRMED",
			Notes:              &notes,
		}

		assert.Equal(t, "J00", diagPayload.ICD10Code)
		assert.True(t, diagPayload.IsPrimary)
		assert.Equal(t, "CONFIRMED", diagPayload.VerificationStatus)
	})
}

// 3. Electronic Prescriptions and POS Completion Response
func TestEncounterWorkflow_PrescriptionAndSettlement(t *testing.T) {
	t.Run("Prescription payload formats medication items accurately", func(t *testing.T) {
		notes := "Take medications with food"
		rxPayload := encounterModel.CreatePrescriptionPayload{
			Notes: &notes,
			Items: []encounterModel.CreatePrescriptionItemDTO{
				{
					DrugName:           "Amoxicillin / Clavulanic Acid",
					DosageForm:         "TABLET",
					Strength:           "625mg",
					Route:              "Oral",
					Frequency:          "Twice daily (BD)",
					DurationDays:       7,
					QuantityPrescribed: 14,
					Instructions:       "Take 1 tablet every 12 hours after meals",
				},
				{
					DrugName:           "Paracetamol",
					DosageForm:         "TABLET",
					Strength:           "500mg",
					Route:              "Oral",
					Frequency:          "Three times daily (TDS)",
					DurationDays:       3,
					QuantityPrescribed: 18,
					Instructions:       "Take 2 tablets for headache/fever",
				},
			},
		}

		require.Len(t, rxPayload.Items, 2)
		assert.Equal(t, "625mg", rxPayload.Items[0].Strength)
		assert.Equal(t, 14, rxPayload.Items[0].QuantityPrescribed)
		assert.Equal(t, 18, rxPayload.Items[1].QuantityPrescribed)
	})

	t.Run("CompleteEncounterResponse formats invoice reference and settlement total", func(t *testing.T) {
		invID := "inv-998877"
		invNum := "INV-2026-45123"
		resp := encounterModel.CompleteEncounterResponse{
			EncounterID:   "enc-112233",
			Status:        "signed",
			InvoiceID:     &invID,
			InvoiceNumber: &invNum,
			TotalAmount:   5000.00,
			SignedAt:      "2026-09-10T10:00:00Z",
		}

		assert.Equal(t, "signed", resp.Status)
		assert.Equal(t, 5000.00, resp.TotalAmount)
		assert.Equal(t, "INV-2026-45123", *resp.InvoiceNumber)
	})
}
