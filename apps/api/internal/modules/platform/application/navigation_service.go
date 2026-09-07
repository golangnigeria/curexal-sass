package application

import (
	"context"
	"strings"

	kernelAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	orgDomain "github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/modules/platform/domain"
	subApp "github.com/golangnigeria/curexal/internal/modules/subscription/application"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NavigationService struct {
	navRepo        domain.NavigationRepository
	orgRepo        orgDomain.OrganizationRepository
	tenantRepo     orgDomain.TenantRepository
	subRepo        orgDomain.SubscriptionRepository
	entitlementSvc *subApp.EntitlementService
	dbPool         *pgxpool.Pool
}

func NewNavigationService(
	navRepo domain.NavigationRepository,
	orgRepo orgDomain.OrganizationRepository,
	tenantRepo orgDomain.TenantRepository,
	subRepo orgDomain.SubscriptionRepository,
) *NavigationService {
	return &NavigationService{
		navRepo:    navRepo,
		orgRepo:    orgRepo,
		tenantRepo: tenantRepo,
		subRepo:    subRepo,
	}
}

func (s *NavigationService) SetEntitlementService(svc *subApp.EntitlementService) {
	s.entitlementSvc = svc
}

func (s *NavigationService) SetDBPool(pool *pgxpool.Pool) {
	s.dbPool = pool
}

func (s *NavigationService) GetNavigation(
	ctx context.Context,
	principal *kernelAuth.AuthenticatedPrincipal,
	reqHost string,
	branchSlug string,
	reqScope string,
) (*domain.NavigationResponse, error) {
	cleanHost := strings.ToLower(reqHost)
	if idx := strings.Index(cleanHost, ":"); idx != -1 {
		cleanHost = cleanHost[:idx]
	}

	isPlatformStaff := false
	isSuperAdminOrOwner := false
	var effectivePermissions []string
	var enabledModules []string
	var userRole string

	if principal != nil {
		isPlatformStaff = principal.Platform.IsPlatformAdmin || principal.Platform.IsSuperAdmin || principal.Platform.IsPlatformStaff || principal.Role == "super_admin" || principal.Role == "platform_admin" || principal.Role == "platform_staff"
		userRole = principal.Role
		if userRole == "" && principal.Organization.OrganizationRole != "" {
			userRole = principal.Organization.OrganizationRole
		}
		if userRole == "" && principal.Platform.PlatformRole != "" {
			userRole = principal.Platform.PlatformRole
		}
		effectivePermissions = principal.Permissions
		isSuperAdminOrOwner = isPlatformStaff || userRole == "owner" || userRole == "org_admin"
	}

	// Sanitize pseudo branch slugs passed by route path extractors
	rawBranch := strings.TrimSpace(branchSlug)
	cleanBranch := rawBranch
	lowerBranch := strings.ToLower(rawBranch)
	if lowerBranch == "platform" || lowerBranch == "organization" || lowerBranch == "workspace" || lowerBranch == "login" || lowerBranch == "portal" || lowerBranch == "api" {
		cleanBranch = ""
	}

	// 1. Explicit Scope Precedence Evaluation
	currentScope := "workspace"
	requestedScope := strings.ToLower(strings.TrimSpace(reqScope))

	if requestedScope != "" {
		switch requestedScope {
		case "platform":
			if !isPlatformStaff {
				return nil, domain.ErrUnauthorizedScope
			}
			currentScope = "platform"
		case "organization":
			if !isSuperAdminOrOwner {
				hasOrgAccess := userRole == "owner" || userRole == "org_admin" || userRole == "admin" || userRole == "org_regional_manager" || userRole == "org_quality_manager" || userRole == "org_finance_manager" || userRole == "org_hr_manager"
				if !hasOrgAccess {
					for _, p := range effectivePermissions {
						if p == "organization:read" || p == "organization:manage" || p == "organization:*" || p == "*" {
							hasOrgAccess = true
							break
						}
					}
				}
				if !hasOrgAccess {
					return nil, domain.ErrUnauthorizedScope
				}
			}
			currentScope = "organization"
		case "workspace":
			currentScope = "workspace"
		case "patient":
			currentScope = "patient"
		default:
			return nil, domain.ErrInvalidScope
		}
	} else {
		isExec := isSuperAdminOrOwner || userRole == "owner" || userRole == "org_admin" || userRole == "admin" || userRole == "org_regional_manager" || userRole == "org_quality_manager" || userRole == "org_finance_manager" || userRole == "org_hr_manager"
		// Context Inference by Host, Path and Session Vectors
		if isPatientHost(cleanHost) {
			currentScope = "patient"
		} else if lowerBranch == "organization" {
			if isExec {
				currentScope = "organization"
			} else {
				currentScope = "workspace"
			}
		} else if isOrgDomain(cleanHost) {
			if cleanBranch != "" {
				currentScope = "workspace"
			} else if isExec {
				currentScope = "organization"
			} else {
				currentScope = "workspace"
			}
		} else if isPlatformHost(cleanHost) {
			if isPlatformStaff && cleanBranch == "" {
				currentScope = "platform"
			} else if principal != nil && principal.OrganizationID != "" && cleanBranch == "" && isExec {
				currentScope = "organization"
			} else if cleanBranch != "" {
				currentScope = "workspace"
			} else if isPlatformStaff {
				currentScope = "platform"
			} else {
				currentScope = "workspace"
			}
		} else {
			// Localhost / Development IP fallback
			if isPlatformStaff && cleanBranch == "" {
				currentScope = "platform"
			} else if isExec && cleanBranch == "" {
				currentScope = "organization"
			} else if cleanBranch != "" {
				currentScope = "workspace"
			} else {
				currentScope = "workspace"
			}
		}
	}

	// 2. Resolve Active Branch Slug, Physical Branch Type & Dynamic Facility Capabilities from Database
	activeBranch := cleanBranch
	var branchType string
	var branchCapabilities []string
	if currentScope == "workspace" {
		if activeBranch == "" && principal != nil && principal.Workspace.WorkspaceName != "" {
			activeBranch = principal.Workspace.WorkspaceName
		}
		if activeBranch == "" {
			activeBranch = "main"
		}

		// Resolve physical facility type and capabilities from database
		if s.dbPool != nil {
			orgIDFilter := ""
			if principal != nil && principal.OrganizationID != "" {
				orgIDFilter = principal.OrganizationID
			}

			// 2a. Query physical facility type code and name
			var ftCode string
			_ = s.dbPool.QueryRow(ctx, `
				SELECT COALESCE(ft.code, 'clinic'), ft.name
				FROM organization.facility_branches b
				JOIN platform.facility_types ft ON ft.id = b.facility_type_id
				WHERE (b.slug = $1 OR b.code = $1 OR b.id::text = $1 OR (b.is_headquarters = TRUE AND ($2 != '' AND b.organization_id::text = $2)))
				  AND ($2 = '' OR b.organization_id::text = $2)
				LIMIT 1
			`, activeBranch, orgIDFilter).Scan(&ftCode, &branchType)

			// 2b. Query authoritative facility type capabilities from platform.facility_capabilities
			fcRows, fcErr := s.dbPool.Query(ctx, `
				SELECT DISTINCT c.code
				FROM platform.facility_capabilities fc
				JOIN subscription.capabilities c ON c.id = fc.capability_id
				JOIN organization.facility_branches b ON b.facility_type_id = fc.facility_type_id
				WHERE (b.slug = $1 OR b.code = $1 OR b.id::text = $1 OR (b.is_headquarters = TRUE AND ($2 != '' AND b.organization_id::text = $2)))
				  AND ($2 = '' OR b.organization_id::text = $2)
			`, activeBranch, orgIDFilter)
			if fcErr == nil {
				defer fcRows.Close()
				for fcRows.Next() {
					var capCode string
					if errScan := fcRows.Scan(&capCode); errScan == nil {
						branchCapabilities = append(branchCapabilities, capCode)
					}
				}
			}

			if len(branchCapabilities) == 0 && ftCode != "" {
				branchCapabilities = resolveBranchCapabilitiesByCode(ftCode)
			}
		}

		if len(branchCapabilities) == 0 {
			norm := strings.ToLower(activeBranch)
			if strings.Contains(norm, "clinic") || strings.Contains(norm, "outpatient") {
				branchCapabilities = resolveBranchCapabilitiesByCode("clinic")
				if branchType == "" {
					branchType = "Outpatient Clinic"
				}
			} else if strings.Contains(norm, "lab") {
				branchCapabilities = resolveBranchCapabilitiesByCode("laboratory")
				if branchType == "" {
					branchType = "Medical Laboratory"
				}
			} else if strings.Contains(norm, "radiology") {
				branchCapabilities = resolveBranchCapabilitiesByCode("radiology")
				if branchType == "" {
					branchType = "Radiology Center"
				}
			} else if strings.Contains(norm, "pharmacy") {
				branchCapabilities = resolveBranchCapabilitiesByCode("pharmacy")
				if branchType == "" {
					branchType = "Community Pharmacy"
				}
			}
		}
	} else {
		activeBranch = ""
	}

	// 3. Resolve Enabled Modules & Subscribed Capabilities
	var capabilities []string
	if s.entitlementSvc != nil && principal != nil && principal.OrganizationID != "" {
		if parsedOrgUUID, err := uuid.Parse(principal.OrganizationID); err == nil {
			activeCaps, err := s.entitlementSvc.GetEffectiveCapabilities(ctx, parsedOrgUUID)
			if err == nil {
				capabilities = activeCaps
				for _, capCode := range activeCaps {
					parts := strings.Split(capCode, ".")
					if len(parts) > 0 {
						enabledModules = append(enabledModules, parts[0])
					}
				}
			}
		}
	}

	if len(enabledModules) == 0 {
		enabledModules = []string{"dashboard", "reception", "care_desk", "clinical", "laboratory", "pharmacy", "radiology", "hospital", "billing"}
	}

	// Intersect with Branch-Type Allowed Product Modules for Workspace Scope
	branchCapSet := make(map[string]bool)
	for _, c := range branchCapabilities {
		branchCapSet[c] = true
	}

	if currentScope == "workspace" {
		if len(branchCapabilities) > 0 {
			var branchModules []string
			branchModules = append(branchModules, "dashboard")
			if branchCapSet["core.patient"] || branchCapSet["core.customer_care"] {
				branchModules = append(branchModules, "reception", "customer_care")
			}
			if branchCapSet["clinical.basic"] {
				branchModules = append(branchModules, "care_desk", "clinical")
			}
			if branchCapSet["core.billing"] {
				branchModules = append(branchModules, "billing")
			}
			if branchCapSet["laboratory.basic"] {
				branchModules = append(branchModules, "laboratory")
			}
			if branchCapSet["pharmacy.basic"] {
				branchModules = append(branchModules, "pharmacy")
			}
			if branchCapSet["radiology.basic"] {
				branchModules = append(branchModules, "radiology")
			}
			if branchCapSet["clinical.inpatient_wards"] {
				branchModules = append(branchModules, "hospital")
			}
			enabledModules = branchModules
		} else if branchType != "" {
			// Fallback to in-memory blueprint mapper for test environments without live DB
			branchTypeAllowed := resolveBranchTypeModules(branchType)
			btSet := make(map[string]bool)
			for _, m := range branchTypeAllowed {
				btSet[m] = true
			}
			var filteredModules []string
			for _, m := range enabledModules {
				if btSet[m] {
					filteredModules = append(filteredModules, m)
				}
			}
			enabledModules = filteredModules
		}
	}

	isWorkspaceAdmin := isSuperAdminOrOwner || (currentScope == "workspace" && (userRole == "branch_admin" || userRole == "branch_manager"))

	// 4. Fetch Authoritative Items from Database
	var rawItems []domain.NavigationItem
	var err error
	if s.navRepo != nil {
		rawItems, err = s.navRepo.GetNavigationItemsByScope(ctx, currentScope, enabledModules, effectivePermissions, isWorkspaceAdmin)
		if err != nil {
			return nil, err
		}
	}

	// 5. Capability and Permission Filtering + Dynamic Branch Route Interpolation
	capabilitySet := make(map[string]bool)
	for _, c := range capabilities {
		capabilitySet[c] = true
	}

	var finalItems []domain.NavigationItem
	for _, item := range rawItems {
		// Filter out items requiring specific capability if facility branch or organization lacks entitlement
		// This applies to ALL personas including owners/admins: an Outpatient Clinic never shows LIS or RIS.
		if item.RequiredCapability != nil && *item.RequiredCapability != "" {
			reqCap := *item.RequiredCapability
			if currentScope == "workspace" && len(branchCapSet) > 0 && !branchCapSet[reqCap] {
				continue
			}
			if !isWorkspaceAdmin && len(capabilities) > 0 && !capabilitySet[reqCap] {
				continue
			}
		}

		itemCopy := item
		if currentScope == "workspace" && activeBranch != "" {
			itemCopy.Path = strings.Replace(itemCopy.Path, "/:branch", "/"+activeBranch, 1)
			if strings.HasPrefix(itemCopy.Path, "/workspace/") {
				itemCopy.Path = strings.Replace(itemCopy.Path, "/workspace/", "/"+activeBranch+"/", 1)
			}
		}
		finalItems = append(finalItems, itemCopy)
	}

	// 6. Assemble Response
	navContext := domain.NavigationContext{
		Type:         currentScope,
		Role:         &userRole,
		Capabilities: capabilities,
	}
	if activeBranch != "" {
		navContext.BranchSlug = &activeBranch
	}
	if branchType != "" {
		navContext.BranchType = &branchType
	}
	if principal != nil && principal.OrganizationID != "" {
		orgID := principal.OrganizationID
		navContext.OrganizationID = &orgID
	}

	return &domain.NavigationResponse{
		Context: navContext,
		Items:   finalItems,
	}, nil
}

func resolveBranchCapabilitiesByCode(code string) []string {
	norm := strings.ToLower(code)
	if strings.Contains(norm, "lab") || strings.Contains(norm, "pathology") {
		return []string{"core.organization", "core.patient", "core.billing", "laboratory.basic"}
	}
	if strings.Contains(norm, "radiology") || strings.Contains(norm, "imaging") || strings.Contains(norm, "pacs") {
		return []string{"core.organization", "core.patient", "core.billing", "radiology.basic"}
	}
	if strings.Contains(norm, "pharmacy") || strings.Contains(norm, "dispensary") {
		return []string{"core.organization", "core.patient", "core.billing", "pharmacy.basic"}
	}
	if strings.Contains(norm, "hospital") || strings.Contains(norm, "inpatient") || strings.Contains(norm, "his") {
		return []string{"core.organization", "core.patient", "core.billing", "clinical.basic", "laboratory.basic", "radiology.basic", "pharmacy.basic", "clinical.inpatient_wards"}
	}
	// Default to canonical outpatient clinic
	return []string{"core.organization", "core.patient", "core.billing", "clinical.basic"}
}

func resolveBranchTypeModules(facilityType string) []string {
	norm := strings.ToLower(facilityType)
	if strings.Contains(norm, "lab") || strings.Contains(norm, "pathology") {
		return []string{"dashboard", "reception", "laboratory", "billing"}
	}
	if strings.Contains(norm, "radiology") || strings.Contains(norm, "imaging") || strings.Contains(norm, "pacs") {
		return []string{"dashboard", "reception", "radiology", "billing"}
	}
	if strings.Contains(norm, "pharmacy") || strings.Contains(norm, "dispensary") {
		return []string{"dashboard", "reception", "pharmacy", "billing"}
	}
	if strings.Contains(norm, "hospital") || strings.Contains(norm, "inpatient") || strings.Contains(norm, "his") {
		return []string{"dashboard", "reception", "care_desk", "clinical", "laboratory", "pharmacy", "radiology", "hospital", "billing"}
	}
	// Canonical Outpatient Clinic: excludes LIS, RIS/PACS, and Inpatient Hospital
	if strings.Contains(norm, "clinic") || strings.Contains(norm, "outpatient") || strings.Contains(norm, "emr") || strings.Contains(norm, "opd") || strings.Contains(norm, "specialty") {
		return []string{"dashboard", "reception", "care_desk", "clinical", "billing"}
	}
	return []string{"dashboard", "reception", "care_desk", "clinical", "billing"}
}

func isPatientHost(host string) bool {
	return host == "patient.localhost" || host == "patient.curexal.space" || host == "patient.curexal.internal" || strings.HasPrefix(host, "patient.")
}

func isPlatformHost(host string) bool {
	return host == "app.localhost" || host == "app.curexal.space" || host == "app.curexal.internal"
}

func isOrgDomain(host string) bool {
	if isPlatformHost(host) || isPatientHost(host) {
		return false
	}
	return strings.HasSuffix(host, ".localhost") || strings.HasSuffix(host, ".curexal.space")
}
