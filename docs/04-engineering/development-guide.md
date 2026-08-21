# Curexal V2 Development Standards & AI Rulebook

This document defines the mandatory engineering rules, design constraints, and quality standards for all contributors—both human software engineers and AI coding agents. Violating these standards will lead to automatic pull request rejection.

---

## 1. AI Agent Core Directives & Architectural Governance

AI coding assistants and human developers contributing to Curexal V2 must strictly follow these non-negotiable rules:

1. **Never Violate Architecture**: Always follow Clean Architecture, Domain-Driven Design (DDD), and Vertical Slice isolation.
2. **Never Invent Unapproved Modules**: All new endpoints and capabilities must be declared in the Capability Map (`docs/02-product/capability-map.md`) and Master Roadmap (`docs/00-overview/roadmap.md`).
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
├── docs/                      # Master blueprint documents
│   ├── 00-overview/           # Vision, Mission, Philosophy, Roadmap
│   ├── 01-business/           # Business Model, Pricing, Customer Personas
│   ├── 02-product/            # Capability Map, Domain Map, Workflows, Events, Domain Models
│   ├── 03-architecture/       # System Architecture, Tenancy, Edge Agent, Sharding
│   └── 04-engineering/        # Development Guide, Standards, AI Rules
```

### Layer Isolation Rules
- **Domain Layer (`domain/`)**: Pure business logic. Must have **zero dependencies** on external frameworks, Bun ORM, Hertz HTTP, or SQL drivers.
- **Application Layer (`app/`)**: Orchestrates domain objects to fulfill use cases. Cross-module calls must use defined service interfaces, not raw database queries.
- **Presentation Layer (`api/`)**: Translates HTTP JSON payloads into DTOs, validates inputs, invokes use cases, and handles RFC 7807 problem details.
- **Infrastructure Layer (`infra/`)**: Implements interfaces defined in the domain layer using Bun ORM, Redis, MinIO, or NATS.
