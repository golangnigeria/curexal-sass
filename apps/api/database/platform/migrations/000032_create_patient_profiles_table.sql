-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: PATIENT PROFILES SCHEMA
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS patient;

CREATE TABLE IF NOT EXISTS patient.patient_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES identity.users(id) ON DELETE CASCADE,
    phone VARCHAR(50),
    date_of_birth DATE,
    gender VARCHAR(20),
    blood_group VARCHAR(10),
    genotype VARCHAR(10),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100) DEFAULT 'Nigeria',
    emergency_contact_name VARCHAR(255),
    emergency_contact_phone VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_patient_profiles_user ON patient.patient_profiles(user_id);

-- +goose Down
DROP TABLE IF EXISTS patient.patient_profiles CASCADE;
