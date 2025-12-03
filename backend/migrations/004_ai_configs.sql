-- Migration 004: AI Configurations Table
-- Stores AI configuration and business context for each tenant

CREATE TABLE IF NOT EXISTS ai_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- Business context for AI
    business_context TEXT,
    business_hours VARCHAR(255),
    
    -- Menu and pricing data (stored as JSONB for flexibility)
    menu_data JSONB DEFAULT '[]'::jsonb,
    
    -- FAQ data
    faq_data JSONB DEFAULT '[]'::jsonb,
    
    -- AI behavior settings
    active_mode VARCHAR(50) DEFAULT 'general' CHECK (active_mode IN ('general', 'order_taking', 'faq', 'complaint', 'custom')),
    auto_reply_enabled BOOLEAN DEFAULT FALSE,
    auto_reply_message TEXT,
    
    -- AI model settings
    ai_model VARCHAR(50) DEFAULT 'gpt-3.5-turbo',
    ai_temperature DECIMAL(3,2) DEFAULT 0.7 CHECK (ai_temperature >= 0 AND ai_temperature <= 2),
    max_tokens INTEGER DEFAULT 500,
    
    -- Custom system prompt (optional override)
    custom_system_prompt TEXT,
    
    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure one config per tenant
    CONSTRAINT unique_tenant_ai_config UNIQUE(tenant_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_ai_configs_tenant_id ON ai_configs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ai_configs_active_mode ON ai_configs(active_mode);

-- Example JSONB structure for menu_data:
-- [
--   {
--     "id": "1",
--     "name": "Nasi Goreng Spesial",
--     "price": 25000,
--     "category": "main_course",
--     "description": "Nasi goreng dengan telur, ayam, dan sayuran",
--     "available": true
--   }
-- ]

-- Example JSONB structure for faq_data:
-- [
--   {
--     "question": "Jam buka?",
--     "answer": "Kami buka setiap hari pukul 10:00 - 22:00",
--     "keywords": ["jam", "buka", "operasional"]
--   }
-- ]

COMMENT ON TABLE ai_configs IS 'AI configuration and business context per tenant';
COMMENT ON COLUMN ai_configs.menu_data IS 'Product/service menu in JSONB format';
COMMENT ON COLUMN ai_configs.faq_data IS 'Frequently asked questions in JSONB format';
COMMENT ON COLUMN ai_configs.active_mode IS 'Current AI operation mode';
COMMENT ON COLUMN ai_configs.ai_temperature IS 'OpenAI temperature parameter (0-2)';
