-- Migration 002: Tenants Table
-- This table stores business information for each UMKM tenant

CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    business_name VARCHAR(255) NOT NULL,
    business_type VARCHAR(100),
    business_description TEXT,
    business_phone VARCHAR(50),
    business_address TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_tenants_user_id ON tenants(user_id);
CREATE INDEX IF NOT EXISTS idx_tenants_is_active ON tenants(is_active);

-- Ensure one user can have multiple tenants (for future multi-business support)
-- But typically one user = one tenant for MVP
COMMENT ON TABLE tenants IS 'Stores UMKM business information for each tenant';
COMMENT ON COLUMN tenants.user_id IS 'Reference to the user who owns this tenant';
COMMENT ON COLUMN tenants.is_active IS 'Whether this tenant subscription is active';
