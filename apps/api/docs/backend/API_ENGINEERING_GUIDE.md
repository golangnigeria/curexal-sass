# CUREXAL API ENGINEERING & AUDIT GUIDE v2.0
**Canonical Handbook for Testing, Auditing, Reviewing, and Validating Curexal Kernel APIs**

---

## Executive Summary & Engineering Charter

You are reading the official API Engineering Guide for the **Curexal Platform Kernel**.

Curexal is an enterprise-grade healthcare operating system responsible for running hospitals, diagnostic laboratories, radiology centers, clinical networks, and pharmacies. In healthcare software, API failures, unhandled concurrency, authorization bypasses, or data corruptions do not merely cause UI bugs—they threaten patient safety, violate clinical regulatory standards, and compromise sensitive medical data.

### The Architectural Axioms
1. **The Backend IS the Product**: The Go backend kernel owns 100% of the business logic, state machines, permission matrices, multi-tenant isolation, and data governance.
2. **The Frontend IS a Client**: Frontend applications (`web-platform`, `web-organization`, `web-workspace`, `web-patient`, `web-public`, mobile/desktop apps) are presentation layers. Frontends compute **zero** business rules, zero permissions, and zero role logic.
3. **Database IS Law**: PostgreSQL schemas are the ultimate source of truth. Every business rule originates from database state.
4. **OpenAPI IS Contract**: Every endpoint, request payload, response schema, and error envelope must be explicitly defined in OpenAPI 3.0+ specifications.

```
                 PostgreSQL
           (Single Source of Truth)
                     │
                     ▼
        Go Domain Application Services
                     │
                     ▼
            REST/OpenAPI Contract
                     │
                     ▼
          Generated Type-safe SDK
                     │
                     ▼
      Web / Mobile / Desktop Clients
```

---

## 1. API Philosophy

### Why APIs Exist in Enterprise Healthcare
APIs exist to establish a secure, deterministic, and audited boundary around enterprise state. An API is not a convenience wrapper over a database; it is a transactional contract that guarantees data integrity, tenant isolation, and regulatory compliance regardless of which client initiates the request.

### Why Frontends Must NEVER Access the Database Directly
Direct database access from frontend clients introduces catastrophic vulnerabilities:
- **Corrupted State**: UI code can bypass validation rules, creating orphan records (e.g., lab accessions without patient MRNs).
- **Security Leaks**: Database credentials exposed in browser client code allow full database dumping.
- **Bypassed Audit Trails**: Healthcare compliance requires strict audit logging for every read/write of Protected Health Information (PHI). Direct database queries bypass backend audit middleware.
- **Architectural Drift**: Business logic becomes fragmented across web apps, mobile apps, and third-party integrations.

### Request Flow Architecture
Every HTTP request follows a strict unidirectional data flow through five isolated layers:

```
Browser / Client (HTTP Request)
       │
       ▼
Echo Router & Middleware Pipeline (Auth, Tenant Resolution, Rate Limiting, CORS)
       │
       ▼
HTTP Handler (Parses DTO, validates payload syntax, invokes service)
       │
       ▼
Application Service (Orchestrates domain logic, executes ACID transaction, checks business rules)
       │
       ▼
Repository (Performs raw SQL persistence against PostgreSQL)
       │
       ▼
PostgreSQL (Enforces schema constraints, search_path isolation, indexes)
       │
       ▼
Response Builder (Formats JSON envelope: { data, meta, links, errors })
```

#### Layer Responsibilities
- **Client**: Renders UI based on backend JSON payloads. Never calculates roles, feature availability, or prices.
- **Router & Middleware**: Intercepts HTTP traffic, authenticates tokens/cookies, resolves active tenant schema via `search_path`, enforces CORS, and tracks metrics.
- **Handler**: Adapts HTTP to Go structs. Parses JSON payloads, validates input syntax (e.g. email format, string lengths), and calls Application Services. **Zero business logic allowed here.**
- **Application Service**: Executes business use cases. Coordinates multiple repositories, manages database transactions (`tx.Begin`), publishes audit events, and evaluates domain rules.
- **Repository**: Performs database persistence. Executes SQL queries (`pgx`), scans rows into domain structs, and returns raw data. **Zero HTTP or business logic allowed here.**
- **PostgreSQL**: Stores data cleanly in context-owned schemas (`identity`, `organization`, `workspace`, `patient`, `laboratory`, `clinical`, `billing`, etc.).

