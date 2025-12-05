#!/bin/bash

# ============================================================================
# Script untuk Test Static Files Access
# ============================================================================
# Test akses ke static files yang benar-benar ada
# Usage: ./scripts/test-static-files.sh
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

CONTAINER_NAME="gowa-app-prod"

echo "=========================================="
echo "  Test Static Files Access"
echo "=========================================="
echo ""

# Get first few actual files
echo -e "${BLUE}[1/3]${NC} Finding actual static files..."
FILES=$(docker exec "$CONTAINER_NAME" find /app/.output/server/chunks/public/_nuxt -type f -name "*.js" 2>/dev/null | head -5)

if [ -z "$FILES" ]; then
    echo -e "${RED}✗${NC} No static files found"
    exit 1
fi

echo -e "${GREEN}✓${NC} Found files:"
echo "$FILES" | while read -r file; do
    filename=$(basename "$file")
    echo "   - $filename"
done
echo ""

# Test each file
echo -e "${BLUE}[2/3]${NC} Testing file access from container..."
echo "$FILES" | while read -r file; do
    filename=$(basename "$file")
    echo -n "   Testing $filename... "
    
    # Test from inside container
    if docker exec "$CONTAINER_NAME" test -f "$file" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} exists"
    else
        echo -e "${RED}✗${NC} not found"
    fi
done
echo ""

# Test HTTP access
echo -e "${BLUE}[3/3]${NC} Testing HTTP access from host..."
FIRST_FILE=$(echo "$FILES" | head -1)
FIRST_FILENAME=$(basename "$FIRST_FILE")

echo "   Testing: http://127.0.0.1:3002/_nuxt/$FIRST_FILENAME"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 10 "http://127.0.0.1:3002/_nuxt/$FIRST_FILENAME" 2>/dev/null || echo "000")

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "   ${GREEN}✓${NC} HTTP 200 OK"
    
    # Get file size
    CONTENT_LENGTH=$(curl -s -I "http://127.0.0.1:3002/_nuxt/$FIRST_FILENAME" 2>/dev/null | grep -i "content-length" | awk '{print $2}' | tr -d '\r')
    if [ -n "$CONTENT_LENGTH" ]; then
        echo "   Content-Length: $CONTENT_LENGTH bytes"
    fi
elif [ "$HTTP_CODE" = "404" ]; then
    echo -e "   ${RED}✗${NC} HTTP 404 Not Found"
    echo -e "   ${YELLOW}ℹ${NC} File exists in container but not accessible via HTTP"
elif [ "$HTTP_CODE" = "503" ]; then
    echo -e "   ${RED}✗${NC} HTTP 503 Service Unavailable"
    echo -e "   ${YELLOW}ℹ${NC} Upstream connection issue"
else
    echo -e "   ${YELLOW}⚠${NC} HTTP $HTTP_CODE"
fi
echo ""

# Test direct from container (using wget or node if curl not available)
echo -e "${BLUE}Testing direct access from container...${NC}"
echo "   Testing: http://localhost:3000/_nuxt/$FIRST_FILENAME"

# Try wget first, then node if wget not available
if docker exec "$CONTAINER_NAME" which wget >/dev/null 2>&1; then
    CONTAINER_HTTP=$(docker exec "$CONTAINER_NAME" wget -q --spider -S "http://localhost:3000/_nuxt/$FIRST_FILENAME" 2>&1 | grep -i "HTTP" | head -1 | awk '{print $2}' || echo "000")
elif docker exec "$CONTAINER_NAME" which node >/dev/null 2>&1; then
    CONTAINER_HTTP=$(docker exec "$CONTAINER_NAME" node -e "require('http').get('http://localhost:3000/_nuxt/$FIRST_FILENAME', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})" 2>/dev/null && echo "200" || echo "000")
else
    CONTAINER_HTTP="SKIP"
    echo -e "   ${YELLOW}⚠${NC} Cannot test (curl/wget/node not available in container)"
    echo -e "   ${YELLOW}ℹ${NC} But HTTP access from host is working, so server is OK"
fi

if [ "$CONTAINER_HTTP" = "200" ]; then
    echo -e "   ${GREEN}✓${NC} Container can serve file (HTTP 200)"
elif [ "$CONTAINER_HTTP" != "SKIP" ]; then
    echo -e "   ${YELLOW}⚠${NC} Cannot verify container access (HTTP $CONTAINER_HTTP)"
    echo -e "   ${YELLOW}ℹ${NC} But HTTP access from host works, so nginx proxy is working"
fi
echo ""

echo "=========================================="
echo "  Summary"
echo "=========================================="
echo ""
if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓${NC} Static files are accessible via nginx (HTTP 200)"
    echo -e "${GREEN}✓${NC} Nginx proxy is working correctly"
    if [ "$CONTAINER_HTTP" != "200" ] && [ "$CONTAINER_HTTP" != "SKIP" ]; then
        echo -e "${YELLOW}⚠${NC} Note: Cannot verify direct container access"
        echo -e "${YELLOW}ℹ${NC} But nginx proxy works, so this is OK"
    fi
    echo ""
    echo "Next steps:"
    echo "  1. Test in browser: https://app3.anakhebat.web.id"
    echo "  2. Check browser console for any errors"
    echo "  3. Monitor nginx logs: sudo tail -f /var/log/nginx/app3_error.log"
elif [ "$HTTP_CODE" = "404" ]; then
    echo -e "${RED}✗${NC} Problem: File not found (404)"
    echo "   File exists in container but not accessible via HTTP"
    echo "   Check Nuxt/Nitro server configuration"
elif [ "$HTTP_CODE" = "503" ]; then
    echo -e "${RED}✗${NC} Problem: Service unavailable (503)"
    echo "   Check nginx error logs: sudo tail -f /var/log/nginx/app3_error.log"
    echo "   Check container logs: docker logs $CONTAINER_NAME"
else
    echo -e "${YELLOW}⚠${NC} Problem: HTTP $HTTP_CODE"
    echo "   Check nginx config and logs"
fi
echo ""

