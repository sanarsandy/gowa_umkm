-- Migration 012: Multi-Provider AI Support
-- Allows users to choose different AI providers and use their own API keys

-- Add new columns to ai_configs
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS ai_provider VARCHAR(50) DEFAULT 'gemini';
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS user_api_key TEXT; -- Encrypted
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS api_key_set BOOLEAN DEFAULT FALSE;
ALTER TABLE ai_configs ADD COLUMN IF NOT EXISTS use_system_key BOOLEAN DEFAULT TRUE;

-- Update model column to allow longer model names
ALTER TABLE ai_configs ALTER COLUMN model TYPE VARCHAR(100);

-- Add constraint for valid providers
ALTER TABLE ai_configs DROP CONSTRAINT IF EXISTS valid_ai_provider;
ALTER TABLE ai_configs ADD CONSTRAINT valid_ai_provider 
    CHECK (ai_provider IN ('gemini', 'openai', 'groq', 'anthropic'));

-- Create index for provider lookup
CREATE INDEX IF NOT EXISTS idx_ai_configs_provider ON ai_configs(ai_provider);

-- Comments
COMMENT ON COLUMN ai_configs.ai_provider IS 'AI provider: gemini, openai, groq, anthropic';
COMMENT ON COLUMN ai_configs.user_api_key IS 'User-provided API key (encrypted)';
COMMENT ON COLUMN ai_configs.api_key_set IS 'Whether user has set their own API key';
COMMENT ON COLUMN ai_configs.use_system_key IS 'Use system API key instead of user key';

