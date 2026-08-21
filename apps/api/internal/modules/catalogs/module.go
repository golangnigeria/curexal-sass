package catalogs

import (
	auditPostgres "github.com/golangnigeria/curexal/internal/modules/audit/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/modules/catalogs/api"
	"github.com/golangnigeria/curexal/internal/modules/catalogs/application"
	"github.com/golangnigeria/curexal/internal/modules/catalogs/domain"
	catalogPostgres "github.com/golangnigeria/curexal/internal/modules/catalogs/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Repo           domain.CatalogRepository
	CatalogService *application.MasterCatalogService
	CatalogHandler *api.MasterCatalogHandler
}

func NewModule(s *server.Server) *Module {
	repo := catalogPostgres.NewCatalogRepository(s)
	var auditRepo *auditPostgres.AuditRepository
	if s != nil && s.DB != nil {
		auditRepo = auditPostgres.NewAuditRepository(s)
	}

	service := application.NewMasterCatalogService(repo, auditRepo)
	handler := api.NewMasterCatalogHandler(service)

	return &Module{
		Repo:           repo,
		CatalogService: service,
		CatalogHandler: handler,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group) {
	if m.CatalogHandler != nil {
		// Legacy compatibility catalog endpoints
		apiGroup.GET("/catalogs", m.CatalogHandler.GetMasterCatalogs)
		apiGroup.GET("/catalogs/icd10", m.CatalogHandler.SearchICD10)

		// Platform Master Reference Catalog Endpoints (EPIC-010)
		pltGroup := apiGroup.Group("/platform")
		pltGroup.GET("/catalogs/:domain", m.CatalogHandler.ListItems)
		pltGroup.GET("/catalogs/:domain/search", m.CatalogHandler.SearchItems)
		pltGroup.POST("/catalogs/:domain", m.CatalogHandler.CreateItem)
		pltGroup.PUT("/catalogs/:domain/:id", m.CatalogHandler.UpdateItem)
	}
}
