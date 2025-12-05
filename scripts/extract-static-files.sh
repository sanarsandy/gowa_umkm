#!/bin/bash

# ============================================================================
# Script untuk Extract Static Files dari Container ke Host
# ============================================================================
# Script ini mengekstrak static files (_nuxt/*) dari container Nuxt
# ke direktori host untuk di-serve langsung oleh Nginx
#
# Usage: ./scripts/extract-static-files.sh
# ============================================================================

set -e  # Exit on error

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CONTAINER_NAME="gowa-app-prod"
STATIC_ROOT="${STATIC_ROOT:-/var/www/gowa_umkm/static}"
CONTAINER_STATIC_PATH="/app/.output/public"

echo "=========================================="
echo "  Extract Static Files from Container"
echo "=========================================="
echo ""

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo -e "${RED}✗${NC} Container $CONTAINER_NAME is not running"
    echo -e "${YELLOW}ℹ${NC} Starting container first..."
    docker compose -f docker-compose.prod.yml up -d app
    sleep 5
fi

# Check if container exists
if ! docker ps -a | grep -q "$CONTAINER_NAME"; then
    echo -e "${RED}✗${NC} Container $CONTAINER_NAME does not exist"
    echo -e "${YELLOW}ℹ${NC} Please build the container first:"
    echo "   docker compose -f docker-compose.prod.yml build app"
    exit 1
fi

# Create static root directory
echo -e "${BLUE}[1/4]${NC} Creating static root directory..."
sudo mkdir -p "$STATIC_ROOT"
sudo chown -R $USER:$USER "$STATIC_ROOT"
echo -e "${GREEN}✓${NC} Directory created: $STATIC_ROOT"
echo ""

# Check if static files exist in container
echo -e "${BLUE}[2/4]${NC} Checking static files in container..."
if docker exec "$CONTAINER_NAME" test -d "$CONTAINER_STATIC_PATH/_nuxt"; then
    echo -e "${GREEN}✓${NC} Static files found in container"
else
    echo -e "${YELLOW}⚠${NC} Static files not found in container at $CONTAINER_STATIC_PATH/_nuxt"
    echo -e "${YELLOW}ℹ${NC} Listing container directory structure:"
    docker exec "$CONTAINER_NAME" ls -la "$CONTAINER_STATIC_PATH" || true
    echo ""
    echo -e "${YELLOW}ℹ${NC} Trying alternative path..."
    # Try to find _nuxt directory
    docker exec "$CONTAINER_NAME" find /app -name "_nuxt" -type d 2>/dev/null || true
    exit 1
fi
echo ""

# Extract static files
echo -e "${BLUE}[3/4]${NC} Extracting static files from container..."
echo -e "${YELLOW}ℹ${NC} This may take a moment..."

# Extract _nuxt directory
if docker cp "$CONTAINER_NAME:$CONTAINER_STATIC_PATH/_nuxt" "$STATIC_ROOT/" 2>/dev/null; then
    echo -e "${GREEN}✓${NC} _nuxt directory extracted"
else
    echo -e "${RED}✗${NC} Failed to extract _nuxt directory"
    exit 1
fi

# Extract other static files if they exist
if docker exec "$CONTAINER_NAME" test -d "$CONTAINER_STATIC_PATH"; then
    # Copy all files from public directory (excluding _nuxt which we already copied)
    docker exec "$CONTAINER_NAME" find "$CONTAINER_STATIC_PATH" -maxdepth 1 -type f -exec cp {} "$STATIC_ROOT/" \; 2>/dev/null || true
    docker exec "$CONTAINER_NAME" find "$CONTAINER_STATIC_PATH" -mindepth 1 -maxdepth 1 -type d ! -name "_nuxt" -exec cp -r {} "$STATIC_ROOT/" \; 2>/dev/null || true
fi

echo ""

# Set proper permissions
echo -e "${BLUE}[4/4]${NC} Setting file permissions..."
sudo chown -R www-data:www-data "$STATIC_ROOT"
sudo chmod -R 755 "$STATIC_ROOT"
echo -e "${GREEN}✓${NC} Permissions set"
echo ""

# Verify extraction
echo -e "${BLUE}Verification:${NC}"
if [ -d "$STATIC_ROOT/_nuxt" ]; then
    FILE_COUNT=$(find "$STATIC_ROOT/_nuxt" -type f | wc -l)
    echo -e "${GREEN}✓${NC} Static files extracted successfully"
    echo -e "${GREEN}✓${NC} Found $FILE_COUNT files in _nuxt directory"
    echo ""
    echo "Directory structure:"
    ls -lh "$STATIC_ROOT/_nuxt" | head -10
    echo ""
else
    echo -e "${RED}✗${NC} Extraction failed - _nuxt directory not found"
    exit 1
fi

echo "=========================================="
echo -e "${GREEN}Extraction completed!${NC}"
echo "=========================================="
echo ""
echo "Static files location: $STATIC_ROOT"
echo ""
echo "Next steps:"
echo "1. Update nginx.conf.production with STATIC_ROOT=$STATIC_ROOT"
echo "2. Test nginx config: sudo nginx -t"
echo "3. Reload nginx: sudo systemctl reload nginx"
echo ""
echo -e "${YELLOW}Note:${NC} Run this script after every frontend rebuild"
echo ""


