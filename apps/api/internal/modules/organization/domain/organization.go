package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrOrganizationNotFound          = errors.New("organization not found")
	ErrInvalidSetupStateTransition   = errors.New("invalid setup wizard state transition")
	ErrOptimisticLockingConflict     = errors.New("optimistic concurrency conflict: organization record was modified concurrently")
	ErrUnauthorizedTenantAccess      = errors.New("unauthorized: caller does not have permission for active tenant organization")
	ErrInvalidOrganizationProfile    = errors.New("invalid organization profile data")
)

type SetupState string

const (
	SetupStatePendingRegistration SetupState = "PENDING_REGISTRATION"
	SetupStateProfileCompleted    SetupState = "PROFILE_COMPLETED"
	SetupStateDocumentsSubmitted SetupState = "DOCUMENTS_SUBMITTED"
	SetupStateUnderReview        SetupState = "UNDER_REVIEW"
	SetupStateVerified           SetupState = "VERIFIED"
	SetupStateRejected           SetupState = "REJECTED"
)

func IsValidSetupTransition(from, to SetupState) bool {
	if from == to {
		return true // Idempotent same-state transition
	}
	switch from {
	case SetupStatePendingRegistration:
		return to == SetupStateProfileCompleted
	case SetupStateProfileCompleted:
		return to == SetupStateDocumentsSubmitted
	case SetupStateDocumentsSubmitted:
		return to == SetupStateUnderReview
	case SetupStateUnderReview:
		return to == SetupStateVerified || to == SetupStateRejected
	case SetupStateRejected:
		return to == SetupStateDocumentsSubmitted || to == SetupStateUnderReview
	default:
		return false
	}
}

type Organization struct {
	ID                 uuid.UUID      `json:"id"                 db:"id"`
	Name               string         `json:"name"               db:"name"`
	Slug               string         `json:"slug"               db:"slug"`
	Status             string         `json:"status"             db:"status"`
	Plan               string         `json:"plan"               db:"plan"`
	LogoURL            *string        `json:"logoUrl"            db:"logo_url"`
	CustomDomain       *string        `json:"customDomain"       db:"custom_domain"`
	RegistrationNumber *string        `json:"registrationNumber" db:"registration_number"`
	LicenseNumber      *string        `json:"licenseNumber"      db:"license_number"`
	TaxID              *string        `json:"taxId"              db:"tax_id"`
	Email              *string        `json:"email"              db:"email"`
	Phone              *string        `json:"phone"              db:"phone"`
	Address            *string        `json:"address"            db:"address"`
	City               *string        `json:"city"               db:"city"`
	State              *string        `json:"state"              db:"state"`
	LGA                *string        `json:"lga"                db:"lga"`
	Country            string         `json:"country"            db:"country"`
	SetupState         SetupState     `json:"setupState"         db:"setup_state"`
	SetupStep          int            `json:"setupStep"          db:"setup_step"`
	CompletedAt        *time.Time     `json:"completedAt"        db:"completed_at"`
	Settings           map[string]any `json:"settings"           db:"settings"`
	Version            int            `json:"version"            db:"version"`
	CreatedAt          time.Time      `json:"createdAt"          db:"created_at"`
	UpdatedAt          time.Time      `json:"updatedAt"          db:"updated_at"`
	UpdatedBy          *uuid.UUID     `json:"updatedBy"          db:"updated_by"`
}

type UpdateOrganizationProfilePayload struct {
	Name               *string `json:"name"`
	RegistrationNumber *string `json:"registrationNumber"`
	LicenseNumber      *string `json:"licenseNumber"`
	TaxID              *string `json:"taxId"`
	Email              *string `json:"email"`
	Phone              *string `json:"phone"`
	Address            *string `json:"address"`
	City               *string `json:"city"`
	State              *string `json:"state"`
	LGA                *string `json:"lga"`
	Country            *string `json:"country"`
	LogoURL            *string `json:"logoUrl"`
	CustomDomain       *string `json:"customDomain"`
	Version            int     `json:"version"`
}

func (o *Organization) Validate() error {
	if o.Name == "" || o.Slug == "" {
		return ErrInvalidOrganizationProfile
	}
	return nil
}

type OrganizationSettings struct {
	ID             string  `json:"id"             db:"id"`
	OrganizationID string  `json:"organizationId" db:"organization_id"`
	LogoURL        *string `json:"logoUrl"        db:"logo_url"`
	ThemeBranding  string  `json:"themeBranding"  db:"theme_branding"`
	CustomDomain   *string `json:"customDomain"   db:"custom_domain"`
	SupportEmail   *string `json:"supportEmail"   db:"support_email"`
	SupportPhone   *string `json:"supportPhone"   db:"support_phone"`
	CACNumber      *string `json:"cacNumber"      db:"cac_number"`
	TINNumber      *string `json:"tinNumber"      db:"tin_number"`
	TaxNumber      *string `json:"taxNumber"      db:"tax_number"`
	BusinessType   *string `json:"businessType"   db:"business_type"`
	Address        *string `json:"address"        db:"address"`
	Timezone       *string `json:"timezone"       db:"timezone"`
	Currency       *string `json:"currency"       db:"currency"`
	Language       *string `json:"language"       db:"language"`
}

type OrganizationNavigationItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Path  string `json:"path"`
	Order int    `json:"order"`
}

type OrganizationDomainProvider struct{}

func NewOrganizationDomainProvider() *OrganizationDomainProvider {
	return &OrganizationDomainProvider{}
}

func (p *OrganizationDomainProvider) GetOrganizationNavigation() []OrganizationNavigationItem {
	return []OrganizationNavigationItem{
		{ID: "nav_org_dashboard", Title: "Executive HQ Dashboard", Icon: "LayoutDashboard", Path: "/organization/dashboard", Order: 1},
		{ID: "nav_org_branches", Title: "Branch Facilities", Icon: "Building2", Path: "/organization/branches", Order: 2},
		{ID: "nav_org_members", Title: "Staff & Members", Icon: "Users", Path: "/organization/members", Order: 3},
		{ID: "nav_org_roles", Title: "Roles & Permissions", Icon: "Shield", Path: "/organization/roles", Order: 4},
		{ID: "nav_org_catalogs", Title: "Catalogs & Pricing", Icon: "BookOpen", Path: "/organization/catalogs", Order: 5},
		{ID: "nav_org_billing", Title: "Corporate Subscription", Icon: "CreditCard", Path: "/organization/billing", Order: 6},
		{ID: "nav_org_branding", Title: "Branding & Customization", Icon: "Palette", Path: "/organization/branding", Order: 7},
		{ID: "nav_org_notifications", Title: "Notification Settings", Icon: "Bell", Path: "/organization/notifications", Order: 8},
		{ID: "nav_org_integrations", Title: "APIs & Webhooks", Icon: "Cpu", Path: "/organization/integrations", Order: 9},
		{ID: "nav_org_audit", Title: "Corporate Audit Ledger", Icon: "History", Path: "/organization/audit", Order: 10},
		{ID: "nav_org_settings", Title: "Organization Settings", Icon: "Settings", Path: "/organization/settings", Order: 11},
	}
}
