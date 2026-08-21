package clinical

import (
	"github.com/golangnigeria/curexal/internal/modules/clinical/api"
	"github.com/golangnigeria/curexal/internal/modules/clinical/application"
	"github.com/golangnigeria/curexal/internal/modules/clinical/domain"
	"github.com/golangnigeria/curexal/internal/modules/clinical/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type Module struct {
	CatalogRepo         domain.CatalogRepository
	AppService          *application.ClinicalApplicationService
	CatalogHandler      *api.CatalogHandler
	PatientVisitHandler *api.PatientVisitHandler
	LimsHandler         *api.LimsHandler
}

func NewModule(s *server.Server) *Module {
	catalogRepo := postgres.NewCatalogRepository(s)
	appService := application.NewClinicalApplicationService(s, catalogRepo)
	catalogHandler := api.NewCatalogHandler(s, appService)
	patientVisitHandler := api.NewPatientVisitHandler(s)
	limsHandler := api.NewLimsHandler(s)

	return &Module{
		CatalogRepo:         catalogRepo,
		AppService:          appService,
		CatalogHandler:      catalogHandler,
		PatientVisitHandler: patientVisitHandler,
		LimsHandler:         limsHandler,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group, entitlementSvc middleware.CapabilityChecker) {
	if m.PatientVisitHandler != nil {
		apiGroup.POST("/clinical/patient-visits", m.PatientVisitHandler.RegisterPatientVisit, middleware.RequireCapability("clinical.basic", entitlementSvc))
	}
	if m.LimsHandler != nil {
		apiGroup.POST("/lims/orders", m.LimsHandler.CreateOrder, middleware.RequireCapability("laboratory.basic", entitlementSvc))
		apiGroup.POST("/lims/specimens/accession", m.LimsHandler.AccessionSpecimen, middleware.RequireCapability("laboratory.basic", entitlementSvc))
		apiGroup.POST("/lims/results", m.LimsHandler.EnterResults, middleware.RequireCapability("laboratory.basic", entitlementSvc))
		apiGroup.POST("/lims/authorizations", m.LimsHandler.AuthorizeOrder, middleware.RequireCapability("laboratory.basic", entitlementSvc))
	}
}

