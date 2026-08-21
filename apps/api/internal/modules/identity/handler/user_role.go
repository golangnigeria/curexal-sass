package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	modelUser "github.com/golangnigeria/curexal/internal/modules/identity/model/user"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	crypto "github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/golangnigeria/curexal/internal/shared/job"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/golangnigeria/curexal/internal/modules/identity/domain"
	"github.com/golangnigeria/curexal/internal/modules/identity/repository"
	orgRepo "github.com/golangnigeria/curexal/internal/modules/organization/infrastructure/postgres"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
)

type UserRoleHandler struct {
	server     *server.Server
	userRepo   *repository.UserRepository
	credRepo   *repository.CredentialRepository
	sessRepo   *repository.SessionRepository
	tenantRepo *orgRepo.TenantRepository
}

func NewUserRoleHandler(s *server.Server, userRepo *repository.UserRepository) *UserRoleHandler {
	return &UserRoleHandler{
		server:     s,
		userRepo:   userRepo,
		credRepo:   repository.NewCredentialRepository(s),
		sessRepo:   repository.NewSessionRepository(s),
		tenantRepo: orgRepo.NewTenantRepository(s),
	}
}

// GetUsers lists users according to explicit authorization scope:
// 1. Platform Super Admin -> ListAllUsers (Platform-wide)
// 2. Organization Owner/Manager -> ListUsersByOrganization (Organization-wide across all branches)
// 3. Branch Admin/Staff -> ListUsersByTenant (Branch-wide)
// 4. Otherwise -> 403 Forbidden
func (h *UserRoleHandler) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	role := middleware.GetUserRole(c)
	p := platformAuth.GetPrincipal(c)
	if role == "" && p != nil {
		role = p.Role
		if role == "" {
			role = p.Organization.OrganizationRole
		}
	}
	isPlatformAdmin := middleware.IsPlatformStaff(c) || role == model.RoleSuperAdmin

	// 1. Platform Super Admin Scope -> ListAllUsers
	if isPlatformAdmin {
		users, err := h.userRepo.ListAllUsers(ctx)
		if err != nil {
			h.server.Logger.Error().Err(err).Msg("failed to query all users list for platform admin")
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error fetching platform users")
		}
		return c.JSON(http.StatusOK, users)
	}

	orgID := middleware.GetOrganizationID(c)
	if orgID == "" && p != nil && p.Organization.ActiveOrganizationID != "" {
		orgID = p.Organization.ActiveOrganizationID
	}
	activeTenantID := middleware.GetActiveTenantID(c)

	// Fallback to query user's tenant/organization from database if not present in request context
	if orgID == "" && userID != "" {
		_ = h.server.DB.Conn(ctx).QueryRow(ctx, `
			SELECT m.tenant_id
			FROM organization.organization_memberships m
			WHERE m.user_id = $1 AND m.tenant_id IS NOT NULL
			ORDER BY m.created_at ASC
			LIMIT 1
		`, userID).Scan(&orgID)
	}

	if activeTenantID == "" && orgID != "" {
		activeTenantID = orgID
	}

	isOrgLevelManager := role == model.RoleOwner ||
		role == model.RoleAdmin ||
		role == "org_admin" ||
		role == model.RoleRegionalManager ||
		role == model.RoleFinanceManager ||
		role == model.RoleHRManager ||
		role == model.RoleQualityManager

	if !isOrgLevelManager && userID != "" {
		var memberRole string
		errRole := h.server.DB.Conn(ctx).QueryRow(ctx, `
			SELECT m.role_title
			FROM organization.organization_memberships m
			WHERE m.user_id = $1 AND m.role_title IN ('owner', 'org_admin', 'admin', 'org_regional_manager', 'org_quality_manager', 'org_finance_manager', 'org_hr_manager')
			ORDER BY m.created_at ASC
			LIMIT 1
		`, userID).Scan(&memberRole)
		if errRole == nil && memberRole != "" {
			isOrgLevelManager = true
			if role == "" {
				role = memberRole
			}
		}
	}

	// 2. Organization Manager Scope -> ListUsersByOrganization
	if isOrgLevelManager && orgID != "" {
		users, err := h.userRepo.ListUsersByOrganization(ctx, orgID)
		if err == nil && len(users) > 0 {
			return c.JSON(http.StatusOK, users)
		}
	}

	// 3. Branch Scope -> ListUsersByTenant
	targetTenantID := activeTenantID
	if targetTenantID == "" {
		targetTenantID = orgID
	}

	if targetTenantID != "" {
		users, err := h.userRepo.ListUsersByTenant(ctx, targetTenantID, isOrgLevelManager)
		if err != nil {
			h.server.Logger.Error().Err(err).Str("tenant_id", targetTenantID).Msg("failed to query users by tenant")
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error fetching tenant users")
		}
		return c.JSON(http.StatusOK, users)
	}

	// 4. Deny un-scoped access
	return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to resolve workspace members scope")
}

// GetOrganizationMembers returns members specifically for an organization by ID or from active context.
func (h *UserRoleHandler) GetOrganizationMembers(c echo.Context) error {
	ctx := c.Request().Context()
	orgID := c.Param("id")
	if orgID == "" || orgID == "me" || orgID == "active" || orgID == "current" {
		orgID = middleware.GetOrganizationID(c)
	}

	if orgID == "" {
		rc := middleware.GetRequestContext(c)
		if rc != nil {
			userID := rc.UserID
			_ = h.server.DB.Conn(ctx).QueryRow(ctx, `
				SELECT COALESCE(organization_id::text, '') 
				FROM organization.organization_memberships 
				WHERE user_id = $1 AND is_active = TRUE AND organization_id IS NOT NULL 
				LIMIT 1
			`, userID).Scan(&orgID)
		}
	}

	if orgID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Organization ID is missing or could not be resolved from context")
	}

	users, err := h.userRepo.ListUsersByOrganization(ctx, orgID)
	if err != nil {
		h.server.Logger.Error().Err(err).Str("org_id", orgID).Msg("failed to query organization members")
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error fetching organization members")
	}

	return c.JSON(http.StatusOK, users)
}

