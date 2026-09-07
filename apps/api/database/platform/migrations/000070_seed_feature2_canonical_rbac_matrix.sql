-- +goose Up
-- ==============================================================================
-- CUREXAL MIGRATION 000070: FEATURE #2 CANONICAL RBAC ROLES & PERMISSION MATRIX
-- ==============================================================================

-- 0. Reconcile legacy platform_support role code to canonical super_support_agent
UPDATE identity.users
SET platform_role = 'super_support_agent'
WHERE platform_role = 'platform_support';

-- If super_support_agent already exists, remove legacy platform_support to prevent name collision on roles_name_key
DELETE FROM "authorization".roles
WHERE code = 'platform_support'
  AND EXISTS (SELECT 1 FROM "authorization".roles WHERE code = 'super_support_agent');

-- Migrate legacy platform_support code to canonical super_support_agent
UPDATE "authorization".roles
SET code = 'super_support_agent',
    name = 'Platform Support Specialist',
    context_scope = 'platform',
    description = 'Platform support agent for diagnostics and troubleshooting'
WHERE code = 'platform_support' OR (name = 'Platform Support Specialist' AND code != 'super_support_agent');

-- 1. Ensure Canonical Clinic & Platform Roles Exist
INSERT INTO "authorization".roles (name, code, context_scope, description)
VALUES
    -- Canonical Platform Roles
    ('Platform Super Admin', 'super_admin', 'platform', 'Root operational authority over platform infrastructure'),
    ('Platform Support Specialist', 'super_support_agent', 'platform', 'Platform support agent for diagnostics and troubleshooting'),
    ('Platform Sales Specialist', 'super_sales_staff', 'platform', 'Platform commercial and demo onboarding specialist'),

    -- Canonical Clinic Roles (Day 1 Section 2.2 A)
    ('Clinic Owner / Medical Director', 'owner', 'organization', 'Full legal, organizational, and financial authority'),
    ('Practice Manager / Clinic Administrator', 'org_admin', 'organization', 'Manages branches, staff invites, catalogs, and schedules'),
    ('Attending Medical Doctor / Specialist', 'doctor', 'workspace', 'Consultation canvas, SOAP notes, ICD-10 coding, e-prescriptions, diagnostic orders'),
    ('Triage Nurse / Clinical Assistant', 'nurse', 'workspace', 'Patient queue check-in, triage intake, vital signs, allergy tagging'),
    ('Front Desk Officer', 'receptionist', 'workspace', 'Master Patient Index (MPI) search, walk-in registration, appointment booking'),
    ('Billing & Accounts Officer', 'cashier', 'workspace', 'Service fee settlement, POS card/cash/transfer collections, receipt generation, daily shift reconciliation')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    context_scope = EXCLUDED.context_scope,
    description = EXCLUDED.description;

