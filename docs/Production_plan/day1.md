# CUREXAL CLINIC OS — DAY 1 PRODUCTION SPECIFICATION & IMPLEMENTATION PLAN

**Document**: `docs/Production_plan/day1.md`  
**Execution Day**: Day 1 of 30 (Thursday, September 3, 2026)  
**Target Milestone**: Foundation Lockdown & Security Hardening  
**Compliance Standard**: HIPAA Security Rule (§ 164.312), NDPR 2019, OWASP ASVS Level 2, Curexal Constitution v2.0  
**Strict Policy**: **NO SHORTCUTS.** Zero hardcoded mock credentials, zero inline bypass headers, zero unhandled errors.

---

## 1. OBJECTIVE & DELIVERABLES

Day 1 establishes the rock-solid, production-grade security and multi-tenant operational backbone required for the Saturday live launch:
1. **Feature #1**: Secure Authentication Engine (Argon2id, RS256/HS256 JWT, Refresh Rotation, Session Lifecycle).
2. **Feature #2**: Role-Based Access Control (Granular Permission Registry, Clinic Personas, Strict Middleware Enforcer).
3. **Feature #3**: Organization & Clinic Management (Multi-Branch Governance, Tenant Schema Isolation, Context Guards).
4. **Feature #24**: Production Audit Logging (Immutable, Append-Only Event Ledger with Actor Attribution).

---

## 2. TECHNICAL SPECIFICATION (DATA SCHEMAS & CONTRACTS)

### 2.1 Feature #1: Secure Authentication Specification

#### A. Cryptographic Standards
* **Password Hashing Algorithm**: `Argon2id`
  * Memory: `64 MB` ($65536 \text{ KiB}$)
  * Iterations: `3`
  * Parallelism: `2 threads`
  * Salt Length: `16 bytes` (cryptographically secure pseudorandom via `crypto/rand`)
  * Key Length: `32 bytes`
  * Format: Standard encoded string: `$argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>`
* **Access Token (JWT)**:
  * Algorithm: `HS256` / `RS256`
  * Lifespan: Exactly `15 minutes` (900 seconds)
  * Claims Payload:
    ```json
    {
      "sub": "usr_uuid_v4",
      "email": "doctor@clinic.com",
      "org_id": "org_uuid_v4",
      "branch_id": "branch_uuid_v4",
      "roles": ["doctor"],
      "perm_hash": "sha256_hash_of_active_permissions",
      "iat": 1788447600,
      "exp": 1788448500,
      "iss": "curexal-auth-engine",
      "aud": "curexal-clinic-os"
    }
    ```
* **Refresh Token & Family Rotation**:
  * Type: High-entropy opaque token (256-bit cryptographically secure string, stored as SHA-256 hash in DB).
  * Lifespan: Exactly `7 days` (sliding expiration).
  * Single-Use Rotation: When a refresh token is used, it is immediately invalidated, and a new refresh token is issued.
  * **Reuse Detection (Critical Security Rule)**: If an already-invalidated refresh token is submitted, the system flags token theft, invalidates the **entire token family**, and terminates all sessions for that user immediately.

#### B. Database Schema: `identity.sessions`
```sql
CREATE TABLE IF NOT EXISTS identity.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    active_branch_id UUID REFERENCES organization.facility_branches(id) ON DELETE SET NULL,
    refresh_token_hash VARCHAR(64) NOT NULL UNIQUE,
    token_family_id UUID NOT NULL DEFAULT gen_random_uuid(),
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT NOT NULL,
    device_fingerprint VARCHAR(64),
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    revocation_reason VARCHAR(255),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_active_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON identity.sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_org_id ON identity.sessions(organization_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token_hash ON identity.sessions(refresh_token_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_family_id ON identity.sessions(token_family_id);
```

