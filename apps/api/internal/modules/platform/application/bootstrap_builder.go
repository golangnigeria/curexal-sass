package application

import (
	"context"
	"encoding/json"
	"strings"

	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	orgDomain "github.com/golangnigeria/curexal/internal/modules/organization/domain"
	platformDomain "github.com/golangnigeria/curexal/internal/modules/platform/domain"
	subApp "github.com/golangnigeria/curexal/internal/modules/subscription/application"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	ID             string `json:"id"`
	Name           string `json:"name"`
	Slug           string `json:"slug,omitempty"`
	Logo           string `json:"logo,omitempty"`
	Role           string `json:"role,omitempty"`
	Subscription   string `json:"subscription"`
	Status         string `json:"status,omitempty"`
	SetupState     string `json:"setupState,omitempty"`
	Hostname       string `json:"hostname,omitempty"`
	IsCustomDomain bool   `json:"isCustomDomain,omitempty"`
}

type BranchPayload struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Slug           string                 `json:"slug"`
	Code           string                 `json:"code"`
	FacilityType   string                 `json:"facilityType"`
	IsHeadquarters bool                   `json:"isHeadquarters"`
	City           string                 `json:"city,omitempty"`
	State          string                 `json:"state,omitempty"`
	OperatingHours map[string]interface{} `json:"operatingHours,omitempty"`
	Status         string                 `json:"status,omitempty"`
}

type BranchSummaryPayload struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Code           string `json:"code"`
	FacilityType   string `json:"facilityType"`
	IsHeadquarters bool   `json:"isHeadquarters"`
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
	Branch               *BranchPayload              `json:"branch,omitempty"`
	Workspace            WorkspacePayload            `json:"workspace"`
	AvailableBranches    []BranchSummaryPayload      `json:"availableBranches,omitempty"`
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
	return b.BuildBootstrapWithContext(ctx, principal, "", "")
}

