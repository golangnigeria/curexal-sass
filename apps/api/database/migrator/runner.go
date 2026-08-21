package migrator

import (
	"context"
	"fmt"

	"github.com/golangnigeria/curexal/database/platform"
	"github.com/golangnigeria/curexal/database/tenant"
	"github.com/rs/zerolog"
)

type Runner struct {
	logger    *zerolog.Logger
	registry  *MigrationRegistry
	goose     *GooseRunner
	validator *SchemaValidator
}

func NewRunner(log *zerolog.Logger) *Runner {
	reg := NewRegistry()

	// Register Core Platform Module
	reg.RegisterPlatform("platform", platform.MigrationsFS, "migrations", platform.SeedersFS, "seeders", "schema_migrations")

	// Register Core Tenant Module
	reg.RegisterTenant("tenant", tenant.MigrationsFS, "migrations", tenant.SeedersFS, "seeders", "tenant_schema_migrations")

	return &Runner{
		logger:    log,
		registry:  reg,
		goose:     NewGooseRunner(log),
		validator: NewValidator(log),
	}
}

// RunPlatform executes all registered Platform schema migrations and seeders sequentially.
func (r *Runner) RunPlatform(ctx context.Context, dsn string) error {
	r.logger.Info().Msg("starting platform database migration pipeline")

	if err := r.registry.Validate(); err != nil {
		return fmt.Errorf("migration registry validation failed: %w", err)
	}

	for _, bundle := range r.registry.PlatformBundles() {
		r.logger.Info().Str("bundle", bundle.Name).Msg("executing platform migration bundle")

		// 1. Pre-flight health check
		if err := r.validator.ValidateSchemaHealth(dsn, bundle.TableName); err != nil {
			r.logger.Warn().Err(err).Str("bundle", bundle.Name).Msg("pre-flight schema check notice")
		}

		// 2. Run Goose migrations UP
		if err := r.goose.RunUp(dsn, bundle.MigrationsFS, bundle.MigrationsDir, bundle.TableName); err != nil {
			return fmt.Errorf("failed to run platform migration bundle %s: %w", bundle.Name, err)
		}

		// 3. Run Seeders
		if err := r.goose.RunSeeders(dsn, bundle.SeedersFS, bundle.SeedersDir); err != nil {
			r.logger.Warn().Err(err).Str("bundle", bundle.Name).Msg("seeder notice")
		}
	}

	r.logger.Info().Msg("platform database migration pipeline completed successfully")
	return nil
}

// RunTenant executes all registered Tenant schema migrations and seeders sequentially.
func (r *Runner) RunTenant(ctx context.Context, dsn string) error {
	r.logger.Info().Msg("starting tenant database migration pipeline")

	for _, bundle := range r.registry.TenantBundles() {
		r.logger.Info().Str("bundle", bundle.Name).Msg("executing tenant migration bundle")

		if err := r.goose.RunUp(dsn, bundle.MigrationsFS, bundle.MigrationsDir, bundle.TableName); err != nil {
			return fmt.Errorf("failed to run tenant migration bundle %s: %w", bundle.Name, err)
		}

		if err := r.goose.RunSeeders(dsn, bundle.SeedersFS, bundle.SeedersDir); err != nil {
			r.logger.Warn().Err(err).Str("bundle", bundle.Name).Msg("tenant seeder notice")
		}
	}

	r.logger.Info().Msg("tenant database migration pipeline completed successfully")
	return nil
}
