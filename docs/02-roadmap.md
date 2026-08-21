# Curexal V2 Product Roadmap

This document outlines the systematic, phase-by-phase evolution of the Curexal V2 platform, detailing both backend and frontend applications along with targeted launch schedules.

---

## Master Milestone Summary

| Phase | Milestone Name | Frontend Applications | Core Backend Capabilities | Target Launch Window |
| :--- | :--- | :--- | :--- | :--- |
| **Phase 1** | Platform Foundation & Marketplace Skeleton | `apps/portal`, `apps/admin`, `apps/patient` | Multi-Tenant Schemas, Auth, RBAC, Patient Directory | **Q3 2026** (Aug – Sep 2026) |
| **Phase 2** | Clinic Module & Referrals | `apps/workspace` (Clinic), `apps/patient` | Consultations, E-Prescriptions, B2B Referral Engine | **Q4 2026** (Oct – Nov 2026) |
| **Phase 3** | Laboratory Module & Barcoding | `apps/workspace` (LIMS), `apps/patient` | Phlebotomy Queue, Barcodes, Results, WhatsApp Dispatch | **Q1 2027** (Dec 2026 – Feb 2027) |
| **Phase 4** | Pharmacy, Inventory, & Payments | `apps/workspace` (POS/Pharmacy), `apps/patient` | Dispensing, Reagent Control, Settlement Wallets, Paystack | **Q2 2027** (Mar – May 2027) |
| **Phase 5** | Radiology & AI Platform | `apps/workspace` (RIS), `apps/patient` | PACS DICOM Listener, AI Report Translation & Copilots | **Q3 2027** (Jun – Aug 2027) |
| **Phase 6** | Enterprise Expansion & SDK | `apps/workspace` (HMS), `apps/admin` | Inpatient Management, FHIR/HL7, Marketplace SDK | **Q4 2027+** (Sep 2027 onwards) |

---

## Phase 1: Platform Foundation & Marketplace Skeleton
**Target Launch Date**: Q3 2026 (August – September 2026)

Establish the multi-tenant architecture, identity platform, public marketing portal, platform administration suite, and universal patient identity.

### Backend Infrastructure & Micro-Services
- **Identity Platform**: Universal SSO directory in `public` schema, JWT token rotation, session revocation, and password reset.
- **Tenant Provisioning Engine**: Dynamic PostgreSQL schema generator, automated migration runner, and initial seeders.
- **Role-Based Access Control (RBAC)**: Casbin RBAC integration with dynamic policy enforcement.
- **Organization & Branch Registry**: Organization profiles, multi-branch addresses, working hours, and branding configurations.
- **Marketplace Discovery Catalog**: Public search index, laboratory profiles, and test catalog price list templates.

### Frontend Applications & Packages
- **`apps/portal` (Public Portal)**:
  - Landing page, feature highlights, and interactive pricing calculator.
  - Self-service Organization Onboarding Registration & Step-by-Step Setup Wizard.
- **`apps/admin` (Platform Admin)**:
  - Multi-tenant workspace status monitor (real-time schema provisioning status).
  - Platform health, telemetry, global ICD-10 catalogs, and security audit log viewer.
- **`apps/patient` (Patient Portal)**:
  - Patient self-registration, account verification, and security profile manager.
- **`packages/ui` & `packages/api`**:
  - Shared Enterprise Design System (shadcn/ui + Tailwind CSS with strict accessibility & dark/light mode).
  - Generated TypeScript API bindings from backend Hertz OpenAPI specs.

---

## Phase 2: Clinic Module & Referrals
**Target Launch Date**: Q4 2026 (October – November 2026)

Deliver operational tools for outpatient clinics and connect them to the diagnostic referral marketplace.

### Backend Infrastructure & Micro-Services
- **Consultation Engine**: Doctor scheduling, consultation records, diagnoses (ICD-10/11), and SOAP notes.
- **E-Prescription Engine**: Digital prescription generator with dosage validation.
- **Referral Engine (B2B)**: State machine router for generating, dispatching, and tracking referral orders to target laboratories.

### Frontend Applications & Packages
- **`apps/workspace` (Clinic Workspace)**:
  - Clinician Dashboard: Active appointments queue, patient consultation interface with rich SOAP note editor.
  - Digital E-Prescription Form builder.
  - B2B Referral Dispatch Hub: Laboratory discovery search, pricing comparison matrix, turnaround time badges, and digital referral order placement.
  - Referral Order Tracking Board: Real-time status cards (Dispatched, Accepted, Sample Collected, Results Ready).
- **`apps/patient`**:
  - Digital referral notification cards with embedded appointment booking links.