func (b *BootstrapBuilder) BuildBootstrapWithContext(ctx context.Context, principal *middleware.AuthenticatedPrincipal, reqHost, reqBranchSlug string) (*BootstrapContractResponse, error) {
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
	var activeOrgStatus string
	var activeOrgSetupState string
	var activeMembershipRole string
	var activeMembershipTenantID string
	var activePrimaryColor string
	var activeThemeBrandingJSON []byte
	var activeCustomDomain string
	var activeLogoURL string
	var resolvedHostname string
	var isCustomDomain bool
	var activeTenant *orgDomain.Tenant
	var activeBranch *BranchPayload
	var availableBranches []BranchSummaryPayload

	cleanHost := strings.ToLower(strings.TrimSpace(reqHost))
	if colonIdx := strings.Index(cleanHost, ":"); colonIdx != -1 {
		cleanHost = cleanHost[:colonIdx]
	}

	// Step 1a: Attempt to resolve specific organization by incoming Host header if provided
	var targetOrgIDFromHost string
	if b.dbPool != nil && cleanHost != "" && cleanHost != "localhost" && cleanHost != "127.0.0.1" && cleanHost != "app.curexal.space" && cleanHost != "curexal.space" {
		if strings.HasSuffix(cleanHost, ".localhost") {
			sub := strings.TrimSuffix(cleanHost, ".localhost")
			if sub != "app" && sub != "api" && sub != "public" {
				_ = b.dbPool.QueryRow(ctx, `SELECT id::text FROM organization.organizations WHERE slug = $1 LIMIT 1`, sub).Scan(&targetOrgIDFromHost)
				resolvedHostname = cleanHost
			}
		} else if strings.HasSuffix(cleanHost, ".curexal.space") || strings.HasSuffix(cleanHost, ".curexal.internal") {
			sub := strings.TrimSuffix(cleanHost, ".curexal.space")
			sub = strings.TrimSuffix(sub, ".curexal.internal")
			if sub != "app" && sub != "api" && sub != "public" && sub != "admin" {
				_ = b.dbPool.QueryRow(ctx, `
					SELECT o.id::text FROM organization.organizations o
					LEFT JOIN organization.organization_domains od ON od.organization_id = o.id
					WHERE od.hostname = $1 OR o.slug = $2 LIMIT 1
				`, cleanHost, sub).Scan(&targetOrgIDFromHost)
				resolvedHostname = cleanHost
			}
		} else {
			// Custom domain lookup
			_ = b.dbPool.QueryRow(ctx, `
				SELECT o.id::text FROM organization.organization_domains od
				JOIN organization.organizations o ON o.id = od.organization_id
				WHERE od.hostname = $1 AND od.is_verified = TRUE LIMIT 1
			`, cleanHost).Scan(&targetOrgIDFromHost)
			if targetOrgIDFromHost != "" {
				resolvedHostname = cleanHost
				isCustomDomain = true
			}
		}
	}

	if b.dbPool != nil && (userID != "" || userEmail != "") {
		var row pgx.Row
		if targetOrgIDFromHost != "" {
			row = b.dbPool.QueryRow(ctx, `
				SELECT 
					m.organization_id::text, 
					COALESCE(m.tenant_id::text, ''), 
					CASE 
						WHEN m.role IN ('owner', 'org_admin', 'admin', 'org_regional_manager', 'org_quality_manager', 'org_finance_manager', 'org_hr_manager') THEN m.role
						ELSE COALESCE(NULLIF(m.role, 'member'), m.role, 'member')
					END, 
					o.name, 
					o.slug, 
					COALESCE(o.plan, 'smart'),
					COALESCE(o.status, 'active'),
					COALESCE(o.setup_state, ''),
					COALESCE(o.primary_color, '#0284c7'),
					COALESCE(o.theme_branding, '{}'::jsonb),
					COALESCE(o.custom_domain, ''),
					COALESCE(o.logo_url, '')
				FROM organization.organization_memberships m
				JOIN organization.organizations o ON o.id = m.organization_id
				LEFT JOIN identity.users u ON u.id = m.user_id
				WHERE (m.user_id::text = $1 OR u.id::text = $1 OR u.email = $2 OR ($2 != '' AND u.email ILIKE $2)) 
				  AND m.organization_id::text = $3 
				  AND m.is_active = TRUE
				LIMIT 1
			`, userID, userEmail, targetOrgIDFromHost)
		} else {
			row = b.dbPool.QueryRow(ctx, `
				SELECT 
					m.organization_id::text, 
					COALESCE(m.tenant_id::text, ''), 
					CASE 
						WHEN m.role IN ('owner', 'org_admin', 'admin', 'org_regional_manager', 'org_quality_manager', 'org_finance_manager', 'org_hr_manager') THEN m.role
						ELSE COALESCE(NULLIF(m.role, 'member'), m.role, 'member')
					END, 
					o.name, 
					o.slug, 
					COALESCE(o.plan, 'smart'),
					COALESCE(o.status, 'active'),
					COALESCE(o.setup_state, ''),
					COALESCE(o.primary_color, '#0284c7'),
					COALESCE(o.theme_branding, '{}'::jsonb),
					COALESCE(o.custom_domain, ''),
					COALESCE(o.logo_url, '')
				FROM organization.organization_memberships m
				JOIN organization.organizations o ON o.id = m.organization_id
				LEFT JOIN identity.users u ON u.id = m.user_id
				WHERE (m.user_id::text = $1 OR u.id::text = $1 OR u.email = $2 OR ($2 != '' AND u.email ILIKE $2)) 
				  AND m.is_active = TRUE
				ORDER BY (m.role IN ('owner', 'org_admin')) DESC, m.created_at ASC
				LIMIT 1
			`, userID, userEmail)
		}

		_ = row.Scan(
			&activeOrgID,
			&activeMembershipTenantID,
			&activeMembershipRole,
			&activeOrgName,
			&activeOrgSlug,
			&activeOrgPlan,
			&activeOrgStatus,
			&activeOrgSetupState,
			&activePrimaryColor,
			&activeThemeBrandingJSON,
			&activeCustomDomain,
			&activeLogoURL,
		)

		if activeOrgID == "" && targetOrgIDFromHost != "" {
			_ = b.dbPool.QueryRow(ctx, `
				SELECT id::text, name, slug, COALESCE(plan, 'smart'), COALESCE(status, 'active'), COALESCE(setup_state, ''), COALESCE(primary_color, '#0284c7'), COALESCE(theme_branding, '{}'::jsonb), COALESCE(custom_domain, ''), COALESCE(logo_url, '')
				FROM organization.organizations WHERE id::text = $1 LIMIT 1
			`, targetOrgIDFromHost).Scan(&activeOrgID, &activeOrgName, &activeOrgSlug, &activeOrgPlan, &activeOrgStatus, &activeOrgSetupState, &activePrimaryColor, &activeThemeBrandingJSON, &activeCustomDomain, &activeLogoURL)
		}

		if !isPlatformStaff && userRole == "" {
			var dbIsAdmin bool
			var dbPlatformRole *string
			if errUser := b.dbPool.QueryRow(ctx, `SELECT is_platform_admin, platform_role FROM identity.users WHERE id::text = $1 OR email = $2 LIMIT 1`, userID, userEmail).Scan(&dbIsAdmin, &dbPlatformRole); errUser == nil {
				if dbIsAdmin || (dbPlatformRole != nil && (*dbPlatformRole == "super_admin" || *dbPlatformRole == "platform_admin" || *dbPlatformRole == "platform_staff")) {
					isPlatformStaff = true
					if dbPlatformRole != nil && *dbPlatformRole != "" {
						userRole = *dbPlatformRole
					} else {
						userRole = "super_admin"
					}
				}
			}
		}
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
			if orgs[0].Status != "" {
				activeOrgStatus = string(orgs[0].Status)
			}
			if orgs[0].SetupState != "" {
				activeOrgSetupState = string(orgs[0].SetupState)
			}
			if activeMembershipRole == "" && userRole != "" && userRole != "user" && userRole != "member" {
				activeMembershipRole = userRole
			}
		}
	}

	if resolvedHostname == "" && activeOrgSlug != "" {
		resolvedHostname = activeOrgSlug + ".curexal.space"
	}

	// 2. Resolve Active Context
	currentContext := "workspace"
	availableContexts := []string{"workspace"}

	if isPlatformStaff {
		currentContext = "platform"
		availableContexts = []string{"platform", "organization", "workspace"}
	} else if activeOrgID != "" {
		isExecRole := activeMembershipRole == "owner" || activeMembershipRole == "org_admin" || activeMembershipRole == "admin" || activeMembershipRole == "org_regional_manager" || activeMembershipRole == "org_quality_manager" || activeMembershipRole == "org_finance_manager" || activeMembershipRole == "org_hr_manager" || userRole == "owner" || userRole == "org_admin"
		if reqBranchSlug != "" {
			currentContext = "workspace"
			if isExecRole {
				availableContexts = []string{"organization", "workspace"}
			} else {
				availableContexts = []string{"workspace"}
			}
		} else if isExecRole {
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

	// 3. Resolve Facility Branches & Active Branch / Workspace details
	targetTenantID := ""
	if principal != nil && principal.TenantID != "" {
		targetTenantID = principal.TenantID
	} else if activeMembershipTenantID != "" {
		targetTenantID = activeMembershipTenantID
	}

	var activeBranchThemeBrandingJSON []byte
	if b.dbPool != nil && activeOrgID != "" {
		bRows, bErr := b.dbPool.Query(ctx, `
			SELECT b.id::text, b.name, b.code, COALESCE(b.slug, b.code), ft.name, b.is_headquarters, COALESCE(b.city, ''), COALESCE(b.state, ''), b.operating_hours, COALESCE(b.theme_branding, '{}'::jsonb), b.status
			FROM organization.facility_branches b
			JOIN platform.facility_types ft ON ft.id = b.facility_type_id
			WHERE b.organization_id = $1
			ORDER BY b.is_headquarters DESC, b.name ASC
		`, activeOrgID)
		if bErr == nil {
			defer bRows.Close()
			for bRows.Next() {
				var (
					bID, bName, bCode, bSlug, ftName, bCity, bState, bStatus string
					isHQ                                                     bool
					opHoursJSON, themeJSON                                   []byte
				)
				if errScan := bRows.Scan(&bID, &bName, &bCode, &bSlug, &ftName, &isHQ, &bCity, &bState, &opHoursJSON, &themeJSON, &bStatus); errScan == nil {
					if bStatus == "ACTIVE" {
						availableBranches = append(availableBranches, BranchSummaryPayload{
							ID:             bID,
							Name:           bName,
							Slug:           bSlug,
							Code:           bCode,
							FacilityType:   ftName,
							IsHeadquarters: isHQ,
						})
					}

					// Check if this branch matches the requested branch slug, code, or tenant ID
					isTargetBranch := false
					if reqBranchSlug != "" && (strings.EqualFold(bSlug, reqBranchSlug) || strings.EqualFold(bCode, reqBranchSlug)) {
						isTargetBranch = true
					} else if targetTenantID != "" && bID == targetTenantID {
						isTargetBranch = true
					} else if activeBranch == nil && isHQ && bStatus == "ACTIVE" {
						isTargetBranch = true
					}

					if isTargetBranch || (activeBranch == nil && bStatus == "ACTIVE") {
						var parsedHours map[string]interface{}
						if len(opHoursJSON) > 0 {
							_ = json.Unmarshal(opHoursJSON, &parsedHours)
						}
						activeBranch = &BranchPayload{
							ID:             bID,
							Name:           bName,
							Slug:           bSlug,
							Code:           bCode,
							FacilityType:   ftName,
							IsHeadquarters: isHQ,
							City:           bCity,
							State:          bState,
							OperatingHours: parsedHours,
							Status:         bStatus,
						}
						activeBranchThemeBrandingJSON = themeJSON
					}
				}
			}
		}
	}

	isOrgAdminOrOwner := isPlatformStaff || userRole == "owner" || userRole == "org_admin" || activeMembershipRole == "owner" || activeMembershipRole == "org_admin" || activeMembershipRole == "org_regional_manager" || activeMembershipRole == "admin"

	// Filter available branches for branch-only staff members to ensure multi-tenant facility isolation
	if !isPlatformStaff && !isOrgAdminOrOwner && targetTenantID != "" {
		filteredBranches := make([]BranchSummaryPayload, 0)
		for _, b := range availableBranches {
			if b.ID == targetTenantID {
				filteredBranches = append(filteredBranches, b)
			}
		}
		if len(filteredBranches) > 0 {
			availableBranches = filteredBranches
			if activeBranch != nil && activeBranch.ID != targetTenantID {
				for _, fb := range filteredBranches {
					if fb.ID == targetTenantID {
						activeBranch.ID = fb.ID
						activeBranch.Name = fb.Name
						activeBranch.Slug = fb.Slug
						activeBranch.Code = fb.Code
						activeBranch.FacilityType = fb.FacilityType
						activeBranch.IsHeadquarters = fb.IsHeadquarters
						break
					}
				}
			}
		}
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
	facilityTypeVal := "Laboratory"
	currency := "NGN"
	enabledModules := []string{"laboratory", "clinical", "pharmacy", "billing", "inventory", "customer_care", "qms"}

	if activeBranch != nil {
		tenantID = activeBranch.ID
		tenantName = activeBranch.Name
		tenantSlug = activeBranch.Slug
		if activeBranch.FacilityType != "" {
			facilityTypeVal = activeBranch.FacilityType
		}
		// Adjust enabled modules according to facility type blueprint
		normFT := strings.ToLower(activeBranch.FacilityType)
		if strings.Contains(normFT, "clinic") || strings.Contains(normFT, "outpatient") || strings.Contains(normFT, "emr") {
			enabledModules = []string{"clinical", "customer_care", "billing", "pharmacy"}
		} else if strings.Contains(normFT, "pharmacy") {
			enabledModules = []string{"pharmacy", "inventory", "billing"}
		} else if strings.Contains(normFT, "radiology") || strings.Contains(normFT, "imaging") {
			enabledModules = []string{"radiology", "customer_care", "billing"}
		} else if strings.Contains(normFT, "hospital") {
			enabledModules = []string{"hospital", "clinical", "laboratory", "pharmacy", "radiology", "billing", "inventory", "customer_care", "qms"}
		} else if strings.Contains(normFT, "lab") || strings.Contains(normFT, "diagnostic") {
			enabledModules = []string{"laboratory", "customer_care", "billing", "qms"}
		}
	} else if activeTenant != nil {
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
	isWorkspaceAdmin := isSuperAdminOrOwner || (currentContext == "workspace" && (userRole == "branch_admin" || activeMembershipRole == "branch_admin" || userRole == "branch_manager"))

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
			JOIN organization.organization_memberships m ON (m.role = r.code OR m.role = r.name)
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
		if currentContext == "organization" {
			effectivePermissions = []string{
				"organization:read",
				"organization:view",
				"organization:manage",
				"organization:branch:read",
				"organization:branch:write",
				"users:read",
				"users:write",
				"organization:catalog:read",
				"organization:catalog:write",
				"organization:branding:read",
				"organization:branding:write",
				"organization:notifications:read",
				"organization:notifications:write",
				"organization:integrations:read",
				"organization:integrations:write",
				"organization:audit:read",
				"organization:settings:read",
				"organization:settings:write",
				"organization:document:upload",
				"organization:document:read",
				"audit:read",
			}
		} else {
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
	}

	// 5. Context-Bound Navigation Tree (100% DB-Driven & Permission Filtered)
	var navigationItems []NavigationItemPayload
	if b.dbPool != nil {
		activeBranchSlug := reqBranchSlug
		if activeBranchSlug == "" {
			activeBranchSlug = tenantSlug
		}
		if activeBranchSlug == "" {
			activeBranchSlug = "main"
		}

		branchCapSet := make(map[string]bool)
		if currentContext == "workspace" && activeBranch != nil && activeBranch.ID != "" {
			fRows, fErr := b.dbPool.Query(ctx, `
				SELECT c.code
				FROM platform.facility_capabilities fc
				JOIN subscription.capabilities c ON c.id = fc.capability_id
				WHERE fc.facility_type_id = (
					SELECT facility_type_id FROM organization.facility_branches WHERE id = $1
				)
			`, activeBranch.ID)
			if fErr == nil {
				defer fRows.Close()
				for fRows.Next() {
					var capCode string
					if errScan := fRows.Scan(&capCode); errScan == nil {
						branchCapSet[capCode] = true
					}
				}
			}
		}

		rows, err := b.dbPool.Query(ctx, `
			SELECT id, title, icon, path, sort_order, required_capability
			FROM navigation_item
			WHERE context_scope = $1
			  AND is_active = true
			  AND is_visible = true
			  AND (module_code IS NULL OR module_code = ANY($2) OR $4 = TRUE)
			  AND (required_permission IS NULL OR required_permission = ANY($3) OR $4 = TRUE)
			ORDER BY sort_order ASC
		`, currentContext, enabledModules, effectivePermissions, isWorkspaceAdmin)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var item NavigationItemPayload
				var reqCap *string
				if err := rows.Scan(&item.ID, &item.Title, &item.Icon, &item.Path, &item.Order, &reqCap); err == nil {
					if currentContext == "workspace" && reqCap != nil && *reqCap != "" && len(branchCapSet) > 0 && !branchCapSet[*reqCap] {
						continue
					}
					if currentContext == "workspace" {
						item.Path = strings.Replace(item.Path, "/:branch", "/"+activeBranchSlug, 1)
						if strings.HasPrefix(item.Path, "/workspace/") {
							item.Path = strings.Replace(item.Path, "/workspace/", "/"+activeBranchSlug+"/", 1)
						}
					}
					navigationItems = append(navigationItems, item)
				}
			}
		}
	} else {
		// In-memory domain fallbacks for unit test suites without a live PostgreSQL connection
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
	orgStatus := activeOrgStatus
	orgSetupState := activeOrgSetupState
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
			if selectedOrg.Status != "" {
				orgStatus = string(selectedOrg.Status)
			}
			if selectedOrg.SetupState != "" {
				orgSetupState = string(selectedOrg.SetupState)
			}
		}
	}

	if orgID == "" || orgID == "org_default" {
		if b.dbPool != nil {
			var firstOrgID, firstOrgName, firstOrgSlug, firstOrgPlan, firstOrgStatus, firstOrgSetupState string
			errFirst := b.dbPool.QueryRow(ctx, `SELECT id::text, name, slug, COALESCE(plan, 'smart'), COALESCE(status, 'active'), COALESCE(setup_state, 'VERIFIED') FROM organization.organizations ORDER BY created_at ASC LIMIT 1`).Scan(&firstOrgID, &firstOrgName, &firstOrgSlug, &firstOrgPlan, &firstOrgStatus, &firstOrgSetupState)
			if errFirst == nil && firstOrgID != "" {
				orgID = firstOrgID
				orgName = firstOrgName
				orgSlug = firstOrgSlug
				orgSubPlan = firstOrgPlan
				orgStatus = firstOrgStatus
				orgSetupState = firstOrgSetupState
				if orgRole == "" {
					orgRole = "owner"
				}
			}
		}
	}
	if orgID == "" {
		orgID = "00000000-0000-0000-0000-000000000001"
		orgName = "Curexal Health Network"
		orgSlug = "curexal"
		orgSubPlan = "smart"
		orgStatus = "active"
		orgSetupState = "VERIFIED"
		orgRole = userRole
	}

	// 7. Structured Navigation (Primary, Topbar, QuickActions, Breadcrumbs)
	var breadcrumbItems []BreadcrumbPayload
	var topbarItems []NavigationItemPayload
	var quickActionItems []NavigationItemPayload

	if currentContext == "platform" {
		breadcrumbItems = []BreadcrumbPayload{
			{Title: "Platform Console", Path: "/platform/dashboard"},
		}
	} else if currentContext == "organization" {
		breadcrumbItems = []BreadcrumbPayload{
			{Title: orgName, Path: "/organization/dashboard"},
		}
		topbarItems = []NavigationItemPayload{
			{ID: "top_branches", Title: "Branch Facilities", Icon: "Building2", Path: "/organization/branches"},
			{ID: "top_members", Title: "Staff Roster", Icon: "Users", Path: "/organization/members"},
			{ID: "top_billing", Title: "Subscription", Icon: "CreditCard", Path: "/organization/billing"},
		}
		quickActionItems = []NavigationItemPayload{
			{ID: "qa_provision_branch", Title: "Provision Branch", Icon: "Building2", Path: "/organization/branches"},
			{ID: "qa_invite_member", Title: "Invite Staff Member", Icon: "Users", Path: "/organization/members"},
		}
	} else {
		// Workspace Context
		if isOrgAdminOrOwner {
			breadcrumbItems = []BreadcrumbPayload{
				{Title: orgName, Path: "/organization/dashboard"},
				{Title: tenantName, Path: "/" + tenantSlug + "/dashboard"},
			}
			topbarItems = []NavigationItemPayload{
				{ID: "top_return_hq", Title: "Executive Org HQ", Icon: "Building2", Path: "/organization/dashboard"},
				{ID: "top_facility_settings", Title: "Facility Settings", Icon: "Settings", Path: "/organization/branches"},
			}
		} else {
			breadcrumbItems = []BreadcrumbPayload{
				{Title: tenantName, Path: "/" + tenantSlug + "/dashboard"},
			}
		}

		if enabledSet["laboratory"] {
			quickActionItems = append(quickActionItems, NavigationItemPayload{
				ID: "qa_accession_sample", Title: "Accession Specimen", Icon: "Activity", Path: "/" + tenantSlug + "/laboratory",
			})
		}
		if enabledSet["clinical"] {
			quickActionItems = append(quickActionItems, NavigationItemPayload{
				ID: "qa_queue_patient", Title: "Consultation Queue", Icon: "Stethoscope", Path: "/" + tenantSlug + "/clinical",
			})
		}
		if enabledSet["pharmacy"] {
			quickActionItems = append(quickActionItems, NavigationItemPayload{
				ID: "qa_dispense_rx", Title: "Dispense Medication", Icon: "Pill", Path: "/" + tenantSlug + "/pharmacy",
			})
		}
		if enabledSet["billing"] {
			quickActionItems = append(quickActionItems, NavigationItemPayload{
				ID: "qa_pos_bill", Title: "Cashier POS Invoice", Icon: "CreditCard", Path: "/" + tenantSlug + "/billing",
			})
		}
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

	// Intersect with Facility Blueprint Capabilities if in facility workspace context
	if activeBranch != nil && activeBranch.ID != "" && b.dbPool != nil {
		fRows, fErr := b.dbPool.Query(ctx, `
			SELECT c.code
			FROM platform.facility_capabilities fc
			JOIN subscription.capabilities c ON c.id = fc.capability_id
			WHERE fc.facility_type_id = (
				SELECT facility_type_id FROM organization.facility_branches WHERE id = $1
			)
		`, activeBranch.ID)
		if fErr == nil {
			defer fRows.Close()
			blueprintCapSet := map[string]bool{
				"core.organization":  true,
				"core.patient":       true,
				"core.billing":       true,
				"core.customer_care": true,
			}
			for fRows.Next() {
				var capCode string
				if errScan := fRows.Scan(&capCode); errScan == nil {
					blueprintCapSet[capCode] = true
				}
			}

			intersected := make([]string, 0)
			for _, capCode := range effectiveCapabilities {
				if blueprintCapSet[capCode] {
					intersected = append(intersected, capCode)
				}
			}
			if len(intersected) > 0 {
				effectiveCapabilities = intersected
			}
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

	// Apply Facility-level Theme Overrides on top of Organization theme
	if activeBranchThemeBrandingJSON != nil && len(activeBranchThemeBrandingJSON) > 0 {
		var facilityTheme map[string]interface{}
		if errUn := json.Unmarshal(activeBranchThemeBrandingJSON, &facilityTheme); errUn == nil {
			useOrgBranding, _ := facilityTheme["useOrganizationBranding"].(bool)
			if !useOrgBranding {
				if prim, ok := facilityTheme["primaryColor"].(string); ok && prim != "" {
					resolvedPrimaryColor = prim
				}
				if sec, ok := facilityTheme["secondaryColor"].(string); ok && sec != "" {
					resolvedSecondaryColor = sec
				}
				if acc, ok := facilityTheme["accentColor"].(string); ok && acc != "" {
					resolvedAccentColor = acc
				}
				if font, ok := facilityTheme["fontFamily"].(string); ok && font != "" {
					resolvedFontFamily = font
				}
				if rad, ok := facilityTheme["borderRadius"].(string); ok && rad != "" {
					resolvedBorderRadius = rad
				}
			}
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
			ID:             orgID,
			Name:           orgName,
			Slug:           orgSlug,
			Role:           orgRole,
			Subscription:   orgSubPlan,
			Status:         orgStatus,
			SetupState:     orgSetupState,
			Hostname:       resolvedHostname,
			IsCustomDomain: isCustomDomain,
		},
		Branch: activeBranch,
		Workspace: WorkspacePayload{
			ID:           tenantID,
			Name:         tenantName,
			FacilityType: facilityTypeVal,
			Slug:         tenantSlug,
			Timezone:     "Africa/Lagos",
			Currency:     currency,
		},
		AvailableBranches: availableBranches,
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
			Primary:      navigationItems,
			Topbar:       topbarItems,
			QuickActions: quickActionItems,
			Breadcrumbs:  breadcrumbItems,
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
