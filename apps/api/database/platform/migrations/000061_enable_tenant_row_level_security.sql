-- +goose Up
-- ==============================================================================
-- CUREXAL PLATFORM — SUPABASE ROW-LEVEL SECURITY (RLS) POLICIES
-- ==============================================================================

-- 0. Conditionally create auth schema fallback functions ONLY if not running on Supabase Cloud
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_proc JOIN pg_namespace ON pg_proc.pronamespace = pg_namespace.oid WHERE pg_namespace.nspname = 'auth' AND pg_proc.proname = 'uid') THEN
        CREATE SCHEMA IF NOT EXISTS auth;
        EXECUTE 'CREATE OR REPLACE FUNCTION auth.uid() RETURNS uuid AS $func$ BEGIN RETURN NULLIF(COALESCE(current_setting(''request.jwt.claim.sub'', true), (current_setting(''request.jwt.claims'', true)::jsonb ->> ''sub'')), '''')::uuid; EXCEPTION WHEN OTHERS THEN RETURN NULL; END; $func$ LANGUAGE plpgsql STABLE;';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_proc JOIN pg_namespace ON pg_proc.pronamespace = pg_namespace.oid WHERE pg_namespace.nspname = 'auth' AND pg_proc.proname = 'role') THEN
        CREATE SCHEMA IF NOT EXISTS auth;
        EXECUTE 'CREATE OR REPLACE FUNCTION auth.role() RETURNS text AS $func$ BEGIN RETURN COALESCE(current_setting(''request.jwt.claim.role'', true), (current_setting(''request.jwt.claims'', true)::jsonb ->> ''role''), ''anon''); EXCEPTION WHEN OTHERS THEN RETURN ''anon''; END; $func$ LANGUAGE plpgsql STABLE;';
    END IF;
END $$;
-- +goose StatementEnd

-- Helper Function to Extract Current Authenticated Organization ID from JWT Claims
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION platform.current_organization_id()
RETURNS uuid AS $$
BEGIN
    RETURN NULLIF(
        COALESCE(
            current_setting('request.jwt.claim.organization_id', true),
            (current_setting('request.jwt.claims', true)::jsonb -> 'user_metadata' ->> 'organization_id'),
            (current_setting('request.jwt.claims', true)::jsonb ->> 'organization_id')
        ),
        ''
    )::uuid;
EXCEPTION WHEN OTHERS THEN
    RETURN NULL;
END;
$$ LANGUAGE plpgsql STABLE SECURITY DEFINER;
-- +goose StatementEnd

-- Helper Function to Extract Current Authenticated User ID
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION platform.current_user_id()
RETURNS uuid AS $$
BEGIN
    RETURN NULLIF(
        COALESCE(
            current_setting('request.jwt.claim.sub', true),
            (current_setting('request.jwt.claims', true)::jsonb ->> 'sub')
        ),
        ''
    )::uuid;
EXCEPTION WHEN OTHERS THEN
    RETURN NULL;
END;
$$ LANGUAGE plpgsql STABLE SECURITY DEFINER;
-- +goose StatementEnd

-- 1. Patient Profiles RLS
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'patient' AND table_name = 'patients') THEN
        ALTER TABLE patient.patients ENABLE ROW LEVEL SECURITY;

        DROP POLICY IF EXISTS tenant_isolation_patient_select ON patient.patients;
        CREATE POLICY tenant_isolation_patient_select ON patient.patients
            FOR SELECT
            USING (
                organization_id = platform.current_organization_id()
                OR auth.uid()::text = user_id::text
                OR auth.role() = 'service_role'
            );

        DROP POLICY IF EXISTS tenant_isolation_patient_insert ON patient.patients;
        CREATE POLICY tenant_isolation_patient_insert ON patient.patients
            FOR INSERT
            WITH CHECK (
                organization_id = platform.current_organization_id()
                OR auth.role() = 'service_role'
            );

        DROP POLICY IF EXISTS tenant_isolation_patient_update ON patient.patients;
        CREATE POLICY tenant_isolation_patient_update ON patient.patients
            FOR UPDATE
            USING (
                organization_id = platform.current_organization_id()
                OR auth.role() = 'service_role'
            );
    END IF;
END $$;
-- +goose StatementEnd

-- 2. Organization Documents RLS
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'organization' AND table_name = 'organization_documents') THEN
        ALTER TABLE organization.organization_documents ENABLE ROW LEVEL SECURITY;

        DROP POLICY IF EXISTS tenant_isolation_documents_policy ON organization.organization_documents;
        CREATE POLICY tenant_isolation_documents_policy ON organization.organization_documents
            FOR ALL
            USING (
                organization_id = platform.current_organization_id()
                OR auth.role() = 'service_role'
            )
            WITH CHECK (
                organization_id = platform.current_organization_id()
                OR auth.role() = 'service_role'
            );
    END IF;
END $$;
-- +goose StatementEnd

-- 3. Billing Invoices RLS
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'billing' AND table_name = 'invoices') THEN
        ALTER TABLE billing.invoices ENABLE ROW LEVEL SECURITY;

        DROP POLICY IF EXISTS tenant_isolation_invoices_policy ON billing.invoices;
        CREATE POLICY tenant_isolation_invoices_policy ON billing.invoices
            FOR ALL
            USING (
                organization_id = platform.current_organization_id()
                OR auth.role() = 'service_role'
            )
            WITH CHECK (
                organization_id = platform.current_organization_id()
                OR auth.role() = 'service_role'
            );
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- ==============================================================================
-- ROLLBACK SUPABASE ROW-LEVEL SECURITY (RLS) POLICIES
-- ==============================================================================

-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'patient' AND table_name = 'patients') THEN
        ALTER TABLE patient.patients DISABLE ROW LEVEL SECURITY;
        DROP POLICY IF EXISTS tenant_isolation_patient_select ON patient.patients;
        DROP POLICY IF EXISTS tenant_isolation_patient_insert ON patient.patients;
        DROP POLICY IF EXISTS tenant_isolation_patient_update ON patient.patients;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'organization' AND table_name = 'organization_documents') THEN
        ALTER TABLE organization.organization_documents DISABLE ROW LEVEL SECURITY;
        DROP POLICY IF EXISTS tenant_isolation_documents_policy ON organization.organization_documents;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'billing' AND table_name = 'invoices') THEN
        ALTER TABLE billing.invoices DISABLE ROW LEVEL SECURITY;
        DROP POLICY IF EXISTS tenant_isolation_invoices_policy ON billing.invoices;
    END IF;
END $$;
-- +goose StatementEnd

DROP FUNCTION IF EXISTS platform.current_organization_id();
DROP FUNCTION IF EXISTS platform.current_user_id();
