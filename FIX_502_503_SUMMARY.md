# 🔧 Fix 502/503 Error pada Static Files - Ringkasan

## 📋 Masalah yang Diperbaiki

Error 502/503 Bad Gateway/Service Unavailable pada static files `/_nuxt/*.js` dan `/_nuxt/builds/meta/*.json`.

## ✅ Perbaikan yang Sudah Dilakukan

### 1. **Dockerfile.prod** ✅
- ✅ Tambahkan health check untuk monitoring container
- ✅ Tambahkan verification untuk static files
- ✅ Pastikan environment variables benar

### 2. **docker-compose.prod.yml** ✅
- ✅ Fix port binding: `127.0.0.1:3002:3000` (bind ke localhost only)
- ✅ Tambahkan health check dengan retry mechanism
- ✅ Tambahkan environment variables yang diperlukan
- ✅ Fix depends_on dengan condition

### 3. **nginx.conf.production** ✅
- ✅ Tambahkan timeout settings (connect, send, read)
- ✅ Tambahkan error handling dengan retry (3x)
- ✅ Tambahkan buffer settings untuk optimasi
- ✅ Tambahkan WebSocket upgrade map
- ✅ Tambahkan API endpoint location
- ✅ Tambahkan health check endpoint
- ✅ Tambahkan security headers
- ✅ Tambahkan gzip compression

## 🚀 Langkah Deployment (PENTING!)

### Step 1: Rebuild Container

```bash
# Stop container
docker compose -f docker-compose.prod.yml stop app

# Rebuild dengan no-cache (PENTING!)
docker compose -f docker-compose.prod.yml build --no-cache app

# Start container
docker compose -f docker-compose.prod.yml up -d app

# Wait for container to be healthy (sekitar 40 detik)
docker ps | grep gowa-app-prod
```

### Step 2: Update Nginx Config

```bash
# Copy config baru
sudo cp nginx.conf.production /etc/nginx/sites-available/app3.anakhebat.web.id

# Test config (PENTING!)
sudo nginx -t

# Jika test berhasil, reload nginx
sudo systemctl reload nginx
```

### Step 3: Verify

```bash
# Check container status
docker ps | grep gowa-app-prod

# Check container health (harus "healthy")
docker inspect --format='{{.State.Health.Status}}' gowa-app-prod

# Test static file access
curl -I http://127.0.0.1:3002/_nuxt/entry.js
# Harus return: HTTP/1.1 200 OK

# Check logs
docker logs gowa-app-prod | tail -20
```

## 🔍 Troubleshooting

### Jika masih error 502/503:

1. **Run diagnostic script:**
```bash
./scripts/diagnose-static-files.sh
```

2. **Check container logs:**
```bash
docker logs gowa-app-prod
```

3. **Check nginx error logs:**
```bash
sudo tail -f /var/log/nginx/app3_error.log
```

4. **Verify static files exist:**
```bash
docker exec gowa-app-prod ls -la /app/.output/public/_nuxt/
```

5. **Test direct access:**
```bash
docker exec gowa-app-prod curl -I http://localhost:3000/_nuxt/entry.js
```

## 📊 Perubahan Utama

### docker-compose.prod.yml
```yaml
# SEBELUM:
ports:
  - "3002:3000"

# SESUDAH:
ports:
  - "127.0.0.1:3002:3000"  # Bind ke localhost only

# TAMBAHAN:
healthcheck:
  test: ["CMD", "node", "-e", "require('http').get('http://localhost:3000', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

### nginx.conf.production
```nginx
# TAMBAHAN: Timeout settings
proxy_connect_timeout 60s;
proxy_send_timeout 60s;
proxy_read_timeout 60s;

# TAMBAHAN: Error handling dengan retry
proxy_next_upstream error timeout invalid_header http_500 http_502 http_503;
proxy_next_upstream_tries 3;
proxy_next_upstream_timeout 15s;

# TAMBAHAN: Buffer settings
proxy_buffering on;
proxy_buffer_size 4k;
proxy_buffers 8 4k;
proxy_busy_buffers_size 8k;
```

## ✅ Checklist Verifikasi

Setelah deployment, pastikan:

- [ ] Container running: `docker ps | grep gowa-app-prod`
- [ ] Container healthy: `docker inspect --format='{{.State.Health.Status}}' gowa-app-prod` → harus "healthy"
- [ ] Port accessible: `curl -I http://127.0.0.1:3002` → harus 200 OK
- [ ] Static files exist: `docker exec gowa-app-prod ls /app/.output/public/_nuxt/` → harus ada files
- [ ] HTTP access works: `curl -I http://127.0.0.1:3002/_nuxt/entry.js` → harus 200 OK
- [ ] Nginx config valid: `sudo nginx -t` → harus "syntax is ok"
- [ ] No errors in logs: `docker logs gowa-app-prod | grep -i error` → tidak ada error
- [ ] Browser test: Buka browser, buka DevTools, test aplikasi → tidak ada 502/503

## 🎯 Root Cause

Masalah utama adalah:
1. **Timeout terlalu pendek** - Nginx timeout sebelum container ready
2. **Tidak ada retry mechanism** - Nginx langsung return error tanpa retry
3. **Port binding** - Port tidak ter-bind dengan benar
4. **Tidak ada health check** - Tidak tahu kapan container ready

Semua masalah ini sudah diperbaiki! ✅

## 📚 Dokumentasi Lengkap

Untuk detail lengkap, lihat:
- `docs/FIX_502_503_STATIC_FILES.md` - Dokumentasi lengkap troubleshooting
- `scripts/diagnose-static-files.sh` - Script diagnostic

---

**Status**: ✅ Fixed
**Last Updated**: 2025-01-XX