#### C. API Contracts for Authentication & Branch Resolution
* `POST /api/v1/auth/login`:
  * **Request**:
    ```json
    {
      "email": "dr.adebayo@curexal.com",
      "password": "SecurePassword123!",
      "organizationSlug": "st-nicholas-clinic",
      "branchId": "branch_uuid" 
    }
    ```
    > [!IMPORTANT]
    > **Deterministic Branch Resolution (Clinical Attribution Rule)**:
    > The system **never silently guesses or defaults to Headquarters**, preventing clinical and legal charting misattributions.
    > 1. **Explicit Branch Supplied**: If `branchId` (or `branchCode`) is in the login payload, backend verifies the user holds an active staff membership in that specific branch. If unauthorized, returns `403 Forbidden` (`unauthorized_branch_access`).
    > 2. **Single Assigned Branch**: If omitted and the user is assigned to **exactly 1 branch**, the backend automatically resolves that branch as `activeBranch`.
    > 3. **Multiple Assigned Branches (Prompt Required)**: If omitted and the user is authorized across **multiple branches**, backend returns `requireBranchSelection: true` alongside the list of `assignedBranches` and a short-lived `selectionToken` (valid 5 mins). The UI presents a **"Select Your Active Clinic Location"** modal before clinical canvas access is granted.
    > 4. **Zero Assigned Branches**: If user has no operational branch assignments, access to clinical workspaces is blocked with `403 Forbidden` (`unassigned_facility_branch`).

  * **Response Scenario A: Single Branch / Explicit Branch Resolved (200 OK)**:
    ```json
    {
      "data": {
        "accessToken": "eyJhbGciOiJIUzI1NiIs...",
        "tokenType": "Bearer",
        "expiresIn": 900,
        "user": {
          "id": "usr_uuid",
          "email": "dr.adebayo@curexal.com",
          "fullName": "Dr. Musbau Adebayo",
          "avatarUrl": null
        },
        "organization": {
          "id": "org_uuid",
          "name": "St. Nicholas Clinic",
          "slug": "st-nicholas-clinic",
          "plan": "optimize"
        },
        "activeBranch": {
          "id": "branch_uuid_1",
          "name": "Victoria Island Main Campus",
          "code": "VI-01",
          "facilityType": "clinic"
        },
        "assignedBranches": [
          {
            "id": "branch_uuid_1",
            "name": "Victoria Island Main Campus",
            "code": "VI-01",
            "role": "doctor"
          },
          {
            "id": "branch_uuid_2",
            "name": "Ikeja Outpatient Center",
            "code": "IK-02",
            "role": "doctor"
          }
        ],
        "roles": ["doctor"],
        "permissions": [
          "workspace:clinical:read",
          "workspace:clinical:write",
          "workspace:patient:read"
        ]
      },
      "meta": { "timestamp": "2026-09-03T17:00:00Z" }
    }
    ```
  * Sets HttpOnly, Secure, SameSite=Strict cookie: `curexal_refresh_token`.

  * **Response Scenario B: Multiple Branches Assigned & No branchId Supplied (200 OK — Branch Selection Required)**:
    ```json
    {
      "data": {
        "requireBranchSelection": true,
        "selectionToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "assignedBranches": [
          {
            "id": "branch_uuid_1",
            "name": "Victoria Island Main Campus",
            "code": "VI-01",
            "role": "doctor"
          },
          {
            "id": "branch_uuid_2",
            "name": "Ikeja Outpatient Center",
            "code": "IK-02",
            "role": "doctor"
          }
        ]
      },
      "meta": { "timestamp": "2026-09-03T17:00:00Z" }
    }
    ```

* `POST /api/v1/auth/select-branch` (Complete Login after Branch Selection Modal):
  * **Request**:
    ```json
    {
      "selectionToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "branchId": "branch_uuid_2"
    }
    ```
  * **Business Logic**:
    1. Validates `selectionToken` signature and expiry.
    2. Verifies user has active membership in `branchId`.
    3. Creates session in `identity.sessions` with `active_branch_id = branchId`.
    4. Issues full `accessToken` and sets refresh cookie.
  * **Response (200 OK)**: Returns full session payload as in Scenario A.

