-- +goose Up
-- Migration: 000065_add_updated_at_to_memberships.sql
-- Description: Add updated_at column to organization.organization_memberships and workspace.workspace_memberships

ALTER TABLE organization.organization_memberships 
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE workspace.workspace_memberships 
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- +goose Down
ALTER TABLE organization.organization_memberships DROP COLUMN IF EXISTS updated_at;
ALTER TABLE workspace.workspace_memberships DROP COLUMN IF EXISTS updated_at;
