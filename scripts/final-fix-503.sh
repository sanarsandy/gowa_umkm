#!/bin/bash

# ============================================================================
# FINAL FIX - Solusi Lengkap untuk 503 Error
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "=========================================="
echo "  FINAL FIX - 503 Error"
echo "=========================================="
echo ""

# 1. Check container logs untuk error sebenarnya
echo -e "${BLUE}[1/6]${NC} Checking container logs for errors..."
echo ""
docker logs gowa-app-prod 2>&1 | grep -i "enoent\|error\|503\|500" | tail -10 | sed 's/^/   /' || echo "   No obvious errors in logs"
echo ""

# 2. Check apakah file yang diminta browser ada
echo -e "${BLUE}[2/6]${NC} Checking if requested files exist..."
MISSING_FILES=("DSzMySIv.js" "BcxIJ4tb.js" "Czrdq3ge.js" "CqpLj_Ac.js" "yOviaWq-.js" "BK-TQHNH.js" "DlAUqK2U.js" "C6qEwHNt.js")

ALL_EXIST=true
for file in "${MISSING_FILES[@]}"; do
    if ! docker exec gowa-app-prod test -f "/app/.output/server/chunks/public/_nuxt/$file" 2>/dev/null && \
       ! docker exec gowa-app-prod test -f "/app/.output/public/_nuxt/$file" 2>/dev/null; then
        echo "   ✗ $file - NOT FOUND"
        ALL_EXIST=false
    fi
done

if [ "$ALL_EXIST" = false ]; then
    echo ""
    echo -e "${RED}PROBLEM FOUND:${NC} Files don't exist in container"
    echo -e "${YELLOW}Solution:${NC} Rebuild container to regenerate HTML with correct file references"
    echo ""
    echo "Rebuilding container..."
    cd "$PROJECT_DIR"
    docker compose -f docker-compose.prod.yml stop app
    docker compose -f docker-compose.prod.yml build --no-cache app
    docker compose -f docker-compose.prod.yml up -d app
    sleep 45
    echo -e "${GREEN}✓${NC} Container rebuilt"
    echo ""
fi
echo ""

# 3. Fix static files (copy semua file)
echo -e "${BLUE}[3/6]${NC} Ensuring all static files are in correct location..."
docker exec -u root gowa-app-prod sh -c "
    if [ -d /app/.output/public/_nuxt ] && [ -d /app/.output/server/chunks ]; then
        mkdir -p /app/.output/server/chunks/public
        if [ -d /app/.output/public/_nuxt ]; then
            cp -r /app/.output/public/_nuxt /app/.output/server/chunks/public/ 2>/dev/null || true
        fi
        if [ -d /app/.output/public/builds ]; then
            cp -r /app/.output/public/builds /app/.output/server/chunks/public/ 2>/dev/null || true
        fi
        chown -R nuxtjs:nodejs /app/.output/server/chunks/public
    fi
" 2>/dev/null || true
echo -e "${GREEN}✓${NC} Static files fixed"
echo ""

# 4. Restart container
echo -e "${BLUE}[4/6]${NC} Restarting container..."
cd "$PROJECT_DIR"
docker compose -f docker-compose.prod.yml restart app
sleep 10
echo -e "${GREEN}✓${NC} Container restarted"
echo ""

# 5. Clear nginx error log dan restart
echo -e "${BLUE}[5/6]${NC} Restarting nginx..."
sudo truncate -s 0 /var/log/nginx/app3_error.log 2>/dev/null || true
sudo systemctl restart nginx
sleep 2
echo -e "${GREEN}✓${NC} Nginx restarted"
echo ""

# 6. Test
echo -e "${BLUE}[6/6]${NC} Testing..."
sleep 3

# Get actual file from container
ACTUAL_FILE=$(docker exec gowa-app-prod find /app/.output/server/chunks/public/_nuxt -type f -name "*.js" 2>/dev/null | head -1)
if [ -n "$ACTUAL_FILE" ]; then
    FILENAME=$(basename "$ACTUAL_FILE")
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 10 "http://127.0.0.1:3002/_nuxt/$FILENAME" 2>/dev/null || echo "000")
    
    if [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓${NC} HTTP access: 200 OK"
    else
        echo -e "${RED}✗${NC} HTTP access: $HTTP_CODE"
    fi
fi

# Check for new errors
NEW_ERRORS=$(sudo tail -10 /var/log/nginx/app3_error.log 2>/dev/null | grep -i "upstream prematurely closed\|503" | wc -l)
if [ "$NEW_ERRORS" -eq 0 ]; then
    echo -e "${GREEN}✓${NC} No new errors in nginx logs"
else
    echo -e "${YELLOW}⚠${NC} Found $NEW_ERRORS new errors"
fi

echo ""
echo "=========================================="
echo -e "${GREEN}Fix completed!${NC}"
echo "=========================================="
echo ""
echo "IMPORTANT:"
echo "  1. Clear browser cache or use incognito mode"
echo "  2. Test: https://app3.anakhebat.web.id"
echo "  3. If still errors, check: docker logs gowa-app-prod"
echo ""

