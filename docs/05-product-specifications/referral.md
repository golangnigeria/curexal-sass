# Curexal V2 Product Requirement Document (PRD): B2B Referral Module

This document serves as the absolute product specification for the B2B Referral Module in Curexal V2.

---

## 1. Module Overview & Purpose

The B2B Referral Module enables outpatient clinics, general practitioners, and diagnostic providers to digitally refer patients to partner laboratories, imaging centers, and specialist providers across the Curexal marketplace.

### Core Objectives
- Replace paper referral notes with cryptographically verifiable digital referral envelopes.
- Eliminate lost diagnostic orders and provide real-time status tracking for referring doctors.
- Automate referral commission accounting and split-payment settlements between healthcare partners.

---

## 2. Actors & User Roles

| Actor | Role | Permissions Required | Primary System Responsibilities |
| :--- | :--- | :--- | :--- |
| **Referring Clinician** | Clinic Doctor / GP | `referrals:create`, `referrals:read` | Selects target lab, picks diagnostic tests, writes clinical summary, dispatches order. |
| **Lab Intake Officer** | Lab Frontdesk / Receptionist | `referrals:read`, `referrals:accept`, `referrals:reject` | Reviews incoming referrals, accepts case with estimated TAT, rejects invalid requests. |
| **Lab Scientist** | Medical Lab Technologist | `referrals:read`, `specimens:collect` | Links specimen container, prints barcode, logs analytical test results. |
| **Patient** | Healthcare Consumer | `patient:portal_access` | Receives SMS/WhatsApp alert, views lab location, completes Paystack payment. |
| **Platform Administrator** | Curexal System Admin | `referrals:admin_override` | Resolves disputed referrals, audits commission split ledgers. |

---

## 3. End-to-End Referral Lifecycle & State Machine

```
               [Draft]
                  │
                  ▼
            [Dispatched] ───────────────┐
                  │                     │
       ┌──────────┴──────────┐          │ (Patient Cancels)
       ▼                     ▼          ▼
   [Accepted]           [Rejected]  [Cancelled]
       │
       ▼
[SpecimenCollected]
       │
       ▼
   [Testing]
       │
       ▼
  [Completed] (Results Released & Commission Settled)
```

### State Definitions & Transitions
1. **`Draft`**: Referral order is being composed by the clinician. Not visible to target lab.
2. **`Dispatched`**: Referral published to global directory. Target lab inbox receives notification (`referral.dispatched`).
3. **`Accepted`**: Target lab accepts order, providing an estimated Turnaround Time (TAT).
4. **`Rejected`**: Target lab rejects order (e.g., test unavailable, specimen requirements not met). Requires rejection reason.
5. **`SpecimenCollected`**: Phlebotomist scans specimen tube and links Code128 barcode.
6. **`Testing`**: Sample processed on bench / analyzer.
7. **`Completed`**: Scientific results authorized, PDF released, commission ledger settled (`results.released`).
8. **`Cancelled`**: Clinician or patient cancels order before sample collection.

---

## 4. Detailed Business & Validation Rules

### Referral Dispatch Rules
- **Rule REF-001**: A clinic cannot dispatch a referral to its own organization (`ReferringOrgID != TargetOrgID`).
- **Rule REF-002**: Every test code included in the referral order must exist and be marked active in the target laboratory's public pricing catalog.
- **Rule REF-003**: The clinical summary field must not exceed 2,000 characters and must automatically strip HTML/script tags to prevent XSS attacks.

### Rejection & Timeout Rules
- **Rule REF-004**: If a target laboratory does not respond to a `Dispatched` referral within 24 hours, the order automatically transitions to `Expired`, and the referring clinic is alerted to select an alternative provider.

### Payment & Wallet Rules
- **Rule REF-005**: If the patient pays online via Paystack, the referral is flagged `Prepaid` and funds are held in escrow until result release.
- **Rule REF-006**: If the patient pays cash at the laboratory desk, the laboratory platform wallet is debited by the referral commission percentage upon result release.

---

## 5. Event Storming & Audit Specifications

### Triggered Events (NATS JetStream)
- `referral.dispatched`: Payload contains `referral_id`, `referring_org_id`, `target_org_id`, `patient_id`, `tests`.
- `referral.accepted`: Payload contains `referral_id`, `estimated_tat_hours`, `accepted_by_user_id`.
- `referral.completed`: Payload contains `referral_id`, `diagnostic_report_id`, `pdf_s3_key`.

### Immutable Audit Log Specifications
Every state transition records a hash-chained log entry in `internal/core/audit`:
- `LogID`, `ReferralID`, `PreviousStatus`, `NewStatus`, `ActorUserID`, `IPAddress`, `PreviousHash`, `CurrentHash`.

---

## 6. Acceptance Criteria

- [x] Must execute zero cross-schema SQL joins during referral dispatch and intake.
- [x] Must dispatch instant SMS/WhatsApp notifications to patient with Paystack payment links.
- [x] Rejection of a referral must enforce selection of a valid reason code from standard dropdown.
- [x] Commission split must automatically update wallet ledgers upon result release.
