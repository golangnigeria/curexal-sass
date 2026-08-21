package bootstrap_test

import (
	"testing"

	"github.com/golangnigeria/curexal/internal/bootstrap"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModuleRegistryInitialization(t *testing.T) {
	e := echo.New()
	logger := zerolog.Nop()
	srv := &server.Server{
		Echo:   e,
		Logger: &logger,
	}

	reg := bootstrap.InitModules(srv)
	require.NotNil(t, reg, "ModuleRegistry should not be nil")
	assert.NotNil(t, reg.Identity, "Identity module should be registered")
	assert.NotNil(t, reg.Organization, "Organization module should be registered")
	assert.NotNil(t, reg.Authorization, "Authorization module should be registered")
	assert.NotNil(t, reg.Clinical, "Clinical module should be registered")
	assert.NotNil(t, reg.Catalogs, "Catalogs module should be registered")
}
