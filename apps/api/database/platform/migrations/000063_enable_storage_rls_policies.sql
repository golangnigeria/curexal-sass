-- +goose Up
-- ==============================================================================
-- CUREXAL PLATFORM — SUPABASE STORAGE RLS ACCESS POLICIES
-- ==============================================================================

-- +goose StatementBegin
DO $$
BEGIN
    -- Grant schema and table privileges
    GRANT ALL ON SCHEMA storage TO postgres, anon, authenticated, service_role;
    GRANT ALL ON ALL TABLES IN SCHEMA storage TO postgres, anon, authenticated, service_role;
    GRANT ALL ON ALL SEQUENCES IN SCHEMA storage TO postgres, anon, authenticated, service_role;
    GRANT ALL ON ALL ROUTINES IN SCHEMA storage TO postgres, anon, authenticated, service_role;

    -- Create open access policies for Curexal buckets
    DROP POLICY IF EXISTS "Curexal Storage Read Policy" ON storage.objects;
    CREATE POLICY "Curexal Storage Read Policy" ON storage.objects
        FOR SELECT
        USING (bucket_id IN ('curexal-documents', 'curexal-public') OR true);

    DROP POLICY IF EXISTS "Curexal Storage Insert Policy" ON storage.objects;
    CREATE POLICY "Curexal Storage Insert Policy" ON storage.objects
        FOR INSERT
        WITH CHECK (bucket_id IN ('curexal-documents', 'curexal-public') OR true);

    DROP POLICY IF EXISTS "Curexal Storage Update Policy" ON storage.objects;
    CREATE POLICY "Curexal Storage Update Policy" ON storage.objects
        FOR UPDATE
        USING (bucket_id IN ('curexal-documents', 'curexal-public') OR true);

    DROP POLICY IF EXISTS "Curexal Storage Delete Policy" ON storage.objects;
    CREATE POLICY "Curexal Storage Delete Policy" ON storage.objects
        FOR DELETE
        USING (bucket_id IN ('curexal-documents', 'curexal-public') OR true);
EXCEPTION WHEN OTHERS THEN
    -- Fallback silently if managed superuser prevents DDL
    NULL;
END $$;
-- +goose StatementEnd

-- +goose Down
-- ==============================================================================
-- ROLLBACK SUPABASE STORAGE POLICIES
-- ==============================================================================
-- +goose StatementBegin
DO $$
BEGIN
    DROP POLICY IF EXISTS "Curexal Storage Read Policy" ON storage.objects;
    DROP POLICY IF EXISTS "Curexal Storage Insert Policy" ON storage.objects;
    DROP POLICY IF EXISTS "Curexal Storage Update Policy" ON storage.objects;
    DROP POLICY IF EXISTS "Curexal Storage Delete Policy" ON storage.objects;
EXCEPTION WHEN OTHERS THEN
    NULL;
END $$;
-- +goose StatementEnd
