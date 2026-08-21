-- +goose Up
-- ==============================================================================
-- ADD DELETED_AT SOFT DELETE COLUMN TO IDENTITY.USERS
-- ==============================================================================

ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;

-- +goose Down
SELECT 1;
