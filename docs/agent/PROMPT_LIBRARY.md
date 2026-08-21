# Enterprise AI Agent Prompt Library

> **Purpose**: Reusable battle-tested prompt templates for AI-assisted development across backend, frontend, database, refactoring, security, and release tasks.  
> **Owner**: Principal Systems Architect  
> **Status**: APPROVED / REUSABLE  
> **Last Updated**: 2026-07-27

---

## 1. Prompt 01: Create New Backend Bounded Context Module

```text
Build a new backend bounded context module named "[module_name]" in `internal/modules/[module_name]`.
Requirements:
1. Follow strict DDD layers: `domain/`, `app/`, `infra/`, `api/`, `module.go`.
2. Write Bun ORM models and repository implementations.
3. Use 26-character ULID primary keys and soft delete timestamps.
4. Implement Hertz REST handlers returning RFC7807 problem details payloads on error.
5. Register Casbin RBAC permissions for roles: [role_1], [role_2].
6. Write unit tests for application services in `app/*_test.go`.
7. Update `docs/api/API_REFERENCE.md` and `docs/project/CHANGELOG.md`.
```

---

## 2. Prompt 02: Wire Frontend View to Backend via `@curexal/api-sdk`

```text
Wire the React view component `[ViewName].tsx` in `apps/[app_name]/src/views/` to live backend APIs.
Requirements:
1. Use `@curexal/api-sdk` exclusively. Zero mock data or static JSON arrays.
2. Handle all required UI states: LoadingSkeleton, ErrorState (formatting RFC7807 problem details with Retry action), EmptyState, and Data State.
3. Apply Curexal Enterprise theme tokens (#266210, #90B800, #00E1E1, #063B00).
4. Update progress matrix in `docs/implementation/PLATFORM_UI_IMPLEMENTATION.md`.
```

---

## 3. Prompt 03: 8-Step Log-Driven Bug Fix

```text
Investigate and fix the following bug report:
[Insert Error Message or Steps to Reproduce]

Follow the 8-step debugging protocol in `docs/agent/DEBUGGING_GUIDE.md`:
1. Collect full logs and stack trace.
2. Identify verified root cause in source code.
3. Implement root-cause fix without symptom hacking or disabling security.
4. Run `go test ./...` or `npm run build` to verify regression-free success.
5. Update `docs/project/CHANGELOG.md`.
```
