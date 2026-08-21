-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: ENSURE ALL AUDIT EVENT COLUMNS & COMPATIBILITY VIEW
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS audit;

CREATE TABLE IF NOT EXISTS audit.audit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    organization_id UUID REFERENCES organization.organizations(id) ON DELETE SET NULL,
    tenant_id UUID,
    workspace_id UUID REFERENCES workspace.workspaces(id) ON DELETE SET NULL,
    actor_id UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    actor_name VARCHAR(255),
    actor_role VARCHAR(100),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id VARCHAR(255),
    resource_name VARCHAR(255),
    event_category VARCHAR(100),
    severity VARCHAR(50) DEFAULT 'INFO',
    status VARCHAR(50) DEFAULT 'SUCCESS',
    ip_address VARCHAR(50),
    device VARCHAR(100),
    operating_system VARCHAR(100),
    browser VARCHAR(100),
    user_agent TEXT,
    hostname VARCHAR(255),
    request_id VARCHAR(255),
    session_id VARCHAR(255),
    trace_id VARCHAR(255),
    payload JSONB DEFAULT '{}'::jsonb,
    before_state JSONB,
    after_state JSONB,
    reason TEXT,
    approval_reference VARCHAR(255),
    digital_signature TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS actor_name VARCHAR(255);
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS actor_role VARCHAR(100);
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS resource_name VARCHAR(255);
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS event_category VARCHAR(100);
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS severity VARCHAR(50) DEFAULT 'INFO';
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'SUCCESS';
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS device VARCHAR(100);
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS operating_system VARCHAR(100);
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS browser VARCHAR(100);
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS hostname VARCHAR(255);
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS request_id VARCHAR(255);
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS session_id VARCHAR(255);
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS trace_id VARCHAR(255);
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS before_state JSONB;
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS after_state JSONB;
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS reason TEXT;
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS approval_reference VARCHAR(255);
ALTER TABLE audit.audit_events ADD COLUMN IF NOT EXISTS digital_signature TEXT;

-- +goose Down
SELECT 1;
