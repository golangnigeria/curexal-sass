package application

import (
	"context"
	"encoding/json"
	"strings"

	orgDomain "github.com/golangnigeria/curexal/internal/modules/organization/domain"
	platformDomain "github.com/golangnigeria/curexal/internal/modules/platform/domain"
	subApp "github.com/golangnigeria/curexal/internal/modules/subscription/application"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IdentityPayload struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar,omitempty"`
	Locale      string `json:"locale"`
	Timezone    string `json:"timezone"`
}

type PlatformPayload struct {
	IsStaff bool   `json:"isStaff"`
	Role    string `json:"role"`
}

type OrganizationPayload struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug,omitempty"`
	Logo         string `json:"logo,omitempty"`
	Role         string `json:"role,omitempty"`
	Subscription string `json:"subscription"`
}

type WorkspacePayload struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	FacilityType string `json:"facilityType"`
	Slug         string `json:"slug"`
	Timezone     string `json:"timezone"`
	Currency     string `json:"currency"`
}

type SubscriptionPayload struct {
	Plan   string         `json:"plan"`
	Status string         `json:"status"`
	Limits map[string]int `json:"limits"`
}

type ModuleCapabilityPayload struct {
	Code             string   `json:"code"`
	Enabled          bool     `json:"enabled"`
	Licensed         bool     `json:"licensed"`
	Visible          bool     `json:"visible"`
	UpgradeAvailable bool     `json:"upgradeAvailable"`
	Actions          []string `json:"actions"`
}

type NavigationItemPayload struct {
	ID       string                  `json:"id"`
	Title    string                  `json:"title"`
	Icon     string                  `json:"icon"`
	Path     string                  `json:"path"`
	Order    int                     `json:"order,omitempty"`
	Children []NavigationItemPayload `json:"children,omitempty"`
}

type BreadcrumbPayload struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type StructuredNavigationPayload struct {
	Primary      []NavigationItemPayload `json:"primary"`
	Secondary    []NavigationItemPayload `json:"secondary,omitempty"`
	Topbar       []NavigationItemPayload `json:"topbar,omitempty"`
	QuickActions []NavigationItemPayload `json:"quickActions,omitempty"`
	Breadcrumbs  []BreadcrumbPayload     `json:"breadcrumbs,omitempty"`
}

type DashboardSectionPayload struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Type  string   `json:"type"`
	Items []string `json:"items,omitempty"`
}

type DashboardPayload struct {
	Widgets  []string                  `json:"widgets"`
	Sections []DashboardSectionPayload `json:"sections,omitempty"`
}

type ContextsPayload struct {
	Current   string   `json:"current"`
	Available []string `json:"available"`
	Default   string   `json:"default"`
}

type BrandingPayload struct {
	LogoURL          string                 `json:"logoUrl,omitempty"`
	FaviconURL       string                 `json:"faviconUrl,omitempty"`
	PrimaryColor     string                 `json:"primaryColor"`
	SecondaryColor   string                 `json:"secondaryColor,omitempty"`
	AccentColor      string                 `json:"accentColor,omitempty"`
	FontFamily       string                 `json:"fontFamily,omitempty"`
	BorderRadius     string                 `json:"borderRadius,omitempty"`
	CustomDomain     string                 `json:"customDomain,omitempty"`
	HideCurexalBadge bool                   `json:"hideCurexalBadge,omitempty"`
	ThemeBranding    map[string]interface{} `json:"themeBranding,omitempty"`
}

type PreferencesPayload struct {
	Theme    string `json:"theme"`
	Language string `json:"language"`
	Timezone string `json:"timezone"`
}

type MetadataPayload struct {
	Version     string `json:"version"`
	GeneratedAt string `json:"generatedAt"`
	TTL         int    `json:"ttl"`
	ETag        string `json:"etag"`
}

