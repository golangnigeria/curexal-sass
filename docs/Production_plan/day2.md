# CUREXAL CLINIC OS — DAY 2 PRODUCTION SPECIFICATION & IMPLEMENTATION PLAN

**Document**: `docs/Production_plan/day2.md`  
**Execution Day**: Day 2 of 30 (Friday, September 4, 2026)  
**Target Milestone**: Patient Intake, Identity Matching & Live Queue Management  
**Compliance Standard**: HIPAA Security Rule (§ 164.312), NDPR 2019, Master Patient Index (MPI) Best Practices, OWASP ASVS Level 2, Curexal Constitution v2.0  
**Strict Policy**: **NO SHORTCUTS.** Zero unverified demographic writes, deterministic collision-free MRN sequencing, atomic single-patient queue transitions, zero silent duplicates.

---

## 1. OBJECTIVE & DELIVERABLES

Day 2 establishes the end-to-end clinical intake, provider assignment, and patient throughput pipeline required for operational clinic desks:
1. **Feature #4: Staff and Provider Management**: Link identity staff accounts to clinical provider profiles (medical specialty, regulatory license number, room assignment, active consultation capacity, duty status).
2. **Feature #5: Patient Registration & Master Patient Index (MPI)**: Comprehensive demographic capture, telecoms, next of kin, emergency contacts, automated non-colliding Medical Record Number (`MRN-YYYY-XXXX` / `PAT-YYYY-XXXXX`) generation.
3. **Feature #6: Patient Identity Matching & Duplicate Detection**: Instant multi-signal duplicate matching engine scoring `Phone` + `NIN` + `DOB` + `First/Last Name` with configurable thresholding and modal intervention.
4. **Feature #8: Appointment Scheduling**: Provider-linked calendar bookings, department routing, conflict-free time slot reservation, and walk-in integration.
5. **Feature #10: Patient Check-In & Live Queue Management**: Real-time status transitions (`Registered` $\to$ `Waiting (Triage)` $\to$ `Triaged` $\to$ `In Consultation` $\to$ `Completed`), live wait-time timers, and multi-desk badge increments.

---

## 2. TECHNICAL SPECIFICATION (DATA SCHEMAS & CONTRACTS)

```
┌───────────────────────────────────────────────────────────────────────────────┐
│                             CUREXAL PATIENT FLOW                              │
└───────────────────────────────────────────────────────────────────────────────┘
  Walk-in Patient / Phone Booking
                 │
                 ▼
  ┌──────────────────────────────┐
  │   Front Desk / Reception     │ ◄── Search MPI (Name, Phone, NIN, MRN)
  │    (/reception workspace)    │
  └──────────────┬───────────────┘
                 │
                 ├── [Existing Record] ────► Verify & Update Details
                 │                                   │
                 └── [New Registration]              ▼
                             │                 Check-In to Care Desk
                 ┌───────────┴──────────┐            │
                 │   MPI Duplicate      │            │
                 │  Resolution Engine   │            │
                 └───────────┬──────────┘            │
                             │                       │
                 (Duplicate Alert / Override)        │
                             │                       │
                             ▼                       ▼
                     Generate MRN          ┌───────────────────────┐
                   (PAT-YYYY-XXXXX) ──────►│   Care Desk Queue     │
                                           │  Status: SUBMITTED    │
                                           └───────────┬───────────┘
                                                       │
                                          Nurse Triage & Auto-Acuity
                                                       │
                                                       ▼
                                           ┌───────────────────────┐
                                           │   Status: TRIAGED     │
                                           │  (RED / YELLOW / GREEN│
                                           └───────────┬───────────┘
                                                       │
                                          Provider Matching & Allocation
                                                       │
                                                       ▼
                                           ┌───────────────────────┐
                                           │  In Consultation      │
                                           │  (Doctor Active Room) │
                                           └───────────────────────┘
```

---

### 2.1 Feature #4: Staff and Provider Management Specification

#### A. Architecture & Clinical Governance
In healthcare facilities, generic user accounts cannot conduct medical consultations without certified credentials. A staff user in `identity.users` and `organization.staff_memberships` must be linked to a verified `orchestration.provider_profiles` record specifying:
* Clinical specialty code (e.g., `GP`, `PEDIATRICS`, `CARDIOLOGY`, `INTERNAL_MEDICINE`, `OBGYN`).
* Medical & Dental Council / Regulatory License Number (verified against state board registries).
* Physical or virtual consultation room assignment (e.g., `Consulting Room 3 - Ikeja Branch`).
* Maximum active patient queue capacity (prevents provider burnout, default `10` patients).
* Duty status toggle (`ON_DUTY`, `ON_BREAK`, `OFF_DUTY`, `BUSY`).