---

## Phase 3: Laboratory Module & Barcoding
**Target Launch Date**: Q1 2027 (December 2026 – February 2027)

Support diagnostic laboratory operations and delivery of verified diagnostic results.

### Backend Infrastructure & Micro-Services
- **Phlebotomy & Sample Engine**: Specimen container management, barcode generator (Code128), and sample tracking.
- **Laboratory Workflow Engine**: Worklists, bench assignments, scientific validation rules, and consultant review paths.
- **Result Dispatcher**: PDF report generator (MinIO storage), automated email worker, and WhatsApp Business API dispatcher.
- **ISO 15189 Quality Control**: Standard deviation charting (Levey-Jennings) and equipment calibration logs.

### Frontend Applications & Packages
- **`apps/workspace` (Laboratory Workspace)**:
  - Reception Queue & Patient Check-in.
  - Phlebotomy Station: Specimen logging and thermal barcode label printing interface.
  - Laboratory Bench View: Analytical worklists, result entry forms, and critical delta alert badges.
  - Scientific Validation & Consultant Review Workspace: Dual-approval interface for report authorization and live PDF preview.
- **`apps/patient`**:
  - Diagnostic Results Hub: Secure PDF report viewer, result history timeline, and direct download links.

---

## Phase 4: Pharmacy, Inventory, & Payments
**Target Launch Date**: Q2 2027 (March – May 2027)

Add pharmacy dispensing workflows, multi-warehouse inventory logistics, and dynamic referral commission clearing.

### Backend Infrastructure & Micro-Services
- **Pharmacy Dispensing**: Electronic prescription matching, stock deduction, and dispensing verification.
- **Inventory Logistics**: Reagent and consumable stock tracking, expiry monitoring, purchase orders, and stock alerts.
- **B2B Settlement Engine**: Financial ledger managing split-payment commission transactions, platform fees, and virtual wallets.
- **Payment Gateways**: Paystack and Flutterwave API integration for online checkouts and POS card transactions.

### Frontend Applications & Packages
- **`apps/workspace` (Pharmacy & Financial Workspace)**:
  - Pharmacy Station: E-prescription queue, drug match checker, and dispensing logger.
  - Inventory & Reagent Manager: Stock level indicators, expiry warning boards, purchase order manager.
  - Cashier POS & Financial Workspace: Invoice builder, cash/card POS checkout modal, settlement wallet overview, and commission payout manager.
- **`apps/patient`**:
  - Integrated Paystack/Flutterwave online checkout modal for test and prescription bill payments.

---

## Phase 5: Radiology & AI Platform
**Target Launch Date**: Q3 2027 (June – August 2027)

Integrate medical imaging pipelines and clinical AI assistance services.

### Backend Infrastructure & Micro-Services
- **PACS & DICOM Infrastructure**: DICOM listener, PACS metadata storage, and image streaming proxy.
- **Radiologist Interface Engine**: Audio dictation processor, structured template engine, and digital signature verifier.
- **AI Assistive Platform**: Generative AI for report translation into local vernacular (Pidgin, Yoruba, Hausa, Igbo) and radiologist report draft generation.

### Frontend Applications & Packages
- **`apps/workspace` (Radiology Workspace)**:
  - RIS Worklist: Modality schedules (X-Ray, CT, MRI, Ultrasound).
  - Integrated Web-Based DICOM Image Viewer (pan, zoom, measurements, windowing).
  - Radiologist Dictation Studio: Voice-to-text dictation panel and structured report template editor.
- **`apps/patient`**:
  - AI Patient Explainer: Tap-to-translate clinical reports into easy-to-understand vernacular summaries.

---

## Phase 6: Enterprise Expansion & SDK
**Target Launch Date**: Q4 2027+ (September 2027 onwards)

Transition Curexal V2 into a full enterprise-grade Healthcare Operating System.

### Backend Infrastructure & Micro-Services
- **Hospital Management System (HMS)**: ICU, theatre, inpatient admissions, bed allocation, and nursing care records.
- **Healthcare Interoperability**: HL7 v2/v3, FHIR resources, and ASTM analyzer connectors.
- **Developer Platform & SDK**: Sandbox workspaces, developer keys, webhook delivery engine, and marketplace plugin SDK.

### Frontend Applications & Packages
- **`apps/workspace` (Hospital Inpatient Workspace)**:
  - Ward & Bed Management Grid, Nursing Flowsheet, Operating Theatre Scheduler.
- **`apps/admin` (Developer Portal)**:
  - Developer dashboard for API key management, webhook event inspection, and third-party marketplace plugin publishing.