* `POST /api/v1/auth/switch-branch` (Dynamic Shift / Multi-Branch Rotation):
  * **Request**:
    ```json
    {
      "branchId": "branch_uuid_2"
    }
    ```
  * **Business Logic**:
    1. Validates that caller has an active, non-suspended staff membership in `branch_uuid_2`.
    2. Updates `identity.sessions.active_branch_id = branch_uuid_2`.
    3. Re-issues fresh Access Token with new `branch_id` embedded in claims.
    4. Records `auth.branch_switched` event in `audit.audit_events`.
  * **Response (200 OK)**:
    ```json
    {
      "data": {
        "accessToken": "eyJhbGciOiJIUzI1NiIs...",
        "activeBranch": {
          "id": "branch_uuid_2",
          "name": "Ikeja Outpatient Center",
          "code": "IK-02",
          "facilityType": "clinic"
        }
      },
      "meta": { "timestamp": "2026-09-03T17:05:00Z" }
    }
    ```

* `POST /api/v1/auth/refresh`:
  * Reads `curexal_refresh_token` cookie.
  * Validates session table, verifies not revoked, rotates token hash.
  * Preserves `active_branch_id` in new access token.
  * Returns new `accessToken` and sets updated cookie.

* `POST /api/v1/auth/logout`:
  * Invalidates the active session in `identity.sessions`.
  * Clears the refresh token cookie.

---

### 2.2 Feature #2: Role-Based Access Control (RBAC) Specification

#### A. The Six Canonical Clinic Roles
1. **`owner`** (Clinic Owner / Medical Director): Full legal, organizational, and financial authority.
2. **`org_admin`** (Practice Manager / Clinic Administrator): Manages branches, staff invites, catalogs, and schedules.
3. **`doctor`** (Attending Medical Doctor / Specialist): Consultation canvas, SOAP notes, ICD-10 coding, e-prescriptions, diagnostic orders.
4. **`nurse`** (Triage Nurse / Clinical Assistant): Patient queue check-in, triage intake, vital signs, allergy tagging.
5. **`receptionist`** (Front Desk Officer): Master Patient Index (MPI) search, walk-in registration, appointment booking.
6. **`cashier`** (Billing & Accounts Officer): Service fee settlement, POS card/cash/transfer collections, receipt generation, daily shift reconciliation.


#### B. Permission Registry & Mapping Matrix
Every action must map to an explicit permission code in `"authorization".permissions`. Hardcoded role string checks (`if role == "doctor"`) are strictly forbidden.

