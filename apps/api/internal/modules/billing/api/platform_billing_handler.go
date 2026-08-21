package api

import (
	"errors"
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/billing/application"
	"github.com/golangnigeria/curexal/internal/modules/billing/domain"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/golangnigeria/curexal/internal/shared/response"
	"github.com/labstack/echo/v4"
)

type PlatformPricingHandler struct {
	pricingService *application.PlatformPricingService
}

func NewPlatformPricingHandler(service *application.PlatformPricingService) *PlatformPricingHandler {
	return &PlatformPricingHandler{
		pricingService: service,
	}
}

func (h *PlatformPricingHandler) ListPricingRules(c echo.Context) error {
	rules, err := h.pricingService.ListPricingRules(c.Request().Context())
	if err != nil {
		return response.InternalErrorEcho(c, "Failed to list platform pricing rules: "+err.Error())
	}
	return response.SuccessEcho(c, http.StatusOK, rules)
}

func (h *PlatformPricingHandler) UpdatePricingRule(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.PricingRule
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid pricing rule payload")
	}

	updated, err := h.pricingService.UpdatePricingRule(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedPlatformAdmin) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrOptimisticLockingConflict) {
			return response.ErrorEcho(c, http.StatusConflict, "CONFLICT", err.Error())
		}
		if errors.Is(err, domain.ErrInvalidPricingRule) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, updated)
}

type PaymentGatewayVaultHandler struct {
	vaultService *application.PaymentGatewayVaultService
}

func NewPaymentGatewayVaultHandler(service *application.PaymentGatewayVaultService) *PaymentGatewayVaultHandler {
	return &PaymentGatewayVaultHandler{
		vaultService: service,
	}
}

func (h *PaymentGatewayVaultHandler) ListGateways(c echo.Context) error {
	gateways, err := h.vaultService.ListGateways(c.Request().Context())
	if err != nil {
		return response.InternalErrorEcho(c, "Failed to list payment gateways: "+err.Error())
	}
	return response.SuccessEcho(c, http.StatusOK, gateways)
}

func (h *PaymentGatewayVaultHandler) GetGateway(c echo.Context) error {
	provider := c.Param("provider")
	gateway, err := h.vaultService.GetGateway(c.Request().Context(), provider)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentGatewayNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}
	return response.SuccessEcho(c, http.StatusOK, gateway)
}

func (h *PaymentGatewayVaultHandler) UpdateGateway(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	provider := c.Param("provider")
	var payload application.UpdateGatewayPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid gateway configuration payload")
	}

	updated, err := h.vaultService.UpdateGateway(c.Request().Context(), principal, provider, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedPlatformAdmin) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrOptimisticLockingConflict) {
			return response.ErrorEcho(c, http.StatusConflict, "CONFLICT", err.Error())
		}
		if errors.Is(err, domain.ErrInvalidGatewayConfig) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, updated)
}
