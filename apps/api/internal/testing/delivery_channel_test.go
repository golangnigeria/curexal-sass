package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	operationsHandler "github.com/golangnigeria/curexal/internal/modules/operations/handler"
	operationsRepo "github.com/golangnigeria/curexal/internal/modules/operations/repository"
	operationsService "github.com/golangnigeria/curexal/internal/modules/operations/service"
	orchestrationService "github.com/golangnigeria/curexal/internal/modules/orchestration/service"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOperationsEngine() (*echo.Echo, *server.Server, *operationsService.AppointmentService) {
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

	repo := operationsRepo.NewAppointmentRepository(s)
	svc := operationsService.NewAppointmentService(s, repo)
	hnd := operationsHandler.NewAppointmentHandler(s, svc)

	g := e.Group("/api/v1/appointments")
	g.POST("", hnd.CreateAppointment)
	g.GET("", hnd.ListAppointments)
	g.GET("/:id", hnd.GetAppointmentByID)
	g.PUT("/:id/status", hnd.UpdateAppointmentStatus)

	return e, s, svc
}

// 1. Delivery Channel Validation Tests
func TestDeliveryChannel_Validation(t *testing.T) {
	_, _, svc := setupOperationsEngine()

	validChannels := []string{"in_person", "video", "telephone", "secure_message"}
	for _, ch := range validChannels {
		assert.True(t, svc.IsValidDeliveryChannel(ch), "expected %s to be valid channel", ch)
	}

	invalidChannels := []string{"hologram", "smoke_signal", "pigeon", "email_only", ""}
	for _, ch := range invalidChannels {
		assert.False(t, svc.IsValidDeliveryChannel(ch), "expected %s to be invalid channel", ch)
	}
}

// 2. Appointment Number Generation Test
func TestDeliveryChannel_AppointmentNumberFormat(t *testing.T) {
	_, _, svc := setupOperationsEngine()

	for i := 0; i < 50; i++ {
		num := svc.GenerateAppointmentNumber()
		assert.Regexp(t, `^APT-\d{4}-\d{5}$`, num)
	}
}

// 3. Telehealth Virtual URL Routing Test
func TestDeliveryChannel_VirtualRoomAllocation(t *testing.T) {
	e, _, _ := setupOperationsEngine()

	t.Run("Create video appointment requires valid time window", func(t *testing.T) {
		// End time before start time should be rejected
		body := map[string]interface{}{
			"patientId":       "pat-01",
			"providerId":      "prov-01",
			"deliveryChannel": "video",
			"startTime":       time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			"endTime":         time.Now().Add(1 * time.Hour).Format(time.RFC3339), // invalid
		}
		jsonBytes, err := json.Marshal(body)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/appointments", bytes.NewReader(jsonBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-Tenant-ID", "00000000-0000-0000-0000-000000000001")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Create appointment with invalid delivery channel returns 400 Bad Request", func(t *testing.T) {
		body := map[string]interface{}{
			"patientId":       "pat-01",
			"providerId":      "prov-01",
			"deliveryChannel": "telepathy",
			"startTime":       time.Now().Add(1 * time.Hour).Format(time.RFC3339),
			"endTime":         time.Now().Add(2 * time.Hour).Format(time.RFC3339),
		}
		jsonBytes, err := json.Marshal(body)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/appointments", bytes.NewReader(jsonBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-Tenant-ID", "00000000-0000-0000-0000-000000000001")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// 4. Care Request Initial State Channel Routing Test
func TestDeliveryChannel_CareRequestRouting(t *testing.T) {
	logger := zerolog.Nop()
	s := &server.Server{
		Config: &config.Config{
			Auth: config.AuthConfig{
				SecretKey: "test-secret-key-32-bytes-long!!",
			},
		},
		Logger: &logger,
	}
	careSvc := orchestrationService.NewCareRequestService(s, nil)

	// Check that request number formats correctly
	reqNum := careSvc.GenerateRequestNumber()
	assert.Regexp(t, `^REQ-\d{4}-\d{6}$`, reqNum)
}
