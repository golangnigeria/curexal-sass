# CUREXAL CLINIC OS — DAY 3 PRODUCTION SPECIFICATION & IMPLEMENTATION PLAN

**Document**: `docs/Production_plan/day3.md`  
**Execution Day**: Day 3 of 35 (Saturday, September 5, 2026)  
**Target Milestone**: 🚀 Saturday Live Launch — Unified Clinical Care Loop & Cashier POS Checkout  
**Compliance Standard**: HIPAA Security Rule (§ 164.312), NDPR 2019, HL7 FHIR R4 (Encounter, Condition, MedicationRequest), OWASP ASVS Level 2, Curexal Constitution v2.0  
**Strict Policy**: **NO SHORTCUTS.** Zero unverified clinical writes, single canonical encounter model for both in-person and telehealth channels, deterministic ICD-10 diagnostic coding, atomic cashier settlement with immutable audit trail.

[![Open 35-Day Master Plan](https://img.shields.io/badge/◀_Back_to-35--Day_Master_Plan-0284c7?style=for-the-badge)](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/docs/Production_plan/30-DAY-PRODUCTION-PLAN.md)
[![Open Day 2 Spec](https://img.shields.io/badge/◀_View-Day_2_Spec-64748b?style=for-the-badge)](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/docs/Production_plan/day2.md)

---

## 1. OBJECTIVE & DELIVERABLES

Day 3 establishes the operational heart of Curexal Clinic OS: the end-to-end clinical consultation encounter, medical charting, diagnostic coding, electronic prescribing, and automated cashier point-of-sale (POS) settlement across both **In-Person** and **Telehealth** delivery channels:

1. **Feature #11: Clinical Encounter Lifecycle**: Start, pause, resume, and sign/complete clinical encounters with auto-assigned UUIDs, linking back to `appointments` and `care_requests` without fragmented parallel medical records.
2. **Feature #12: Structured Clinical Documentation (SOAP Canvas)**: Unified Physician charting canvas:
   - **Subjective**: Chief complaint, history of presenting illness (HPI), review of systems.
   - **Objective**: Physical examination, triage vitals, systemic observations.
   - **Assessment**: Clinical impressions and physician rationale.
   - **Plan**: Diagnostic investigations, therapeutic orders, patient advice.
3. **Feature #13: Vital Signs & Provenance-Aware Observations**: Automated BMI calculation ($\text{Weight} / \text{Height}^2$), pediatric and adult vitals ranges, and provenance classification (`staff_measured`, `patient_reported`, `device_reported`, `provider_observed`).
4. **Feature #14: Diagnosis Management & ICD-10 Coding**: Regulatory diagnostic classification with primary condition, secondary comorbidities, provisional vs. confirmed verification status.
5. **Feature #16: Digital Medication & e-Prescriptions**: Structured prescription orders (Drug Name, Dosage Form, Strength, Route, Frequency, Duration, Dispensing Instructions, Refill count).
6. **Feature #26: Automated Billing & Cashier POS Checkout**: Auto-generation of patient invoice from consultation fees + billable procedures + diagnostic orders; cashier shift settlement with multi-tender split payment (Cash, POS Card, Bank Transfer).

---

## 2. TECHNICAL SPECIFICATION (DATA SCHEMAS & CONTRACTS)

### Unified Encounter & Clinical Delivery Core

```text
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                   CUREXAL UNIFIED ENCOUNTER & CLINICAL CARE LOOP                       │
└────────────────────────────────────────────────────────────────────────────────────────┘

        [CARE DESK LIVE QUEUE BOARD] (from Day 2 Intake)
                     │
         ┌───────────┴───────────┐
         │ (In-Person Room 104)  │ (Virtual WebRTC Room)
         ▼                       ▼
    Admit Patient           Join Patient
         │                       │
         └───────────┬───────────┘
                     ▼
  ┌─────────────────────────────────────────────────────────────────────────────────────┐
  │                           CANONICAL ENCOUNTER CREATION                              │
  │   POST /api/v1/encounters                                                           │
  │   encounter_channel: in_person | video | telephone | secure_message                 │
  │   status: open -> in_progress                                                       │
  └──────────────────────────────────────┬──────────────────────────────────────────────┘
                                         │
                                         ▼
  ┌─────────────────────────────────────────────────────────────────────────────────────┐
  │                        UNIFIED CLINICAL CONSULTATION CANVAS                         │
  │   Left Viewport (Input Channel Adapter):                                            │
  │     • In-Person: Triage Vitals Summary, Physical Exam Checklist                     │
  │     • Telehealth: Encrypted WebRTC Two-Way Video Stream, Connection Watchdog        │
  │                                                                                     │
  │   Main Charting Pane (Shared EMR Core):                                             │
  │     • Subjective: HPI, Review of Systems, Chief Complaint                           │
  │     • Objective: Blood Pressure, Pulse, Temp, SpO2, BMI, Observations               │
  │     • Assessment: ICD-10 Diagnoses (Primary & Secondary Comorbidities)              │
  │     • Plan: e-Prescriptions, Lab Requisitions, Radiology Studies, Follow-up         │
  └──────────────────────────────────────┬──────────────────────────────────────────────┘
                                         │
                             Provider Signs & Completes
                                         │
                                         ▼
  ┌─────────────────────────────────────────────────────────────────────────────────────┐
  │                      AUTOMATED BILLING INVOICE GENERATOR                            │
  │   Auto-computes: Consultation Fee + Clinic Procedures + e-Prescriptions             │
  │   Status: UNPAID / AWAITING_PAYMENT                                                 │
  └──────────────────────────────────────┬──────────────────────────────────────────────┘
                                         │
                                         ▼
  ┌─────────────────────────────────────────────────────────────────────────────────────┐
  │                            CASHIER POINT-OF-SALE (POS)                              │
  │   Route: /:branchSlug/billing                                                       │
  │   Tenders: Cash + POS Card Terminal + Bank Transfer (Split-Tender Capable)          │
  │   Status -> PAID • Printable Thermal Receipt • Digital Copy to Patient Portal       │
  └─────────────────────────────────────────────────────────────────────────────────────┘
```

---

### 2.1 Feature #11: Clinical Encounter Lifecycle Specification

#### Database Schema: `encounter.encounters`
```sql
CREATE SCHEMA IF NOT EXISTS encounter;

CREATE TYPE encounter_channel AS ENUM ('in_person', 'video', 'telephone', 'secure_message');
CREATE TYPE encounter_status AS ENUM ('open', 'in_progress', 'awaiting_documentation', 'signed', 'closed', 'amended');

CREATE TABLE IF NOT EXISTS encounter.encounters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    appointment_id UUID REFERENCES operations.appointments(id) ON DELETE SET NULL,
    care_request_id UUID REFERENCES orchestration.care_requests(id) ON DELETE SET NULL,
    provider_id UUID NOT NULL REFERENCES orchestration.provider_profiles(id) ON DELETE RESTRICT,
    encounter_channel encounter_channel NOT NULL DEFAULT 'in_person',
    status encounter_status NOT NULL DEFAULT 'open',
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
```

---

### 2.2 Feature #13: Clinical Observations & Vital Signs

#### Database Schema: `encounter.observations`
```sql
CREATE TABLE IF NOT EXISTS encounter.observations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID NOT NULL REFERENCES encounter.encounters(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL, -- e.g. "systolic_bp", "pulse_rate", "temperature", "bmi"
    value_numeric NUMERIC(8,2),
    value_text TEXT,
    unit VARCHAR(50),
    source VARCHAR(50) NOT NULL DEFAULT 'staff_measured', -- staff_measured, patient_reported, device_reported, provider_observed
    observed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    recorded_by UUID REFERENCES identity.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_observations_encounter ON encounter.observations(encounter_id);
CREATE INDEX IF NOT EXISTS idx_observations_patient_code ON encounter.observations(patient_id, code);
```

---

### 2.3 Feature #14: Diagnosis Management & ICD-10 Classification

#### Database Schema: `encounter.diagnoses`
```sql
CREATE TABLE IF NOT EXISTS encounter.diagnoses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID NOT NULL REFERENCES encounter.encounters(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    icd10_code VARCHAR(30) NOT NULL,
    icd10_title VARCHAR(255) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    clinical_status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, RECURRENCE, RESOLVED
    verification_status VARCHAR(30) NOT NULL DEFAULT 'CONFIRMED', -- PROVISIONAL, DIFFERENTIAL, CONFIRMED
    notes TEXT,
    diagnosed_by UUID REFERENCES identity.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_diagnoses_encounter ON encounter.diagnoses(encounter_id);
CREATE INDEX IF NOT EXISTS idx_diagnoses_patient ON encounter.diagnoses(patient_id);
```

---

### 2.4 Feature #16: Digital Medication & e-Prescriptions

#### Database Schema: `encounter.prescriptions` & `encounter.prescription_items`
```sql
CREATE TABLE IF NOT EXISTS encounter.prescriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID NOT NULL REFERENCES encounter.encounters(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    prescriber_id UUID NOT NULL REFERENCES orchestration.provider_profiles(id) ON DELETE RESTRICT,
    prescription_number VARCHAR(64) NOT NULL UNIQUE,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING_DISPENSE', -- PENDING_DISPENSE, PARTIALLY_DISPENSED, DISPENSED, CANCELLED
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS encounter.prescription_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prescription_id UUID NOT NULL REFERENCES encounter.prescriptions(id) ON DELETE CASCADE,
    drug_name VARCHAR(255) NOT NULL,
    dosage_form VARCHAR(50) NOT NULL, -- TABLET, CAPSULE, SYRUP, INJECTION, OINTMENT
    strength VARCHAR(50),             -- e.g. "500mg", "10mg/5ml"
    route VARCHAR(50) DEFAULT 'ORAL', -- ORAL, INTRAVENOUS, INTRAMUSCULAR, TOPICAL
    frequency VARCHAR(50) NOT NULL,   -- e.g. "BD (twice daily)", "TDS (thrice daily)"
    duration_days INT NOT NULL DEFAULT 3,
    quantity_prescribed INT NOT NULL,
    instructions TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_prescriptions_patient ON encounter.prescriptions(patient_id);
```

---

### 2.5 Feature #26: Automated Invoicing & Cashier POS Settlement

#### Database Schema: `billing.invoices` & `billing.payments`
```sql
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
    status VARCHAR(30) NOT NULL DEFAULT 'UNPAID', -- UNPAID, PARTIALLY_PAID, PAID, VOIDED
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS billing.invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES billing.invoices(id) ON DELETE CASCADE,
    description VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL, -- CONSULTATION, PROCEDURE, LAB, PHARMACY, OTHER
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
    tender_type VARCHAR(30) NOT NULL, -- CASH, CARD_POS, BANK_TRANSFER, SPLIT
    tender_breakdown JSONB DEFAULT '{}'::jsonb,
    cashier_id UUID REFERENCES identity.users(id),
    status VARCHAR(30) NOT NULL DEFAULT 'SETTLED',
    paid_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 3. STEP-BY-STEP IMPLEMENTATION RUNBOOK

### Step 1: Database Migration & Schema Expansion
1. Apply migration `000075_day3_clinical_encounter_prescriptions_and_pos_schema.sql`.
2. Verify foreign key references from `encounters` to `patients`, `appointments`, and `care_requests`.
3. Verify `billing.invoices` can link back to `encounters`.

### Step 2: Backend Modules & Handlers
1. In `apps/api/internal/modules/encounter`:
   - Enforce Start Encounter command handling `in_person` and `video` channels.
   - Enforce Save SOAP Notes and ICD-10 diagnostic tagging.
   - Enforce e-Prescribing with quantity calculation.
   - Enforce Complete Encounter which triggers auto-invoice generation.
2. In `apps/api/internal/modules/billing`:
   - Create invoice automatically upon encounter conclusion or service ordering.
   - Handle Cashier POS settlement with split tenders.

### Step 3: Frontend Integration
1. In `apps/web-platform/src/pages/workspace/clinical/index.tsx`:
   - Connect active encounter to live Care Desk queue item.
   - Support SOAP note auto-save.
   - Enable ICD-10 search and prescription builder.
   - Add one-click "Complete & Invoice".
2. In `apps/web-platform/src/pages/workspace/billing/index.tsx`:
   - Live cashier register listing unpaid consultation invoices.
   - Tender breakdown calculator (Cash, Card, Transfer).
   - Printable receipt dialog.

---

## 4. ACCEPTANCE CRITERIA & TESTING PROTOCOL

1. **Encounter State Machine Test**: Open $\to$ In Progress $\to$ Signed/Closed.
2. **Clinical SOAP Note Test**: Persist Subjective, Objective, Assessment, Plan notes and verify audit log.
3. **ICD-10 Diagnostic Tagging Test**: Primary diagnosis saved with ICD-10 code.
4. **e-Prescription Dispatch Test**: Prescriptions created with itemized dosages and instructions.
5. **Cashier POS Checkout Test**: Auto-invoice created for patient, settled via Cash/Card POS with status `PAID`.

---

**Day 3 Plan Certified for Saturday Live Launch.**
