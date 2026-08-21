package subscription

import (
	billingDomain "github.com/golangnigeria/curexal/internal/modules/billing/domain"
	billingInfra "github.com/golangnigeria/curexal/internal/modules/billing/infrastructure"
	"github.com/golangnigeria/curexal/internal/modules/subscription/api"
	"github.com/golangnigeria/curexal/internal/modules/subscription/application"
	"github.com/golangnigeria/curexal/internal/modules/subscription/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Service           *application.EntitlementService
	Handler           *api.EntitlementHandler
	CommercialService *application.CommercialService
	CommercialHandler *api.CommercialHandler
	ProviderRegistry  *billingDomain.PaymentProviderRegistry
}

func NewModule(s *server.Server) *Module {
	repo := postgres.NewEntitlementRepository(s)
	svc := application.NewEntitlementService(s, repo)
	hnd := api.NewEntitlementHandler(svc)

	commRepo := postgres.NewCommercialRepository(s)
	providerReg := billingDomain.NewPaymentProviderRegistry()
	providerReg.Register(billingInfra.NewMockPaymentProvider())
	providerReg.Register(billingInfra.NewPaystackProvider(""))
	providerReg.Register(billingInfra.NewFlutterwaveProvider(""))
	providerReg.Register(billingInfra.NewStripeProvider(""))

	commSvc := application.NewCommercialService(s, repo, commRepo, svc, providerReg)
	commHnd := api.NewCommercialHandler(commSvc)

	return &Module{
		Service:           svc,
		Handler:           hnd,
		CommercialService: commSvc,
		CommercialHandler: commHnd,
		ProviderRegistry:  providerReg,
	}
}

func (m *Module) RegisterRoutes(e *echo.Echo, apiGroup *echo.Group, pltGroup *echo.Group) {
	if m.Handler == nil {
		return
	}

	// Organization & Marketplace Entitlements
	apiGroup.GET("/organizations/:id/capabilities", m.Handler.GetEffectiveCapabilities)
	apiGroup.GET("/organizations/:id/capabilities/trace", m.Handler.GetEntitlementTrace)
	apiGroup.GET("/organizations/:id/entitlements", m.Handler.GetOrganizationEntitlements)
	apiGroup.GET("/organizations/:id/marketplace/catalog", m.Handler.GetCapabilityCatalog)
	apiGroup.GET("/marketplace/capabilities", m.Handler.GetCapabilityCatalog)
	apiGroup.POST("/organizations/:id/marketplace/subscribe", m.Handler.PurchaseCapabilityAddOn, middleware.RequirePermission("organization:write"))

	// Commercial Order Creation (Authenticated + org:write permission)
	if m.CommercialHandler != nil {
		apiGroup.POST("/organizations/:id/marketplace/orders", m.CommercialHandler.CreateCommercialOrder, middleware.RequirePermission("organization:write"))
	}

	// Unauthenticated Cryptographically Verified Payment Webhooks (No middleware.Authenticate)
	if m.CommercialHandler != nil && e != nil {
		e.POST("/api/v1/billing/webhooks/:provider", m.CommercialHandler.HandlePaymentWebhook)
	}

	// Platform Admin Capability Operations
	if pltGroup != nil {
		pltGroup.POST("/organizations/:id/capabilities", m.Handler.GrantCapabilityAddOn, middleware.RequirePermission("platform:admin"))
		pltGroup.DELETE("/organizations/:id/capabilities/:capability", m.Handler.RevokeCapabilityAddOn, middleware.RequirePermission("platform:admin"))
		pltGroup.POST("/organizations/:id/capabilities/:capability/trial", m.Handler.StartTrialCapability, middleware.RequirePermission("platform:admin"))
	}
}
