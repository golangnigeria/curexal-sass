# Curexal V2 UX Journeys & Permissions Matrix

This document maps out the user experience journeys and permission matrix across all platform actors and frontend web applications.

---

## 1. Role-Permissions Matrix

| Role | Application Host | Key Capabilities Granted | API Permissions |
| :--- | :--- | :--- | :--- |
| **Clinician / Doctor** | `apps/workspace` | Consultations, E-Prescriptions, B2B Referral Dispatch | `encounters:create`, `prescriptions:create`, `referrals:create` |
| **Lab Scientist** | `apps/workspace` | Sample Intake, Barcoding, Result Validation | `specimens:collect`, `worklists:write`, `results:validate` |
| **Lab Director** | `apps/workspace` | Consultant Approval, Report Release, QC Management | `reports:release`, `qc:manage` |
| **Pharmacist** | `apps/workspace` | Prescription Dispensing, Stock Deduction | `dispensing:write`, `inventory:update` |
| **Patient** | `apps/patient` | Self-Registration, Appointment Booking, Result Portal | `patient:portal_access`, `results:download` |
| **Platform Admin** | `apps/admin` | Schema Provisioning Monitoring, Global Catalogs | `admin:all`, `tenants:manage` |

---

## 2. Actor Journey Maps

### Doctor Journey (`apps/workspace`)
```text
Login ──► Appointments Queue ──► Start Encounter ──► Input SOAP Notes ──► Issue E-Prescription ──► Dispatch Referral ──► Complete Encounter
```

### Lab Scientist Journey (`apps/workspace`)
```text
Laboratory Inbox ──► Accept Referral ──► Specimen Intake ──► Print Barcode ──► Ingest Analyzer Data ──► Validate Results ──► Submit for Sign-off
```

### Patient Journey (`apps/patient`)
```text
Receive Referral Alert ──► View Lab Location ──► Pay online via Paystack ──► Present Barcode ──► Download PDF Results
```
