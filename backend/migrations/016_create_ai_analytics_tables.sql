-- Migration 016: Create Missing AI and Analytics Tables
-- Creates ai_conversation_logs and daily_analytics tables that are missing

-- AI Conversation Logs table (for analytics and AI tracking)
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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

-- Indexes for ai_conversation_logs
CREATE INDEX IF NOT EXISTS idx_ai_logs_tenant_id ON ai_conversation_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ai_logs_customer_id ON ai_conversation_logs(customer_id);
CREATE INDEX IF NOT EXISTS idx_ai_logs_created_at ON ai_conversation_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ai_logs_action_taken ON ai_conversation_logs(action_taken);
CREATE INDEX IF NOT EXISTS idx_ai_logs_detected_intent ON ai_conversation_logs(detected_intent);

-- Indexes for daily_analytics
CREATE INDEX IF NOT EXISTS idx_daily_analytics_tenant_date ON daily_analytics(tenant_id, date DESC);

-- Comments
COMMENT ON TABLE ai_conversation_logs IS 'Logs of all AI-powered conversations for analytics and tracking';
COMMENT ON TABLE daily_analytics IS 'Aggregated daily analytics per tenant';
COMMENT ON COLUMN ai_conversation_logs.confidence_score IS 'AI confidence in the response (0.0-1.0)';
COMMENT ON COLUMN ai_conversation_logs.knowledge_used IS 'Array of knowledge base entry IDs used to generate response';
