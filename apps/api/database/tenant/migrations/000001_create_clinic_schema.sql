-- +goose Up
-- SQL tenant migration for Phase 6 Clinic EMR (BSD-006)

CREATE TABLE IF NOT EXISTS patient (
    id TEXT PRIMARY KEY,
    mrn TEXT UNIQUE NOT NULL,
    first_name TEXT NOT NULL,
    middle_name TEXT,
    last_name TEXT NOT NULL,
    dob DATE NOT NULL,
    gender TEXT NOT NULL,
    phone TEXT,
    email TEXT,
    address TEXT,
    blood_group TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    updated_by TEXT,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_patient_mrn ON patient(mrn) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_patient_name ON patient(last_name, first_name);
CREATE INDEX IF NOT EXISTS idx_patient_created_at ON patient(created_at DESC);

CREATE TABLE IF NOT EXISTS patient_identifier (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL REFERENCES patient(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    value TEXT NOT NULL,
    issuer TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    updated_by TEXT
);

CREATE INDEX IF NOT EXISTS idx_patient_identifier_patient ON patient_identifier(patient_id);
CREATE INDEX IF NOT EXISTS idx_patient_identifier_type_value ON patient_identifier(type, value);

CREATE TABLE IF NOT EXISTS provider (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    license_number TEXT NOT NULL,
    specialty TEXT NOT NULL,
    title TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    updated_by TEXT
);

CREATE INDEX IF NOT EXISTS idx_provider_user_id ON provider(user_id);
CREATE INDEX IF NOT EXISTS idx_provider_license ON provider(license_number);

CREATE TABLE IF NOT EXISTS appointment (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL REFERENCES patient(id) ON DELETE CASCADE,
    provider_id TEXT REFERENCES provider(id) ON DELETE SET NULL,
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    reason TEXT,
    status TEXT NOT NULL DEFAULT 'scheduled',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    updated_by TEXT
);

CREATE INDEX IF NOT EXISTS idx_appointment_patient ON appointment(patient_id);
CREATE INDEX IF NOT EXISTS idx_appointment_scheduled_at ON appointment(scheduled_at);

CREATE TABLE IF NOT EXISTS visit (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL REFERENCES patient(id) ON DELETE CASCADE,
    provider_id TEXT REFERENCES provider(id) ON DELETE SET NULL,
    appointment_id TEXT REFERENCES appointment(id) ON DELETE SET NULL,
    visit_number TEXT UNIQUE NOT NULL,
    visit_type TEXT NOT NULL DEFAULT 'outpatient',
    status TEXT NOT NULL DEFAULT 'in_progress',
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    updated_by TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_visit_number ON visit(visit_number);
CREATE INDEX IF NOT EXISTS idx_visit_patient ON visit(patient_id);
CREATE INDEX IF NOT EXISTS idx_visit_status ON visit(status);

-- +goose Down
DROP TABLE IF EXISTS visit;
DROP TABLE IF EXISTS appointment;
DROP TABLE IF EXISTS provider;
DROP TABLE IF EXISTS patient_identifier;
DROP TABLE IF EXISTS patient;
