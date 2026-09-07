-- +goose Up
-- ==============================================================================
-- MIGRATION 000067: Canonical Workspace Navigation SSOT for Clinics & Facilities
-- ==============================================================================

-- 1. Remove legacy duplicate items
DELETE FROM navigation_item WHERE id IN ('nav_wsp_patients', 'nav_wsp_settings');

-- 2. Update and align canonical workspace navigation items with SSOT capabilities and permissions
INSERT INTO navigation_item (
    id, key, context_scope, module_code, title, description, icon, path, sort_order, required_permission, required_capability, status, is_visible, is_active
) VALUES
    ('nav_wsp_dashboard', 'workspace.dashboard', 'workspace', 'dashboard', 'Workspace Overview', 'Facility operational metrics, patient throughput and active alerts', 'LayoutDashboard', '/:branch/dashboard', 1, NULL, NULL, 'active', true, true),
    ('nav_wsp_reception', 'workspace.reception', 'workspace', 'reception', 'Patient Intake (MPI)', 'Master Patient Index, new patient registration & check-in queue', 'Users', '/:branch/reception', 2, 'workspace:patient:read', 'core.patient', 'active', true, true),
    ('nav_wsp_care_desk', 'workspace.care_desk', 'workspace', 'care_desk', 'Care & Triage Desk', 'Acuity vitals triage and queue orchestration', 'Activity', '/:branch/care-desk', 3, 'workspace:patient:read', 'clinical.basic', 'active', true, true),
    ('nav_wsp_clinical', 'workspace.clinical', 'workspace', 'clinical', 'Clinical & EMR', 'Doctor consultations, clinical SOAP notes, prescriptions & care plans', 'Stethoscope', '/:branch/clinical', 4, 'workspace:clinical:read', 'clinical.basic', 'active', true, true),
    ('nav_wsp_laboratory', 'workspace.laboratory', 'workspace', 'laboratory', 'Laboratory (LIS)', 'Specimen accessioning, worklists & diagnostic result authorization', 'Microscope', '/:branch/laboratory', 5, 'workspace:sample:receive', 'laboratory.basic', 'active', true, true),
    ('nav_wsp_pharmacy', 'workspace.pharmacy', 'workspace', 'pharmacy', 'Pharmacy Dispensary', 'Prescription fulfillment, FEFO batch dispensing and drug inventory', 'Pill', '/:branch/pharmacy', 6, 'workspace:clinical:read', 'pharmacy.basic', 'active', true, true),
    ('nav_wsp_radiology', 'workspace.radiology', 'workspace', 'radiology', 'Radiology (RIS & PACS)', 'Diagnostic imaging worklists, DICOM viewing & radiology reports', 'Layers', '/:branch/radiology', 7, 'workspace:clinical:read', 'radiology.basic', 'active', true, true),
    ('nav_wsp_hospital', 'workspace.hospital', 'workspace', 'hospital', 'Inpatient Hospital (HIS)', 'Inpatient ward admissions, bed tracking and clinical handovers', 'Building2', '/:branch/hospital', 8, 'workspace:clinical:read', 'clinical.inpatient_wards', 'active', true, true),
    ('nav_wsp_billing', 'workspace.billing', 'workspace', 'billing', 'Billing & Cashier POS', 'POS invoice generation, payment collections, insurance & receipts', 'CreditCard', '/:branch/billing', 9, 'workspace:billing:create', 'core.billing', 'active', true, true)
ON CONFLICT (id) DO UPDATE SET
    key = EXCLUDED.key,
    context_scope = EXCLUDED.context_scope,
    module_code = EXCLUDED.module_code,
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    icon = EXCLUDED.icon,
    path = EXCLUDED.path,
    sort_order = EXCLUDED.sort_order,
    required_permission = EXCLUDED.required_permission,
    required_capability = EXCLUDED.required_capability,
    status = EXCLUDED.status,
    is_visible = EXCLUDED.is_visible,
    is_active = EXCLUDED.is_active;

-- 3. Ensure branch_admin role has all required workspace permissions
INSERT INTO "authorization".role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM "authorization".roles r
CROSS JOIN "authorization".permissions p
WHERE r.code = 'branch_admin'
  AND p.code IN (
      'workspace:patient:read',
      'workspace:patient:create',
      'workspace:clinical:read',
      'workspace:billing:create',
      'workspace:settings:manage',
      'workspace:sample:receive',
      'workspace:result:authorize',
      'workspace:worksheet:update'
  )
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM navigation_item WHERE id IN ('nav_wsp_dashboard', 'nav_wsp_reception', 'nav_wsp_care_desk', 'nav_wsp_clinical', 'nav_wsp_laboratory', 'nav_wsp_pharmacy', 'nav_wsp_radiology', 'nav_wsp_hospital', 'nav_wsp_billing');