##### 1. Clinic Workspace Permissions
| Permission Code | Category | Module | Granted To | Real-World Operational Scope |
| :--- | :--- | :--- | :--- | :--- |
| `organization:view` | core | organization | `owner`, `org_admin` | View clinic profile, branches, and analytics |
| `organization:manage` | core | organization | `owner`, `org_admin` | Edit clinic settings, business hours, and VAT rules |
| `organization:branch:manage` | core | organization | `owner`, `org_admin` | Provision and manage clinic branch locations |
| `users:read` | core | organization | `owner`, `org_admin`, `doctor` | View staff roster and provider availability |
| `users:write` | core | organization | `owner`, `org_admin` | Invite, assign roles, and deactivate staff |
| `audit:read` | core | audit | `owner`, `org_admin` | Inspect immutable audit and compliance trail |
| `workspace:patient:create` | clinical | patient | `owner`, `org_admin`, `doctor`, `nurse`, `receptionist` | Register new patients and assign MRN |
| `workspace:patient:read` | clinical | patient | `owner`, `org_admin`, `doctor`, `nurse`, `receptionist`, `cashier` | Search Master Patient Index (MPI) |
| `workspace:patient:update` | clinical | patient | `owner`, `org_admin`, `doctor`, `nurse`, `receptionist` | Update patient demographics and contacts |
| `workspace:appointment:read` | operations | customer_care | `owner`, `org_admin`, `doctor`, `nurse`, `receptionist` | View doctor calendars and booked appointments |
| `workspace:appointment:write`| operations | customer_care | `owner`, `org_admin`, `doctor`, `nurse`, `receptionist` | Book, reschedule, and cancel appointments |
| `workspace:queue:manage` | operations | customer_care | `owner`, `org_admin`, `doctor`, `nurse`, `receptionist` | Patient check-in and queue status management |
| `workspace:triage:create` | clinical | clinical | `owner`, `doctor`, `nurse` | Record vitals (BP, Pulse, Temp, SpO2, Weight) |
| `workspace:triage:read` | clinical | clinical | `owner`, `doctor`, `nurse` | View historical vitals and trend charts |
| `workspace:clinical:read` | clinical | clinical | `owner`, `doctor`, `nurse` | View past visit history, allergies, and care plans |
| `workspace:clinical:write` | clinical | clinical | `owner`, `doctor` | Document SOAP notes, exam findings, and ICD-10 |
| `workspace:clinical:sign` | clinical | clinical | `owner`, `doctor` | Digitally sign and permanently lock encounter |
| `workspace:prescription:write`| clinical | clinical | `owner`, `doctor` | Prescribe medications (drug, dose, duration) |
| `workspace:prescription:read` | clinical | clinical | `owner`, `doctor`, `nurse`, `cashier` | Cashier bills drugs; nurse administers meds |
| `workspace:diagnostic:order` | clinical | clinical | `owner`, `doctor` | Order laboratory tests and imaging procedures |
| `workspace:diagnostic:read` | clinical | clinical | `owner`, `doctor`, `nurse`, `cashier` | Cashier bills tests; clinicians review results |
| `workspace:document:upload` | clinical | documents | `owner`, `doctor`, `nurse`, `receptionist` | Upload referral letters, paper scans, and IDs |
| `workspace:document:read` | clinical | documents | `owner`, `doctor`, `nurse`, `receptionist` | Inspect uploaded documents and test PDFs |
| `workspace:billing:read` | financial | billing | `owner`, `org_admin`, `cashier` | Access invoice ledgers and balance statements |
| `workspace:billing:charge` | financial | billing | `owner`, `org_admin`, `cashier` | Generate itemized invoices from consults/tests |
| `workspace:billing:refund` | financial | billing | `owner`, `org_admin` | Authorize payment reversals and credit notes |
| `workspace:pos:settle` | financial | billing | `owner`, `cashier` | Settle invoices via Cash, POS Card, or Transfer |
| `workspace:pos:shift_close` | financial | billing | `owner`, `cashier` | Reconcile cash drawer and close cashier shift |

##### 2. Platform Super-Admin Permissions
| Permission Code | Category | Module | Granted To | Real-World Operational Scope |
| :--- | :--- | :--- | :--- | :--- |
| `platform:admin` | platform | platform | `super_admin` | Root authority over platform, schemas, and keys |
| `platform:view` | platform | platform | `super_admin`, `super_support_agent`, `super_sales_staff` | Access platform dashboard metrics & clinic counts |
| `platform:manage` | platform | platform | `super_admin` | Global settings, feature flags, and tenant lifecycle |
| `platform:impersonate` | platform | platform | `super_admin`, `super_support_agent` | Break-glass support access to debug clinic issues |
| `platform:catalogs:manage`| platform | platform | `super_admin` | Maintain global ICD-10 registry & standard catalogs |
| `platform:pricing:manage` | platform | platform | `super_admin` | Configure subscription rates and gateway vault keys |
| `demo:manage` | platform | platform | `super_admin`, `super_sales_staff` | Process inbound demo requests & provision pilots |


