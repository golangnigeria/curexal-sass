-- +goose Up
-- Migration: 000037_identity_verification_tokens_metadata.sql
-- Description: Add metadata JSONB column to identity.verification_tokens for flexible token payloads.

ALTER TABLE identity.verification_tokens ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}'::jsonb;

-- +goose Down
ALTER TABLE identity.verification_tokens DROP COLUMN IF EXISTS metadata;
