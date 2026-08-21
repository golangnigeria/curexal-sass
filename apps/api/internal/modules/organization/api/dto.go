package api

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
)

type TenantSettingsPayload struct {
	OrgType           *string           `json:"orgType,omitempty"`
	LabHeaderTitle    *string           `json:"labHeaderTitle,omitempty"`
	LabHeaderSubtitle *string           `json:"labHeaderSubtitle,omitempty"`
	PrimaryColor      *string           `json:"primaryColor,omitempty"`
	SecondaryColor    *string           `json:"secondaryColor,omitempty"`
	FontFamily        *string           `json:"fontFamily,omitempty"`
	ClinicalRanges    map[string]string `json:"clinicalRanges,omitempty"`
}

type TenantBranding struct {
	PrimaryColor      string            `json:"primaryColor"`
	SecondaryColor    string            `json:"secondaryColor"`
	FontFamily        string            `json:"fontFamily"`
	OrgType           string            `json:"orgType,omitempty"`
	LabHeaderTitle    string            `json:"labHeaderTitle,omitempty"`
	LabHeaderSubtitle string            `json:"labHeaderSubtitle,omitempty"`
	ClinicalRanges    map[string]string `json:"clinicalRanges,omitempty"`
}

type TenantResponse struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Slug           string         `json:"slug"`
	LogoURL        *string        `json:"logoUrl,omitempty"`
	Branding       TenantBranding `json:"branding"`
	OrganizationID string         `json:"organizationId,omitempty"`
	Currency       string         `json:"currency"`
	EnabledModules []string       `json:"enabledModules"`
}

type OwnerPayload struct {
	Email string  `json:"email" validate:"required,email"`
	Name  *string `json:"name,omitempty" validate:"omitempty,max=100"`
}

type CreateOrganizationPayload struct {
	Name               string        `json:"name" validate:"required,min=1,max=100"`
	Slug               string        `json:"slug" validate:"required,min=1,max=100"`
	Plan               *string       `json:"plan,omitempty" validate:"omitempty,oneof=smart optimize pro enterprise"`
	Address            *string       `json:"address,omitempty"`
	City               *string       `json:"city,omitempty"`
	State              *string       `json:"state,omitempty"`
	LGA                *string       `json:"lga,omitempty"`
	Country            *string       `json:"country,omitempty"`
	Phone              *string       `json:"phone,omitempty"`
	Email              *string       `json:"email,omitempty" validate:"omitempty,email"`
	RegistrationNumber *string       `json:"registrationNumber,omitempty"`
	LicenseNumber      *string       `json:"licenseNumber,omitempty"`
	TaxID              *string       `json:"taxId,omitempty"`
	Owner              *OwnerPayload `json:"owner,omitempty"`
	OwnerEmail         *string       `json:"ownerEmail,omitempty" validate:"omitempty,email"`
	OwnerName          *string       `json:"ownerName,omitempty" validate:"omitempty,max=100"`
}