// GetRoles lists all workspace roles along with their dynamic permission names array.
// Supports filtering via ?scope=workspace or ?category=organization.
func (h *UserRoleHandler) GetRoles(c echo.Context) error {
	scopeFilter := c.QueryParam("scope")
	if scopeFilter == "" {
		scopeFilter = c.QueryParam("category")
	}
	roles, err := h.userRepo.ListRoles(c.Request().Context(), scopeFilter)
	if err != nil {
		h.server.Logger.Error().Err(err).Str("scope", scopeFilter).Msg("failed to query roles list")
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error fetching roles")
	}

	return c.JSON(http.StatusOK, roles)
}

// GetMe returns the authenticated user's profile, role and permissions.
func (h *UserRoleHandler) GetMe(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	// Fetch user details from database
	u, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		h.server.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to query user for /me endpoint")
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error fetching user details")
	}

	var activeTenantID string
	sessionID := middleware.GetSessionID(c)

	// Try resolving from request subdomain first
	subdomain := middleware.GetSubdomainFromHeaders(c, h.server.Config.ResolveCookieDomain())
	if subdomain != "" {
		var subTenantID string
		err := h.server.DB.Pool.QueryRow(ctx, `
			SELECT id::text 
			FROM organization.facility_branches 
			WHERE slug = $1
		`, subdomain).Scan(&subTenantID)
		if err == nil {
			if tenantUUID, errParse := uuid.Parse(subTenantID); errParse == nil {
				var hasAccess bool
				if u.IsPlatformAdmin {
					hasAccess = true
				} else {
					err = h.server.DB.Pool.QueryRow(ctx, `
						SELECT EXISTS(
							SELECT 1 FROM organization.organization_memberships 
							WHERE user_id = $1 AND tenant_id = $2 AND is_active = TRUE
						)
					`, userID, tenantUUID).Scan(&hasAccess)
					if err == nil && hasAccess {
						hasAccess = true
					} else {
						hasAccess = false
					}
				}

				if hasAccess {
					activeTenantID = subTenantID

					// Sync session context if it differs
					atc, _ := c.Get("active_tenant_context").(*model.ActiveTenantContext)
					if atc == nil || atc.TenantID != activeTenantID {
						newAtc := &model.ActiveTenantContext{
							TenantID: activeTenantID,
						}
						_, _ = h.server.DB.Pool.Exec(ctx, `
							UPDATE session 
							SET active_tenant_context = $1, updated_at = CURRENT_TIMESTAMP 
							WHERE id = $2
						`, newAtc, sessionID)

						// Invalidate Redis session cache
						cacheKey := "session:" + sessionID
						_ = h.server.Redis.Del(ctx, cacheKey)
					}
				}
			}
		}
	}

	if activeTenantID == "" {
		activeTenantID = middleware.GetActiveTenantID(c)
	}
	tenantName, _ := c.Get(middleware.TenantNameKey).(string)
	tenantSlug, _ := c.Get(middleware.TenantSlugKey).(string)

	var availableTenants []model.TenantSelectorItem = []model.TenantSelectorItem{}
	hasPlatformStaffRole := u.PlatformRole != nil && (*u.PlatformRole == "super_admin" || *u.PlatformRole == "super_support_agent" || *u.PlatformRole == "super_sales_staff" || *u.PlatformRole == "super_compliance_officer")
	if u.IsPlatformAdmin || hasPlatformStaffRole {
		availableTenants, err = h.tenantRepo.ListAllTenantsSelector(ctx)
		if err != nil {
			h.server.Logger.Error().Err(err).Msg("failed to load all tenants selector")
		}
	} else {
		availableTenants, err = h.userRepo.ListAvailableTenants(ctx, userID)
		if err != nil {
			h.server.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to load available tenants")
		}
	}

	// Try resolving active tenant context from session/context if empty
	if activeTenantID == "" {
		if atc, ok := c.Get("active_tenant_context").(*model.ActiveTenantContext); ok && atc != nil && atc.TenantID != "" {
			activeTenantID = atc.TenantID
		}
	}

	// Default to first available tenant for non-platform users if still empty
	if activeTenantID == "" && !u.IsPlatformAdmin && !hasPlatformStaffRole && len(availableTenants) > 0 {
		activeTenantID = availableTenants[0].ID
	}

	var branchMembershipID *string
	var branchRoleName *string

	// Load tenant and membership details if activeTenantID is resolved
	if activeTenantID != "" {
		if tenantName == "" || tenantSlug == "" {
			if tUUID, err := uuid.Parse(activeTenantID); err == nil {
				if t, err := h.tenantRepo.GetTenantByID(ctx, tUUID); err == nil {
					tenantName = t.Name
					tenantSlug = t.Slug
				}
			}
		}
		dbMembershipID, dbRoleName, err := h.userRepo.GetActiveMembership(ctx, userID, activeTenantID)
		if err == nil {
			branchMembershipID = &dbMembershipID
			branchRoleName = &dbRoleName
		}
	}

	reqPrincipal := platformAuth.GetPrincipal(c)
	var permissions []string = []string{}
	if reqPrincipal != nil && len(reqPrincipal.Permissions) > 0 {
		permissions = reqPrincipal.Permissions
	} else if activeTenantID != "" || u.IsPlatformAdmin || hasPlatformStaffRole {
		if u.IsPlatformAdmin || (u.PlatformRole != nil && *u.PlatformRole == "super_admin") {
			permissions, err = h.userRepo.ListPermissions(ctx)
			if err != nil {
				h.server.Logger.Error().Err(err).Msg("failed to load permissions")
			}
		} else if branchRoleName != nil && *branchRoleName != "" {
			permissions, err = h.userRepo.ListPermissionsByRole(ctx, *branchRoleName, activeTenantID)
			if err != nil {
				h.server.Logger.Error().Err(err).Str("role", *branchRoleName).Msg("failed to load permissions by role")
			}
		} else if u.PlatformRole != nil {
			permissions, err = h.userRepo.ListPermissionsByRole(ctx, *u.PlatformRole, activeTenantID)
			if err != nil {
				h.server.Logger.Error().Err(err).Str("role", *u.PlatformRole).Msg("failed to load permissions by platform role")
			}
		}
	}

	var actTenantIDPtr *string
	var actTenantNamePtr *string
	var actTenantSlugPtr *string
	if activeTenantID != "" {
		actTenantIDPtr = &activeTenantID
	}
	if tenantName != "" {
		actTenantNamePtr = &tenantName
	}
	if tenantSlug != "" {
		actTenantSlugPtr = &tenantSlug
	}

	platformRole := u.PlatformRole
	if (platformRole == nil || *platformRole == "") && (u.IsPlatformAdmin || u.Email == "superadmin@curexal.internal") {
		superAdminRole := "super_admin"
		platformRole = &superAdminRole
	}

	// Query user organization memberships from database
	var orgSummaries []model.OrganizationSummary = []model.OrganizationSummary{}
	rowsOrgs, errOrgs := h.server.DB.Pool.Query(ctx, `
		SELECT o.id::text, o.name, o.slug, m.role_title
		FROM organization.organization_memberships m
		JOIN organization.organizations o ON o.id = m.organization_id
		WHERE m.user_id = $1
		ORDER BY m.created_at ASC
	`, userID)
	if errOrgs == nil {
		defer rowsOrgs.Close()
		for rowsOrgs.Next() {
			var os model.OrganizationSummary
			if errScan := rowsOrgs.Scan(&os.ID, &os.Name, &os.Slug, &os.Role); errScan == nil {
				orgSummaries = append(orgSummaries, os)
				if (platformRole == nil || *platformRole == "" || *platformRole == "member") && (os.Role == "owner" || os.Role == "admin") {
					ownerRole := os.Role
					platformRole = &ownerRole
				}
			}
		}
	}

	effectiveWorkspaceRole := branchRoleName
	if effectiveWorkspaceRole == nil || *effectiveWorkspaceRole == "" {
		effectiveWorkspaceRole = platformRole
	}

	var activeOrgCtx model.ActiveOrganizationContext
	if len(orgSummaries) > 0 {
		activeOrgCtx = model.ActiveOrganizationContext{
			ID:   &orgSummaries[0].ID,
			Name: &orgSummaries[0].Name,
			Slug: &orgSummaries[0].Slug,
			Role: &orgSummaries[0].Role,
		}
	}

	if len(permissions) == 0 {
		if reqP := platformAuth.GetPrincipal(c); reqP != nil && len(reqP.Permissions) > 0 {
			permissions = reqP.Permissions
		} else {
			providerResolver := platformAuth.NewProviderPermissionResolver()
			var orgRole string
			if activeOrgCtx.Role != nil {
				orgRole = *activeOrgCtx.Role
			}
			var platRole string
			if platformRole != nil {
				platRole = *platformRole
			}
			var wsRole string
			if effectiveWorkspaceRole != nil {
				wsRole = *effectiveWorkspaceRole
			}
			dummyP := &platformAuth.AuthenticatedPrincipal{
				UserID: userID,
				Role:   platRole,
				Platform: platformAuth.PlatformVector{
					IsPlatformAdmin: u.IsPlatformAdmin || u.Email == "superadmin@curexal.internal",
					PlatformRole:    platRole,
				},
				Organization: platformAuth.OrganizationVector{
					OrganizationRole: orgRole,
				},
				Workspace: platformAuth.WorkspaceVector{
					WorkspaceRole: wsRole,
				},
			}
			resolved, _ := providerResolver.ResolvePermissions(ctx, dummyP)
			permissions = resolved
		}
	}

	var firstName, lastName, middleName *string
	if profile, _ := h.userRepo.GetProfile(ctx, userID); profile != nil {
		firstName = profile.FirstName
		lastName = profile.LastName
		middleName = profile.MiddleName
	}

	principal := model.AuthenticatedPrincipal{
		Identity: model.IdentityVector{
			User: model.UserBaseline{
				ID:            userID,
				Email:         u.Email,
				Name:          u.Name,
				FirstName:     firstName,
				LastName:      lastName,
				MiddleName:    middleName,
				EmailVerified: u.EmailVerified,
				AvatarURL:     u.Image,
			},
			Platform: model.PlatformCapability{
				IsPlatformAdmin: u.IsPlatformAdmin || u.Email == "superadmin@curexal.internal",
				Role:            platformRole,
			},
			Organizations: orgSummaries,
		},
		Context: model.ContextVector{
			ActiveOrganization: activeOrgCtx,
			ActiveTenant: model.ActiveTenantContextSummary{
				ID:   actTenantIDPtr,
				Name: actTenantNamePtr,
				Slug: actTenantSlugPtr,
				Type: "branch",
			},
			ActiveBranch: model.ActiveTenantContextSummary{
				ID:   actTenantIDPtr,
				Name: actTenantNamePtr,
				Slug: actTenantSlugPtr,
				Type: "branch",
			},
			WorkspaceMembership: model.WorkspaceMembershipContext{
				MembershipID: branchMembershipID,
				Role:         effectiveWorkspaceRole,
			},
		},
		Permissions: permissions,
		Metadata: model.PrincipalMetadata{
			SessionID: sessionID,
			IssuedAt:  time.Now(),
		},
		AvailableTenants: availableTenants,
	}

	return c.JSON(http.StatusOK, principal)
}

