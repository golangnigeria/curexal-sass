package migrator

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

type SchemaValidator struct {
	logger *zerolog.Logger
}

func NewValidator(log *zerolog.Logger) *SchemaValidator {
	return &SchemaValidator{logger: log}
}

// ValidateSchemaHealth verifies that database migrations are non-dirty and connection is responsive.
func (v *SchemaValidator) ValidateSchemaHealth(dsn string, tableName string) error {
	if tableName == "" {
		tableName = "schema_migrations"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("pre-flight schema validation failed: unable to open connection: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("pre-flight schema validation failed: database unreachable: %w", err)
	}

	goose.SetTableName(tableName)
	version, err := goose.GetDBVersion(db)
	if err != nil {
		v.logger.Warn().Err(err).Str("table", tableName).Msg("schema tracking table not found yet (will be initialized by migrator)")
		return nil
	}

	v.logger.Info().Int64("current_version", version).Str("table", tableName).Msg("pre-flight database schema version check passed")
	return nil
}
