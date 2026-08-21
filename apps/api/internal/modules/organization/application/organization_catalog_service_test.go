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

type MockCatalogRepo struct {
	mock.Mock
}

func (m *MockCatalogRepo) ListCatalogItems(ctx context.Context, orgID uuid.UUID, domainType string) ([]domain.OrganizationCatalogItem, error) {
	args := m.Called(ctx, orgID, domainType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.OrganizationCatalogItem), args.Error(1)
}

func (m *MockCatalogRepo) GetCatalogItemByID(ctx context.Context, orgID, itemID uuid.UUID) (*domain.OrganizationCatalogItem, error) {
	args := m.Called(ctx, orgID, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.OrganizationCatalogItem), args.Error(1)
}

func (m *MockCatalogRepo) CreateCatalogItem(ctx context.Context, item *domain.OrganizationCatalogItem, actorID uuid.UUID) (*domain.OrganizationCatalogItem, error) {
	args := m.Called(ctx, item, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.OrganizationCatalogItem), args.Error(1)
}

func (m *MockCatalogRepo) UpdateCatalogItem(ctx context.Context, item *domain.OrganizationCatalogItem, actorID uuid.UUID) (*domain.OrganizationCatalogItem, error) {
	args := m.Called(ctx, item, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.OrganizationCatalogItem), args.Error(1)
}

func (m *MockCatalogRepo) SetBranchPriceOverride(ctx context.Context, override *domain.BranchPriceOverride, actorID uuid.UUID) (*domain.BranchPriceOverride, error) {
	args := m.Called(ctx, override, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BranchPriceOverride), args.Error(1)
}

func (m *MockCatalogRepo) GetBranchPriceOverride(ctx context.Context, orgID, branchID, itemID uuid.UUID) (*domain.BranchPriceOverride, error) {
	args := m.Called(ctx, orgID, branchID, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BranchPriceOverride), args.Error(1)
}

func (m *MockCatalogRepo) ListBranchPriceOverrides(ctx context.Context, orgID, branchID uuid.UUID) ([]domain.BranchPriceOverride, error) {
	args := m.Called(ctx, orgID, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.BranchPriceOverride), args.Error(1)
}

func (m *MockCatalogRepo) CreateInsuranceProvider(ctx context.Context, provider *domain.InsuranceProvider, actorID uuid.UUID) (*domain.InsuranceProvider, error) {
	args := m.Called(ctx, provider, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.InsuranceProvider), args.Error(1)
}

func (m *MockCatalogRepo) ListInsuranceProviders(ctx context.Context, orgID uuid.UUID) ([]domain.InsuranceProvider, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.InsuranceProvider), args.Error(1)
}

func TestOrganizationCatalogService_CreateCatalogItem_Success(t *testing.T) {
	mockCatalogRepo := new(MockCatalogRepo)
	mockBranchRepo := new(MockBranchRepo)
	mockOrgRepo := new(MockOrgRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewOrganizationCatalogService(mockCatalogRepo, mockBranchRepo, mockOrgRepo, mockAuditRepo)

	orgID := uuid.New()
	actorID := uuid.New()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	payload := &domain.CreateCatalogItemPayload{
		DomainType: "LABORATORY",
		Code:       "fbc-local",
		Name:       "Full Blood Count (Local Override)",
		BasePrice:  4500.00,
	}

	createdItem := &domain.OrganizationCatalogItem{
		ID:             uuid.New(),
		OrganizationID: orgID,
		DomainType:     "LABORATORY",
		Code:           "fbc-local",
		Name:           "Full Blood Count (Local Override)",
		BasePrice:      4500.00,
		Currency:       "NGN",
		IsActive:       true,
		Version:        1,
	}

	mockCatalogRepo.On("CreateCatalogItem", mock.Anything, mock.Anything, actorID).Return(createdItem, nil)
	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "CATALOG_ITEM_CREATED"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.CreateCatalogItem(context.Background(), principal, payload)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "LABORATORY", res.DomainType)
	assert.Equal(t, 4500.00, res.BasePrice)

	mockCatalogRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestOrganizationCatalogService_CreateCatalogItem_InvalidDomain_Fails(t *testing.T) {
	mockCatalogRepo := new(MockCatalogRepo)
	svc := application.NewOrganizationCatalogService(mockCatalogRepo, nil, nil, nil)

	orgID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	payload := &domain.CreateCatalogItemPayload{
		DomainType: "UNRECOGNIZED_DOMAIN",
		Code:       "invalid-code",
		Name:       "Invalid Test Item",
	}

	res, err := svc.CreateCatalogItem(context.Background(), principal, payload)

	assert.ErrorIs(t, err, domain.ErrInvalidCatalogDomain)
	assert.Nil(t, res)

	mockCatalogRepo.AssertNotCalled(t, "CreateCatalogItem", mock.Anything, mock.Anything, mock.Anything)
}

func TestOrganizationCatalogService_SetBranchPriceOverride_Success(t *testing.T) {
	mockCatalogRepo := new(MockCatalogRepo)
	mockBranchRepo := new(MockBranchRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewOrganizationCatalogService(mockCatalogRepo, mockBranchRepo, nil, mockAuditRepo)

	orgID := uuid.New()
	actorID := uuid.New()
	itemID := uuid.New()
	branchID := uuid.New()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	payload := &domain.SetBranchPricePayload{
		FacilityBranchID: branchID,
		OverridePrice:    6000.00,
	}

	mockBranchRepo.On("GetBranchByID", mock.Anything, orgID, branchID).Return(&domain.FacilityBranch{ID: branchID, OrganizationID: orgID}, nil)
	mockCatalogRepo.On("GetCatalogItemByID", mock.Anything, orgID, itemID).Return(&domain.OrganizationCatalogItem{ID: itemID, OrganizationID: orgID, BasePrice: 4500.00}, nil)

	overrideRes := &domain.BranchPriceOverride{
		ID:               uuid.New(),
		OrganizationID:   orgID,
		FacilityBranchID: branchID,
		CatalogItemID:    itemID,
		OverridePrice:    6000.00,
	}

	mockCatalogRepo.On("SetBranchPriceOverride", mock.Anything, mock.Anything, actorID).Return(overrideRes, nil)
	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "BRANCH_PRICE_OVERRIDDEN"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.SetBranchPriceOverride(context.Background(), principal, itemID, payload)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 6000.00, res.OverridePrice)

	mockBranchRepo.AssertExpectations(t)
	mockCatalogRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestOrganizationCatalogService_CreateInsuranceProvider_Success(t *testing.T) {
	mockCatalogRepo := new(MockCatalogRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewOrganizationCatalogService(mockCatalogRepo, nil, nil, mockAuditRepo)

	orgID := uuid.New()
	actorID := uuid.New()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	covVal := 80.00
	payload := &domain.CreateInsuranceProviderPayload{
		Name:               "Hygeia HMO",
		Code:               "hygeia-hmo",
		CoveragePercentage: &covVal,
	}

	createdProvider := &domain.InsuranceProvider{
		ID:                 uuid.New(),
		OrganizationID:     orgID,
		Name:               "Hygeia HMO",
		Code:               "hygeia-hmo",
		CoveragePercentage: 80.00,
		IsActive:           true,
	}

	mockCatalogRepo.On("CreateInsuranceProvider", mock.Anything, mock.Anything, actorID).Return(createdProvider, nil)
	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "INSURANCE_PROVIDER_ADDED"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.CreateInsuranceProvider(context.Background(), principal, payload)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Hygeia HMO", res.Name)
	assert.Equal(t, 80.00, res.CoveragePercentage)

	mockCatalogRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}
