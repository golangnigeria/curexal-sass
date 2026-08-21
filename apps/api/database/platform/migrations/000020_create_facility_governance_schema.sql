-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: PLATFORM FACILITY GOVERNANCE SCHEMA
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS platform;

-- 1. Facility Types Table
CREATE TABLE IF NOT EXISTS platform.facility_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'clinical',
    icon_key VARCHAR(100) DEFAULT 'Building2',
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL
);

-- 2. Facility Capability Mapping Table
CREATE TABLE IF NOT EXISTS platform.facility_capabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    facility_type_id UUID NOT NULL REFERENCES platform.facility_types(id) ON DELETE CASCADE,
    capability_id UUID NOT NULL REFERENCES subscription.capabilities(id) ON DELETE CASCADE,
    is_default BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_facility_capability UNIQUE (facility_type_id, capability_id)
);

-- Seed Baseline Facility Types
INSERT INTO platform.facility_types (code, name, category, icon_key, description)
VALUES
    ('hospital', 'Hospital', 'clinical', 'Building2', 'Full service inpatient and outpatient tertiary hospital'),
    ('clinic', 'Outpatient Clinic', 'clinical', 'Stethoscope', 'Primary healthcare and outpatient clinical practice'),
    ('diagnostic_center', 'Diagnostic Center', 'diagnostic', 'Activity', 'Diagnostic pathology and imaging center'),
    ('pharmacy', 'Community Pharmacy', 'retail', 'Pill', 'Retail and clinical prescription pharmacy'),
    ('laboratory', 'Medical Laboratory', 'diagnostic', 'Microscope', 'Pathology, clinical biochemistry and microbiology laboratory')
ON CONFLICT (code) DO NOTHING;

-- Bind Facility Type Default Capabilities
INSERT INTO platform.facility_capabilities (facility_type_id, capability_id, is_default)
SELECT f.id, c.id, TRUE
FROM platform.facility_types f
CROSS JOIN subscription.capabilities c
WHERE (f.code = 'laboratory' AND c.code IN ('core.organization', 'core.patient', 'core.billing', 'laboratory.basic'))
   OR (f.code = 'diagnostic_center' AND c.code IN ('core.organization', 'core.patient', 'core.billing', 'laboratory.basic', 'radiology.basic'))
   OR (f.code = 'clinic' AND c.code IN ('core.organization', 'core.patient', 'core.billing', 'clinical.basic'))
   OR (f.code = 'pharmacy' AND c.code IN ('core.organization', 'core.patient', 'core.billing', 'pharmacy.basic'))
   OR (f.code = 'hospital' AND c.code IN ('core.organization', 'core.patient', 'core.billing', 'clinical.basic', 'laboratory.basic', 'radiology.basic', 'pharmacy.basic'))
ON CONFLICT DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS platform.facility_capabilities;
DROP TABLE IF EXISTS platform.facility_types;
