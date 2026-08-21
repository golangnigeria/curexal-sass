# 15-Step Feature Development Pipeline & Workflow

> **Purpose**: Mandatory 15-step pipeline for implementing any new feature, bounded context, or API capability across the Curexal V2 monorepo.  
> **Owner**: Principal Systems Architect  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. Master Development Pipeline

```text
 1. Business Requirement Analysis
                │
                ▼
 2. Architecture & Domain Design
                │
                ▼
 3. Database DDL Migration & Schema
                │
                ▼
 4. Domain Entities & Value Objects (`domain/`)
                │
                ▼
 5. Bun ORM Repository Implementation (`infra/`)
                │
                ▼
 6. Application Use Cases & Services (`app/`)
                │
                ▼
 7. Domain Events & Notification Triggers (`event/`)
                │
                ▼
 8. Hertz REST Handlers & Route Wiring (`api/`)
                │
                ▼
 9. Casbin RBAC Permission Registration (`authz/`)
                │
                ▼
10. Go Unit & Integration Tests (`*_test.go`)
                │
                ▼
11. API SDK Client Update (`@curexal/api-sdk`)
                │
                ▼
12. React Frontend Application View (`apps/*`)
                │
                ▼
13. End-to-End Flow Verification (`UI -> SDK -> Hertz -> DB`)
                │
                ▼
14. System Documentation Sync (`docs/`)
                │
                ▼
15. 13 Quality Production Gates Sign-Off
```

---

## 2. Definition of Done (DoD) per Phase

| Pipeline Phase | Mandatory Deliverables | Verification Command |
| :--- | :--- | :--- |
| **Phase 3: Database Migration** | `.sql` migration file in `infra/migrations/`, Bun ORM struct tags, foreign keys, indexes. | `task db-migrate` |
| **Phase 4-6: Backend DDD** | `domain/models.go`, `infra/repository.go`, `app/service.go`. Clean layer separation. | `go build ./...` |
| **Phase 8-9: REST API & Authz** | Hertz handler, RFC7807 error responses, Casbin permission rule. | `go test -v ./internal/modules/...` |
| **Phase 11: API SDK Client** | Exported client method in `@curexal/api-sdk` (`packages/api-sdk/src/`). | `bun run build` |
| **Phase 12: React Frontend** | Screen/view in `apps/*`, skeleton loaders, error states, empty states, zero mock data. | `npm run build` |
| **Phase 14: Documentation** | Update `docs/api/API_REFERENCE.md`, `docs/project/CHANGELOG.md`, `FEATURE_MATRIX.md`. | Verified doc links |
