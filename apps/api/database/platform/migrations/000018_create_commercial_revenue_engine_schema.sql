-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: COMMERCIAL REVENUE ENGINE (MILESTONE 4)
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS subscription;

-- 1. Commercial Orders Table
CREATE TABLE IF NOT EXISTS subscription.commercial_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(100) NOT NULL UNIQUE,
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_payment', -- pending_payment, paid, cancelled, expired, refunded
    currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    tax DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    discount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    total DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    created_by UUID,
    paid_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    refunded_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_order_status CHECK (status IN ('pending_payment', 'paid', 'cancelled', 'expired', 'refunded'))
);

CREATE INDEX IF NOT EXISTS idx_commercial_orders_org ON subscription.commercial_orders(organization_id);
CREATE INDEX IF NOT EXISTS idx_commercial_orders_status ON subscription.commercial_orders(status);

-- 2. Commercial Order Items Table
CREATE TABLE IF NOT EXISTS subscription.commercial_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES subscription.commercial_orders(id) ON DELETE CASCADE,
    capability_id UUID NOT NULL REFERENCES subscription.capabilities(id) ON DELETE CASCADE,
    capability_code VARCHAR(100) NOT NULL,
    billing_cycle VARCHAR(20) NOT NULL DEFAULT 'monthly', -- monthly, annual
    unit_price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    amount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_order_item_cycle CHECK (billing_cycle IN ('monthly', 'annual'))
);

CREATE INDEX IF NOT EXISTS idx_commercial_order_items_order ON subscription.commercial_order_items(order_id);

-- 3. Payment Transactions Table
CREATE TABLE IF NOT EXISTS subscription.payment_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES subscription.commercial_orders(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- paystack, flutterwave, stripe, mock
    provider_transaction_id VARCHAR(255) NOT NULL UNIQUE,
    provider_reference VARCHAR(255) NOT NULL,
    payment_url TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, successful, failed, refunded
    amount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
    metadata JSONB DEFAULT '{}'::jsonb,
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_payment_status CHECK (status IN ('pending', 'successful', 'failed', 'refunded'))
);

CREATE INDEX IF NOT EXISTS idx_payment_tx_order ON subscription.payment_transactions(order_id);
CREATE INDEX IF NOT EXISTS idx_payment_tx_org ON subscription.payment_transactions(organization_id);
CREATE INDEX IF NOT EXISTS idx_payment_tx_provider_tx ON subscription.payment_transactions(provider, provider_transaction_id);

-- 4. Cryptographic Webhook Events Log (Enforces Provider/Event Idempotency)
CREATE TABLE IF NOT EXISTS subscription.webhook_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider VARCHAR(50) NOT NULL,
    event_id VARCHAR(255) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_webhook_provider_event UNIQUE (provider, event_id)
);

CREATE INDEX IF NOT EXISTS idx_webhook_events_provider_event ON subscription.webhook_events(provider, event_id);

-- 5. Commercial Invoices Table
CREATE TABLE IF NOT EXISTS subscription.commercial_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(100) NOT NULL UNIQUE,
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES subscription.commercial_orders(id) ON DELETE CASCADE,
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    tax DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    total DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
    status VARCHAR(50) NOT NULL DEFAULT 'issued', -- issued, paid, cancelled, refunded
    issued_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    due_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '14 days',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_commercial_invoices_org ON subscription.commercial_invoices(organization_id);
CREATE INDEX IF NOT EXISTS idx_commercial_invoices_order ON subscription.commercial_invoices(order_id);

-- 6. Commercial Receipts Table
CREATE TABLE IF NOT EXISTS subscription.commercial_receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receipt_number VARCHAR(100) NOT NULL UNIQUE,
    payment_id UUID NOT NULL REFERENCES subscription.payment_transactions(id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES subscription.commercial_orders(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organization.organizations(id) ON DELETE CASCADE,
    amount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
    provider_reference VARCHAR(255) NOT NULL,
    paid_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_commercial_receipts_org ON subscription.commercial_receipts(organization_id);
CREATE INDEX IF NOT EXISTS idx_commercial_receipts_payment ON subscription.commercial_receipts(payment_id);

-- +goose Down
DROP TABLE IF EXISTS subscription.commercial_receipts;
DROP TABLE IF EXISTS subscription.commercial_invoices;
DROP TABLE IF EXISTS subscription.webhook_events;
DROP TABLE IF EXISTS subscription.payment_transactions;
DROP TABLE IF EXISTS subscription.commercial_order_items;
DROP TABLE IF EXISTS subscription.commercial_orders;
