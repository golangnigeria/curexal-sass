package api

import (
	"github.com/go-playground/validator/v10"
)

type UpdateSettingsSectionPayload struct {
	Section string         `param:"section" validate:"required,oneof=general financial inventory integrations notifications document_header patient lims consultation staff"`
	Config  map[string]any `json:"config" validate:"required"`
}

func (p *UpdateSettingsSectionPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
