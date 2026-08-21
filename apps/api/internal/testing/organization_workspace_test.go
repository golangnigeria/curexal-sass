package testing

import (
	"context"
	"testing"

	orgDomain "github.com/golangnigeria/curexal/internal/modules/organization/domain"
	platformApp "github.com/golangnigeria/curexal/internal/modules/platform/application"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Organization Repository
type MockOrgRepo struct {
	mock.Mock
}

func (m *MockOrgRepo) GetByID(ctx context.Context, id uuid.UUID) (*orgDomain.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orgDomain.Organization), args.Error(1)
}

func (m *MockOrgRepo) Create(ctx context.Context, name, slug, plan string) (*orgDomain.Organization, error) {
	args := m.Called(ctx, name, slug, plan)
	return args.Get(0).(*orgDomain.Organization), args.Error(1)
}

func (m *MockOrgRepo) Update(ctx context.Context, id uuid.UUID, name, slug, plan, customDomain *string, settings map[string]any) (*orgDomain.Organization, error) {
	args := m.Called(ctx, id, name, slug, plan, customDomain, settings)
	return args.Get(0).(*orgDomain.Organization), args.Error(1)
}

func (m *MockOrgRepo) UpdateProfile(ctx context.Context, orgID uuid.UUID, payload *orgDomain.UpdateOrganizationProfilePayload, actorID uuid.UUID) (*orgDomain.Organization, error) {
	args := m.Called(ctx, orgID, payload, actorID)
	return args.Get(0).(*orgDomain.Organization), args.Error(1)
}

func (m *MockOrgRepo) UpdateSetupState(ctx context.Context, orgID uuid.UUID, newState orgDomain.SetupState, step int, actorID uuid.UUID) (*orgDomain.Organization, error) {
	args := m.Called(ctx, orgID, newState, step, actorID)
	return args.Get(0).(*orgDomain.Organization), args.Error(1)
}

func (m *MockOrgRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrgRepo) List(ctx context.Context, userID string, isPlatformAdmin bool) ([]orgDomain.Organization, error) {
	args := m.Called(ctx, userID, isPlatformAdmin)
	return args.Get(0).([]orgDomain.Organization), args.Error(1)
}

func (m *MockOrgRepo) GetSettings(ctx context.Context, orgID uuid.UUID) (*orgDomain.OrganizationSettings, error) {
	args := m.Called(ctx, orgID)
	return args.Get(0).(*orgDomain.OrganizationSettings), args.Error(1)
}

func (m *MockOrgRepo) UpdateSettings(ctx context.Context, orgID uuid.UUID, logoURL, themeBranding, customDomain, supportEmail, supportPhone, cacNumber, tinNumber, taxNumber, businessType, address, timezone, currency, language *string) (*orgDomain.OrganizationSettings, error) {
	args := m.Called(ctx, orgID, logoURL, themeBranding, customDomain, supportEmail, supportPhone, cacNumber, tinNumber, taxNumber, businessType, address, timezone, currency, language)
	return args.Get(0).(*orgDomain.OrganizationSettings), args.Error(1)
}

// Mock Tenant Repository
type MockTenantRepo struct {
	mock.Mock
}

func (m *MockTenantRepo) GetTenantByID(ctx context.Context, id uuid.UUID) (*orgDomain.Tenant, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*orgDomain.Tenant), args.Error(1)
}

func (m *MockTenantRepo) GetTenantBySlug(ctx context.Context, slug string) (*orgDomain.Tenant, error) {
	args := m.Called(ctx, slug)
	return args.Get(0).(*orgDomain.Tenant), args.Error(1)
}

func (m *MockTenantRepo) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	args := m.Called(ctx, slug)
	return args.Bool(0), args.Error(1)
}

func (m *MockTenantRepo) CountActiveMembers(ctx context.Context, tenantID uuid.UUID) (int, error) {
	args := m.Called(ctx, tenantID)
	return args.Int(0), args.Error(1)
}

func (m *MockTenantRepo) ListBranchesByOrgID(ctx context.Context, orgID uuid.UUID) ([]orgDomain.Tenant, error) {
	args := m.Called(ctx, orgID)
	return args.Get(0).([]orgDomain.Tenant), args.Error(1)
}

func (m *MockTenantRepo) CreateTenant(ctx context.Context, userID string, name, slug, orgID, location, phone, address string, logoURL, currency *string, modules []string) (*orgDomain.Tenant, error) {
	args := m.Called(ctx, userID, name, slug, orgID, location, phone, address, logoURL, currency, modules)
	return args.Get(0).(*orgDomain.Tenant), args.Error(1)
}

func (m *MockTenantRepo) UpdateTenant(ctx context.Context, id uuid.UUID, name, slug, logoURL, currency *string, settingsJSON *string) (*orgDomain.Tenant, error) {
	args := m.Called(ctx, id, name, slug, logoURL, currency, settingsJSON)
	return args.Get(0).(*orgDomain.Tenant), args.Error(1)
}

func (m *MockTenantRepo) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTenantRepo) ListTenants(ctx context.Context) ([]orgDomain.Tenant, error) {
	args := m.Called(ctx)
	return args.Get(0).([]orgDomain.Tenant), args.Error(1)
}

// TestOrganizationBuildAndWorkspaceBootstrap tests Organization creation and Bootstrap Contract Generation.
func TestOrganizationBuildAndWorkspaceBootstrap(t *testing.T) {
	ctx := context.Background()
	mockOrgRepo := new(MockOrgRepo)
	mockTenantRepo := new(MockTenantRepo)

	orgID := uuid.New()
	tenantID := uuid.New()

	testOrg := &orgDomain.Organization{
		ID:   orgID,
		Name: "Everight Healthcare Network",
		Slug: "everight",
		Plan: "enterprise",
	}

	testTenant := &orgDomain.Tenant{
		ID:             tenantID,
		OrganizationID: orgID.String(),
		Name:           "Everight Main Diagnostic Facility",
		Slug:           "main-facility",
		Currency:       "NGN",
		EnabledModules: []string{"laboratory", "clinical", "pharmacy", "billing"},
	}

	mockOrgRepo.On("List", ctx, "usr_org_owner_001", false).Return([]orgDomain.Organization{*testOrg}, nil)
	mockTenantRepo.On("ListTenants", ctx).Return([]orgDomain.Tenant{*testTenant}, nil)

	// Test BootstrapBuilder
	builder := platformApp.NewBootstrapBuilder(mockOrgRepo, mockTenantRepo, nil)

	principal := &middleware.AuthenticatedPrincipal{
		UserID:         "usr_org_owner_001",
		TenantID:       tenantID.String(),
		OrganizationID: orgID.String(),
		Role:           "owner",
		Identity: middleware.IdentityVector{
			Email:    "org.owner@test.curexal.local",
			FullName: "Dr. Chidi Okezie",
		},
	}

	bootstrap, err := builder.BuildBootstrap(ctx, principal)

	assert.NoError(t, err)
	assert.NotNil(t, bootstrap)
	assert.Equal(t, "usr_org_owner_001", bootstrap.Identity.ID)
	assert.Equal(t, "org.owner@test.curexal.local", bootstrap.Identity.Email)
	assert.Equal(t, "organization", bootstrap.Contexts.Current)
	assert.Equal(t, orgID.String(), bootstrap.Organization.ID)
	assert.Equal(t, "Everight Healthcare Network", bootstrap.Organization.Name)
	assert.Equal(t, tenantID.String(), bootstrap.Workspace.ID)
	assert.Equal(t, "main-facility", bootstrap.Workspace.Slug)
	assert.True(t, len(bootstrap.Navigation) > 0)
}
