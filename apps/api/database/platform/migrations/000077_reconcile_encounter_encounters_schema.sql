-- +goose Up
-- ==============================================================================
-- CUREXAL MIGRATION 000077: RECONCILE ENCOUNTER.ENCOUNTERS SCHEMA
-- (Safe forward-only schema evolution resolving Day 3 schema delta)
-- ==============================================================================

-- 1. Add missing application contract columns to encounter.encounters
ALTER TABLE encounter.encounters 
    ADD COLUMN IF NOT EXISTS appointment_id UUID REFERENCES operations.appointments(id) ON DELETE SET NULL;

ALTER TABLE encounter.encounters 
    ADD COLUMN IF NOT EXISTS encounter_channel VARCHAR(30) NOT NULL DEFAULT 'in_person';

ALTER TABLE encounter.encounters 
    ADD COLUMN IF NOT EXISTS signed_at TIMESTAMPTZ;

ALTER TABLE encounter.encounters 
    ADD COLUMN IF NOT EXISTS signed_by UUID REFERENCES identity.users(id) ON DELETE SET NULL;

ALTER TABLE encounter.encounters 
    ADD COLUMN IF NOT EXISTS closed_at TIMESTAMPTZ;

-- 2. Enforce check constraint for strongly-typed delivery channel
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'chk_encounters_channel'
    ) THEN
        ALTER TABLE encounter.encounters 
        ADD CONSTRAINT chk_encounters_channel 
        CHECK (encounter_channel IN ('in_person', 'video', 'telephone', 'secure_message'));
    END IF;
END $$;
-- +goose StatementEnd

-- 3. Create high-performance foreign-key & lifecycle indexes
CREATE INDEX IF NOT EXISTS idx_encounters_appointment ON encounter.encounters(appointment_id);
CREATE INDEX IF NOT EXISTS idx_encounters_channel ON encounter.encounters(encounter_channel);
CREATE INDEX IF NOT EXISTS idx_encounters_signed_at ON encounter.encounters(signed_at);
CREATE INDEX IF NOT EXISTS idx_encounters_closed_at ON encounter.encounters(closed_at);

-- +goose Down
DROP INDEX IF EXISTS encounter.idx_encounters_closed_at;
DROP INDEX IF EXISTS encounter.idx_encounters_signed_at;
DROP INDEX IF EXISTS encounter.idx_encounters_channel;
DROP INDEX IF EXISTS encounter.idx_encounters_appointment;
ALTER TABLE encounter.encounters DROP CONSTRAINT IF EXISTS chk_encounters_channel;
ALTER TABLE encounter.encounters DROP COLUMN IF EXISTS closed_at;
ALTER TABLE encounter.encounters DROP COLUMN IF EXISTS signed_by;
ALTER TABLE encounter.encounters DROP COLUMN IF EXISTS signed_at;
ALTER TABLE encounter.encounters DROP COLUMN IF EXISTS encounter_channel;
ALTER TABLE encounter.encounters DROP COLUMN IF EXISTS appointment_id;
