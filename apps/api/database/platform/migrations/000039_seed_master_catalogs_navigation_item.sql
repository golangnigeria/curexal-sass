-- +goose Up
-- ==============================================================================
-- MIGRATION 000039: Idempotent Platform Navigation Items & Master Catalogs
-- ==============================================================================

INSERT INTO navigation_item (id, context_scope, module_code, title, icon, path, sort_order, required_permission)
VALUES
    ('nav_plat_dashboard', 'platform', NULL, 'Platform Dashboard', 'LayoutDashboard', '/platform/dashboard', 1, 'platform:view'),
    ('nav_plat_orgs', 'platform', NULL, 'Organizations', 'Building2', '/platform/organizations', 2, 'platform:manage'),
    ('nav_plat_users', 'platform', NULL, 'User Directory', 'Users', '/platform/users', 3, 'platform:manage'),
    ('nav_plat_marketplace', 'platform', NULL, 'B2B Marketplace', 'Store', '/platform/marketplace', 4, 'platform:manage'),
    ('nav_plat_pricing', 'platform', NULL, 'Pricing & Billing', 'CreditCard', '/platform/pricing', 5, 'platform:manage'),
    ('nav_plat_facility_types', 'platform', NULL, 'Facility Types', 'Layers', '/platform/facility-types', 6, 'platform:manage'),
    ('nav_plat_catalogs', 'platform', NULL, 'Master Catalogs', 'BookOpen', '/platform/catalogs', 7, 'platform:manage'),
    ('nav_plat_audit', 'platform', NULL, 'Audit Trail', 'History', '/platform/audit', 8, 'platform:manage'),
    ('nav_plat_diag', 'platform', NULL, 'Diagnostics & Gate', 'Cpu', '/platform/diagnostics', 9, 'platform:manage'),
    ('nav_plat_demo', 'platform', NULL, 'Demo Requests', 'Inbox', '/platform/demo-requests', 10, 'platform:manage'),
    ('nav_plat_settings', 'platform', NULL, 'Console Settings', 'Settings', '/platform/settings', 11, 'platform:manage')
ON CONFLICT (id) DO UPDATE SET
    context_scope = EXCLUDED.context_scope,
    module_code = EXCLUDED.module_code,
    title = EXCLUDED.title,
    icon = EXCLUDED.icon,
    path = EXCLUDED.path,
    sort_order = EXCLUDED.sort_order,
    required_permission = EXCLUDED.required_permission;

-- +goose Down
DELETE FROM navigation_item WHERE id = 'nav_plat_catalogs';
