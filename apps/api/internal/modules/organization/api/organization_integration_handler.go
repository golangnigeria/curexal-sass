package api

import (
	"errors"
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/organization/application"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/golangnigeria/curexal/internal/shared/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrganizationIntegrationHandler struct {
	integrationService *application.OrganizationIntegrationService
}

func NewOrganizationIntegrationHandler(integrationService *application.OrganizationIntegrationService) *OrganizationIntegrationHandler {
	return &OrganizationIntegrationHandler{
		integrationService: integrationService,
	}
}

func (h *OrganizationIntegrationHandler) ListAPIKeys(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	keys, err := h.integrationService.ListAPIKeys(c.Request().Context(), principal)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to list API keys: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, keys)
}

func (h *OrganizationIntegrationHandler) CreateAPIKey(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.CreateAPIKeyPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid API key payload")
	}

	res, err := h.integrationService.CreateAPIKey(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrInvalidIPWhitelist) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to create API key: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusCreated, res)
}

func (h *OrganizationIntegrationHandler) RevokeAPIKey(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	keyIDParam := c.Param("id")
	keyUUID, errParse := uuid.Parse(keyIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid API key ID format")
	}

	err := h.integrationService.RevokeAPIKey(c.Request().Context(), principal, keyUUID)
	if err != nil {
		if errors.Is(err, domain.ErrAPIKeyNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to revoke API key: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]string{"message": "API key revoked successfully"})
}

func (h *OrganizationIntegrationHandler) ListWebhookSubscriptions(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	subs, err := h.integrationService.ListWebhookSubscriptions(c.Request().Context(), principal)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to list webhook subscriptions: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, subs)
}

func (h *OrganizationIntegrationHandler) CreateWebhookSubscription(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.CreateWebhookSubscriptionPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid webhook subscription payload")
	}

	created, err := h.integrationService.CreateWebhookSubscription(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrSSRFURLForbidden) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrDuplicateWebhookTarget) {
			return response.ErrorEcho(c, http.StatusConflict, "DUPLICATE_WEBHOOK_TARGET", err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to create webhook subscription: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusCreated, created)
}

func (h *OrganizationIntegrationHandler) DeleteWebhookSubscription(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	subIDParam := c.Param("id")
	subUUID, errParse := uuid.Parse(subIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid webhook subscription ID format")
	}

	err := h.integrationService.DeleteWebhookSubscription(c.Request().Context(), principal, subUUID)
	if err != nil {
		if errors.Is(err, domain.ErrWebhookSubscriptionNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to delete webhook subscription: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]string{"message": "Webhook subscription deleted successfully"})
}

func (h *OrganizationIntegrationHandler) ListWebhookDeliveries(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	deliveries, err := h.integrationService.ListWebhookDeliveries(c.Request().Context(), principal, 50)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to list webhook deliveries: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, deliveries)
}
