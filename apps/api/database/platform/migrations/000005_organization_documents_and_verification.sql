-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: ORGANIZATION DOCUMENTS & VERIFICATION SCHEMA
-- ==============================================================================

CREATE TABLE IF NOT EXISTS organization.organization_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    document_type VARCHAR(100) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    storage_key VARCHAR(500) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    checksum_sha256 VARCHAR(64) NOT NULL,
    uploaded_by UUID NOT NULL REFERENCES identity.users(id) ON DELETE RESTRICT,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    version INT NOT NULL DEFAULT 1,
    reviewed_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMP WITH TIME ZONE,
    rejection_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_org_docs_org_id ON organization.organization_documents(organization_id);
CREATE INDEX IF NOT EXISTS idx_org_docs_status ON organization.organization_documents(status);
CREATE INDEX IF NOT EXISTS idx_org_docs_doc_type ON organization.organization_documents(document_type);

-- Ensure identity.users has credential tracking columns if missing
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS credential_status VARCHAR(50) DEFAULT 'ACTIVE';
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS failed_login_attempts INT DEFAULT 0;
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS locked_until TIMESTAMP WITH TIME ZONE;
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS password_changed_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS password_expires_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS last_successful_login_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE identity.users ADD COLUMN IF NOT EXISTS last_failed_login_at TIMESTAMP WITH TIME ZONE;

-- Ensure authorization.permissions has module and category columns if missing
ALTER TABLE "authorization".permissions ADD COLUMN IF NOT EXISTS module VARCHAR(100) DEFAULT 'organization';
ALTER TABLE "authorization".permissions ADD COLUMN IF NOT EXISTS category VARCHAR(50) DEFAULT 'core';

-- Register Document & Verification Permissions in "authorization".permissions
INSERT INTO "authorization".permissions (id, code, module, category, description)
VALUES 
    (gen_random_uuid(), 'organization:document:upload', 'organization', 'core', 'Upload regulatory organization documents'),
    (gen_random_uuid(), 'organization:document:read', 'organization', 'core', 'Read organization documents and metadata'),
    (gen_random_uuid(), 'organization:document:review', 'organization', 'core', 'Review regulatory organization documents'),
    (gen_random_uuid(), 'organization:document:approve', 'organization', 'core', 'Approve regulatory organization documents'),
    (gen_random_uuid(), 'organization:document:reject', 'organization', 'core', 'Reject regulatory organization documents'),
    (gen_random_uuid(), 'organization:verify', 'organization', 'core', 'Approve healthcare organization activation status')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS organization.organization_documents CASCADE;
