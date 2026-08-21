package api

import (
	"net/http"
	"strings"

	"github.com/golangnigeria/curexal/internal/modules/organization/application"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	subApp "github.com/golangnigeria/curexal/internal/modules/subscription/application"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TenantHandler struct {
	server         *server.Server
	appService     *application.OrganizationApplicationService
	tenantRepo     domain.TenantRepository
	entitlementSvc *subApp.EntitlementService
}

func NewTenantHandler(s *server.Server, appService *application.OrganizationApplicationService, tenantRepo domain.TenantRepository) *TenantHandler {
	return &TenantHandler{
		server:     s,
		appService: appService,
		tenantRepo: tenantRepo,
	}
}

func (h *TenantHandler) SetEntitlementService(svc *subApp.EntitlementService) {
	h.entitlementSvc = svc
}


func (h *TenantHandler) CreateTenant(c echo.Context) error {
	var payload CreateTenantPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	activeTenantID := middleware.GetActiveTenantID(c)
	if activeTenantID != "" {
		role := middleware.GetUserRole(c)
		if role != "owner" && role != "super_admin" {
			return echo.NewHTTPError(http.StatusForbidden, "Only organization owners can register new branches")
		}
	}

	userID := middleware.GetUserID(c)
	orgID, err := uuid.Parse(payload.OrganizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID format")
	}

	t, err := h.appService.CreateBranch(c, orgID, userID, payload.Name, payload.Slug, payload.Location, payload.Phone, payload.Address, payload.LogoURL, payload.Currency, payload.Modules)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, MapTenantToResponse(t))
}



func (h *TenantHandler) GetActiveTenant(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, _ := c.Get("tenant_id").(string)
	tenantName, _ := c.Get("tenant_name").(string)
	tenantSlug, _ := c.Get("tenant_slug").(string)
	currency, _ := c.Get("tenant_currency").(string)

	var logoURL *string
	if val, ok := c.Get("tenant_logo_url").(string); ok && val != "" {
		logoURL = &val
	}

	// 1. Fallback from request headers
	if tenantSlug == "" {
		tenantSlug = c.Request().Header.Get("X-Tenant-Slug")
	}
	if tenantID == "" {
		tenantID = c.Request().Header.Get("X-Active-Tenant-ID")
	}

	// 2. Fallback from subdomain
	if tenantSlug == "" && tenantID == "" {
		host := c.Request().Host
		parts := strings.Split(host, ".")
		if len(parts) > 2 || (strings.Contains(host, "localhost") && len(parts) > 1 && parts[0] != "localhost") {
			sub := parts[0]
			if sub != "" && sub != "admin" && sub != "app" && sub != "www" {
				tenantSlug = sub
			}
		}
	}

	// 3. Resolve tenant by slug or ID
	if (tenantSlug != "" || tenantID != "") && (tenantName == "" || tenantID == "") {
		target := tenantSlug
		if target == "" {
			target = tenantID
		}
		if t, err := h.tenantRepo.GetTenantBySlug(ctx, target); err == nil && t != nil {
			tenantID = t.ID.String()
			tenantName = t.Name
			tenantSlug = t.Slug
			currency = t.Currency
			logoURL = t.LogoURL
		}
	}

	// 4. Fallback: Query first available tenant in DB
	if tenantID == "" {
		tenants, err := h.tenantRepo.ListTenants(ctx)
		if err == nil && len(tenants) > 0 {
			t := tenants[0]
			tenantID = t.ID.String()
			tenantName = t.Name
			tenantSlug = t.Slug
			currency = t.Currency
			logoURL = t.LogoURL
		}
	}

	if tenantID == "" {
		tenantID = "default-tenant"
		tenantName = "Default Workspace"
		tenantSlug = "default"
	}

	branding := TenantBranding{
		PrimaryColor:   "#0284c7",
		SecondaryColor: "#0f172a",
		FontFamily:     "Outfit",
	}

	if currency == "" {
		currency = "NGN"
	}

	return c.JSON(http.StatusOK, &TenantResponse{
		ID:       tenantID,
		Name:     tenantName,
		Slug:     tenantSlug,
		LogoURL:  logoURL,
		Branding: branding,
		Currency: currency,
	})
}

