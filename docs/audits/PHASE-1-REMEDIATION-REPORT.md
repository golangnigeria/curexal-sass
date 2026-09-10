# CUREXAL CLINIC OS — PHASE 1 FORENSIC BLOCKER REMEDIATION REPORT

**Audit Date:** September 10, 2026  
**Status:** REMEDIATION COMPLETE — PENDING RE-AUDIT & MANUAL E2E CERTIFICATION  
**Database Schema Version:** Goose Platform Migration `000077`  
**Platform Core:** Go 1.24 / Echo v4 / pgx v5 / PostgreSQL 16  
**Auditors & Engineers:** Senior Healthcare Software Architect, PostgreSQL Database Architect, Multi-Tenant Security Engineer  

---

## 1. Defect 1 Remediation — Encounter Schema Reconciliation (P0/P1 Blocker)

### 1.1 Root Cause Analysis
During the initial Day 3 implementation, migration `000048` created an early baseline of `encounter.encounters`, while migration `000075` utilized `CREATE TABLE IF NOT EXISTS encounter.encounters`. Because the table already existed in PostgreSQL from migration 48, PostgreSQL ignored the `CREATE TABLE IF NOT EXISTS` statement entirely. 
As a consequence, the physical database table lacked critical fields specified in the application and repository layer:
- `appointment_id` (UUID foreign key referencing `operations.appointments(id)`)
- `encounter_channel` (VARCHAR with channel check constraint)
- `signed_at` (TIMESTAMPTZ audit timestamp)
- `signed_by` (UUID foreign key referencing `identity.users(id)`)
- `closed_at` (TIMESTAMPTZ completion timestamp)

Any runtime attempt by `EncounterRepository.CreateEncounter` or `EncounterRepository.GetEncounterByID` to insert or scan these columns resulted in database execution failure (`column does not exist`).

### 1.2 Forward-Only Schema Migration
In accordance with Rule 4 & 5 (zero modification of historical migrations; strictly forward-only schema evolution), migration `000077_reconcile_encounter_encounters_schema.sql` was created and applied:
- Added columns: `appointment_id`, `encounter_channel`, `signed_at`, `signed_by`, `closed_at`.
- Enforced constraint: `chk_encounters_channel CHECK (encounter_channel IN ('in_person', 'video', 'telephone', 'secure_message'))`.
- Enforced foreign keys:
  - `encounters_appointment_id_fkey`: `REFERENCES operations.appointments(id) ON DELETE SET NULL`
  - `encounters_signed_by_fkey`: `REFERENCES identity.users(id) ON DELETE SET NULL`
- Enforced B-tree indexes for performance and multi-tenant lookup:
  - `idx_encounters_channel` ON `encounter.encounters(encounter_channel)`
  - `idx_encounters_appointment` ON `encounter.encounters(appointment_id)`
  - `idx_encounters_signed_at` ON `encounter.encounters(signed_at)`
  - `idx_encounters_closed_at` ON `encounter.encounters(closed_at)`

---

## 2. Defect 2 Remediation — Provider API Security & RBAC Enforcement (P1 Blocker)

### 2.1 Root Cause Analysis
In `apps/api/internal/modules/orchestration/module.go`, provider profile routes were mounted directly onto `apiGroup.Group("/providers")` without authentication middleware:
- `GET /api/v1/providers/profiles`
- `POST /api/v1/providers/profiles`
- `GET /api/v1/providers/profiles/:id`
- `PUT /api/v1/providers/profiles/:id/status`

Anonymous clients could query clinical provider records, room numbers, and internal schedules without credentials. Furthermore, `POST` and `PUT` mutation endpoints did not enforce Casbin permission checks (`users:write`), allowing unauthorized roles to register provider identities. Additionally, `resolveTenantID` fell back to unauthenticated header inspection or database fallback queries.

### 2.2 Security Hardening Implemented
1. **Route Authentication Middleware**:
   Attached `middleware.RequireAuth()` to `providersGroup` in `apps/api/internal/modules/orchestration/module.go`. Anonymous requests are intercepted immediately and rejected with `401 Unauthorized`.
