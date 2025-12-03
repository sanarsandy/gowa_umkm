#!/bin/bash

# Script untuk backup database
# Usage: ./scripts/backup-db.sh [backup-name]

set -e

BACKUP_DIR="backups"
mkdir -p $BACKUP_DIR

# Get DB credentials from environment or use defaults
DB_USER=${DB_USER:-{{PROJECT_NAME}}_user}
DB_NAME=${DB_NAME:-{{PROJECT_NAME}}_db}

# Generate backup filename
if [ -n "$1" ]; then
    BACKUP_FILE="$BACKUP_DIR/$1.sql"
else
    BACKUP_FILE="$BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql"
fi

echo "Backing up database: $DB_NAME"
echo "Backup file: $BACKUP_FILE"

# Try to detect container name
DB_CONTAINER=$(docker ps --format "{{.Names}}" | grep -E "postgres|{{PROJECT_NAME}}.*db" | head -1)

if [ -z "$DB_CONTAINER" ]; then
    echo "Error: Database container not found"
    exit 1
fi

# Create backup
if docker exec -T $DB_CONTAINER pg_dump -U $DB_USER $DB_NAME > "$BACKUP_FILE"; then
    echo "✓ Backup completed successfully"
    echo "File size: $(du -h "$BACKUP_FILE" | cut -f1)"
else
    echo "✗ Backup failed"
    exit 1
fi

# Optional: Compress backup
if command -v gzip &> /dev/null; then
    echo "Compressing backup..."
    gzip "$BACKUP_FILE"
    echo "✓ Compressed backup: ${BACKUP_FILE}.gz"
fi

