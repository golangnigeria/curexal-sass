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

type OrganizationProfileHandler struct {
	setupService *application.OrganizationSetupService
}

func NewOrganizationProfileHandler(setupService *application.OrganizationSetupService) *OrganizationProfileHandler {
	return &OrganizationProfileHandler{
		setupService: setupService,
	}
}

func (h *OrganizationProfileHandler) GetProfile(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	profile, err := h.setupService.GetProfile(c.Request().Context(), principal)
	if err != nil {
		if errors.Is(err, domain.ErrOrganizationNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to load organization profile: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, profile)
}

func (h *OrganizationProfileHandler) UpdateProfile(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.UpdateOrganizationProfilePayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid organization profile update payload")
	}

	updated, err := h.setupService.UpdateProfile(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrOptimisticLockingConflict) {
			return response.ErrorEcho(c, http.StatusConflict, "CONFLICT", err.Error())
		}
		if errors.Is(err, domain.ErrInvalidOrganizationProfile) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to update organization profile: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, updated)
}

func (h *OrganizationProfileHandler) SubmitForReview(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	updated, err := h.setupService.SubmitForReview(c.Request().Context(), principal)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidSetupStateTransition) {
			return response.BadRequestEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, updated)
}

type VerifyOrganizationPayload struct {
	Approve bool   `json:"approve"`
	Reason  string `json:"reason"`
}

func (h *OrganizationProfileHandler) VerifyOrganization(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	orgIDParam := c.Param("id")
	orgUUID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid organization ID format")
	}

	var payload VerifyOrganizationPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid verification payload")
	}

	updated, err := h.setupService.VerifyOrganization(c.Request().Context(), principal, orgUUID, payload.Approve, payload.Reason)
	if err != nil {
		return response.ForbiddenEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, updated)
}
