# Curexal V2 Product Requirement Document (PRD): Laboratory Module (LIMS)

This document serves as the absolute product specification for the Laboratory Information Management System (LIMS) Module in Curexal V2.

---

## 1. Module Overview & Purpose

The Laboratory Module automates diagnostic laboratory operations, including phlebotomy specimen intake, thermal barcode label generation, bench worklist management, analyzer instrument connectivity, scientific validation, consultant approval, PDF report generation, WhatsApp delivery, and ISO 15189 Quality Control.

---

## 2. Actors & User Roles

| Actor | Role | Permissions Required | System Responsibilities |
| :--- | :--- | :--- | :--- |
| **Phlebotomist** | Specimen Collector | `specimens:collect`, `barcodes:print` | Verifies patient, logs container type, generates & prints Code128 thermal barcode labels. |
| **Bench Scientist** | Lab Technologist | `worklists:read`, `results:write` | Operates testing bench, validates analyzer telemetry, inputs parameter values, checks Delta rules. |
| **Consultant Pathologist**| Lab Director | `results:validate`, `reports:release` | Performs dual-authorization sign-off on scientific reports, releases final PDFs. |
| **Quality Officer** | QMS Manager | `qc:manage`, `equipment:log` | Logs QC control lots, monitors Levey-Jennings charts, enforces Westgard rules. |

---

## 3. Laboratory Specimen & Result Lifecycle

```
[Referral Accepted / Walk-in Intake]
                 │
                 ▼
     [Phlebotomy & Barcoding] ──► (Prints Code128 Label: e.g., LAB-99281-EDTA)
                 │
                 ▼
      [Bench / Analyzer Sync] ──► (Ingests ASTM/HL7 strings from Sysmex/Mindray)
                 │
                 ▼
      [Scientific Validation] ──► (Delta Check & Reference Range Validation)
                 │
                 ▼
     [Consultant Sign-off]    ──► (Dual-Authorization Approval)
                 │
                 ▼
     [PDF Report & WhatsApp] ──► (Dispatches to Patient & Referring Doctor)
```

---

## 4. ISO 15189 Quality Control & Westgard Rules

### Quality Control Rules
- **Rule LAB-QC-001**: Every analytical bench (Chemistry, Hematology, Microbiology) must run at least two levels of Quality Control (Normal & Abnormal) before processing clinical patient samples.
- **Rule LAB-QC-002 (Westgard Validation)**:
  - `1-2s` Warning: Control exceeds 2 standard deviations ($2SD$).
  - `1-3s` Reject: Control exceeds 3 standard deviations ($3SD$). Auto-locks analyzer channel.
  - `2-2s` Reject: Two consecutive controls exceed $2SD$.
  - `R-4s` Reject: Range between controls within run exceeds $4SD$.

---

## 5. Analyzer Data Ingestion (Edge Agent Specification)

The **Curexal Edge Agent** listens to laboratory instruments:
- Parses ASTM E1394 and HL7 v2.x result messages over RS-232 serial or TCP connections.
- Matches `Barcode128` string against active `Specimen` registry in local SQLite cache.
- Transmits encrypted JSON payload to cloud workspace over mTLS WebSockets.

---

## 6. Acceptance Criteria

- [x] Barcode labels must generate standard Code128 format compatible with thermal printers.
- [x] Analyzers transmitting data must automatically link results to matching barcode numbers.
- [x] Westgard $1-3s$ failures must immediately lock the testing bench and send alerts to Lab Director.
- [x] Report releases must generate PDF documents stored in MinIO/S3 and dispatch WhatsApp delivery links.
