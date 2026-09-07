# CUREXAL PLATFORM — ARCHITECTURE SPECIFICATION
**Document**: `specs/01-platform/architecture.md`  
**Status**: APPROVED ARCHITECTURAL BASELINE  

---

## 1. System Overview & The Kernel/Module Boundary

Curexal is architected as two decoupled layers:
1. **The Platform Kernel**: Reusable, vertical-agnostic healthcare infrastructure.
2. **Business Modules**: Domain-specific healthcare verticals running on top of the kernel.

```text
                                CUREXAL PLATFORM
                                       │
                ┌──────────────────────┴──────────────────────┐
                │                                             │
         PLATFORM KERNEL                               BUSINESS MODULES
  - Identity & Authentication                   - Clinic (HMS / Outpatient EMR) [MVP]
  - Multi-Tenant & Branch Context               - Laboratory (LIS / Diagnostics) [Future]
  - Permissions & Casbin RBAC                   - Pharmacy (Dispensary & Stock) [Future]
  - Master Patient Index (MPI)                  - Radiology (RIS / PACS) [Future]
  - Billing & Invoicing Engine                  - Hospital (Inpatient / ADT) [Future]
  - Audit Logging & Observability
  - Notification Channels
  - Commercial Entitlements
```

---

## 2. Invariants

1. **Kernel Blindness**: The Platform Kernel has zero awareness of outpatient clinic consultation forms, SOAP notes, or lab test analyzers. It manages identity, tenancy, patient identifiers, billing primitives, and authorization.
2. **Module Autonomy**: Business modules (Clinic, LIS, Pharmacy) own their domain logic, workflows, and private data schemas.
3. **No Cross-Module Database Access**: A business module must never directly query or write to another business module's private database tables.
4. **Contract-Driven Communication**: Inter-module communication occurs exclusively via shared kernel contracts, application service interfaces, or domain events.

---

## 3. Technology Stack

- **Backend**: Go (Echo Framework) + PostgreSQL 16 (Relational SSOT) + Casbin v2 (Multi-Tenant Authorization Engine) + Redis (Session & Query Cache).
- **Frontend**: React 18 + TypeScript + Vite + TailwindCSS + React Query + React Router v6.
- **Monorepo**: Turborepo + Bun.
