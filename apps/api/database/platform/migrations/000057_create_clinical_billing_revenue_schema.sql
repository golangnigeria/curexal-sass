-- +goose Up
-- Migration 000057: Complete Clinical Encounter Billing & Revenue Cycle Schemas

CREATE SCHEMA IF NOT EXISTS clinical_billing;

-- Encounter Invoices
CREATE TABLE IF NOT EXISTS clinical_billing.invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    invoice_number VARCHAR(64) NOT NULL,
    patient_id UUID NOT NULL,
    encounter_id UUID,
    status VARCHAR(30) NOT NULL DEFAULT 'ISSUED', -- DRAFT, ISSUED, PARTIALLY_PAID, PAID, VOIDED, REFUNDED
    subtotal NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    discount_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    tax_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    total_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    amount_paid NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    balance_due NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    hmo_claim_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    patient_copay NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    due_date TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_clinical_invoice_number UNIQUE (tenant_id, invoice_number)
);

CREATE INDEX IF NOT EXISTS idx_clinical_invoices_tenant_patient ON clinical_billing.invoices(tenant_id, patient_id);
CREATE INDEX IF NOT EXISTS idx_clinical_invoices_status ON clinical_billing.invoices(tenant_id, status);

-- Invoice Line Items (Consultations, Labs, Imaging, Pharmacy, Procedures)
CREATE TABLE IF NOT EXISTS clinical_billing.invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    invoice_id UUID NOT NULL REFERENCES clinical_billing.invoices(id) ON DELETE CASCADE,
    service_type VARCHAR(50) NOT NULL, -- CONSULTATION, LAB_TEST, RADIOLOGY, PHARMACY, PROCEDURE, ROOM_STAY
    service_code VARCHAR(64),
    description VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    total_price NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    discount_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_clinical_invoice_items ON clinical_billing.invoice_items(invoice_id);

-- Payment Transactions & Receipts
CREATE TABLE IF NOT EXISTS clinical_billing.payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    invoice_id UUID NOT NULL REFERENCES clinical_billing.invoices(id) ON DELETE CASCADE,
    receipt_number VARCHAR(64) NOT NULL,
    patient_id UUID NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL, -- CASH, POS_CARD, BANK_TRANSFER, HMO, WALLET
    payment_reference VARCHAR(128),
    cashier_id UUID NOT NULL,
    cashier_name VARCHAR(255) NOT NULL,
    paid_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_clinical_receipt_number UNIQUE (tenant_id, receipt_number)
);

CREATE INDEX IF NOT EXISTS idx_clinical_payments_invoice ON clinical_billing.payments(invoice_id);

-- +goose Down
DROP SCHEMA IF EXISTS clinical_billing CASCADE;
