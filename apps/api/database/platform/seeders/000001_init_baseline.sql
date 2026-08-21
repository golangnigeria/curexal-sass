-- ==============================================================================
-- CLEANUP RECENT SEEDED TEST DATA
-- Purges all test organizations, tenants, memberships, and user accounts
-- ==============================================================================

DELETE FROM organization.organization_memberships WHERE user_id::text LIKE 'usr_%' OR id::text LIKE 'mem_%';
DELETE FROM organization.facility_branches WHERE id::text LIKE '00000000%' OR code IN ('main-facility', 'owerri-branch', 'port-harcourt-branch', 'lifecare-hq');
DELETE FROM organization.organizations WHERE id::text LIKE '00000000%' OR slug IN ('everight', 'lifecare', 'primelab');
DELETE FROM identity.users WHERE id::text LIKE 'usr_%' OR email LIKE '%@test.curexal.local' OR email IN ('owner@everight.com', 'admin@curexal.com', 'manager@everight.com');
