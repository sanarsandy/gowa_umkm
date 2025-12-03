-- Migration 015: Fix Customer Tags and Notes Tables
-- Creates missing tables with correct foreign keys to customer_insights

-- Customer Tags table (must be created first)
CREATE TABLE IF NOT EXISTS customer_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    color VARCHAR(20) DEFAULT '#6366f1',
    description TEXT,
    auto_apply_rules JSONB, -- Rules for auto-applying this tag
    customer_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

-- Customer-Tag relationship (many-to-many)
CREATE TABLE IF NOT EXISTS customer_tag_assignments (
    customer_id UUID NOT NULL REFERENCES customer_insights(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES customer_tags(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT NOW(),
    assigned_by VARCHAR(50) DEFAULT 'manual', -- 'manual', 'auto', 'ai'
    PRIMARY KEY (customer_id, tag_id)
);

-- Customer notes
CREATE TABLE IF NOT EXISTS customer_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customer_insights(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_by UUID, -- user who created the note
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_customer_tags_tenant ON customer_tags(tenant_id);
CREATE INDEX IF NOT EXISTS idx_customer_tag_assignments_customer ON customer_tag_assignments(customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_tag_assignments_tag ON customer_tag_assignments(tag_id);
CREATE INDEX IF NOT EXISTS idx_customer_notes_customer ON customer_notes(customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_notes_tenant ON customer_notes(tenant_id);

COMMENT ON TABLE customer_tags IS 'Custom tags for customer segmentation';
COMMENT ON TABLE customer_tag_assignments IS 'Many-to-many relationship between customers and tags';
COMMENT ON TABLE customer_notes IS 'Notes/comments about customers';
