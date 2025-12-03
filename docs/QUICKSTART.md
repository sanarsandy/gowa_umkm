# 🚀 Quick Start - Development Mode

Selamat datang di Gowa UMKM! Untuk development lokal, gunakan mode development yang jauh lebih ringan dan cepat.

## ⚡ Development Mode (Recommended)

```bash
# Start development mode dengan hot-reload
make dev

# Atau manual:
docker compose -f docker-compose.dev.yml up -d
```

**Keuntungan Development Mode:**
- ✅ **Hot-reload** - Perubahan code langsung terlihat tanpa rebuild
- ✅ **Cepat** - Startup ~30 detik (vs ~10 menit production)
- ✅ **Ringan** - Tidak perlu build production bundle
- ✅ **Efisien** - Source code di-mount sebagai volume

**Akses Aplikasi:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Database: localhost:5432

## 📝 Common Commands

```bash
# Development
make dev              # Start development mode
make dev-logs         # View logs
make dev-restart      # Restart services
make dev-down         # Stop services

# Production (untuk testing production build)
make up               # Start production mode
make down             # Stop production mode

# Utilities
make shell-api        # Shell ke backend container
make shell-app        # Shell ke frontend container
make shell-db         # Access database
```

## 📚 Full Documentation

Lihat [DOCKER_GUIDE.md](./DOCKER_GUIDE.md) untuk dokumentasi lengkap.

## 🐛 Troubleshooting

**Port sudah digunakan?**
```bash
# Stop semua container
make dev-down

# Atau cek port yang digunakan
lsof -i :3000
lsof -i :8080
```

**Perlu clean rebuild?**
```bash
make dev-clean        # Stop dan hapus volumes
make dev-build        # Rebuild containers
```

## 💡 Tips

1. **Selalu gunakan development mode** untuk local development
2. **Jangan restart container** setiap kali ubah code - hot-reload otomatis!
3. **Gunakan `make dev-logs`** untuk debug
4. **Test production build** sebelum deploy dengan `make up`
