package api

import (
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/organization/application"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DemoHandler struct {
	server     *server.Server
	appService *application.OrganizationApplicationService
}

func NewDemoHandler(s *server.Server, appService *application.OrganizationApplicationService) *DemoHandler {
	return &DemoHandler{
		server:     s,
		appService: appService,
	}
}

func (h *DemoHandler) CreateDemoRequest(c echo.Context) error {
	var payload CreateDemoRequestPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("invalid payload", false, nil, nil, nil)
	}
	if err := payload.Validate(); err != nil {
		return errs.NewBadRequestError(err.Error(), false, nil, nil, nil)
	}

	req, err := h.appService.CreateDemoRequest(c.Request().Context(), payload.LaboratoryName, payload.ContactName, payload.Email, payload.Phone, payload.DailyVolume, payload.Notes)
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to create demo request")
		return errs.NewInternalServerError()
	}

	return c.JSON(http.StatusCreated, req)
}

func (h *DemoHandler) ListDemoRequests(c echo.Context) error {
	reqs, err := h.appService.ListDemoRequests(c.Request().Context())
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to list demo requests")
		return c.JSON(http.StatusOK, []interface{}{})
	}
	if reqs == nil {
		return c.JSON(http.StatusOK, []interface{}{})
	}
	return c.JSON(http.StatusOK, reqs)
}

func (h *DemoHandler) UpdateDemoRequest(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errs.NewBadRequestError("invalid id format", false, nil, nil, nil)
	}

	var payload UpdateDemoRequestPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("invalid payload", false, nil, nil, nil)
	}
	if err := payload.Validate(); err != nil {
		return errs.NewBadRequestError(err.Error(), false, nil, nil, nil)
	}

	req, err := h.appService.UpdateDemoRequest(c.Request().Context(), id, payload.Status, payload.Notes)
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to update demo request")
		return errs.NewInternalServerError()
	}

	return c.JSON(http.StatusOK, req)
}
