#!/bin/bash

# ============================================================================
# Test Static Files Access
# ============================================================================
# Test apakah static files bisa diakses dengan benar
# Usage: ./scripts/test-static-access.sh
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

STATIC_ROOT="${STATIC_ROOT:-/var/www/gowa_umkm/static}"
DOMAIN="app3.anakhebat.web.id"

echo "=========================================="
echo "  Test Static Files Access"
echo "=========================================="
echo ""

# Test 1: Check filesystem
echo -e "${BLUE}[1/5]${NC} Checking filesystem..."
if [ -d "$STATIC_ROOT/_nuxt" ]; then
    FILE_COUNT=$(find "$STATIC_ROOT/_nuxt" -type f | wc -l)
    echo -e "${GREEN}✓${NC} Found $FILE_COUNT files in $STATIC_ROOT/_nuxt"
    
    FIRST_FILE=$(find "$STATIC_ROOT/_nuxt" -type f -name "*.js" | head -1)
    if [ -n "$FIRST_FILE" ]; then
        FILENAME=$(basename "$FIRST_FILE")
        echo -e "${GREEN}✓${NC} Sample file: $FILENAME"
    fi
else
    echo -e "${RED}✗${NC} Directory not found: $STATIC_ROOT/_nuxt"
    exit 1
fi
echo ""

# Test 2: Check permissions
echo -e "${BLUE}[2/5]${NC} Checking permissions..."
if [ -r "$STATIC_ROOT/_nuxt" ]; then
    echo -e "${GREEN}✓${NC} Directory is readable"
else
    echo -e "${RED}✗${NC} Directory is not readable"
fi

if [ -r "$FIRST_FILE" ]; then
    echo -e "${GREEN}✓${NC} File is readable"
else
    echo -e "${RED}✗${NC} File is not readable"
fi
echo ""

# Test 3: Check nginx config
echo -e "${BLUE}[3/5]${NC} Checking nginx configuration..."
NGINX_CONFIG="/etc/nginx/sites-available/$DOMAIN"
if sudo grep -q "alias.*_nuxt" "$NGINX_CONFIG" 2>/dev/null; then
    echo -e "${GREEN}✓${NC} Nginx uses direct file serving (alias)"
    sudo grep "alias.*_nuxt" "$NGINX_CONFIG" | head -1
else
    echo -e "${RED}✗${NC} Nginx might not be using direct file serving"
fi
echo ""

# Test 4: Test HTTP access (localhost)
echo -e "${BLUE}[4/5]${NC} Testing HTTP access (localhost)..."
if [ -n "$FIRST_FILE" ]; then
    HTTP_CODE=$(curl -sL -o /dev/null -w "%{http_code}" --max-time 5 "http://127.0.0.1/_nuxt/$FILENAME" 2>/dev/null || echo "000")
    if [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓${NC} HTTP access: HTTP 200 OK"
    elif [ "$HTTP_CODE" = "301" ] || [ "$HTTP_CODE" = "302" ]; then
        echo -e "${YELLOW}⚠${NC} HTTP redirects (normal if HTTPS is configured)"
    else
        echo -e "${RED}✗${NC} HTTP access failed: HTTP $HTTP_CODE"
    fi
fi
echo ""

# Test 5: Test HTTPS access (domain)
echo -e "${BLUE}[5/5]${NC} Testing HTTPS access (domain)..."
if [ -n "$FIRST_FILE" ]; then
    HTTPS_CODE=$(curl -sLk -o /dev/null -w "%{http_code}" --max-time 10 "https://$DOMAIN/_nuxt/$FILENAME" 2>/dev/null || echo "000")
    if [ "$HTTPS_CODE" = "200" ]; then
        echo -e "${GREEN}✓${NC} HTTPS access: HTTP 200 OK"
        echo -e "${GREEN}✓${NC} Static file is being served correctly!"
    else
        echo -e "${RED}✗${NC} HTTPS access failed: HTTP $HTTPS_CODE"
        echo "   Check nginx logs: sudo tail -f /var/log/nginx/app3_error.log"
    fi
fi
echo ""

echo "=========================================="
echo -e "${GREEN}Test completed!${NC}"
echo "=========================================="
echo ""
echo "If all tests passed, static files are working correctly."
echo "If not, check:"
echo "  1. Static files exist: ls -la $STATIC_ROOT/_nuxt/"
echo "  2. Nginx config: sudo nginx -t"
echo "  3. Nginx logs: sudo tail -f /var/log/nginx/app3_error.log"
echo ""

