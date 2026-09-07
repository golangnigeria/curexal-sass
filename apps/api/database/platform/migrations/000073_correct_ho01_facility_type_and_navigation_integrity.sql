-- +goose Up
-- ==============================================================================
-- CUREXAL MIGRATION 000073: HO-01 FACILITY TYPE & NAVIGATION INTEGRITY CORRECTION
-- ==============================================================================

-- 1. Ensure Canonical Facility Types Exist with Accurate Names and Codes
INSERT INTO platform.facility_types (code, name, category, icon_key, description)
VALUES
    ('clinic', 'Outpatient Clinic', 'clinical', 'Stethoscope', 'Primary healthcare, general practice, and outpatient medical center'),
    ('laboratory', 'Diagnostic Laboratory', 'diagnostic', 'Microscope', 'Clinical pathology, biochemistry, and microbiology laboratory'),
    ('radiology', 'Radiology & Imaging Center', 'diagnostic', 'Layers', 'Diagnostic imaging, ultrasound, X-ray, CT, and PACS viewing'),
    ('pharmacy', 'Community Pharmacy', 'retail', 'Pill', 'Retail and prescription fulfillment dispensary'),
    ('hospital', 'General Hospital', 'clinical', 'Building2', 'Multi-specialty tertiary inpatient and outpatient healthcare facility')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    category = EXCLUDED.category,
    icon_key = EXCLUDED.icon_key,
    description = EXCLUDED.description;

-- 2. Correct HO-01: Point strictly to canonical facility type 'clinic' (Outpatient Clinic)
UPDATE organization.facility_branches fb
SET facility_type_id = ft.id,
    updated_at = NOW()
FROM platform.facility_types ft
WHERE (fb.code = 'HO-01' OR fb.slug = 'curexal-clinic')
  AND ft.code = 'clinic';

-- 3. Ensure Default Capabilities for Canonical Facility Types
-- Outpatient Clinic: Only patient, clinical basic, billing, organization
DELETE FROM platform.facility_capabilities fc
USING platform.facility_types ft, subscription.capabilities c
WHERE fc.facility_type_id = ft.id
  AND fc.capability_id = c.id
  AND ft.code = 'clinic'
  AND c.code IN ('laboratory.basic', 'radiology.basic', 'clinical.inpatient_wards');

INSERT INTO platform.facility_capabilities (facility_type_id, capability_id, is_default)
SELECT ft.id, c.id, TRUE
FROM platform.facility_types ft
CROSS JOIN subscription.capabilities c
WHERE ft.code = 'clinic'
  AND c.code IN ('core.organization', 'core.patient', 'core.billing', 'clinical.basic')
ON CONFLICT (facility_type_id, capability_id) DO NOTHING;

-- Diagnostic Laboratory: Only patient, billing, laboratory basic, organization
DELETE FROM platform.facility_capabilities fc
USING platform.facility_types ft, subscription.capabilities c
WHERE fc.facility_type_id = ft.id
  AND fc.capability_id = c.id
  AND ft.code = 'laboratory'
  AND c.code IN ('clinical.basic', 'radiology.basic', 'pharmacy.basic', 'clinical.inpatient_wards');

INSERT INTO platform.facility_capabilities (facility_type_id, capability_id, is_default)
SELECT ft.id, c.id, TRUE
FROM platform.facility_types ft
CROSS JOIN subscription.capabilities c
WHERE ft.code = 'laboratory'
  AND c.code IN ('core.organization', 'core.patient', 'core.billing', 'laboratory.basic')
ON CONFLICT (facility_type_id, capability_id) DO NOTHING;

-- Radiology Center: Only patient, billing, radiology basic, organization
DELETE FROM platform.facility_capabilities fc
USING platform.facility_types ft, subscription.capabilities c
WHERE fc.facility_type_id = ft.id
  AND fc.capability_id = c.id
  AND ft.code = 'radiology'
  AND c.code IN ('clinical.basic', 'laboratory.basic', 'pharmacy.basic', 'clinical.inpatient_wards');

INSERT INTO platform.facility_capabilities (facility_type_id, capability_id, is_default)
SELECT ft.id, c.id, TRUE
FROM platform.facility_types ft
CROSS JOIN subscription.capabilities c
WHERE ft.code = 'radiology'
  AND c.code IN ('core.organization', 'core.patient', 'core.billing', 'radiology.basic')
ON CONFLICT (facility_type_id, capability_id) DO NOTHING;

-- Community Pharmacy: Only patient, billing, pharmacy basic, organization
DELETE FROM platform.facility_capabilities fc
USING platform.facility_types ft, subscription.capabilities c
WHERE fc.facility_type_id = ft.id
  AND fc.capability_id = c.id
  AND ft.code = 'pharmacy'
  AND c.code IN ('clinical.basic', 'laboratory.basic', 'radiology.basic', 'clinical.inpatient_wards');

INSERT INTO platform.facility_capabilities (facility_type_id, capability_id, is_default)
SELECT ft.id, c.id, TRUE
FROM platform.facility_types ft
CROSS JOIN subscription.capabilities c
WHERE ft.code = 'pharmacy'
  AND c.code IN ('core.organization', 'core.patient', 'core.billing', 'pharmacy.basic')
ON CONFLICT (facility_type_id, capability_id) DO NOTHING;

-- General Hospital: All capabilities
INSERT INTO platform.facility_capabilities (facility_type_id, capability_id, is_default)
SELECT ft.id, c.id, TRUE
FROM platform.facility_types ft
CROSS JOIN subscription.capabilities c
WHERE ft.code = 'hospital'
  AND c.code IN ('core.organization', 'core.patient', 'core.billing', 'clinical.basic', 'laboratory.basic', 'radiology.basic', 'pharmacy.basic', 'clinical.inpatient_wards')
ON CONFLICT (facility_type_id, capability_id) DO NOTHING;

-- 4. Reconcile Navigation Items: Align required_permission and required_capability
UPDATE navigation_item
SET required_permission = 'workspace:patient:create',
    required_capability = 'core.patient'
WHERE id = 'nav_wsp_reception';

UPDATE navigation_item
SET required_permission = 'workspace:queue:manage',
    required_capability = 'clinical.basic'
WHERE id = 'nav_wsp_care_desk';

UPDATE navigation_item
SET required_permission = 'workspace:clinical:read',
    required_capability = 'clinical.basic'
WHERE id = 'nav_wsp_clinical';

UPDATE navigation_item
SET required_permission = 'workspace:billing:read',
    required_capability = 'core.billing'
WHERE id = 'nav_wsp_billing';

UPDATE navigation_item
SET required_permission = 'workspace:sample:receive',
    required_capability = 'laboratory.basic'
WHERE id = 'nav_wsp_laboratory';

UPDATE navigation_item
SET required_permission = 'workspace:prescription:read',
    required_capability = 'pharmacy.basic'
WHERE id = 'nav_wsp_pharmacy';

UPDATE navigation_item
SET required_permission = 'workspace:diagnostic:read',
    required_capability = 'radiology.basic'
WHERE id = 'nav_wsp_radiology';

UPDATE navigation_item
SET required_permission = 'workspace:clinical:read',
    required_capability = 'clinical.inpatient_wards'
WHERE id = 'nav_wsp_hospital';

-- +goose Down
-- Revert to original settings if needed
