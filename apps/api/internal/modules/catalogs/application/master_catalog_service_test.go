package application_test

import (
	"context"
	"testing"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/catalogs/application"
	"github.com/golangnigeria/curexal/internal/modules/catalogs/domain"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCatalogRepo struct {
	mock.Mock
}

func (m *MockCatalogRepo) ListItems(ctx context.Context, catDomain domain.CatalogDomain, category string, activeOnly bool) ([]domain.CatalogItem, error) {
	args := m.Called(ctx, catDomain, category, activeOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CatalogItem), args.Error(1)
}

func (m *MockCatalogRepo) SearchItems(ctx context.Context, catDomain domain.CatalogDomain, query string) ([]domain.CatalogItem, error) {
	args := m.Called(ctx, catDomain, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CatalogItem), args.Error(1)
}

func (m *MockCatalogRepo) GetItemByCode(ctx context.Context, catDomain domain.CatalogDomain, code string) (*domain.CatalogItem, error) {
	args := m.Called(ctx, catDomain, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CatalogItem), args.Error(1)
}

func (m *MockCatalogRepo) GetItemByID(ctx context.Context, catDomain domain.CatalogDomain, id uuid.UUID) (*domain.CatalogItem, error) {
	args := m.Called(ctx, catDomain, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CatalogItem), args.Error(1)
}

func (m *MockCatalogRepo) CreateItem(ctx context.Context, item *domain.CatalogItem, updatedBy uuid.UUID) (*domain.CatalogItem, error) {
	args := m.Called(ctx, item, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CatalogItem), args.Error(1)
}

func (m *MockCatalogRepo) UpdateItem(ctx context.Context, item *domain.CatalogItem, updatedBy uuid.UUID) (*domain.CatalogItem, error) {
	args := m.Called(ctx, item, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CatalogItem), args.Error(1)
}

func (m *MockCatalogRepo) GetSpecimenTypes(ctx context.Context) ([]domain.SpecimenType, error) {
	return nil, nil
}
func (m *MockCatalogRepo) GetTestCategories(ctx context.Context) ([]domain.TestCategory, error) {
	return nil, nil
}
func (m *MockCatalogRepo) GetUnitsOfMeasure(ctx context.Context) ([]domain.UnitOfMeasure, error) {
	return nil, nil
}
func (m *MockCatalogRepo) GetSpecialties(ctx context.Context) ([]domain.MedicalSpecialty, error) {
	return nil, nil
}
func (m *MockCatalogRepo) SearchICD10(ctx context.Context, query string) ([]domain.ICD10Code, error) {
	return nil, nil
}

type MockAuditRepo struct {
	mock.Mock
}

func (m *MockAuditRepo) Create(ctx context.Context, payload *auditDomain.CreateAuditLogPayload) (*auditDomain.AuditLog, error) {
	args := m.Called(ctx, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auditDomain.AuditLog), args.Error(1)
}
func (m *MockAuditRepo) GetByID(ctx context.Context, id uuid.UUID) (*auditDomain.AuditLog, error) {
	return nil, nil
}
func (m *MockAuditRepo) ListTenantLogs(ctx context.Context, tenantID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]auditDomain.AuditLog, error) {
	return nil, nil
}
func (m *MockAuditRepo) ListPlatformLogs(ctx context.Context, orgID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]auditDomain.AuditLog, error) {
	return nil, nil
}
func (m *MockAuditRepo) GetStats(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID) (*auditDomain.AdminStats, error) {
	return nil, nil
}
func (m *MockAuditRepo) ListAll(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID, limit int, offset int) ([]auditDomain.AuditLog, error) {
	return nil, nil
}

func TestMasterCatalogService_ListItems_Success(t *testing.T) {
	mockRepo := new(MockCatalogRepo)
	svc := application.NewMasterCatalogService(mockRepo, nil)

	expected := []domain.CatalogItem{
		{Code: "E11.9", Name: "Type 2 Diabetes", Domain: domain.ClinicalDomain, Category: "icd10"},
	}
	mockRepo.On("ListItems", mock.Anything, domain.ClinicalDomain, "icd10", true).Return(expected, nil)

	res, err := svc.ListItems(context.Background(), domain.ClinicalDomain, "icd10", true)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	mockRepo.AssertExpectations(t)
}