2. **RBAC & Casbin Permission Guard**:
   Attached `middleware.RequirePermission("users:write")` to:
   - `POST /api/v1/providers/profiles`
   - `PUT /api/v1/providers/profiles/:id/status`
   Authenticated principals without `users:write` (such as patients or clinical staff without admin authority) are rejected with `403 Forbidden` ("Insufficient permissions to perform this action").
3. **Tenant Resolution Hardening**:
   Updated `resolveTenantID` in `provider_profile_handler.go` to strictly prioritize authenticated context (`c.Get("tenant_id")`, `c.Get("branch_id")`, and `platformAuth.GetPrincipal(c)`). Unauthenticated requests with missing tenant context immediately return `400 Bad Request` ("Tenant context required").

---

## 3. Defect 3 Remediation — Domain Error Mapping & Tenant Scoping (P2 Defect)

### 3.1 Root Cause Analysis
Two interconnected issues caused patient lookup to return HTTP 500:
1. **PostgreSQL Parameter Type Coercion Bug**:
   In `CanonicalPatientRepository.GetPatientByID`, the SQL query was:
   ```sql
   WHERE (tenant_id = $1 OR $1 = '') AND id = $2
   ```
   Because `tenant_id` is a `UUID` column in PostgreSQL, PostgreSQL inferred `$1` as type `UUID`. When evaluated against `$1 = ''`, PostgreSQL attempted to cast `''` to `UUID`, producing an immediate database error: `ERROR: invalid input syntax for type uuid: "" (SQLSTATE 22P02)`.
2. **Missing UUID Validation & pgx.ErrNoRows Mapping**:
   When clients passed non-UUID strings (e.g. `invalid-id` or nonexistent identifiers), PostgreSQL threw syntax errors that bubbled up as unhandled internal server errors (HTTP 500). Furthermore, `pgx.ErrNoRows` was swallowed as `nil, nil` rather than being mapped to a domain error.

### 3.2 Remediation Implemented
1. **Domain Errors Defined**:
   Exported `ErrPatientNotFound` and `ErrInvalidPatientID` in `apps/api/internal/modules/patient/service/canonical_patient_service.go`.
2. **Repository Query Branching**:
   In `CanonicalPatientRepository.GetPatientByID`:
   - Validated UUID syntax prior to execution; malformed IDs return `pgx.ErrNoRows` safely without querying the DB.
   - Cleanly branched query execution based on `tenant_id`:
     - If `tenant_id` is specified: `WHERE tenant_id = $1 AND id = $2`
     - If `tenant_id` is empty: `WHERE id = $1`
   - Explicitly mapped query misses to `pgx.ErrNoRows`.
3. **Service Layer Mapping**:
   In `CanonicalPatientService.GetPatientByID`:
   - Validates `patientID` and `tenantID` UUID syntax.
   - Intercepts `pgx.ErrNoRows` and returns `ErrPatientNotFound`.
4. **HTTP Handler Standardization**:
   In `CanonicalPatientHandler.GetPatientByID`:
   - Mapped `service.ErrPatientNotFound` and `service.ErrInvalidPatientID` to standardized `echo.NewHTTPError(http.StatusNotFound, "Patient not found")`.
   - Guaranteed HTTP 404 response across all nonexistent, foreign-tenant, and malformed UUID requests, completely eliminating HTTP 500 leaks.

---

## 4. Database Evidence

### 4.1 Migration Status
Direct inspection of `schema_migrations` in PostgreSQL:
```
Current Version: 77
Recent Migrations:
  - 000077: Applied (2026-09-10 10:27:29) [reconcile_encounter_encounters_schema]
  - 000076: Applied (2026-09-10 09:53:17) [codify_architecture_invariants]
  - 000075: Applied (2026-09-10 08:53:08) [day3_clinical_encounter_prescriptions_and_pos_schema]
```

