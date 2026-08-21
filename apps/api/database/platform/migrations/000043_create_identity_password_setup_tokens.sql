-- +goose Up
-- ==============================================================================
-- CUREXAL IDENTITY: PASSWORD SETUP TOKENS (OWNER & STAFF INVITATIONS)
-- ==============================================================================

CREATE TABLE IF NOT EXISTS identity.password_setup_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organization.organizations(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL,
    token_type VARCHAR(50) NOT NULL DEFAULT 'OWNER_INVITATION',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_password_setup_tokens_hash ON identity.password_setup_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_password_setup_tokens_user ON identity.password_setup_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_password_setup_tokens_org ON identity.password_setup_tokens(organization_id);

ALTER TABLE identity.credentials ALTER COLUMN password_hash DROP NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS identity.password_setup_tokens CASCADE;