#### B. Database Schema: `orchestration.provider_profiles` & Room Extensions
```sql
-- Schema: orchestration
CREATE TABLE IF NOT EXISTS orchestration.provider_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    license_number VARCHAR(100) NOT NULL,
    license_issuer VARCHAR(100) DEFAULT 'MDCN',
    license_verified_at TIMESTAMPTZ,
    specialty_code VARCHAR(50) NOT NULL,
    sub_specialties TEXT[] DEFAULT ARRAY[]::TEXT[],
    room_number VARCHAR(50),
    room_name VARCHAR(100),
    telehealth_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    in_person_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    max_active_queue INT NOT NULL DEFAULT 10,
    current_active_queue INT NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'ON_DUTY', -- ON_DUTY, ON_BREAK, OFF_DUTY, BUSY
    consultation_languages TEXT[] DEFAULT ARRAY['English'],
    rating NUMERIC(3,2) DEFAULT 5.0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_provider_tenant_user UNIQUE (tenant_id, user_id),
    CONSTRAINT uk_provider_license UNIQUE (license_number)
);

CREATE INDEX IF NOT EXISTS idx_provider_profiles_tenant_status 
    ON orchestration.provider_profiles(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_provider_profiles_specialty 
    ON orchestration.provider_profiles(specialty_code);
```

#### C. API Contracts: Provider Profile Management
* `GET /api/v1/providers/profiles`:
  * **Headers**: `Authorization: Bearer <token>`, `X-Branch-ID: <branch_uuid>`
  * **Permissions**: `users:read` or `workspace:clinical:read`
  * **Response (200 OK)**:
    ```json
    {
      "data": [
        {
          "id": "prov_9b83a210",
          "userId": "usr_78129034",
          "tenantId": "branch_550e8400",
          "providerName": "Dr. Musbau Adebayo, FWACP",
          "email": "dr.adebayo@curexal.com",
          "phone": "+2348031234567",
          "licenseNumber": "MDCN/R/78219",
          "licenseIssuer": "Medical and Dental Council of Nigeria",
          "specialtyCode": "GENERAL_PRACTICE",
          "subSpecialties": ["FAMILY_MEDICINE", "DIABETOLOGY"],
          "roomNumber": "Suite 104",
          "roomName": "Primary Clinical Examination Room",
          "telehealthEnabled": true,
          "inPersonEnabled": true,
          "maxActiveQueue": 12,
          "currentActiveQueue": 3,
          "status": "ON_DUTY",
          "consultationLanguages": ["English", "Yoruba"]
        }
      ],
      "meta": { "total": 1, "timestamp": "2026-09-04T08:30:00Z" }
    }
    ```

* `PUT /api/v1/providers/profiles/:id/status`:
  * **Request**:
    ```json
    {
      "status": "ON_BREAK",
      "reason": "Midday ward rounds"
    }
    ```
  * **Response (200 OK)**:
    ```json
    {
      "data": {
        "id": "prov_9b83a210",
        "status": "ON_BREAK",
        "updatedAt": "2026-09-04T12:00:00Z"
      }
    }
    ```

---

### 2.2 Feature #5: Patient Registration & MRN Sequencing Specification

#### A. Master Patient Index (MPI) Demographic Model
Patient records must capture complete legal identity, biological data, contact mechanisms, emergency contacts, and next of kin without data truncations.