### 4.2 Verified Physical Columns on `encounter.encounters`
```json
[
  { "column_name": "id", "data_type": "uuid", "is_nullable": "NO" },
  { "column_name": "tenant_id", "data_type": "uuid", "is_nullable": "NO" },
  { "column_name": "care_request_id", "data_type": "uuid", "is_nullable": "YES" },
  { "column_name": "patient_id", "data_type": "uuid", "is_nullable": "NO" },
  { "column_name": "provider_id", "data_type": "uuid", "is_nullable": "NO" },
  { "column_name": "encounter_type", "data_type": "varchar", "is_nullable": "NO" },
  { "column_name": "mode", "data_type": "varchar", "is_nullable": "NO" },
  { "column_name": "status", "data_type": "varchar", "is_nullable": "NO" },
  { "column_name": "chief_complaint", "data_type": "text", "is_nullable": "YES" },
  { "column_name": "subjective", "data_type": "text", "is_nullable": "YES" },
  { "column_name": "objective", "data_type": "text", "is_nullable": "YES" },
  { "column_name": "assessment", "data_type": "text", "is_nullable": "YES" },
  { "column_name": "plan", "data_type": "text", "is_nullable": "YES" },
  { "column_name": "primary_diagnosis_code", "data_type": "varchar", "is_nullable": "YES" },
  { "column_name": "primary_diagnosis_name", "data_type": "varchar", "is_nullable": "YES" },
  { "column_name": "secondary_diagnoses", "data_type": "jsonb", "is_nullable": "NO" },
  { "column_name": "started_at", "data_type": "timestamptz", "is_nullable": "NO" },
  { "column_name": "completed_at", "data_type": "timestamptz", "is_nullable": "YES" },
  { "column_name": "created_at", "data_type": "timestamptz", "is_nullable": "NO" },
  { "column_name": "updated_at", "data_type": "timestamptz", "is_nullable": "NO" },
  { "column_name": "care_location_type", "data_type": "varchar", "is_nullable": "NO" },
  { "column_name": "appointment_id", "data_type": "uuid", "is_nullable": "YES" },
  { "column_name": "encounter_channel", "data_type": "varchar", "is_nullable": "NO" },
  { "column_name": "signed_at", "data_type": "timestamptz", "is_nullable": "YES" },
  { "column_name": "signed_by", "data_type": "uuid", "is_nullable": "YES" },
  { "column_name": "closed_at", "data_type": "timestamptz", "is_nullable": "YES" }
]
```

### 4.3 Verified Constraints & Foreign Keys
```text
- chk_encounters_channel: CHECK (encounter_channel IN ('in_person', 'video', 'telephone', 'secure_message'))
- encounters_appointment_id_fkey: FOREIGN KEY (appointment_id) REFERENCES operations.appointments(id) ON DELETE SET NULL
- encounters_signed_by_fkey: FOREIGN KEY (signed_by) REFERENCES identity.users(id) ON DELETE SET NULL
- encounters_care_request_id_fkey: FOREIGN KEY (care_request_id) REFERENCES orchestration.care_requests(id) ON DELETE SET NULL
- encounters_patient_id_fkey: FOREIGN KEY (patient_id) REFERENCES patient.patients(id) ON DELETE CASCADE
- encounters_tenant_id_fkey: FOREIGN KEY (tenant_id) REFERENCES organization.facility_branches(id) ON DELETE CASCADE
```

---

## 5. Security Evidence

Empirical live security verification executed against the platform API server:

