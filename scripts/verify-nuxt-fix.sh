#!/bin/bash
# Verification script for Nuxt static asset fix
# Run this script to verify the fix is working correctly

set -e

echo "🔍 Nuxt Static Asset Verification"
echo "=================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if container is running
echo "1. Checking if frontend container is running..."
if docker ps | grep -q "gowa-app-prod"; then
    echo -e "${GREEN}✅ Container is running${NC}"
else
    echo -e "${RED}❌ Container is not running${NC}"
    exit 1
fi

# Check static assets in container
echo ""
echo "2. Checking static assets in container..."
if docker exec gowa-app-prod ls /app/public/_nuxt/ > /dev/null 2>&1; then
    JS_COUNT=$(docker exec gowa-app-prod find /app/public/_nuxt/ -name "*.js" 2>/dev/null | wc -l)
    if [ "$JS_COUNT" -gt 0 ]; then
        echo -e "${GREEN}✅ Found $JS_COUNT JavaScript files${NC}"
    else
        echo -e "${RED}❌ No JavaScript files found${NC}"
        exit 1
    fi
else
    echo -e "${RED}❌ _nuxt directory not found${NC}"
    exit 1
fi

# Check file permissions
echo ""
echo "3. Checking file permissions..."
PERMS=$(docker exec gowa-app-prod stat -c %a /app/public 2>/dev/null || echo "000")
if [ "$PERMS" = "755" ] || [ "$PERMS" = "775" ]; then
    echo -e "${GREEN}✅ Permissions are correct ($PERMS)${NC}"
else
    echo -e "${YELLOW}⚠️  Permissions may need adjustment ($PERMS)${NC}"
fi

# Test Nuxt server response
echo ""
echo "4. Testing Nuxt server response..."
if docker exec gowa-app-prod wget -q -O /dev/null http://localhost:3000/ 2>/dev/null; then
    echo -e "${GREEN}✅ Nuxt server is responding${NC}"
else
    echo -e "${RED}❌ Nuxt server is not responding${NC}"
    echo "Checking logs..."
    docker logs gowa-app-prod --tail 20
    exit 1
fi

# Check if port 3002 is accessible from host
echo ""
echo "5. Testing port accessibility from host..."
if curl -s -o /dev/null -w "%{http_code}" http://localhost:3002/ | grep -q "200\|301\|302"; then
    echo -e "${GREEN}✅ Port 3002 is accessible${NC}"
else
    echo -e "${YELLOW}⚠️  Port 3002 may not be accessible (check docker-compose port mapping)${NC}"
fi

# Check container logs for errors
echo ""
echo "6. Checking container logs for errors..."
ERROR_COUNT=$(docker logs gowa-app-prod --tail 50 2>&1 | grep -i "error\|fail\|exception" | wc -l)
if [ "$ERROR_COUNT" -eq 0 ]; then
    echo -e "${GREEN}✅ No errors in recent logs${NC}"
else
    echo -e "${YELLOW}⚠️  Found $ERROR_COUNT potential errors in logs${NC}"
    echo "Recent errors:"
    docker logs gowa-app-prod --tail 50 2>&1 | grep -i "error\|fail\|exception" | tail -5
fi

echo ""
echo "=================================="
echo -e "${GREEN}✅ Verification complete!${NC}"
echo ""
echo "Summary:"
echo "  - Container: Running"
echo "  - Static assets: $JS_COUNT files found"
echo "  - Nuxt server: Responding"
echo "  - Recent errors: $ERROR_COUNT"
echo ""
echo "Next steps:"
echo "1. If all checks passed, update Nginx configuration on your VPS"
echo "2. Test the application in browser: https://app2.anakhebat.web.id/"
echo "3. Check browser console for any remaining 503 errors"
echo ""
echo "To view full logs:"
echo "  docker logs -f gowa-app-prod"