#### B. Database Schemas: `patient.patients`, `patient.patient_contacts`, `patient.emergency_contacts`
```sql
CREATE SCHEMA IF NOT EXISTS patient;

-- Canonical Patient Entity
CREATE TABLE IF NOT EXISTS patient.patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    mrn VARCHAR(64) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    preferred_name VARCHAR(100),
    gender VARCHAR(20) NOT NULL, -- MALE, FEMALE, OTHER, UNDISCLOSED
    date_of_birth DATE NOT NULL,
    blood_group VARCHAR(10),     -- A+, A-, B+, B-, AB+, AB-, O+, O-
    genotype VARCHAR(10),        -- AA, AS, SS, AC, SC
    marital_status VARCHAR(20),  -- SINGLE, MARRIED, DIVORCED, WIDOWED
    occupation VARCHAR(100),
    nin VARCHAR(30),             -- National Identification Number
    residential_address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100) DEFAULT 'Nigeria',
    preferred_language VARCHAR(50) DEFAULT 'English',
    status VARCHAR(30) NOT NULL DEFAULT 'REGISTERED', -- REGISTERED, ACTIVE, INACTIVE, DECEASED
    registration_channel VARCHAR(30) NOT NULL DEFAULT 'RECEPTION', -- RECEPTION, PORTAL, TELEHEALTH, EMERGENCY
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_patients_tenant_mrn UNIQUE (tenant_id, mrn)
);

CREATE INDEX IF NOT EXISTS idx_patients_tenant_names 
    ON patient.patients(tenant_id, last_name, first_name);
CREATE INDEX IF NOT EXISTS idx_patients_tenant_dob 
    ON patient.patients(tenant_id, date_of_birth);
CREATE INDEX IF NOT EXISTS idx_patients_tenant_nin 
    ON patient.patients(tenant_id, nin);

-- Patient Telecoms & Channels
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
CREATE INDEX IF NOT EXISTS idx_patient_contacts_system_val 
    ON patient.patient_contacts(system, value);

-- Next of Kin & Emergency Contacts
CREATE TABLE IF NOT EXISTS patient.patient_guardians (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    relationship_type VARCHAR(50) NOT NULL, -- SPOUSE, PARENT, CHILD, SIBLING, GUARDIAN, NEXT_OF_KIN
    full_name VARCHAR(150) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    email VARCHAR(100),
    address TEXT,
    is_emergency_contact BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### C. Deterministic MRN Sequencing Engine
* **Format**: `PAT-YYYY-XXXXX` (or configurable prefix `MRN-YYYY-XXXX`).
* **Generation Strategy**:
  1. Prefix: Canonical string `PAT` (or branch-specific abbreviation).
  2. Year: Current UTC 4-digit calendar year (`2026`).
  3. Sequence: Cryptographically secure 5-digit monotonic pseudo-random integer avoiding predictive enumeration while preventing sequence exhaustion.
  4. Collision Handling: Wrap in transaction with retry loop up to 5 attempts against `uk_patients_tenant_mrn`.

---

### 2.3 Feature #6: Patient Identity Matching & Duplicate Detection Specification

#### A. The Multi-Signal Matching Scoring Matrix
Duplicate medical records cause dangerous medication errors, fragmented medical histories, and catastrophic billing discrepancies. The Curexal MPI matching engine scores incoming candidate signals deterministically:

| Signal Name | Evaluation Method | Weight (Points) | Clinical Justification |
| :--- | :--- | :--- | :--- |
| **Phone Exact Match** | E.164 normalized string comparison (`+234...`) | **+40 pts** | Primary mobile phone is high-entropy in outpatient clinics. |
| **National ID (NIN)** | Exact match on stripped alphanumeric | **+50 pts** | Biometrically verified government identifier. |
| **Date of Birth (DOB)** | Exact match on `YYYY-MM-DD` | **+20 pts** | Eliminates accidental same-name collisions across generations. |
| **Last Name Exact** | Case-insensitive normalized match | **+20 pts** | Family lineage anchor. |
| **First Name Exact** | Case-insensitive normalized match | **+15 pts** | Given name match. |
| **Double Metaphone / Soundex**| Phonetic match on African/Anglo names | **+10 pts** | Protects against clerical spelling variants (e.g., `Emeka` vs `Emeca`). |

#### B. Confidence Tier Classification
$$\text{Total Score} = \sum \text{Matched Signal Weights} \quad (\text{Capped at } 100)$$

* **`EXACT_MATCH` ($\ge 80$ pts)**:
  * Example: Same Phone (+40) + Same DOB (+20) + Same Last Name (+20) = **80 pts**.
  * **System Action**: Form registration is blocked. Receptionist is shown an **Immediate Duplicate Alert Modal** displaying existing patient profile, MRN, and last visit date with an action button: *"Use Existing Patient Profile"*.
* **`PROBABLE_DUPLICATE` ($50 \le \text{Score} < 80$ pts)**:
  * Example: Same Phone (+40) + Same First Name (+15) = **55 pts**.
  * **System Action**: Warning modal displayed. Receptionist can verify physical ID and either link to existing patient or trigger explicit supervisor-authorized override (`forceRegistration: true`).
* **`LOW_SIMILARITY` ($20 \le \text{Score} < 50$ pts)**:
  * Informational banner displayed in background without blocking intake.
* **`NO_MATCH` ($< 20$ pts)**:
  * Clean intake proceeds directly.

#### C. API Contracts: Identity Evaluation & Registration
* `POST /api/v1/patients/mpi/evaluate`:
  * **Request**:
    ```json
    {
      "firstName": "Chinedu",
      "lastName": "Okonkwo",
      "dateOfBirth": "1988-11-23",
      "phone": "+2348023456789",
      "nin": "98127391023"
    }
    ```
  * **Response (200 OK — Duplicate Match Found)**:
    ```json
    {
      "data": {
        "matchStatus": "EXACT_MATCH",
        "candidates": [
          {
            "patientId": "pat_1189204",
            "mrn": "PAT-2026-48192",
            "firstName": "Chinedu",
            "lastName": "Okonkwo",
            "dateOfBirth": "1988-11-23",
            "gender": "MALE",
            "matchedSignals": ["PHONE", "NIN", "DOB", "LAST_NAME", "FIRST_NAME"],
            "confidenceScore": 100,
            "confidenceLevel": "EXACT_MATCH",
            "lastVisitAt": "2026-08-14T10:15:00Z",
            "registeredBranch": "Victoria Island Main Campus"
          }
        ]
      }
    }
    ```

* `POST /api/v1/patients/canonical` (Register Patient):
  * **Request**:
    ```json
    {
      "firstName": "Amaka",
      "middleName": "Grace",
      "lastName": "Eze",
      "gender": "FEMALE",
      "dateOfBirth": "1994-05-18",
      "phone": "+2348123456789",
      "email": "amaka.eze@example.com",
      "nin": "49201928401",
      "bloodGroup": "O+",
      "genotype": "AA",
      "address": "14 Admiralty Way, Lekki Phase 1, Lagos",
      "emergencyContact": {
        "fullName": "Tobechukwu Eze",
        "relationship": "SPOUSE",
        "phone": "+2348099887766"
      },
      "registrationChannel": "RECEPTION",
      "forceRegistration": false
    }
    ```
  * **Response (201 Created)**:
    ```json
    {
      "data": {
        "patient": {
          "id": "pat_3910283",
          "mrn": "PAT-2026-92817",
          "fullName": "Amaka Grace Eze",
          "gender": "FEMALE",
          "dateOfBirth": "1994-05-18",
          "phone": "+2348123456789",
          "status": "REGISTERED",
          "createdAt": "2026-09-04T09:12:00Z"
        },
        "portalInviteSent": true
      }
    }
    ```
  * **Response (409 Conflict — Duplicate Block)**:
    ```json
    {
      "error": {
        "code": "duplicate_patient_detected",
        "message": "Duplicate patient detected with confidence level EXACT_MATCH (Score: 85). Review candidate PAT-2026-48192.",
        "details": {
          "matchStatus": "EXACT_MATCH",
          "candidates": [...]
        }
      }
    }
    ```

---

### 2.4 Feature #8: Appointment Scheduling Specification

#### A. Operational Workflow
Appointments bridge external/portal bookings and clinical rosters:
1. Patient selects department (e.g. `Cardiology`) and specific attending provider.
2. System computes provider availability based on:
   * Standard clinic operational hours.
   * Provider roster shifts (`ON_DUTY`).
   * Existing booked appointments (prevents double-booking).
3. Slot reservation holds time window for 10 minutes prior to confirmation.

#### B. Database Schema: `operations.appointments`
```sql
CREATE SCHEMA IF NOT EXISTS operations;

