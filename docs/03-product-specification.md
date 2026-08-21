# Curexal V2 Product Specification & Capability Map

This document establishes the functional requirements, capability mappings, frontend application bindings, and release windows for all modules of the Curexal V2 platform.

---

## 1. Capability Map & Monorepo Architecture

Every service in Curexal is organized around a business capability. The following map outlines the bounded contexts and their primary frontend application hosts:

```
Healthcare Platform Capability Matrix
│
├── Identity & Tenant Context [Phase 1 | Q3 2026] ──► Host: apps/portal, apps/admin, apps/patient
├── Patient EHR Context        [Phase 1 | Q3 2026] ──► Host: apps/workspace, apps/patient
├── Clinic Context             [Phase 2 | Q4 2026] ──► Host: apps/workspace (Clinic)
├── Referral Marketplace       [Phase 2 | Q4 2026] ──► Host: apps/workspace, apps/patient
├── Laboratory (LIMS) Context  [Phase 3 | Q1 2027] ──► Host: apps/workspace (LIMS), apps/patient
├── Pharmacy & POS Context     [Phase 4 | Q2 2027] ──► Host: apps/workspace (POS/Pharmacy)
├── Inventory Logistics        [Phase 4 | Q2 2027] ──► Host: apps/workspace
├── Billing & Ledger Wallets   [Phase 4 | Q2 2027] ──► Host: apps/workspace, apps/patient
├── Radiology (RIS/PACS)       [Phase 5 | Q3 2027] ──► Host: apps/workspace (RIS)
├── AI Copilots & Translation  [Phase 5 | Q3 2027] ──► Host: apps/workspace, apps/patient
└── Hospital (HMS) & SDK       [Phase 6 | Q4 2027+] ─► Host: apps/workspace (HMS), apps/admin
```

---

## 2. Core Module Specifications

### 2.1. Identity & Multi-Tenant Platform Context
- **Release Phase**: Phase 1
- **Target Launch Window**: Q3 2026 (August – September 2026)
- **Frontend Applications**: `apps/portal` (Public Marketing & Onboarding Wizard), `apps/admin` (Platform Management), `apps/patient` (Patient Auth)
- **Shared Packages**: `packages/ui` (Auth Layouts, Form Inputs), `packages/api` (Auth Client Hooks)

#### Purpose
Manages universal single sign-on (SSO), dynamic PostgreSQL schema provisioning, organization workspace registration, session rotation, branding customization, and platform security administration.

#### User Roles & Permissions
- `SuperAdmin` (`admin:all`): Full platform administration.
- `OrgOwner` (`workspace:manage`): Full workspace & organization settings management.
- `StaffUser` (`workspace:read`): Workspace staff member access.

#### Frontend UI Layouts & Component Tree
1. **Public Portal Onboarding (`apps/portal/src/pages/onboarding`)**:
   - Facility Type Selection Cards (Clinic, Lab, Pharmacy, Radiology, Hospital).
   - Organization Slug Availability Checker.
   - Owner Profile & Credentials Step Form.
   - Real-time Provisioning Status Progress Bar.
2. **Platform Admin Studio (`apps/admin/src/pages/tenants`)**:
   - Tenant Directory Table with real-time status badges (Provisioning, Ready, Suspended).
   - Database Migration Execution Console.
   - Global ICD-10 and LOINC Catalog Manager.

---

### 2.2. Patient Identity & EHR Module
- **Release Phase**: Phase 1 & Phase 2
- **Target Launch Window**: Q3 2026 – Q4 2026
- **Frontend Applications**: `apps/workspace` (Staff Directory Search), `apps/patient` (Patient Self-Service Portal)
- **Shared Packages**: `packages/ui` (Patient Avatar, Record Cards), `packages/api` (Patient Hooks)

#### Purpose
Provides a single source of truth for patient clinical identity across the healthcare ecosystem. Enables search-matching, tracking of medical history, allergies, and diagnostic timelines.

#### User Roles & Permissions
- `patients:create`: Register new patient identity.
- `patients:read`: View patient cards and clinical record timelines.
- `patients:write`: Update allergies, immunizations, and chronic conditions.

#### Frontend UI Layouts & Component Tree
1. **Clinic Patient Directory (`apps/workspace/src/features/patients`)**:
   - Fast Search Bar (Phonetic matching by Name, Phone, MRN).
   - Patient Profile Header Card (Age, Gender, Blood Group, Vital Warning Badges).
   - Interactive Clinical Record Timeline (Visits, Diagnoses, Prescriptions, Test Results).
2. **Patient Health Portal (`apps/patient/src/pages/medical-history`)**:
   - Personal Medical Record Summary.
   - Diagnostic Test Results History with downloadable PDF reports.

---

### 2.3. Referral Module (B2B Marketplace Engine)
- **Release Phase**: Phase 2
- **Target Launch Window**: Q4 2026 (October – November 2026)
- **Frontend Applications**: `apps/workspace` (Clinic Dispatcher & Lab Inbox), `apps/patient` (Referral Notifications)
- **Shared Packages**: `packages/ui` (Status Badges, Price Comparison Tables)

#### Purpose
Enables outpatient clinics to digitize patient referrals to third-party diagnostic laboratories and imaging centers, tracks sample collections, updates order states in real-time, and transfers validated results back to the referring physician.

