package testing

import (
	"context"
	"testing"

	"github.com/golangnigeria/curexal/internal/bootstrap"
	"github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/platform/application"
	platformDomain "github.com/golangnigeria/curexal/internal/modules/platform/domain"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNavigationRepository struct {
	mock.Mock
}

func (m *MockNavigationRepository) GetNavigationItemsByScope(
	ctx context.Context,
	scope string,
	enabledModules []string,
	userPermissions []string,
	isSuperAdminOrOwner bool,
) ([]platformDomain.NavigationItem, error) {
	args := m.Called(ctx, scope, enabledModules, userPermissions, isSuperAdminOrOwner)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]platformDomain.NavigationItem), args.Error(1)
}

func (m *MockNavigationRepository) GetAllActiveNavigationItems(ctx context.Context) ([]platformDomain.NavigationItem, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]platformDomain.NavigationItem), args.Error(1)
}

func TestCanonicalNavigationAPI_RouteAndScoping(t *testing.T) {
	e := echo.New()
	logger := zerolog.Nop()
	s := &server.Server{
		Echo:   e,
		Logger: &logger,
	}
	reg := bootstrap.InitModules(s)

	assert.NotNil(t, reg)
	assert.NotNil(t, reg.Platform)
	assert.NotNil(t, reg.Platform.NavigationHandler)

	// Verify route is registered
	hasNavRoute := false
	for _, route := range e.Routes() {
		if route.Method == "GET" && route.Path == "/api/v1/navigation" {
			hasNavRoute = true
			break
		}
	}
	assert.True(t, hasNavRoute, "GET /api/v1/navigation must be registered on Echo router")
}

func TestNavigationService_WorkspaceInterpolation(t *testing.T) {
	ctx := context.Background()
	mockNavRepo := new(MockNavigationRepository)

	sampleItems := []platformDomain.NavigationItem{
		{
			ID:           "nav_wsp_dashboard",
			Key:          "workspace.dashboard",
			ContextScope: "workspace",
			Title:        "Overview",
			Icon:         "LayoutDashboard",
			Path:         "/:branch/dashboard",
			Order:        1,
			Status:       platformDomain.NavigationStatusActive,
			IsVisible:    true,
			IsActive:     true,
		},
		{
			ID:           "nav_wsp_lab",
			Key:          "workspace.laboratory",
			ContextScope: "workspace",
			Title:        "Laboratory",
			Icon:         "Microscope",
			Path:         "/:branch/laboratory",
			Order:        2,
			Status:       platformDomain.NavigationStatusActive,
			IsVisible:    true,
			IsActive:     true,
		},
	}

	mockNavRepo.On(
		"GetNavigationItemsByScope",
		ctx,
		"workspace",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(sampleItems, nil)

	principal := &auth.AuthenticatedPrincipal{
		UserID:         "usr_clinician_01",
		TenantID:       "branch_owerri_01",
		OrganizationID: "org_everight_01",
		Role:           "clinician",
		Workspace: auth.WorkspaceVector{
			ActiveWorkspaceID: "branch_owerri_01",
			WorkspaceName:     "owerri",
		},
	}

	navService := application.NewNavigationService(mockNavRepo, nil, nil, nil)
	res, err := navService.GetNavigation(ctx, principal, "app.curexal.space", "owerri", "workspace")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "workspace", res.Context.Type)
	assert.NotNil(t, res.Context.BranchSlug)
	assert.Equal(t, "owerri", *res.Context.BranchSlug)
	assert.Len(t, res.Items, 2)
	assert.Equal(t, "/owerri/dashboard", res.Items[0].Path)
	assert.Equal(t, "/owerri/laboratory", res.Items[1].Path)
}

func TestNavigationService_PlatformPrivilegeEscalation_Rejected(t *testing.T) {
	ctx := context.Background()
	mockNavRepo := new(MockNavigationRepository)

	// An organization admin attempts to request platform navigation
	orgAdminPrincipal := &auth.AuthenticatedPrincipal{
		UserID:         "usr_org_admin_01",
		OrganizationID: "org_everight_01",
		Role:           "org_admin",
		Platform: auth.PlatformVector{
			IsPlatformAdmin: false,
			IsSuperAdmin:    false,
			IsPlatformStaff: false,
		},
	}

	navService := application.NewNavigationService(mockNavRepo, nil, nil, nil)
	res, err := navService.GetNavigation(ctx, orgAdminPrincipal, "app.localhost:5002", "", "platform")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, platformDomain.ErrUnauthorizedScope, err, "Organization admin MUST NOT receive platform navigation on request")
}

