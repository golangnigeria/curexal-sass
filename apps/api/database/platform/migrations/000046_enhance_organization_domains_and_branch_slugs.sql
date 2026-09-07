-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: ENHANCE ORGANIZATION DOMAINS & BRANCH SLUGS
-- ==============================================================================

-- 1. Ensure slug and uniqueness on organization.facility_branches
ALTER TABLE organization.facility_branches
ADD COLUMN IF NOT EXISTS slug VARCHAR(100);

UPDATE organization.facility_branches
SET slug = LOWER(code)
WHERE slug IS NULL OR slug = '';

ALTER TABLE organization.facility_branches
ALTER COLUMN slug SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uk_facility_branches_org_slug
ON organization.facility_branches(organization_id, slug);

-- 2. Upgrade existing organization.organization_domains for multi-domain routing & verification
-- The table was originally created in 000001 with (id, organization_id, domain_name, is_verified, created_at).
-- We rename domain_name -> hostname and add the new columns needed for org-centric routing.

-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'organization' AND table_name = 'organization_domains' AND column_name = 'domain_name'
    ) THEN
        ALTER TABLE organization.organization_domains RENAME COLUMN domain_name TO hostname;
    END IF;
END $$;
-- +goose StatementEnd

ALTER TABLE organization.organization_domains
ADD COLUMN IF NOT EXISTS domain_type VARCHAR(50) NOT NULL DEFAULT 'curexal_subdomain',
ADD COLUMN IF NOT EXISTS is_primary BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS verification_status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
ADD COLUMN IF NOT EXISTS verification_token VARCHAR(255),
ADD COLUMN IF NOT EXISTS ssl_status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
ADD COLUMN IF NOT EXISTS status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_domain_type') THEN
        ALTER TABLE organization.organization_domains
        ADD CONSTRAINT chk_domain_type CHECK (domain_type IN ('curexal_subdomain', 'custom'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_verification_status') THEN
        ALTER TABLE organization.organization_domains
        ADD CONSTRAINT chk_verification_status CHECK (verification_status IN ('PENDING', 'VERIFIED', 'FAILED'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_ssl_status') THEN
        ALTER TABLE organization.organization_domains
        ADD CONSTRAINT chk_ssl_status CHECK (ssl_status IN ('PENDING', 'ACTIVE', 'FAILED'));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_domain_status') THEN
        ALTER TABLE organization.organization_domains
        ADD CONSTRAINT chk_domain_status CHECK (status IN ('ACTIVE', 'INACTIVE', 'SUSPENDED'));
    END IF;
END $$;
-- +goose StatementEnd

CREATE INDEX IF NOT EXISTS idx_org_domains_org_id ON organization.organization_domains(organization_id);
CREATE INDEX IF NOT EXISTS idx_org_domains_hostname ON organization.organization_domains(hostname);
CREATE INDEX IF NOT EXISTS idx_org_domains_verified ON organization.organization_domains(is_verified) WHERE is_verified = TRUE;

INSERT INTO organization.organization_domains (organization_id, hostname, domain_type, is_primary, is_verified, verification_status, ssl_status, status)
SELECT id, slug || '.curexal.space', 'curexal_subdomain', TRUE, TRUE, 'VERIFIED', 'ACTIVE', 'ACTIVE'
FROM organization.organizations
WHERE slug IS NOT NULL AND slug != ''
ON CONFLICT (hostname) DO NOTHING;

-- +goose Down
DROP INDEX IF EXISTS organization.idx_org_domains_verified;
DROP INDEX IF EXISTS organization.idx_org_domains_hostname;
DROP INDEX IF EXISTS organization.idx_org_domains_org_id;

ALTER TABLE organization.organization_domains
DROP CONSTRAINT IF EXISTS chk_domain_status,
DROP CONSTRAINT IF EXISTS chk_ssl_status,
DROP CONSTRAINT IF EXISTS chk_verification_status,
DROP CONSTRAINT IF EXISTS chk_domain_type;

ALTER TABLE organization.organization_domains
DROP COLUMN IF EXISTS updated_at,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS ssl_status,
DROP COLUMN IF EXISTS verification_token,
DROP COLUMN IF EXISTS verification_status,
DROP COLUMN IF EXISTS is_primary,
DROP COLUMN IF EXISTS domain_type;

-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'organization' AND table_name = 'organization_domains' AND column_name = 'hostname'
    ) THEN
        ALTER TABLE organization.organization_domains RENAME COLUMN hostname TO domain_name;
    END IF;
END $$;
-- +goose StatementEnd

DROP INDEX IF EXISTS organization.uk_facility_branches_org_slug;
ALTER TABLE organization.facility_branches DROP COLUMN IF EXISTS slug;
