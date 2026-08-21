# Curexal V2 Event Catalog (Event Storming Architecture)

This document contains the complete Event Catalog for the Curexal V2 platform. Every domain event represents an immutable state change within a bounded context that drives asynchronous workflows, NATS messaging, audit trails, notification triggers, and AI contextual intelligence.

---

## Master Event Catalog Matrix

| Event Name | Publisher Bounded Context | Triggering Business Action | Key Data Payload Fields | Primary Subscribers & Consumers |
| :--- | :--- | :--- | :--- | :--- |
| `user.registered` | Identity | User account registration | `user_id`, `email`, `role`, `timestamp` | Audit Log, Email Worker (Verification Link) |
| `user.logged_in` | Identity | Session authentication | `user_id`, `session_id`, `ip_address`, `user_agent` | Audit Log, Security Monitor |
| `workspace.provisioned` | Tenant | Schema & seed execution complete | `organization_id`, `tenant_slug`, `schema_name` | Admin Dashboard, Email Worker (Welcome Mail) |
| `patient.registered` | Patient | Patient identity creation | `patient_id`, `phone`, `national_id_hash`, `tenant_slug` | Master Patient Index (MPI), Audit Log |
| `encounter.started` | Clinic | Consultation initialized | `encounter_id`, `patient_id`, `clinician_id`, `start_time` | Patient Medical Timeline, Audit Log |
| `prescription.created` | Clinic | E-prescription signed | `prescription_id`, `encounter_id`, `patient_id`, `medications` | Pharmacy Dispensing Queue, Patient Portal |
| `referral.dispatched` | Referral | B2B referral order generated | `referral_id`, `referring_org_id`, `target_org_id`, `tests` | Target Lab Inbox, SMS/WhatsApp Worker |
| `referral.accepted` | Referral | Receiving lab accepts order | `referral_id`, `receiving_staff_id`, `estimated_tat` | Referring Clinic Dashboard, Patient Portal |
| `referral.rejected` | Referral | Receiving lab rejects order | `referral_id`, `reason_code`, `rejected_at` | Referring Clinic Alert, Audit Log |
| `specimen.collected` | Lab LIMS | Phlebotomist logs specimen | `specimen_id`, `barcode_128`, `container_type`, `sample_type` | Lab Bench Worklist, Referral Tracker |
| `analyzer.data_received` | Lab LIMS | Edge agent posts ASTM data | `analyzer_id`, `barcode_128`, `raw_results_json` | Lab Worklist Result Matcher |
| `results.validated` | Lab LIMS | Scientist validates parameter | `test_order_id`, `parameter_code`, `validated_value` | Consultant Review Queue |
| `results.released` | Lab LIMS | Consultant authorizes report | `diagnostic_report_id`, `referral_id`, `pdf_s3_key` | Clinic Dashboard, Patient Portal, WhatsApp Worker |
| `payment.processed` | Financial | Paystack/Flutterwave webhook | `transaction_id`, `order_id`, `amount`, `payment_method` | Split Settlement Ledger, POS Receipts |
| `commission.settled` | Financial | B2B referral split executed | `ledger_id`, `referring_org_id`, `performing_org_id`, `split_amount` | Wallet Ledger, Clinic Wallet Manager |
| `qc.violation_detected` | Quality (ISO) | Westgard rule failure | `equipment_id`, `control_lot`, `violated_rule` | Lab Director Alert, Bench Channel Lock |

---

## Detailed Event Payload Specifications

### 1. `referral.dispatched`
- **Topic Name**: `curexal.events.referral.dispatched`
- **Publisher**: Clinic Bounded Context
- **Delivery**: At-Least-Once (NATS JetStream)
- **Schema Specification**:
```json
{
  "event_id": "evt_018e-4aef-92b1",
  "event_type": "referral.dispatched",
  "producer_tenant": "everight-clinic",
  "timestamp": "2026-07-21T17:46:00Z",
  "schema_version": "v1",
  "data": {
    "referral_id": "ref_99201",
    "referring_organization_id": "org_clinic_a",
    "target_organization_id": "org_synlab_b",
    "patient_id": "pat_88102",
    "patient_checksum": "sha256_a9b8...",
    "test_codes": ["FBC", "LIPID", "LFT"],
    "clinical_summary": "Patient presents with persistent fatigue and elevated blood pressure.",
    "urgent_flag": false
  }
}
```

### 2. `results.released`
- **Topic Name**: `curexal.events.results.released`
- **Publisher**: Laboratory (LIMS) Bounded Context
- **Schema Specification**:
```json
{
  "event_id": "evt_018e-4aef-9999",
  "event_type": "results.released",
  "producer_tenant": "synlab-main",
  "timestamp": "2026-07-21T17:50:00Z",
  "schema_version": "v1",
  "data": {
    "diagnostic_report_id": "rep_77102",
    "referral_id": "ref_99201",
    "patient_id": "pat_88102",
    "consultant_signoff_by": "usr_pathologist_ibrahim",
    "pdf_document_s3_key": "vault/synlab-main/2026/07/rep_77102.pdf",
    "critical_flag": false
  }
}
```
