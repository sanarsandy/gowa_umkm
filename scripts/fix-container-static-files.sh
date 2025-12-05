#!/bin/bash

# ============================================================================
# Script untuk Fix Static Files di Container yang Sudah Running
# ============================================================================
# Script ini akan copy files ke lokasi yang dicari Nitro di container yang sudah running
# Usage: ./scripts/fix-container-static-files.sh
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
echo "  Fix Static Files in Running Container"
echo "=========================================="
echo ""

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo -e "${RED}✗${NC} Container $CONTAINER_NAME is not running"
    exit 1
fi

echo -e "${BLUE}[1/4]${NC} Checking source files..."
SRC_COUNT=$(docker exec "$CONTAINER_NAME" find /app/.output/public/_nuxt -type f 2>/dev/null | wc -l)
if [ "$SRC_COUNT" -eq 0 ]; then
    echo -e "${RED}✗${NC} No source files found in /app/.output/public/_nuxt"
    exit 1
fi
echo -e "${GREEN}✓${NC} Found $SRC_COUNT source files"
echo ""

echo -e "${BLUE}[2/4]${NC} Creating destination directory..."
docker exec -u root "$CONTAINER_NAME" mkdir -p /app/.output/server/chunks/public
echo -e "${GREEN}✓${NC} Directory created"
echo ""

echo -e "${BLUE}[3/4]${NC} Copying files..."
# Copy _nuxt directory
docker exec -u root "$CONTAINER_NAME" cp -r /app/.output/public/_nuxt /app/.output/server/chunks/public/ 2>/dev/null || {
    echo -e "${YELLOW}⚠${NC} Direct copy failed, trying with tar..."
    docker exec "$CONTAINER_NAME" tar -C /app/.output/public -cf - _nuxt | \
    docker exec -u root -i "$CONTAINER_NAME" tar -C /app/.output/server/chunks/public -xf -
}

# Copy builds directory if exists
if docker exec "$CONTAINER_NAME" test -d "/app/.output/public/builds" 2>/dev/null; then
    docker exec -u root "$CONTAINER_NAME" cp -r /app/.output/public/builds /app/.output/server/chunks/public/ 2>/dev/null || \
    docker exec "$CONTAINER_NAME" tar -C /app/.output/public -cf - builds | \
    docker exec -u root -i "$CONTAINER_NAME" tar -C /app/.output/server/chunks/public -xf -
fi

echo -e "${GREEN}✓${NC} Files copied"
echo ""

echo -e "${BLUE}[4/4]${NC} Setting permissions..."
docker exec -u root "$CONTAINER_NAME" chown -R nuxtjs:nodejs /app/.output/server/chunks/public
echo -e "${GREEN}✓${NC} Permissions set"
echo ""

# Verify
echo -e "${BLUE}Verification:${NC}"
DEST_COUNT=$(docker exec "$CONTAINER_NAME" find /app/.output/server/chunks/public/_nuxt -type f 2>/dev/null | wc -l)
echo "Source files: $SRC_COUNT"
echo "Destination files: $DEST_COUNT"

if [ "$SRC_COUNT" -eq "$DEST_COUNT" ]; then
    echo -e "${GREEN}✓${NC} All files copied successfully!"
else
    echo -e "${YELLOW}⚠${NC} File count mismatch. Expected $SRC_COUNT, got $DEST_COUNT"
    echo -e "${YELLOW}ℹ${NC} You may need to rebuild the container"
fi
echo ""

echo "=========================================="
echo -e "${GREEN}Fix completed!${NC}"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Test static file access: curl -I http://127.0.0.1:3002/_nuxt/entry.js"
echo "2. Check logs: docker logs gowa-app-prod | grep -i enoent"
echo "3. If still issues, rebuild container: docker compose -f docker-compose.prod.yml build --no-cache app"
echo ""

