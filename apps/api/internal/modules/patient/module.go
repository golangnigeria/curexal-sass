package patient

import (
	patientHandler "github.com/golangnigeria/curexal/internal/modules/patient/handler"
	patientRepo "github.com/golangnigeria/curexal/internal/modules/patient/repository"
	patientService "github.com/golangnigeria/curexal/internal/modules/patient/service"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Repo    *patientRepo.PatientRepository
	Service *patientService.PatientService
	Handler *patientHandler.PatientHandler
}

func NewModule(s *server.Server, userRepo patientService.UserIdentityRepo) *Module {
	repo := patientRepo.NewPatientRepository(s)
	svc := patientService.NewPatientService(s, repo, userRepo)
	hnd := patientHandler.NewPatientHandler(svc)

	return &Module{
		Repo:    repo,
		Service: svc,
		Handler: hnd,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group) {
	if m.Handler != nil {
		patientGroup := apiGroup.Group("/patient", middleware.PatientGuard)
		patientGroup.GET("/profile", m.Handler.GetProfile)
		patientGroup.PUT("/profile", m.Handler.UpdateProfile)
		patientGroup.GET("/results", m.Handler.GetResults)
		patientGroup.GET("/orders", m.Handler.GetOrders)
		patientGroup.GET("/appointments", m.Handler.GetAppointments)
	}
}