---

## 2. API Lifecycle

Below is the execution sequence for a single request (`POST /api/v1/auth/sign-in`):

```
1. HTTP POST Request Received (/api/v1/auth/sign-in)
   │
2. Router Matching (Echo framework matches route pattern)
   │
3. Global Middleware Execution (RequestID, Logger, Recover, CORS)
   │
4. Auth & Context Middleware (Extracts session cookie / Bearer token if present)
   │
5. Request Body Binding & Syntax Validation (Struct tag validation: email, password length)
   │
6. Handler Execution (Extracts DTO, invokes AuthService.SignIn)
   │
7. Application Service Execution (AuthService opens tx, queries identity.users & identity.credentials)
   │
8. Password Hash Verification (bcrypt / Argon2 verification)
   │
9. Session / Token Generation (Creates session record in DB, generates HTTP-only cookie)
   │
10. Audit Log Dispatch (Asynchronously logs identity.login_success event)
    │
11. Transaction Commit (ACID tx committed to PostgreSQL)
    │
12. Response Formatting (Handler wraps user DTO in { "data": {...}, "meta": {...} })
    │
13. HTTP Response Transmission (200 OK + Set-Cookie header)
```

---

## 3. The 25 Questions Every Endpoint Must Answer

Before any endpoint is considered production-ready, engineers must answer all 25 questions:

### 1. Why does this endpoint exist?
It satisfies a specific enterprise capability (e.g., allowing lab technicians to authorize accession test results).

### 2. Which business capability does it serve?
Maps directly to a medical, administrative, or operational capability (e.g., Laboratory LIS Worksheet Authorization).

### 3. Which bounded context owns it?
Belongs to exactly one domain module (e.g., `laboratory`, `identity`, `patient`, `clinical`, `billing`).

### 4. Which database tables does it read?
Enumerates all tables read during execution (e.g., `laboratory.results`, `laboratory.accessions`, `patient.patients`).

### 5. Which database tables does it write?
Enumerates all tables modified (e.g., `laboratory.results`, `audit.audit_events`).

### 6. Does it create records?
Identifies primary entities created (e.g., creates `laboratory.results` row).

### 7. Does it update records?
Identifies entities updated (e.g., updates `status = 'authorized'` and `authorized_at = NOW()`).

### 8. Does it delete records?
Confirms whether soft-delete (`deleted_at`) or hard-delete is performed. Hard-deletes are forbidden for medical records.

### 9. Which transaction boundary is used?
Explicitly defines the PostgreSQL `pgx.Tx` transaction scope across all modified tables.

### 10. Is rollback supported?
Verifies that `defer tx.Rollback(ctx)` guarantees partial write rollback if any query or event fails.

### 11. Which permissions are required?
Lists mandatory authorization scopes (e.g., `workspace:result:authorize`).

### 12. Which roles can call it?
Lists authorized system roles (e.g., `lab_technician`, `pathologist`, `branch_admin`).

### 13. Which middleware executes first?
Traces middleware order: `Logger` -> `CORS` -> `Authenticate` -> `TenantResolver` -> `RequirePermission`.

### 14. Which validation rules run?
Details input validation (e.g., UUID format, non-empty result values, reference range bounds).

### 15. Which errors can occur?
Catalogues expected error codes (e.g., `RESULT_NOT_FOUND`, `ALREADY_AUTHORIZED`, `UNAUTHORIZED_ROLE`).

