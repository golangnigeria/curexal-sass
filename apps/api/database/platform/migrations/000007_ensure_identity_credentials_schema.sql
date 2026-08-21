-- +goose Up
-- ==============================================================================
-- ENSURE ACCOUNT_ID COLUMN IN IDENTITY.CREDENTIALS
-- ==============================================================================

ALTER TABLE identity.credentials ADD COLUMN IF NOT EXISTS account_id VARCHAR(255);

-- +goose Down
SELECT 1;
