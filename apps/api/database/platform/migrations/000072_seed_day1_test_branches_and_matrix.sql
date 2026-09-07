-- +goose Up
-- Migration 000072: Seed Day 1 Test Branches and Verification Matrix

-- +goose StatementBegin
DO $$
DECLARE
    v_org_id UUID;
    v_facility_type_clinic_id UUID;
    v_facility_type_specialty_id UUID;
    v_facility_type_diag_id UUID;
    v_branch_main_id UUID := '00000000-0000-0000-0000-000000000001';
    v_branch_ikeja_id UUID := '00000000-0000-0000-0000-000000000002';
    v_branch_lekki_id UUID := '00000000-0000-0000-0000-000000000003';
    
    v_owner_user_id UUID := '00000000-0000-0000-0000-000000000009';
    v_dr_emeka_id UUID := '00000000-0000-0000-0000-000000000010';
    v_dr_sarah_id UUID := '00000000-0000-0000-0000-000000000011';
    v_nurse_chioma_id UUID := '00000000-0000-0000-0000-000000000012';
    v_unassigned_id UUID := '00000000-0000-0000-0000-000000000013';
    v_admin_user_id UUID := '00000000-0000-0000-0000-000000000008';

    v_mem_owner_id UUID;
    v_mem_emeka_id UUID;
    v_mem_sarah_id UUID;
    v_mem_chioma_id UUID;
    v_mem_unassigned_id UUID;