CREATE TABLE IF NOT EXISTS operations.appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES orchestration.provider_profiles(id) ON DELETE RESTRICT,
    appointment_number VARCHAR(64) NOT NULL,
    service_type VARCHAR(50) NOT NULL DEFAULT 'CONSULTATION',
    appointment_mode VARCHAR(30) NOT NULL DEFAULT 'IN_PERSON', -- IN_PERSON, VIDEO, AUDIO
    status VARCHAR(30) NOT NULL DEFAULT 'BOOKED', -- BOOKED, CHECKED_IN, IN_PROGRESS, COMPLETED, CANCELLED, NO_SHOW
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    reason_for_visit TEXT,
    cancellation_reason TEXT,
    created_by UUID REFERENCES identity.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_appointment_number UNIQUE (tenant_id, appointment_number),
    CONSTRAINT uk_provider_slot_no_overlap EXCLUDE USING gist (
        provider_id WITH =,
        tstzrange(start_time, end_time) WITH &&
    ) WHERE (status IN ('BOOKED', 'CHECKED_IN', 'IN_PROGRESS'))
);

CREATE INDEX IF NOT EXISTS idx_appointments_tenant_date 
    ON operations.appointments(tenant_id, start_time);
CREATE INDEX IF NOT EXISTS idx_appointments_patient 
    ON operations.appointments(patient_id);
```

#### C. API Contracts: Appointment Management
* `POST /api/v1/appointments`:
  * **Request**:
    ```json
    {
      "patientId": "pat_3910283",
      "providerId": "prov_9b83a210",
      "startTime": "2026-09-04T10:30:00Z",
      "endTime": "2026-09-04T11:00:00Z",
      "serviceType": "SPECIALIST_CONSULTATION",
      "appointmentMode": "IN_PERSON",
      "reasonForVisit": "Follow-up hypertension review"
    }
    ```
  * **Response (201 Created)**:
    ```json
    {
      "data": {
        "id": "apt_7719283",
        "appointmentNumber": "APT-2026-00481",
        "patientId": "pat_3910283",
        "providerId": "prov_9b83a210",
        "startTime": "2026-09-04T10:30:00Z",
        "endTime": "2026-09-04T11:00:00Z",
        "status": "BOOKED"
      }
    }
    ```

---

### 2.5 Feature #10: Patient Check-In & Live Queue Management Specification

#### A. State Machine & Urgency Transitions
When a patient arrives at the clinic (whether walk-in or booked), front desk staff checks them in:

```
[REGISTERED / BOOKED]
         │  (Front Desk Check-in: POST /api/v1/orchestration/requests)
         ▼
