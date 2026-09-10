# CUREXAL PHASE 1 — PRODUCTION FORENSIC AUDIT MATRIX

**Audit Date**: September 10, 2026  
**Environment**: Local Production Staging (`PostgreSQL 16`, `Go 1.25 + Echo`, `Bun / React 19`)  
**Auditor Personas**: Senior Healthcare Architect, Senior Go Backend Engineer, PostgreSQL Architect, Security & RBAC Auditor  
**Target Milestone**: Phase 1 Live Launch (Days 1 – 3)

---

## 1. Feature Audit Matrix

| Feature | Backend | DB | API | Frontend | Auth | Tenant Isolation | Audit | E2E | Status |
|---------|---------|----|-----|----------|------|------------------|-------|-----|--------|
| **1. Authentication (Argon2id)** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **2. Login & Credential Verification** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **3. Session Management & Tokens** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **4. Logout & Token Invalidation** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **5. Organization Management** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **6. Clinic & Branch Facility Governance** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **7. Role-Based Access Control (RBAC)** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **8. Scope-Aware Permissions** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **9. Audit Logging Engine & Immutability** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **10. Multi-Tenant Isolation** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **11. Staff Management & Memberships** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **12. Provider Profiles & Management** | PASS | PASS | FAIL | PASS | FAIL | PASS | PASS | PARTIAL | **PARTIAL** |
| **13. Patient Demographic Registration** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **14. Monotonic MRN Generation** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **15. Duplicate Detection & MPI Engine** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **16. Patient Lookup & Search** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **17. Appointment Scheduling** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **18. Provider Availability Constraints** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **19. Patient Check-in & Live Queue** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **20. Encounter Creation** | PASS | FAIL | FAIL | PASS | PASS | PASS | PASS | FAIL | **FAIL** |
| **21. Encounter Lifecycle Transitions** | PASS | FAIL | FAIL | PASS | PASS | PASS | PASS | FAIL | **FAIL** |
| **22. SOAP Notes Canvas** | PASS | FAIL | FAIL | PASS | PASS | PASS | PASS | FAIL | **PARTIAL** |
| **23. Triage Vitals & Observations** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **24. ICD-10 Diagnoses Management** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **25. Medication & e-Prescribing** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **26. Clinical Record Finalization & Sign-off** | PASS | FAIL | FAIL | PASS | PASS | PASS | PASS | FAIL | **FAIL** |
| **27. POS Billing Engine** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **28. Invoice Generation & Line Items** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **29. Payment Recording & Split Tenders** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **30. Patient Visit Summary & POS Receipt** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **31. Clinical & Financial Audit Trail** | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | **PASS** |
| **32. DB Schema Migration Integrity** | PASS | FAIL | N/A | N/A | N/A | PASS | PASS | N/A | **PARTIAL** |
| **33. API Gateway Authentication Enforcer** | PASS | PASS | FAIL | N/A | FAIL | PASS | PASS | N/A | **PARTIAL** |

---

## 2. Forensic Defect Reports (PARTIAL / FAIL Items)

### Defect 1: Missing Columns on `encounter.encounters` (Schema Delta Bug)
- **Status**: **FAIL**
- **Affected Features**: Feature #20 (Encounter Creation), Feature #21 (Encounter Lifecycle), Feature #22 (SOAP Notes), Feature #26 (Clinical Finalization), Feature #32 (DB Migration Integrity).
- **File**: `apps/api/database/platform/migrations/000075_day3_clinical_encounter_prescriptions_and_pos_schema.sql` (Line 10).
- **Function / Component**: `EncounterRepository.CreateEncounter`, `EncounterRepository.GetEncounterByID` in `apps/api/internal/modules/encounter/repository/encounter_repository.go`.
- **Endpoint**: `POST /api/v1/encounters`, `GET /api/v1/encounters/:id`, `PUT /api/v1/encounters/:id/soap`.
- **Database Table**: `encounter.encounters`.
- **Migration**: `000075_day3_clinical_encounter_prescriptions_and_pos_schema.sql`.
- **Exact Problem**: 
  The table `encounter.encounters` was originally created in migration `000048_create_patient_access_and_care_orchestration_schema.sql`. Migration `000075` declared:
  ```sql
  CREATE TABLE IF NOT EXISTS encounter.encounters (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      tenant_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
      patient_id UUID NOT NULL REFERENCES patient.patients(id) ON DELETE CASCADE,
      appointment_id UUID REFERENCES operations.appointments(id) ON DELETE SET NULL,
      care_request_id UUID REFERENCES orchestration.care_requests(id) ON DELETE SET NULL,
      provider_id UUID REFERENCES orchestration.provider_profiles(id) ON DELETE SET NULL,
      encounter_channel VARCHAR(30) NOT NULL DEFAULT 'in_person',
      status VARCHAR(30) NOT NULL DEFAULT 'in_progress',
      ...
      signed_at TIMESTAMPTZ,
      signed_by UUID REFERENCES identity.users(id),
      closed_at TIMESTAMPTZ,
      ...
  );
  ```
  Because the table already existed, PostgreSQL skipped the statement due to `IF NOT EXISTS`. As a result, the following columns were **never added to PostgreSQL**:
  1. `appointment_id` (UUID)
  2. `encounter_channel` (VARCHAR(30))
  3. `signed_at` (TIMESTAMPTZ)
  4. `signed_by` (UUID)
  5. `closed_at` (TIMESTAMPTZ)
