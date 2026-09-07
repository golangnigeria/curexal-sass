package handler

import (
	"net/http"
	"strings"

	patientModel "github.com/golangnigeria/curexal/internal/modules/patient/model"
	patientService "github.com/golangnigeria/curexal/internal/modules/patient/service"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
)

type CanonicalPatientHandler struct {
	server           *server.Server
	canonicalService *patientService.CanonicalPatientService
	mpiService       *patientService.MPIService
}

func NewCanonicalPatientHandler(
	s *server.Server,
	canonicalService *patientService.CanonicalPatientService,
	mpiService *patientService.MPIService,
) *CanonicalPatientHandler {
	return &CanonicalPatientHandler{
		server:           s,
		canonicalService: canonicalService,
		mpiService:       mpiService,
	}
}

// resolveTenantID extracts active branch/tenant ID from context or header
func (h *CanonicalPatientHandler) resolveTenantID(c echo.Context) string {
	if tid := c.Request().Header.Get("X-Tenant-ID"); tid != "" && tid != "undefined" {
		return tid
	}
	if bid := c.Request().Header.Get("X-Branch-ID"); bid != "" && bid != "undefined" {
		return bid
	}
	if orgID := c.Request().Header.Get("X-Organization-ID"); orgID != "" && orgID != "undefined" {
		if h.server != nil && h.server.DB != nil {
			var branchID string
			_ = h.server.DB.Conn(c.Request().Context()).QueryRow(
				c.Request().Context(),
				`SELECT id::text FROM organization.facility_branches WHERE organization_id = $1 OR id = $1 ORDER BY is_headquarters DESC, created_at ASC LIMIT 1`,
				orgID,
			).Scan(&branchID)
			if branchID != "" {
				return branchID
			}
		}
		return orgID
	}

	if h.server != nil && h.server.DB != nil {
		var defaultBranchID string
		_ = h.server.DB.Conn(c.Request().Context()).QueryRow(
			c.Request().Context(),
			`SELECT id::text FROM organization.facility_branches ORDER BY is_headquarters DESC, created_at ASC LIMIT 1`,
		).Scan(&defaultBranchID)
		if defaultBranchID != "" {
			return defaultBranchID
		}
	}

	return ""
}

// ResolveDuplicates evaluates identity signals for existing duplicate patients
func (h *CanonicalPatientHandler) ResolveDuplicates(c echo.Context) error {
	var payload patientModel.DuplicateEvaluationRequest
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid duplicate evaluation payload")
	}

	tenantID := h.resolveTenantID(c)
	result, err := h.mpiService.EvaluateDuplicates(c.Request().Context(), tenantID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to evaluate duplicates")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Duplicate evaluation completed",
		"data":    result,
	})
}

// RegisterPatient creates a canonical patient with Master Patient Index duplicate detection
func (h *CanonicalPatientHandler) RegisterPatient(c echo.Context) error {
	var payload patientModel.RegisterCanonicalPatientPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid registration payload")
	}

	tenantID := h.resolveTenantID(c)
	patient, accessToken, refreshToken, duplicateRes, err := h.canonicalService.RegisterCanonicalPatient(
		c.Request().Context(), tenantID, payload, c.RealIP(), c.Request().UserAgent(),
	)
	if err != nil {
		if duplicateRes != nil {
			return c.JSON(http.StatusConflict, echo.Map{
				"code":       "PATIENT_DUPLICATE_SUSPECT",
				"message":    err.Error(),
				"duplicates": duplicateRes,
			})
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Set HttpOnly session cookies if access token was generated
	if accessToken != "" && h.server != nil && h.server.Config != nil {
		platformAuth.SetSessionCookies(c, h.server.Config, accessToken, refreshToken)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success":     true,
		"message":     "Patient registered successfully",
		"data":        patient,
		"token":       accessToken,
		"accessToken": accessToken,
	})
}

// ListPatients retrieves a paginated patient directory
func (h *CanonicalPatientHandler) ListPatients(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	var filter patientModel.PatientListFilter
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid list query parameters")
	}

	res, err := h.canonicalService.ListPatients(c.Request().Context(), tenantID, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve patients")
	}

	return c.JSON(http.StatusOK, res)
}

// GetPatientByID retrieves a patient's canonical profile
func (h *CanonicalPatientHandler) GetPatientByID(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	patientID := c.Param("id")
	if strings.TrimSpace(patientID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Patient ID is required")
	}

	patient, err := h.canonicalService.GetPatientByID(c.Request().Context(), tenantID, patientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve patient")
	}
	if patient == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Patient not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    patient,
	})
}

// SendPortalOTP triggers passwordless OTP dispatch
func (h *CanonicalPatientHandler) SendPortalOTP(c echo.Context) error {
	var payload patientModel.SendPortalOTPPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid OTP payload")
	}

	tenantID := h.resolveTenantID(c)
	code, err := h.canonicalService.SendPortalOTP(c.Request().Context(), tenantID, payload.Identifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Portal verification OTP sent successfully",
		"data": echo.Map{
			"identifier": payload.Identifier,
			"debugCode":  code, // Accessible for dev/demo testing
		},
	})
}

// VerifyPortalOTP verifies the 6-digit OTP code and establishes an authenticated session
func (h *CanonicalPatientHandler) VerifyPortalOTP(c echo.Context) error {
	var payload patientModel.VerifyPortalOTPPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid OTP verification payload")
	}

	tenantID := h.resolveTenantID(c)
	res, err := h.canonicalService.VerifyPortalOTP(
		c.Request().Context(), tenantID, payload.Identifier, payload.Code, c.RealIP(), c.Request().UserAgent(),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if res.AccessToken != "" && h.server != nil && h.server.Config != nil {
		platformAuth.SetSessionCookies(c, h.server.Config, res.AccessToken, res.RefreshToken)
	}

	return c.JSON(http.StatusOK, res)
}

// SetPortalPIN sets a quick 4-6 digit PIN
func (h *CanonicalPatientHandler) SetPortalPIN(c echo.Context) error {
	var payload patientModel.SetPortalPINPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid PIN payload")
	}

	tenantID := h.resolveTenantID(c)
	patientID := c.Param("id")
	if err := h.canonicalService.SetPortalPIN(c.Request().Context(), tenantID, patientID, payload.PIN); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Portal PIN configured successfully",
	})
}
