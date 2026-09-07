# CUREXAL CLINIC OS — MVP DEFINITION
**Document**: `specs/00-product/mvp-definition.md`  
**Status**: APPROVED BASELINE  
**Scope**: 14-Day Production Release  

---

## 1. Purpose & Vision

Enable a small-to-medium outpatient clinic to digitally manage its complete patient lifecycle from registration through appointment, triage, consultation, diagnosis, prescription, invoice, and payment.

---

## 2. Included MVP Capabilities

```text
[1. PATIENT REGISTRATION] ──► [2. APPOINTMENT & QUEUE] ──► [3. TRIAGE VITALS]
                                                                  │
[6. PATIENT PORTAL] ◄── [5. INVOICE & POS PAYMENT] ◄── [4. CONSULTATION & SOAP]
                                                                  │
                                                        [4b. ICD-10 & RX]
```

### Core MVP Capabilities
1. **Organization & Facility Setup**: Multi-tenant organization creation, branch facility provisioning.
2. **Staff Onboarding & RBAC**: Granular roles (Org Admin, Branch Manager, Doctor, Nurse, Receptionist, Cashier).
3. **Master Patient Index (MPI) Intake**: Unique MRN assignment, demographics, contact details, duplicate search.
4. **Appointment Scheduling & Live Queue**: Provider calendars, patient check-in, real-time waiting list.
5. **Nurse Triage Intake**: Vital signs recording (BP, Pulse, Temp, SpO2, Weight, Height), allergy tagging.
6. **Doctor Consultation Canvas**: Electronic encounter notes (SOAP), chief complaint, clinical examination.
7. **ICD-10 Diagnostic Coding**: Standardized diagnosis selection, primary/secondary condition tagging.
8. **E-Prescribing Contract**: Digital prescription creation (Drug, Dose, Route, Frequency, Duration).
9. **Billing & Invoicing Engine**: Auto-generated billable service invoices, fee schedules.
10. **POS Cashier Checkout**: Cash, Card, and Transfer payment recording, instant receipt generation.
11. **Patient Health Summary**: Digital encounter summary accessible to the patient.
12. **Care Notification Triggers**: Automated appointment and billing alerts via email/SMS.

---

## 3. Explicit MVP Non-Goals (Strictly Excluded from MVP Release)

The following systems are **explicitly excluded** from the first 14-day production release:
- Diagnostic Laboratory Information System (LIS) & analyzer serial telemetry.
- Radiology Information System (RIS), PACS, and DICOM imaging viewers.
- Inpatient Hospital Information System (HIS) & ward bed occupancy tracking.
- Complex pharmacy stock batching and supplier purchase orders.
- Automated HMO insurance claims clearinghouse and reconciliation.
- Multi-facility diagnostic referral marketplace.
- Autonomous AI diagnosis or automated clinical actions.
- Large-scale event streaming brokers (Kafka/RabbitMQ).

---

## 4. End-to-End Acceptance Criteria

The MVP is complete when a pilot clinic successfully executes:
1. Registration of a new patient $\to$ MRN created in `patient.patients`.
2. Booking an appointment $\to$ Checked in $\to$ Moved to Triage Queue.
3. Nurse records vitals $\to$ Patient appears in Doctor Consultation Queue.
4. Doctor opens encounter $\to$ Records SOAP note $\to$ Selects ICD-10 code $\to$ Issues E-Prescription.
5. Invoice automatically generated $\to$ Cashier receives payment at POS $\to$ Invoice marked `PAID`.
6. Patient receives receipt and digital clinical encounter summary on their portal.
7. Entire lifecycle audited in `audit.events` with strict tenant isolation.