- **Reproduction Steps**:
  1. Send `POST /api/v1/encounters` with a valid clinical payload:
     ```json
     {
       "patientId": "aa7bc2b1-94b6-4e8e-bc00-6bd1ccfd2122",
       "providerId": "00000000-0000-0000-0000-000000000001",
       "encounterChannel": "in_person",
       "chiefComplaint": "Severe migraine"
     }
     ```
  2. The server crashes the query and responds:
     ```json
     {
       "code": "Bad Request",
       "message": "failed to start encounter: ERROR: column \"appointment_id\" of relation \"encounters\" does not exist (SQLSTATE 42703)"
     }
     ```
- **Severity**: **CRITICAL (Blocker)** — Prevents clinical encounters from being created or finalized in production.
- **Recommended Fix**:
  Create an idempotent forward migration (e.g. `000077_add_missing_encounter_columns.sql`):
  ```sql
  ALTER TABLE encounter.encounters ADD COLUMN IF NOT EXISTS appointment_id UUID REFERENCES operations.appointments(id) ON DELETE SET NULL;
  ALTER TABLE encounter.encounters ADD COLUMN IF NOT EXISTS encounter_channel VARCHAR(30) NOT NULL DEFAULT 'in_person';
  ALTER TABLE encounter.encounters ADD COLUMN IF NOT EXISTS signed_at TIMESTAMPTZ;
  ALTER TABLE encounter.encounters ADD COLUMN IF NOT EXISTS signed_by UUID REFERENCES identity.users(id);
  ALTER TABLE encounter.encounters ADD COLUMN IF NOT EXISTS closed_at TIMESTAMPTZ;
  ```

---

### Defect 2: Unauthenticated Public Exposure of Provider Profiles
- **Status**: **PARTIAL**
- **Affected Features**: Feature #12 (Provider Profiles & Management), Feature #33 (API Gateway Security).
- **File**: `apps/api/internal/modules/orchestration/module.go` (Lines 40–47).
- **Function / Component**: `Module.RegisterRoutes(apiGroup *echo.Group)`
- **Endpoint**: `GET /api/v1/providers/profiles`, `GET /api/v1/providers/profiles/:id`, `POST /api/v1/providers/profiles`, `PUT /api/v1/providers/profiles/:id/status`.
- **Database Table**: `orchestration.provider_profiles`.
- **Migration**: `000074_day2_provider_profiles_appointments_and_patient_schema.sql`.
- **Exact Problem**:
  `providersGroup` was attached directly to `apiGroup` without requiring `authMiddleware` or permission checks:
  ```go
  // Day 2 Feature #4: Provider Profile Management
  providersGroup := apiGroup.Group("/providers/profiles")
  providersGroup.GET("", m.ProviderHandler.ListProviders)
  providersGroup.GET("/:id", m.ProviderHandler.GetProviderByID)
  providersGroup.POST("", m.ProviderHandler.CreateProvider)
  providersGroup.PUT("/:id/status", m.ProviderHandler.UpdateStatus)
  ```
  An anonymous unauthenticated client without a JWT token can execute `GET /api/v1/providers/profiles` and receive HTTP 200 with full provider license numbers, duty statuses, room assignments, and internal IDs.
- **Reproduction Steps**:
  1. Execute:
     ```powershell
     curl -i http://localhost:8080/api/v1/providers/profiles
     ```
  2. Notice that the server returns HTTP 200 with full provider directory data without requiring authentication headers or cookies.
- **Severity**: **HIGH (Security & Privacy Vulnerability)**.
- **Recommended Fix**:
  Attach auth middleware to `providersGroup` or mount it inside the authenticated route group:
  ```go
  providersGroup := apiGroup.Group("/providers/profiles", authMiddleware)
  ```

---

### Defect 3: Cross-Tenant 500 Panic on Foreign Tenant Lookups
- **Status**: **PARTIAL**
- **Affected Features**: Feature #10 (Multi-Tenant Isolation error handling).
- **File**: `apps/api/internal/modules/patient/service/canonical_patient_service.go` / `canonical_patient_handler.go`.
- **Function / Component**: `GetPatientByID`.
- **Endpoint**: `GET /api/v1/patients/canonical/:id`.
- **Database Table**: `patient.patients`.
- **Exact Problem**:
  When a user authenticated in Tenant A requests a patient with `X-Tenant-ID: <Tenant B>`, the service queries the database with Tenant B. When no row is returned (`pgx.ErrNoRows`), instead of returning a clean `404 Not Found` or `403 Forbidden`, the handler formats an internal error that propagates as `500 Internal Server Error`.
- **Reproduction Steps**:
  1. Login as Nurse in Tenant A.
  2. Send `GET /api/v1/patients/canonical/<valid-patient-id>` with header `X-Tenant-ID: <foreign-tenant-uuid>`.
  3. The response is `HTTP 500 Internal Server Error` instead of `HTTP 404 Patient Not Found`.
- **Severity**: **MEDIUM (Error Handling & Leaking Internal Exception)**.
- **Recommended Fix**:
  Map `pgx.ErrNoRows` in `CanonicalPatientService.GetPatientByID` to domain error `domain.ErrPatientNotFound`, and map that to `404 Not Found` in the HTTP handler.
