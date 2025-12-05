-- Message Templates for quick replies
CREATE TABLE IF NOT EXISTS message_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) DEFAULT 'general', -- general, greeting, promo, follow_up, closing
    content TEXT NOT NULL,
    variables JSONB DEFAULT '[]', -- e.g., ["nama", "produk", "harga"]
    is_active BOOLEAN DEFAULT true,
    usage_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_message_templates_tenant ON message_templates(tenant_id);
CREATE INDEX IF NOT EXISTS idx_message_templates_category ON message_templates(tenant_id, category);

-- Broadcast Messages
CREATE TABLE IF NOT EXISTS broadcasts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(200) NOT NULL,
    message_content TEXT NOT NULL,
    template_id UUID REFERENCES message_templates(id),
    status VARCHAR(20) DEFAULT 'draft', -- draft, scheduled, sending, completed, cancelled
    scheduled_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    total_recipients INTEGER DEFAULT 0,
    sent_count INTEGER DEFAULT 0,
    delivered_count INTEGER DEFAULT 0,
    failed_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_broadcasts_tenant ON broadcasts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_broadcasts_status ON broadcasts(tenant_id, status);

-- Broadcast Recipients
CREATE TABLE IF NOT EXISTS broadcast_recipients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    broadcast_id UUID NOT NULL REFERENCES broadcasts(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES customer_insights(id),
    customer_jid VARCHAR(100) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, sent, delivered, failed
    message_id VARCHAR(100), -- WhatsApp message ID
    sent_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_broadcast_recipients_broadcast ON broadcast_recipients(broadcast_id);
CREATE INDEX IF NOT EXISTS idx_broadcast_recipients_status ON broadcast_recipients(broadcast_id, status);