### 16. Which HTTP status codes are returned?
Defines exact return codes: `200 OK`, `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `409 Conflict`, `422 Unprocessable Entity`, `500 Internal Error`.

### 17. What response contract is expected?
Matches the standard JSON envelope: `{ "data": {}, "meta": {}, "links": {}, "errors": [] }`.

### 18. Does it require authentication?
Specifies if the endpoint is public (e.g., `POST /auth/sign-in`) or requires an active session (`GET /users/me`).

### 19. Does it require tenant context?
Confirms if schema-per-tenant (`SET search_path = tenant_slug`) must be isolated.

### 20. Which audit logs are created?
Identifies compliance audit rows inserted into `audit.audit_events` (Actor, Action, Resource, Payload, IP).

### 21. Does it publish domain events?
Lists domain events emitted for async processing (e.g., `LabResultAuthorizedEvent`).

### 22. Which OpenAPI definition documents it?
Verifies exact match in [static/openapi.json](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/backend/static/openapi.json).

### 23. Which frontend client consumes it?
Identifies consuming applications (`web-workspace`, `web-patient`, etc.).

### 24. How should it be tested?
Defines required automated unit, repository, service, integration, and manual Postman tests.

### 25. How do we know it is production ready?
Checks off the Definition of Done (Passed `go test`, `go vet`, OpenAPI validation, zero hardcoded values).

---

## 4. Manual Postman Testing Guide

### 1. Setting Up Postman Collections & Environments
1. Open **Postman Desktop**.
2. Create a new collection named **`Curexal Production Kernel`**.
3. Create an Environment named **`Local Backend`** with variables:
   - `baseUrl` = `http://localhost:8080/api/v1`
   - `sessionToken` = *(leave blank; auto-filled by tests)*
   - `activeTenantId` = `default-tenant-id`

### 2. Automated Session & Cookie Management
When calling `POST {{baseUrl}}/auth/sign-in`, add this script to the **Tests** tab of the request:

```javascript
// Automatically extract session & CSRF tokens into Postman Environment
const response = pm.response.json();
if (response.data && response.data.token) {
    pm.environment.set("sessionToken", response.data.token);
}

// Extract CSRF token from headers if present
const csrfHeader = pm.response.headers.get("X-CSRF-Token");
if (csrfHeader) {
    pm.environment.set("csrfToken", csrfHeader);
}
```

### 3. Collection Authorization Setup
Select the collection root -> **Authorization** tab:
- **Type**: `Bearer Token`
- **Token**: `{{sessionToken}}`

Every request inside the collection inherits authentication automatically.

---

## 5. Response Inspection Guide

Engineers must inspect every API response against standard HTTP status codes:

| HTTP Status Code | Meaning | Required Envelope Shape | When to Expect |
| :--- | :--- | :--- | :--- |
| **`200 OK`** | Request succeeded | `{ "data": {...}, "meta": {} }` | GET, PUT, PATCH read/update operations |
| **`201 Created`** | Entity created | `{ "data": {...}, "meta": {} }` | POST resource creations (e.g. Register Patient) |
| **`204 No Content`** | Succeeded with no body | *Empty Body* | DELETE operations or context switches |
| **`400 Bad Request`** | Malformed JSON / Syntax error | `{ "errors": [{"code": "INVALID_JSON"}] }` | Unparseable body or invalid data type |
| **`401 Unauthorized`** | Missing or expired auth token | `{ "errors": [{"code": "UNAUTHORIZED"}] }` | Expired session cookie or missing Bearer header |
| **`403 Forbidden`** | Insufficient permissions | `{ "errors": [{"code": "FORBIDDEN"}] }` | User lacks required permission scope |
| **`404 Not Found`** | Resource does not exist | `{ "errors": [{"code": "NOT_FOUND"}] }` | Invalid UUID or non-existent record |
| **`409 Conflict`** | Duplicate key / Unique constraint | `{ "errors": [{"code": "DUPLICATE_SLUG"}]}`| Email or tenant slug already registered |
| **`422 Unprocessable`**| Business rule validation failed | `{ "errors": [{"code": "INVALID_MRN"}] }` | Syntax valid but fails domain validation |
| **`429 Too Many`** | Rate limit exceeded | `{ "errors": [{"code": "RATE_LIMITED"}] }` | Exceeded API rate limits |
| **`500 Internal Error`**| Unexpected backend panic / DB fail | `{ "errors": [{"code": "INTERNAL_ERROR"}]}`| **UNACCEPTABLE IN PRODUCTION**. Investigate log. |