func (p *CreateOrganizationPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type OrganizationInvitationInfo struct {
	Sent  bool   `json:"sent"`
	Email string `json:"email"`
}

type CreateOrganizationResponse struct {
	Message      string                     `json:"message"`
	Organization *domain.Organization       `json:"organization"`
	Invitation   OrganizationInvitationInfo `json:"invitation"`
}

type TransferOrganizationOwnershipPayload struct {
	NewOwnerEmail string  `json:"newOwnerEmail" validate:"required,email"`
	NewOwnerName  string  `json:"newOwnerName" validate:"required,min=1,max=100"`
	Notes         *string `json:"notes,omitempty"`
}

func (p *TransferOrganizationOwnershipPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type UpdateOrganizationPayload struct {
	ID           string         `param:"id" validate:"required,uuid"`
	Name         *string        `json:"name" validate:"omitempty,min=1,max=100"`
	Slug         *string        `json:"slug" validate:"omitempty,min=1,max=100"`
	Plan         *string        `json:"plan" validate:"omitempty,oneof=smart optimize pro enterprise"`
	CustomDomain *string        `json:"customDomain" validate:"omitempty,max=253"`
	Settings     map[string]any `json:"settings" validate:"omitempty"`
}

func (p *UpdateOrganizationPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type CreateTenantPayload struct {
	Name           string   `json:"name" validate:"required,min=1,max=100"`
	Slug           string   `json:"slug" validate:"required,min=1,max=100"`
	LogoURL        *string  `json:"logoUrl,omitempty"`
	OrganizationID string   `json:"organizationId" param:"orgId" validate:"required,uuid"`
	Location       string   `json:"location,omitempty"`
	Phone          string   `json:"phone,omitempty"`
	Address        string   `json:"address,omitempty"`
	Currency       *string  `json:"currency,omitempty" validate:"omitempty,len=3"`
	Modules        []string `json:"modules,omitempty"`
}

func (p *CreateTenantPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}





type CreateDemoRequestPayload struct {
	LaboratoryName string  `json:"laboratoryName" validate:"required,min=1"`
	ContactName    string  `json:"contactName"    validate:"required,min=1"`
	Email          string  `json:"email"          validate:"required,email"`
	Phone          *string `json:"phone"          validate:"omitempty"`
	DailyVolume    *string `json:"dailyVolume"    validate:"omitempty"`
	Notes          *string `json:"notes"          validate:"omitempty"`
}

func (p *CreateDemoRequestPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type UpdateDemoRequestPayload struct {
	Status *string `json:"status" validate:"omitempty,oneof=pending scheduled completed"`
	Notes  *string `json:"notes"  validate:"omitempty"`
}

func (p *UpdateDemoRequestPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type UpdateOrganizationSettingsPayload struct {
	LogoURL       *string `json:"logoUrl"       validate:"omitempty"`
	ThemeBranding *string `json:"themeBranding" validate:"omitempty,json"`
	CustomDomain  *string `json:"customDomain"  validate:"omitempty"`
	SupportEmail  *string `json:"supportEmail"  validate:"omitempty,email"`
	SupportPhone  *string `json:"supportPhone"  validate:"omitempty,max=20"`
	CACNumber     *string `json:"cacNumber"     validate:"omitempty,max=100"`
	TINNumber     *string `json:"tinNumber"     validate:"omitempty,max=100"`
	TaxNumber     *string `json:"taxNumber"     validate:"omitempty,max=100"`
	BusinessType  *string `json:"businessType"  validate:"omitempty,max=100"`
	Address       *string `json:"address"       validate:"omitempty"`
	Timezone      *string `json:"timezone"      validate:"omitempty,max=100"`
	Currency      *string `json:"currency"      validate:"omitempty,max=10"`
	Language      *string `json:"language"      validate:"omitempty,max=10"`
}

func (p *UpdateOrganizationSettingsPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func MapTenantToResponse(t *domain.Tenant) *TenantResponse {
	var branding TenantBranding
	if t.Settings != "" {
		_ = json.Unmarshal([]byte(t.Settings), &branding)
	}
	if branding.PrimaryColor == "" {
		branding.PrimaryColor = "#0284c7"
	}
	if branding.SecondaryColor == "" {
		branding.SecondaryColor = "#0f172a"
	}
	if branding.FontFamily == "" {
		branding.FontFamily = "Outfit"
	}

	var logoURL *string
	if t.LogoURL != nil && *t.LogoURL != "" {
		logoURL = t.LogoURL
	}

	currency := t.Currency
	if currency == "" {
		currency = "NGN"
	}

	enabledModules := t.EnabledModules
	if enabledModules == nil {
		enabledModules = []string{}
	}

	return &TenantResponse{
		ID:             t.ID.String(),
		Name:           t.Name,
		Slug:           t.Slug,
		LogoURL:        logoURL,
		Branding:       branding,
		OrganizationID: t.OrganizationID,
		Currency:       currency,
		EnabledModules: enabledModules,
	}
}

type ReviewDocumentPayload struct {
	Status          string  `json:"status" validate:"required,oneof=approved rejected"`
	RejectionReason *string `json:"rejectionReason,omitempty" validate:"required_if=Status rejected"`
}

func (p *ReviewDocumentPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type RejectOrganizationPayload struct {
	Reason string `json:"reason" validate:"required,min=1,max=500"`
}

func (p *RejectOrganizationPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
