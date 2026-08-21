package audit

import (
	"github.com/golangnigeria/curexal/internal/modules/audit/api"
	"github.com/golangnigeria/curexal/internal/modules/audit/application"
	"github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/audit/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type Module struct {
	AuditRepo    domain.AuditRepository
	AppService   *application.AuditApplicationService
	AuditHandler *api.AuditHandler
}

func NewModule(s *server.Server) *Module {
	auditRepo := postgres.NewAuditRepository(s)
	appService := application.NewAuditApplicationService(s, auditRepo)
	auditHandler := api.NewAuditHandler(s, appService)

	return &Module{
		AuditRepo:    auditRepo,
		AppService:   appService,
		AuditHandler: auditHandler,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group) {
	if m.AuditHandler != nil {
		apiGroup.GET("/audit-logs/platform", m.AuditHandler.ListPlatformLogs, middleware.RequirePermission("audit:read"))
		apiGroup.GET("/audit-logs/tenant", m.AuditHandler.ListTenantLogs, middleware.RequirePermission("audit:read"))
		apiGroup.GET("/audit-logs/stats", m.AuditHandler.GetAdminStats)
	}
}


