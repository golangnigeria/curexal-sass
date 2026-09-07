# Day 1: Frontend-to-Backend Testing & Verification Walkthrough

The **Day 1 Production Foundation** is now locked down and fully testable end-to-end from Frontend (`apps/web-platform`) to Backend (`apps/api`), adhering strictly to our company design principles (rich aesthetics, modern typography, glassmorphism, responsive micro-animations, and Zero-Trust clinical governance).

---

## 1. Day 1 Features Implemented & Integrated

| Feature | Scope | Backend Implementation | Frontend UI Component |
| :--- | :--- | :--- | :--- |
| **Feature #1** | Secure Authentication Engine & Deterministic Branch Selection | Argon2id/Bcrypt credential verification, 5-minute branch selection tokens, session family rotation. | [LoginPage](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/apps/web-platform/src/pages/auth/login/index.tsx) with **Day 1 Persona Matrix** and **Branch Selection Modal**. |
| **Feature #2** | Canonical Role-Based Access Control | Migration 70 seeded 9 canonical roles (`owner`, `clinic_doctor`, `clinic_nurse`, etc.) with `RequirePermission()` guards. | Canonical route guards, dynamic permissions query, and role-specific workspace redirection. |
| **Feature #3** | Organization & Clinic Management | Migration 72 seeded 3 operational branches; `RequireActiveBranch` enforces HTTP 428 zero-fallback and HTTP 403 staff isolation. | [BranchSwitcher](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/apps/web-platform/src/components/design-system/branch-switcher/index.tsx) in Topbar with real-time `/auth/switch-branch` session re-issuance. |
| **Feature #24** | Immutable Audit Ledger & HIPAA Accounting | Migration 71 trigger `trg_prevent_audit_tampering`, SHA-256 hash chains, patient disclosure accounting. | [AuditLogsPage](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/apps/web-platform/src/pages/platform/audit/index.tsx) with SHA-256 cryptographic verification proofs & break-glass badges. |

---

## 2. Interactive Day 1 Test Matrix

The login screen now includes an interactive **Day 1 Verification Matrix** panel on the left side of the credentials form. Clicking any persona automatically populates the email and password, displays the scenario's expected flow, and enables 1-click test sign in.

All test accounts share the default password: **`password`**.

```
┌────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                       DAY 1 TEST MATRIX PERSONAS                                       │
├──────────────────────────┬─────────────────────────────┬───────────────────┬───────────────────────────┤
│ Persona                  │ Email                       │ Role & Branches   │ Expected Day 1 Flow       │
├──────────────────────────┼─────────────────────────────┼───────────────────┼───────────────────────────┤
│ 👑 Clinic Owner          │ owner@curexal.com           │ owner             │ Multi-Branch Modal ->     │
│                          │                             │ (3 Branches)      │ Topbar Branch Switcher    │
├──────────────────────────┼─────────────────────────────┼───────────────────┼───────────────────────────┤
│ 🩺 Dr. Emeka Okonkwo     │ dr.emeka@curexal.com        │ clinic_doctor     │ Multi-Branch Modal ->     │
│    (Attending Physician) │                             │ (HO-01 & IKJ-02)  │ Clinical EMR Workspace    │
├──────────────────────────┼─────────────────────────────┼───────────────────┼───────────────────────────┤
│ 👩‍⚕️ Dr. Sarah Alabi       │ dr.sarah@curexal.com        │ clinic_doctor     │ Direct Bypass (No modal)  │
│    (Consultant)          │                             │ (HO-01 only)      │ -> Clinical EMR Workspace │
├──────────────────────────┼─────────────────────────────┼───────────────────┼───────────────────────────┤
│ 📋 Nurse Chioma Eze      │ nurse.chioma@curexal.com    │ clinic_nurse      │ Direct Bypass (No modal)  │
│    (Triage Lead)         │                             │ (HO-01 only)      │ -> Care Desk / Triage MPI │
├──────────────────────────┼─────────────────────────────┼───────────────────┼───────────────────────────┤
│ 🚫 Provisional Staff     │ unassigned.staff@curexal.com│ clinic_nurse      │ HTTP 403 Zero-Trust       │
│    (Unassigned)          │                             │ (0 Branches)      │ Guidance Resolution Modal │
├──────────────────────────┼─────────────────────────────┼───────────────────┼───────────────────────────┤
│ 🛡️ Platform Super Admin │ admin@curexal.com           │ super_admin       │ Direct entrance to        │
│                          │                             │ (Platform)        │ Immutable Audit Ledger    │
└──────────────────────────┴─────────────────────────────┴───────────────────┴───────────────────────────┘
```

