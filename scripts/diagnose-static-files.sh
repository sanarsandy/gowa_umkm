#!/bin/bash

# ============================================================================
# Diagnostic Script untuk Troubleshooting Static Files 502/503 Error
# ============================================================================
# Script ini membantu diagnose masalah static files di production
# Usage: ./scripts/diagnose-static-files.sh
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

CONTAINER_NAME="gowa-app-prod"
NGINX_PORT="3002"

echo "=========================================="
echo "  Diagnostic: Static Files 502/503 Error"
echo "=========================================="
echo ""

# 1. Check if container is running
echo -e "${BLUE}[1/8]${NC} Checking container status..."
if docker ps | grep -q "$CONTAINER_NAME"; then
    echo -e "${GREEN}✓${NC} Container $CONTAINER_NAME is running"
    CONTAINER_STATUS=$(docker ps --filter "name=$CONTAINER_NAME" --format "{{.Status}}")
    echo "   Status: $CONTAINER_STATUS"
else
    echo -e "${RED}✗${NC} Container $CONTAINER_NAME is NOT running"
    echo -e "${YELLOW}ℹ${NC} Start container: docker compose -f docker-compose.prod.yml up -d app"
    exit 1
fi
echo ""

# 2. Check container health
echo -e "${BLUE}[2/8]${NC} Checking container health..."
HEALTH=$(docker inspect --format='{{.State.Health.Status}}' "$CONTAINER_NAME" 2>/dev/null || echo "no-healthcheck")
if [ "$HEALTH" = "healthy" ]; then
    echo -e "${GREEN}✓${NC} Container is healthy"
elif [ "$HEALTH" = "no-healthcheck" ]; then
    echo -e "${YELLOW}⚠${NC} No health check configured"
else
    echo -e "${RED}✗${NC} Container health: $HEALTH"
fi
echo ""

# 3. Check if port is accessible
echo -e "${BLUE}[3/8]${NC} Checking port accessibility..."
if curl -s -o /dev/null -w "%{http_code}" --max-time 5 http://127.0.0.1:$NGINX_PORT > /dev/null 2>&1; then
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 http://127.0.0.1:$NGINX_PORT)
    echo -e "${GREEN}✓${NC} Port $NGINX_PORT is accessible (HTTP $HTTP_CODE)"
else
    echo -e "${RED}✗${NC} Port $NGINX_PORT is NOT accessible"
    echo -e "${YELLOW}ℹ${NC} Check if container is listening on port 3000"
fi
echo ""

# 4. Check static files in container
echo -e "${BLUE}[4/8]${NC} Checking static files in container..."
if docker exec "$CONTAINER_NAME" test -d "/app/.output/public/_nuxt" 2>/dev/null; then
    FILE_COUNT=$(docker exec "$CONTAINER_NAME" find /app/.output/public/_nuxt -type f 2>/dev/null | wc -l)
    echo -e "${GREEN}✓${NC} Static files directory exists"
    echo "   Found $FILE_COUNT files in _nuxt directory"
    
    # List first few files
    echo "   Sample files:"
    docker exec "$CONTAINER_NAME" ls -lh /app/.output/public/_nuxt/ 2>/dev/null | head -5 | tail -4
else
    echo -e "${RED}✗${NC} Static files directory NOT found"
    echo -e "${YELLOW}ℹ${NC} Check Dockerfile.prod build process"
fi
echo ""

# 5. Test static file access from container
echo -e "${BLUE}[5/8]${NC} Testing static file access from container..."
FIRST_FILE=$(docker exec "$CONTAINER_NAME" find /app/.output/public/_nuxt -type f -name "*.js" 2>/dev/null | head -1)
if [ -n "$FIRST_FILE" ]; then
    FILE_NAME=$(basename "$FIRST_FILE")
    if docker exec "$CONTAINER_NAME" test -f "$FIRST_FILE" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} Can access static file: $FILE_NAME"
        FILE_SIZE=$(docker exec "$CONTAINER_NAME" stat -c%s "$FIRST_FILE" 2>/dev/null || docker exec "$CONTAINER_NAME" stat -f%z "$FIRST_FILE" 2>/dev/null)
        echo "   File size: $FILE_SIZE bytes"
    else
        echo -e "${RED}✗${NC} Cannot access static file: $FILE_NAME"
    fi
