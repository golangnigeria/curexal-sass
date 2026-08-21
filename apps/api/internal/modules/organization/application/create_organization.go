package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type CreateOrganizationCommand struct {
	Name               string
	Slug               string
	Plan               *string
	Address            *string
	City               *string
	State              *string
	LGA                *string
	Country            *string
	Phone              *string
	Email              *string
	RegistrationNumber *string
	LicenseNumber      *string
	TaxID              *string
	OwnerEmail         *string
	OwnerName          *string
}

type CreateOrganizationResult struct {
	Organization   *domain.Organization
	OwnerEmail     string
	InvitationSent bool
	SetupToken     string // populated for audit/internal dispatch
}

func (s *OrganizationApplicationService) CreateOrganization(ctx echo.Context, userID string, cmd *CreateOrganizationCommand) (*CreateOrganizationResult, error) {
	var (
		createdOrg *domain.Organization
		targetOwnerEmail string
		rawSetupToken string
	)

	err := s.server.DB.RunInTx(ctx.Request().Context(), func(txCtx context.Context) error {
		dbExec := s.server.DB.Conn(txCtx)

		cleanSlug := strings.ToLower(strings.TrimSpace(cmd.Slug))
		if cleanSlug == "" {
			cleanSlug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(cmd.Name), " ", "-"))
		}

		// 1. Validate slug uniqueness in organization.organizations
		var slugExists bool
		errSlug := dbExec.QueryRow(txCtx, `SELECT EXISTS(SELECT 1 FROM organization.organizations WHERE slug = $1)`, cleanSlug).Scan(&slugExists)
		if errSlug != nil {
			return fmt.Errorf("failed to check organization slug uniqueness: %w", errSlug)
		}
		if slugExists {
			return errs.NewBadRequestError("Organization slug/URL is already taken")
		}

		// 2. Resolve Owner Identity in identity.users
		cleanOwnerEmail := ""
		if cmd.OwnerEmail != nil && strings.TrimSpace(*cmd.OwnerEmail) != "" {
			cleanOwnerEmail = strings.ToLower(strings.TrimSpace(*cmd.OwnerEmail))
		} else if cmd.Email != nil && strings.TrimSpace(*cmd.Email) != "" {
			cleanOwnerEmail = strings.ToLower(strings.TrimSpace(*cmd.Email))
		} else {
			cleanOwnerEmail = fmt.Sprintf("owner@%s.curexal.com", cleanSlug)
		}
		targetOwnerEmail = cleanOwnerEmail

		ownerName := cmd.Name + " Administrator"
		if cmd.OwnerName != nil && strings.TrimSpace(*cmd.OwnerName) != "" {
			ownerName = strings.TrimSpace(*cmd.OwnerName)
		}

		var ownerUserID string
		var isNewUser bool
		errExists := dbExec.QueryRow(txCtx, `SELECT id::text FROM identity.users WHERE email = $1`, cleanOwnerEmail).Scan(&ownerUserID)
		if errExists != nil {
			if errors.Is(errExists, pgx.ErrNoRows) {
				ownerUserID = uuid.New().String()
				isNewUser = true
				stmtInsertUser := `
					INSERT INTO identity.users (id, name, email, email_verified, is_platform_admin)
					VALUES ($1, $2, $3, FALSE, FALSE)
				`
				_, txErr := dbExec.Exec(txCtx, stmtInsertUser, ownerUserID, ownerName, cleanOwnerEmail)
				if txErr != nil {
					return fmt.Errorf("failed to create owner user: %w", txErr)
				}

				// Insert credentials entry with NULL password_hash (Setup Required)
				stmtInsertAccount := `
					INSERT INTO identity.credentials (id, account_id, auth_provider, user_id, password_hash)
					VALUES ($1, $2, 'credential', $3, NULL)
				`
				_, txErr = dbExec.Exec(txCtx, stmtInsertAccount, uuid.New().String(), cleanOwnerEmail, ownerUserID)
				if txErr != nil {
					return fmt.Errorf("failed to create owner credentials record: %w", txErr)
				}
			} else {
				return errExists
			}
		}

		// 3. Resolve & Validate Plan against subscription.plans
		planToUse := "smart"
		if cmd.Plan != nil && strings.TrimSpace(*cmd.Plan) != "" {
			planToUse = strings.ToLower(strings.TrimSpace(*cmd.Plan))
		}

		var planID string
		errPlan := dbExec.QueryRow(txCtx, `SELECT id::text FROM subscription.plans WHERE code = $1 LIMIT 1`, planToUse).Scan(&planID)
		if errPlan != nil {
			if errors.Is(errPlan, pgx.ErrNoRows) {
				return errs.NewBadRequestError(fmt.Sprintf("Invalid subscription plan code '%s'", planToUse))
			}
			return fmt.Errorf("failed to query subscription plan: %w", errPlan)
		}

		// 4. Create Organization record with full profile metadata
		orgID := uuid.New().String()
		country := "Nigeria"
		if cmd.Country != nil && strings.TrimSpace(*cmd.Country) != "" {
			country = strings.TrimSpace(*cmd.Country)
		}

		stmtInsertOrg := `
			INSERT INTO organization.organizations (
				id, name, slug, plan, status,
				address, city, state, lga, country, phone, email,
				registration_number, license_number, tax_id,
				setup_state, setup_step, version
			)
			VALUES (
				$1, $2, $3, $4, 'pending_verification',
				$5, $6, $7, $8, $9, $10, $11,
				$12, $13, $14,
				'PENDING_REGISTRATION', 1, 1
			)
			RETURNING id, name, slug, plan, status, logo_url, custom_domain,
			          registration_number, license_number, tax_id, email, phone, address, city, state, lga, country,
			          setup_state, setup_step, completed_at, settings, version, created_at, updated_at, updated_by
		`
		rows, txErr := dbExec.Query(txCtx, stmtInsertOrg,
			orgID, cmd.Name, cleanSlug, planToUse,
			cmd.Address, cmd.City, cmd.State, cmd.LGA, country, cmd.Phone, cmd.Email,
			cmd.RegistrationNumber, cmd.LicenseNumber, cmd.TaxID,
		)
		if txErr != nil {
			return fmt.Errorf("failed to create organization: %w", txErr)
		}

		org, txErr := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[domain.Organization])
		if txErr != nil {
			return fmt.Errorf("failed to collect organization row: %w", txErr)
		}
		createdOrg = &org

		// 5. Generate 6-Character Verification Setup Code for New Owner Identity
		if isNewUser {
			// Generate 6-character uppercase alphanumeric code
			verificationCode, codeErr := crypto.GenerateAlphanumericCode(6)
			if codeErr != nil {
				return fmt.Errorf("failed to generate verification code: %w", codeErr)
			}
			rawSetupToken = verificationCode

			expiresAt := time.Now().Add(72 * time.Hour) // 72 hours validity

			// Insert into identity.verification_tokens
			stmtVerifyToken := `
				INSERT INTO identity.verification_tokens (id, user_id, token, token_type, metadata, expires_at)
				VALUES ($1, $2, $3, 'OWNER_INVITATION', jsonb_build_object('organization_id', $4::text, 'email', $5::text), $6)
			`
			_, txErr = dbExec.Exec(txCtx, stmtVerifyToken, uuid.New().String(), ownerUserID, verificationCode, createdOrg.ID.String(), cleanOwnerEmail, expiresAt)
			if txErr != nil {
				return fmt.Errorf("failed to store verification code: %w", txErr)
			}

			// Redundantly hash into password_setup_tokens for backwards-compatibility
			digest := sha256.Sum256([]byte(verificationCode))
			tokenHash := hex.EncodeToString(digest[:])
			stmtToken := `
				INSERT INTO identity.password_setup_tokens (id, user_id, organization_id, token_hash, token_type, expires_at)
				VALUES ($1, $2, $3, $4, 'OWNER_INVITATION', $5)
			`
			_, txErr = dbExec.Exec(txCtx, stmtToken, uuid.New().String(), ownerUserID, createdOrg.ID.String(), tokenHash, expiresAt)
			if txErr != nil {
				return fmt.Errorf("failed to store password setup token: %w", txErr)
			}
		}

		// 6. Create Active Subscription Record
		stmtSub := `
			INSERT INTO subscription.subscriptions (id, organization_id, plan_id, plan, status)
			VALUES ($1, $2, $3, $4, 'active')
		`
		_, txErr = dbExec.Exec(txCtx, stmtSub, uuid.New().String(), createdOrg.ID.String(), planID, planToUse)
		if txErr != nil {
			return fmt.Errorf("failed to create subscription: %w", txErr)
		}

		// 7. Create Owner Membership in organization.organization_memberships
		membershipID := uuid.New().String()
		stmtInsertOrgMembership := `
			INSERT INTO organization.organization_memberships (id, organization_id, user_id, role, role_title, is_active)
			VALUES ($1, $2, $3, 'owner', 'owner', TRUE)
			ON CONFLICT (organization_id, user_id) DO UPDATE SET role = 'owner', role_title = 'owner', is_active = TRUE
		`
		_, txErr = dbExec.Exec(txCtx, stmtInsertOrgMembership, membershipID, createdOrg.ID.String(), ownerUserID)
		if txErr != nil {
			return fmt.Errorf("failed to create owner membership: %w", txErr)
		}

		// 8. Create Primary Workspace in workspace.workspaces
		workspaceID := uuid.New().String()
		stmtInsertWorkspace := `
			INSERT INTO workspace.workspaces (id, organization_id, name, slug, facility_type)
			VALUES ($1, $2, $3, $4, 'laboratory')
			ON CONFLICT (slug) DO NOTHING
		`
		_, txErr = dbExec.Exec(txCtx, stmtInsertWorkspace, workspaceID, createdOrg.ID.String(), cmd.Name+" Main Branch", cleanSlug)
		if txErr != nil {
			return fmt.Errorf("failed to create primary workspace: %w", txErr)
		}

		// 9. Create Workspace Owner Membership in workspace.workspace_memberships
		wsMembershipID := uuid.New().String()
		stmtInsertWsMembership := `
			INSERT INTO workspace.workspace_memberships (id, workspace_id, user_id, role, role_title, is_active)
			VALUES ($1, $2, $3, 'owner', 'owner', TRUE)
			ON CONFLICT (workspace_id, user_id) DO UPDATE SET role = 'owner', role_title = 'owner', is_active = TRUE
		`
		_, txErr = dbExec.Exec(txCtx, stmtInsertWsMembership, wsMembershipID, workspaceID, ownerUserID)
		if txErr != nil {
			return fmt.Errorf("failed to create workspace membership: %w", txErr)
		}

		// 10. Write Audit Event into audit.audit_events
		var actorUUID string
		if userID != "" {
			if _, parseErr := uuid.Parse(userID); parseErr == nil {
				var userExists bool
				errUser := dbExec.QueryRow(txCtx, `SELECT EXISTS(SELECT 1 FROM identity.users WHERE id = $1::uuid)`, userID).Scan(&userExists)
				if errUser == nil && userExists {
					actorUUID = userID
				}
			}
		}
		if actorUUID == "" {
			actorUUID = ownerUserID
		}

		stmtAudit := `
			INSERT INTO audit.audit_events (id, action, resource_type, resource_id, actor_id, organization_id, workspace_id, payload)
			VALUES ($1::uuid, 'ORGANIZATION_REGISTERED', 'ORGANIZATION', $2, $3::uuid, $4::uuid, $5::uuid, jsonb_build_object('name', $6::text, 'slug', $7::text, 'status', 'pending_verification', 'owner_email', $8::text))
		`
		_, txErr = dbExec.Exec(txCtx, stmtAudit,
			uuid.New().String(),
			createdOrg.ID.String(),
			actorUUID,
			createdOrg.ID.String(),
			workspaceID,
			cmd.Name,
			cleanSlug,
			cleanOwnerEmail,
		)
		if txErr != nil {
			return fmt.Errorf("failed to write audit event: %w", txErr)
		}

		return nil
	})

	if err != nil {
		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		loweredMsg := strings.ToLower(err.Error())
		if strings.Contains(loweredMsg, "slug") || strings.Contains(loweredMsg, "already taken") {
			return nil, errs.NewBadRequestError("Organization slug/URL is already taken")
		}
		s.server.Logger.Error().Err(err).Str("org_name", cmd.Name).Str("slug", cmd.Slug).Msg("failed to create organization")
		return nil, errs.NewInternalServerError()
	}

	if rawSetupToken != "" {
		s.server.Logger.Info().
			Str("email", targetOwnerEmail).
			Str("org", createdOrg.Name).
			Str("verification_code", rawSetupToken).
			Msg("owner 6-character verification code generated")

		// Dispatch real email via Mailer
		if s.server.Mailer != nil {
			ownerDisplayName := createdOrg.Name + " Administrator"
			if cmd.OwnerName != nil && strings.TrimSpace(*cmd.OwnerName) != "" {
				ownerDisplayName = strings.TrimSpace(*cmd.OwnerName)
			}
			if mailErr := s.server.Mailer.SendOwnerInvitationEmail(ctx.Request().Context(), targetOwnerEmail, ownerDisplayName, createdOrg.Name, rawSetupToken); mailErr != nil {
				s.server.Logger.Error().Err(mailErr).Str("email", targetOwnerEmail).Msg("failed to send owner invitation email via Mailer")
			}
		}
	}

	return &CreateOrganizationResult{
		Organization:   createdOrg,
		OwnerEmail:     targetOwnerEmail,
		InvitationSent: true,
		SetupToken:     rawSetupToken,
	}, nil
}

