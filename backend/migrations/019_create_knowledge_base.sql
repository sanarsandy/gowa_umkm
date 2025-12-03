-- Migration 019: Create Knowledge Base Table
-- Creates knowledge_base table for AI auto-reply knowledge management

-- Knowledge Base table (simplified version without vector extension)
CREATE TABLE IF NOT EXISTS knowledge_base (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Content
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(100), -- faq, product, policy, pricing, shipping, hours, location
    
    -- Metadata for search
    keywords TEXT[], -- manual keywords for quick search
    tags TEXT[], -- custom tags for organization
    priority INTEGER DEFAULT 5 CHECK (priority BETWEEN 1 AND 10), -- 1=low, 10=high
    
    -- Usage tracking
    usage_count INTEGER DEFAULT 0,
    last_used_at TIMESTAMP,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for knowledge_base
CREATE INDEX IF NOT EXISTS idx_knowledge_base_tenant_id ON knowledge_base(tenant_id);
CREATE INDEX IF NOT EXISTS idx_knowledge_base_category ON knowledge_base(category);
CREATE INDEX IF NOT EXISTS idx_knowledge_base_is_active ON knowledge_base(is_active) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_knowledge_base_keywords ON knowledge_base USING GIN(keywords);
CREATE INDEX IF NOT EXISTS idx_knowledge_base_tags ON knowledge_base USING GIN(tags);

-- Comments
COMMENT ON TABLE knowledge_base IS 'Knowledge base for AI auto-reply system';
COMMENT ON COLUMN knowledge_base.priority IS 'Priority for ranking search results (1-10)';
COMMENT ON COLUMN knowledge_base.keywords IS 'Keywords for quick search matching';
COMMENT ON COLUMN knowledge_base.tags IS 'Custom tags for organization';