-- 2. Ensure All Canonical Permissions Exist (Day 1 Section 2.2 B)
INSERT INTO "authorization".permissions (code, category, module, description)
VALUES
    -- Clinic Workspace Governance & Identity
    ('organization:view', 'core', 'organization', 'View clinic profile, branches, and analytics'),
    ('organization:manage', 'core', 'organization', 'Edit clinic settings, business hours, and VAT rules'),
    ('organization:branch:manage', 'core', 'organization', 'Provision and manage clinic branch locations'),
    ('users:read', 'core', 'organization', 'View staff roster and provider availability'),
    ('users:write', 'core', 'organization', 'Invite, assign roles, and deactivate staff'),
    ('audit:read', 'core', 'audit', 'Inspect immutable audit and compliance trail'),

    -- Patient & MPI
    ('workspace:patient:create', 'clinical', 'patient', 'Register new patients and assign MRN'),
    ('workspace:patient:read', 'clinical', 'patient', 'Search Master Patient Index (MPI)'),
    ('workspace:patient:update', 'clinical', 'patient', 'Update patient demographics and contacts'),

    -- Appointments & Operations
    ('workspace:appointment:read', 'operations', 'customer_care', 'View doctor calendars and booked appointments'),
    ('workspace:appointment:write', 'operations', 'customer_care', 'Book, reschedule, and cancel appointments'),
    ('workspace:queue:manage', 'operations', 'customer_care', 'Patient check-in and queue status management'),

    -- Clinical & Triage
    ('workspace:triage:create', 'clinical', 'clinical', 'Record vitals (BP, Pulse, Temp, SpO2, Weight)'),
    ('workspace:triage:read', 'clinical', 'clinical', 'View historical vitals and trend charts'),
    ('workspace:clinical:read', 'clinical', 'clinical', 'View past visit history, allergies, and care plans'),
    ('workspace:clinical:write', 'clinical', 'clinical', 'Document SOAP notes, exam findings, and ICD-10'),
    ('workspace:clinical:sign', 'clinical', 'clinical', 'Digitally sign and permanently lock encounter'),
    ('workspace:prescription:write', 'clinical', 'clinical', 'Prescribe medications (drug, dose, duration)'),
    ('workspace:prescription:read', 'clinical', 'clinical', 'Cashier bills drugs; nurse administers meds'),
    ('workspace:diagnostic:order', 'clinical', 'clinical', 'Order laboratory tests and imaging procedures'),
    ('workspace:diagnostic:read', 'clinical', 'clinical', 'Cashier bills tests; clinicians review results'),

    -- Documents
    ('workspace:document:upload', 'clinical', 'documents', 'Upload referral letters, paper scans, and IDs'),
    ('workspace:document:read', 'clinical', 'documents', 'Inspect uploaded documents and test PDFs'),

    -- Financial, POS & Billing
    ('workspace:billing:read', 'financial', 'billing', 'Access invoice ledgers and balance statements'),
    ('workspace:billing:charge', 'financial', 'billing', 'Generate itemized invoices from consults/tests'),
    ('workspace:billing:refund', 'financial', 'billing', 'Authorize payment reversals and credit notes'),
    ('workspace:pos:settle', 'financial', 'billing', 'Settle invoices via Cash, POS Card, or Transfer'),
    ('workspace:pos:shift_close', 'financial', 'billing', 'Reconcile cash drawer and close cashier shift'),

    -- Platform Super-Admin
    ('platform:admin', 'platform', 'platform', 'Root authority over platform, schemas, and keys'),
    ('platform:view', 'platform', 'platform', 'Access platform dashboard metrics & clinic counts'),
    ('platform:manage', 'platform', 'platform', 'Global settings, feature flags, and tenant lifecycle'),
    ('platform:impersonate', 'platform', 'platform', 'Break-glass support access to debug clinic issues'),
    ('platform:catalogs:manage', 'platform', 'platform', 'Maintain global ICD-10 registry & standard catalogs'),
    ('platform:pricing:manage', 'platform', 'platform', 'Configure subscription rates and gateway vault keys'),
    ('demo:manage', 'platform', 'platform', 'Process inbound demo requests & provision pilots')
ON CONFLICT (code) DO UPDATE SET
    category = EXCLUDED.category,
    module = EXCLUDED.module,
    description = EXCLUDED.description;

