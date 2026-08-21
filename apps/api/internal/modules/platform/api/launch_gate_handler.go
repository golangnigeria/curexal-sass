package api

import (
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/platform/application"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/golangnigeria/curexal/internal/shared/response"
	"github.com/labstack/echo/v4"
)

type LaunchGateHandler struct {
	launchGateService *application.LaunchGateService
}

func NewLaunchGateHandler(launchGateService *application.LaunchGateService) *LaunchGateHandler {
	return &LaunchGateHandler{
		launchGateService: launchGateService,
	}
}

func (h *LaunchGateHandler) GetStatus(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	audit, err := h.launchGateService.GetStatus(c.Request().Context(), principal)
	if err != nil {
		return response.ForbiddenEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, audit)
}

func (h *LaunchGateHandler) VerifyProductionReadiness(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	audit, err := h.launchGateService.VerifyProductionReadiness(c.Request().Context(), principal)
	if err != nil {
		return response.ForbiddenEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, audit)
}

func (h *LaunchGateHandler) ListHealthMetrics(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	metrics, err := h.launchGateService.ListHealthMetrics(c.Request().Context(), principal)
	if err != nil {
		return response.ForbiddenEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, metrics)
}
