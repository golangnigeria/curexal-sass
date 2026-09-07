package auth

import (
	"github.com/go-playground/validator/v10"
)

type SignInPayload struct {
	Email            string  `json:"email" validate:"required,email"`
	Password         string  `json:"password" validate:"required,min=8"`
	OrganizationSlug *string `json:"organizationSlug,omitempty"`
	BranchID         *string `json:"branchId,omitempty"`
	BranchCode       *string `json:"branchCode,omitempty"`
}

func (p *SignInPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type SelectBranchPayload struct {
	SelectionToken string `json:"selectionToken" validate:"required"`
	BranchID       string `json:"branchId" validate:"required"`
}

func (p *SelectBranchPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type SwitchBranchPayload struct {
	BranchID string `json:"branchId" validate:"required"`
}

func (p *SwitchBranchPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type SignUpPayload struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	OrgName  string `json:"orgName"`
	OrgType  string `json:"orgType"`
}

func (p *SignUpPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type SignOutPayload struct{}

func (p *SignOutPayload) Validate() error {
	return nil
}

type VerifyOTPPayload struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

func (p *VerifyOTPPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type RequestPasswordPayload struct {
	Email string `json:"email" validate:"required,email"`
}

func (p *RequestPasswordPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type ForgotPasswordPayload struct {
	Email string `json:"email" validate:"required,email"`
}

func (p *ForgotPasswordPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

