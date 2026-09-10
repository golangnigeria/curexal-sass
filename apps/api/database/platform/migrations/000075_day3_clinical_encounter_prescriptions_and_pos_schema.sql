-- +goose Up
-- ==============================================================================
-- CUREXAL MIGRATION 000075: DAY 3 CLINICAL ENCOUNTERS, PRESCRIPTIONS & POS SCHEMA
-- ==============================================================================

-- 1. SCHEMA: encounter
CREATE SCHEMA IF NOT EXISTS encounter;

-- Ensure encounter.encounters table
CREATE TABLE IF NOT EXISTS encounter.encounters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    appointment_id UUID REFERENCES operations.appointments(id) ON DELETE SET NULL,
    care_request_id UUID REFERENCES orchestration.care_requests(id) ON DELETE SET NULL,
    provider_id UUID REFERENCES orchestration.provider_profiles(id) ON DELETE SET NULL,
    encounter_channel VARCHAR(30) NOT NULL DEFAULT 'in_person',
    status VARCHAR(30) NOT NULL DEFAULT 'in_progress',
    chief_complaint TEXT,
    subjective TEXT,
    objective TEXT,
    assessment TEXT,
    plan TEXT,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    signed_at TIMESTAMPTZ,
    signed_by UUID REFERENCES identity.users(id),
    closed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_encounters_tenant_status ON encounter.encounters(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_encounters_patient ON encounter.encounters(patient_id);
CREATE INDEX IF NOT EXISTS idx_encounters_provider ON encounter.encounters(provider_id);
CREATE INDEX IF NOT EXISTS idx_encounters_care_request ON encounter.encounters(care_request_id);

-- 2. SCHEMA: encounter.observations (Vitals & Clinical Observations)
CREATE TABLE IF NOT EXISTS encounter.observations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID NOT NULL REFERENCES encounter.encounters(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL,
    value_numeric NUMERIC(8,2),
    value_text TEXT,
    unit VARCHAR(50),
    source VARCHAR(50) NOT NULL DEFAULT 'staff_measured',
    observed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    recorded_by UUID REFERENCES identity.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_observations_encounter ON encounter.observations(encounter_id);
CREATE INDEX IF NOT EXISTS idx_observations_patient_code ON encounter.observations(patient_id, code);

-- 3. SCHEMA: encounter.diagnoses (ICD-10 Diagnostic Classification)
CREATE TABLE IF NOT EXISTS encounter.diagnoses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID NOT NULL REFERENCES encounter.encounters(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    icd10_code VARCHAR(30) NOT NULL,
    icd10_title VARCHAR(255) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    clinical_status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    verification_status VARCHAR(30) NOT NULL DEFAULT 'CONFIRMED',
    notes TEXT,
    diagnosed_by UUID REFERENCES identity.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_diagnoses_encounter ON encounter.diagnoses(encounter_id);
CREATE INDEX IF NOT EXISTS idx_diagnoses_patient ON encounter.diagnoses(patient_id);

-- 4. SCHEMA: encounter.prescriptions & prescription_items
CREATE TABLE IF NOT EXISTS encounter.prescriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID NOT NULL REFERENCES encounter.encounters(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    prescriber_id UUID REFERENCES orchestration.provider_profiles(id) ON DELETE SET NULL,
    prescription_number VARCHAR(64) NOT NULL UNIQUE,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING_DISPENSE',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS encounter.prescription_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prescription_id UUID NOT NULL REFERENCES encounter.prescriptions(id) ON DELETE CASCADE,
    drug_name VARCHAR(255) NOT NULL,
    dosage_form VARCHAR(50) NOT NULL,
    strength VARCHAR(50),
    route VARCHAR(50) DEFAULT 'ORAL',
    frequency VARCHAR(50) NOT NULL,
    duration_days INT NOT NULL DEFAULT 3,
    quantity_prescribed INT NOT NULL DEFAULT 1,
    instructions TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_prescriptions_patient ON encounter.prescriptions(patient_id);
CREATE INDEX IF NOT EXISTS idx_prescription_items_rx ON encounter.prescription_items(prescription_id);

-- 5. SCHEMA: billing.invoices & billing.payments
CREATE SCHEMA IF NOT EXISTS billing;

CREATE TABLE IF NOT EXISTS billing.invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    encounter_id UUID REFERENCES encounter.encounters(id) ON DELETE SET NULL,
    invoice_number VARCHAR(64) NOT NULL UNIQUE,
    subtotal NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    discount_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    tax_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    total_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    amount_paid NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    balance_due NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    status VARCHAR(30) NOT NULL DEFAULT 'UNPAID',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS billing.invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES billing.invoices(id) ON DELETE CASCADE,
    description VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL,
    unit_price NUMERIC(12,2) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    total_price NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS billing.payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    invoice_id UUID NOT NULL REFERENCES billing.invoices(id) ON DELETE RESTRICT,
    receipt_number VARCHAR(64) NOT NULL UNIQUE,
    amount NUMERIC(12,2) NOT NULL,
    tender_type VARCHAR(30) NOT NULL,
    tender_breakdown JSONB DEFAULT '{}'::jsonb,
    cashier_id UUID REFERENCES identity.users(id),
    status VARCHAR(30) NOT NULL DEFAULT 'SETTLED',
    paid_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_invoices_tenant_status ON billing.invoices(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_invoices_patient ON billing.invoices(patient_id);
CREATE INDEX IF NOT EXISTS idx_payments_invoice ON billing.payments(invoice_id);

-- +goose Down
DROP TABLE IF EXISTS billing.payments CASCADE;
DROP TABLE IF EXISTS billing.invoice_items CASCADE;
DROP TABLE IF EXISTS billing.invoices CASCADE;
DROP TABLE IF EXISTS encounter.prescription_items CASCADE;
DROP TABLE IF EXISTS encounter.prescriptions CASCADE;
DROP TABLE IF EXISTS encounter.diagnoses CASCADE;
DROP TABLE IF EXISTS encounter.observations CASCADE;
DROP TABLE IF EXISTS encounter.encounters CASCADE;
