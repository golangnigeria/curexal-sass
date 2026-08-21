-- +goose Up
-- ==============================================================================
-- MIGRATION 000003: Create organization.demo_requests table for platform demo leads
-- ==============================================================================

CREATE TABLE IF NOT EXISTS organization.demo_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    organization_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(50),
    message TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS organization.demo_requests CASCADE;

