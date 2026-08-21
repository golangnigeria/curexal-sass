-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: BRANDING, NOTIFICATIONS & IN-APP ENGINE SCHEMA
-- ==============================================================================

-- Extend organization.organizations with primary_color, theme_branding, and custom_domain_status
ALTER TABLE organization.organizations
ADD COLUMN IF NOT EXISTS primary_color VARCHAR(20) DEFAULT '#0F172A',
ADD COLUMN IF NOT EXISTS theme_branding JSONB DEFAULT '{}'::jsonb,
ADD COLUMN IF NOT EXISTS custom_domain_status VARCHAR(50) DEFAULT 'PENDING';

-- Unique Custom Domain Constraint
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'uk_org_custom_domain'
    ) THEN
        ALTER TABLE organization.organizations
        ADD CONSTRAINT uk_org_custom_domain UNIQUE (custom_domain);
    END IF;
END $$;
-- +goose StatementEnd

-- Notification Provider Configurations (Storing AEAD Encrypted Passwords & API Keys)
CREATE TABLE IF NOT EXISTS organization.notification_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    channel VARCHAR(50) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    sender_email VARCHAR(255),
    sender_name VARCHAR(255),
    host VARCHAR(255),
    port INT,
    username VARCHAR(255),
    encrypted_password TEXT,
    encrypted_api_key TEXT,
    config_metadata JSONB DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    CONSTRAINT uk_org_channel_provider UNIQUE (organization_id, channel, provider)
);

CREATE INDEX IF NOT EXISTS idx_notification_configs_org ON organization.notification_configs(organization_id);

-- Channel-Specific Notification Templates
CREATE TABLE IF NOT EXISTS organization.notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    template_key VARCHAR(100) NOT NULL,
    channel VARCHAR(50) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    body_template TEXT NOT NULL,
    allowed_variables JSONB DEFAULT '[]'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    CONSTRAINT uk_org_template_key_channel UNIQUE (organization_id, template_key, channel)
);

CREATE INDEX IF NOT EXISTS idx_notification_templates_org ON organization.notification_templates(organization_id);

-- Persistent User In-App Notifications
CREATE TABLE IF NOT EXISTS organization.user_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    notification_type VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    data JSONB DEFAULT '{}'::jsonb,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_user_notifications_user ON organization.user_notifications(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_notifications_unread ON organization.user_notifications(user_id, read_at);

-- Notification Delivery Tracking Logs
CREATE TABLE IF NOT EXISTS organization.notification_deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    notification_id UUID,
    channel VARCHAR(50) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    template_key VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    provider_message_id VARCHAR(255),
    attempt_count INT NOT NULL DEFAULT 0,
    last_error TEXT,
    queued_at TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_deliveries_org_status ON organization.notification_deliveries(organization_id, status);

-- Register Branding & Notification Permissions
INSERT INTO "authorization".permissions (id, code, module, category, description)
VALUES 
    (gen_random_uuid(), 'organization:branding:read', 'organization', 'core', 'Read organization white-label branding configuration'),
    (gen_random_uuid(), 'organization:branding:write', 'organization', 'core', 'Update organization white-label branding configuration'),
    (gen_random_uuid(), 'organization:notifications:read', 'organization', 'core', 'Read organization notification providers and templates'),
    (gen_random_uuid(), 'organization:notifications:write', 'organization', 'core', 'Update organization notification providers and templates')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS organization.notification_deliveries CASCADE;
DROP TABLE IF EXISTS organization.user_notifications CASCADE;
DROP TABLE IF EXISTS organization.notification_templates CASCADE;
DROP TABLE IF EXISTS organization.notification_configs CASCADE;
ALTER TABLE organization.organizations DROP CONSTRAINT IF EXISTS uk_org_custom_domain;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS custom_domain_status;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS theme_branding;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS primary_color;
