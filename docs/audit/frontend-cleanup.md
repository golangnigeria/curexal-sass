# Curexal — Frontend Cleanup & Clinic MVP Reset Audit

## 1. Executive Summary

This forensic audit establishes the target state for the Curexal frontend across all three web applications:
- `apps/web-platform` (Organization & Facility Operational Portal)
- `apps/web-patient` (Patient Self-Service Health Portal)
- `apps/web-public` (Public Marketing & Product Showcase)

The architectural goal is a strict, intentional reset to:
> **CUREXAL PLATFORM KERNEL + CLINIC OS MVP**

All residual non-clinic vertical artifacts (LIS, RIS, PACS, DICOM, HIS Inpatient, Wards, Beds, Pharmacy Dispensary/FEFO Inventory, Public Healthcare Marketplace) are identified for complete removal rather than visual hiding or conditional flagging.

---

## 2. Frontend Architecture Mapping

```text
apps/
├── web-platform/
│   ├── src/
│   │   ├── api/               # API clients, hooks & query keys
│   │   ├── components/        # Design system, layout, auth guards, data-display
│   │   ├── domains/           # Domain-specific models & types [CLEANUP TARGET]
│   │   ├── features/          # Feature components (appointments, patients, triage, billing)
│   │   ├── pages/             # Route page views (auth, organization, platform, workspace)
│   │   ├── products/          # Stub product engines [DELETE TARGET]
│   │   ├── workspaces/        # Workspace route wrappers (reception, care-desk, clinical, billing)
│   │   └── router.tsx         # Canonical route tree
├── web-patient/
│   ├── src/
│   │   ├── api/               # Patient self-service API client
│   │   ├── components/        # Patient UI design system
│   │   ├── pages/             # Patient login, dashboard, consultations
│   │   └── router.tsx         # Patient portal routes
└── web-public/
    ├── src/
    │   ├── components/        # Marketing navbar, hero, layouts
    │   ├── pages/             # Solutions, pricing, waitlist, book-demo
    │   └── router.tsx         # Public marketing route tree
```

---

## 3. Comprehensive Frontend Inventory Table

| Area | Component / File | Purpose | Vertical | Used By | Action |
| :--- | :--- | :--- | :--- | :--- | :---: |
| **Products** | `src/products/lis/` | LIS Accessioning stub | LIS / Lab | None (Orphaned) | **DELETE** |
| **Products** | `src/products/ris/` | RIS Modality queue stub | RIS / Radiology | None (Orphaned) | **DELETE** |
| **Products** | `src/products/pacs/` | PACS DICOM viewer stub | PACS / Imaging | None (Orphaned) | **DELETE** |
| **Products** | `src/products/pharmacy/` | Pharmacy dispensing stub | Pharmacy | None (Orphaned) | **DELETE** |
| **Products** | `src/products/hms/` | HMS inpatient documentation stub | HIS / Inpatient | `features/encounters` re-export | **DELETE** |
| **Products** | `src/products/billing/` | POS Billing stub | Billing / Clinic | None (Orphaned) | **DELETE** |
| **Domains** | `src/domains/laboratory/` | Specimen & lab result interfaces | LIS / Lab | None (Orphaned) | **DELETE** |
| **Domains** | `src/domains/radiology/` | Radiology modality interfaces | RIS / Radiology | None (Orphaned) | **DELETE** |
| **Domains** | `src/domains/pharmacy/` | Pharmacy inventory/FEFO interfaces | Pharmacy | None (Orphaned) | **DELETE** |
| **Domains** | `src/domains/clinical/` | Clinical consultation & SOAP models | Clinic OS | `features/encounters` | **KEEP** |
| **Domains** | `src/domains/patient/` | Patient demographics & MPI models | Platform Kernel | `features/patients` | **KEEP** |
| **Domains** | `src/domains/billing/` | Invoicing & POS payment models | Clinic OS / Billing | `features/billing` | **KEEP** |
| **Domains** | `src/domains/organization/`| Organization & Branch models | Platform Kernel | `pages/organization` | **KEEP** |
| **Features** | `src/features/encounters/` | Encounter barrel re-exports | Clinic OS | Router / Workspaces | **REFACTOR** |
| **Features** | `src/features/search/` | Global Command Palette | Platform Kernel | App Layout | **REFACTOR** |
| **Pages (Platform)** | `pages/platform/marketplace/` | Platform capability marketplace | Platform Kernel | Super Admin | **REFACTOR** |
| **Pages (Org)** | `pages/organization/catalogs/`| Org catalog management | Platform Kernel | Org Admin | **REFACTOR** |
| **Pages (Org)** | `pages/organization/branches/`| Branch facility blueprints | Platform Kernel | Org Admin | **REFACTOR** |
| **Pages (Org)** | `pages/organization/roles/`| RBAC role configuration | Platform Kernel | Org Admin | **REFACTOR** |
| **Pages (Org)** | `pages/organization/dashboard/`| Org executive dashboard | Platform Kernel | Org Owner | **REFACTOR** |
| **Pages (Workspace)** | `pages/workspace/dashboard/` | Facility operations dashboard | Clinic OS | Branch Staff | **REFACTOR** |
| **Pages (Workspace)** | `pages/workspace/clinical/` | Outpatient Doctor Practice Canvas | Clinic OS | Attending Doctor | **REFACTOR** |
| **Pages (Workspace)** | `pages/workspace/care-desk/` | Triage & Nursing Desk | Clinic OS | Triage Nurse | **KEEP** |
| **Pages (Workspace)** | `pages/workspace/reception/` | Patient Reception & MPI Intake | Clinic OS | Receptionist | **KEEP** |
| **Pages (Workspace)** | `pages/workspace/billing/` | Cashier POS Register | Clinic OS | Cashier | **KEEP** |
| **Web Patient** | `apps/web-patient/src/pages/dashboard/` | Patient Care Journey milestones | Clinic OS | Patient | **REFACTOR** |
| **Web Public** | `apps/web-public/src/pages/patient-marketplace-page.tsx` | Healthcare marketplace directory | Marketplace | Public Router | **DELETE** |
| **Web Public** | `apps/web-public/src/components/layouts/marketing-navbar.tsx` | Solutions navigation dropdown | Public Marketing | Public Web | **REFACTOR** |

