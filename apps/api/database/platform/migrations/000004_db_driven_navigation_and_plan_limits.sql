-- +goose Up
-- ==============================================================================
-- MIGRATION 000004: DB-Driven Navigation Items & Subscription Plan Limits
-- ==============================================================================

CREATE TABLE IF NOT EXISTS navigation_item (
    id VARCHAR(100) PRIMARY KEY,
    context_scope VARCHAR(50) NOT NULL,
    module_code VARCHAR(100),
    title VARCHAR(255) NOT NULL,
    icon VARCHAR(100) NOT NULL,
    path VARCHAR(255) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    required_permission VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Seed Baseline Navigation Items
INSERT INTO navigation_item (id, context_scope, module_code, title, icon, path, sort_order, required_permission)
VALUES
    ('nav_plt_dashboard', 'platform', NULL, 'Platform Overview', 'LayoutDashboard', '/platform/dashboard', 1, 'platform:view'),
    ('nav_plt_organizations', 'platform', NULL, 'Organizations', 'Building2', '/platform/organizations', 2, 'platform:manage'),
    ('nav_plt_users', 'platform', NULL, 'User Directory', 'Users', '/platform/users', 3, 'platform:manage'),
    ('nav_plt_settings', 'platform', NULL, 'System Settings', 'Settings', '/platform/settings', 4, 'platform:manage'),

    ('nav_org_dashboard', 'organization', NULL, 'Organization Dashboard', 'LayoutDashboard', '/organization/dashboard', 1, 'organization:view'),
    ('nav_org_workspaces', 'organization', NULL, 'Facilities & Workspaces', 'Building', '/organization/workspaces', 2, 'organization:manage'),
    ('nav_org_members', 'organization', NULL, 'Team & Access Control', 'UserCheck', '/organization/members', 3, 'organization:manage'),
    ('nav_org_billing', 'organization', NULL, 'Subscription & Billing', 'CreditCard', '/organization/billing', 4, 'organization:manage'),

    ('nav_wsp_dashboard', 'workspace', NULL, 'Workspace Dashboard', 'LayoutDashboard', '/workspace/dashboard', 1, NULL),
    ('nav_wsp_patients', 'workspace', 'customer_care', 'Patient Reception', 'UserPlus', '/workspace/patients', 2, 'workspace:patient:read'),
    ('nav_wsp_laboratory', 'workspace', 'laboratory', 'Laboratory LIS', 'Activity', '/workspace/laboratory/accessioning', 3, 'workspace:sample:receive'),
    ('nav_wsp_clinical', 'workspace', 'clinical', 'Clinical & EMR', 'Stethoscope', '/workspace/clinical/tests', 4, 'workspace:clinical:read'),
    ('nav_wsp_billing', 'workspace', 'billing', 'Billing POS', 'CreditCard', '/workspace/billing', 5, 'workspace:billing:create'),
    ('nav_wsp_settings', 'workspace', NULL, 'Facility Settings', 'Settings', '/workspace/settings', 6, 'workspace:settings:manage')
ON CONFLICT (id) DO NOTHING;

-- Seed Baseline Subscription Plans
INSERT INTO subscription.plans (code, name, limits, features)
VALUES
    ('starter', 'Starter', '{"maxBranches": 2, "maxMembers": 10, "storageGb": 20}'::jsonb, '["customer_care", "billing"]'::jsonb),
    ('pro', 'Pro', '{"maxBranches": 10, "maxMembers": 100, "storageGb": 200}'::jsonb, '["customer_care", "laboratory", "clinical", "billing"]'::jsonb),
    ('enterprise', 'Enterprise', '{"maxBranches": 100, "maxMembers": 1000, "storageGb": 2000}'::jsonb, '["customer_care", "laboratory", "clinical", "pharmacy", "inventory", "radiology", "billing"]'::jsonb)
ON CONFLICT (code) DO NOTHING;

-- Seed Baseline Permissions
INSERT INTO "authorization".permissions (code, module, description)
VALUES
    ('workspace:patient:read', 'customer_care', 'Read patient profile records'),
    ('workspace:patient:create', 'customer_care', 'Register new patients'),
    ('workspace:sample:receive', 'laboratory', 'Receive and accession specimen samples'),
    ('workspace:worksheet:update', 'laboratory', 'Enter laboratory worksheet results'),
    ('workspace:result:authorize', 'laboratory', 'Authorize lab results'),
    ('workspace:billing:create', 'billing', 'Generate invoices and collect payments'),
    ('workspace:clinical:read', 'clinical', 'Access patient clinical history'),
    ('workspace:settings:manage', 'workspace', 'Manage workspace facility configuration')
ON CONFLICT (code) DO NOTHING;

-- Seed Default Member Role Permissions
INSERT INTO "authorization".roles (name, code, context_scope, description)
VALUES
    ('Member', 'member', 'workspace', 'Standard workspace operational staff')
ON CONFLICT (name) DO NOTHING;

INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'member'
ON CONFLICT DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS navigation_item;

