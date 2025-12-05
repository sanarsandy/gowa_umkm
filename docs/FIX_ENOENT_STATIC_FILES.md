# 🔧 Fix ENOENT Error - Static Files Path Issue

## 🔴 Masalah

Error yang muncul:
```
ENOENT: no such file or directory, open '/app/.output/server/chunks/public/_nuxt/DeXWKxqL.js'
```

## 🔍 Root Cause

**Nitro dengan preset `node-server` mencari static files di path yang berbeda:**

- **Path yang dicari Nitro:** `/app/.output/server/chunks/public/_nuxt/`
- **Path file sebenarnya:** `/app/.output/public/_nuxt/`

Ini adalah masalah konfigurasi path di Nitro preset `node-server`.

## ✅ Solusi

### Fix Dockerfile.prod

Tambahkan step untuk copy atau symlink files ke lokasi yang dicari Nitro:

```dockerfile
# CRITICAL FIX: Nitro mencari static files di server/chunks/public/
# tapi file ada di public/, jadi kita copy atau symlink
RUN if [ -d ".output/public" ] && [ -d ".output/server/chunks" ]; then \
        mkdir -p .output/server/chunks/public && \
        if [ -d ".output/public/_nuxt" ]; then \
            cp -r .output/public/_nuxt .output/server/chunks/public/ || \
            ln -sf ../../../../public/_nuxt .output/server/chunks/public/_nuxt; \
        fi && \
        if [ -d ".output/public/builds" ]; then \
            cp -r .output/public/builds .output/server/chunks/public/ || \
            ln -sf ../../../../public/builds .output/server/chunks/public/builds; \
        fi; \
    fi
```

## 🚀 Deployment Steps

### Step 1: Update Dockerfile

File sudah diperbaiki di `frontend/Dockerfile.prod`.

### Step 2: Rebuild Container

```bash
# Stop container
docker compose -f docker-compose.prod.yml stop app

# Rebuild dengan no-cache (PENTING!)
docker compose -f docker-compose.prod.yml build --no-cache app

# Start container
docker compose -f docker-compose.prod.yml up -d app
```

### Step 3: Verify Fix

```bash
# Check if files exist in both locations
docker exec gowa-app-prod ls -la /app/.output/public/_nuxt/ | head -5
docker exec gowa-app-prod ls -la /app/.output/server/chunks/public/_nuxt/ | head -5

# Both should show files
```

### Step 4: Test

```bash
# Test static file access
curl -I http://127.0.0.1:3002/_nuxt/entry.js
# Expected: HTTP/1.1 200 OK (bukan 500 ENOENT)
```

## 🔍 Diagnostic Commands

### Check Container Structure

```bash
# Run diagnostic script
./scripts/check-container-structure.sh
```

### Manual Check

```bash
# Check if public/_nuxt exists
docker exec gowa-app-prod ls -la /app/.output/public/_nuxt/

# Check if server/chunks/public/_nuxt exists
docker exec gowa-app-prod ls -la /app/.output/server/chunks/public/_nuxt/

# Check if symlink exists
docker exec gowa-app-prod ls -la /app/.output/server/chunks/public/
```

## 📊 Expected Structure After Fix

```
/app/.output/
├── public/
│   └── _nuxt/
│       ├── entry.js
│       ├── *.js files
│       └── builds/
│           └── meta/
└── server/
    └── chunks/
        └── public/
            └── _nuxt/  (symlink atau copy dari public/_nuxt)
                ├── entry.js
                ├── *.js files
                └── builds/
                    └── meta/
```

## ✅ Verification Checklist

- [ ] Files exist in `/app/.output/public/_nuxt/`
- [ ] Files exist in `/app/.output/server/chunks/public/_nuxt/`
- [ ] No ENOENT errors in logs
- [ ] Static files return 200 OK (bukan 500)
- [ ] Browser can load static files

## 🎯 Alternative Solutions

### Option 1: Copy Files (Current Solution)
- ✅ Reliable
- ✅ Works in all cases
- ❌ Uses more disk space

### Option 2: Symlink
- ✅ Saves disk space
- ✅ Always in sync
- ❌ May not work in all environments

### Option 3: Change Nitro Config
- ✅ Clean solution
- ❌ Requires Nuxt config changes
- ❌ May not be supported

**Current implementation uses copy with symlink fallback for best compatibility.**

---

**Last Updated**: 2025-01-XX
**Status**: ✅ Fixed

