-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: MEMBERSHIP & SESSION COMPATIBILITY VIEWS
-- ==============================================================================

CREATE OR REPLACE VIEW public.membership AS 
SELECT * FROM organization.organization_memberships;

CREATE OR REPLACE VIEW organization.memberships AS 
SELECT * FROM organization.organization_memberships;

CREATE OR REPLACE VIEW public.session AS 
SELECT * FROM identity.sessions;

-- +goose Down
DROP VIEW IF EXISTS public.session CASCADE;
DROP VIEW IF EXISTS organization.memberships CASCADE;
DROP VIEW IF EXISTS public.membership CASCADE;