---

## 6. Database Verification Guide

To verify an API endpoint's database operations:

### 1. Pre-Execution Database Snapshot
Run SQL query in `psql` before sending Postman request:
```sql
-- Snapshot target table row count & state
SELECT count(*), max(created_at) FROM laboratory.results;
SELECT * FROM audit.audit_events ORDER BY created_at DESC LIMIT 1;
```

### 2. Execute Request in Postman
Send `POST` or `PUT` request.

### 3. Post-Execution Database Verification
Run verification SQL queries:
```sql
-- 1. Verify primary record creation
SELECT id, accession_id, status, authorized_by, authorized_at 
FROM laboratory.results 
WHERE id = 'target-uuid';

-- 2. Verify foreign key integrity
SELECT r.id, a.accession_number, p.mrn 
FROM laboratory.results r
JOIN laboratory.accessions a ON a.id = r.accession_id
JOIN patient.patients p ON p.id = a.patient_id
WHERE r.id = 'target-uuid';

-- 3. Verify audit trail generation
SELECT actor_id, action, resource_type, resource_id, payload, created_at 
FROM audit.audit_events 
ORDER BY created_at DESC LIMIT 1;
```

---

## 7. Security Review Checklist

Engineers auditing endpoints must test for 10 security vulnerabilities:

1. **Authentication Enforcement**: Call endpoint without cookies/headers -> Must return `401 Unauthorized`.
2. **Authorization & RBAC**: Call endpoint using a `member` role token -> Must return `403 Forbidden` if `admin` permission is required.
3. **Tenant Context Isolation**: Request tenant resource `A` using credentials from tenant `B` -> Must return `404 Not Found` or `403 Forbidden` (no cross-tenant leakage).
4. **SQL Injection Prevention**: Pass `' OR 1=1 --` into search query parameters -> Must execute parameterized query cleanly without syntax error.
5. **CSRF Protection**: Verify state-changing requests (`POST`, `PUT`, `DELETE`) validate `X-CSRF-Token` headers.
6. **XSS Protection**: Pass `<script>alert(1)</script>` in input fields -> Must sanitize or escape HTML content in JSON outputs.
7. **Mass Assignment Prevention**: Attempt to inject `"is_platform_admin": true` in a profile update payload -> Payload binding must ignore unmapped fields.
8. **IDOR (Insecure Direct Object Reference)**: Replace record UUID in URL path with another user's record -> Must enforce ownership check.
9. **Sensitive Data Exposure**: Verify user endpoints never return `password_hash`, `secret_keys`, or raw tokens in responses.
10. **Rate Limiting**: Execute 100 rapid requests within 5 seconds -> Must trigger `429 Too Many Requests`.

---

## 8. Performance Review Guide

### Metrics & Thresholds
- **Max Response Time**: Target < 50ms for reads, < 150ms for transactional writes.
- **Query Count**: Max 1 to 3 queries per HTTP request. Zero N+1 query loops permitted.
- **Database Indexing**: All `WHERE`, `JOIN`, and `ORDER BY` columns must use PostgreSQL indexes (`btree`, `gin`).

### Query Verification Command
```sql
EXPLAIN ANALYZE
SELECT r.id, r.result_value, a.accession_number
FROM laboratory.results r
JOIN laboratory.accessions a ON a.id = r.accession_id
WHERE r.accession_id = 'c4b8b6a1-0000-0000-0000-000000000000';
```
*Verify output uses `Index Scan` rather than `Seq Scan`.*

---

## 9. Architecture Review Standards

