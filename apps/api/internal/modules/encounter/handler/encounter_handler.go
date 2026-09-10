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

// ListEncounters lists encounters for a tenant with optional status/provider/patient filters
func (h *EncounterHandler) ListEncounters(c echo.Context) error {
	tenantID := resolveTenantID(c)

	var statusPtr *string
	if s := c.QueryParam("status"); s != "" {
		statusPtr = &s
	}

	var providerIDPtr *string
	if p := c.QueryParam("provider_id"); p != "" {
		providerIDPtr = &p
	}

	var patientIDPtr *string
	if pt := c.QueryParam("patient_id"); pt != "" {
		patientIDPtr = &pt
	}

	encounters, err := h.service.ListEncounters(c.Request().Context(), tenantID, statusPtr, providerIDPtr, patientIDPtr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    encounters,
		"count":   len(encounters),
	})
}

// StartEncounter starts a clinical consultation (unified in-person or telehealth)
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

// GetEncounterByID retrieves active encounter details including SOAP, diagnoses, prescriptions
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

// SaveSOAPNotes saves physician clinical notes and primary diagnosis
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

// AddDiagnosis records an ICD-10 diagnosis for the encounter
func (h *EncounterHandler) AddDiagnosis(c echo.Context) error {
	tenantID := resolveTenantID(c)
	encounterID := c.Param("id")
	if strings.TrimSpace(encounterID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Encounter ID is required")
	}

	var payload model.AddDiagnosisPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid diagnosis payload")
	}

	diag, err := h.service.AddDiagnosis(c.Request().Context(), tenantID, encounterID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Diagnosis recorded successfully",
		"data":    diag,
	})
}

// ListDiagnoses returns all diagnoses recorded for an encounter
func (h *EncounterHandler) ListDiagnoses(c echo.Context) error {
	encounterID := c.Param("id")
	if strings.TrimSpace(encounterID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Encounter ID is required")
	}

	diagnoses, err := h.service.ListDiagnoses(c.Request().Context(), encounterID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    diagnoses,
	})
}

// CreatePrescription creates an electronic prescription with line items
func (h *EncounterHandler) CreatePrescription(c echo.Context) error {
	tenantID := resolveTenantID(c)
	encounterID := c.Param("id")
	if strings.TrimSpace(encounterID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Encounter ID is required")
	}

	var payload model.CreatePrescriptionPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid prescription payload")
	}

	prescriberID := ""
	principal := auth.GetPrincipal(c)
	if principal != nil {
		prescriberID = principal.UserID
	}
	if prescriberID == "" {
		prescriberID = c.Request().Header.Get("X-Provider-ID")
	}

	rx, err := h.service.CreatePrescription(c.Request().Context(), tenantID, encounterID, prescriberID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Prescription order created successfully",
		"data":    rx,
	})
}

// ListPrescriptions returns all prescriptions for an encounter
func (h *EncounterHandler) ListPrescriptions(c echo.Context) error {
	encounterID := c.Param("id")
	if strings.TrimSpace(encounterID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Encounter ID is required")
	}

	prescriptions, err := h.service.ListPrescriptions(c.Request().Context(), encounterID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    prescriptions,
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

type completeEncounterReq struct {
	ConsultationFee float64 `json:"consultationFee"`
}

// CompleteEncounter discharges/finalizes encounter and generates POS invoice
func (h *EncounterHandler) CompleteEncounter(c echo.Context) error {
	tenantID := resolveTenantID(c)
	encounterID := c.Param("id")
	if strings.TrimSpace(encounterID) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Encounter ID is required")
	}

	var req completeEncounterReq
	_ = c.Bind(&req)

	res, err := h.service.CompleteEncounter(c.Request().Context(), tenantID, encounterID, req.ConsultationFee)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Encounter finalized and cashier invoice generated successfully",
		"data":    res,
	})
}
