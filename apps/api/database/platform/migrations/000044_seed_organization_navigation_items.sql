-- +goose Up
-- ==============================================================================
-- MIGRATION 000044: Idempotent Organization Navigation Items (DB-Driven Source of Truth)
-- ==============================================================================

-- Remove legacy obsolete navigation items if present
DELETE FROM navigation_item
WHERE context_scope = 'organization'
  AND id IN ('nav_org_workspaces');

-- Seed Comprehensive Organization Navigation Items
INSERT INTO navigation_item (id, context_scope, module_code, title, icon, path, sort_order, required_permission)
VALUES
    ('nav_org_dashboard', 'organization', NULL, 'Executive HQ Dashboard', 'LayoutDashboard', '/organization/dashboard', 1, 'organization:read'),
    ('nav_org_branches', 'organization', NULL, 'Branch Facilities', 'Building2', '/organization/branches', 2, 'organization:branch:read'),
    ('nav_org_members', 'organization', NULL, 'Staff & Members', 'Users', '/organization/members', 3, 'users:read'),
    ('nav_org_roles', 'organization', NULL, 'Roles & Permissions', 'Shield', '/organization/roles', 4, 'organization:manage'),
    ('nav_org_catalogs', 'organization', NULL, 'Catalogs & Pricing', 'BookOpen', '/organization/catalogs', 5, 'organization:catalog:read'),
    ('nav_org_billing', 'organization', NULL, 'Corporate Subscription', 'CreditCard', '/organization/billing', 6, 'organization:manage'),
    ('nav_org_branding', 'organization', NULL, 'Branding & Customization', 'Palette', '/organization/branding', 7, 'organization:branding:read'),
    ('nav_org_notifications', 'organization', NULL, 'Notification Settings', 'Bell', '/organization/notifications', 8, 'organization:notifications:read'),
    ('nav_org_integrations', 'organization', NULL, 'APIs & Webhooks', 'Cpu', '/organization/integrations', 9, 'organization:integrations:read'),
    ('nav_org_audit', 'organization', NULL, 'Corporate Audit Ledger', 'History', '/organization/audit', 10, 'organization:audit:read'),
    ('nav_org_settings', 'organization', NULL, 'Organization Settings', 'Settings', '/organization/settings', 11, 'organization:settings:read')
ON CONFLICT (id) DO UPDATE SET
    context_scope = EXCLUDED.context_scope,
    module_code = EXCLUDED.module_code,
    title = EXCLUDED.title,
    icon = EXCLUDED.icon,
    path = EXCLUDED.path,
    sort_order = EXCLUDED.sort_order,
    required_permission = EXCLUDED.required_permission;

-- +goose Down
DELETE FROM navigation_item
WHERE context_scope = 'organization'
  AND id IN (
    'nav_org_dashboard',
    'nav_org_branches',
    'nav_org_members',
    'nav_org_roles',
    'nav_org_catalogs',
    'nav_org_billing',
    'nav_org_branding',
    'nav_org_notifications',
    'nav_org_integrations',
    'nav_org_audit',
    'nav_org_settings'
  );
