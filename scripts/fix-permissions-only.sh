#!/bin/bash

# ============================================================================
# Fix Permissions Only - Quick Fix
# ============================================================================
# Script untuk fix permissions saja (jika file sudah ter-copy)
# ============================================================================

set -e

CONTAINER_NAME="gowa-app-prod"

echo "=========================================="
echo "  Fix Permissions Only"
echo "=========================================="
echo ""

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "❌ Container $CONTAINER_NAME is not running"
    exit 1
fi

echo "Fixing permissions as root..."
docker exec -u root "$CONTAINER_NAME" sh -c "
    if [ -d '/app/.output/server/chunks/public' ]; then
        chown -R 1001:1001 /app/.output/server/chunks/public && \
        chmod -R 755 /app/.output/server/chunks/public && \
        echo '✓ Permissions fixed'
    else
        echo '✗ Directory does not exist'
        exit 1
    fi
"

echo ""
echo "Verifying..."
docker exec "$CONTAINER_NAME" sh -c "
    if [ -f '/app/.output/server/chunks/public/_nuxt/F_WmEIm6.js' ]; then
        echo '✓ File exists and is readable'
        ls -lh /app/.output/server/chunks/public/_nuxt/F_WmEIm6.js
    else
        echo '✗ File does not exist'
    fi
"

echo ""
echo "=========================================="
echo "✓ Done!"
echo "=========================================="

