-- +goose Up
-- ==============================================================================
-- CUREXAL MIGRATION 000071: FEATURE #24 PRODUCTION AUDIT LEDGER IMMUTABILITY (HIPAA § 164.312)
-- ==============================================================================

-- 1. Ensure HIPAA § 164.312 / § 164.528 Audit Columns Exist
ALTER TABLE audit.audit_events
ADD COLUMN IF NOT EXISTS facility_branch_id UUID,
ADD COLUMN IF NOT EXISTS patient_id UUID,
ADD COLUMN IF NOT EXISTS actor_email VARCHAR(255),
ADD COLUMN IF NOT EXISTS is_break_glass BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS prev_record_hash CHAR(64),
ADD COLUMN IF NOT EXISTS record_hash CHAR(64);

-- 2. Strict Immutability Trigger: Reject any UPDATE or DELETE operations on the audit ledger
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION audit.prevent_tampering()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Security violation: audit.audit_events records are append-only and cannot be updated or deleted.';
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS trg_prevent_audit_tampering ON audit.audit_events;
CREATE TRIGGER trg_prevent_audit_tampering
BEFORE UPDATE OR DELETE ON audit.audit_events
FOR EACH ROW EXECUTE FUNCTION audit.prevent_tampering();

-- 3. High-Performance Query & Compliance Disclosure Indexes
CREATE INDEX IF NOT EXISTS idx_audit_patient_id ON audit.audit_events(patient_id) WHERE patient_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_audit_branch_id ON audit.audit_events(facility_branch_id) WHERE facility_branch_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_audit_actor_id ON audit.audit_events(actor_id);
CREATE INDEX IF NOT EXISTS idx_audit_break_glass ON audit.audit_events(is_break_glass) WHERE is_break_glass = TRUE;
CREATE INDEX IF NOT EXISTS idx_audit_occurred_at ON audit.audit_events(occurred_at DESC);

-- +goose Down
DROP TRIGGER IF EXISTS trg_prevent_audit_tampering ON audit.audit_events;
DROP FUNCTION IF EXISTS audit.prevent_tampering();

DROP INDEX IF EXISTS audit.idx_audit_patient_id;
DROP INDEX IF EXISTS audit.idx_audit_branch_id;
DROP INDEX IF EXISTS audit.idx_audit_actor_id;
DROP INDEX IF EXISTS audit.idx_audit_break_glass;
DROP INDEX IF EXISTS audit.idx_audit_occurred_at;

ALTER TABLE audit.audit_events
DROP COLUMN IF EXISTS facility_branch_id,
DROP COLUMN IF EXISTS patient_id,
DROP COLUMN IF EXISTS actor_email,
DROP COLUMN IF EXISTS is_break_glass,
DROP COLUMN IF EXISTS prev_record_hash,
DROP COLUMN IF EXISTS record_hash;