---

## 4. Audit Findings: Non-Clinic Residual Terms

1. **`src/products/`**: Contains orphaned legacy modules (`lis`, `ris`, `pacs`, `pharmacy`, `hms`, `billing`). Must be deleted in full.
2. **`src/domains/`**: Contains orphaned domain type folders (`laboratory`, `pharmacy`, `radiology`). Must be deleted in full.
3. **`apps/web-public/src/pages/patient-marketplace-page.tsx`**: Public marketplace search page. Must be removed along with `/marketplace` route.
4. **Command Palette (`src/features/search/command-palette.tsx`)**: Had hardcoded reference to `/${activeBranchSlug}/laboratory`.
5. **Branch Switcher & Organization Guard**: Default fallback module was hardcoded to `"laboratory"`. Must be updated to `"clinical"`.
6. **Patient Care Journey (`apps/web-patient/src/pages/dashboard/index.tsx`)**: Contained `LAB_WORKLIST` stage code. Must be updated to canonical Clinic OS stages.

---

## 5. Execution Plan

1. **Delete Obsolete Directories**:
   - `apps/web-platform/src/products/` (Entire folder)
   - `apps/web-platform/src/domains/laboratory/`
   - `apps/web-platform/src/domains/pharmacy/`
   - `apps/web-platform/src/domains/radiology/`
   - `apps/web-public/src/pages/patient-marketplace-page.tsx`
2. **Refactor Features & Barrels**:
   - Clean `apps/web-platform/src/features/encounters/index.ts`.
   - Update `command-palette.tsx` search routes to point only to Clinic OS workspaces.
3. **Refactor Organization & Platform Views**:
   - Clean default module fallbacks in `organization-guard.tsx`, `branch-switcher/index.tsx`, `login/index.tsx`.
   - Clean mock data in `pages/organization/catalogs/index.tsx`, `branches/index.tsx`, `roles/index.tsx`.
4. **Refactor Public Web**:
   - Update `apps/web-public/src/router.tsx` (remove `/marketplace` route).
   - Update `marketing-navbar.tsx` and `marketing-footer.tsx` (remove `/marketplace` link, focus on Clinic OS).
5. **Refactor Patient Portal**:
   - Update `apps/web-patient/src/pages/dashboard/index.tsx` care journey milestones to canonical Clinic stages.
6. **Build & Test Validation**:
   - `bun x tsc --noEmit` across all 3 web apps.
   - `bun run build` in all 3 web apps.
   - `go test -v ./internal/testing/...` in `apps/api`.
