# Database Migrations Guide

## Overview

Database migrations are automatically executed when the application starts. All SQL files in `backend/migrations/` are run in alphabetical order.

## How It Works

1. **Automatic Execution**: Migrations run automatically on application startup
2. **Tracking**: A `schema_migrations` table tracks which migrations have been applied
3. **Idempotent**: Already-applied migrations are skipped
4. **Ordered**: Migrations are executed in alphabetical order (001, 002, 003, etc.)

## Migration Files

All migration files are located in `backend/migrations/`:

```
backend/migrations/
├── 001_init_schema.sql
├── 002_tenants_table.sql
├── 003_whatsapp_devices.sql
├── 004_ai_configs.sql
├── 005_customer_insights_messages.sql
├── 006_jid_mapping.sql
├── 007_message_templates.sql
├── 008_recurring_broadcasts.sql
├── 009_ai_knowledge_base.sql
├── 010_knowledge_base_simple.sql
├── 011_alter_ai_configs.sql
├── 012_multi_provider_ai.sql
├── 013_customer_tags_analytics.sql
├── 014_create_customers_table.sql
├── 015_fix_customer_tags_notes.sql
├── 016_create_ai_analytics_tables.sql
└── 017_update_ai_configs_schema.sql
```

## Creating New Migrations

1. **Naming Convention**: Use format `XXX_description.sql` where XXX is a sequential number
   ```
   018_add_new_feature.sql
   019_alter_table_xyz.sql
   ```

2. **Content**: Write idempotent SQL using `IF NOT EXISTS` or `IF EXISTS`
   ```sql
   -- Good: Idempotent
   CREATE TABLE IF NOT EXISTS my_table (...);
   ALTER TABLE my_table ADD COLUMN IF NOT EXISTS my_column VARCHAR(255);
   
   -- Bad: Will fail if run twice
   CREATE TABLE my_table (...);
   ALTER TABLE my_table ADD COLUMN my_column VARCHAR(255);
   ```

3. **Testing**: Test locally before deploying
   ```bash
   # Run migrations manually
   docker compose exec db psql -U gowa_user -d gowa_db -f /app/migrations/018_your_migration.sql
   ```

## Deployment Process

### Development
```bash
# Migrations run automatically on startup
docker compose up

# Check migration status
docker compose exec db psql -U gowa_user -d gowa_db -c "SELECT * FROM schema_migrations ORDER BY version;"
```

### Production
```bash
# Using docker-compose.prod.yml
docker compose -f docker-compose.prod.yml up -d

# Migrations will run automatically on first startup
# Check logs
docker compose -f docker-compose.prod.yml logs api | grep migration
```

## Migration Status

### Check Applied Migrations
```bash
docker compose exec db psql -U gowa_user -d gowa_db -c "
  SELECT version, applied_at 
  FROM schema_migrations 
  ORDER BY version;
"
```

### Check Pending Migrations
```bash
# List all migration files
ls -1 backend/migrations/*.sql

# Compare with applied migrations
docker compose exec db psql -U gowa_user -d gowa_db -c "SELECT version FROM schema_migrations;"
```

## Troubleshooting

### Migration Failed

If a migration fails:

1. **Check logs**:
   ```bash
   docker compose logs api | grep -A 10 "migration"
   ```

2. **Fix the migration file** and restart:
   ```bash
   # Edit the migration file
   nano backend/migrations/XXX_failed_migration.sql
   
   # Remove from tracking (if partially applied)
   docker compose exec db psql -U gowa_user -d gowa_db -c "
     DELETE FROM schema_migrations WHERE version = 'XXX_failed_migration.sql';
   "
   
   # Restart to retry
   docker compose restart api
   ```

### Reset All Migrations (DANGER!)

**⚠️ WARNING: This will delete all data!**

```bash
# Drop all tables
docker compose exec db psql -U gowa_user -d gowa_db -c "
  DROP SCHEMA public CASCADE;
  CREATE SCHEMA public;
"

# Restart to re-run all migrations
docker compose restart api
```

## Best Practices

1. **Always use IF NOT EXISTS / IF EXISTS**
   - Makes migrations idempotent
   - Safe to run multiple times

2. **Test locally first**
   - Run migration manually before committing
   - Verify it works on fresh database

3. **Never modify existing migrations**
   - Once deployed, create a new migration instead
   - Existing migrations may have already run in production

4. **Use descriptive names**
   - `018_add_user_preferences.sql` ✅
   - `018_update.sql` ❌

5. **Add comments**
   ```sql
   -- Migration 018: Add User Preferences
   -- Adds user_preferences table for storing user settings
   
   CREATE TABLE IF NOT EXISTS user_preferences (...);
   ```

6. **Keep migrations small**
   - One logical change per migration
   - Easier to debug and rollback

## Production Deployment Checklist

- [ ] All migrations tested locally
- [ ] Migration files committed to git
- [ ] Backup database before deployment
- [ ] Deploy new code
- [ ] Verify migrations ran successfully
- [ ] Check application logs for errors
- [ ] Test critical features

## Rollback Strategy

If you need to rollback:

1. **Create a rollback migration**:
   ```sql
   -- 019_rollback_user_preferences.sql
   DROP TABLE IF EXISTS user_preferences;
   ```

2. **Deploy the rollback migration**:
   ```bash
   docker compose restart api
   ```

## Monitoring

### Check Migration Status on Startup
```bash
# Watch logs during startup
docker compose logs -f api | grep migration
```

Expected output:
```
🔄 Running database migrations...
  ✅ Applied: 001_init_schema.sql
  ✅ Applied: 002_tenants_table.sql
  ...
✅ Migrations complete: 17 applied, 0 skipped, 17 total
```

### Verify Database Schema
```bash
# List all tables
docker compose exec db psql -U gowa_user -d gowa_db -c "\dt"

# Check specific table structure
docker compose exec db psql -U gowa_user -d gowa_db -c "\d ai_configs"
```
