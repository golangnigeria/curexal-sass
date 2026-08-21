-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: ENHANCE ORGANIZATION PROFILE & SETUP WIZARD SCHEMA
-- ==============================================================================

ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS registration_number VARCHAR(100);
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS license_number VARCHAR(100);
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS tax_id VARCHAR(100);
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS email VARCHAR(255);
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS phone VARCHAR(50);
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS address TEXT;
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS city VARCHAR(100);
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS state VARCHAR(100);
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS lga VARCHAR(100);
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS country VARCHAR(100) DEFAULT 'Nigeria';
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS setup_state VARCHAR(50) NOT NULL DEFAULT 'PENDING_REGISTRATION';
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS setup_step INT NOT NULL DEFAULT 1;
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS version INT NOT NULL DEFAULT 1;
ALTER TABLE organization.organizations ADD COLUMN IF NOT EXISTS updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_organizations_setup_state ON organization.organizations(setup_state);
CREATE INDEX IF NOT EXISTS idx_organizations_tax_id ON organization.organizations(tax_id);

-- +goose Down
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS updated_by;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS version;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS completed_at;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS setup_step;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS setup_state;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS country;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS lga;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS state;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS city;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS address;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS phone;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS email;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS tax_id;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS license_number;
ALTER TABLE organization.organizations DROP COLUMN IF EXISTS registration_number;
