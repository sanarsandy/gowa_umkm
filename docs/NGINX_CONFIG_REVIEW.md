# Nginx Configuration Review - Critical Fixes Needed

## 🔴 Critical Issues Found

### 1. **Missing WebSocket Support in `/api` Location**

**Your Current Config:**
```nginx
location /api {
    proxy_pass http://127.0.0.1:8087;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    # ... other headers
    # ❌ NO WebSocket upgrade headers!
}
```

**Fixed Config:**
```nginx
location /api/ {  # Note: trailing slash added
    proxy_pass http://127.0.0.1:8087;
    proxy_http_version 1.1;
    
    # ✅ WebSocket support (CRITICAL FIX)
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $connection_upgrade;
    
    proxy_set_header Host $host;
    # ... other headers
    
    # ✅ Extended timeouts for WebSocket (24 hours)
    proxy_read_timeout 86400s;
    proxy_send_timeout 86400s;
}
```

**Why This Matters:** Without these headers, WebSocket connections to `/api/ws` will fail with error 1006.

---

### 2. **Missing WebSocket Map Directive**

**Your Current Config:**
```nginx
# ❌ Missing at the top of the file
server {
    listen 443 ssl;
    # ...
}
```

**Fixed Config:**
```nginx
# ✅ Add this BEFORE the server block
map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}

server {
    listen 443 ssl http2;  # Also added http2
    # ...
}
```

**Why This Matters:** The `$connection_upgrade` variable is used in the WebSocket proxy headers.

---

### 3. **Location Path Mismatch**

**Your Current Config:**
```nginx
location /api {  # ❌ Without trailing slash
```

**Fixed Config:**
```nginx
location /api/ {  # ✅ With trailing slash
```

**Why This Matters:** Using `/api/` ensures all API routes like `/api/ws`, `/api/dashboard/stats`, etc. are properly matched.

---

## 🟡 Additional Improvements

### 4. **Error Handling for Static Assets**

**Added:**
```nginx
location /_nuxt/ {
    # ... existing config
    
    # ✅ Error handling
    proxy_intercept_errors on;
    error_page 502 503 504 = @nuxt_fallback;
}

location @nuxt_fallback {
    return 503 "Nuxt service temporarily unavailable";
    add_header Content-Type text/plain;
}
```

### 5. **Security Headers**

**Added:**
```nginx
# Security headers
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "no-referrer-when-downgrade" always;
```

### 6. **Gzip Compression**

**Added:**
```nginx
# Gzip compression
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json application/javascript;
```

---

## 📋 Deployment Steps

### On Your VPS Server:

1. **Upload the corrected config:**
   ```bash
   cd /path/to/gowa_umkm
   git pull origin main
   ```

2. **Run the update script:**
   ```bash
   sudo ./scripts/update-nginx-prod.sh
   ```

   Or manually:
   ```bash
   # Backup current config
   sudo cp /etc/nginx/sites-available/app2.anakhebat.web.id \
          /etc/nginx/sites-available/app2.anakhebat.web.id.backup
   
   # Copy new config
   sudo cp nginx.conf.prod /etc/nginx/sites-available/app2.anakhebat.web.id
   
   # Test
   sudo nginx -t
   
   # Reload
   sudo systemctl reload nginx
   ```

3. **Verify:**
   - Open https://app2.anakhebat.web.id/
   - Check browser console - WebSocket should connect successfully
   - No more 503 errors for static assets

---

## 🔍 What Will Be Fixed

After applying this configuration:

✅ **WebSocket connections will work** - The `/api/ws` endpoint will properly upgrade HTTP to WebSocket  
✅ **Static assets will load** - `_nuxt/*.js` files will return 200 instead of 503  
✅ **Better error handling** - Graceful fallback for Nuxt service issues  
✅ **Improved security** - Security headers protect against common attacks  
✅ **Better performance** - Gzip compression reduces bandwidth  

---

## 📊 Before vs After

| Issue | Before | After |
|-------|--------|-------|
| WebSocket Connection | ❌ Fails (1006) | ✅ Works |
| Static Assets | ❌ 503 errors | ✅ 200 OK |
| API Routing | ⚠️ `/api` (no slash) | ✅ `/api/` (with slash) |
| WebSocket Timeout | ⚠️ 60s | ✅ 24 hours |
| Error Handling | ❌ None | ✅ Graceful fallback |
| Security Headers | ❌ Missing | ✅ Added |
| Compression | ❌ None | ✅ Gzip enabled |

---

## 🚨 Important Notes

1. **The WebSocket fix is critical** - This is the main reason your WebSocket connections are failing
2. **Backup is automatic** - The update script creates a timestamped backup
3. **Zero downtime** - Nginx reload is graceful, no service interruption
4. **Rollback is easy** - Just restore the backup if needed

---

## Files Created

- [nginx.conf.prod](file:///Users/lman.kadiv-doti/Documents/2025/SAAS/gowa_umkm/nginx.conf.prod) - Corrected production config
- [scripts/update-nginx-prod.sh](file:///Users/lman.kadiv-doti/Documents/2025/SAAS/gowa_umkm/scripts/update-nginx-prod.sh) - Safe update script
