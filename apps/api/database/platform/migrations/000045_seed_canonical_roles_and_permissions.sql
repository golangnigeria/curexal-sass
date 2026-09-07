-- +goose Up
-- ==============================================================================
-- MIGRATION 000045: Canonical Roles, Permissions & Role-Permission Matrix (DB SSOT)
-- ==============================================================================

-- 1. Ensure required columns on "authorization".roles
ALTER TABLE "authorization".roles ADD COLUMN IF NOT EXISTS code VARCHAR(100);
ALTER TABLE "authorization".roles ADD COLUMN IF NOT EXISTS context_scope VARCHAR(50) DEFAULT 'workspace';
ALTER TABLE "authorization".roles ADD COLUMN IF NOT EXISTS description TEXT;

-- Ensure code lookup uniqueness index
CREATE UNIQUE INDEX IF NOT EXISTS idx_authorization_roles_code ON "authorization".roles(code);

-- 2. Seed All Canonical System Roles (Platform, Organization & Workspace)
INSERT INTO "authorization".roles (id, name, code, context_scope, description)
VALUES
    -- Platform Control Plane Roles
    ('a0000000-0000-0000-0000-000000000001', 'Platform Super Admin', 'super_admin', 'platform', 'Full operational control over the multi-tenant cluster'),
    ('a0000000-0000-0000-0000-000000000002', 'Platform Administrator', 'platform_admin', 'platform', 'Platform administrative governance and tenant onboarding'),
    ('a0000000-0000-0000-0000-000000000003', 'Platform Staff', 'platform_staff', 'platform', 'Platform operational staff and diagnostics observer'),

    -- Organization Control Plane Roles (Executive HQ)
    ('b0000000-0000-0000-0000-000000000001', 'Organization Owner', 'owner', 'organization', 'Owner of the organization healthcare network with full administrative control'),
    ('b0000000-0000-0000-0000-000000000002', 'Organization Administrator', 'org_admin', 'organization', 'Administrative access across all operational branches and staff'),
    ('b0000000-0000-0000-0000-000000000003', 'Regional Operations Manager', 'org_regional_manager', 'organization', 'Oversees operational performance across regional facility clusters'),
    ('b0000000-0000-0000-0000-000000000004', 'Quality & Compliance Manager', 'org_quality_manager', 'organization', 'Enforces clinical compliance, accreditation, and document audits'),
    ('b0000000-0000-0000-0000-000000000005', 'Finance & Billing Manager', 'org_finance_manager', 'organization', 'Corporate financial management, billing POS, and subscriptions'),
    ('b0000000-0000-0000-0000-000000000006', 'Human Resources Manager', 'org_hr_manager', 'organization', 'Staff roster onboarding, credentialing, and access control'),

    -- Workspace / Facility Branch Operational Roles
    ('c0000000-0000-0000-0000-000000000001', 'Branch Facility Administrator', 'branch_admin', 'workspace', 'Facility-level operational management and staff oversight'),
    ('c0000000-0000-0000-0000-000000000002', 'Attending Physician', 'clinician', 'workspace', 'Attending physician / doctor managing clinical patient consultations'),
    ('c0000000-0000-0000-0000-000000000003', 'Laboratory Technician', 'technician', 'workspace', 'Diagnostic technologist performing sample accessioning and testing'),
    ('c0000000-0000-0000-0000-000000000004', 'Front Desk & Reception', 'customer_care', 'workspace', 'Front desk patient reception, registration, and queue coordination'),
    ('c0000000-0000-0000-0000-000000000005', 'Cashier & Billing Officer', 'cashier', 'workspace', 'Point of sale cashier and invoice reconciliation officer'),
    ('c0000000-0000-0000-0000-000000000006', 'Standard Member', 'member', 'workspace', 'Standard workspace operational staff')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    context_scope = EXCLUDED.context_scope,
    description = EXCLUDED.description;