Ensure strict adherence to Clean Architecture layers:

```
[Forbidden Practices]
❌ Writing SQL queries inside HTTP Handlers
❌ Rendering JSON directly inside Application Services
❌ Importing Echo framework packages in Repositories or Domain Services
❌ Modifying another module's database tables directly (must call module interface)
❌ Creating God packages (internal/services, internal/utils, internal/helpers)
```

---

## 10. OpenAPI Review Guide

Every endpoint must match [static/openapi.json](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/backend/static/openapi.json):

1. **Path Alignment**: Endpoint path matches OpenAPI path exact casing (`/api/v1/tenant/active`).
2. **Method**: Correct HTTP verb (`GET`, `POST`, `PUT`, `DELETE`).
3. **Tags**: Tagged with appropriate domain context (`Identity`, `Tenant`, `Laboratory`, `Platform`).
4. **OperationId**: Unique camelCase operation name (e.g. `getBootstrap`).
5. **Request Schema**: Every request body properties and required fields documented.
6. **Response Envelope**: Documented under `200`, `201`, `400`, `401`, `403`, `500`.

---

## 11. Production Readiness Checklist (100+ Items)

### A. Authentication & Authorization (1–10)
- [ ] 1. Authentication middleware applied to protected routes
- [ ] 2. Passwords hashed using bcrypt/Argon2
- [ ] 3. JWT/Session tokens expire within configured TTL
- [ ] 4. Session revocation endpoint (`POST /auth/sign-out`) active
- [ ] 5. Permission scopes validated on every handler
- [ ] 6. Role-based access control (RBAC) enforced via Casbin/Service
- [ ] 7. Cross-tenant access strictly blocked
- [ ] 8. Service-to-service internal calls use `x-service-token`
- [ ] 9. Platform admin actions require dual-factor verification
- [ ] 10. Impersonation events generate audit logs

### B. Validation & Input Sanitation (11–20)
- [ ] 11. Struct tag validation on all incoming DTOs
- [ ] 12. UUID strings validated via `uuid.Parse`
- [ ] 13. String lengths bounded (`maxLength`)
- [ ] 14. Email format validated via regex
- [ ] 15. HTML tags stripped from text inputs
- [ ] 16. JSON payload max body size limited (e.g. 2MB)
- [ ] 17. Empty strings rejected on mandatory fields
- [ ] 18. Numeric parameters bounds-checked (e.g. positive prices)
- [ ] 19. Enums validated against allowed value lists
- [ ] 20. Unknown JSON fields stripped during binding

### C. Database & Transactions (21–35)
- [ ] 21. Schema owned exclusively by bounded context
- [ ] 22. Migration files exist in `database/platform/migrations` or `database/tenant/migrations`
- [ ] 23. `SET search_path` executed per tenant request
- [ ] 24. Multi-table writes execute inside PostgreSQL transaction (`pgx.Tx`)
- [ ] 25. Defer `tx.Rollback(ctx)` configured on all transactions
- [ ] 26. Indexes present on foreign keys and search columns
- [ ] 27. Unique constraints handle race conditions
- [ ] 28. Soft deletes (`deleted_at`) implemented for clinical entities
- [ ] 29. Prepared statements reused by connection pool
- [ ] 30. DB connection pool max size tuned (`max_conns`)
- [ ] 31. DB queries use context timeout (`context.WithTimeout`)
- [ ] 32. Raw SQL uses named/positional binding parameters (no string formatting)
- [ ] 33. Cascading deletes configured on parent-child entities
- [ ] 34. TIMESTAMP WITH TIME ZONE used for all temporal fields
- [ ] 35. Database health check verified via `/status`

