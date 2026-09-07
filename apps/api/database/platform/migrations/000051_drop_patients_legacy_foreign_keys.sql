-- +goose Up
-- Migration 000051: Drop legacy foreign key constraints and nullable alignment on patient.patients

ALTER TABLE patient.patients DROP CONSTRAINT IF EXISTS patients_organization_id_fkey;
ALTER TABLE patient.patients DROP CONSTRAINT IF EXISTS uk_patients_mrn;
ALTER TABLE patient.patients ALTER COLUMN organization_id DROP NOT NULL;
ALTER TABLE patient.patients ALTER COLUMN dob DROP NOT NULL;

-- +goose Down
-- Reversible
