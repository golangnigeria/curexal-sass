# Curexal V2 Domain Map & Bounded Contexts

This document defines the Domain-Driven Design (DDD) Bounded Contexts, aggregates, entities, and relationship mappings for Curexal V2.

---

## 1. Domain Map & Bounded Contexts Overview

```
                          +──────────────────────────────────────+
                          |        Identity Bounded Context      |
                          | (Users, Sessions, MFA, Memberships)  |
                          +----------------──┬───────────────────+
                                             │
                                             v
                          +──────────────────────────────────────+
                          |         Tenant Bounded Context       |
                          | (Organizations, Schemas, Features)   |
                          +----------------──┬───────────────────+
                                             │
      ┌──────────────────────────────────────┼──────────────────────────────────────┐
      │                                      │                                      │
      v                                      v                                      v
+──────────────────────────+   +──────────────────────────+   +──────────────────────────+
|  Patient Bounded Context |   | Referral Bounded Context |   | Marketplace Context      |
| (Demographics, Timeline) |   | (State Router, Envelopes)|   | (Discovery, Directories) |
+────────────┬─────────────+   +─────────────┬────────────+   +─────────────┬────────────+
             │                               │                              │
             └───────────────────────┬───────┴──────────────────────────────┘
                                     │
    ┌────────────────────────────────┼────────────────────────────────┐
    │                                │                                │
    v                                v                                v
+──────────────────────────+   +──────────────────────────+   +──────────────────────────+
|   Clinic Bounded Context |   |      Lab LIMS Context    |   |  Financial Ledger Context|
| (Encounters, SOAP, Rx)   |   | (Samples, Results, QC)   |   | (Wallets, Split Payments)|
+──────────────────────────+   +──────────────────────────+   +──────────────────────────+
```

---

## 2. Bounded Context Specifications & Aggregates

### 2.1. Identity Bounded Context
- **Primary Aggregates**:
  - `User`: Credential, PasswordHash, MFASecret, Status.
  - `Session`: AccessToken, RefreshToken, ExpiresAt, RevokedAt.
  - `Membership`: UserID, OrganizationID, Role, BranchID.
- **Ubiquitous Language**: User, Principal, Credential, JWT, Refresh Token, Session Rotation, Active Context.

### 2.2. Tenant Bounded Context
- **Primary Aggregates**:
  - `Organization`: Slug, Name, FacilityType, ProvisioningStatus.
  - `Branch`: Address, Phone, OperatingHours, Timezone.
  - `FeatureToggle`: ModuleName, IsEnabled, MaxVolumeLimit.
- **Ubiquitous Language**: Tenant, Schema, Provisioning, Workspace, Facility Type, Branch.

### 2.3. Patient Bounded Context
- **Primary Aggregates**:
  - `Patient`: Demographics, NationalID, ContactPhone, EmergencyContact.
  - `MedicalTimeline`: EntryID, EntryType (Visit, Diagnosis, Prescription, Result), Timestamp, ProviderID.
- **Ubiquitous Language**: Master Patient Index (MPI), Demographic Match, Medical Timeline, Allergy Alert.

### 2.4. Clinic Bounded Context
- **Primary Aggregates**:
  - `Encounter`: PatientID, ClinicianID, StartTime, EndTime, Status.
  - `SOAPNote`: Subjective, Objective, Assessment, Plan.
  - `Prescription`: MedicationList, DosageRules, Instructions.
- **Ubiquitous Language**: Encounter, Consultation, SOAP, Diagnosis, ICD-10 Code, E-Prescription.

### 2.5. Referral Bounded Context
- **Primary Aggregates**:
  - `ReferralOrder`: ReferralID, ReferringOrgID, TargetOrgID, PatientID, TestList, Status, ClinicalNotes.
- **Ubiquitous Language**: Referral Dispatch, Referral Intake, Referral Envelope, Target Laboratory, Order State Router.

### 2.6. Laboratory (LIMS) Bounded Context
- **Primary Aggregates**:
  - `Specimen`: BarcodeCode128, ContainerType, CollectedAt, PhlebotomistID, Status.
  - `TestOrder`: OrderID, SpecimenID, TestCatalogCode, BenchID, Status.
  - `DiagnosticResult`: TestOrderID, ParameterValues, MeasuredAt, ScientistID, ApprovedBy, ReleaseStatus.
  - `QCRun`: ControlLot, TargetMean, StandardDeviation, MeasuredValue, WestgardViolation.
- **Ubiquitous Language**: Phlebotomy, Specimen, Barcode, Bench Worklist, Scientific Validation, Consultant Release, Levey-Jennings.

### 2.7. Financial Ledger Bounded Context
- **Primary Aggregates**:
  - `Invoice`: InvoiceNumber, PatientID, LineItems, TotalAmount, Tax, Status.
  - `PlatformWallet`: OrganizationID, Balance, Currency, CreditLimit.
  - `LedgerTransaction`: TransactionID, SourceWalletID, TargetWalletID, Amount, TransactionType (Split, Commission, CashDebit, Payout).
- **Ubiquitous Language**: Ledger Journal, Double-Entry, Split Payment, Platform Fee, Referral Commission, Wallet Credit Limit.

---

## 3. Context Mapping & Relationships

| Source Context | Target Context | Relationship Type | Communication Channel |
| :--- | :--- | :--- | :--- |
| **Identity Context** | **Tenant Context** | Customer-Supplier | Direct Shared Schema Lookup |
| **Clinic Context** | **Referral Context** | Publisher-Subscriber | NATS Event (`referral.dispatched`) |
| **Referral Context** | **Lab LIMS Context** | Publisher-Subscriber | NATS Event (`referral.accepted`) |
| **Lab LIMS Context** | **Financial Ledger Context**| Downstream Consumer | Event (`results.released` -> Split Trigger) |
| **Lab LIMS Context** | **Patient Context** | Synchronizer | Asynchronous Patient Timeline Event |
