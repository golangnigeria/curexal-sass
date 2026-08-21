package database

import (
	"context"
	"fmt"

	"github.com/golangnigeria/curexal/database/migrator"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// DBExecutor represents common database operation interface for pgxpool and transactions.
type DBExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type txKey struct{}

// Database wraps the PostgreSQL connection pool.
type Database struct {
	Pool *pgxpool.Pool
}

// DB is an alias for Database for backwards compatibility.
type DB = Database

// New initializes and returns a new Database connection pool.
func New(cfg *config.Config, log *zerolog.Logger, _ interface{}) (*Database, error) {
	connStr := cfg.Database.DSN()
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgxpool config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.Database.MaxOpenConns)
	if cfg.Database.MaxIdleConns > 0 {
		poolConfig.MinConns = int32(cfg.Database.MaxIdleConns)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}

	return &Database{Pool: pool}, nil
}

// Conn returns active transaction if context contains one, otherwise connection pool.
func (d *Database) Conn(ctx context.Context) DBExecutor {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return d.Pool
}

// RunInTx executes a function within a database transaction.
func (d *Database) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := d.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	txCtx := context.WithValue(ctx, txKey{}, tx)
	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Tx is an alias for RunInTx.
func (d *Database) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.RunInTx(ctx, fn)
}

func (d *Database) ProvisionTenant(ctx context.Context, slug string) error {
	return nil
}

// Migrate executes versioned database migrations for Platform and Tenant schemas using Goose.
func Migrate(ctx context.Context, log *zerolog.Logger, cfg *config.Config) error {
	log.Info().Msg("running platform and tenant versioned database migrations via Goose")

	runner := migrator.NewRunner(log)
	dsn := cfg.Database.DSN()

	if err := runner.RunPlatform(ctx, dsn); err != nil {
		return fmt.Errorf("platform schema migration error: %w", err)
	}

	if err := runner.RunTenant(ctx, dsn); err != nil {
		return fmt.Errorf("tenant schema migration error: %w", err)
	}

	return nil
}

