-- +goose Up
-- Migration 000048: Patient Access, Care Orchestration & Clinical Entry Schemas

-- ============================================================================
-- 1. SCHEMA: patient (Patient Identity, MPI & Portal Accounts)
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS patient;

-- Canonical Patient Registry (Tenant Isolated)
CREATE TABLE IF NOT EXISTS patient.patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    mrn VARCHAR(64),
    first_name VARCHAR(100),
    middle_name VARCHAR(100),
    last_name VARCHAR(100),
    gender VARCHAR(20),
    date_of_birth DATE,
    blood_group VARCHAR(10),
    genotype VARCHAR(10),
    nin VARCHAR(30),
    status VARCHAR(30) DEFAULT 'REGISTERED',
    registration_channel VARCHAR(30) DEFAULT 'RECEPTION',
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Ensure all canonical columns exist if table was created in an earlier migration
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS tenant_id UUID;
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS mrn VARCHAR(64);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS first_name VARCHAR(100);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS middle_name VARCHAR(100);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS last_name VARCHAR(100);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS gender VARCHAR(20);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS date_of_birth DATE;
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS blood_group VARCHAR(10);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS genotype VARCHAR(10);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS nin VARCHAR(30);
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS status VARCHAR(30) DEFAULT 'REGISTERED';
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS registration_channel VARCHAR(30) DEFAULT 'RECEPTION';
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}'::jsonb;
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT NOW();
ALTER TABLE patient.patients ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();

-- Backfill missing columns from legacy columns if present
-- +goose StatementBegin
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='patient' AND table_name='patients' AND column_name='dob') THEN
        UPDATE patient.patients SET date_of_birth = dob WHERE date_of_birth IS NULL AND dob IS NOT NULL;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='patient' AND table_name='patients' AND column_name='organization_id') THEN
        UPDATE patient.patients SET tenant_id = organization_id WHERE tenant_id IS NULL AND organization_id IS NOT NULL;
    END IF;
END $$;
-- +goose StatementEnd

CREATE INDEX IF NOT EXISTS idx_patients_tenant_name ON patient.patients(tenant_id, last_name, first_name);
CREATE INDEX IF NOT EXISTS idx_patients_tenant_dob ON patient.patients(tenant_id, date_of_birth);
CREATE INDEX IF NOT EXISTS idx_patients_tenant_nin ON patient.patients(tenant_id, nin);

-- Patient Contacts & Telecoms
CREATE TABLE IF NOT EXISTS patient.patient_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    system VARCHAR(20) NOT NULL, -- PHONE, EMAIL, WHATSAPP
    value VARCHAR(255) NOT NULL,
    use_type VARCHAR(20) NOT NULL DEFAULT 'MOBILE', -- MOBILE, HOME, WORK, EMERGENCY
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_patient_contacts_lookup ON patient.patient_contacts(system, value);

-- Patient Portal Accounts (Authentication & Session Bindings)
CREATE TABLE IF NOT EXISTS patient.portal_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    identifier VARCHAR(255) NOT NULL, -- Phone or Email
    password_hash VARCHAR(255),       -- Optional if passwordless/PIN enabled
    pin_hash VARCHAR(255),            -- 4-6 digit quick login PIN
    status VARCHAR(30) NOT NULL DEFAULT 'INVITED', -- INVITED, PENDING_VERIFICATION, ACTIVE, LOCKED, SUSPENDED
    mfa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_portal_account_tenant_id UNIQUE (tenant_id, identifier)
);

-- Patient Consents & Access Directives
CREATE TABLE IF NOT EXISTS patient.consents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    consent_type VARCHAR(50) NOT NULL, -- TELEHEALTH, DATA_SHARING, RESEARCH, PROXY_ACCESS
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, REVOKED, EXPIRED
    granted_by VARCHAR(100) NOT NULL, -- SELF, LEGAL_GUARDIAN, PROXY
    granted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb
);

-- ============================================================================
-- 2. SCHEMA: orchestration (Care Requests, Triage & Provider Matching)
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS orchestration;

-- Care Requests
CREATE TABLE IF NOT EXISTS orchestration.care_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    request_number VARCHAR(64) NOT NULL,
    service_type VARCHAR(50) NOT NULL DEFAULT 'GENERAL_CONSULTATION', -- GENERAL_CONSULTATION, SPECIALIST, LAB_TEST, REFILL, TELEHEALTH
    preferred_mode VARCHAR(30) NOT NULL DEFAULT 'IN_PERSON', -- IN_PERSON, VIDEO, AUDIO, ASYNC_CHAT
    urgency VARCHAR(20) NOT NULL DEFAULT 'ROUTINE', -- ROUTINE, URGENT, EMERGENCY
    status VARCHAR(30) NOT NULL DEFAULT 'SUBMITTED', -- SUBMITTED, TRIAGED, ASSIGNED_AGENT, MATCHED, IN_PROGRESS, COMPLETED, CANCELLED
    chief_complaint TEXT,
    symptoms_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    preferred_time_window JSONB,
    assigned_care_agent_id UUID,
    matched_provider_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_care_request_number UNIQUE (tenant_id, request_number)
);

