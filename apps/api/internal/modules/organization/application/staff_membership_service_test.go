package application_test

import (
	"context"
	"testing"
	"time"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/organization/application"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStaffRepo struct {
	mock.Mock
}

func (m *MockStaffRepo) ListMembers(ctx context.Context, orgID uuid.UUID) ([]domain.StaffMemberDTO, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.StaffMemberDTO), args.Error(1)
}

func (m *MockStaffRepo) GetMemberByID(ctx context.Context, orgID, membershipID uuid.UUID) (*domain.StaffMemberDTO, error) {
	args := m.Called(ctx, orgID, membershipID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StaffMemberDTO), args.Error(1)
}

func (m *MockStaffRepo) CountActiveMembers(ctx context.Context, orgID uuid.UUID) (int, error) {
	args := m.Called(ctx, orgID)
	return args.Int(0), args.Error(1)
}

func (m *MockStaffRepo) CreateInvitation(ctx context.Context, invite *domain.StaffInvitation) (*domain.StaffInvitation, error) {
	args := m.Called(ctx, invite)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StaffInvitation), args.Error(1)
}

func (m *MockStaffRepo) ListInvitations(ctx context.Context, orgID uuid.UUID) ([]domain.StaffInvitation, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.StaffInvitation), args.Error(1)
}

func (m *MockStaffRepo) GetInvitationByHash(ctx context.Context, hash string) (*domain.StaffInvitation, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StaffInvitation), args.Error(1)
}

func (m *MockStaffRepo) RevokeInvitation(ctx context.Context, orgID, inviteID uuid.UUID) error {
	args := m.Called(ctx, orgID, inviteID)
	return args.Error(0)
}

func (m *MockStaffRepo) AcceptInvitation(ctx context.Context, inviteID uuid.UUID) error {
	args := m.Called(ctx, inviteID)
	return args.Error(0)
}

