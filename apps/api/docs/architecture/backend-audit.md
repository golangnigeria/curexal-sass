# CUREXAL BACKEND — PRODUCTION ARCHITECTURE AUDIT REPORT
**Document Status**: COMPLETED (Phase 1 — Audit)  
**Author**: Lead Principal Backend Architect  
**Target Workspace**: `curexal-backend`  

---

## EXECUTIVE SUMMARY

A comprehensive pre-code architecture audit of the **Curexal Backend** codebase has been performed. This audit evaluates all packages, database migrations, authentication mechanisms, authorization implementations, multi-tenancy controls, handlers, application services, tests, and OpenAPI definitions against the **Curexal Backend Constitution (v2.0)**.

---

## 1. CURRENT ARCHITECTURE OVERVIEW

### Tech Stack & Runtime Components
- **Language & Runtime**: Go 1.25.7
- **HTTP Engine**: Echo `v4.13.4`
- **Database Driver & Pool**: `pgx/v5` (`github.com/jackc/pgx/v5`)
- **Database Migrations**: Goose `v3.27.2`
- **Caching & KV**: Redis `v9.7.0`
- **Logging**: Zerolog `v1.34.0`
- **Token Handling**: `golang-jwt/jwt/v5`
- **Password Hashing**: Argon2id (`golang.org/x/crypto`)

### Structural Layout
The current repository structure exhibits a hybrid, non-uniform architectural pattern:
- Core entry point: `cmd/CUREXAL/main.go` initializes server, config, logging, migrations, and calls `bootstrap.InitModules(srv)`.
- Central composition root: `internal/bootstrap/modules.go` instantiates modules and registers HTTP routes.
- Modules in `internal/modules/`: `audit`, `authorization`, `billing`, `catalogs`, `clinical`, `facility_config`, `identity`, `notification`, `organization`, `patient`, `platform`, `settings`.
- Mixed architectural patterns across modules:
  - **Clean Architecture Pattern** (`domain/`, `application/`, `infrastructure/`, `api/`): Used by `organization`, `platform`, `clinical`, `facility_config`, `settings`.
  - **Legacy MVC Pattern** (`handler/`, `model/`, `repository/`, `service/`): Used by `identity`, `authorization`, `patient`, `catalogs`, `notification`.

---

## 2. CONSTITUTION & ARCHITECTURAL VIOLATIONS

### 2.1 Handlers Containing Business & Persistence Logic
- **Violation**: HTTP handlers in `internal/modules/identity/handler/user_role.go`, `auth.go`, `invite.go`, and `internal/modules/clinical/api/catalog_handler.go` execute raw SQL queries, perform multi-step tenant provisioning, and handle complex domain orchestration directly inside handler methods.
- **Constitution Rule (Section 16)**: Handlers must be thin wrappers that parse requests, invoke application services, and format responses. Handlers must NOT execute SQL or perform multi-step workflows.

### 2.2 Hardcoded Role Comparisons
- **Violation**: Inline string role comparisons are scattered across 12+ codebase files. Examples:
  ```go
  // In internal/shared/middleware/middleware.go:
  isStaff := isPlatformAdmin || platformRole == "super_admin" || platformRole == "platform_staff" || platformRole == "super_support_agent" || platformRole == "super_sales_staff"
  
  // In internal/platform/auth/permission_provider.go:
  if principal.Role == "owner" || principal.Role == "org_admin" || principal.Role == "org_regional_manager" ...
  ```
- **Constitution Rule (Section 9)**: Authorization decisions must be permission-driven (`authorizationService.Require(principal, permission.OrganizationManage)`). Scattered `if role == "admin"` checks bypass DB-backed authorization.

### 2.3 Non-Standard API Response Envelopes
- **Violation**: Handlers across `identity`, `organization`, `clinical`, and `catalogs` return varied HTTP response structures (raw JSON structs, custom status wrappers, or bare arrays).
- **Constitution Rule (Section 15)**: Every API endpoint must return a single standardized envelope:
  ```json
  {
    "data": {},
    "meta": {},
    "links": {},
    "errors": []
  }
  ```

---

## 3. SECURITY RISKS & VULNERABILITIES

