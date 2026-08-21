# Curexal V2 Capability Architecture & Map

This document establishes the official Capability Map of the Curexal V2 platform. System capabilities define **what the platform can do** at a business level, completely abstracted from user interfaces, pages, sidebars, or databases.

---

## Master Capability Hierarchy

```text
Curexal Healthcare Platform
│
├── 1. Identity & Access Capability
│   ├── Universal Authentication
│   ├── Session Lifecycle Management
│   ├── Multi-Factor Authentication (MFA)
│   ├── Access Policy Enforcement (RBAC)
│   └── User Credential Management
│
├── 2. Multi-Tenant Platform Capability
│   ├── Tenant Workspace Provisioning
│   ├── Workspace Context Resolution
│   ├── Organization & Branch Registry
│   ├── Feature Flag & Module Licensing
│   └── Organization Customization & Branding
│
├── 3. Patient Identity & Health Record Capability
│   ├── Master Patient Indexing (MPI)
│   ├── Demographic Record Management
│   ├── Medical History & Timeline Tracking
│   ├── Allergy & Chronic Condition Alerting
│   └── Patient Access Authorization
│
├── 4. B2B Marketplace Capability
│   ├── Provider Discovery & Search
│   ├── Service & Test Catalog Management
│   ├── Dynamic Pricing & SLA Matrix
│   ├── Provider Reputation & Rating
│   └── Marketplace Referral Router
│
├── 5. Clinical Encounter Capability
│   ├── Appointment Scheduling
│   ├── Clinical Consultation (SOAP)
│   ├── Diagnostic Request Ordering
│   └── Electronic Prescription Generation
│
├── 6. Diagnostic Laboratory (LIMS) Capability
│   ├── Phlebotomy & Specimen Intake
│   ├── Barcode Label Generation & Tracking
│   ├── Laboratory Bench Worklist Management
│   ├── Analyzer Data Ingestion (ASTM/HL7)
│   ├── Scientific Result Verification
│   ├── Consultant Approval Workflow
│   └── Diagnostic Report Generation
│
├── 7. Radiology (RIS/PACS) Capability
│   ├── Modality Worklist Management
│   ├── DICOM Image Metadata Storage
│   ├── Medical Image Streaming & Viewing
│   ├── Radiologist Voice Dictation
│   └── Structured Radiology Reporting
│
├── 8. Pharmacy & Dispensing Capability
│   ├── E-Prescription Matching
│   ├── Medication Dispensing Verification
│   ├── Drug Dosage Validation
│   └── Medication Batch Allocation
│
├── 9. Inventory & Supply Logistics Capability
│   ├── Multi-Warehouse Inventory Tracking
│   ├── Reagent & Consumable Stock Management
│   ├── Expiry Date & Low-Stock Alerting
│   ├── Purchase Order Generation
│   └── Goods Received Verification
│
├── 10. Financial & Settlement Capability
│   ├── Point-of-Sale (POS) Cashiering
│   ├── Invoice & Receipt Generation
│   ├── Automated Split-Payment Clearing
│   ├── Referral Commission Ledger Tracking
│   └── Virtual Balance Wallet Management
│
├── 11. Notification & Communication Capability
│   ├── Event Message Routing
│   ├── Email Template Dispatch
│   ├── SMS Alert Delivery
│   └── WhatsApp Business API Report Delivery
│
├── 12. Artificial Intelligence & Clinical Assistance Capability
│   ├── AI Clinical Report Summarization
│   ├── AI Patient Vernacular Report Translation
│   ├── AI Radiologist Draft Assistant
│   └── AI Operational Insights Engine
│
├── 13. Quality Management (ISO 15189) Capability
│   ├── Quality Control (QC) Data Ingestion
│   ├── Levey-Jennings & Westgard Validation
│   ├── Equipment Maintenance Logging
│   └── Non-Conformance Reporting (NCR)
│
├── 14. Observability & Security Audit Capability
│   ├── Correlation ID Request Tracing
│   ├── Telemetry Metrics Collection
│   ├── Cryptographic Hash-Chained Audit Logging
│   └── Security Event Violation Alerting
│
└── 15. Platform Integration & Developer Capability
    ├── REST & WebSocket API Integration
    ├── Webhook Event Dispatching
    ├── Interoperability Connector (FHIR / HL7)
    └── Developer Sandbox & API Key Provisioning
```

---

## Detailed Capability Definitions

### 1. Identity & Access Capability
- **Universal Authentication**: Authenticates user identity across the global registry.
- **Session Lifecycle Management**: Manages access token generation, rotation, and revocation.
- **Access Policy Enforcement**: Evaluates role-based and context-scoped permissions.

### 2. Multi-Tenant Platform Capability
- **Tenant Workspace Provisioning**: Orchestrates dynamic database schema creation and seed execution.
- **Workspace Context Resolution**: Resolves tenant scope from host headers and user memberships.

### 3. Patient Identity & Health Record Capability
- **Master Patient Indexing**: Prevents duplicate patient profile creation across independent providers.
- **Medical History & Timeline Tracking**: Aggregates visits, diagnoses, prescriptions, and test results chronologically.

### 4. B2B Marketplace Capability
- **Provider Discovery & Search**: Enables clinics to search for labs and imaging centers by location, turnaround time, ratings, and price catalogs.
- **Marketplace Referral Router**: Dispatches digital order envelopes across tenant boundaries.

### 5. Clinical Encounter Capability
- **Clinical Consultation**: Captures Subjective, Objective, Assessment, and Plan notes during patient visits.
- **Diagnostic Request Ordering**: Generates digital laboratory and radiology test orders.

### 6. Diagnostic Laboratory (LIMS) Capability
- **Barcode Label Generation**: Generates standard Code128 barcodes for specimen tubes.
- **Analyzer Data Ingestion**: Receives raw scientific output strings from analyzer instruments.

### 7. Radiology (RIS/PACS) Capability
- **DICOM Image Streaming**: Streams medical DICOM study files to web viewers.
- **Structured Radiology Reporting**: Allows radiologists to sign off on structured diagnostic reports.

### 8. Financial & Settlement Capability
- **Automated Split-Payment Clearing**: Executes real-time payment splits between laboratory, clinic, and platform.
- **Virtual Balance Wallet Management**: Manages pre-funded balances and cash referral debits.
