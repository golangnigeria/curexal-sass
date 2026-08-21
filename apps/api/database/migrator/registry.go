package migrator

import (
	"embed"
	"fmt"
)

type MigrationBundle struct {
	Name         string
	MigrationsFS embed.FS
	MigrationsDir string
	SeedersFS    embed.FS
	SeedersDir   string
	TableName    string
}

type MigrationRegistry struct {
	platformBundles []MigrationBundle
	tenantBundles   []MigrationBundle
}

func NewRegistry() *MigrationRegistry {
	return &MigrationRegistry{
		platformBundles: make([]MigrationBundle, 0),
		tenantBundles:   make([]MigrationBundle, 0),
	}
}

func (r *MigrationRegistry) RegisterPlatform(name string, migrationsFS embed.FS, migrationsDir string, seedersFS embed.FS, seedersDir string, tableName string) {
	if tableName == "" {
		tableName = "schema_migrations"
	}
	r.platformBundles = append(r.platformBundles, MigrationBundle{
		Name:          name,
		MigrationsFS:  migrationsFS,
		MigrationsDir: migrationsDir,
		SeedersFS:     seedersFS,
		SeedersDir:    seedersDir,
		TableName:     tableName,
	})
}

func (r *MigrationRegistry) RegisterTenant(name string, migrationsFS embed.FS, migrationsDir string, seedersFS embed.FS, seedersDir string, tableName string) {
	if tableName == "" {
		tableName = "tenant_schema_migrations"
	}
	r.tenantBundles = append(r.tenantBundles, MigrationBundle{
		Name:          name,
		MigrationsFS:  migrationsFS,
		MigrationsDir: migrationsDir,
		SeedersFS:     seedersFS,
		SeedersDir:    seedersDir,
		TableName:     tableName,
	})
}

func (r *MigrationRegistry) PlatformBundles() []MigrationBundle {
	return r.platformBundles
}

func (r *MigrationRegistry) TenantBundles() []MigrationBundle {
	return r.tenantBundles
}

func (r *MigrationRegistry) Validate() error {
	if len(r.platformBundles) == 0 {
		return fmt.Errorf("migration registry has zero platform bundles registered")
	}
	return nil
}
