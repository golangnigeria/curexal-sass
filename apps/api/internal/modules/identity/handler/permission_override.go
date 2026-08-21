package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// GetPermissionOverrides returns active overrides, inherited permissions, and role context for a membership.
func (h *UserRoleHandler) GetPermissionOverrides(c echo.Context) error {
	ctx := c.Request().Context()
	dbExec := h.server.DB.Conn(ctx)

	membershipID := c.Param("id")
	if _, err := uuid.Parse(membershipID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid membership ID format")
	}

	// 1. Resolve membership user_id and tenant_id from membership table
	var userID, tenantID string
	var roleName string
	err := dbExec.QueryRow(ctx, `
		SELECT m.user_id::text, COALESCE(m.tenant_id::text, ''), m.role_title
		FROM organization.memberships m
		WHERE m.id = $1 AND m.is_active = TRUE
	`, membershipID).Scan(&userID, &tenantID, &roleName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Membership not found or inactive")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error querying membership")
	}

	// 2. Resolve default inherited permissions for this role/tenant
	var inheritedPermissions []string
	permQuery := `
		SELECT COALESCE(
			ARRAY_REMOVE(ARRAY_AGG(DISTINCT p.code), NULL),
			ARRAY[]::text[]
		) AS permissions
		FROM "authorization".roles r
		LEFT JOIN "authorization".role_permissions rp ON rp.role_id = r.id
		LEFT JOIN "authorization".permissions p       ON p.id = rp.permission_id
		WHERE r.name = $1 AND (r.tenant_id = $2 OR r.tenant_id IS NULL)
	`
	err = dbExec.QueryRow(ctx, permQuery, roleName, tenantID).Scan(&inheritedPermissions)
	if err != nil {
		inheritedPermissions = []string{}
	}

	// 3. Resolve active permission overrides from repository
	overrides, err := h.userRepo.GetPermissionOverrides(ctx, userID, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error fetching permission overrides")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"userId":               userID,
		"role":                 roleName,
		"inheritedPermissions": inheritedPermissions,
		"overrides":            overrides,
	})
}

// CreatePermissionOverride assigns a direct grant or deny override to a user membership.
func (h *UserRoleHandler) CreatePermissionOverride(c echo.Context) error {
	ctx := c.Request().Context()
	dbExec := h.server.DB.Conn(ctx)
	actorID := middleware.GetUserID(c)
	actorRole := middleware.GetUserRole(c)

	membershipID := c.Param("id")
	if _, err := uuid.Parse(membershipID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid membership ID format")
	}

	var payload struct {
		Permission      string `json:"permission"`
		OverrideType    string `json:"overrideType"`
		DurationSeconds *int   `json:"durationSeconds"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if payload.Permission == "" || (payload.OverrideType != "grant" && payload.OverrideType != "deny") {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields (permission and overrideType)")
	}

	// 1. Resolve membership user_id and tenant_id
	var targetUserID, tenantID string
	err := dbExec.QueryRow(ctx, `
		SELECT user_id::text, COALESCE(tenant_id::text, '')
		FROM organization.memberships
		WHERE id = $1 AND is_active = TRUE
	`, membershipID).Scan(&targetUserID, &tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Target membership not found or inactive")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error querying membership")
	}

	// 2. Prevent self-elevation or security escalation checks (Privilege Escalation check)
	actorPermissions := middleware.GetPermissions(c)
	hasPermission := false
	for _, p := range actorPermissions {
		if p == payload.Permission {
			hasPermission = true
			break
		}
	}
	rc := middleware.GetRequestContext(c)
	isPlatformAdmin := rc != nil && rc.IsPlatformAdmin
	if !hasPermission && !isPlatformAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Privilege escalation blocked: you cannot assign permissions you do not hold yourself")
	}

	// 3. Resolve expiration if durationSeconds is provided
	var expiresAt *time.Time
	if payload.DurationSeconds != nil && *payload.DurationSeconds > 0 {
		t := time.Now().Add(time.Duration(*payload.DurationSeconds) * time.Second)
		expiresAt = &t
	}

	// 4. Save override to repository
	overrideID, err := h.userRepo.CreatePermissionOverride(
		ctx, targetUserID, tenantID, payload.Permission, payload.OverrideType, actorID, expiresAt,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error creating permission override")
	}

	// 5. Invalidate active permission Redis cache for the target user immediately
	rowsSess, errSess := dbExec.Query(ctx, `SELECT id FROM identity.sessions WHERE user_id = $1 OR user_id = $2`, targetUserID, actorID)
	if errSess == nil {
		for rowsSess.Next() {
			var sessID string
			if errScan := rowsSess.Scan(&sessID); errScan == nil {
				h.server.Redis.Del(ctx, "session_permissions:"+sessID+":"+tenantID)
			}
		}
		rowsSess.Close()
	}

	// 6. Write Audit Trail record (HIPAA compliance)
	reasonStr := fmt.Sprintf("Explicit %s override for permission %s created by %s", payload.OverrideType, payload.Permission, actorRole)
	resType := "permission_override"
	eventCategory := "identity"
	actorName := "Administrator"

	h.server.Logger.Info().
		Str("action", "user.permission.override_created").
		Str("resource_type", resType).
		Str("category", eventCategory).
		Str("actor_name", actorName).
		Str("reason", reasonStr).
		Str("target", targetUserID).
		Msg("Permission override created")

	return c.JSON(http.StatusCreated, map[string]any{
		"id": overrideID,
	})
}

// DeletePermissionOverride revokes an active override immediately.
func (h *UserRoleHandler) DeletePermissionOverride(c echo.Context) error {
	ctx := c.Request().Context()
	dbExec := h.server.DB.Conn(ctx)
	actorID := middleware.GetUserID(c)
	actorRole := middleware.GetUserRole(c)

	membershipID := c.Param("id")
	overrideID := c.Param("overrideId")
	if _, err := uuid.Parse(membershipID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid membership ID format")
	}

	// 1. Resolve membership user_id and tenant_id
	var targetUserID, tenantID string
	err := dbExec.QueryRow(ctx, `
		SELECT user_id::text, COALESCE(tenant_id::text, '')
		FROM organization.memberships
		WHERE id = $1 AND is_active = TRUE
	`, membershipID).Scan(&targetUserID, &tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Target membership not found")
	}

	// 2. Delete override
	err = h.userRepo.DeletePermissionOverride(ctx, overrideID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error deleting permission override")
	}

	// 3. Invalidate Redis cache
	rowsSess, errSess := dbExec.Query(ctx, `SELECT id FROM identity.sessions WHERE user_id = $1`, targetUserID)
	if errSess == nil {
		for rowsSess.Next() {
			var sessID string
			if errScan := rowsSess.Scan(&sessID); errScan == nil {
				h.server.Redis.Del(ctx, "session_permissions:"+sessID+":"+tenantID)
			}
		}
		rowsSess.Close()
	}

	// 4. Audit Trail record
	reasonStr := fmt.Sprintf("Permission override revoked for user %s", targetUserID)
	resType := "permission_override"
	eventCategory := "identity"
	actorName := "Administrator"

	h.server.Logger.Info().
		Str("action", "user.permission.override_deleted").
		Str("actor_id", actorID).
		Str("actor_name", actorName).
		Str("actor_role", actorRole).
		Str("resource_type", resType).
		Str("category", eventCategory).
		Str("reason", reasonStr).
		Str("target", targetUserID).
		Msg("Permission override deleted")

	return c.NoContent(http.StatusNoContent)
}