func TestNavigationService_PlatformStaff_Authorized(t *testing.T) {
	ctx := context.Background()
	mockNavRepo := new(MockNavigationRepository)

	sampleItems := []platformDomain.NavigationItem{
		{
			ID:           "nav_plat_dashboard",
			Key:          "platform.dashboard",
			ContextScope: "platform",
			Title:        "Platform Dashboard",
			Icon:         "LayoutDashboard",
			Path:         "/platform/dashboard",
			Order:        1,
			Status:       platformDomain.NavigationStatusActive,
			IsVisible:    true,
			IsActive:     true,
		},
	}

	mockNavRepo.On(
		"GetNavigationItemsByScope",
		ctx,
		"platform",
		mock.Anything,
		mock.Anything,
		true,
	).Return(sampleItems, nil)

	platformPrincipal := &auth.AuthenticatedPrincipal{
		UserID: "usr_superadmin_01",
		Role:   "super_admin",
		Platform: auth.PlatformVector{
			IsSuperAdmin:    true,
			IsPlatformAdmin: true,
			IsPlatformStaff: true,
		},
	}

	navService := application.NewNavigationService(mockNavRepo, nil, nil, nil)
	res, err := navService.GetNavigation(ctx, platformPrincipal, "app.localhost:5002", "", "platform")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "platform", res.Context.Type)
	assert.Len(t, res.Items, 1)
	assert.Equal(t, "/platform/dashboard", res.Items[0].Path)
}

func TestNavigationService_OrganizationScoping_OnPlatformHost(t *testing.T) {
	ctx := context.Background()
	mockNavRepo := new(MockNavigationRepository)

	sampleItems := []platformDomain.NavigationItem{
		{
			ID:           "nav_org_dashboard",
			Key:          "organization.dashboard",
			ContextScope: "organization",
			Title:        "Executive HQ Dashboard",
			Icon:         "LayoutDashboard",
			Path:         "/organization/dashboard",
			Order:        1,
			Status:       platformDomain.NavigationStatusActive,
			IsVisible:    true,
			IsActive:     true,
		},
		{
			ID:           "nav_org_branches",
			Key:          "organization.branches",
			ContextScope: "organization",
			Title:        "Branch Facilities",
			Icon:         "Building2",
			Path:         "/organization/branches",
			Order:        2,
			Status:       platformDomain.NavigationStatusActive,
			IsVisible:    true,
			IsActive:     true,
		},
	}

	mockNavRepo.On(
		"GetNavigationItemsByScope",
		ctx,
		"organization",
		mock.Anything,
		mock.Anything,
		true,
	).Return(sampleItems, nil)

	orgPrincipal := &auth.AuthenticatedPrincipal{
		UserID:         "usr_owner_01",
		OrganizationID: "org_everight_01",
		Role:           "owner",
		Platform: auth.PlatformVector{
			IsPlatformAdmin: false,
			IsSuperAdmin:    false,
		},
	}

	navService := application.NewNavigationService(mockNavRepo, nil, nil, nil)
	// Testing organization request when visiting app.localhost:5002/organization/dashboard
	res, err := navService.GetNavigation(ctx, orgPrincipal, "app.localhost:5002", "organization", "organization")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "organization", res.Context.Type)
	assert.Len(t, res.Items, 2)
	assert.Equal(t, "/organization/dashboard", res.Items[0].Path)
	assert.Equal(t, "/organization/branches", res.Items[1].Path)
}

func TestNavigationService_BranchDoctor_RejectedFromOrganizationScope(t *testing.T) {
	ctx := context.Background()
	mockNavRepo := new(MockNavigationRepository)

	// An employed branch physician attempts to request organization scope navigation
	doctorPrincipal := &auth.AuthenticatedPrincipal{
		UserID:         "usr_doctor_01",
		OrganizationID: "org_everight_01",
		TenantID:       "branch_owerri_01",
		Role:           "doctor",
		Permissions:    []string{"workspace:clinical:read", "workspace:patient:read"},
		Platform: auth.PlatformVector{
			IsPlatformAdmin: false,
			IsSuperAdmin:    false,
		},
	}

	navService := application.NewNavigationService(mockNavRepo, nil, nil, nil)
	res, err := navService.GetNavigation(ctx, doctorPrincipal, "curexal-clinic.localhost:5002", "", "organization")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, platformDomain.ErrUnauthorizedScope, err, "Branch doctor MUST NOT receive organization scope navigation")
}

