package facility_config

import (
	auditPostgres "github.com/golangnigeria/curexal/internal/modules/audit/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/modules/facility_config/api"
	"github.com/golangnigeria/curexal/internal/modules/facility_config/application"
	"github.com/golangnigeria/curexal/internal/modules/facility_config/domain"
	"github.com/golangnigeria/curexal/internal/modules/facility_config/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Repo       domain.FacilityConfigRepository
	AppService *application.FacilityConfigApplicationService
	Handler    *api.FacilityConfigHandler
}

func NewModule(s *server.Server) *Module {
	repo := postgres.NewFacilityConfigRepository(s)
	appService := application.NewFacilityConfigApplicationService(s, repo)
	if s != nil && s.DB != nil {
		auditRepo := auditPostgres.NewAuditRepository(s)
		appService.SetAuditRepository(auditRepo)
	}
	handler := api.NewFacilityConfigHandler(s, appService)

	return &Module{
		Repo:       repo,
		AppService: appService,
		Handler:    handler,
	}
}

func (m *Module) RegisterRoutes(pltGroup *echo.Group) {
	if pltGroup != nil && m.Handler != nil {
		pltGroup.GET("/facility-types", m.Handler.GetActiveFacilityTypes)
		pltGroup.POST("/facility-types", m.Handler.CreateFacilityType)
		pltGroup.PUT("/facility-types/:typeId", m.Handler.UpdateFacilityType)
		pltGroup.GET("/facility-types/:typeId/form", m.Handler.GetRegistrationForm)
		pltGroup.GET("/facility-types/:typeId/navigation", m.Handler.GetNavigationMenu)
		pltGroup.GET("/facility-types/:typeId/setup-steps", m.Handler.GetSetupSteps)
		pltGroup.GET("/facility-types/:typeId/dashboard", m.Handler.GetDashboard)
	}
}