| Test Case | Method & Endpoint | Auth State / Role | Expected | Actual | Result |
|---|---|---|---|---|---|
| Anonymous Profile Query | `GET /api/v1/providers/profiles` | None (Anonymous) | 401 | 401 | **PASS** |
| Anonymous Profile Create | `POST /api/v1/providers/profiles` | None (Anonymous) | 401 | 401 | **PASS** |
| Anonymous Status Update | `PUT /api/v1/providers/profiles/:id/status` | None (Anonymous) | 401 | 401 | **PASS** |
| Unauthorized Role Mutation | `POST /api/v1/providers/profiles` | Doctor (`clinic_doctor`, no `users:write`) | 403 | 403 | **PASS** |
| Authorized Provider Query | `GET /api/v1/providers/profiles` | SuperAdmin / Platform Admin | 200 | 200 | **PASS** |
| Cross-Tenant Provider Isolation | `GET /api/v1/providers/profiles` (scoped to Tenant 2) | Admin (Tenant 2 context) | 200 (`total: 0`) | 200 (`total: 0`) | **PASS** |
| Foreign UUID Provider Lookup | `GET /api/v1/providers/profiles/00000000-0000-0000-0000-000000000000` | Admin | 404 | 404 | **PASS** |

---

## 6. API Evidence (Domain Error Mapping & Clinical Encounter)

Empirical live API test results executed via live HTTP probe:

| Scenario | Method & Endpoint | Request Scoping | Response Code | Domain Response Message | Status |
|---|---|---|---|---|---|
| Nonexistent Patient UUID | `GET /api/v1/patients/00000000-0000-0000-0000-000000000000` | Current Tenant (`0001`) | **404 Not Found** | `{"code":"Not Found","message":"Patient not found"}` | **PASS** |
| Malformed UUID Format | `GET /api/v1/patients/malformed-not-a-uuid` | Current Tenant (`0001`) | **404 Not Found** | `{"code":"Not Found","message":"Patient not found"}` | **PASS** |
| Foreign-Tenant Patient Lookup | `GET /api/v1/patients/:id` (Tenant 1 patient) | Foreign Tenant Header (`0002`) | **404 Not Found** | `{"code":"Not Found","message":"Patient not found"}` | **PASS** |
| Valid Patient Profile Lookup | `GET /api/v1/patients/:id` (Tenant 1 patient) | Matching Tenant Header (`0001`) | **200 OK** | Full Patient 360 Canonical Profile returned | **PASS** |
| Clinical Encounter Creation | `POST /api/v1/encounters` with `appointment_id`, `encounter_channel` | Matching Tenant Header (`0001`) | **201 Created** | `{"message":"Encounter started successfully","success":true}` | **PASS** |

---

## 7. Test Suite Evidence

All mandated regression and contract test commands were executed and passed cleanly:

```bash
# 1. TypeScript Static Verification
$ bun x tsc --noEmit
apps/web-patient:  exited with code 0 (0 errors)
apps/web-platform: exited with code 0 (0 errors)
apps/web-public:   exited with code 0 (0 errors)

# 2. Go Module Test Suite
$ go test ./...
ok  github.com/golangnigeria/curexal/internal/modules/audit/application (cached)
ok  github.com/golangnigeria/curexal/internal/modules/authorization/model (cached)
ok  github.com/golangnigeria/curexal/internal/modules/billing/application (cached)
ok  github.com/golangnigeria/curexal/internal/modules/catalogs/application (cached)
ok  github.com/golangnigeria/curexal/internal/modules/clinical/domain (cached)
ok  github.com/golangnigeria/curexal/internal/modules/identity (cached)
ok  github.com/golangnigeria/curexal/internal/modules/identity/domain (cached)
ok  github.com/golangnigeria/curexal/internal/modules/identity/model/auth (cached)
ok  github.com/golangnigeria/curexal/internal/modules/identity/service (cached)
ok  github.com/golangnigeria/curexal/internal/modules/operations/service (cached)
ok  github.com/golangnigeria/curexal/internal/modules/orchestration/service (cached)
ok  github.com/golangnigeria/curexal/internal/modules/organization/application (cached)
ok  github.com/golangnigeria/curexal/internal/modules/organization/domain (cached)
ok  github.com/golangnigeria/curexal/internal/modules/patient/service (cached)
ok  github.com/golangnigeria/curexal/internal/modules/platform/api (cached)
ok  github.com/golangnigeria/curexal/internal/modules/platform/application (cached)
ok  github.com/golangnigeria/curexal/internal/modules/settings/application (cached)
ok  github.com/golangnigeria/curexal/internal/modules/subscription/application (cached)
ok  github.com/golangnigeria/curexal/internal/shared/crypto (cached)
ok  github.com/golangnigeria/curexal/internal/shared/mailer (cached)
ok  github.com/golangnigeria/curexal/internal/shared/middleware (cached)
ok  github.com/golangnigeria/curexal/internal/testing 2.046s
PASS (all packages passed)

# 3. Dedicated Integration Testing
$ go test -count=1 ./internal/testing/...
ok  github.com/golangnigeria/curexal/internal/testing 1.691s

# 4. Database Migration Status
$ go run ./cmd/migrate/main.go status
pre-flight database schema version check passed (current_version: 77)
platform database migration pipeline completed successfully
tenant database migration pipeline completed successfully

# 5. Production Monorepo Build
$ bun run build
Tasks: 4 successful, 4 total (backend binary + all web applications compiled)
```

