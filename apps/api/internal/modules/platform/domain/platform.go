package domain

// PlatformNavigationItem defines navigation items within the Platform Bounded Context.
type PlatformNavigationItem struct {
	ID       string                   `json:"id"`
	Title    string                   `json:"title"`
	Icon     string                   `json:"icon"`
	Path     string                   `json:"path"`
	Order    int                      `json:"order,omitempty"`
	Children []PlatformNavigationItem `json:"children,omitempty"`
}

// PlatformDomainProvider encapsulates core domain logic and rules for the Platform Bounded Context.
type PlatformDomainProvider struct{}

func NewPlatformDomainProvider() *PlatformDomainProvider {
	return &PlatformDomainProvider{}
}

// CanonicalPlatformRoles defines the exhaustive set of valid platform administrator/staff roles.
var CanonicalPlatformRoles = map[string]bool{
	"super_admin":              true,
	"platform_admin":           true,
	"platform_staff":           true,
	"super_support_agent":      true,
	"super_sales_staff":        true,
	"super_compliance_officer": true,
}

// IsPlatformContext returns true if and only if the principal possesses an explicit platform staff role.
func (p *PlatformDomainProvider) IsPlatformContext(isStaff bool, role string, email string) bool {
	if isStaff {
		return true
	}
	return CanonicalPlatformRoles[role]
}

// GetPlatformNavigation returns the domain-driven platform navigation items.
func (p *PlatformDomainProvider) GetPlatformNavigation() []PlatformNavigationItem {
	return []PlatformNavigationItem{
		{ID: "nav_plat_dashboard", Title: "Platform Dashboard", Icon: "LayoutDashboard", Path: "/platform/dashboard", Order: 1},
		{ID: "nav_plat_orgs", Title: "Organizations", Icon: "Building2", Path: "/platform/organizations", Order: 2},
		{ID: "nav_plat_users", Title: "User Directory", Icon: "Users", Path: "/platform/users", Order: 3},
		{ID: "nav_plat_marketplace", Title: "B2B Marketplace", Icon: "Store", Path: "/platform/marketplace", Order: 4},
		{ID: "nav_plat_pricing", Title: "Pricing & Billing", Icon: "CreditCard", Path: "/platform/pricing", Order: 5},
		{ID: "nav_plat_facility_types", Title: "Facility Types", Icon: "Layers", Path: "/platform/facility-types", Order: 6},
		{ID: "nav_plat_catalogs", Title: "Master Catalogs", Icon: "BookOpen", Path: "/platform/catalogs", Order: 7},
		{ID: "nav_plat_audit", Title: "Audit Trail", Icon: "History", Path: "/platform/audit", Order: 8},
		{ID: "nav_plat_diag", Title: "Diagnostics & Gate", Icon: "Cpu", Path: "/platform/diagnostics", Order: 9},
		{ID: "nav_plat_demo", Title: "Demo Requests", Icon: "Inbox", Path: "/platform/demo-requests", Order: 10},
		{ID: "nav_plat_settings", Title: "Console Settings", Icon: "Settings", Path: "/platform/settings", Order: 11},
	}
}
