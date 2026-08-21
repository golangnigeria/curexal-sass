-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: PLATFORM CONFIGURATION KERNEL SCHEMA
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS platform;

-- 1. Platform General Settings Table
CREATE TABLE IF NOT EXISTS platform.general_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    platform_name VARCHAR(255) NOT NULL DEFAULT 'Curexal',
    platform_logo_url TEXT,
    platform_favicon_url TEXT,
    platform_description TEXT DEFAULT 'Enterprise Healthcare Operating Platform',
    support_email VARCHAR(255) NOT NULL DEFAULT 'support@curexal.space',
    support_phone VARCHAR(50) NOT NULL DEFAULT '+234800CUREXAL',
    support_website VARCHAR(255) DEFAULT 'https://curexal.space',
    default_country VARCHAR(10) NOT NULL DEFAULT 'NG',
    default_currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
    supported_countries JSONB NOT NULL DEFAULT '["NG", "GH", "KE", "ZA"]'::jsonb,
    supported_currencies JSONB NOT NULL DEFAULT '["NGN", "USD", "GHS", "KES"]'::jsonb,
    default_timezone VARCHAR(50) NOT NULL DEFAULT 'Africa/Lagos',
    default_locale VARCHAR(10) NOT NULL DEFAULT 'en',
    date_format VARCHAR(30) NOT NULL DEFAULT 'YYYY-MM-DD',
    time_format VARCHAR(30) NOT NULL DEFAULT 'HH:mm',
    number_format VARCHAR(30) NOT NULL DEFAULT 'standard',
    measurement_units VARCHAR(20) NOT NULL DEFAULT 'metric',
    maintenance_mode BOOLEAN NOT NULL DEFAULT FALSE,
    announcement_banner TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL
);

-- 2. Platform Identity & Security Policies Table
CREATE TABLE IF NOT EXISTS platform.identity_security_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    min_password_length INT NOT NULL DEFAULT 8,
    password_require_uppercase BOOLEAN NOT NULL DEFAULT TRUE,
    password_require_number BOOLEAN NOT NULL DEFAULT TRUE,
    password_require_symbol BOOLEAN NOT NULL DEFAULT TRUE,
    password_expiration_days INT NOT NULL DEFAULT 90,
    email_verification_required BOOLEAN NOT NULL DEFAULT TRUE,
    mfa_policy VARCHAR(50) NOT NULL DEFAULT 'OPTIONAL',
    max_failed_login_attempts INT NOT NULL DEFAULT 5,
    account_lockout_duration_minutes INT NOT NULL DEFAULT 30,
    session_max_duration_hours INT NOT NULL DEFAULT 12,
    refresh_token_duration_days INT NOT NULL DEFAULT 30,
    max_active_sessions INT NOT NULL DEFAULT 5,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL
);

-- Seed Default Singleton Rows
INSERT INTO platform.general_settings (id, platform_name, default_country, default_currency, default_timezone)
SELECT '00000000-0000-0000-0000-000000000001'::uuid, 'Curexal', 'NG', 'NGN', 'Africa/Lagos'
WHERE NOT EXISTS (SELECT 1 FROM platform.general_settings);

INSERT INTO platform.identity_security_policies (id, min_password_length, max_failed_login_attempts)
SELECT '00000000-0000-0000-0000-000000000002'::uuid, 8, 5
WHERE NOT EXISTS (SELECT 1 FROM platform.identity_security_policies);

-- +goose Down
DROP TABLE IF EXISTS platform.identity_security_policies;
DROP TABLE IF EXISTS platform.general_settings;
