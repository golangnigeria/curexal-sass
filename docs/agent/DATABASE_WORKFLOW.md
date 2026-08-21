# Database-First Development Workflow

> **Purpose**: Standard Operating Procedure for database migrations, Bun ORM model definitions, indexing, foreign keys, and Schema-per-Tenant isolation.  
> **Owner**: Lead Database Architect  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. Database-First Development Pipeline

```text
Step 1: Write DDL SQL Migration File (`infra/migrations/00X_create_table.sql`)
            │
            ▼
Step 2: Run Migration Verification (`task db-migrate`)
            │
            ▼
Step 3: Define Bun ORM Struct Model (`domain/models.go`)
            │
            ▼
Step 4: Implement Bun Repository Methods (`infra/repository.go`)
            │
            ▼
Step 5: Write Repository Integration Test (`infra/repository_test.go`)
```

---

## 2. Mandatory Database Rules

1. **Schema Target Choice**:
   - `public`: Shared universal data (`users`, `sessions`, `platform_settings`, `commercial_plans`).
   - `tenant_<slug>`: Workspace clinical data (`patients`, `encounters`, `lab_orders`, `specimens`, `invoices`).
2. **Primary Keys & ULID**: Always use 26-character ULID strings as primary keys (`bun:",pk,type:varchar(26)"`) for sorting efficiency and collision resistance.
3. **Timestamps & Audit Fields**: Every table MUST include `created_at`, `updated_at`, and `deleted_at` (soft delete timestamp).
4. **Indexes & Foreign Keys**: Foreign keys MUST be explicitly defined with `ON DELETE RESTRICT` or `CASCADE` where appropriate. Compound indexes MUST be created for frequently queried filter columns.
