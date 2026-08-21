-- +goose Up
-- SQL tenant migration for Phase 7 Laboratory LIMS (BSD-007)

CREATE TABLE IF NOT EXISTS test_catalog (
    id TEXT PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    category_code TEXT NOT NULL,
    specimen_type_code TEXT NOT NULL,
    turnaround_time_hours INT NOT NULL DEFAULT 24,
    price NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    updated_by TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_test_catalog_code ON test_catalog(code);
CREATE INDEX IF NOT EXISTS idx_test_catalog_category ON test_catalog(category_code);

CREATE TABLE IF NOT EXISTS "order" (
    id TEXT PRIMARY KEY,
    order_number TEXT UNIQUE NOT NULL,
    patient_id TEXT NOT NULL REFERENCES patient(id) ON DELETE CASCADE,
    visit_id TEXT REFERENCES visit(id) ON DELETE SET NULL,
    ordering_provider_id TEXT REFERENCES provider(id) ON DELETE SET NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    priority TEXT NOT NULL DEFAULT 'routine',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    updated_by TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_order_number ON "order"(order_number);
CREATE INDEX IF NOT EXISTS idx_order_patient ON "order"(patient_id);
CREATE INDEX IF NOT EXISTS idx_order_status ON "order"(status);

CREATE TABLE IF NOT EXISTS specimen (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL REFERENCES "order"(id) ON DELETE CASCADE,
    barcode TEXT UNIQUE NOT NULL,
    specimen_type_code TEXT NOT NULL,
    collected_at TIMESTAMP WITH TIME ZONE,
    collected_by TEXT,
    accessioned_at TIMESTAMP WITH TIME ZONE,
    accessioned_by TEXT,
    status TEXT NOT NULL DEFAULT 'collected',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    updated_by TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_specimen_barcode ON specimen(barcode);
CREATE INDEX IF NOT EXISTS idx_specimen_order ON specimen(order_id);

CREATE TABLE IF NOT EXISTS result (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL REFERENCES "order"(id) ON DELETE CASCADE,
    test_catalog_id TEXT NOT NULL REFERENCES test_catalog(id) ON DELETE CASCADE,
    parameter_name TEXT NOT NULL,
    value TEXT NOT NULL,
    unit TEXT,
    reference_range TEXT,
    flag TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    updated_by TEXT
);

CREATE INDEX IF NOT EXISTS idx_result_order ON result(order_id);
CREATE INDEX IF NOT EXISTS idx_result_test_catalog ON result(test_catalog_id);

CREATE TABLE IF NOT EXISTS lab_authorization (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL REFERENCES "order"(id) ON DELETE CASCADE,
    authorized_by TEXT NOT NULL,
    authorized_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    stage TEXT NOT NULL DEFAULT 'final',
    signature_hash TEXT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    updated_by TEXT
);

CREATE INDEX IF NOT EXISTS idx_lab_authorization_order ON lab_authorization(order_id);

-- +goose Down
DROP TABLE IF EXISTS lab_authorization;
DROP TABLE IF EXISTS result;
DROP TABLE IF EXISTS specimen;
DROP TABLE IF EXISTS "order";
DROP TABLE IF EXISTS test_catalog;