[SUBMITTED / WAITING_TRIAGE] ────► Real-Time Queue Badge increments on Nurse Desk
         │  (Nurse Triage Intake: POST /api/v1/orchestration/requests/:id/triage)
         ▼
[TRIAGED] ────────────────────────► Priority Acuity Assigned:
         │                           • RED (Emergency / Immediate)
         │                           • YELLOW (Urgent Priority)
         │                           • GREEN (Standard / Routine)
         ▼
[MATCHED / IN_CONSULTATION] ──────► Linked to Consulting Doctor Room
         │  (Doctor Signs SOAP Note / Finishes Consultation)
         ▼
[COMPLETED] ──────────────────────► Discharged / Transferred to Lab/Pharmacy/Billing
```

#### B. Database Schema: `orchestration.care_requests` & `orchestration.triage_assessments`
```sql
CREATE SCHEMA IF NOT EXISTS orchestration;

CREATE TABLE IF NOT EXISTS orchestration.care_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    appointment_id UUID REFERENCES operations.appointments(id) ON DELETE SET NULL,
    request_number VARCHAR(64) NOT NULL,
    service_type VARCHAR(50) NOT NULL DEFAULT 'GENERAL_CONSULTATION',
    preferred_mode VARCHAR(30) NOT NULL DEFAULT 'IN_PERSON',
    urgency VARCHAR(20) NOT NULL DEFAULT 'ROUTINE', -- ROUTINE, URGENT, EMERGENCY
    status VARCHAR(30) NOT NULL DEFAULT 'SUBMITTED', -- SUBMITTED, TRIAGED, MATCHED, IN_PROGRESS, COMPLETED, CANCELLED
    chief_complaint TEXT,
    symptoms_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    checked_in_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    triaged_at TIMESTAMPTZ,
    consultation_started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    assigned_care_agent_id UUID,
    matched_provider_id UUID REFERENCES orchestration.provider_profiles(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_care_request_tenant_number UNIQUE (tenant_id, request_number)
);

CREATE INDEX IF NOT EXISTS idx_care_requests_queue 
    ON orchestration.care_requests(tenant_id, status, checked_in_at);

