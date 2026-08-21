# Curexal Healthcare Operating System — Master Documentation Portal

> **Purpose**: Central navigation portal for the Curexal V2 Enterprise Architecture, Bounded Contexts, Frontend Applications, Monorepo Packages, API Contracts, Operating Runbooks, and Production Audit Scorecards.  
> **Owner**: Principal Systems Architect  
> **Status**: APPROVED / PRODUCTION READY  
> **Last Updated**: 2026-07-27  
> **Verification Criteria**: Verified against live source code implementation across Go Hertz backend, Bun ORM PostgreSQL multi-tenancy, and React/Vite applications.

---

## 1. Executive Summary

Curexal is an enterprise-grade **Healthcare Operating System (HOS)** engineered as a **Modular Monolith** using **Go (Hertz Framework)**, **Bun ORM**, **PostgreSQL (Schema-per-Tenant isolation)**, **Redis**, and **Casbin RBAC**.

It delivers unified diagnostic, clinical, revenue cycle, compliance, and platform management operations across laboratories, radiology centers, hospitals, and multi-network diagnostic groups.

---

## 2. Master Documentation Directory Map

```text
docs/
├── architecture/          # Monorepo topology, C4 model, DDD context isolation, Multi-tenancy
├── backend/               # Go Hertz backend engine, Bun ORM schema, Domain events, Asynq workers
├── frontend/              # 4 Web Applications (web-public, web-workspace, web-patient, web-admin)
├── packages/              # 10 Shared monorepo npm packages (@curexal/api-sdk, @curexal/ui-core, etc.)
├── modules/               # 15 Bounded context functional specifications
├── api/                   # REST API specs, OpenAPI reference, RFC7807 problem details
├── ui/                    # Curexal Design System (#266210, #90B800, #00E1E1, #063B00)
├── deployment/            # Docker Compose, Taskfile automation, Production Kubernetes topologies
├── operations/            # Platform bootstrap pipeline, Backup, Restore, Break-Glass runbooks
├── environment/           # Environment profiles, Secrets management, Canonical .env.example
├── roadmap/               # Product evolution phases, Feature release plans
├── project/               # Living sprint status, Backlog, Known blockers, Technical debt, Changelog
├── audit/                 # Verified repository audits, Master dashboard, Security & Database scorecards
├── testing/               # QA strategy, Unit, Integration, E2E, and Performance testing
├── adr/                   # Architecture Decision Records (ADR-0001 through ADR-0007)
├── business/              # Commercial vision, Pricing strategy, Sales funnel, Compliance alignment
└── archive/               # Historical documentation snapshots
```

---

## 3. Quick Start & Operating Shortcuts

### Local Development Setup
```bash
# 1. Start Docker Infrastructure (PostgreSQL 16, Redis 7)
task infra-up

# 2. Run Database Schema Migrations & Reference Data Seeders
task db-migrate
task db-seed

# 3. Execute 10-Step Platform Initialization (Creates first Platform Owner)
task bootstrap

# 4. Start Go Backend Monolith (Port 8080)
task run

# 5. Start Frontend Applications (in separate terminal)
bun run dev
```

### Key Ports & Local URLs
- **Backend API Server**: `http://127.0.0.1:8080`
- **Platform Administration Portal (`web-admin`)**: `http://localhost:3003`
- **Organization Workspace (`web-workspace`)**: `http://localhost:3002`
- **Public Onboarding Portal (`web-public`)**: `http://localhost:3000`
- **Patient Portal (`web-patient`)**: `http://localhost:3001`
- **Swagger OpenAPI Docs**: `http://127.0.0.1:8080/swagger/index.html`

---

## 4. Architectural Guiding Principles

1. **Schema-per-Tenant Isolation**: Every enterprise customer organization gets a dedicated PostgreSQL database schema (`tenant_<slug>`). Cross-tenant SQL querying is physically impossible.
2. **Strict Platform Decoupling**: Platform Management & Bootstrap (`internal/modules/platform`) are isolated from Authentication (`internal/modules/auth`). Authentication only verifies credentials; Platform governs platform settings, roles, and telemetry.
3. **Zero Mock Data Policy**: Every frontend application consumes live APIs strictly through `@curexal/api-sdk`.
4. **Cookie-First Session Strategy**: HttpOnly, SameSite=Lax credentials cookies for 15-minute sessions with automated background refresh token rotation.
# CUREXAL_IDEA_DOCS 
