-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: PLATFORM HEALTH & LAUNCH GATE SCHEMA
-- ==============================================================================

-- Launch Gate Audit Reports Table
CREATE TABLE IF NOT EXISTS platform.launch_gate_audits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gate_name VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    check_results JSONB DEFAULT '[]'::jsonb,
    evaluated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    evaluated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_launch_gate_audits_gate ON platform.launch_gate_audits(gate_name, evaluated_at DESC);

-- System Health Metrics Table
CREATE TABLE IF NOT EXISTS platform.system_health_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    component_name VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    metrics JSONB DEFAULT '{}'::jsonb,
    checked_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_health_metrics_comp ON platform.system_health_metrics(component_name, checked_at DESC);

-- Register Launch Gate Permissions
INSERT INTO "authorization".permissions (id, code, module, category, description)
VALUES 
    (gen_random_uuid(), 'platform:launch_gate:read', 'platform', 'system', 'Read platform launch gate audit reports'),
    (gen_random_uuid(), 'platform:launch_gate:execute', 'platform', 'system', 'Execute automated launch gate verification audit')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS platform.system_health_metrics CASCADE;
DROP TABLE IF EXISTS platform.launch_gate_audits CASCADE;
