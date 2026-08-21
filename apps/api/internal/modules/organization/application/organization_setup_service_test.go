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

type MockOrgRepo struct {
	mock.Mock
}

func (m *MockOrgRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrgRepo) Create(ctx context.Context, name, slug, plan string) (*domain.Organization, error) {
	args := m.Called(ctx, name, slug, plan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrgRepo) Update(ctx context.Context, id uuid.UUID, name, slug, plan, customDomain *string, settings map[string]any) (*domain.Organization, error) {
	args := m.Called(ctx, id, name, slug, plan, customDomain, settings)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrgRepo) UpdateProfile(ctx context.Context, orgID uuid.UUID, payload *domain.UpdateOrganizationProfilePayload, actorID uuid.UUID) (*domain.Organization, error) {
	args := m.Called(ctx, orgID, payload, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrgRepo) UpdateSetupState(ctx context.Context, orgID uuid.UUID, newState domain.SetupState, step int, actorID uuid.UUID) (*domain.Organization, error) {
	args := m.Called(ctx, orgID, newState, step, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrgRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *MockOrgRepo) GetSettings(ctx context.Context, orgID uuid.UUID) (*domain.OrganizationSettings, error) {
	return nil, nil
}

func (m *MockOrgRepo) UpdateSettings(ctx context.Context, orgID uuid.UUID, logoURL, themeBranding, customDomain, supportEmail, supportPhone, cacNumber, tinNumber, taxNumber, businessType, address, timezone, currency, language *string) (*domain.OrganizationSettings, error) {
	return nil, nil
}

func (m *MockOrgRepo) List(ctx context.Context, userID string, isPlatformAdmin bool) ([]domain.Organization, error) {
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

func TestOrganizationSetupService_GetProfile_Success(t *testing.T) {
	mockOrgRepo := new(MockOrgRepo)
	svc := application.NewOrganizationSetupService(mockOrgRepo, nil)

	orgID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	expectedOrg := &domain.Organization{
		ID:         orgID,
		Name:       "Curexal Medical Center",
		SetupState: domain.SetupStatePendingRegistration,
		Plan:       "smart",
	}

	mockOrgRepo.On("GetByID", mock.Anything, orgID).Return(expectedOrg, nil)

	res, err := svc.GetProfile(context.Background(), principal)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Curexal Medical Center", res.Name)
	assert.Equal(t, "smart", res.Plan) // Plan immutable
	mockOrgRepo.AssertExpectations(t)
}

func TestOrganizationSetupService_UpdateProfile_AdvancesSetupState(t *testing.T) {
	mockOrgRepo := new(MockOrgRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewOrganizationSetupService(mockOrgRepo, mockAuditRepo)

	orgID := uuid.New()
	userID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: userID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
		Identity: platformAuth.IdentityVector{
			FullName: "Org Owner",
		},
	}

	existingOrg := &domain.Organization{
		ID:         orgID,
		Name:       "Old Name",
		SetupState: domain.SetupStatePendingRegistration,
		Version:    1,
	}

	mockOrgRepo.On("GetByID", mock.Anything, orgID).Return(existingOrg, nil)

	newName := "Curexal Hospital HQ"
	payload := &domain.UpdateOrganizationProfilePayload{
		Name:    &newName,
		Version: 1,
	}

	updatedProfile := *existingOrg
	updatedProfile.Name = newName
	updatedProfile.Version = 2

	mockOrgRepo.On("UpdateProfile", mock.Anything, orgID, payload, userID).Return(&updatedProfile, nil)

	advancedProfile := updatedProfile
	advancedProfile.SetupState = domain.SetupStateProfileCompleted
	advancedProfile.SetupStep = 2

	mockOrgRepo.On("UpdateSetupState", mock.Anything, orgID, domain.SetupStateProfileCompleted, 2, userID).Return(&advancedProfile, nil)

	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "ORGANIZATION_PROFILE_UPDATED"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.UpdateProfile(context.Background(), principal, payload)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Curexal Hospital HQ", res.Name)
	assert.Equal(t, domain.SetupStateProfileCompleted, res.SetupState)

	mockOrgRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestOrganizationSetupService_SetupWizard_InvalidStateTransition_Fails(t *testing.T) {
	mockOrgRepo := new(MockOrgRepo)
	svc := application.NewOrganizationSetupService(mockOrgRepo, nil)

	orgID := uuid.New()
	userID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: userID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	// Currently in PENDING_REGISTRATION state
	existingOrg := &domain.Organization{
		ID:         orgID,
		SetupState: domain.SetupStatePendingRegistration,
	}

	mockOrgRepo.On("GetByID", mock.Anything, orgID).Return(existingOrg, nil)

	// Attempting to jump directly to UNDER_REVIEW without PROFILE_COMPLETED and DOCUMENTS_SUBMITTED
	res, err := svc.SubmitForReview(context.Background(), principal)

	assert.ErrorIs(t, err, domain.ErrInvalidSetupStateTransition)
	assert.Nil(t, res)
	mockOrgRepo.AssertNotCalled(t, "UpdateSetupState", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestOrganizationSetupService_VerifyOrganization_PlatformAdminOnly(t *testing.T) {
	mockOrgRepo := new(MockOrgRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewOrganizationSetupService(mockOrgRepo, mockAuditRepo)

	adminID := uuid.New()
	adminPrincipal := &middleware.AuthenticatedPrincipal{
		UserID: adminID.String(),
		Role:   "platform_admin",
		Platform: platformAuth.PlatformVector{
			IsPlatformAdmin: true,
		},
	}

	targetOrgID := uuid.New()
	verifiedOrg := &domain.Organization{
		ID:         targetOrgID,
		SetupState: domain.SetupStateVerified,
		SetupStep:  5,
	}

	mockOrgRepo.On("UpdateSetupState", mock.Anything, targetOrgID, domain.SetupStateVerified, 5, adminID).Return(verifiedOrg, nil)
	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "ORGANIZATION_VERIFIED"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.VerifyOrganization(context.Background(), adminPrincipal, targetOrgID, true, "Approved regulatory docs")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, domain.SetupStateVerified, res.SetupState)

	// Non-admin call should fail
	userPrincipal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Role:   "regular_user",
	}

	resUser, errUser := svc.VerifyOrganization(context.Background(), userPrincipal, targetOrgID, true, "Reason")
	assert.Error(t, errUser)
	assert.Nil(t, resUser)
}

func TestOrganizationSetupService_StateTransitions_Validations(t *testing.T) {
	assert.True(t, domain.IsValidSetupTransition(domain.SetupStatePendingRegistration, domain.SetupStateProfileCompleted))
	assert.True(t, domain.IsValidSetupTransition(domain.SetupStateProfileCompleted, domain.SetupStateDocumentsSubmitted))
	assert.True(t, domain.IsValidSetupTransition(domain.SetupStateDocumentsSubmitted, domain.SetupStateUnderReview))
	assert.True(t, domain.IsValidSetupTransition(domain.SetupStateUnderReview, domain.SetupStateVerified))
	assert.True(t, domain.IsValidSetupTransition(domain.SetupStateUnderReview, domain.SetupStateRejected))

	// Invalid state jumps
	assert.False(t, domain.IsValidSetupTransition(domain.SetupStatePendingRegistration, domain.SetupStateVerified))
	assert.False(t, domain.IsValidSetupTransition(domain.SetupStatePendingRegistration, domain.SetupStateUnderReview))
}
