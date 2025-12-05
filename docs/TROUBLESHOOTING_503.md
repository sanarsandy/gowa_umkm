# 🔧 Troubleshooting 503 Service Unavailable

## 🔴 Gejala

Error di browser:
```
GET https://app3.anakhebat.web.id/_nuxt/xxx.js net::ERR_ABORTED 503 (Service Unavailable)
```

## 🔍 Penyebab Umum

1. **Container tidak running atau crash**
2. **Container tidak healthy**
3. **Port tidak accessible**
4. **Static files tidak lengkap**
5. **Nginx tidak bisa connect ke container**

## 🚀 Langkah Troubleshooting

### Step 1: Run Diagnostic Script

```bash
./scripts/diagnose-503-error.sh
```

Script ini akan check:
- Container status
- Container health
- Port accessibility
- Static files
- Container logs
- Nginx logs

### Step 2: Check Container Status

```bash
# Check if container is running
docker ps | grep gowa-app-prod

# Check container health
docker inspect --format='{{.State.Health.Status}}' gowa-app-prod
# Expected: "healthy"
```

### Step 3: Check Container Logs

```bash
# View recent logs
docker logs --tail 50 gowa-app-prod

# Follow logs in real-time
docker logs -f gowa-app-prod
```

Look for:
- ENOENT errors
- Port binding errors
- Application crashes
- Health check failures

### Step 4: Test Container Directly

```bash
# Test from inside container
docker exec gowa-app-prod curl -I http://localhost:3000

# Test from host
curl -I http://127.0.0.1:3002
```

### Step 5: Check Static Files

```bash
# Check if files exist
docker exec gowa-app-prod ls -la /app/.output/server/chunks/public/_nuxt/ | head -10

# Count files
docker exec gowa-app-prod find /app/.output/server/chunks/public/_nuxt -type f | wc -l
# Expected: ~46 files
```

## 🔧 Solusi

### Solution 1: Restart Container

```bash
docker compose -f docker-compose.prod.yml restart app
sleep 45
docker inspect --format='{{.State.Health.Status}}' gowa-app-prod
```

### Solution 2: Fix Static Files

```bash
./scripts/fix-container-static-files.sh
```

### Solution 3: Rebuild Container

```bash
docker compose -f docker-compose.prod.yml stop app
docker compose -f docker-compose.prod.yml build --no-cache app
docker compose -f docker-compose.prod.yml up -d app
sleep 45
```

### Solution 4: Check Nginx Config

```bash
# Test nginx config
sudo nginx -t

# Check nginx error logs
sudo tail -f /var/log/nginx/app3_error.log

# Reload nginx
sudo systemctl reload nginx
```

## 📊 Common Issues & Solutions

### Issue: Container not healthy

**Symptoms:**
- Health check returns "unhealthy" or "starting"
- Container keeps restarting

**Solution:**
```bash
# Check logs
docker logs gowa-app-prod

# Check if port 3000 is accessible from inside
docker exec gowa-app-prod curl http://localhost:3000

# Restart container
docker compose -f docker-compose.prod.yml restart app
```

### Issue: Port not accessible

**Symptoms:**
- `curl http://127.0.0.1:3002` returns connection refused
- Port 3002 not listening

**Solution:**
```bash
# Check port mapping
docker ps | grep gowa-app-prod

# Check docker-compose config
cat docker-compose.prod.yml | grep -A 2 "ports:"

# Restart container
docker compose -f docker-compose.prod.yml restart app
```

### Issue: Static files incomplete

**Symptoms:**
- Only 2 files instead of 46
- ENOENT errors in logs

**Solution:**
```bash
# Run fix script
./scripts/fix-container-static-files.sh

# Or rebuild
docker compose -f docker-compose.prod.yml build --no-cache app
docker compose -f docker-compose.prod.yml up -d app
```

### Issue: Nginx can't connect

**Symptoms:**
- 502/503 errors in nginx logs
- "upstream connection failed" errors

**Solution:**
```bash
# Check if container is accessible
curl -I http://127.0.0.1:3002

# Check nginx config
sudo nginx -t

# Check nginx error logs
sudo tail -f /var/log/nginx/app3_error.log

# Reload nginx
sudo systemctl reload nginx
```

## ✅ Verification Checklist

After fixing, verify:

- [ ] Container running: `docker ps | grep gowa-app-prod`
- [ ] Container healthy: `docker inspect --format='{{.State.Health.Status}}' gowa-app-prod` → "healthy"
- [ ] Port accessible: `curl -I http://127.0.0.1:3002` → 200 OK
- [ ] Static files exist: `docker exec gowa-app-prod find /app/.output/server/chunks/public/_nuxt -type f | wc -l` → ~46
- [ ] No errors in logs: `docker logs gowa-app-prod | grep -i error | tail -5` → minimal errors
- [ ] Browser test: No 503 errors in browser console

## 🎯 Quick Fix Commands

```bash
# Full diagnostic
./scripts/diagnose-503-error.sh

# Quick fix static files
./scripts/fix-container-static-files.sh

# Restart everything
docker compose -f docker-compose.prod.yml restart app && sleep 45

# Check status
docker ps | grep gowa-app-prod && curl -I http://127.0.0.1:3002
```

---

**Last Updated**: 2025-01-XX

