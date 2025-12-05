#!/bin/bash

# ============================================================================
# Setup Production Final - Solusi Permanen
# ============================================================================
# Script ini mengimplementasikan best practice untuk Nuxt.js production:
# - Extract static files dari container ke filesystem
# - Serve static files langsung dari Nginx (tidak melalui proxy)
# - Hanya proxy HTML pages ke Node.js
#
# Usage: ./scripts/setup-production-final.sh
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STATIC_ROOT="${STATIC_ROOT:-/var/www/gowa_umkm/static}"
DOMAIN="app3.anakhebat.web.id"
NGINX_CONFIG="/etc/nginx/sites-available/$DOMAIN"

if [ "$EUID" -ne 0 ]; then 
    SUDO="sudo"
else
    SUDO=""
fi

echo "=========================================="
echo "  Setup Production Final - Best Practice"
echo "=========================================="
echo "Domain: $DOMAIN"
echo "Static Root: $STATIC_ROOT"
echo ""

# Step 1: Ensure container is running
echo -e "${BLUE}[1/5]${NC} Ensuring container is running..."
cd "$PROJECT_DIR"
if ! docker ps | grep -q "gowa-app-prod"; then
    docker compose -f docker-compose.prod.yml up -d app
    sleep 10
fi
echo -e "${GREEN}✓${NC} Container is running"
echo ""

# Step 2: Extract static files
echo -e "${BLUE}[2/5]${NC} Extracting static files from container..."
if [ -f "$PROJECT_DIR/scripts/extract-static-files.sh" ]; then
    chmod +x "$PROJECT_DIR/scripts/extract-static-files.sh"
    STATIC_ROOT="$STATIC_ROOT" "$PROJECT_DIR/scripts/extract-static-files.sh"
else
    echo -e "${YELLOW}⚠${NC} extract-static-files.sh not found, creating manually..."
    $SUDO mkdir -p "$STATIC_ROOT"
    docker cp gowa-app-prod:/app/.output/public/_nuxt "$STATIC_ROOT/" 2>/dev/null || true
    docker cp gowa-app-prod:/app/.output/public/builds "$STATIC_ROOT/" 2>/dev/null || true
    $SUDO chown -R www-data:www-data "$STATIC_ROOT"
    $SUDO chmod -R 755 "$STATIC_ROOT"
fi
echo ""

# Step 3: Verify static files
echo -e "${BLUE}[3/5]${NC} Verifying static files..."
if [ -d "$STATIC_ROOT/_nuxt" ]; then
    FILE_COUNT=$(find "$STATIC_ROOT/_nuxt" -type f | wc -l)
    echo -e "${GREEN}✓${NC} Found $FILE_COUNT static files"
    if [ "$FILE_COUNT" -lt 10 ]; then
        echo -e "${RED}✗${NC} Too few files! Something is wrong."
        exit 1
    fi
else
    echo -e "${RED}✗${NC} Static files directory not found!"
    exit 1
fi
echo ""

# Step 4: Update nginx config
echo -e "${BLUE}[4/5]${NC} Updating nginx configuration..."
if [ -f "$PROJECT_DIR/nginx.conf.production.final" ]; then
    # Backup current config
    if [ -f "$NGINX_CONFIG" ]; then
        $SUDO cp "$NGINX_CONFIG" "${NGINX_CONFIG}.backup.$(date +%Y%m%d_%H%M%S)"
    fi
    
    # Copy new config
    $SUDO cp "$PROJECT_DIR/nginx.conf.production.final" "$NGINX_CONFIG"
    
    # Replace placeholders
    $SUDO sed -i "s|{{STATIC_ROOT}}|$STATIC_ROOT|g" "$NGINX_CONFIG"
    
    # Test config
    if $SUDO nginx -t; then
        echo -e "${GREEN}✓${NC} Nginx configuration is valid"
    else
        echo -e "${RED}✗${NC} Nginx configuration test failed"
        exit 1
    fi
else
    echo -e "${RED}✗${NC} nginx.conf.production.final not found"
    exit 1
fi
echo ""

# Step 5: Reload nginx
echo -e "${BLUE}[5/5]${NC} Reloading nginx..."
$SUDO systemctl reload nginx
echo -e "${GREEN}✓${NC} Nginx reloaded"
echo ""

# Verification
echo -e "${BLUE}Verification:${NC}"
echo ""

# Test static file access
FIRST_FILE=$(find "$STATIC_ROOT/_nuxt" -type f -name "*.js" | head -1)
if [ -n "$FIRST_FILE" ]; then
    FILENAME=$(basename "$FIRST_FILE")
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 "http://127.0.0.1/_nuxt/$FILENAME" 2>/dev/null || echo "000")
    if [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓${NC} Static file served directly by Nginx: HTTP 200"
    else
        echo -e "${YELLOW}⚠${NC} Static file access: HTTP $HTTP_CODE"
    fi
fi

echo ""
echo "=========================================="
echo -e "${GREEN}Setup completed!${NC}"
echo "=========================================="
echo ""
echo "Static files are now served directly by Nginx (not proxied)"
echo "This is the recommended best practice for maximum performance"
echo ""
echo "Test: https://$DOMAIN"
echo ""

