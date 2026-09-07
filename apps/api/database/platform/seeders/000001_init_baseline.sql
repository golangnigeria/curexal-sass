-- ==============================================================================
-- PLATFORM SEEDER BASELINE
-- Idempotently maintains platform and organization baseline entities
-- ==============================================================================

-- 1. Ensure default organizations exist
INSERT INTO organization.organizations (id, name, slug, status, plan, setup_state)
VALUES 
    ('00000000-0000-0000-0000-000000000001', 'Everight Diagnostic & Speciality Hospital', 'everight', 'active', 'enterprise', 'VERIFIED'),
    ('00000000-0000-0000-0000-000000000002', 'Curexal Premier Medical Center', 'curexal-clinic', 'active', 'enterprise', 'VERIFIED')
ON CONFLICT (slug) DO UPDATE SET 
    name = EXCLUDED.name, 
    status = 'active', 
    setup_state = 'VERIFIED';

-- 2. Ensure default facility branches exist
INSERT INTO organization.facility_branches (organization_id, facility_type_id, name, code, slug, is_headquarters, status)
SELECT 
    o.id, 
    COALESCE((SELECT id FROM platform.facility_types WHERE name ILIKE '%hospital%' OR name ILIKE '%diagnostic%' LIMIT 1), (SELECT id FROM platform.facility_types LIMIT 1)),
    o.name || ' (Main Branch)', 
    'main', 
    'main', 
    TRUE, 
    'ACTIVE'
FROM organization.organizations o
WHERE o.slug IN ('everight', 'curexal-clinic')
ON CONFLICT (organization_id, code) DO UPDATE SET 
    is_headquarters = TRUE, 
    status = 'ACTIVE';