BEGIN
    -- 1. Ensure Organization 'curexal-clinic' exists
    SELECT id INTO v_org_id FROM organization.organizations WHERE slug = 'curexal-clinic' LIMIT 1;
    IF v_org_id IS NULL THEN
        v_org_id := '00000000-0000-0000-0000-000000000001';
        INSERT INTO organization.organizations (
            id, name, slug, status, plan, settings, created_at, updated_at
        ) VALUES (
            v_org_id,
            'Curexal Clinic',
            'curexal-clinic',
            'active',
            'enterprise',
            '{"theme": {"primaryColor": "#0F766E", "fontFamily": "Inter"}, "features": ["clinical", "telehealth", "laboratory", "pharmacy", "radiology"]}'::jsonb,
            NOW(),
            NOW()
        );
    END IF;

    -- 2. Ensure Facility Types
    SELECT id INTO v_facility_type_clinic_id FROM platform.facility_types WHERE code = 'clinic' LIMIT 1;
    IF v_facility_type_clinic_id IS NULL THEN
        INSERT INTO platform.facility_types (code, name, category, icon_key, description)
        VALUES ('clinic', 'Outpatient Clinic', 'clinical', 'Stethoscope', 'Primary care and outpatient clinical practice')
        RETURNING id INTO v_facility_type_clinic_id;
    END IF;

    SELECT id INTO v_facility_type_specialty_id FROM platform.facility_types WHERE code = 'specialty' LIMIT 1;
    IF v_facility_type_specialty_id IS NULL THEN
        INSERT INTO platform.facility_types (code, name, category, icon_key, description)
        VALUES ('specialty', 'Specialty Center', 'clinical', 'Activity', 'Specialist consultation and ambulatory clinical care')
        RETURNING id INTO v_facility_type_specialty_id;
    END IF;

    SELECT id INTO v_facility_type_diag_id FROM platform.facility_types WHERE code = 'diagnostic' LIMIT 1;
    IF v_facility_type_diag_id IS NULL THEN
        INSERT INTO platform.facility_types (code, name, category, icon_key, description)
        VALUES ('diagnostic', 'Diagnostic & Emergency', 'clinical', 'Microscope', 'Urgent care, laboratory, and advanced diagnostics')
        RETURNING id INTO v_facility_type_diag_id;
    END IF;

    -- 3. Seed / Ensure 3 Operational Facility Branches for Curexal Clinic
    -- Branch 1: Main Campus (HQ)
    SELECT id INTO v_branch_main_id FROM organization.facility_branches WHERE organization_id = v_org_id AND code = 'HO-01' LIMIT 1;
    IF v_branch_main_id IS NULL THEN
        v_branch_main_id := '00000000-0000-0000-0000-000000000001';
        INSERT INTO organization.facility_branches (
            id, organization_id, facility_type_id, code, slug, name, is_headquarters,
            email, phone, address, city, state, country, status, created_at, updated_at
        ) VALUES (
            v_branch_main_id,
            v_org_id,
            v_facility_type_clinic_id,
            'HO-01',
            'curexal-clinic',
            'Curexal Clinic Main Campus',
            TRUE,
            'main@curexal.com',
            '+234 800 CUREXAL',
            'Plot 104, Curexal Medical Boulevard, Victoria Island',
            'Lagos',
            'Lagos State',
            'Nigeria',
            'ACTIVE',
            NOW(),
            NOW()
        );
    ELSE
        UPDATE organization.facility_branches SET
            name = 'Curexal Clinic Main Campus',
            code = 'HO-01',
            slug = 'curexal-clinic',
            is_headquarters = TRUE,
            status = 'ACTIVE',
            updated_at = NOW()
        WHERE id = v_branch_main_id;
    END IF;

    -- Branch 2: Ikeja Ambulatory Specialty Center
    SELECT id INTO v_branch_ikeja_id FROM organization.facility_branches WHERE organization_id = v_org_id AND code = 'IKJ-02' LIMIT 1;
    IF v_branch_ikeja_id IS NULL THEN
        v_branch_ikeja_id := '00000000-0000-0000-0000-000000000002';
        INSERT INTO organization.facility_branches (
            id, organization_id, facility_type_id, code, slug, name, is_headquarters,
            email, phone, address, city, state, country, status, created_at, updated_at
        ) VALUES (
            v_branch_ikeja_id,
            v_org_id,
            v_facility_type_specialty_id,
            'IKJ-02',
            'ikeja-specialty',
            'Curexal Ikeja Specialty Center',
            FALSE,
            'ikeja@curexal.com',
            '+234 801 CUREXAL',
            '12 Mobolaji Bank Anthony Way, Ikeja',
            'Lagos',
            'Lagos State',
            'Nigeria',
            'ACTIVE',
            NOW(),
            NOW()
        );
    ELSE
        UPDATE organization.facility_branches SET
            name = 'Curexal Ikeja Specialty Center',
            code = 'IKJ-02',
            slug = 'ikeja-specialty',
            is_headquarters = FALSE,
            status = 'ACTIVE',
            updated_at = NOW()
        WHERE id = v_branch_ikeja_id;
    END IF;

    -- Branch 3: Lekki Diagnostic & Emergency
    SELECT id INTO v_branch_lekki_id FROM organization.facility_branches WHERE organization_id = v_org_id AND code = 'LEK-03' LIMIT 1;
    IF v_branch_lekki_id IS NULL THEN
        v_branch_lekki_id := '00000000-0000-0000-0000-000000000003';
        INSERT INTO organization.facility_branches (
            id, organization_id, facility_type_id, code, slug, name, is_headquarters,
            email, phone, address, city, state, country, status, created_at, updated_at
        ) VALUES (
            v_branch_lekki_id,
            v_org_id,
            v_facility_type_diag_id,
            'LEK-03',
            'lekki-diagnostic',
            'Curexal Lekki Diagnostic & Emergency',
            FALSE,
            'lekki@curexal.com',
            '+234 802 CUREXAL',
            'Block 5, Admiralty Way, Lekki Phase 1',
            'Lagos',
            'Lagos State',
            'Nigeria',
            'ACTIVE',
            NOW(),
            NOW()
        );
    ELSE
        UPDATE organization.facility_branches SET
            name = 'Curexal Lekki Diagnostic & Emergency',
            code = 'LEK-03',
            slug = 'lekki-diagnostic',
            is_headquarters = FALSE,
            status = 'ACTIVE',
            updated_at = NOW()
        WHERE id = v_branch_lekki_id;
    END IF;

    -- 4. Ensure Users
    -- Owner
    INSERT INTO identity.users (id, email, name, email_verified, is_platform_admin, credential_status, created_at, updated_at)
    VALUES (v_owner_user_id, 'owner@curexal.com', 'Curexal Clinic Owner', TRUE, FALSE, 'ACTIVE', NOW(), NOW())
    ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name, email_verified = TRUE, credential_status = 'ACTIVE', updated_at = NOW();

    -- Dr. Emeka (Multi-Branch Doctor)
    INSERT INTO identity.users (id, email, name, email_verified, is_platform_admin, credential_status, created_at, updated_at)
    VALUES (v_dr_emeka_id, 'dr.emeka@curexal.com', 'Dr. Emeka Okonkwo', TRUE, FALSE, 'ACTIVE', NOW(), NOW())
    ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name, email_verified = TRUE, credential_status = 'ACTIVE', updated_at = NOW();

    -- Dr. Sarah (Single-Branch Doctor)
    INSERT INTO identity.users (id, email, name, email_verified, is_platform_admin, credential_status, created_at, updated_at)
    VALUES (v_dr_sarah_id, 'dr.sarah@curexal.com', 'Dr. Sarah Alabi', TRUE, FALSE, 'ACTIVE', NOW(), NOW())
    ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name, email_verified = TRUE, credential_status = 'ACTIVE', updated_at = NOW();

    -- Nurse Chioma (Single-Branch Triage Nurse)
    INSERT INTO identity.users (id, email, name, email_verified, is_platform_admin, credential_status, created_at, updated_at)
    VALUES (v_nurse_chioma_id, 'nurse.chioma@curexal.com', 'Nurse Chioma Eze', TRUE, FALSE, 'ACTIVE', NOW(), NOW())
    ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name, email_verified = TRUE, credential_status = 'ACTIVE', updated_at = NOW();

    -- Unassigned Staff (Zero assigned branches for HIPAA guard test)
    INSERT INTO identity.users (id, email, name, email_verified, is_platform_admin, credential_status, created_at, updated_at)
    VALUES (v_unassigned_id, 'unassigned.staff@curexal.com', 'Provisional Staff Member', TRUE, FALSE, 'ACTIVE', NOW(), NOW())
    ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name, email_verified = TRUE, credential_status = 'ACTIVE', updated_at = NOW();

    -- Platform Super Admin
    INSERT INTO identity.users (id, email, name, email_verified, is_platform_admin, credential_status, created_at, updated_at)
    VALUES (v_admin_user_id, 'admin@curexal.com', 'Curexal Super Administrator', TRUE, TRUE, 'ACTIVE', NOW(), NOW())
    ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name, email_verified = TRUE, is_platform_admin = TRUE, credential_status = 'ACTIVE', updated_at = NOW();

    -- Resolve actual IDs in case of pre-existing records with different UUIDs
    SELECT id INTO v_owner_user_id FROM identity.users WHERE email = 'owner@curexal.com';
    SELECT id INTO v_dr_emeka_id FROM identity.users WHERE email = 'dr.emeka@curexal.com';
    SELECT id INTO v_dr_sarah_id FROM identity.users WHERE email = 'dr.sarah@curexal.com';
    SELECT id INTO v_nurse_chioma_id FROM identity.users WHERE email = 'nurse.chioma@curexal.com';
    SELECT id INTO v_unassigned_id FROM identity.users WHERE email = 'unassigned.staff@curexal.com';
    SELECT id INTO v_admin_user_id FROM identity.users WHERE email = 'admin@curexal.com';

    -- 5. Seed Credentials (Password: 'password')
    INSERT INTO identity.credentials (id, account_id, auth_provider, user_id, password_hash, created_at, updated_at)
    SELECT gen_random_uuid(), u.email, 'credential', u.id, '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW()
    FROM identity.users u
    WHERE u.email IN ('owner@curexal.com', 'dr.emeka@curexal.com', 'dr.sarah@curexal.com', 'nurse.chioma@curexal.com', 'unassigned.staff@curexal.com', 'admin@curexal.com')
    ON CONFLICT (user_id, auth_provider) DO UPDATE SET
        password_hash = EXCLUDED.password_hash,
        account_id = EXCLUDED.account_id,
        updated_at = NOW();

    -- 6. Seed / Update Canonical Organization Memberships
    -- Owner
    INSERT INTO organization.organization_memberships (id, user_id, organization_id, role, role_title, is_active, joined_at, created_at)
    VALUES (gen_random_uuid(), v_owner_user_id, v_org_id, 'owner', 'Clinic Owner & Managing Director', TRUE, NOW(), NOW())
    ON CONFLICT (organization_id, user_id) DO UPDATE SET role = 'owner', role_title = EXCLUDED.role_title, is_active = TRUE;

    -- Dr. Emeka (clinic_doctor)
    INSERT INTO organization.organization_memberships (id, user_id, organization_id, role, role_title, is_active, joined_at, created_at)
    VALUES (gen_random_uuid(), v_dr_emeka_id, v_org_id, 'clinic_doctor', 'Chief Medical Officer / Attending Physician', TRUE, NOW(), NOW())
    ON CONFLICT (organization_id, user_id) DO UPDATE SET role = 'clinic_doctor', role_title = EXCLUDED.role_title, is_active = TRUE;

    -- Dr. Sarah (clinic_doctor)
    INSERT INTO organization.organization_memberships (id, user_id, organization_id, role, role_title, is_active, joined_at, created_at)
    VALUES (gen_random_uuid(), v_dr_sarah_id, v_org_id, 'clinic_doctor', 'Consultant Specialist', TRUE, NOW(), NOW())
    ON CONFLICT (organization_id, user_id) DO UPDATE SET role = 'clinic_doctor', role_title = EXCLUDED.role_title, is_active = TRUE;

    -- Nurse Chioma (clinic_nurse)
    INSERT INTO organization.organization_memberships (id, user_id, organization_id, role, role_title, is_active, joined_at, created_at)
    VALUES (gen_random_uuid(), v_nurse_chioma_id, v_org_id, 'clinic_nurse', 'Lead Triage Nurse', TRUE, NOW(), NOW())
    ON CONFLICT (organization_id, user_id) DO UPDATE SET role = 'clinic_nurse', role_title = EXCLUDED.role_title, is_active = TRUE;

    -- Unassigned Staff (clinic_nurse with NO branch assignment)
    INSERT INTO organization.organization_memberships (id, user_id, organization_id, role, role_title, is_active, joined_at, created_at)
    VALUES (gen_random_uuid(), v_unassigned_id, v_org_id, 'clinic_nurse', 'Provisional Care Coordinator', TRUE, NOW(), NOW())
    ON CONFLICT (organization_id, user_id) DO UPDATE SET role = 'clinic_nurse', role_title = EXCLUDED.role_title, is_active = TRUE;

    -- Retrieve membership IDs
    SELECT id INTO v_mem_owner_id FROM organization.organization_memberships WHERE organization_id = v_org_id AND user_id = v_owner_user_id LIMIT 1;
    SELECT id INTO v_mem_emeka_id FROM organization.organization_memberships WHERE organization_id = v_org_id AND user_id = v_dr_emeka_id LIMIT 1;
    SELECT id INTO v_mem_sarah_id FROM organization.organization_memberships WHERE organization_id = v_org_id AND user_id = v_dr_sarah_id LIMIT 1;
    SELECT id INTO v_mem_chioma_id FROM organization.organization_memberships WHERE organization_id = v_org_id AND user_id = v_nurse_chioma_id LIMIT 1;
    SELECT id INTO v_mem_unassigned_id FROM organization.organization_memberships WHERE organization_id = v_org_id AND user_id = v_unassigned_id LIMIT 1;

    -- 7. Multi-Branch Assignments in organization.membership_branches
    -- Dr. Emeka -> Main Campus (HO-01) and Ikeja Specialty Center (IKJ-02) [Multi-Branch Doctor]
    DELETE FROM organization.membership_branches WHERE membership_id = v_mem_emeka_id;
    INSERT INTO organization.membership_branches (id, membership_id, facility_branch_id, created_at, created_by)
    VALUES 
        (gen_random_uuid(), v_mem_emeka_id, v_branch_main_id, NOW(), v_owner_user_id),
        (gen_random_uuid(), v_mem_emeka_id, v_branch_ikeja_id, NOW(), v_owner_user_id)
    ON CONFLICT (membership_id, facility_branch_id) DO NOTHING;

    -- Dr. Sarah -> Main Campus (HO-01) ONLY [Single-Branch Doctor]
    DELETE FROM organization.membership_branches WHERE membership_id = v_mem_sarah_id;
    INSERT INTO organization.membership_branches (id, membership_id, facility_branch_id, created_at, created_by)
    VALUES 
        (gen_random_uuid(), v_mem_sarah_id, v_branch_main_id, NOW(), v_owner_user_id)
    ON CONFLICT (membership_id, facility_branch_id) DO NOTHING;

    -- Nurse Chioma -> Main Campus (HO-01) ONLY [Single-Branch Nurse]
    DELETE FROM organization.membership_branches WHERE membership_id = v_mem_chioma_id;
    INSERT INTO organization.membership_branches (id, membership_id, facility_branch_id, created_at, created_by)
    VALUES 
        (gen_random_uuid(), v_mem_chioma_id, v_branch_main_id, NOW(), v_owner_user_id)
    ON CONFLICT (membership_id, facility_branch_id) DO NOTHING;

    -- Unassigned Staff -> Explicitly Zero Branches (Ensure clean state)
    DELETE FROM organization.membership_branches WHERE membership_id = v_mem_unassigned_id;

    -- 8. Seed Sample Immutable Audit Ledger Entries for Feature #24 Verification
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'audit' AND table_name = 'audit_events') THEN
        INSERT INTO audit.audit_events (
            id, organization_id, facility_branch_id, event_category, action,
            actor_id, actor_role, actor_email, resource_type, resource_id,
            reason, ip_address, status, severity, is_break_glass,
            prev_record_hash, record_hash, created_at
        ) VALUES (
            gen_random_uuid(),
            v_org_id,
            v_branch_main_id,
            'FACILITY_PROVISIONED',
            'facility:branch:create',
            v_owner_user_id,
            'STAFF',
            'owner@curexal.com',
            'facility_branch',
            v_branch_main_id::text,
            'Operational facility branch HO-01 provisioned with HIPAA ledger genesis block',
            '127.0.0.1',
            'SUCCESS',
            'INFO',
            FALSE,
            '0000000000000000000000000000000000000000000000000000000000000000',
            encode(sha256('GENESIS_BLOCK_HO01'::bytea), 'hex'),
            NOW() - INTERVAL '1 day'
        ),
        (
            gen_random_uuid(),
            v_org_id,
            v_branch_ikeja_id,
            'FACILITY_PROVISIONED',
            'facility:branch:create',
            v_owner_user_id,
            'STAFF',
            'owner@curexal.com',
            'facility_branch',
            v_branch_ikeja_id::text,
            'Operational specialty facility branch IKJ-02 provisioned for ambulatory care',
            '127.0.0.1',
            'SUCCESS',
            'INFO',
            FALSE,
            encode(sha256('GENESIS_BLOCK_HO01'::bytea), 'hex'),
            encode(sha256('GENESIS_BLOCK_IKJ02'::bytea), 'hex'),
            NOW() - INTERVAL '12 hours'
        ) ON CONFLICT DO NOTHING;
    END IF;

END $$;
-- +goose StatementEnd

-- +goose Down
-- Reversible rollback
