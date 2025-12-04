# Troubleshooting Error 503 pada File _nuxt/*.js

## Masalah
Error 503 (Service Unavailable) pada file-file static Nuxt seperti:
- `/_nuxt/DQUWNHGA.js`
- `/_nuxt/CcVpZqDF.js`
- `/_nuxt/BHZQsgze.js`

## Analisis Penyebab

### 1. Nuxt Server Tidak Berjalan atau Crash
**Cek:**
```bash
# Cek apakah container berjalan
docker ps | grep gowa-app

# Cek log container
docker logs gowa-app-prod

# Cek apakah port 3000 accessible
curl http://localhost:3000
```

### 2. Static Files Tidak Tersedia
**Struktur yang diharapkan setelah build:**
```
/app/
  .output/
    public/
      _nuxt/
        DQUWNHGA.js
        CcVpZqDF.js
        BHZQsgze.js
        ...
    server/
      index.mjs
      ...
```

**Cek di container:**
```bash
# Masuk ke container
docker exec -it gowa-app-prod sh

# Cek struktur directory
ls -la /app/
ls -la /app/public/
ls -la /app/public/_nuxt/

# Cek apakah file-file JS ada
ls -la /app/public/_nuxt/*.js
```

### 3. Nginx Tidak Bisa Connect ke Nuxt
**Cek:**
```bash
# Dari host, cek apakah port 3000 accessible
curl http://127.0.0.1:3000

# Cek log Nginx
sudo tail -f /var/log/nginx/error.log

# Test koneksi
telnet 127.0.0.1 3000
```

### 4. Masalah dengan Symlink
Symlink yang dibuat di Dockerfile:
```dockerfile
RUN mkdir -p /app/server/chunks && \
    ln -s /app/public /app/server/chunks/public
```

**Cek:**
```bash
# Di dalam container
ls -la /app/server/chunks/
ls -la /app/server/chunks/public/
ls -la /app/server/chunks/public/_nuxt/
```

## Langkah Troubleshooting

### Step 1: Cek Status Container dan Logs
```bash
# Cek status
docker ps -a | grep gowa-app

# Cek logs untuk error
docker logs --tail 100 gowa-app-prod

# Cek logs real-time
docker logs -f gowa-app-prod
```

### Step 2: Verifikasi Struktur File
```bash
# Masuk ke container
docker exec -it gowa-app-prod sh

# Cek apakah public directory ada dan berisi _nuxt/
cd /app
ls -la
ls -la public/
ls -la public/_nuxt/

# Jika _nuxt/ tidak ada, masalah ada di build process
```

### Step 3: Test Nuxt Server Langsung
```bash
# Di dalam container, test apakah server bisa serve static files
curl http://localhost:3000/_nuxt/DQUWNHGA.js

# Atau dari host (jika port exposed)
curl http://localhost:3000/_nuxt/DQUWNHGA.js
```

### Step 4: Cek Nginx Configuration
```bash
# Cek apakah Nginx config benar
sudo nginx -t

# Cek apakah proxy_pass benar
grep -A 10 "location /" /etc/nginx/sites-available/your-site-config

# Reload Nginx jika diperlukan
sudo systemctl reload nginx
```

### Step 5: Rebuild Container
Jika struktur file tidak benar, rebuild:
```bash
cd /path/to/gowa_umkm/frontend
docker build -f Dockerfile.prod -t gowa-frontend:latest .

# Atau rebuild via docker-compose
docker-compose -f docker-compose.prod.yml build app
docker-compose -f docker-compose.prod.yml up -d app
```

## Solusi yang Mungkin

### Solusi 1: Pastikan Static Files Di-copy dengan Benar
Verifikasi di Dockerfile bahwa `.output/public/` di-copy dengan benar:
```dockerfile
COPY --from=builder /app/.output ./
```

### Solusi 2: Periksa Build Process
Pastikan build menghasilkan file static:
```bash
# Build lokal untuk test
cd frontend
npm run build

# Cek apakah .output/public/_nuxt/ ada
ls -la .output/public/_nuxt/
```

### Solusi 3: Verifikasi Nuxt Config
Pastikan tidak ada konfigurasi yang mengganggu:
- `nitro.preset` harus `node-server`
- Tidak ada konfigurasi yang memindahkan static files

### Solusi 4: Cek Permissions
Pastikan file memiliki permission yang benar:
```bash
# Di dalam container
chmod -R 755 /app/public
```

## Checklist Debugging

- [ ] Container Nuxt berjalan dengan baik
- [ ] Log container tidak menunjukkan error
- [ ] Port 3000 accessible dari host
- [ ] Directory `/app/public/_nuxt/` ada dan berisi file JS
- [ ] Symlink `/app/server/chunks/public` mengarah ke `/app/public`
- [ ] Nginx bisa connect ke port 3000
- [ ] Nginx config sudah benar (proxy_pass ke 127.0.0.1:3000)
- [ ] Build process menghasilkan file static dengan benar

## Next Steps

Jika setelah semua langkah di atas masih error, cek:
1. Resource server (CPU, memory, disk space)
2. Network connectivity antara Nginx dan container
3. Firewall rules
4. Nuxt version compatibility issues

