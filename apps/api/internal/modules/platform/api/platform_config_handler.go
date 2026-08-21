package api

import (
	"errors"

	"github.com/golangnigeria/curexal/internal/modules/platform/application"
	"github.com/golangnigeria/curexal/internal/modules/platform/domain"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/golangnigeria/curexal/internal/shared/response"
	"github.com/labstack/echo/v4"
)

type PlatformConfigHandler struct {
	configService *application.PlatformConfigService
}

func NewPlatformConfigHandler(service *application.PlatformConfigService) *PlatformConfigHandler {
	return &PlatformConfigHandler{
		configService: service,
	}
}

func (h *PlatformConfigHandler) GetPlatformConfig(c echo.Context) error {
	config, err := h.configService.GetGeneralSettings(c.Request().Context())
	if err != nil {
		return response.InternalErrorEcho(c, err.Error())
	}
	return response.SuccessEcho(c, 200, config)
}

func (h *PlatformConfigHandler) UpdatePlatformConfig(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.PlatformGeneralSettings
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid request payload")
	}

	updated, err := h.configService.UpdateGeneralSettings(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedPlatformAdmin) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrOptimisticLockingConflict) {
			return response.ErrorEcho(c, 409, "CONFLICT", err.Error())
		}
		if errors.Is(err, domain.ErrInvalidPlatformConfig) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, 200, updated)
}

func (h *PlatformConfigHandler) GetSecurityPolicy(c echo.Context) error {
	policy, err := h.configService.GetSecurityPolicy(c.Request().Context())
	if err != nil {
		return response.InternalErrorEcho(c, err.Error())
	}
	return response.SuccessEcho(c, 200, policy)
}

func (h *PlatformConfigHandler) UpdateSecurityPolicy(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.IdentitySecurityPolicy
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid request payload")
	}

	updated, err := h.configService.UpdateSecurityPolicy(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedPlatformAdmin) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrOptimisticLockingConflict) {
			return response.ErrorEcho(c, 409, "CONFLICT", err.Error())
		}
		if errors.Is(err, domain.ErrInvalidSecurityPolicy) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, 200, updated)
}
