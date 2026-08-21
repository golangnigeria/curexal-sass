-- +goose Up
-- ==============================================================================
-- ENSURE ORGANIZATION COLUMNS IN ORGANIZATION.ORGANIZATIONS
-- ==============================================================================

ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS plan VARCHAR(50) DEFAULT 'enterprise';
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'pending_verification';
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS logo_url TEXT;
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS custom_domain VARCHAR(255);
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS settings JSONB DEFAULT '{}'::jsonb;

-- +goose Down
SELECT 1;
