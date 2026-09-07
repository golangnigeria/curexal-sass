-- +goose Up
-- Migration 000050: Seed Curexal Clinic Enterprise Owner, Plan & Staff Credentials

-- ============================================================================
-- 1. Schema Constraint Alignment on patient.patients
-- ============================================================================
ALTER TABLE patient.patients ALTER COLUMN organization_id DROP NOT NULL;
ALTER TABLE patient.patients ALTER COLUMN dob DROP NOT NULL;

-- ============================================================================
-- 2. Seed Default Staff & Owner User Identities
-- ============================================================================
-- Primary Organization Owner: owner@curexal.com
INSERT INTO identity.users (
    id, email, name, email_verified, is_platform_admin, credential_status, created_at, updated_at
) VALUES (
    '00000000-0000-0000-0000-000000000009',
    'owner@curexal.com',
    'Curexal Clinic Owner',
    TRUE,
    FALSE,
    'ACTIVE',
    NOW(),
    NOW()
) ON CONFLICT (email) DO UPDATE SET 
    name = EXCLUDED.name, 
    email_verified = TRUE,
    credential_status = 'ACTIVE',
    updated_at = NOW();

-- Dr. Emeka Okonkwo (Chief Medical Officer / Telehealth Doctor)
INSERT INTO identity.users (
    id, email, name, email_verified, is_platform_admin, credential_status, created_at, updated_at
) VALUES (
    '00000000-0000-0000-0000-000000000010',
    'dr.emeka@curexal.com',
    'Dr. Emeka Okonkwo',
    TRUE,
    FALSE,
    'ACTIVE',
    NOW(),
    NOW()
) ON CONFLICT (email) DO UPDATE SET 
    name = EXCLUDED.name, 
    email_verified = TRUE,
    credential_status = 'ACTIVE',
    updated_at = NOW();

-- Dr. Sarah Adebayo (Consultant Specialist / Internal Medicine)
INSERT INTO identity.users (
    id, email, name, email_verified, is_platform_admin, credential_status, created_at, updated_at
) VALUES (
    '00000000-0000-0000-0000-000000000011',
    'dr.sarah@curexal.com',
    'Dr. Sarah Adebayo',
    TRUE,
    FALSE,
    'ACTIVE',
    NOW(),
    NOW()
) ON CONFLICT (email) DO UPDATE SET 
    name = EXCLUDED.name, 
    email_verified = TRUE,
    credential_status = 'ACTIVE',
    updated_at = NOW();

-- Nurse Chioma Eze (Lead Triage & Care Coordinator)
INSERT INTO identity.users (
    id, email, name, email_verified, is_platform_admin, credential_status, created_at, updated_at
) VALUES (
    '00000000-0000-0000-0000-000000000012',
    'nurse.chioma@curexal.com',
    'Nurse Chioma Eze',
    TRUE,
    FALSE,
    'ACTIVE',
    NOW(),
    NOW()
) ON CONFLICT (email) DO UPDATE SET 
    name = EXCLUDED.name, 
    email_verified = TRUE,
    credential_status = 'ACTIVE',
    updated_at = NOW();

-- Set Password Credentials for all seeded users (Password: 'password')
-- bcrypt hash: $2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi
INSERT INTO identity.credentials (id, account_id, auth_provider, user_id, password_hash, created_at, updated_at)
SELECT gen_random_uuid(), u.email, 'credential', u.id, '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW()
FROM identity.users u
WHERE u.email IN ('owner@curexal.com', 'dr.emeka@curexal.com', 'dr.sarah@curexal.com', 'nurse.chioma@curexal.com')
ON CONFLICT (user_id, auth_provider) DO UPDATE SET 
    password_hash = EXCLUDED.password_hash,
    account_id = EXCLUDED.account_id,
    updated_at = NOW();

-- ============================================================================
-- 3. Dynamic Seed of Organization, Facility Branch, Memberships & Providers
-- ============================================================================
-- +goose StatementBegin
DO $$
DECLARE
    v_org_id UUID;
    v_branch_id UUID;
    v_facility_type_id UUID;
    v_plan_id UUID;
    v_owner_user_id UUID;
    v_dr_emeka_id UUID;
    v_dr_sarah_id UUID;
    v_nurse_chioma_id UUID;
