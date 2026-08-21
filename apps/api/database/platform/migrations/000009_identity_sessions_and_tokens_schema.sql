-- +goose Up
-- ==============================================================================
-- IDENTITY SESSIONS & PASSWORD RESET TOKENS
-- ==============================================================================

CREATE TABLE IF NOT EXISTS identity.sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_identity_sessions_user_id ON identity.sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_identity_sessions_token ON identity.sessions(token);

CREATE TABLE IF NOT EXISTS identity.password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_identity_reset_tokens_user ON identity.password_reset_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_identity_reset_tokens_hash ON identity.password_reset_tokens(token_hash);

-- Drop legacy public tables to eliminate foreign key constraint conflicts
DROP TABLE IF EXISTS public.session CASCADE;
DROP TABLE IF EXISTS public.sessions CASCADE;
DROP TABLE IF EXISTS public.verification_token CASCADE;
DROP TABLE IF EXISTS public.password_reset_tokens CASCADE;

-- +goose Down
DROP TABLE IF EXISTS identity.sessions CASCADE;
DROP TABLE IF EXISTS identity.password_reset_tokens CASCADE;
