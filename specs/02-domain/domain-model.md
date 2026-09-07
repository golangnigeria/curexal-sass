# CUREXAL CLINIC MVP — DOMAIN MODEL SPECIFICATION
**Document**: `specs/02-domain/domain-model.md`  
**Status**: APPROVED BASELINE  

---

## 1. Domain Entities & Ownership Boundaries

```text
Organization (Core)
    └── Facility (Core)
          └── Patient (Core / MPI)
                └── Appointment (Clinic)
                      └── Queue Entry (Clinic)
                            └── Encounter (Clinic)
                                  ├── Vitals (Clinic)
                                  ├── SOAP Note (Clinic)
                                  ├── Diagnosis (Clinic)
                                  ├── Prescription (Clinic)
                                  └── Invoice (Billing)
                                        └── Payment (Billing)
```

---

## 2. Canonical Entity Definitions

### Patient (Core / MPI)
- **Purpose**: Canonical identity of an individual receiving healthcare.
- **Ownership**: Platform Core (`patient.patients`).
- **Key Fields**: `id (UUID)`, `organization_id (UUID)`, `mrn (String, Unique)`, `first_name`, `last_name`, `date_of_birth`, `gender`, `phone`, `email`, `blood_group`, `genotype`, `allergies (JSONB)`.

### Appointment & Queue (Clinic)
- **Purpose**: Outpatient visit booking and live facility waiting workflow.
- **States**: `SCHEDULED` $\to$ `CHECKED_IN` $\to$ `IN_TRIAGE` $\to$ `WAITING_FOR_DOCTOR` $\to$ `IN_CONSULTATION` $\to$ `COMPLETED` (or `CANCELLED`, `NO_SHOW`).

### Encounter & SOAP Note (Clinic)
- **Purpose**: Attending physician's clinical documentation for a visit.
- **Structure**:
  - `Subjective`: Chief complaint, history of presenting illness.
  - `Objective`: Physical examination, measured vital signs.
  - `Assessment`: ICD-10 clinical diagnoses and severity.
  - `Plan`: Medical management, e-prescriptions, lifestyle counseling.

### Prescription (Clinic / Future Pharmacy Contract)
- **Purpose**: Ordered medication regimen for patient.
- **Fields**: `id`, `encounter_id`, `patient_id`, `medication_name`, `dosage`, `route`, `frequency`, `duration_days`, `instructions`, `status (PENDING, DISPENSED, CANCELLED)`.

### Invoice & Payment (Billing)
- **Purpose**: Financial itemization and receipting for rendered care.
- **States**: `DRAFT` $\to$ `ISSUED` $\to$ `PAID` (or `PARTIALLY_PAID`, `VOID`).
