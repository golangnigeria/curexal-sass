# Schema-per-Tenant Multi-Tenancy Architecture

> **Purpose**: Deep-dive technical specification of Curexal's Schema-per-Tenant multi-tenant PostgreSQL database isolation model.  
> **Owner**: Lead Database Architect  
> **Status**: APPROVED / PRODUCTION READY  
> **Last Updated**: 2026-07-27  
> **Verification Criteria**: Verified against `internal/core/middleware/search_path_middleware.go` and `internal/modules/provisioning`.

---

## 1. Schema-per-Tenant Isolation Model

Curexal enforces physical schema isolation for every enterprise customer organization:

```text
PostgreSQL Database (`curexal`)
├── public                       # Universal shared tables (users, sessions, platform_settings, reference_*)
├── tenant_stjude                # St. Jude Diagnostic Center (LIMS, RIS, Patients, Encounters)
├── tenant_everight_labs         # Everight Pathology Network (LIMS, Billing, Inventory)
└── tenant_apex_pathology        # Apex LIMS (Specimens, Analyzers, Results)
```

---

## 2. Dynamic `search_path` Middleware Execution

When a request targeting a tenant endpoint is received:
1. `TenantResolverMiddleware` extracts the target tenant slug from the Host subdomain (e.g. `everight.curexal.com`), query parameter `?tenant=everight-labs`, or `X-Tenant-Slug` header.
2. `SearchPathMiddleware` executes a lightweight schema switch query before running handler queries:
   ```sql
   SET LOCAL search_path TO "tenant_everight_labs", "public";
   ```
3. All subsequent Bun ORM queries execute within the context of `"tenant_everight_labs"`, making cross-tenant data leaks impossible.

---

## 3. Tenant Schema Provisioning Runner

When a customer verification is approved (`organization.verified` event), the **Provisioning Workflow Engine** (`internal/modules/provisioning`):
1. Sanitizes the organization slug to generate schema name `tenant_<slug>`.
2. Executes DDL statement `CREATE SCHEMA IF NOT EXISTS "tenant_<slug>";`.
3. Runs all tenant migrations creating clinical tables (`patients`, `encounters`, `lab_orders`, `specimens`, `radiology_studies`, `invoices`).
4. Seeds default Casbin roles (`tenant_admin`, `pathologist`, `radiologist`, `lab_technician`, `receptionist`, `biller`).
5. Marks workspace status `ACTIVE` and fires `tenant.provisioned` event.
