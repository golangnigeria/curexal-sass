-- +goose Up
-- Migration: 000035_identity_name_source_of_truth.sql
-- Description: Establish identity.user_profiles as the single source of truth for structured user names,
-- relax identity.users.name constraints, synchronize existing records, and install automated sync trigger.

-- 1. Relax NOT NULL constraint on identity.users.name
ALTER TABLE identity.users ALTER COLUMN name DROP NOT NULL;
ALTER TABLE identity.users ALTER COLUMN name SET DEFAULT '';

-- 2. One-time data synchronization: sync identity.users.name from authoritative user_profiles
UPDATE identity.users u
SET name = COALESCE(
    NULLIF(TRIM(CONCAT_WS(' ',
        NULLIF(TRIM(p.first_name), ''),
        NULLIF(TRIM(p.middle_name), ''),
        NULLIF(TRIM(p.last_name), '')
    )), ''),
    u.name,
    u.email
),
avatar_url = COALESCE(p.avatar_url, u.avatar_url),
updated_at = CURRENT_TIMESTAMP
FROM identity.user_profiles p
WHERE p.user_id = u.id;

-- 3. PostgreSQL Trigger Function to automatically synchronize identity.users.name and avatar_url from user_profiles
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION identity.sync_user_name_from_profile()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE identity.users
    SET name = COALESCE(
        NULLIF(TRIM(CONCAT_WS(' ',
            NULLIF(TRIM(NEW.first_name), ''),
            NULLIF(TRIM(NEW.middle_name), ''),
            NULLIF(TRIM(NEW.last_name), '')
        )), ''),
        email
    ),
    avatar_url = COALESCE(NEW.avatar_url, avatar_url),
    updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.user_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS trg_sync_user_name_from_profile ON identity.user_profiles;
CREATE TRIGGER trg_sync_user_name_from_profile
AFTER INSERT OR UPDATE OF first_name, middle_name, last_name, avatar_url ON identity.user_profiles
FOR EACH ROW
EXECUTE FUNCTION identity.sync_user_name_from_profile();

-- 4. Create display view for reporting / query projection
CREATE OR REPLACE VIEW identity.user_profile_display AS
SELECT
    u.id AS user_id,
    u.email,
    p.first_name,
    p.middle_name,
    p.last_name,
    TRIM(CONCAT_WS(' ',
        NULLIF(TRIM(p.first_name), ''),
        NULLIF(TRIM(p.middle_name), ''),
        NULLIF(TRIM(p.last_name), '')
    )) AS display_name,
    COALESCE(p.avatar_url, u.avatar_url) AS avatar_url,
    p.phone_number,
    u.is_platform_admin,
    u.platform_role,
    u.created_at,
    u.updated_at
FROM identity.users u
LEFT JOIN identity.user_profiles p ON p.user_id = u.id;

-- +goose Down
DROP VIEW IF EXISTS identity.user_profile_display;
DROP TRIGGER IF EXISTS trg_sync_user_name_from_profile ON identity.user_profiles;
DROP FUNCTION IF EXISTS identity.sync_user_name_from_profile();
