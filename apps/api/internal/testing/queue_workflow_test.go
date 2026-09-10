package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	orchestrationHandler "github.com/golangnigeria/curexal/internal/modules/orchestration/handler"
	orchestrationModel "github.com/golangnigeria/curexal/internal/modules/orchestration/model"
	orchestrationRepo "github.com/golangnigeria/curexal/internal/modules/orchestration/repository"
	orchestrationService "github.com/golangnigeria/curexal/internal/modules/orchestration/service"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCareEngine() (*echo.Echo, *server.Server, *orchestrationService.CareRequestService) {
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

	repo := orchestrationRepo.NewCareRequestRepository(s)
	svc := orchestrationService.NewCareRequestService(s, repo)
	hnd := orchestrationHandler.NewCareRequestHandler(s, svc)

	g := e.Group("/api/v1/orchestration/requests")
	g.POST("", hnd.CreateCareRequest)
	g.GET("", hnd.ListCareRequests)
	g.GET("/:id", hnd.GetCareRequestByID)
	g.POST("/:id/triage", hnd.SubmitTriage)
	g.GET("/:id/match-provider", hnd.MatchProviders)
	g.POST("/:id/assign-provider", hnd.AssignProvider)

	return e, s, svc
}

// 1. Clinical Acuity Evaluation Test
func TestQueueWorkflow_EvaluateAcuity(t *testing.T) {
	_, _, svc := setupCareEngine()

	t.Run("Critical vitals evaluate to RED urgency", func(t *testing.T) {
		temp := 40.1
		spo2 := 89
		pain := 9
		bp := 195

		payload := orchestrationModel.SubmitTriagePayload{
			Temperature: &temp,
			SpO2:        &spo2,
			PainScore:   &pain,
			SystolicBP:  &bp,
		}

		acuity := svc.EvaluateAcuity(payload)
		assert.Equal(t, "RED", acuity)
	})

	t.Run("Moderate fever & mild tachycardia evaluate to YELLOW urgency", func(t *testing.T) {
		temp := 38.6
		spo2 := 94
		pain := 6
		bp := 145

		payload := orchestrationModel.SubmitTriagePayload{
			Temperature: &temp,
			SpO2:        &spo2,
			PainScore:   &pain,
			SystolicBP:  &bp,
		}

		acuity := svc.EvaluateAcuity(payload)
		assert.Equal(t, "YELLOW", acuity)
	})

	t.Run("Stable normal vitals evaluate to GREEN routine", func(t *testing.T) {
		temp := 36.8
		spo2 := 99
		pain := 1
		bp := 115

		payload := orchestrationModel.SubmitTriagePayload{
			Temperature: &temp,
			SpO2:        &spo2,
			PainScore:   &pain,
			SystolicBP:  &bp,
		}

		acuity := svc.EvaluateAcuity(payload)
		assert.Equal(t, "GREEN", acuity)
	})

	t.Run("Explicit clinician override takes precedence", func(t *testing.T) {
		temp := 36.8
		override := "RED"
		payload := orchestrationModel.SubmitTriagePayload{
			Temperature:    &temp,
			AcuityOverride: &override,
		}

		acuity := svc.EvaluateAcuity(payload)
		assert.Equal(t, "RED", acuity)
	})
}

// 2. Queue Check-In Endpoint Validations
func TestQueueWorkflow_CheckIn_EndpointValidation(t *testing.T) {
	e, _, _ := setupCareEngine()

	t.Run("CheckIn without PatientID fails with 400 Bad Request", func(t *testing.T) {
		body := map[string]interface{}{
			"serviceType": "GENERAL_CONSULTATION",
		}
		jsonBytes, err := json.Marshal(body)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/orchestration/requests", bytes.NewReader(jsonBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-Tenant-ID", "00000000-0000-0000-0000-000000000001")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Triage without CareRequest ID fails with 400 Bad Request", func(t *testing.T) {
		body := map[string]interface{}{
			"systolicBp": 120,
		}
		jsonBytes, err := json.Marshal(body)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/orchestration/requests/%20/triage", bytes.NewReader(jsonBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-Tenant-ID", "00000000-0000-0000-0000-000000000001")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
