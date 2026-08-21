-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: CAPABILITY COMPOSITION & ENTITLEMENTS SCHEMAS
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS subscription;

-- 1. Capabilities Catalog
CREATE TABLE IF NOT EXISTS subscription.capabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    module VARCHAR(100) NOT NULL,
    tier_level VARCHAR(50) NOT NULL DEFAULT 'core',
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 2. Plan Base Capabilities
CREATE TABLE IF NOT EXISTS subscription.plan_capabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES subscription.plans(id) ON DELETE CASCADE,
    capability_id UUID NOT NULL REFERENCES subscription.capabilities(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_plan_capability UNIQUE (plan_id, capability_id)
);

-- 3. Capability Dependencies
CREATE TABLE IF NOT EXISTS subscription.capability_dependencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    capability_id UUID NOT NULL REFERENCES subscription.capabilities(id) ON DELETE CASCADE,
    depends_on_capability_id UUID NOT NULL REFERENCES subscription.capabilities(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_capability_dependency UNIQUE (capability_id, depends_on_capability_id)
);

-- 4. Organization Add-on Entitlements
CREATE TABLE IF NOT EXISTS subscription.organization_entitlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    capability_id UUID NOT NULL REFERENCES subscription.capabilities(id) ON DELETE CASCADE,
    source VARCHAR(50) NOT NULL DEFAULT 'purchase',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    starts_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    granted_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_org_capability_entitlement UNIQUE (organization_id, capability_id, source)
);

-- Seed Core Platform Capabilities
INSERT INTO subscription.capabilities (code, module, tier_level, name, description)
VALUES
    ('core.organization', 'organization', 'core', 'Organization Management', 'Base organization, branch and department management'),
    ('core.patient', 'patient', 'core', 'Patient Foundation', 'Patient registration, search, MRN and demographics'),
    ('core.customer_care', 'customer_care', 'core', 'Customer Care & Queue', 'Reception, queue management, appointments and referrals'),
    ('core.billing', 'billing', 'core', 'Billing & Financials', 'Service pricing, invoicing, payment recording and receipts'),
    ('core.audit', 'audit', 'core', 'Audit & Compliance', 'System audit logs and compliance history'),
    ('core.documents', 'documents', 'core', 'Document Management', 'Document uploads, status tracking and verification'),
    ('laboratory.basic', 'laboratory', 'core', 'Basic Laboratory Operations', 'Test catalog, lab orders, sample collection and basic result entry'),
    ('radiology.basic', 'radiology', 'core', 'Basic Radiology Operations', 'Procedure catalog, radiology requests, scheduling and PDF reports'),
    ('clinical.basic', 'clinical', 'core', 'Basic Clinical Consultations', 'Encounters, chief complaint, vitals, clinical notes and prescriptions'),
    ('pharmacy.basic', 'pharmacy', 'core', 'Basic Pharmacy Operations', 'Drug catalog, prescription dispensing and basic stock'),
    ('inventory.basic', 'inventory', 'core', 'Basic Inventory Management', 'Item catalog, stock levels and low-stock alerts'),
    ('qms.basic', 'qms', 'core', 'Basic Quality Management', 'SOP repository, document management and incident recording'),
    
    -- Advanced Add-on Capabilities
    ('laboratory.analyzer_integration', 'laboratory', 'addon', 'Analyzer & Middleware Integration', 'ASTM/HL7 instrument interfaces and analyzer worklists'),
    ('laboratory.advanced_qc', 'laboratory', 'addon', 'Advanced Quality Control', 'Levey-Jennings charts, Westgard rules and control material tracking'),
    ('laboratory.pathology_microbiology', 'laboratory', 'addon', 'Pathology & Microbiology', 'Culture/sensitivity, molecular results and histopathology'),
    ('radiology.pacs_dicom', 'radiology', 'addon', 'PACS & DICOM Integration', 'DICOM worklists, PACS image routing and viewer integration'),
    ('radiology.voice_dictation', 'radiology', 'addon', 'Voice Dictation & Structured Reporting', 'Radiologist voice dictation and structured reporting macros'),
    ('clinical.inpatient_wards', 'clinical', 'addon', 'Inpatient & Ward Management', 'Bed management, ward transfers, inpatient billing and nursing notes'),
    ('pharmacy.advanced_inventory', 'pharmacy', 'addon', 'Advanced Pharmacy & FEFO', 'Purchase orders, batch management, FEFO and supplier management'),
    ('inventory.multi_warehouse', 'inventory', 'addon', 'Multi-Warehouse Inventory', 'Multi-location warehouses, lot/serial tracking and reorder automation'),
    ('qms.iso15189', 'qms', 'addon', 'ISO 15189 Quality Workflows', 'ISO 15189 compliance, CAPA, internal/external audits and risk management')
