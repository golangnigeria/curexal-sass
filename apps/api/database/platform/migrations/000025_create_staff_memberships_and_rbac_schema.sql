-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: STAFF MEMBERSHIPS & RBAC ENFORCEMENT SCHEMA
-- ==============================================================================

-- Multi-Branch Membership Assignments
CREATE TABLE IF NOT EXISTS organization.membership_branches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    membership_id UUID NOT NULL REFERENCES organization.organization_memberships(id) ON DELETE CASCADE,
    facility_branch_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    CONSTRAINT uk_membership_branch UNIQUE (membership_id, facility_branch_id)
);

CREATE INDEX IF NOT EXISTS idx_membership_branches_mem ON organization.membership_branches(membership_id);
CREATE INDEX IF NOT EXISTS idx_membership_branches_branch ON organization.membership_branches(facility_branch_id);

-- Departmental Memberships
CREATE TABLE IF NOT EXISTS organization.departmental_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    membership_id UUID NOT NULL REFERENCES organization.organization_memberships(id) ON DELETE CASCADE,
    facility_branch_id UUID NOT NULL REFERENCES organization.facility_branches(id) ON DELETE CASCADE,
    department_code VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    CONSTRAINT uk_dept_membership UNIQUE (membership_id, facility_branch_id, department_code)
);

CREATE INDEX IF NOT EXISTS idx_dept_memberships_mem ON organization.departmental_memberships(membership_id);
CREATE INDEX IF NOT EXISTS idx_dept_memberships_branch ON organization.departmental_memberships(facility_branch_id);
CREATE INDEX IF NOT EXISTS idx_dept_memberships_code ON organization.departmental_memberships(department_code);

-- Staff Invitations (Storing SHA-256 Hashed Tokens)
CREATE TABLE IF NOT EXISTS organization.staff_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    facility_branch_id UUID REFERENCES organization.facility_branches(id) ON DELETE SET NULL,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    role_title VARCHAR(100) NOT NULL DEFAULT 'member',
    invite_token_hash VARCHAR(64) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    invited_by UUID NOT NULL REFERENCES identity.users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_staff_invites_org ON organization.staff_invitations(organization_id);
CREATE INDEX IF NOT EXISTS idx_staff_invites_email ON organization.staff_invitations(email);
CREATE INDEX IF NOT EXISTS idx_staff_invites_pending_expiry ON organization.staff_invitations (expires_at) WHERE status = 'PENDING';

-- +goose Down
DROP TABLE IF EXISTS organization.staff_invitations CASCADE;
DROP TABLE IF EXISTS organization.departmental_memberships CASCADE;
DROP TABLE IF EXISTS organization.membership_branches CASCADE;
