# Monorepo Repository Map & Directory Inventory

> **Purpose**: Complete directory map and dependency inventory for all applications, packages, core libraries, CLI runners, and configuration manifests in the Curexal V2 monorepo.  
> **Owner**: Principal Systems Architect  
> **Status**: APPROVED / VERIFIED  
> **Last Updated**: 2026-07-27  
> **Verification Criteria**: Audited from ground truth source code directory listings.

---

## 1. Monorepo Topology Overview

```text
curexalV2/
├── apps/               # Frontend Single-Page Applications (Vite + React)
├── packages/           # Shared Monorepo npm Packages
├── internal/           # Go Backend Monolith Source Code
│   ├── core/           # Shared System Utilities, Database, Security, Middleware
│   └── modules/        # Bounded Context Domain Modules
├── cmd/                # Go Application Binaries & CLI Entrypoints
├── docs/               # Architecture & Operations Documentation System
├── scripts/            # Database Seeder & Utility Scripts
├── Taskfile.yml        # Task Runner Orchestration File
├── turbo.json          # Turborepo Build Pipeline Config
├── package.json        # Root Bun Workspace Manifest
├── go.mod              # Backend Go Module Definition
└── docker-compose.yml  # Local Development Infrastructure Manifest
```

---

## 2. Applications Directory Inventory (`/apps`)

| Directory | Name | Purpose | Primary Tech Stack | Status |
| :--- | :--- | :--- | :--- | :--- |
| `apps/web-admin` | `@curexal/web-admin` | Platform Admin Portal control center | React 18, Vite, Tailwind CSS, `@curexal/api-sdk` | ✅ Active |
| `apps/web-workspace` | `@curexal/web-workspace` | Multi-tenant organization workspace (LIMS/RIS/EMR) | React 18, Vite, Tailwind CSS, `@curexal/api-sdk` | ✅ Active |
| `apps/web-patient` | `@curexal/web-patient` | Patient portal for lab results & appointment booking | React 18, Vite, Tailwind CSS, `@curexal/api-sdk` | ✅ Active |
| `apps/web-public` | `@curexal/web-public` | Public marketing landing page & onboarding registration | React 18, Vite, Tailwind CSS, `@curexal/api-sdk` | ✅ Active |

---

## 3. Shared Packages Inventory (`/packages`)

| Package Directory | Package Name | Purpose | Consumers | Status |
| :--- | :--- | :--- | :--- | :--- |
| `packages/api-sdk` | `@curexal/api-sdk` | Typed HTTP client for Hertz backend endpoints | All `apps/*` | ✅ Active |
| `packages/auth` | `@curexal/auth` | Auth context provider & session token management | All `apps/*` | ✅ Active |
| `packages/authorization` | `@curexal/authorization` | Client-side Casbin permission hooks & guard components | `web-workspace`, `web-admin` | ✅ Active |
| `packages/types` | `@curexal/types` | Shared TypeScript domain interfaces & DTOs | All `apps/*`, `packages/*` | ✅ Active |
| `packages/ui-core` | `@curexal/ui-core` | Base UI components (Buttons, Inputs, Modals, Cards) | All `apps/*` | ✅ Active |
| `packages/ui-healthcare` | `@curexal/ui-healthcare` | Medical UI components (Specimen tracker, DICOM viewer) | `web-workspace`, `web-patient` | ✅ Active |
| `packages/design-tokens` | `@curexal/design-tokens` | Tailwind color tokens (`#266210`, `#90B800`, `#00E1E1`) | All `apps/*`, `packages/*` | ✅ Active |
| `packages/i18n` | `@curexal/i18n` | Multi-language translation dictionaries | All `apps/*` | ✅ Active |
| `packages/state` | `@curexal/state` | Global Zustand state stores | `web-workspace`, `web-admin` | ✅ Active |
| `packages/config` | `@curexal/config` | Shared tsconfig, ESLint, and Vite configurations | Monorepo root | ✅ Active |

---

## 4. Backend Bounded Context Modules (`/internal/modules`)

| Module Directory | Bounded Context | Responsibilities | Status |
| :--- | :--- | :--- | :--- |
| `internal/modules/platform` | Platform Management | 10-step bootstrap pipeline, platform settings, governance roles | ✅ Active |
| `internal/modules/auth` | Identity & Authentication | User credentials, password reset, session rotation, email verification | ✅ Active |
| `internal/modules/organization` | Workspace Management | Tenant registration, compliance verification, workspace subscription | ✅ Active |
| `internal/modules/provisioning` | Schema Runner | PostgreSQL DDL schema creation (`tenant_<slug>`), role seeding | ✅ Active |
| `internal/modules/reference` | Catalogs & Terminology | Reference data (Countries, Currencies, Timezones, LOINC, SNOMED) | ✅ Active |
| `internal/modules/lead` | Commercial CRM | Demo requests, sales qualification pipeline, customer conversion | ✅ Active |
| `internal/modules/authz` | Authorization | Casbin domain policy enforcement & RBAC rules | ✅ Active |
| `internal/modules/notification` | Notification Engine | Email notifications (Resend API) and verification tokens | ✅ Active |