-- 3. Seed Comprehensive System Permissions
INSERT INTO "authorization".permissions (code, module, category, description)
VALUES
    -- Identity & Auth
    ('identity:password:write', 'identity', 'core', 'Change personal and staff passwords'),

    -- Platform Console Permissions
    ('platform:admin', 'platform', 'core', 'Root platform administration'),
    ('platform:health', 'platform', 'core', 'Access system health and diagnostics'),
    ('platform:view', 'platform', 'core', 'View platform console dashboard'),
    ('platform:manage', 'platform', 'core', 'Manage platform-wide resources'),
    ('platform:impersonate', 'platform', 'core', 'Impersonate tenant context for support'),
    ('platform:catalogs:read', 'platform', 'core', 'View platform master catalogs'),
    ('demo:read', 'platform', 'core', 'View inbound platform demo requests'),

    -- Organization Governance Permissions
    ('organization:read', 'organization', 'core', 'Read organization profile and dashboard'),
    ('organization:view', 'organization', 'core', 'View organization overview and analytics'),
    ('organization:manage', 'organization', 'core', 'Full management of organization settings and resources'),
    ('organization:create', 'organization', 'core', 'Create new child organizations and legal entities'),
    ('organization:update', 'organization', 'core', 'Update organization company details and profile'),
    ('organization:delete', 'organization', 'core', 'Decommission organization account'),
    ('organization:settings:read', 'organization', 'core', 'View organization settings and configuration'),
    ('organization:settings:write', 'organization', 'core', 'Update organization settings and configuration'),
    ('organization:branch:read', 'organization', 'operations', 'View organization facility branches'),
    ('organization:branch:write', 'organization', 'operations', 'Create and modify facility branches'),
    ('organization:catalog:read', 'organization', 'catalogs', 'View organization test catalogs and pricing'),
    ('organization:catalog:write', 'organization', 'catalogs', 'Modify organization test catalogs and pricing'),
    ('organization:branding:read', 'organization', 'branding', 'View organization custom white-label branding'),
    ('organization:branding:write', 'organization', 'branding', 'Update organization custom white-label branding'),
    ('organization:notifications:read', 'organization', 'notifications', 'View organization notification settings'),
    ('organization:notifications:write', 'organization', 'notifications', 'Update organization notification settings'),
    ('organization:integrations:read', 'organization', 'security', 'View developer API keys and webhooks'),
    ('organization:integrations:write', 'organization', 'security', 'Create and rotate API keys and webhooks'),
    ('organization:audit:read', 'organization', 'audit', 'Inspect corporate immutable audit logs'),
    ('organization:document:upload', 'organization', 'compliance', 'Upload legal and verification documents'),
    ('organization:document:read', 'organization', 'compliance', 'View uploaded verification documents'),
    ('organization:document:review', 'organization', 'compliance', 'Review pending organization documents'),
    ('organization:document:approve', 'organization', 'compliance', 'Approve organization verification documents'),
    ('organization:document:reject', 'organization', 'compliance', 'Reject organization verification documents'),
    ('organization:verify', 'organization', 'compliance', 'Verify and certify organization accreditation'),

    -- Staff & RBAC Permissions
    ('users:read', 'organization', 'staff', 'View staff directory and member profiles'),
    ('users:write', 'organization', 'staff', 'Invite, modify, and deactivate staff members'),
    ('audit:read', 'audit', 'core', 'View system audit logs'),

    -- Workspace Operational Permissions
    ('workspace:patient:read', 'customer_care', 'clinical', 'Read patient profiles and reception queue'),
    ('workspace:patient:create', 'customer_care', 'clinical', 'Register new patients at reception'),
    ('workspace:sample:receive', 'laboratory', 'clinical', 'Receive and accession specimen samples'),
    ('workspace:worksheet:update', 'laboratory', 'clinical', 'Enter laboratory worksheet results'),
    ('workspace:result:authorize', 'laboratory', 'clinical', 'Authorize and publish lab results'),
    ('workspace:billing:create', 'billing', 'finance', 'Generate invoices and collect POS payments'),
    ('workspace:clinical:read', 'clinical', 'clinical', 'Access patient clinical history and charts'),
    ('workspace:settings:manage', 'workspace', 'core', 'Manage workspace facility configuration'),

    -- Clinical, Diagnostic & Financial Granular Permissions
    ('patient:view', 'patient', 'clinical', 'View detailed patient records'),
    ('patient:create', 'patient', 'clinical', 'Create new patient records'),
    ('patient:update', 'patient', 'clinical', 'Update existing patient records'),
    ('laboratory:create_order', 'laboratory', 'clinical', 'Create laboratory orders'),
    ('laboratory:accession', 'laboratory', 'clinical', 'Accession diagnostic laboratory samples'),
    ('laboratory:enter_result', 'laboratory', 'clinical', 'Input laboratory diagnostic results'),
    ('laboratory:authorize_result', 'laboratory', 'clinical', 'Authorize diagnostic laboratory results'),
    ('billing:read', 'billing', 'finance', 'Read invoices and financial transaction logs'),
    ('billing:write', 'billing', 'finance', 'Create invoices, discounts, and billing adjustments'),
    ('billing:invoice', 'billing', 'finance', 'Issue formal billing statements and invoices'),
    ('billing:payment', 'billing', 'finance', 'Process and record customer payment transactions'),
    ('billing:refund', 'billing', 'finance', 'Process customer refunds and credit notes'),
    ('consultation:write', 'clinical', 'clinical', 'Conduct and record clinical doctor consultations'),
    ('prescription:write', 'clinical', 'clinical', 'Prescribe medications and diagnostic orders'),
    ('appointments:write', 'customer_care', 'clinical', 'Schedule and manage patient appointments'),
    ('settings:read', 'workspace', 'core', 'Read workspace facility configuration'),
    ('settings:write', 'workspace', 'core', 'Update workspace facility configuration')
