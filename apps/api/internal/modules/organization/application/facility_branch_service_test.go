package application_test

import (
	"context"
	"testing"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/organization/application"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBranchRepo struct {
	mock.Mock
}

func (m *MockBranchRepo) ListBranches(ctx context.Context, orgID uuid.UUID) ([]domain.FacilityBranch, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.FacilityBranch), args.Error(1)
}

func (m *MockBranchRepo) GetBranchByID(ctx context.Context, orgID, branchID uuid.UUID) (*domain.FacilityBranch, error) {
	args := m.Called(ctx, orgID, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FacilityBranch), args.Error(1)
}

func (m *MockBranchRepo) GetBranchByCode(ctx context.Context, orgID uuid.UUID, code string) (*domain.FacilityBranch, error) {
	args := m.Called(ctx, orgID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FacilityBranch), args.Error(1)
}

func (m *MockBranchRepo) CreateBranch(ctx context.Context, branch *domain.FacilityBranch, actorID uuid.UUID) (*domain.FacilityBranch, error) {
	args := m.Called(ctx, branch, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FacilityBranch), args.Error(1)
}

func (m *MockBranchRepo) UpdateBranch(ctx context.Context, branch *domain.FacilityBranch, actorID uuid.UUID) (*domain.FacilityBranch, error) {
	args := m.Called(ctx, branch, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FacilityBranch), args.Error(1)
}

func (m *MockBranchRepo) DeactivateBranch(ctx context.Context, orgID, branchID uuid.UUID, actorID uuid.UUID) error {
	args := m.Called(ctx, orgID, branchID, actorID)
	return args.Error(0)
}

func (m *MockBranchRepo) CountActiveBranches(ctx context.Context, orgID uuid.UUID) (int, error) {
	args := m.Called(ctx, orgID)
	return args.Int(0), args.Error(1)
}

func (m *MockBranchRepo) HasActiveHeadquarters(ctx context.Context, orgID uuid.UUID) (bool, error) {
	args := m.Called(ctx, orgID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBranchRepo) CheckFacilityTypeActive(ctx context.Context, facilityTypeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, facilityTypeID)
	return args.Bool(0), args.Error(1)
}

