package model_test

import (
	"testing"

	"github.com/golangnigeria/curexal/internal/modules/authorization/model"
	"github.com/stretchr/testify/assert"
)

func TestPermissionScopes(t *testing.T) {
	assert.Equal(t, model.PermissionScope("patient:view"), model.PermPatientView)
	assert.Equal(t, model.PermissionScope("laboratory:create_order"), model.PermLabOrderCreate)
	assert.Equal(t, model.PermissionScope("laboratory:authorize_result"), model.PermLabResultAuth)
	assert.Equal(t, model.PermissionScope("organization:manage"), model.PermOrgManage)
}

func TestEnforceRequest(t *testing.T) {
	req := model.EnforceRequest{
		Subject:  "usr_01HGB72A9X",
		Tenant:   "tenant_apollo",
		Resource: "laboratory:order",
		Action:   "create",
	}

	assert.Equal(t, "usr_01HGB72A9X", req.Subject)
	assert.Equal(t, "tenant_apollo", req.Tenant)
	assert.Equal(t, "laboratory:order", req.Resource)
	assert.Equal(t, "create", req.Action)
}
