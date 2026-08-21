-- +goose Up
-- ==============================================================================
-- ENSURE SCHEMAS AND MEMBERSHIP, WORKSPACE AND SUBSCRIPTION COLUMNS
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS organization;
CREATE SCHEMA IF NOT EXISTS workspace;
CREATE SCHEMA IF NOT EXISTS subscription;
CREATE SCHEMA IF NOT EXISTS marketplace;

ALTER TABLE organization.organization_memberships ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'member';
ALTER TABLE organization.organization_memberships ADD COLUMN IF NOT EXISTS role_title VARCHAR(100) DEFAULT 'member';
ALTER TABLE organization.organization_memberships ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;
ALTER TABLE organization.organization_memberships ADD COLUMN IF NOT EXISTS permissions JSONB DEFAULT '[]'::jsonb;

ALTER TABLE workspace.workspace_memberships ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'member';
ALTER TABLE workspace.workspace_memberships ADD COLUMN IF NOT EXISTS role_title VARCHAR(100) DEFAULT 'member';
ALTER TABLE workspace.workspace_memberships ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;

ALTER TABLE workspace.workspaces ADD COLUMN IF NOT EXISTS facility_type VARCHAR(100) DEFAULT 'laboratory';

CREATE TABLE IF NOT EXISTS subscription.subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    plan_id UUID,
    plan VARCHAR(50) NOT NULL DEFAULT 'enterprise',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    starts_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE subscription.subscriptions ADD COLUMN IF NOT EXISTS plan_id UUID;
ALTER TABLE subscription.subscriptions ADD COLUMN IF NOT EXISTS plan VARCHAR(50) DEFAULT 'enterprise';
ALTER TABLE subscription.subscriptions ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'active';

-- +goose Down
SELECT 1;
