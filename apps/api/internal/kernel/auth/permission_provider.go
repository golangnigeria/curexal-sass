package auth

import (
	"context"
)

// PermissionProvider defines a pluggable source of permissions for an authenticated principal.
type PermissionProvider interface {
	Name() string
	Permissions(ctx context.Context, principal *AuthenticatedPrincipal) ([]string, error)
}

// MemoryRoleMapPermissionProvider resolves permissions using a role-to-permissions map for specific principal extractors.
type MemoryRoleMapPermissionProvider struct {
	name            string
	rolePermissions map[string][]string
	roleExtractor   func(principal *AuthenticatedPrincipal) []string
}

func NewPlatformPermissionProvider(rolePermissions map[string][]string) PermissionProvider {
	return &MemoryRoleMapPermissionProvider{
		name:            "PlatformPermissionProvider",
		rolePermissions: rolePermissions,
		roleExtractor: func(principal *AuthenticatedPrincipal) []string {
			var roles []string
			if principal.Platform.PlatformRole != "" {
				roles = append(roles, principal.Platform.PlatformRole)
			}
			return roles
		},
	}
}

func NewOrganizationPermissionProvider(rolePermissions map[string][]string) PermissionProvider {
	return &MemoryRoleMapPermissionProvider{
		name:            "OrganizationPermissionProvider",
		rolePermissions: rolePermissions,
		roleExtractor: func(principal *AuthenticatedPrincipal) []string {
			var roles []string
			if principal.Organization.OrganizationRole != "" {
				roles = append(roles, principal.Organization.OrganizationRole)
			}
			if principal.Role != "" {
				roles = append(roles, principal.Role)
			}
			if len(roles) == 0 && principal.UserID != "" {
				roles = append(roles, "owner", "member")
			}
			return roles
		},
	}
}

func NewWorkspacePermissionProvider(rolePermissions map[string][]string) PermissionProvider {
	return &MemoryRoleMapPermissionProvider{
		name:            "WorkspacePermissionProvider",
		rolePermissions: rolePermissions,
		roleExtractor: func(principal *AuthenticatedPrincipal) []string {
			var roles []string
			if principal.Workspace.WorkspaceRole != "" {
				roles = append(roles, principal.Workspace.WorkspaceRole)
			}
			if principal.Role != "" {
				roles = append(roles, principal.Role)
			}
			return roles
		},
	}
}

func (p *MemoryRoleMapPermissionProvider) Name() string {
	return p.name
}

func (p *MemoryRoleMapPermissionProvider) Permissions(ctx context.Context, principal *AuthenticatedPrincipal) ([]string, error) {
	if principal == nil {
		return nil, nil
	}

	roles := p.roleExtractor(principal)
	var perms []string
	seen := make(map[string]bool)

	for _, role := range roles {
		if rolePerms, exists := p.rolePermissions[role]; exists {
			for _, perm := range rolePerms {
				if !seen[perm] {
					seen[perm] = true
					perms = append(perms, perm)
				}
			}
		}
	}

	return perms, nil
}

// DirectGrantPermissionProvider resolves explicit direct user permission grants/overrides.
type DirectGrantPermissionProvider struct {
	name string
}

func NewDirectGrantPermissionProvider() PermissionProvider {
	return &DirectGrantPermissionProvider{
		name: "DirectGrantPermissionProvider",
	}
}

func (p *DirectGrantPermissionProvider) Name() string {
	return p.name
}

func (p *DirectGrantPermissionProvider) Permissions(ctx context.Context, principal *AuthenticatedPrincipal) ([]string, error) {
	if principal == nil {
		return nil, nil
	}
	return principal.Permissions, nil
}

// ProviderPermissionResolver aggregates permissions across multiple registered providers.
type ProviderPermissionResolver struct {
	providers []PermissionProvider
}

func NewProviderPermissionResolver(providers ...PermissionProvider) *ProviderPermissionResolver {
	if len(providers) == 0 {
		mem := NewMemoryPermissionResolver()
		providers = []PermissionProvider{
			NewPlatformPermissionProvider(mem.rolePermissions),
			NewOrganizationPermissionProvider(mem.rolePermissions),
			NewWorkspacePermissionProvider(mem.rolePermissions),
			NewDirectGrantPermissionProvider(),
		}
	}
	return &ProviderPermissionResolver{
		providers: providers,
	}
}

