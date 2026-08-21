# EPIC-011: Enterprise Demo Booking, CRM Lead Pipeline & Compliance Onboarding

This document defines the specification for **EPIC-011**, replacing self-service public workspace registration with an Enterprise B2B Healthcare lifecycle (Book a Demo → CRM Sales Qualification → Demo & Proposal → Contract/Payment Approval → Platform Admin Conversion → Owner Invitation → 8-Step Setup Wizard → Platform Verification → Verified Tenant Schema Provisioning).

---

## 1. Executive Summary

Curexal is an enterprise Healthcare Operating System (comparable to Epic Systems, Oracle Health/Cerner, LabWare LIMS, STARLIMS, and InterSystems TrakCare).

Public visitors are **never allowed** to automatically create an operational organization, workspace, owner account, or database schema through a public self-service registration form.

Tenant PostgreSQL database schemas (`tenant_<slug>`) are provisioned **strictly after**:
1. Sales Lead Qualification & Contract/Payment Approval by Platform Admins.
2. Conversion of CRM Lead to Customer Organization (`status = INVITED`).
3. Owner Invitation Acceptance & Password Creation.
4. Completion of the 8-Step Setup & Regulatory Compliance Wizard.
5. Platform Verification Approval by Curexal Compliance Auditors.

---

## 2. Bounded Context Architecture

```text
Visitor (web-public)
   │
   ▼ POST /api/v1/demo/book
Lead Bounded Context (public.leads & public.demo_bookings)
   │
   ▼ Admin Reviews & Converts Lead in web-admin
Organization Bounded Context (public.organizations, status: INVITED)
   │
   ▼ Owner Accepts Token & Sets Password in web-workspace
Auth & Organization Context (Owner created, redirects to /workspace/setup)
   │
   ▼ Owner Submits 8-Step Setup Wizard
Organization Compliance Forms & Documents (public.organization_setup_profiles)
   │
   ▼ Admin Verification Approval in web-admin
Provisioning Bounded Context (Executes tenant_<slug> DDL migrations & seeders)
   │
   ▼
Organization & Workspace Activated (status: ACTIVE)
```

---

## 3. Key Endpoints

### Public Demo Request
- `POST /api/v1/demo/book`
- `GET /api/v1/demo/availability`

### Admin CRM Lead Management
- `GET /api/v1/leads`
- `PATCH /api/v1/leads/:id`
- `POST /api/v1/leads/:id/schedule-demo`
- `POST /api/v1/leads/:id/convert`

### Owner Invitation & Setup Wizard
- `POST /api/v1/invitations/accept`
- `POST /api/v1/setup/profile`
- `POST /api/v1/setup/documents`
- `POST /api/v1/setup/branches`
- `POST /api/v1/setup/departments`
- `POST /api/v1/setup/modules`
- `POST /api/v1/setup/invitations`
- `GET /api/v1/setup/status`

### Platform Verification & Schema Provisioning
- `POST /api/v1/verification/approve`
- `POST /api/v1/verification/reject`
- `POST /api/v1/verification/request-information`
