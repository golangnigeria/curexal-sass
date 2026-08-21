# CUREXAL BACKEND — TARGET ARCHITECTURE SPECIFICATION
**Document Status**: APPROVED BLUEPRINT (Phase 2 — Target Architecture)  
**Author**: Lead Principal Backend Architect  
**Target Workspace**: `curexal-backend`  

---

## 1. OVERVIEW & SYSTEM ARCHITECTURE MODEL

The Curexal Backend is designed as an independent, enterprise-grade, multi-tenant healthcare kernel. PostgreSQL serves as the Single Source of Truth (SSOT). All business rules, role definitions, permissions, navigation items, plan limits, and tenant boundaries originate strictly from PostgreSQL.

### Layered Architecture Flow
```text
PostgreSQL (SSOT)
       ↓
Domain Models & Interfaces
       ↓
Application Services (Use Cases & ACID Orchestration)
       ↓
Infrastructure Repositories (pgx/v5 & Redis)
       ↓
Thin HTTP API Handlers (Echo v4)
       ↓
OpenAPI 3.0 Contract (static/openapi.json)
       ↓
Downstream API Consumers (Web, Mobile, Postman, Integrations)
```

---

## 2. BOUNDED CONTEXT ARCHITECTURE

Every domain context is self-contained under `internal/modules/<context_name>/` and structured into clear internal layers:

```text
internal/modules/<context_name>/
├── domain/            # Entities, Value Objects, Domain Errors, Domain Events
├── application/       # Use Cases, Application Services, Command/Query DTOs
├── infrastructure/    # PostgreSQL Repositories, Redis Cache, External Adapters
└── api/               # Thin HTTP Handlers, Request/Response Mappers, Route Wiring
```

### Canonical Bounded Context Inventory
1. `identity`: Manages user accounts, authentication credentials, profiles, password history, and verification tokens.
2. `organization`: Manages organizations, custom domains, organization memberships, branches, and organization approval states.
3. `workspace`: Manages facilities/workspaces, workspace memberships, and facility operating settings.
4. `authorization`: Manages database-backed roles, permissions, role-permission mappings, and authorization evaluation.
5. `subscription`: Manages subscription plans, plan limits, active organization subscriptions, and feature availability.
6. `marketplace`: Manages available platform modules and subscription add-ons.
7. `patient`: Manages patient master index (EMR/MRN), patient demographics, and patient portal access.
8. `clinical`: Manages clinical visits, encounters, orders, and clinical documentation.
9. `laboratory`: Manages LIS test catalogs, accessioning, specimen processing, worksheets, and result authorization.
10. `radiology`: Manages radiology order workflows, imaging accessions, and diagnostic reports.
11. `pharmacy`: Manages medication prescriptions, dispensing, and order fulfillment.
12. `inventory`: Manages medical inventory stocks, item catalogs, and stock movements.
13. `audit`: Records immutable audit log events across all operational mutations.

### Inter-Module Communication Rule
Modules **MUST NOT** directly access or manipulate the persistence tables of another context. Communication across contexts occurs exclusively through Go interfaces (Ports & Adapters) defined in application services.

---

## 3. SINGLE SOURCE OF TRUTH DATABASE OWNERSHIP

All entity schemas are assigned to explicit PostgreSQL schemas. Duplication between the `public` schema and context schemas is completely removed.

### Canonical Database Schema Mapping

| Context | PostgreSQL Schema | Canonical Tables |
| :--- | :--- | :--- |
| `identity` | `identity` | `users`, `credentials`, `user_profiles`, `password_histories`, `verification_tokens` |
| `organization` | `organization` | `organizations`, `organization_memberships`, `organization_domains`, `branches` |
| `workspace` | `workspace` | `workspaces`, `workspace_memberships`, `facility_settings` |
| `authorization` | `authorization` | `roles`, `permissions`, `role_permissions` |
| `subscription` | `subscription` | `plans`, `subscriptions` |
| `marketplace` | `marketplace` | `modules` |
| `patient` | `patient` | `patients` |
| `laboratory` | `laboratory` | `test_catalogs`, `accessions`, `results` |
| `audit` | `audit` | `audit_events` |

---

## 4. AUTHENTICATION & MULTI-TENANT RESOLUTION

### 4.1 Authentication Pipeline
Authentication establishes **WHO YOU ARE**.
- Tokens are issued via JWT cookies or Authorization headers.
- Unauthenticated header hints (`X-User-ID`, `X-User-Role`) are **STRICTLY PROHIBITED** and ignored by authentication middleware in production.

