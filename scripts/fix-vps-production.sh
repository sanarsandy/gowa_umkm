#!/bin/bash

# ============================================================================
# Script Perbaikan VPS Production - Fix 502/503 Static Files Error
# ============================================================================
# Script ini akan melakukan perbaikan lengkap di VPS production
# Usage: ./scripts/fix-vps-production.sh
# ============================================================================

set -e  # Exit on error

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DOMAIN="app3.anakhebat.web.id"
NGINX_CONFIG="/etc/nginx/sites-available/$DOMAIN"
STATIC_ROOT="/var/www/gowa_umkm/static"

echo "=========================================="
echo "  Fix VPS Production - Static Files 502/503"
echo "=========================================="
echo "Domain: $DOMAIN"
echo "Project Dir: $PROJECT_DIR"
echo ""

# Check if running as root for certain operations
if [ "$EUID" -ne 0 ]; then 
    SUDO="sudo"
else
    SUDO=""
fi

# ============================================================================
# Step 1: Backup Current Setup
# ============================================================================
echo -e "${BLUE}[1/10]${NC} Creating backup..."
BACKUP_DIR="$PROJECT_DIR/backups"
mkdir -p "$BACKUP_DIR"
BACKUP_FILE="$BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql"

# Backup database
if docker compose -f "$PROJECT_DIR/docker-compose.prod.yml" ps | grep -q "gowa-db-prod"; then
    if docker compose -f "$PROJECT_DIR/docker-compose.prod.yml" exec -T db pg_dump -U gowa_user gowa_db > "$BACKUP_FILE" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} Database backup saved: $BACKUP_FILE"
    else
        echo -e "${YELLOW}⚠${NC} Database backup failed, continuing..."
    fi
fi

# Backup nginx config
if [ -f "$NGINX_CONFIG" ]; then
    $SUDO cp "$NGINX_CONFIG" "${NGINX_CONFIG}.backup.$(date +%Y%m%d_%H%M%S)"
    echo -e "${GREEN}✓${NC} Nginx config backed up"
fi
echo ""

# ============================================================================
# Step 2: Pull Latest Changes from Git
# ============================================================================
echo -e "${BLUE}[2/10]${NC} Pulling latest changes from git..."
cd "$PROJECT_DIR"
if [ -d ".git" ]; then
    if git pull origin main; then
        echo -e "${GREEN}✓${NC} Git pull completed"
    else
        echo -e "${RED}✗${NC} Git pull failed"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠${NC} Not a git repository, skipping pull"
fi
echo ""

# ============================================================================
# Step 3: Stop Services
# ============================================================================
echo -e "${BLUE}[3/10]${NC} Stopping services..."
cd "$PROJECT_DIR"
docker compose -f docker-compose.prod.yml stop app
echo -e "${GREEN}✓${NC} Services stopped"
echo ""

# ============================================================================
# Step 4: Rebuild Container (No Cache)
# ============================================================================
echo -e "${BLUE}[4/10]${NC} Rebuilding container (this may take a few minutes)..."
cd "$PROJECT_DIR"
docker compose -f docker-compose.prod.yml build --no-cache app
echo -e "${GREEN}✓${NC} Container rebuilt"
echo ""

# ============================================================================
# Step 5: Start Services
# ============================================================================
echo -e "${BLUE}[5/10]${NC} Starting services..."
cd "$PROJECT_DIR"
docker compose -f docker-compose.prod.yml up -d app
echo -e "${GREEN}✓${NC} Services started"
echo ""

# ============================================================================
# Step 6: Wait for Container to be Healthy
# ============================================================================
echo -e "${BLUE}[6/10]${NC} Waiting for container to be healthy..."
MAX_WAIT=60
WAIT_COUNT=0
while [ $WAIT_COUNT -lt $MAX_WAIT ]; do
    HEALTH=$(docker inspect --format='{{.State.Health.Status}}' gowa-app-prod 2>/dev/null || echo "starting")
    if [ "$HEALTH" = "healthy" ]; then
        echo -e "${GREEN}✓${NC} Container is healthy"
        break
    fi
    WAIT_COUNT=$((WAIT_COUNT + 5))
    echo -e "${YELLOW}ℹ${NC} Waiting... ($WAIT_COUNT/$MAX_WAIT seconds) - Status: $HEALTH"
    sleep 5
done

if [ "$HEALTH" != "healthy" ]; then
    echo -e "${YELLOW}⚠${NC} Container health check timeout, but continuing..."
fi
echo ""