func TestNavigationService_PatientScoping(t *testing.T) {
	ctx := context.Background()
	mockNavRepo := new(MockNavigationRepository)

	sampleItems := []platformDomain.NavigationItem{
		{
			ID:           "nav_pat_dashboard",
			Key:          "patient.dashboard",
			ContextScope: "patient",
			Title:        "Care Journey",
			Icon:         "LayoutDashboard",
			Path:         "/dashboard",
			Order:        1,
			Status:       platformDomain.NavigationStatusActive,
			IsVisible:    true,
			IsActive:     true,
		},
	}

	mockNavRepo.On(
		"GetNavigationItemsByScope",
		ctx,
		"patient",
		mock.Anything,
		mock.Anything,
		false,
	).Return(sampleItems, nil)

	patientPrincipal := &auth.AuthenticatedPrincipal{
		UserID: "usr_patient_01",
		Role:   "patient",
	}

	navService := application.NewNavigationService(mockNavRepo, nil, nil, nil)
	res, err := navService.GetNavigation(ctx, patientPrincipal, "patient.localhost:5003", "", "patient")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "patient", res.Context.Type)
	assert.Len(t, res.Items, 1)
	assert.Equal(t, "/dashboard", res.Items[0].Path)
}

func TestNavigationService_UserPermissionDifferentiation(t *testing.T) {
	ctx := context.Background()

	labItems := []platformDomain.NavigationItem{
		{
			ID:                 "nav_wsp_lab",
			Key:                "workspace.laboratory",
			ContextScope:       "workspace",
			Title:              "Laboratory",
			Path:               "/:branch/laboratory",
			RequiredPermission: stringPtr("workspace:sample:receive"),
			Status:             platformDomain.NavigationStatusActive,
			IsVisible:          true,
			IsActive:           true,
		},
	}

	directorItems := []platformDomain.NavigationItem{
		{
			ID:                 "nav_wsp_lab",
			Key:                "workspace.laboratory",
			ContextScope:       "workspace",
			Title:              "Laboratory",
			Path:               "/:branch/laboratory",
			RequiredPermission: stringPtr("workspace:sample:receive"),
			Status:             platformDomain.NavigationStatusActive,
			IsVisible:          true,
			IsActive:           true,
		},
		{
			ID:                 "nav_wsp_billing",
			Key:                "workspace.billing",
			ContextScope:       "workspace",
			Title:              "Billing",
			Path:               "/:branch/billing",
			RequiredPermission: stringPtr("workspace:billing:create"),
			Status:             platformDomain.NavigationStatusActive,
			IsVisible:          true,
			IsActive:           true,
		},
	}

	// 1. Test Scientist (no billing permission)
	mockNavRepo1 := new(MockNavigationRepository)
	mockNavRepo1.On(
		"GetNavigationItemsByScope",
		ctx,
		"workspace",
		mock.Anything,
		[]string{"workspace:sample:receive"},
		false,
	).Return(labItems, nil)

	scientistPrincipal := &auth.AuthenticatedPrincipal{
		UserID:         "usr_scientist_01",
		TenantID:       "branch_owerri_01",
		OrganizationID: "org_everight_01",
		Role:           "lab_scientist",
		Permissions:    []string{"workspace:sample:receive"},
	}

	navService1 := application.NewNavigationService(mockNavRepo1, nil, nil, nil)
	res1, err1 := navService1.GetNavigation(ctx, scientistPrincipal, "app.curexal.space", "owerri", "workspace")

	assert.NoError(t, err1)
	assert.NotNil(t, res1)
	assert.Len(t, res1.Items, 1)
	assert.Equal(t, "/owerri/laboratory", res1.Items[0].Path)

	// 2. Test Director (with billing permission)
	mockNavRepo2 := new(MockNavigationRepository)
	mockNavRepo2.On(
		"GetNavigationItemsByScope",
		ctx,
		"workspace",
		mock.Anything,
		[]string{"workspace:sample:receive", "workspace:billing:create"},
		false,
	).Return(directorItems, nil)

	directorPrincipal := &auth.AuthenticatedPrincipal{
		UserID:         "usr_director_01",
		TenantID:       "branch_owerri_01",
		OrganizationID: "org_everight_01",
		Role:           "lab_director",
		Permissions:    []string{"workspace:sample:receive", "workspace:billing:create"},
	}

	navService2 := application.NewNavigationService(mockNavRepo2, nil, nil, nil)
	res2, err2 := navService2.GetNavigation(ctx, directorPrincipal, "app.curexal.space", "owerri", "workspace")

	assert.NoError(t, err2)
	assert.NotNil(t, res2)
	assert.Len(t, res2.Items, 2)
	assert.Equal(t, "/owerri/laboratory", res2.Items[0].Path)
	assert.Equal(t, "/owerri/billing", res2.Items[1].Path)
}

