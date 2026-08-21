-- +goose Up
-- ==============================================================================
-- ADD CREDENTIAL TRACKING COLUMNS TO IDENTITY.USERS
-- ==============================================================================

-- Ensure identity.users has all baseline columns if missing
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS is_platform_admin BOOLEAN DEFAULT FALSE;
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS platform_role VARCHAR(100);
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS credential_status VARCHAR(50) DEFAULT 'ACTIVE';
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS failed_login_attempts INT DEFAULT 0;
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS locked_until TIMESTAMP WITH TIME ZONE;
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS password_changed_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS password_expires_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS last_successful_login_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS last_failed_login_at TIMESTAMP WITH TIME ZONE;

-- Ensure identity.credentials has account_id column if missing
ALTER TABLE identity.credentials ADD COLUMN IF NOT EXISTS account_id VARCHAR(255);

-- +goose Down
SELECT 1;
