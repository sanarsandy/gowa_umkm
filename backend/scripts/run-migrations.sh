#!/bin/bash

# Script untuk menjalankan migrations secara manual
# Usage: ./run-migrations.sh [docker|local]
#   docker: Run migrations via Docker (default if psql not found)
#   local: Run migrations directly using psql

set -e

echo "Running database migrations..."

# List migrations in order (add your migrations here)
MIGRATIONS=(
    "001_init_schema.sql"
)

# Database connection (adjust as needed)
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-gowa_user}
DB_NAME=${DB_NAME:-gowa_db}
DB_PASSWORD=${DB_PASSWORD:-change_me_secure_password}

# Docker container name (auto-detect or use environment variable)
if [ -z "$DB_CONTAINER" ]; then
    # Try to find postgres container automatically
    DB_CONTAINER=$(docker ps --format "{{.Names}}" | grep -E "postgres|gowa.*db" | head -1)
    if [ -z "$DB_CONTAINER" ]; then
        # Fallback to default
        DB_CONTAINER="gowa-db"
    fi
fi
echo "Using Docker container: $DB_CONTAINER"

# Detect if we should use Docker
USE_DOCKER=false
if [ "$1" = "docker" ]; then
    USE_DOCKER=true
elif [ "$1" = "local" ]; then
    USE_DOCKER=false
else
    # Auto-detect: check if psql is available
    if ! command -v psql &> /dev/null; then
        echo "psql not found, using Docker..."
        USE_DOCKER=true
    else
        USE_DOCKER=false
    fi
fi

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MIGRATIONS_DIR="$SCRIPT_DIR/../migrations"

# Run each migration
for migration in "${MIGRATIONS[@]}"; do
    if [ ! -f "$MIGRATIONS_DIR/$migration" ]; then
        echo "⚠ Migration file not found: $migration (skipping...)"
        continue
    fi
    
    echo "Running migration: $migration"
    
    if [ "$USE_DOCKER" = true ]; then
        # Run via Docker
        if docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME < "$MIGRATIONS_DIR/$migration" 2>&1; then
            echo "✓ Migration $migration completed"
        else
            # Check if it's just a "already exists" error (which is OK for idempotent migrations)
            if docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME < "$MIGRATIONS_DIR/$migration" 2>&1 | grep -q "already exists\|duplicate"; then
                echo "  (Migration already applied, skipping...)"
            else
                echo "✗ Migration $migration failed"
            fi
        fi
    else
        # Run directly using psql
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$MIGRATIONS_DIR/$migration" 2>&1 || {
            echo "✗ Migration $migration failed"
        }
    fi
done

echo "All migrations completed!"

