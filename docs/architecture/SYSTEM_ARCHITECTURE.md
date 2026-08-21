# Enterprise System Architecture Specification

> **Purpose**: Technical specification of Curexal's Modular Monolith architecture, Domain Driven Design (DDD) boundaries, layer decoupling, dependency injection, and multi-tenant operating model.  
> **Owner**: Principal Systems Architect  
> **Status**: APPROVED / PRODUCTION READY  
> **Last Updated**: 2026-07-27  
> **Verification Criteria**: Verified against live backend Go modules and Hertz HTTP routers.

---

## 1. High-Level Monolith Architecture

```text
 ┌────────────────────────────────────────────────────────────────────────────────────────┐
 │                              API Gateway / Hertz HTTP Router                           │
 └───────────────────────────┬──────────────────────────────── font-bold ─────────────────┘
                             │
         ┌───────────────────┼───────────────────┬───────────────────┐
         │                   │                   │                   │
         ▼                   ▼                   ▼                   ▼
┌──────────────────┐┌──────────────────┐┌──────────────────┐┌──────────────────┐
│ Platform Context ││ Identity Context ││ Organization Context││ Provisioning Context│
│ (internal/platform)││ (internal/auth)  ││ (internal/org)   ││ (internal/prov)  │
└────────┬─────────┘└────────┬─────────┘└────────┬─────────┘└────────┬─────────┘
         │                   │                   │                   │
         └───────────────────┼───────────────────┴───────────────────┘
                             │
                             ▼
 ┌────────────────────────────────────────────────────────────────────────────────────────┐
 │                      Core Infrastructure Engine (`internal/core`)                      │
 ├───────────────────────────┬───────────────────────────────────┬────────────────────────┤
 │ Bun ORM PostgreSQL Pool   │ Redis Session & Cache Store       │ Casbin RBAC Engine     │
 └───────────────────────────┴───────────────────────────────────┴────────────────────────┘
```

---

## 2. Strict Bounded Context Decoupling

In Curexal V2, every domain capability is isolated inside its own bounded context (`internal/modules/<bounded_context>`):

```text
internal/modules/<context>/
├── app/          # Application Use Cases & Command Handlers
├── domain/       # Domain Entities, Value Objects & Repository Interfaces
├── infra/        # Infrastructure Implementations (Bun ORM Repositories & SQL Migrations)
├── api/          # Hertz REST Handlers & Route Registration
└── module.go     # Dependency Injection Container & Module Lifecycle Interface
```

### Architectural Constraints
- **Repositories Never Invoke HTTP**: Repository implementations (`infra/`) execute database queries exclusively and never execute network HTTP requests.
- **Handlers Never Query Database Directly**: HTTP REST handlers (`api/`) validate DTOs and delegate execution to Application Use Cases (`app/`).
- **No Direct Table Cross-Querying**: Modules must interact with other modules through exported Service interfaces or Domain Events.

---

## 3. Dual Isolated Middleware Chains

To enforce strict isolation between global platform administration and multi-tenant customer operations, Curexal employs two independent middleware chains:

```text
                            Incoming HTTP Request
                                     │
                 ┌───────────────────┴───────────────────┐
                 │ Path Router Matching                  │
                 └─────────┬───────────────────┬─────────┘
                           │                   │
          /api/v1/platform/*                   /api/v1/tenant/* (or /api/v1/orgs/*)
                           │                   │
                           ▼                   ▼
              ┌──────────────────────────┐   ┌──────────────────────────┐
              │ RequirePlatformSession   │   │ ResolveTenant            │
              │ (Validates public.sessions│   │ (Extracts Subdomain/Header│
              │  & PlatformUserMembership│   │  Sets tenant_slug ctx)   │
              └────────────┬─────────────┘   └────────────┬─────────────┘
                           │                              │
                           ▼                              ▼
              ┌──────────────────────────┐   ┌──────────────────────────┐
              │ RequirePlatformRBAC      │   │ RequireTenantMembership  │
              │ (Casbin domain: platform)│   │ (Validates tenant membership)
              └────────────┬─────────────┘   └────────────┬─────────────┘
                           │                              │
                           ▼                              ▼
              ┌──────────────────────────┐   ┌──────────────────────────┐
              │ Platform API Handlers    │   │ Tenant Operational Handlers│
              └──────────────────────────┘   └──────────────────────────┘
```
