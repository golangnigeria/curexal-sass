package patient

import (
	"github.com/golangnigeria/curexal/internal/kernel/server"
	patientHandler "github.com/golangnigeria/curexal/internal/modules/patient/handler"
	patientRepo "github.com/golangnigeria/curexal/internal/modules/patient/repository"
	patientService "github.com/golangnigeria/curexal/internal/modules/patient/service"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Repo             *patientRepo.PatientRepository
	CanonicalRepo    *patientRepo.CanonicalPatientRepository
	Service          *patientService.PatientService
	CanonicalService *patientService.CanonicalPatientService
	MPIService       *patientService.MPIService
	Handler          *patientHandler.PatientHandler
	CanonicalHandler *patientHandler.CanonicalPatientHandler
}

func NewModule(s *server.Server, userRepo patientService.UserIdentityRepo) *Module {
	repo := patientRepo.NewPatientRepository(s)
	canonicalRepo := patientRepo.NewCanonicalPatientRepository(s)

	mpiSvc := patientService.NewMPIService(canonicalRepo)
	canonicalSvc := patientService.NewCanonicalPatientService(s, canonicalRepo, mpiSvc)
	svc := patientService.NewPatientService(s, repo, userRepo)

	hnd := patientHandler.NewPatientHandler(svc)
	canonicalHnd := patientHandler.NewCanonicalPatientHandler(s, canonicalSvc, mpiSvc)

	return &Module{
		Repo:             repo,
		CanonicalRepo:    canonicalRepo,
		Service:          svc,
		CanonicalService: canonicalSvc,
		MPIService:       mpiSvc,
		Handler:          hnd,
		CanonicalHandler: canonicalHnd,
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

	if m.CanonicalHandler != nil {
		// Canonical Patient Directory & Intake Routes
		patientsGroup := apiGroup.Group("/patients")
		patientsGroup.POST("/resolve", m.CanonicalHandler.ResolveDuplicates)
		patientsGroup.POST("", m.CanonicalHandler.RegisterPatient)
		patientsGroup.GET("", m.CanonicalHandler.ListPatients)
		patientsGroup.GET("/:id", m.CanonicalHandler.GetPatientByID)

		// Patient Portal Authentication & Onboarding Routes
		portalAuthGroup := apiGroup.Group("/portal/auth")
		portalAuthGroup.POST("/send-otp", m.CanonicalHandler.SendPortalOTP)
		portalAuthGroup.POST("/verify-otp", m.CanonicalHandler.VerifyPortalOTP)
		portalAuthGroup.POST("/patients/:id/pin", m.CanonicalHandler.SetPortalPIN)
	}
}
