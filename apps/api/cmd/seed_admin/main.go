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
	for _, envPath := range []string{".env", "../.env", "../../.env"} {
		_ = godotenv.Overload(envPath)
	}

	dsn := os.Getenv("CUREXAL_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		log.Fatal("CUREXAL_DB_DSN or DATABASE_URL environment variable is not set in .env")
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
		orgRole      string
		orgSlug      string
		orgName      string
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
		{
			id:           "user_org_owner_curexal_clinic",
			name:         "Dr. Alexander Vance (Owner)",
			email:        "owner@curexal.space",
			platformRole: "",
			orgRole:      "owner",
			orgSlug:      "curexal-clinic",
			orgName:      "Curexal Premier Medical Center",
		},
		{
			id:           "user_org_owner_everight",
			name:         "Dr. Everett Right (Owner)",
			email:        "owner@everight.com",
			platformRole: "",
			orgRole:      "owner",
			orgSlug:      "everight",
			orgName:      "Everight Diagnostic & Speciality Hospital",
		},
	}

	password := "password"
	// Hash password using Argon2id
	hash, err := crypto.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v\n", err)
	}

	fmt.Println("Seeding platform and organization users...")

	for _, u := range usersToSeed {
		userUUID := uuid.NewMD5(uuid.NameSpaceDNS, []byte(u.id)).String()
		isPAdmin := u.platformRole == "super_admin"

		var exists bool
		err = conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM identity.users WHERE email = $1)`, u.email).Scan(&exists)
		if err != nil {
			log.Fatalf("Failed to check if user %s exists: %v\n", u.email, err)
		}

		if !exists {
			var pRole *string
			if u.platformRole != "" {
				pRole = &u.platformRole
			}
			err = conn.QueryRow(ctx, `
				INSERT INTO identity.users (id, name, email, email_verified, is_platform_admin, platform_role)
				VALUES ($1, $2, $3, TRUE, $4, $5)
				RETURNING id
			`, userUUID, u.name, u.email, isPAdmin, pRole).Scan(&userUUID)
			if err != nil {
				log.Fatalf("Failed to insert user %s: %v\n", u.email, err)
			}
			fmt.Printf("Created user record (%s - %s)\n", u.name, u.email)
		} else {
			var pRole *string
			if u.platformRole != "" {
				pRole = &u.platformRole
			}
			err = conn.QueryRow(ctx, `
				UPDATE identity.users SET is_platform_admin = $2, platform_role = $3, name = $4 WHERE email = $1 RETURNING id
			`, u.email, isPAdmin, pRole, u.name).Scan(&userUUID)
			if err != nil {
				log.Fatalf("Failed to update user %s: %v\n", u.email, err)
			}
			fmt.Printf("Ensured user record for %s\n", u.email)
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

		// If user belongs to an organization as Owner/Admin, seed organization and membership
		if u.orgSlug != "" && u.orgRole != "" {
			var orgID string
			errOrg := conn.QueryRow(ctx, `
				INSERT INTO organization.organizations (name, slug, status, plan, setup_state)
				VALUES ($1, $2, 'active', 'enterprise', 'VERIFIED')
				ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name, status = 'active', setup_state = 'VERIFIED'
				RETURNING id::text
			`, u.orgName, u.orgSlug).Scan(&orgID)
			if errOrg != nil {
				log.Fatalf("Failed to upsert organization %s: %v\n", u.orgSlug, errOrg)
			}

			// Ensure facility branch exists
			var branchID string
			_ = conn.QueryRow(ctx, `
				INSERT INTO organization.facility_branches (organization_id, name, code, slug, is_headquarters, status)
				VALUES ($1, $2, 'main', 'main', TRUE, 'ACTIVE')
				ON CONFLICT (organization_id, code) DO UPDATE SET is_headquarters = TRUE, status = 'ACTIVE'
				RETURNING id::text
			`, orgID, u.orgName+" (Main Branch)").Scan(&branchID)

			// Ensure organization membership exists
			_, errMem := conn.Exec(ctx, `
				INSERT INTO organization.organization_memberships (user_id, organization_id, role, role_title, is_active)
				VALUES ($1, $2, $3, $3, TRUE)
				ON CONFLICT (organization_id, user_id) DO UPDATE SET role = EXCLUDED.role, role_title = EXCLUDED.role_title, is_active = TRUE
			`, userUUID, orgID, u.orgRole)
			if errMem != nil {
				log.Fatalf("Failed to upsert organization membership for %s in %s: %v\n", u.email, u.orgSlug, errMem)
			}
			fmt.Printf("Configured %s as %s of %s (%s)\n", u.email, u.orgRole, u.orgName, u.orgSlug)
		}
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("All users configured. Default password: %s\n", password)
	fmt.Println("--------------------------------------------------")
}
