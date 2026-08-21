package authorization

import (
	"github.com/golangnigeria/curexal/internal/modules/authorization/handler"
	"github.com/golangnigeria/curexal/internal/modules/authorization/service"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Enforcer     *service.CasbinEnforcer
	AuthzHandler *handler.AuthzHandler
}

func NewModule(s *server.Server) *Module {
	enforcer := service.NewCasbinEnforcer(s)
	authzHandler := handler.NewAuthzHandler(s, enforcer)
	return &Module{
		Enforcer:     enforcer,
		AuthzHandler: authzHandler,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group) {
	if m.AuthzHandler != nil {
		apiGroup.POST("/authorization/enforce", m.AuthzHandler.EvaluatePermission)
		apiGroup.GET("/authorization/permissions", m.AuthzHandler.GetUserPermissions)
	}
}

