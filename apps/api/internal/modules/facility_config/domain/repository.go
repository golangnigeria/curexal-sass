package domain

import (
	"context"

	"github.com/google/uuid"
)

type FacilityConfigRepository interface {
	GetActiveFacilityTypes(ctx context.Context) ([]FacilityTypeDTO, error)
	ListFacilityTypeEntities(ctx context.Context) ([]FacilityTypeEntity, error)
	CreateFacilityType(ctx context.Context, ft *FacilityTypeEntity, updatedBy uuid.UUID) (*FacilityTypeEntity, error)
	UpdateFacilityType(ctx context.Context, ft *FacilityTypeEntity, updatedBy uuid.UUID) (*FacilityTypeEntity, error)
	GetRegistrationForm(ctx context.Context, facilityTypeID string) (*RegistrationFormDTO, error)
	GetNavigationMenu(ctx context.Context, facilityTypeID string) (*NavigationMenuDTO, error)
	GetSetupSteps(ctx context.Context, facilityTypeID string) ([]SetupStepDTO, error)
	GetDashboard(ctx context.Context, facilityTypeID string) (*DashboardDTO, error)
	GetTenantOverrides(ctx context.Context, tenantID uuid.UUID, branchID *uuid.UUID) (map[string]any, error)
}
