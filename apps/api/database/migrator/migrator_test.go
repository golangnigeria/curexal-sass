package migrator_test

import (
	"testing"

	"github.com/golangnigeria/curexal/database/migrator"
	"github.com/golangnigeria/curexal/database/platform"
	"github.com/golangnigeria/curexal/database/tenant"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigratorRegistryInitialization(t *testing.T) {
	log := zerolog.Nop()
	runner := migrator.NewRunner(&log)
	require.NotNil(t, runner, "runner should not be nil")
}

func TestMigrationRegistryBundles(t *testing.T) {
	reg := migrator.NewRegistry()
	reg.RegisterPlatform("platform", platform.MigrationsFS, "migrations", platform.SeedersFS, "seeders", "schema_migrations")
	reg.RegisterTenant("tenant", tenant.MigrationsFS, "migrations", tenant.SeedersFS, "seeders", "tenant_schema_migrations")

	assert.Len(t, reg.PlatformBundles(), 1)
	assert.Len(t, reg.TenantBundles(), 1)

	err := reg.Validate()
	assert.NoError(t, err, "registered bundles should validate cleanly")
}
