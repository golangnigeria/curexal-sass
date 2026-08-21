# Curexal V2 Product Requirement Document (PRD): Financial Clearing, Billing, & Wallet Module

This document serves as the absolute product specification for the Financial Clearing, Invoicing, POS Cashiering, and Settlement Wallet Module in Curexal V2.

---

## 1. Module Overview & Purpose

Manages patient billing, POS cash checkouts, invoice generation, split-payment clearing via Paystack/Flutterwave, and automated B2B referral commission settlement ledgers.

---

## 2. Business Rules & Financial Accounting Mechanics

### Double-Entry Ledger System
All wallet mutations are executed via double-entry journal logs (`LedgerTransaction`).

```
Cash Referral Payment (NGN 50,000)
    │
    ├─► Debit: Lab Platform Wallet (NGN 6,000)
    │      ├─► NGN 5,000 (10% Referral Commission)
    │      └─► NGN 1,000 (2% Platform Transaction Fee)
    │
    ├─► Credit: Clinic Commission Wallet (NGN 5,000)
    │
    └─► Credit: Curexal Platform Vault (NGN 1,000)
```

### Wallet Balance Limits & Referral Locking
- **Rule FIN-001**: Every organization has a pre-configured `CreditLimit` (e.g., NGN -50,000).
- **Rule FIN-002**: If cash referral debits drop a laboratory's wallet balance below `(-1 * CreditLimit)`, the system automatically transitions the wallet state to `Overdrawn` and pauses incoming marketplace referrals until pre-funded via card top-up.

---

## 3. Webhook Integration Specifications (Paystack & Flutterwave)

- **Signature Verification**: Webhooks MUST verify the HMAC SHA512 signature in the `X-Paystack-Signature` header using the secret key before processing.
- **Idempotency**: Webhook handlers check the transaction reference against the `ledger_transactions` table to prevent duplicate crediting.

---

## 4. Acceptance Criteria

- [x] All cash transactions must record balanced double-entry ledger entries.
- [x] Webhook signatures must be verified cryptographically before mutating financial state.
- [x] Reaching wallet credit limits must immediately lock marketplace intake.
