package operations

import (
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/operations/handler"
	"github.com/golangnigeria/curexal/internal/modules/operations/repository"
	"github.com/golangnigeria/curexal/internal/modules/operations/service"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Repo    *repository.AppointmentRepository
	Service *service.AppointmentService
	Handler *handler.AppointmentHandler
}

func NewModule(s *server.Server) *Module {
	repo := repository.NewAppointmentRepository(s)
	svc := service.NewAppointmentService(s, repo)
	hnd := handler.NewAppointmentHandler(s, svc)

	return &Module{
		Repo:    repo,
		Service: svc,
		Handler: hnd,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group) {
	if m.Handler != nil {
		appointmentsGroup := apiGroup.Group("/appointments")
		appointmentsGroup.POST("", m.Handler.CreateAppointment)
		appointmentsGroup.GET("", m.Handler.ListAppointments)
		appointmentsGroup.GET("/:id", m.Handler.GetAppointmentByID)
		appointmentsGroup.PUT("/:id/status", m.Handler.UpdateAppointmentStatus)
	}
}
