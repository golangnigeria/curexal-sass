package organization

import (
	"fmt"

	auditPostgres "github.com/golangnigeria/curexal/internal/modules/audit/infrastructure/postgres"
	identityHdl "github.com/golangnigeria/curexal/internal/modules/identity/handler"
	"github.com/golangnigeria/curexal/internal/modules/organization/api"
	"github.com/golangnigeria/curexal/internal/modules/organization/application"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/modules/organization/infrastructure/postgres"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/kernel/storage"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/labstack/echo/v4"
)

type Module struct {
	OrgRepo            domain.OrganizationRepository
	TenantRepo         domain.TenantRepository
	SubRepo            domain.SubscriptionRepository
	DemoRepo           domain.DemoRepository
	DocRepo            domain.DocumentRepository
	BranchRepo         domain.FacilityBranchRepository
	StaffRepo          domain.StaffMembershipRepository
	CatalogRepo        domain.OrganizationCatalogRepository
	BrandingRepo       domain.OrganizationBrandingRepository
	IntegrationRepo    domain.OrganizationIntegrationRepository
	Storage            storage.ObjectStorageService
	AppService         *application.OrganizationApplicationService
	DocService         *application.OrganizationDocumentApplicationService
	SetupService       *application.OrganizationSetupService
	BranchService      *application.FacilityBranchService
	StaffService       *application.StaffMembershipService
	CatalogService     *application.OrganizationCatalogService
	BrandingService    *application.OrganizationBrandingService
	IntegrationService *application.OrganizationIntegrationService
	OrgHandler         *api.OrganizationHandler
	TenantHandler      *api.TenantHandler
	DemoHandler        *api.DemoHandler
	DocHandler         *api.DocumentHandler
	ProfileHandler     *api.OrganizationProfileHandler
	BranchHandler      *api.FacilityBranchHandler
	StaffHandler       *api.StaffMembershipHandler
	CatalogHandler     *api.OrganizationCatalogHandler
	BrandingHandler    *api.OrganizationBrandingHandler
	IntegrationHandler *api.OrganizationIntegrationHandler
}

func NewModule(s *server.Server) *Module {
	orgRepo := postgres.NewOrganizationRepository(s)
	tenantRepo := postgres.NewTenantRepository(s)
	subRepo := postgres.NewSubscriptionRepository(s)
	demoRepo := postgres.NewDemoRepository(s)
	docRepo := postgres.NewDocumentRepository(s)
	branchRepo := postgres.NewFacilityBranchRepository(s)
	staffRepo := postgres.NewStaffMembershipRepository(s)
	catalogRepo := postgres.NewOrganizationCatalogRepository(s)
	brandingRepo := postgres.NewOrganizationBrandingRepository(s)
	integrationRepo := postgres.NewOrganizationIntegrationRepository(s)

	var auditRepo *auditPostgres.AuditRepository
	if s != nil && s.DB != nil {
		auditRepo = auditPostgres.NewAuditRepository(s)
	}

	masterKey := ""
	baseURL := ""
	if s != nil && s.Config != nil {
		masterKey = s.Config.Auth.SecretKey
		if s.Config.Server.Port != "" {
			baseURL = fmt.Sprintf("http://localhost:%s", s.Config.Server.Port)
		}
	}
	var storageService storage.ObjectStorageService
	if s != nil && s.Storage != nil {
		storageService = s.Storage
	} else {
		storageService, _ = storage.NewLocalStorageService("./storage/documents", baseURL)
	}

	appService := application.NewOrganizationApplicationService(s, orgRepo, tenantRepo, subRepo, demoRepo)
	docService := application.NewOrganizationDocumentApplicationService(s, docRepo, orgRepo, storageService)
	setupService := application.NewOrganizationSetupService(orgRepo, auditRepo)
	branchService := application.NewFacilityBranchService(branchRepo, orgRepo, auditRepo)
	staffService := application.NewStaffMembershipService(staffRepo, orgRepo, branchRepo, auditRepo)
	catalogService := application.NewOrganizationCatalogService(catalogRepo, branchRepo, orgRepo, auditRepo)
	brandingService := application.NewOrganizationBrandingService(brandingRepo, orgRepo, auditRepo, masterKey)
	integrationService := application.NewOrganizationIntegrationService(integrationRepo, orgRepo, auditRepo, masterKey)

	orgHandler := api.NewOrganizationHandler(s, appService)
	tenantHandler := api.NewTenantHandler(s, appService, tenantRepo)
	demoHandler := api.NewDemoHandler(s, appService)
	docHandler := api.NewDocumentHandler(s, docService)
	profileHandler := api.NewOrganizationProfileHandler(setupService)
	branchHandler := api.NewFacilityBranchHandler(branchService)
	staffHandler := api.NewStaffMembershipHandler(staffService)
	catalogHandler := api.NewOrganizationCatalogHandler(catalogService)
	brandingHandler := api.NewOrganizationBrandingHandler(brandingService)
	integrationHandler := api.NewOrganizationIntegrationHandler(integrationService)

	return &Module{
		OrgRepo:            orgRepo,
		TenantRepo:         tenantRepo,
		SubRepo:            subRepo,
		DemoRepo:           demoRepo,
		DocRepo:            docRepo,
		BranchRepo:         branchRepo,
		StaffRepo:          staffRepo,
		CatalogRepo:        catalogRepo,
		BrandingRepo:       brandingRepo,
		IntegrationRepo:    integrationRepo,
		Storage:            storageService,
		AppService:         appService,
		DocService:         docService,
		SetupService:       setupService,
		BranchService:      branchService,
		StaffService:       staffService,
		CatalogService:     catalogService,
		BrandingService:    brandingService,
		IntegrationService: integrationService,
		OrgHandler:         orgHandler,
		TenantHandler:      tenantHandler,
		DemoHandler:        demoHandler,
		DocHandler:         docHandler,
		ProfileHandler:     profileHandler,
		BranchHandler:      branchHandler,
		StaffHandler:       staffHandler,
		CatalogHandler:     catalogHandler,
		BrandingHandler:    brandingHandler,
		IntegrationHandler: integrationHandler,
	}
}

