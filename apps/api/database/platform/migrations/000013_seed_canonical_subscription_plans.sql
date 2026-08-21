-- +goose Up
-- ==============================================================================
-- SEED CANONICAL SUBSCRIPTION PLANS & ENSURE AUDIT EVENT COLUMNS
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS subscription;
CREATE SCHEMA IF NOT EXISTS audit;

CREATE TABLE IF NOT EXISTS subscription.plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    limits JSONB NOT NULL DEFAULT '{}'::jsonb,
    features JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO subscription.plans (code, name, limits, features)
VALUES
    ('smart', 'Smart', '{"maxBranches": 1, "maxMembers": 5, "storageGb": 10}'::jsonb, '["customer_care", "billing"]'::jsonb),
    ('optimize', 'Optimize', '{"maxBranches": 3, "maxMembers": 25, "storageGb": 50}'::jsonb, '["customer_care", "laboratory", "billing"]'::jsonb),
    ('pro', 'Pro', '{"maxBranches": 10, "maxMembers": 100, "storageGb": 200}'::jsonb, '["customer_care", "laboratory", "clinical", "billing"]'::jsonb),
    ('enterprise', 'Enterprise', '{"maxBranches": 1000, "maxMembers": 10000, "storageGb": 5000}'::jsonb, '["customer_care", "laboratory", "clinical", "pharmacy", "inventory", "radiology", "billing"]'::jsonb)
ON CONFLICT (code) DO NOTHING;

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