-- Clinical Triage Assessments
CREATE TABLE IF NOT EXISTS orchestration.triage_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    care_request_id UUID NOT NULL REFERENCES orchestration.care_requests(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    assessor_id UUID, -- Staff user id or NULL if system-rule triaged
    acuity_level VARCHAR(20) NOT NULL DEFAULT 'GREEN', -- RED (Immediate), YELLOW (Urgent), GREEN (Standard)
    systolic_bp INT,
    diastolic_bp INT,
    pulse_rate INT,
    temperature NUMERIC(4,1),
    spo2 INT,
    respiratory_rate INT,
    pain_score INT,
    triage_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Care Coordination Cases & Agent Logs
CREATE TABLE IF NOT EXISTS orchestration.care_coordination_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    care_request_id UUID NOT NULL REFERENCES orchestration.care_requests(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    agent_id UUID NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, ESCALATED, RESOLVED
    intake_verified BOOLEAN NOT NULL DEFAULT FALSE,
    insurance_verified BOOLEAN NOT NULL DEFAULT FALSE,
    summary_notes TEXT,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);

-- Provider Availability & Matching Profiles
CREATE TABLE IF NOT EXISTS orchestration.provider_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    specialty_code VARCHAR(50) NOT NULL,
    sub_specialties TEXT[],
    telehealth_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    in_person_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    max_active_queue INT NOT NULL DEFAULT 10,
    current_active_queue INT NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'ON_DUTY', -- ON_DUTY, ON_BREAK, OFF_DUTY, BUSY
    consultation_languages TEXT[] DEFAULT ARRAY['English'],
    rating NUMERIC(3,2) DEFAULT 5.0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_provider_tenant_user UNIQUE (tenant_id, user_id)
);

-- ============================================================================
-- 3. SCHEMA: encounter (Unified Clinical Encounters & Care Journeys)
-- ============================================================================
CREATE SCHEMA IF NOT EXISTS encounter;

-- Clinical Encounters
CREATE TABLE IF NOT EXISTS encounter.encounters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    care_request_id UUID REFERENCES orchestration.care_requests(id) ON DELETE SET NULL,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL,
    encounter_type VARCHAR(50) NOT NULL DEFAULT 'OUTPATIENT', -- OUTPATIENT, EMERGENCY, TELEHEALTH, INPATIENT_ROUND
    mode VARCHAR(30) NOT NULL DEFAULT 'IN_PERSON', -- IN_PERSON, VIDEO, AUDIO, ASYNC_CHAT
    status VARCHAR(30) NOT NULL DEFAULT 'IN_PROGRESS', -- WAITING, IN_PROGRESS, ON_HOLD, COMPLETED, CANCELLED
    chief_complaint TEXT,
    subjective TEXT,
    objective TEXT,
    assessment TEXT,
    plan TEXT,
    primary_diagnosis_code VARCHAR(50),
    primary_diagnosis_name VARCHAR(255),
    secondary_diagnoses JSONB NOT NULL DEFAULT '[]'::jsonb,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Telehealth Virtual Sessions
CREATE TABLE IF NOT EXISTS encounter.telehealth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID NOT NULL REFERENCES encounter.encounters(id) ON DELETE CASCADE,
    room_sid VARCHAR(255) NOT NULL,
    channel_type VARCHAR(30) NOT NULL DEFAULT 'WEBRTC_VIDEO', -- WEBRTC_VIDEO, WEBRTC_AUDIO, SECURE_CHAT
    status VARCHAR(30) NOT NULL DEFAULT 'WAITING_ROOM', -- WAITING_ROOM, CONNECTED, DISCONNECTED, ENDED
    patient_joined_at TIMESTAMPTZ,
    provider_joined_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    recording_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Curexal Care Journey Ledger
CREATE TABLE IF NOT EXISTS encounter.care_journey_milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    encounter_id UUID REFERENCES encounter.encounters(id) ON DELETE SET NULL,
    stage_code VARCHAR(50) NOT NULL, -- INTAKE, TRIAGE, CONSULTATION, LAB_WORKLIST, RADIOLOGY_STUDY, PHARMACY_DISPENSE, SETTLEMENT, FOLLOW_UP
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING', -- PENDING, IN_PROGRESS, COMPLETED, SKIPPED, BLOCKED
    blocking_reason TEXT,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP SCHEMA IF EXISTS encounter CASCADE;
DROP SCHEMA IF EXISTS orchestration CASCADE;
DROP SCHEMA IF EXISTS patient CASCADE;
