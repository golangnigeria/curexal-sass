-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: UPGRADE IDENTITY SESSIONS SCHEMA (FEATURE #1)
-- ==============================================================================

DROP TABLE IF EXISTS identity.sessions CASCADE;

CREATE TABLE IF NOT EXISTS identity.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    active_branch_id UUID REFERENCES organization.facility_branches(id) ON DELETE SET NULL,
    refresh_token_hash VARCHAR(64) NOT NULL UNIQUE,
    token_family_id UUID NOT NULL DEFAULT gen_random_uuid(),
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT NOT NULL,
    device_fingerprint VARCHAR(64),
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    revocation_reason VARCHAR(255),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_active_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON identity.sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_org_id ON identity.sessions(organization_id);
CREATE INDEX IF NOT EXISTS idx_sessions_branch_id ON identity.sessions(active_branch_id) WHERE active_branch_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_sessions_token_hash ON identity.sessions(refresh_token_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_family_id ON identity.sessions(token_family_id);

-- +goose Down
DROP TABLE IF EXISTS identity.sessions CASCADE;

CREATE TABLE IF NOT EXISTS identity.sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
