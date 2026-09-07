-- +goose Up
-- ==============================================================================
-- CUREXAL PLATFORM — SUPABASE STORAGE BUCKETS
-- ==============================================================================

-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = 'storage') THEN
        -- Insert curexal-documents bucket
        INSERT INTO storage.buckets (id, name, public, file_size_limit, allowed_mime_types)
        VALUES ('curexal-documents', 'curexal-documents', true, 52428800, null)
        ON CONFLICT (id) DO UPDATE SET public = true;

        -- Insert curexal-public bucket for logos, branding, assets
        INSERT INTO storage.buckets (id, name, public, file_size_limit, allowed_mime_types)
        VALUES ('curexal-public', 'curexal-public', true, 10485760, null)
        ON CONFLICT (id) DO UPDATE SET public = true;
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- ==============================================================================
-- ROLLBACK SUPABASE STORAGE BUCKETS
-- ==============================================================================
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = 'storage') THEN
        DELETE FROM storage.buckets WHERE id IN ('curexal-documents', 'curexal-public');
    END IF;
END $$;
-- +goose StatementEnd
