package application

import (
	"context"
	"fmt"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/facility_config/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
)

type FacilityConfigApplicationService struct {
	server    *server.Server
	repo      domain.FacilityConfigRepository
	auditRepo auditDomain.AuditRepository
}

func NewFacilityConfigApplicationService(server *server.Server, repo domain.FacilityConfigRepository) *FacilityConfigApplicationService {
	return &FacilityConfigApplicationService{
		server: server,
		repo:   repo,
	}
}

func (s *FacilityConfigApplicationService) SetAuditRepository(auditRepo auditDomain.AuditRepository) {
	s.auditRepo = auditRepo
}

func (s *FacilityConfigApplicationService) IsPlatformAdmin(principal *middleware.AuthenticatedPrincipal) bool {
	if principal == nil {
		return false
	}
	if principal.Platform.IsPlatformAdmin || principal.Platform.IsSuperAdmin || principal.Platform.IsPlatformStaff {
		return true
	}
	if principal.Role == "super_admin" || principal.Role == "platform_admin" || principal.Role == "platform_staff" {
		return true
	}
	return false
}

func (s *FacilityConfigApplicationService) GetActiveFacilityTypes(ctx context.Context) ([]domain.FacilityTypeDTO, error) {
	return s.repo.GetActiveFacilityTypes(ctx)
}

func (s *FacilityConfigApplicationService) ListFacilityTypeEntities(ctx context.Context) ([]domain.FacilityTypeEntity, error) {
	return s.repo.ListFacilityTypeEntities(ctx)
}

func (s *FacilityConfigApplicationService) CreateFacilityType(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	ft *domain.FacilityTypeEntity,
) (*domain.FacilityTypeEntity, error) {
	if !s.IsPlatformAdmin(principal) {
		return nil, fmt.Errorf("unauthorized: operation requires platform admin privileges")
	}
	if err := ft.Validate(); err != nil {
		return nil, err
	}

	actorUUID, err := uuid.Parse(principal.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", err)
	}

	created, err := s.repo.CreateFacilityType(ctx, ft, actorUUID)
	if err != nil {
		return nil, err
	}

	if s.auditRepo != nil {
		action := "FACILITY_TYPE_CREATED"
		resType := "platform.facility_types"
		resID := created.ID.String()
		eventCat := "ADMINISTRATIVE"
		severity := "HIGH"
		status := "SUCCESS"

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    true,
			ActorID:       &principal.UserID,
			ActorName:     &principal.Identity.FullName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	return created, nil
}

func (s *FacilityConfigApplicationService) UpdateFacilityType(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	ft *domain.FacilityTypeEntity,
) (*domain.FacilityTypeEntity, error) {
	if !s.IsPlatformAdmin(principal) {
		return nil, fmt.Errorf("unauthorized: operation requires platform admin privileges")
	}
	if err := ft.Validate(); err != nil {
		return nil, err
	}

	actorUUID, err := uuid.Parse(principal.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", err)
	}

	updated, err := s.repo.UpdateFacilityType(ctx, ft, actorUUID)
	if err != nil {
		return nil, err
	}

	if s.auditRepo != nil {
		action := "FACILITY_TYPE_UPDATED"
		resType := "platform.facility_types"
		resID := updated.ID.String()
		eventCat := "ADMINISTRATIVE"
		severity := "HIGH"
		status := "SUCCESS"

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    true,
			ActorID:       &principal.UserID,
			ActorName:     &principal.Identity.FullName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	return updated, nil
}

func (s *FacilityConfigApplicationService) GetRegistrationForm(ctx context.Context, facilityTypeID string) (*domain.RegistrationFormDTO, error) {
	return s.repo.GetRegistrationForm(ctx, facilityTypeID)
}

func (s *FacilityConfigApplicationService) GetNavigationMenu(ctx context.Context, facilityTypeID string) (*domain.NavigationMenuDTO, error) {
	return s.repo.GetNavigationMenu(ctx, facilityTypeID)
}

func (s *FacilityConfigApplicationService) GetSetupSteps(ctx context.Context, facilityTypeID string) ([]domain.SetupStepDTO, error) {
	return s.repo.GetSetupSteps(ctx, facilityTypeID)
}

func (s *FacilityConfigApplicationService) GetDashboard(ctx context.Context, facilityTypeID string) (*domain.DashboardDTO, error) {
	return s.repo.GetDashboard(ctx, facilityTypeID)
}

func (s *FacilityConfigApplicationService) GetTenantOverrides(ctx context.Context, tenantID uuid.UUID, branchID *uuid.UUID) (map[string]any, error) {
	return s.repo.GetTenantOverrides(ctx, tenantID, branchID)
}
