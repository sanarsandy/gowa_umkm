-- Migration 003: WhatsApp Devices Table
-- This table stores WhatsApp session data for whatsmeow library
-- Replaces the default SQLite storage with PostgreSQL for multi-tenant SaaS

CREATE TABLE IF NOT EXISTS whatsapp_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- WhatsApp JID (Jabber ID) - unique identifier for the device
    jid VARCHAR(255) UNIQUE NOT NULL,
    
    -- Session data (encrypted keys for whatsmeow)
    registration_id INTEGER,
    adv_secret_key BYTEA,
    next_pre_key_id INTEGER,
    first_unuploaded_pre_key_id INTEGER,
    account_sync_timestamp BIGINT,
    account_settings BYTEA,
    
    -- Identity keys
    identity_key BYTEA,
    
    -- Connection status
    is_connected BOOLEAN DEFAULT FALSE,
    last_connected_at TIMESTAMP,
    connection_error TEXT,
    
    -- Platform information
    platform VARCHAR(50) DEFAULT 'web',
    business_name VARCHAR(255),
    push_name VARCHAR(255),
    
    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure one tenant can only have one active WhatsApp device
    CONSTRAINT unique_tenant_device UNIQUE(tenant_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_whatsapp_devices_tenant_id ON whatsapp_devices(tenant_id);
CREATE INDEX IF NOT EXISTS idx_whatsapp_devices_jid ON whatsapp_devices(jid);
CREATE INDEX IF NOT EXISTS idx_whatsapp_devices_is_connected ON whatsapp_devices(is_connected);

-- Additional tables for whatsmeow session storage
-- These mirror the structure used by whatsmeow's sqlstore

CREATE TABLE IF NOT EXISTS whatsapp_identity_keys (
    our_jid VARCHAR(255) NOT NULL,
    their_id VARCHAR(255) NOT NULL,
    identity BYTEA NOT NULL,
    PRIMARY KEY (our_jid, their_id)
);

CREATE TABLE IF NOT EXISTS whatsapp_pre_keys (
    jid VARCHAR(255) NOT NULL,
    key_id INTEGER NOT NULL,
    key BYTEA NOT NULL,
    uploaded BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (jid, key_id)
);

CREATE TABLE IF NOT EXISTS whatsapp_sessions (
    our_jid VARCHAR(255) NOT NULL,
    their_id VARCHAR(255) NOT NULL,
    session BYTEA NOT NULL,
    PRIMARY KEY (our_jid, their_id)
);

CREATE TABLE IF NOT EXISTS whatsapp_sender_keys (
    our_jid VARCHAR(255) NOT NULL,
    chat_id VARCHAR(255) NOT NULL,
    sender_id VARCHAR(255) NOT NULL,
    sender_key BYTEA NOT NULL,
    PRIMARY KEY (our_jid, chat_id, sender_id)
);

CREATE TABLE IF NOT EXISTS whatsapp_app_state_sync_keys (
    jid VARCHAR(255) NOT NULL,
    key_id BYTEA NOT NULL,
    key_data BYTEA NOT NULL,
    timestamp BIGINT NOT NULL,
    fingerprint BYTEA NOT NULL,
    PRIMARY KEY (jid, key_id)
);

CREATE TABLE IF NOT EXISTS whatsapp_app_state_version (
    jid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    version BIGINT NOT NULL,
    hash BYTEA NOT NULL,
    PRIMARY KEY (jid, name)
);

CREATE TABLE IF NOT EXISTS whatsapp_app_state_mutations (
    jid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    version BIGINT NOT NULL,
    index_mac BYTEA,
    value_mac BYTEA,
    PRIMARY KEY (jid, name, version, index_mac)
);

CREATE TABLE IF NOT EXISTS whatsapp_contacts (
    our_jid VARCHAR(255) NOT NULL,
    their_jid VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    full_name VARCHAR(255),
    push_name VARCHAR(255),
    business_name VARCHAR(255),
    PRIMARY KEY (our_jid, their_jid)
);

CREATE TABLE IF NOT EXISTS whatsapp_chat_settings (
    our_jid VARCHAR(255) NOT NULL,
    chat_jid VARCHAR(255) NOT NULL,
    muted_until BIGINT NOT NULL DEFAULT 0,
    pinned BOOLEAN NOT NULL DEFAULT FALSE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (our_jid, chat_jid)
);

-- Indexes for whatsmeow tables
CREATE INDEX IF NOT EXISTS idx_whatsapp_identity_keys_our_jid ON whatsapp_identity_keys(our_jid);
CREATE INDEX IF NOT EXISTS idx_whatsapp_sessions_our_jid ON whatsapp_sessions(our_jid);
CREATE INDEX IF NOT EXISTS idx_whatsapp_sender_keys_our_jid ON whatsapp_sender_keys(our_jid);

COMMENT ON TABLE whatsapp_devices IS 'Main table for WhatsApp device sessions per tenant';
COMMENT ON TABLE whatsapp_identity_keys IS 'Stores identity keys for E2E encryption';
COMMENT ON TABLE whatsapp_sessions IS 'Stores session data for encrypted conversations';
COMMENT ON TABLE whatsapp_sender_keys IS 'Stores sender keys for group chats';
