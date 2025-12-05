# 🚀 Langkah Perbaikan VPS Production - Fix 502/503 Static Files

Panduan lengkap untuk melakukan perbaikan di VPS production.

## 📋 Prasyarat

- Akses SSH ke VPS
- User memiliki permission untuk docker dan nginx
- Git repository sudah di-clone di VPS

## 🎯 Opsi 1: Menggunakan Script Otomatis (Recommended)

### Step 1: SSH ke VPS

```bash
ssh user@your-vps-ip
```

### Step 2: Navigate ke Project Directory

```bash
cd /path/to/gowa_umkm
# atau
cd ~/gowa_umkm
```

### Step 3: Pull Latest Changes

```bash
git pull origin main
```

### Step 4: Berikan Permission Execute

```bash
chmod +x scripts/fix-vps-production.sh
```

### Step 5: Jalankan Script

```bash
./scripts/fix-vps-production.sh
```

Script akan otomatis:
- ✅ Backup database dan nginx config
- ✅ Pull latest changes
- ✅ Rebuild container
- ✅ Update nginx config
- ✅ Reload services
- ✅ Verifikasi fixes

---

## 🔧 Opsi 2: Manual Step-by-Step

Jika ingin melakukan secara manual atau script tidak berjalan:

### Step 1: Backup

```bash
# Backup database
docker compose -f docker-compose.prod.yml exec -T db pg_dump -U gowa_user gowa_db > backups/backup_$(date +%Y%m%d_%H%M%S).sql

# Backup nginx config
sudo cp /etc/nginx/sites-available/app3.anakhebat.web.id /etc/nginx/sites-available/app3.anakhebat.web.id.backup
```

### Step 2: Pull Latest Changes

```bash
cd /path/to/gowa_umkm
git pull origin main
```

### Step 3: Stop Container

```bash
docker compose -f docker-compose.prod.yml stop app
```

### Step 4: Rebuild Container

```bash
# Rebuild dengan no-cache (PENTING!)
docker compose -f docker-compose.prod.yml build --no-cache app
```

**Note:** Proses ini mungkin memakan waktu 5-10 menit.

### Step 5: Start Container

```bash
docker compose -f docker-compose.prod.yml up -d app
```

### Step 6: Wait for Health Check

```bash
# Tunggu sekitar 40-60 detik untuk health check
sleep 45

# Check health status
docker inspect --format='{{.State.Health.Status}}' gowa-app-prod
# Expected: "healthy"
```

### Step 7: Verify Static Files

```bash
# Check if static files exist
docker exec gowa-app-prod ls -la /app/.output/public/_nuxt/

# Count files
docker exec gowa-app-prod find /app/.output/public/_nuxt -type f | wc -l
# Expected: Should return number > 0
```

### Step 8: Update Nginx Config

```bash
# Copy new config
sudo cp nginx.conf.production /etc/nginx/sites-available/app3.anakhebat.web.id

# Test config (PENTING!)
sudo nginx -t
# Expected: "syntax is ok" dan "test is successful"

# If test fails, restore backup:
# sudo cp /etc/nginx/sites-available/app3.anakhebat.web.id.backup /etc/nginx/sites-available/app3.anakhebat.web.id
```

### Step 9: Reload Nginx

```bash
sudo systemctl reload nginx
```

### Step 10: Verify

```bash
# Test port accessibility
curl -I http://127.0.0.1:3002
# Expected: HTTP/1.1 200 OK

# Test static file access
curl -I http://127.0.0.1:3002/_nuxt/entry.js
# Expected: HTTP/1.1 200 OK (atau 404 jika file tidak ada, tapi bukan 502/503)

# Check container logs
docker logs gowa-app-prod | tail -20
```

---

## ✅ Verifikasi Setelah Perbaikan

### 1. Check Container Status

```bash
docker ps | grep gowa-app-prod
# Expected: Container harus running
```

### 2. Check Container Health

```bash
docker inspect --format='{{.State.Health.Status}}' gowa-app-prod
# Expected: "healthy"
```

### 3. Test Port Access

```bash
curl -I http://127.0.0.1:3002
# Expected: HTTP/1.1 200 OK
```

### 4. Test Static Files

