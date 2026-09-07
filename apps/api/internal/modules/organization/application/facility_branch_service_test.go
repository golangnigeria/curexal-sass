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

func (m *MockBranchRepo) GetBranchBySlug(ctx context.Context, orgID uuid.UUID, slug string) (*domain.FacilityBranch, error) {
	args := m.Called(ctx, orgID, slug)
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

func (m *MockBranchRepo) SetHeadquarters(ctx context.Context, orgID, branchID, actorID uuid.UUID) error {
	args := m.Called(ctx, orgID, branchID, actorID)
	return args.Error(0)
}

func (m *MockBranchRepo) GetFacilityTypeByCode(ctx context.Context, code string) (*domain.FacilityType, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FacilityType), args.Error(1)
}

func (m *MockBranchRepo) ListFacilityTypes(ctx context.Context) ([]domain.FacilityType, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.FacilityType), args.Error(1)
}

func (m *MockBranchRepo) VerifyUserFacilityAccess(ctx context.Context, orgID, branchID uuid.UUID, userID string, isOrgAdmin bool) (bool, error) {
	args := m.Called(ctx, orgID, branchID, userID, isOrgAdmin)
	return args.Bool(0), args.Error(1)
}

func (m *MockBranchRepo) IsUserAssignedToBranch(ctx context.Context, userID, orgID, branchID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, orgID, branchID)
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

	mockOrgRepo.On("VerifyMembership", mock.Anything, orgID, actorID.String()).Return(true, nil)

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
	mockOrgRepo := new(MockOrgRepo)
	svc := application.NewFacilityBranchService(mockBranchRepo, mockOrgRepo, nil)

	orgID := uuid.New()
	facilityTypeID := uuid.New()
	userID := uuid.New().String()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: userID,
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	mockOrgRepo.On("VerifyMembership", mock.Anything, orgID, userID).Return(true, nil)

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
	userID := uuid.New().String()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: userID,
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	mockOrgRepo.On("VerifyMembership", mock.Anything, orgID, userID).Return(true, nil)

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
	userID := uuid.New().String()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: userID,
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	mockOrgRepo.On("VerifyMembership", mock.Anything, orgID, userID).Return(true, nil)

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

func TestTenantIsolation_CrossOrgIDOR_Rejected(t *testing.T) {
	mockBranchRepo := new(MockBranchRepo)
	mockOrgRepo := new(MockOrgRepo)
	svc := application.NewFacilityBranchService(mockBranchRepo, mockOrgRepo, nil)

	victimOrgID := uuid.New()
	attackerUserID := uuid.New().String()

	// Attacker provides victim's organization ID in request header
	principal := &middleware.AuthenticatedPrincipal{
		UserID: attackerUserID,
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: victimOrgID.String(),
		},
		Platform: platformAuth.PlatformVector{
			IsPlatformStaff: false,
		},
	}

	// Membership check fails for attacker in victim org
	mockOrgRepo.On("VerifyMembership", mock.Anything, victimOrgID, attackerUserID).Return(false, nil)
	mockOrgRepo.On("List", mock.Anything, attackerUserID, false).Return([]domain.Organization{}, nil)

	res, err := svc.ListBranches(context.Background(), principal)

	assert.ErrorIs(t, err, domain.ErrUnauthorizedTenantAccess)
	assert.Nil(t, res)
	mockBranchRepo.AssertNotCalled(t, "ListBranches", mock.Anything, mock.Anything)
}

func TestTenantIsolation_PlatformAdmin_Allowed(t *testing.T) {
	mockBranchRepo := new(MockBranchRepo)
	mockOrgRepo := new(MockOrgRepo)
	svc := application.NewFacilityBranchService(mockBranchRepo, mockOrgRepo, nil)

	targetOrgID := uuid.New()
	adminUserID := uuid.New().String()

	// Platform Super Admin accesses target organization
	principal := &middleware.AuthenticatedPrincipal{
		UserID: adminUserID,
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: targetOrgID.String(),
		},
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: true,
			IsSuperAdmin:    true,
		},
	}

	mockBranchRepo.On("ListBranches", mock.Anything, targetOrgID).Return([]domain.FacilityBranch{
		{ID: uuid.New(), Name: "Main Branch", Code: "MAIN"},
	}, nil)

	res, err := svc.ListBranches(context.Background(), principal)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "Main Branch", res[0].Name)
}

