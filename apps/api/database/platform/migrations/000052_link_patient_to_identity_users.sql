-- +goose Up
-- Migration 000052: Link patient.patients to identity.users schema & establish patient profiles

ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES identity.users(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_patients_user_id ON patient.patients(user_id);

-- Ensure patient.patient_profiles exists
CREATE TABLE IF NOT EXISTS patient.patient_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES identity.users(id) ON DELETE CASCADE,
    phone VARCHAR(50),
    nin VARCHAR(50),
    gender VARCHAR(20),
    date_of_birth DATE,
    blood_group VARCHAR(10),
    genotype VARCHAR(10),
    address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_patient_profiles_user ON patient.patient_profiles(user_id);

-- +goose Down
DROP INDEX IF EXISTS patient.idx_patients_user_id;
ALTER TABLE patient.patients DROP COLUMN IF EXISTS user_id;
