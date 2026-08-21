-- +goose Up
-- ==============================================================================
-- ENSURE SUBSCRIPTION.PLANS TABLE SCHEMAS
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS subscription;

CREATE TABLE IF NOT EXISTS subscription.plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    limits JSONB NOT NULL DEFAULT '{}'::jsonb,
    features JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
SELECT 1;
