# Curexal V2 Development Standards & AI Rulebook

This document defines the mandatory engineering rules, design constraints, and quality standards for all contributors—both human software engineers and AI coding agents. Violating these standards will lead to automatic pull request rejection.

---

## 1. AI Agent Core Directives & Architectural Governance

AI coding assistants and human developers contributing to Curexal V2 must strictly follow these non-negotiable rules:

1. **Never Violate Architecture**: Always follow Clean Architecture, Domain-Driven Design (DDD), and Vertical Slice isolation.
2. **Never Invent Unapproved Modules**: All new endpoints and capabilities must be declared in the Capability Map (`docs/03-product-specification.md`) and Master Roadmap (`docs/02-roadmap.md`).
3. **Never Bypass RBAC or Tenant Authorization**: Every protected endpoint must verify authentication, tenant context (`X-Tenant-Slug`), organization/branch membership, and Casbin permissions.
4. **Never Write Placeholder Code or Mock Implementations**:
   - **PROHIBITED**: `// TODO`, fake data, dummy repositories, mock APIs, swallowed exceptions, or sample fallback objects.
   - **REQUIRED**: Real Bun ORM models, SQL migrations, repositories, application services, Hertz handlers, route registration, Casbin permissions, NATS events, audit logs, and tests.
5. **Never Introduce Unannounced Breaking Changes**: Identify impacts on existing modules, database schemas, permissions, events, and frontend consumers before proposing structural edits.

---

## 2. Monorepo Directory & Architecture Rules

```
Curexal V2 (Monorepo)
├── cmd/
│   └── server/                # Main application entrypoint
├── docs/                      # Master blueprint documents (01-05, openapi.yaml, adr/)
├── internal/
│   ├── core/                  # Core infrastructure (authz, database, event, rest, audit)
│   └── modules/               # Bounded contexts (auth, tenant, referral, etc.)
│       ├── api/               # Presentation (Hertz HTTP handlers, DTOs, route registration)
│       ├── app/               # Application Layer (Use cases, command/query handlers)
│       ├── domain/            # Domain Layer (Pure models, value objects, repo interfaces)
│       ├── infra/             # Infrastructure Layer (Bun ORM repos, S3, NATS adapters)
│       └── module.go          # Dependency injection wireframe
└── frontend/
    ├── apps/                  # portal, workspace, patient, admin
    └── packages/              # ui (shadcn/ui), api (bindings), utils
```

### Layer Isolation Rules
- **Domain Layer (`domain/`)**: Pure business logic. Must have **zero dependencies** on external frameworks, Bun ORM, Hertz HTTP, or SQL drivers.
- **Application Layer (`app/`)**: Orchestrates domain objects to fulfill use cases. Cross-module calls must use defined service interfaces, not raw database queries.
- **Presentation Layer (`api/`)**: Translates HTTP JSON payloads into DTOs, validates inputs, invokes use cases, and handles RFC 7807 problem details.
- **Infrastructure Layer (`infra/`)**: Implements interfaces defined in the domain layer using Bun ORM, Redis, MinIO, or NATS.

---

## 3. Database Isolation & Transaction Rules

- **Zero Cross-Schema SQL Joins**: Queries scoped to `tenant_a` must never join or reference tables in `tenant_b` or the `public` schema.
- **Transaction Safety**: All write operations within a tenant context must use the DB Manager's `RunInTenantTx` helper:
  ```go
  err := mgr.dbManager.RunInTenantTx(ctx, tenantSlug, func(ctx context.Context, tx bun.Tx) error {
      return repo.CreateSample(ctx, tx, sample)
  })
  ```
- **Migrations Only**: Database schema changes must be written as versioned SQL migration files under the target module's migration directory. Never run ad-hoc DDL.

---

## 4. Curexal Enterprise Frontend Design Standards

All frontend components, layout configurations, and screens generated across `apps/portal`, `apps/workspace`, `apps/patient`, and `apps/admin` must conform to these design constraints:

### Style Inspiration & Visual Standards
- **Aesthetic Benchmark**: High-fidelity, Dribbble-quality enterprise SaaS UI inspired by **Stripe Dashboard**, **Linear**, **Vercel**, **Notion**, and **Ramp**.
- **Visual Tone**: Clean, modern, hyper-professional SaaS platform.
- **Prohibited Patterns**: Generic admin templates, default Material UI / Bootstrap themes, cartoonish rounded buttons, and heavy dark gradients.

### Design System & Layout Rules
- **Color Palette**: Primary Blue (`#2563EB`), Accents in Emerald/Indigo, Neutrals in Slate/Zinc. Enforce accessible WCAG contrast standards.
- **Typography & Grid**: Inter font family, clear heading hierarchy, 8px grid spacing system (Tailwind CSS friendly).

### Mandatory Component States
Every generated UI component or page view must explicitly handle all 5 states:
1. **Loading State**: Render responsive skeleton loaders (e.g., `Skeleton` cards/tables).
2. **Empty State**: Render clean illustration or icon with descriptive text and primary action callout when no data exists.
3. **Error State**: Display clear problem detail feedback with retry triggers (conforming to RFC 7807 backend error responses).
4. **Success State**: Smooth transition with toast notifications or visual indicators.
5. **Permission Denied State**: Render unauthorized fallback card with clear guidance.

---

## 5. API & Event Messaging Standards

- **RFC 7807 Problem Details**: HTTP error responses must use `application/problem+json` formatting:
  ```json
  {
    "type": "https://curexal.com/errors/unauthorized",
    "title": "Permission Denied",
    "status": 403,
    "detail": "User lacks the 'referrals:accept' permission for this tenant branch.",
    "instance": "/api/v1/referrals/ref-99201/accept"
  }
  ```
- **NATS Event Naming**: Domain events must follow the format `[domain].[action]` (e.g., `referral.dispatched`, `results.released`).
- **Audit Logging**: Sensitive entity operations (patient records, referrals, billing actions) must record an immutable log entry in `internal/core/audit`.

---

## 6. PR Verification Checklist & Standards

At the end of every feature implementation, include a completed verification checklist in your pull request summary:

```markdown
## Verification Checklist

### Files Created
- `internal/modules/referral/domain/models.go`
- `internal/modules/referral/infra/repository.go`

### Files Modified
- `internal/router/router.go`

### Database Migrations
- `003_create_referrals_table.sql`

### API Endpoints Added
- `POST /api/v1/referrals` (201 Created)
- `GET /api/v1/referrals/{id}` (200 OK)

### Permissions Registered
- `referrals:create`
- `referrals:read`

### Events Registered
- `referral.dispatched`

### Postman Test Cases
- [x] Positive: Dispatches referral with valid payload (201 Created)
- [x] Negative: Rejects request with missing target_organization_id (400 Bad Request)
- [x] Negative: Returns 403 Forbidden when user lacks `referrals:create` permission

### Risks & Rollback Plan
- Risk: High volume referral dispatches could spike NATS memory.
- Rollback: Revert migration `003_create_referrals_table.sql` and restart server binary.
```