ON CONFLICT (code) DO NOTHING;

-- Map Smart Plan Base Capabilities
INSERT INTO subscription.plan_capabilities (plan_id, capability_id)
SELECT p.id, c.id
FROM subscription.plans p
CROSS JOIN subscription.capabilities c
WHERE p.code = 'smart'
  AND c.code IN ('core.organization', 'core.patient', 'core.customer_care', 'core.billing', 'core.audit', 'core.documents', 'laboratory.basic', 'radiology.basic', 'clinical.basic', 'pharmacy.basic', 'inventory.basic', 'qms.basic')
ON CONFLICT DO NOTHING;

-- Map Optimize Plan Base Capabilities
INSERT INTO subscription.plan_capabilities (plan_id, capability_id)
SELECT p.id, c.id
FROM subscription.plans p
CROSS JOIN subscription.capabilities c
WHERE p.code = 'optimize'
  AND c.code IN ('core.organization', 'core.patient', 'core.customer_care', 'core.billing', 'core.audit', 'core.documents', 'laboratory.basic', 'radiology.basic', 'clinical.basic', 'pharmacy.basic', 'inventory.basic', 'qms.basic', 'laboratory.analyzer_integration', 'laboratory.advanced_qc')
ON CONFLICT DO NOTHING;

-- Map Pro Plan Base Capabilities
INSERT INTO subscription.plan_capabilities (plan_id, capability_id)
SELECT p.id, c.id
FROM subscription.plans p
CROSS JOIN subscription.capabilities c
WHERE p.code = 'pro'
  AND c.code IN ('core.organization', 'core.patient', 'core.customer_care', 'core.billing', 'core.audit', 'core.documents', 'laboratory.basic', 'radiology.basic', 'clinical.basic', 'pharmacy.basic', 'inventory.basic', 'qms.basic', 'laboratory.analyzer_integration', 'laboratory.advanced_qc', 'clinical.inpatient_wards', 'pharmacy.advanced_inventory')
ON CONFLICT DO NOTHING;

-- Map Enterprise Plan Base Capabilities (All Capabilities)
INSERT INTO subscription.plan_capabilities (plan_id, capability_id)
SELECT p.id, c.id
FROM subscription.plans p
CROSS JOIN subscription.capabilities c
WHERE p.code = 'enterprise'
ON CONFLICT DO NOTHING;

-- Seed Capability Dependencies
INSERT INTO subscription.capability_dependencies (capability_id, depends_on_capability_id)
SELECT c1.id, c2.id
FROM subscription.capabilities c1
CROSS JOIN subscription.capabilities c2
WHERE (c1.code = 'laboratory.analyzer_integration' AND c2.code = 'laboratory.basic')
   OR (c1.code = 'laboratory.advanced_qc' AND c2.code = 'laboratory.basic')
   OR (c1.code = 'radiology.pacs_dicom' AND c2.code = 'radiology.basic')
   OR (c1.code = 'clinical.inpatient_wards' AND c2.code = 'clinical.basic')
   OR (c1.code = 'pharmacy.advanced_inventory' AND c2.code = 'pharmacy.basic')
   OR (c1.code = 'inventory.multi_warehouse' AND c2.code = 'inventory.basic')
   OR (c1.code = 'qms.iso15189' AND c2.code = 'qms.basic')
ON CONFLICT DO NOTHING;

-- +goose Down
SELECT 1;
