package orchestration

import (
	"github.com/golangnigeria/curexal/internal/modules/orchestration/handler"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/repository"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/service"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Repo    *repository.CareRequestRepository
	Service *service.CareRequestService
	Handler *handler.CareRequestHandler
}

func NewModule(s *server.Server) *Module {
	repo := repository.NewCareRequestRepository(s)
	svc := service.NewCareRequestService(s, repo)
	hnd := handler.NewCareRequestHandler(s, svc)

	return &Module{
		Repo:    repo,
		Service: svc,
		Handler: hnd,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group) {
	if m.Handler != nil {
		// Staff & Operational Care Requests
		careRequestsGroup := apiGroup.Group("/care-requests")
		careRequestsGroup.POST("", m.Handler.CreateCareRequest)
		careRequestsGroup.GET("", m.Handler.ListCareRequests)
		careRequestsGroup.GET("/:id", m.Handler.GetCareRequestByID)
		careRequestsGroup.POST("/:id/triage", m.Handler.SubmitTriage)
		careRequestsGroup.GET("/:id/match-provider", m.Handler.MatchProviders)
		careRequestsGroup.POST("/:id/assign-provider", m.Handler.AssignProvider)

		// Patient Portal Workspace Routes
		portalGroup := apiGroup.Group("/portal")
		portalGroup.GET("/my-care-requests", m.Handler.GetMyCareRequests)
		portalGroup.GET("/my-care-journey", m.Handler.GetMyCareJourney)
	}
}
