# Curexal V2 Epics & Execution Roadmap

This document outlines the master Epics for Curexal V2, grouping capabilities into deliverable engineering work packages.

---

## Master Epic Matrix

| Epic ID | Title | Bounded Context | Core Deliverables | Target Phase |
| :--- | :--- | :--- | :--- | :--- |
| **EPIC-001** | Identity Platform | Identity | Universal SSO, User Directory, JWT Session Rotation, Password Reset, Email Verification. | Phase 1 |
| **EPIC-002** | Multi-Tenant Provisioning | Tenant | Dynamic Schema Creator, Migration Runner, Organization & Branch Settings, Branding. | Phase 1 |
| **EPIC-003** | RBAC & Security Audit | Identity & Audit | Casbin RBAC Integration, Context Middleware, Hash-Chained Audit Logs. | Phase 1 |
| **EPIC-004** | Marketplace Discovery Registry | Marketplace | Public Provider Index, Search Engine, Public Test Price List Catalogs. | Phase 1 |
| **EPIC-005** | Master Patient Index & Portal | Patient | Patient Demographics, Phonetic Matching, Medical Timeline, Patient Portal Self-Service. | Phase 1 |
| **EPIC-006** | Clinic Consultation & E-Prescriptions | Clinic | Doctor Appointments Queue, SOAP Notes, Diagnosis (ICD-10), E-Prescription Form. | Phase 2 |
| **EPIC-007** | B2B Referral State Router | Referral | Digital Referral Dispatcher, Order State Machine, Lab Inbox, Patient Notification. | Phase 2 |
| **EPIC-008** | Laboratory LIMS & Barcode Engine | Lab LIMS | Specimen Intake, Code128 Thermal Labeling, Bench Worklists, Consultant Sign-off. | Phase 3 |
| **EPIC-009** | Curexal Edge Instrument Sync | Lab LIMS | On-Premises Go Edge Agent, ASTM/HL7 Parser, SQLite Local Cache, mTLS Sync. | Phase 3 |
| **EPIC-010** | Financial Settlement & Wallets | Financial | Double-Entry Wallet Ledgers, Automated Split Clearing, Paystack/Flutterwave Webhooks. | Phase 4 |
| **EPIC-011** | Enterprise Demo Booking & Compliance Onboarding | Lead & Tenant | CRM Lead Pipeline, Demo Scheduler, Admin Customer Conversion, 8-Step Setup Wizard, Verified Provisioning. | Phase 1 |

---

## Granular Epic Descriptions

### EPIC-001: Identity Platform
- **Scope**: Implements universal user authentication across the global `public` schema directory.
- **Key Features**:
  - `POST /api/v1/auth/register` (Register user/patient)
  - `POST /api/v1/auth/login` (Authenticate credentials, create DB session, set HTTP-Only cookies)
  - `POST /api/v1/auth/refresh` (Rotate refresh token and issue new 15-min access JWT)
  - `POST /api/v1/auth/logout` (Revoke session token)

### EPIC-002: Multi-Tenant Provisioning
- **Scope**: Orchestrates organization onboarding and dynamic PostgreSQL schema creation.
- **Key Features**:
  - `POST /api/v1/workspaces` (Accept workspace registration)
  - Asynchronous schema runner (`CREATE SCHEMA tenant_<slug>;` + DDL migrations + seeders)
  - `GET /api/v1/workspaces/{id}/status` (Provisioning progress monitor)
