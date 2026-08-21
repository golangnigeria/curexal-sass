-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: PLATFORM PRICING & PAYMENT GATEWAY VAULT
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS platform;

-- 1. Platform Pricing Rules Table
CREATE TABLE IF NOT EXISTS platform.pricing_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_type VARCHAR(50) NOT NULL, -- 'plan' or 'capability'
    target_code VARCHAR(100) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    monthly_price NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    annual_price NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    vat_percentage NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL,
    CONSTRAINT uk_pricing_target UNIQUE (target_type, target_code, currency)
);

-- 2. Platform Payment Gateways Encrypted Vault Table
CREATE TABLE IF NOT EXISTS platform.payment_gateways (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    priority INT NOT NULL DEFAULT 1,
    supported_currencies JSONB NOT NULL DEFAULT '["NGN", "USD"]'::jsonb,
    encrypted_secret_key TEXT NOT NULL,
    public_key VARCHAR(255),
    webhook_secret VARCHAR(255),
    retry_rules JSONB DEFAULT '{"maxRetries": 3, "backoffSeconds": 5}'::jsonb,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL
);

-- Seed Baseline Commercial Pricing Rules
INSERT INTO platform.pricing_rules (target_type, target_code, currency, monthly_price, annual_price, vat_percentage)
VALUES
    ('plan', 'smart', 'NGN', 50000.00, 500000.00, 7.50),
    ('plan', 'smart', 'USD', 100.00, 1000.00, 0.00),
    ('plan', 'optimize', 'NGN', 150000.00, 1500000.00, 7.50),
    ('plan', 'optimize', 'USD', 300.00, 3000.00, 0.00),
    ('plan', 'pro', 'NGN', 350000.00, 3500000.00, 7.50),
    ('plan', 'pro', 'USD', 700.00, 7000.00, 0.00),
    ('capability', 'laboratory.analyzer_integration', 'NGN', 25000.00, 250000.00, 7.50),
    ('capability', 'laboratory.analyzer_integration', 'USD', 50.00, 500.00, 0.00),
    ('capability', 'radiology.pacs_dicom', 'NGN', 40000.00, 400000.00, 7.50),
    ('capability', 'radiology.pacs_dicom', 'USD', 80.00, 800.00, 0.00)
ON CONFLICT (target_type, target_code, currency) DO NOTHING;

-- Seed Baseline Encrypted Payment Gateways
INSERT INTO platform.payment_gateways (provider_code, name, is_enabled, priority, supported_currencies, encrypted_secret_key, public_key)
VALUES
    ('paystack', 'Paystack Payments', TRUE, 1, '["NGN", "GHS", "USD"]'::jsonb, 'ae_mock_ciphertext_paystack', 'pk_live_paystack_public_placeholder'),
    ('flutterwave', 'Flutterwave Payments', TRUE, 2, '["NGN", "KES", "USD"]'::jsonb, 'ae_mock_ciphertext_flutterwave', 'FLWPUBK_LIVE-flutterwave-placeholder'),
    ('stripe', 'Stripe Payments', TRUE, 3, '["USD", "EUR", "GBP"]'::jsonb, 'ae_mock_ciphertext_stripe', 'pk_live_stripe_public_placeholder'),
    ('mock', 'Curexal Test Mock Provider', TRUE, 99, '["NGN", "USD", "GHS", "KES"]'::jsonb, 'ae_mock_ciphertext_mock', 'pk_test_mock_public')
ON CONFLICT (provider_code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS platform.payment_gateways;
DROP TABLE IF EXISTS platform.pricing_rules;