func (m *MockStaffRepo) CheckPendingInviteExists(ctx context.Context, orgID uuid.UUID, email string) (bool, error) {
	args := m.Called(ctx, orgID, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockStaffRepo) AssignBranch(ctx context.Context, membershipID, branchID uuid.UUID, actorID uuid.UUID) (*domain.MembershipBranch, error) {
	args := m.Called(ctx, membershipID, branchID, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MembershipBranch), args.Error(1)
}

func (m *MockStaffRepo) RemoveBranchAssignment(ctx context.Context, membershipID, branchID uuid.UUID) error {
	args := m.Called(ctx, membershipID, branchID)
	return args.Error(0)
}

func (m *MockStaffRepo) AssignDepartment(ctx context.Context, membershipID, branchID uuid.UUID, deptCode string, actorID uuid.UUID) (*domain.DepartmentalMembership, error) {
	args := m.Called(ctx, membershipID, branchID, deptCode, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DepartmentalMembership), args.Error(1)
}

func (m *MockStaffRepo) RemoveDepartmentAssignment(ctx context.Context, membershipID, branchID uuid.UUID, deptCode string) error {
	args := m.Called(ctx, membershipID, branchID, deptCode)
	return args.Error(0)
}

func (m *MockStaffRepo) UpdateMemberRole(ctx context.Context, orgID, membershipID uuid.UUID, role, roleTitle string, actorID uuid.UUID) (*domain.StaffMemberDTO, error) {
	args := m.Called(ctx, orgID, membershipID, role, roleTitle, actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StaffMemberDTO), args.Error(1)
}

func TestStaffMembershipService_CreateInvitation_SHA256HashToken_Success(t *testing.T) {
	mockStaffRepo := new(MockStaffRepo)
	mockOrgRepo := new(MockOrgRepo)
	mockBranchRepo := new(MockBranchRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewStaffMembershipService(mockStaffRepo, mockOrgRepo, mockBranchRepo, mockAuditRepo)

	orgID := uuid.New()
	actorID := uuid.New()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	payload := &domain.CreateStaffInvitationPayload{
		Email:     "new.clinician@curexal.health",
		Role:      "clinician",
		RoleTitle: "Medical Doctor",
	}

	// 1. Check pending invite -> false
	mockStaffRepo.On("CheckPendingInviteExists", mock.Anything, orgID, "new.clinician@curexal.health").Return(false, nil)

	// 2. Org plan -> smart (max 10 staff)
	mockOrgRepo.On("GetByID", mock.Anything, orgID).Return(&domain.Organization{ID: orgID, Plan: "smart"}, nil)
	mockStaffRepo.On("CountActiveMembers", mock.Anything, orgID).Return(2, nil)

	mockStaffRepo.On("CreateInvitation", mock.Anything, mock.Anything).Return(&domain.StaffInvitation{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		Email:           "new.clinician@curexal.health",
		Role:            "clinician",
		RoleTitle:       "Medical Doctor",
		InviteTokenHash: "hashed_token_64_chars",
		Status:          "PENDING",
		ExpiresAt:       time.Now().Add(7 * 24 * time.Hour),
		InvitedBy:       actorID,
	}, nil)

	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "STAFF_INVITED"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.CreateInvitation(context.Background(), principal, payload)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.RawToken)
	assert.Equal(t, 64, len(res.RawToken)) // 32 bytes hex = 64 chars
	assert.NotEqual(t, res.RawToken, res.Invitation.InviteTokenHash) // Plaintext token differs from stored hash

	mockStaffRepo.AssertExpectations(t)
	mockOrgRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestStaffMembershipService_CreateInvitation_MaxStaffExceeded_Fails(t *testing.T) {
	mockStaffRepo := new(MockStaffRepo)
	mockOrgRepo := new(MockOrgRepo)
	svc := application.NewStaffMembershipService(mockStaffRepo, mockOrgRepo, nil, nil)

	orgID := uuid.New()
	principal := &middleware.AuthenticatedPrincipal{
		UserID: uuid.New().String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	payload := &domain.CreateStaffInvitationPayload{
		Email: "extra.staff@curexal.health",
		Role:  "cashier",
	}

	mockStaffRepo.On("CheckPendingInviteExists", mock.Anything, orgID, "extra.staff@curexal.health").Return(false, nil)
	// Smart plan allows max 10 staff
	mockOrgRepo.On("GetByID", mock.Anything, orgID).Return(&domain.Organization{ID: orgID, Plan: "smart"}, nil)
	mockStaffRepo.On("CountActiveMembers", mock.Anything, orgID).Return(10, nil)

	res, err := svc.CreateInvitation(context.Background(), principal, payload)

	assert.ErrorIs(t, err, domain.ErrMaxStaffExceeded)
	assert.Nil(t, res)

	mockStaffRepo.AssertNotCalled(t, "CreateInvitation", mock.Anything, mock.Anything)
}

func TestStaffMembershipService_AssignDepartment_ValidAndInvalidCodes(t *testing.T) {
	mockStaffRepo := new(MockStaffRepo)
	mockOrgRepo := new(MockOrgRepo)
	mockBranchRepo := new(MockBranchRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewStaffMembershipService(mockStaffRepo, mockOrgRepo, mockBranchRepo, mockAuditRepo)

	orgID := uuid.New()
	actorID := uuid.New()
	memID := uuid.New()
	branchID := uuid.New()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	// 1. Invalid department code rejection
	resInv, errInv := svc.AssignDepartment(context.Background(), principal, memID, branchID, "invalid_dept_code")
	assert.ErrorIs(t, errInv, domain.ErrInvalidDepartmentCode)
	assert.Nil(t, resInv)

	// 2. Valid department code assignment ('laboratory')
	mockBranchRepo.On("GetBranchByID", mock.Anything, orgID, branchID).Return(&domain.FacilityBranch{ID: branchID, OrganizationID: orgID}, nil)
	mockStaffRepo.On("GetMemberByID", mock.Anything, orgID, memID).Return(&domain.StaffMemberDTO{MembershipID: memID, OrganizationID: orgID}, nil)
	mockStaffRepo.On("AssignDepartment", mock.Anything, memID, branchID, "laboratory", actorID).Return(&domain.DepartmentalMembership{
		ID:               uuid.New(),
		MembershipID:     memID,
		FacilityBranchID: branchID,
		DepartmentCode:   "laboratory",
	}, nil)

	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "STAFF_DEPARTMENT_ASSIGNED"
	})).Return(&auditDomain.AuditLog{}, nil)

	resValid, errValid := svc.AssignDepartment(context.Background(), principal, memID, branchID, "laboratory")
	assert.NoError(t, errValid)
	assert.NotNil(t, resValid)
	assert.Equal(t, "laboratory", resValid.DepartmentCode)

	mockStaffRepo.AssertExpectations(t)
	mockBranchRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestStaffMembershipService_UpdateMemberRole_Success(t *testing.T) {
	mockStaffRepo := new(MockStaffRepo)
	mockAuditRepo := new(MockAuditRepo)
	svc := application.NewStaffMembershipService(mockStaffRepo, nil, nil, mockAuditRepo)

	orgID := uuid.New()
	actorID := uuid.New()
	memID := uuid.New()

	principal := &middleware.AuthenticatedPrincipal{
		UserID: actorID.String(),
		Organization: platformAuth.OrganizationVector{
			ActiveOrganizationID: orgID.String(),
		},
	}

	updatedDTO := &domain.StaffMemberDTO{
		MembershipID:   memID,
		OrganizationID: orgID,
		Role:           "branch_admin",
		RoleTitle:      "Branch Administrator",
	}

	mockStaffRepo.On("UpdateMemberRole", mock.Anything, orgID, memID, "branch_admin", "Branch Administrator", actorID).Return(updatedDTO, nil)
	mockAuditRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *auditDomain.CreateAuditLogPayload) bool {
		return p.Action == "STAFF_ROLE_UPDATED"
	})).Return(&auditDomain.AuditLog{}, nil)

	res, err := svc.UpdateMemberRole(context.Background(), principal, memID, "branch_admin", "Branch Administrator")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "branch_admin", res.Role)
	assert.Equal(t, "Branch Administrator", res.RoleTitle)

	mockStaffRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}
