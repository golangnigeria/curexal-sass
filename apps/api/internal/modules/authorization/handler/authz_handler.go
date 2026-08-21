package handler

import (
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/authorization/model"
	"github.com/golangnigeria/curexal/internal/modules/authorization/service"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type AuthzHandler struct {
	server   *server.Server
	enforcer *service.CasbinEnforcer
}

func NewAuthzHandler(s *server.Server, e *service.CasbinEnforcer) *AuthzHandler {
	return &AuthzHandler{
		server:   s,
		enforcer: e,
	}
}

// EvaluatePermission handles explicit permission enforcement requests.
func (h *AuthzHandler) EvaluatePermission(c echo.Context) error {
	var req model.EnforceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if req.Subject == "" {
		req.Subject = middleware.GetUserID(c)
	}
	if req.Tenant == "" {
		req.Tenant = middleware.GetActiveTenantID(c)
	}

	allowed, reason := h.enforcer.Enforce(c.Request().Context(), req.Subject, req.Tenant, req.Resource, req.Action)

	return c.JSON(http.StatusOK, model.EnforceResponse{
		Allowed: allowed,
		Reason:  reason,
	})
}

// GetUserPermissions returns current permissions in active tenant context.
func (h *AuthzHandler) GetUserPermissions(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetActiveTenantID(c)

	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	perms := h.enforcer.ListUserPermissions(ctx, userID, tenantID)

	return c.JSON(http.StatusOK, model.UserPermissionsResponse{
		Subject:     userID,
		Tenant:      tenantID,
		Permissions: perms,
	})
}