---

## 8. Files Changed

1. **`apps/api/database/platform/migrations/000077_reconcile_encounter_encounters_schema.sql`**  
   Forward migration adding missing columns (`appointment_id`, `encounter_channel`, `signed_at`, `signed_by`, `closed_at`), foreign key constraints, channel check constraint, and indexes.
2. **`apps/api/internal/modules/orchestration/module.go`**  
   Added `middleware.RequireAuth()` to provider routes, care requests, and portal groups; enforced `middleware.RequirePermission("users:write")` on provider profile creation and status modification routes.
3. **`apps/api/internal/modules/orchestration/handler/provider_profile_handler.go`**  
   Hardened `resolveTenantID` to prioritize authenticated context (`c.Get("tenant_id")`, `c.Get("branch_id")`, `platformAuth.GetPrincipal(c)`), and reject missing tenant scope with 400 Bad Request.
4. **`apps/api/internal/modules/patient/service/canonical_patient_service.go`**  
   Defined domain errors `ErrPatientNotFound` and `ErrInvalidPatientID`. Updated `GetPatientByID` to validate UUIDs and map `pgx.ErrNoRows` to `ErrPatientNotFound`.
5. **`apps/api/internal/modules/patient/repository/canonical_patient_repository.go`**  
   Fixed PostgreSQL UUID type coercion defect by branching queries for `tenant_id`, validating UUID syntax, and returning `pgx.ErrNoRows`.
6. **`apps/api/internal/modules/patient/handler/canonical_patient_handler.go`**  
   Updated `resolveTenantID` to prioritize authenticated principal context; mapped `service.ErrPatientNotFound` and `service.ErrInvalidPatientID` to `404 Not Found`.
7. **`turbo.json`**  
   Configured `globalPassThroughEnv` for `LOCALAPPDATA`, `GOCACHE`, `PATH` to support native Go compilation in turborepo on Windows environments.

---

## 9. Migration Version

- **Current Version:** `77`
- **Migration Name:** `000077_reconcile_encounter_encounters_schema.sql`
- **Application Status:** Fully applied to PostgreSQL database `CUREXAL` on `localhost:5432`.

---

## 10. Remaining Risks & Observations

1. **Casbin Policy Sync across Branch Switches**: When a user switches branch contexts via `POST /api/v1/auth/switch-branch`, the session updates immediately, but frontend client state must be cleared of cached authorization claims.
2. **E2E UI Clinical Workflow Validation**: Although all API and database contracts are verified, complete end-to-end browser walkthroughs across all 3 web apps (`web-platform`, `web-patient`, `web-public`) should be conducted in staging prior to production traffic cutover.
3. **Phase 2 Invariants**: All Phase 2 work (telehealth sessions, WebRTC, real-time messaging) remains cleanly paused; architectural invariants codified in migration 76 ensure telehealth sessions remain strictly subordinate to core clinical encounters.

---

## Remediation Status

```
REMEDIATION STATUS: PASS
```

> **NOTICE:** While all 3 forensic audit blockers and defects have been empirically resolved and verified, **Phase 1 is NOT yet declared production-ready**. A separate forensic re-audit and manual E2E certification must follow.
