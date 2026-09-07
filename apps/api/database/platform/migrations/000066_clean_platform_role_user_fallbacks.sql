-- +goose Up
-- Migration: 000066_clean_platform_role_user_fallbacks.sql
-- Description: Clean legacy platform_role values ('user', 'member') in identity.users so that platform_role is strictly reserved for platform staff.

UPDATE identity.users
SET platform_role = NULL
WHERE platform_role IN ('user', 'member', '');

-- +goose Down
-- No-op reversal
