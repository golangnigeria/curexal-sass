package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/golangnigeria/curexal/internal/kernel/database"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("CUREXAL_DB_DSN")
	if dsn == "" {
		log.Fatal("CUREXAL_DB_DSN environment variable is not set in .env")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	fmt.Println("Querying all non-system database schemas to drop...")

	rows, err := conn.Query(ctx, `
		SELECT nspname 
		FROM pg_namespace 
		WHERE nspname NOT LIKE 'pg_%' 
		  AND nspname != 'information_schema'
	`)
	if err != nil {
		log.Fatalf("Failed to query database schemas: %v\n", err)
	}

	var schemas []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err == nil {
			schemas = append(schemas, s)
		}
	}
	rows.Close()

	for _, s := range schemas {
		dropQuery := fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", pgx.Identifier{s}.Sanitize())
		if _, err := conn.Exec(ctx, dropQuery); err != nil {
			log.Fatalf("Failed to drop schema %s: %v\n", s, err)
		}
	}

	// Recreate public schema
	if _, err := conn.Exec(ctx, "CREATE SCHEMA public;"); err != nil {
		log.Fatalf("Failed to recreate public schema: %v\n", err)
	}
	if _, err := conn.Exec(ctx, "GRANT ALL ON SCHEMA public TO public;"); err != nil {
		log.Printf("Warning: failed to grant on public schema: %v\n", err)
	}

	fmt.Println("All database schemas dropped and reset successfully.")

	// Load configuration and run migrations
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config for migration: %v\n", err)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	fmt.Println("Running up migrations...")
	if err := database.Migrate(ctx, &logger, cfg); err != nil {
		log.Fatalf("Failed to migrate database after wipe: %v\n", err)
	}

	fmt.Println("Database reset and migrations completed successfully!")
}
