-- +goose Up
-- ==============================================================================
-- MIGRATION 000040: Align organization.demo_requests schema with platform model
-- ==============================================================================

-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'organization' 
        AND table_name = 'demo_requests' 
        AND column_name = 'organization_name'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'organization' 
        AND table_name = 'demo_requests' 
        AND column_name = 'laboratory_name'
    ) THEN
        ALTER TABLE organization.demo_requests RENAME COLUMN organization_name TO laboratory_name;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'organization' 
        AND table_name = 'demo_requests' 
        AND column_name = 'full_name'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'organization' 
        AND table_name = 'demo_requests' 
        AND column_name = 'contact_name'
    ) THEN
        ALTER TABLE organization.demo_requests RENAME COLUMN full_name TO contact_name;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'organization' 
        AND table_name = 'demo_requests' 
        AND column_name = 'phone_number'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'organization' 
        AND table_name = 'demo_requests' 
        AND column_name = 'phone'
    ) THEN
        ALTER TABLE organization.demo_requests RENAME COLUMN phone_number TO phone;
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'organization' 
        AND table_name = 'demo_requests' 
        AND column_name = 'message'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'organization' 
        AND table_name = 'demo_requests' 
        AND column_name = 'notes'
    ) THEN
        ALTER TABLE organization.demo_requests RENAME COLUMN message TO notes;
    END IF;
END $$;
-- +goose StatementEnd

ALTER TABLE organization.demo_requests ADD COLUMN IF NOT EXISTS laboratory_name VARCHAR(255);
ALTER TABLE organization.demo_requests ADD COLUMN IF NOT EXISTS contact_name VARCHAR(255);
ALTER TABLE organization.demo_requests ADD COLUMN IF NOT EXISTS phone VARCHAR(50);
ALTER TABLE organization.demo_requests ADD COLUMN IF NOT EXISTS daily_volume VARCHAR(100);
ALTER TABLE organization.demo_requests ADD COLUMN IF NOT EXISTS notes TEXT;
ALTER TABLE organization.demo_requests ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_demo_requests_status ON organization.demo_requests(status);
CREATE INDEX IF NOT EXISTS idx_demo_requests_created_at ON organization.demo_requests(created_at DESC);

-- +goose Down
-- ==============================================================================
DROP INDEX IF EXISTS organization.idx_demo_requests_created_at;
DROP INDEX IF EXISTS organization.idx_demo_requests_status;
ALTER TABLE organization.demo_requests DROP COLUMN IF EXISTS updated_at;
ALTER TABLE organization.demo_requests DROP COLUMN IF EXISTS daily_volume;
