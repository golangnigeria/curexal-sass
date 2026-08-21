package identity

import (
	"context"

	"github.com/golangnigeria/curexal/internal/modules/identity/handler"
	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/identity/repository"
	"github.com/golangnigeria/curexal/internal/modules/identity/service"
	patientModel "github.com/golangnigeria/curexal/internal/modules/patient/model"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type Module struct {
	UserRepo        *repository.UserRepository
	AuthService     *service.AuthService
	InviteService   *service.InviteService
	AuthHandler     *handler.AuthHandler
	InviteHandler   *handler.InviteHandler
	UserRoleHandler *handler.UserRoleHandler
}

type PatientServiceIntf interface {
	RegisterPatient(ctx context.Context, payload patientModel.RegisterPatientPayload, origin string) error
	LoadPatientContext(ctx context.Context, userID string) (*model.PatientContext, error)
}

type PatientRepoIntf interface {
	ProfileExists(ctx context.Context, userID string) (bool, string, error)
}

func NewModule(s *server.Server, patientSvc PatientServiceIntf, patientRepo PatientRepoIntf, tenantLookup model.TenantLookup) *Module {
	userRepo := repository.NewUserRepository(s)
	authService := service.NewAuthService(s)
	inviteService := service.NewInviteService(s, authService, tenantLookup)

	return &Module{
		UserRepo:        userRepo,
		AuthService:     authService,
		InviteService:   inviteService,
		AuthHandler:     handler.NewAuthHandler(s, authService, userRepo, patientSvc, patientRepo),
		InviteHandler:   handler.NewInviteHandler(s, inviteService, authService),
		UserRoleHandler: handler.NewUserRoleHandler(s, userRepo),
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group, orgGroup *echo.Group, wspGroup *echo.Group) {
	if m.UserRoleHandler != nil {
		apiGroup.POST("/context/switch", m.UserRoleHandler.SwitchTenant)
	}

	if m.AuthHandler != nil {
		authGroup := apiGroup.Group("/auth")
		authGroup.POST("/sign-in", m.AuthHandler.SignIn)
		authGroup.POST("/sign-up", m.AuthHandler.SignUp)
		authGroup.POST("/sign-out", m.AuthHandler.SignOut)
		authGroup.POST("/switch-context", m.AuthHandler.SwitchContext)
		authGroup.POST("/verify-email", m.AuthHandler.VerifyEmail)
		authGroup.GET("/csrf", m.AuthHandler.GetCSRFToken)
		authGroup.POST("/impersonate", m.AuthHandler.ImpersonateTenant)
		if m.InviteHandler != nil {
			authGroup.POST("/accept-invite", m.InviteHandler.AcceptInvite)
		}
		authGroup.POST("/request-password", m.AuthHandler.RequestPassword)
		authGroup.POST("/forgot-password", m.AuthHandler.ForgotPassword)
		authGroup.POST("/set-password", m.AuthHandler.SetPassword)
		authGroup.POST("/resend-verification", m.AuthHandler.ResendVerification)
		authGroup.GET("/verify-email-change", m.AuthHandler.VerifyEmailChangeGet)
		authGroup.POST("/verify-email-change", m.AuthHandler.VerifyEmailChange)

		apiGroup.POST("/users/me/request-email-change", m.AuthHandler.RequestEmailChange)
	}

	if m.UserRoleHandler != nil {
		apiGroup.GET("/users/me", m.UserRoleHandler.GetMe)
		apiGroup.GET("/users/me/profile", m.UserRoleHandler.GetUserProfile)
		apiGroup.PATCH("/users/me/profile", m.UserRoleHandler.UpdateUserProfile)
		apiGroup.GET("/users/me/employment", m.UserRoleHandler.GetTenantEmployment)
		apiGroup.PUT("/users/me/employment", m.UserRoleHandler.UpdateTenantEmployment)
		apiGroup.PUT("/users/me/password", m.UserRoleHandler.ChangePassword)
		apiGroup.GET("/users/me/professional", m.UserRoleHandler.GetProfessionalProfiles)
		apiGroup.POST("/users/me/professional", m.UserRoleHandler.CreateProfessionalProfile)
		apiGroup.GET("/users/me/signatures", m.UserRoleHandler.GetSignatures)
		apiGroup.POST("/users/me/signatures", m.UserRoleHandler.CreateSignature)
		apiGroup.GET("/users", m.UserRoleHandler.GetUsers, middleware.RequirePermission("users:read"))
		apiGroup.GET("/roles", m.UserRoleHandler.GetRoles)

		if orgGroup != nil {
			orgGroup.GET("/members", m.UserRoleHandler.GetUsers)
			orgGroup.POST("/members", m.UserRoleHandler.CreateWorkspaceMember)
			orgGroup.GET("/roles", m.UserRoleHandler.GetRoles)
		}

		if wspGroup != nil {
			wspGroup.GET("/patients", m.UserRoleHandler.GetUsers)
		}
	}
}