func TestFacilityIsolation_UserBranchAccess_Enforced(t *testing.T) {
	mockBranchRepo := new(MockBranchRepo)
	mockOrgRepo := new(MockOrgRepo)
	svc := application.NewFacilityBranchService(mockBranchRepo, mockOrgRepo, nil)

	orgID := uuid.New()
	branchID := uuid.New()
	userID := uuid.New().String()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: userID,
		Role:   "receptionist",
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
			OrganizationRole:     "receptionist",
		},
	}

	mockOrgRepo.On("VerifyMembership", mock.Anything, orgID, userID).Return(true, nil)
	// User is not assigned to this branch
	mockBranchRepo.On("VerifyUserFacilityAccess", mock.Anything, orgID, branchID, userID, false).Return(false, nil)

	res, err := svc.GetBranchByID(context.Background(), principal, branchID)

	assert.ErrorIs(t, err, domain.ErrUnauthorizedTenantAccess)
	assert.Nil(t, res)
}

func TestFacilityBranchService_SetHeadquarters_Success(t *testing.T) {
	mockBranchRepo := new(MockBranchRepo)
	mockOrgRepo := new(MockOrgRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewFacilityBranchService(mockBranchRepo, mockOrgRepo, mockAuditRepo)

	orgID := uuid.New()
	branchID := uuid.New()
	adminID := uuid.New()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: adminID.String(),
		Role:   "owner",
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
			OrganizationRole:     "owner",
		},
	}

	mockOrgRepo.On("VerifyMembership", mock.Anything, orgID, adminID.String()).Return(true, nil)
	mockBranchRepo.On("GetBranchByID", mock.Anything, orgID, branchID).Return(&domain.FacilityBranch{
		ID:             branchID,
		OrganizationID: orgID,
		Name:           "Ikeja Branch",
		Status:         "ACTIVE",
		IsHeadquarters: false,
	}, nil).Once()

	mockBranchRepo.On("SetHeadquarters", mock.Anything, orgID, branchID, adminID).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "HEADQUARTERS_CHANGED"
	})).Return(&auditDomain.AuditLog{}, nil)

	mockBranchRepo.On("GetBranchByID", mock.Anything, orgID, branchID).Return(&domain.FacilityBranch{
		ID:             branchID,
		OrganizationID: orgID,
		Name:           "Ikeja Branch",
		Status:         "ACTIVE",
		IsHeadquarters: true,
	}, nil).Once()

	res, err := svc.SetHeadquarters(context.Background(), principal, branchID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.IsHeadquarters)
}

func TestFacilityBranchService_UpdateBranch_OptimisticConcurrencyConflict(t *testing.T) {
	mockBranchRepo := new(MockBranchRepo)
	mockOrgRepo := new(MockOrgRepo)
	svc := application.NewFacilityBranchService(mockBranchRepo, mockOrgRepo, nil)

	orgID := uuid.New()
	branchID := uuid.New()
	actorID := uuid.New()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Role:   "admin",
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	mockOrgRepo.On("VerifyMembership", mock.Anything, orgID, actorID.String()).Return(true, nil)
	mockBranchRepo.On("GetBranchByID", mock.Anything, orgID, branchID).Return(&domain.FacilityBranch{
		ID:             branchID,
		OrganizationID: orgID,
		Name:           "Main Branch",
		Status:         "ACTIVE",
		Version:        2,
	}, nil)

	// Stale client sends version 1 while DB is at version 2
	staleName := "Updated Name by User B"
	payload := &domain.UpdateFacilityBranchPayload{
		Name:    &staleName,
		Version: 1,
	}

	mockBranchRepo.On("UpdateBranch", mock.Anything, mock.MatchedBy(func(b *domain.FacilityBranch) bool {
		return b.Version == 1
	}), actorID).Return(nil, domain.ErrOptimisticLockingConflict)

	res, err := svc.UpdateBranch(context.Background(), principal, branchID, payload)

	assert.ErrorIs(t, err, domain.ErrOptimisticLockingConflict)
	assert.Nil(t, res)
}
