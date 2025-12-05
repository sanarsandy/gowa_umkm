# Setup Production VPS - Panduan Lengkap

Dokumentasi ini menjelaskan cara setup aplikasi Gowa UMKM di VPS production dengan optimasi untuk serving static files menggunakan best practices.

## 📋 Daftar Isi

1. [Persyaratan Sistem](#persyaratan-sistem)
2. [Persiapan VPS](#persiapan-vps)
3. [Instalasi Dependencies](#instalasi-dependencies)
4. [Konfigurasi Aplikasi](#konfigurasi-aplikasi)
5. [Deployment](#deployment)
6. [Setup Nginx untuk Static Files](#setup-nginx-untuk-static-files)
7. [Setup SSL/TLS](#setup-ssltls)
8. [Optimasi & Best Practices](#optimasi--best-practices)
9. [Monitoring & Maintenance](#monitoring--maintenance)
10. [Troubleshooting](#troubleshooting)

---

## 🖥️ Persyaratan Sistem

### Minimum Requirements
- **OS**: Ubuntu 20.04 LTS atau lebih baru / Debian 11+
- **RAM**: 2GB (4GB recommended)
- **CPU**: 2 cores (4 cores recommended)
- **Storage**: 20GB SSD
- **Network**: 1 Gbps connection

### Software Requirements
- Docker 20.10+
- Docker Compose 2.0+
- Nginx 1.18+
- Git
- Certbot (untuk SSL)

---

## 🚀 Persiapan VPS

### 1. Update Sistem

```bash
sudo apt update && sudo apt upgrade -y
```

### 2. Install Dependencies Dasar

```bash
# Install tools dasar
sudo apt install -y curl wget git ufw

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Logout dan login kembali untuk apply docker group
```

### 3. Konfigurasi Firewall

```bash
# Enable firewall
sudo ufw enable

# Allow SSH (penting! lakukan sebelum enable firewall)
sudo ufw allow 22/tcp

# Allow HTTP dan HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Check status
sudo ufw status
```

---

## 📦 Instalasi Dependencies

### 1. Install Nginx

```bash
sudo apt install -y nginx

# Enable dan start nginx
sudo systemctl enable nginx
sudo systemctl start nginx
```

### 2. Install Certbot (untuk SSL)

```bash
sudo apt install -y certbot python3-certbot-nginx
```

### 3. Install Nginx Modules (Optional - untuk Brotli compression)

```bash
# Install build dependencies
sudo apt install -y build-essential libpcre3-dev zlib1g-dev libssl-dev

# Clone nginx-brotli module
cd /tmp
git clone https://github.com/google/ngx_brotli.git
cd ngx_brotli
git submodule update --init

# Note: Brotli memerlukan rebuild nginx dari source
# Untuk production, gzip sudah cukup
```

---

## ⚙️ Konfigurasi Aplikasi

### 1. Clone Repository

```bash
cd /opt
sudo git clone <repository-url> gowa_umkm
cd gowa_umkm
sudo chown -R $USER:$USER .
```

### 2. Setup Environment Variables

```bash
# Copy template
cp .env.example .env

# Edit dengan editor favorit
nano .env
```

**Variabel penting yang harus di-set:**

```bash
# Database
DB_USER=gowa_user
DB_PASSWORD=<strong-password>
DB_NAME=gowa_db
DB_HOST=db

# JWT Secret (minimal 32 karakter)
JWT_SECRET=<generate-random-32-char-string>

# Google OAuth
GOOGLE_CLIENT_ID=<your-client-id>
GOOGLE_CLIENT_SECRET=<your-client-secret>
GOOGLE_REDIRECT_URL=https://yourdomain.com/api/auth/google/callback

# Frontend URL
FRONTEND_URL=https://yourdomain.com
NUXT_PUBLIC_API_BASE=https://yourdomain.com/api

# CORS
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# AI Provider
GEMINI_API_KEY=<your-gemini-api-key>
GEMINI_MODEL=gemini-1.5-flash

# Environment
ENV=production
```

### 3. Generate JWT Secret

```bash
# Generate random 32 character string
openssl rand -base64 32
```

---

## 🚢 Deployment

### Opsi 1: Menggunakan Script Deployment (Recommended)

```bash
# Berikan permission execute
chmod +x scripts/deploy-production.sh

# Jalankan deployment
./scripts/deploy-production.sh yourdomain.com
```

Script ini akan:
- ✅ Backup database
- ✅ Pull latest code
- ✅ Build containers
- ✅ Start services
- ✅ Extract static files
- ✅ Setup nginx configuration
- ✅ Reload nginx

### Opsi 2: Manual Deployment

```bash
# 1. Build containers
docker compose -f docker-compose.prod.yml build

# 2. Start services
docker compose -f docker-compose.prod.yml up -d

# 3. Wait for services to be ready
sleep 10

# 4. Extract static files
chmod +x scripts/extract-static-files.sh
./scripts/extract-static-files.sh

# 5. Setup nginx (lihat section berikutnya)
```

---

## 🌐 Setup Nginx untuk Static Files

### Best Practice: Direct File Serving

Untuk performa optimal, kita akan serve static files langsung dari filesystem Nginx, bukan melalui proxy ke Node.js.

### 1. Copy dan Konfigurasi Nginx

```bash
# Copy template
sudo cp nginx.conf.production /etc/nginx/sites-available/yourdomain.com

# Edit dan replace placeholders
sudo nano /etc/nginx/sites-available/yourdomain.com
```

**Replace:**
- `{{DOMAIN}}` → `yourdomain.com`
- `{{STATIC_ROOT}}` → `/var/www/gowa_umkm/static`

### 2. Create Static Files Directory

```bash
# Create directory
sudo mkdir -p /var/www/gowa_umkm/static

# Set permissions
sudo chown -R www-data:www-data /var/www/gowa_umkm/static
sudo chmod -R 755 /var/www/gowa_umkm/static
```

### 3. Extract Static Files dari Container

```bash
# Set STATIC_ROOT environment variable
export STATIC_ROOT=/var/www/gowa_umkm/static

# Run extraction script
chmod +x scripts/extract-static-files.sh
./scripts/extract-static-files.sh
```

### 4. Enable Site dan Test

```bash
# Create symlink
sudo ln -s /etc/nginx/sites-available/yourdomain.com /etc/nginx/sites-enabled/

# Remove default site (optional)
sudo rm /etc/nginx/sites-enabled/default

# Test configuration
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

### 5. Verify Static Files

```bash
# Check if files exist
ls -lh /var/www/gowa_umkm/static/_nuxt/

# Test from command line
curl -I http://localhost/_nuxt/entry.js
# Should return 200 OK
```

---

## 🔒 Setup SSL/TLS

### Menggunakan Let's Encrypt (Certbot)

```bash
# Obtain certificate
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Certbot akan:
# 1. Obtain certificate
# 2. Automatically update nginx config
# 3. Setup auto-renewal

# Test auto-renewal
sudo certbot renew --dry-run
```

### Manual SSL Setup

Jika sudah punya certificate:

```bash
# Copy certificates
sudo mkdir -p /etc/letsencrypt/live/yourdomain.com
sudo cp your-cert.pem /etc/letsencrypt/live/yourdomain.com/fullchain.pem
sudo cp your-key.pem /etc/letsencrypt/live/yourdomain.com/privkey.pem

# Update nginx config
sudo nano /etc/nginx/sites-available/yourdomain.com
# Uncomment SSL configuration

# Test dan reload
sudo nginx -t
sudo systemctl reload nginx
```

---

## ⚡ Optimasi & Best Practices

### 1. Nginx Cache untuk Static Files

Konfigurasi sudah termasuk di `nginx.conf.production`:

```nginx
proxy_cache_path /var/cache/nginx/static levels=1:2 keys_zone=static_cache:10m max_size=1g 
                 inactive=60m use_temp_path=off;
```

### 2. Gzip Compression

Sudah dikonfigurasi di nginx config. Verify:

```bash
# Test gzip
curl -H "Accept-Encoding: gzip" -I https://yourdomain.com/_nuxt/entry.js
# Should see: Content-Encoding: gzip
```

### 3. Browser Caching

Static files sudah dikonfigurasi dengan cache headers:
- `_nuxt/*`: 1 year (immutable)
- Images/Fonts: 6 months
- HTML: No cache

### 4. Rate Limiting

Sudah dikonfigurasi:
- API: 10 requests/second
- General: 30 requests/second

### 5. Security Headers

Sudah termasuk:
- X-Frame-Options
- X-Content-Type-Options
- X-XSS-Protection
- Referrer-Policy
- HSTS (optional)

### 6. Monitoring Disk Space

```bash
# Check static files size
du -sh /var/www/gowa_umkm/static

# Check nginx cache size
du -sh /var/cache/nginx/static

# Setup log rotation
sudo nano /etc/logrotate.d/nginx
```

### 7. Auto-extract Static Files setelah Rebuild

Tambahkan ke cron atau CI/CD:

```bash
# Edit crontab
crontab -e

# Add (setelah rebuild frontend)
0 2 * * * cd /opt/gowa_umkm && ./scripts/extract-static-files.sh
```

---

## 📊 Monitoring & Maintenance

### 1. Check Service Status

```bash
# Docker containers
docker compose -f docker-compose.prod.yml ps

# Nginx
sudo systemctl status nginx

# Check logs
docker compose -f docker-compose.prod.yml logs -f
sudo tail -f /var/log/nginx/yourdomain.com-access.log
```

### 2. Health Checks

```bash
# API health
curl https://yourdomain.com/api/health

# Frontend
curl -I https://yourdomain.com

# Static files
curl -I https://yourdomain.com/_nuxt/entry.js
```

### 3. Database Backup

```bash
# Manual backup
./scripts/backup-db.sh

# Setup automated backup (crontab)
0 2 * * * cd /opt/gowa_umkm && ./scripts/backup-db.sh
```

### 4. Update Aplikasi

```bash
# Pull latest code
git pull origin main

# Rebuild dan restart
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up -d

# Extract static files (penting!)
./scripts/extract-static-files.sh

# Reload nginx
sudo systemctl reload nginx
```

### 5. Resource Monitoring

```bash
# Container stats
docker stats

# Disk usage
df -h

# Memory usage
free -h

# Nginx connections
sudo netstat -an | grep :443 | wc -l
```

---

## 🔧 Troubleshooting

### Static Files Tidak Muncul

```bash
# 1. Check if files exist
ls -la /var/www/gowa_umkm/static/_nuxt/

# 2. Check permissions
sudo chown -R www-data:www-data /var/www/gowa_umkm/static
sudo chmod -R 755 /var/www/gowa_umkm/static

# 3. Re-extract files
./scripts/extract-static-files.sh

# 4. Check nginx config
sudo nginx -t
sudo tail -f /var/log/nginx/error.log
```

### 502 Bad Gateway

```bash
# Check if containers are running
docker compose -f docker-compose.prod.yml ps

# Check container logs
docker compose -f docker-compose.prod.yml logs app
docker compose -f docker-compose.prod.yml logs api

# Check if ports are accessible
curl http://localhost:3002
curl http://localhost:8087/api/health
```

### SSL Certificate Issues

```bash
# Check certificate
sudo certbot certificates

# Renew certificate
sudo certbot renew

# Check nginx SSL config
sudo nginx -t
```

### High Memory Usage

```bash
# Check container memory
docker stats

# Limit container resources in docker-compose.prod.yml
services:
  app:
    deploy:
      resources:
        limits:
          memory: 1G
        reservations:
          memory: 512M
```

### Static Files Masih di-Proxy (Tidak dari Filesystem)

```bash
# Verify nginx config
grep -A 10 "location /_nuxt/" /etc/nginx/sites-available/yourdomain.com
# Should show: alias /var/www/gowa_umkm/static/_nuxt/;

# Check nginx error log
sudo tail -f /var/log/nginx/error.log

# Test direct file access
curl -I https://yourdomain.com/_nuxt/entry.js
# Should return 200 OK with proper cache headers
```

---

## 📚 Referensi Best Practices

### Nginx Static File Serving
- [Nginx Performance Tuning](https://www.nginx.com/blog/tuning-nginx/)
- [Serving Static Content](https://docs.nginx.com/nginx/admin-guide/web-server/serving-static-content/)

### Docker Production
- [Docker Production Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Docker Compose Production](https://docs.docker.com/compose/production/)

### Security
- [OWASP Security Headers](https://owasp.org/www-project-secure-headers/)
- [SSL Labs SSL Test](https://www.ssllabs.com/ssltest/)

### Performance
- [Web.dev Performance](https://web.dev/performance/)
- [Lighthouse CI](https://github.com/GoogleChrome/lighthouse-ci)

---

## ✅ Checklist Deployment

- [ ] VPS sudah di-setup dengan OS yang sesuai
- [ ] Docker dan Docker Compose terinstall
- [ ] Nginx terinstall dan running
- [ ] Firewall dikonfigurasi (ports 22, 80, 443)
- [ ] Repository di-clone ke VPS
- [ ] Environment variables dikonfigurasi (.env)
- [ ] Database credentials aman
- [ ] JWT_SECRET sudah di-generate
- [ ] Containers berhasil di-build
- [ ] Services berjalan dengan baik
- [ ] Static files berhasil di-extract
- [ ] Nginx configuration sudah di-setup
- [ ] SSL certificate sudah di-install
- [ ] DNS records mengarah ke VPS
- [ ] Aplikasi bisa diakses via domain
- [ ] Static files di-serve dari filesystem (bukan proxy)
- [ ] Health checks passing
- [ ] Monitoring sudah di-setup
- [ ] Backup strategy sudah di-setup

---

**Last Updated**: 2025-01-XX
**Maintained by**: Gowa UMKM Team


