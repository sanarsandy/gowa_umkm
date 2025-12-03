-- Migration 017: Update AI Configs Table Schema
-- Adds missing columns required by the AI handler

-- Add new columns for AI configuration
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS ai_provider VARCHAR(50) DEFAULT 'gemini';
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS model VARCHAR(100) DEFAULT 'gemini-2.0-flash';
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS user_api_key TEXT;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS api_key_set BOOLEAN DEFAULT FALSE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS use_system_key BOOLEAN DEFAULT TRUE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS confidence_threshold DECIMAL(5,4) DEFAULT 0.80;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS language VARCHAR(10) DEFAULT 'id';
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS system_prompt TEXT;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS business_name VARCHAR(255);
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS business_type VARCHAR(100);
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS business_description TEXT;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS business_address TEXT;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS payment_methods TEXT;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS escalate_low_confidence BOOLEAN DEFAULT TRUE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS escalate_complaint BOOLEAN DEFAULT TRUE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS escalate_order BOOLEAN DEFAULT FALSE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS escalate_urgent BOOLEAN DEFAULT TRUE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS notify_whatsapp BOOLEAN DEFAULT TRUE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS notify_email BOOLEAN DEFAULT FALSE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS total_requests INTEGER DEFAULT 0;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS total_tokens_used BIGINT DEFAULT 0;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS total_cost_usd DECIMAL(10,6) DEFAULT 0;

-- Rename old columns to avoid conflicts (keep for backward compatibility)
ALTER TABLE ai_configs RENAME COLUMN ai_model TO old_ai_model;
ALTER TABLE ai_configs RENAME COLUMN auto_reply_enabled TO old_auto_reply_enabled;
ALTER TABLE ai_configs RENAME COLUMN ai_temperature TO old_ai_temperature;

-- Update max_tokens default if needed
ALTER TABLE ai_configs ALTER COLUMN max_tokens SET DEFAULT 200;

-- Create indexes for new columns
CREATE INDEX IF NOT EXISTS idx_ai_configs_enabled ON ai_configs(enabled) WHERE enabled = TRUE;
CREATE INDEX IF NOT EXISTS idx_ai_configs_ai_provider ON ai_configs(ai_provider);

COMMENT ON COLUMN ai_configs.enabled IS 'Whether AI auto-reply is enabled for this tenant';
COMMENT ON COLUMN ai_configs.ai_provider IS 'AI provider (gemini, openai, etc)';
COMMENT ON COLUMN ai_configs.model IS 'AI model to use';
COMMENT ON COLUMN ai_configs.use_system_key IS 'Whether to use system API key or user-provided key';
COMMENT ON COLUMN ai_configs.confidence_threshold IS 'Minimum confidence score to auto-reply (0.0-1.0)';