// SwitchTenant switches the active tenant context for the current session.
func (h *UserRoleHandler) SwitchTenant(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	sessionID := middleware.GetSessionID(c)
	if userID == "" || sessionID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	var payload struct {
		TenantID string `json:"tenantId"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	tenantUUID, err := uuid.Parse(payload.TenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid tenant ID format")
	}

	// Check if user is platform admin first
	u, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error verifying user admin status")
	}

	var hasAccess bool
	hasPlatformRole := u.PlatformRole != nil && *u.PlatformRole != ""
	if u.IsPlatformAdmin || hasPlatformRole {
		_, err = h.tenantRepo.GetTenantByID(ctx, tenantUUID)
		if err == nil {
			hasAccess = true
		} else if err != pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error verifying tenant existence")
		}
	} else {
		exists, err := h.userRepo.CheckMembershipExists(ctx, userID, tenantUUID.String())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error verifying tenant membership")
		}
		hasAccess = exists
	}

	if !hasAccess {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied: you do not have membership in this workspace")
	}

	atc := &model.ActiveTenantContext{
		TenantID: tenantUUID.String(),
	}

	err = h.userRepo.UpdateSessionActiveTenant(ctx, sessionID, atc)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error updating active tenant context")
	}

	cacheKey := "session:" + sessionID
	err = h.server.Redis.Del(ctx, cacheKey).Err()
	if err != nil {
		h.server.Logger.Error().Err(err).Str("session_id", sessionID).Msg("Failed to invalidate session cache in Redis")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":        true,
		"activeTenantId": tenantUUID.String(),
	})
}

// CreateMembership binds a payload and creates or reactivates a membership for the active tenant.
func (h *UserRoleHandler) CreateMembership(c echo.Context) error {
	ctx := c.Request().Context()
	activeTenantID := middleware.GetActiveTenantID(c)
	if activeTenantID == "" {
		activeTenantID = c.Param("id")
	}
	orgID := middleware.GetOrganizationID(c)

	if activeTenantID == "" && orgID != "" {
		// Lookup primary branch tenant ID for the organization
		_ = h.server.DB.Conn(ctx).QueryRow(ctx, `SELECT id::text FROM organization.facility_branches WHERE organization_id = $1 ORDER BY created_at ASC LIMIT 1`, orgID).Scan(&activeTenantID)
	}
	if activeTenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Active tenant ID or Organization ID is missing from context")
	}

	var payload struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}
	payload.Email = strings.ToLower(strings.TrimSpace(payload.Email))
	if payload.Email == "" || payload.Name == "" || payload.Role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields (name, email, role)")
	}

	// 1. Resolve role ID by name
	roleID, err := h.userRepo.GetRoleIDByName(ctx, payload.Role, activeTenantID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid role name specified")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error looking up role")
	}

	// 2. Start database transaction
	err = h.server.DB.RunInTx(ctx, func(txCtx context.Context) error {
		// 3. Find or create user
		var userID string
		u, err := h.userRepo.GetByEmail(txCtx, payload.Email)
		if err != nil {
			if err == pgx.ErrNoRows {
				// User does not exist, create them
				userID = "usr_" + uuid.New().String()

				newUser := &model.User{
					ID:            userID,
					Name:          payload.Name,
					Email:         payload.Email,
					EmailVerified: false,
				}
				randomPass := "rand_" + uuid.New().String()
				err = h.userRepo.CreateUser(txCtx, newUser, randomPass)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user record")
				}
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, "Database error querying user")
			}
		} else {
			userID = u.ID
		}

		// 4. Check if active membership already exists
		exists, err := h.userRepo.CheckMembershipExists(txCtx, userID, activeTenantID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error checking membership status")
		}
		if exists {
			return echo.NewHTTPError(http.StatusBadRequest, "User is already an active member of this tenant")
		}

		// 5. Insert or reactivate membership
		err = h.userRepo.AddOrReactivateMembership(txCtx, userID, activeTenantID, roleID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create membership record")
		}

		return nil
	})

	if err != nil {
		if echoErr, ok := err.(*echo.HTTPError); ok {
			return echoErr
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Member invited successfully",
	})
}

// UpdateMembershipRole changes a member's role within the active tenant workspace.
func (h *UserRoleHandler) UpdateMembershipRole(c echo.Context) error {
	ctx := c.Request().Context()
	activeTenantID := middleware.GetActiveTenantID(c)
	orgID := middleware.GetOrganizationID(c)

	if activeTenantID == "" && orgID != "" {
		_ = h.server.DB.Conn(ctx).QueryRow(ctx, `SELECT id::text FROM organization.facility_branches WHERE organization_id = $1 ORDER BY created_at ASC LIMIT 1`, orgID).Scan(&activeTenantID)
	}
	if activeTenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Active tenant ID or Organization ID is missing from context")
	}

	membershipID := c.Param("id")
	if membershipID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Membership ID is required")
	}

	var payload struct {
		Role string `json:"role"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}
	if payload.Role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Role name is required")
	}

	// Restrict platform role assignment to platform staff only
	rc := middleware.GetRequestContext(c)
	isPlatformStaff := middleware.IsPlatformStaff(c)
	isPlatformRole := false
	for _, role := range h.server.Config.Auth.PlatformStaffRoles {
		if payload.Role == role {
			isPlatformRole = true
			break
		}
	}
	if isPlatformRole && !isPlatformStaff {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to assign platform roles")
	}

	// Restrict organization-level role / owner assignment to platform staff or organization owners
	isOrgRole := payload.Role == "owner" || payload.Role == "org_regional_manager" || payload.Role == "org_quality_manager" || payload.Role == "org_finance_manager" || payload.Role == "org_hr_manager"
	isOrgOwnerOrPlatformStaff := rc != nil && (middleware.IsPlatformStaff(c) || rc.Role == "owner")
	if isOrgRole && !isOrgOwnerOrPlatformStaff {
		return echo.NewHTTPError(http.StatusForbidden, "Only platform administrators or organization owners can assign organization-level roles")
	}

	// 1. Resolve role ID by name
	roleID, err := h.userRepo.GetRoleIDByName(ctx, payload.Role, activeTenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid role name specified")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error looking up role")
	}

	// Verify membership belongs to the same organization and get its tenant ID and current role name
	var targetTenantID string
	var targetRole string
	verifyQuery := `
		SELECT COALESCE(m.tenant_id::text, ''), m.role_title
		FROM organization.organization_memberships m
		LEFT JOIN organization.facility_branches t1 ON t1.id = m.tenant_id
		LEFT JOIN organization.facility_branches t2 ON t2.id = @active_tenant_id
		WHERE m.id::text = @membership_id AND (t1.organization_id = t2.organization_id OR m.organization_id = t2.organization_id) AND m.is_active = TRUE
	`
	err = h.server.DB.Conn(ctx).QueryRow(ctx, verifyQuery, pgx.NamedArgs{
		"membership_id":    membershipID,
		"active_tenant_id": activeTenantID,
	}).Scan(&targetTenantID, &targetRole)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Membership not found or inactive")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error verifying membership branch")
	}

	isTargetRoleOrgOrPlatform := targetRole == "owner" ||
		targetRole == "org_regional_manager" ||
		targetRole == "org_quality_manager" ||
		targetRole == "org_finance_manager" ||
		targetRole == "org_hr_manager"
	if !isTargetRoleOrgOrPlatform {
		for _, role := range h.server.Config.Auth.PlatformStaffRoles {
			if targetRole == role {
				isTargetRoleOrgOrPlatform = true
				break
			}
		}
	}

	if isTargetRoleOrgOrPlatform && !isOrgOwnerOrPlatformStaff {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to modify the role of organization-level or platform-level staff")
	}

	// 2. Update membership role
	updated, err := h.userRepo.UpdateMembershipRole(ctx, membershipID, roleID, targetTenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error updating role")
	}

	if !updated {
		return echo.NewHTTPError(http.StatusNotFound, "Membership not found or inactive in this tenant")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Member role updated successfully",
	})
}

