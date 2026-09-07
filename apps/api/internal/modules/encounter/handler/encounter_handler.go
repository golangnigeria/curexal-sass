package handler

import (
	"net/http"
	"strings"

	"github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/modules/encounter/model"
	"github.com/golangnigeria/curexal/internal/modules/encounter/service"
	"github.com/labstack/echo/v4"
)

type EncounterHandler struct {
	service *service.EncounterService
}

func NewEncounterHandler(s *service.EncounterService) *EncounterHandler {
	return &EncounterHandler{service: s}
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

// StartEncounter starts a clinical consultation
func (h *EncounterHandler) StartEncounter(c echo.Context) error {
	tenantID := resolveTenantID(c)

	var payload model.StartEncounterPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid start encounter payload")
	}

	enc, err := h.service.StartEncounter(c.Request().Context(), tenantID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Encounter started successfully",
		"data":    enc,
	})
}

// GetEncounterByID retrieves active encounter details
func (h *EncounterHandler) GetEncounterByID(c echo.Context) error {
	tenantID := resolveTenantID(c)
	encounterID := c.Param("id")
	if strings.TrimSpace(encounterID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Encounter ID is required")
	}

	enc, err := h.service.GetEncounterByID(c.Request().Context(), tenantID, encounterID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if enc == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Encounter not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    enc,
	})
}

// SaveSOAPNotes saves physician clinical notes
func (h *EncounterHandler) SaveSOAPNotes(c echo.Context) error {
	tenantID := resolveTenantID(c)
	encounterID := c.Param("id")
	if strings.TrimSpace(encounterID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Encounter ID is required")
	}

	var payload model.UpdateSOAPNotesPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid SOAP notes payload")
	}

	if err := h.service.SaveSOAPNotes(c.Request().Context(), tenantID, encounterID, payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "SOAP notes saved successfully",
	})
}

// DispatchOrders routes lab, radiology, and prescription orders
func (h *EncounterHandler) DispatchOrders(c echo.Context) error {
	tenantID := resolveTenantID(c)
	encounterID := c.Param("id")
	if strings.TrimSpace(encounterID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Encounter ID is required")
	}

	var payload model.DispatchOrdersPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid orders payload")
	}

	if err := h.service.DispatchOrders(c.Request().Context(), tenantID, encounterID, payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Diagnostic and pharmacy orders dispatched successfully",
	})
}

// CompleteEncounter discharges/finalizes encounter
func (h *EncounterHandler) CompleteEncounter(c echo.Context) error {
	tenantID := resolveTenantID(c)
	encounterID := c.Param("id")
	if strings.TrimSpace(encounterID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Encounter ID is required")
	}

	if err := h.service.CompleteEncounter(c.Request().Context(), tenantID, encounterID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Encounter finalized successfully",
	})
}
