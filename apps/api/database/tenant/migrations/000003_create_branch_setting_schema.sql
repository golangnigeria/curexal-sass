-- +goose Up
-- SQL tenant migration for Branch Settings (BSD-010)

CREATE TABLE IF NOT EXISTS branch_setting (
    id TEXT PRIMARY KEY,
    tenant_id TEXT UNIQUE NOT NULL,
    general_config JSONB NOT NULL DEFAULT '{}',
    financial_config JSONB NOT NULL DEFAULT '{}',
    inventory_config JSONB NOT NULL DEFAULT '{}',
    integrations_config JSONB NOT NULL DEFAULT '{}',
    notifications_config JSONB NOT NULL DEFAULT '{}',
    document_header_config JSONB NOT NULL DEFAULT '{}',
    patient_config JSONB NOT NULL DEFAULT '{}',
    lims_config JSONB NOT NULL DEFAULT '{}',
    consultation_config JSONB NOT NULL DEFAULT '{}',
    staff_config JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_branch_setting_tenant ON branch_setting(tenant_id);

-- +goose Down
DROP TABLE IF EXISTS branch_setting;
