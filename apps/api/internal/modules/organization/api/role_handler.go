package api

import (
	"net/http"
	"strings"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/golangnigeria/curexal/internal/shared/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrganizationRoleHandler struct {
	server *server.Server
}

func NewOrganizationRoleHandler(s *server.Server) *OrganizationRoleHandler {
	return &OrganizationRoleHandler{server: s}
}

type CreateRoleRequest struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type RoleResponseDTO struct {
	ID          string   `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IsSystem    bool     `json:"isSystem"`
	Permissions []string `json:"permissions"`
	MemberCount int      `json:"memberCount"`
}

func (h *OrganizationRoleHandler) CreateRole(c echo.Context) error {
	ctx := c.Request().Context()
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var req CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequestEcho(c, "Invalid role creation payload")
	}

	req.Code = strings.TrimSpace(strings.ToLower(req.Code))
	req.Name = strings.TrimSpace(req.Name)
	if req.Code == "" || req.Name == "" {
		return response.BadRequestEcho(c, "Role name and unique identifier are required")
	}

	dbPool := h.server.DB.Pool
	roleID := uuid.New().String()

	// 1. Insert role
	stmtRole := `
		INSERT INTO "authorization".roles (id, name, code, context_scope, description)
		VALUES ($1, $2, $3, 'organization', $4)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description
		RETURNING id::text
	`
	var actualRoleID string
	err := dbPool.QueryRow(ctx, stmtRole, roleID, req.Name, req.Code, req.Description).Scan(&actualRoleID)
	if err != nil {
		return response.InternalErrorEcho(c, "Failed to persist role: "+err.Error())
	}

	// 2. Insert permission bindings
	if len(req.Permissions) > 0 {
		for _, permCode := range req.Permissions {
			stmtPerm := `
				INSERT INTO "authorization".role_permissions (role_id, permission_id)
				SELECT $1::uuid, p.id
				FROM "authorization".permissions p
				WHERE p.code = $2
				ON CONFLICT DO NOTHING
			`
			_, _ = dbPool.Exec(ctx, stmtPerm, actualRoleID, permCode)
		}
	}

	return response.SuccessEcho(c, http.StatusCreated, RoleResponseDTO{
		ID:          actualRoleID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    false,
		Permissions: req.Permissions,
		MemberCount: 0,
	})
}

func (h *OrganizationRoleHandler) UpdateRole(c echo.Context) error {
	ctx := c.Request().Context()
	roleID := c.Param("id")
	if roleID == "" {
		return response.BadRequestEcho(c, "Role ID is required")
	}

	var req CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequestEcho(c, "Invalid role update payload")
	}

	dbPool := h.server.DB.Pool

	stmtUpdate := `
		UPDATE "authorization".roles
		SET name = COALESCE(NULLIF($1, ''), name),
		    description = COALESCE($2, description)
		WHERE id::text = $3 OR code = $3
	`
	_, err := dbPool.Exec(ctx, stmtUpdate, req.Name, req.Description, roleID)
	if err != nil {
		return response.InternalErrorEcho(c, "Failed to update role: "+err.Error())
	}

	// Sync permissions if provided
	if len(req.Permissions) > 0 {
		_, _ = dbPool.Exec(ctx, `DELETE FROM "authorization".role_permissions WHERE role_id::text = $1`, roleID)
		for _, permCode := range req.Permissions {
			stmtPerm := `
				INSERT INTO "authorization".role_permissions (role_id, permission_id)
				SELECT r.id, p.id
				FROM "authorization".roles r
				CROSS JOIN "authorization".permissions p
				WHERE (r.id::text = $1 OR r.code = $1) AND p.code = $2
				ON CONFLICT DO NOTHING
			`
			_, _ = dbPool.Exec(ctx, stmtPerm, roleID, permCode)
		}
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Role updated successfully",
	})
}

func (h *OrganizationRoleHandler) ListRoles(c echo.Context) error {
	ctx := c.Request().Context()
	dbPool := h.server.DB.Pool

	stmt := `
		SELECT r.id::text, r.code, r.name, COALESCE(r.description, ''),
		       CASE WHEN r.context_scope = 'platform' OR r.code IN ('owner', 'super_admin') THEN TRUE ELSE FALSE END AS is_system,
		       COALESCE(
		           (SELECT array_agg(p.code) 
		            FROM "authorization".role_permissions rp 
		            JOIN "authorization".permissions p ON p.id = rp.permission_id 
		            WHERE rp.role_id = r.id), 
		           '{}'::text[]
		       ) AS permissions,
		       (SELECT COUNT(*) FROM organization.organization_memberships m WHERE m.role_title = r.code OR m.role = r.code) AS member_count
		FROM "authorization".roles r
		ORDER BY is_system DESC, r.name ASC
	`
	rows, err := dbPool.Query(ctx, stmt)
	if err != nil {
		return response.InternalErrorEcho(c, "Failed to query roles: "+err.Error())
	}
	defer rows.Close()

	var roles []RoleResponseDTO
	for rows.Next() {
		var role RoleResponseDTO
		var perms []string
		if errScan := rows.Scan(&role.ID, &role.Code, &role.Name, &role.Description, &role.IsSystem, &perms, &role.MemberCount); errScan == nil {
			role.Permissions = perms
			roles = append(roles, role)
		}
	}

	return response.SuccessEcho(c, http.StatusOK, roles)
}