func TestFacilityBranchService_CreateBranch_Success(t *testing.T) {
	mockBranchRepo := new(MockBranchRepo)
	mockOrgRepo := new(MockOrgRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewFacilityBranchService(mockBranchRepo, mockOrgRepo, mockAuditRepo)

	orgID := uuid.New()
	actorID := uuid.New()
	facilityTypeID := uuid.New()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	payload := &domain.CreateFacilityBranchPayload{
		FacilityTypeID: facilityTypeID,
		Code:           "main-clinic",
		Name:           "Main Outpatient Branch",
		IsHeadquarters: true,
	}

	// 1. Facility type active check -> true
	mockBranchRepo.On("CheckFacilityTypeActive", mock.Anything, facilityTypeID).Return(true, nil)

	// 2. Org plan lookup -> smart (max 3 branches)
	mockOrgRepo.On("GetByID", mock.Anything, orgID).Return(&domain.Organization{
		ID:   orgID,
		Plan: "smart",
	}, nil)
	mockBranchRepo.On("CountActiveBranches", mock.Anything, orgID).Return(0, nil)

	// 3. Single HQ check -> false
	mockBranchRepo.On("HasActiveHeadquarters", mock.Anything, orgID).Return(false, nil)

	createdBranch := &domain.FacilityBranch{
		ID:             uuid.New(),
		OrganizationID: orgID,
		FacilityTypeID: facilityTypeID,
		Code:           "main-clinic",
		Name:           "Main Outpatient Branch",
		IsHeadquarters: true,
		Status:         "ACTIVE",
		Version:        1,
	}

	mockBranchRepo.On("CreateBranch", mock.Anything, mock.Anything, actorID).Return(createdBranch, nil)
	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "BRANCH_CREATED"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.CreateBranch(context.Background(), principal, payload)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "main-clinic", res.Code)
	assert.True(t, res.IsHeadquarters)

	mockBranchRepo.AssertExpectations(t)
	mockOrgRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestFacilityBranchService_CreateBranch_InactiveFacilityType_Fails(t *testing.T) {
	mockBranchRepo := new(MockBranchRepo)
	svc := application.NewFacilityBranchService(mockBranchRepo, nil, nil)

	orgID := uuid.New()
	facilityTypeID := uuid.New()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	payload := &domain.CreateFacilityBranchPayload{
		FacilityTypeID: facilityTypeID,
		Code:           "deactivated-type-branch",
		Name:           "Invalid Branch",
	}

	// Facility type is inactive
	mockBranchRepo.On("CheckFacilityTypeActive", mock.Anything, facilityTypeID).Return(false, nil)

	res, err := svc.CreateBranch(context.Background(), principal, payload)

	assert.ErrorIs(t, err, domain.ErrInactiveFacilityType)
	assert.Nil(t, res)

	mockBranchRepo.AssertNotCalled(t, "CreateBranch", mock.Anything, mock.Anything, mock.Anything)
}

func TestFacilityBranchService_CreateBranch_SingleHeadquartersConflict(t *testing.T) {
	mockBranchRepo := new(MockBranchRepo)
	mockOrgRepo := new(MockOrgRepo)
	svc := application.NewFacilityBranchService(mockBranchRepo, mockOrgRepo, nil)

	orgID := uuid.New()
	facilityTypeID := uuid.New()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	payload := &domain.CreateFacilityBranchPayload{
		FacilityTypeID: facilityTypeID,
		Code:           "second-hq",
		Name:           "Second Headquarters Attempt",
		IsHeadquarters: true,
	}

	mockBranchRepo.On("CheckFacilityTypeActive", mock.Anything, facilityTypeID).Return(true, nil)
	mockOrgRepo.On("GetByID", mock.Anything, orgID).Return(&domain.Organization{ID: orgID, Plan: "pro"}, nil)
	mockBranchRepo.On("CountActiveBranches", mock.Anything, orgID).Return(1, nil)

	// HQ already exists
	mockBranchRepo.On("HasActiveHeadquarters", mock.Anything, orgID).Return(true, nil)

	res, err := svc.CreateBranch(context.Background(), principal, payload)

	assert.ErrorIs(t, err, domain.ErrHeadquartersConflict)
	assert.Nil(t, res)

	mockBranchRepo.AssertNotCalled(t, "CreateBranch", mock.Anything, mock.Anything, mock.Anything)
}

func TestFacilityBranchService_CreateBranch_MaxBranchesExceeded(t *testing.T) {
	mockBranchRepo := new(MockBranchRepo)
	mockOrgRepo := new(MockOrgRepo)
	svc := application.NewFacilityBranchService(mockBranchRepo, mockOrgRepo, nil)

	orgID := uuid.New()
	facilityTypeID := uuid.New()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	payload := &domain.CreateFacilityBranchPayload{
		FacilityTypeID: facilityTypeID,
		Code:           "exceeded-branch",
		Name:           "Branch 4 on Smart Plan",
	}

	mockBranchRepo.On("CheckFacilityTypeActive", mock.Anything, facilityTypeID).Return(true, nil)
	// Smart plan allows max 3 branches
	mockOrgRepo.On("GetByID", mock.Anything, orgID).Return(&domain.Organization{ID: orgID, Plan: "smart"}, nil)
	mockBranchRepo.On("CountActiveBranches", mock.Anything, orgID).Return(3, nil)

	res, err := svc.CreateBranch(context.Background(), principal, payload)

	assert.ErrorIs(t, err, domain.ErrMaxBranchesExceeded)
	assert.Nil(t, res)

	mockBranchRepo.AssertNotCalled(t, "CreateBranch", mock.Anything, mock.Anything, mock.Anything)
}
