-- +goose Up
-- ==============================================================================
-- ENSURE AUDIT.AUDIT_EVENTS COLUMNS: ORGANIZATION_ID & WORKSPACE_ID
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS audit;

CREATE TABLE IF NOT EXISTS audit.audit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    organization_id UUID REFERENCES organization.organizations(id) ON DELETE SET NULL,
    workspace_id UUID REFERENCES workspace.workspaces(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id VARCHAR(255),
    payload JSONB DEFAULT '{}'::jsonb,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organization.organizations(id) ON DELETE SET NULL;
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS workspace_id UUID REFERENCES workspace.workspaces(id) ON DELETE SET NULL;

-- +goose Down
SELECT 1;