### D. Architecture & Code Quality (36–50)
- [ ] 36. Zero business logic inside HTTP Handlers
- [ ] 37. Zero SQL inside Handlers or Middleware
- [ ] 38. Zero Echo framework references in Repositories
- [ ] 39. Shared packages contain only infrastructure utilities (`internal/shared`)
- [ ] 40. Interfaces defined at consumption sites
- [ ] 41. Code formatted via `go fmt`
- [ ] 42. Code passes `go vet ./...` without warnings
- [ ] 43. Linter passes (`golangci-lint`)
- [ ] 44. Public symbols documented with Go doc comments
- [ ] 45. Errors wrapped with contextual message (`fmt.Errorf("failed to fetch user: %w", err)`)
- [ ] 46. Domain events emitted asynchronously for side effects
- [ ] 47. Dependency injection used in constructors
- [ ] 48. Zero hardcoded UUIDs, emails, or credentials in code
- [ ] 49. Configuration loaded via environment variables (`internal/shared/config`)
- [ ] 50. Application shuts down gracefully (`sigterm` handling)

### E. Performance & Scaling (51–60)
- [ ] 51. Read requests complete under 50ms (p95)
- [ ] 52. Write requests complete under 150ms (p95)
- [ ] 53. N+1 queries eliminated via SQL `JOIN` or `ANY($1)` batching
- [ ] 54. Pagination implemented for list endpoints (`limit`, `offset`/`cursor`)
- [ ] 55. Heavy responses compressed with Gzip
- [ ] 56. Statics/Assets served via CDN or cached headers
- [ ] 57. Redis cache implemented for read-heavy reference tables
- [ ] 58. Allocations minimized in critical execution paths
- [ ] 59. Unnecessary reflection avoided
- [ ] 60. Heavy background processing offloaded to async job queue

### F. Security & Vulnerability Defense (61–75)
- [ ] 61. HTTPS enforced in non-local environments
- [ ] 62. Security headers configured (HSTS, CSP, X-Frame-Options, X-Content-Type-Options)
- [ ] 63. CORS policy restricted to authorized origins
- [ ] 64. Rate limiting active on sensitive endpoints (auth, registration)
- [ ] 65. CSRF tokens validated on session-based write requests
- [ ] 66. Password reset tokens cryptographically secure (`crypto/rand`)
- [ ] 67. Password reset tokens expire within short TTL (e.g. 15 mins)
- [ ] 68. Insignificant error details masked from client responses
- [ ] 69. Stack traces never exposed in HTTP error envelopes
- [ ] 70. Dependency vulnerabilities scanned via `govulncheck`
- [ ] 71. Secrets stored in environment/vault (never committed to git)
- [ ] 72. Sensitive fields in logs masked (passwords, tokens, PHI)
- [ ] 73. Request payloads sanitized before logging
- [ ] 74. Cookies configured with `HttpOnly`, `Secure`, and `SameSite=Lax/Strict`
- [ ] 75. Session IDs regenerated upon login to prevent session fixation

### G. Logging & Observability (76–85)
- [ ] 76. Structured JSON logging active (`zerolog` / `slog`)
- [ ] 77. Unique `X-Request-ID` generated and passed across log contexts
- [ ] 78. Log levels properly categorized (`DEBUG`, `INFO`, `WARN`, `ERROR`)
- [ ] 79. Audit logs emitted to `audit.audit_events` table for compliance
- [ ] 80. Metrics endpoint exposed for Prometheus tracking
- [ ] 81. Database connection pool stats monitored
- [ ] 82. Panic recovery middleware catches unhandled panics and returns 500
- [ ] 83. Slow queries (>200ms) logged with execution parameters
- [ ] 84. External service calls logged with duration and status
- [ ] 85. Distributed tracing headers propagated (`traceparent`)

### H. OpenAPI & Testing (86–100)
- [ ] 86. OpenAPI spec updated in `static/openapi.json`
- [ ] 87. OpenAPI response schemas match actual API response JSON
- [ ] 88. SDK client regenerates cleanly without build errors
- [ ] 89. Unit tests written for Application Services (`*_test.go`)
- [ ] 90. Integration tests written for HTTP API Handlers
- [ ] 91. Repository tests run against PostgreSQL test container/database
- [ ] 92. Multi-tenant isolation test verified
- [ ] 93. Transaction rollback test verified
- [ ] 94. Code coverage exceeds target threshold (>80%)
- [ ] 95. Postman manual verification completed cleanly
- [ ] 96. Health check route `/status` returns `200 OK`
- [ ] 97. API versioning strategy enforced (`/api/v1`)
- [ ] 98. Deprecation notices documented if replacing legacy routes
- [ ] 99. Deployment playbook documented
- [ ] 100. Production Audit Report generated and signed off

