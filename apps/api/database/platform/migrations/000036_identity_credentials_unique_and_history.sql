-- +goose Up
-- Migration: 000036_identity_credentials_unique_and_history.sql
-- Description: Deduplicate identity.credentials, enforce UNIQUE(user_id, auth_provider), and index password_histories.

-- 1. Deduplicate identity.credentials: keep the single most recently updated row per (user_id, auth_provider)
DELETE FROM identity.credentials
WHERE id NOT IN (
    SELECT DISTINCT ON (user_id, auth_provider) id
    FROM identity.credentials
    ORDER BY user_id, auth_provider, updated_at DESC, created_at DESC, id DESC
);

-- 2. Add UNIQUE constraint to enforce exactly one active credential per (user_id, auth_provider)
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'uk_identity_credentials_user_provider'
    ) THEN
        ALTER TABLE identity.credentials
        ADD CONSTRAINT uk_identity_credentials_user_provider UNIQUE (user_id, auth_provider);
    END IF;
END $$;
-- +goose StatementEnd

-- 3. Ensure identity.password_histories has user_agent column if missing
ALTER TABLE identity.password_histories ADD COLUMN IF NOT EXISTS user_agent TEXT;

-- 4. Add index on identity.password_histories for high-performance history retrieval
CREATE INDEX IF NOT EXISTS idx_password_histories_user_created 
ON identity.password_histories(user_id, created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS identity.idx_password_histories_user_created;
ALTER TABLE identity.password_histories DROP COLUMN IF EXISTS user_agent;
ALTER TABLE identity.credentials DROP CONSTRAINT IF EXISTS uk_identity_credentials_user_provider;
