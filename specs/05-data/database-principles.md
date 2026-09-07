# CUREXAL PLATFORM — DATABASE PRINCIPLES SPECIFICATION
**Document**: `specs/05-data/database-principles.md`  
**Status**: APPROVED BASELINE  

---

## 1. Core Principles

1. **PostgreSQL as SSOT**: PostgreSQL is the authoritative transactional data store.
2. **Referential Integrity**: All relationships must be enforced via foreign keys, unique constraints, and check constraints at the database level.
3. **Tenant Scoping**: Every tenant-owned table must include an `organization_id UUID NOT NULL` indexed foreign key, and facility-specific tables must include `facility_branch_id UUID NOT NULL`.
4. **Deterministic Timestamps & Audit Fields**:
   - `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
   - `updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
   - `created_by UUID REFERENCES identity.users(id)`
5. **Durable Clinical Records**: Clinical encounters, SOAP notes, and diagnoses are never hard-deleted. Soft deletes use `deleted_at TIMESTAMPTZ` with immutable audit log entries.
