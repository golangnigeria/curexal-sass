package domain_test

import (
	"testing"

	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOrganizationDomainEntity(t *testing.T) {
	orgID := uuid.New()
	org := domain.Organization{
		ID:   orgID,
		Name: "Apollo Healthcare Network",
		Slug: "apollo-health",
		Plan: "enterprise",
	}

	assert.Equal(t, orgID, org.ID)
	assert.Equal(t, "Apollo Healthcare Network", org.Name)
	assert.Equal(t, "apollo-health", org.Slug)
	assert.Equal(t, "enterprise", org.Plan)
}

func TestTenantDomainEntity(t *testing.T) {
	tenantID := uuid.New()
	ten := domain.Tenant{
		ID:             tenantID,
		Name:           "Apollo Diagnostics Central Lab",
		Slug:           "apollo_central",
		OrganizationID: uuid.New().String(),
		Currency:       "NGN",
	}

	assert.Equal(t, tenantID, ten.ID)
	assert.Equal(t, "apollo_central", ten.Slug)
	assert.Equal(t, "NGN", ten.Currency)
}
