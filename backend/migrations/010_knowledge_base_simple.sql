-- Migration 010: Knowledge Base (Simple - No pgvector)
-- This is a fallback migration for knowledge_base without vector embeddings

-- Knowledge Base table (without vector column)
CREATE TABLE IF NOT EXISTS knowledge_base (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Content
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(100), -- faq, product, policy, pricing, shipping, hours, location, payment
    
    -- Metadata for search
    keywords TEXT[], -- manual keywords for quick search
    tags TEXT[], -- custom tags for organization
    priority INTEGER DEFAULT 5 CHECK (priority BETWEEN 1 AND 10), -- 1=low, 10=high
    
    -- Usage tracking
    usage_count INTEGER DEFAULT 0,
    last_used_at TIMESTAMPTZ,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- AI Configuration table
CREATE TABLE IF NOT EXISTS ai_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL UNIQUE REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Basic settings
    enabled BOOLEAN DEFAULT FALSE,
    model VARCHAR(100) DEFAULT 'gemini-1.5-flash',
    confidence_threshold DECIMAL(3,2) DEFAULT 0.80,
    max_tokens INTEGER DEFAULT 200,
    language VARCHAR(10) DEFAULT 'id',
    
    -- System prompt
    system_prompt TEXT,
    
    -- Business context
    business_name VARCHAR(255),
    business_type VARCHAR(100),
    business_hours VARCHAR(255),
    business_description TEXT,
    business_address TEXT,
    payment_methods VARCHAR(500),
    
    -- Escalation settings
    escalate_low_confidence BOOLEAN DEFAULT TRUE,
    escalate_complaint BOOLEAN DEFAULT TRUE,
    escalate_order BOOLEAN DEFAULT FALSE,
    escalate_urgent BOOLEAN DEFAULT TRUE,
    notify_whatsapp BOOLEAN DEFAULT TRUE,
    notify_email BOOLEAN DEFAULT FALSE,
    
    -- Usage tracking
    total_requests INTEGER DEFAULT 0,
    total_tokens_used BIGINT DEFAULT 0,
    total_cost_usd DECIMAL(12,6) DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for knowledge_base
CREATE INDEX IF NOT EXISTS idx_knowledge_base_tenant_id ON knowledge_base(tenant_id);
CREATE INDEX IF NOT EXISTS idx_knowledge_base_category ON knowledge_base(category);
CREATE INDEX IF NOT EXISTS idx_knowledge_base_is_active ON knowledge_base(is_active) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_knowledge_base_keywords ON knowledge_base USING GIN(keywords);
CREATE INDEX IF NOT EXISTS idx_knowledge_base_priority ON knowledge_base(priority DESC);

-- Indexes for ai_configs
CREATE INDEX IF NOT EXISTS idx_ai_configs_tenant_id ON ai_configs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ai_configs_enabled ON ai_configs(enabled) WHERE enabled = TRUE;

-- Update function for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
DROP TRIGGER IF EXISTS update_knowledge_base_updated_at ON knowledge_base;
CREATE TRIGGER update_knowledge_base_updated_at
    BEFORE UPDATE ON knowledge_base
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_ai_configs_updated_at ON ai_configs;
CREATE TRIGGER update_ai_configs_updated_at
    BEFORE UPDATE ON ai_configs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE knowledge_base IS 'Knowledge base for AI auto-reply system';
COMMENT ON TABLE ai_configs IS 'AI configuration per tenant';
COMMENT ON COLUMN knowledge_base.priority IS 'Priority for ranking search results (1-10, 10=highest)';
COMMENT ON COLUMN ai_configs.confidence_threshold IS 'Minimum confidence to auto-reply (0.00-1.00)';

