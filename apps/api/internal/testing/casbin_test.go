package testing

import (
	"context"
	"testing"

	"github.com/golangnigeria/curexal/internal/shared/authz"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCasbinEngine_RolePermissions(t *testing.T) {
	engine, err := authz.NewCasbinEngine()
	require.NoError(t, err)
	require.NotNil(t, engine)

	service := authz.NewService(engine)
	ctx := context.Background()

	orgID := uuid.New()
	branchID := uuid.New()

	// 1. DOCTOR: Authorized for clinical consultation, rejected from organization governance
	doctorPrincipal := &authz.Principal{
		UserID:      uuid.New(),
		Roles:       []string{"doctor"},
		Permissions: []string{"hms.consultation.create", "hms.prescription.create", "core.patient.read"},
		Credentials: []string{"physician"},
		Context: authz.ScopeContext{
			Scope:          authz.ScopeWorkspace,
			OrganizationID: &orgID,
			BranchID:       &branchID,
		},
	}

	assert.True(t, service.Can(ctx, doctorPrincipal, "hms.consultation.create", doctorPrincipal.Context))
	assert.True(t, service.Can(ctx, doctorPrincipal, "core.patient.read", doctorPrincipal.Context))
	assert.False(t, service.Can(ctx, doctorPrincipal, "organization.governance", authz.ScopeContext{Scope: authz.ScopeOrganization}))

	// 2. RECEPTIONIST: Authorized for shared Core/MPI patient registration, rejected from clinical consultation
	receptionistPrincipal := &authz.Principal{
		UserID:      uuid.New(),
		Roles:       []string{"receptionist"},
		Permissions: []string{"core.patient.create", "core.patient.search", "core.patient.read", "hms.appointment.create"},
		Context: authz.ScopeContext{
			Scope:          authz.ScopeWorkspace,
			OrganizationID: &orgID,
			BranchID:       &branchID,
		},
	}

	assert.True(t, service.Can(ctx, receptionistPrincipal, "core.patient.create", receptionistPrincipal.Context), "Receptionist MUST be able to register patient via shared Core MPI")
	assert.True(t, service.Can(ctx, receptionistPrincipal, "core.patient.search", receptionistPrincipal.Context))
	assert.False(t, service.Can(ctx, receptionistPrincipal, "hms.consultation.create", receptionistPrincipal.Context), "Receptionist MUST NOT be able to create consultation")
	assert.False(t, service.Can(ctx, receptionistPrincipal, "hms.prescription.create", receptionistPrincipal.Context), "Receptionist MUST NOT be able to prescribe")

	// 3. ORGANIZATION OWNER: Authorized for governance, rejected from clinical encounter rooms
	ownerPrincipal := &authz.Principal{
		UserID:      uuid.New(),
		Roles:       []string{"owner"},
		Permissions: []string{"organization.governance", "organization.billing", "branch.operations.overview"},
		Context: authz.ScopeContext{
			Scope:          authz.ScopeOrganization,
			OrganizationID: &orgID,
		},
	}

	assert.True(t, service.Can(ctx, ownerPrincipal, "organization.governance", ownerPrincipal.Context))
	assert.False(t, service.Can(ctx, ownerPrincipal, "hms.consultation.create", authz.ScopeContext{Scope: authz.ScopeWorkspace}), "Owner without physician credentials cannot write SOAP notes")

	// 4. NURSE: Authorized for triage intake & vitals, rejected from consultation signing
	nursePrincipal := &authz.Principal{
		UserID:      uuid.New(),
		Roles:       []string{"nurse"},
		Permissions: []string{"hms.triage.create", "core.patient.read"},
		Context: authz.ScopeContext{
			Scope:          authz.ScopeWorkspace,
			OrganizationID: &orgID,
			BranchID:       &branchID,
		},
	}

	assert.True(t, service.Can(ctx, nursePrincipal, "hms.triage.create", nursePrincipal.Context))
	assert.False(t, service.Can(ctx, nursePrincipal, "hms.consultation.sign", nursePrincipal.Context), "Nurse cannot sign doctor clinical consultation")

	// 5. CASHIER: Authorized for POS payment collection, rejected from doctor consultation
	cashierPrincipal := &authz.Principal{
		UserID:      uuid.New(),
		Roles:       []string{"cashier"},
		Permissions: []string{"billing.payment.create", "billing.invoice.read", "core.patient.read"},
		Context: authz.ScopeContext{
			Scope:          authz.ScopeWorkspace,
			OrganizationID: &orgID,
			BranchID:       &branchID,
		},
	}

	assert.True(t, service.Can(ctx, cashierPrincipal, "billing.payment.create", cashierPrincipal.Context))
	assert.False(t, service.Can(ctx, cashierPrincipal, "hms.consultation.create", cashierPrincipal.Context))

	// 6. FAIL-CLOSED: User with empty roles is denied all permissions (no owner escalation)
	emptyPrincipal := &authz.Principal{
		UserID:      uuid.New(),
		Roles:       []string{},
		Permissions: []string{},
		Context: authz.ScopeContext{
			Scope:          authz.ScopeWorkspace,
			OrganizationID: &orgID,
			BranchID:       &branchID,
		},
	}
	assert.False(t, service.Can(ctx, emptyPrincipal, "organization.governance", authz.ScopeContext{Scope: authz.ScopeOrganization}), "Empty roles must not escalate to owner")
	assert.False(t, service.Can(ctx, emptyPrincipal, "hms.consultation.create", emptyPrincipal.Context))
	assert.False(t, service.Can(ctx, emptyPrincipal, "core.patient.create", emptyPrincipal.Context))

	// 7. CROSS-BRANCH ISOLATION: User in Branch A cannot access Branch B
	branchB_ID := uuid.New()
	assert.False(t, service.Can(ctx, doctorPrincipal, "organization.governance", authz.ScopeContext{
		Scope:          authz.ScopeOrganization,
		OrganizationID: &orgID,
		BranchID:       &branchB_ID,
	}))
}

