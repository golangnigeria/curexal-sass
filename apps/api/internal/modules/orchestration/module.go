package orchestration

import (
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/handler"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/repository"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/service"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Repo            *repository.CareRequestRepository
	Service         *service.CareRequestService
	Handler         *handler.CareRequestHandler
	ProviderRepo    *repository.ProviderProfileRepository
	ProviderService *service.ProviderProfileService
	ProviderHandler *handler.ProviderProfileHandler
}

func NewModule(s *server.Server) *Module {
	repo := repository.NewCareRequestRepository(s)
	svc := service.NewCareRequestService(s, repo)
	hnd := handler.NewCareRequestHandler(s, svc)

	providerRepo := repository.NewProviderProfileRepository(s)
	providerSvc := service.NewProviderProfileService(s, providerRepo)
	providerHnd := handler.NewProviderProfileHandler(s, providerSvc)

	return &Module{
		Repo:            repo,
		Service:         svc,
		Handler:         hnd,
		ProviderRepo:    providerRepo,
		ProviderService: providerSvc,
		ProviderHandler: providerHnd,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group) {
	if m.ProviderHandler != nil {
		// Day 2 Feature #4: Provider Profile Management (Secured)
		providersGroup := apiGroup.Group("/providers/profiles", middleware.RequireAuth())
		providersGroup.GET("", m.ProviderHandler.ListProviders)
		providersGroup.GET("/:id", m.ProviderHandler.GetProviderByID)
		providersGroup.POST("", m.ProviderHandler.CreateProvider, middleware.RequirePermission("users:write"))
		providersGroup.PUT("/:id/status", m.ProviderHandler.UpdateStatus, middleware.RequirePermission("users:write"))
	}

	if m.Handler != nil {
		// Care Requests (Canonical /care-requests and /orchestration/requests - Secured)
		registerCareRequestRoutes := func(g *echo.Group) {
			g.Use(middleware.RequireAuth())
			g.POST("", m.Handler.CreateCareRequest)
			g.GET("", m.Handler.ListCareRequests)
			g.GET("/:id", m.Handler.GetCareRequestByID)
			g.POST("/:id/triage", m.Handler.SubmitTriage)
			g.GET("/:id/match-provider", m.Handler.MatchProviders)
			g.POST("/:id/assign-provider", m.Handler.AssignProvider)
		}

		registerCareRequestRoutes(apiGroup.Group("/care-requests"))
		registerCareRequestRoutes(apiGroup.Group("/orchestration/requests"))

		// Patient Portal Workspace Routes
		portalGroup := apiGroup.Group("/portal", middleware.RequireAuth())
		portalGroup.GET("/my-care-requests", m.Handler.GetMyCareRequests)
		portalGroup.GET("/my-care-journey", m.Handler.GetMyCareJourney)
	}
}
