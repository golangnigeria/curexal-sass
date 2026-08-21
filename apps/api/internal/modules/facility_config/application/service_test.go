package application_test

import (
	"context"
	"testing"

	"github.com/golangnigeria/curexal/internal/modules/facility_config/application"
	"github.com/golangnigeria/curexal/internal/modules/facility_config/domain"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFacilityConfigRepo struct {
	mock.Mock
}

func (m *MockFacilityConfigRepo) GetActiveFacilityTypes(ctx context.Context) ([]domain.FacilityTypeDTO, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.FacilityTypeDTO), args.Error(1)
}

func (m *MockFacilityConfigRepo) ListFacilityTypeEntities(ctx context.Context) ([]domain.FacilityTypeEntity, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.FacilityTypeEntity), args.Error(1)
}

func (m *MockFacilityConfigRepo) CreateFacilityType(ctx context.Context, ft *domain.FacilityTypeEntity, updatedBy uuid.UUID) (*domain.FacilityTypeEntity, error) {
	args := m.Called(ctx, ft, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FacilityTypeEntity), args.Error(1)
}

func (m *MockFacilityConfigRepo) UpdateFacilityType(ctx context.Context, ft *domain.FacilityTypeEntity, updatedBy uuid.UUID) (*domain.FacilityTypeEntity, error) {
	args := m.Called(ctx, ft, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FacilityTypeEntity), args.Error(1)
}

func (m *MockFacilityConfigRepo) GetRegistrationForm(ctx context.Context, facilityTypeID string) (*domain.RegistrationFormDTO, error) {
	return nil, nil
}
func (m *MockFacilityConfigRepo) GetNavigationMenu(ctx context.Context, facilityTypeID string) (*domain.NavigationMenuDTO, error) {
	return nil, nil
}
func (m *MockFacilityConfigRepo) GetSetupSteps(ctx context.Context, facilityTypeID string) ([]domain.SetupStepDTO, error) {
	return nil, nil
}
func (m *MockFacilityConfigRepo) GetDashboard(ctx context.Context, facilityTypeID string) (*domain.DashboardDTO, error) {
	return nil, nil
}
func (m *MockFacilityConfigRepo) GetTenantOverrides(ctx context.Context, tenantID uuid.UUID, branchID *uuid.UUID) (map[string]any, error) {
	return nil, nil
}

func TestFacilityConfigService_CreateFacilityType_PlatformAdmin(t *testing.T) {
	mockRepo := new(MockFacilityConfigRepo)
	service := application.NewFacilityConfigApplicationService(nil, mockRepo)

	adminID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: adminID.String(),
		Role:   "platform_admin",
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: true,
		},
	}

	ftInput := &domain.FacilityTypeEntity{
		Code:     "research_center",
		Name:     "Research Center",
		Category: "research",
	}

	expected := *ftInput
	expected.ID = uuid.New()
	expected.Version = 1

	mockRepo.On("CreateFacilityType", mock.Anything, ftInput, adminID).Return(&expected, nil)

	res, err := service.CreateFacilityType(context.Background(), principal, ftInput)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "research_center", res.Code)
	mockRepo.AssertExpectations(t)
}

func TestFacilityConfigService_CreateFacilityType_Unauthorized(t *testing.T) {
	mockRepo := new(MockFacilityConfigRepo)
	service := application.NewFacilityConfigApplicationService(nil, mockRepo)

	principal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Role:   "regular_user",
	}

	ftInput := &domain.FacilityTypeEntity{
		Code:     "research_center",
		Name:     "Research Center",
		Category: "research",
	}

	res, err := service.CreateFacilityType(context.Background(), principal, ftInput)

	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "CreateFacilityType", mock.Anything, mock.Anything, mock.Anything)
}
