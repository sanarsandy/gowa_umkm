#!/bin/bash

# ============================================================================
# Diagnostic Script untuk 503 Service Unavailable Error
# ============================================================================
# Script ini membantu diagnose masalah 503 pada static files
# Usage: ./scripts/diagnose-503-error.sh
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
echo "  Diagnostic: 503 Service Unavailable"
echo "=========================================="
echo ""

# 1. Check container status
echo -e "${BLUE}[1/8]${NC} Checking container status..."
if docker ps | grep -q "$CONTAINER_NAME"; then
    STATUS=$(docker ps --filter "name=$CONTAINER_NAME" --format "{{.Status}}")
    echo -e "${GREEN}✓${NC} Container is running"
    echo "   Status: $STATUS"
else
    echo -e "${RED}✗${NC} Container is NOT running"
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
    echo -e "${YELLOW}ℹ${NC} Container may not be ready yet"
fi
echo ""

# 3. Check if container is listening on port 3000
echo -e "${BLUE}[3/8]${NC} Checking if container is listening on port 3000..."
if docker exec "$CONTAINER_NAME" netstat -tuln 2>/dev/null | grep -q ":3000" || \
   docker exec "$CONTAINER_NAME" ss -tuln 2>/dev/null | grep -q ":3000"; then
    echo -e "${GREEN}✓${NC} Container is listening on port 3000"
else
    echo -e "${RED}✗${NC} Container is NOT listening on port 3000"
    echo -e "${YELLOW}ℹ${NC} Container may have crashed or not started properly"
fi
echo ""