#### C. Hardened Production Enforcement Middleware (`internal/shared/middleware/enforcer.go`)
```go
type AuditLogger interface {
    LogSecurityViolation(ctx context.Context, principal *Principal, permCode string, ip string, userAgent string)
}

type MembershipValidator interface {
    IsMembershipActive(ctx context.Context, userID, orgID string) (bool, error)
}

type RBACEnforcer struct {
    auditService AuditLogger
    membership   MembershipValidator
}

func NewRBACEnforcer(audit AuditLogger, membership MembershipValidator) *RBACEnforcer {
    return &RBACEnforcer{
        auditService: audit,
        membership:   membership,
    }
}

// RequirePermission enforces exact granular permissions with tenant isolation and audit logging.
func (e *RBACEnforcer) RequirePermission(permCode string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            principal := GetPrincipal(c)
            if principal == nil {
                return response.UnauthorizedEcho(c, "Authentication credentials missing or invalid")
            }

            // 1. Strict Tenant Active Verification (Prevents revoked staff from using unexpired tokens)
            ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
            defer cancel()

            isActive, err := e.membership.IsMembershipActive(ctx, principal.UserID, principal.OrganizationID)
            if err != nil || !isActive {
                return response.ForbiddenEcho(c, "Organization membership is inactive or suspended")
            }

            // 2. Permission Check (O(1) fast lookup - NO un-audited super-admin backdoors)
            if !principal.HasPermission(permCode) {
                ip := c.RealIP()
                ua := c.Request().UserAgent()
                go e.auditService.LogSecurityViolation(context.Background(), principal, permCode, ip, ua)

                return response.ForbiddenEcho(c, fmt.Sprintf("Access denied: missing required permission '%s'", permCode))
            }

            // 3. Inject verified tenant and branch context into Echo Context
            c.Set("tenant_id", principal.OrganizationID)
            c.Set("branch_id", principal.ActiveBranchID)

            return next(c)
        }
    }
}
```

---


### 2.3 Feature #3: Organization & Clinic Management Specification (Healthcare Production-Grade)

#### A. Database Schema: `organization.facility_branches`
```sql
CREATE TABLE IF NOT EXISTS organization.facility_branches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE RESTRICT,
    facility_type_id UUID NOT NULL REFERENCES platform.facility_types(id),
    code VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    lga VARCHAR(100),
    country VARCHAR(100) NOT NULL DEFAULT 'Nigeria',
    operating_hours JSONB NOT NULL DEFAULT '{"mon": {"open": "08:00", "close": "18:00"}}'::jsonb,
    theme_branding JSONB NOT NULL DEFAULT '{}'::jsonb,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'SUSPENDED')),
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    
    CONSTRAINT uk_facility_org_code UNIQUE (organization_id, code),
    CONSTRAINT uk_facility_branches_org_slug UNIQUE (organization_id, slug)
);

-- Performance & Isolation Indexes
CREATE INDEX IF NOT EXISTS idx_facility_branches_org_id ON organization.facility_branches(organization_id);
CREATE INDEX IF NOT EXISTS idx_facility_branches_type ON organization.facility_branches(facility_type_id);
CREATE INDEX IF NOT EXISTS idx_facility_branches_status ON organization.facility_branches(status);
```

#### B. Tenant & Facility Isolation Invariants
1. **Tenant Organization Boundary**:
   * Resolved from verified JWT `org_id` claim and cross-verified against `organization.organization_memberships` (`status = 'ACTIVE'`).
   * Cross-tenant data access is strictly blocked at both the API middleware and PostgreSQL Row-Level Security (RLS) layers.
2. **Branch Context & Staff Assignment Verification**:
   * Facility-scoped clinical and operational requests require an explicit `branch_id` context.
   * **Assignment Guard**: Unless the caller has an organization-level governance role (`owner`, `org_admin`), the user's `membership_id` must have an active assignment in `organization.membership_branches` matching the requested `branch_id`. Unauthorized cross-branch access attempts are logged as security violations (`rbac.permission_denied`).
