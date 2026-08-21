package api

import (
	"errors"
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/facility_config/application"
	"github.com/golangnigeria/curexal/internal/modules/facility_config/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/golangnigeria/curexal/internal/shared/response"
	"github.com/labstack/echo/v4"
)

type FacilityConfigHandler struct {
	server     *server.Server
	appService *application.FacilityConfigApplicationService
}

func NewFacilityConfigHandler(s *server.Server, appService *application.FacilityConfigApplicationService) *FacilityConfigHandler {
	return &FacilityConfigHandler{
		server:     s,
		appService: appService,
	}
}

func (h *FacilityConfigHandler) GetActiveFacilityTypes(c echo.Context) error {
	ctx := c.Request().Context()
	types, err := h.appService.GetActiveFacilityTypes(ctx)
	if err != nil {
		return response.InternalErrorEcho(c, "Failed to fetch facility types: "+err.Error())
	}
	return response.SuccessEcho(c, http.StatusOK, types)
}

func (h *FacilityConfigHandler) CreateFacilityType(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.FacilityTypeEntity
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid request payload")
	}

	created, err := h.appService.CreateFacilityType(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidFacilityType) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.ForbiddenEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusCreated, created)
}

func (h *FacilityConfigHandler) UpdateFacilityType(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.FacilityTypeEntity
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid request payload")
	}

	updated, err := h.appService.UpdateFacilityType(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrOptimisticLockingConflict) {
			return response.ErrorEcho(c, http.StatusConflict, "CONFLICT", err.Error())
		}
		if errors.Is(err, domain.ErrInvalidFacilityType) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.ForbiddenEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, updated)
}

func (h *FacilityConfigHandler) GetRegistrationForm(c echo.Context) error {
	id := c.Param("typeId")
	ctx := c.Request().Context()
	form, err := h.appService.GetRegistrationForm(ctx, id)
	if err != nil {
		return response.NotFoundEcho(c, "Registration form not found for facility type: "+err.Error())
	}
	return response.SuccessEcho(c, http.StatusOK, form)
}

func (h *FacilityConfigHandler) GetNavigationMenu(c echo.Context) error {
	id := c.Param("typeId")
	ctx := c.Request().Context()
	nav, err := h.appService.GetNavigationMenu(ctx, id)
	if err != nil {
		return response.NotFoundEcho(c, "Navigation menu not found for facility type: "+err.Error())
	}
	return response.SuccessEcho(c, http.StatusOK, nav)
}

func (h *FacilityConfigHandler) GetSetupSteps(c echo.Context) error {
	id := c.Param("typeId")
	ctx := c.Request().Context()
	steps, err := h.appService.GetSetupSteps(ctx, id)
	if err != nil {
		return response.InternalErrorEcho(c, "Failed to fetch setup steps: "+err.Error())
	}
	return response.SuccessEcho(c, http.StatusOK, steps)
}

func (h *FacilityConfigHandler) GetDashboard(c echo.Context) error {
	id := c.Param("typeId")
	ctx := c.Request().Context()
	dash, err := h.appService.GetDashboard(ctx, id)
	if err != nil {
		return response.NotFoundEcho(c, "Dashboard configuration not found for facility type: "+err.Error())
	}
	return response.SuccessEcho(c, http.StatusOK, dash)
}
