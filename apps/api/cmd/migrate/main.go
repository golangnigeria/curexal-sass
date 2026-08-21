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
	_ = godotenv.Load()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	ctx := context.Background()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	fmt.Println("Running database migrations...")
	if err := database.Migrate(ctx, &logger, cfg); err != nil {
		log.Fatalf("Failed to migrate database: %v\n", err)
	}

	fmt.Println("Database migrations run successfully!")
}
