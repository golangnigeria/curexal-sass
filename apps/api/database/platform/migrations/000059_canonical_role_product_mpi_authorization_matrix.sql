-- +goose Up
-- ==============================================================================
-- MIGRATION 000059: Canonical Role, Product & Shared Core/MPI Authorization Matrix
-- ==============================================================================

-- 1. Ensure columns and constraints on "authorization".roles
ALTER TABLE "authorization".roles ADD COLUMN IF NOT EXISTS code VARCHAR(100);
ALTER TABLE "authorization".roles ADD COLUMN IF NOT EXISTS context_scope VARCHAR(50) DEFAULT 'workspace';
ALTER TABLE "authorization".roles ADD COLUMN IF NOT EXISTS description TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_authorization_roles_code ON "authorization".roles(code);

-- 2. Seed All Canonical Roles across Scopes (Platform, Organization, Branch, Products)
INSERT INTO "authorization".roles (name, code, context_scope, description)
VALUES
    -- Platform Control Plane
    ('Platform Super Admin', 'super_admin', 'platform', 'Root operational authority over platform infrastructure'),
    ('Platform Administrator', 'platform_admin', 'platform', 'Platform tenant governance and subscription management'),
    ('Platform Staff', 'platform_staff', 'platform', 'Platform operational staff and diagnostics observer'),
    ('Platform Support Specialist', 'platform_support', 'platform', 'Customer support and diagnostic troubleshooting'),
    ('Platform Security Auditor', 'platform_security', 'platform', 'Security compliance and audit inspection'),
    ('Platform Billing Officer', 'platform_billing', 'platform', 'Platform revenue, billing tiers and gateway reconciliation'),
    ('Platform Compliance Officer', 'platform_compliance', 'platform', 'Accreditation vetting and legal review'),

    -- Organization Control Plane (Executive HQ)
    ('Organization Owner', 'owner', 'organization', 'Owner of the healthcare network with full administrative control'),
    ('Organization Administrator', 'org_admin', 'organization', 'Enterprise administrator managing multi-branch facilities'),
    ('Regional Operations Manager', 'org_regional_manager', 'organization', 'Oversees operational performance across regional facilities'),
    ('Quality & Compliance Manager', 'org_quality_manager', 'organization', 'Enforces corporate clinical accreditation and audits'),
    ('Finance & Billing Manager', 'org_finance_manager', 'organization', 'Corporate financial management and subscription billing'),
    ('Human Resources Manager', 'org_hr_manager', 'organization', 'Staff roster onboarding, credentialing and access control'),

    -- Branch Facility Operations
    ('Branch Operations Manager', 'branch_manager', 'workspace', 'Facility-level operations, staff scheduling and throughput'),
    ('Branch Facility Administrator', 'branch_admin', 'workspace', 'Branch facility administration and configuration'),
    ('Front Desk Receptionist', 'receptionist', 'workspace', 'Patient intake, registration, check-in and appointments'),
    ('Cashier & Billing Officer', 'cashier', 'workspace', 'POS payment collection, invoices and receipts'),

    -- HMS / Clinical Practice
    ('Attending Physician (Doctor)', 'doctor', 'workspace', 'Consultations, clinical notes, prescriptions and lab orders'),
    ('Attending Physician', 'clinician', 'workspace', 'Attending physician / doctor managing clinical patient consultations'),
    ('Clinical Specialist', 'specialist', 'workspace', 'Specialized patient care and medical consultations'),
    ('Registered Staff Nurse', 'nurse', 'workspace', 'Vital signs, triage acuity, care plans and nursing notes'),
    ('Lead Triage Nurse', 'triage_nurse', 'workspace', 'Auto-acuity triage and care desk doctor allocation'),
    ('Clinical Officer', 'clinical_officer', 'workspace', 'General medical consultations and primary care'),
    ('Medical Assistant', 'medical_assistant', 'workspace', 'Clinical intake support, vitals recording and prep'),

    -- LIS / Diagnostic Pathology
    ('Laboratory Director / Manager', 'lab_manager', 'workspace', 'Laboratory operations, QC review, validation and instruments'),
    ('Senior Medical Lab Scientist', 'senior_scientist', 'workspace', 'Complex test analysis, result entry and validation'),
    ('Medical Laboratory Scientist', 'scientist', 'workspace', 'Specimen accessioning, testing, result entry and QC read'),
    ('Laboratory Technician', 'technician', 'workspace', 'Diagnostic technologist performing sample accessioning and testing'),
    ('Phlebotomist', 'phlebotomist', 'workspace', 'Sample collection, specimen labeling and tracking'),
    ('Sample Receptionist', 'sample_receptionist', 'workspace', 'Diagnostic specimen intake, barcoding and accessioning'),
    ('Result Validator / Pathologist', 'result_validator', 'workspace', 'Diagnostic review and final result authorization'),

    -- Pharmacy & Dispensary
    ('Pharmacy Manager', 'pharmacy_manager', 'workspace', 'Pharmacy operations, formulary and inventory management'),
    ('Pharmacist', 'pharmacist', 'workspace', 'Prescription fulfillment, medication dispensing and counseling'),
    ('Pharmacy Technician', 'pharmacy_technician', 'workspace', 'Prescription batch preparation and stock management'),
    ('Pharmacy Assistant', 'pharmacy_assistant', 'workspace', 'Counter dispensing support and inventory intake'),
    ('Inventory Officer', 'inventory_officer', 'workspace', 'Batch expiration tracking, FEFO management and orders'),
    ('Storekeeper', 'storekeeper', 'workspace', 'Physical stock warehousing and supply distribution'),

    -- Radiology / RIS / PACS
    ('Consultant Radiologist', 'radiologist', 'workspace', 'Imaging study review, PACS interpretation and reports'),
    ('Radiographer', 'radiographer', 'workspace', 'Modality acquisition (X-Ray, Ultrasound, CT, MRI)'),
    ('Radiology Technician', 'radiology_technician', 'workspace', 'Imaging setup, patient positioning and processing'),
    ('Radiology Manager', 'radiology_manager', 'workspace', 'Imaging department workflows and equipment maintenance'),
    ('Radiology Receptionist', 'radiology_receptionist', 'workspace', 'Imaging appointments, study intake and patient prep'),

    -- HIS / Inpatient Hospital
    ('Ward Manager / Matron', 'ward_manager', 'workspace', 'Inpatient ward management, bed allocation and handovers'),
    ('Inpatient Staff Nurse', 'staff_nurse', 'workspace', 'Ward patient care, medication administration and rounds'),
    ('Ward Clerk', 'ward_clerk', 'workspace', 'Inpatient record keeping and chart maintenance'),
    ('Admission Officer', 'admission_officer', 'workspace', 'Hospital admissions, transfers and discharge processing'),

    -- Billing & Revenue
    ('Billing & Claims Officer', 'billing_officer', 'workspace', 'HMO claims processing, corporate tariff management'),
    ('Finance Officer', 'finance_officer', 'workspace', 'Cash reconciliation, receipts ledger and accounting'),

    -- Patient & Family Care Portal
    ('Patient (Self)', 'patient', 'patient', 'Patient accessing personal health records and care journey'),
    ('Family Guardian', 'guardian', 'patient', 'Guardian authorized for minor or dependent healthcare records'),
    ('Authorized Caregiver', 'caregiver', 'patient', 'Caregiver authorized for designated patient consultations')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    context_scope = EXCLUDED.context_scope,
    description = EXCLUDED.description;

