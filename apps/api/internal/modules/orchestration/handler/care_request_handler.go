package handler

import (
	"net/http"
	"strings"

	"github.com/golangnigeria/curexal/internal/modules/orchestration/model"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/service"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
)

type CareRequestHandler struct {
	server  *server.Server
	service *service.CareRequestService
}

func NewCareRequestHandler(s *server.Server, service *service.CareRequestService) *CareRequestHandler {
	return &CareRequestHandler{
		server:  s,
		service: service,
	}
}

// resolveTenantID extracts active branch/tenant ID from context or header
func (h *CareRequestHandler) resolveTenantID(c echo.Context) string {
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

// resolvePatientID extracts authenticated patient ID from headers, params, or identity session
func (h *CareRequestHandler) resolvePatientID(c echo.Context) string {
	if pid := c.Request().Header.Get("X-Patient-ID"); pid != "" && pid != "undefined" {
		return pid
	}
	if pid := c.QueryParam("patientId"); pid != "" && pid != "undefined" {
		return pid
	}

	// 1. Check Authenticated Principal / User ID from token/session
	var userID string
	if p := platformAuth.GetPrincipal(c); p != nil && p.UserID != "" {
		userID = p.UserID
	} else if u := c.Get(platformAuth.UserIDKey); u != nil {
		if s, ok := u.(string); ok && s != "" {
			userID = s
		}
	}

	if userID != "" && h.server != nil && h.server.DB != nil {
		var patientID string
		_ = h.server.DB.Conn(c.Request().Context()).QueryRow(
			c.Request().Context(),
			`SELECT id::text FROM patient.patients WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`,
			userID,
		).Scan(&patientID)
		if patientID != "" {
			return patientID
		}
	}

	return ""
}

// CreateCareRequest creates a new care request and initializes the Care Journey
func (h *CareRequestHandler) CreateCareRequest(c echo.Context) error {
	var payload model.CreateCareRequestPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid care request payload")
	}

	tenantID := h.resolveTenantID(c)
	if payload.PatientID == "" {
		payload.PatientID = h.resolvePatientID(c)
	}

	if payload.PatientID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "patientId is required to submit a care request")
	}

	req, err := h.service.CreateCareRequest(c.Request().Context(), tenantID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Care request submitted successfully",
		"data":    req,
	})
}

// ListCareRequests lists care requests with filters
func (h *CareRequestHandler) ListCareRequests(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	var filter model.CareRequestFilter
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid query parameters")
	}

	res, err := h.service.ListCareRequests(c.Request().Context(), tenantID, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve care requests")
	}

	return c.JSON(http.StatusOK, res)
}

// GetCareRequestByID retrieves single request by ID
func (h *CareRequestHandler) GetCareRequestByID(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	requestID := c.Param("id")
	if strings.TrimSpace(requestID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Request ID is required")
	}

	req, err := h.service.GetCareRequestByID(c.Request().Context(), tenantID, requestID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve care request")
	}
	if req == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Care request not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    req,
	})
}

// GetMyCareRequests returns care requests submitted by the logged-in patient
func (h *CareRequestHandler) GetMyCareRequests(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	patientID := h.resolvePatientID(c)
	if patientID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authenticated patient identity not provided")
	}

	requests, err := h.service.ListPatientCareRequests(c.Request().Context(), tenantID, patientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve patient care requests")
	}
	if requests == nil {
		requests = make([]model.CareRequest, 0)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    requests,
	})
}

// GetMyCareJourney returns live milestones along the patient's Care Journey
func (h *CareRequestHandler) GetMyCareJourney(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	patientID := h.resolvePatientID(c)
	if patientID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authenticated patient identity not provided")
	}

	milestones, err := h.service.GetPatientCareJourney(c.Request().Context(), tenantID, patientID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve patient care journey")
	}
	if milestones == nil {
		milestones = make([]model.CareJourneyMilestone, 0)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    milestones,
	})
}

// SubmitTriage records clinical vitals and updates acuity
func (h *CareRequestHandler) SubmitTriage(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	requestID := c.Param("id")
	if strings.TrimSpace(requestID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Request ID is required")
	}

	var payload model.SubmitTriagePayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid triage vitals payload")
	}

	assessorID := "staff_nurse_01" // Fallback assessor
	triage, err := h.service.SubmitTriage(c.Request().Context(), tenantID, requestID, assessorID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Triage vitals recorded successfully",
		"data":    triage,
	})
}

// MatchProviders calculates ranked provider candidates
func (h *CareRequestHandler) MatchProviders(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	requestID := c.Param("id")
	if strings.TrimSpace(requestID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Request ID is required")
	}

	res, err := h.service.MatchProviders(c.Request().Context(), tenantID, requestID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}

// AssignProvider assigns doctor to the request and launches consultation
func (h *CareRequestHandler) AssignProvider(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	requestID := c.Param("id")
	if strings.TrimSpace(requestID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Request ID is required")
	}

	var payload model.AssignProviderPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid provider assignment payload")
	}

	if err := h.service.AssignProvider(c.Request().Context(), tenantID, requestID, payload.ProviderID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Provider assigned successfully",
	})
}
