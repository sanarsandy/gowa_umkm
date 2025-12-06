-- Add media columns to knowledge_base table
ALTER TABLE knowledge_base ADD COLUMN IF NOT EXISTS media_url TEXT;
ALTER TABLE knowledge_base ADD COLUMN IF NOT EXISTS media_type TEXT; -- image, document, video, etc.

-- Add media columns to whatsapp_messages table
ALTER TABLE whatsapp_messages ADD COLUMN IF NOT EXISTS media_url TEXT;
ALTER TABLE whatsapp_messages ADD COLUMN IF NOT EXISTS caption TEXT;
-- message_type already exists in whatsapp_messages