# 4. Test direct access to container
echo -e "${BLUE}[4/8]${NC} Testing direct access to container..."
if docker exec "$CONTAINER_NAME" wget -q --spider http://localhost:3000 2>/dev/null || \
   docker exec "$CONTAINER_NAME" curl -s -o /dev/null -w "%{http_code}" http://localhost:3000 2>/dev/null | grep -q "200"; then
    HTTP_CODE=$(docker exec "$CONTAINER_NAME" curl -s -o /dev/null -w "%{http_code}" http://localhost:3000 2>/dev/null || echo "000")
    echo -e "${GREEN}✓${NC} Container responds on localhost:3000 (HTTP $HTTP_CODE)"
else
    echo -e "${RED}✗${NC} Container does NOT respond on localhost:3000"
    echo -e "${YELLOW}ℹ${NC} Check container logs: docker logs $CONTAINER_NAME"
fi
echo ""

# 5. Check port mapping
echo -e "${BLUE}[5/8]${NC} Checking port mapping..."
if netstat -tuln 2>/dev/null | grep -q ":3002" || ss -tuln 2>/dev/null | grep -q ":3002"; then
    echo -e "${GREEN}✓${NC} Port 3002 is listening on host"
    PORT_INFO=$(netstat -tuln 2>/dev/null | grep ":3002" || ss -tuln 2>/dev/null | grep ":3002")
    echo "   $PORT_INFO"
else
    echo -e "${RED}✗${NC} Port 3002 is NOT listening on host"
    echo -e "${YELLOW}ℹ${NC} Check docker-compose.prod.yml port mapping"
fi
echo ""

# 6. Test access from host
echo -e "${BLUE}[6/8]${NC} Testing access from host..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 http://127.0.0.1:$NGINX_PORT 2>/dev/null || echo "000")
if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓${NC} Host can access container (HTTP 200)"
elif [ "$HTTP_CODE" = "000" ]; then
    echo -e "${RED}✗${NC} Cannot connect to container (connection refused)"
    echo -e "${YELLOW}ℹ${NC} Container may not be running or port mapping is wrong"
else
    echo -e "${YELLOW}⚠${NC} Host access returned HTTP $HTTP_CODE"
fi
echo ""

# 7. Check static files
echo -e "${BLUE}[7/8]${NC} Checking static files..."
if docker exec "$CONTAINER_NAME" test -d "/app/.output/server/chunks/public/_nuxt" 2>/dev/null; then
    FILE_COUNT=$(docker exec "$CONTAINER_NAME" find /app/.output/server/chunks/public/_nuxt -type f 2>/dev/null | wc -l)
    if [ "$FILE_COUNT" -gt 0 ]; then
        echo -e "${GREEN}✓${NC} Static files found: $FILE_COUNT files"
        if [ "$FILE_COUNT" -lt 40 ]; then
            echo -e "${YELLOW}⚠${NC} File count seems low (expected ~46 files)"
            echo -e "${YELLOW}ℹ${NC} Run: ./scripts/fix-container-static-files.sh"
        fi
    else
        echo -e "${RED}✗${NC} No static files found"
        echo -e "${YELLOW}ℹ${NC} Run: ./scripts/fix-container-static-files.sh"
    fi
else
    echo -e "${RED}✗${NC} Static files directory does not exist"
    echo -e "${YELLOW}ℹ${NC} Run: ./scripts/fix-container-static-files.sh"
fi
echo ""

# 8. Check container logs for errors
echo -e "${BLUE}[8/8]${NC} Checking container logs for errors..."
ERROR_COUNT=$(docker logs "$CONTAINER_NAME" 2>&1 | grep -i "error\|enoent\|503\|500" | tail -10 | wc -l)
if [ "$ERROR_COUNT" -gt 0 ]; then
    echo -e "${YELLOW}⚠${NC} Found errors in container logs:"
    docker logs "$CONTAINER_NAME" 2>&1 | grep -i "error\|enoent\|503\|500" | tail -5 | sed 's/^/   /'
else
    echo -e "${GREEN}✓${NC} No recent errors in container logs"
fi
echo ""

# 9. Check nginx error logs
echo -e "${BLUE}[9/8]${NC} Checking nginx error logs..."
if [ -f "/var/log/nginx/app3_error.log" ]; then
    NGINX_ERRORS=$(sudo tail -20 /var/log/nginx/app3_error.log 2>/dev/null | grep -i "502\|503\|upstream" | wc -l)
    if [ "$NGINX_ERRORS" -gt 0 ]; then
        echo -e "${YELLOW}⚠${NC} Found errors in nginx logs:"
        sudo tail -20 /var/log/nginx/app3_error.log 2>/dev/null | grep -i "502\|503\|upstream" | tail -5 | sed 's/^/   /'
    else
        echo -e "${GREEN}✓${NC} No recent errors in nginx logs"
    fi
else
    echo -e "${YELLOW}⚠${NC} Nginx error log not found"
fi
echo ""

# Summary and recommendations
echo "=========================================="
echo "  Summary & Recommendations"
echo "=========================================="
echo ""

if [ "$HTTP_CODE" != "200" ] || [ "$HEALTH" != "healthy" ]; then
    echo -e "${RED}Problem detected:${NC} Container may not be ready or accessible"
    echo ""
    echo "Recommended actions:"
    echo ""
    echo "1. ${YELLOW}Check container logs${NC}"
    echo "   docker logs -f $CONTAINER_NAME"
    echo ""
    echo "2. ${YELLOW}Restart container${NC}"
    echo "   docker compose -f docker-compose.prod.yml restart app"
    echo ""
    echo "3. ${YELLOW}Wait for health check${NC}"
    echo "   sleep 45"
    echo "   docker inspect --format='{{.State.Health.Status}}' $CONTAINER_NAME"
    echo ""
    echo "4. ${YELLOW}Fix static files if needed${NC}"
    echo "   ./scripts/fix-container-static-files.sh"
    echo ""
elif [ "$FILE_COUNT" -lt 40 ]; then
    echo -e "${YELLOW}Problem detected:${NC} Static files may be incomplete"
    echo ""
    echo "Recommended action:"
    echo "   ./scripts/fix-container-static-files.sh"
    echo ""
else
    echo -e "${GREEN}No obvious issues detected${NC}"
    echo ""
    echo "If problems persist:"
    echo "  - Check nginx config: sudo nginx -t"
    echo "  - Check nginx logs: sudo tail -f /var/log/nginx/app3_error.log"
    echo "  - Check container logs: docker logs -f $CONTAINER_NAME"
fi
echo ""