func (h *TenantHandler) ListTenants(c echo.Context) error {
	tenants, err := h.tenantRepo.ListTenants(c.Request().Context())
	if err != nil {
		return err
	}
	responses := make([]*TenantResponse, len(tenants))
	for i, t := range tenants {
		responses[i] = MapTenantToResponse(&t)
	}
	return c.JSON(http.StatusOK, responses)
}



func (h *TenantHandler) GetWorkspaceContext(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, _ := c.Get("tenant_id").(string)
	tenantName, _ := c.Get("tenant_name").(string)
	tenantSlug, _ := c.Get("tenant_slug").(string)
	currency, _ := c.Get("tenant_currency").(string)

	if tenantSlug == "" {
		tenantSlug = c.Request().Header.Get("X-Tenant-Slug")
	}
	if tenantID == "" {
		tenantID = c.Request().Header.Get("X-Active-Tenant-ID")
	}

	var foundTenant *domain.Tenant
	if tenantSlug != "" || tenantID != "" {
		target := tenantSlug
		if target == "" {
			target = tenantID
		}
		if t, err := h.tenantRepo.GetTenantBySlug(ctx, target); err == nil && t != nil {
			foundTenant = t
			tenantID = t.ID.String()
			tenantName = t.Name
			tenantSlug = t.Slug
			currency = t.Currency
		}
	}

	if foundTenant == nil {
		tenants, err := h.tenantRepo.ListTenants(ctx)
		if err == nil && len(tenants) > 0 {
			foundTenant = &tenants[0]
			tenantID = foundTenant.ID.String()
			tenantName = foundTenant.Name
			tenantSlug = foundTenant.Slug
			currency = foundTenant.Currency
		}
	}

	if tenantID == "" {
		tenantID = "default-tenant-id"
		tenantName = "Diagnostic Laboratory Workspace"
		tenantSlug = "default-lab"
	}
	if currency == "" {
		currency = "NGN"
	}

	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)
	if userRole == "" {
		userRole = "branch_admin"
	}

	enabledModules := []string{"laboratory", "clinical", "pharmacy", "billing", "inventory", "customer_care", "qms"}
	if foundTenant != nil && len(foundTenant.EnabledModules) > 0 {
		enabledModules = foundTenant.EnabledModules
	}

	orgID := ""
	orgName := "Curexal Health Network"
	if foundTenant != nil {
		orgID = foundTenant.OrganizationID
	}

	activePlan := "smart"
	effectiveCapabilities := enabledModules

	if orgID != "" {
		if orgUUID, errP := uuid.Parse(orgID); errP == nil && h.entitlementSvc != nil {
			if realPlan, errPlan := h.entitlementSvc.GetOrganizationPlan(c.Request().Context(), orgUUID); errPlan == nil && realPlan != "" {
				activePlan = realPlan
			}
			if effCaps, errCaps := h.entitlementSvc.GetEffectiveCapabilities(c.Request().Context(), orgUUID); errCaps == nil && len(effCaps) > 0 {
				effectiveCapabilities = effCaps
			}
		}
	}

	response := map[string]interface{}{
		"workspace": map[string]interface{}{
			"id":           tenantID,
			"name":         tenantName,
			"slug":         tenantSlug,
			"facilityType": "Laboratory",
			"currency":     currency,
		},
		"organization": map[string]interface{}{
			"id":   orgID,
			"name": orgName,
		},
		"subscription": map[string]interface{}{
			"plan": activePlan,
		},
		"capabilities": effectiveCapabilities,
		"user": map[string]interface{}{
			"id":   userID,
			"role": userRole,
		},
	}

	return c.JSON(http.StatusOK, response)

}