else
    echo -e "${YELLOW}⚠${NC} No static files found to test"
fi
echo ""

# 6. Test HTTP access to static files
echo -e "${BLUE}[6/8]${NC} Testing HTTP access to static files..."
if [ -n "$FIRST_FILE" ]; then
    FILE_NAME=$(basename "$FIRST_FILE")
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 "http://127.0.0.1:$NGINX_PORT/_nuxt/$FILE_NAME" 2>/dev/null || echo "000")
    if [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓${NC} HTTP access successful (200 OK)"
    elif [ "$HTTP_CODE" = "502" ] || [ "$HTTP_CODE" = "503" ]; then
        echo -e "${RED}✗${NC} HTTP access failed: $HTTP_CODE (Bad Gateway/Service Unavailable)"
        echo -e "${YELLOW}ℹ${NC} This is the error you're experiencing"
    else
        echo -e "${YELLOW}⚠${NC} HTTP access returned: $HTTP_CODE"
    fi
else
    echo -e "${YELLOW}⚠${NC} Cannot test HTTP access (no files found)"
fi
echo ""

# 7. Check container logs
echo -e "${BLUE}[7/8]${NC} Checking recent container logs..."
echo "   Last 10 lines of container logs:"
docker logs --tail 10 "$CONTAINER_NAME" 2>&1 | sed 's/^/   /'
echo ""

# 8. Check nginx error logs
echo -e "${BLUE}[8/8]${NC} Checking nginx error logs..."
if [ -f "/var/log/nginx/app3_error.log" ]; then
    echo "   Recent nginx errors:"
    sudo tail -10 /var/log/nginx/app3_error.log 2>/dev/null | sed 's/^/   /' || echo "   (Cannot read nginx logs - may need sudo)"
else
    echo -e "${YELLOW}⚠${NC} Nginx error log not found"
fi
echo ""

# Summary and recommendations
echo "=========================================="
echo "  Summary & Recommendations"
echo "=========================================="
echo ""

if [ "$HTTP_CODE" = "502" ] || [ "$HTTP_CODE" = "503" ]; then
    echo -e "${RED}Problem detected:${NC} Static files returning 502/503"
    echo ""
    echo "Possible causes and solutions:"
    echo ""
    echo "1. ${YELLOW}Container not ready${NC}"
    echo "   Solution: Wait for container to be healthy"
    echo "   Check: docker ps | grep $CONTAINER_NAME"
    echo ""
    echo "2. ${YELLOW}Port binding issue${NC}"
    echo "   Solution: Check docker-compose.prod.yml port mapping"
    echo "   Should be: 127.0.0.1:3002:3000"
    echo ""
    echo "3. ${YELLOW}Static files not built${NC}"
    echo "   Solution: Rebuild container"
    echo "   Run: docker compose -f docker-compose.prod.yml build app"
    echo ""
    echo "4. ${YELLOW}Nginx timeout too short${NC}"
    echo "   Solution: Check nginx.conf.production timeout settings"
    echo "   Should have: proxy_connect_timeout, proxy_send_timeout, proxy_read_timeout"
    echo ""
    echo "5. ${YELLOW}Container crash${NC}"
    echo "   Solution: Check container logs"
    echo "   Run: docker logs $CONTAINER_NAME"
    echo ""
    echo "Quick fix commands:"
    echo "  docker compose -f docker-compose.prod.yml restart app"
    echo "  docker compose -f docker-compose.prod.yml logs -f app"
    echo "  sudo nginx -t && sudo systemctl reload nginx"
else
    echo -e "${GREEN}No obvious issues detected${NC}"
    echo "If problems persist, check:"
    echo "  - Container logs: docker logs $CONTAINER_NAME"
    echo "  - Nginx logs: sudo tail -f /var/log/nginx/app3_error.log"
    echo "  - Network connectivity: curl -v http://127.0.0.1:$NGINX_PORT"
fi
echo ""

