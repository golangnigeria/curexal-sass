package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/identity/repository"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	crypto "github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/golangnigeria/curexal/internal/shared/job"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type InviteService struct {
	server       *server.Server
	authService  *AuthService
	userRepo     *repository.UserRepository
	inviteRepo   *repository.InviteRepository
	tenantLookup model.TenantLookup
}

func NewInviteService(s *server.Server, authService *AuthService, tenantLookup model.TenantLookup) *InviteService {
	return &InviteService{
		server:       s,
		authService:  authService,
		userRepo:     repository.NewUserRepository(s),
		inviteRepo:   repository.NewInviteRepository(s),
		tenantLookup: tenantLookup,
	}
}

// InviteMember creates an invitation for the given email to join the tenant with the specified role.
func (s *InviteService) InviteMember(ctx context.Context, tenantID uuid.UUID, email, roleName, inviterUserID, origin string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return errors.New("email is required")
	}

	// 1. Resolve role ID
	roleID, err := s.userRepo.GetRoleIDByName(ctx, roleName, tenantID.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("invalid role name specified")
		}
		return fmt.Errorf("failed to look up role: %w", err)
	}

	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		return fmt.Errorf("invalid role ID format: %w", err)
	}

	// 2. Check if user already has active membership in this tenant
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		exists, err := s.userRepo.CheckMembershipExists(ctx, existingUser.ID, tenantID.String())
		if err != nil {
			return fmt.Errorf("failed to check membership: %w", err)
		}
		if exists {
			return errors.New("user is already an active member of this workspace")
		}
	}

	// 3. Check for existing pending invitation
	existingInvite, err := s.inviteRepo.GetPendingInvitationByEmail(ctx, tenantID, email)
	if err == nil && existingInvite != nil {
		// Check if it's expired
		if time.Now().After(existingInvite.ExpiresAt) {
			// Clean up expired invite so we can re-issue
			_ = s.inviteRepo.DeleteExpiredInvitation(ctx, tenantID, email)
		} else {
			return errors.New("an invitation has already been sent to this email address")
		}
	}

	// 4. Generate secure invite token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("failed to generate invite token: %w", err)
	}
	inviteToken := "invite_" + hex.EncodeToString(tokenBytes)

	// 5. Create invitation record
	inv := &model.Invitation{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Email:     email,
		RoleID:    roleUUID,
		Token:     inviteToken,
		InvitedBy: inviterUserID,
		Status:    model.InvitationStatusPending,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.inviteRepo.CreateInvitation(ctx, inv); err != nil {
		return fmt.Errorf("failed to create invitation: %w", err)
	}

	// 6. Resolve inviter name and tenant details
	inviterName := "A team member"
	if inviter, err := s.userRepo.GetByID(ctx, inviterUserID); err == nil {
		inviterName = inviter.Name
	}

	tenantName := "your workspace"
	tenantSlug := ""
	if s.tenantLookup != nil {
		if t, err := s.tenantLookup.GetTenantByID(ctx, tenantID); err == nil && t != nil {
			tenantName = t.Name
			tenantSlug = t.Slug
		}
	}

	// 7. Resolve invite link (target port with the tenant's slug subdomain)
	workspaceBaseURL := fmt.Sprintf("http://%s:%s", s.server.Config.Server.Domain, s.server.Config.Server.WorkspacePort)
	if tenantSlug != "" {
		workspaceBaseURL = fmt.Sprintf("http://%s.%s:%s", tenantSlug, s.server.Config.Server.Domain, s.server.Config.Server.WorkspacePort)
	}
	inviteLink := fmt.Sprintf("%s/auth/accept-invite?token=%s", workspaceBaseURL, inviteToken)

	// 8. Enqueue invite email task
	task, err := job.NewInviteMemberEmailTask(email, inviterName, tenantName, inviteLink)
	if err == nil {
		if s.server.Job != nil && s.server.Job.Client != nil {
			if _, enqueueErr := s.server.Job.Client.Enqueue(task); enqueueErr != nil {
				s.server.Logger.Error().Err(enqueueErr).Msg("failed to enqueue invite member email task")
			} else {
				s.server.Logger.Info().Str("email", email).Str("tenant", tenantName).Msg("enqueued invite member email task")
			}
		} else {
			s.server.Logger.Warn().Str("email", email).Msg("job service not available, skipping invite email task")
		}
	} else {
		s.server.Logger.Error().Err(err).Msg("failed to create invite member email task")
	}

	// 9. Audit log
	s.logInviteEvent(ctx, &tenantID, &inviterUserID, "invite:sent",
		fmt.Sprintf(`{"email":"%s","tenant_id":"%s","role":"%s"}`, email, tenantID.String(), roleName), "info")

	return nil
}

