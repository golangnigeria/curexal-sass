package api

import (
	"io"
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/subscription/application"
	"github.com/golangnigeria/curexal/internal/modules/subscription/domain"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/golangnigeria/curexal/internal/shared/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CommercialHandler struct {
	commercialSvc *application.CommercialService
}

func NewCommercialHandler(commercialSvc *application.CommercialService) *CommercialHandler {
	return &CommercialHandler{
		commercialSvc: commercialSvc,
	}
}

func (h *CommercialHandler) CreateCommercialOrder(c echo.Context) error {
	orgIDParam := c.Param("id")
	orgID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid organization ID format")
	}

	var req domain.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("invalid request payload")
	}

	provider := c.QueryParam("provider")
	if provider == "" {
		provider = "mock"
	}

	actorID := middleware.GetUserID(c)
	var actorUUID *uuid.UUID
	if actorID != "" {
		if parsed, errP := uuid.Parse(actorID); errP == nil {
			actorUUID = &parsed
		}
	}

	res, err := h.commercialSvc.CreateCommercialOrder(c.Request().Context(), orgID, actorUUID, req, provider)
	if err != nil {
		return errs.NewBadRequestError(err.Error())
	}

	return response.SuccessEcho(c, http.StatusCreated, res)
}

func (h *CommercialHandler) HandlePaymentWebhook(c echo.Context) error {
	provider := c.Param("provider")
	if provider == "" {
		return errs.NewBadRequestError("provider parameter required")
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return errs.NewBadRequestError("failed to read webhook body")
	}

	if err := h.commercialSvc.ProcessWebhookEvent(c.Request().Context(), provider, c.Request(), body); err != nil {
		if err.Error() == "invalid webhook signature" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid cryptographic signature"})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "processed"})
}
