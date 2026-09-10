-- +goose Up
-- ==============================================================================
-- CUREXAL MIGRATION 000074: DAY 2 PROVIDER PROFILES, APPOINTMENTS & PATIENT EXPANSION
-- ==============================================================================

-- 1. EXTENSIONS
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- 2. ENHANCE: orchestration.provider_profiles (Feature #4)
ALTER TABLE orchestration.provider_profiles ADD COLUMN IF NOT EXISTS license_number VARCHAR(100);
ALTER TABLE orchestration.provider_profiles ADD COLUMN IF NOT EXISTS license_issuer VARCHAR(100) DEFAULT 'MDCN';
ALTER TABLE orchestration.provider_profiles ADD COLUMN IF NOT EXISTS license_verified_at TIMESTAMPTZ;
ALTER TABLE orchestration.provider_profiles ADD COLUMN IF NOT EXISTS room_number VARCHAR(50);
ALTER TABLE orchestration.provider_profiles ADD COLUMN IF NOT EXISTS room_name VARCHAR(100);
ALTER TABLE orchestration.provider_profiles ADD COLUMN IF NOT EXISTS virtual_room_url TEXT;

-- Drop existing check constraint on provider status if exists and re-add standard set
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'chk_provider_status' AND table_schema = 'orchestration'
    ) THEN
        ALTER TABLE orchestration.provider_profiles DROP CONSTRAINT chk_provider_status;
    END IF;
END $$;
-- +goose StatementEnd

ALTER TABLE orchestration.provider_profiles
    ADD CONSTRAINT chk_provider_status CHECK (status IN ('ON_DUTY', 'ON_BREAK', 'OFF_DUTY', 'BUSY'));

-- Seed missing license numbers for existing provider records
UPDATE orchestration.provider_profiles
SET license_number = 'MDCN/R/' || SUBSTRING(id::text, 1, 8),
    room_number = COALESCE(room_number, 'Suite 104'),
    room_name = COALESCE(room_name, 'Primary Clinical Examination Room'),
    virtual_room_url = COALESCE(virtual_room_url, 'https://telehealth.curexal.com/room/' || id::text)
WHERE license_number IS NULL OR license_number = '';

CREATE UNIQUE INDEX IF NOT EXISTS uk_provider_license 
    ON orchestration.provider_profiles(license_number) WHERE license_number IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_provider_profiles_tenant_status 
    ON orchestration.provider_profiles(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_provider_profiles_specialty 
    ON orchestration.provider_profiles(specialty_code);
CREATE INDEX IF NOT EXISTS idx_provider_profiles_channels 
    ON orchestration.provider_profiles(in_person_enabled, telehealth_enabled);

-- 3. ENHANCE: patient.patients & patient.patient_guardians (Feature #5)
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS preferred_name VARCHAR(100);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS marital_status VARCHAR(20);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS occupation VARCHAR(100);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS residential_address TEXT;
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS city VARCHAR(100);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS state VARCHAR(100);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS country VARCHAR(100) DEFAULT 'Nigeria';
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS preferred_language VARCHAR(50) DEFAULT 'English';

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'uk_patients_tenant_mrn' AND table_schema = 'patient'
    ) THEN
        -- Ensure unique mrn per tenant
        CREATE UNIQUE INDEX IF NOT EXISTS uk_patients_tenant_mrn ON patient.patients(tenant_id, mrn) WHERE mrn IS NOT NULL;
    END IF;
END $$;
-- +goose StatementEnd

CREATE INDEX IF NOT EXISTS idx_patients_channel ON patient.patients(registration_channel);

-- Table: patient.patient_guardians (Next of Kin & Emergency Contacts)
CREATE TABLE IF NOT EXISTS patient.patient_guardians (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    relationship_type VARCHAR(50) NOT NULL, -- SPOUSE, PARENT, CHILD, SIBLING, GUARDIAN, NEXT_OF_KIN, OTHER
    full_name VARCHAR(150) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    email VARCHAR(100),
    address TEXT,
    is_emergency_contact BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_patient_guardians_patient_id ON patient.patient_guardians(patient_id);

-- 4. SCHEMA: operations & operations.appointments (Feature #8)
CREATE SCHEMA IF NOT EXISTS operations;

CREATE TABLE IF NOT EXISTS operations.appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES orchestration.provider_profiles(id) ON DELETE RESTRICT,
    appointment_number VARCHAR(64) NOT NULL,
    service_type VARCHAR(50) NOT NULL DEFAULT 'CONSULTATION',
    delivery_channel VARCHAR(30) NOT NULL DEFAULT 'in_person',
    status VARCHAR(30) NOT NULL DEFAULT 'BOOKED',
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    reason_for_visit TEXT,
    cancellation_reason TEXT,
    virtual_meeting_url TEXT,
    created_by UUID REFERENCES identity.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_appointment_number UNIQUE (tenant_id, appointment_number),
    CONSTRAINT chk_appointment_channel CHECK (delivery_channel IN ('in_person', 'video', 'telephone', 'secure_message')),
    CONSTRAINT chk_appointment_status CHECK (status IN ('BOOKED', 'CHECKED_IN', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED', 'NO_SHOW'))
);

CREATE INDEX IF NOT EXISTS idx_appointments_tenant_date 
    ON operations.appointments(tenant_id, start_time);
CREATE INDEX IF NOT EXISTS idx_appointments_patient 
    ON operations.appointments(patient_id);
CREATE INDEX IF NOT EXISTS idx_appointments_provider_slot 
    ON operations.appointments(provider_id, start_time, end_time);
CREATE INDEX IF NOT EXISTS idx_appointments_channel 
    ON operations.appointments(delivery_channel);

-- 5. ENHANCE: orchestration.care_requests (Feature #10)
ALTER TABLE orchestration.care_requests ADD COLUMN IF NOT EXISTS appointment_id UUID REFERENCES operations.appointments(id) ON DELETE SET NULL;
ALTER TABLE orchestration.care_requests ADD COLUMN IF NOT EXISTS delivery_channel VARCHAR(30) NOT NULL DEFAULT 'in_person';
ALTER TABLE orchestration.care_requests ADD COLUMN IF NOT EXISTS checked_in_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE orchestration.care_requests ADD COLUMN IF NOT EXISTS triaged_at TIMESTAMPTZ;
ALTER TABLE orchestration.care_requests ADD COLUMN IF NOT EXISTS consultation_started_at TIMESTAMPTZ;
ALTER TABLE orchestration.care_requests ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ;

-- Backfill delivery_channel from preferred_mode if preferred_mode exists
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema='orchestration' AND table_name='care_requests' AND column_name='preferred_mode'
    ) THEN
        UPDATE orchestration.care_requests 
        SET delivery_channel = CASE 
            WHEN LOWER(preferred_mode) = 'video' THEN 'video'
            WHEN LOWER(preferred_mode) = 'audio' THEN 'telephone'
            WHEN LOWER(preferred_mode) = 'async_chat' THEN 'secure_message'
            ELSE 'in_person'
        END
        WHERE delivery_channel = 'in_person' AND preferred_mode IS NOT NULL;
    END IF;
END $$;
-- +goose StatementEnd

CREATE INDEX IF NOT EXISTS idx_care_requests_queue_active 
    ON orchestration.care_requests(tenant_id, status, checked_in_at);
CREATE INDEX IF NOT EXISTS idx_care_requests_channel 
    ON orchestration.care_requests(delivery_channel);

-- +goose Down
DROP TABLE IF EXISTS operations.appointments CASCADE;
DROP TABLE IF EXISTS patient.patient_guardians CASCADE;
