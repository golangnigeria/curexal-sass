package tenant

import "embed"

//go:embed migrations/*.sql
var MigrationsFS embed.FS

//go:embed seeders/*.sql
var SeedersFS embed.FS
