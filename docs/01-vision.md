# Curexal V2 Product Vision & Strategy

This document serves as the absolute, unchanging foundation of the Curexal V2 platform. All design decisions, database designs, API contracts, and user experience paradigms must align with the core philosophy and strategy detailed herein.

---

## 1. The Problem Statement

Healthcare delivery across Africa and emerging markets is severely constrained by structural, technological, and infrastructural fragmentation. Existing healthcare software solutions fail because they treat healthcare providers as isolated islands rather than nodes in an interconnected network.

### Key Operational Challenges Addressed by Curexal V2

#### 1. Fragmented Healthcare Silos & Paper Referrals
- Clinics, laboratories, imaging centers, and pharmacies operate in complete isolation.
- Referrals are handled via paper notes carried manually by patients. This leads to lost referrals, untracked patient histories, miscommunicated orders, and severe diagnostic delays.
- Referring doctors have zero visibility into whether a patient completed a diagnostic test or what the results were, directly compromising patient care.

#### 2. Infrastructural & Connectivity Vulnerability
- Western cloud SaaS platforms assume 99.99% internet uptime and high bandwidth.
- In Nigeria and across Africa, power grid instability and intermittent network access cause traditional cloud-only software to collapse, bringing critical laboratory phlebotomy, sample logging, and emergency care to a halt.

#### 3. Opaque Financial Settlements & B2B Commission Friction
- B2B referral commissions between clinics and diagnostic centers are tracked manually on paper logbooks or spreadsheets.
- This creates accounting errors, delayed payouts, lack of transparency, lost revenue, and mistrust between healthcare partners.

#### 4. Cost & Complexity of Fragmented Vendor Systems
- Healthcare facilities are forced to purchase, integrate, and maintain separate legacy systems for LIS, RIS, EMR, Billing, and Inventory.
- Small-to-medium diagnostic centers and clinics cannot afford multi-vendor enterprise software suites, trapping them in manual, paper-heavy operations.

#### 5. Poor Patient Experience & Delayed Diagnostic Results
- Patients must make multiple physical trips back and forth between clinics and labs to hand-deliver requests and collect paper test results.
- Delivery of critical diagnostic reports lacks modern digital accessibility (e.g., instant WhatsApp/email delivery or direct patient portal access).

---

## 2. Mission Statement

> **To build the digital backbone of healthcare delivery in Africa by connecting providers, patients, and partners on a single intelligent operating network.**

Curexal V2 resolves the systemic friction of fragmented, offline, and siloed healthcare systems in emerging markets. By providing a unified software operating system combined with an active collaborative marketplace, Curexal V2 enables independent providers to digitize their internal workflows while participating in a shared digital economy.

---

## 3. Core Philosophy

### A Healthcare Operating System, Not an HMS
Curexal is not a standard, isolated Hospital Management System (HMS). It is an **Operating System (HOS)**.
- Traditional HMS solutions focus exclusively on internal hospital boundaries.
- Curexal operates as a **unified cloud platform** where multiple independent tenants (Clinics, Laboratories, Radiologists, Pharmacies) collaborate dynamically.
- Data flows securely across organizational boundaries via digital referrals, automated reports, and split-payment transactions.

### Schema-per-Tenant Multi-Tenancy
To guarantee absolute data privacy and comply with regional regulations (such as the Nigeria Data Protection Regulation - NDPR and HIPAA standards), every tenant's operational data is isolated within its own dedicated database schema. No tenant can ever perform a database-level query against another tenant's schema.

### Collaboration by Default
The platform is designed around the **referral lifecycle**. A clinic should be able to order diagnostic tests from a third-party laboratory, track the sample collection, receive the scientific results, and automatically pay the laboratory—all without leaving the Curexal ecosystem.

### Localization for Africa
Curexal is designed to solve real-world infrastructural bottlenecks in Africa:
- **Low-Bandwidth & Offline Capabilities**: Applications continue functioning during network outages, synchronizing data when online.
- **African Payment Gateways**: Native integration with Paystack and Flutterwave for automated split payouts, cash-reconciliation ledgers, and patient billing.
- **WhatsApp Integration**: Using WhatsApp as the primary communication channel for appointment scheduling, patient notifications, and diagnostic report delivery.

---

## 4. Business & Monetization Model

Curexal V2 employs a hybrid monetization strategy that combines SaaS subscription revenue with marketplace transaction fees.

### SaaS Subscription (Operational Software)
Tenants pay a subscription fee scaled by:
- **Active Modules**: A clinic only pays for the Clinic/EHR module. A laboratory pays for the LIMS module.
- **Volume Metrics**: Number of registered patients, branches, or processed diagnostic test orders.
- **Enterprise Features**: Custom branding, advanced analytics, analyzer integrations, and multi-branch networks.

### B2B Marketplace Transactions (Financial Clearing)
When a referral transaction occurs through the marketplace:
- **Platform Fee**: Curexal charges a small transaction percentage (e.g., 2%) on the total cost of diagnostic orders referred through the system.
- **Commission Clearing**: The platform manages the settlement ledger between referring clinics and performing laboratories, automatically deducting commission amounts.

---

## 5. Target Customer Segments

| Tenant Type | Core Value Proposition | Primary Features Used |
| :--- | :--- | :--- |
| **Outpatient Clinics** | Streamline patient workflows, dispatch referrals to top diagnostic centers, and track outcomes. | SOAP notes, prescriptions, digital referral portal, patient billing. |
| **Diagnostic Laboratories** | Automated sample tracking, analyzer connectivity, scientific validation, and marketplace discovery. | LIMS phlebotomy queue, barcode printing, analyzer sync, WhatsApp report delivery. |
| **Radiology Centers** | DICOM file linking, PACS compatibility, radiologist reporting interfaces, and voice dictation. | RIS dashboard, DICOM metadata linking, PACS viewer integrations. |
| **Pharmacies** | E-prescription intake, inventory management, dispensing queues, and stock alerts. | Dispensing queue, POS billing, expiry monitor. |

---

## 6. The Long-Term Vision

As Curexal V2 scales across Africa:
- **National Health Insurance Integrations**: Connecting directly with national networks (like NHIA in Nigeria) for automated claims submission and verification.
- **Developer Platform (Marketplace SDK)**: Exposing secure, sandboxed APIs to allow third-party developers to build custom modules, patient monitoring integrations, and niche diagnostics tools.
- **Ecosystem Network Effects**: Transitioning the diagnostic healthcare landscape from a collection of manual, paper-dependent actors into an interconnected, real-time digital grid.
