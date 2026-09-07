-- +goose Up
-- Migration 000053: Add phone column and index to identity.users
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS phone VARCHAR(50);
CREATE INDEX IF NOT EXISTS idx_identity_users_phone ON identity.users(phone);

-- +goose Down
DROP INDEX IF EXISTS identity.idx_identity_users_phone;
ALTER TABLE identity.users DROP COLUMN IF EXISTS phone;