---

## 12. API Audit Template

Use this markdown template when reviewing any Curexal endpoint:

```markdown
# API Audit Report: [ENDPOINT_NAME]

- **Path**: `METHOD /api/v1/...`
- **Bounded Context**: `[module_name]`
- **Auditor**: [Engineer Name]
- **Date**: [YYYY-MM-DD]
- **Verdict**: [ APPROVED / REJECTED ]

### 1. Business Purpose & Context
[Describe why the endpoint exists and what clinical/operational use case it serves]

### 2. Database Impact
- **Schemas**: `[schema_name]`
- **Tables Read**: `[table_1]`, `[table_2]`
- **Tables Modified**: `[table_3]`
- **Transaction Scope**: `[ACID Tx details or N/A for read]`

### 3. Contract & OpenAPI Match
- [ ] OpenAPI spec matched in `static/openapi.json`
- [ ] Standard Envelope returned (`{ data, meta, links, errors }`)

### 4. Security & Multi-Tenancy Review
- **Auth Required**: [ Yes / No ]
- **Permissions Required**: `[permission_code]`
- **Tenant Isolation Method**: `SET search_path = [tenant]`

### 5. Automated & Manual Test Coverage
- Unit Tests: `PASS`
- Repository Tests: `PASS`
- Postman Manual Test: `PASS`

### 6. Performance & Index Analysis
- Execution Time: `XX ms`
- Query Count: `X queries`
- Index Status: `Index Scan verified via EXPLAIN ANALYZE`

### 7. Auditor Sign-Off
Approved for production merge: **[ YES / NO ]**
```

---

## 13. Complete Example Walkthrough: `GET /api/v1/bootstrap`

Let's trace a real endpoint from client request to database to Postman verification:

### Step 1: Client Issues HTTP GET Request
```http
GET /api/v1/bootstrap HTTP/1.1
Host: localhost:8080
Cookie: curexal_session=sess_abc123xyz
Accept: application/json
```

### Step 2: Echo Router & Middleware Interception
1. `Echo` matches `/api/v1/bootstrap` to `BootstrapHandler.GetBootstrap`.
2. `middleware.Authenticate` reads `curexal_session` cookie, verifies active session in database, and attaches `AuthenticatedPrincipal` to `c.Request().Context()`.
3. `middleware.TenantResolver` inspects tenant context and issues `SET search_path = 'tenant_slug', public`.

### Step 3: Handler Execution
`BootstrapHandler.GetBootstrap(c echo.Context)` extracts `principal` from context and calls Application Service:
```go
principal := middleware.GetPrincipal(c)
contract, err := h.bootstrapBuilder.BuildBootstrap(c.Request().Context(), principal)
```

### Step 4: Application Service Execution (`BootstrapBuilder`)
`BootstrapBuilder.BuildBootstrap` executes domain orchestration:
1. Queries `organization` and `tenant` repositories for active facility metadata.
2. Queries PostgreSQL `permission` and `role_permission` tables for user's effective permissions:
   ```sql
   SELECT DISTINCT p.code
   FROM permission p
   JOIN role_permission rp ON rp.permission_id = p.id
   JOIN role r ON r.id = rp.role_id
   JOIN membership m ON (m.role_title = r.code OR m.role = r.code)
   WHERE m.user_id = $1 AND m.is_active = TRUE;
   ```
3. Queries PostgreSQL `navigation_item` table for context-bound navigation menu:
   ```sql
   SELECT id, title, icon, path, sort_order
   FROM navigation_item
   WHERE context_scope = $1 AND (module_code IS NULL OR module_code = ANY($2))
   ORDER BY sort_order ASC;
   ```
