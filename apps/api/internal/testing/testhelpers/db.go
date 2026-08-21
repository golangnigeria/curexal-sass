//go:build integration

// Package testhelpers provides shared test fixtures for integration tests.
// It spins up a real PostgreSQL container via testcontainers-go and runs all
// Goose migrations so that repository tests operate against the actual schema.
package testhelpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/golangnigeria/curexal/database/migrator"
	"github.com/golangnigeria/curexal/internal/kernel/database"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testDBUser     = "curexal_test"
	testDBPassword = "curexal_test_secret"
	testDBName     = "curexal_test"
)

// TestDB bundles the live database connection and a cleanup function.
type TestDB struct {
	DB      *database.Database
	DSN     string
	Cleanup func()
}

// NewTestDB starts a PostgreSQL container, runs all migrations, and returns
// a ready-to-use *database.Database. Cleanup() stops the container.
//
// This function is intended to be called once per TestMain (or once per
// test that cannot share state). It skips if Docker is unavailable.
func NewTestDB(t *testing.T) *TestDB {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     testDBUser,
			"POSTGRES_PASSWORD": testDBPassword,
			"POSTGRES_DB":       testDBName,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Skipf("Docker unavailable, skipping integration test: %v", err)
		return nil
	}

	host, err := container.Host(ctx)
	require.NoError(t, err)

	mappedPort, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		testDBUser, testDBPassword, host, mappedPort.Port(), testDBName)

	// Run all Goose platform migrations against the fresh container.
	log := zerolog.Nop()
	runner := migrator.NewRunner(&log)
	require.NoError(t, runner.RunPlatform(ctx, dsn), "platform migrations must succeed")

	// Build the pgxpool Database connection.
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			MaxOpenConns: 5,
			MaxIdleConns: 2,
		},
	}
	// Override DSN via the env variable path that config.DSN() checks.
	t.Setenv("CUREXAL_DB_DSN", dsn)

	db, err := database.New(cfg, &log, nil)
	require.NoError(t, err, "database.New must connect to the test container")

	cleanup := func() {
		if db.Pool != nil {
			db.Pool.Close()
		}
		_ = container.Terminate(context.Background())
	}

	return &TestDB{
		DB:      db,
		DSN:     dsn,
		Cleanup: cleanup,
	}
}
