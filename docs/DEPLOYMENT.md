# Production Deployment Guide

## Pre-Deployment Checklist

### 1. Environment Configuration

```bash
# Validate environment variables
./scripts/validate-env.sh

# Required variables:
✅ JWT_SECRET (min 32 characters)
✅ DB_PASSWORD (strong password)
✅ GOOGLE_CLIENT_ID
✅ GOOGLE_CLIENT_SECRET
✅ GEMINI_API_KEY
✅ CORS_ALLOWED_ORIGINS (production domain)
```

### 2. Security Audit

```bash
# Check for hardcoded secrets
grep -r "password\|secret\|api_key" backend/ --exclude-dir=vendor

# Verify .gitignore
cat .gitignore

# Test environment validation
docker compose -f docker-compose.prod.yml config
```

### 3. Database Migrations

All 19 migrations will run automatically on startup:

```
✅ 001_init_schema.sql
✅ 002_tenants_table.sql
✅ 003_whatsapp_devices.sql
✅ 004_ai_configs.sql
✅ 005_customer_insights_messages.sql
✅ 006_jid_mapping.sql
✅ 007_message_templates.sql
✅ 008_recurring_broadcasts.sql
✅ 009_ai_knowledge_base.sql
✅ 010_knowledge_base_simple.sql
✅ 011_alter_ai_configs.sql
✅ 012_multi_provider_ai.sql
✅ 013_customer_tags_analytics.sql
✅ 014_create_customers_table.sql
✅ 015_fix_customer_tags_notes.sql
✅ 016_create_ai_analytics_tables.sql
✅ 017_update_ai_configs_schema.sql
✅ 018_add_recurring_broadcast_columns.sql
✅ 019_create_knowledge_base.sql
```

## Deployment Steps

### Option 1: Docker Compose (Recommended)

```bash
# 1. Clone repository
git clone <repository-url>
cd gowa_umkm

# 2. Setup environment
cp .env.example .env
nano .env  # Configure all variables

# 3. Build and start
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up -d

# 4. Verify migrations
docker compose -f docker-compose.prod.yml logs api | grep migration

# 5. Check health
curl http://localhost:8080/health
```

### Option 2: Manual Deployment

```bash
# Backend
cd backend
go build -o gowa-api
./gowa-api

# Frontend
cd frontend
npm run build
npm run start

# Database
# Ensure PostgreSQL 15+ is running
# Run migrations manually if needed
```

## Post-Deployment Verification

### 1. Health Checks

```bash
# API health
curl https://your-domain.com/api/health

# Database connection
docker compose exec db psql -U gowa_user -d gowa_db -c "SELECT COUNT(*) FROM schema_migrations;"

# Redis connection
docker compose exec redis redis-cli ping
```

### 2. Feature Testing

```bash
# Test authentication
curl -X POST https://your-domain.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'

# Test WhatsApp status
curl https://your-domain.com/api/whatsapp/status \
  -H "Authorization: Bearer <token>"

# Test broadcast scheduler
docker compose logs api | grep Scheduler
```

### 3. Migration Verification

```bash
# Check all migrations applied
docker compose exec db psql -U gowa_user -d gowa_db -c "
  SELECT version, applied_at 
  FROM schema_migrations 
  ORDER BY version;
"

# Should show 19 migrations
```

## SSL/TLS Setup

### Using Let's Encrypt

```bash
# Install Certbot
sudo apt-get install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# Auto-renewal
sudo certbot renew --dry-run
```

### Nginx Configuration

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Backend API
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Monitoring

### Logs

```bash
# All services
docker compose -f docker-compose.prod.yml logs -f

# Specific service
docker compose -f docker-compose.prod.yml logs -f api

# Last 100 lines
docker compose -f docker-compose.prod.yml logs --tail=100 api
```

### Metrics

```bash
# Container stats
docker stats

# Database size
docker compose exec db psql -U gowa_user -d gowa_db -c "
  SELECT pg_size_pretty(pg_database_size('gowa_db'));
"

# Active connections
docker compose exec db psql -U gowa_user -d gowa_db -c "
  SELECT count(*) FROM pg_stat_activity;
"
```

## Backup & Recovery

### Database Backup

```bash
# Create backup
docker compose exec db pg_dump -U gowa_user gowa_db > backup_$(date +%Y%m%d).sql

# Restore backup
docker compose exec -T db psql -U gowa_user gowa_db < backup_20231203.sql
```

### Automated Backups

```bash
# Add to crontab
0 2 * * * /path/to/backup-script.sh
```

## Rollback Procedure

### If deployment fails:

```bash
# 1. Stop new version
docker compose -f docker-compose.prod.yml down

# 2. Restore database backup
docker compose exec -T db psql -U gowa_user gowa_db < backup_before_deploy.sql

# 3. Start previous version
git checkout <previous-commit>
docker compose -f docker-compose.prod.yml up -d

# 4. Verify
curl http://localhost:8080/health
```

## Troubleshooting

### Migration Errors

```bash
# Check migration status
docker compose exec db psql -U gowa_user -d gowa_db -c "
  SELECT * FROM schema_migrations ORDER BY version;
"

# Remove failed migration (if needed)
docker compose exec db psql -U gowa_user -d gowa_db -c "
  DELETE FROM schema_migrations WHERE version = 'XXX_failed.sql';
"

# Restart to retry
docker compose restart api
```

### Connection Issues

```bash
# Check network
docker network ls
docker network inspect gowa_umkm_default

# Check ports
netstat -tulpn | grep -E '8080|3000|5432|6379'

# Check firewall
sudo ufw status
```

### Performance Issues

```bash
# Check resource usage
docker stats

# Increase resources in docker-compose.prod.yml
services:
  api:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
```

## Security Hardening

### 1. Firewall Rules

```bash
# Allow only necessary ports
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 22/tcp
sudo ufw enable
```

### 2. Database Security

```bash
# Restrict database access
# Edit postgresql.conf
listen_addresses = 'localhost'

# Edit pg_hba.conf
host    all    all    127.0.0.1/32    scram-sha-256
```

### 3. Regular Updates

```bash
# Update system
sudo apt update && sudo apt upgrade

# Update Docker images
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
```

## Support

For issues:
1. Check logs: `docker compose logs -f`
2. Review documentation in `/docs`
3. Create GitHub issue
4. Contact support team

---

**Last Updated:** 2025-12-03
