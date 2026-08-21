package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/golangnigeria/curexal/internal/modules/audit/domain"
)

type CreateAuditLogPayload = domain.CreateAuditLogPayload

func ValidateCreatePayload(p *CreateAuditLogPayload) error {
	validate := validator.New()
	return validate.Struct(p)
}

type ListAuditLogsPayload struct {
	Limit          int     `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset         int     `query:"offset" validate:"omitempty,min=0"`
	Category       *string `query:"category"`
	Severity       *string `query:"severity"`
	Status         *string `query:"status"`
	ActorID        *string `query:"actorId"`
	Action         *string `query:"action"`
	ResourceType   *string `query:"resourceType"`
	ResourceID     *string `query:"resourceId"`
	OrganizationID *string `query:"organizationId"`
	StartDate      *string `query:"startDate"`
	EndDate        *string `query:"endDate"`
	Search         *string `query:"search"`
}

func (p *ListAuditLogsPayload) Validate() error {
	if p.Limit == 0 {
		p.Limit = 50
	}
	validate := validator.New()
	return validate.Struct(p)
}
