# PLATFORM AUTHORIZATION & MANAGEMENT AUDIT REPORT

## Executive Summary
This document outlines the architectural refactoring of Platform Management into a dedicated bounded context (`internal/modules/platform`) separate from the Authentication module (`internal/modules/auth`).

---

## 1. Architectural Separation of Concerns

```text
BEFORE (Tightly Coupled):
Auth Module ──► Creates Platform Owner ──► Authenticates Users

AFTER (Enterprise Clean Architecture):
Platform Management (`internal/modules/platform`)
   ├── Manages `public.platform_settings`
   ├── Executes 10-Step Bootstrap Pipeline
   ├── Governs Platform Roles (`platform_owner`, `platform_admin`, `platform_support`, `platform_sales`)
   └── Enforces Platform Domain Scope (`domain: platform`)

Authentication Module (`internal/modules/auth`)
   └── Solely handles credential validation, password hashing, and session tokens.
```

---

## 2. Decoupled User Routing Logic (`GET /api/v1/auth/me`)

```text
GET /api/v1/auth/me
        │
        ▼
Is User a Platform User? (Has `platform_user_memberships` record)
        │
  ┌─────┴─────────────────────────────────────┐
  │                                           │
 YES                                          NO
  │                                           │
  ▼                                           ▼
Redirect to Platform Admin               Does user belong to active tenants?
(http://localhost:3001)                       │
                                        ┌─────┴──────────────────┐
                                        │                        │
                                       YES                       NO
                                        │                        │
                                        ▼                        ▼
                              Show Workspace Selector    Show Invitation Waiting Page
```

---

## 3. Platform Roles Matrix

| Role | Scope / Capabilities |
| :--- | :--- |
| **`platform_owner`** | Full System Governance, Bootstrap, Security Policy & Emergency Overrides |
| **`platform_admin`** | Organization Management, Verification Review Queue, Schema Provisioning |
| **`platform_support`** | Tenant Impersonation, Diagnostic Telemetry, Support Tickets |
| **`platform_sales`** | CRM Lead Pipeline, Demo Scheduling, Commercial Quotes |
| **`platform_finance`** | Subscriptions, Invoicing, ARR / MRR Reporting |
| **`platform_auditor`** | Security Audit Logs, HIPAA & ISO 27001 Compliance Telemetry |
| **`platform_devops`** | PostgreSQL Schema Provisioning Jobs, Redis Cache & NATS Telemetry |

---

## 4. Key Takeaway

Platform Administrators are **not** tenant members and must **never** be subjected to tenant middleware or `X-Tenant-ID` header checks. Platform routes (`/api/v1/platform/*`) run on an isolated Platform Middleware Chain.
