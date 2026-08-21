# Curexal Healthcare Organization Verification & Document Submission Architecture

## 1. Architectural Overview

Curexal requires all healthcare providers (Clinics, Laboratories, Hospitals, Pharmacies, Diagnostic Centers) to undergo formal regulatory document verification before being granted operational access to active healthcare workflows.

Registration and Verification are distinct lifecycle stages:
- **Registration** (`POST /api/v1/organizations`): Creates core organization, owner identity, primary workspace, and subscription records in an atomic ACID transaction. The organization is created with status `pending_verification`.
- **Document Submission** (`POST /api/v1/organizations/{id}/documents`): Organization owners submit required regulatory documentation.
- **Platform Review** (`PATCH /api/v1/platform/documents/{docID}/review`): Authorized Curexal platform staff review submitted documents.
- **Organization Activation** (`POST /api/v1/platform/organizations/{id}/approve`): The backend verifies that all mandatory document types for the organization category are approved before updating status to `active`.

---

## 2. Organization & Document State Machine

```
Self-Service Registration
           │
           ▼
Status: pending_verification
           │
           ▼
Document Upload (Owner)
  - Filename & MIME validation
  - SHA-256 calculation
  - Stored in ObjectStorageService
  - Metadata in organization.organization_documents (status='pending', version=1)
           │
           ▼
Platform Staff Review
  - Requires permission 'organization:document:review'
  - Self-review blocked (uploader cannot review own document)
  - Action: APPROVE or REJECT
           │
   ┌───────┴───────────────────────────────────┐
   │                                           │
   ▼ APPROVED                                  ▼ REJECTED
Check Org Document Requirements             Owner Resubmission
   │                                           - Uploads version N+1
   ├─► ALL REQUIRED APPROVED?                  - Historical v1 records preserved
   │     │                                     - Metadata status='pending'
   │     ▼ YES                                 - Re-triggers Platform Review
   │   Platform Approval (`POST /platform/organizations/{id}/approve`)
   │   - Requires permission 'organization:verify'
   │   - Org status -> ACTIVE
   │
   └─► MISSING REQUIRED DOCS ──► Org remains 'pending_verification'
```

---

## 3. Storage Security & Binary Isolation

1. **Database Metadata & Binary Isolation**:
   - Binary document files are NEVER stored in PostgreSQL.
   - Files are stored in `ObjectStorageService` (`internal/platform/storage`).
   - PostgreSQL table `organization.organization_documents` stores metadata, SHA-256 checksum, version, and review records.
2. **Presigned URL Security**:
   - Bucket objects are private.
   - Read access is granted exclusively through short-lived presigned URLs (15-minute expiration).
3. **Storage Key Determinism & Path Traversal Protection**:
   - Keys are deterministically built server-side: `organizations/{organizationID}/documents/{docType}/{documentID}/v{version}`.
   - Client filenames are sanitized via `filepath.Clean`.
4. **Non-ACID Storage/Database Rollback**:
   - If PostgreSQL metadata persistence fails after object storage upload succeeds, the application executes an automatic cleanup to delete the orphaned storage object.

---

## 4. Required Document Catalog per Healthcare Category

| Healthcare Category | Required Regulatory Documents |
| :--- | :--- |
| **Laboratory** | `operating_license`, `registration_certificate` |
| **Clinic** | `operating_license`, `facility_license` |
| **Pharmacy** | `operating_license`, `pharmacy_license` |
| **Hospital** | `operating_license`, `facility_license`, `registration_certificate` |
| **Radiology** | `operating_license`, `accreditation_certificate` |
| **General Provider** | `operating_license` |

---

## 5. Security & Authorization Matrix

| Endpoint | Method | Required Permission / Role | Security Controls |
| :--- | :--- | :--- | :--- |
| `/api/v1/organizations/{id}/documents` | `POST` | `organization:document:upload` | Verify Org Membership, MIME type check, SHA-256 |
| `/api/v1/organizations/{id}/documents` | `GET` | `organization:document:read` | Verify Org Membership, Presigned URL generation |
| `/api/v1/platform/documents/{docID}/review` | `PATCH` | `organization:document:review` | Platform Staff only, Self-review blocked |
| `/api/v1/platform/organizations/{id}/approve` | `POST` | `organization:verify` | Platform Staff only, Checks all required docs approved |
| `/api/v1/platform/organizations/{id}/reject` | `POST` | `organization:verify` | Platform Staff only, Requires rejection reason |