type TransferOwnershipCommand struct {
	OrganizationID uuid.UUID
	NewOwnerEmail  string
	NewOwnerName   string
	ActorUserID    string
	Notes          *string
}

func (s *OrganizationApplicationService) TransferOwnership(ctx echo.Context, cmd *TransferOwnershipCommand) error {
	var (
		targetCode string
		orgName    string
	)
	cleanEmail := strings.ToLower(strings.TrimSpace(cmd.NewOwnerEmail))

	err := s.server.DB.RunInTx(ctx.Request().Context(), func(txCtx context.Context) error {
		dbExec := s.server.DB.Conn(txCtx)

		// 1. Verify Organization exists
		errOrg := dbExec.QueryRow(txCtx, `SELECT name FROM organization.organizations WHERE id = $1`, cmd.OrganizationID).Scan(&orgName)
		if errOrg != nil {
			if errors.Is(errOrg, pgx.ErrNoRows) {
				return errs.NewNotFoundError("Organization not found")
			}
			return fmt.Errorf("failed to verify organization: %w", errOrg)
		}

		// 2. Resolve or Provision new owner identity
		var newOwnerUserID string
		errUser := dbExec.QueryRow(txCtx, `SELECT id::text FROM identity.users WHERE email = $1`, cleanEmail).Scan(&newOwnerUserID)
		if errUser != nil {
			if errors.Is(errUser, pgx.ErrNoRows) {
				newOwnerUserID = uuid.New().String()
				stmtInsertUser := `
					INSERT INTO identity.users (id, name, email, email_verified, is_platform_admin)
					VALUES ($1, $2, $3, FALSE, FALSE)
				`
				_, txErr := dbExec.Exec(txCtx, stmtInsertUser, newOwnerUserID, cmd.NewOwnerName, cleanEmail)
				if txErr != nil {
					return fmt.Errorf("failed to create new owner user: %w", txErr)
				}

				// Insert credentials entry with NULL password_hash
				stmtInsertAccount := `
					INSERT INTO identity.credentials (id, account_id, auth_provider, user_id, password_hash)
					VALUES ($1, $2, 'credential', $3, NULL)
				`
				_, _ = dbExec.Exec(txCtx, stmtInsertAccount, uuid.New().String(), cleanEmail, newOwnerUserID)

				// Generate 6-character verification code
				code, _ := crypto.GenerateAlphanumericCode(6)
				targetCode = code
				expiresAt := time.Now().Add(72 * time.Hour)

				// Store in verification_tokens
				_, _ = dbExec.Exec(txCtx, `
					INSERT INTO identity.verification_tokens (id, user_id, token, token_type, metadata, expires_at)
					VALUES ($1, $2, $3, 'OWNER_INVITATION', jsonb_build_object('organization_id', $4::text, 'email', $5::text), $6)
				`, uuid.New().String(), newOwnerUserID, code, cmd.OrganizationID.String(), cleanEmail, expiresAt)

				digest := sha256.Sum256([]byte(code))
				tokenHash := hex.EncodeToString(digest[:])
				_, _ = dbExec.Exec(txCtx, `
					INSERT INTO identity.password_setup_tokens (id, user_id, organization_id, token_hash, token_type, expires_at)
					VALUES ($1, $2, $3, $4, 'OWNER_INVITATION', $5)
				`, uuid.New().String(), newOwnerUserID, cmd.OrganizationID.String(), tokenHash, expiresAt)
			} else {
				return errUser
			}
		}

		// 3. Find current owner(s)
		var previousOwnerID string
		_ = dbExec.QueryRow(txCtx, `
			SELECT user_id::text 
			FROM organization.organization_memberships 
			WHERE organization_id = $1 AND role = 'owner' AND is_active = TRUE 
			LIMIT 1
		`, cmd.OrganizationID).Scan(&previousOwnerID)

		// 4. Demote former owner to admin
		if previousOwnerID != "" && previousOwnerID != newOwnerUserID {
			_, txErr := dbExec.Exec(txCtx, `
				UPDATE organization.organization_memberships 
				SET role = 'admin', role_title = 'Administrator' 
				WHERE organization_id = $1 AND user_id = $2
			`, cmd.OrganizationID, previousOwnerID)
			if txErr != nil {
				return fmt.Errorf("failed to demote previous owner: %w", txErr)
			}
		}

		// 5. Grant or promote new owner membership
		stmtPromote := `
			INSERT INTO organization.organization_memberships (id, organization_id, user_id, role, role_title, is_active)
			VALUES ($1, $2, $3, 'owner', 'owner', TRUE)
			ON CONFLICT (organization_id, user_id) DO UPDATE 
			SET role = 'owner', role_title = 'owner', is_active = TRUE
		`
		_, txErr := dbExec.Exec(txCtx, stmtPromote, uuid.New().String(), cmd.OrganizationID, newOwnerUserID)
		if txErr != nil {
			return fmt.Errorf("failed to promote new owner in organization: %w", txErr)
		}

		// 6. Update workspace memberships across all workspaces of the organization
		wsRows, errWs := dbExec.Query(txCtx, `SELECT id FROM workspace.workspaces WHERE organization_id = $1`, cmd.OrganizationID)
		if errWs == nil {
			var wsIDs []uuid.UUID
			for wsRows.Next() {
				var wID uuid.UUID
				if errScan := wsRows.Scan(&wID); errScan == nil {
					wsIDs = append(wsIDs, wID)
				}
			}
			wsRows.Close()

			for _, wID := range wsIDs {
				_, _ = dbExec.Exec(txCtx, `
					INSERT INTO workspace.workspace_memberships (id, workspace_id, user_id, role, role_title, is_active)
					VALUES ($1, $2, $3, 'owner', 'owner', TRUE)
					ON CONFLICT (workspace_id, user_id) DO UPDATE 
					SET role = 'owner', role_title = 'owner', is_active = TRUE
				`, uuid.New().String(), wID, newOwnerUserID)
			}
		}

		// 7. Write Audit Event
		actorID := cmd.ActorUserID
		if actorID == "" {
			actorID = newOwnerUserID
		}
		notes := "Ownership transferred"
		if cmd.Notes != nil && *cmd.Notes != "" {
			notes = *cmd.Notes
		}

		stmtAudit := `
			INSERT INTO audit.audit_events (id, action, resource_type, resource_id, actor_id, organization_id, payload)
			VALUES ($1::uuid, 'ORGANIZATION_OWNERSHIP_TRANSFERRED', 'ORGANIZATION', $2, $3::uuid, $4::uuid, jsonb_build_object('previous_owner_id', $5::text, 'new_owner_id', $6::text, 'new_owner_email', $7::text, 'notes', $8::text))
		`
		_, _ = dbExec.Exec(txCtx, stmtAudit,
			uuid.New().String(),
			cmd.OrganizationID.String(),
			actorID,
			cmd.OrganizationID.String(),
			previousOwnerID,
			newOwnerUserID,
			cleanEmail,
			notes,
		)

		return nil
	})

	if err != nil {
		return err
	}

	if targetCode != "" {
		s.server.Logger.Info().
			Str("email", cleanEmail).
			Str("verification_code", targetCode).
			Msg("new owner verification code generated")

		if s.server.Mailer != nil {
			_ = s.server.Mailer.SendOwnerInvitationEmail(ctx.Request().Context(), cleanEmail, cmd.NewOwnerName, orgName, targetCode)
		}
	}

	return nil
}

