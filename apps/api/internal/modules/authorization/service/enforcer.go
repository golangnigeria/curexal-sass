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

	// Canonical Clinic Roles (Day 1 Section 2.2)
	e.rolePermissions["owner"] = []string{
		"organization:view", "organization:manage", "organization:branch:manage",
		"users:read", "users:write", "audit:read",
		"workspace:patient:create", "workspace:patient:read", "workspace:patient:update",
		"workspace:appointment:read", "workspace:appointment:write", "workspace:queue:manage",
		"workspace:triage:create", "workspace:triage:read",
		"workspace:clinical:read", "workspace:clinical:write", "workspace:clinical:sign",
		"workspace:prescription:write", "workspace:prescription:read",
		"workspace:diagnostic:order", "workspace:diagnostic:read",
		"workspace:document:upload", "workspace:document:read",
		"workspace:billing:read", "workspace:billing:charge", "workspace:billing:refund",
		"workspace:pos:settle", "workspace:pos:shift_close",
	}
	e.rolePermissions["org_admin"] = []string{
		"organization:view", "organization:manage", "organization:branch:manage",
		"users:read", "users:write", "audit:read",
		"workspace:patient:create", "workspace:patient:read", "workspace:patient:update",
		"workspace:appointment:read", "workspace:appointment:write", "workspace:queue:manage",
		"workspace:billing:read", "workspace:billing:charge", "workspace:billing:refund",
	}
	e.rolePermissions["doctor"] = []string{
		"users:read",
		"workspace:patient:create", "workspace:patient:read", "workspace:patient:update",
		"workspace:appointment:read", "workspace:appointment:write", "workspace:queue:manage",
		"workspace:triage:create", "workspace:triage:read",
		"workspace:clinical:read", "workspace:clinical:write", "workspace:clinical:sign",
		"workspace:prescription:write", "workspace:prescription:read",
		"workspace:diagnostic:order", "workspace:diagnostic:read",
		"workspace:document:upload", "workspace:document:read",
	}
	e.rolePermissions["nurse"] = []string{
		"workspace:patient:create", "workspace:patient:read", "workspace:patient:update",
		"workspace:appointment:read", "workspace:appointment:write", "workspace:queue:manage",
		"workspace:triage:create", "workspace:triage:read",
		"workspace:clinical:read",
		"workspace:prescription:read",
		"workspace:diagnostic:read",
		"workspace:document:upload", "workspace:document:read",
	}
	e.rolePermissions["receptionist"] = []string{
		"workspace:patient:create", "workspace:patient:read", "workspace:patient:update",
		"workspace:appointment:read", "workspace:appointment:write", "workspace:queue:manage",
		"workspace:document:upload", "workspace:document:read",
	}
	e.rolePermissions["cashier"] = []string{
		"workspace:patient:read",
		"workspace:prescription:read",
		"workspace:diagnostic:read",
		"workspace:billing:read", "workspace:billing:charge",
		"workspace:pos:settle", "workspace:pos:shift_close",
	}

	// Legacy / Compatibility aliases
	e.rolePermissions["clinician"] = e.rolePermissions["doctor"]
	e.rolePermissions["customer_care"] = e.rolePermissions["receptionist"]
	e.rolePermissions["branch_admin"] = e.rolePermissions["org_admin"]
	e.rolePermissions["technician"] = []string{
		"workspace:patient:read", "workspace:diagnostic:read",
		"laboratory:accession", "laboratory:enter_result", "laboratory:authorize_result",
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
			SELECT COALESCE(NULLIF(m.role, ''), NULLIF(m.role_title, ''), 'member')
			FROM organization.organization_memberships m
			LEFT JOIN organization.organizations o ON o.id = m.organization_id
			LEFT JOIN organization.facility_branches fb ON fb.organization_id = m.organization_id
			WHERE m.user_id = @subject 
			  AND (m.organization_id::text = @tenant OR o.slug = @tenant OR fb.id::text = @tenant OR fb.slug = @tenant) 
			  AND m.is_active = TRUE
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
		JOIN organization.organization_memberships m ON (m.role = r.code OR m.role_title = r.code OR m.role_title = r.name)
		LEFT JOIN organization.organizations o ON o.id = m.organization_id
		LEFT JOIN organization.facility_branches fb ON fb.organization_id = m.organization_id
		WHERE m.user_id = $1 
		  AND ($2 = '' OR $2 IS NULL OR m.organization_id::text = $2 OR o.slug = $2 OR fb.id::text = $2 OR fb.slug = $2) 
		  AND m.is_active = TRUE
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
			SELECT COALESCE(NULLIF(m.role, ''), NULLIF(m.role_title, ''), 'member')
			FROM organization.organization_memberships m
			LEFT JOIN organization.organizations o ON o.id = m.organization_id
			LEFT JOIN organization.facility_branches fb ON fb.organization_id = m.organization_id
			WHERE m.user_id = @subject 
			  AND (m.organization_id::text = @tenant OR o.slug = @tenant OR fb.id::text = @tenant OR fb.slug = @tenant) 
			  AND m.is_active = TRUE
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
