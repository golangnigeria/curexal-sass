# CUREXAL — NON-CLINIC VERTICAL PURGE
# CLEAN WORKING DIRECTORY RESET

## ROLE

You are the Lead Software Architect and Repository Migration Engineer
responsible for performing a controlled reset of the Curexal codebase.

You are working on an existing repository.

Your objective is to remove obsolete/non-MVP healthcare verticals and leave
behind a clean, working Curexal Clinic MVP codebase.

The target is:

> PLATFORM KERNEL + CLINIC / OUTPATIENT LOGIC ONLY

Everything else must be removed unless it is genuinely required by the
platform kernel or Clinic MVP.

---

# 1. PRIMARY OBJECTIVE

Clean the entire repository so that the working implementation contains:

1. Platform kernel
2. Clinic / outpatient workflows
3. Shared infrastructure required by Clinic
4. Billing/payment primitives required by Clinic
5. Patient portal functionality required by Clinic
6. Authentication
7. Organization/facility/membership/RBAC
8. Patient identity / MPI foundation
9. Audit/security infrastructure
10. Required notification infrastructure

The repository must NOT contain active business implementations for:

- LIS (Laboratory Information System)
- RIS (Radiology / PACS / DICOM)
- HIS (Hospital Inpatient / Bed / Ward Management)
- Pharmacy inventory & dispensing (Keep clinic prescriptions)
- Diagnostic marketplace & referral networks

---

# 2. EXECUTION ORDER

PHASE 1: Repository inspection & inventory
PHASE 2: Database migration cleanup
PHASE 3: Backend Go modules purge
PHASE 4: Frontend React workspace purge
PHASE 5: Routing & RBAC purge
PHASE 6: Verification & Test validation