-- 3. Seed Canonical Permissions (Shared Core/MPI + Product Families)
INSERT INTO "authorization".permissions (code, module, category, description)
VALUES
    -- =========================================================================
    -- SHARED CORE / MPI (Master Patient Index)
    -- =========================================================================
    ('core.patient.read', 'core', 'patient', 'View patient demographics and records'),
    ('core.patient.create', 'core', 'patient', 'Register a new patient into the Master Patient Index'),
    ('core.patient.update', 'core', 'patient', 'Update patient demographic and contact details'),
    ('core.patient.search', 'core', 'patient', 'Search the Master Patient Index across branches'),
    ('core.patient.merge', 'core', 'patient', 'Merge duplicate patient MPI identities'),
    ('core.patient.identifiers.read', 'core', 'patient', 'View national, insurance and MRN identifiers'),
    ('core.patient.identifiers.manage', 'core', 'patient', 'Assign and update patient MRN and insurance IDs'),
    ('core.appointment.read', 'core', 'appointments', 'View facility appointment schedule'),
    ('core.appointment.create', 'core', 'appointments', 'Schedule a patient appointment'),
    ('core.appointment.manage', 'core', 'appointments', 'Reschedule or cancel appointments'),

    -- =========================================================================
    -- LIS / DIAGNOSTIC PATHOLOGY
    -- =========================================================================
    ('lis.order.read', 'lis', 'laboratory', 'View laboratory diagnostic orders'),
    ('lis.order.create', 'lis', 'laboratory', 'Create a new laboratory test order'),
    ('lis.sample.read', 'lis', 'laboratory', 'View specimen samples and accession logs'),
    ('lis.sample.collect', 'lis', 'laboratory', 'Record specimen sample collection'),
    ('lis.sample.label', 'lis', 'laboratory', 'Print barcodes and label specimen containers'),
    ('lis.sample.track', 'lis', 'laboratory', 'Track specimen transport and chain of custody'),
    ('lis.sample.receive', 'lis', 'laboratory', 'Receive and verify specimen at lab reception'),
    ('lis.sample.accession', 'lis', 'laboratory', 'Accession specimen and assign lab accession number'),
    ('lis.sample.manage', 'lis', 'laboratory', 'Reject, aliquot or re-route specimen samples'),
    ('lis.result.read', 'lis', 'laboratory', 'View laboratory test results and parameter values'),
    ('lis.result.create', 'lis', 'laboratory', 'Enter laboratory test results on worksheets'),
    ('lis.result.update', 'lis', 'laboratory', 'Modify preliminary laboratory test results'),
    ('lis.result.validate', 'lis', 'laboratory', 'Technically validate laboratory test results'),
    ('lis.result.authorize', 'lis', 'laboratory', 'Clinically authorize and release lab results'),
    ('lis.qc.read', 'lis', 'laboratory', 'View quality control charts (Levey-Jennings)'),
    ('lis.qc.review', 'lis', 'laboratory', 'Review QC runs and Westgard rule violations'),
    ('lis.qc.manage', 'lis', 'laboratory', 'Configure QC lots, target values and standard deviations'),
    ('lis.instrument.read', 'lis', 'laboratory', 'View lab analyzer interfaces and telemetry'),
    ('lis.instrument.manage', 'lis', 'laboratory', 'Configure bidirectional analyzer connections'),
    ('lis.worklist.manage', 'lis', 'laboratory', 'Manage laboratory batch worksheets and worklists'),
    ('lis.settings.manage', 'lis', 'laboratory', 'Manage laboratory reference ranges and panic thresholds'),

    -- =========================================================================
    -- HMS / CLINICAL & OUTPATIENT EMR
    -- =========================================================================
    ('hms.consultation.read', 'hms', 'clinical', 'View clinical consultation records'),
    ('hms.consultation.create', 'hms', 'clinical', 'Conduct and document patient clinical consultation'),
    ('hms.consultation.update', 'hms', 'clinical', 'Amend clinical encounter notes'),
    ('hms.vitals.read', 'hms', 'clinical', 'View patient vital signs and biometric history'),
    ('hms.vitals.record', 'hms', 'clinical', 'Record vital signs (BP, Pulse, Temp, SpO2, BMI)'),
    ('hms.diagnosis.create', 'hms', 'clinical', 'Record ICD-10 clinical diagnoses and problem list'),
    ('hms.prescription.create', 'hms', 'clinical', 'Issue outpatient e-prescriptions'),
    ('hms.lab_order.create', 'hms', 'clinical', 'Order diagnostic laboratory investigations'),
    ('hms.radiology_order.create', 'hms', 'clinical', 'Order diagnostic radiology imaging scans'),
    ('hms.nursing_notes.create', 'hms', 'clinical', 'Record nursing observations and shift notes'),
    ('hms.care_plan.read', 'hms', 'clinical', 'View multidisciplinary patient care plans'),
    ('hms.care_plan.update', 'hms', 'clinical', 'Update care plans and nursing interventions'),
    ('hms.checkin.create', 'hms', 'clinical', 'Check in patient into facility waiting room'),
    ('hms.triage.evaluate', 'hms', 'clinical', 'Perform emergency triage acuity scoring'),

    -- =========================================================================
    -- PHARMACY & DISPENSARY
    -- =========================================================================
    ('pharmacy.prescription.read', 'pharmacy', 'pharmacy', 'View prescription queue and medication orders'),
    ('pharmacy.dispense.create', 'pharmacy', 'pharmacy', 'Dispense medications against valid prescription'),
    ('pharmacy.inventory.read', 'pharmacy', 'pharmacy', 'View pharmacy inventory and stock levels'),
    ('pharmacy.inventory.adjust', 'pharmacy', 'pharmacy', 'Perform stock adjustments and batch reconciliations'),
    ('pharmacy.pricing.manage', 'pharmacy', 'pharmacy', 'Manage medication tariffs and pricing markups'),
    ('pharmacy.medication.read', 'pharmacy', 'pharmacy', 'View drug formulary and interactions'),
    ('pharmacy.purchase.manage', 'pharmacy', 'pharmacy', 'Issue purchase orders to pharmaceutical suppliers'),

    -- =========================================================================
    -- RIS & PACS / RADIOLOGY
    -- =========================================================================
    ('radiology.worklist.read', 'radiology', 'radiology', 'View modality worklist (MWL)'),
    ('radiology.study.read', 'radiology', 'radiology', 'View imaging studies and DICOM metadata'),
    ('radiology.study.create', 'radiology', 'radiology', 'Acquire and upload imaging study to PACS'),
    ('radiology.report.create', 'radiology', 'radiology', 'Author diagnostic radiology report'),
    ('radiology.report.authorize', 'radiology', 'radiology', 'Authorize and publish diagnostic imaging report'),
    ('radiology.pacs.view', 'radiology', 'radiology', 'Launch web PACS DICOM image viewer'),

    -- =========================================================================
    -- HIS / INPATIENT HOSPITAL
    -- =========================================================================
    ('his.ward.read', 'his', 'hospital', 'View inpatient ward occupancy and bed status'),
    ('his.bed.manage', 'his', 'hospital', 'Allocate, transfer or release hospital beds'),
    ('his.admission.create', 'his', 'hospital', 'Process inpatient admission from ER or clinic'),
    ('his.discharge.create', 'his', 'hospital', 'Process patient discharge and summary'),

    -- =========================================================================
    -- BILLING & FINANCIAL POS
    -- =========================================================================
    ('billing.invoice.read', 'billing', 'finance', 'View patient invoices and billing statements'),
    ('billing.invoice.create', 'billing', 'finance', 'Generate patient invoices for services'),
    ('billing.payment.create', 'billing', 'finance', 'Collect cash, card or transfer payments at POS'),
    ('billing.receipt.create', 'billing', 'finance', 'Issue official payment receipts'),
    ('billing.refund.request', 'billing', 'finance', 'Initiate customer refund request'),
    ('billing.branch.read', 'billing', 'finance', 'View branch daily cashier reconciliation reports'),

    -- =========================================================================
    -- BRANCH FACILITY GOVERNANCE
    -- =========================================================================
    ('branch.dashboard.read', 'branch', 'operations', 'View branch operational dashboard'),
    ('branch.staff.read', 'branch', 'operations', 'View branch staff roster and duty shifts'),
    ('branch.operations.read', 'branch', 'operations', 'Inspect branch patient queue and turnaround time'),
    ('branch.reports.read', 'branch', 'operations', 'Generate branch operational analytics')
