# Setup Static Files - Solusi Permanen

## Nama File Nginx Config
**File:** `/etc/nginx/sites-available/app3.anakhebat.web.id`

## Langkah Setup di VPS

### Step 1: Pull Latest Changes
```bash
cd /var/rumah_afiat/gowa_umkm
git pull origin main
```

### Step 2: Extract Static Files
```bash
# Extract static files dari container ke filesystem
chmod +x scripts/extract-static-simple.sh
./scripts/extract-static-simple.sh
```

Script ini akan:
- Membuat directory `/var/www/gowa_umkm/static/`
- Extract `_nuxt` directory dari container `gowa-app-prod`
- Set permissions yang benar

### Step 3: Verify Static Files
```bash
# Check apakah path sudah ada
ls -la /var/www/gowa_umkm/static/_nuxt/

# Harus menampilkan list files .js
```

### Step 4: Setup Nginx Config
```bash
# Copy config final
sudo cp nginx.conf.production.final /etc/nginx/sites-available/app3.anakhebat.web.id

# Replace placeholder dengan path yang benar
sudo sed -i 's|{{STATIC_ROOT}}|/var/www/gowa_umkm/static|g' /etc/nginx/sites-available/app3.anakhebat.web.id

# Test config
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

### Step 5: Test
```bash
# Test static file access
curl -I http://127.0.0.1/_nuxt/entry.js
# Expected: HTTP/1.1 200 OK (dari filesystem, bukan proxy)
```

## Atau Gunakan Script Otomatis (Recommended)

```bash
# Script ini akan extract files DAN setup nginx sekaligus
chmod +x scripts/setup-production-final.sh
./scripts/setup-production-final.sh
```

## Verifikasi

1. **Check static files:**
   ```bash
   ls -la /var/www/gowa_umkm/static/_nuxt/ | head -10
   ```

2. **Check nginx config:**
   ```bash
   sudo nginx -t
   ```

3. **Check nginx location block:**
   ```bash
   sudo grep -A 5 "location /_nuxt/" /etc/nginx/sites-available/app3.anakhebat.web.id
   ```
   
   Harus menampilkan:
   ```nginx
   location /_nuxt/ {
       alias /var/www/gowa_umkm/static/_nuxt/;
       ...
   }
   ```

4. **Test di browser:**
   - Buka: `https://app3.anakhebat.web.id`
   - Clear cache atau gunakan incognito
   - Check Network tab di DevTools
   - Static files harus return HTTP 200 (bukan 503/502)

## Troubleshooting

### Jika path tidak ada:
```bash
# Buat manual
sudo mkdir -p /var/www/gowa_umkm/static
sudo chown -R www-data:www-data /var/www/gowa_umkm/static

# Extract lagi
./scripts/extract-static-simple.sh
```

### Jika nginx error:
```bash
# Check error
sudo nginx -t

# Check logs
sudo tail -f /var/log/nginx/error.log
```

### Jika masih 503:
```bash
# Check apakah static files ada
ls -la /var/www/gowa_umkm/static/_nuxt/

# Check nginx config
sudo cat /etc/nginx/sites-available/app3.anakhebat.web.id | grep -A 10 "_nuxt"
```