ON CONFLICT (code) DO UPDATE SET
    module = EXCLUDED.module,
    category = EXCLUDED.category,
    description = EXCLUDED.description;

-- 4. Map Canonical Role-Permission Bindings in "authorization".role_permissions

-- Super Admin: All Permissions
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'super_admin'
ON CONFLICT DO NOTHING;

-- Platform Admin: All platform and view permissions
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code IN ('platform_admin', 'platform_staff')
  AND (p.module = 'platform' OR p.code IN ('identity:password:write', 'organization:read', 'audit:read', 'users:read'))
ON CONFLICT DO NOTHING;

-- Organization Owner: Full Organization Governance, Staff, Catalogs, Billing & Clinical Workspace Access
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'owner'
  AND p.module IN ('organization', 'identity', 'customer_care', 'laboratory', 'clinical', 'billing', 'patient', 'workspace', 'audit')
ON CONFLICT DO NOTHING;

-- Organization Admin: Organization Management, Branches, Staff, Catalogs, Billing, Compliance
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'org_admin'
  AND p.code IN (
    'identity:password:write',
    'organization:read', 'organization:view', 'organization:manage', 'organization:update',
    'organization:settings:read', 'organization:settings:write',
    'organization:branch:read', 'organization:branch:write',
    'organization:catalog:read', 'organization:catalog:write',
    'organization:branding:read', 'organization:branding:write',
    'organization:notifications:read', 'organization:notifications:write',
    'organization:integrations:read', 'organization:integrations:write',
    'organization:audit:read', 'audit:read',
    'organization:document:upload', 'organization:document:read', 'organization:document:review',
    'users:read', 'users:write',
    'billing:read', 'billing:write'
  )
ON CONFLICT DO NOTHING;