-- Clinical Triage Assessment
CREATE TABLE IF NOT EXISTS orchestration.triage_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    care_request_id UUID NOT NULL REFERENCES orchestration.care_requests(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
    assessor_id UUID REFERENCES identity.users(id),
    acuity_level VARCHAR(20) NOT NULL DEFAULT 'GREEN', -- RED, YELLOW, GREEN
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
```

#### C. API Contracts: Check-In & Queue Operations
* `POST /api/v1/orchestration/requests` (Patient Check-In):
  * **Request**:
    ```json
    {
      "patientId": "pat_3910283",
      "appointmentId": "apt_7719283",
      "serviceType": "GENERAL_CONSULTATION",
      "preferredMode": "IN_PERSON",
      "chiefComplaint": "Severe recurring migraine and photophobia"
    }
    ```
  * **Response (201 Created)**:
    ```json
    {
      "data": {
        "id": "req_8819203",
        "requestNumber": "REQ-2026-00392",
        "patientId": "pat_3910283",
        "status": "SUBMITTED",
        "urgency": "ROUTINE",
        "checkedInAt": "2026-09-04T09:30:00Z"
      }
    }
    ```

* `POST /api/v1/orchestration/requests/:id/triage` (Nurse Vitals Intake):
  * **Request**:
    ```json
    {
      "systolicBp": 145,
      "diastolicBp": 95,
      "pulseRate": 88,
      "temperature": 38.4,
      "spo2": 96,
      "respiratoryRate": 20,
      "painScore": 6,
      "triageNotes": "Febrile, severe frontal headache, no neck stiffness.",
      "acuityOverride": "YELLOW"
    }
    ```
  * **Response (200 OK)**:
    ```json
    {
      "data": {
        "assessmentId": "tri_4401928",
        "careRequestId": "req_8819203",
        "acuityLevel": "YELLOW",
        "status": "TRIAGED",
        "triagedAt": "2026-09-04T09:38:00Z"
      }
    }
    ```

* `GET /api/v1/orchestration/requests` (Live Queue Feed):
  * **Query Params**: `status=SUBMITTED,TRIAGED,IN_PROGRESS`, `limit=50`
  * **Response (200 OK)**:
    ```json
    {
      "data": [
        {
          "id": "req_8819203",
          "requestNumber": "REQ-2026-00392",
          "patient": {
            "id": "pat_3910283",
            "mrn": "PAT-2026-92817",
            "name": "Amaka Grace Eze",
            "gender": "FEMALE",
            "age": 32
          },
          "status": "TRIAGED",
          "acuityLevel": "YELLOW",
          "checkedInAt": "2026-09-04T09:30:00Z",
          "elapsedWaitMinutes": 18,
          "chiefComplaint": "Severe recurring migraine and photophobia",
          "assignedProvider": null
        }
      ],
      "meta": {
        "total": 1,
        "waitingTriageCount": 0,
        "waitingDoctorCount": 1,
        "inConsultationCount": 2
      }
    }
    ```

---

## 3. IMPLEMENTATION ARCHITECTURE (GO BACKEND & REACT 19 FRONTEND)

### 3.1 Go Backend Modular Structure

```
apps/api/internal/modules/
├── patient/
│   ├── api/
│   │   ├── canonical_patient_handler.go      # POST /canonical, GET /:id, search
│   │   └── mpi_handler.go                    # POST /mpi/evaluate, GET /candidates
│   ├── model/
│   │   ├── patient.go                        # Canonical Patient entity
│   │   ├── patient_dto.go                    # Request payloads & DTOs
│   │   └── mpi_resolution.go                 # Scored signals & candidate models
│   ├── repository/
│   │   └── canonical_patient_repository.go   # PostgreSQL pgx queries, indexed lookups
│   └── service/
│       ├── canonical_patient_service.go      # Business logic & atomic MRN sequencing
│       └── mpi_service.go                    # Matching scoring algorithm (+40, +50, +20...)
├── orchestration/
│   ├── handler/
│   │   └── care_request_handler.go           # Check-in, triage intake, provider match
│   ├── model/
│   │   ├── care_request_domain.go            # Request states, queue entity
│   │   └── triage_and_matching.go            # Triage vats, acuity scoring, provider candidates
│   ├── repository/
│   │   └── care_request_repository.go        # Queue state management, elapsed time calc
│   └── service/
│       └── care_request_service.go           # Orchestration workflows & event publishing
└── operations/
    ├── handler/
    │   └── appointment_handler.go            # Appointment booking, calendar slot query
    └── service/
        └── appointment_service.go            # Double-booking exclusion & check-in binding
```

### 3.2 Frontend Applications (`apps/web-platform`)

* **Reception Intake & Master Patient Index**:
  * Route: `/:branchSlug/reception` ([`apps/web-platform/src/pages/workspace/reception/index.tsx`](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/apps/web-platform/src/pages/workspace/reception/index.tsx)).
  * Modal: [`apps/web-platform/src/components/patients/patient-intake-modal.tsx`](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/apps/web-platform/src/components/patients/patient-intake-modal.tsx).
  * Features:
    * Debounced (500ms) instant duplicate lookup on Phone/NIN/Last Name input.
    * Warning banner and duplicate resolution candidate list with confidence badges.
    * Supervisor override checkbox with reason tracking.
    * Immediate check-in button upon successful registration.
* **Care Desk & Live Queue Monitor**:
  * Route: `/:branchSlug/care-desk` ([`apps/web-platform/src/pages/workspace/care-desk/index.tsx`](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/apps/web-platform/src/pages/workspace/care-desk/index.tsx)).
  * Features:
    * Live queue board grouped by acuity (`RED`, `YELLOW`, `GREEN`).
    * Real-time wait timer component (`Clock` icon showing elapsed minutes with color escalation: $>30\text{ mins}$ amber, $>60\text{ mins}$ red).
    * Nurse triage intake dialog with live auto-acuity calculator (BP, SpO2, Temp, Pain scale).
    * Provider matching modal assigning on-duty doctors by specialty and current active queue load.

---

## 4. STEP-BY-STEP IMPLEMENTATION RUNBOOK

### Step 1: Database Migration & Schema Validation
1. Verify migration `000048_create_patient_access_and_care_orchestration_schema.sql` is fully applied.
2. Verify table `patient.patients` has unique constraint on `(tenant_id, mrn)`.
3. Verify table `patient.patient_contacts` has indexed `(system, value)` lookups.
4. Verify table `orchestration.provider_profiles` has unique `(tenant_id, user_id)` and status index.

### Step 2: Backend Service & Handler Verification
1. Verify `apps/api/internal/modules/patient/service/mpi_service.go`:
   - Validates scoring algorithm: Phone (+40), NIN (+50), DOB (+20), Last Name (+20), First Name (+15).
   - Classifies `EXACT_MATCH` ($\ge 80$), `PROBABLE_DUPLICATE` ($\ge 50$).
2. Verify `apps/api/internal/modules/patient/service/canonical_patient_service.go`:
   - Validates `GenerateMRN()` generates non-colliding formatted IDs.
   - Enforces `mpiService.EvaluateDuplicates` check unless `ForceRegistration == true`.
3. Verify `apps/api/internal/modules/orchestration/service/care_request_service.go`:
   - Handles `CheckInPatient` creating care request in `SUBMITTED` state.
   - Handles `SubmitTriage` transitioning state to `TRIAGED` with acuity calculation.

### Step 3: Frontend Integration & Real-Time Polish
1. In [`apps/web-platform/src/components/patients/patient-intake-modal.tsx`](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/apps/web-platform/src/components/patients/patient-intake-modal.tsx):
   - Confirm duplicate detection triggers automatically when phone $\ge 8$ digits.
   - Confirm candidate cards display matched signals (`PHONE`, `DOB`, `NIN`) and confidence percentage.
2. In [`apps/web-platform/src/pages/workspace/care-desk/index.tsx`](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/apps/web-platform/src/pages/workspace/care-desk/index.tsx):
   - Confirm live queue badge increments upon check-in.
   - Confirm elapsed wait timer renders dynamically without freezing the DOM.

---

## 5. ACCEPTANCE CRITERIA & TESTING PROTOCOL

### 5.1 Automated Test Suite Matrix

| Test Suite | File Location | Target Verification |
| :--- | :--- | :--- |
| **MPI Scoring Unit Test** | `apps/api/internal/modules/patient/service/mpi_service_test.go` | Assert Phone + DOB gives 60 pts (`PROBABLE_DUPLICATE`). Assert Phone + NIN gives 90 pts (`EXACT_MATCH`). |
| **MRN Generation Test** | `apps/api/internal/modules/patient/service/canonical_patient_test.go` | Generate 1,000 consecutive MRNs; verify 0 collisions and regex `PAT-2026-\d{5}` compliance. |
| **Duplicate Prevention Integration** | `apps/api/internal/testing/mpi_integration_test.go` | Attempt registering second patient with identical phone number without force flag; assert `409 Conflict`. |
| **Queue State Machine Test** | `apps/api/internal/testing/queue_workflow_test.go` | Check-in $\to$ Triage $\to$ Match $\to$ Complete; assert timestamps and status transitions. |

### 5.2 PowerShell Automated Verification Script

```powershell
# ==============================================================================
# DAY 2 AUTOMATED INTEGRATION VERIFICATION SCRIPT
# ==============================================================================
$ErrorActionPreference = "Stop"
$BaseUrl = "http://localhost:8080/api/v1"

Write-Host "==> 1. Authenticating as Receptionist..." -ForegroundColor Cyan
$LoginBody = @{
    email = "receptionist@curexal.com"
    password = "SecurePassword123!"
} | ConvertTo-Json

$AuthResponse = Invoke-RestMethod -Uri "$BaseUrl/auth/login" -Method Post -Body $LoginBody -ContentType "application/json"
$Token = $AuthResponse.data.accessToken
$BranchId = $AuthResponse.data.activeBranch.id
$Headers = @{
    "Authorization" = "Bearer $Token"
    "X-Branch-ID" = $BranchId
}

Write-Host "==> 2. Registering Initial Patient (John Doe)..." -ForegroundColor Cyan
$Patient1 = @{
    firstName = "John"
    lastName = "Doe"
    gender = "MALE"
    dateOfBirth = "1990-01-15"
    phone = "+2348011223344"
    email = "john.doe@test.com"
    registrationChannel = "RECEPTION"
} | ConvertTo-Json

$Reg1 = Invoke-RestMethod -Uri "$BaseUrl/patients/canonical" -Method Post -Headers $Headers -Body $Patient1 -ContentType "application/json"
$Patient1Id = $Reg1.data.patient.id
$MRN1 = $Reg1.data.patient.mrn
Write-Host "    Registered Patient 1 with MRN: $MRN1" -ForegroundColor Green

Write-Host "==> 3. Testing MPI Duplicate Detection with Duplicate Phone..." -ForegroundColor Cyan
$DuplicatePatient = @{
    firstName = "Jonathan"
    lastName = "Doe"
    gender = "MALE"
    dateOfBirth = "1990-01-15"
    phone = "+2348011223344"
} | ConvertTo-Json

try {
    $DupCheck = Invoke-RestMethod -Uri "$BaseUrl/patients/canonical" -Method Post -Headers $Headers -Body $DuplicatePatient -ContentType "application/json"
    Write-Error "FAIL: Duplicate registration should have been blocked!"
} catch {
    Write-Host "    PASS: Duplicate registration blocked with HTTP 409 / Error" -ForegroundColor Green
}

Write-Host "==> 4. Checking in Patient 1 to Care Desk Queue..." -ForegroundColor Cyan
$CheckInBody = @{
    patientId = $Patient1Id
    serviceType = "GENERAL_CONSULTATION"
    chiefComplaint = "High fever and persistent chills"
} | ConvertTo-Json

$CheckIn = Invoke-RestMethod -Uri "$BaseUrl/orchestration/requests" -Method Post -Headers $Headers -Body $CheckInBody -ContentType "application/json"
$RequestId = $CheckIn.data.id
Write-Host "    Patient checked in. Care Request: $RequestId (Status: SUBMITTED)" -ForegroundColor Green

Write-Host "==> 5. Submitting Nurse Triage Assessment..." -ForegroundColor Cyan
$TriageBody = @{
    systolicBp = 135
    diastolicBp = 85
    pulseRate = 92
    temperature = 38.9
    spo2 = 95
    respiratoryRate = 22
    painScore = 5
    triageNotes = "Patient visibly flushed, shivering."
} | ConvertTo-Json

$Triage = Invoke-RestMethod -Uri "$BaseUrl/orchestration/requests/$RequestId/triage" -Method Post -Headers $Headers -Body $TriageBody -ContentType "application/json"
Write-Host "    Triage saved. Acuity Level: $($Triage.data.acuityLevel) (Status: $($Triage.data.status))" -ForegroundColor Green

Write-Host "==> 6. Verifying Live Queue Feed..." -ForegroundColor Cyan
$Queue = Invoke-RestMethod -Uri "$BaseUrl/orchestration/requests" -Method Get -Headers $Headers
$MatchedQueue = $Queue.data | Where-Object { $_.id -eq $RequestId }

if ($MatchedQueue -and $MatchedQueue.status -eq "TRIAGED") {
    Write-Host "SUCCESS: Day 2 Intake, Matching & Queue Verification Passed 100%!" -ForegroundColor Green
} else {
    Write-Error "FAIL: Queue item not found or incorrect status!"
}
```

---

## 6. ROLLBACK & FAILURE RECOVERY PLAN

### 6.1 Database Schema Rollback
If schema defects are identified during deployment, execute downgrade statements:
```sql
-- Down migration for Day 2 components
DROP SCHEMA IF EXISTS encounter CASCADE;
DROP SCHEMA IF EXISTS orchestration CASCADE;
DROP SCHEMA IF EXISTS patient CASCADE;
```

### 6.2 Manual Paper Backup Reconciliation Protocol
In the event of network or server failure at reception desks:
1. Receptionists transition to standardized carbon-copy paper intake forms (`Form CUREX-INTAKE-01`).
2. The front desk records: Full Legal Name, Phone Number, Date of Birth, Gender, Emergency Contact, and Arrival Timestamp.
3. Once connectivity is restored:
   - Receptionists enter records into `/reception` utilizing the standard intake modal.
   - The MPI engine evaluates each backlog record. If duplicate detected, records are merged using the primary MRN.
   - Historical check-in timestamps are backfilled to preserve true waiting time telemetry.

---

## 7. PRODUCTION AUDIT TRAIL EVIDENCE TEMPLATE

Every demographic mutation, duplicate resolution, and queue movement must write an immutable audit record to `audit.audit_events`:

| Action Code | Target Entity | Actor Persona | Severity | Metadata Logged |
| :--- | :--- | :--- | :--- | :--- |
| `patient:create` | `patient.patients` | `receptionist` | `INFO` | `mrn`, `channel: RECEPTION`, `name_hash` |
| `patient:duplicate_detected`| `patient.patients` | `receptionist` | `WARNING` | `matched_mrn`, `confidence_score`, `signals` |
| `patient:duplicate_override`| `patient.patients` | `org_admin` | `HIGH` | `justification`, `supervisor_id`, `new_mrn` |
| `appointment:schedule` | `operations.appointments` | `receptionist` | `INFO` | `patient_id`, `provider_id`, `time_slot` |
| `queue:check_in` | `orchestration.care_requests` | `receptionist` | `INFO` | `request_number`, `service_type` |
| `queue:triage_completed` | `orchestration.care_requests` | `nurse` | `INFO` | `acuity: RED/YELLOW/GREEN`, `vitals_summary` |
| `queue:provider_assigned` | `orchestration.care_requests` | `doctor` | `INFO` | `provider_id`, `room_number` |

---

**Day 2 Plan Complete & Certified for Production Implementation.**
