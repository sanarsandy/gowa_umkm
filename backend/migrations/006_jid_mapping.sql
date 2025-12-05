-- JID Mapping table to link @lid and @s.whatsapp.net JIDs
CREATE TABLE IF NOT EXISTS jid_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    lid_jid VARCHAR(100) NOT NULL,  -- @lid format
    phone_jid VARCHAR(100) NOT NULL, -- @s.whatsapp.net format
    phone_number VARCHAR(20) NOT NULL, -- Just the phone number
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, lid_jid),
    UNIQUE(tenant_id, phone_jid)
);

CREATE INDEX IF NOT EXISTS idx_jid_mappings_tenant ON jid_mappings(tenant_id);
CREATE INDEX IF NOT EXISTS idx_jid_mappings_lid ON jid_mappings(lid_jid);
CREATE INDEX IF NOT EXISTS idx_jid_mappings_phone ON jid_mappings(phone_jid);







