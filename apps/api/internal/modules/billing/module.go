package billing

import (
	"github.com/golangnigeria/curexal/internal/kernel/server"
	auditPostgres "github.com/golangnigeria/curexal/internal/modules/audit/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/modules/billing/api"
	"github.com/golangnigeria/curexal/internal/modules/billing/application"
	billingPostgres "github.com/golangnigeria/curexal/internal/modules/billing/infrastructure/postgres"
	"github.com/labstack/echo/v4"
)

type Module struct {
	AppService          *application.BillingApplicationService
	MarketplaceHandler  *api.MarketplaceHandler
	PricingService      *application.PlatformPricingService
	PricingHandler      *api.PlatformPricingHandler
	VaultService        *application.PaymentGatewayVaultService
	GatewayVaultHandler *api.PaymentGatewayVaultHandler
	POSService          *application.POSService
	POSHandler          *api.POSHandler
}

func NewModule(s *server.Server) *Module {
	appService := application.NewBillingApplicationService(s)
	mpHandler := api.NewMarketplaceHandler(s)

	billingRepo := billingPostgres.NewPlatformBillingRepository(s)
	var auditRepo *auditPostgres.AuditRepository
	masterKey := ""
	if s != nil {
		if s.DB != nil {
			auditRepo = auditPostgres.NewAuditRepository(s)
		}
		if s.Config != nil {
			masterKey = s.Config.Auth.SecretKey
		}
	}

	pricingService := application.NewPlatformPricingService(billingRepo, auditRepo)
	pricingHandler := api.NewPlatformPricingHandler(pricingService)

	vaultService := application.NewPaymentGatewayVaultService(billingRepo, auditRepo, masterKey)
	vaultHandler := api.NewPaymentGatewayVaultHandler(vaultService)

	posRepo := billingPostgres.NewPOSRepository(s)
	posService := application.NewPOSService(s, posRepo)
	posHandler := api.NewPOSHandler(posService)

	return &Module{
		AppService:          appService,
		MarketplaceHandler:  mpHandler,
		PricingService:      pricingService,
		PricingHandler:      pricingHandler,
		VaultService:        vaultService,
		GatewayVaultHandler: vaultHandler,
		POSService:          posService,
		POSHandler:          posHandler,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group) {
	if m.MarketplaceHandler != nil {
		apiGroup.GET("/marketplace/modules", m.MarketplaceHandler.ListModules)
		apiGroup.POST("/marketplace/subscribe", m.MarketplaceHandler.Subscribe)
	}

	pltGroup := apiGroup.Group("/platform")
	if m.PricingHandler != nil {
		pltGroup.GET("/pricing", m.PricingHandler.ListPricingRules)
		pltGroup.PUT("/pricing", m.PricingHandler.UpdatePricingRule)
	}
	if m.GatewayVaultHandler != nil {
		pltGroup.GET("/payment-gateways", m.GatewayVaultHandler.ListGateways)
		pltGroup.GET("/payment-gateways/:provider", m.GatewayVaultHandler.GetGateway)
		pltGroup.PUT("/payment-gateways/:provider", m.GatewayVaultHandler.UpdateGateway)
	}

	// Clinic Cashier POS Settlement routes
	if m.POSHandler != nil {
		billingGroup := apiGroup.Group("/billing")
		billingGroup.GET("/invoices", m.POSHandler.ListInvoices)
		billingGroup.GET("/invoices/:id", m.POSHandler.GetInvoiceByID)
		billingGroup.POST("/invoices/:id/payments", m.POSHandler.ProcessPayment)
		billingGroup.GET("/receipts/:receiptNumber", m.POSHandler.GetPaymentReceipt)
	}
}
