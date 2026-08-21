package auth_test

import (
	"context"
	"testing"

	"github.com/golangnigeria/curexal/internal/kernel/auth"
)

func TestPermissionResolver_HasPermission(t *testing.T) {
	resolver := auth.NewMemoryPermissionResolver()
	ctx := context.Background()

	t.Run("nil principal returns false", func(t *testing.T) {
		allowed, err := resolver.HasPermission(ctx, nil, "organization:read")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if allowed {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("super admin bypass grants permission", func(t *testing.T) {
		p := &auth.AuthenticatedPrincipal{
			UserID: "admin_1",
			Platform: auth.PlatformVector{
				IsSuperAdmin: true,
			},
		}
		allowed, err := resolver.HasPermission(ctx, p, "any:custom:permission")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed {
			t.Errorf("expected true for super admin, got false")
		}
	})

	t.Run("platform staff bypass grants permission", func(t *testing.T) {
		p := &auth.AuthenticatedPrincipal{
			UserID: "staff_1",
			Platform: auth.PlatformVector{
				IsPlatformStaff: true,
			},
		}
		allowed, err := resolver.HasPermission(ctx, p, "organization:read")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed {
			t.Errorf("expected true for platform staff, got false")
		}
	})

	t.Run("role based permission granted", func(t *testing.T) {
		p := &auth.AuthenticatedPrincipal{
			UserID: "user_1",
			Role:   "branch_admin",
		}
		allowed, err := resolver.HasPermission(ctx, p, "organization:read")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed {
			t.Errorf("expected true for branch_admin organization:read, got false")
		}
	})

	t.Run("missing permission returns false", func(t *testing.T) {
		p := &auth.AuthenticatedPrincipal{
			UserID: "user_2",
			Role:   "member",
		}
		allowed, err := resolver.HasPermission(ctx, p, "organization:settings:write")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if allowed {
			t.Errorf("expected false for member organization:settings:write, got true")
		}
	})
}

func TestProviderPermissionResolver_MultiScopeUnion(t *testing.T) {
	resolver := auth.NewProviderPermissionResolver()
	ctx := context.Background()

	t.Run("union of organization owner and workspace clinician permissions", func(t *testing.T) {
		p := &auth.AuthenticatedPrincipal{
			UserID: "usr_multi_scope",
			Organization: auth.OrganizationVector{
				OrganizationRole: "owner",
			},
			Workspace: auth.WorkspaceVector{
				WorkspaceRole: "clinician",
			},
		}

		perms, err := resolver.ResolvePermissions(ctx, p)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		permMap := make(map[string]bool)
		for _, perm := range perms {
			permMap[perm] = true
		}

		// Check org owner permission
		if !permMap[auth.PermissionOrganizationSettingsWrite] {
			t.Errorf("expected organization:settings:write from owner role")
		}

		// Check workspace clinician permission
		if !permMap[auth.PermissionPatientView] {
			t.Errorf("expected patient:view from clinician workspace role")
		}
	})

	t.Run("super admin returns explicit all permissions", func(t *testing.T) {
		p := &auth.AuthenticatedPrincipal{
			UserID: "admin_super",
			Platform: auth.PlatformVector{
				IsSuperAdmin: true,
			},
		}

		perms, err := resolver.ResolvePermissions(ctx, p)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(perms) == 0 || perms[0] == "*" {
			t.Errorf("expected explicit permission list without wildcard, got %v", perms)
		}
	})
}

func TestDirectGrantPermissionProvider(t *testing.T) {
	provider := auth.NewDirectGrantPermissionProvider()
	ctx := context.Background()

	p := &auth.AuthenticatedPrincipal{
		UserID:      "usr_direct",
		Permissions: []string{"custom:direct:grant", "reports:export"},
	}

	perms, err := provider.Permissions(ctx, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(perms) != 2 || perms[0] != "custom:direct:grant" || perms[1] != "reports:export" {
		t.Errorf("expected explicit direct grants, got %v", perms)
	}
}

