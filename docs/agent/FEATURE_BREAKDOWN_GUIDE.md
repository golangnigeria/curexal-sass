# Feature Breakdown & Epic Decomposition Guide

> **Purpose**: Guidelines for breaking down large enterprise healthcare epics into single-turn development tasks.  
> **Owner**: Lead Systems Architect  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. Hierarchy of Task Decomposition

```text
Large Healthcare Epic (e.g., LIMS Laboratory Specimen & Analyzer Integration)
                  │
                  ▼
Capability 1: Specimen Accessioning & Barcode Generation
                  │
                  ▼
Feature 1.1: Specimen Accessioning Backend DDL & Domain Service
                  │
                  ▼
Use Case 1.1.1: Create Specimen Accession Record (Task)
                  │
                  ▼
Single-Turn Pull Request Implementation
```

---

## 2. Real Curexal Module Decomposition Examples

### Example A: LIMS Module (Laboratory Information System)
- **Epic**: Diagnostic Specimen Processing & Result Verification Engine.
- **Task Breakdown**:
  1. `Task 1 (Backend DDL & Models)`: Create `specimens` and `test_orders` tables in `tenant_<slug>` schema.
  2. `Task 2 (Repository & Use Case)`: Build `BunSpecimenRepository` and `AccessionSpecimenUseCase`.
  3. `Task 3 (Hertz REST Handler)`: Wire `POST /api/v1/lims/specimens/accession` with RFC7807 problem details.
  4. `Task 4 (API SDK)`: Add `limsApi.accessionSpecimen()` to `@curexal/api-sdk`.
  5. `Task 5 (React UI)`: Build `SpecimenAccessionForm.tsx` in `apps/web-workspace` with barcode generation.

### Example B: RIS Module (Radiology Information System)
- **Epic**: DICOM Worklist & Image Viewer Integration.
- **Task Breakdown**:
  1. `Task 1 (Backend)`: Create `radiology_studies` schema, MinIO S3 DICOM file attachment handler.
  2. `Task 2 (Frontend)`: Integrate CornerstoneJS / OHIF viewer component in `@curexal/ui-healthcare`.
