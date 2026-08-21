# Curexal V2 Dependency Architecture & DAG

This document defines the Capability Dependency Graph (Directed Acyclic Graph - DAG) ensuring modules are built in the exact architectural order.

---

## Capability Dependency Graph

```
[Identity & Auth (EPIC-001)] 
         │
         ▼
[Tenant Provisioning (EPIC-002)] ──► [RBAC & Audit (EPIC-003)]
         │
         ├────────────────────────────────────────┐
         ▼                                        ▼
[Marketplace Registry (EPIC-004)]        [Patient Identity (EPIC-005)]
         │                                        │
         └───────────────────┬────────────────────┘
                             │
                             ▼
               [Clinic Consultations (EPIC-006)]
                             │
                             ▼
               [B2B Referral Router (EPIC-007)]
                             │
                             ▼
               [Laboratory LIMS (EPIC-008)] ──► [Edge Instrument Sync (EPIC-009)]
                             │
                             ▼
               [Financial Settlement Ledger (EPIC-010)]
```

---

## Detailed Dependency Matrix

| Module / Feature | Direct Dependencies (Prerequisites) | Blocking Downstream Features |
| :--- | :--- | :--- |
| **Identity & Auth** | *None (Base Platform)* | Tenant Provisioning, Patient Identity, All Modules |
| **Tenant Provisioning** | Identity & Auth | Workspace Context Resolver, RBAC, All Tenant Schemas |
| **RBAC & Authorization** | Tenant Provisioning, Identity | Protected Endpoints, Casbin Middleware |
| **Marketplace Registry** | Tenant Provisioning | B2B Referral Router, Provider Discovery |
| **Patient Identity** | Identity & Auth | Clinic Encounters, Lab Intake, Results Portal |
| **Clinic Consultations** | Patient Identity, RBAC | Referral Dispatcher, E-Prescriptions |
| **B2B Referral Router** | Clinic Consultations, Marketplace Registry | Laboratory Intake, Result Delivery |
| **Laboratory LIMS** | B2B Referral Router, Patient Identity | Analyzer Edge Sync, Financial Settlement |
| **Financial Ledger** | Laboratory LIMS, B2B Referral Router | Wallet Payouts, Paystack Webhooks |
