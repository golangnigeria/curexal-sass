package api

import (
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/settings/application"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SettingsHandler struct {
	server     *server.Server
	appService *application.SettingsApplicationService
}

func NewSettingsHandler(s *server.Server, appService *application.SettingsApplicationService) *SettingsHandler {
	return &SettingsHandler{
		server:     s,
		appService: appService,
	}
}

func (h *SettingsHandler) GetBranchSettings(c echo.Context) error {
	tenantIDStr := middleware.GetActiveTenantID(c)
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid tenant context")
	}

	settings, err := h.appService.GetByTenantID(c.Request().Context(), tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve branch settings: "+err.Error())
	}

	return c.JSON(http.StatusOK, settings)
}

func (h *SettingsHandler) UpdateBranchSettingsSection(c echo.Context) error {
	tenantIDStr := middleware.GetActiveTenantID(c)
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid tenant context")
	}

	section := c.Param("section")

	var payload map[string]any
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid settings payload")
	}

	settings, err := h.appService.UpdateSection(c.Request().Context(), tenantID, section, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update branch settings: "+err.Error())
	}

	return c.JSON(http.StatusOK, settings)
}

func (h *SettingsHandler) ResetBranchSettings(c echo.Context) error {
	tenantIDStr := middleware.GetActiveTenantID(c)
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid tenant context")
	}

	sectionParam := c.QueryParam("section")
	var section *string
	if sectionParam != "" {
		section = &sectionParam
	}

	settings, err := h.appService.ResetToDefaults(c.Request().Context(), tenantID, section)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to reset branch settings: "+err.Error())
	}

	return c.JSON(http.StatusOK, settings)
}