```bash
# Get first static file
FIRST_FILE=$(docker exec gowa-app-prod find /app/.output/public/_nuxt -type f -name "*.js" | head -1)
FILE_NAME=$(basename "$FIRST_FILE")

# Test access
curl -I "http://127.0.0.1:3002/_nuxt/$FILE_NAME"
# Expected: HTTP/1.1 200 OK (bukan 502/503)
```

### 5. Browser Test

1. Buka browser
2. Buka DevTools (F12)
3. Buka tab Network
4. Akses: `https://app3.anakhebat.web.id`
5. Check:
   - Request ke `/_nuxt/*.js` harus return 200 OK
   - Tidak ada 502/503 error
   - Console tidak ada error terkait static files

---

## 🔍 Troubleshooting

### Problem: Container tidak healthy

```bash
# Check logs
docker logs gowa-app-prod

# Check if container is running
docker ps -a | grep gowa-app-prod

# Restart container
docker compose -f docker-compose.prod.yml restart app

# Wait and check again
sleep 45
docker inspect --format='{{.State.Health.Status}}' gowa-app-prod
```

### Problem: Static files tidak ada

```bash
# Check if build was successful
docker logs gowa-app-prod | grep -i "build\|error"

# Rebuild again
docker compose -f docker-compose.prod.yml build --no-cache app
docker compose -f docker-compose.prod.yml up -d app

# Verify
docker exec gowa-app-prod ls -la /app/.output/public/_nuxt/
```

### Problem: Nginx config test failed

```bash
# Check error
sudo nginx -t

# Restore backup
sudo cp /etc/nginx/sites-available/app3.anakhebat.web.id.backup /etc/nginx/sites-available/app3.anakhebat.web.id

# Test again
sudo nginx -t
```

### Problem: Masih error 502/503

```bash
# Run diagnostic script
./scripts/diagnose-static-files.sh

# Check nginx error log
sudo tail -f /var/log/nginx/app3_error.log

# Check container logs
docker logs -f gowa-app-prod
```

### Problem: Port tidak accessible

```bash
# Check if port is in use
sudo netstat -tulpn | grep 3002

# Check container port mapping
docker ps | grep gowa-app-prod

# Check docker-compose config
cat docker-compose.prod.yml | grep -A 5 "ports:"
```

---

## 📊 Monitoring Setelah Perbaikan

### Check Logs

```bash
# Container logs
docker logs -f gowa-app-prod

# Nginx access logs
sudo tail -f /var/log/nginx/app3_access.log

# Nginx error logs
sudo tail -f /var/log/nginx/app3_error.log
```

### Check Performance

```bash
# Container stats
docker stats gowa-app-prod

# Network connections
sudo netstat -an | grep 3002 | wc -l
```

---

## ✅ Checklist Final

Setelah semua langkah, pastikan:

- [ ] Container running: `docker ps | grep gowa-app-prod`
- [ ] Container healthy: `docker inspect --format='{{.State.Health.Status}}' gowa-app-prod` → "healthy"
- [ ] Port accessible: `curl -I http://127.0.0.1:3002` → 200 OK
- [ ] Static files exist: `docker exec gowa-app-prod ls /app/.output/public/_nuxt/` → ada files
- [ ] Static files accessible: `curl -I http://127.0.0.1:3002/_nuxt/entry.js` → 200 OK (bukan 502/503)
- [ ] Nginx config valid: `sudo nginx -t` → "syntax is ok"
- [ ] Browser test: Tidak ada 502/503 error di browser console

---

## 🎯 Quick Reference

### Commands yang Sering Digunakan

```bash
# Rebuild dan restart
docker compose -f docker-compose.prod.yml build --no-cache app
docker compose -f docker-compose.prod.yml up -d app

# Check status
docker ps | grep gowa-app-prod
docker inspect --format='{{.State.Health.Status}}' gowa-app-prod

# View logs
docker logs -f gowa-app-prod
sudo tail -f /var/log/nginx/app3_error.log

# Test access
curl -I http://127.0.0.1:3002
curl -I http://127.0.0.1:3002/_nuxt/entry.js

# Reload nginx
sudo nginx -t && sudo systemctl reload nginx
```

---

**Last Updated**: 2025-01-XX
**Status**: ✅ Ready for Production

