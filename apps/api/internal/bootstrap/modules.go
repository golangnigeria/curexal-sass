package bootstrap

import (
	"context"
	"net/http"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/audit"
	"github.com/golangnigeria/curexal/internal/modules/authorization"
	"github.com/golangnigeria/curexal/internal/modules/billing"
	"github.com/golangnigeria/curexal/internal/modules/catalogs"
	"github.com/golangnigeria/curexal/internal/modules/clinical"
	"github.com/golangnigeria/curexal/internal/modules/facility_config"
	"github.com/golangnigeria/curexal/internal/modules/identity"
	identityHdl "github.com/golangnigeria/curexal/internal/modules/identity/handler"
	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/identity/repository"
	"github.com/golangnigeria/curexal/internal/modules/notification"
	"github.com/golangnigeria/curexal/internal/modules/organization"
	"github.com/golangnigeria/curexal/internal/modules/organization/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/modules/patient"
	"github.com/golangnigeria/curexal/internal/modules/platform"
	"github.com/golangnigeria/curexal/internal/modules/settings"
	"github.com/golangnigeria/curexal/internal/modules/subscription"
	subAPI "github.com/golangnigeria/curexal/internal/modules/subscription/api"
	subApp "github.com/golangnigeria/curexal/internal/modules/subscription/application"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type tenantLookupAdapter struct {
	tenantRepo *postgres.TenantRepository
}

func (a *tenantLookupAdapter) GetTenantByID(ctx context.Context, id uuid.UUID) (*model.TenantInfo, error) {
	t, err := a.tenantRepo.GetTenantByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.TenantInfo{
		ID:   t.ID,
		Name: t.Name,
		Slug: t.Slug,
	}, nil
}

// ModuleRegistry holds initialized bounded context modules.
type ModuleRegistry struct {
	Identity       *identity.Module
	Platform       *platform.Module
	Organization   *organization.Module
	Authorization  *authorization.Module
	Clinical       *clinical.Module
	Catalogs       *catalogs.Module
	Audit          *audit.Module
	FacilityConfig *facility_config.Module
	Billing        *billing.Module
	Settings       *settings.Module
	Patient        *patient.Module
	Notification   *notification.Module
	Subscription   *subscription.Module
	EntitlementSvc *subApp.EntitlementService
	EntitlementHdl *subAPI.EntitlementHandler
}

// InitModules initializes all bounded context modules using the server composition handle.
func InitModules(s *server.Server) *ModuleRegistry {
	userRepo := repository.NewUserRepository(s)
	patientMod := patient.NewModule(s, userRepo)
	orgTenantRepo := postgres.NewTenantRepository(s)
	lookupAdapter := &tenantLookupAdapter{tenantRepo: orgTenantRepo}
	identityMod := identity.NewModule(s, patientMod.Service, patientMod.Repo, lookupAdapter)
	notifMod := notification.NewModule(s, userRepo)

	subMod := subscription.NewModule(s)
	entitlementSvc := subMod.Service
	entitlementHdl := subMod.Handler

	platMod := platform.NewModule(s)
	if platMod != nil && platMod.BootstrapHandler != nil {
		platMod.BootstrapHandler.SetEntitlementService(entitlementSvc)
	}

	orgMod := organization.NewModule(s)
	if orgMod != nil && orgMod.TenantHandler != nil {
		orgMod.TenantHandler.SetEntitlementService(entitlementSvc)
	}

	reg := &ModuleRegistry{
		Identity:       identityMod,
		Patient:        patientMod,
		Notification:   notifMod,
		Platform:       platMod,
		Organization:   orgMod,
		Subscription:   subMod,
		Authorization:  authorization.NewModule(s),
		Clinical:       clinical.NewModule(s),
		Catalogs:       catalogs.NewModule(s),
		Audit:          audit.NewModule(s),
		FacilityConfig: facility_config.NewModule(s),
		Billing:        billing.NewModule(s),
		Settings:       settings.NewModule(s),
		EntitlementSvc: entitlementSvc,
		EntitlementHdl: entitlementHdl,
	}
	reg.RegisterRoutes(s)
	return reg
}


func (r *ModuleRegistry) RegisterRoutes(s *server.Server) {
	if s.Echo == nil {
		return
	}

	// Public Health Status Endpoints
	statusHandler := func(c echo.Context) error {
		dbStatus := "healthy"
		if s.DB == nil || s.DB.Pool == nil {
			dbStatus = "unhealthy"
		} else if err := s.DB.Pool.Ping(c.Request().Context()); err != nil {
			dbStatus = "unhealthy"
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":      "healthy",
			"timestamp":   time.Now().Format(time.RFC3339),
			"environment": s.Config.Primary.Env,
			"checks": map[string]interface{}{
				"database": map[string]interface{}{
					"status": dbStatus,
				},
			},
		})
	}
	s.Echo.GET("/status", statusHandler)
	s.Echo.GET("/api/v1/status", statusHandler)

	api := s.Echo.Group("/api/v1")
	api.Use(middleware.Authenticate(s.Config))

	plt := api.Group("/platform")
	org := api.Group("/organization")
	wsp := api.Group("/workspace")
	wsp.Use(middleware.HostAndSessionTenantResolver())

	// Bounded-Context Module Route Registrations
	if r.Identity != nil {
		r.Identity.RegisterRoutes(api, org, wsp)
	}

	if r.Platform != nil {
		r.Platform.RegisterRoutes(api, plt)
	}

	if r.Subscription != nil {
		r.Subscription.RegisterRoutes(s.Echo, api, plt)
	}

	if r.Organization != nil {
		var identityUserHdl *identityHdl.UserRoleHandler
		if r.Identity != nil {
			identityUserHdl = r.Identity.UserRoleHandler
		}
		r.Organization.RegisterRoutes(api, plt, org, wsp, identityUserHdl)
	}

	if r.Clinical != nil {
		r.Clinical.RegisterRoutes(api, r.EntitlementSvc)
	}

	if r.Audit != nil {
		r.Audit.RegisterRoutes(api)
	}

	if r.FacilityConfig != nil {
		r.FacilityConfig.RegisterRoutes(plt)
	}

	if r.Billing != nil {
		r.Billing.RegisterRoutes(api)
	}

	if r.Authorization != nil {
		r.Authorization.RegisterRoutes(api)
	}

	if r.Patient != nil {
		r.Patient.RegisterRoutes(api)
	}

	if r.Settings != nil {
		r.Settings.RegisterRoutes(api)
	}

	if r.Catalogs != nil {
		r.Catalogs.RegisterRoutes(api)
	}

	if r.Notification != nil {
		r.Notification.RegisterRoutes(s.Echo, middleware.Authenticate(s.Config))
	}
}