// AcceptInvite processes an invitation token: creates/finds the user, creates the membership,
// marks the invitation accepted, and returns a session for the user.
func (s *InviteService) AcceptInvite(ctx context.Context, token, name, password, ip, ua string) (*model.User, string, string, error) {
	// 1. Fetch invitation by token
	inv, err := s.inviteRepo.GetInvitationByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", "", errors.New("invalid or expired invitation token")
		}
		return nil, "", "", fmt.Errorf("failed to fetch invitation: %w", err)
	}

	// 2. Validate status and expiry
	if inv.Status != model.InvitationStatusPending {
		return nil, "", "", fmt.Errorf("this invitation has already been %s", inv.Status)
	}
	if time.Now().After(inv.ExpiresAt) {
		return nil, "", "", errors.New("this invitation has expired")
	}

	// 3. Run user creation + membership in a transaction
	var user *model.User
	err = s.server.DB.RunInTx(ctx, func(txCtx context.Context) error {
		// Find or create user
		existingUser, findErr := s.userRepo.GetByEmail(txCtx, inv.Email)
		if findErr != nil {
			if !errors.Is(findErr, pgx.ErrNoRows) {
				return fmt.Errorf("failed to check user existence: %w", findErr)
			}

			// User does not exist — create new user
			if name == "" {
				return errors.New("name is required for new accounts")
			}
			if password == "" {
				return errors.New("password is required for new accounts")
			}

			userID := "usr_" + uuid.New().String()
			newUser := &model.User{
				ID:            userID,
				Name:          name,
				Email:         inv.Email,
				EmailVerified: true, // Invited users are pre-verified (email confirmed by invite delivery)
			}

			// Hash password
			passwordHash, hashErr := crypto.HashPassword(password)
			if hashErr != nil {
				return fmt.Errorf("failed to hash password: %w", hashErr)
			}

			if err := s.userRepo.CreateUser(txCtx, newUser, passwordHash); err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}

			user = newUser
		} else {
			// User already exists — use their existing account
			user = existingUser
		}

		// Create/reactivate membership
		if err := s.userRepo.AddOrReactivateMembership(txCtx, user.ID, inv.TenantID.String(), inv.RoleID.String()); err != nil {
			return fmt.Errorf("failed to create membership: %w", err)
		}

		// Mark invitation as accepted
		if err := s.inviteRepo.MarkInvitationAccepted(txCtx, inv.ID); err != nil {
			return fmt.Errorf("failed to mark invitation accepted: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, "", "", err
	}

	// 4. Create session + generate tokens
	sess, refreshToken, err := s.authService.CreateSession(ctx, user.ID, ip, ua, true)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create session: %w", err)
	}

	// Set the active tenant context to the invited tenant
	atc := &model.ActiveTenantContext{
		TenantID: inv.TenantID.String(),
	}
	if err := s.userRepo.UpdateSessionActiveTenant(ctx, sess.ID, atc); err != nil {
		s.server.Logger.Error().Err(err).Msg("failed to set active tenant for invite session")
	}

	accessToken, err := s.authService.GenerateToken(user.ID, sess.ID)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// 5. Audit log
	s.logInviteEvent(ctx, &inv.TenantID, &user.ID, "invite:accepted",
		fmt.Sprintf(`{"tenant_id":"%s","email":"%s"}`, inv.TenantID.String(), inv.Email), "info")

	return user, accessToken, refreshToken, nil
}

// logInviteEvent writes an invite-related audit event.
func (s *InviteService) logInviteEvent(ctx context.Context, tenantID *uuid.UUID, userID *string, action, details, severity string) {
	sevUpper := "INFO"
	switch severity {
	case "info":
		sevUpper = "INFO"
	case "warn":
		sevUpper = "WARNING"
	case "error":
		sevUpper = "HIGH"
	case "critical":
		sevUpper = "CRITICAL"
	default:
		sevUpper = "INFO"
	}

	categoryStr := "Administration"
	statusVal := "SUCCESS"
	if severity == "error" || severity == "critical" {
		statusVal = "FAILED"
	}

	s.server.Logger.Info().
		Str("action", action).
		Str("severity", sevUpper).
		Str("status", statusVal).
		Str("category", categoryStr).
		Str("details", details).
		Msg("LogInviteEvent")
}
