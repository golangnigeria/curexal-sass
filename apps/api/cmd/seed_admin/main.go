package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env configuration
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

	type seedUser struct {
		id           string
		name         string
		email        string
		platformRole string
	}

	usersToSeed := []seedUser{
		{
			id:           "user_owner_seed_id_000",
			name:         "Platform Owner",
			email:        "admin@curexal.com",
			platformRole: "super_admin",
		},
		{
			id:           "user_3DfYn1lrFG48OsZFqQX8zVEu154",
			name:         "Super Admin",
			email:        "superadmin@curexal.internal",
			platformRole: "super_admin",
		},
		{
			id:           "user_support_agent_seed_id_123",
			name:         "Support Agent",
			email:        "support@curexal.internal",
			platformRole: "super_support_agent",
		},
		{
			id:           "user_sales_staff_seed_id_456",
			name:         "Sales Staff",
			email:        "sales@curexal.internal",
			platformRole: "super_sales_staff",
		},
		{
			id:           "user_compliance_officer_seed_id",
			name:         "Compliance Officer",
			email:        "compliance@curexal.internal",
			platformRole: "super_compliance_officer",
		},
	}

	password := "password"
	// Hash password using Argon2id
	hash, err := crypto.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v\n", err)
	}

	fmt.Println("Seeding platform users...")

	for _, u := range usersToSeed {
		// Parse string ID to UUID format or generate predictable UUID
		userUUID := uuid.NewMD5(uuid.NameSpaceDNS, []byte(u.id)).String()

		// Ensure user exists and is configured with the correct platform role
		var exists bool
		err = conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM identity.users WHERE email = $1)`, u.email).Scan(&exists)
		if err != nil {
			log.Fatalf("Failed to check if user %s exists: %v\n", u.email, err)
		}

		isPAdmin := u.platformRole == "super_admin"

		if !exists {
			err = conn.QueryRow(ctx, `
				INSERT INTO identity.users (id, name, email, email_verified, is_platform_admin, platform_role)
				VALUES ($1, $2, $3, TRUE, $4, $5)
				RETURNING id
			`, userUUID, u.name, u.email, isPAdmin, u.platformRole).Scan(&userUUID)
			if err != nil {
				log.Fatalf("Failed to insert user %s: %v\n", u.email, err)
			}
			fmt.Printf("Created %s user record (%s)\n", u.name, u.email)
		} else {
			err = conn.QueryRow(ctx, `
				UPDATE identity.users SET is_platform_admin = $2, platform_role = $3, name = $4 WHERE email = $1 RETURNING id
			`, u.email, isPAdmin, u.platformRole, u.name).Scan(&userUUID)
			if err != nil {
				log.Fatalf("Failed to update user %s to platform role %s: %v\n", u.email, u.platformRole, err)
			}
			fmt.Printf("Ensured user record role is set to %s for %s\n", u.platformRole, u.email)
		}

		// Insert or update credential record in identity.credentials table
		var accountExists bool
		err = conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM identity.credentials WHERE user_id = $1 AND auth_provider = 'credential')`, userUUID).Scan(&accountExists)
		if err != nil {
			log.Fatalf("Failed to check if account exists for user %s: %v\n", u.email, err)
		}

		cleanEmail := strings.ToLower(strings.TrimSpace(u.email))
		if !accountExists {
			accountID := uuid.New().String()
			_, err = conn.Exec(ctx, `
				INSERT INTO identity.credentials (id, account_id, auth_provider, user_id, password_hash)
				VALUES ($1, $2, 'credential', $3, $4)
			`, accountID, cleanEmail, userUUID, hash)
			if err != nil {
				log.Fatalf("Failed to insert account credential for user %s: %v\n", u.email, err)
			}
			fmt.Printf("Created local credential account record for %s\n", u.email)
		} else {
			_, err = conn.Exec(ctx, `
				UPDATE identity.credentials SET password_hash = $1, account_id = $3 WHERE user_id = $2 AND auth_provider = 'credential'
			`, hash, userUUID, cleanEmail)
			if err != nil {
				log.Fatalf("Failed to update account credential for user %s: %v\n", u.email, err)
			}
			fmt.Printf("Updated local credential account password for %s\n", u.email)
		}
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Successfully set platform users' password to: %s\n", password)
	fmt.Println("--------------------------------------------------")
}
