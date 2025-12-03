-- Migration 013: Customer Tags & Analytics
-- Adds customer segmentation and analytics support

-- Customer Tags table
CREATE TABLE IF NOT EXISTS customer_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    color VARCHAR(20) DEFAULT '#6366f1',
    description TEXT,
    auto_apply_rules JSONB, -- Rules for auto-applying this tag
    customer_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

-- Customer-Tag relationship (many-to-many)
CREATE TABLE IF NOT EXISTS customer_tag_assignments (
    customer_id UUID NOT NULL REFERENCES customer_insights(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES customer_tags(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT NOW(),
    assigned_by VARCHAR(50) DEFAULT 'manual', -- 'manual', 'auto', 'ai'
    PRIMARY KEY (customer_id, tag_id)
);

-- Customer notes
CREATE TABLE IF NOT EXISTS customer_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customer_insights(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_by UUID, -- user who created the note
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add lead_score and other columns to customer_insights (not customers)
ALTER TABLE customer_insights ADD COLUMN IF NOT EXISTS lead_score INT DEFAULT 0;
ALTER TABLE customer_insights ADD COLUMN IF NOT EXISTS lead_status VARCHAR(20) DEFAULT 'new';
ALTER TABLE customer_insights ADD COLUMN IF NOT EXISTS last_interaction_at TIMESTAMP;
ALTER TABLE customer_insights ADD COLUMN IF NOT EXISTS total_messages INT DEFAULT 0;
ALTER TABLE customer_insights ADD COLUMN IF NOT EXISTS first_message_at_new TIMESTAMP;

-- Update total_messages from message_count if not set
UPDATE customer_insights SET total_messages = message_count WHERE total_messages = 0 OR total_messages IS NULL;

-- Update first_message_at_new from first_message_at if not set
UPDATE customer_insights SET first_message_at_new = first_message_at WHERE first_message_at_new IS NULL;

-- Daily analytics aggregation table
CREATE TABLE IF NOT EXISTS daily_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    -- Message metrics
    messages_received INT DEFAULT 0,
    messages_sent INT DEFAULT 0,
    unique_customers INT DEFAULT 0,
    new_customers INT DEFAULT 0,
    -- AI metrics
    ai_responses INT DEFAULT 0,
    ai_escalations INT DEFAULT 0,
    ai_avg_confidence DECIMAL(5,4) DEFAULT 0,
    ai_tokens_used INT DEFAULT 0,
    ai_cost_usd DECIMAL(10,6) DEFAULT 0,
    -- Broadcast metrics
    broadcasts_sent INT DEFAULT 0,
    broadcast_recipients INT DEFAULT 0,
    -- Customer metrics
    active_customers INT DEFAULT 0,
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, date)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_customer_tags_tenant ON customer_tags(tenant_id);
CREATE INDEX IF NOT EXISTS idx_customer_tag_assignments_customer ON customer_tag_assignments(customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_tag_assignments_tag ON customer_tag_assignments(tag_id);
CREATE INDEX IF NOT EXISTS idx_customer_notes_customer ON customer_notes(customer_id);
CREATE INDEX IF NOT EXISTS idx_daily_analytics_tenant_date ON daily_analytics(tenant_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_customer_insights_lead_score ON customer_insights(tenant_id, lead_score DESC);
CREATE INDEX IF NOT EXISTS idx_customer_insights_lead_status ON customer_insights(tenant_id, lead_status);

-- Create default tags for existing tenants
INSERT INTO customer_tags (tenant_id, name, color, description) 
SELECT id, 'VIP', '#f59e0b', 'Pelanggan VIP/premium'
FROM tenants
ON CONFLICT DO NOTHING;

INSERT INTO customer_tags (tenant_id, name, color, description) 
SELECT id, 'Hot Lead', '#ef4444', 'Prospek dengan minat tinggi'
FROM tenants
ON CONFLICT DO NOTHING;

INSERT INTO customer_tags (tenant_id, name, color, description) 
SELECT id, 'Repeat Customer', '#10b981', 'Pelanggan yang sudah pernah beli'
FROM tenants
ON CONFLICT DO NOTHING;

INSERT INTO customer_tags (tenant_id, name, color, description) 
SELECT id, 'New', '#3b82f6', 'Pelanggan baru'
FROM tenants
ON CONFLICT DO NOTHING;

-- Comments
COMMENT ON TABLE customer_tags IS 'Custom tags for customer segmentation';
COMMENT ON TABLE customer_tag_assignments IS 'Many-to-many relationship between customers and tags';
COMMENT ON TABLE customer_notes IS 'Notes/comments about customers';
COMMENT ON TABLE daily_analytics IS 'Aggregated daily analytics per tenant';