### 3.1 Unauthenticated Header Identity Trust
- **Severity**: HIGH
- **Location**: `internal/shared/middleware/middleware.go` (`resolvePrincipal`) and `internal/platform/auth/resolver.go` (`ResolvePrincipalWithProvider`)
- **Detail**: When JWT tokens are absent or `AllowTestHeaders` is enabled, the backend reads `X-User-ID`, `X-Tenant-ID`, and `X-User-Role` directly from HTTP headers and populates the `AuthenticatedPrincipal`. An external client can send:
  ```http
  X-User-ID: admin-user-uuid
  X-Tenant-ID: victim-tenant-uuid
  X-User-Role: super_admin
  ```
  and bypass authentication entirely to gain unauthorized access to any tenant.

### 3.2 Unverified Multi-Tenancy Resolution
- **Severity**: HIGH
- **Location**: `internal/shared/middleware/middleware.go` (`HostAndSessionTenantResolver`)
- **Detail**: The tenant ID extracted from `X-Tenant-ID` or host subdomains is set as the active tenant context without verifying against PostgreSQL that the authenticated user actually holds an active membership in that organization or workspace.

### 3.3 Default / Weak Hardcoded Credentials in Seeders
- **Severity**: MEDIUM
- **Location**: `cmd/seed_admin/main.go`
- **Detail**: Default administrative accounts (`admin@curexal.com`, `superadmin@curexal.internal`) are seeded with hardcoded password `"password"`.

---

## 4. DATABASE DUPLICATION & SCHEMA DUALITY

### 4.1 Dual Entity Tables (`public` vs Schema-Qualified)
- **Location**: `database/platform/migrations/000001_database_constitution.sql`
- **Detail**: Migration `000001` creates tables in bounded context schemas (`identity.users`, `organization.organizations`, `workspace.workspaces`, `"authorization".roles`, `"authorization".permissions`, `"authorization".role_permissions`, `patient.patients`, `audit.audit_events`) AND ALSO creates legacy tables in the `public` schema (`public."user"`, `public.organization`, `public.tenant`, `public.membership`, `public.role`, `public.permission`, `public.role_permission`, `public.audit_log`, `public.subscription`, `public.module`).
- **Impact**: Domain repositories (`user.go`, `tenant_repository.go`, `enforcer.go`) query `public.*` tables, while new schemas remain unpopulated or partially shadowed, violating the **Single Source of Truth** principle.

| Canonical Entity | Public Schema Table | Bounded Context Table | Status |
| :--- | :--- | :--- | :--- |
| Users | `public."user"` | `identity.users` | DUPLICATED |
| Profiles | `public.user_profile` | `identity.user_profiles` | DUPLICATED |
| Organizations | `public.organization` | `organization.organizations` | DUPLICATED |
| Workspaces / Tenants | `public.tenant` | `workspace.workspaces` | DUPLICATED |
| Memberships | `public.membership` | `organization.organization_memberships` | DUPLICATED |
| Roles | `public.role` | `"authorization".roles` | DUPLICATED |
| Permissions | `public.permission` | `"authorization".permissions` | DUPLICATED |
| Role Permissions | `public.role_permission` | `"authorization".role_permissions` | DUPLICATED |
| Audit Logs | `public.audit_log` | `audit.audit_events` | DUPLICATED |

---

## 5. BOUNDED CONTEXT VIOLATIONS

### 5.1 Cross-Context Repository & Direct Database Access
- **Violation**: In `internal/bootstrap/modules.go`, `tenantLookupAdapter` directly instantiates `postgres.TenantRepository` from the `organization` module to serve `identity` lookup queries.
- **Violation**: `identity.UserRoleHandler` manages organization memberships, workspace memberships, professional profiles, digital signatures, employment records, and role assignments in one monolithic handler struct.
- **Constitution Rule (Section 5)**: A module MUST NOT directly manipulate another module's private persistence. Cross-context communication must occur through explicit application/domain interfaces.

### 5.2 Missing Bounded Contexts
- **Missing Bounded Context Packages**: `subscription`, `marketplace`, `radiology`, `pharmacy`, `inventory`. Standard capabilities for these contexts are currently stubbed or mixed inside `billing` and `facility_config`.

---

## 6. AUTHORIZATION RISKS

### 6.1 In-Memory Permission Providers
- **Location**: `internal/platform/auth/permission_provider.go`
- **Detail**: Permissions are resolved primarily using `MemoryRoleMapPermissionProvider` which maps static strings in Go code rather than querying PostgreSQL `"authorization".role_permissions` per request.

### 6.2 Scope-Agnostic Permission Checks
- **Location**: `internal/shared/middleware/middleware.go` (`RequirePermission`)
- **Detail**: Permission checks verify global permission string strings without validating tenant context scope (`context_scope`), allowing a permission granted in Workspace A to potentially authorize actions in Workspace B if context switching is bypassed.

