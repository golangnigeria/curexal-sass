# ENTERPRISE PLATFORM MANAGEMENT & BOOTSTRAP PIPELINE

## Executive Overview
In Curexal Healthcare Operating System, **Platform Management** (`internal/modules/platform`) is a dedicated bounded context decoupled from the Authentication module. Authentication solely authenticates credentials; Platform Management governs initial platform initialization, platform roles, Casbin policies, reference catalogs, marketplace configuration, commercial plans, and multi-tenant telemetry.

---

## Architecture: Decoupled Platform Bounded Context & Pipeline Interface

Platform Bootstrap executes **exclusively via infrastructure CLI** (`cmd/bootstrap` or `task bootstrap`). No REST endpoint (`/api/v1/platform/bootstrap`) is registered in production.

```text
                               CLI / Task Execution (`cmd/bootstrap`)
                                                  │
                                                  ▼
                               PlatformBootstrapPipeline (Modular Go Interface)
                                                  │
                                                  ├── Step 0: Validate System Prerequisites (DB, Redis, Env)
                                                  ├── Step 1: Create Platform Settings (`public.platform_settings`)
                                                  ├── Step 2: Create Platform Owner Account (`user_type: platform_owner`)
                                                  ├── Step 3: Seed Platform Governance Roles (7 Roles)
                                                  ├── Step 4: Seed Casbin Domain Policies (`domain: platform`)
                                                  ├── Step 5: Seed Reference Catalogs (Countries, Currencies, Timezones)
                                                  ├── Step 6: Seed Marketplace Module Configurations
                                                  ├── Step 7: Seed Commercial Plans (Starter, Professional, Enterprise)
                                                  ├── Step 8: Execute Platform Bootstrap Integrity Verification
                                                  └── Step 9: Mark platform_initialized = true
```

---

## 7 Platform Governance Roles

1. `platform_owner`: Full platform administration, system management, global settings.
2. `platform_admin`: Customer management, subscriptions, provisioning.
3. `platform_sales`: Enterprise CRM, lead qualification, demo scheduling.
4. `platform_support`: Customer support, tenant ticket resolution, impersonation.
5. `platform_finance`: Revenue cycle, invoices, MRR/ARR analytics.
6. `platform_auditor`: Immutable audit log review, compliance monitoring.
7. `platform_devops`: Infrastructure metrics, provisioning jobs, system telemetry.

---

## Middleware Architecture: Platform vs Tenant Isolation

Platform Administrators operate in global governance scope and **NEVER** require tenant context (`X-Tenant-Slug` or `X-Tenant-ID`).

### 1. Platform Middleware Chain
```text
Platform Router (/api/v1/platform/*)
       │
       ▼
RequirePlatformSession (Validates public.sessions & PlatformUserMembership)
       │
       ▼
RequirePlatformRBAC (Validates role privileges on domain: platform)
       │
       ▼
Platform API Handlers (Dashboard micro-aggregator, Infrastructure Telemetry)
```

### 2. Tenant Middleware Chain
```text
Tenant Router (/api/v1/tenant/*)
       │
       ▼
ResolveTenant (Extracts subdomain or X-Tenant-Slug header)
       │
       ▼
RequireTenantMembership (Validates membership in tenant_<slug>)
       │
       ▼
Tenant RBAC (Casbin domain: tenant_<slug>)
       │
       ▼
Tenant Operational Handlers (LIMS, RIS, EMR, Billing)
```

---

## Expanded Identity Contract (`GET /api/v1/auth/me`)

```json
{
  "user": {
    "id": "usr_powner_01H...",
    "email": "admin@curexal.com",
    "name": "System Architect",
    "is_verified": true
  },
  "session": {
    "mfa_enabled": false
  },
  "platform": {
    "is_platform_user": true,
    "roles": ["platform_owner"],
    "permissions": ["platform:*"]
  },
  "memberships": [],
  "default_destination": "platform"
}
```

---

## 10-Milestone Production Roadmap

1. **Milestone 1: Platform Bootstrap** – Modular 10-step pipeline interface, pre-flight checks, seed catalogs, CLI runner.
2. **Milestone 2: Platform Authentication** – Identity destination routing, platform session cookies, MFA support.
3. **Milestone 3: Platform Dashboard** – Dashboard micro-aggregator, health metrics, telemetry APIs.
4. **Milestone 4: CRM & Book Demo Pipeline** – Demo requests, sales qualification, lead pipeline.
5. **Milestone 5: Invitation & Activation** – Cryptographic token links, owner setup wizard, compliance doc upload.
6. **Milestone 6: Compliance Verification** – Accreditation review queue, document approval.
7. **Milestone 7: Tenant Provisioning** – PostgreSQL schema runner (`tenant_<slug>`), schema migrations, tenant activation.
8. **Milestone 8: Organization Workspace** – Subdomain workspace, LIMS/RIS/EMR module dashboards.
9. **Milestone 9: Patient Portal** – Patient results, appointment scheduling, medical history.
10. **Milestone 10: Marketplace** – Add-on module activation, commercial subscriptions, plan upgrades.
