package api

import (
	"net/http"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/subscription/application"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/golangnigeria/curexal/internal/shared/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type EntitlementHandler struct {
	service *application.EntitlementService
}

func NewEntitlementHandler(service *application.EntitlementService) *EntitlementHandler {
	return &EntitlementHandler{
		service: service,
	}
}

type GrantCapabilityPayload struct {
	CapabilityCode string     `json:"capabilityCode" validate:"required"`
	Source         string     `json:"source"`
	ExpiresAt      *time.Time `json:"expiresAt"`
}

type TrialCapabilityPayload struct {
	DurationDays int `json:"durationDays"`
}

func (h *EntitlementHandler) GetEffectiveCapabilities(c echo.Context) error {
	orgIDParam := c.Param("id")
	orgID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid organization ID format")
	}

	caps, err := h.service.GetEffectiveCapabilities(c.Request().Context(), orgID)
	if err != nil {
		return err
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]interface{}{
		"organizationId": orgID.String(),
		"capabilities":   caps,
	})
}

func (h *EntitlementHandler) GetOrganizationEntitlements(c echo.Context) error {
	orgIDParam := c.Param("id")
	orgID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid organization ID format")
	}

	entitlements, err := h.service.GetOrganizationEntitlements(c.Request().Context(), orgID)
	if err != nil {
		return err
	}

	return response.SuccessEcho(c, http.StatusOK, entitlements)
}

func (h *EntitlementHandler) GrantCapabilityAddOn(c echo.Context) error {
	orgIDParam := c.Param("id")
	orgID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid organization ID format")
	}

	var payload GrantCapabilityPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("invalid request payload")
	}
	if payload.CapabilityCode == "" {
		return errs.NewBadRequestError("capabilityCode is required")
	}

	actorID := middleware.GetUserID(c)
	errGrant := h.service.GrantCapabilityAddOn(c.Request().Context(), actorID, orgID, payload.CapabilityCode, payload.Source, payload.ExpiresAt)
	if errGrant != nil {
		return errGrant
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]string{
		"message": "Capability add-on entitlement granted successfully",
	})
}

func (h *EntitlementHandler) RevokeCapabilityAddOn(c echo.Context) error {
	orgIDParam := c.Param("id")
	orgID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid organization ID format")
	}

	capCode := c.Param("capability")
	if capCode == "" {
		return errs.NewBadRequestError("capability code parameter required")
	}

	actorID := middleware.GetUserID(c)
	errRevoke := h.service.RevokeCapabilityAddOn(c.Request().Context(), actorID, orgID, capCode)
	if errRevoke != nil {
		return errRevoke
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]string{
		"message": "Capability entitlement revoked successfully",
	})
}

func (h *EntitlementHandler) StartTrialCapability(c echo.Context) error {
	orgIDParam := c.Param("id")
	orgID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid organization ID format")
	}

	capCode := c.Param("capability")
	if capCode == "" {
		return errs.NewBadRequestError("capability code parameter required")
	}

	var payload TrialCapabilityPayload
	_ = c.Bind(&payload)

	actorID := middleware.GetUserID(c)
	errTrial := h.service.StartTrialCapability(c.Request().Context(), actorID, orgID, capCode, payload.DurationDays)
	if errTrial != nil {
		return errTrial
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]string{
		"message": "Trial capability entitlement activated successfully",
	})
}

func (h *EntitlementHandler) GetEntitlementTrace(c echo.Context) error {
	orgIDParam := c.Param("id")
	orgID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid organization ID format")
	}

	capCode := c.QueryParam("capability")
	if capCode == "" {
		return errs.NewBadRequestError("query parameter 'capability' is required")
	}

	trace, err := h.service.GetEntitlementTrace(c.Request().Context(), orgID, capCode)
	if err != nil {
		return err
	}

	return response.SuccessEcho(c, http.StatusOK, trace)
}

type PurchaseCapabilityPayload struct {
	CapabilityCode string `json:"capabilityCode" validate:"required"`
	Currency       string `json:"currency"`
	BillingCycle   string `json:"billingCycle"`
}

func (h *EntitlementHandler) GetCapabilityCatalog(c echo.Context) error {
	orgIDParam := c.Param("id")
	var orgID uuid.UUID
	if orgIDParam != "" {
		parsed, errP := uuid.Parse(orgIDParam)
		if errP == nil {
			orgID = parsed
		}
	}

	catalog, err := h.service.GetCapabilityCatalog(c.Request().Context(), orgID)
	if err != nil {
		return err
	}

	return response.SuccessEcho(c, http.StatusOK, catalog)
}

func (h *EntitlementHandler) PurchaseCapabilityAddOn(c echo.Context) error {
	orgIDParam := c.Param("id")
	orgID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid organization ID format")
	}

	var payload PurchaseCapabilityPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("invalid request payload")
	}
	if payload.CapabilityCode == "" {
		return errs.NewBadRequestError("capabilityCode is required")
	}

	actorID := middleware.GetUserID(c)
	sub, errSub := h.service.PurchaseCapabilityAddOn(c.Request().Context(), actorID, orgID, payload.CapabilityCode, payload.Currency, payload.BillingCycle)
	if errSub != nil {
		return errSub
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]interface{}{
		"message":                "Commercial capability subscription activated successfully",
		"capabilitySubscription": sub,
	})
}
