#!/bin/bash

# ============================================================================
# FIX ALL - Rebuild dan Fix Static Files
# ============================================================================
# Script lengkap untuk fix semua masalah static files
# Usage: ./scripts/fix-all-static-files.sh
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
echo "  FIX ALL - Static Files"
echo "=========================================="
echo ""

# 1. Stop container
echo -e "${BLUE}[1/5]${NC} Stopping container..."
cd "$PROJECT_DIR"
docker compose -f docker-compose.prod.yml stop app
echo -e "${GREEN}✓${NC} Container stopped"
echo ""

# 2. Rebuild dengan no-cache
echo -e "${BLUE}[2/5]${NC} Rebuilding container (this will take a few minutes)..."
docker compose -f docker-compose.prod.yml build --no-cache app
echo -e "${GREEN}✓${NC} Container rebuilt"
echo ""

# 3. Start container
echo -e "${BLUE}[3/5]${NC} Starting container..."
docker compose -f docker-compose.prod.yml up -d app
echo -e "${GREEN}✓${NC} Container started"
echo ""

# 4. Wait for health check
echo -e "${BLUE}[4/5]${NC} Waiting for container to be healthy..."
sleep 45
HEALTH=$(docker inspect --format='{{.State.Health.Status}}' gowa-app-prod 2>/dev/null || echo "starting")
if [ "$HEALTH" = "healthy" ]; then
    echo -e "${GREEN}✓${NC} Container is healthy"
else
    echo -e "${YELLOW}⚠${NC} Container health: $HEALTH (continuing anyway)"
fi
echo ""

# 5. Fix static files
echo -e "${BLUE}[5/5]${NC} Fixing static files..."
if [ -f "$PROJECT_DIR/scripts/fix-container-static-files.sh" ]; then
    chmod +x "$PROJECT_DIR/scripts/fix-container-static-files.sh"
    "$PROJECT_DIR/scripts/fix-container-static-files.sh"
else
    echo -e "${YELLOW}⚠${NC} fix-container-static-files.sh not found, skipping..."
fi
echo ""

# Verify
echo -e "${BLUE}Verification:${NC}"
echo ""

# Check file count
FILE_COUNT=$(docker exec gowa-app-prod find /app/.output/server/chunks/public/_nuxt -type f 2>/dev/null | wc -l)
echo "Static files in container: $FILE_COUNT"

# Test access
FIRST_FILE=$(docker exec gowa-app-prod find /app/.output/server/chunks/public/_nuxt -type f -name "*.js" 2>/dev/null | head -1)
if [ -n "$FIRST_FILE" ]; then
    FILENAME=$(basename "$FIRST_FILE")
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 10 "http://127.0.0.1:3002/_nuxt/$FILENAME" 2>/dev/null || echo "000")
    if [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓${NC} HTTP access: OK (200)"
    else
        echo -e "${RED}✗${NC} HTTP access: $HTTP_CODE"
    fi
fi

echo ""
echo "=========================================="
echo -e "${GREEN}Fix completed!${NC}"
echo "=========================================="
echo ""
echo "IMPORTANT: Clear browser cache or use incognito mode!"
echo "The HTML page now references the correct files."
echo ""
echo "Test: https://app3.anakhebat.web.id"
echo ""