func TestMasterCatalogService_CreateItem_PlatformAdmin_Success(t *testing.T) {
	mockRepo := new(MockCatalogRepo)
	mockAudit := new(MockAuditRepo)
	svc := application.NewMasterCatalogService(mockRepo, mockAudit)

	adminID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: adminID.String(),
		Role:   "platform_admin",
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: true,
		},
		Identity: platformAuth.IdentityVector{
			FullName: "Catalog Admin",
		},
	}

	input := &domain.CatalogItem{
		Domain:    domain.LabDomain,
		Category:  "specimen",
		Code:      "BLOOD_EDTA",
		Name:      "EDTA Whole Blood",
		BasePrice: 4500.0,
		IsActive:  true,
	}

	created := *input
	created.ID = uuid.New()
	created.Version = 1

	mockRepo.On("CreateItem", mock.Anything, input, adminID).Return(&created, nil)
	mockAudit.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "MASTER_CATALOG_ITEM_CREATED" && p.IsPlatform == true
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.CreateItem(context.Background(), principal, input)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "BLOOD_EDTA", res.Code)
	assert.Equal(t, 4500.0, res.BasePrice)

	mockRepo.AssertExpectations(t)
	mockAudit.AssertExpectations(t)
}

func TestMasterCatalogService_CreateItem_Unauthorized(t *testing.T) {
	mockRepo := new(MockCatalogRepo)
	svc := application.NewMasterCatalogService(mockRepo, nil)

	principal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Role:   "regular_user",
	}

	input := &domain.CatalogItem{
		Domain:   domain.LabDomain,
		Category: "specimen",
		Code:     "BLOOD_EDTA",
		Name:     "EDTA Whole Blood",
	}

	res, err := svc.CreateItem(context.Background(), principal, input)

	assert.ErrorIs(t, err, domain.ErrUnauthorizedPlatformAdmin)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "CreateItem", mock.Anything, mock.Anything, mock.Anything)
}

func TestMasterCatalogService_UpdateItem_Success(t *testing.T) {
	mockRepo := new(MockCatalogRepo)
	mockAudit := new(MockAuditRepo)
	svc := application.NewMasterCatalogService(mockRepo, mockAudit)

	adminID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: adminID.String(),
		Role:   "platform_admin",
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: true,
		},
		Identity: platformAuth.IdentityVector{
			FullName: "Catalog Admin",
		},
	}

	itemID := uuid.New()
	existing := &domain.CatalogItem{
		ID:        itemID,
		Domain:    domain.LabDomain,
		Category:  "test",
		Code:      "FBC-001",
		Name:      "Full Blood Count",
		BasePrice: 4500.0,
		IsActive:  true,
		Version:   1,
	}

	newName := "Full Blood Count (Updated)"
	newPrice := 5000.0
	payload := &domain.UpdateCatalogItemPayload{
		Name:      &newName,
		BasePrice: &newPrice,
	}

	updated := *existing
	updated.Name = newName
	updated.BasePrice = newPrice
	updated.Version = 2

	mockRepo.On("GetItemByID", mock.Anything, domain.LabDomain, itemID).Return(existing, nil)
	mockRepo.On("UpdateItem", mock.Anything, mock.MatchedBy(func(item *domain.CatalogItem) bool {
		return item.Name == newName && item.BasePrice == 5000.0
	}), adminID).Return(&updated, nil)
	mockAudit.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "MASTER_CATALOG_ITEM_UPDATED" && p.IsPlatform == true
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.UpdateItem(context.Background(), principal, itemID, domain.LabDomain, payload)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, newName, res.Name)
	assert.Equal(t, 5000.0, res.BasePrice)

	mockRepo.AssertExpectations(t)
	mockAudit.AssertExpectations(t)
}

func TestMasterCatalogService_UpdateItem_OptimisticLockingConflict(t *testing.T) {
	mockRepo := new(MockCatalogRepo)
	svc := application.NewMasterCatalogService(mockRepo, nil)

	adminID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: adminID.String(),
		Role:   "platform_admin",
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: true,
		},
	}

	itemID := uuid.New()
	existing := &domain.CatalogItem{
		ID:       itemID,
		Domain:   domain.RadiologyDomain,
		Category: "modality",
		Code:     "CT",
		Name:     "CT Scan",
		Version:  1,
	}

	payload := &domain.UpdateCatalogItemPayload{
		Version: 1,
	}

	mockRepo.On("GetItemByID", mock.Anything, domain.RadiologyDomain, itemID).Return(existing, nil)
	mockRepo.On("UpdateItem", mock.Anything, mock.Anything, adminID).Return(nil, domain.ErrOptimisticLockingConflict)

	res, err := svc.UpdateItem(context.Background(), principal, itemID, domain.RadiologyDomain, payload)

	assert.ErrorIs(t, err, domain.ErrOptimisticLockingConflict)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}
