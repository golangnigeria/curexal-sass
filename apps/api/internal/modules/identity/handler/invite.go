package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/identity/repository"
	"github.com/golangnigeria/curexal/internal/modules/identity/service"
	orgRepo "github.com/golangnigeria/curexal/internal/modules/organization/infrastructure/postgres"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type InviteHandler struct {
	server        *server.Server
	inviteService *service.InviteService
	authService   *service.AuthService
	userRepo      *repository.UserRepository
	tenantRepo    *orgRepo.TenantRepository
}

func NewInviteHandler(s *server.Server, inviteService *service.InviteService, authService *service.AuthService) *InviteHandler {
	return &InviteHandler{
		server:        s,
		inviteService: inviteService,
		authService:   authService,
		userRepo:      repository.NewUserRepository(s),
		tenantRepo:    orgRepo.NewTenantRepository(s),
	}
}

// InviteMember handles POST /organizations/:id/invite
// Requires auth + tenant context + users:write permission.
func (h *InviteHandler) InviteMember(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	activeTenantID := middleware.GetActiveTenantID(c)
	if activeTenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Active tenant ID is missing from context")
	}

	tenantUUID, err := uuid.Parse(activeTenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid tenant ID format")
	}

	var payload model.InviteMemberPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	payload.Email = strings.ToLower(strings.TrimSpace(payload.Email))
	if payload.Email == "" || payload.Role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields (email, role)")
	}

	// Prevent self-invite
	if inviter, err := h.userRepo.GetByID(ctx, userID); err == nil && inviter.Email == payload.Email {
		return echo.NewHTTPError(http.StatusBadRequest, "You cannot invite yourself")
	}

	// Resolve origin for invite link
	origin := c.Request().Header.Get("Origin")
	if origin == "" {
		if ref := c.Request().Header.Get("Referer"); ref != "" {
			if idx := strings.Index(ref, "/"); idx > 8 { // skip protocol://
				origin = ref[:idx]
			} else {
				origin = ref
			}
		}
	}

	err = h.inviteService.InviteMember(ctx, tenantUUID, payload.Email, payload.Role, userID, origin)
	if err != nil {
		h.server.Logger.Error().Err(err).
			Str("email", payload.Email).
			Str("tenant_id", activeTenantID).
			Msg("failed to invite member")

		// Return user-facing errors as 400
		errMsg := err.Error()
		if strings.Contains(errMsg, "already") || strings.Contains(errMsg, "invalid role") {
			return echo.NewHTTPError(http.StatusBadRequest, errMsg)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send invitation")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Invitation sent to %s", payload.Email),
	})
}

// AcceptInvite handles POST /auth/accept-invite
// Public endpoint — no authentication required.
func (h *InviteHandler) AcceptInvite(c echo.Context) error {
	ip := c.RealIP()
	ua := c.Request().UserAgent()
	ctx := c.Request().Context()

	var payload model.AcceptInvitePayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if payload.Token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invitation token is required")
	}

	user, accessToken, refreshToken, err := h.inviteService.AcceptInvite(ctx, payload.Token, payload.Name, payload.Password, ip, ua)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "expired") || strings.Contains(errMsg, "already been") || strings.Contains(errMsg, "required") {
			return echo.NewHTTPError(http.StatusBadRequest, errMsg)
		}
		h.server.Logger.Error().Err(err).Msg("failed to accept invite")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to accept invitation")
	}

	// Set session cookies
	platformAuth.SetSessionCookies(c, h.server.Config, accessToken, refreshToken)

	// Resolve full Me profile for frontend
	me, err := h.resolveInviteMe(c, user.ID)
	if err != nil {
		// Still return success — session is created, cookies are set
		h.server.Logger.Error().Err(err).Msg("failed to resolve me after invite accept")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Invitation accepted successfully. Please refresh to load your profile.",
		})
	}

	return c.JSON(http.StatusOK, me)
}

