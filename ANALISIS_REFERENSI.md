# Analisis Referensi Config yang Berhasil

## Perbedaan Kunci antara Referensi (Berhasil) vs Config Kita (Error)

### 1. **Nginx Config - Static Files Handling**

#### Referensi (BERHASIL):
```nginx
location /seniman/_nuxt/ {
    proxy_pass http://10.100.246.202:3068/seniman/_nuxt/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    
    expires 1y;
    add_header Cache-Control "public, immutable";
    add_header Vary "Accept-Encoding";
    
    gzip_vary on;  # ← PENTING: gzip_vary, bukan gzip_static
}
```

**Key Points:**
- ✅ **Proxy ke container** (bukan direct file serving)
- ✅ **Tidak ada `try_files`** atau fallback
- ✅ **`gzip_vary on`** (bukan `gzip_static on`)
- ✅ **Buffer settings di location utama**, bukan di static location

#### Config Kita (ERROR):
```nginx
location /_nuxt/ {
    alias {{STATIC_ROOT}}/_nuxt/;  # ← Direct file serving
    try_files $uri @nuxt_proxy;     # ← Fallback
    gzip_static on;                 # ← gzip_static (butuh pre-compressed files)
}
```

**Masalah:**
- ❌ Direct file serving memerlukan extract files manual
- ❌ `gzip_static` butuh files sudah di-compress sebelumnya
- ❌ `try_files` bisa menyebabkan loop atau error

---

### 2. **Dockerfile - Static Files Handling**

#### Referensi (BERHASIL):
```dockerfile
# Copy built application
COPY --from=builder --chown=nuxtjs:nodejs /app/.output /app/.output
COPY --from=builder --chown=nuxtjs:nodejs /app/node_modules /app/node_modules
COPY --from=builder --chown=nuxtjs:nodejs /app/package.json /app/package.json

# Copy public assets
COPY --from=builder --chown=nuxtjs:nodejs /app/public /app/public
```

**Key Points:**
- ✅ **Tidak ada copy manual ke `server/chunks/public`**
- ✅ **Nitro handle static files secara otomatis**
- ✅ **Copy `node_modules` secara eksplisit**
- ✅ **Copy `public` directory terpisah**

#### Config Kita (ERROR):
```dockerfile
# Copy manual ke server/chunks/public
RUN cp -rv .output/public/_nuxt .output/server/chunks/public/
```

**Masalah:**
- ❌ Copy manual mungkin tidak diperlukan
- ❌ Bisa menyebabkan file tidak sync
- ❌ Tidak copy `node_modules` secara eksplisit

---

### 3. **Docker Compose - Port & Health Check**

#### Referensi (BERHASIL):
```yaml
ports:
  - "3068:3000"  # ← Exposed ke host (bukan 127.0.0.1)

healthcheck:
  test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3000/seniman/api/health"]
```

**Key Points:**
- ✅ **Port exposed ke host** (bukan hanya localhost)
- ✅ **Health check menggunakan `wget`** (lebih reliable)
- ✅ **Health check ke endpoint spesifik** (`/api/health`)

#### Config Kita (ERROR):
```yaml
ports:
  - "127.0.0.1:3002:3000"  # ← Hanya localhost

healthcheck:
  test: ["CMD", "node", "-e", "require('http').get('http://localhost:3000', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"]
```

**Masalah:**
- ⚠️ Port hanya localhost (ini OK untuk security)
- ⚠️ Health check menggunakan Node.js (bisa lebih lambat)

---

### 4. **Nginx - Buffer Settings**

#### Referensi (BERHASIL):
```nginx
# Di location utama (bukan di static location)
proxy_buffering on;
proxy_buffer_size 8k;
proxy_buffers 16 8k;
proxy_busy_buffers_size 16k;
```

#### Config Kita:
```nginx
# Di location / (HTML pages)
proxy_buffer_size 4k;  # ← Lebih kecil
proxy_buffers 8 4k;    # ← Lebih kecil
```

---

## Solusi yang Harus Diterapkan

### Solusi 1: Proxy Static Files (Seperti Referensi)
**Ini adalah solusi TERBAIK** karena:
1. ✅ Tidak perlu extract files manual
2. ✅ Nitro handle static files dengan benar
3. ✅ Lebih simple dan reliable
4. ✅ Sudah terbukti bekerja di referensi

### Solusi 2: Perbaiki Dockerfile
1. ✅ Copy `node_modules` secara eksplisit
2. ✅ Copy `public` directory terpisah
3. ✅ Hapus copy manual ke `server/chunks/public` (Nitro handle ini)

### Solusi 3: Perbaiki Nginx Config
1. ✅ Gunakan `proxy_pass` untuk static files
2. ✅ Gunakan `gzip_vary on` (bukan `gzip_static`)
3. ✅ Hapus `try_files` dan fallback
4. ✅ Increase buffer sizes

---

## Kesimpulan

**Root Cause:**
- Kita mencoba serve static files langsung dari filesystem, tapi Nitro dengan preset `node-server` sudah handle static files dengan benar melalui proxy
- Copy manual ke `server/chunks/public` mungkin tidak diperlukan atau malah menyebabkan masalah

**Solusi Permanen:**
- **Kembali ke proxy approach** (seperti referensi yang berhasil)
- **Perbaiki Dockerfile** untuk copy dependencies dengan benar
- **Perbaiki Nginx config** dengan buffer settings yang lebih besar

