-- +goose Up
-- ==============================================================================
-- CUREXAL DATABASE MIGRATION: MASTER HEALTHCARE REFERENCE CATALOGS SCHEMA
-- ==============================================================================

CREATE SCHEMA IF NOT EXISTS platform;

-- 1. Platform Clinical Reference Catalogs (ICD-10, Specialties)
CREATE TABLE IF NOT EXISTS platform.clinical_catalogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category VARCHAR(50) NOT NULL DEFAULT 'icd10',
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    system_group VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_clinical_catalogs_category_code ON platform.clinical_catalogs(category, code);

-- 2. Platform Laboratory Reference Catalogs (Tests, Specimens, Containers, UOM)
CREATE TABLE IF NOT EXISTS platform.lab_catalogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category VARCHAR(50) NOT NULL DEFAULT 'test',
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_lab_catalogs_category_code ON platform.lab_catalogs(category, code);

-- 3. Platform Radiology Reference Catalogs (DICOM Modalities, Procedures)
CREATE TABLE IF NOT EXISTS platform.radiology_catalogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category VARCHAR(50) NOT NULL DEFAULT 'modality',
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_radiology_catalogs_category_code ON platform.radiology_catalogs(category, code);

-- 4. Platform Pharmacy Reference Catalogs (Drug Categories, Dosages)
CREATE TABLE IF NOT EXISTS platform.pharmacy_catalogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category VARCHAR(50) NOT NULL DEFAULT 'drug_class',
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID REFERENCES identity.users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_pharmacy_catalogs_category_code ON platform.pharmacy_catalogs(category, code);

-- Seed Baseline Reference Records
INSERT INTO platform.clinical_catalogs (category, code, name, description, system_group)
VALUES
    ('icd10', 'E11.9', 'Type 2 diabetes mellitus without complications', 'Endocrine disorder', 'Endocrine'),
    ('icd10', 'I10', 'Essential (primary) hypertension', 'High blood pressure', 'Cardiovascular'),
    ('icd10', 'B54', 'Unspecified malaria', 'Plasmodium parasitic infection', 'Infectious'),
    ('icd10', 'D64.9', 'Anemia, unspecified', 'Low hemoglobin count', 'Hematology'),
    ('icd10', 'J06.9', 'Acute upper respiratory infection, unspecified', 'Common cold / URI', 'Respiratory'),
    ('specialty', 'GEN', 'General Practice / Internal Medicine', 'Primary care and internal medicine', 'General'),
    ('specialty', 'PATH', 'Pathology & Laboratory Medicine', 'Laboratory diagnostic medicine', 'Diagnostic'),
    ('specialty', 'RAD', 'Radiology & Diagnostic Imaging', 'Imaging diagnostic medicine', 'Diagnostic')
ON CONFLICT (code) DO NOTHING;

INSERT INTO platform.lab_catalogs (category, code, name, description, metadata)
VALUES
    ('specimen', 'EDTA', 'Whole Blood (EDTA)', 'Anticoagulated whole blood sample', '{"container": "Purple Cap EDTA Tube"}'::jsonb),
    ('specimen', 'SERUM', 'Serum (Clot Activator)', 'Clotted blood serum sample', '{"container": "Red/Yellow SST Tube"}'::jsonb),
    ('specimen', 'URINE', 'Urine (Midstream)', 'Midstream clean catch urine', '{"container": "Sterile Urine Cup"}'::jsonb),
    ('test_category', 'HEM', 'Hematology', 'Blood cell counts and coagulation testing', '{}'::jsonb),
    ('test_category', 'CHM', 'Clinical Chemistry', 'Metabolic panels and electrolytes', '{}'::jsonb),
    ('uom', 'MG_DL', 'Milligrams per Deciliter', 'mg/dL', '{"symbol": "mg/dL"}'::jsonb),
    ('uom', 'MMOL_L', 'Millimoles per Liter', 'mmol/L', '{"symbol": "mmol/L"}'::jsonb)
ON CONFLICT (code) DO NOTHING;

INSERT INTO platform.radiology_catalogs (category, code, name, description, metadata)
VALUES
    ('modality', 'DX', 'Digital Radiography (X-Ray)', 'Diagnostic Projection Radiography', '{"modalityType": "X-Ray"}'::jsonb),
    ('modality', 'US', 'Ultrasound', 'Diagnostic Sonography Imaging', '{"modalityType": "Ultrasound"}'::jsonb),
    ('modality', 'CT', 'Computed Tomography (CT Scan)', 'Cross-sectional X-Ray Tomography', '{"modalityType": "CT"}'::jsonb),
    ('modality', 'MR', 'Magnetic Resonance Imaging (MRI)', 'Magnetic Resonance Tomography', '{"modalityType": "MRI"}'::jsonb)
ON CONFLICT (code) DO NOTHING;

INSERT INTO platform.pharmacy_catalogs (category, code, name, description, metadata)
VALUES
    ('drug_class', 'ABX', 'Antibiotics & Antimicrobials', 'Antibacterial pharmaceutical agents', '{}'::jsonb),
    ('drug_class', 'ANALGESIC', 'Analgesics & Antipyretics', 'Pain and fever relief medications', '{}'::jsonb),
    ('dosage_form', 'TAB', 'Tablet', 'Oral compressed tablet', '{}'::jsonb),
    ('dosage_form', 'CAP', 'Capsule', 'Oral gelatin capsule', '{}'::jsonb),
    ('dosage_form', 'INJ', 'Injection', 'Parenteral injectable solution', '{}'::jsonb)
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS platform.pharmacy_catalogs;
DROP TABLE IF EXISTS platform.radiology_catalogs;
DROP TABLE IF EXISTS platform.lab_catalogs;
DROP TABLE IF EXISTS platform.clinical_catalogs;
