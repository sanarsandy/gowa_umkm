-- Migration 011: Alter AI Configs for Gemini Integration
-- Adds new columns for Gemini AI and auto-reply settings

-- Add enabled column (rename from auto_reply_enabled)
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS enabled BOOLEAN DEFAULT FALSE;
UPDATE ai_configs SET enabled = auto_reply_enabled WHERE enabled IS NULL;

-- Add model column (use gemini instead of openai)
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS model VARCHAR(100) DEFAULT 'gemini-1.5-flash';

-- Add confidence threshold
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS confidence_threshold DECIMAL(3,2) DEFAULT 0.80;

-- Add language
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS language VARCHAR(10) DEFAULT 'id';

-- Rename/add system_prompt column
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS system_prompt TEXT;
UPDATE ai_configs SET system_prompt = custom_system_prompt WHERE system_prompt IS NULL AND custom_system_prompt IS NOT NULL;

-- Add business context columns
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS business_name VARCHAR(255);
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS business_type VARCHAR(100);
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS business_description TEXT;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS business_address TEXT;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS payment_methods VARCHAR(500);

-- Add escalation settings
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS escalate_low_confidence BOOLEAN DEFAULT TRUE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS escalate_complaint BOOLEAN DEFAULT TRUE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS escalate_order BOOLEAN DEFAULT FALSE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS escalate_urgent BOOLEAN DEFAULT TRUE;

-- Add notification settings
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS notify_whatsapp BOOLEAN DEFAULT TRUE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS notify_email BOOLEAN DEFAULT FALSE;

-- Add usage tracking
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS total_requests INTEGER DEFAULT 0;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS total_tokens_used BIGINT DEFAULT 0;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS total_cost_usd DECIMAL(12,6) DEFAULT 0;

-- Create index for enabled configs
CREATE INDEX IF NOT EXISTS idx_ai_configs_enabled ON ai_configs(enabled) WHERE enabled = TRUE;

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS update_ai_configs_updated_at ON ai_configs;
CREATE TRIGGER update_ai_configs_updated_at
    BEFORE UPDATE ON ai_configs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Update comments
COMMENT ON COLUMN ai_configs.enabled IS 'Whether AI auto-reply is enabled';
COMMENT ON COLUMN ai_configs.model IS 'Gemini model to use (gemini-1.5-flash, gemini-1.5-pro)';
COMMENT ON COLUMN ai_configs.confidence_threshold IS 'Minimum confidence to auto-reply (0.50-0.95)';