ON CONFLICT (code) DO UPDATE SET
    module = EXCLUDED.module,
    category = EXCLUDED.category,
    description = EXCLUDED.description;

-- 4. Seed Canonical Role-Permission Matrix (DB SSOT)
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM (VALUES
    -- DOCTOR (Physician): Shared Core + HMS Clinical + Lab Orders
    ('doctor', 'core.patient.read'),
    ('doctor', 'core.patient.search'),
    ('doctor', 'hms.consultation.read'),
    ('doctor', 'hms.consultation.create'),
    ('doctor', 'hms.consultation.update'),
    ('doctor', 'hms.vitals.read'),
    ('doctor', 'hms.vitals.record'),
    ('doctor', 'hms.diagnosis.create'),
    ('doctor', 'hms.prescription.create'),
    ('doctor', 'hms.lab_order.create'),
    ('doctor', 'hms.radiology_order.create'),
    ('doctor', 'lis.order.read'),
    ('doctor', 'lis.result.read'),

    -- NURSE: Shared Core + Triage + Vitals + Nursing Notes
    ('nurse', 'core.patient.read'),
    ('nurse', 'core.patient.search'),
    ('nurse', 'hms.vitals.read'),
    ('nurse', 'hms.vitals.record'),
    ('nurse', 'hms.nursing_notes.create'),
    ('nurse', 'hms.care_plan.read'),
    ('nurse', 'hms.care_plan.update'),
    ('nurse', 'hms.triage.evaluate'),

    -- RECEPTIONIST: Shared Core (Full Intake) + Appointments + Check-in
    ('receptionist', 'core.patient.read'),
    ('receptionist', 'core.patient.create'),
    ('receptionist', 'core.patient.update'),
    ('receptionist', 'core.patient.search'),
    ('receptionist', 'core.appointment.read'),
    ('receptionist', 'core.appointment.create'),
    ('receptionist', 'core.appointment.manage'),
    ('receptionist', 'hms.checkin.create'),

    -- SCIENTIST: Shared Core (MPI Intake & Search) + LIS Testing
    ('scientist', 'core.patient.read'),
    ('scientist', 'core.patient.create'),
    ('scientist', 'core.patient.search'),
    ('scientist', 'lis.order.read'),
    ('scientist', 'lis.sample.read'),
    ('scientist', 'lis.sample.receive'),
    ('scientist', 'lis.sample.accession'),
    ('scientist', 'lis.result.read'),
    ('scientist', 'lis.result.create'),
    ('scientist', 'lis.qc.read'),

    -- PHLEBOTOMIST: Shared Core Search + Specimen Collection & Barcoding
    ('phlebotomist', 'core.patient.read'),
    ('phlebotomist', 'core.patient.search'),
    ('phlebotomist', 'lis.order.read'),
    ('phlebotomist', 'lis.sample.read'),
    ('phlebotomist', 'lis.sample.collect'),
    ('phlebotomist', 'lis.sample.label'),
    ('phlebotomist', 'lis.sample.track'),

    -- LAB MANAGER: Shared Core + Full LIS Management & Result Authorization
    ('lab_manager', 'core.patient.read'),
    ('lab_manager', 'core.patient.create'),
    ('lab_manager', 'core.patient.search'),
    ('lab_manager', 'lis.order.read'),
    ('lab_manager', 'lis.order.create'),
    ('lab_manager', 'lis.sample.read'),
    ('lab_manager', 'lis.sample.manage'),
    ('lab_manager', 'lis.result.read'),
    ('lab_manager', 'lis.result.validate'),
    ('lab_manager', 'lis.result.authorize'),
    ('lab_manager', 'lis.qc.read'),
    ('lab_manager', 'lis.qc.review'),
    ('lab_manager', 'lis.qc.manage'),
    ('lab_manager', 'lis.instrument.read'),
    ('lab_manager', 'lis.instrument.manage'),
    ('lab_manager', 'lis.worklist.manage'),
    ('lab_manager', 'lis.settings.manage'),

    -- PHARMACIST: Shared Core + Prescription Queue + Dispensing + Inventory
    ('pharmacist', 'core.patient.read'),
    ('pharmacist', 'core.patient.search'),
    ('pharmacist', 'pharmacy.prescription.read'),
    ('pharmacist', 'pharmacy.dispense.create'),
    ('pharmacist', 'pharmacy.inventory.read'),
    ('pharmacist', 'pharmacy.medication.read'),

    -- CASHIER: Shared Core Search + Billing POS
    ('cashier', 'core.patient.read'),
    ('cashier', 'core.patient.search'),
    ('cashier', 'billing.invoice.read'),
    ('cashier', 'billing.payment.create'),
    ('cashier', 'billing.receipt.create'),

    -- BRANCH MANAGER: Shared Core + Branch Operations + Branch Billing
    ('branch_manager', 'core.patient.read'),
    ('branch_manager', 'core.patient.search'),
    ('branch_manager', 'branch.dashboard.read'),
    ('branch_manager', 'branch.staff.read'),
    ('branch_manager', 'branch.operations.read'),
    ('branch_manager', 'branch.reports.read'),
    ('branch_manager', 'billing.branch.read'),

    -- RADIOLOGIST: Shared Core + PACS + Study Reports
    ('radiologist', 'core.patient.read'),
    ('radiologist', 'core.patient.search'),
    ('radiologist', 'radiology.worklist.read'),
    ('radiologist', 'radiology.study.read'),
    ('radiologist', 'radiology.report.create'),
    ('radiologist', 'radiology.report.authorize'),
    ('radiologist', 'radiology.pacs.view')
) AS mapping(role_code, perm_code)
JOIN "authorization".roles r ON r.code = mapping.role_code
JOIN "authorization".permissions p ON p.code = mapping.perm_code
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- +goose Down
-- Revert migration 000059 additions
DELETE FROM "authorization".permissions WHERE module IN ('core', 'lis', 'hms', 'pharmacy', 'radiology', 'his', 'branch');