#### Workflows
```
[Clinic: Create Referral] 
         │
         ▼
[Marketplace: Select Lab & Test Catalog]
         │
         ▼
[Referral Engine: Dispatch Order Event] ──► (Notify Lab & Patient via SMS/WhatsApp)
         │
         ▼
[Lab: Accept / Reject Referral] 
         │
         ├─► Reject ──► [Referral Terminated]
         │
         └─► Accept ──► [Specimen Collected (Barcode)] ──► [Scientific Result Released]
```

#### Frontend UI Layouts & Component Tree
1. **Clinic Referral Dispatcher (`apps/workspace/src/features/referrals/create`)**:
   - Laboratory Search Filter (Filter by Proximity, Rating, Pricing, Turnaround Time).
   - Test Catalog Selection List with dynamic total price calculator.
   - Clinical History & Attachment Uploader.
2. **Laboratory Referral Inbox (`apps/workspace/src/features/referrals/inbox`)**:
   - Incoming Referral Kanban Board (Pending Acceptance, Accepted, In-Progress, Completed).
   - Quick Action Dialogs (Accept with Turnaround Estimate, Reject with Reason Code).
3. **Patient Referral Hub (`apps/patient/src/pages/referrals`)**:
   - Referral Order Details & Laboratory Directions.
   - Integrated Paystack Checkout Modal for prepaying test fees.

#### API Endpoints & NATS Events
- APIs: `POST /api/v1/referrals`, `GET /api/v1/referrals/incoming`, `POST /api/v1/referrals/{id}/accept`
- Events: `referral.dispatched`, `referral.accepted`, `referral.completed`

---

### 2.4. Laboratory Module (LIMS) & Barcoding
- **Release Phase**: Phase 3
- **Target Launch Window**: Q1 2027 (December 2026 – February 2027)
- **Frontend Applications**: `apps/workspace` (LIMS Station), `apps/patient` (Results Download)

#### Purpose
Automates diagnostic lab operations, including phlebotomy specimen intake, thermal barcode generation, analyzer worklist routing, scientific result entry, consultant review, and WhatsApp report delivery.

#### Frontend UI Layouts & Component Tree
1. **Phlebotomy Station (`apps/workspace/src/features/lims/phlebotomy`)**:
   - Specimen Intake Queue with barcode printing triggers (Code128 thermal labels).
   - Container Verification Checklist (EDTA, Serum, Sodium Fluoride tubes).
2. **Analytical Worklists (`apps/workspace/src/features/lims/worklists`)**:
   - Bench Worklist (Chemistry, Hematology, Microbiology).
   - Parameter Result Entry Table with auto-calculation fields and Delta Checks.
3. **Consultant Review Studio (`apps/workspace/src/features/lims/review`)**:
   - Dual-authorization sign-off panel.
   - Live PDF Report Previewer with dynamic header/footer branding templates.

---

### 2.5. Pharmacy & Financial Settlement Module
- **Release Phase**: Phase 4
- **Target Launch Window**: Q2 2027 (March – May 2027)
- **Frontend Applications**: `apps/workspace` (Pharmacy POS & Financial Workspace), `apps/patient` (Online Paystack Checkout)

#### Purpose
Supports electronic prescription matching, drug dispensing verification, multi-warehouse stock management, and automated B2B split-payment referral commission settlements via Paystack and Flutterwave.

#### Frontend UI Layouts & Component Tree
1. **Pharmacy Dispensing Desk (`apps/workspace/src/features/pharmacy/dispense`)**:
   - E-Prescription Search & Prescription Matcher.
   - Batch Selection & Expiry Expiration Warning Checks.
2. **POS Cashier Station (`apps/workspace/src/features/billing/pos`)**:
   - Cash / Card / Transfer Payment Checkout Modal.
   - Instant Digital Receipt Generator.
3. **Settlement Wallet Overview (`apps/workspace/src/features/finance/wallet`)**:
   - Platform Commission Balance Card.
   - B2B Referral Payout Ledger Table with automated Paystack bank account transfer button.

---

### 2.6. Radiology Module (RIS/PACS) & AI Services
- **Release Phase**: Phase 5
- **Target Launch Window**: Q3 2027 (June – August 2027)
- **Frontend Applications**: `apps/workspace` (RIS Studio), `apps/patient` (AI Vernacular Explainer)

#### Purpose
Integrates medical imaging pipelines (DICOM/PACS viewer), radiologist voice dictation, and generative AI services for automated patient report translations into local African languages (Pidgin, Yoruba, Hausa, Igbo).

#### Frontend UI Layouts & Component Tree
1. **RIS Dashboard (`apps/workspace/src/features/ris`)**:
   - Modality Worklist (X-Ray, CT, MRI, Ultrasound).
   - Embedded Cornerstone.js / OHIF Web-based DICOM Viewer with pan, zoom, and measurement tools.
   - Voice Dictation & Audio Structuring Panel.
2. **Patient AI Explainer (`apps/patient/src/pages/reports/explain`)**:
   - Interactive Report Card with "Translate to Pidgin/Yoruba/Hausa/Igbo" toggle button powered by Curexal AI Copilot.