func TestSecurityMatrix_CrossScopeRejections(t *testing.T) {
	ctx := context.Background()
	mockNavRepo := new(MockNavigationRepository)

	testCases := []struct {
		name          string
		role          string
		reqScope      string
		expectedError error
	}{
		{"Doctor cannot access Executive HQ", "doctor", "organization", platformDomain.ErrUnauthorizedScope},
		{"Nurse cannot access Executive HQ", "nurse", "organization", platformDomain.ErrUnauthorizedScope},
		{"Scientist cannot access Executive HQ", "scientist", "organization", platformDomain.ErrUnauthorizedScope},
		{"Pharmacist cannot access Executive HQ", "pharmacist", "organization", platformDomain.ErrUnauthorizedScope},
		{"Cashier cannot access Executive HQ", "cashier", "organization", platformDomain.ErrUnauthorizedScope},
		{"Phlebotomist cannot access Executive HQ", "phlebotomist", "organization", platformDomain.ErrUnauthorizedScope},
		{"Branch Manager cannot access Platform Console", "branch_manager", "platform", platformDomain.ErrUnauthorizedScope},
		{"Organization Owner cannot access Platform Console", "owner", "platform", platformDomain.ErrUnauthorizedScope},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			principal := &auth.AuthenticatedPrincipal{
				UserID:         "usr_test_01",
				OrganizationID: "org_everight_01",
				TenantID:       "branch_owerri_01",
				Role:           tc.role,
				Platform: auth.PlatformVector{
					IsPlatformAdmin: false,
					IsSuperAdmin:    false,
					IsPlatformStaff: false,
				},
			}

			navService := application.NewNavigationService(mockNavRepo, nil, nil, nil)
			res, err := navService.GetNavigation(ctx, principal, "curexal-clinic.localhost:5002", "", tc.reqScope)

			assert.Error(t, err)
			assert.Nil(t, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestSecurityMatrix_SharedCoreMPI_Access(t *testing.T) {
	ctx := context.Background()

	// 1. Shared Core Patient Registry Item
	corePatientItem := platformDomain.NavigationItem{
		ID:                 "nav_core_patient",
		Key:                "core.patient",
		ContextScope:       "workspace",
		Title:              "Patient Registry (MPI)",
		Path:               "/:branch/reception",
		RequiredPermission: stringPtr("core.patient.create"),
		Status:             platformDomain.NavigationStatusActive,
		IsVisible:          true,
		IsActive:           true,
	}

	// 2. Billing POS Item
	billingItem := platformDomain.NavigationItem{
		ID:                 "nav_billing_pos",
		Key:                "billing.pos",
		ContextScope:       "workspace",
		Title:              "Billing (POS)",
		Path:               "/:branch/billing",
		RequiredPermission: stringPtr("billing.payment.create"),
		Status:             platformDomain.NavigationStatusActive,
		IsVisible:          true,
		IsActive:           true,
	}

	mockNavRepo := new(MockNavigationRepository)
	mockNavRepo.On(
		"GetNavigationItemsByScope",
		ctx,
		"workspace",
		mock.Anything,
		[]string{"core.patient.create", "billing.payment.create"},
		false,
	).Return([]platformDomain.NavigationItem{corePatientItem, billingItem}, nil)

	// Reception Principal with Core MPI + Billing (No clinical consultation permissions)
	receptionPrincipal := &auth.AuthenticatedPrincipal{
		UserID:         "usr_reception_01",
		TenantID:       "branch_owerri_01",
		OrganizationID: "org_everight_01",
		Role:           "receptionist",
		Permissions:    []string{"core.patient.create", "billing.payment.create"},
	}

	navService := application.NewNavigationService(mockNavRepo, nil, nil, nil)
	res, err := navService.GetNavigation(ctx, receptionPrincipal, "app.curexal.space", "owerri", "workspace")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Items, 2)
	assert.Equal(t, "/owerri/reception", res.Items[0].Path)
	assert.Equal(t, "/owerri/billing", res.Items[1].Path)

	// Verify that receptionist does NOT have consultation permissions
	hasConsultationPerm := false
	for _, p := range receptionPrincipal.Permissions {
		if p == "hms.consultation.create" {
			hasConsultationPerm = true
			break
		}
	}
	assert.False(t, hasConsultationPerm, "Receptionist must not have clinical consultation permissions")
}

func TestNavigationService_BranchAdmin_ClinicWorkspaceLevel(t *testing.T) {
	ctx := context.Background()

	dashboardItem := platformDomain.NavigationItem{
		ID:           "nav_wsp_dashboard",
		ContextScope: "workspace",
		ModuleCode:   stringPtr("dashboard"),
		Title:        "Workspace Overview",
		Path:         "/:branch/dashboard",
		Order:        1,
		Status:       platformDomain.NavigationStatusActive,
		IsVisible:    true,
		IsActive:     true,
	}
	receptionItem := platformDomain.NavigationItem{
		ID:                 "nav_wsp_reception",
		ContextScope:       "workspace",
		ModuleCode:         stringPtr("reception"),
		Title:              "Patient Intake (MPI)",
		Path:               "/:branch/reception",
		Order:              2,
		RequiredPermission: stringPtr("workspace:patient:read"),
		RequiredCapability: stringPtr("core.patient"),
		Status:             platformDomain.NavigationStatusActive,
		IsVisible:          true,
		IsActive:           true,
	}
	careDeskItem := platformDomain.NavigationItem{
		ID:                 "nav_wsp_care_desk",
		ContextScope:       "workspace",
		ModuleCode:         stringPtr("care_desk"),
		Title:              "Care & Triage Desk",
		Path:               "/:branch/care-desk",
		Order:              3,
		RequiredPermission: stringPtr("workspace:patient:read"),
		RequiredCapability: stringPtr("clinical.basic"),
		Status:             platformDomain.NavigationStatusActive,
		IsVisible:          true,
		IsActive:           true,
	}
	clinicalItem := platformDomain.NavigationItem{
		ID:                 "nav_wsp_clinical",
		ContextScope:       "workspace",
		ModuleCode:         stringPtr("clinical"),
		Title:              "Clinical & EMR",
		Path:               "/:branch/clinical",
		Order:              4,
		RequiredPermission: stringPtr("workspace:clinical:read"),
		RequiredCapability: stringPtr("clinical.basic"),
		Status:             platformDomain.NavigationStatusActive,
		IsVisible:          true,
		IsActive:           true,
	}
	billingItem := platformDomain.NavigationItem{
		ID:                 "nav_wsp_billing",
		ContextScope:       "workspace",
		ModuleCode:         stringPtr("billing"),
		Title:              "Billing & Cashier POS",
		Path:               "/:branch/billing",
		Order:              5,
		RequiredPermission: stringPtr("workspace:billing:create"),
		RequiredCapability: stringPtr("core.billing"),
		Status:             platformDomain.NavigationStatusActive,
		IsVisible:          true,
		IsActive:           true,
	}

	mockNavRepo := new(MockNavigationRepository)
	mockNavRepo.On(
		"GetNavigationItemsByScope",
		ctx,
		"workspace",
		mock.Anything,
		mock.Anything,
		true, // isWorkspaceAdmin must be true for branch_admin in workspace scope
	).Return([]platformDomain.NavigationItem{dashboardItem, receptionItem, careDeskItem, clinicalItem, billingItem}, nil)

	branchAdminPrincipal := &auth.AuthenticatedPrincipal{
		UserID:         "usr_branch_admin_01",
		TenantID:       "branch_edl_01",
		OrganizationID: "org_everight_01",
		Role:           "branch_admin",
		Permissions: []string{
			"workspace:patient:read",
			"workspace:patient:create",
			"workspace:clinical:read",
			"workspace:billing:create",
			"workspace:settings:manage",
		},
	}

	navService := application.NewNavigationService(mockNavRepo, nil, nil, nil)
	res, err := navService.GetNavigation(ctx, branchAdminPrincipal, "everight.localhost", "edl-01", "workspace")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Items, 5)
	assert.Equal(t, "/edl-01/dashboard", res.Items[0].Path)
	assert.Equal(t, "/edl-01/reception", res.Items[1].Path)
	assert.Equal(t, "/edl-01/care-desk", res.Items[2].Path)
	assert.Equal(t, "/edl-01/clinical", res.Items[3].Path)
	assert.Equal(t, "/edl-01/billing", res.Items[4].Path)
}

func stringPtr(s string) *string {
	return &s
}
