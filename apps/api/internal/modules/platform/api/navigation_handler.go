package api

import (
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/platform/application"
	"github.com/golangnigeria/curexal/internal/modules/platform/domain"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type NavigationHandler struct {
	navService *application.NavigationService
}

func NewNavigationHandler(navService *application.NavigationService) *NavigationHandler {
	return &NavigationHandler{
		navService: navService,
	}
}

// GetNavigation serves the canonical GET /api/v1/navigation endpoint
func (h *NavigationHandler) GetNavigation(c echo.Context) error {
	principal := middleware.GetPrincipal(c)

	reqHost := c.Request().Header.Get("X-Forwarded-Host")
	if reqHost == "" {
		reqHost = c.Request().Host
	}
	if qHost := c.QueryParam("host"); qHost != "" {
		reqHost = qHost
	}

	branchSlug := c.QueryParam("branch")
	if branchSlug == "" {
		branchSlug = c.Request().Header.Get("X-Branch-Slug")
	}

	reqScope := c.QueryParam("scope")
	if reqScope == "" {
		reqScope = c.Request().Header.Get("X-Context-Scope")
	}

	res, err := h.navService.GetNavigation(c.Request().Context(), principal, reqHost, branchSlug, reqScope)
	if err != nil {
		if err == domain.ErrUnauthorizedScope {
			return echo.NewHTTPError(http.StatusForbidden, "Unauthorized navigation context scope")
		}
		if err == domain.ErrInvalidScope {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid navigation context scope parameter")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to resolve navigation state: "+err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
