package platform

import (
	auditPostgres "github.com/golangnigeria/curexal/internal/modules/audit/infrastructure/postgres"
	orgRepo "github.com/golangnigeria/curexal/internal/modules/organization/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/modules/platform/api"
	"github.com/golangnigeria/curexal/internal/modules/platform/application"
	platformPostgres "github.com/golangnigeria/curexal/internal/modules/platform/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type Module struct {
	AppService         *application.PlatformApplicationService
	ConfigService      *application.PlatformConfigService
	LaunchGateService  *application.LaunchGateService
	BootstrapHandler   *api.BootstrapHandler
	DiagnosticsHandler *api.DiagnosticsHandler
	ConfigHandler      *api.PlatformConfigHandler
	LaunchGateHandler  *api.LaunchGateHandler
}

func NewModule(s *server.Server) *Module {
	appService := application.NewPlatformApplicationService()

	oRepo := orgRepo.NewOrganizationRepository(s)
	tRepo := orgRepo.NewTenantRepository(s)
	sRepo := orgRepo.NewSubscriptionRepository(s)

	bootstrapHandler := api.NewBootstrapHandler(s, oRepo, tRepo, sRepo)
	diagnosticsHandler := api.NewDiagnosticsHandler(s)

	configRepo := platformPostgres.NewPlatformConfigRepository(s)
	launchGateRepo := platformPostgres.NewLaunchGateRepository(s)
	auditRepo := auditPostgres.NewAuditRepository(s)

	configService := application.NewPlatformConfigService(configRepo, auditRepo)
	launchGateService := application.NewLaunchGateService(launchGateRepo, auditRepo)

	configHandler := api.NewPlatformConfigHandler(configService)
	launchGateHandler := api.NewLaunchGateHandler(launchGateService)

	return &Module{
		AppService:         appService,
		ConfigService:      configService,
		LaunchGateService:  launchGateService,
		BootstrapHandler:   bootstrapHandler,
		DiagnosticsHandler: diagnosticsHandler,
		ConfigHandler:      configHandler,
		LaunchGateHandler:  launchGateHandler,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group, pltGroup *echo.Group) {
	if m.BootstrapHandler != nil {
		apiGroup.GET("/bootstrap", m.BootstrapHandler.GetBootstrap)
	}
	if pltGroup != nil {
		if m.DiagnosticsHandler != nil {
			pltGroup.GET("/diagnostics", m.DiagnosticsHandler.GetDiagnostics)
		}
		if m.ConfigHandler != nil {
			pltGroup.GET("/config", m.ConfigHandler.GetPlatformConfig)
			pltGroup.PUT("/config", m.ConfigHandler.UpdatePlatformConfig)
			pltGroup.PATCH("/config", m.ConfigHandler.UpdatePlatformConfig)
			pltGroup.GET("/security-policy", m.ConfigHandler.GetSecurityPolicy)
			pltGroup.PUT("/security-policy", m.ConfigHandler.UpdateSecurityPolicy)
			pltGroup.PATCH("/security-policy", m.ConfigHandler.UpdateSecurityPolicy)
		}
		if m.LaunchGateHandler != nil {
			pltGroup.GET("/launch-gate/status", m.LaunchGateHandler.GetStatus)
			pltGroup.POST("/launch-gate/verify", m.LaunchGateHandler.VerifyProductionReadiness, middleware.RequirePermission("platform:launch_gate:execute"))
			pltGroup.GET("/health/metrics", m.LaunchGateHandler.ListHealthMetrics)
		}
	}
}