-- 3. Populate Canonical Role-Permission Mappings Matrix
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM (
    VALUES
        -- 1. OWNER: Full authority across all 27 clinic workspace permissions
        ('owner', 'organization:view'),
        ('owner', 'organization:manage'),
        ('owner', 'organization:branch:manage'),
        ('owner', 'users:read'),
        ('owner', 'users:write'),
        ('owner', 'audit:read'),
        ('owner', 'workspace:patient:create'),
        ('owner', 'workspace:patient:read'),
        ('owner', 'workspace:patient:update'),
        ('owner', 'workspace:appointment:read'),
        ('owner', 'workspace:appointment:write'),
        ('owner', 'workspace:queue:manage'),
        ('owner', 'workspace:triage:create'),
        ('owner', 'workspace:triage:read'),
        ('owner', 'workspace:clinical:read'),
        ('owner', 'workspace:clinical:write'),
        ('owner', 'workspace:clinical:sign'),
        ('owner', 'workspace:prescription:write'),
        ('owner', 'workspace:prescription:read'),
        ('owner', 'workspace:diagnostic:order'),
        ('owner', 'workspace:diagnostic:read'),
        ('owner', 'workspace:document:upload'),
        ('owner', 'workspace:document:read'),
        ('owner', 'workspace:billing:read'),
        ('owner', 'workspace:billing:charge'),
        ('owner', 'workspace:billing:refund'),
        ('owner', 'workspace:pos:settle'),
        ('owner', 'workspace:pos:shift_close'),

        -- 2. ORG_ADMIN: Practice Manager / Clinic Administrator
        ('org_admin', 'organization:view'),
        ('org_admin', 'organization:manage'),
        ('org_admin', 'organization:branch:manage'),
        ('org_admin', 'users:read'),
        ('org_admin', 'users:write'),
        ('org_admin', 'audit:read'),
        ('org_admin', 'workspace:patient:create'),
        ('org_admin', 'workspace:patient:read'),
        ('org_admin', 'workspace:patient:update'),
        ('org_admin', 'workspace:appointment:read'),
        ('org_admin', 'workspace:appointment:write'),
        ('org_admin', 'workspace:queue:manage'),
        ('org_admin', 'workspace:billing:read'),
        ('org_admin', 'workspace:billing:charge'),
        ('org_admin', 'workspace:billing:refund'),

        -- 3. DOCTOR: Attending Medical Doctor / Specialist
        ('doctor', 'users:read'),
        ('doctor', 'workspace:patient:create'),
        ('doctor', 'workspace:patient:read'),
        ('doctor', 'workspace:patient:update'),
        ('doctor', 'workspace:appointment:read'),
        ('doctor', 'workspace:appointment:write'),
        ('doctor', 'workspace:queue:manage'),
        ('doctor', 'workspace:triage:create'),
        ('doctor', 'workspace:triage:read'),
        ('doctor', 'workspace:clinical:read'),
        ('doctor', 'workspace:clinical:write'),
        ('doctor', 'workspace:clinical:sign'),
        ('doctor', 'workspace:prescription:write'),
        ('doctor', 'workspace:prescription:read'),
        ('doctor', 'workspace:diagnostic:order'),
        ('doctor', 'workspace:diagnostic:read'),
        ('doctor', 'workspace:document:upload'),
        ('doctor', 'workspace:document:read'),

        -- 4. NURSE: Triage Nurse / Clinical Assistant
        ('nurse', 'workspace:patient:create'),
        ('nurse', 'workspace:patient:read'),
        ('nurse', 'workspace:patient:update'),
        ('nurse', 'workspace:appointment:read'),
        ('nurse', 'workspace:appointment:write'),
        ('nurse', 'workspace:queue:manage'),
        ('nurse', 'workspace:triage:create'),
        ('nurse', 'workspace:triage:read'),
        ('nurse', 'workspace:clinical:read'),
        ('nurse', 'workspace:prescription:read'),
        ('nurse', 'workspace:diagnostic:read'),
        ('nurse', 'workspace:document:upload'),
        ('nurse', 'workspace:document:read'),

        -- 5. RECEPTIONIST: Front Desk Officer
        ('receptionist', 'workspace:patient:create'),
        ('receptionist', 'workspace:patient:read'),
        ('receptionist', 'workspace:patient:update'),
        ('receptionist', 'workspace:appointment:read'),
        ('receptionist', 'workspace:appointment:write'),
        ('receptionist', 'workspace:queue:manage'),
        ('receptionist', 'workspace:document:upload'),
        ('receptionist', 'workspace:document:read'),

        -- 6. CASHIER: Billing & Accounts Officer
        ('cashier', 'workspace:patient:read'),
        ('cashier', 'workspace:prescription:read'),
        ('cashier', 'workspace:diagnostic:read'),
        ('cashier', 'workspace:billing:read'),
        ('cashier', 'workspace:billing:charge'),
        ('cashier', 'workspace:pos:settle'),
        ('cashier', 'workspace:pos:shift_close'),

        -- Platform Roles
        ('super_admin', 'platform:admin'),
        ('super_admin', 'platform:view'),
        ('super_admin', 'platform:manage'),
        ('super_admin', 'platform:impersonate'),
        ('super_admin', 'platform:catalogs:manage'),
        ('super_admin', 'platform:pricing:manage'),
        ('super_admin', 'demo:manage'),

        ('super_support_agent', 'platform:view'),
        ('super_support_agent', 'platform:impersonate'),

        ('super_sales_staff', 'platform:view'),
        ('super_sales_staff', 'demo:manage')
) AS matrix(role_code, perm_code)
JOIN "authorization".roles r ON r.code = matrix.role_code
JOIN "authorization".permissions p ON p.code = matrix.perm_code
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- +goose Down
DELETE FROM "authorization".role_permissions rp
USING "authorization".permissions p
WHERE rp.permission_id = p.id
  AND p.code IN (
    'organization:view', 'organization:manage', 'organization:branch:manage', 'users:read', 'users:write', 'audit:read',
    'workspace:patient:create', 'workspace:patient:read', 'workspace:patient:update',
    'workspace:appointment:read', 'workspace:appointment:write', 'workspace:queue:manage',
    'workspace:triage:create', 'workspace:triage:read', 'workspace:clinical:read', 'workspace:clinical:write', 'workspace:clinical:sign',
    'workspace:prescription:write', 'workspace:prescription:read', 'workspace:diagnostic:order', 'workspace:diagnostic:read',
    'workspace:document:upload', 'workspace:document:read',
    'workspace:billing:read', 'workspace:billing:charge', 'workspace:billing:refund', 'workspace:pos:settle', 'workspace:pos:shift_close',
    'platform:admin', 'platform:view', 'platform:manage', 'platform:impersonate', 'platform:catalogs:manage', 'platform:pricing:manage', 'demo:manage'
  );
