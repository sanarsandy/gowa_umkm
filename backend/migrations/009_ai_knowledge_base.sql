-- Migration 009: AI Knowledge Base & Conversation Logs
-- Stores knowledge base for AI auto-reply and conversation logs

-- NOTE: Vector extension disabled - using simplified version in migration 019
-- Enable pgvector extension for vector similarity search
-- CREATE EXTENSION IF NOT EXISTS vector;

-- Knowledge Base table
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
    
    -- Vector embedding for semantic search (OpenAI embedding dimension)
    -- embedding vector(1536),  -- DISABLED: requires pgvector extension
    
    -- Usage tracking
    usage_count INTEGER DEFAULT 0,
    last_used_at TIMESTAMPTZ,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- AI Conversation Logs table
CREATE TABLE IF NOT EXISTS ai_conversation_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id UUID REFERENCES customer_insights(id) ON DELETE SET NULL,
    
    -- Message info
    customer_message TEXT NOT NULL,
    ai_response TEXT,
    
    -- AI processing
    detected_intent VARCHAR(100), -- price_inquiry, location_inquiry, etc
    confidence_score DECIMAL(5,4), -- 0.0000 to 1.0000
    knowledge_used UUID[], -- array of knowledge_base IDs used
    
    -- Decision
    action_taken VARCHAR(50) CHECK (action_taken IN ('auto_replied', 'escalated', 'queued', 'failed')),
    escalation_reason TEXT,
    
    -- Performance metrics
    response_time_ms INTEGER,
    tokens_used INTEGER,
    input_tokens INTEGER,
    output_tokens INTEGER,
    cost_usd DECIMAL(10,6),
    model_used VARCHAR(100), -- gemini-1.5-flash, gemini-1.5-pro, etc
    
    -- Feedback
    customer_satisfaction INTEGER CHECK (customer_satisfaction BETWEEN 1 AND 5),
    human_override BOOLEAN DEFAULT FALSE,
    human_feedback TEXT,
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Prompt Templates table
CREATE TABLE IF NOT EXISTS prompt_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Template info
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100), -- greeting, pricing, shipping, complaint, general
    
    -- Prompt content
    prompt_text TEXT NOT NULL,
    variables JSONB DEFAULT '{}'::jsonb, -- {"product_name": "string", "price": "number"}
    
    -- Usage
    is_default BOOLEAN DEFAULT FALSE,
    usage_count INTEGER DEFAULT 0,
    
    -- A/B Testing
    variant_of UUID REFERENCES prompt_templates(id) ON DELETE SET NULL,
    conversion_rate DECIMAL(5,4),
    avg_confidence DECIMAL(5,4),
    
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for knowledge_base
CREATE INDEX idx_knowledge_base_tenant_id ON knowledge_base(tenant_id);
CREATE INDEX idx_knowledge_base_category ON knowledge_base(category);
CREATE INDEX idx_knowledge_base_is_active ON knowledge_base(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_knowledge_base_keywords ON knowledge_base USING GIN(keywords);
CREATE INDEX idx_knowledge_base_tags ON knowledge_base USING GIN(tags);
-- Vector similarity search index (IVFFlat for faster approximate search)
-- CREATE INDEX idx_knowledge_base_embedding ON knowledge_base USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);  -- DISABLED: requires pgvector

-- Indexes for ai_conversation_logs
CREATE INDEX idx_ai_logs_tenant_id ON ai_conversation_logs(tenant_id);
CREATE INDEX idx_ai_logs_customer_id ON ai_conversation_logs(customer_id);
CREATE INDEX idx_ai_logs_created_at ON ai_conversation_logs(created_at DESC);
CREATE INDEX idx_ai_logs_action_taken ON ai_conversation_logs(action_taken);
CREATE INDEX idx_ai_logs_detected_intent ON ai_conversation_logs(detected_intent);

-- Indexes for prompt_templates
CREATE INDEX idx_prompt_templates_tenant_id ON prompt_templates(tenant_id);
CREATE INDEX idx_prompt_templates_category ON prompt_templates(category);
CREATE INDEX idx_prompt_templates_is_default ON prompt_templates(is_default) WHERE is_default = TRUE;

-- Comments
COMMENT ON TABLE knowledge_base IS 'Knowledge base for AI auto-reply system';
COMMENT ON TABLE ai_conversation_logs IS 'Logs of all AI-powered conversations';
COMMENT ON TABLE prompt_templates IS 'Reusable prompt templates for AI responses';

-- COMMENT ON COLUMN knowledge_base.embedding IS 'Vector embedding for semantic similarity search';  -- DISABLED
COMMENT ON COLUMN knowledge_base.priority IS 'Priority for ranking search results (1-10)';
COMMENT ON COLUMN ai_conversation_logs.confidence_score IS 'AI confidence in the response (0.0-1.0)';
COMMENT ON COLUMN ai_conversation_logs.knowledge_used IS 'Array of knowledge base entry IDs used to generate response';
