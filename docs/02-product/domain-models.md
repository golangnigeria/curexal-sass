# Curexal V2 Domain Models & Aggregates

This document defines the pure domain models, aggregates, value objects, and business invariants of the Curexal V2 platform. These models contain **zero database imports, zero SQL statements, and zero framework dependencies**.

---

## 1. Patient Aggregate

### `Patient` (Aggregate Root)
- **Attributes**:
  - `ID`: Unique Identifier (UUID)
  - `TenantID`: Tenant Identifier
  - `Demographics`: `PatientName`, `DateOfBirth`, `Gender`, `Phone`, `Email` (Value Objects)
  - `NationalID`: National ID Number (Value Object)
  - `EmergencyContact`: `ContactName`, `Relationship`, `Phone` (Value Object)
  - `BloodGroup`: Enum (`A+`, `A-`, `B+`, `B-`, `AB+`, `AB-`, `O+`, `O-`)
  - `Genotype`: Enum (`AA`, `AS`, `SS`, `AC`)
  - `Status`: Enum (`Active`, `Archived`, `Merged`)
  - `CreatedAt`: Timestamp
- **Business Invariants**:
  - `DateOfBirth` must be in the past.
  - `Phone` must conform to E.164 regional phone formatting.
  - A patient cannot be deleted if active referral orders or unfulfilled diagnostic requests exist.

---

## 2. Referral Order Aggregate

### `ReferralOrder` (Aggregate Root)
- **Attributes**:
  - `ID`: Referral Identifier (UUID)
  - `ReferringOrganizationID`: Tenant Identifier of referring clinic
  - `TargetOrganizationID`: Tenant Identifier of diagnostic lab
  - `PatientID`: Patient Identifier
  - `TestItems`: List of `ReferralTestItem` (Entity)
  - `ClinicalSummary`: Text String
  - `Status`: Enum (`Dispatched`, `Accepted`, `Rejected`, `SpecimenCollected`, `Testing`, `Completed`, `Cancelled`)
  - `RejectionReason`: Optional Text String
  - `EstimatedTurnaroundTime`: Duration (Hours)
  - `DispatchedAt`: Timestamp
  - `CompletedAt`: Optional Timestamp
- **Business Invariants**:
  - `ReferringOrganizationID` cannot equal `TargetOrganizationID`.
  - A referral order in `Completed` status cannot be `Cancelled` or `Rejected`.
  - At least one valid diagnostic test item must be present upon dispatch.

---

## 3. Specimen Aggregate (Laboratory LIMS)

### `Specimen` (Aggregate Root)
- **Attributes**:
  - `ID`: Specimen Identifier (UUID)
  - `ReferralOrderID`: Optional Referral Identifier
  - `Barcode`: Code128 Barcode String (Value Object)
  - `SampleType`: Enum (`WholeBlood`, `Serum`, `Plasma`, `Urine`, `CSF`, `Sputum`, `Swab`, `Tissue`)
  - `ContainerType`: Enum (`EDTACap`, `SerumSeparator`, `SodiumCitrate`, `SterileCup`)
  - `CollectedBy`: Staff User Identifier
  - `CollectedAt`: Timestamp
  - `Status`: Enum (`Collected`, `InTransit`, `ReceivedAtBench`, `Processing`, `Analyzed`, `Rejected`, `Disposed`)
- **Business Invariants**:
  - `Barcode` must be globally unique across the laboratory tenant schema.
  - A specimen in `Analyzed` or `Disposed` state cannot return to `Collected`.

---

## 4. Diagnostic Result Aggregate

### `DiagnosticResult` (Aggregate Root)
- **Attributes**:
  - `ID`: Result Identifier (UUID)
  - `SpecimenID`: Specimen Identifier
  - `TestCatalogCode`: Standardized LOINC / Test Code
  - `ParameterResults`: List of `ParameterResult` (Value Object: `ParameterName`, `MeasuredValue`, `Unit`, `ReferenceRange`, `IsAbnormal`, `IsCritical`)
  - `ScientistValidatorID`: Staff User Identifier
  - `ScientistValidatedAt`: Optional Timestamp
  - `ConsultantApproverID`: Staff User Identifier
  - `ConsultantApprovedAt`: Optional Timestamp
  - `PDFDocumentKey`: Object Storage S3 Key String
  - `ReleaseStatus`: Enum (`Draft`, `ScientistValidated`, `ConsultantApproved`, `Released`)
- **Business Invariants**:
  - `ReleaseStatus` can only transition to `Released` if `ConsultantApprovedAt` is present.
  - If any parameter result is flagged `IsCritical`, an automatic critical notification trigger must fire.

---

## 5. Wallet Ledger Aggregate (Financial Clearing)

### `PlatformWallet` (Aggregate Root)
- **Attributes**:
  - `ID`: Wallet Identifier (UUID)
  - `OrganizationID`: Tenant Identifier
  - `Balance`: Currency Amount (Decimal)
  - `CreditLimit`: Maximum Allowed Negative Balance (Decimal)
  - `Currency`: Enum (`NGN`, `USD`, `GHS`, `KES`)
  - `Status`: Enum (`Active`, `Frozen`, `Overdrawn`)
- **Business Invariants**:
  - `Balance` cannot drop below `(-1 * CreditLimit)`. If exceeded, `Status` transitions to `Overdrawn` and marketplace intake is paused.
  - All balance changes must occur via double-entry `LedgerTransaction` records (never direct mutation).
