-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: COMMERCIAL CAPABILITY PRICING & SUBSCRIPTIONS
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS subscription;

-- 1. Multi-Currency Capability Pricing Catalog (DB is Law - Prices are NOT hardcoded in Go)
CREATE TABLE IF NOT EXISTS subscription.capability_prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    capability_id UUID NOT NULL REFERENCES subscription.capabilities(id) ON DELETE CASCADE,
    currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
    billing_period VARCHAR(20) NOT NULL DEFAULT 'monthly', -- monthly, annual
    price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_capability_price UNIQUE (capability_id, currency, billing_period)
);

-- 2. Commercial Capability Subscriptions (Purchases)
CREATE TABLE IF NOT EXISTS subscription.capability_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    capability_id UUID NOT NULL REFERENCES subscription.capabilities(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_payment', -- pending_payment, active, past_due, cancelled
    billing_cycle VARCHAR(20) NOT NULL DEFAULT 'monthly',
    price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    current_period_start TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    current_period_end TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '30 days',
    cancelled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 3. Extend Entitlements with Source Reference Tracing
ALTER TABLE subscription.organization_entitlements
ADD COLUMN IF NOT EXISTS source_type VARCHAR(50) DEFAULT 'add_on',
ADD COLUMN IF NOT EXISTS source_id UUID;

-- 4. Seed Multi-Currency Pricing for Capabilities (African Market NGN & Global USD)
INSERT INTO subscription.capability_prices (capability_id, currency, billing_period, price)
SELECT c.id, p.currency, p.billing_period, p.price
FROM subscription.capabilities c
CROSS JOIN (
    VALUES
        -- Basic LIS
        ('laboratory.basic', 'NGN', 'monthly', 35000.00),
        ('laboratory.basic', 'USD', 'monthly', 30.00),
        -- LIS Analyzer Integration
        ('laboratory.analyzer_integration', 'NGN', 'monthly', 25000.00),
        ('laboratory.analyzer_integration', 'USD', 'monthly', 20.00),
        -- Advanced QC
        ('laboratory.advanced_qc', 'NGN', 'monthly', 15000.00),
        ('laboratory.advanced_qc', 'USD', 'monthly', 12.00),
        -- Basic RIS
        ('radiology.basic', 'NGN', 'monthly', 30000.00),
        ('radiology.basic', 'USD', 'monthly', 25.00),
        -- PACS DICOM Integration
        ('radiology.pacs_dicom', 'NGN', 'monthly', 45000.00),
        ('radiology.pacs_dicom', 'USD', 'monthly', 40.00),
        -- Clinical Inpatient Wards
        ('clinical.inpatient_wards', 'NGN', 'monthly', 40000.00),
        ('clinical.inpatient_wards', 'USD', 'monthly', 35.00),
        -- Pharmacy Basic & Advanced Inventory
        ('pharmacy.basic', 'NGN', 'monthly', 25000.00),
        ('pharmacy.basic', 'USD', 'monthly', 20.00),
        ('pharmacy.advanced_inventory', 'NGN', 'monthly', 20000.00),
        ('pharmacy.advanced_inventory', 'USD', 'monthly', 18.00),
        -- Inventory Basic & Multi-Warehouse
        ('inventory.basic', 'NGN', 'monthly', 20000.00),
        ('inventory.basic', 'USD', 'monthly', 15.00),
        ('inventory.multi_warehouse', 'NGN', 'monthly', 25000.00),
        ('inventory.multi_warehouse', 'USD', 'monthly', 20.00),
        -- Core Customer Care & Billing
        ('core.customer_care', 'NGN', 'monthly', 15000.00),
        ('core.customer_care', 'USD', 'monthly', 12.00),
        ('core.billing', 'NGN', 'monthly', 30000.00),
        ('core.billing', 'USD', 'monthly', 25.00),
        -- ISO 15189 QMS
        ('qms.iso15189', 'NGN', 'monthly', 20000.00),
        ('qms.iso15189', 'USD', 'monthly', 18.00)
) AS p(code, currency, billing_period, price)
WHERE c.code = p.code
ON CONFLICT (capability_id, currency, billing_period) DO NOTHING;

-- +goose Down
SELECT 1;
