# Task Execution Checklists

> **Purpose**: Step-by-step checklists to verify completeness during backend feature implementation, frontend view development, refactoring, and release verification.  
> **Owner**: Lead Quality Assurance Architect  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. Backend Feature Checklist

- [ ] SQL DDL Migration created in `infra/migrations/`.
- [ ] Migration verified via `task db-migrate`.
- [ ] Bun ORM Struct Model defined in `domain/models.go`.
- [ ] Repository interface declared in `domain/models.go`.
- [ ] Bun Repository implemented in `infra/repository.go`.
- [ ] Use case service implemented in `app/service.go`.
- [ ] Hertz REST Handler implemented in `api/handler.go`.
- [ ] Route registered in `api/router.go`.
- [ ] Casbin RBAC permissions registered in `authz`.
- [ ] RFC7807 problem details returned on error.
- [ ] Go unit tests created in `app/*_test.go` and passing (`go test ./...`).

---

## 2. Frontend View Checklist

- [ ] API method added to `@curexal/api-sdk`.
- [ ] View component created in `apps/*/src/views/`.
- [ ] `LoadingSkeleton` component implemented.
- [ ] `ErrorState` component implemented (formatting RFC7807 with Retry button).
- [ ] `EmptyState` component implemented.
- [ ] Curexal theme colors (`#266210`, `#90B800`, `#00E1E1`, `#063B00`) applied.
- [ ] Zero mock data used.
- [ ] TypeScript strict mode (`tsc --noEmit`) passes cleanly.
- [ ] Application build (`npm run build`) passes cleanly.
