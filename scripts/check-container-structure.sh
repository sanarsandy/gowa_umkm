#!/bin/bash

# ============================================================================
# Script untuk Check Struktur File di Container
# ============================================================================
# Script ini membantu diagnose masalah ENOENT dengan melihat struktur file
# Usage: ./scripts/check-container-structure.sh
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
echo "  Check Container File Structure"
echo "=========================================="
echo ""

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo -e "${RED}✗${NC} Container $CONTAINER_NAME is not running"
    exit 1
fi

echo -e "${BLUE}[1/6]${NC} Checking .output directory structure..."
echo ""
docker exec "$CONTAINER_NAME" ls -la /app/.output/ 2>/dev/null | head -20
echo ""

echo -e "${BLUE}[2/6]${NC} Checking .output/public directory..."
echo ""
if docker exec "$CONTAINER_NAME" test -d "/app/.output/public" 2>/dev/null; then
    echo -e "${GREEN}✓${NC} /app/.output/public exists"
    docker exec "$CONTAINER_NAME" ls -la /app/.output/public/ 2>/dev/null | head -10
else
    echo -e "${RED}✗${NC} /app/.output/public does NOT exist"
fi
echo ""

echo -e "${BLUE}[3/6]${NC} Checking .output/public/_nuxt directory..."
echo ""
if docker exec "$CONTAINER_NAME" test -d "/app/.output/public/_nuxt" 2>/dev/null; then
    echo -e "${GREEN}✓${NC} /app/.output/public/_nuxt exists"
    FILE_COUNT=$(docker exec "$CONTAINER_NAME" find /app/.output/public/_nuxt -type f 2>/dev/null | wc -l)
    echo "   Found $FILE_COUNT files"
    echo "   Sample files:"
    docker exec "$CONTAINER_NAME" ls -lh /app/.output/public/_nuxt/ 2>/dev/null | head -5 | tail -4
else
    echo -e "${RED}✗${NC} /app/.output/public/_nuxt does NOT exist"
fi
echo ""

echo -e "${BLUE}[4/6]${NC} Checking .output/server directory..."
echo ""
if docker exec "$CONTAINER_NAME" test -d "/app/.output/server" 2>/dev/null; then
    echo -e "${GREEN}✓${NC} /app/.output/server exists"
    docker exec "$CONTAINER_NAME" ls -la /app/.output/server/ 2>/dev/null | head -10
else
    echo -e "${RED}✗${NC} /app/.output/server does NOT exist"
fi
echo ""

echo -e "${BLUE}[5/6]${NC} Checking .output/server/chunks directory..."
echo ""
if docker exec "$CONTAINER_NAME" test -d "/app/.output/server/chunks" 2>/dev/null; then
    echo -e "${GREEN}✓${NC} /app/.output/server/chunks exists"
    docker exec "$CONTAINER_NAME" ls -la /app/.output/server/chunks/ 2>/dev/null | head -10
else
    echo -e "${RED}✗${NC} /app/.output/server/chunks does NOT exist"
fi
echo ""

echo -e "${BLUE}[6/6]${NC} Checking if Nitro is looking in wrong path..."
echo ""
# Check if the problematic path exists
PROBLEM_PATH="/app/.output/server/chunks/public/_nuxt"
if docker exec "$CONTAINER_NAME" test -d "$PROBLEM_PATH" 2>/dev/null; then
    echo -e "${GREEN}✓${NC} $PROBLEM_PATH exists (this is where Nitro is looking)"
    docker exec "$CONTAINER_NAME" ls -la "$PROBLEM_PATH" 2>/dev/null | head -5
else
    echo -e "${RED}✗${NC} $PROBLEM_PATH does NOT exist"
    echo -e "${YELLOW}ℹ${NC} This is the problem! Nitro is looking for files here but they don't exist"
    echo ""
    echo "Solution: Create symlink or copy files to this location"
fi
echo ""

echo "=========================================="
echo "  Summary"
echo "=========================================="
echo ""
echo "Expected structure:"
echo "  /app/.output/public/_nuxt/*.js  (should exist)"
echo ""
echo "Nitro is looking for:"
echo "  /app/.output/server/chunks/public/_nuxt/*.js  (doesn't exist)"
echo ""
echo "This is why you're getting ENOENT errors!"
echo ""

