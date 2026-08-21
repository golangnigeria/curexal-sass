-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: ORGANIZATION BRANCHES & OPERATIONS SCHEMA
-- ==============================================================================

CREATE TABLE IF NOT EXISTS organization.facility_branches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    facility_type_id UUID NOT NULL REFERENCES platform.facility_types(id),
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_headquarters BOOLEAN NOT NULL DEFAULT FALSE,
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    lga VARCHAR(100),
    country VARCHAR(100) DEFAULT 'Nigeria',
    operating_hours JSONB NOT NULL DEFAULT '{}'::jsonb,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    CONSTRAINT uk_facility_org_code UNIQUE (organization_id, code)
);

-- Database-Level Single Headquarters Guarantee per Organization
CREATE UNIQUE INDEX IF NOT EXISTS uk_facility_single_headquarters
ON organization.facility_branches (organization_id)
WHERE is_headquarters = TRUE;

CREATE INDEX IF NOT EXISTS idx_facility_branches_org_id ON organization.facility_branches(organization_id);
CREATE INDEX IF NOT EXISTS idx_facility_branches_type ON organization.facility_branches(facility_type_id);
CREATE INDEX IF NOT EXISTS idx_facility_branches_status ON organization.facility_branches(status);

-- Register Facility Branch Permissions in "authorization".permissions
INSERT INTO "authorization".permissions (id, code, module, category, description)
VALUES 
    (gen_random_uuid(), 'organization:branch:create', 'organization', 'core', 'Create operational facility branches'),
    (gen_random_uuid(), 'organization:branch:read', 'organization', 'core', 'Read operational facility branches'),
    (gen_random_uuid(), 'organization:branch:update', 'organization', 'core', 'Update operational facility branches'),
    (gen_random_uuid(), 'organization:branch:deactivate', 'organization', 'core', 'Deactivate operational facility branches')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS organization.facility_branches CASCADE;
