-- Migration 014: Create Customers Table
-- Creates the main customers table that was missing

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Customer identification
    jid VARCHAR(255) NOT NULL, -- WhatsApp JID
    name VARCHAR(255),
    phone VARCHAR(50),
    email VARCHAR(255),
    
    -- Status and classification
    status VARCHAR(50) DEFAULT 'new' CHECK (status IN ('new', 'hot_lead', 'warm_lead', 'cold_lead', 'customer', 'complaint', 'spam')),
    lead_score INT DEFAULT 0,
    lead_status VARCHAR(20) DEFAULT 'new',
    sentiment VARCHAR(50) CHECK (sentiment IN ('positive', 'neutral', 'negative', 'mixed')),
    
    -- Engagement metrics
    total_messages INT DEFAULT 0,
    last_message_at TIMESTAMP,
    first_message_at TIMESTAMP,
    last_interaction_at TIMESTAMP,
    
    -- Follow-up tracking
    needs_follow_up BOOLEAN DEFAULT FALSE,
    follow_up_scheduled_at TIMESTAMP,
    
    -- Additional info
    tags JSONB DEFAULT '[]'::jsonb,
    notes TEXT,
    
    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint per tenant-customer pair
    CONSTRAINT unique_tenant_customer_jid UNIQUE(tenant_id, jid)
);

-- Indexes for customers table
CREATE INDEX IF NOT EXISTS idx_customers_tenant_id ON customers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_customers_jid ON customers(jid);
CREATE INDEX IF NOT EXISTS idx_customers_status ON customers(status);
CREATE INDEX IF NOT EXISTS idx_customers_lead_score ON customers(tenant_id, lead_score DESC);
CREATE INDEX IF NOT EXISTS idx_customers_lead_status ON customers(tenant_id, lead_status);
CREATE INDEX IF NOT EXISTS idx_customers_last_message_at ON customers(last_message_at DESC);
CREATE INDEX IF NOT EXISTS idx_customers_needs_follow_up ON customers(needs_follow_up) WHERE needs_follow_up = TRUE;

-- Migrate data from customer_insights to customers if exists
INSERT INTO customers (
    id, tenant_id, jid, name, phone, status, sentiment, 
    total_messages, last_message_at, first_message_at, 
    needs_follow_up, follow_up_scheduled_at, tags, 
    created_at, updated_at
)
SELECT 
    id, tenant_id, customer_jid as jid, customer_name as name, 
    customer_phone as phone, status, sentiment,
    message_count as total_messages, last_message_at, first_message_at,
    needs_follow_up, follow_up_scheduled_at, tags,
    created_at, updated_at
FROM customer_insights
ON CONFLICT (tenant_id, jid) DO UPDATE SET
    name = EXCLUDED.name,
    phone = EXCLUDED.phone,
    status = EXCLUDED.status,
    sentiment = EXCLUDED.sentiment,
    total_messages = EXCLUDED.total_messages,
    last_message_at = EXCLUDED.last_message_at,
    first_message_at = EXCLUDED.first_message_at,
    needs_follow_up = EXCLUDED.needs_follow_up,
    follow_up_scheduled_at = EXCLUDED.follow_up_scheduled_at,
    tags = EXCLUDED.tags,
    updated_at = NOW();

COMMENT ON TABLE customers IS 'Main customers table with contact info and engagement metrics';
COMMENT ON COLUMN customers.jid IS 'WhatsApp JID (unique identifier)';
COMMENT ON COLUMN customers.lead_score IS 'Lead scoring (0-100)';
COMMENT ON COLUMN customers.status IS 'Customer status/classification';