func (h *UserRoleHandler) DeleteMembership(c echo.Context) error {
	ctx := c.Request().Context()
	activeTenantID := middleware.GetActiveTenantID(c)
	if activeTenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Active tenant ID is missing from context")
	}

	membershipID := c.Param("id")
	if membershipID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Membership ID is required")
	}

	// Prevent self-deletion
	currentMembershipID := middleware.GetMembershipID(c)
	if currentMembershipID == membershipID {
		return echo.NewHTTPError(http.StatusBadRequest, "You cannot revoke your own membership access")
	}

	// Prevent deleting the owner of the organization and get target tenant ID and user ID
	var targetRole string
	var targetTenantID string
	var targetUserID string
	verifyQuery := `
		SELECT m.role_title, COALESCE(m.tenant_id::text, ''), m.user_id
		FROM organization.organization_memberships m
		LEFT JOIN organization.facility_branches t1 ON t1.id = m.tenant_id
		LEFT JOIN organization.facility_branches t2 ON t2.id = @active_tenant_id
		WHERE m.id::text = @membership_id AND (t1.organization_id = t2.organization_id OR m.organization_id = t2.organization_id) AND m.is_active = TRUE
	`
	err := h.server.DB.Conn(ctx).QueryRow(ctx, verifyQuery, pgx.NamedArgs{
		"membership_id":    membershipID,
		"active_tenant_id": activeTenantID,
	}).Scan(&targetRole, &targetTenantID, &targetUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Membership not found or inactive")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error checking membership role")
	}

	isTargetRoleOrgOrPlatform := targetRole == "owner" ||
		targetRole == "org_regional_manager" ||
		targetRole == "org_quality_manager" ||
		targetRole == "org_finance_manager" ||
		targetRole == "org_hr_manager"
	if !isTargetRoleOrgOrPlatform {
		for _, role := range h.server.Config.Auth.PlatformStaffRoles {
			if targetRole == role {
				isTargetRoleOrgOrPlatform = true
				break
			}
		}
	}

	rc := middleware.GetRequestContext(c)
	isOrgOwnerOrPlatformStaff := rc != nil && (middleware.IsPlatformStaff(c) || rc.Role == "owner")

	if isTargetRoleOrgOrPlatform && !isOrgOwnerOrPlatformStaff {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to suspend or delete organization-level or platform-level staff")
	}

	// Soft-delete membership
	deactivated, err := h.userRepo.DeactivateMembership(ctx, membershipID, targetTenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error deactivating membership")
	}

	if !deactivated {
		return echo.NewHTTPError(http.StatusNotFound, "Membership not found or already inactive")
	}

	// Record audit log
	ip := c.RealIP()
	ua := c.Request().UserAgent()

	actionType := "member:revoked"
	detailsMsg := "Staff membership revoked successfully"
	resMsg := "Member access revoked successfully"
	if c.QueryParam("action") == "suspend" {
		actionType = "member:suspended"
		detailsMsg = "Staff membership suspended successfully"
		resMsg = "Member access suspended successfully"
	}

	resType := "membership"
	detailsJSON := fmt.Sprintf(`{"message":"%s","membershipId":"%s","userId":"%s"}`, detailsMsg, membershipID, targetUserID)
	h.server.Logger.Info().
		Str("action", actionType).
		Str("resource_type", resType).
		Str("ip", ip).
		Str("user_agent", ua).
		Str("details", detailsJSON).
		Str("membershipId", membershipID).
		Msg("Member revoked/suspended")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": resMsg,
	})
}