func (m *Module) RegisterRoutes(apiGroup *echo.Group, pltGroup *echo.Group, orgGroup *echo.Group, wspGroup *echo.Group, identityHandler *identityHdl.UserRoleHandler) {
	if m.IntegrationHandler != nil {
		apiGroup.GET("/organization/api-keys", m.IntegrationHandler.ListAPIKeys)
		apiGroup.POST("/organization/api-keys", m.IntegrationHandler.CreateAPIKey, middleware.RequirePermission("organization:integrations:write"))
		apiGroup.DELETE("/organization/api-keys/:id", m.IntegrationHandler.RevokeAPIKey, middleware.RequirePermission("organization:integrations:write"))
		apiGroup.GET("/organization/webhooks", m.IntegrationHandler.ListWebhookSubscriptions)
		apiGroup.POST("/organization/webhooks", m.IntegrationHandler.CreateWebhookSubscription, middleware.RequirePermission("organization:integrations:write"))
		apiGroup.DELETE("/organization/webhooks/:id", m.IntegrationHandler.DeleteWebhookSubscription, middleware.RequirePermission("organization:integrations:write"))
		apiGroup.GET("/organization/webhooks/deliveries", m.IntegrationHandler.ListWebhookDeliveries)
	}

	if m.BrandingHandler != nil {
		apiGroup.GET("/organization/branding", m.BrandingHandler.GetBranding)
		apiGroup.PUT("/organization/branding", m.BrandingHandler.UpdateBranding, middleware.RequirePermission("organization:branding:write"))
		apiGroup.GET("/organization/notifications/configs", m.BrandingHandler.ListNotificationConfigs)
		apiGroup.PUT("/organization/notifications/configs", m.BrandingHandler.SaveNotificationConfig, middleware.RequirePermission("organization:notifications:write"))
		apiGroup.GET("/organization/notifications/templates", m.BrandingHandler.ListNotificationTemplates)
		apiGroup.PUT("/organization/notifications/templates/:key", m.BrandingHandler.SaveNotificationTemplate, middleware.RequirePermission("organization:notifications:write"))
		apiGroup.GET("/organization/notifications", m.BrandingHandler.ListUserNotifications)
		apiGroup.POST("/organization/notifications/:id/read", m.BrandingHandler.MarkNotificationRead)
		apiGroup.POST("/organization/notifications/read-all", m.BrandingHandler.MarkAllNotificationsRead)
		apiGroup.GET("/organization/notifications/deliveries", m.BrandingHandler.ListNotificationDeliveries)
	}

	if m.CatalogHandler != nil {
		apiGroup.GET("/organization/catalogs", m.CatalogHandler.ListCatalogItems)
		apiGroup.POST("/organization/catalogs", m.CatalogHandler.CreateCatalogItem, middleware.RequirePermission("organization:catalog:write"))
		apiGroup.GET("/organization/catalogs/:id", m.CatalogHandler.GetCatalogItem)
		apiGroup.PUT("/organization/catalogs/:id", m.CatalogHandler.UpdateCatalogItem, middleware.RequirePermission("organization:catalog:write"))
		apiGroup.POST("/organization/catalogs/:id/branch-prices", m.CatalogHandler.SetBranchPriceOverride, middleware.RequirePermission("organization:catalog:write"))
		apiGroup.GET("/organization/insurance-providers", m.CatalogHandler.ListInsuranceProviders)
		apiGroup.POST("/organization/insurance-providers", m.CatalogHandler.CreateInsuranceProvider, middleware.RequirePermission("organization:write"))
	}

	if m.StaffHandler != nil {
		apiGroup.GET("/organization/members", m.StaffHandler.ListMembers)
		apiGroup.POST("/organization/invitations", m.StaffHandler.CreateInvitation, middleware.RequirePermission("users:write"))
		apiGroup.GET("/organization/invitations", m.StaffHandler.ListInvitations)
		apiGroup.DELETE("/organization/invitations/:id", m.StaffHandler.RevokeInvitation, middleware.RequirePermission("users:write"))
		apiGroup.POST("/organization/members/:id/branches", m.StaffHandler.AssignBranch, middleware.RequirePermission("users:write"))
		apiGroup.POST("/organization/members/:id/departments", m.StaffHandler.AssignDepartment, middleware.RequirePermission("users:write"))
		apiGroup.PUT("/organization/members/:id/role", m.StaffHandler.UpdateRole, middleware.RequirePermission("users:write"))
	}

	if m.BranchHandler != nil {
		apiGroup.GET("/organization/branches", m.BranchHandler.ListBranches)
		apiGroup.POST("/organization/branches", m.BranchHandler.CreateBranch, middleware.RequirePermission("organization:branch:create"))
		apiGroup.GET("/organization/branches/:id", m.BranchHandler.GetBranch)
		apiGroup.PUT("/organization/branches/:id", m.BranchHandler.UpdateBranch, middleware.RequirePermission("organization:branch:update"))
		apiGroup.DELETE("/organization/branches/:id", m.BranchHandler.DeactivateBranch, middleware.RequirePermission("organization:branch:deactivate"))
	}

	if m.ProfileHandler != nil {
		apiGroup.GET("/organization/profile", m.ProfileHandler.GetProfile)
		apiGroup.PUT("/organization/profile", m.ProfileHandler.UpdateProfile)
		apiGroup.POST("/organization/setup/submit-review", m.ProfileHandler.SubmitForReview)
	}

	if m.TenantHandler != nil {
		apiGroup.GET("/tenant/active", m.TenantHandler.GetActiveTenant)
		apiGroup.GET("/tenants", m.TenantHandler.ListTenants)
		apiGroup.POST("/tenants", m.TenantHandler.CreateTenant, middleware.RequirePermission("organization:write"))
	}
	if m.OrgHandler != nil {
		apiGroup.GET("/organizations", m.OrgHandler.ListOrganizations, middleware.RequirePermission("organization:read"))
		apiGroup.POST("/organizations", m.OrgHandler.CreateOrganization)
		apiGroup.GET("/organizations/:id", m.OrgHandler.GetOrganization)
		apiGroup.PUT("/organizations/:id", m.OrgHandler.UpdateOrganization)
		apiGroup.POST("/organizations/:id/transfer-ownership", m.OrgHandler.TransferOwnership, middleware.RequirePermission("organization:write"))
		apiGroup.POST("/organizations/:id/resend-invite", m.OrgHandler.ResendOwnerInvite, middleware.RequirePermission("organization:write"))
		apiGroup.GET("/organizations/:id/settings", m.OrgHandler.GetOrganizationSettings)
		apiGroup.PUT("/organizations/:id/settings", m.OrgHandler.UpdateOrganizationSettings, middleware.RequirePermission("organization:settings:write"))
	}

	if identityHandler != nil {
		apiGroup.GET("/organizations/:id/members", identityHandler.GetOrganizationMembers, middleware.RequirePermission("organization:read"))
		apiGroup.POST("/organizations/:id/members", identityHandler.CreateWorkspaceMember, middleware.RequirePermission("users:write"))
	}

	if m.DemoHandler != nil {
		apiGroup.POST("/demo-requests", m.DemoHandler.CreateDemoRequest)
		apiGroup.GET("/demo-requests", m.DemoHandler.ListDemoRequests)
		apiGroup.PUT("/demo-requests/:id", m.DemoHandler.UpdateDemoRequest)
	}

	if m.DocHandler != nil {
		apiGroup.POST("/organizations/:id/documents", m.DocHandler.UploadDocument, middleware.RequirePermission("organization:document:upload"))
		apiGroup.GET("/organizations/:id/documents", m.DocHandler.ListDocuments, middleware.RequirePermission("organization:document:read"))
	}

	if pltGroup != nil {
		if m.OrgHandler != nil {
			pltGroup.GET("/organizations", m.OrgHandler.ListOrganizations)
		}
		if m.ProfileHandler != nil {
			pltGroup.POST("/organizations/:id/verify", m.ProfileHandler.VerifyOrganization)
		}
		if m.DocHandler != nil {
			pltGroup.PATCH("/documents/:docID/review", m.DocHandler.ReviewDocument, middleware.RequirePermission("organization:document:review"))
			pltGroup.POST("/organizations/:id/approve", m.DocHandler.ApproveOrganization, middleware.RequirePermission("organization:verify"))
			pltGroup.POST("/organizations/:id/reject", m.DocHandler.RejectOrganization, middleware.RequirePermission("organization:verify"))
		}
	}

	if orgGroup != nil && m.BranchHandler != nil {
		orgGroup.GET("/branches", m.BranchHandler.ListBranches)
		orgGroup.POST("/branches", m.BranchHandler.CreateBranch)
	}

	if wspGroup != nil && m.TenantHandler != nil {
		wspGroup.GET("/dashboard", m.TenantHandler.GetActiveTenant)
		wspGroup.GET("/context", m.TenantHandler.GetWorkspaceContext)
	}
}
