# Backend Architecture & Bounded Contexts Specification

> **Purpose**: Detailed technical reference for all 15 backend bounded contexts in the Curexal Go Hertz Monolith.  
> **Owner**: Principal Backend Engineer  
> **Status**: APPROVED / VERIFIED  
> **Last Updated**: 2026-07-27  
> **Verification Criteria**: Audited directly from `internal/modules/*` and `internal/core/*`.

---

## 1. Backend Bounded Contexts Inventory

```text
internal/
├── core/
│   ├── audit/          # Core Audit Logger & Diff Calculator
│   ├── authz/          # Casbin Authorization Engine & Enforcer Provider
│   ├── config/         # Environment Configuration Loader
│   ├── database/       # Bun ORM Connection Pool & Migration Runner
│   ├── email/          # Resend Email Client
│   ├── errors/         # RFC 7807 Problem Details Error Framework
│   ├── event/          # In-Memory & Redis Domain Event Bus
│   ├── health/         # Infrastructure Liveness & Readiness Probes
│   ├── middleware/     # Platform & Tenant Middleware Chains
│   ├── primitive/      # ULID Generator & Clock Interfaces
│   ├── rest/           # Hertz JSON Helpers & Swagger/Redoc Handlers
│   └── security/       # Argon2id Password Hashing Engine
└── modules/
    ├── platform/       # 10-Step Bootstrap Pipeline, Settings, Admin Telemetry
    ├── auth/           # SSO Registration, Credentials Login, Session Rotation, Password Reset
    ├── organization/   # Organization Registration, Verification Queue, Subscription
    ├── provisioning/   # Asynchronous PostgreSQL Schema Runner Pipeline
    ├── reference/      # Terminology & Reference Catalogs (LOINC, SNOMED, Countries)
    ├── lead/           # Commercial CRM Lead Pipeline & Demo Scheduling
    ├── authz/          # Domain Casbin Role & Permission Registries
    └── notification/   # Email Notification Dispatch Handlers
```

---

## 2. Context Technical Breakdown

### 1. Platform Management (`internal/modules/platform`)
- **Entities**: `PlatformSettings`, `PlatformUserMembership`, `MarketplaceSetting`, `CommercialPlan`.
- **Pipeline Runner**: `PlatformBootstrapPipeline` executing 10 modular steps in `cmd/bootstrap`.
- **API Endpoints**: `GET /api/v1/platform/dashboard`, `GET /api/v1/platform/health`.

### 2. Identity & Authentication (`internal/modules/auth`)
- **Entities**: `User`, `UserProfile`, `Session`, `EmailVerificationToken`, `PasswordResetToken`, `LoginAttempt`.
- **Session Strategy**: HttpOnly Cookie (`curexal_access_token` 15m, `curexal_refresh_token` 7d) with Argon2id password hashing.
- **Identity Routing**: `GET /api/v1/auth/me` returns `user_type` and `default_destination` (`"platform"`, `"tenant"`, `"invitation_pending"`).

### 3. Organization Workspace (`internal/modules/organization`)
- **Entities**: `Organization`, `Workspace`, `WorkspaceBranch`, `WorkspaceSubscription`, `ComplianceVerification`.
- **Workflow**: Registration -> Compliance document upload -> Admin verification approval -> Tenant schema runner trigger.

### 4. Provisioning Runner (`internal/modules/provisioning`)
- **Entities**: `ProvisioningJob`, `ProvisioningStepLog`.
- **Runner**: Executes `CREATE SCHEMA IF NOT EXISTS "tenant_<slug>"`, tenant DDL migrations, and Casbin role seeders.
