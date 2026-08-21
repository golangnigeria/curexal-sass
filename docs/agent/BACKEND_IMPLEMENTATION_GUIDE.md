# Backend Implementation Guide & Layering Rules

> **Purpose**: Standard Operating Procedure defining the exact implementation sequence for backend Go modules.  
> **Owner**: Principal Backend Engineer  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. Required Layer Implementation Order

```text
Step 1: SQL Migration (`infra/migrations/00X_table.sql`)
            │
            ▼
Step 2: Domain Entities & Repositories (`domain/models.go`)
            │
            ▼
Step 3: Bun ORM Repository Implementation (`infra/repository.go`)
            │
            ▼
Step 4: Application Use Cases & Services (`app/service.go`)
            │
            ▼
Step 5: Domain Events & Triggers (`internal/core/event`)
            │
            ▼
Step 6: Hertz REST Handlers & Router (`api/handler.go` & `api/router.go`)
            │
            ▼
Step 7: Dependency Injection Container (`module.go`)
            │
            ▼
Step 8: Unit & Integration Tests (`*_test.go`)
```

---

## 2. Mandatory Architectural Constraints

1. **Strict Package Imports**:
   - `domain/` MAY NOT import `app/`, `infra/`, or `api/`.
   - `app/` imports `domain/` interfaces ONLY.
   - `infra/` implements `domain/` interfaces.
   - `api/` validates request DTOs and calls `app/` service methods.

2. **Error Responses**:
   - Handlers MUST return standard RFC 7807 problem details payloads via `coreerr.MapToProblemDetails(err, ctx.Request.URI().String())`.
   - Never write raw text error strings or return 500 without problem detail context.
