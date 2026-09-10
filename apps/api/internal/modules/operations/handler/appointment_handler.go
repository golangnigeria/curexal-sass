package handler

import (
	"net/http"
	"strings"
	"time"

	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/operations/model"
	"github.com/golangnigeria/curexal/internal/modules/operations/service"
	"github.com/labstack/echo/v4"
)

type AppointmentHandler struct {
	server  *server.Server
	service *service.AppointmentService
}

func NewAppointmentHandler(s *server.Server, svc *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		server:  s,
		service: svc,
	}
}

// resolveTenantID extracts active branch/tenant ID from context or header
func (h *AppointmentHandler) resolveTenantID(c echo.Context) string {
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

// CreateAppointment handles POST /api/v1/appointments
func (h *AppointmentHandler) CreateAppointment(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	var req model.CreateAppointmentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid appointment payload")
	}

	var createdBy *string
	if p := platformAuth.GetPrincipal(c); p != nil && p.UserID != "" {
		createdBy = &p.UserID
	}

	apt, err := h.service.CreateAppointment(c.Request().Context(), tenantID, createdBy, req)
	if err != nil {
		if strings.Contains(err.Error(), "overlapping") || strings.Contains(err.Error(), "existing booking") {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Appointment scheduled successfully",
		"data":    apt,
	})
}

// ListAppointments handles GET /api/v1/appointments
func (h *AppointmentHandler) ListAppointments(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	var filter model.AppointmentFilter
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid appointment query parameters")
	}

	appointments, err := h.service.ListAppointments(c.Request().Context(), tenantID, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve appointments")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": appointments,
		"meta": map[string]interface{}{
			"total":     len(appointments),
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
}

// GetAppointmentByID handles GET /api/v1/appointments/:id
func (h *AppointmentHandler) GetAppointmentByID(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Appointment ID is required")
	}

	apt, err := h.service.GetAppointmentByID(c.Request().Context(), tenantID, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve appointment")
	}
	if apt == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Appointment not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": apt,
	})
}

// UpdateAppointmentStatus handles PUT /api/v1/appointments/:id/status
func (h *AppointmentHandler) UpdateAppointmentStatus(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Appointment ID is required")
	}

	var req model.UpdateAppointmentStatusRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid appointment status payload")
	}

	if err := h.service.UpdateAppointmentStatus(c.Request().Context(), tenantID, id, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Appointment status updated successfully",
		"data": map[string]interface{}{
			"id":     id,
			"status": strings.ToUpper(req.Status),
		},
	})
}
