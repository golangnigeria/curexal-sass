-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: ORGANIZATION LOCAL CATALOGS & PRICING SCHEMA
-- ==============================================================================

CREATE TABLE IF NOT EXISTS organization.catalog_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    master_catalog_id UUID,
    domain_type VARCHAR(50) NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    base_price NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    CONSTRAINT uk_org_catalog_code UNIQUE (organization_id, domain_type, code)
);

CREATE INDEX IF NOT EXISTS idx_org_catalogs_org_domain ON organization.catalog_items(organization_id, domain_type);
CREATE INDEX IF NOT EXISTS idx_org_catalogs_code ON organization.catalog_items(code);

-- Branch Price Overrides
CREATE TABLE IF NOT EXISTS organization.branch_price_overrides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    facility_branch_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    catalog_item_id UUID NOT NULL REFERENCES organization.catalog_items(id) ON DELETE CASCADE,
    override_price NUMERIC(15,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    CONSTRAINT uk_branch_catalog_price UNIQUE (facility_branch_id, catalog_item_id)
);

CREATE INDEX IF NOT EXISTS idx_branch_prices_branch ON organization.branch_price_overrides(facility_branch_id);

-- Insurance Providers & Copay Rules
CREATE TABLE IF NOT EXISTS organization.insurance_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL,
    coverage_percentage NUMERIC(5,2) NOT NULL DEFAULT 100.00,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    CONSTRAINT uk_org_insurance_code UNIQUE (organization_id, code)
);

CREATE INDEX IF NOT EXISTS idx_insurance_providers_org ON organization.insurance_providers(organization_id);

-- Register Catalog Permissions in "authorization".permissions
INSERT INTO "authorization".permissions (id, code, module, category, description)
VALUES 
    (gen_random_uuid(), 'organization:catalog:read', 'organization', 'core', 'Read organization catalog items and price lists'),
    (gen_random_uuid(), 'organization:catalog:write', 'organization', 'core', 'Create and update organization catalog items and price lists')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS organization.insurance_providers CASCADE;
DROP TABLE IF EXISTS organization.branch_price_overrides CASCADE;
DROP TABLE IF EXISTS organization.catalog_items CASCADE;