### 4.2 Multi-Tenant Verification Pipeline
Tenant header `X-Tenant-ID` is treated strictly as a **REQUEST HINT**.
```text
HTTP Request
    ↓
1. Extract Bearer Token / Session Cookie -> Validate Cryptographic Signature -> Extract UserID
    ↓
2. Extract X-Tenant-ID Request Hint
    ↓
3. Execute PostgreSQL Database Lookup:
   SELECT role, is_active FROM organization.organization_memberships 
   WHERE user_id = $1 AND organization_id = $2
    ↓
4. If Membership Exists & Is Active -> Construct AuthenticatedPrincipal & Set Context
   If Membership Missing / Inactive -> Return 403 FORBIDDEN
```

---

## 5. DATABASE-BACKED AUTHORIZATION MODEL

Authorization determines **WHAT YOU ARE ALLOWED TO DO**.
- Hardcoded `if role == "owner"` logic in Go handlers is forbidden.
- Permission evaluation follows the database permission graph:
  ```text
  User -> Membership (organization_id / workspace_id) -> Role -> Role Permissions -> Permission Code
  ```
- Services invoke:
  ```go
  if err := authorizationService.Require(ctx, principal, "organization:write"); err != nil {
      return errs.ErrForbidden
  }
  ```

---

## 6. TRANSACTIONAL LIFECYCLE & ACID BOUNDARIES

Operations modifying multiple domain entities MUST execute within explicit database transactions (`pgx.Tx`).

### Organization & Workspace Provisioning Workflow
```text
BEGIN TRANSACTION (pgx.Tx)

1. INSERT INTO organization.organizations (name, slug, status) VALUES (...) RETURNING id;
2. INSERT INTO organization.organization_memberships (user_id, organization_id, role) VALUES (owner_id, org_id, 'owner');
3. INSERT INTO workspace.workspaces (organization_id, name, slug, facility_type) VALUES (org_id, 'Main Workspace', ...);
4. INSERT INTO subscription.subscriptions (organization_id, plan_id, status) VALUES (org_id, plan_id, 'active');
5. INSERT INTO audit.audit_events (actor_id, action, resource_type, resource_id) VALUES (...);

COMMIT TRANSACTION
```
*If any step fails, the entire operation is rolled back with zero partial state left in PostgreSQL.*

---

## 7. STANDARDIZED API RESPONSE ENVELOPE

Every HTTP endpoint returns the unified response envelope:

```json
{
  "data": {},
  "meta": {},
  "links": {},
  "errors": []
}
```

### Success Response Example (200 OK)
```json
{
  "data": {
    "id": "org_123",
    "name": "Everight Diagnostics",
    "slug": "everight"
  },
  "meta": {
    "requestId": "req_abc123",
    "timestamp": "2026-08-10T19:50:00Z"
  },
  "links": {
    "self": "/api/v1/organizations/org_123"
  },
  "errors": []
}
```

### Error Response Example (403 Forbidden)
```json
{
  "data": null,
  "meta": {
    "requestId": "req_abc124",
    "timestamp": "2026-08-10T19:50:05Z"
  },
  "links": {},
  "errors": [
    {
      "code": "FORBIDDEN",
      "message": "User does not have permission 'organization:manage' on the requested tenant."
    }
  ]
}
```

---

## 8. THIN HTTP HANDLER SPECIFICATION

Handlers are restricted to transport orchestration:
1. Parse HTTP headers and JSON body into DTO structs.
2. Validate DTO structural requirements.
3. Pass input DTO and `context.Context` to the corresponding application service use case.
4. Map domain application errors to standard HTTP response envelopes.

**Handlers MUST NOT execute SQL queries directly or perform multi-step business orchestration.**

---

## 9. TARGET DOCUMENTATION ARTIFACTS

As part of the architecture reset, the following documentation structure will be created:

```text
docs/
├── architecture/
│   ├── backend-audit.md            # [COMPLETED Phase 1]
│   ├── target-architecture.md      # [COMPLETED Phase 2]
│   ├── database-ownership.md       # Canonical schema & table mappings
│   ├── multi-tenancy.md            # Tenant verification & search_path specs
│   └── transaction-boundaries.md   # ACID workflow specifications
│
├── testing/
│   ├── postman-test-plan.md        # Comprehensive Postman test scenarios
│   └── security-tests.md           # Header spoofing & isolation test suite
│
└── api/
    └── api-behavior.md             # Standard error codes & response envelopes
```

---

## 10. REBUILD VERIFICATION GATEWAY

Before completing the architectural reset, the backend MUST pass:

```bash
go test ./... -v
go vet ./...
go build ./...
```
Along with containerized database integration tests proving:
- Multi-tenant isolation enforcement (Cross-tenant access returns HTTP 403).
- Header spoofing rejection (Unauthenticated `X-User-ID` returns HTTP 401).
- ACID transaction rollback on partial failure.
- OpenAPI 3.0 specification alignment.
