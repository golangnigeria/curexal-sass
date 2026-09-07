-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: ADD FACILITY THEME BRANDING
-- ==============================================================================

ALTER TABLE organization.facility_branches
ADD COLUMN IF NOT EXISTS theme_branding JSONB NOT NULL DEFAULT '{}'::jsonb;

-- +goose Down
ALTER TABLE organization.facility_branches
DROP COLUMN IF EXISTS theme_branding;
