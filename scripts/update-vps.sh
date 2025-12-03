#!/bin/bash

# Script untuk update aplikasi di VPS
# Usage: ./scripts/update-vps.sh

set -e  # Exit on error

echo "=========================================="
echo "  Update Application on VPS"
echo "=========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. Backup Database
echo -e "${YELLOW}[1/7]${NC} Backup database..."
BACKUP_DIR="backups"
BACKUP_FILE="$BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql"
mkdir -p $BACKUP_DIR

# Get DB credentials from environment or use defaults
DB_USER=${DB_USER:-{{PROJECT_NAME}}_user}
DB_NAME=${DB_NAME:-{{PROJECT_NAME}}_db}

if docker compose -f docker-compose.prod.yml exec -T db pg_dump -U $DB_USER $DB_NAME > "$BACKUP_FILE" 2>/dev/null; then
    echo -e "${GREEN}✓${NC} Backup saved to: $BACKUP_FILE"
else
    echo -e "${YELLOW}⚠${NC} Backup failed, continuing..."
fi
echo ""

# 2. Pull dari Git
echo -e "${YELLOW}[2/7]${NC} Pull perubahan dari git..."
if git pull origin main; then
    echo -e "${GREEN}✓${NC} Git pull completed"
else
    echo -e "${RED}✗${NC} Git pull failed"
    exit 1
fi
echo ""

# 3. Run Database Migration
echo -e "${YELLOW}[3/7]${NC} Run database migration..."
if [ -f "backend/scripts/run-migrations.sh" ]; then
    chmod +x backend/scripts/run-migrations.sh
    cd backend && ./scripts/run-migrations.sh docker && cd ..
    echo -e "${GREEN}✓${NC} Migration completed"
else
    echo -e "${YELLOW}⚠${NC} Migration script not found, skipping..."
fi
echo ""

# 4. Rebuild Backend
echo -e "${YELLOW}[4/7]${NC} Rebuild backend container..."
docker compose -f docker-compose.prod.yml build api
echo -e "${GREEN}✓${NC} Backend rebuild completed"
echo ""

# 5. Rebuild Frontend
echo -e "${YELLOW}[5/7]${NC} Rebuild frontend container..."
docker compose -f docker-compose.prod.yml build app
echo -e "${GREEN}✓${NC} Frontend rebuild completed"
echo ""

# 6. Restart Services
echo -e "${YELLOW}[6/7]${NC} Restart services..."
docker compose -f docker-compose.prod.yml up -d
echo -e "${GREEN}✓${NC} Services restarted"
echo ""

# 7. Wait and Verify
echo -e "${YELLOW}[7/7]${NC} Waiting for services to be ready..."
sleep 10

echo "Service Status:"
docker compose -f docker-compose.prod.yml ps
echo ""

# Check API health
echo "Checking API health..."
API_PORT=${API_PORT:-8085}
if curl -s http://localhost:$API_PORT/health > /dev/null; then
    echo -e "${GREEN}✓${NC} API is healthy"
else
    echo -e "${RED}✗${NC} API health check failed"
fi
echo ""

echo "=========================================="
echo -e "${GREEN}Update completed!${NC}"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Check logs: docker compose -f docker-compose.prod.yml logs -f"
echo "2. Test application in browser"
echo "3. Verify all features are working"
echo ""

