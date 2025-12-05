#!/bin/bash

# ============================================================================
# Apply Referensi Solution - Proxy Approach
# ============================================================================
# Script ini mengimplementasikan solusi berdasarkan referensi yang berhasil:
# - Proxy static files ke container (bukan direct file serving)
# - Update Dockerfile seperti referensi
# - Update nginx config dengan proxy approach
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DOMAIN="app3.anakhebat.web.id"
NGINX_CONFIG="/etc/nginx/sites-available/$DOMAIN"

if [ "$EUID" -ne 0 ]; then 
    SUDO="sudo"
else
    SUDO=""
fi

echo "=========================================="
echo "  Apply Referensi Solution - Proxy Approach"
echo "=========================================="
echo "Domain: $DOMAIN"
echo ""

# Step 1: Backup current configs
echo -e "${BLUE}[1/6]${NC} Backing up current configurations..."
cd "$PROJECT_DIR"

# Backup nginx config
if [ -f "$NGINX_CONFIG" ]; then
    $SUDO cp "$NGINX_CONFIG" "${NGINX_CONFIG}.backup.$(date +%Y%m%d_%H%M%S)"
    echo -e "${GREEN}✓${NC} Nginx config backed up"
fi

# Backup Dockerfile
if [ -f "$PROJECT_DIR/frontend/Dockerfile.prod" ]; then
    cp "$PROJECT_DIR/frontend/Dockerfile.prod" "$PROJECT_DIR/frontend/Dockerfile.prod.backup.$(date +%Y%m%d_%H%M%S)"
    echo -e "${GREEN}✓${NC} Dockerfile backed up"
fi
echo ""

# Step 2: Update Dockerfile (optional - user can choose)
echo -e "${BLUE}[2/6]${NC} Dockerfile update..."
if [ -f "$PROJECT_DIR/frontend/Dockerfile.prod.referensi" ]; then
    echo -e "${YELLOW}⚠${NC} Dockerfile referensi found: frontend/Dockerfile.prod.referensi"
    echo -e "${YELLOW}ℹ${NC} To use it, run:"
    echo "   cp frontend/Dockerfile.prod.referensi frontend/Dockerfile.prod"
    echo -e "${YELLOW}ℹ${NC} Then rebuild: docker compose -f docker-compose.prod.yml build app"
    echo ""
    read -p "Do you want to use the referensi Dockerfile now? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cp "$PROJECT_DIR/frontend/Dockerfile.prod.referensi" "$PROJECT_DIR/frontend/Dockerfile.prod"
        echo -e "${GREEN}✓${NC} Dockerfile updated"
        echo -e "${YELLOW}ℹ${NC} You need to rebuild the container after this script completes"
    else
        echo -e "${YELLOW}⚠${NC} Skipping Dockerfile update (using current)"
    fi
else
    echo -e "${YELLOW}⚠${NC} Dockerfile referensi not found, skipping"
fi
echo ""

# Step 3: Update nginx config
echo -e "${BLUE}[3/6]${NC} Updating nginx configuration..."
if [ -f "$PROJECT_DIR/nginx.conf.production.proxy" ]; then
    $SUDO cp "$PROJECT_DIR/nginx.conf.production.proxy" "$NGINX_CONFIG"
    echo -e "${GREEN}✓${NC} Nginx config updated (proxy approach)"
else
    echo -e "${RED}✗${NC} nginx.conf.production.proxy not found"
    exit 1
fi
echo ""

# Step 4: Test nginx config
echo -e "${BLUE}[4/6]${NC} Testing nginx configuration..."
if $SUDO nginx -t; then
    echo -e "${GREEN}✓${NC} Nginx configuration is valid"
else
    echo -e "${RED}✗${NC} Nginx configuration test failed"
    echo -e "${YELLOW}ℹ${NC} Restoring backup..."
    $SUDO cp "${NGINX_CONFIG}.backup."* "$NGINX_CONFIG" 2>/dev/null || true
    exit 1
fi
echo ""

# Step 5: Rebuild container (if Dockerfile was updated)
if [ -f "$PROJECT_DIR/frontend/Dockerfile.prod" ] && grep -q "Dockerfile.prod.referensi" "$PROJECT_DIR/frontend/Dockerfile.prod" 2>/dev/null || \
   [ -f "$PROJECT_DIR/frontend/Dockerfile.prod" ] && ! grep -q "server/chunks/public" "$PROJECT_DIR/frontend/Dockerfile.prod" 2>/dev/null; then
    echo -e "${BLUE}[5/6]${NC} Rebuilding container..."
    read -p "Do you want to rebuild the container now? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        docker compose -f docker-compose.prod.yml build app
        docker compose -f docker-compose.prod.yml up -d app
        echo -e "${GREEN}✓${NC} Container rebuilt and restarted"
        echo -e "${YELLOW}ℹ${NC} Waiting for container to be ready..."
        sleep 10
    else
        echo -e "${YELLOW}⚠${NC} Skipping rebuild (you need to rebuild manually)"
    fi
else
    echo -e "${BLUE}[5/6]${NC} Container rebuild..."
    echo -e "${YELLOW}ℹ${NC} Dockerfile not changed, skipping rebuild"
    echo -e "${YELLOW}ℹ${NC} If you want to rebuild: docker compose -f docker-compose.prod.yml build app"
fi
echo ""

# Step 6: Reload nginx
echo -e "${BLUE}[6/6]${NC} Reloading nginx..."
$SUDO systemctl reload nginx
echo -e "${GREEN}✓${NC} Nginx reloaded"
echo ""

# Verification
echo -e "${BLUE}Verification:${NC}"
echo ""

# Test static file access
echo "Testing static file access..."
HTTP_CODE=$(curl -sL -o /dev/null -w "%{http_code}" --max-time 10 "https://$DOMAIN/_nuxt/entry.js" 2>/dev/null || echo "000")
if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓${NC} Static file accessible: HTTP 200"
else
    echo -e "${YELLOW}⚠${NC} Static file access: HTTP $HTTP_CODE"
    echo -e "${YELLOW}ℹ${NC} This might be normal if container is still starting"
fi

# Check container logs
echo ""
echo "Checking container status..."
if docker ps | grep -q "gowa-app-prod"; then
    echo -e "${GREEN}✓${NC} Container is running"
    echo ""
    echo "Recent container logs:"
    docker logs --tail 20 gowa-app-prod 2>&1 | tail -5
else
    echo -e "${RED}✗${NC} Container is not running"
fi

echo ""
echo "=========================================="
echo -e "${GREEN}Setup completed!${NC}"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Test in browser: https://$DOMAIN"
echo "2. Check Network tab for static files (should be HTTP 200)"
echo "3. If errors persist, check logs:"
echo "   - Container: docker logs gowa-app-prod"
echo "   - Nginx: sudo tail -f /var/log/nginx/app3_error.log"
echo ""
echo "If you updated Dockerfile, make sure to rebuild:"
echo "   docker compose -f docker-compose.prod.yml build app"
echo "   docker compose -f docker-compose.prod.yml up -d app"
echo ""