4. Assembles `BootstrapContractResponse` struct.

### Step 5: Handler Renders JSON Response
`BootstrapHandler` returns `200 OK` with JSON envelope:
```json
{
  "identity": {
    "id": "usr_1001",
    "email": "admin@curexal.com",
    "displayName": "Dr. Chief Medical Officer",
    "locale": "en",
    "timezone": "Africa/Lagos"
  },
  "platform": {
    "isStaff": false,
    "role": "branch_admin"
  },
  "organization": {
    "id": "org_5001",
    "name": "Curexal Health Network",
    "subscription": "Enterprise"
  },
  "workspace": {
    "id": "wsp_9001",
    "name": "Central Diagnostic Facility",
    "facilityType": "Laboratory",
    "slug": "main-facility",
    "timezone": "Africa/Lagos",
    "currency": "NGN"
  },
  "subscription": {
    "plan": "Enterprise",
    "status": "active",
    "limits": {
      "maxBranches": 100,
      "maxMembers": 1000
    }
  },
  "permissions": [
    "workspace:patient:read",
    "workspace:patient:create",
    "workspace:sample:receive",
    "workspace:worksheet:update",
    "workspace:result:authorize",
    "workspace:billing:create"
  ],
  "navigation": [
    { "id": "nav_wsp_dashboard", "title": "Workspace Dashboard", "icon": "LayoutDashboard", "path": "/workspace/dashboard", "order": 1 },
    { "id": "nav_wsp_patients", "title": "Patient Reception", "icon": "UserPlus", "path": "/workspace/patients", "order": 2 },
    { "id": "nav_wsp_laboratory", "title": "Laboratory LIS", "icon": "Activity", "path": "/workspace/laboratory/accessioning", "order": 3 }
  ]
}
```

---

## 14. Backend Engineering Career Roadmap

```
Junior Backend Engineer
       │ (Mastering Go syntax, SQL queries, HTTP handlers, basic Postman testing)
       ▼
Mid-Level Backend Engineer
       │ (Domain module design, ACID transaction management, Clean Architecture, unit/integration testing)
       ▼
Senior Backend Engineer
       │ (Schema per tenant multi-tenancy, OpenAPI specification design, performance indexing, security audits)
       ▼
Staff Backend Engineer
       │ (System architecture design, microservice boundary planning, platform performance, engineering standards)
       ▼
Principal Backend Engineer
       │ (Enterprise system kernel design, security constitution, multi-app SDK generation strategy, organization mentoring)
       ▼
Distinguished Engineer
       │ (Healthcare platform industry architecture, high-availability zero-downtime distributed kernel systems)
```

| Career Level | Core Responsibilities | Architecture & DB Expectations | Testing & Audit Expectations |
| :--- | :--- | :--- | :--- |
| **Junior Engineer** | Implements CRUD handlers and simple queries | Writes clean Go structs; understands SQL SELECT/INSERT | Writes basic unit tests; performs manual Postman tests |
| **Mid-Level Engineer** | Implements bounded context services and repositories | Manages ACID transactions; prevents SQL injection; indexes FKs | Writes service integration tests; enforces validation |
| **Senior Engineer** | Audits API production readiness; designs DB migrations | Enforces Clean Architecture; implements schema-per-tenant isolation | 100% production audit checklist compliance |
| **Staff Engineer** | Owns multi-module architecture and performance | Eliminates query bottlenecks; tunes DB pools & Redis caches | Designs end-to-end integration and load testing suites |
| **Principal Engineer** | Sets company-wide backend constitution & standards | Architects platform kernels, security frameworks, and OpenAPI SDKs | Defines enterprise audit handbook and quality standards |
| **Distinguished** | Strategic infrastructure & high-availability design | Zero-downtime migrations; multi-region distributed DB systems | Sets industry healthcare compliance and engineering excellence |
