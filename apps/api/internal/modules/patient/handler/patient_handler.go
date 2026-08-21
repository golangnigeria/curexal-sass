package handler

import (
	"net/http"

	patientModel "github.com/golangnigeria/curexal/internal/modules/patient/model"
	patientService "github.com/golangnigeria/curexal/internal/modules/patient/service"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type PatientHandler struct {
	patientService *patientService.PatientService
}

func NewPatientHandler(patientService *patientService.PatientService) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
	}
}

// GetProfile returns the PatientProfile for the authenticated user.
func (h *PatientHandler) GetProfile(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	profile, err := h.patientService.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load patient profile")
	}

	if profile == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Patient profile not found")
	}

	return c.JSON(http.StatusOK, profile)
}

// UpdateProfile updates the PatientProfile for the authenticated user.
func (h *PatientHandler) UpdateProfile(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	var payload patientModel.UpdatePatientProfilePayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := h.patientService.UpdateProfile(c.Request().Context(), userID, payload); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Patient profile updated successfully",
	})
}

// GetResults is a placeholder endpoint returning lab results for the patient.
func (h *PatientHandler) GetResults(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"results": []interface{}{},
		"total":   0,
	})
}

// GetOrders is a placeholder endpoint returning lab orders for the patient.
func (h *PatientHandler) GetOrders(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"orders": []interface{}{},
		"total":  0,
	})
}

// GetAppointments is a placeholder endpoint returning appointments for the patient.
func (h *PatientHandler) GetAppointments(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"appointments": []interface{}{},
		"total":        0,
	})
}
