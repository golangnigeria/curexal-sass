# Curexal V2 Feature Flags & Module Matrix

This document defines the platform feature flag configuration matrix controlling module availability across tenant types and plan tiers.

---

## Master Feature Flag Matrix

| Feature Flag Key | Capability Module | Clinic Workspace | Laboratory (LIMS) | Pharmacy | Radiology (RIS) | Enterprise Plan Only |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `module.clinic.soap` | Outpatient Consultations | **Enabled** | Disabled | Disabled | Disabled | No |
| `module.clinic.referrals` | B2B Marketplace Referral Dispatch | **Enabled** | **Enabled** | Disabled | **Enabled** | No |
| `module.lims.phlebotomy` | Specimen Intake & Barcode Printing | Disabled | **Enabled** | Disabled | Disabled | No |
| `module.lims.analyzer_sync` | Edge Agent ASTM/HL7 Integration | Disabled | **Enabled** | Disabled | Disabled | **Yes** |
| `module.lims.iso15189` | Levey-Jennings & Westgard QC | Disabled | **Enabled** | Disabled | Disabled | **Yes** |
| `module.pharmacy.dispense` | E-Prescription Dispensing Desk | Disabled | Disabled | **Enabled** | Disabled | No |
| `module.ris.pacs_viewer` | Web-based Cornerstone.js DICOM Viewer | Disabled | Disabled | Disabled | **Enabled** | **Yes** |
| `module.ai.explainer` | AI Patient Vernacular Report Translation | **Enabled** | **Enabled** | **Enabled** | **Enabled** | No |
| `module.ai.radiology_draft` | AI Radiologist Draft Assistant | Disabled | Disabled | Disabled | **Enabled** | **Yes** |
| `module.finance.split_pay` | Automated Paystack/Flutterwave Settlement | **Enabled** | **Enabled** | **Enabled** | **Enabled** | No |
