//go:build integration

package testing

import (
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/database"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/rs/zerolog"
)

// CreateTestServer creates a server instance for testing
func CreateTestServer(logger *zerolog.Logger, db *TestDB) *server.Server {
	// Set up observability config with defaults if not present
	if db.Config.Observability == nil {
		db.Config.Observability = &config.ObservabilityConfig{
			ServiceName: "CUREXAL-test",
			Environment: "test",
			Logging: config.LoggingConfig{
				Level:              "info",
				Format:             "json",
				SlowQueryThreshold: int(100 * time.Millisecond),
			},
			NewRelic: config.NewRelicConfig{
				LicenseKey:                "",    // Empty for tests
				AppLogForwardingEnabled:   false, // Disabled for tests
				DistributedTracingEnabled: false, // Disabled for tests
				DebugLogging:              false, // Disabled for tests
			},
			HealthChecks: config.HealthChecksConfig{
				Enabled: false,
			},
		}
	}

	testServer := &server.Server{
		Logger: logger,
		DB: &database.Database{
			Pool: db.Pool,
		},
		Config: db.Config,
	}

	return testServer
}
