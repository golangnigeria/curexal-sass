package api

import (
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/audit/application"
	"github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuditHandler struct {
	server     *server.Server
	appService *application.AuditApplicationService
}

func NewAuditHandler(s *server.Server, appService *application.AuditApplicationService) *AuditHandler {
	return &AuditHandler{
		server:     s,
		appService: appService,
	}
}

func (h *AuditHandler) CreateAuditLog(c echo.Context) error {
	var payload CreateAuditLogPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("invalid request payload", false, nil, nil, nil)
	}

	logEntry, err := h.appService.LogEvent(c.Request().Context(), &payload)
	if err != nil {
		return errs.NewInternalServerError()
	}

	return c.JSON(http.StatusCreated, logEntry)
}

func (h *AuditHandler) GetAuditLog(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errs.NewBadRequestError("invalid audit log id format", false, nil, nil, nil)
	}

	log, err := h.appService.GetLogByID(c.Request().Context(), id)
	if err != nil {
		return errs.NewNotFoundError("Audit log not found")
	}

	return c.JSON(http.StatusOK, log)
}

func (h *AuditHandler) ListTenantLogs(c echo.Context) error {
	var payload ListAuditLogsPayload
	_ = c.Bind(&payload)

	activeTenantID := middleware.GetActiveTenantID(c)
	var tenantUUID *uuid.UUID
	if activeTenantID != "" {
		if parsed, err := uuid.Parse(activeTenantID); err == nil {
			tenantUUID = &parsed
		}
	}

	logs, err := h.appService.ListTenantLogs(
		c.Request().Context(),
		tenantUUID,
		payload.Category,
		payload.Severity,
		payload.Status,
		payload.ActorID,
		payload.Action,
		payload.ResourceType,
		payload.ResourceID,
		payload.StartDate,
		payload.EndDate,
		payload.Search,
		payload.Limit,
		payload.Offset,
	)
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to query tenant audit logs")
		return c.JSON(http.StatusOK, []interface{}{})
	}
	if logs == nil {
		logs = []domain.AuditLog{}
	}

	return c.JSON(http.StatusOK, logs)
}

func (h *AuditHandler) ListPlatformLogs(c echo.Context) error {
	var payload ListAuditLogsPayload
	_ = c.Bind(&payload)

	var orgUUID *uuid.UUID
	if payload.OrganizationID != nil && *payload.OrganizationID != "" {
		if parsed, err := uuid.Parse(*payload.OrganizationID); err == nil {
			orgUUID = &parsed
		}
	}

	logs, err := h.appService.ListPlatformLogs(
		c.Request().Context(),
		orgUUID,
		payload.Category,
		payload.Severity,
		payload.Status,
		payload.ActorID,
		payload.Action,
		payload.ResourceType,
		payload.ResourceID,
		payload.StartDate,
		payload.EndDate,
		payload.Search,
		payload.Limit,
		payload.Offset,
	)
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to query platform audit logs")
		return c.JSON(http.StatusOK, []interface{}{})
	}
	if logs == nil {
		logs = []domain.AuditLog{}
	}

	return c.JSON(http.StatusOK, logs)
}

func (h *AuditHandler) GetAdminStats(c echo.Context) error {
	stats, err := h.appService.GetAdminStats(c.Request().Context(), nil, nil)
	if err != nil {
		return errs.NewInternalServerError()
	}
	return c.JSON(http.StatusOK, stats)
}