type BootstrapContractResponse struct {
	Identity             IdentityPayload             `json:"identity"`
	Platform             PlatformPayload             `json:"platform"`
	Organization         OrganizationPayload         `json:"organization"`
	Workspace            WorkspacePayload            `json:"workspace"`
	Subscription         SubscriptionPayload         `json:"subscription"`
	Modules              []ModuleCapabilityPayload   `json:"modules"`
	Capabilities         []string                    `json:"capabilities"`
	Permissions          []string                    `json:"permissions"`
	Navigation           []NavigationItemPayload     `json:"navigation"`
	StructuredNavigation StructuredNavigationPayload `json:"structuredNavigation"`
	Dashboard            DashboardPayload            `json:"dashboard"`
	Contexts             ContextsPayload             `json:"contexts"`
	Branding             BrandingPayload             `json:"branding"`
	Preferences          PreferencesPayload          `json:"preferences"`
	FeatureFlags         map[string]bool             `json:"featureFlags"`
	Limits               map[string]int              `json:"limits"`
	Metadata             MetadataPayload             `json:"metadata"`
}

type BootstrapBuilder struct {
	orgRepo        orgDomain.OrganizationRepository
	tenantRepo     orgDomain.TenantRepository
	subRepo        orgDomain.SubscriptionRepository
	platformDomain *platformDomain.PlatformDomainProvider
	orgDomain      *orgDomain.OrganizationDomainProvider
	entitlementSvc *subApp.EntitlementService
	dbPool         *pgxpool.Pool
}

func NewBootstrapBuilder(orgRepo orgDomain.OrganizationRepository, tenantRepo orgDomain.TenantRepository, subRepo orgDomain.SubscriptionRepository) *BootstrapBuilder {
	return &BootstrapBuilder{
		orgRepo:        orgRepo,
		tenantRepo:     tenantRepo,
		subRepo:        subRepo,
		platformDomain: platformDomain.NewPlatformDomainProvider(),
		orgDomain:      orgDomain.NewOrganizationDomainProvider(),
	}
}

func (b *BootstrapBuilder) SetDBPool(pool *pgxpool.Pool) {
	b.dbPool = pool
}

