# Fix 502/503 Error pada Static Files (_nuxt/*)

Dokumen ini menjelaskan masalah 502/503 Bad Gateway/Service Unavailable pada static files dan solusinya.

## 🔴 Gejala Masalah

Error yang muncul di browser console:
```
GET https://app3.anakhebat.web.id/_nuxt/xxx.js 502 (Bad Gateway)
GET https://app3.anakhebat.web.id/_nuxt/xxx.js 503 (Service Unavailable)
```

## 🔍 Root Cause Analysis

### Masalah Utama

1. **Container tidak ready** - Container app belum fully started saat nginx mencoba connect
2. **Port binding issue** - Port tidak ter-bind dengan benar
3. **Timeout terlalu pendek** - Nginx timeout sebelum container merespons
4. **Static files tidak ada** - Build process tidak menghasilkan static files
5. **Container crash** - Container crash setelah start

### Analisis Detail

#### 1. Dockerfile.prod
- ✅ Sudah benar: Copy `.output` directory yang sudah include static files
- ✅ Health check sudah ditambahkan
- ✅ Working directory sudah benar

#### 2. docker-compose.prod.yml
- ✅ Port binding: `127.0.0.1:3002:3000` (sudah diperbaiki)
- ✅ Health check sudah ditambahkan
- ✅ Environment variables sudah benar

#### 3. nginx.conf.production
- ✅ Timeout settings sudah ditambahkan
- ✅ Error handling sudah ditambahkan
- ✅ Buffer settings sudah dikonfigurasi

## ✅ Solusi yang Sudah Diimplementasikan

### 1. Fix Dockerfile.prod

**Perubahan:**
- Tambahkan health check
- Tambahkan verification untuk static files
- Pastikan working directory benar

```dockerfile
# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
  CMD node -e "require('http').get('http://localhost:3000', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"
```

### 2. Fix docker-compose.prod.yml

**Perubahan:**
- Port binding: `127.0.0.1:3002:3000` (bind ke localhost only)
- Tambahkan health check
- Tambahkan environment variables

```yaml
ports:
  - "127.0.0.1:3002:3000"
healthcheck:
  test: ["CMD", "node", "-e", "require('http').get('http://localhost:3000', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

### 3. Fix nginx.conf.production

**Perubahan:**
- Tambahkan timeout settings
- Tambahkan error handling dengan retry
- Tambahkan buffer settings
- Tambahkan WebSocket upgrade map

```nginx
# Timeouts - CRITICAL untuk mencegah 502/503
proxy_connect_timeout 60s;
proxy_send_timeout 60s;
proxy_read_timeout 60s;

# Error handling dengan retry
proxy_next_upstream error timeout invalid_header http_500 http_502 http_503;
proxy_next_upstream_tries 3;
proxy_next_upstream_timeout 15s;
```

## 🚀 Langkah Deployment

### 1. Rebuild Container

```bash
# Stop container
docker compose -f docker-compose.prod.yml stop app

# Rebuild dengan no-cache
docker compose -f docker-compose.prod.yml build --no-cache app

# Start container
docker compose -f docker-compose.prod.yml up -d app
```

### 2. Update Nginx Config

```bash
# Copy config baru
sudo cp nginx.conf.production /etc/nginx/sites-available/app3.anakhebat.web.id

# Test config
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

### 3. Verify

```bash
# Check container status
docker ps | grep gowa-app-prod

# Check container health
docker inspect --format='{{.State.Health.Status}}' gowa-app-prod

# Test static file access
curl -I http://127.0.0.1:3002/_nuxt/entry.js

# Check logs
docker logs gowa-app-prod
```

## 🔧 Troubleshooting

### Problem: Container tidak start

**Solution:**
```bash
# Check logs
docker logs gowa-app-prod

# Check if port is in use
sudo netstat -tulpn | grep 3002

# Restart container
docker compose -f docker-compose.prod.yml restart app
```

### Problem: Static files masih 502/503

**Solution:**
```bash
# Run diagnostic script
./scripts/diagnose-static-files.sh

# Check if static files exist in container
docker exec gowa-app-prod ls -la /app/.output/public/_nuxt/

# Test direct access
docker exec gowa-app-prod curl -I http://localhost:3000/_nuxt/entry.js
```

### Problem: Container crash setelah start

**Solution:**
```bash
# Check logs
docker logs gowa-app-prod

# Check memory
docker stats gowa-app-prod

# Increase memory limit in docker-compose.prod.yml if needed
```

### Problem: Nginx timeout

**Solution:**
```bash
# Check nginx error log
sudo tail -f /var/log/nginx/app3_error.log

# Increase timeout in nginx.conf.production
proxy_connect_timeout 120s;
proxy_send_timeout 120s;
proxy_read_timeout 120s;
```

## 📊 Monitoring

### Health Check

```bash
# Check container health
docker inspect --format='{{.State.Health.Status}}' gowa-app-prod

# Should return: healthy
```

### Logs Monitoring

```bash
# Container logs
docker logs -f gowa-app-prod

# Nginx access logs
sudo tail -f /var/log/nginx/app3_access.log

# Nginx error logs
sudo tail -f /var/log/nginx/app3_error.log
```

### Performance Monitoring

```bash
# Container stats
docker stats gowa-app-prod

# Network connections
sudo netstat -an | grep 3002

# Response time test
time curl -o /dev/null -s http://127.0.0.1:3002/_nuxt/entry.js
```

## ✅ Checklist Verifikasi

- [ ] Container running: `docker ps | grep gowa-app-prod`
- [ ] Container healthy: `docker inspect --format='{{.State.Health.Status}}' gowa-app-prod`
- [ ] Port accessible: `curl -I http://127.0.0.1:3002`
- [ ] Static files exist: `docker exec gowa-app-prod ls /app/.output/public/_nuxt/`
- [ ] HTTP access works: `curl -I http://127.0.0.1:3002/_nuxt/entry.js`
- [ ] Nginx config valid: `sudo nginx -t`
- [ ] No errors in logs: `docker logs gowa-app-prod | grep -i error`
- [ ] Browser can load: Test di browser dengan DevTools open

## 🎯 Best Practices

1. **Always wait for health check** - Jangan langsung test setelah start container
2. **Monitor logs** - Check logs setelah setiap deployment
3. **Use health checks** - Enable health check di docker-compose
4. **Proper timeouts** - Set timeout yang cukup di nginx
5. **Error handling** - Enable retry mechanism di nginx
6. **Port binding** - Bind ke 127.0.0.1 untuk security

## 📚 Referensi

- [Nginx Proxy Timeout](https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_read_timeout)
- [Docker Health Checks](https://docs.docker.com/engine/reference/builder/#healthcheck)
- [Nuxt Static Files](https://nuxt.com/docs/getting-started/deployment#static-hosting)

---

**Last Updated**: 2025-01-XX
**Status**: ✅ Fixed

