-- +goose Up
-- Migration: Add emergency_contact column to identity.user_profiles
-- Reason: The identity user profile repository queries p.emergency_contact but the column was never created.

ALTER TABLE identity.user_profiles
    ADD COLUMN IF NOT EXISTS emergency_contact TEXT;

COMMENT ON COLUMN identity.user_profiles.emergency_contact IS 'Free-text emergency contact information (name, phone, relationship)';

-- +goose Down
ALTER TABLE identity.user_profiles
    DROP COLUMN IF EXISTS emergency_contact;