-- Regional Manager
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'org_regional_manager'
  AND p.code IN (
    'identity:password:write',
    'organization:read', 'organization:view',
    'organization:branch:read',
    'organization:catalog:read',
    'organization:audit:read', 'audit:read',
    'organization:document:upload', 'organization:document:read',
    'users:read'
  )
ON CONFLICT DO NOTHING;

-- Quality & Compliance Manager
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'org_quality_manager'
  AND p.code IN (
    'identity:password:write',
    'organization:read', 'organization:view',
    'organization:audit:read', 'audit:read',
    'organization:document:upload', 'organization:document:read', 'organization:document:review', 'organization:verify',
    'users:read'
  )
ON CONFLICT DO NOTHING;

-- Finance & Billing Manager
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'org_finance_manager'
  AND p.code IN (
    'identity:password:write',
    'organization:read', 'organization:view',
    'organization:catalog:read', 'organization:catalog:write',
    'billing:read', 'billing:write', 'billing:invoice', 'billing:payment', 'billing:refund',
    'workspace:billing:create'
  )
ON CONFLICT DO NOTHING;

-- HR Manager
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'org_hr_manager'
  AND p.code IN (
    'identity:password:write',
    'organization:read', 'organization:view',
    'users:read', 'users:write'
  )
ON CONFLICT DO NOTHING;

-- Branch Facility Admin
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'branch_admin'
  AND p.code IN (
    'identity:password:write',
    'organization:read',
    'organization:settings:read',
    'organization:branch:read',
    'users:read', 'users:write',
    'patient:view', 'patient:create', 'patient:update',
    'workspace:patient:read', 'workspace:patient:create',
    'laboratory:create_order', 'laboratory:accession',
    'workspace:sample:receive', 'workspace:worksheet:update', 'workspace:result:authorize',
    'workspace:clinical:read', 'workspace:settings:manage',
    'workspace:billing:create', 'billing:read', 'billing:invoice', 'billing:payment',
    'organization:document:upload', 'organization:document:read'
  )
ON CONFLICT DO NOTHING;

-- Clinician (Doctor / Attending Physician)
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'clinician'
  AND p.code IN (
    'identity:password:write',
    'patient:view', 'patient:create', 'patient:update',
    'workspace:patient:read', 'workspace:clinical:read',
    'consultation:write', 'prescription:write', 'appointments:write',
    'laboratory:create_order'
  )
ON CONFLICT DO NOTHING;

-- Technician (Diagnostic Laboratory Technologist)
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'technician'
  AND p.code IN (
    'identity:password:write',
    'patient:view',
    'workspace:patient:read',
    'laboratory:create_order', 'laboratory:accession', 'laboratory:enter_result', 'laboratory:authorize_result',
    'workspace:sample:receive', 'workspace:worksheet:update', 'workspace:result:authorize'
  )
ON CONFLICT DO NOTHING;

-- Front Desk / Customer Care Reception
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'customer_care'
  AND p.code IN (
    'identity:password:write',
    'patient:view', 'patient:create', 'patient:update',
    'workspace:patient:read', 'workspace:patient:create',
    'appointments:write'
  )
ON CONFLICT DO NOTHING;

-- Cashier
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'cashier'
  AND p.code IN (
    'identity:password:write',
    'patient:view',
    'billing:read', 'billing:invoice', 'billing:payment',
    'workspace:billing:create'
  )
ON CONFLICT DO NOTHING;

-- Member
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'member'
  AND p.code IN (
    'identity:password:write',
    'workspace:patient:read'
  )
ON CONFLICT DO NOTHING;

-- 5. Normalize existing organization memberships where role_title is owner/org_admin
UPDATE organization.organization_memberships
SET role = 'owner'
WHERE role_title = 'owner' AND role != 'owner';

UPDATE organization.organization_memberships
SET role = 'org_admin'
WHERE role_title = 'org_admin' AND role != 'org_admin';

-- +goose Down
DELETE FROM "authorization".role_permissions;