# ============================================================================
# Step 7: Verify Static Files in Container
# ============================================================================
echo -e "${BLUE}[7/10]${NC} Verifying static files in container..."
if docker exec gowa-app-prod test -d "/app/.output/public/_nuxt" 2>/dev/null; then
    FILE_COUNT=$(docker exec gowa-app-prod find /app/.output/public/_nuxt -type f 2>/dev/null | wc -l)
    echo -e "${GREEN}✓${NC} Static files found: $FILE_COUNT files"
else
    echo -e "${RED}✗${NC} Static files NOT found in container"
    echo -e "${YELLOW}ℹ${NC} Check container logs: docker logs gowa-app-prod"
fi
echo ""

# ============================================================================
# Step 8: Update Nginx Configuration
# ============================================================================
echo -e "${BLUE}[8/10]${NC} Updating nginx configuration..."
if [ -f "$PROJECT_DIR/nginx.conf.production" ]; then
    # Replace domain placeholder if needed
    $SUDO cp "$PROJECT_DIR/nginx.conf.production" "$NGINX_CONFIG"
    
    # Replace static root if needed (if using direct file serving)
    if grep -q "{{STATIC_ROOT}}" "$NGINX_CONFIG"; then
        $SUDO sed -i "s|{{STATIC_ROOT}}|$STATIC_ROOT|g" "$NGINX_CONFIG"
    fi
    
    # Test nginx configuration
    if $SUDO nginx -t; then
        echo -e "${GREEN}✓${NC} Nginx configuration is valid"
    else
        echo -e "${RED}✗${NC} Nginx configuration test failed"
        echo -e "${YELLOW}ℹ${NC} Restoring backup..."
        $SUDO cp "${NGINX_CONFIG}.backup."* "$NGINX_CONFIG" 2>/dev/null || true
        exit 1
    fi
else
    echo -e "${RED}✗${NC} nginx.conf.production not found"
    exit 1
fi
echo ""

# ============================================================================
# Step 9: Reload Nginx
# ============================================================================
echo -e "${BLUE}[9/10]${NC} Reloading nginx..."
if $SUDO systemctl reload nginx; then
    echo -e "${GREEN}✓${NC} Nginx reloaded"
else
    echo -e "${RED}✗${NC} Failed to reload nginx"
    exit 1
fi
echo ""

# ============================================================================
# Step 10: Verification
# ============================================================================
echo -e "${BLUE}[10/10]${NC} Verifying fixes..."
echo ""

# Check container status
echo "Container Status:"
docker ps --filter "name=gowa-app-prod" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
echo ""

# Test port accessibility
echo "Testing port accessibility..."
if curl -s -o /dev/null -w "%{http_code}" --max-time 5 http://127.0.0.1:3002 > /dev/null 2>&1; then
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 http://127.0.0.1:3002)
    if [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓${NC} Port 3002 is accessible (HTTP $HTTP_CODE)"
    else
        echo -e "${YELLOW}⚠${NC} Port 3002 returned HTTP $HTTP_CODE"
    fi
else
    echo -e "${RED}✗${NC} Port 3002 is not accessible"
fi
echo ""

# Test static file access
echo "Testing static file access..."
FIRST_FILE=$(docker exec gowa-app-prod find /app/.output/public/_nuxt -type f -name "*.js" 2>/dev/null | head -1)
if [ -n "$FIRST_FILE" ]; then
    FILE_NAME=$(basename "$FIRST_FILE")
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 "http://127.0.0.1:3002/_nuxt/$FILE_NAME" 2>/dev/null || echo "000")
    if [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓${NC} Static file access successful (HTTP 200)"
    elif [ "$HTTP_CODE" = "502" ] || [ "$HTTP_CODE" = "503" ]; then
        echo -e "${RED}✗${NC} Static file access failed: HTTP $HTTP_CODE"
        echo -e "${YELLOW}ℹ${NC} Run diagnostic: ./scripts/diagnose-static-files.sh"
    else
        echo -e "${YELLOW}⚠${NC} Static file access returned: HTTP $HTTP_CODE"
    fi
else
    echo -e "${YELLOW}⚠${NC} Cannot test static file access (no files found)"
fi
echo ""

# ============================================================================
# Summary
# ============================================================================
echo "=========================================="
echo -e "${GREEN}Fix completed!${NC}"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Test application in browser: https://$DOMAIN"
echo "2. Check browser console (F12) for any errors"
echo "3. Monitor logs: docker compose -f docker-compose.prod.yml logs -f app"
echo "4. If issues persist, run: ./scripts/diagnose-static-files.sh"
echo ""
echo "Useful commands:"
echo "  - View container logs: docker logs -f gowa-app-prod"
echo "  - Check nginx logs: sudo tail -f /var/log/nginx/app3_error.log"
echo "  - Restart services: docker compose -f docker-compose.prod.yml restart app"
echo ""