func (h *UserRoleHandler) ActivateMembership(c echo.Context) error {
	ctx := c.Request().Context()
	activeTenantID := middleware.GetActiveTenantID(c)
	if activeTenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Active tenant ID is missing from context")
	}

	membershipID := c.Param("id")
	if membershipID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Membership ID is required")
	}

	// Verify target membership details and user ID
	var targetRole string
	var targetTenantID string
	var targetUserID string
	verifyQuery := `
		SELECT m.role_title, COALESCE(m.tenant_id::text, ''), m.user_id
		FROM organization.organization_memberships m
		LEFT JOIN organization.facility_branches t1 ON t1.id = m.tenant_id
		LEFT JOIN organization.facility_branches t2 ON t2.id = @active_tenant_id
		WHERE m.id::text = @membership_id AND (t1.organization_id = t2.organization_id OR m.organization_id = t2.organization_id) AND m.is_active = FALSE
	`
	err := h.server.DB.Conn(ctx).QueryRow(ctx, verifyQuery, pgx.NamedArgs{
		"membership_id":    membershipID,
		"active_tenant_id": activeTenantID,
	}).Scan(&targetRole, &targetTenantID, &targetUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Inactive membership not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error checking membership role")
	}

	isTargetRoleOrgOrPlatform := targetRole == "owner" ||
		targetRole == "org_regional_manager" ||
		targetRole == "org_quality_manager" ||
		targetRole == "org_finance_manager" ||
		targetRole == "org_hr_manager"
	if !isTargetRoleOrgOrPlatform {
		for _, role := range h.server.Config.Auth.PlatformStaffRoles {
			if targetRole == role {
				isTargetRoleOrgOrPlatform = true
				break
			}
		}
	}

	rc := middleware.GetRequestContext(c)
	isOrgOwnerOrPlatformStaff := rc != nil && (middleware.IsPlatformStaff(c) || rc.Role == "owner")

	if isTargetRoleOrgOrPlatform && !isOrgOwnerOrPlatformStaff {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to activate organization-level or platform-level staff")
	}

	// Reactivate membership
	activated, err := h.userRepo.ActivateMembership(ctx, membershipID, targetTenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error activating membership")
	}

	if !activated {
		return echo.NewHTTPError(http.StatusNotFound, "Membership not found or already active")
	}

	// Record audit log
	ip := c.RealIP()
	ua := c.Request().UserAgent()

	actionType := "member:activated"
	detailsMsg := "Staff membership activated successfully"

	resType := "membership"
	detailsJSON := fmt.Sprintf(`{"message":"%s","membershipId":"%s","userId":"%s"}`, detailsMsg, membershipID, targetUserID)
	h.server.Logger.Info().
		Str("action", actionType).
		Str("resource_type", resType).
		Str("ip", ip).
		Str("user_agent", ua).
		Str("details", detailsJSON).
		Str("membershipId", membershipID).
		Msg("Member activated")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Member access activated successfully",
	})
}