---

## 3. Step-by-Step Testing Guide

### Test Scenario A: Multi-Branch Doctor Login (`dr.emeka@curexal.com`)
1. Open the web platform login page (`/login`).
2. In the **Day 1 Verification Personas** grid, click **Dr. Emeka Okonkwo** (or enter `dr.emeka@curexal.com` / `password`).
3. Click **"Sign In to Secure Workspace"**.
4. **Verified Result**:
   - The backend returns `requireBranchSelection: true` with a 5-minute JWT `selectionToken`.
   - The **Branch Resolution Modal** slides into view with glassmorphic cards for:
     - `HO-01` — **Curexal Clinic Main Campus** (HQ • Victoria Island)
     - `IKJ-02` — **Curexal Ikeja Specialty Center** (Ikeja)
   - Select **Curexal Ikeja Specialty Center** and click **"Enter Clinical Workspace"**.
   - The app exchanges the token via `POST /api/v1/auth/select-branch` and routes directly to `/ikeja-specialty/clinical`.

### Test Scenario B: Single-Branch Direct Login (`dr.sarah@curexal.com`)
1. Select **Dr. Sarah Alabi** from the persona grid.
2. Click **"Sign In"**.
3. **Verified Result**:
   - Backend detects exactly 1 assigned branch (`HO-01`).
   - With zero modal interruption, directly issues the access token and navigates directly to `/curexal-clinic/clinical`.

### Test Scenario C: Zero-Trust HIPAA Unassigned Guard (`unassigned.staff@curexal.com`)
1. Select **Provisional Care Staff** from the persona grid.
2. Click **"Sign In"**.
3. **Verified Result**:
   - The Go backend blocks authorization returning `HTTP 403` with error code `UNASSIGNED_FACILITY_BRANCH`.
   - The **Facility Location Assignment Required** modal pops up with HIPAA multi-center compliance instructions.
   - Click **"Check Assignment Status"** to re-poll after an administrator assigns branches.

### Test Scenario D: Active Branch Switching in Workspace Topbar
1. Log in as **Curexal Clinic Owner** (`owner@curexal.com`).
2. In the Branch Resolution Modal, pick **Curexal Clinic Main Campus** (`HO-01`).
3. Inside the workspace, locate the **Branch Switcher** in the top navigation bar.
4. Click the dropdown to see all 3 branches (`HO-01`, `IKJ-02`, `LEK-03`).
5. Select **Curexal Lekki Diagnostic & Emergency**.
6. **Verified Result**:
   - Frontend issues `POST /api/v1/auth/switch-branch` with `{ branchId: "..." }`.
   - Backend re-issues session cookies and access token.
   - TanStack query cache is invalidated, and the browser cleanly navigates to `/lekki-diagnostic/dashboard`.

### Test Scenario E: Feature #24 Immutable Audit Ledger
1. Log in as **Platform Super Admin** (`admin@curexal.com`).
2. Navigate to **Platform Console > Audit Logs** (`/platform/audit`).
3. **Verified Result**:
   - The audit log shows immutable facility and administrative records.
   - Click **"Inspect"** on any log event.
   - The inspect dialog displays the **SHA-256 Cryptographic Integrity Proof** with the verified record hash and previous block hash.

---

## 4. Key Files Created & Modified
- [000072_seed_day1_test_branches_and_matrix.sql](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/apps/api/database/platform/migrations/000072_seed_day1_test_branches_and_matrix.sql): Seeded 3 clinic branches, test accounts, canonical RBAC memberships, multi-branch links, and genesis audit block.
- [login/index.tsx](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/apps/web-platform/src/pages/auth/login/index.tsx): Complete redesign of login and branch selection experiences adhering to company aesthetics, framer-motion micro-interactions, and 1-click persona testing.
- [contracts/src/index.ts](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/packages/contracts/src/index.ts): Enhanced `AuditLog` contract with Feature #24 compliance fields (`facilityBranchId`, `patientId`, `isBreakGlass`, `recordHash`, `prevRecordHash`).
- [platform/audit/index.tsx](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/apps/web-platform/src/pages/platform/audit/index.tsx): Integrated cryptographic SHA-256 proof badge and break-glass indicators in inspect dialog.