func NewDefaultProviderPermissionResolver(rolePermissions map[string][]string) *ProviderPermissionResolver {
	if rolePermissions == nil {
		rolePermissions = NewMemoryPermissionResolver().rolePermissions
	}
	return &ProviderPermissionResolver{
		providers: []PermissionProvider{
			NewPlatformPermissionProvider(rolePermissions),
			NewOrganizationPermissionProvider(rolePermissions),
			NewWorkspacePermissionProvider(rolePermissions),
			NewDirectGrantPermissionProvider(),
		},
	}
}

func NewDatabaseProviderPermissionResolver(finder RolePermissionsFinder) *ProviderPermissionResolver {
	return &ProviderPermissionResolver{
		providers: []PermissionProvider{
			NewDatabaseRolePermissionProvider(finder),
			NewDirectGrantPermissionProvider(),
		},
	}
}

func (r *ProviderPermissionResolver) ResolvePermissions(ctx context.Context, principal *AuthenticatedPrincipal) ([]string, error) {
	if principal == nil {
		return nil, nil
	}
	if principal.Platform.IsSuperAdmin || principal.Platform.IsPlatformStaff || principal.Platform.IsPlatformAdmin {
		return GetAllPermissions(), nil
	}

	seen := make(map[string]bool)
	var effective []string

	for _, provider := range r.providers {
		perms, err := provider.Permissions(ctx, principal)
		if err != nil {
			continue
		}
		for _, p := range perms {
			if p == "*" {
				for _, allP := range GetAllPermissions() {
					if !seen[allP] {
						seen[allP] = true
						effective = append(effective, allP)
					}
				}
				continue
			}
			if !seen[p] {
				seen[p] = true
				effective = append(effective, p)
			}
		}
	}
	return effective, nil
}

func (r *ProviderPermissionResolver) HasPermission(ctx context.Context, principal *AuthenticatedPrincipal, permission string) (bool, error) {
	if principal == nil {
		return false, nil
	}
	if principal.Platform.IsSuperAdmin || principal.Platform.IsPlatformStaff || principal.Platform.IsPlatformAdmin {
		return true, nil
	}

	if principal.HasPermission(permission) {
		return true, nil
	}

	perms, err := r.ResolvePermissions(ctx, principal)
	if err != nil {
		return false, err
	}

	for _, p := range perms {
		if p == "*" || p == permission {
			return true, nil
		}
	}
	return false, nil
}

type RolePermissionsFinder interface {
	ListPermissionsByRole(ctx context.Context, roleName, tenantID string) ([]string, error)
}

type DatabaseRolePermissionProvider struct {
	finder RolePermissionsFinder
}

func NewDatabaseRolePermissionProvider(finder RolePermissionsFinder) PermissionProvider {
	return &DatabaseRolePermissionProvider{finder: finder}
}

func (p *DatabaseRolePermissionProvider) Name() string {
	return "DatabaseRolePermissionProvider"
}

func (p *DatabaseRolePermissionProvider) Scope() string {
	return "database"
}

func (p *DatabaseRolePermissionProvider) Permissions(ctx context.Context, principal *AuthenticatedPrincipal) ([]string, error) {
	if principal == nil || p.finder == nil {
		return nil, nil
	}
	var roles []string
	if principal.Role != "" {
		roles = append(roles, principal.Role)
	}
	if principal.Organization.OrganizationRole != "" && principal.Organization.OrganizationRole != principal.Role {
		roles = append(roles, principal.Organization.OrganizationRole)
	}
	if principal.Platform.PlatformRole != "" && principal.Platform.PlatformRole != principal.Role {
		roles = append(roles, principal.Platform.PlatformRole)
	}

	var all []string
	seen := make(map[string]bool)
	for _, r := range roles {
		perms, err := p.finder.ListPermissionsByRole(ctx, r, principal.TenantID)
		if err == nil {
			for _, perm := range perms {
				if !seen[perm] {
					seen[perm] = true
					all = append(all, perm)
				}
			}
		}
	}
	return all, nil
}

