# CUREXAL CLINIC OS — PHASE 1 MANUAL TEST PLAN & PRODUCTION CERTIFICATION CHECKLIST

**Document Identifier**: `docs/testing/PHASE-1-MANUAL-TEST-PLAN.md`  
**Target Milestone**: Phase 1 Release Certification (Days 1–3 Baseline)  
**System Under Test**: Curexal Clinic OS (Unified Clinic Operating Platform)  
**Execution Type**: Human Manual Execution & Security Verification  
**Database Migration Baseline**: Goose Platform Migration `000077`  
**Author / Role**: Senior QA & Release Engineering Agent  
**Classification**: Controlled Medical Device & Clinical Platform Release Documentation  

---

## TABLE OF CONTENTS

1. [SECTION 1 — Test Execution Setup](#section-1--test-execution-setup)
2. [SECTION 2 — Test Accounts / Personas](#section-2--test-accounts--personas)
3. [SECTION 3 — Platform / Organization / Branch Setup](#section-3--platform--organization--branch-setup)
4. [SECTION 4 — Authentication](#section-4--authentication)
5. [SECTION 5 — RBAC / Casbin / Authorization](#section-5--rbac--casbin--authorization)
6. [SECTION 6 — Multi-Tenant Isolation](#section-6--multi-tenant-isolation)
7. [SECTION 7 — Patient Registration / MPI](#section-7--patient-registration--mpi)
8. [SECTION 8 — Appointments](#section-8--appointments)
9. [SECTION 9 — Patient Check-In / Queue](#section-9--patient-check-in--queue)
10. [SECTION 10 — Triage](#section-10--triage)
11. [SECTION 11 — Encounter / Clinical Core](#section-11--encounter--clinical-core)
12. [SECTION 12 — SOAP / Clinical Documentation](#section-12--soap--clinical-documentation)
13. [SECTION 13 — Diagnosis](#section-13--diagnosis)
14. [SECTION 14 — Prescriptions](#section-14--prescriptions)
15. [SECTION 15 — Billing / POS](#section-15--billing--pos)
16. [SECTION 16 — Patient Portal](#section-16--patient-portal)
17. [SECTION 17 — Audit Logging](#section-17--audit-logging)
18. [SECTION 18 — Frontend / UI Testing](#section-18--frontend--ui-testing)
19. [SECTION 19 — API Testing](#section-19--api-testing)
20. [SECTION 20 — Database Verification](#section-20--database-verification)
21. [SECTION 21 — Security / Attack Tests](#section-21--security--attack-tests)
22. [SECTION 22 — Golden Path](#section-22--golden-path)
23. [SECTION 23 — Negative Path](#section-23--negative-path)
24. [SECTION 24 — Phase 1 Feature Matrix](#section-24--phase-1-feature-matrix)
25. [SECTION 25 — Evidence Requirements](#section-25--evidence-requirements)
26. [SECTION 26 — Defect Reporting](#section-26--defect-reporting)
27. [SECTION 27 — Final Certification Gate](#section-27--final-certification-gate)

---

## SECTION 1 — TEST EXECUTION SETUP

Follow these instructions in exact sequence to instantiate a clean, deterministic local staging environment prior to initiating manual test execution.

### 1.1 Infrastructure & Service Topology

| Service | Target Port | Protocol / Technology | Runtime Process / Host |
|---|---|---|---|
| **PostgreSQL 16** | `5432` | `postgresql://postgres:postgres@localhost:5432/CUREXAL` | Docker (`curexal-postgres`) |
| **Redis 7** | `6379` | `redis://localhost:6379` | Docker (`curexal-redis`) |
| **Native API Server** | `8080` | HTTP / Go 1.25 + Echo v4 | `apps/api/cmd/CUREXAL` |
| **Web Platform (Staff/Console)** | `5002` | HTTP / React 19 + Vite (Bun) | `apps/web-platform` |
| **Web Patient Portal** | `5003` | HTTP / React 19 + Vite (Bun) | `apps/web-patient` |
| **Web Public (Marketing/Demo)** | `5001` | HTTP / React 19 + Vite (Bun) | `apps/web-public` |

### 1.2 Step-by-Step Environment Preparation

#### Step 1: Start Containerized Database & Cache
Execute from the project workspace root:
```powershell
# Using Taskfile
task db:up

# Or directly via Docker Compose
docker compose up -d
```
*Verification*: Run `docker compose ps` and ensure `curexal-postgres` and `curexal-redis` report `Up (healthy)`.

#### Step 2: Database Migration Status & Forward Application
Run the database migration tool to ensure the PostgreSQL schema is at Goose version **77**:
```powershell
# Check current migration version
cd apps/api
go run ./cmd/migrate/main.go status

# If unapplied migrations exist, apply up migrations:
go run ./cmd/migrate/main.go
```
*Pass Criteria*: The console reports: `pre-flight database schema version check passed (current_version: 77)`.

#### Step 3: Seed Platform Admin & Baseline Organization Data
Seed the standard platform roles, default organizations, and branch accounts:
```powershell
cd apps/api
go run ./cmd/seed_admin
```
*Pass Criteria*: Console outputs confirmation that platform admin (`admin@curexal.com`) and organization owners (`owner@curexal.space`, `owner@everight.com`) have been created or updated.

#### Step 4: Verify Backend Health Endpoint
Send an unauthenticated probe to the health endpoint:
```powershell
curl -i http://localhost:8080/health
```
*Expected Response*:
```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "checks": {
    "database": {
      "status": "healthy"
    }
  },
  "environment": "development",
  "status": "healthy",
  "timestamp": "..."
}
```

#### Step 5: Launch Native API Server
In a dedicated terminal:
```powershell
cd apps/api
go run ./cmd/CUREXAL
```
*Expected Log*: `⇨ http server started on [::]:8080`

#### Step 6: Launch Frontend Applications
In a second terminal:
```powershell
# Run staff platform workspace (Port 5002)
cd apps/web-platform
bun run dev

# In another terminal: Run patient portal (Port 5003)
cd apps/web-patient
bun run dev
```

#### Step 7: Verify Web Access in Browser
- Open Chrome or Firefox at `http://localhost:5002/login`. Verify the Curexal Platform Sign-in page renders with branding and input fields.
- Open `http://localhost:5003/login`. Verify the Curexal Patient Portal Sign-in page renders.

#### Step 8: Test Reset Procedure (If fresh slate required)
To completely wipe transactional test data and restore seed baselines:
```powershell
cd apps/api
go run ./cmd/clean_db
go run ./cmd/seed_admin
```

---

## SECTION 2 — TEST ACCOUNTS / PERSONAS

The test plan utilizes 12 pre-configured personas representing distinct privilege tiers, administrative scopes, and medical specialties.

| Persona ID | Account Email | Initial Password | Scope / Level | Role Code | Allowed Actions | Strictly Prohibited Actions |
|---|---|---|---|---|---|---|
| **PER-01** | `admin@curexal.com` | `password` | Platform | `super_admin` | Full platform diagnostics, tenant provisioning, system config, global audit view. | Direct clinical entry into tenant charts without active clinical role. |
| **PER-02** | `support@curexal.internal` | `password` | Platform | `super_support_agent` | View platform tenant directory, diagnostics, impersonate tenant for support. | Mutate clinical records, dispense prescriptions, settle POS. |
| **PER-03** | `owner@curexal.space` | `password` | Organization | `owner` (Curexal Clinic) | Full executive management of Curexal Clinic: create branches, invite staff, view all financial/clinical data. | Access Tenant B (Everight Hospital) data or platform-only configs. |
| **PER-04** | `admin@curexal.space` | `password` | Organization | `org_admin` | Manage organization branches, invite staff, assign departments, configure catalogs. | Clinical prescription signing, deleting audit records. |
| **PER-05** | `branchadmin@curexal.space` | `password` | Branch | `branch_admin` (Main Branch) | Manage branch operations, patient intake, appointment scheduling, local staff roster. | Access other branches without explicit assignment; sign doctor prescriptions. |
| **PER-06** | `dr.emeka@curexal.com` | `password` | Branch | `doctor` / `clinic_doctor` | Open encounters, record SOAP notes, add ICD-10 diagnoses, sign e-prescriptions. | Financial cashier POS settlement, platform administration. |
| **PER-07** | `dr.sarah@curexal.com` | `password` | Multi-Branch | `doctor` (Main & Victoria Island) | Perform consultations across multiple assigned branches. | Access patient records of branches not assigned to her profile. |
| **PER-08** | `nurse.chioma@curexal.com` | `password` | Branch | `nurse` / `clinic_nurse` | Access triage queue, record vital signs, document observations, assign acuity. | Sign physician encounters, issue prescription medications, settle POS cash. |
| **PER-09** | `receptionist@curexal.com` | `password` | Branch | `receptionist` | Register patient demographics, search MPI, book appointments, check-in patients. | Enter clinical SOAP notes, triage acuity, sign prescriptions. |
| **PER-10** | `cashier@curexal.com` | `password` | Branch | `cashier` | Query pending patient invoices, record cash/card/transfer payments, print receipts. | View protected clinical SOAP notes, modify provider diagnoses, sign charts. |
| **PER-11** | `owner@everight.com` | `password` | Organization | `owner` (Everight Hospital) | Manage Tenant B (Everight Hospital) organization and branches. | Access any patient, appointment, encounter, or invoice belonging to Curexal Clinic. |
| **PER-12** | `patient.adanna@curexal.com` | *(OTP/PIN)* | Patient | `patient` | View own profile, appointments, care journey, and finalized prescriptions. | Access other patients' records, view staff internal audit logs, modify bills. |

---

## SECTION 3 — PLATFORM / ORGANIZATION / BRANCH SETUP

### TC-SETUP-01: Platform Admin Login & Console Bootstrap
- **Feature**: Platform Console Management
- **Objective**: Verify Platform Super Admin can log in, receive authorized JWT, and view platform dashboard.
- **Preconditions**: Backend running on `8080`, database seeded.
- **Actor/Role**: Platform Admin (`admin@curexal.com`)
- **Test Data**: Email `admin@curexal.com`, Password `password`.
- **Steps**:
  1. Open browser to `http://localhost:5002/login`.
  2. Enter email `admin@curexal.com` and password `password`.
  3. Click **Sign In**.
  4. Inspect browser network console for call to `POST /api/v1/auth/sign-in`.
  5. Observe token redirect and exchange at `POST /api/v1/auth/exchange`.
- **Expected Result**: User is authenticated and automatically redirected to `http://localhost:5002/platform/dashboard`. Platform menu items (Organizations, Users, Marketplace, Pricing, Diagnostics, Audit) are visible in the sidebar.
- **Pass Criteria**: HTTP 200 on exchange, `jwt` cookie established, platform dashboard loads with zero console errors.
- **Fail Criteria**: Redirect loop, 401 Unauthorized, or blank dashboard.
- **Evidence**: Screenshot of platform dashboard with user badge "Platform Owner / Super Admin".
- **Severity**: Critical.

### TC-SETUP-02: Organization & Branch Hierarchy Provisioning
- **Feature**: Multi-Branch Facility Management
- **Objective**: Verify creation of an Organization and multiple facility branches.
- **Preconditions**: Logged in as `admin@curexal.com` or `owner@curexal.space`.
- **Actor/Role**: Organization Owner (`owner@curexal.space`)
- **Test Data**:
  - Organization Name: `Curexal Premier Medical Center` (Slug: `curexal-clinic`)
  - Branch 1: `Main Clinical Center` (Code: `main`, Is Headquarter: `true`)
  - Branch 2: `Victoria Island Annex` (Code: `vi-01`, Is Headquarter: `false`)
- **Steps**:
  1. Navigate to `/organization/branches` in the web platform.
  2. Click **Add Facility Branch**.
  3. Enter Name: `Victoria Island Annex`, Branch Code: `vi-01`, Address: `14 Adeola Odeku St, Victoria Island, Lagos`.
  4. Submit form (`POST /api/v1/organization/branches`).
  5. Verify branch appears in the branches table with status `ACTIVE`.
- **Expected Result**: Branch is persisted in `organization.facility_branches` table with valid UUID.
- **Pass Criteria**: HTTP 201 response, branch list increments count, database row contains correct `organization_id`.
- **Fail Criteria**: 500 error, duplicate code constraint violation unhandled.
- **Evidence**: DB Query: `SELECT id, name, code, is_headquarters FROM organization.facility_branches WHERE code = 'vi-01';`.
- **Severity**: High.

### TC-SETUP-03: Multi-Branch Staff Assignment & Dynamic Context Switching
- **Feature**: Multi-Branch Physician Assignment
- **Objective**: Verify a provider assigned to two branches must explicitly select branch context during sign-in.
- **Preconditions**: User `dr.sarah@curexal.com` is assigned memberships in both `main` and `vi-01` branches.
- **Actor/Role**: Doctor (`dr.sarah@curexal.com`)
- **Test Data**: Email `dr.sarah@curexal.com`, Password `password`.
- **Steps**:
  1. Navigate to `http://localhost:5002/login`.
  2. Enter credentials and click **Sign In**.
  3. Observe API response from `POST /api/v1/auth/sign-in`.
  4. Verify response body contains: `{"status":"branch_selection_required","selectionToken":"..."}`.
  5. Verify UI modal prompts: "Select Active Clinical Branch" showing `Main Clinical Center` and `Victoria Island Annex`.
  6. Select `Victoria Island Annex` (`vi-01`).
  7. Client submits `POST /api/v1/auth/select-branch` with `selectionToken` and `branchId`.
  8. Client exchanges token at `POST /api/v1/auth/exchange`.
- **Expected Result**: Browser redirects to `http://localhost:5002/vi-01/clinical`. Active header displays "Victoria Island Annex".
- **Pass Criteria**: Branch selection enforced by backend; user cannot bypass branch selection without valid token.
- **Fail Criteria**: System auto-assigns branch arbitrarily or crashes on selection.
- **Evidence**: Network tab recording showing the `branch_selection_required` response payload.
- **Severity**: High.

---

## SECTION 4 — AUTHENTICATION

### TC-AUTH-01: Valid Credential Authentication (Argon2id)
- **Feature**: Credential Authentication
- **Objective**: Verify standard sign-in succeeds with valid Argon2id hashed password.
- **Actor/Role**: Receptionist (`receptionist@curexal.com`)
- **Steps**:
  1. Send `POST /api/v1/auth/sign-in` with:
     ```json
     {
       "email": "receptionist@curexal.com",
       "password": "password"
     }
     ```
  2. Verify response status is HTTP 200.
  3. Verify response payload returns `redirect_required` or `branch_selection_required`.
  4. Perform token exchange: `POST /api/v1/auth/exchange`.
- **Expected Result**: Server sets HTTP-only `jwt` cookie with `SameSite=Lax` or `Strict` and `Path=/`.
- **Pass Criteria**: Cookie contains valid signature; subsequent calls to `GET /api/v1/users/me` return HTTP 200.
- **Fail Criteria**: 401 Unauthorized or plain-text password stored in database.
- **Severity**: Critical.

### TC-AUTH-02: Invalid Password Rejection
- **Feature**: Authentication Security
- **Objective**: Verify system denies entry on incorrect password without leaking user existence timing.
- **Actor/Role**: Unauthenticated Client
- **Steps**:
  1. Send `POST /api/v1/auth/sign-in` with `email: "admin@curexal.com"`, `password: "WrongPassword123!"`.
- **Expected Result**: HTTP 401 Unauthorized. Error message: `"Invalid email or password"`.
- **Pass Criteria**: HTTP 401 returned, no session cookie issued.
- **Fail Criteria**: HTTP 200, HTTP 500, or verbose message distinguishing password error from email error.
- **Severity**: High.

### TC-AUTH-03: Nonexistent Email Rejection
- **Feature**: Authentication Security
- **Objective**: Verify non-existent user account cannot trigger internal exception.
- **Actor/Role**: Unauthenticated Client
- **Steps**:
  1. Send `POST /api/v1/auth/sign-in` with `email: "ghost.user.999@curexal.com"`, `password: "anypassword"`.
- **Expected Result**: HTTP 401 Unauthorized. Error message: `"Invalid email or password"`.
- **Pass Criteria**: Consistent HTTP 401.
- **Fail Criteria**: HTTP 404, HTTP 500, or stack trace disclosure.
- **Severity**: High.

### TC-AUTH-04: Session Revocation & Sign-Out
- **Feature**: Session Management
- **Objective**: Verify `POST /api/v1/auth/sign-out` immediately terminates session.
- **Actor/Role**: Doctor (`dr.emeka@curexal.com`)
- **Steps**:
  1. Authenticate and obtain active `jwt` cookie.
  2. Send `POST /api/v1/auth/sign-out`.
  3. Verify `Set-Cookie` header clears the `jwt` cookie (Max-Age=0 / expired).
  4. Immediately send `GET /api/v1/users/me` using the previous token.
- **Expected Result**: `GET /api/v1/users/me` returns HTTP 401 Unauthorized.
- **Pass Criteria**: Immediate invalidation across all subsequent protected endpoints.
- **Fail Criteria**: Old token continues to authenticate requests.
- **Severity**: Critical.

### TC-AUTH-05: Anonymous API Access Rejection
- **Feature**: API Gateway Guard
- **Objective**: Verify protected clinical endpoints reject unauthenticated calls.
- **Actor/Role**: Anonymous Requestor
- **Steps**:
  1. Execute via curl without headers or cookies:
     ```powershell
     curl -i -X GET http://localhost:8080/api/v1/encounters
     curl -i -X GET http://localhost:8080/api/v1/patients
     curl -i -X GET http://localhost:8080/api/v1/billing/invoices
     curl -i -X GET http://localhost:8080/api/v1/providers/profiles
     ```
- **Expected Result**: Every request returns `HTTP 401 Unauthorized`.
- **Pass Criteria**: 100% rejection rate with HTTP 401.
- **Fail Criteria**: Any endpoint returns HTTP 200 or executes business logic.
- **Severity**: Critical.

---

## SECTION 5 — RBAC / CASBIN / AUTHORIZATION

Curexal enforces strict role separation across administrative, clinical, and financial domains. Frontend visibility alone does not constitute authorization.

### Complete Role-Permission Matrix Reference

| Role | Organization Scope | Clinical Read | Clinical Write / Sign | Prescribe | Triage | Billing / POS | Manage Staff |
|---|---|---|---|---|---|---|---|
| **super_admin** | Global | Yes | No | No | No | Read | Yes |
| **owner** | Organization | Yes | Yes | Yes | Yes | Yes | Yes |
| **org_admin** | Organization | Read | No | No | No | Read | Yes |
| **branch_admin**| Branch | Read | No | No | No | No | Branch Only|
| **doctor** | Branch | Yes | Yes | Yes | Yes | No | No |
| **nurse** | Branch | Yes | Draft Only | No | Yes | No | No |
| **receptionist** | Branch | Demographics Only| No | No | No | No | No |
| **cashier** | Branch | Demographics Only| No | No | No | Settle Only | No |
| **patient** | Self | Self Only | No | No | No | View Self | No |

### TC-RBAC-01: Receptionist Prescription Creation Denial (Negative Test)
- **Feature**: Role Privilege Boundaries
- **Objective**: Verify Front Desk Receptionist cannot generate or sign e-prescriptions.
- **Actor/Role**: Receptionist (`receptionist@curexal.com`)
- **Preconditions**: Active valid session for Receptionist. Existing encounter ID `ENC-001`.
- **Steps**:
  1. Using Receptionist credentials, send `POST /api/v1/encounters/ENC-001/prescriptions` with payload:
     ```json
     {
       "medicationName": "Ciprofloxacin 500mg",
       "dosage": "500mg",
       "frequency": "BID",
       "duration": "5 days"
     }
     ```
- **Expected Result**: `HTTP 403 Forbidden` (`Insufficient permissions to perform this action`).
- **Pass Criteria**: HTTP 403 returned; no row inserted into `clinical.prescriptions`.
- **Fail Criteria**: HTTP 200 / 201 Created or prescription saved.
- **Severity**: Critical.

### TC-RBAC-02: Cashier Clinical Record Modification Denial (Negative Test)
- **Feature**: Clinical Documentation Protection
- **Objective**: Verify Cashier cannot update clinical SOAP notes or alter diagnoses.
- **Actor/Role**: Cashier (`cashier@curexal.com`)
- **Steps**:
  1. Using Cashier session, send `PUT /api/v1/encounters/ENC-001/soap` with payload:
     ```json
     {
       "subjective": "Altered by cashier",
       "assessment": "Tampered note"
     }
     ```
- **Expected Result**: `HTTP 403 Forbidden`.
- **Pass Criteria**: HTTP 403 returned; clinical record unaltered.
- **Fail Criteria**: HTTP 200 or notes updated.
- **Severity**: Critical.

### TC-RBAC-03: Doctor Financial POS Settlement Denial (Negative Test)
- **Feature**: Financial Segregation of Duties
- **Objective**: Verify Physician cannot process POS cashier payments or close shifts.
- **Actor/Role**: Doctor (`dr.emeka@curexal.com`)
- **Steps**:
  1. Using Doctor session, send `POST /api/v1/billing/invoices/INV-001/payments` with payload:
     ```json
     {
       "amount": 5000.00,
       "tenderType": "CASH"
     }
     ```
- **Expected Result**: `HTTP 403 Forbidden` (`Insufficient permissions to perform this action`).
- **Pass Criteria**: HTTP 403 returned; invoice remains UNPAID.
- **Fail Criteria**: HTTP 201 Created or payment recorded.
- **Severity**: High.

### TC-RBAC-04: Patient Cross-Patient Record Modification Denial (Negative Test)
- **Feature**: Patient Portal Security
- **Objective**: Verify Patient cannot access or update another patient's medical chart.
- **Actor/Role**: Patient A (`patient.adanna@curexal.com`)
- **Steps**:
  1. Using Patient A portal token, send `GET /api/v1/patients/{patient_B_uuid}`.
  2. Send `PUT /api/v1/patients/{patient_B_uuid}`.
- **Expected Result**: `HTTP 403 Forbidden` or `HTTP 404 Not Found`.
- **Pass Criteria**: Absolute denial; zero disclosure of Patient B identity or clinical notes.
- **Fail Criteria**: Patient B record returned.
- **Severity**: Critical.

---

## SECTION 6 — MULTI-TENANT ISOLATION

This section verifies that data belonging to **Tenant A** (`Curexal Clinic`) is mathematically isolated from **Tenant B** (`Everight Hospital`), and that Branch A1 cannot access Branch A2 data where branch isolation applies.

### TC-TENANT-01: Cross-Tenant Patient Demographic Isolation
- **Feature**: Tenant Isolation
- **Objective**: Verify Tenant A user cannot view Tenant B patient records via direct ID lookup.
- **Preconditions**:
  - Tenant A patient: `Adaeze Nwosu` (`ID: 11111111-1111-1111-1111-111111111111`) in Tenant A.
  - Tenant B patient: `Emeka Okonkwo` (`ID: 22222222-2222-2222-2222-222222222222`) in Tenant B.
- **Actor/Role**: Doctor in Tenant A (`dr.emeka@curexal.com`)
- **Steps**:
  1. Authenticate as Doctor in Tenant A.
  2. Execute `GET /api/v1/patients/22222222-2222-2222-2222-222222222222`.
- **Expected Result**: `HTTP 404 Not Found` (`{"code":"Not Found","message":"Patient not found"}`).
- **Pass Criteria**: Clean HTTP 404 domain error; zero data returned; no HTTP 500 panic.
- **Fail Criteria**: HTTP 200 with patient details or HTTP 500 internal server error.
- **Evidence**: HTTP response body and status code.
- **Severity**: Critical.

### TC-TENANT-02: Header Injection & Tenant ID Spoofing (Attack Test)
- **Feature**: Tenant Context Enforcer
- **Objective**: Verify that passing a foreign `X-Tenant-ID` or `X-Branch-ID` header is ignored when authenticated.
- **Actor/Role**: Receptionist in Tenant A (`receptionist@curexal.com`)
- **Steps**:
  1. Authenticate as Receptionist in Tenant A.
  2. Send `GET /api/v1/patients` with injected header:
     `X-Tenant-ID: {Tenant_B_UUID}`
- **Expected Result**: The server enforces the tenant bound to the authenticated JWT principal. The returned patient list contains ONLY Tenant A patients.
- **Pass Criteria**: Zero records from Tenant B returned; server logs warning of tenant mismatch if applicable.
- **Fail Criteria**: System returns Tenant B patients.
- **Severity**: Critical.

### TC-TENANT-03: Cross-Tenant Invoice Settlement Tampering (Negative Test)
- **Feature**: Billing Isolation
- **Objective**: Verify Cashier in Tenant A cannot settle an invoice belonging to Tenant B.
- **Actor/Role**: Cashier in Tenant A (`cashier@curexal.com`)
- **Steps**:
  1. Obtain Invoice ID `INV-B-999` from Tenant B.
  2. Using Tenant A Cashier session, send `POST /api/v1/billing/invoices/INV-B-999/payments`:
     ```json
     {
       "amount": 10000.00,
       "tenderType": "CASH"
     }
     ```
- **Expected Result**: `HTTP 404 Not Found` or `HTTP 403 Forbidden` (`Invoice not found`).
- **Pass Criteria**: Payment rejected; Tenant B invoice balance unchanged in DB.
- **Fail Criteria**: HTTP 201 Created or payment recorded on Tenant B invoice.
- **Severity**: Critical.

---

## SECTION 7 — PATIENT REGISTRATION / MPI

The Master Patient Index (MPI) ensures unique patient identification, deterministic MRN formatting (`PAT-YYYY-XXXXX`), and duplicate resolution.

### TC-MPI-01: Canonical Patient Registration & Monotonic MRN Generation
- **Feature**: Patient Demographic Intake
- **Objective**: Register a new patient and verify issuance of compliant MRN.
- **Actor/Role**: Receptionist (`receptionist@curexal.com`)
- **Test Data**:
  - First Name: `Chinedu`
  - Last Name: `Eze`
  - Date of Birth: `1988-04-12`
  - Gender: `MALE`
  - Phone Number: `+2348031234567`
  - National Identity Number (NIN): `12345678901`
- **Steps**:
  1. Navigate to `http://localhost:5002/main/reception`.
  2. Click **Register New Patient**.
  3. Fill all fields with test data.
  4. Submit form (`POST /api/v1/patients`).
- **Expected Result**:
  - Server returns `HTTP 201 Created`.
  - Response contains generated `mrn` matching regex `^PAT-\d{4}-\d{5}$` (e.g. `PAT-2026-00001`).
  - Patient record is linked to current tenant ID.
- **Pass Criteria**: HTTP 201, MRN matches format, patient appears in live search.
- **Fail Criteria**: Malformed MRN (e.g. UUID used as MRN), missing tenant ID, or 500 crash.
- **Evidence**: DB Query: `SELECT id, mrn, first_name, last_name, phone_number FROM patient.patients WHERE phone_number = '+2348031234567';`.
- **Severity**: High.

### TC-MPI-02: Probable Duplicate Detection via MPI Engine
- **Feature**: MPI Duplicate Evaluation
- **Objective**: Verify system detects matching demographic records and prevents silent duplicate creation.
- **Actor/Role**: Receptionist (`receptionist@curexal.com`)
- **Test Data**: Same phone (`+2348031234567`) and DOB (`1988-04-12`) as `Chinedu Eze`.
- **Steps**:
  1. In Reception workspace, enter Registration details for `Chinedu Eze` with matching phone number and DOB.
  2. Click **Validate / Check Duplicate** (`POST /api/v1/patients/mpi/evaluate`).
- **Expected Result**:
  - API returns `hasMatch: true`, `matchConfidence: "HIGH"`, `status: "PROBABLE_DUPLICATE"`.
  - UI displays modal: "Potential Existing Patient Record Found" displaying existing MRN.
- **Pass Criteria**: System flags high confidence duplicate; offers "Select Existing" or "Proceed with Override".
- **Fail Criteria**: System auto-creates duplicate record without notification.
- **Severity**: High.

### TC-MPI-03: Nonexistent & Foreign Patient Lookup (Remediation Verification)
- **Feature**: Error Handling & Boundary Defense
- **Objective**: Verify lookup of nonexistent or foreign patient IDs returns standardized `404 Not Found` (never HTTP 500).
- **Actor/Role**: Doctor (`dr.emeka@curexal.com`)
- **Steps**:
  1. Send `GET /api/v1/patients/00000000-0000-0000-0000-000000000000`.
  2. Send `GET /api/v1/patients/invalid-uuid-string-format`.
- **Expected Result**: Both requests return `HTTP 404 Not Found` with body:
  ```json
  {
    "code": "Not Found",
    "message": "Patient not found"
  }
  ```
- **Pass Criteria**: Both calls return clean 404; zero HTTP 500 exceptions in server log.
- **Fail Criteria**: HTTP 500 or SQL syntax error exposed in response.
- **Severity**: Medium.

---

## SECTION 8 — APPOINTMENTS

### TC-APPT-01: In-Person Appointment Booking & Interval Validation
- **Feature**: Appointment Scheduling
- **Objective**: Schedule an in-person consultation and verify time constraint enforcement.
- **Actor/Role**: Receptionist (`receptionist@curexal.com`)
- **Test Data**:
  - Patient ID: Existing patient UUID
  - Provider ID: `dr.emeka` provider UUID
  - Delivery Channel: `in_person`
  - Start Time: Tomorrow at `09:00:00Z`
  - End Time: Tomorrow at `09:30:00Z`
- **Steps**:
  1. Send `POST /api/v1/appointments` with valid payload.
  2. Verify response status is `HTTP 201 Created`.
  3. Attempt invalid appointment where `start_time` is AFTER `end_time` (Start: `10:00`, End: `09:00`).
- **Expected Result**:
  - Valid appointment returns HTTP 201 with status `SCHEDULED`.
  - Invalid appointment returns `HTTP 400 Bad Request` (`Start time must precede end time`).
- **Pass Criteria**: Valid booking persisted; inverted interval strictly blocked.
- **Fail Criteria**: Inverted time appointment accepted into database.
- **Severity**: High.

### TC-APPT-02: Multi-Channel Appointment Support
- **Feature**: Delivery Channel Typing
- **Objective**: Verify appointments accept all 4 strongly typed channels: `in_person`, `video`, `telephone`, `secure_message`.
- **Actor/Role**: Receptionist (`receptionist@curexal.com`)
- **Steps**:
  1. Create 4 separate appointments using each delivery channel code.
  2. Attempt creation with invalid channel `smoke_signal`.
- **Expected Result**:
  - All 4 valid channels return HTTP 201.
  - Invalid channel returns `HTTP 400 Bad Request` (database check constraint or validation failure).
- **Pass Criteria**: Strict enum enforcement.
- **Fail Criteria**: Unsupported channel accepted.
- **Severity**: Medium.

---

## SECTION 9 — PATIENT CHECK-IN / QUEUE

> **ARCHITECTURAL INVARIANT**: Patient check-in creates an operational care request / queue entry (`WAITING_TRIAGE`). It does **NOT** create a clinical encounter.

### TC-QUEUE-01: Scheduled Patient Arrival & Live Queue Placement
- **Feature**: Care Orchestration Queue
- **Objective**: Check in an arrived patient and verify appearance in the Triage Queue.
- **Actor/Role**: Receptionist (`receptionist@curexal.com`)
- **Steps**:
  1. In Reception workspace (`/main/reception`), locate scheduled appointment.
  2. Click **Check In Patient** (`POST /api/v1/orchestration/requests`).
  3. Select Mode: `WALK_IN` or `SCHEDULED_ARRIVAL`.
  4. Submit check-in.
- **Expected Result**:
  - Care request is created with status `WAITING_TRIAGE`.
  - Appointment status updates to `ARRIVED`.
  - No row is created in `encounter.encounters`.
- **Pass Criteria**: Care request created; encounter table count remains unchanged.
- **Fail Criteria**: Clinical encounter created prematurely upon front desk arrival.
- **Severity**: Critical.

### TC-QUEUE-02: Acuity Priority Ordering (RED > YELLOW > GREEN)
- **Feature**: Emergency Acuity Triage Queue
- **Objective**: Verify that high-acuity patients sort to the top of the provider clinical queue.
- **Actor/Role**: Triage Nurse (`nurse.chioma@curexal.com`)
- **Steps**:
  1. Patient 1 checked in at 09:00, triaged as `GREEN` (Routine).
  2. Patient 2 checked in at 09:10, triaged as `YELLOW` (Urgent).
  3. Patient 3 checked in at 09:15, triaged as `RED` (Emergency).
  4. Log in as Doctor (`dr.emeka@curexal.com`) and open `/main/clinical`.
  5. Inspect order of patients in queue.
- **Expected Result**: Queue displays Patient 3 (`RED`) first, followed by Patient 2 (`YELLOW`), then Patient 1 (`GREEN`).
- **Pass Criteria**: Strict clinical acuity sorting regardless of arrival timestamp.
- **Fail Criteria**: Strict FIFO sorting that strands emergency patients behind routine visits.
- **Severity**: High.

---

## SECTION 10 — TRIAGE

### TC-TRIAGE-01: Nurse Vitals Recording & Acuity Assignment
- **Feature**: Clinical Triage Intake
- **Objective**: Record baseline physiological parameters and submit triage assessment.
- **Actor/Role**: Triage Nurse (`nurse.chioma@curexal.com`)
- **Test Data**:
  - Systolic BP: `120` mmHg
  - Diastolic BP: `80` mmHg
  - Heart Rate: `74` bpm
  - Temperature: `36.8` °C
  - Respiratory Rate: `16` breaths/min
  - SpO2: `98` %
  - Weight: `70.5` kg
  - Height: `175` cm
  - Acuity: `YELLOW`
- **Steps**:
  1. In Care Desk (`/main/care-desk`), select patient in `WAITING_TRIAGE`.
  2. Click **Start Triage**.
  3. Enter vitals and observations.
  4. Select Acuity: `YELLOW`.
  5. Click **Submit Triage** (`POST /api/v1/orchestration/requests/{id}/triage`).
- **Expected Result**:
  - HTTP 200 / 201 response.
  - Care request transitions to status `TRIAGED` or `WAITING_CONSULTATION`.
  - Vitals are persisted in `orchestration.care_requests` (or triage observations table).
- **Pass Criteria**: Successful triage submission; patient appears in Doctor's consultation queue.
- **Fail Criteria**: Validation error on valid vitals or 500 error.
- **Severity**: High.

### TC-TRIAGE-02: Physiological Range Boundary Validation (Negative Test)
- **Feature**: Clinical Safety Guard
- **Objective**: Verify system rejects physically impossible vital signs.
- **Actor/Role**: Triage Nurse (`nurse.chioma@curexal.com`)
- **Test Data**:
  - Heart Rate: `-5` bpm
  - Temperature: `150.0` °C
  - Systolic BP: `1000` mmHg
- **Steps**:
  1. Enter impossible vitals and submit triage.
- **Expected Result**: Client-side and server-side validation error (`HTTP 400 Bad Request`).
- **Pass Criteria**: Form blocked with descriptive error message ("Temperature must be between 30°C and 45°C").
- **Fail Criteria**: Out-of-bounds vitals persisted to patient chart.
- **Severity**: Medium.

---

## SECTION 11 — ENCOUNTER / CLINICAL CORE

> **CRITICAL PRODUCTION VERIFICATION**: Verifies that table `encounter.encounters` properly persists `appointment_id`, `encounter_channel`, `signed_at`, `signed_by`, and `closed_at` as reconciled in migration 77.

### TC-ENC-01: Doctor Starts Outpatient Clinical Encounter
- **Feature**: Clinical Consultation Lifecycle
- **Objective**: Initiate an active clinical encounter linked to patient and appointment.
- **Actor/Role**: Doctor (`dr.emeka@curexal.com`)
- **Preconditions**: Patient triaged; appointment exists.
- **Test Data**:
  - Patient ID: Existing patient UUID
  - Provider ID: Provider profile UUID
  - Appointment ID: Appointment UUID
  - Channel: `in_person`
  - Chief Complaint: "Persistent cough and mild fever for 4 days"
- **Steps**:
  1. In Doctor workspace (`/main/clinical`), select triaged patient.
  2. Click **Start Consultation** (`POST /api/v1/encounters`).
  3. Inspect HTTP response.
- **Expected Result**:
  - Server returns `HTTP 201 Created`.
  - Response contains generated encounter UUID and status `in_progress`.
  - Database row exists in `encounter.encounters` with non-null `encounter_channel = 'in_person'`.
- **Pass Criteria**: HTTP 201; no SQL column errors (reconciled schema verified).
- **Fail Criteria**: "column 'appointment_id' does not exist" or HTTP 500.
- **Evidence**: DB Query:
  ```sql
  SELECT id, tenant_id, patient_id, appointment_id, encounter_channel, status, started_at
  FROM encounter.encounters WHERE patient_id = '{patient_uuid}' ORDER BY started_at DESC LIMIT 1;
  ```
- **Severity**: Critical (Blocker).

### TC-ENC-02: Telehealth Session Subordination Invariant
- **Feature**: Architecture Invariant Verification
- **Objective**: Verify that accessing a telehealth link or virtual waiting room does NOT create an encounter.
- **Actor/Role**: Patient / Provider
- **Steps**:
  1. Open virtual waiting room URL: `http://localhost:5003/consultations/{requestId}/room`.
  2. Inspect database count of `encounter.encounters`.
  3. Refresh the page multiple times.
- **Expected Result**: No clinical encounter is inserted merely by loading or joining the video waiting room. Encounters are only instantiated when the clinician explicitly commences care.
- **Pass Criteria**: Encounter count unchanged.
- **Fail Criteria**: Phantom encounter row created upon link click.
- **Severity**: Critical.

---

## SECTION 12 — SOAP / CLINICAL DOCUMENTATION

### TC-SOAP-01: Clinical SOAP Documentation Canvas Draft & Save
- **Feature**: Clinical Notes
- **Objective**: Enter Subjective, Objective, Assessment, and Plan notes, and verify persistence.
- **Actor/Role**: Doctor (`dr.emeka@curexal.com`)
- **Test Data**:
  - Subjective: "Patient reports productive cough with yellowish sputum, sore throat, and intermittent chills. No shortness of breath."
  - Objective: "Pharyngeal erythema with tonsillar hypertrophy (Grade 2). Lungs clear to auscultation bilaterally. Vitals stable."
  - Assessment: "Acute upper respiratory tract infection (URTI) with secondary pharyngitis."
  - Plan: "Prescribe oral Amoxicillin/Clavulanate 625mg BID x 7 days. Paracetamol 1g TID PRN for fever. Increase oral fluid intake."
- **Steps**:
  1. In Encounter Room (`/main/clinical/encounters/{id}`), enter the 4 SOAP quadrants.
  2. Click **Save Draft** (`PUT /api/v1/encounters/{id}/soap`).
  3. Reload the browser page.
- **Expected Result**: Notes reload perfectly from `encounter.encounters` table. Status remains `in_progress`.
- **Pass Criteria**: HTTP 200 response; exact text persisted.
- **Fail Criteria**: Data loss on reload or 500 error.
- **Severity**: High.

---

## SECTION 13 — DIAGNOSIS

### TC-DIAG-01: ICD-10 Primary & Secondary Diagnosis Assignment
- **Feature**: Diagnostic Coding
- **Objective**: Assign standard ICD-10 primary and secondary diagnosis codes to the encounter.
- **Actor/Role**: Doctor (`dr.emeka@curexal.com`)
- **Test Data**:
  - Primary Diagnosis: `J06.9` — *Acute upper respiratory infection, unspecified*
  - Secondary Diagnosis 1: `R50.9` — *Fever, unspecified*
  - Secondary Diagnosis 2: `R05` — *Cough*
- **Steps**:
  1. In Encounter Room, open Diagnosis section.
  2. Search for `J06.9` in ICD-10 search bar. Select as Primary.
  3. Add `R50.9` and `R05` as Secondary.
  4. Submit diagnoses (`POST /api/v1/encounters/{id}/diagnoses`).
- **Expected Result**:
  - Server returns HTTP 200 / 201.
  - `primary_diagnosis_code` and `secondary_diagnoses` (JSONB) are updated on `encounter.encounters`.
- **Pass Criteria**: Diagnoses saved and visible in encounter summary.
- **Fail Criteria**: Primary diagnosis code null or JSONB parsing error.
- **Severity**: High.

---

## SECTION 14 — PRESCRIPTIONS

### TC-RX-01: Physician e-Prescribing Workflow
- **Feature**: Medication Prescribing
- **Objective**: Authorize medication prescription with complete dosage, route, frequency, and instructions.
- **Actor/Role**: Doctor (`dr.emeka@curexal.com`)
- **Test Data**:
  - Medication: `Amoxicillin/Clavulanic Acid 625mg Tab`
  - Dosage: `625mg`
  - Route: `Oral`
  - Frequency: `Twice daily (every 12 hours)`
  - Duration: `7 days`
  - Dispense Quantity: `14 tablets`
  - Instructions: `Take with meals to avoid gastric discomfort`
- **Steps**:
  1. In Encounter Room, open Prescriptions tab.
  2. Click **Add Medication**.
  3. Fill all fields and click **Authorize Prescription** (`POST /api/v1/encounters/{id}/prescriptions`).
- **Expected Result**:
  - Prescription is stored with status `PENDING` or `ACTIVE`.
  - Linked to current encounter ID, patient ID, and provider ID.
- **Pass Criteria**: HTTP 201 Created; record stored in database.
- **Fail Criteria**: 500 error or missing prescription details.
- **Severity**: High.

### TC-RX-02: Clinical Finalization & Immutable Record Sign-Off
- **Feature**: Encounter Completion & Integrity Locking
- **Objective**: Sign and finalize encounter; verify record becomes read-only.
- **Actor/Role**: Doctor (`dr.emeka@curexal.com`)
- **Steps**:
  1. In Encounter Room, click **Complete & Sign Consultation** (`POST /api/v1/encounters/{id}/complete`).
  2. Inspect database record for `signed_at`, `signed_by`, and `closed_at`.
  3. Attempt to send `PUT /api/v1/encounters/{id}/soap` to modify signed notes.
- **Expected Result**:
  - Status transitions to `completed` or `closed`.
  - `signed_at` timestamp and `signed_by` doctor UUID are populated.
  - Modification attempt returns `HTTP 400 Bad Request` or `HTTP 403 Forbidden` (`Cannot modify closed encounter`).
- **Pass Criteria**: Record locked; audit timestamps verified in PostgreSQL.
- **Fail Criteria**: Signed encounter accepts arbitrary note overwrites.
- **Severity**: Critical.

---

## SECTION 15 — BILLING / POS

> **CLINICAL-FINANCIAL INTEGRATION**: Completing an encounter automatically instantiates an invoice in `billing.invoices` with category `CONSULTATION` and default fee (e.g. 5,000 NGN).

### TC-BILL-01: Automated Invoice Generation upon Consultation Completion
- **Feature**: Encounter Billing Trigger
- **Objective**: Verify invoice generation following encounter sign-off.
- **Actor/Role**: Doctor (`dr.emeka@curexal.com`) / System
- **Steps**:
  1. Finalize encounter as tested in TC-RX-02.
  2. Inspect `billing.invoices` table for new invoice linked to `encounter_id`.
- **Expected Result**:
  - Invoice created with unique `invoice_number` (format `INV-YYYY-XXXXX`).
  - `status = 'UNPAID'`.
  - `total_amount = 5000.00`, `balance_due = 5000.00`, `amount_paid = 0.00`.
  - Line item in `billing.invoice_items`: Description: `Physician Clinical Consultation`, Unit Price: `5000.00`.
- **Pass Criteria**: Invoice and line item generated accurately without cashier intervention.
- **Fail Criteria**: No invoice generated, or balance due mathematically incorrect.
- **Evidence**: DB Query: `SELECT id, invoice_number, total_amount, balance_due, status FROM billing.invoices WHERE encounter_id = '{id}';`.
- **Severity**: High.

### TC-BILL-02: Cashier Single-Tender Payment Settlement (Cash)
- **Feature**: Cashier POS Settlement
- **Objective**: Settle unpaid consultation invoice with Cash tender and issue receipt.
- **Actor/Role**: Cashier (`cashier@curexal.com`)
- **Test Data**: Invoice ID from TC-BILL-01, Amount: `5000.00`, Tender: `CASH`.
- **Steps**:
  1. In Cashier Billing Workspace (`/main/billing`), locate patient invoice.
  2. Click **Settle Invoice**.
  3. Enter Payment Amount: `5000.00`, Select Tender: `CASH`.
  4. Submit payment (`POST /api/v1/billing/invoices/{id}/payments`).
- **Expected Result**:
  - Server returns `HTTP 201 Created`.
  - Response contains `receiptNumber` (format `REC-YYYY-XXXXX`).
  - Invoice status transitions to `PAID`. `amount_paid = 5000.00`, `balance_due = 0.00`.
- **Pass Criteria**: Invoice cleared; receipt generated with matching total.
- **Fail Criteria**: Status remains UNPAID or mathematical balance due is non-zero.
- **Severity**: High.

### TC-BILL-03: Split Tender Payment Processing (Cash + POS Card)
- **Feature**: Split Tender Settlement
- **Objective**: Settle an invoice using multiple payment tenders (e.g. 3,000 NGN Cash + 2,000 NGN Card).
- **Actor/Role**: Cashier (`cashier@curexal.com`)
- **Test Data**: Total: `5000.00`. Tender: `SPLIT` (Cash: 3000, Card: 2000).
- **Steps**:
  1. Submit payment with payload:
     ```json
     {
       "amount": 5000.00,
       "tenderType": "SPLIT",
       "tenderBreakdown": {
         "cash": 3000.00,
         "pos_card": 2000.00
       }
     }
     ```
- **Expected Result**: Payment recorded; breakdown stored in `tender_breakdown` JSONB. Invoice marked `PAID`.
- **Pass Criteria**: Sum of tender breakdown exactly matches `amount_paid`; status is `PAID`.
- **Fail Criteria**: Mismatched totals allowed or tender breakdown dropped.
- **Severity**: High.

### TC-BILL-04: Partial Payment Handling
- **Feature**: Partial Invoice Settlement
- **Objective**: Verify partial payment updates balance due and sets status to `PARTIALLY_PAID`.
- **Actor/Role**: Cashier (`cashier@curexal.com`)
- **Test Data**: Total: `5000.00`. Payment: `2000.00`.
- **Steps**:
  1. Submit payment of `2000.00` against `5000.00` invoice.
- **Expected Result**:
  - Status becomes `PARTIALLY_PAID`.
  - `amount_paid = 2000.00`.
  - `balance_due = 3000.00`.
- **Pass Criteria**: Exact subtraction verified (`5000.00 - 2000.00 = 3000.00`).
- **Fail Criteria**: Invoice marked PAID prematurely or balance due miscalculated.
- **Severity**: High.

---

## SECTION 16 — PATIENT PORTAL

### TC-PORTAL-01: Patient Portal Authentication & Profile View
- **Feature**: Patient Self-Service Access
- **Objective**: Patient logs in via portal authentication and views own health record.
- **Actor/Role**: Patient (`patient.adanna@curexal.com`)
- **Steps**:
  1. Open `http://localhost:5003/login`.
  2. Request login OTP / enter Portal PIN.
  3. View Patient Portal Dashboard (`/dashboard`).
  4. Verify patient demographics, MRN, and visit history render.
- **Expected Result**: Patient views personal profile and care journey.
- **Pass Criteria**: Patient data matches database; portal token established.
- **Fail Criteria**: Blank screen or unauthorized redirect.
- **Severity**: High.

### TC-PORTAL-02: Patient Cross-Record Enumeration Defense (Attack Test)
- **Feature**: PHI Perimeter Security
- **Objective**: Verify that Patient A cannot query Patient B records by manipulating URLs or API parameters.
- **Actor/Role**: Patient A
- **Steps**:
  1. While logged in as Patient A on port 5003, inspect network requests.
  2. Issue direct fetch to `GET /api/v1/patient/profile` with Patient A token: returns Patient A.
  3. Attempt to request `GET /api/v1/patients/{patient_B_uuid}` using Patient A token.
- **Expected Result**: Request is rejected with `HTTP 403 Forbidden` by `PatientGuard`.
- **Pass Criteria**: 100% isolation; Patient A cannot view any other patient's demographics, prescriptions, or notes.
- **Fail Criteria**: Foreign patient data returned.
- **Severity**: Critical (HIPAA / Privacy Blocker).

---

## SECTION 17 — AUDIT LOGGING

### TC-AUDIT-01: Comprehensive Clinical & Financial Mutation Auditing
- **Feature**: Regulatory Audit Trail
- **Objective**: Verify that all critical operations generate immutable audit records.
- **Actor/Role**: System / Platform Compliance Auditor
- **Steps**:
  1. Perform actions: Login, Patient Creation, SOAP Update, Prescription Sign, Payment Settlement.
  2. Query `audit.audit_logs` table in PostgreSQL.
- **Expected Result**: Audit entries exist with:
  - `actor_id`: Authenticated user UUID
  - `actor_role`: Role at time of action
  - `action`: e.g. `auth.login`, `patient.created`, `encounter.soap_updated`, `billing.payment_processed`
  - `resource_id`: Target entity UUID
  - `tenant_id`: Organization / branch UUID
  - `timestamp`: ISO-8601 UTC timestamp
- **Pass Criteria**: Every action produces a corresponding audit entry.
- **Fail Criteria**: Silent mutations without audit trail.
- **Evidence**: DB Query:
  ```sql
  SELECT action, actor_role, resource_type, resource_id, created_at
  FROM audit.audit_logs ORDER BY created_at DESC LIMIT 10;
  ```
- **Severity**: Critical.

### TC-AUDIT-02: Audit Log Immutability Verification (Security Test)
- **Feature**: Tamper-Proof Audit Vault
- **Objective**: Verify that audit logs cannot be updated or deleted.
- **Actor/Role**: Database Operator / Rogue Actor
- **Steps**:
  1. Attempt direct SQL update:
     ```sql
     UPDATE audit.audit_logs SET action = 'tampered' WHERE id = (SELECT id FROM audit.audit_logs LIMIT 1);
     ```
  2. Attempt direct SQL deletion:
     ```sql
     DELETE FROM audit.audit_logs WHERE id = (SELECT id FROM audit.audit_logs LIMIT 1);
     ```
- **Expected Result**: Operation is rejected by database trigger or security rule, or immediately flagged in security log.
- **Pass Criteria**: Database prevents mutation or preserves audit integrity.
- **Fail Criteria**: Silent deletion of audit logs without trace.
- **Severity**: Critical.

---

## SECTION 18 — FRONTEND / UI TESTING

### TC-UI-01: Full Responsive Navigation & Role-Based Sidebar Rendering
- **Feature**: Workspace UI Shell
- **Objective**: Verify sidebar displays only authorized navigation nodes per role.
- **Steps**:
  1. Log in as **Receptionist**: Sidebar shows Reception, Patients, Appointments, Queue. Does NOT show Clinical EMR or POS Settlement.
  2. Log in as **Doctor**: Sidebar shows Care Desk, Clinical Queue, Consultation Room. Does NOT show Billing or Platform Config.
  3. Log in as **Cashier**: Sidebar shows Billing, Invoices, POS Settlement, Receipts. Does NOT show Clinical SOAP or Triage.
  4. Log in as **Platform Admin**: Sidebar shows Platform Console, Organizations, Diagnostics, Global Audit.
- **Expected Result**: Dynamic navigation menus tailored strictly to active persona.
- **Pass Criteria**: No unauthorized buttons rendered; no broken routing links.
- **Fail Criteria**: Unauthorized menu items exposed or UI crashes on navigation.
- **Severity**: Medium.

### TC-UI-02: Browser History, Refresh & Direct Deep Linking
- **Feature**: SPA Route Hydration
- **Objective**: Verify refreshing the page on deep clinical routes (`/:branchSlug/clinical/encounters/:id`) preserves context.
- **Steps**:
  1. Navigate to `/main/clinical/encounters/{valid_encounter_id}` as Doctor.
  2. Hit browser Refresh (F5).
  3. Click browser Back button, then Forward button.
- **Expected Result**: Page reloads smoothly; session remains active; encounter data re-populates without redirect to login.
- **Pass Criteria**: Zero white-screens; zero loss of active form state on reload.
- **Fail Criteria**: Blank screen, redirect loop, or 404 page not found.
- **Severity**: High.

---

## SECTION 19 — API TESTING

The following table documents the core Phase 1 API endpoints discovered from the codebase for direct curl / Postman regression testing.

| # | HTTP Method | Endpoint Route | Auth Required? | Authorized Roles | Expected Success | Expected Denial |
|---|---|---|---|---|---|---|
| **1** | `POST` | `/api/v1/auth/sign-in` | No | Public | 200 OK | 401 Unauthorized |
| **2** | `POST` | `/api/v1/auth/select-branch` | Yes (Selection Token) | Multi-Branch Users | 200 OK | 400 Bad Request |
| **3** | `POST` | `/api/v1/auth/exchange` | No | Exchange Token | 200 OK | 400 Bad Request |
| **4** | `POST` | `/api/v1/auth/sign-out` | Yes | All Authenticated | 200 OK | 401 Unauthorized |
| **5** | `GET` | `/api/v1/users/me` | Yes | All Authenticated | 200 OK | 401 Unauthorized |
| **6** | `GET` | `/api/v1/patients` | Yes | Staff | 200 OK | 401 / 403 |
| **7** | `POST` | `/api/v1/patients` | Yes | Receptionist, Nurse, Doctor, Admin | 201 Created | 400 Bad Request |
| **8** | `GET` | `/api/v1/patients/:id` | Yes | Staff | 200 OK | 404 Not Found |
| **9** | `POST` | `/api/v1/patients/mpi/evaluate` | Yes | Receptionist, Admin | 200 OK | 400 Bad Request |
| **10** | `POST` | `/api/v1/appointments` | Yes | Receptionist, Doctor, Admin | 201 Created | 400 Bad Request |
| **11** | `GET` | `/api/v1/appointments` | Yes | Staff | 200 OK | 401 Unauthorized |
| **12** | `POST` | `/api/v1/orchestration/requests` | Yes | Receptionist, Staff | 201 Created | 400 Bad Request |
| **13** | `POST` | `/api/v1/orchestration/requests/:id/triage` | Yes | Nurse, Doctor | 200 OK | 403 Forbidden |
| **14** | `POST` | `/api/v1/encounters` | Yes | Doctor, Clinical Admin | 201 Created | 403 Forbidden |
| **15** | `GET` | `/api/v1/encounters/:id` | Yes | Doctor, Nurse | 200 OK | 404 Not Found |
| **16** | `PUT` | `/api/v1/encounters/:id/soap` | Yes | Doctor | 200 OK | 403 Forbidden |
| **17** | `POST` | `/api/v1/encounters/:id/diagnoses` | Yes | Doctor | 201 Created | 403 Forbidden |
| **18** | `POST` | `/api/v1/encounters/:id/prescriptions` | Yes | Doctor | 201 Created | 403 Forbidden |
| **19** | `POST` | `/api/v1/encounters/:id/complete` | Yes | Doctor | 200 OK | 403 Forbidden |
| **20** | `GET` | `/api/v1/billing/invoices` | Yes | Cashier, Finance, Owner | 200 OK | 403 Forbidden |
| **21** | `GET` | `/api/v1/billing/invoices/:id` | Yes | Cashier, Finance, Owner | 200 OK | 404 Not Found |
| **22** | `POST` | `/api/v1/billing/invoices/:id/payments` | Yes | Cashier, Finance | 201 Created | 403 Forbidden |
| **23** | `GET` | `/api/v1/billing/receipts/:receiptNumber`| Yes | Cashier, Finance | 200 OK | 404 Not Found |
| **24** | `GET` | `/api/v1/providers/profiles` | Yes | Staff (`users:read`) | 200 OK | 401 Unauthorized |
| **25** | `POST` | `/api/v1/providers/profiles` | Yes | Admin (`users:write`) | 201 Created | 403 Forbidden |
| **26** | `GET` | `/api/v1/audit-logs/tenant` | Yes | Owner, Compliance (`audit:read`) | 200 OK | 403 Forbidden |

---

## SECTION 20 — DATABASE VERIFICATION

Direct SQL verification ensures that application logic maps 1:1 to physical database constraints, schema migrations, and relational integrity.

### TV-DB-01: Encounter Schema Columns & Integrity Check
Run via `psql`:
```sql
SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns
WHERE table_schema = 'encounter' AND table_name = 'encounters'
ORDER BY ordinal_position;
```
*Mandatory Columns Checked*:
- `appointment_id` (`uuid`, nullable)
- `encounter_channel` (`varchar`, non-nullable, default `'in_person'`)
- `signed_at` (`timestamptz`, nullable)
- `signed_by` (`uuid`, nullable)
- `closed_at` (`timestamptz`, nullable)

### TV-DB-02: Foreign Key Constraints Validation
```sql
SELECT conname, contype, pg_get_constraintdef(oid)
FROM pg_constraint
WHERE conrelid = 'encounter.encounters'::regclass;
```
*Verified Constraints*:
- `encounters_appointment_id_fkey`: `FOREIGN KEY (appointment_id) REFERENCES operations.appointments(id) ON DELETE SET NULL`
- `encounters_signed_by_fkey`: `FOREIGN KEY (signed_by) REFERENCES identity.users(id) ON DELETE SET NULL`
- `encounters_patient_id_fkey`: `FOREIGN KEY (patient_id) REFERENCES patient.patients(id) ON DELETE CASCADE`
- `encounters_tenant_id_fkey`: `FOREIGN KEY (tenant_id) REFERENCES organization.facility_branches(id) ON DELETE CASCADE`
- `chk_encounters_channel`: `CHECK (encounter_channel IN ('in_person', 'video', 'telephone', 'secure_message'))`

### TV-DB-03: Billing Ledger Mathematical Invariant Verification
```sql
SELECT 
    invoice_number,
    subtotal,
    discount_amount,
    tax_amount,
    total_amount,
    amount_paid,
    balance_due,
    status
FROM billing.invoices
WHERE (subtotal - discount_amount + tax_amount) != total_amount
   OR (amount_paid + balance_due) != total_amount;
```
*Pass Criteria*: Query returns **0 rows**. Any returned row represents an illegal financial discrepancy.

---

## SECTION 21 — SECURITY / ATTACK TESTS

Execute this dedicated attack-path checklist to stress-test system boundary controls.

| Attack ID | Threat Vector / Description | Attack Action / Method | Expected Server Behavior | Actual Result | Status |
|---|---|---|---|---|---|
| **ATK-01** | Anonymous Provider Harvester | `GET /api/v1/providers/profiles` without token | Rejection with HTTP 401 | 401 Unauthorized | [ ] PASS |
| **ATK-02** | Anonymous Patient Record Query | `GET /api/v1/patients` without token | Rejection with HTTP 401 | 401 Unauthorized | [ ] PASS |
| **ATK-03** | Privilege Escalation (Prescription) | Receptionist sends `POST /api/v1/encounters/{id}/prescriptions` | Rejection with HTTP 403 | 403 Forbidden | [ ] PASS |
| **ATK-04** | Unauthorized Clinical Modification | Cashier sends `PUT /api/v1/encounters/{id}/soap` | Rejection with HTTP 403 | 403 Forbidden | [ ] PASS |
| **ATK-05** | Cross-Branch Infiltration | Doctor at Branch 1 queries Branch 2 unassigned patient | Rejection with HTTP 403 / 404 | 404 / 403 | [ ] PASS |
| **ATK-06** | Cross-Tenant Data Siphon | Tenant A user queries Tenant B patient UUID | Rejection with HTTP 404 domain error | 404 Not Found | [ ] PASS |
| **ATK-07** | Cross-Tenant Billing Tampering | Tenant A cashier settles Tenant B invoice | Rejection with HTTP 404 / 403 | 404 / 403 | [ ] PASS |
| **ATK-08** | Patient Portal Boundary Break | Patient A queries Patient B chart via portal API | Rejection with HTTP 403 (`PatientGuard`) | 403 Forbidden | [ ] PASS |
| **ATK-09** | Malformed UUID Injection | `GET /api/v1/patients/not-a-valid-uuid` | Handled cleanly as HTTP 404 (No 500) | 404 Not Found | [ ] PASS |
| **ATK-10** | Header Tenant Spoofing | Pass `X-Tenant-ID: {Foreign_UUID}` with valid JWT | Ignored; server binds to JWT claims | 200 (Scoped) | [ ] PASS |
| **ATK-11** | Header Branch Spoofing | Pass `X-Branch-ID: {Unassigned_Branch}` with valid JWT | Rejection with HTTP 403 | 403 Forbidden | [ ] PASS |
| **ATK-12** | Expired JWT Replay | Send request with expired session cookie | Rejection with HTTP 401 | 401 Unauthorized | [ ] PASS |
| **ATK-13** | Tampered JWT Signature | Alter signature bits in `jwt` cookie | Cryptographic validation fails; HTTP 401 | 401 Unauthorized | [ ] PASS |
| **ATK-14** | Direct UI Route Bypass | Nurse navigates directly to `/platform/dashboard` | `PlatformGuard` blocks; redirects to login | Redirected | [ ] PASS |
| **ATK-15** | Direct API Negative Quantities | Send `quantity: -5` in invoice or prescription | Validation rejects with HTTP 400 | 400 Bad Request | [ ] PASS |
| **ATK-16** | Inverted Appointment Time Range | Send `start_time > end_time` | Validation rejects with HTTP 400 | 400 Bad Request | [ ] PASS |
| **ATK-17** | SQL Injection in Search Bar | Send `' OR '1'='1` in patient search query | Parameterized query protects DB; 0 rows | 200 (0 results) | [ ] PASS |

---

## SECTION 22 — GOLDEN PATH

Follow this sequential, end-to-end operational script from organization creation through clinical consultation to financial settlement and audit verification.

```
[Platform Admin] Creates Organization & Facility Branch
      │
      ▼
[Organization Owner] Invites Receptionist, Nurse, Doctor, Cashier
      │
      ▼
[Receptionist] Registers Patient → Generates MRN "PAT-2026-00042"
      │
      ▼
[Receptionist] Books Appointment → Checks In Patient
      │
      ▼
[Patient Queue] Transitions to "WAITING_TRIAGE"
      │
      ▼
[Triage Nurse] Records Vitals (BP, Temp, Pulse) → Assigns Acuity "YELLOW"
      │
      ▼
[Doctor] Opens Consultation Room → Starts Encounter ("in_person")
      │
      ▼
[Doctor] Enters SOAP Notes → Assigns ICD-10 Diagnosis "J02.9"
      │
      ▼
[Doctor] Authorizes e-Prescription (Amoxicillin 500mg) → Signs Encounter
      │
      ▼
[Billing Engine] Automatically Generates Invoice "INV-2026-00101" (5,000 NGN)
      │
      ▼
[Cashier] Opens POS Terminal → Collects Cash Payment → Issues Receipt
      │
      ▼
[Patient] Logs into Patient Portal → Views Summary & Prescription
      │
      ▼
[Auditor] Inspects Immutable Audit Trail → 12 Consecutive Events Verified
```

### Detailed Golden Path Execution Protocol

#### Phase A: Administrative Scoping
1. Log into `http://localhost:5002/login` as Platform Admin (`admin@curexal.com`).
2. Navigate to `/platform/organizations`. Confirm `Curexal Premier Medical Center` is active.
3. Verify Branch `Main Clinical Center` (`main`) exists.
4. Log out.

#### Phase B: Patient Intake & MPI
5. Log into `/login` as Receptionist (`receptionist@curexal.com`).
6. Navigate to `/main/reception`. Click **Register Patient**.
7. Enter:
   - First Name: `Adaeze`
   - Last Name: `Nwosu`
   - DOB: `1992-08-24`
   - Gender: `FEMALE`
   - Phone: `+2348029988776`
8. Click **Submit**. Verify MRN generated: e.g. `PAT-2026-00042`.
9. Click **Schedule Appointment**:
   - Provider: `Dr. Emeka`
   - Channel: `in_person`
   - Date: Today, 30-minute block.
10. Click **Check In**. Verify patient appears in Live Queue with status `WAITING_TRIAGE`.
11. Log out.

#### Phase C: Clinical Triage
12. Log into `/login` as Triage Nurse (`nurse.chioma@curexal.com`).
13. Navigate to `/main/care-desk`.
14. Locate `Adaeze Nwosu` in `WAITING_TRIAGE`. Click **Start Triage**.
15. Enter Vitals:
    - BP: `118/78` mmHg, Pulse: `72` bpm, Temp: `36.7` °C, SpO2: `99` %, Resp: `16`.
    - Acuity: `YELLOW` (Urgent).
16. Click **Submit Triage Assessment**.
17. Verify queue transitions to `WAITING_CONSULTATION`.
18. Log out.

#### Phase D: Doctor Consultation & Clinical Finalization
19. Log into `/login` as Doctor (`dr.emeka@curexal.com`).
20. Navigate to `/main/clinical`. Locate `Adaeze Nwosu`.
21. Click **Start Encounter**. Verify encounter room loads.
22. Fill SOAP Notes Canvas:
    - **S**: Patient complains of severe sore throat, painful swallowing for 3 days.
    - **O**: Bilateral tonsillar enlargement with purulent exudate. Afebrile currently.
    - **A**: Acute tonsillopharyngitis (presumed bacterial).
    - **P**: Amoxicillin 500mg PO TID x 7 days. Warm saline gargles. Return if dyspnea occurs.
23. Add ICD-10 Diagnosis:
    - Search: `J03.90` (*Acute tonsillitis, unspecified*). Select as **Primary Diagnosis**.
24. Create Prescription:
    - Medication: `Amoxicillin 500mg Capsules`
    - Dose: `500mg`, Frequency: `TID (three times daily)`, Duration: `7 days`, Quantity: `21 capsules`.
25. Click **Complete & Sign Consultation**.
26. Verify encounter status shows `closed` / `completed`.
27. Log out.

#### Phase E: Cashier POS Settlement
28. Log into `/login` as Cashier (`cashier@curexal.com`).
29. Navigate to `/main/billing`.
30. Locate invoice for `Adaeze Nwosu` (Amount: `5,000.00 NGN`, Status: `UNPAID`).
31. Click **Settle Invoice**. Select Tender: `CASH`, Amount: `5000.00`.
32. Click **Process Payment**.
33. Verify receipt dialog displays Receipt Number (e.g. `REC-2026-00018`), Amount Paid: `5,000.00 NGN`, Balance Due: `0.00 NGN`.
34. Click **Print Receipt**.
35. Log out.

#### Phase F: Patient Portal Verification
36. Open `http://localhost:5003/login`. Log in as `Adaeze Nwosu`.
37. Navigate to `/dashboard`.
38. Verify Consultation visit appears in Visit History.
39. Verify Prescription for `Amoxicillin 500mg` is displayed with Doctor's instructions.
40. Verify Invoice is displayed as `PAID`.
41. Log out.

#### Phase G: Audit Trail Verification
42. Connect to database and execute:
    ```sql
    SELECT action, actor_role, resource_type, created_at 
    FROM audit.audit_logs 
    WHERE tenant_id = (SELECT id FROM organization.facility_branches WHERE code = 'main')
    ORDER BY created_at DESC LIMIT 12;
    ```
43. Verify all steps from registration to payment are recorded.

---

## SECTION 23 — NEGATIVE PATH

| Step | Persona | Attempted Forbidden Action | Expected Status | Expected Error Response | Verified? |
|---|---|---|---|---|---|
| **NP-01** | Anonymous | `GET /api/v1/patients` | `401 Unauthorized` | `{"message":"User not authenticated"}` | [ ] |
| **NP-02** | Anonymous | `GET /api/v1/providers/profiles` | `401 Unauthorized` | `{"message":"User not authenticated"}` | [ ] |
| **NP-03** | Receptionist | `POST /api/v1/encounters/{id}/prescriptions` | `403 Forbidden` | `{"message":"Insufficient permissions to perform this action"}` | [ ] |
| **NP-04** | Nurse | `POST /api/v1/encounters/{id}/prescriptions` | `403 Forbidden` | `{"message":"Insufficient permissions to perform this action"}` | [ ] |
| **NP-05** | Cashier | `PUT /api/v1/encounters/{id}/soap` | `403 Forbidden` | `{"message":"Insufficient permissions to perform this action"}` | [ ] |
| **NP-06** | Cashier | `POST /api/v1/encounters/{id}/diagnoses` | `403 Forbidden` | `{"message":"Insufficient permissions to perform this action"}` | [ ] |
| **NP-07** | Doctor | `POST /api/v1/billing/invoices/{id}/payments` | `403 Forbidden` | `{"message":"Insufficient permissions to perform this action"}` | [ ] |
| **NP-08** | Doctor Branch A | Query patient of Branch B (unassigned) | `404 / 403` | `{"message":"Patient not found"}` | [ ] |
| **NP-09** | Tenant A User | Query patient of Tenant B | `404 Not Found` | `{"message":"Patient not found"}` | [ ] |
| **NP-10** | Tenant A Cashier| Settle invoice of Tenant B | `404 Not Found` | `{"message":"Invoice not found"}` | [ ] |
| **NP-11** | Patient A | Query Patient B profile | `403 Forbidden` | `{"message":"Access denied"}` | [ ] |
| **NP-12** | Doctor | Modify notes of a `closed` encounter | `400 / 403` | `{"message":"Cannot modify finalized clinical record"}` | [ ] |
| **NP-13** | Receptionist | Book appointment with `start_time > end_time`| `400 Bad Request` | `{"message":"Start time must precede end time"}` | [ ] |
| **NP-14** | Cashier | Settle payment with `amount <= 0` | `400 Bad Request` | `{"message":"Payment amount must be greater than zero"}` | [ ] |

---

## SECTION 24 — PHASE 1 FEATURE MATRIX

Complete manual test tracking table for all 33 Phase 1 features. **Do not mark anything PASS prior to physical test execution.**

| Feature # | Feature Description | Applicable Test IDs | Execution Status | Supporting Evidence | Notes / Observations |
|---|---|---|---|---|---|
| **1** | Authentication (Argon2id) | TC-AUTH-01, TC-AUTH-02 | NOT TESTED | | |
| **2** | Login & Credential Verification | TC-AUTH-01, TC-AUTH-03 | NOT TESTED | | |
| **3** | Session Management & Tokens | TC-AUTH-01, TC-AUTH-04 | NOT TESTED | | |
| **4** | Logout & Token Invalidation | TC-AUTH-04 | NOT TESTED | | |
| **5** | Organization Management | TC-SETUP-01, TC-SETUP-02 | NOT TESTED | | |
| **6** | Clinic & Branch Facility Governance | TC-SETUP-02, TC-SETUP-03 | NOT TESTED | | |
| **7** | Role-Based Access Control (RBAC) | TC-RBAC-01..04 | NOT TESTED | | |
| **8** | Scope-Aware Permissions | TC-RBAC-01..04, TC-TENANT-01 | NOT TESTED | | |
| **9** | Audit Logging Engine & Immutability | TC-AUDIT-01, TC-AUDIT-02 | NOT TESTED | | |
| **10** | Multi-Tenant Isolation | TC-TENANT-01..03, ATK-06..07| NOT TESTED | | |
| **11** | Staff Management & Memberships | TC-SETUP-02, TC-SETUP-03 | NOT TESTED | | |
| **12** | Provider Profiles & Management | ATK-01, TC-API-24..25 | NOT TESTED | | Post-remediation verification |
| **13** | Patient Demographic Registration | TC-MPI-01 | NOT TESTED | | |
| **14** | Monotonic MRN Generation | TC-MPI-01 | NOT TESTED | | Format: `PAT-YYYY-XXXXX` |
| **15** | Duplicate Detection & MPI Engine | TC-MPI-02 | NOT TESTED | | |
| **16** | Patient Lookup & Search | TC-MPI-01, TC-MPI-03 | NOT TESTED | | Nonexistent lookup returns 404 |
| **17** | Appointment Scheduling | TC-APPT-01, TC-APPT-02 | NOT TESTED | | |
| **18** | Provider Availability Constraints | TC-APPT-01 | NOT TESTED | | |
| **19** | Patient Check-in & Live Queue | TC-QUEUE-01, TC-QUEUE-02 | NOT TESTED | | |
| **20** | Encounter Creation | TC-ENC-01, TC-ENC-02 | NOT TESTED | | Post-remediation verification |
| **21** | Encounter Lifecycle Transitions | TC-ENC-01, TC-RX-02 | NOT TESTED | | |
| **22** | SOAP Notes Canvas | TC-SOAP-01 | NOT TESTED | | |
| **23** | Triage Vitals & Observations | TC-TRIAGE-01, TC-TRIAGE-02 | NOT TESTED | | |
| **24** | ICD-10 Diagnoses Management | TC-DIAG-01 | NOT TESTED | | |
| **25** | Medication & e-Prescribing | TC-RX-01 | NOT TESTED | | |
| **26** | Clinical Record Finalization & Sign-off | TC-RX-02 | NOT TESTED | | Post-remediation verification |
| **27** | POS Billing Engine | TC-BILL-01..04 | NOT TESTED | | |
| **28** | Invoice Generation & Line Items | TC-BILL-01 | NOT TESTED | | Automated on encounter close |
| **29** | Payment Recording & Split Tenders | TC-BILL-02, TC-BILL-03 | NOT TESTED | | |
| **30** | Patient Visit Summary & POS Receipt | TC-BILL-02, TC-PORTAL-01 | NOT TESTED | | |
| **31** | Clinical & Financial Audit Trail | TC-AUDIT-01 | NOT TESTED | | |
| **32** | DB Schema Migration Integrity | TV-DB-01..03 | NOT TESTED | | Goose version 77 verified |
| **33** | API Gateway Authentication Enforcer | TC-AUTH-05, ATK-01..02 | NOT TESTED | | |

*Permitted Status Values*: `NOT TESTED` | `PASS` | `PARTIAL` | `FAIL` | `BLOCKED`

---

## SECTION 25 — EVIDENCE REQUIREMENTS

To maintain clinical audit integrity, capture the following artifacts for every test case executed. Store all test artifacts in `docs/testing/evidence/{TEST_ID}/`.

1. **Browser Screenshots**: Full-window screenshot capturing browser address bar, active user persona, timestamp, and visual UI confirmation.
2. **HTTP API Responses**: Raw HTTP response header and JSON body output saved as `.json` or `.txt`.
3. **HTTP Status Code**: Document exact numeric status (e.g. `200`, `201`, `400`, `401`, `403`, `404`).
4. **Database Query Output**: Direct SQL console output from PostgreSQL validating row persistence, updated timestamps, and foreign key references.
5. **Audit Log Snapshot**: JSON snapshot of `audit.audit_logs` record generated by the test.
6. **Unique Entity Identifiers**:
   - Patient MRN (`PAT-YYYY-XXXXX`)
   - Encounter ID (`UUID`)
   - Invoice Number (`INV-YYYY-XXXXX`)
   - Payment Receipt Number (`REC-YYYY-XXXXX`)

---

## SECTION 26 — DEFECT REPORTING

Use this standardized defect template when logging any deviation or failure discovered during manual execution.

```markdown
## BUG-[XXX]: [Short Descriptive Title]

- **Feature**: [Feature name from Section 24]
- **Test ID**: [e.g. TC-ENC-01 / ATK-04]
- **Environment**: Local Staging (PostgreSQL 16, Go 1.25, React 19)
- **Actor / Persona**: [e.g. Doctor / Receptionist]
- **Tenant ID / Branch**: [e.g. curexal-clinic / main]
- **Severity**: [CRITICAL | HIGH | MEDIUM | LOW]
- **Security Impact**: [YES (PHI Leak / Escalation) | NO]
- **Reproducibility**: [ALWAYS | INTERMITTENT | ONCE]

### Steps to Reproduce
1. Log in as ...
2. Navigate to ...
3. Submit payload ...

### Expected Result
[Document exact expected behavior per test plan]

### Actual Result
[Document actual system behavior observed]

### Diagnostic Evidence
- **HTTP Status**: [e.g. 500 Internal Server Error]
- **API Response**:
  ```json
  [Paste raw JSON payload]
  ```
- **Database Query Result**:
  ```text
  [Paste SQL output]
  ```
- **Server Log Snippet**:
  ```text
  [Paste error log from apps/api]
  ```
- **Screenshot Link**: [evidence/BUG-XXX/screenshot.png]

### Suggested Investigation / Root Cause Hypothesis
[Developer hint or suspected code path]
```

---

## SECTION 27 — FINAL CERTIFICATION GATE

Before Curexal Clinic OS Phase 1 can be certified for production deployment or Phase 2 initiation, the QA Lead and Release Engineer must physically sign off on every item below.

### Mandatory Release Gate Criteria

- [ ] **Zero Critical Severity Defects** open in defect tracker.
- [ ] **Zero High Severity Defects** open in defect tracker.
- [ ] **Zero Cross-Tenant Data Leaks** verified empirically across multi-tenant test suites.
- [ ] **Zero Unauthorized Clinical Access** (Nurses/Receptionists/Cashiers strictly barred from prescribing).
- [ ] **Zero Unauthorized Financial Access** (Clinicians strictly barred from POS cash settlement).
- [ ] **Golden Path Test Suite Passes 100%** sequentially without manual database intervention.
- [ ] **Negative Security Path Passes 100%** with proper 401/403/404 rejections.
- [ ] **Authentication & Argon2id Verification Passes** across all personas.
- [ ] **Role-Based Access Control (RBAC) Passes** across all 27 clinic workspace permissions.
- [ ] **Branch & Organization Scoping Passes** without cross-branch context bleed.
- [ ] **Patient Registration & MPI Passes** with monotonic `PAT-YYYY-XXXXX` MRN formatting.
- [ ] **Appointment Scheduling & Constraint Checks Pass** (interval validation enforced).
- [ ] **Queue Orchestration & Acuity Sorting Passes** (`RED` > `YELLOW` > `GREEN`).
- [ ] **Triage Vitals & Observation Capture Passes** with boundary validation.
- [ ] **Encounter Lifecycle & Reconciled Schema Passes** (`appointment_id`, `encounter_channel`, `signed_at`, `signed_by`, `closed_at` physically confirmed).
- [ ] **Clinical SOAP Documentation Canvas Passes** draft and finalization stages.
- [ ] **ICD-10 Diagnostic Management Passes** primary and secondary assignments.
- [ ] **Medication e-Prescribing Passes** with full clinical safety parameters.
- [ ] **Billing & Cashier POS Passes** with automatic invoice generation and split-tender settlement.
- [ ] **Patient Self-Service Portal Passes** with strict PHI perimeter defense.
- [ ] **Regulatory Audit Logging Passes** with verifiable actor, tenant, and timestamp tracking.
- [ ] **Database Integrity Verification Passes** at Goose migration version `77`.
- [ ] **All Automated Backend Go Tests Pass** (`go test ./...` exits with code 0).
- [ ] **All Frontend TypeScript Compilations Pass** (`bun x tsc --noEmit` exits with 0 errors across all 3 web apps).
- [ ] **No Known Forensic Audit Blocker Remains Unresolved**.

---

### Final Release Determination

```
======================================================================
PHASE 1 CERTIFICATION STATUS: [   ] GO    /    [   ] NO-GO
======================================================================
```

**Sign-off Rationale & Release Engineer Statement**:  
*(To be completed by human QA Engineer following manual test execution)*

- **Lead QA Engineer Signature**: ___________________________ **Date**: ______________
- **Release Architect Signature**: ___________________________ **Date**: ______________
- **Security Auditor Signature**: ___________________________ **Date**: ______________
