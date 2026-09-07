-- +goose Up
-- ==============================================================================
-- CUREXAL PLATFORM — SUPABASE PG_CRON SCHEDULED JOBS
-- ==============================================================================

-- +goose StatementBegin
DO $$
BEGIN
    -- 1. Conditionally enable pg_cron extension if available on this Postgres host
    CREATE EXTENSION IF NOT EXISTS pg_cron WITH SCHEMA pg_catalog;
EXCEPTION WHEN OTHERS THEN
    RAISE NOTICE 'pg_cron extension not available or permission denied (will rely on application-level cron scheduler)';
END $$;
-- +goose StatementEnd

-- 2. Stored Procedure for Expired Token Cleanup
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION platform.cleanup_expired_tokens_and_sessions()
RETURNS void AS $$
BEGIN
    -- Delete expired identity verification tokens
    DELETE FROM identity.verification_tokens
    WHERE expires_at < NOW() - INTERVAL '24 hours';

    -- Delete expired password setup & reset tokens
    DELETE FROM identity.password_requests
    WHERE expires_at < NOW() - INTERVAL '24 hours';
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- 3. Stored Procedure for Outbox Pruning
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION platform.cleanup_completed_outbox_events()
RETURNS void AS $$
BEGIN
    DELETE FROM notification.outbox_events
    WHERE status = 'completed'
      AND updated_at < NOW() - INTERVAL '14 days';
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- 4. Register cron jobs with pg_cron if pg_cron is installed
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_cron') THEN
        -- Hourly cleanup of expired tokens
        PERFORM cron.unschedule('curexal-hourly-token-cleanup')
        WHERE EXISTS (SELECT 1 FROM cron.job WHERE jobname = 'curexal-hourly-token-cleanup');
        
        PERFORM cron.schedule(
            'curexal-hourly-token-cleanup',
            '0 * * * *',
            'SELECT platform.cleanup_expired_tokens_and_sessions();'
        );

        -- Daily outbox pruning at 03:00 UTC
        PERFORM cron.unschedule('curexal-daily-outbox-prune')
        WHERE EXISTS (SELECT 1 FROM cron.job WHERE jobname = 'curexal-daily-outbox-prune');

        PERFORM cron.schedule(
            'curexal-daily-outbox-prune',
            '0 3 * * *',
            'SELECT platform.cleanup_completed_outbox_events();'
        );
    END IF;
EXCEPTION WHEN OTHERS THEN
    RAISE NOTICE 'Skipping cron.schedule registration: %', SQLERRM;
END $$;
-- +goose StatementEnd

-- +goose Down
-- ==============================================================================
-- ROLLBACK SUPABASE PG_CRON SCHEDULED JOBS
-- ==============================================================================

-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_cron') THEN
        PERFORM cron.unschedule('curexal-hourly-token-cleanup')
        WHERE EXISTS (SELECT 1 FROM cron.job WHERE jobname = 'curexal-hourly-token-cleanup');

        PERFORM cron.unschedule('curexal-daily-outbox-prune')
        WHERE EXISTS (SELECT 1 FROM cron.job WHERE jobname = 'curexal-daily-outbox-prune');
    END IF;
EXCEPTION WHEN OTHERS THEN
    NULL;
END $$;
-- +goose StatementEnd

DROP FUNCTION IF EXISTS platform.cleanup_completed_outbox_events();
DROP FUNCTION IF EXISTS platform.cleanup_expired_tokens_and_sessions();
