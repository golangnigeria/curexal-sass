package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/golangnigeria/curexal/internal/kernel/database"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	_ = godotenv.Load("../../.env", "../.env", ".env")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	ctx := context.Background()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	dsn := cfg.Database.DSN()
	fmt.Printf("Running database migrations on target: %s\n", dsn)
	if err := database.Migrate(ctx, &logger, cfg); err != nil {
		log.Fatalf("Failed to migrate database: %v\n", err)
	}

	fmt.Println("Database migrations run successfully!")
}
