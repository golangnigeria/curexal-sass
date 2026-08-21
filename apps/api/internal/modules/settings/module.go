package settings

import (
	"github.com/golangnigeria/curexal/internal/modules/settings/api"
	"github.com/golangnigeria/curexal/internal/modules/settings/application"
	"github.com/golangnigeria/curexal/internal/modules/settings/domain"
	"github.com/golangnigeria/curexal/internal/modules/settings/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Repo       domain.BranchSettingsRepository
	AppService *application.SettingsApplicationService
	Handler    *api.SettingsHandler
}

func NewModule(s *server.Server) *Module {
	repo := postgres.NewBranchSettingsRepository(s)
	appService := application.NewSettingsApplicationService(s, repo)
	handler := api.NewSettingsHandler(s, appService)

	return &Module{
		Repo:       repo,
		AppService: appService,
		Handler:    handler,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group) {
	if m.Handler != nil {
		apiGroup.GET("/settings/branch", m.Handler.GetBranchSettings)
		apiGroup.PUT("/settings/branch/:section", m.Handler.UpdateBranchSettingsSection)
		apiGroup.POST("/settings/branch/reset", m.Handler.ResetBranchSettings)
	}
}

