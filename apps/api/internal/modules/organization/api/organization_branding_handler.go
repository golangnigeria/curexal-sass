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

type OrganizationBrandingHandler struct {
	brandingService *application.OrganizationBrandingService
}

func NewOrganizationBrandingHandler(brandingService *application.OrganizationBrandingService) *OrganizationBrandingHandler {
	return &OrganizationBrandingHandler{
		brandingService: brandingService,
	}
}

func (h *OrganizationBrandingHandler) GetBranding(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	branding, err := h.brandingService.GetBranding(c.Request().Context(), principal)
	if err != nil {
		if errors.Is(err, domain.ErrOrganizationNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to get organization branding: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, branding)
}

func (h *OrganizationBrandingHandler) UpdateBranding(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.UpdateBrandingPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid branding update payload")
	}

	updated, err := h.brandingService.UpdateBranding(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrDuplicateCustomDomain) {
			return response.ErrorEcho(c, http.StatusConflict, "DUPLICATE_CUSTOM_DOMAIN", err.Error())
		}
		if errors.Is(err, domain.ErrOptimisticLockingConflict) {
			return response.ErrorEcho(c, http.StatusConflict, "CONFLICT", err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to update branding: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, updated)
}

func (h *OrganizationBrandingHandler) SaveNotificationConfig(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.SaveNotificationConfigPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid notification config payload")
	}

	saved, err := h.brandingService.SaveNotificationConfig(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrInvalidNotificationChannel) || errors.Is(err, domain.ErrInvalidNotificationProvider) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to save notification config: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, saved)
}

func (h *OrganizationBrandingHandler) ListNotificationConfigs(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	configs, err := h.brandingService.ListNotificationConfigs(c.Request().Context(), principal)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to list notification configs: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, configs)
}

func (h *OrganizationBrandingHandler) SaveNotificationTemplate(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	templateKey := c.Param("key")
	if templateKey == "" {
		return response.BadRequestEcho(c, "template key parameter is required")
	}

	var payload domain.SaveNotificationTemplatePayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid notification template payload")
	}

	saved, err := h.brandingService.SaveNotificationTemplate(c.Request().Context(), principal, templateKey, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrInvalidNotificationChannel) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to save notification template: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, saved)
}

func (h *OrganizationBrandingHandler) ListNotificationTemplates(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	templates, err := h.brandingService.ListNotificationTemplates(c.Request().Context(), principal)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to list notification templates: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, templates)
}

func (h *OrganizationBrandingHandler) ListUserNotifications(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	notifs, err := h.brandingService.ListUserNotifications(c.Request().Context(), principal, 50)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to list in-app notifications: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, notifs)
}

func (h *OrganizationBrandingHandler) MarkNotificationRead(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	notifIDParam := c.Param("id")
	notifUUID, errParse := uuid.Parse(notifIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid notification ID format")
	}

	err := h.brandingService.MarkNotificationRead(c.Request().Context(), principal, notifUUID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotificationNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]string{"message": "In-app notification marked as read"})
}

func (h *OrganizationBrandingHandler) MarkAllNotificationsRead(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	err := h.brandingService.MarkAllNotificationsRead(c.Request().Context(), principal)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]string{"message": "All in-app notifications marked as read"})
}

func (h *OrganizationBrandingHandler) ListNotificationDeliveries(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	deliveries, err := h.brandingService.ListNotificationDeliveries(c.Request().Context(), principal, 50)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to list notification deliveries: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, deliveries)
}
