package api

import (
	"net/http"

	orgDomain "github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/modules/platform/application"
	subApp "github.com/golangnigeria/curexal/internal/modules/subscription/application"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type BootstrapHandler struct {
	server           *server.Server
	bootstrapBuilder *application.BootstrapBuilder
}

func NewBootstrapHandler(s *server.Server, orgRepo orgDomain.OrganizationRepository, tenantRepo orgDomain.TenantRepository, subRepo orgDomain.SubscriptionRepository) *BootstrapHandler {
	builder := application.NewBootstrapBuilder(orgRepo, tenantRepo, subRepo)
	if s != nil && s.DB != nil && s.DB.Pool != nil {
		builder.SetDBPool(s.DB.Pool)
	}
	return &BootstrapHandler{
		server:           s,
		bootstrapBuilder: builder,
	}
}

func (h *BootstrapHandler) SetEntitlementService(svc *subApp.EntitlementService) {
	if h.bootstrapBuilder != nil {
		h.bootstrapBuilder.SetEntitlementService(svc)
	}
}


func (h *BootstrapHandler) GetBootstrap(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	contract, err := h.bootstrapBuilder.BuildBootstrap(c.Request().Context(), principal)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to assemble bootstrap contract")
	}

	return c.JSON(http.StatusOK, contract)
}

type SwitchContextPayload struct {
	TargetContext string `json:"targetContext" validate:"required,oneof=platform organization workspace"`
	TargetID      string `json:"targetId,omitempty"`
}

func (h *BootstrapHandler) SwitchContext(c echo.Context) error {
	var payload SwitchContextPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid context switch payload")
	}

	// 204 No Content - frontend re-fetches GET /api/v1/bootstrap immediately
	return c.NoContent(http.StatusNoContent)
}
