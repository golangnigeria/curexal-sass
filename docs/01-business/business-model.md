# Curexal V2 Business & Marketplace Architecture

This document details the commercial model, revenue mechanics, target customer personas, and marketplace dynamics governing Curexal V2.

---

## 1. Hybrid Monetization Model

Curexal V2 employs a dual-revenue engine:
1. **Software-as-a-Service (SaaS) Subscriptions**: Recurring module-based subscription fees paid by healthcare organizations operating workspaces on the platform.
2. **Marketplace Transaction Clearing (Financial Split-Fees)**: Transaction percentages on diagnostic orders referred and settled across organizational boundaries.

```
                      +------------------------------------------+
                      |         Curexal Revenue Streams          |
                      +--------------------+---------------------+
                                           |
                ┌──────────────────────────┴──────────────────────────┐
                ▼                                                     ▼
+───────────────────────────────+                     +───────────────────────────────+
|     SaaS Subscriptions        |                     |   Marketplace Split Fees      |
|                               |                     |                               |
| - Clinic Workspace: $99/mo    |                     | - Platform Fee: 2% per Order  |
| - Laboratory LIMS: $199/mo    |                     | - Automated Clearing Ledger   |
| - Pharmacy Module: $79/mo     |                     | - Commission Payout Engine    |
| - Enterprise Multi-Branch     |                     | - Paystack/Flutterwave Split  |
+───────────────────────────────+                     +───────────────────────────────+
```

---

## 2. Target Customer Personas

| Persona | Role | Key Operational Needs | Curexal Solution |
| :--- | :--- | :--- | :--- |
| **Dr. Chioma** | Outpatient Clinic Physician | Digital SOAP notes, fast prescriptions, reliable diagnostic lab referral options. | E-Prescriptions & B2B Referral Dispatcher with real-time lab result tracking. |
| **Lab Director Ibrahim** | Diagnostic Laboratory Owner | Specimen barcode tracking, analyzer connectivity, fast report release, revenue growth. | LIMS Bench Worklists, Edge Instrument Sync, WhatsApp Report Delivery, Marketplace Discovery. |
| **Pharm. Emeka** | Retail Pharmacy Owner | E-prescription intake, inventory expiry management, point-of-sale checkouts. | Pharmacy Dispensing Station, Expiry Alert System, Integrated POS. |
| **Radiologist Dr. Bola** | Imaging Specialist | PACS viewer integration, voice dictation, structured reporting templates. | RIS Studio, Cornerstone.js DICOM Viewer, Dictation Panel. |
| **Patient Blessing** | Healthcare Consumer | Easy appointment booking, online bill payments, immediate diagnostic result access. | Patient Health Portal, Paystack Checkout, WhatsApp Result Notifications. |

---

## 3. Marketplace Referral & Clearing Dynamics

When Clinic A refers a patient to Laboratory B:
- **Test Order Cost**: NGN 50,000
- **Platform Fee (2%)**: NGN 1,000 (Credited to Curexal Platform Vault)
- **Clinic Commission (10%)**: NGN 5,000 (Credited to Clinic Settlement Wallet)
- **Lab Fee (88%)**: NGN 44,000 (Credited to Laboratory Wallet)

### Settlement Execution
- **Online Checkout**: Paystack/Flutterwave APIs execute immediate split transfers.
- **Counter Cash Transactions**: Double-entry ledger records a debit on the performing laboratory's platform wallet balance. The laboratory maintains a pre-funded wallet or automatic credit-card debit arrangement to ensure continuous marketplace referral reception.
