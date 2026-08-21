package migrator

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

// GooseRunner encapsulates goose migration execution against a database connection.
type GooseRunner struct {
	logger *zerolog.Logger
}

func NewGooseRunner(log *zerolog.Logger) *GooseRunner {
	return &GooseRunner{logger: log}
}

// RunUp executes all pending Goose SQL migrations embedded inside an embed.FS directory.
func (r *GooseRunner) RunUp(dsn string, targetFS embed.FS, dir string, tableName string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection for goose: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database for goose migration: %w", err)
	}

	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if tableName != "" {
		goose.SetTableName(tableName)
	}

	goose.SetBaseFS(targetFS)

	r.logger.Info().Str("dir", dir).Str("table", tableName).Msg("executing goose up migrations")
	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("goose up failed in %s: %w", dir, err)
	}

	return nil
}

// RunSeeders executes non-goose SQL seeder files embedded inside an embed.FS directory.
func (r *GooseRunner) RunSeeders(dsn string, targetFS embed.FS, dir string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection for seeders: %w", err)
	}
	defer db.Close()

	subFS, err := fs.Sub(targetFS, dir)
	if err != nil {
		return nil // directory may not contain files
	}

	entries, err := fs.ReadDir(subFS, ".")
	if err != nil || len(entries) == 0 {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		content, errRead := fs.ReadFile(subFS, entry.Name())
		if errRead != nil {
			r.logger.Warn().Err(errRead).Str("file", entry.Name()).Msg("failed to read seeder file")
			continue
		}

		if len(content) == 0 {
			continue
		}

		r.logger.Info().Str("file", entry.Name()).Msg("running seeder script")
		_, errExec := db.Exec(string(content))
		if errExec != nil {
			r.logger.Warn().Err(errExec).Str("file", entry.Name()).Msg("seeder execution notice")
		}
	}

	return nil
}
