-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: DROP SINGLE HEADQUARTERS UNIQUE INDEX
-- ==============================================================================

DROP INDEX IF EXISTS organization.uk_facility_single_headquarters;

-- +goose Down
-- Recreate unique index if rolled back
CREATE UNIQUE INDEX IF NOT EXISTS uk_facility_single_headquarters
ON organization.facility_branches (organization_id)
WHERE is_headquarters = TRUE;
