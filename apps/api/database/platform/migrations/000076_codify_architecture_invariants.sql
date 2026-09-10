-- +goose Up
-- ==============================================================================
-- CUREXAL MIGRATION 000076: CODIFY ARCHITECTURE INVARIANTS
-- (Encounter Participants, First-Class Consents, Care Continuations & Telehealth Sessions)
-- ==============================================================================

-- 1. FIRST-CLASS ENCOUNTER CONSENTS (Invariant #10)
CREATE TABLE IF NOT EXISTS encounter.encounter_consents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID REFERENCES encounter.encounters(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    consent_type VARCHAR(50) NOT NULL, -- 'treatment', 'telehealth', 'diagnostic_sharing', 'document_disclosure'
    status VARCHAR(30) NOT NULL DEFAULT 'obtained', -- 'obtained', 'refused', 'revoked'
    obtained_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    obtained_by UUID REFERENCES identity.users(id),
    method VARCHAR(50) NOT NULL DEFAULT 'digital_signature', -- 'digital_signature', 'verbal', 'written_form'
    version VARCHAR(20) NOT NULL DEFAULT '1.0',
    notes TEXT,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_encounter_consents_encounter ON encounter.encounter_consents(encounter_id);
CREATE INDEX IF NOT EXISTS idx_encounter_consents_patient ON encounter.encounter_consents(patient_id);
CREATE INDEX IF NOT EXISTS idx_encounter_consents_type ON encounter.encounter_consents(consent_type, status);

-- 2. FIRST-CLASS ENCOUNTER PARTICIPANTS (Invariant #11)
CREATE TABLE IF NOT EXISTS encounter.encounter_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID NOT NULL REFERENCES encounter.encounters(id) ON DELETE CASCADE,
    person_id UUID NOT NULL,
    participant_type VARCHAR(50) NOT NULL DEFAULT 'patient', -- 'patient', 'provider', 'nurse', 'interpreter', 'caregiver', 'observer', 'support_staff'
    role VARCHAR(50) NOT NULL DEFAULT 'ATTENDEE',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    left_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_encounter_participants_encounter ON encounter.encounter_participants(encounter_id);
CREATE INDEX IF NOT EXISTS idx_encounter_participants_person ON encounter.encounter_participants(person_id);

-- 3. ENHANCE TELEHEALTH SESSIONS (Subordinate to Encounter - Invariant #6 & #7)
-- Rule: Never store raw WebRTC credentials in database; store operational references only.
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS room_reference VARCHAR(255);
UPDATE encounter.telehealth_sessions SET room_reference = room_sid WHERE room_reference IS NULL;
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS session_token_ref VARCHAR(255);
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS consent_id UUID REFERENCES encounter.encounter_consents(id) ON DELETE SET NULL;
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ;
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS last_connected_at TIMESTAMPTZ;
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS disconnect_count INT NOT NULL DEFAULT 0;
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS connection_failure_reason TEXT;
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS termination_reason TEXT;
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS network_quality_summary JSONB DEFAULT '{}'::jsonb;
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS fallback_used BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS recording_enabled BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS recording_reference VARCHAR(255);
ALTER TABLE encounter.telehealth_sessions ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_telehealth_sessions_status ON encounter.telehealth_sessions(status);
CREATE INDEX IF NOT EXISTS idx_telehealth_sessions_room ON encounter.telehealth_sessions(room_reference);

-- 4. CARE CONTINUATIONS (Escalation & Longitudinal Linkage - Invariant #16)
CREATE TABLE IF NOT EXISTS encounter.care_continuations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_encounter_id UUID NOT NULL REFERENCES encounter.encounters(id) ON DELETE RESTRICT,
    target_appointment_id UUID REFERENCES operations.appointments(id) ON DELETE SET NULL,
    reason TEXT NOT NULL,
    clinical_priority VARCHAR(30) NOT NULL DEFAULT 'routine',
    created_by UUID REFERENCES identity.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_care_continuations_source ON encounter.care_continuations(source_encounter_id);
CREATE INDEX IF NOT EXISTS idx_care_continuations_target ON encounter.care_continuations(target_appointment_id);

-- 5. CARE LOCATION TYPE ENHANCEMENT (Invariant #12)
ALTER TABLE operations.appointments ADD COLUMN IF NOT EXISTS care_location_type VARCHAR(30) NOT NULL DEFAULT 'physical_facility';
ALTER TABLE encounter.encounters ADD COLUMN IF NOT EXISTS care_location_type VARCHAR(30) NOT NULL DEFAULT 'physical_facility';

-- +goose Down
DROP TABLE IF EXISTS encounter.care_continuations CASCADE;
DROP TABLE IF EXISTS encounter.encounter_participants CASCADE;
DROP TABLE IF EXISTS encounter.encounter_consents CASCADE;