// resolveInviteMe mirrors the AuthHandler.resolveMe logic for building the Me response.
func (h *InviteHandler) resolveInviteMe(c echo.Context, userID string) (*model.Me, error) {
	ctx := c.Request().Context()

	u, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var availableTenants []model.TenantSelectorItem = []model.TenantSelectorItem{}
	hasPlatformStaffRole := u.PlatformRole != nil && (*u.PlatformRole == "super_admin" || *u.PlatformRole == "super_support_agent" || *u.PlatformRole == "super_sales_staff" || *u.PlatformRole == "super_compliance_officer")
	if u.IsPlatformAdmin || hasPlatformStaffRole {
		availableTenants, _ = h.tenantRepo.ListAllTenantsSelector(ctx)
	} else {
		availableTenants, _ = h.userRepo.ListAvailableTenants(ctx, userID)
	}

	// Try to resolve active tenant from the session JWT that was just set
	activeTenantID := ""
	tenantName := ""
	tenantSlug := ""
	membershipID := ""
	role := ""

	// Try to resolve active tenant from the session context
	sessionID := platformAuth.GetSessionID(c)
	if sessionID != "" {
		sess, sessErr := h.userRepo.GetSessionByID(ctx, sessionID)
		if sessErr == nil && sess.ActiveTenantContext != nil {
			activeTenantID = sess.ActiveTenantContext.TenantID
		}
	}

	// If we got the active tenant, resolve details
	if activeTenantID != "" {
		if tUUID, err := uuid.Parse(activeTenantID); err == nil {
			if t, err := h.tenantRepo.GetTenantByID(ctx, tUUID); err == nil {
				tenantName = t.Name
				tenantSlug = t.Slug
			}
		}
		dbMembershipID, dbRoleName, err := h.userRepo.GetActiveMembership(ctx, userID, activeTenantID)
		if err == nil {
			membershipID = dbMembershipID
			role = dbRoleName
		}
	} else if !u.IsPlatformAdmin && !hasPlatformStaffRole && len(availableTenants) > 0 {
		// Fallback to first available for non-platform users
		activeTenantID = availableTenants[0].ID
		if tUUID, err := uuid.Parse(activeTenantID); err == nil {
			if t, err := h.tenantRepo.GetTenantByID(ctx, tUUID); err == nil {
				tenantName = t.Name
				tenantSlug = t.Slug
			}
		}
		dbMembershipID, dbRoleName, err := h.userRepo.GetActiveMembership(ctx, userID, activeTenantID)
		if err == nil {
			membershipID = dbMembershipID
			role = dbRoleName
		}
	}

	effectiveRole := "member"
	if u.PlatformRole != nil && *u.PlatformRole != "" {
		effectiveRole = *u.PlatformRole
	} else if u.IsPlatformAdmin {
		effectiveRole = "super_admin"
	} else if role != "" {
		effectiveRole = role
	}

	var permissions []string = []string{}
	if activeTenantID != "" || u.IsPlatformAdmin || hasPlatformStaffRole {
		if u.IsPlatformAdmin || (u.PlatformRole != nil && *u.PlatformRole == "super_admin") {
			permissions, _ = h.userRepo.ListPermissions(ctx)
		} else if effectiveRole != "" {
			permissions, _ = h.userRepo.ListPermissionsByRole(ctx, effectiveRole, activeTenantID)
		}
	}

	return &model.Me{
		ID:               userID,
		Name:             u.Name,
		Email:            u.Email,
		EmailVerified:    u.EmailVerified,
		Image:            u.Image,
		IsPlatformAdmin:  u.IsPlatformAdmin,
		PlatformRole:     u.PlatformRole,
		ActiveTenantID:   activeTenantID,
		TenantName:       tenantName,
		TenantSlug:       tenantSlug,
		MembershipID:     membershipID,
		Role:             effectiveRole,
		Permissions:      permissions,
		AvailableTenants: availableTenants,
	}, nil
}
