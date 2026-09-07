# CUREXAL PLATFORM — MODULE SYSTEM SPECIFICATION
**Document**: `specs/01-platform/module-system.md`  
**Status**: APPROVED BASELINE  

---

## 1. Modular Architecture

Curexal business modules represent discrete commercial and clinical capabilities.

```text
[Core / MPI] ◄──────┐
      │             │
      ▼             │
   [Clinic] ──► [Billing Engine] ◄── [Future LIS / Pharmacy]
```

---

## 2. Module Rules & Invariants

1. **No God Services**: Never create a monolithic `PatientService` that creates patients, books appointments, writes SOAP notes, accessions lab specimens, and dispenses drugs.
2. **Private Data Schemas**:
   - `core` schema owns `patient.patients`, `organization.*`, `identity.*`.
   - `clinic` schema owns `clinic.appointments`, `clinic.encounters`, `clinic.vitals`, `clinic.diagnoses`, `clinic.prescriptions`.
   - `billing` schema owns `billing.invoices`, `billing.invoice_items`, `billing.payments`.
3. **Event-Driven Integration**: When a clinic consultation completes with billable services, the clinic module emits `encounter.completed`, which the billing engine consumes to generate `billing.invoice`.
