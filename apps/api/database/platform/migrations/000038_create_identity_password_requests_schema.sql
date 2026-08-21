-- +goose Up
-- ==============================================================================
-- CUREXAL IDENTITY: PASSWORD REQUESTS AUDIT TABLE
-- ==============================================================================

CREATE TABLE IF NOT EXISTS identity.password_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'DELIVERED',
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_password_requests_user ON identity.password_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_password_requests_email ON identity.password_requests(email);
CREATE INDEX IF NOT EXISTS idx_password_requests_created ON identity.password_requests(created_at);

-- +goose Down
DROP TABLE IF EXISTS identity.password_requests CASCADE;