BEGIN
    -- Resolve User IDs
    SELECT id INTO v_owner_user_id FROM identity.users WHERE email = 'owner@curexal.com';
    SELECT id INTO v_dr_emeka_id FROM identity.users WHERE email = 'dr.emeka@curexal.com';
    SELECT id INTO v_dr_sarah_id FROM identity.users WHERE email = 'dr.sarah@curexal.com';
    SELECT id INTO v_nurse_chioma_id FROM identity.users WHERE email = 'nurse.chioma@curexal.com';

    -- 1. Ensure or retrieve organization 'curexal-clinic' with ENTERPRISE plan
    SELECT id INTO v_org_id FROM organization.organizations WHERE slug = 'curexal-clinic' LIMIT 1;
    IF v_org_id IS NULL THEN
        INSERT INTO organization.organizations (
            id, name, slug, status, plan, settings, created_at, updated_at
        ) VALUES (
            '00000000-0000-0000-0000-000000000001',
            'Curexal Clinic',
            'curexal-clinic',
            'active',
            'enterprise',
            '{"theme": {"primaryColor": "#0284c7", "fontFamily": "Plus Jakarta Sans"}, "features": ["clinical", "telehealth", "laboratory", "pharmacy", "radiology"]}'::jsonb,
            NOW(),
            NOW()
        ) RETURNING id INTO v_org_id;
    ELSE
        UPDATE organization.organizations SET
            name = 'Curexal Clinic',
            status = 'active',
            plan = 'enterprise',
            settings = '{"theme": {"primaryColor": "#0284c7", "fontFamily": "Plus Jakarta Sans"}, "features": ["clinical", "telehealth", "laboratory", "pharmacy", "radiology"]}'::jsonb,
            updated_at = NOW()
        WHERE id = v_org_id;
    END IF;

    -- Ensure active Enterprise subscription record
    SELECT id INTO v_plan_id FROM subscription.plans WHERE code = 'enterprise' LIMIT 1;
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema='subscription' AND table_name='subscriptions') THEN
        IF EXISTS (SELECT 1 FROM subscription.subscriptions WHERE organization_id = v_org_id) THEN
            UPDATE subscription.subscriptions SET
                plan = 'enterprise',
                plan_id = v_plan_id,
                status = 'active',
                updated_at = NOW()
            WHERE organization_id = v_org_id;
        ELSE
            INSERT INTO subscription.subscriptions (
                id, organization_id, plan_id, plan, status, starts_at, created_at, updated_at
            ) VALUES (
                gen_random_uuid(), v_org_id, v_plan_id, 'enterprise', 'active', NOW(), NOW(), NOW()
            );
        END IF;
    END IF;

    -- 2. Ensure facility type 'clinic' exists
    SELECT id INTO v_facility_type_id FROM platform.facility_types WHERE code IN ('clinic', 'hospital') LIMIT 1;
    IF v_facility_type_id IS NULL THEN
        INSERT INTO platform.facility_types (code, name, category, icon_key, description)
        VALUES ('clinic', 'Outpatient Clinic', 'clinical', 'Stethoscope', 'Primary healthcare and outpatient clinical practice')
        RETURNING id INTO v_facility_type_id;
    END IF;

    -- 3. Ensure or retrieve facility branch for Curexal Clinic
    SELECT id INTO v_branch_id FROM organization.facility_branches WHERE organization_id = v_org_id AND code = 'HO-01' LIMIT 1;
    IF v_branch_id IS NULL THEN
        SELECT id INTO v_branch_id FROM organization.facility_branches WHERE organization_id = v_org_id AND is_headquarters = TRUE LIMIT 1;
    END IF;

    IF v_branch_id IS NULL THEN
        INSERT INTO organization.facility_branches (
            id, organization_id, facility_type_id, code, slug, name, is_headquarters,
            email, phone, address, city, state, country, status, created_at, updated_at
        ) VALUES (
            '00000000-0000-0000-0000-000000000001',
            v_org_id,
            v_facility_type_id,
            'HO-01',
            'curexal-clinic',
            'Curexal Clinic Main Campus',
            TRUE,
            'support@curexal.com',
            '+234 800 CUREXAL',
            'Plot 104, Curexal Medical Boulevard, Victoria Island',
            'Lagos',
            'Lagos State',
            'Nigeria',
            'ACTIVE',
            NOW(),
            NOW()
        ) RETURNING id INTO v_branch_id;
    ELSE
        UPDATE organization.facility_branches SET
            name = 'Curexal Clinic Main Campus',
            code = 'HO-01',
            slug = 'curexal-clinic',
            email = 'support@curexal.com',
            phone = '+234 800 CUREXAL',
            address = 'Plot 104, Curexal Medical Boulevard, Victoria Island',
            city = 'Lagos',
            state = 'Lagos State',
            country = 'Nigeria',
            status = 'ACTIVE',
            updated_at = NOW()
        WHERE id = v_branch_id;
    END IF;

    -- Update slug if column exists on facility_branches
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='organization' AND table_name='facility_branches' AND column_name='slug') THEN
        UPDATE organization.facility_branches SET slug = 'curexal-clinic' WHERE id = v_branch_id;
    END IF;

    -- 4. Purge obsolete owner memberships (e.g. trusteegain@gmail.com)
    DELETE FROM organization.organization_memberships 
    WHERE organization_id = v_org_id 
      AND user_id IN (SELECT id FROM identity.users WHERE email IN ('trusteegain@gmail.com'));

    -- 5. Seed / Update Owner & Staff Memberships for Curexal Clinic
    INSERT INTO organization.organization_memberships (
        id, user_id, organization_id, role, role_title, is_active, joined_at, created_at
    ) VALUES 
        (gen_random_uuid(), v_owner_user_id, v_org_id, 'owner', 'Organization Owner', TRUE, NOW(), NOW()),
        (gen_random_uuid(), v_dr_emeka_id, v_org_id, 'doctor', 'Chief Medical Officer', TRUE, NOW(), NOW()),
        (gen_random_uuid(), v_dr_sarah_id, v_org_id, 'doctor', 'Consultant Specialist', TRUE, NOW(), NOW()),
        (gen_random_uuid(), v_nurse_chioma_id, v_org_id, 'nurse', 'Lead Triage Nurse', TRUE, NOW(), NOW())
    ON CONFLICT (organization_id, user_id) DO UPDATE SET
        role = EXCLUDED.role,
        role_title = EXCLUDED.role_title,
        is_active = TRUE;

    -- 6. Provider Profiles linked to v_branch_id
    INSERT INTO orchestration.provider_profiles (
        id, user_id, tenant_id, specialty_code, sub_specialties,
        telehealth_enabled, in_person_enabled, max_active_queue, current_active_queue,
        status, consultation_languages, rating, created_at, updated_at
    ) VALUES 
        (
            gen_random_uuid(),
            v_dr_emeka_id,
            v_branch_id,
            'GENERAL_CONSULTATION',
            ARRAY['Family Medicine', 'Telehealth General Practice', 'Emergency Care'],
            TRUE,
            TRUE,
            20,
            0,
            'ON_DUTY',
            ARRAY['English', 'Igbo', 'Pidgin'],
            4.95,
            NOW(),
            NOW()
        ),
        (
            gen_random_uuid(),
            v_dr_sarah_id,
            v_branch_id,
            'SPECIALIST',
            ARRAY['Internal Medicine', 'Cardiovascular Health', 'Endocrinology'],
            TRUE,
            TRUE,
            15,
            0,
            'ON_DUTY',
            ARRAY['English', 'Yoruba'],
            4.90,
            NOW(),
            NOW()
        )
    ON CONFLICT (tenant_id, user_id) DO UPDATE SET
        status = 'ON_DUTY',
        telehealth_enabled = TRUE,
        in_person_enabled = TRUE,
        updated_at = NOW();

END $$;
-- +goose StatementEnd

-- +goose Down
-- Reversible rollback
