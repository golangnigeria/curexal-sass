# Feature Request Prompt Template

> **Purpose**: Standardized template to be used whenever requesting new features or capability implementations from AI coding agents.

```markdown
# Feature Request: [Feature Name]

## 1. Business Objective & Goal
- **Problem Statement**: Brief description of the problem or healthcare workflow to address.
- **Target Bounded Context**: `internal/modules/[module_name]` (e.g. `lims`, `ris`, `emr`, `billing`, `platform`).
- **Target Frontend Application**: `apps/[web-admin | web-workspace | web-patient | web-public]`.

## 2. Actors & Permissions
- **User Roles**: `platform_owner`, `tenant_admin`, `pathologist`, `radiologist`, `lab_technician`, `biller`, etc.
- **Casbin Permission Strings**: `[domain]:[resource]:[action]` (e.g. `lims:specimens:create`).

## 3. Database Schema Changes
- **Target Schema**: `tenant_<slug>` (clinical/workspace) OR `public` (platform/global).
- **New Tables / Columns**: Specify table names, column types, foreign keys, and indexes.

## 4. Backend Deliverables
- [ ] Database Migration file in `internal/modules/[module_name]/infra/migrations/`.
- [ ] Domain Models & Repository Interface in `internal/modules/[module_name]/domain/`.
- [ ] Bun ORM Repository Implementation in `internal/modules/[module_name]/infra/`.
- [ ] Application Use Cases in `internal/modules/[module_name]/app/`.
- [ ] Hertz REST Handlers & Route Registration in `internal/modules/[module_name]/api/`.
- [ ] Unit & Integration Tests in `internal/modules/[module_name]/app/*_test.go`.

## 5. API SDK & Frontend Deliverables
- [ ] Exported API client method in `@curexal/api-sdk`.
- [ ] React View Component in `apps/[app_name]/src/views/`.
- [ ] Skeleton Loaders, Error States (RFC7807), and Empty States.

## 6. Definition of Done
- [ ] `go build ./...` succeeds cleanly.
- [ ] `npm run build` succeeds cleanly.
- [ ] Documentation updated in `docs/` and `docs/project/CHANGELOG.md`.
```
