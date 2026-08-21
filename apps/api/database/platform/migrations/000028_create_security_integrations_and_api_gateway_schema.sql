-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: SECURITY INTEGRATIONS & API GATEWAY SCHEMA
-- ==============================================================================

-- External API Keys Table (Persisted ONLY as SHA-256 Hashes)
CREATE TABLE IF NOT EXISTS organization.api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    key_prefix VARCHAR(20) NOT NULL,
    key_hash VARCHAR(64) NOT NULL,
    scopes JSONB DEFAULT '[]'::jsonb,
    ip_whitelist JSONB DEFAULT '[]'::jsonb,
    rate_limit_rpm INT NOT NULL DEFAULT 60,
    expires_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    CONSTRAINT uk_org_api_key_hash UNIQUE (key_hash)
);

CREATE INDEX IF NOT EXISTS idx_api_keys_org ON organization.api_keys(organization_id);

-- Webhook Subscriptions Table (Storing AEAD Encrypted Signing Secrets)
CREATE TABLE IF NOT EXISTS organization.webhook_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    target_url TEXT NOT NULL,
    event_types JSONB DEFAULT '[]'::jsonb,
    encrypted_secret TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    CONSTRAINT uk_org_webhook_target UNIQUE (organization_id, target_url)
);

CREATE INDEX IF NOT EXISTS idx_webhook_subs_org ON organization.webhook_subscriptions(organization_id);

-- Webhook Delivery Tracking Logs
CREATE TABLE IF NOT EXISTS organization.webhook_deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    subscription_id UUID REFERENCES organization.webhook_subscriptions(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB DEFAULT '{}'::jsonb,
    response_status INT,
    response_body TEXT,
    attempt_count INT NOT NULL DEFAULT 1,
    last_error TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_org ON organization.webhook_deliveries(organization_id, status);

-- Register Integration Permissions
INSERT INTO "authorization".permissions (id, code, module, category, description)
VALUES 
    (gen_random_uuid(), 'organization:integrations:read', 'organization', 'core', 'Read organization API keys and webhooks'),
    (gen_random_uuid(), 'organization:integrations:write', 'organization', 'core', 'Create and revoke organization API keys and webhooks')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS organization.webhook_deliveries CASCADE;
DROP TABLE IF EXISTS organization.webhook_subscriptions CASCADE;
DROP TABLE IF EXISTS organization.api_keys CASCADE;
