# Curexal — Final Frontend Cleanup & Clinic MVP Reset Report

## 1. Executive Summary

The frontend codebase across all 3 web applications (`apps/web-platform`, `apps/web-patient`, and `apps/web-public`) has been completely cleaned, unburdened of obsolete non-clinic vertical code, and reset to:

> **CUREXAL PLATFORM KERNEL + CLINIC OS MVP**

No non-clinic code is hidden via CSS (`display: none`), feature flags, or dead stub imports. All obsolete implementations have been removed from the repository.

---

## 2. Removed Artifacts

### 1. Legacy Product Stubs (`apps/web-platform/src/products/` - Entire Directory Deleted)
- `src/products/lis/` (LIS Specimen Verification & Accessioning stub)
- `src/products/ris/` (RIS Modality Queue stub)
- `src/products/pacs/` (PACS DICOM Viewer stub)
- `src/products/pharmacy/` (Pharmacy Dispensary stub)
- `src/products/hms/` (HMS Inpatient Documentation stub)
- `src/products/billing/` (Duplicate Billing stub)

### 2. Legacy Domain Definitions (`apps/web-platform/src/domains/` - Deleted Subdirectories)
- `src/domains/laboratory/` (Specimen, analyzer, and lab test result interfaces)
- `src/domains/radiology/` (Imaging, modality, and PACS interfaces)
- `src/domains/pharmacy/` (Drug inventory, FEFO batch, and dispensary interfaces)

### 3. Legacy Public Web Routes & Pages
- `apps/web-public/src/pages/patient-marketplace-page.tsx` (Healthcare marketplace directory search page)
- `/marketplace` route removed from `apps/web-public/src/router.tsx`

---

## 3. Preserved Architecture (Clinic OS MVP + Platform Kernel)

### Platform Kernel
- **Authentication & Security**: Multi-tenant session resolution, Argon2id auth, CSRF protection, Organization Guard, Platform Admin Guard.
- **Organization Governance**: Organization executive dashboard, branch management, staff memberships, and RBAC role assignments.
- **Master Catalogs**: Clinical consultation, nursing triage, and outpatient procedure fee catalogs.

### Clinic OS Operational Workspaces
- **Patient Reception & MPI Directory**: Walk-in registration, Master Patient Index lookup, appointment check-in (`/reception`).
- **Nursing Care Desk**: Vital signs intake, acuity scoring, and nurse-to-doctor queue handoff (`/care-desk`).
- **Doctor EMR Practice Canvas**: Consultation Room 1, active patient drawer, electronic SOAP notes, ICD-10 coding, digital prescriptions (`/clinical`).
- **Cashier Billing & POS**: Service fee settlement, multi-tender POS receipts (Cash, Card, Transfer), and payment audit trail (`/billing`).

### Patient Health Portal
- **Care Journey Continuum**: Intake $\to$ Triage $\to$ Doctor Consultation $\to$ E-Prescriptions $\to$ Billing & Receipts $\to$ Care Summary.
- **Self-Service Actions**: Prescription views, invoice & payment tracking, telehealth consultations.

---

## 4. Refactored Components & Routes

| Component | Path | Refactoring Summary |
| :--- | :--- | :--- |
| **Command Palette** | `apps/web-platform/src/features/search/command-palette.tsx` | Replaced legacy LIS specimen scanner with Care Desk Triage and Doctor Consultation Queue shortcuts. |
| **Organization Guard** | `apps/web-platform/src/components/auth/organization-guard.tsx` | Cleaned user role-to-workspace redirection to map only to Clinic OS modules. |
| **Branch Switcher** | `apps/web-platform/src/components/design-system/branch-switcher/index.tsx` | Defaulted fallback workspace route from `laboratory` to `clinical`/`dashboard`. |
| **Login Redirect** | `apps/web-platform/src/pages/auth/login/index.tsx` | Mapped login role redirects strictly to Clinic OS personas (`clinical`, `care-desk`, `reception`, `billing`, `dashboard`). |
| **Org Catalogs** | `apps/web-platform/src/pages/organization/catalogs/index.tsx` | Updated default catalog templates to outpatient consultations, triage, and procedures. |
| **Org Branches** | `apps/web-platform/src/pages/organization/branches/index.tsx` | Updated facility blueprints to Outpatient Clinic, Specialist Practice, and Telehealth Center. |
| **Org Roles** | `apps/web-platform/src/pages/organization/roles/index.tsx` | Standardized role cards to Owner, Admin, Doctor, Nurse, Receptionist, and Cashier. |
| **Patient Dashboard** | `apps/web-patient/src/pages/dashboard/index.tsx` | Updated care journey stages from `LAB_WORKLIST` to Clinic OS care continuum. |
| **Marketing Nav & Footer** | `apps/web-public/src/components/layouts/` | Aligned marketing links to Clinic OS, Doctor EMR, Patient Portal, and Billing POS. |

---

## 5. Build & Validation Results

| Test Suite / Build Target | Command | Result |
| :--- | :--- | :---: |
| **Web Platform Typecheck** | `bun x tsc --noEmit` (`apps/web-platform`) | **PASS (0 Errors)** |
| **Web Patient Typecheck** | `bun x tsc --noEmit` (`apps/web-patient`) | **PASS (0 Errors)** |
| **Web Public Typecheck** | `bun x tsc --noEmit` (`apps/web-public`) | **PASS (0 Errors)** |
| **Web Platform Production Build** | `bun run build` (`apps/web-platform`) | **BUILT (1m 36s)** |
| **Web Patient Production Build** | `bun run build` (`apps/web-patient`) | **BUILT (20.70s)** |
| **Web Public Production Build** | `bun run build` (`apps/web-public`) | **BUILT (46.87s)** |
| **Backend Go Test Suite** | `go test -v ./internal/testing/...` (`apps/api`) | **30 / 30 PASS (0 Failures)** |
