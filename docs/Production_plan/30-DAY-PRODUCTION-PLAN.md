# CUREXAL CLINIC SOFTWARE — 35-DAY DAILY PRODUCTION DEPLOYMENT PLAN

[![Open Day 1 Spec](https://img.shields.io/badge/▶_Open-Day_1_Detailed_Spec-0284c7?style=for-the-badge)](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/docs/Production_plan/day1.md)
[![Open Day 2 Spec](https://img.shields.io/badge/▶_Open-Day_2_Detailed_Spec-0284c7?style=for-the-badge)](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/docs/Production_plan/day2.md)
[![Jump to Saturday Launch](https://img.shields.io/badge/🚀_Jump_to-Saturday_Launch_Milestone-16a34a?style=for-the-badge)](#day-3--saturday-live-launch--clinical-care-loop--pos-checkout-saturday-sep-5-2026)
[![Platform Suite](https://img.shields.io/badge/🛡️_Jump_to-Platform_Super--Admin_Suite-7c3aed?style=for-the-badge)](#phase-5-platform-operations--super-admin-suite-days-25--28)
[![Deployment Protocol](https://img.shields.io/badge/⚡_Jump_to-Daily_Deploy_Protocol-ea580c?style=for-the-badge)](#daily-production-deployment-protocol)

---

**Target Launch Date**: Saturday, September 5, 2026 (Day 3 / Milestone 1)  
**Execution Cycle**: September 3, 2026 – October 7, 2026 (35 Continuous Production Days)  
**Scope**: 40 Production Clinical Features + Platform Super-Admin Governance Suite Across 7 Phased Modules  
**Stack**: Go 1.25 + Echo (`apps/api`), PostgreSQL 16 (`platform` + tenant schemas), React 19 + TypeScript (`apps/web-platform`, `apps/web-patient`, `apps/web-public`)

---

## MILESTONE 1 (THE SATURDAY LIVE LAUNCH LOOP)

```text
[1. SECURE LOGIN] ──► [2. CLINIC SELECT] ──► [3. PATIENT INTAKE] ──► [4. APPOINTMENT / QUEUE]
                                                                            │
[8. AUDIT LOG] ◄── [7. POS BILLING] ◄── [6. RX & FINALIZATION] ◄── [5. SOAP CONSULT & ICD-10]
```

---

## PHASE 1: FOUNDATION & THE SATURDAY LAUNCH (DAYS 1 – 3)

### Day 1: Foundation Lockdown & Security Verification (Thursday, Sep 3, 2026)
* **Features**:
  - [x] **#1. Secure authentication**: Argon2id password hashing, JWT session expiration, CSRF protection, active session token management.
  - [x] **#2. Role-based access control (RBAC)**: Enforce granular permissions across Doctor, Nurse, Receptionist, Cashier, Admin, Owner, and Platform Super-Admin roles (`super_admin`, `super_support_agent`, `super_sales_staff`).
  - [x] **#3. Organization & clinic management**: Organization governance, branch facilities, departments, settings isolation.
  - [x] **#24. Audit logging**: Capture auth events, tenant context, IP, user-agent, and failed access attempts in `audit.audit_events`.
* **Daily Release Scope**:
  1. Verify token signing and refresh rotation in `apps/api/internal/modules/identity`.
  2. Verify DB-driven permission provider (`authorization.permissions`).
  3. Validate tenant schema isolation between test organizations.
* **Daily Release Verification**:
  ```powershell
  cd apps/api; go test -v ./internal/testing/...
  bun x tsc --noEmit
  ```
[![View Day 1 Spec](https://img.shields.io/badge/Inspect-Day_1_Spec_File-0284c7?style=flat-square)](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/docs/Production_plan/day1.md)

---

### Day 2: Patient Intake, Identity Matching & Live Queue (Friday, Sep 4, 2026)
* **Features**:
  - [x] **#4. Staff and provider management**: Link staff accounts to clinical provider profiles (specialty, license, room assignment).
  - [x] **#5. Patient registration**: Full demographic capture, contact details, emergency contacts, next of kin, MRN generation (`MRN-YYYY-XXXX`).
  - [x] **#6. Patient identity matching & duplicate detection**: Instant duplicate match on `Phone` + `DOB` + `First/Last Name`.
  - [x] **#8. Appointment scheduling**: Provider-linked calendar booking with department and time slot selection.
  - [x] **#10. Patient check-in & queue management**: Real-time status transitions: `Registered` $\to$ `Waiting (Triage)` $\to$ `In Consultation` $\to$ `Completed`.
* **Daily Release Scope**:
  1. Test walk-in patient registration in `/reception` workspace.
  2. Verify duplicate detection modal appears when duplicate phone number is entered.
  3. Check-in patient and verify real-time queue badge increments on Care Desk.
* **Daily Release Verification**:
  ```powershell
  # Register duplicate patient and confirm validation warning
  # Verify queue displays patient with elapsed wait timer
  ```
[![View Day 2 Spec](https://img.shields.io/badge/Inspect-Day_2_Spec_File-0284c7?style=flat-square)](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/docs/Production_plan/day2.md)

---

### Day 3: 🚀 SATURDAY LIVE LAUNCH — Clinical Care Loop & POS Checkout (Saturday, Sep 5, 2026)
* **Features**:
  - [x] **#11. Encounter management**: Start, pause, resume, and complete clinical encounters with auto-assigned encounter UUIDs.
  - [x] **#12. Clinical documentation**: Structured SOAP canvas (Subjective, Objective, Assessment, Plan) with clinical notes.
  - [x] **#13. Vital signs & observations**: Triage vitals intake (BP, Pulse, Temp, SpO2, Weight, Height, automated BMI calculation).
  - [x] **#14. Diagnosis management**: ICD-10 coding modal with primary/secondary condition tagging.
  - [x] **#16. Medication & prescription management**: Digital e-prescription creation (Drug, Dose, Route, Frequency, Duration).
  - [x] **#26. Billing and invoicing**: Auto-generate invoice from consultation fee + clinic procedures; cashier POS settlement (Cash, Card, Transfer).
* **Daily Release Scope**:
  1. Complete live clinical encounter from Triage to Doctor SOAP notes.
  2. Issue e-prescription and add ICD-10 code.
  3. Cashier settles invoice at POS register and generates printable receipt.
  4. Patient receives encounter summary and receipt on `/patient` portal.
* **Daily Release Verification**:
  - Run live pilot encounter end-to-end.
  - Confirm invoice status is `PAID` and encounter status is `COMPLETED`.
  - Tag production release: `git tag -a "v1.0.0-launch" -m "Saturday Launch Milestone 1"`.

---

## PHASE 2: CLINICAL INTEGRITY & SAFETY (DAYS 4 – 10)

### Day 4: Post-Launch Stability & Disaster Recovery (Sunday, Sep 6, 2026)
* **Features**:
  - [ ] **#29. Backup, recovery, and business continuity**: Automated PostgreSQL daily backup dump, WAL point-in-time recovery (PITR) verification, `/healthz` and `/readyz` uptime probes.
* **Release Scope**: Test database restore script against a staging replica; verify 100% data recovery.

---

### Day 5: Unified Longitudinal Patient Profile (Monday, Sep 7, 2026)
* **Features**:
  - [ ] **#7. Longitudinal patient profile**: Single-pane-of-glass patient summary showing encounter history, past diagnoses, active medications, vital sign trend charts, and payment logs.
* **Release Scope**: Deploy chronological visit timeline and interactive Recharts vital sign trend graphs.

---

### Day 6: Allergy & Adverse Reaction Management (Tuesday, Sep 8, 2026)
* **Features**:
  - [ ] **#15. Allergy and adverse-reaction management**: Drug, food, and environmental allergy tracking with severity levels (`Mild`, `Moderate`, `Life-Threatening`) and verification status (`Suspected` vs `Confirmed`).
* **Release Scope**: Persistent high-visibility red allergy banner across patient drawer, triage desk, and consultation canvas.

---

### Day 7: Clinical Safety Warnings & Decision Support (Wednesday, Sep 9, 2026)
* **Features**:
  - [ ] **#17. Clinical safety warnings**: Drug-allergy cross-check on e-prescribing, duplicate therapeutic medication warnings, abnormal vitals alert flags (e.g. BP $> 160/100$, SpO2 $< 92\%$).
* **Release Scope**: Prescription intercept modal preventing accidental medication orders matching active patient allergies.

---

### Day 8: Record Finalization & Controlled Amendments (Thursday, Sep 10, 2026)
* **Features**:
  - [ ] **#18. Clinical record finalization & amendment control**: Digital doctor sign-off locking completed SOAP notes into read-only mode; mandatory amendment reason logging for post-signoff addenda.
* **Release Scope**: Immutable encounter locking with cryptographic signature hash and append-only addendum table.

---

### Day 9: Document & Clinical Attachment Vault (Friday, Sep 11, 2026)
* **Features**:
  - [ ] **#19. Document and attachment management**: Upload referral letters, paper lab scans, ID cards, and consent forms with secure presigned S3/MinIO URLs and category tagging.
* **Release Scope**: Encrypted attachment storage with in-browser PDF/image preview modal.

---

### Day 10: Provider Availability & Calendar Capacity (Saturday, Sep 12, 2026)
* **Features**:
  - [ ] **#9. Provider availability and calendar management**: Doctor working hours, leave schedules, recurring time blocks, consult duration presets (15m, 30m, 60m).
* **Release Scope**: Real-time slot conflict prevention ensuring appointments can only be booked during active provider availability.

---

## PHASE 3: CONNECTED CARE & DIAGNOSTICS (DAYS 11 – 17)

### Day 11: Electronic Diagnostic Lab Orders (Sunday, Sep 13, 2026)
* **Features**:
  - [ ] **#20. Laboratory and diagnostic requests**: Doctor orders tests directly from consult canvas with clinical indication and requisition ID (`REQ-YYYY-XXXX`).
* **Release Scope**: Laboratory test requisition drawer linked to standard fee catalog.

---

### Day 12: Diagnostic Results Ingestion & Chart Sync (Monday, Sep 14, 2026)
* **Features**:
  - [ ] **#21. Results management**: Ingest test results, compare against biological reference ranges, doctor verification review, and automated attachment to patient chart.
* **Release Scope**: Diagnostic results review workspace with critical value flagging and PDF report generator.

---

### Day 13: Structured Referrals Engine (Tuesday, Sep 15, 2026)
* **Features**:
  - [ ] **#22. Referral management**: Create, send, and track outbound patient referrals to partner diagnostic centers or specialists with attached encounter summaries.
* **Release Scope**: Digital referral builder with printable referral letter and lifecycle status tracking (`Draft` $\to$ `Dispatched` $\to$ `Accepted` $\to$ `Completed`).

---

### Day 14: Consent & Privacy Governance (Wednesday, Sep 16, 2026)
* **Features**:
  - [ ] **#23. Consent and privacy management**: Explicit consent capture for treatment, diagnostic sharing, WhatsApp notifications, and document disclosure compliant with NDPR & HIPAA.
* **Release Scope**: Patient consent audit log with revocable permission checkboxes.

---

### Day 15: Automated Notifications & Reminders (Thursday, Sep 17, 2026)
* **Features**:
  - [ ] **#25. Notifications and reminders**: SMS, WhatsApp, and email alerts for appointment confirmations, 24h reminders, ready lab results, and digital receipts.
* **Release Scope**: Background notification queue worker with Termii/Twilio SMS integration and webhook event listeners.

---

### Day 16: Advanced Cashier POS & Shift Settlement (Friday, Sep 18, 2026)
* **Features**:
  - [ ] **#26. Billing and invoicing (Advanced POS)**: Split tender payments (Cash + Card + Bank Transfer), itemized procedure billing, discounts with manager override, and daily cashier shift reconciliation.
* **Release Scope**: Shift closeout report with opening cash float, total collections by tender, and discrepancy tracking.

---

### Day 17: Payment Gateway Vault & Automated Webhooks (Saturday, Sep 19, 2026)
* **Features**:
  - [ ] **#27. Payment integration**: Paystack, Flutterwave, and Stripe webhook verification, automatic invoice settlement upon payment confirmation, and refund processing.
* **Release Scope**: Cryptographic webhook signature verification and live online bill settlement.

---

## PHASE 4: TELEHEALTH DELIVERY INFRASTRUCTURE (DAYS 18 – 24)
> **Architectural Law**: Days 18–24 implement the **remote care-delivery infrastructure** (WebRTC browser video rooms, virtual waiting room state machine, connection watchdog, in-call controls) on top of the already channel-neutral clinical core — NOT creating a disconnected parallel clinical system. All clinical charting, orders, prescriptions, and billing continue using the single canonical clinical core.

### Day 18: Virtual Appointment Scheduling (Sunday, Sep 20, 2026)
* **Features**:
  - [ ] **#31. Virtual appointment scheduling**: Dedicated telehealth appointment type with auto-generated secure room tokens and pre-consultation instructions.
* **Release Scope**: Separate virtual care calendar view and automated meeting link generation.

---

### Day 19: Browser Video Consultation Room (Monday, Sep 21, 2026)
* **Features**:
  - [ ] **#32. Secure video consultation**: Encrypted WebRTC browser video/audio room (Daily.co / LiveKit adapter) with zero third-party app downloads required.
* **Release Scope**: In-browser two-way video room with low-latency audio.

---

### Day 20: Virtual Waiting Room & Triage Admission (Tuesday, Sep 22, 2026)
* **Features**:
  - [ ] **#33. Virtual waiting room**: Patient holding screen with queue position; doctor dashboard with notification of waiting virtual patients and one-click admission.
* **Release Scope**: Live waiting room state machine managing patient admittance.

---

### Day 21: In-Call Controls & Connection Recovery (Wednesday, Sep 23, 2026)
* **Features**:
  - [ ] **#34. Video consultation controls and connection recovery**: Mic/cam toggle, audio device selector, connection quality indicator, auto-reconnect on network drop, and low-bandwidth audio fallback.
* **Release Scope**: Telehealth connection watchdog preventing call dropouts during cellular network switches.

---

### Day 22: Dual-Screen Virtual Clinical Canvas (Thursday, Sep 24, 2026)
* **Features**:
  - [ ] **#36. Virtual consultation clinical workspace**: Side-by-side split screen with video call on one side and full SOAP charting, vitals, allergy list, and past visits on the other.
* **Release Scope**: Real-time SOAP charting during active video stream with zero interface stutter.

---

### Day 23: Telehealth Prescriptions & Session Audit (Friday, Sep 25, 2026)
* **Features**:
  - [ ] **#35. Telehealth consent and identity verification**: Pre-call photo ID check and explicit digital telehealth consent.
  - [ ] **#37. Telehealth prescriptions and follow-up**: Digital prescription delivery immediately upon call completion with care instructions.
  - [ ] **#39. Telehealth session records and audit trail**: Comprehensive log of call duration, connection timestamps, participants, and clinical outcome.
* **Release Scope**: Session audit recorder logging every telehealth encounter into `audit.audit_events`.

---

### Day 24: Pre-Payment & Patient Messaging (Saturday, Sep 26, 2026)
* **Features**:
  - [ ] **#38. Secure patient-provider messaging**: Encrypted 48-hour follow-up messaging channel attached to the clinical encounter.
  - [ ] **#40. Telehealth billing and payment workflow**: Pre-payment gate requiring settled card/transfer payment before virtual room admission.
* **Release Scope**: Automated payment lock on virtual room entry and post-encounter messaging tab.

---

## PHASE 5: PLATFORM OPERATIONS & SUPER-ADMIN SUITE (DAYS 25 – 28)

### Day 25: Platform Clinic Lifecycle & KYC Verification (Sunday, Sep 27, 2026)
* **Features**:
  - [ ] **Platform Organization Governance** (`/platform/organizations`, `/platform/demo-requests`):
    * Super-Admin verification workflow: Approve, reject, or suspend clinic organizations based on regulatory credentials (MDCN/HEFAMAA).
    * Inbound demo requests queue: Convert marketing leads into provisioned pilot clinics.
    * Platform-level tenant status switches (`ACTIVE`, `SUSPENDED`, `UNDER_REVIEW`).
* **Release Scope**: Deploy platform organization inspection drawer and instant demo provisioning pipeline.

---

### Day 26: Master Medical Reference Catalogs (Monday, Sep 28, 2026)
* **Features**:
  - [ ] **Master Catalog Management** (`/platform/catalogs`):
    * Seed and maintain the global ICD-10 diagnostic code registry.
    * Master medical procedure & consultation catalog templates pushed to new clinic branches.
    * Clinical specialty mapping (General Practice, Pediatrics, Cardiology, Obstetrics).
* **Release Scope**: Deploy master catalog search API with bulk CSV import and versioning.

---

### Day 27: Platform Pricing Rules & Payment Gateway Vault (Tuesday, Sep 29, 2026)
* **Features**:
  - [ ] **Platform Pricing & Secrets Vault** (`/platform/pricing`):
    * Dynamic pricing rules engine: Set monthly and annual rates for `Smart`, `Optimize`, `Pro`, and `Enterprise`.
    * Platform VAT configuration (default 7.5%).
    * Secure encrypted key vault for Paystack, Flutterwave, and Stripe production API secrets and webhook signing keys.
* **Release Scope**: Deploy platform pricing manager with zero code redeploys required to update subscription rates.

---

### Day 28: Support Impersonation & Break-Glass Audit (Wednesday, Sep 30, 2026)
* **Features**:
  - [ ] **Platform Staff Directory & Support Access** (`/platform/users`, `/platform/audit`):
    * Platform staff role management (`super_admin`, `super_support_agent`, `super_sales_staff`, `platform_compliance_auditor`).
    * Break-glass support access: Time-limited (30 min) audited tenant impersonation to debug clinic customer issues.
    * Cross-tenant platform-wide audit ledger query console with NDPR export controls.
* **Release Scope**: Deploy support impersonation token generator with mandatory audit logging (`platform.support.impersonate`).

---

## PHASE 6: ENTERPRISE OPERATIONS & METRICS (DAYS 29 – 31)

### Day 29: Operational Dashboards & Clinic Reporting (Thursday, Oct 1, 2026)
* **Features**:
  - [ ] **#28. Operational dashboards and reporting**: Executive dashboards tracking daily patient volume, queue wait times, provider consultation hours, revenue by payment tender, and top ICD-10 diagnoses.
* **Release Scope**: Visual clinic metrics console with date range filtering and CSV export.

---

### Day 30: OpenAPI Interoperability & Webhooks (Friday, Oct 2, 2026)
* **Features**:
  - [ ] **#30. Integration and interoperability foundation**: Public OpenAPI v3 documentation (`/api/v1/docs`), external webhook dispatchers, and API key management console.
* **Release Scope**: Outbound webhook engine dispatching events to external partner systems.

---

### Day 31: Provider Workload & Scheduling Analytics (Saturday, Oct 3, 2026)
* **Features**:
  - [ ] **#4. Staff & Provider Management (Advanced Analytics)**: Provider workload balance, average encounter duration, no-show rate analytics, and receptionist intake velocity.
* **Release Scope**: Staff performance reporting dashboard for clinic medical directors.

---

## PHASE 7: SCALE, SYSTEM DIAGNOSTICS & GOVERNANCE (DAYS 32 – 35)

### Day 32: Enterprise Session Revocation & Brute-Force Defense (Sunday, Oct 4, 2026)
* **Features**:
  - [ ] **#1. Secure Authentication (Enterprise Hardening)**: Active session manager allowing users to view and remotely revoke login sessions; automated 15-minute account lockout after 5 consecutive failed passwords.
* **Release Scope**: Security settings console with remote device revocation and brute-force IP rate limiting.

---

### Day 33: Tenant Data Portability & Disaster Recovery Testing (Monday, Oct 5, 2026)
* **Features**:
  - [ ] **#29. Backup, Recovery & Business Continuity (Compliance Audit)**: One-click full clinic data export (JSON/CSV archive), automated database re-indexing, and 99.9% uptime SLA monitoring alerts.
* **Release Scope**: Tenant self-service data portability export and PostgreSQL maintenance automation.

---

### Day 34: Platform System Diagnostics & Health Engine (Tuesday, Oct 6, 2026)
* **Features**:
  - [ ] **Infrastructure Diagnostics** (`/platform/diagnostics`):
    * Real-time PostgreSQL connection pool metrics and slow query monitor.
    * Redis cache hit/miss ratio and background queue latency telemetry.
    * Automated uptime probes (`/healthz`, `/readyz`) with incident webhook triggers.
* **Release Scope**: Deploy platform diagnostics dashboard with live service telemetry.

---

### Day 35: 35-Day Milestone Certification & Phase 2 Diagnostic Highway (Wednesday, Oct 7, 2026)
* **Features**:
  - [ ] **Full 40-Feature + Platform Suite Operational Review**: End-to-end audit of all 40 features and platform administration operating smoothly in live production.
  - [ ] **Phase 2 Diagnostic Highway Activation**: Turn on the B2B electronic requisition highway connecting clinics directly to partner pathology laboratories and radiology centers.
* **Release Scope**: Final 35-day production milestone certification with zero critical defects.

---

## DAILY PRODUCTION DEPLOYMENT PROTOCOL

Execute this exact command sequence every day at release time:

```powershell
# 1. Typecheck validation across all web apps
bun x tsc --noEmit

# 2. Run backend test suite
cd apps/api
go test -v ./internal/testing/...

# 3. Verify database migrations
goose -dir database/platform/migrations status

# 4. Production build test
bun run build

# 5. Tag and push release
git tag -a "deploy-YYYY-MM-DD" -m "Daily Production Release: Day [X]"
git push origin "deploy-YYYY-MM-DD"
```

---

## ARCHITECTURAL CONSTITUTION: UNIFIED CLINICAL CORE & MULTI-CHANNEL CARE DELIVERY

> **Fundamental Architectural Law of Curexal Clinic OS**:  
> **Curexal possesses ONE clinical core with MULTIPLE care-delivery channels — NOT separate clinic and telehealth systems.**
> - **Clinic / Organization** = Where / institutional care context
> - **Delivery Channel** = How care is planned and delivered (`in_person`, `video`, `telephone`, `secure_message`)
> - **Appointment** = Planned care
> - **Attendance / Check-in** = Physical queue or virtual waiting room state
> - **Encounter** = Care actually delivered (created only upon care initiation)
> - **Telehealth Session** = Remote communication layer (subordinate to Encounter)
> - **Clinical Record** = Single longitudinal source of truth shared across all channels

---

### 1. The Four Independent State Machines

Curexal strictly avoids overloading `appointment.status` with attendance, real-time connectivity, or clinical charting state. The system enforces four decoupled, orthogonal state machines:

```text
1. APPOINTMENT (Planned Care)
   ├── scheduled
   ├── confirmed
   ├── cancelled
   ├── no_show
   └── completed

2. CHECK-IN / ATTENDANCE (Waiting & Ingestion)
   ├── not_checked_in
   ├── checked_in
   ├── waiting
   └── admitted

3. ENCOUNTER (Clinical Delivery)
   ├── open
   ├── in_progress
   ├── awaiting_documentation
   ├── signed
   ├── closed
   └── amended

4. TELEHEALTH SESSION (Connection & Stream Telemetry)
   ├── not_started
   ├── patient_waiting
   ├── provider_waiting
   ├── both_connected
   ├── in_progress
   ├── disconnected
   ├── completed
   └── failed
```

---

### 2. Care Initiation Gate: When is an Encounter Created?

An encounter represents **actual care delivered**, never speculative access:

1. **In-Person**: Patient check-in $\to$ Triage queue $\to$ Provider admits patient into consultation room $\to$ **Care Initiation Event** creates `encounters` record.
2. **Telehealth**: Patient check-in verifies consent $\to$ Patient placed in virtual waiting room (`telehealth_sessions.status = patient_waiting`). Provider joins $\to$ Both connected $\to$ **Care Initiation Event** atomically creates `encounters` and transitions `telehealth_sessions.status = in_progress`.
   - If the patient waits in the virtual room but the doctor never joins, **no encounter is created**; the appointment is marked `unfulfilled` / `no_show`.
   - If the call drops immediately before clinical consultation occurs, the session is marked `failed`, preventing erroneous billing or ghost clinical records.

---

### 3. Production Canonical Database Model

All clinical data resides in a single canonical schema per tenant:

```sql
-- Strongly typed Delivery Channel
CREATE TYPE delivery_channel AS ENUM (
    'in_person',
    'video',
    'telephone',
    'secure_message'
);

-- Decoupled state machines
CREATE TYPE appointment_status AS ENUM (
    'scheduled',
    'confirmed',
    'cancelled',
    'no_show',
    'completed'
);

CREATE TYPE attendance_status AS ENUM (
    'not_checked_in',
    'checked_in',
    'waiting',
    'admitted'
);

CREATE TYPE encounter_status AS ENUM (
    'open',
    'in_progress',
    'awaiting_documentation',
    'signed',
    'closed',
    'amended'
);

CREATE TYPE telehealth_session_status AS ENUM (
    'not_started',
    'patient_waiting',
    'provider_waiting',
    'both_connected',
    'in_progress',
    'disconnected',
    'completed',
    'failed'
);

CREATE TYPE care_location_type AS ENUM (
    'physical_facility',
    'virtual_service'
);
```

#### Canonical Clinical Tables:

```sql
-- 1. Appointments (Planned Care)
CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id UUID NOT NULL REFERENCES patients(id),
    provider_id UUID NOT NULL,
    care_location_type care_location_type NOT NULL DEFAULT 'physical_facility',
    facility_id UUID REFERENCES organization.facility_branches(id), -- NULL for virtual service
    appointment_type_id UUID NOT NULL,
    delivery_channel delivery_channel NOT NULL,
    status appointment_status NOT NULL DEFAULT 'scheduled',
    scheduled_start TIMESTAMPTZ NOT NULL,
    scheduled_end TIMESTAMPTZ NOT NULL,
    reason TEXT,
    booking_source VARCHAR(50) NOT NULL,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (scheduled_end > scheduled_start)
);

-- 2. Attendance & Queue Management
CREATE TABLE attendance_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL REFERENCES appointments(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    facility_id UUID, -- NULL for virtual service
    status attendance_status NOT NULL DEFAULT 'not_checked_in',
    priority VARCHAR(30) NOT NULL DEFAULT 'normal',
    checked_in_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    called_at TIMESTAMPTZ,
    admitted_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

-- 3. Encounters (Delivered Care Core)
CREATE TABLE encounters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id UUID NOT NULL REFERENCES patients(id),
    appointment_id UUID NOT NULL REFERENCES appointments(id),
    provider_id UUID NOT NULL,
    care_location_type care_location_type NOT NULL DEFAULT 'physical_facility',
    facility_id UUID REFERENCES organization.facility_branches(id), -- NULL for virtual service
    encounter_channel delivery_channel NOT NULL,
    status encounter_status NOT NULL DEFAULT 'open',
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    ended_at TIMESTAMPTZ,
    signed_at TIMESTAMPTZ,
    signed_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 4. Encounter Participants (Multi-party attendance)
CREATE TABLE encounter_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID NOT NULL REFERENCES encounters(id) ON DELETE CASCADE,
    person_id UUID NOT NULL,
    participant_type VARCHAR(50) NOT NULL, -- 'patient', 'provider', 'nurse', 'interpreter', 'caregiver', 'observer', 'support_staff'
    role VARCHAR(50) NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    left_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 5. Encounter Consents (First-Class Consent Governance)
CREATE TABLE encounter_consents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID REFERENCES encounters(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patients(id),
    consent_type VARCHAR(50) NOT NULL, -- 'treatment', 'telehealth', 'diagnostic_sharing', 'document_disclosure'
    status VARCHAR(30) NOT NULL DEFAULT 'obtained', -- 'obtained', 'refused', 'revoked'
    obtained_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    obtained_by UUID,
    method VARCHAR(50) NOT NULL DEFAULT 'digital_signature', -- 'digital_signature', 'verbal', 'written_form'
    version VARCHAR(20) NOT NULL DEFAULT '1.0',
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 6. Telehealth Sessions (Subordinate to Encounter - Operational Telemetry Only)
-- Rule: Do NOT store raw WebRTC credentials/secrets in database.
CREATE TABLE telehealth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID NOT NULL UNIQUE REFERENCES encounters(id) ON DELETE CASCADE,
    room_reference VARCHAR(255) NOT NULL,
    session_token_ref VARCHAR(255),
    status telehealth_session_status NOT NULL DEFAULT 'not_started',
    consent_id UUID REFERENCES encounter_consents(id),
    patient_joined_at TIMESTAMPTZ,
    provider_joined_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    last_connected_at TIMESTAMPTZ,
    disconnect_count INT NOT NULL DEFAULT 0,
    connection_failure_reason TEXT,
    termination_reason TEXT,
    network_quality_summary JSONB,
    fallback_used BOOLEAN NOT NULL DEFAULT FALSE,
    recording_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    recording_reference VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 7. Care Continuations (Escalation & Care Journey Linkage)
CREATE TABLE care_continuations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_encounter_id UUID NOT NULL REFERENCES encounters(id),
    target_appointment_id UUID REFERENCES appointments(id),
    reason TEXT NOT NULL,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 8. Observations with Provenance
CREATE TABLE observations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    encounter_id UUID NOT NULL REFERENCES encounters(id),
    patient_id UUID NOT NULL REFERENCES patients(id),
    code VARCHAR(100) NOT NULL,
    value_numeric NUMERIC,
    value_text TEXT,
    unit VARCHAR(50),
    source VARCHAR(50) NOT NULL, -- staff_measured, patient_reported, device_reported, provider_observed, external_source
    observed_at TIMESTAMPTZ NOT NULL,
    recorded_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

### 4. Master Patient Index (MPI) as a Platform Utility

Patient identification must happen first as a shared platform utility across all clinical touchpoints:

```text
Booking / Walk-in / Telehealth Portal / Labs / Billing
                          │
                          ▼
            Patient Identity / MPI Engine
                          │
     ┌────────────────────┼────────────────────┐
     ▼                    ▼                    ▼
Exact Match          Probable Match         No Match
(Score ≥ 80)          (Score 50-79)        (Score < 50)
     │                    │                    │
Use Existing          Staff Resolve        Register Patient
```

The MPI prevents fragmented patient profiles across branches and channels, maintaining a single longitudinal health record.

---

### 5. Decoupling the Clinical Core from Operational UI Views

The backend clinical core has zero concept of "telehealth SOAP" vs "clinic SOAP". It exposes uniform, channel-neutral REST APIs:

```http
POST /api/v1/encounters/{id}/notes
POST /api/v1/encounters/{id}/observations
POST /api/v1/encounters/{id}/diagnoses
POST /api/v1/encounters/{id}/treatment-plan
POST /api/v1/encounters/{id}/prescriptions
POST /api/v1/encounters/{id}/diagnostic-requests
POST /api/v1/encounters/{id}/follow-ups
```

The frontend adapts its layout to the operational context:
- **In-Person Workspace**: Full-width longitudinal visit timeline, multi-tab SOAP note editor, physical triage vitals intake.
- **Telehealth Workspace**: Split-screen canvas with WebRTC audio/video stream on the left and identical SOAP charting, e-prescribing, and lab orders on the right.

---

### 6. 35-Day Phased Architecture Strategy

- **Days 1 – 17**: Build all primitives as **channel-neutral foundations** (`Patient`, `Appointment`, `Encounter`, `SOAP Notes`, `Observations with Provenance`, `ICD-10`, `e-Prescriptions`, `POS Settlement`, `MPI`, `Consent`).
- **Days 18 – 24**: Implement the **telehealth delivery infrastructure** (WebRTC browser video rooms, virtual waiting rooms, connection watchdog, in-call controls) on top of this already channel-neutral clinical core — NOT creating a second clinical system.

---

### 7. Final Production Relationship Architecture

```text
                         CUREXAL CLINICAL CORE

                                PATIENT
                                   │
                                   ▼
                              APPOINTMENT
                                   │
                    ┌──────────────┴──────────────┐
                    │                             │
               IN-PERSON                       REMOTE
                    │                             │
              Check-in                       Check-in
                    │                             │
                Queue                        Consent
                    │                             │
                Triage                    Waiting Room
                    │                             │
                    └──────────────┬──────────────┘
                                   │
                                   ▼
                              ENCOUNTER
                                   │
                 ┌─────────────────┼─────────────────┐
                 │                 │                 │
                 ▼                 ▼                 ▼
          Clinical Record    Telehealth Session   Participants
                 │
       ┌─────────┼─────────┬──────────┬──────────┐
       ▼         ▼         ▼          ▼          ▼
      SOAP    Diagnosis     Rx      Orders    Follow-up
       │         │         │          │
       └─────────┴─────────┴──────────┘
                         │
                         ▼
                  Care Continuation
                         │
             ┌───────────┴───────────┐
             ▼                       ▼
       In-person visit         Remote follow-up
             │
             ▼
         New Encounter
```

```text
                    ENCOUNTER
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
       Billing       Orders      Notifications
          │            │
       Invoice      Laboratory
          │          Radiology
       Payment          │
          │             ▼
          └──────► Results
                       │
                       ▼
                Patient Record
```

> **Architectural Law**: Appointment determines planned delivery channel. Encounter represents actual care delivered. Telehealth Session represents the remote session layer. Clinical records belong to the encounter and are identical regardless of delivery channel.