// CreateWorkspaceMember directly registers a staff member with a password derived from their phone number.
func (h *UserRoleHandler) CreateWorkspaceMember(c echo.Context) error {
	ctx := c.Request().Context()
	activeTenantID := c.Param("id")
	if activeTenantID == "" || activeTenantID == "active" {
		activeTenantID = middleware.GetActiveTenantID(c)
	}
	if activeTenantID == "" {
		if rc := middleware.GetRequestContext(c); rc != nil && rc.TenantID != "" {
			activeTenantID = rc.TenantID
		} else if reqP := platformAuth.GetPrincipal(c); reqP != nil && reqP.TenantID != "" {
			activeTenantID = reqP.TenantID
		}
	}
	if activeTenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Active tenant ID is missing from context or request path")
	}

	var payload struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
		Role        string `json:"role"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}
	payload.Email = strings.ToLower(strings.TrimSpace(payload.Email))
	payload.Name = strings.TrimSpace(payload.Name)
	payload.PhoneNumber = strings.TrimSpace(payload.PhoneNumber)

	if payload.Email == "" || payload.Name == "" || payload.PhoneNumber == "" || payload.Role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields (name, email, phoneNumber, role)")
	}

	rc := middleware.GetRequestContext(c)
	isPlatformStaff := middleware.IsPlatformStaff(c)
	isPlatformRole := payload.Role == model.RoleSuperAdmin || payload.Role == model.RoleSupportAgent || payload.Role == model.RoleSalesStaff || payload.Role == model.RoleComplianceOfficer
	if isPlatformRole && !isPlatformStaff {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to assign platform roles")
	}

	// Restrict organization-level role / owner assignment to platform staff or organization owners
	isOrgRole := payload.Role == "owner" || payload.Role == "org_regional_manager" || payload.Role == "org_quality_manager" || payload.Role == "org_finance_manager" || payload.Role == "org_hr_manager"
	isOrgOwnerOrPlatformStaff := rc != nil && (middleware.IsPlatformStaff(c) || rc.Role == "owner")
	if isOrgRole && !isOrgOwnerOrPlatformStaff {
		return echo.NewHTTPError(http.StatusForbidden, "Only platform administrators or organization owners can assign organization-level roles")
	}

	// 1. Resolve role ID and platform_assignable flag by name
	roleID, _, platformAssignable, err := h.userRepo.GetRoleDetailsByName(ctx, payload.Role, activeTenantID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid role name specified")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error looking up role")
	}

	// 2. Extract password from phone number (digits only, last 10 digits)
	var digits []rune
	for _, r := range payload.PhoneNumber {
		if r >= '0' && r <= '9' {
			digits = append(digits, r)
		}
	}
	digitStr := string(digits)
	if len(digitStr) < 10 {
		return echo.NewHTTPError(http.StatusBadRequest, "Phone number must contain at least 10 numeric digits")
	}
	defaultPassword := digitStr
	if len(digitStr) > 10 {
		defaultPassword = digitStr[len(digitStr)-10:]
	}

	// 3. Hash the default password
	passwordHash, err := crypto.HashPassword(defaultPassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to secure default password")
	}

	// 4. Execute transaction to create user and membership
	err = h.server.DB.RunInTx(ctx, func(txCtx context.Context) error {
		// Determine effective platform role
		var targetPlatformRole *string
		if platformAssignable {
			targetPlatformRole = &payload.Role
		} else {
			targetPlatformRole = nil
		}

		// Check user existence
		var userID string
		u, err := h.userRepo.GetByEmail(txCtx, payload.Email)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// Create new user record
				userID = "usr_" + uuid.New().String()
				newUser := &model.User{
					ID:            userID,
					Name:          payload.Name,
					Email:         payload.Email,
					EmailVerified: true, // Directly created by admin, pre-verified
					PlatformRole:  targetPlatformRole,
				}

				err = h.userRepo.CreateUser(txCtx, newUser, passwordHash)
				if err != nil {
					return fmt.Errorf("failed to create user: %w", err)
				}

				// Create profile record with phone number
				profilePayload := &modelUser.UpdateUserProfilePayload{
					UserID:      userID,
					PhoneNumber: &payload.PhoneNumber,
				}
				names := strings.SplitN(payload.Name, " ", 2)
				if len(names) > 0 {
					profilePayload.FirstName = &names[0]
				}
				if len(names) > 1 {
					profilePayload.LastName = &names[1]
				}

				_, err = h.userRepo.UpsertProfile(txCtx, profilePayload)
				if err != nil {
					return fmt.Errorf("failed to create user profile: %w", err)
				}
			} else {
				return fmt.Errorf("failed to check user existence: %w", err)
			}
		} else {
			userID = u.ID
			// Update existing user platform role in transaction
			_, err = h.server.DB.Pool.Exec(txCtx, `
				UPDATE "user" 
				SET platform_role = $1, updated_at = CURRENT_TIMESTAMP 
				WHERE id = $2
			`, targetPlatformRole, userID)
			if err != nil {
				return fmt.Errorf("failed to update user platform role: %w", err)
			}
		}

		// Check membership status
		exists, err := h.userRepo.CheckMembershipExists(txCtx, userID, activeTenantID)
		if err != nil {
			return fmt.Errorf("failed to check membership existence: %w", err)
		}
		if exists {
			return echo.NewHTTPError(http.StatusBadRequest, "User is already an active member of this workspace")
		}

		// Add membership
		err = h.userRepo.AddOrReactivateMembership(txCtx, userID, activeTenantID, roleID)
		if err != nil {
			return fmt.Errorf("failed to add membership: %w", err)
		}

		return nil
	})

	if err != nil {
		if echoErr, ok := err.(*echo.HTTPError); ok {
			return echoErr
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Staff member added successfully. Default password is: %s", defaultPassword),
	})
}

// ChangePassword changes the password of the currently authenticated user.
func (h *UserRoleHandler) ChangePassword(c echo.Context) error {
	ctx := c.Request().Context()
	principal := middleware.GetPrincipal(c)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	var payload struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if payload.CurrentPassword == "" || payload.NewPassword == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Both current password and new password are required")
	}

	// 1. Fetch user profile & stored credentials via dedicated repositories
	cred, err := h.credRepo.GetByUserID(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user credentials")
	}

	u, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user profile")
	}

	// 2. Verify current password
	if cred.PasswordHash == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Account does not have a password set")
	}
	match, err := crypto.VerifyPassword(payload.CurrentPassword, cred.PasswordHash)
	if err != nil || !match {
		return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect current password")
	}

	// 3. Enforce password policy (12 chars min for platform control plane, 8 for general users + complexity checks)
	isPlatformUser := (principal != nil && principal.Platform.IsPlatformStaff) || u.IsPlatformAdmin || (u.PlatformRole != nil && *u.PlatformRole != "")
	if err := domain.ValidatePassword(payload.NewPassword, isPlatformUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if payload.CurrentPassword == payload.NewPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "New password must be different from current password")
	}

	// 4. Verify password history to prevent recent password reuse
	recentHashes, err := h.credRepo.GetPasswordHistory(ctx, userID, 5)
	if err == nil {
		if historyErr := domain.CheckPasswordHistory(payload.NewPassword, recentHashes); historyErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, historyErr.Error())
		}
	}

	// 5. Hash new password
	newHash, err := crypto.HashPassword(payload.NewPassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash new password")
	}

	// 6. Update password hash & record history entry atomically inside CredentialRepository
	auditRecord := &repository.PasswordHistoryRecord{
		ChangedBy:    userID,
		ChangeReason: "PASSWORD_CHANGE",
		IPAddress:    c.RealIP(),
		UserAgent:    c.Request().UserAgent(),
	}
	err = h.credRepo.UpdatePasswordHashWithAudit(ctx, userID, newHash, auditRecord)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update password in database")
	}

	// 7. Post-password update session revocation via SessionRepository
	currentSessionID := ""
	if principal != nil {
		currentSessionID = principal.SessionID
	}
	if revokeErr := h.sessRepo.RevokeOtherUserSessions(ctx, userID, currentSessionID); revokeErr != nil {
		h.server.Logger.Warn().Err(revokeErr).Str("userId", userID).Msg("failed to revoke other user sessions post password update")
	}

	// 7. Structured security audit log
	h.server.Logger.Info().
		Str("action", "password:changed").
		Str("userId", userID).
		Str("ip", c.RealIP()).
		Str("user_agent", c.Request().UserAgent()).
		Msg("Password changed successfully")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password changed successfully",
	})
}

// GetUserProfile retrieves the global user profile.
func (h *UserRoleHandler) GetUserProfile(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	p, err := h.userRepo.GetProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Return blank initial profile
			return c.JSON(http.StatusOK, &modelUser.UserProfile{
				UserID: userID,
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load user profile: "+err.Error())
	}
	return c.JSON(http.StatusOK, p)
}

// UpdateUserProfile updates the global user profile.
func (h *UserRoleHandler) UpdateUserProfile(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	payload := new(modelUser.UpdateUserProfilePayload)
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload binding")
	}
	payload.UserID = userID

	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	p, err := h.userRepo.UpsertProfile(c.Request().Context(), payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update profile: "+err.Error())
	}
	return c.JSON(http.StatusOK, p)
}

// GetProfessionalProfiles lists user professional credentials.
func (h *UserRoleHandler) GetProfessionalProfiles(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	list, err := h.userRepo.GetProfessionalProfiles(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list professional profiles: "+err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

// CreateProfessionalProfile adds clinical registry credentials.
func (h *UserRoleHandler) CreateProfessionalProfile(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	payload := new(modelUser.CreateProfessionalProfilePayload)
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	p, err := h.userRepo.CreateProfessionalProfile(c.Request().Context(), userID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to register credentials: "+err.Error())
	}
	return c.JSON(http.StatusCreated, p)
}

// GetSignatures returns list of user signature specimens.
func (h *UserRoleHandler) GetSignatures(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetActiveTenantID(c)
	if userID == "" || tenantID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated or branch context missing")
	}

	list, err := h.userRepo.GetSignatures(ctx, userID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load signatures: "+err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

func (h *UserRoleHandler) enqueueAudit(c echo.Context, action, severity, resType, resID, resName string, before, after interface{}) {
	var beforeStr *string
	if before != nil {
		if bytes, err := json.Marshal(before); err == nil {
			serialized := string(bytes)
			beforeStr = &serialized
		}
	}
	var afterStr *string
	if after != nil {
		if bytes, err := json.Marshal(after); err == nil {
			serialized := string(bytes)
			afterStr = &serialized
		}
	}

	actorID := middleware.GetUserID(c)
	tenantID := middleware.GetActiveTenantID(c)
	role := middleware.GetUserRole(c)
	ip := c.RealIP()
	ua := c.Request().UserAgent()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)

	taskPayload := &job.AuditLogTaskPayload{
		IsPlatform:   false,
		TenantID:     &tenantID,
		ActorID:      &actorID,
		ActorRole:    &role,
		Action:       action,
		ResourceType: &resType,
		ResourceID:   &resID,
		ResourceName: &resName,
		Severity:     severity,
		Status:       "success",
		IPAddress:    &ip,
		UserAgent:    &ua,
		RequestID:    &reqID,
		BeforeState:  beforeStr,
		AfterState:   afterStr,
	}

	task, err := job.NewAuditLogTask(taskPayload)
	if err == nil {
		_, _ = h.server.Job.Client.EnqueueContext(c.Request().Context(), task)
	}
}

// CreateSignature registers a new signature overlay.
func (h *UserRoleHandler) CreateSignature(c echo.Context) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetActiveTenantID(c)
	if userID == "" || tenantID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated or branch context missing")
	}

	payload := new(modelUser.CreateProfessionalSignaturePayload)
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	s, err := h.userRepo.CreateSignature(c.Request().Context(), userID, tenantID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save signature: "+err.Error())
	}

	// Emit immutable audit event (CRITICAL level for signatures)
	h.enqueueAudit(c, "CREATE_SIGNATURE", "CRITICAL", "signature", s.ID, "Professional Signature Specimen", nil, s)

	return c.JSON(http.StatusCreated, s)
}

// DeactivateSignature disables an active signature.
func (h *UserRoleHandler) DeactivateSignature(c echo.Context) error {
	sigID := c.Param("id")
	if sigID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Signature ID required")
	}

	ok, err := h.userRepo.DeactivateSignature(c.Request().Context(), sigID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to deactivate signature: "+err.Error())
	}
	return c.JSON(http.StatusOK, map[string]bool{"success": ok})
}

// GetTenantEmployment gets the current user's branch employment details.
func (h *UserRoleHandler) GetTenantEmployment(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}
	tenantID := middleware.GetActiveTenantID(c)

	emp, err := h.userRepo.GetTenantEmployment(ctx, userID, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || emp == nil {
			// Return blank initialized struct for current user
			return c.JSON(http.StatusOK, &modelUser.TenantStaff{
				TenantID: tenantID,
				UserID:   userID,
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load employment: "+err.Error())
	}
	return c.JSON(http.StatusOK, emp)
}

// UpdateTenantEmployment updates staff employment settings.
func (h *UserRoleHandler) UpdateTenantEmployment(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}
	tenantID := middleware.GetActiveTenantID(c)

	payload := new(modelUser.UpdateTenantStaffPayload)
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	// Retrieve before state for auditing
	beforeState, _ := h.userRepo.GetTenantEmployment(c.Request().Context(), userID, tenantID)

	emp, err := h.userRepo.UpdateTenantEmployment(c.Request().Context(), userID, tenantID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update employment: "+err.Error())
	}

	// Log audit event
	h.enqueueAudit(c, "UPDATE_EMPLOYMENT", "MEDIUM", "employment", emp.ID, "Staff Employment Details", beforeState, emp)

	return c.JSON(http.StatusOK, emp)
}

// GetCompetencies lists the staff's competencies.
func (h *UserRoleHandler) GetCompetencies(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetActiveTenantID(c)
	if userID == "" || tenantID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated or branch context missing")
	}

	// Resolve staff ID for the user inside the current branch tenant
	emp, err := h.userRepo.GetTenantEmployment(ctx, userID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Staff record not found for active user")
	}

	list, err := h.userRepo.GetCompetencies(ctx, emp.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list competencies: "+err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

// CreateCompetency adds a staff competency entry.
func (h *UserRoleHandler) CreateCompetency(c echo.Context) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetActiveTenantID(c)
	if userID == "" || tenantID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated or branch context missing")
	}

	emp, err := h.userRepo.GetTenantEmployment(c.Request().Context(), userID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Staff record not found for active user")
	}

	payload := new(modelUser.CreateProfessionalCompetencyPayload)
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	comp, err := h.userRepo.CreateCompetency(c.Request().Context(), emp.ID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save competency: "+err.Error())
	}
	return c.JSON(http.StatusCreated, comp)
}

// GetPrivileges returns list of clinical privileges.
func (h *UserRoleHandler) GetPrivileges(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetActiveTenantID(c)
	if userID == "" || tenantID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated or branch context missing")
	}

	emp, err := h.userRepo.GetTenantEmployment(ctx, userID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Staff record not found")
	}

	list, err := h.userRepo.GetPrivileges(ctx, emp.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load privileges: "+err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

// CreatePrivilege registers clinician privilges.
func (h *UserRoleHandler) CreatePrivilege(c echo.Context) error {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetActiveTenantID(c)
	if userID == "" || tenantID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	// Fetch staff ID to assign privileges to
	emp, err := h.userRepo.GetTenantEmployment(c.Request().Context(), userID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Staff record not found")
	}

	payload := new(modelUser.CreateClinicalPrivilegePayload)
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	// Approved by active actor
	p, err := h.userRepo.CreatePrivilege(c.Request().Context(), emp.ID, payload, emp.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to grant privilege: "+err.Error())
	}
	return c.JSON(http.StatusCreated, p)
}

// VerifyProfessionalProfile logs a PSV verification action.
func (h *UserRoleHandler) VerifyProfessionalProfile(c echo.Context) error {
	actorID := middleware.GetUserID(c)
	profileID := c.Param("id")
	if actorID == "" || profileID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Profile ID and active user session required")
	}

	payload := new(modelUser.CreateProfessionalVerificationPayload)
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	v, err := h.userRepo.CreateVerification(c.Request().Context(), profileID, actorID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to complete verification: "+err.Error())
	}
	return c.JSON(http.StatusCreated, v)
}