3. **Data Scoping Separation**:
   * **Organization Scope (Shared)**: Master Patient Index (MPI demographics), global pricing catalog, staff directory, and enterprise billing.
   * **Branch Scope (Strictly Isolated)**: Clinical encounters, vital signs, lab orders/specimens, prescription dispensing, inventory stock levels, physical bed assignments, and financial tills.
4. **Zero Implicit Fallback Policy**:
   * The platform strictly forbids implicit defaulting or fallback to an assumed branch.
   * Requests lacking an active branch must fail fast (`428 Precondition Required` or `400 Bad Request`) rather than silently routing clinical or financial actions to an arbitrary facility.

---


### 2.4 Feature #24: Audit Logging Specification (Healthcare Production-Grade)

#### A. Compliance & Regulatory Standards
* **HIPAA Security Rule § 164.312(b)**: Automated hardware/software audit controls recording all activity in systems containing ePHI.
* **HIPAA Privacy Rule § 164.528**: Mandatory Accounting of Disclosures (instant queryability of every user who accessed or exported a given patient's records).
* **ASTM E2147-18**: Standard Specification for Audit and Disclosure Logs for Use in Health Information Systems (actor persona, patient identity, clinical justification, access type).
* **NDPA 2023 / NDPR 2019**: Non-repudiation, tamper-evidence, and forensic data subject attribution.
* **Retention Invariant**: Minimum 6-year append-only ledger lifespan.

#### B. Database Schema: `audit.audit_events`
```sql
CREATE TABLE IF NOT EXISTS audit.audit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 1. Multi-Tenant & Facility Scoping (Strict Context)
    organization_id UUID NOT NULL,
    facility_branch_id UUID,
    
    -- 2. Patient / Subject of Care (Mandatory for HIPAA Accounting of Disclosures)
    patient_id UUID,
    
    -- 3. Forensic Actor Snapshot (Permanent snapshot - NEVER erased or set NULL on user deletion)
    actor_id UUID,
    actor_name VARCHAR(255) NOT NULL,
    actor_email VARCHAR(255) NOT NULL,
    actor_role VARCHAR(100) NOT NULL,
    
    -- 4. Action, Resource & Event Taxonomy
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id VARCHAR(255),
    resource_name VARCHAR(255),
    event_category VARCHAR(50) NOT NULL, -- 'AUTH', 'CLINICAL_ACCESS', 'DATA_MUTATION', 'SECURITY', 'DISCLOSURE'
    severity VARCHAR(20) NOT NULL DEFAULT 'INFO', -- 'INFO', 'WARN', 'CRITICAL', 'ALERT'
    status VARCHAR(50) NOT NULL DEFAULT 'SUCCESS', -- 'SUCCESS', 'FAILURE', 'DENIED'
    
    -- 5. Clinical Justification & Emergency Break-The-Glass Tracking
    reason TEXT,
    is_break_glass BOOLEAN NOT NULL DEFAULT FALSE,
    approval_reference VARCHAR(255),
    
    -- 6. Network & Request Forensics
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT NOT NULL,
    device VARCHAR(100),
    operating_system VARCHAR(100),
    browser VARCHAR(100),
    hostname VARCHAR(255),
    request_id VARCHAR(100),
    session_id VARCHAR(100),
    trace_id VARCHAR(100),
    
    -- 7. State Diff & Payloads (Sanitized - Sensitive PHI masked; schema diffs only)
    before_state JSONB,
    after_state JSONB,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    
    -- 8. Cryptographic Tamper-Evidence & Non-Repudiation
    prev_record_hash CHAR(64),
    record_hash CHAR(64),
    digital_signature TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- STRICT IMMUTABILITY: Prevent any update or delete operations on the audit ledger
CREATE OR REPLACE FUNCTION audit.prevent_tampering()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Security violation: audit.audit_events records are append-only and cannot be updated or deleted.';
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_prevent_audit_tampering ON audit.audit_events;
CREATE TRIGGER trg_prevent_audit_tampering
BEFORE UPDATE OR DELETE ON audit.audit_events
FOR EACH ROW EXECUTE FUNCTION audit.prevent_tampering();

-- High-performance query indexes for compliance reporting
CREATE INDEX IF NOT EXISTS idx_audit_patient_id ON audit.audit_events(patient_id) WHERE patient_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_audit_org_id ON audit.audit_events(organization_id);
CREATE INDEX IF NOT EXISTS idx_audit_branch_id ON audit.audit_events(facility_branch_id) WHERE facility_branch_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_audit_actor_id ON audit.audit_events(actor_id);
CREATE INDEX IF NOT EXISTS idx_audit_action ON audit.audit_events(action);
CREATE INDEX IF NOT EXISTS idx_audit_category ON audit.audit_events(event_category);
CREATE INDEX IF NOT EXISTS idx_audit_created_at ON audit.audit_events(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_break_glass ON audit.audit_events(is_break_glass) WHERE is_break_glass = TRUE;
```

#### C. Canonical Healthcare Audit Action Catalog
1. **Authentication & Session Lifecycle**:
   * `auth.login.success` / `auth.login.failure`
   * `auth.mfa.challenged` / `auth.mfa.verified`
   * `auth.session.revoked` / `auth.session.expired`
   * `auth.token.refresh_rotated` / `auth.token.theft_detected`
   * `auth.password.reset_requested` / `auth.password.changed`

2. **Clinical Read & Access Auditing (HIPAA Snooping Prevention)**:
   * `patient.chart.viewed` (Records access to a patient record)
   * `clinical.encounter.viewed` (Doctor opens clinical notes)
   * `laboratory.result.viewed` (Viewing lab/pathology results)
   * `pharmacy.prescription.viewed` (Viewing medication regimens)
   * `imaging.study.viewed` (Accessing radiology DICOM / reports)

3. **Emergency & Break-the-Glass Overrides**:
   * `clinical.break_glass.invoked` (Unassigned provider emergency chart access with mandatory `reason`)
   * `rbac.permission_override` (Privileged administrative bypass)

4. **Clinical Data Mutations (Write / Amend / Sign)**:
   * `patient.registered` / `patient.demographics_updated`
   * `encounter.started` / `encounter.triage_completed`
   * `encounter.clinical_note_signed` (Digital clinical signature locked)
   * `prescription.ordered` / `prescription.dispensed` / `prescription.cancelled`
   * `laboratory.order_placed` / `laboratory.result_verified`
   * `patient.admitted` / `patient.transferred` / `patient.discharged`

5. **Data Exfiltration, Disclosures & Exports**:
   * `phi.record.exported` (PDF / CSV export of patient records)
   * `phi.report.printed` (Hardcopy print action)
   * `phi.batch_downloaded` (Bulk export alerting)
   * `disclosure.accounting_generated` (Patient-requested access history disclosure)

6. **Administrative & Governance**:
   * `staff.membership_created` / `staff.membership_deactivated`
   * `staff.role_assigned` / `staff.permissions_modified`
   * `organization.settings_updated`
   * `branch.created` / `branch.deactivated`

---

## 3. STRICT DAY 1 IMPLEMENTATION PLAN (STEP-BY-STEP)

```text
[STEP 1: DB MIGRATION] ──► [STEP 2: ENFORCER REFACTOR] ──► [STEP 3: AUTH ROTATION] ──► [STEP 4: E2E TESTS]
 • Run Goose Migrations     • Remove string checks       • Session Hash Table        • Test 401, 403, 200
 • Apply Trigger Immutability • Wire RequirePermission()  • Rotate refresh token      • Zero Regressions
```

### Step 1: Database Migration & Schema Hardening
1. Verify Goose migration status:
   ```powershell
   goose -dir apps/api/database/platform/migrations postgres "$DATABASE_URL" status
   ```
2. Verify table existence and constraints:
   * `identity.sessions`
   * `audit.audit_events` with `trg_prevent_audit_tampering` trigger.
   * `"authorization".role_permissions` mapping all 6 canonical clinic roles.

### Step 2: Purge Hardcoded Role Strings in Middleware & Handlers
1. Inspect `internal/shared/middleware/middleware.go` and `internal/platform/auth/`.
2. Replace any `if role == "doctor"` with `RequirePermission("workspace:clinical:write")`.
3. Eliminate any bypass logic that trusts unverified `X-User-Role` headers outside testing environments.

### Step 3: Implement Session Hash & Family Rotation in `identity` Module
1. In `internal/modules/identity/handler/auth.go`, ensure:
   * Password verification uses `argon2id.ComparePasswordAndHash(hash, password)`.
   * On successful login, a session record is inserted into `identity.sessions`.
   * On `/refresh`, token reuse triggers complete family revocation.
   * On `/logout`, session record is marked `is_revoked = true`.

### Step 4: Wire Unified Audit Logging Interceptor
1. Ensure all authenticated mutations trigger an audit entry in `audit.audit_events`.
2. Redact sensitive values from JSONB payload (passwords, complete credit card numbers).

---

## 4. VERIFICATION & TEST CRITERIA (THE "NO SHORTCUTS" GATE)

Every test below must execute and pass with zero failures before Day 1 is signed off:

### Automated Test Suite
```powershell
# 1. Backend Integration Tests
cd apps/api
go test -v ./internal/modules/identity/...
go test -v ./internal/modules/organization/...
go test -v ./internal/testing/...

# 2. Frontend Strict Typecheck
cd ../..
bun x tsc --noEmit
```

### Security & Multi-Branch Verification Checklist (Manual Sign-Off)
| Test Case | Method / Scenario | Expected Result | Verified? |
| :--- | :--- | :--- | :---: |
| **Tampered JWT** | Modify payload signature of access token and send request | `401 Unauthorized` | [ ] |
| **Expired Token** | Send request with access token older than 15 mins | `401 Unauthorized` (Token Expired) | [ ] |
| **Token Theft Simulation** | Submit an old, already-rotated refresh token | All sessions revoked; `401` returned | [ ] |
| **Privilege Escalation** | Log in as Cashier; attempt `POST /api/v1/clinical/encounters` | `403 Forbidden` with audit log | [ ] |
| **Unauthorized Branch Access** | Doctor logs in with `branchId` for a facility where they have no membership | `403 Forbidden` (Not authorized in branch) | [ ] |
| **Multi-Branch Selection Prompt** | Doctor assigned to 2 branches logs in without `branchId` | `200 OK` with `requireBranchSelection: true`; completes via `/select-branch` | [ ] |
| **Dynamic Branch Switch** | Doctor invokes `POST /api/v1/auth/switch-branch` with valid Branch B ID | New JWT issued with Branch B in claims; audit event logged | [ ] |
| **Tenant Cross-Bleed** | Send valid token from Clinic A with `X-Tenant-ID: Clinic B` | `403 Forbidden` (Tenant mismatch) | [ ] |
| **Audit Immutability** | Attempt `DELETE FROM audit.audit_events WHERE id = ...` | Database error: trigger exception | [ ] |

---

## 5. DAY 1 PRODUCTION DEPLOYMENT PROCEDURE

At 18:00 UTC, execute the release protocol:
```powershell
# Step 1: Pre-flight check
bun x tsc --noEmit
cd apps/api; go test -v ./internal/testing/...

# Step 2: Push database schema updates to production
goose -dir database/platform/migrations up

# Step 3: Tag release
git add .
git commit -m "release(day-1): Foundation lockdown - Auth, RBAC, Clinic Governance & Audit Ledger"
git tag -a "deploy-2026-09-03" -m "Day 1 Production Release: Features #1, #2, #3, #24"
git push origin main --tags
```
