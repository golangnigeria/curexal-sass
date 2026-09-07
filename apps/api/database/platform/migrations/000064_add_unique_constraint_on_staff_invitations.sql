-- +goose Up
-- ==============================================================================
-- CUREXAL PLATFORM — UNIQUE CONSTRAINT ON STAFF INVITATIONS
-- ==============================================================================
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'uk_staff_invitations_org_email'
    ) THEN
        ALTER TABLE organization.staff_invitations 
        ADD CONSTRAINT uk_staff_invitations_org_email UNIQUE (organization_id, email);
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- ==============================================================================
-- ROLLBACK UNIQUE CONSTRAINT ON STAFF INVITATIONS
-- ==============================================================================
-- +goose StatementBegin
DO $$
BEGIN
    ALTER TABLE organization.staff_invitations DROP CONSTRAINT IF EXISTS uk_staff_invitations_org_email;
END $$;
-- +goose StatementEnd
