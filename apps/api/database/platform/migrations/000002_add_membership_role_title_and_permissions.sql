-- +goose Up
-- ==============================================================================
-- MIGRATION 000002: Ensure Organization Membership and Authorization Attributes
-- ==============================================================================

ALTER TABLE organization.organization_memberships ADD COLUMN IF NOT EXISTS role_title VARCHAR(100) DEFAULT 'member';
ALTER TABLE organization.organization_memberships ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;
ALTER TABLE organization.organization_memberships ADD COLUMN IF NOT EXISTS joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- +goose Down
ALTER TABLE organization.organization_memberships DROP COLUMN IF EXISTS role_title;
ALTER TABLE organization.organization_memberships DROP COLUMN IF EXISTS is_active;
ALTER TABLE organization.organization_memberships DROP COLUMN IF EXISTS joined_at;