func (b *BootstrapBuilder) BuildBootstrap(ctx context.Context, principal *middleware.AuthenticatedPrincipal) (*BootstrapContractResponse, error) {
	userID := "usr_default"
	userEmail := ""
	displayName := ""
	userRole := ""
	isPlatformStaff := false

	if principal != nil {
		if principal.UserID != "" {
			userID = principal.UserID
		}
		if principal.Identity.Email != "" {
			userEmail = principal.Identity.Email
		}
		if principal.Identity.FullName != "" {
			displayName = principal.Identity.FullName
		}
		if principal.Role != "" {
			userRole = principal.Role
		}
		if principal.Platform.PlatformRole != "" {
			userRole = principal.Platform.PlatformRole
		}
		isPlatformStaff = b.platformDomain.IsPlatformContext(
			principal.Platform.IsPlatformStaff || principal.Platform.IsSuperAdmin || principal.Platform.IsPlatformAdmin,
			userRole,
			userEmail,
		)
	}

	// 1. Authoritative Organization Membership & Branding Resolution
	var activeOrgID string
	var activeOrgName string
	var activeOrgSlug string
	var activeOrgPlan string
	var activeMembershipRole string
	var activeMembershipTenantID string
	var activePrimaryColor string
	var activeThemeBrandingJSON []byte
	var activeCustomDomain string
	var activeLogoURL string
	var activeTenant *orgDomain.Tenant

	if b.dbPool != nil && userID != "" && userID != "usr_default" {
		row := b.dbPool.QueryRow(ctx, `
			SELECT 
				m.organization_id::text, 
				COALESCE(m.tenant_id::text, ''), 
				COALESCE(m.role, m.role_title, 'member'), 
				o.name, 
				o.slug, 
				COALESCE(o.plan, 'smart'),
				COALESCE(o.primary_color, '#0284c7'),
				COALESCE(o.theme_branding, '{}'::jsonb),
				COALESCE(o.custom_domain, ''),
				COALESCE(o.logo_url, '')
			FROM organization.organization_memberships m
			JOIN organization.organizations o ON o.id = m.organization_id
			WHERE m.user_id = $1 AND m.is_active = TRUE
			ORDER BY (m.role = 'owner') DESC, m.created_at ASC
			LIMIT 1
		`, userID)
		_ = row.Scan(
			&activeOrgID, 
			&activeMembershipTenantID, 
			&activeMembershipRole, 
			&activeOrgName, 
			&activeOrgSlug, 
			&activeOrgPlan,
			&activePrimaryColor,
			&activeThemeBrandingJSON,
			&activeCustomDomain,
			&activeLogoURL,
		)
	}

	// Fallback to orgRepo if not resolved from direct query
	if activeOrgID == "" && b.orgRepo != nil && userID != "" {
		orgs, err := b.orgRepo.List(ctx, userID, isPlatformStaff)
		if err == nil && len(orgs) > 0 {
			activeOrgID = orgs[0].ID.String()
			activeOrgName = orgs[0].Name
			activeOrgSlug = orgs[0].Slug
			if orgs[0].Plan != "" {
				activeOrgPlan = orgs[0].Plan
			}
			activeMembershipRole = "owner"
		}
	}

	// 2. Resolve Active Context
	currentContext := "workspace"
	availableContexts := []string{"workspace"}

	if isPlatformStaff {
		currentContext = "platform"
		availableContexts = []string{"platform", "organization", "workspace"}
	} else if activeOrgID != "" {
		if activeMembershipRole == "owner" || activeMembershipRole == "org_admin" || activeMembershipRole == "admin" || activeMembershipRole == "org_regional_manager" || userRole == "owner" {
			currentContext = "organization"
			availableContexts = []string{"organization", "workspace"}
			if userRole == "" {
				userRole = activeMembershipRole
			}
		} else {
			currentContext = "workspace"
			availableContexts = []string{"workspace"}
		}
	}

	// 3. Resolve Workspace / Facility details
	targetTenantID := ""
	if principal != nil && principal.TenantID != "" {
		targetTenantID = principal.TenantID
	} else if activeMembershipTenantID != "" {
		targetTenantID = activeMembershipTenantID
	}

	if targetTenantID != "" && b.tenantRepo != nil {
		tenants, err := b.tenantRepo.ListTenants(ctx)
		if err == nil {
			for i := range tenants {
				if tenants[i].ID.String() == targetTenantID {
					activeTenant = &tenants[i]
					break
				}
			}
		}
	}
	if activeTenant == nil && b.tenantRepo != nil {
		tenants, err := b.tenantRepo.ListTenants(ctx)
		if err == nil && len(tenants) > 0 {
			activeTenant = &tenants[0]
		}
	}

	tenantID := "default-tenant-id"
	tenantName := "Main Diagnostic Facility"
	tenantSlug := "main-facility"
	currency := "NGN"
	enabledModules := []string{"laboratory", "clinical", "pharmacy", "billing", "inventory", "customer_care", "qms"}

	if activeTenant != nil {
		tenantID = activeTenant.ID.String()
		tenantName = activeTenant.Name
		tenantSlug = activeTenant.Slug
		if activeTenant.Currency != "" {
			currency = activeTenant.Currency
		}
		if len(activeTenant.EnabledModules) > 0 {
			enabledModules = activeTenant.EnabledModules
		}
	}

	allModuleCodes := []string{"laboratory", "clinical", "customer_care", "billing", "pharmacy", "inventory", "qms", "radiology"}
	enabledSet := make(map[string]bool)
	for _, m := range enabledModules {
		enabledSet[m] = true
	}

	// 4. Resolve Effective Permissions
	effectivePermissions := make([]string, 0)
	isSuperAdminOrOwner := isPlatformStaff || userRole == "owner" || activeMembershipRole == "owner"

	if isPlatformStaff {
		effectivePermissions = platformAuth.GetAllPermissions()
	} else if principal != nil && len(principal.Permissions) > 0 {
		seen := make(map[string]bool)
		for _, p := range principal.Permissions {
			if p == "*" {
				for _, allP := range platformAuth.GetAllPermissions() {
					if !seen[allP] {
						seen[allP] = true
						effectivePermissions = append(effectivePermissions, allP)
					}
				}
			} else if !seen[p] {
				seen[p] = true
				effectivePermissions = append(effectivePermissions, p)
			}
		}
	} else if isSuperAdminOrOwner {
		effectivePermissions = platformAuth.GetAllPermissions()
	}

	if len(effectivePermissions) == 0 && b.dbPool != nil && userID != "" {
		rows, err := b.dbPool.Query(ctx, `
			SELECT DISTINCT p.code
			FROM "authorization".permissions p
			JOIN "authorization".role_permissions rp ON rp.permission_id = p.id
			JOIN "authorization".roles r ON r.id = rp.role_id
			JOIN organization.organization_memberships m ON (m.role_title = r.code OR m.role = r.code OR m.role_title = r.name)
			WHERE m.user_id = $1 AND (m.tenant_id = $2 OR $2 = '' OR $2 IS NULL) AND m.is_active = TRUE
		`, userID, tenantID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var pCode string
				if err := rows.Scan(&pCode); err == nil && pCode != "*" {
					effectivePermissions = append(effectivePermissions, pCode)
				}
			}
		}
	}

	if len(effectivePermissions) == 0 {
		effectivePermissions = []string{
			"organization:read",
			"organization:dashboard:read",
			"workspace:patient:read",
			"workspace:patient:create",
			"workspace:sample:receive",
			"workspace:worksheet:update",
			"workspace:result:authorize",
			"workspace:billing:create",
		}
	}

	// 5. Context-Bound Navigation Tree (DB-Driven & Permission Filtered)
	var navigationItems []NavigationItemPayload
	if b.dbPool != nil {
		rows, err := b.dbPool.Query(ctx, `
			SELECT id, title, icon, path, sort_order
			FROM navigation_item
			WHERE context_scope = $1
			  AND (module_code IS NULL OR module_code = ANY($2))
			  AND (required_permission IS NULL OR required_permission = ANY($3) OR $4 = TRUE)
			ORDER BY sort_order ASC
		`, currentContext, enabledModules, effectivePermissions, isSuperAdminOrOwner)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var item NavigationItemPayload
				if err := rows.Scan(&item.ID, &item.Title, &item.Icon, &item.Path, &item.Order); err == nil {
					navigationItems = append(navigationItems, item)
				}
			}
		}
	}

	if len(navigationItems) == 0 {
		if currentContext == "platform" {
			platNav := b.platformDomain.GetPlatformNavigation()
			for _, item := range platNav {
				navigationItems = append(navigationItems, NavigationItemPayload{
					ID:    item.ID,
					Title: item.Title,
					Icon:  item.Icon,
					Path:  item.Path,
					Order: item.Order,
				})
			}
		} else if currentContext == "organization" {
			orgNav := b.orgDomain.GetOrganizationNavigation()
			for _, item := range orgNav {
				navigationItems = append(navigationItems, NavigationItemPayload{
					ID:    item.ID,
					Title: item.Title,
					Icon:  item.Icon,
					Path:  item.Path,
					Order: item.Order,
				})
			}
		} else {
			// Workspace Context Fallback
			navigationItems = []NavigationItemPayload{
				{ID: "nav_wsp_dashboard", Title: "Workspace Dashboard", Icon: "LayoutDashboard", Path: "/workspace/dashboard", Order: 1},
			}
			if enabledSet["customer_care"] {
				navigationItems = append(navigationItems, NavigationItemPayload{
					ID: "nav_wsp_patients", Title: "Patient Reception", Icon: "UserPlus", Path: "/workspace/patients", Order: 2,
				})
			}
			if enabledSet["laboratory"] {
				navigationItems = append(navigationItems, NavigationItemPayload{
					ID: "nav_wsp_laboratory", Title: "Laboratory LIS", Icon: "Activity", Path: "/workspace/laboratory/accessioning", Order: 3,
				})
			}
			if enabledSet["clinical"] {
				navigationItems = append(navigationItems, NavigationItemPayload{
					ID: "nav_wsp_clinical", Title: "Clinical & EMR", Icon: "Stethoscope", Path: "/workspace/clinical/tests", Order: 4,
				})
			}
			if enabledSet["billing"] {
				navigationItems = append(navigationItems, NavigationItemPayload{
					ID: "nav_wsp_billing", Title: "Billing POS", Icon: "CreditCard", Path: "/workspace/billing", Order: 5,
				})
			}
			navigationItems = append(navigationItems, NavigationItemPayload{
				ID: "nav_wsp_settings", Title: "Facility Settings", Icon: "Settings", Path: "/workspace/settings", Order: 6,
			})
		}
	}


	// 6. Dashboard Widgets Spec
	dashboardWidgets := []string{}
	for _, m := range enabledModules {
		dashboardWidgets = append(dashboardWidgets, m)
	}

	// Resolve Organization & Real Base Plan from DB
	orgID := activeOrgID
	orgName := activeOrgName
	orgSlug := activeOrgSlug
	orgSubPlan := activeOrgPlan
	orgRole := activeMembershipRole

	if orgID == "" && b.orgRepo != nil {
		orgs, err := b.orgRepo.List(ctx, userID, isPlatformStaff)
		if err == nil && len(orgs) > 0 {
			var selectedOrg *orgDomain.Organization
			if principal != nil && principal.OrganizationID != "" {
				for i := range orgs {
					if orgs[i].ID.String() == principal.OrganizationID {
						selectedOrg = &orgs[i]
						break
					}
				}
			}
			if selectedOrg == nil {
				selectedOrg = &orgs[0]
			}

			orgID = selectedOrg.ID.String()
			orgName = selectedOrg.Name
			orgSlug = selectedOrg.Slug
			if selectedOrg.Plan != "" {
				orgSubPlan = strings.ToLower(selectedOrg.Plan)
			}
		}
	}

	if orgID == "" {
		orgID = "org_default"
		orgName = "Curexal Health Network"
		orgSlug = "curexal"
		orgSubPlan = "smart"
		orgRole = userRole
	}

	effectiveCapabilities := []string{
		"core.organization", "core.patient", "core.customer_care", "core.billing",
		"laboratory.basic", "clinical.basic", "pharmacy.basic", "inventory.basic", "qms.basic",
	}

	if orgUUID, errP := uuid.Parse(orgID); errP == nil && b.entitlementSvc != nil {
		if realPlan, errPlan := b.entitlementSvc.GetOrganizationPlan(ctx, orgUUID); errPlan == nil && realPlan != "" {
			orgSubPlan = realPlan
		}
		if effCaps, errCaps := b.entitlementSvc.GetEffectiveCapabilities(ctx, orgUUID); errCaps == nil && len(effCaps) > 0 {
			effectiveCapabilities = effCaps
		}
	}

	// Build dynamic limits based on real organization plan code
	resolvedLimits := map[string]int{
		"maxBranches": 1,
		"maxMembers":  5,
		"storageGb":   10,
	}
	switch orgSubPlan {
	case "optimize":
		resolvedLimits = map[string]int{"maxBranches": 3, "maxMembers": 25, "storageGb": 50}
	case "pro":
		resolvedLimits = map[string]int{"maxBranches": 10, "maxMembers": 100, "storageGb": 200}
	case "enterprise":
		resolvedLimits = map[string]int{"maxBranches": 1000, "maxMembers": 10000, "storageGb": 5000}
	}

	if b.dbPool != nil {
		var limitsJSON []byte
		errLim := b.dbPool.QueryRow(ctx, `SELECT limits FROM subscription.plans WHERE code = $1`, orgSubPlan).Scan(&limitsJSON)
		if errLim == nil && len(limitsJSON) > 0 {
			var dbLimits map[string]int
			if errUnmarshal := json.Unmarshal(limitsJSON, &dbLimits); errUnmarshal == nil && len(dbLimits) > 0 {
				resolvedLimits = dbLimits
			}
		}
	}

	// Update moduleCapabilities based on effective capabilities
	effectiveCapSet := make(map[string]bool)
	for _, c := range effectiveCapabilities {
		effectiveCapSet[c] = true
		parts := strings.Split(c, ".")
		if len(parts) > 0 {
			effectiveCapSet[parts[0]] = true
		}
	}

	moduleCapabilities := make([]ModuleCapabilityPayload, 0, len(allModuleCodes))
	for _, code := range allModuleCodes {
		isEnabled := effectiveCapSet[code]
		moduleCapabilities = append(moduleCapabilities, ModuleCapabilityPayload{
			Code:             code,
			Enabled:          isEnabled,
			Licensed:         isEnabled,
			Visible:          true,
			UpgradeAvailable: !isEnabled,
			Actions:          []string{"read", "write"},
		})
	}

	// Resolve dynamic multi-tenant organization branding tokens
	resolvedPrimaryColor := "#0284c7"
	resolvedSecondaryColor := "#0f172a"
	resolvedAccentColor := "#38bdf8"
	resolvedFontFamily := "Outfit"
	resolvedBorderRadius := "0.5rem"
	resolvedLogoURL := activeLogoURL
	resolvedCustomDomain := activeCustomDomain
	var resolvedThemeBranding map[string]interface{}

	if activePrimaryColor != "" {
		resolvedPrimaryColor = activePrimaryColor
	}

	if len(activeThemeBrandingJSON) > 0 {
		_ = json.Unmarshal(activeThemeBrandingJSON, &resolvedThemeBranding)
		if sec, ok := resolvedThemeBranding["secondaryColor"].(string); ok && sec != "" {
			resolvedSecondaryColor = sec
		}
		if acc, ok := resolvedThemeBranding["accentColor"].(string); ok && acc != "" {
			resolvedAccentColor = acc
		}
		if font, ok := resolvedThemeBranding["fontFamily"].(string); ok && font != "" {
			resolvedFontFamily = font
		}
		if rad, ok := resolvedThemeBranding["borderRadius"].(string); ok && rad != "" {
			resolvedBorderRadius = rad
		}
		if logo, ok := resolvedThemeBranding["logoUrl"].(string); ok && logo != "" && resolvedLogoURL == "" {
			resolvedLogoURL = logo
		}
	}

	return &BootstrapContractResponse{
		Identity: IdentityPayload{
			ID:          userID,
			Email:       userEmail,
			DisplayName: displayName,
			Locale:      "en",
			Timezone:    "Africa/Lagos",
		},
		Platform: PlatformPayload{
			IsStaff: isPlatformStaff,
			Role:    userRole,
		},
		Organization: OrganizationPayload{
			ID:           orgID,
			Name:         orgName,
			Slug:         orgSlug,
			Role:         orgRole,
			Subscription: orgSubPlan,
		},
		Workspace: WorkspacePayload{
			ID:           tenantID,
			Name:         tenantName,
			FacilityType: "Laboratory",
			Slug:         tenantSlug,
			Timezone:     "Africa/Lagos",
			Currency:     currency,
		},
		Subscription: SubscriptionPayload{
			Plan:   orgSubPlan,
			Status: "active",
			Limits: resolvedLimits,
		},
		Modules:      moduleCapabilities,
		Capabilities: effectiveCapabilities,
		Permissions:  effectivePermissions,
		Navigation:   navigationItems,
		StructuredNavigation: StructuredNavigationPayload{
			Primary: navigationItems,
		},
		Dashboard: DashboardPayload{
			Widgets: dashboardWidgets,
			Sections: []DashboardSectionPayload{
				{ID: "sec_overview", Title: "Operational Overview", Type: "cards", Items: dashboardWidgets},
			},
		},
		Contexts: ContextsPayload{
			Current:   currentContext,
			Available: availableContexts,
			Default:   "workspace",
		},
		Branding: BrandingPayload{
			LogoURL:        resolvedLogoURL,
			PrimaryColor:   resolvedPrimaryColor,
			SecondaryColor: resolvedSecondaryColor,
			AccentColor:    resolvedAccentColor,
			FontFamily:     resolvedFontFamily,
			BorderRadius:   resolvedBorderRadius,
			CustomDomain:   resolvedCustomDomain,
			ThemeBranding:  resolvedThemeBranding,
		},
		Preferences: PreferencesPayload{
			Theme:    "dark",
			Language: "en",
			Timezone: "Africa/Lagos",
		},
		FeatureFlags: map[string]bool{
			"aiCatalogMapping": true,
			"betaRadiology":    false,
		},
		Limits: resolvedLimits,
		Metadata: MetadataPayload{
			Version:     "2.0",
			GeneratedAt: "2026-08-04T14:00:00Z",
			TTL:         300,
			ETag:        "W/\"b89a-202608041400\"",
		},
	}, nil
}

func (b *BootstrapBuilder) SetEntitlementService(svc *subApp.EntitlementService) {
	b.entitlementSvc = svc
}

