#!/bin/bash

# ============================================================================
# Extract Static Files - Simple Version
# ============================================================================
# Extract static files dari container ke host filesystem
# Usage: ./scripts/extract-static-simple.sh
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

CONTAINER_NAME="gowa-app-prod"
STATIC_ROOT="${STATIC_ROOT:-/var/www/gowa_umkm/static}"

if [ "$EUID" -ne 0 ]; then 
    SUDO="sudo"
else
    SUDO=""
fi

echo "=========================================="
echo "  Extract Static Files"
echo "=========================================="
echo "Container: $CONTAINER_NAME"
echo "Static Root: $STATIC_ROOT"
echo ""

# Check container
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo -e "${RED}✗${NC} Container $CONTAINER_NAME is not running"
    exit 1
fi

# Create directory
echo -e "${BLUE}[1/4]${NC} Creating directory..."
$SUDO mkdir -p "$STATIC_ROOT"
$SUDO chown -R $USER:$USER "$STATIC_ROOT"
echo -e "${GREEN}✓${NC} Directory created"
echo ""

# Check source files
echo -e "${BLUE}[2/4]${NC} Checking source files in container..."
if docker exec "$CONTAINER_NAME" test -d "/app/.output/public/_nuxt" 2>/dev/null; then
    FILE_COUNT=$(docker exec "$CONTAINER_NAME" find /app/.output/public/_nuxt -type f 2>/dev/null | wc -l)
    echo -e "${GREEN}✓${NC} Found $FILE_COUNT files in container"
else
    echo -e "${RED}✗${NC} _nuxt directory not found in container"
    exit 1
fi
echo ""

# Extract files
echo -e "${BLUE}[3/4]${NC} Extracting files..."
echo "This may take a moment..."

# Remove old files if exist
if [ -d "$STATIC_ROOT/_nuxt" ]; then
    echo "Removing old files..."
    $SUDO rm -rf "$STATIC_ROOT/_nuxt"
fi

# Extract _nuxt
if docker cp "$CONTAINER_NAME:/app/.output/public/_nuxt" "$STATIC_ROOT/" 2>/dev/null; then
    echo -e "${GREEN}✓${NC} _nuxt directory extracted"
else
    echo -e "${RED}✗${NC} Failed to extract _nuxt"
    exit 1
fi

# Extract builds if exists
if docker exec "$CONTAINER_NAME" test -d "/app/.output/public/builds" 2>/dev/null; then
    if docker cp "$CONTAINER_NAME:/app/.output/public/builds" "$STATIC_ROOT/" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} builds directory extracted"
    fi
fi

echo ""

# Set permissions
echo -e "${BLUE}[4/4]${NC} Setting permissions..."
$SUDO chown -R www-data:www-data "$STATIC_ROOT"
$SUDO chmod -R 755 "$STATIC_ROOT"
echo -e "${GREEN}✓${NC} Permissions set"
echo ""

# Verify
echo -e "${BLUE}Verification:${NC}"
if [ -d "$STATIC_ROOT/_nuxt" ]; then
    DEST_COUNT=$(find "$STATIC_ROOT/_nuxt" -type f | wc -l)
    echo -e "${GREEN}✓${NC} Files extracted: $DEST_COUNT files"
    echo "Location: $STATIC_ROOT/_nuxt"
    echo ""
    echo "Sample files:"
    ls -lh "$STATIC_ROOT/_nuxt" | head -5 | tail -4
else
    echo -e "${RED}✗${NC} Extraction failed"
    exit 1
fi

echo ""
echo "=========================================="
echo -e "${GREEN}Extraction completed!${NC}"
echo "=========================================="
echo ""

