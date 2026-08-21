-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: GRANT DOCUMENT UPLOAD PERMISSIONS TO ROLES
-- ==============================================================================

ALTER TABLE "authorization".roles ADD COLUMN IF NOT EXISTS code VARCHAR(100);

INSERT INTO "authorization".permissions (id, code, module, category, description)
VALUES 
    (gen_random_uuid(), 'organization:document:upload', 'organization', 'core', 'Upload regulatory organization documents'),
    (gen_random_uuid(), 'organization:document:read', 'organization', 'core', 'Read organization documents and metadata')
ON CONFLICT (code) DO NOTHING;

-- Grant document upload and read permissions to owner, org_admin, and branch_admin roles
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE (r.name IN ('owner', 'org_admin', 'branch_admin', 'Owner', 'Org Admin', 'Branch Admin') OR COALESCE(r.code, '') IN ('owner', 'org_admin', 'branch_admin'))
  AND p.code IN ('organization:document:upload', 'organization:document:read')
ON CONFLICT DO NOTHING;

-- +goose Down
SELECT 1;