// ResendOwnerInvite resends the owner setup verification code for an organization.
func (s *OrganizationApplicationService) ResendOwnerInvite(ctx echo.Context, orgID uuid.UUID) (string, error) {
	db := s.server.DB.Conn(ctx.Request().Context())

	// Find the current owner user and org name
	var (
		ownerEmail string
		ownerName  string
		ownerID    string
		orgName    string
	)
	err := db.QueryRow(ctx.Request().Context(), `
		SELECT u.id::text, u.name, u.email, o.name
		FROM organization.organization_memberships m
		JOIN identity.users u ON u.id = m.user_id
		JOIN organization.organizations o ON o.id = m.organization_id
		WHERE m.organization_id = $1 AND m.role = 'owner' AND m.is_active = TRUE
		LIMIT 1
	`, orgID).Scan(&ownerID, &ownerName, &ownerEmail, &orgName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errs.NewNotFoundError("Active owner membership not found for this organization")
		}
		return "", fmt.Errorf("failed to find owner: %w", err)
	}

	// Generate 6-character code
	code, err := crypto.GenerateAlphanumericCode(6)
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}

	expiresAt := time.Now().Add(72 * time.Hour)

	// Save to identity.verification_tokens
	_, err = db.Exec(ctx.Request().Context(), `
		INSERT INTO identity.verification_tokens (id, user_id, token, token_type, metadata, expires_at)
		VALUES ($1, $2, $3, 'OWNER_INVITATION', jsonb_build_object('organization_id', $4::text, 'email', $5::text), $6)
	`, uuid.New().String(), ownerID, code, orgID.String(), ownerEmail, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to store verification token: %w", err)
	}

	// Redundantly store in password_setup_tokens
	digest := sha256.Sum256([]byte(code))
	tokenHash := hex.EncodeToString(digest[:])
	_, _ = db.Exec(ctx.Request().Context(), `
		INSERT INTO identity.password_setup_tokens (id, user_id, organization_id, token_hash, token_type, expires_at)
		VALUES ($1, $2, $3, $4, 'OWNER_INVITATION', $5)
	`, uuid.New().String(), ownerID, orgID.String(), tokenHash, expiresAt)

	s.server.Logger.Info().
		Str("email", ownerEmail).
		Str("org_id", orgID.String()).
		Str("verification_code", code).
		Msg("Owner invitation code resent")

	if s.server.Mailer != nil {
		_ = s.server.Mailer.SendOwnerInvitationEmail(ctx.Request().Context(), ownerEmail, ownerName, orgName, code)
	}

	return code, nil
}
