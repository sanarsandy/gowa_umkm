-- Migration 005: Customer Insights & Messages Tables
-- Stores AI analysis results and message history

-- Customer insights from AI analysis
CREATE TABLE IF NOT EXISTS customer_insights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Customer identification
    customer_jid VARCHAR(255) NOT NULL,
    customer_name VARCHAR(255),
    customer_phone VARCHAR(50),
    
    -- AI analysis results
    status VARCHAR(50) DEFAULT 'new' CHECK (status IN ('new', 'hot_lead', 'warm_lead', 'cold_lead', 'customer', 'complaint', 'spam')),
    sentiment VARCHAR(50) CHECK (sentiment IN ('positive', 'neutral', 'negative', 'mixed')),
    intent VARCHAR(100),
    
    -- Detected interests
    product_interest JSONB DEFAULT '[]'::jsonb,
    
    -- Conversation summary
    last_message_summary TEXT,
    conversation_context TEXT,
    
    -- Engagement metrics
    message_count INTEGER DEFAULT 0,
    last_message_at TIMESTAMP,
    first_message_at TIMESTAMP,
    
    -- Follow-up tracking
    needs_follow_up BOOLEAN DEFAULT FALSE,
    follow_up_scheduled_at TIMESTAMP,
    follow_up_completed BOOLEAN DEFAULT FALSE,
    
    -- Tags (for custom categorization)
    tags JSONB DEFAULT '[]'::jsonb,
    
    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint per tenant-customer pair
    CONSTRAINT unique_tenant_customer UNIQUE(tenant_id, customer_jid)
);

-- Message history table
CREATE TABLE IF NOT EXISTS whatsapp_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Message identification
    message_id VARCHAR(255) NOT NULL,
    chat_jid VARCHAR(255) NOT NULL,
    sender_jid VARCHAR(255) NOT NULL,
    
    -- Message content
    message_type VARCHAR(50) DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'video', 'audio', 'document', 'sticker', 'location', 'contact')),
    message_text TEXT,
    media_url TEXT,
    media_mime_type VARCHAR(100),
    
    -- Message metadata
    is_from_me BOOLEAN DEFAULT FALSE,
    is_group BOOLEAN DEFAULT FALSE,
    timestamp BIGINT NOT NULL,
    
    -- AI processing status
    ai_processed BOOLEAN DEFAULT FALSE,
    ai_processed_at TIMESTAMP,
    
    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint to prevent duplicate messages
    CONSTRAINT unique_message UNIQUE(tenant_id, message_id)
);

-- Scheduled messages for auto-follow-up
CREATE TABLE IF NOT EXISTS scheduled_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Target
    recipient_jid VARCHAR(255) NOT NULL,
    
    -- Message content
    message_text TEXT NOT NULL,
    
    -- Scheduling
    scheduled_at TIMESTAMP NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed', 'cancelled')),
    
    -- Execution tracking
    sent_at TIMESTAMP,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    
    -- Reference to insight (optional)
    insight_id UUID REFERENCES customer_insights(id) ON DELETE SET NULL,
    
    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for customer_insights
CREATE INDEX IF NOT EXISTS idx_customer_insights_tenant_id ON customer_insights(tenant_id);
CREATE INDEX IF NOT EXISTS idx_customer_insights_customer_jid ON customer_insights(customer_jid);
CREATE INDEX IF NOT EXISTS idx_customer_insights_status ON customer_insights(status);
CREATE INDEX IF NOT EXISTS idx_customer_insights_needs_follow_up ON customer_insights(needs_follow_up);
CREATE INDEX IF NOT EXISTS idx_customer_insights_last_message_at ON customer_insights(last_message_at DESC);

-- Indexes for whatsapp_messages
CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_tenant_id ON whatsapp_messages(tenant_id);
CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_chat_jid ON whatsapp_messages(chat_jid);
CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_sender_jid ON whatsapp_messages(sender_jid);
CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_timestamp ON whatsapp_messages(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_whatsapp_messages_ai_processed ON whatsapp_messages(ai_processed) WHERE ai_processed = FALSE;

-- Indexes for scheduled_messages
CREATE INDEX IF NOT EXISTS idx_scheduled_messages_tenant_id ON scheduled_messages(tenant_id);
CREATE INDEX IF NOT EXISTS idx_scheduled_messages_status ON scheduled_messages(status);
CREATE INDEX IF NOT EXISTS idx_scheduled_messages_scheduled_at ON scheduled_messages(scheduled_at);

-- Partitioning hint for future optimization (when message volume grows)
-- whatsapp_messages can be partitioned by timestamp (monthly/yearly)

COMMENT ON TABLE customer_insights IS 'AI-generated insights about customers';
COMMENT ON TABLE whatsapp_messages IS 'Message history for all WhatsApp conversations';
COMMENT ON TABLE scheduled_messages IS 'Queue for automated follow-up messages';
COMMENT ON COLUMN customer_insights.status IS 'Lead qualification status from AI analysis';
COMMENT ON COLUMN customer_insights.product_interest IS 'Array of products/services customer showed interest in';
COMMENT ON COLUMN whatsapp_messages.ai_processed IS 'Whether this message has been analyzed by AI';
