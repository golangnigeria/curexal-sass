package encounter

import (
	"github.com/golangnigeria/curexal/internal/modules/encounter/handler"
	"github.com/golangnigeria/curexal/internal/modules/encounter/repository"
	"github.com/golangnigeria/curexal/internal/modules/encounter/service"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Repo    *repository.EncounterRepository
	Service *service.EncounterService
	Handler *handler.EncounterHandler
}

func NewModule(s *server.Server) *Module {
	repo := repository.NewEncounterRepository(s)
	svc := service.NewEncounterService(s, repo)
	hnd := handler.NewEncounterHandler(svc)

	return &Module{
		Repo:    repo,
		Service: svc,
		Handler: hnd,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group) {
	if m.Handler != nil {
		encountersGroup := apiGroup.Group("/encounters")
		encountersGroup.POST("/start", m.Handler.StartEncounter)
		encountersGroup.GET("/:id", m.Handler.GetEncounterByID)
		encountersGroup.PUT("/:id/notes", m.Handler.SaveSOAPNotes)
		encountersGroup.POST("/:id/orders", m.Handler.DispatchOrders)
		encountersGroup.POST("/:id/complete", m.Handler.CompleteEncounter)
	}
}