---

## 7. TENANT ISOLATION RISKS

### 7.1 Search Path Connection Leakage
- **Location**: `database/tenant/migrations/000001_create_clinic_schema.sql` and `internal/platform/database`
- **Detail**: Multi-tenant isolation relies on PostgreSQL schema search paths (`SET search_path TO tenant_<id>`). However:
  - Connection pooling with `pgxpool` does not enforce automatic search path reset (`RESET search_path`) when connections are returned to the pool.
  - Reused pool connections risk executing tenant queries against the schema of the previously served request.

### 7.2 Lack of Multi-Tenant Verification Tests
- **Detail**: No integration test exists in `internal/testing/` that verifies cross-tenant access rejection (`User A -> Tenant B -> 403 Forbidden`).

---

## 8. TRANSACTION RISKS

### 8.1 Non-Transactional Multi-Record Operations
- **Location**: `internal/modules/organization/application/organization_command.go` and `internal/modules/identity/handler/user_role.go`
- **Detail**: Workflows such as "Create Organization -> Create Tenant -> Create Owner Membership -> Provision Workspace" execute multiple separate database `Exec` calls outside a `pgx.Tx` transaction block.
- **Impact**: If workspace creation fails due to a database constraint or connection drop, the organization and owner membership remain created in an incomplete, un-provisioned state without automatic rollback.

---

## 9. TEST GAPS

### 9.1 Missing Database Integration & Security Tests
- `internal/testing/api_contract_test.go` and `e2e_http_verification_test.go` use dummy HTTP recorders and in-memory server structs without database assertion.
- Missing Automated Test Coverage:
  - No authorization enforcement tests.
  - No tenant isolation & header spoofing rejection tests.
  - No transaction rollback failure tests.
  - No migration integrity tests against real PostgreSQL testcontainers.

---

## 10. OPENAPI CONTRACT GAPS

### 10.1 Route Mismatches
- `static/openapi.json` defines routes under `/tenant/active`, `/tenant/{id}`, `/tenant` with different parameters than registered in `internal/bootstrap/modules.go` (`/api/v1/tenant/active`, `/api/v1/tenants`, `/api/v1/organizations`).

### 10.2 Schema Mismatches
- Request and response schemas in `openapi.json` do not reflect the single standardized API response envelope (`{ "data": ..., "meta": ..., "links": ..., "errors": ... }`).

---

## 11. DEAD CODE & OBSOLETE ARTIFACTS

### 11.1 Obsolete CLI Binaries
- `cmd/clean_db/main.go`
- `cmd/convert_migrations/main.go`
- `cmd/migrate_kratos/main.go`
- `cmd/query_db/main.go`
- `cmd/tasker/main.go`

### 11.2 Abandoned ORY Kratos Infrastructure
- `kratos_identity_id` columns in `identity.users` and `public."user"`.
- `ResolvePrincipalWithProvider` ORY Kratos session token verification blocks in `internal/platform/auth/resolver.go`.

---

## SUMMARY OF AUDIT FINDINGS

| Audit Category | Current Status | Action Required |
| :--- | :--- | :--- |
| **Architecture Model** | Hybrid / Inconsistent | Standardize on PostgreSQL -> Domain -> Application -> Infrastructure -> API |
| **Bounded Contexts** | Partial / Mixed | Establish explicit contexts (`domain/`, `application/`, `infrastructure/`, `api/`) |
| **Database SSOT** | Duplicated (`public` vs schemas) | Define canonical schema for every entity & consolidate migrations |
| **Multi-Tenancy** | Header-trusted / Unverified | Implement DB-backed `TenantMembershipVerifier` & sanitize `search_path` |
| **Authentication** | Spoofable via headers | Remove header-based auth fallbacks; enforce JWT/Session validation |
| **Authorization** | Hardcoded strings in Go | Enforce DB-backed `authorization.roles` and `authorization.permissions` |
| **Transactions** | Non-transactional writes | Enforce ACID transaction blocks (`pgx.Tx`) on multi-step mutations |
| **API Envelope** | Inconsistent across handlers | Enforce strict `{ "data", "meta", "links", "errors" }` format |
| **Handlers** | Fat / SQL-executing | Refactor to thin wrappers delegating to application services |
| **OpenAPI** | Desynchronized from routes | Re-generate/align `static/openapi.json` with runtime routes |
| **Testing** | Lacks DB security tests | Implement live containerized multi-tenant & security integration tests |
