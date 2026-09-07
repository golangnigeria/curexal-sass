-- +goose Up
-- ==============================================================================
-- MIGRATION 000058: Enhance DB-Driven Navigation System (Canonical Schema & SSOT)
-- ==============================================================================

-- 1. Enhance navigation_item schema with hierarchy, status, capabilities and machine keys
ALTER TABLE navigation_item
    ADD COLUMN IF NOT EXISTS parent_id VARCHAR(100) REFERENCES navigation_item(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS key VARCHAR(100),
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active',
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS is_visible BOOLEAN NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS badge_key VARCHAR(50),
    ADD COLUMN IF NOT EXISTS required_capability VARCHAR(100);

-- Create index on hierarchy and context lookups
CREATE INDEX IF NOT EXISTS idx_nav_item_context_scope ON navigation_item(context_scope);
CREATE INDEX IF NOT EXISTS idx_nav_item_parent_id ON navigation_item(parent_id);
CREATE INDEX IF NOT EXISTS idx_nav_item_key ON navigation_item(key);

-- 2. Seed Canonical Platform Navigation Items (SSOT)
INSERT INTO navigation_item (id, key, context_scope, module_code, title, description, icon, path, sort_order, required_permission, required_capability, status, is_visible, is_active)
VALUES
    ('nav_plat_dashboard', 'platform.dashboard', 'platform', NULL, 'Platform Dashboard', 'System-wide analytics and platform KPIs', 'LayoutDashboard', '/platform/dashboard', 1, 'platform:view', NULL, 'active', true, true),
    ('nav_plat_orgs', 'platform.organizations', 'platform', NULL, 'Organizations', 'Enterprise tenant organizations directory', 'Building2', '/platform/organizations', 2, 'platform:manage', NULL, 'active', true, true),
    ('nav_plat_users', 'platform.users', 'platform', NULL, 'User Directory', 'Platform-wide identities and credential security', 'Users', '/platform/users', 3, 'platform:manage', NULL, 'active', true, true),
    ('nav_plat_marketplace', 'platform.marketplace', 'platform', NULL, 'B2B Marketplace', 'Commercial capability and add-on catalog', 'Store', '/platform/marketplace', 4, 'platform:manage', NULL, 'active', true, true),
    ('nav_plat_pricing', 'platform.pricing', 'platform', NULL, 'Pricing & Billing', 'Subscription tiers and revenue gateways', 'CreditCard', '/platform/pricing', 5, 'platform:manage', NULL, 'active', true, true),
    ('nav_plat_facility_types', 'platform.facility_types', 'platform', NULL, 'Facility Types', 'Clinical blueprint schemas and defaults', 'Layers', '/platform/facility-types', 6, 'platform:manage', NULL, 'active', true, true),
    ('nav_plat_catalogs', 'platform.catalogs', 'platform', NULL, 'Master Catalogs', 'Canonical reference lab and ICD catalogs', 'BookOpen', '/platform/catalogs', 7, 'platform:manage', NULL, 'active', true, true),
    ('nav_plat_audit', 'platform.audit', 'platform', NULL, 'Audit Trail', 'Platform security and administrative audit ledger', 'History', '/platform/audit', 8, 'platform:manage', NULL, 'active', true, true),
    ('nav_plat_diag', 'platform.diagnostics', 'platform', NULL, 'Diagnostics & Gate', 'Kernel health metrics and launch gate verification', 'Cpu', '/platform/diagnostics', 9, 'platform:manage', NULL, 'active', true, true),
    ('nav_plat_demo', 'platform.demo_requests', 'platform', NULL, 'Demo Requests', 'Inbound enterprise prospect requests', 'Inbox', '/platform/demo-requests', 10, 'platform:manage', NULL, 'active', true, true),
    ('nav_plat_settings', 'platform.settings', 'platform', NULL, 'Console Settings', 'Platform control plane configurations', 'Settings', '/platform/settings', 11, 'platform:manage', NULL, 'active', true, true)
ON CONFLICT (id) DO UPDATE SET
    key = EXCLUDED.key,
    context_scope = EXCLUDED.context_scope,
    module_code = EXCLUDED.module_code,
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    icon = EXCLUDED.icon,
    path = EXCLUDED.path,
    sort_order = EXCLUDED.sort_order,
    required_permission = EXCLUDED.required_permission,
    required_capability = EXCLUDED.required_capability,
    status = EXCLUDED.status,
    is_visible = EXCLUDED.is_visible,
    is_active = EXCLUDED.is_active;

-- 3. Seed Canonical Organization HQ Navigation Items (SSOT)
INSERT INTO navigation_item (id, key, context_scope, module_code, title, description, icon, path, sort_order, required_permission, required_capability, status, is_visible, is_active)
VALUES
    ('nav_org_dashboard', 'organization.dashboard', 'organization', NULL, 'Executive HQ Dashboard', 'Corporate KPIs and multi-branch overview', 'LayoutDashboard', '/organization/dashboard', 1, 'organization:read', NULL, 'active', true, true),
    ('nav_org_branches', 'organization.branches', 'organization', NULL, 'Branch Facilities', 'Manage diagnostic and clinical branches', 'Building2', '/organization/branches', 2, 'organization:branch:read', NULL, 'active', true, true),
    ('nav_org_members', 'organization.members', 'organization', NULL, 'Staff & Members', 'Employee roster and organization access control', 'Users', '/organization/members', 3, 'users:read', NULL, 'active', true, true),
    ('nav_org_roles', 'organization.roles', 'organization', NULL, 'Roles & Permissions', 'RBAC permission matrices and custom roles', 'Shield', '/organization/roles', 4, 'organization:manage', NULL, 'active', true, true),
    ('nav_org_catalogs', 'organization.catalogs', 'organization', NULL, 'Catalogs & Pricing', 'Organization fee schedules and lab test catalogs', 'BookOpen', '/organization/catalogs', 5, 'organization:catalog:read', NULL, 'active', true, true),
    ('nav_org_billing', 'organization.billing', 'organization', NULL, 'Corporate Subscription', 'Subscription plan, invoices and payment methods', 'CreditCard', '/organization/billing', 6, 'organization:manage', NULL, 'active', true, true),
    ('nav_org_branding', 'organization.branding', 'organization', NULL, 'Branding & Customization', 'White-label themes, logos and patient portal styling', 'Palette', '/organization/branding', 7, 'organization:branding:read', NULL, 'active', true, true),
    ('nav_org_notifications', 'organization.notifications', 'organization', NULL, 'Notification Settings', 'SMS/Email gateways and automated care triggers', 'Bell', '/organization/notifications', 8, 'organization:notifications:read', NULL, 'active', true, true),
    ('nav_org_integrations', 'organization.integrations', 'organization', NULL, 'APIs & Webhooks', 'Developer API keys, webhooks and HL7/FHIR gateways', 'Cpu', '/organization/integrations', 9, 'organization:integrations:read', NULL, 'active', true, true),
    ('nav_org_audit', 'organization.audit', 'organization', NULL, 'Corporate Audit Ledger', 'Immutable organizational audit and compliance logs', 'History', '/organization/audit', 10, 'organization:audit:read', NULL, 'active', true, true),
    ('nav_org_compliance', 'organization.compliance', 'organization', NULL, 'Compliance & Documents', 'Regulatory accreditation vault and certifications', 'FileCheck', '/organization/compliance', 11, 'organization:compliance:read', NULL, 'active', true, true),
    ('nav_org_settings', 'organization.settings', 'organization', NULL, 'Organization Settings', 'General organization profile and tax policies', 'Settings', '/organization/settings', 12, 'organization:settings:read', NULL, 'active', true, true)
ON CONFLICT (id) DO UPDATE SET
    key = EXCLUDED.key,
    context_scope = EXCLUDED.context_scope,
    module_code = EXCLUDED.module_code,
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    icon = EXCLUDED.icon,
    path = EXCLUDED.path,
    sort_order = EXCLUDED.sort_order,
    required_permission = EXCLUDED.required_permission,
    required_capability = EXCLUDED.required_capability,
    status = EXCLUDED.status,
    is_visible = EXCLUDED.is_visible,
    is_active = EXCLUDED.is_active;

-- 4. Seed Canonical Workspace Operational Navigation Items (SSOT)
INSERT INTO navigation_item (id, key, context_scope, module_code, title, description, icon, path, sort_order, required_permission, required_capability, status, is_visible, is_active)
VALUES
    ('nav_wsp_dashboard', 'workspace.dashboard', 'workspace', 'dashboard', 'Workspace Overview', 'Facility operational metrics and throughput', 'LayoutDashboard', '/:branch/dashboard', 1, NULL, NULL, 'active', true, true),
    ('nav_wsp_reception', 'workspace.reception', 'workspace', 'reception', 'Patient Intake (MPI)', 'Master Patient Index, registration & portal onboarding', 'Users', '/:branch/reception', 2, 'workspace:patient:read', 'core.intake', 'active', true, true),
    ('nav_wsp_care_desk', 'workspace.care_desk', 'workspace', 'care_desk', 'Care & Triage Desk', 'Auto-acuity triage and intelligent doctor allocation', 'Activity', '/:branch/care-desk', 3, 'workspace:care:manage', 'care.orchestration', 'active', true, true),
    ('nav_wsp_clinical', 'workspace.clinical', 'workspace', 'clinical', 'Clinical & EMR', 'Outpatient consultations, SOAP notes and prescriptions', 'Stethoscope', '/:branch/clinical', 4, 'workspace:clinical:read', 'clinical.basic', 'active', true, true),
    ('nav_wsp_laboratory', 'workspace.laboratory', 'workspace', 'laboratory', 'Laboratory (LIS)', 'Specimen accessioning, barcodes & result authorization', 'Microscope', '/:branch/laboratory', 5, 'workspace:sample:receive', 'laboratory.basic', 'active', true, true),
    ('nav_wsp_pharmacy', 'workspace.pharmacy', 'workspace', 'pharmacy', 'Pharmacy Dispensary', 'FEFO batch dispensing and prescription fulfillment', 'Pill', '/:branch/pharmacy', 6, 'workspace:pharmacy:dispense', 'pharmacy.basic', 'active', true, true),
    ('nav_wsp_radiology', 'workspace.radiology', 'workspace', 'radiology', 'Radiology (RIS & PACS)', 'Imaging modality worklists, PACS viewing & reports', 'Layers', '/:branch/radiology', 7, 'workspace:radiology:read', 'radiology.basic', 'active', true, true),
    ('nav_wsp_hospital', 'workspace.hospital', 'workspace', 'hospital', 'Inpatient Hospital (HIS)', 'Ward bed management, admissions & nursing handovers', 'Building2', '/:branch/hospital', 8, 'workspace:hospital:read', 'clinical.inpatient_wards', 'active', true, true),
    ('nav_wsp_billing', 'workspace.billing', 'workspace', 'billing', 'Billing & Cashier POS', 'POS checkout terminal, HMO claims & receipt generation', 'CreditCard', '/:branch/billing', 9, 'workspace:billing:create', 'core.billing', 'active', true, true)
ON CONFLICT (id) DO UPDATE SET
    key = EXCLUDED.key,
    context_scope = EXCLUDED.context_scope,
    module_code = EXCLUDED.module_code,
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    icon = EXCLUDED.icon,
    path = EXCLUDED.path,
    sort_order = EXCLUDED.sort_order,
    required_permission = EXCLUDED.required_permission,
    required_capability = EXCLUDED.required_capability,
    status = EXCLUDED.status,
    is_visible = EXCLUDED.is_visible,
    is_active = EXCLUDED.is_active;

-- 5. Seed Canonical Patient Portal Navigation Items (SSOT)
INSERT INTO navigation_item (id, key, context_scope, module_code, title, description, icon, path, sort_order, required_permission, required_capability, status, is_visible, is_active)
VALUES
    ('nav_pat_dashboard', 'patient.dashboard', 'patient', NULL, 'Care Journey', 'Unified healthcare progression, prescriptions and lab records', 'LayoutDashboard', '/dashboard', 1, NULL, NULL, 'active', true, true),
    ('nav_pat_consultations', 'patient.consultations', 'patient', NULL, 'Book Consultation', 'Request virtual or in-person physician appointment', 'CalendarPlus', '/consultations/new', 2, NULL, NULL, 'active', true, true)
ON CONFLICT (id) DO UPDATE SET
    key = EXCLUDED.key,
    context_scope = EXCLUDED.context_scope,
    module_code = EXCLUDED.module_code,
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    icon = EXCLUDED.icon,
    path = EXCLUDED.path,
    sort_order = EXCLUDED.sort_order,
    required_permission = EXCLUDED.required_permission,
    required_capability = EXCLUDED.required_capability,
    status = EXCLUDED.status,
    is_visible = EXCLUDED.is_visible,
    is_active = EXCLUDED.is_active;

-- +goose Down
-- Revert navigation schema additions
ALTER TABLE navigation_item
    DROP COLUMN IF EXISTS parent_id,
    DROP COLUMN IF EXISTS key,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS is_visible,
    DROP COLUMN IF EXISTS is_active,
    DROP COLUMN IF EXISTS badge_key,
    DROP COLUMN IF EXISTS required_capability;
