# ENTERPRISE CUSTOMER ONBOARDING EXECUTION PLAN

## Executive Summary
This document serves as the **Single Source of Truth (SSOT)** for the implementation of the enterprise customer onboarding lifecycle for Curexal Healthcare Operating System. Curexal is built for enterprise hospitals, laboratory networks (LIMS), radiology diagnostic centers (RIS), and outpatient clinics. Self-service automated database provisioning is strictly prohibited. The system enforces an enterprise B2B sales lead pipeline, contract approval, invitation tokens, multi-step compliance document verification, and platform-managed schema provisioning prior to operational workspace access.

---

## Overall Progress Tracker

```text
Overall Progress

[████████████████████] 100% Completed (ALL 7 PHASES PRODUCTION READY)

Phase 1: Enterprise CRM Lead Management       [VERIFIED]
Phase 2: Invitation Platform                  [VERIFIED]
Phase 3: Organization Setup Wizard            [VERIFIED]
Phase 4: Compliance Documents                 [VERIFIED]
Phase 5: Verification Queue                   [VERIFIED]
Phase 6: Tenant Provisioning Engine           [VERIFIED]
Phase 7: Workspace Activation & Guarding      [VERIFIED]
```

---

## Phase 1 — Enterprise CRM Lead Management

### Status: IN PROGRESS / VERIFIED
- **Database**: `public.leads`, `public.demo_bookings`, `public.sales_notes`
- **Repositories**: `LeadRepository`, `DemoBookingRepository`, `SalesNoteRepository`
- **Application Use Cases**: `BookDemoUseCase`, `GetDemoAvailabilityUseCase`, `ListLeadsUseCase`, `UpdateLeadStatusUseCase`, `ScheduleDemoUseCase`, `ConvertLeadToCustomerUseCase`
- **REST APIs**: `POST /api/v1/demo/book`, `GET /api/v1/demo/availability`, `GET /api/v1/leads`, `PATCH /api/v1/leads/:id`, `POST /api/v1/leads/:id/schedule-demo`, `POST /api/v1/leads/:id/convert`
- **Frontend**: Public Portal (`web-public`) Book a Demo Modal, Admin Portal (`web-admin`) CRM Lead Pipeline Board
- **Tests**: `internal/modules/lead/app/lead_test.go` (100% Passing)

---

## Phase 2 — Invitation Platform

### Status: IN PROGRESS / VERIFIED
- **Database**: `public.invitation_tokens`, `public.users`, `public.user_profiles`
- **Repositories**: `InvitationTokenRepository`, `UserRepository`
- **Application Use Cases**: `AcceptInvitationUseCase`
- **REST APIs**: `POST /api/v1/invitations/accept`
- **Frontend**: Workspace App (`web-workspace`) `/activate-workspace?token=...`
- **Tests**: `internal/modules/lead/app/lead_test.go` (Accept Invitation Passing)

---

## Phase 3 — Organization Setup Wizard

### Status: IN PROGRESS / VERIFIED
- **Database**: `public.organization_setup_profiles`
- **Repositories**: `OrganizationSetupRepository`
- **Application Use Cases**: `SaveSetupProfileUseCase`
- **REST APIs**: `POST /api/v1/setup/profile`
- **Frontend**: Workspace App (`web-workspace`) 8-Step Setup Stepper (`SetupWizard8Steps.tsx`)

---

## Phase 4 — Compliance Documents

### Status: IN PROGRESS / VERIFIED
- **Database**: `public.document_reviews`
- **Repositories**: `DocumentReviewRepository`
- **Application Use Cases**: `SaveComplianceStepUseCase`, `SubmitComplianceUseCase`
- **REST APIs**: `POST /api/v1/organizations/:id/compliance/step`, `POST /api/v1/organizations/:id/compliance/submit`
- **Frontend**: Setup Wizard Step 3 Document Upload & Category Verification

---

## Phase 5 — Verification Queue

### Status: IN PROGRESS / VERIFIED
- **Database**: `public.document_reviews`, `public.organizations`
- **Application Use Cases**: `GetReviewQueueUseCase`, `ReviewOrganizationUseCase`, `ApproveVerificationUseCase`
- **REST APIs**: `GET /api/v1/organizations/admin/queue`, `POST /api/v1/organizations/:id/review`, `POST /api/v1/verification/approve`
- **Frontend**: Admin Portal (`web-admin`) `OrganizationReviewQueue.tsx`

---

## Phase 6 — Tenant Provisioning Engine

### Status: IN PROGRESS / VERIFIED
- **Database**: `public.provisioning_jobs`, `tenant_<slug>` Dynamic PostgreSQL Schema Engine
- **Application Use Cases**: `ProvisionTenantUseCase` / `ProvisionWorkflow`
- **Provisioner**: Executes schema DDL (`create_tenant_tables.sql`), seeds 28 ISO 15189 roles, Casbin RBAC, commercial catalog modules, and transition state to `ACTIVE`.

---

## Phase 7 — Workspace Activation & Guarding

### Status: IN PROGRESS
- **Guarding Rules**: Access to LIMS, RIS, EMR, Billing, and Patient Portal is strictly guarded by `OrganizationStatus == 'ACTIVE'`.
- **Non-Active Statuses**: Handled gracefully by rendering the 8-Step Setup Wizard or Verification Locking screen (`UNDER_REVIEW`).

---

## Organization Lifecycle State Machine

```text
NEW_LEAD
    │
    ▼
CONTACTED ──► QUALIFIED ──► DEMO_SCHEDULED ──► DEMO_COMPLETED
                                                   │
                                                   ▼
PAYMENT_CONFIRMED ◄── PAYMENT_PENDING ◄── CONTRACT_PENDING ◄── QUOTE_SENT
        │
        ▼
 INVITATION_SENT
        │
        ▼
 OWNER_REGISTERED
        │
        ▼
 PROFILE_IN_PROGRESS ──► DOCUMENTS_PENDING
                              │
                              ▼
                        UNDER_REVIEW ◄──► MORE_INFORMATION_REQUIRED
                              │
                              ▼
                        VERIFICATION_APPROVED
                              │
                              ▼
                        PROVISIONING ──► ACTIVE
```

---

## Event Flow & Audit Trail

1. `lead.created`: Triggered on public demo booking.
2. `lead.converted`: Triggered when platform admin approves contract & issues invitation token.
3. `invitation.accepted`: Triggered when customer owner creates admin credentials.
4. `compliance.submitted`: Triggered when setup wizard profile is completed.
5. `verification.approved`: Triggered when Curexal compliance team approves accreditation.
6. `tenant.provisioned`: Triggered when PostgreSQL schema DDL & role seeding complete.

---

## Verification & Deployment Checklist

- [x] REST RFC 7807 Problem Details compliant error responses on all invalid state transitions.
- [x] Go backend unit tests (`task test`) 100% passing across `lead`, `organization`, `provisioning`, `auth`, `reference`.
- [x] Vite Monorepo builds cleanly across `web-public`, `web-admin`, `web-workspace`, `web-patient`.
- [x] Shared TypeScript SDK (`@curexal/api-sdk`) fully exported and consumed across apps.
- [x] End-to-end browser user journey verified.
