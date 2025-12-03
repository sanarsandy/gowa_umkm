# Gowa UMKM - WhatsApp CRM System

> **Platform CRM berbasis WhatsApp untuk UMKM di Kabupaten Gowa**

Sistem manajemen pelanggan terintegrasi WhatsApp dengan fitur AI auto-reply, broadcast terjadwal, analytics, dan knowledge base untuk membantu UMKM mengelola komunikasi pelanggan secara efisien.

## 🚀 Fitur Utama

### ✅ WhatsApp Integration
- Multi-device WhatsApp Web connection
- QR code pairing
- Real-time message sync
- Message history tracking

### ✅ Customer Management
- Customer insights & analytics
- Lead scoring & segmentation
- Customer tags & notes
- Follow-up tracking

### ✅ AI Auto-Reply (Gemini)
- Intelligent auto-response
- Confidence-based escalation
- Business context awareness
- Knowledge base integration
- Multi-provider support (Gemini, OpenAI)

### ✅ Broadcast System
- One-time scheduled broadcasts
- Recurring broadcasts (hourly, daily, weekly)
- Message personalization ({{nama}}, {{name}})
- Recipient tracking & analytics

### ✅ Analytics & Reporting
- Message analytics
- Customer growth tracking
- AI performance metrics
- Broadcast statistics
- Hourly distribution analysis

### ✅ Knowledge Base
- FAQ management
- Product information
- Business policies
- Searchable by keywords & tags
- Priority-based ranking

## 🛠️ Tech Stack

**Backend:**
- Go (Echo framework)
- PostgreSQL
- Redis
- WhatsApp Web.js (via whmeow)

**Frontend:**
- Nuxt 3
- Vue 3
- TailwindCSS
- Pinia (state management)

**Infrastructure:**
- Docker & Docker Compose
- Nginx (reverse proxy)
- Let's Encrypt (SSL/TLS)

## 📋 Prerequisites

- Docker & Docker Compose
- Node.js 18+ (for local development)
- Go 1.21+ (for local development)
- PostgreSQL 15+
- Redis 7+

## 🚀 Quick Start

### 1. Clone Repository

```bash
git clone <repository-url>
cd gowa_umkm
```

### 2. Setup Environment

```bash
# Copy environment template
cp .env.example .env

# Edit .env and configure:
# - Database credentials
# - JWT secret (min 32 characters)
# - Google OAuth credentials
# - Gemini API key
# - CORS allowed origins
nano .env
```

### 3. Start Development

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Access application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

### 4. Database Migrations

Migrations run automatically on startup. To verify:

```bash
# Check migration status
docker compose exec db psql -U gowa_user -d gowa_db -c "SELECT * FROM schema_migrations ORDER BY version;"

# View logs
docker compose logs api | grep migration
```

## 📚 Documentation

Comprehensive documentation available in `/docs`:

- **[Quick Start Guide](docs/QUICKSTART.md)** - Get started in 5 minutes
- **[Docker Guide](docs/DOCKER_GUIDE.md)** - Docker setup and commands
- **[Migrations Guide](docs/MIGRATIONS.md)** - Database migration management
- **[Setup Checklist](docs/SETUP_CHECKLIST.md)** - Complete setup verification
- **[Feature Analysis](docs/ANALISIS_FITUR.md)** - Detailed feature documentation
- **[Roadmap](docs/ROADMAP_PENGEMBANGAN.md)** - Development roadmap

## 🔒 Security

### Critical Security Features

✅ **JWT Authentication**
- Secure token-based auth
- HttpOnly cookies
- 32+ character secret requirement
- Automatic validation on startup

✅ **Environment Validation**
- Startup checks for critical variables
- Secure secret validation
- Production-ready defaults

✅ **Rate Limiting**
- Stricter limits for auth endpoints
- Configurable per-endpoint limits
- Automatic cleanup

✅ **CORS Configuration**
- Environment-based allowed origins
- Secure production settings

### Security Checklist

Before deploying to production:

1. ✅ Set strong `JWT_SECRET` (min 32 chars)
2. ✅ Set strong `DB_PASSWORD`
3. ✅ Configure `CORS_ALLOWED_ORIGINS`
4. ✅ Set `ENV=production`
5. ✅ Use HTTPS/SSL certificates
6. ✅ Review and set all API keys
7. ✅ Run `./scripts/validate-env.sh`

## 🚢 Production Deployment

### Using Docker Compose

```bash
# Build and start production services
docker compose -f docker-compose.prod.yml up -d

# View logs
docker compose -f docker-compose.prod.yml logs -f

# Stop services
docker compose -f docker-compose.prod.yml down
```

### Environment Variables

**Required:**
- `JWT_SECRET` - JWT signing key (min 32 chars)
- `DB_PASSWORD` - Database password
- `GOOGLE_CLIENT_ID` - Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth secret
- `GEMINI_API_KEY` - Gemini AI API key

**Optional:**
- `CORS_ALLOWED_ORIGINS` - Allowed CORS origins (production)
- `OPENAI_API_KEY` - OpenAI API key (if using OpenAI)
- `ENV` - Environment (development/production)

See `.env.example` for complete list.

## 📊 Database Migrations

All migrations are in `backend/migrations/` and run automatically on startup:

```
001_init_schema.sql              - Initial schema
002_tenants_table.sql            - Multi-tenancy
003_whatsapp_devices.sql         - WhatsApp integration
004_ai_configs.sql               - AI configuration
005_customer_insights_messages.sql - Customer data
...
019_create_knowledge_base.sql    - Knowledge base
```

Migrations are:
- ✅ Idempotent (safe to run multiple times)
- ✅ Tracked in `schema_migrations` table
- ✅ Executed in alphabetical order
- ✅ Automatically run on container start

## 🧪 Testing

```bash
# Run backend tests
cd backend
go test ./...

# Run frontend tests
cd frontend
npm run test

# Integration tests
make test
```

## 📈 Monitoring

### Health Checks

```bash
# Backend health
curl http://localhost:8080/health

# Database connection
docker compose exec db psql -U gowa_user -d gowa_db -c "SELECT 1;"

# Redis connection
docker compose exec redis redis-cli ping
```

### Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f api
docker compose logs -f frontend
docker compose logs -f db
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For issues and questions:
- Create an issue in GitHub
- Check documentation in `/docs`
- Review logs: `docker compose logs -f`

## 🎯 Roadmap

See [ROADMAP](docs/ROADMAP_PENGEMBANGAN.md) for planned features and improvements.

---

**Built with ❤️ for UMKM Kabupaten Gowa**
