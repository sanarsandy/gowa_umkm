#!/bin/bash

# ============================================================================
# Script untuk Restart Nginx dan Test
# ============================================================================
# Restart nginx dan test apakah masih ada error
# Usage: ./scripts/restart-nginx-and-test.sh
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "=========================================="
echo "  Restart Nginx and Test"
echo "=========================================="
echo ""

# Backup current error log
echo -e "${BLUE}[1/4]${NC} Backing up current error log..."
if [ -f "/var/log/nginx/app3_error.log" ]; then
    sudo cp /var/log/nginx/app3_error.log /var/log/nginx/app3_error.log.backup.$(date +%Y%m%d_%H%M%S)
    echo -e "${GREEN}✓${NC} Error log backed up"
else
    echo -e "${YELLOW}⚠${NC} Error log not found"
fi
echo ""

# Clear error log
echo -e "${BLUE}[2/4]${NC} Clearing error log..."
sudo truncate -s 0 /var/log/nginx/app3_error.log 2>/dev/null || sudo sh -c '> /var/log/nginx/app3_error.log'
echo -e "${GREEN}✓${NC} Error log cleared"
echo ""

# Test nginx config
echo -e "${BLUE}[3/4]${NC} Testing nginx configuration..."
if sudo nginx -t; then
    echo -e "${GREEN}✓${NC} Nginx configuration is valid"
else
    echo -e "${RED}✗${NC} Nginx configuration test failed"
    exit 1
fi
echo ""

# Restart nginx
echo -e "${BLUE}[4/4]${NC} Restarting nginx..."
if sudo systemctl restart nginx; then
    echo -e "${GREEN}✓${NC} Nginx restarted"
    sleep 2
else
    echo -e "${RED}✗${NC} Failed to restart nginx"
    exit 1
fi
echo ""

# Test static file access
echo -e "${BLUE}Testing static file access...${NC}"
FIRST_FILE=$(docker exec gowa-app-prod find /app/.output/server/chunks/public/_nuxt -type f -name "*.js" 2>/dev/null | head -1)
if [ -n "$FIRST_FILE" ]; then
    FILENAME=$(basename "$FIRST_FILE")
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 10 "http://127.0.0.1:3002/_nuxt/$FILENAME" 2>/dev/null || echo "000")
    
    if [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓${NC} Static file access: HTTP 200 OK"
    else
        echo -e "${RED}✗${NC} Static file access: HTTP $HTTP_CODE"
    fi
fi
echo ""

# Check for new errors
echo -e "${BLUE}Checking for new errors...${NC}"
sleep 3
NEW_ERRORS=$(sudo tail -20 /var/log/nginx/app3_error.log 2>/dev/null | grep -i "upstream prematurely closed" | wc -l)

if [ "$NEW_ERRORS" -eq 0 ]; then
    echo -e "${GREEN}✓${NC} No new 'upstream prematurely closed' errors"
    echo -e "${GREEN}✓${NC} Fix is working!"
else
    echo -e "${YELLOW}⚠${NC} Found $NEW_ERRORS new errors"
    echo "Recent errors:"
    sudo tail -5 /var/log/nginx/app3_error.log | sed 's/^/   /'
fi
echo ""

echo "=========================================="
echo -e "${GREEN}Restart completed!${NC}"
echo "=========================================="
echo ""
echo "Monitor logs in real-time:"
echo "  sudo tail -f /var/log/nginx/app3_error.log"
echo ""
echo "Test in browser:"
echo "  https://app3.anakhebat.web.id"
echo ""

