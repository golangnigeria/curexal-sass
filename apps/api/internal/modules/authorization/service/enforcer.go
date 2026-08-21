package service

import (
	"context"
	"strings"
	"sync"

	"github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/jackc/pgx/v5"
)

type CasbinEnforcer struct {
	server *server.Server
	mu     sync.RWMutex
	// Role to permissions mapping cache
	rolePermissions map[string][]string
}

func NewCasbinEnforcer(s *server.Server) *CasbinEnforcer {
	e := &CasbinEnforcer{
		server:          s,
		rolePermissions: make(map[string][]string),
	}
	e.bootstrapDefaultPolicies()
	return e
}

func (e *CasbinEnforcer) bootstrapDefaultPolicies() {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Platform Roles
	e.rolePermissions["super_admin"] = auth.GetAllPermissions()
	e.rolePermissions["super_support_agent"] = []string{"users:read", "audit:read", "settings:read", "support:impersonate"}
	e.rolePermissions["super_sales_staff"] = []string{"demo:read", "orgs:read", "orgs:write"}

	// Branch / Clinic Roles
	e.rolePermissions["branch_admin"] = []string{
		"users:read", "users:write", "settings:read", "settings:write",
		"patient:view", "patient:create", "patient:update",
		"laboratory:create_order", "laboratory:accession", "laboratory:enter_result", "laboratory:authorize_result",
		"billing:invoice", "billing:refund",
	}
	e.rolePermissions["clinician"] = []string{
		"patient:view", "patient:create", "patient:update",
		"consultation:write", "prescription:write", "laboratory:create_order",
	}
	e.rolePermissions["technician"] = []string{
		"patient:view", "laboratory:accession", "laboratory:enter_result", "laboratory:authorize_result",
	}
	e.rolePermissions["customer_care"] = []string{
		"patient:view", "patient:create", "patient:update", "appointments:write",
	}
	e.rolePermissions["cashier"] = []string{
		"patient:view", "billing:invoice", "billing:payment",
	}
}

// Enforce evaluates whether subject is allowed to perform action on resource within tenant context.
func (e *CasbinEnforcer) Enforce(ctx context.Context, subject, tenant, resource, action string) (bool, string) {
	if subject == "" {
		return false, "subject is empty"
	}

	// 1. Fast O(1) in-memory evaluation via AuthenticatedPrincipal context
	if p := middleware.GetPrincipalFromContext(ctx); p != nil && p.UserID == subject {
		if p.Platform.IsPlatformStaff {
			return true, "platform super admin bypass"
		}
		targetPerm := strings.ToLower(resource + ":" + action)
		for _, userPerm := range p.Permissions {
			if userPerm == "*" || userPerm == targetPerm {
				return true, "granted by in-memory principal permissions"
			}
		}
	}

	// 2. Database Fallback (for CLI / background jobs)
	var isPlatformAdmin bool
	err := e.server.DB.Pool.QueryRow(ctx, `
		SELECT COALESCE(is_platform_admin, FALSE)
		FROM identity.users
		WHERE id = @subject
	`, pgx.NamedArgs{"subject": subject}).Scan(&isPlatformAdmin)
	if err == nil && isPlatformAdmin {
		return true, "platform super admin bypass"
	}

	// Query user effective role within active tenant
	var roleName string
	if tenant != "" {
		err = e.server.DB.Pool.QueryRow(ctx, `
			SELECT COALESCE(m.role_title, 'member')
			FROM organization.organization_memberships m
			LEFT JOIN workspace.workspaces t ON t.organization_id = m.organization_id
			WHERE m.user_id = @subject AND (t.id = @tenant OR t.slug = @tenant OR m.organization_id = @tenant) AND m.is_active = TRUE
			LIMIT 1
		`, pgx.NamedArgs{"subject": subject, "tenant": tenant}).Scan(&roleName)
	}

	if roleName == "" {
		// Fallback to user platform_role
		_ = e.server.DB.Pool.QueryRow(ctx, `
			SELECT COALESCE(platform_role, 'member')
			FROM identity.users
			WHERE id = @subject
		`, pgx.NamedArgs{"subject": subject}).Scan(&roleName)
	}

	e.mu.RLock()
	perms, exists := e.rolePermissions[roleName]
	e.mu.RUnlock()

	if !exists {
		return false, "role has no permissions assigned"
	}

	targetPerm := strings.ToLower(resource + ":" + action)
	for _, p := range perms {
		if p == "*" || p == targetPerm {
			return true, "granted by role: " + roleName
		}
	}

	return false, "denied: missing permission " + targetPerm
}

// ListUserPermissions returns all permission scopes granted to subject in target tenant.
func (e *CasbinEnforcer) ListUserPermissions(ctx context.Context, subject, tenant string) []string {
	if subject == "" {
		return []string{}
	}

	// 1. Query permissions from DB dynamically via role_permission join
	rows, err := e.server.DB.Pool.Query(ctx, `
		SELECT DISTINCT p.code
		FROM "authorization".permissions p
		JOIN "authorization".role_permissions rp ON rp.permission_id = p.id
		JOIN "authorization".roles r ON r.id = rp.role_id
		JOIN organization.organization_memberships m ON (m.role_title = r.code OR m.role = r.code OR m.role_title = r.name)
		WHERE m.user_id = $1 AND (m.organization_id = $2 OR $2 = '' OR $2 IS NULL) AND m.is_active = TRUE
	`, subject, tenant)
	if err == nil {
		defer rows.Close()
		var perms []string
		for rows.Next() {
			var code string
			if err := rows.Scan(&code); err == nil {
				perms = append(perms, code)
			}
		}
		if len(perms) > 0 {
			return perms
		}
	}

	// 2. Fallback to in-memory role lookup
	e.mu.RLock()
	defer e.mu.RUnlock()

	var roleName string
	if tenant != "" {
		_ = e.server.DB.Pool.QueryRow(ctx, `
			SELECT COALESCE(m.role_title, 'member')
			FROM organization.organization_memberships m
			LEFT JOIN workspace.workspaces t ON t.organization_id = m.organization_id
			WHERE m.user_id = @subject AND (t.id = @tenant OR t.slug = @tenant OR m.organization_id = @tenant) AND m.is_active = TRUE
			LIMIT 1
		`, pgx.NamedArgs{"subject": subject, "tenant": tenant}).Scan(&roleName)
	}

	if perms, exists := e.rolePermissions[roleName]; exists {
		return perms
	}
	return []string{
		"workspace:patient:read",
		"workspace:patient:create",
		"workspace:sample:receive",
		"workspace:worksheet:update",
		"workspace:result:authorize",
		"workspace:billing:create",
	}
}
