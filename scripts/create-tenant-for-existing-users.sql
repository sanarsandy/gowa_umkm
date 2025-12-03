-- Script untuk membuat tenant untuk user yang sudah ada tapi belum punya tenant
-- Run script ini di database jika ada user yang sudah register sebelum auto-create tenant diimplementasi

-- Insert tenant untuk setiap user yang belum punya tenant
INSERT INTO tenants (user_id, business_name, business_type, business_description, business_phone, business_address, is_active, created_at, updated_at)
SELECT 
    u.id as user_id,
    u.full_name || '''s Business' as business_name,
    'UMKM' as business_type,
    '' as business_description,
    '' as business_phone,
    '' as business_address,
    true as is_active,
    NOW() as created_at,
    NOW() as updated_at
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM tenants t WHERE t.user_id = u.id
);

-- Verify hasil
SELECT 
    u.id as user_id,
    u.email,
    u.full_name,
    t.id as tenant_id,
    t.business_name
FROM users u
LEFT JOIN tenants t ON t.user_id = u.id
ORDER BY u.created_at DESC;



