#!/bin/bash

# ============================================================================
# Fix Static Files - Final Comprehensive Solution
# ============================================================================
# Script komprehensif untuk fix static files issue
# 1. Check struktur file
# 2. Fix permissions dengan root
# 3. Verify file accessibility
# 4. Test HTTP access
# ============================================================================

set -e

CONTAINER_NAME="gowa-app-prod"

echo "=========================================="
echo "  Fix Static Files - Final Solution"
echo "=========================================="
echo ""

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "❌ Container $CONTAINER_NAME is not running"
    exit 1
fi

echo "Step 1: Checking current structure..."
echo "----------------------------------------"
SOURCE_COUNT=$(docker exec "$CONTAINER_NAME" sh -c "find /app/.output/public/_nuxt -type f 2>/dev/null | wc -l" || echo "0")
DEST_COUNT=$(docker exec "$CONTAINER_NAME" sh -c "find /app/.output/server/chunks/public/_nuxt -type f 2>/dev/null | wc -l" || echo "0")

echo "Source files (.output/public/_nuxt): $SOURCE_COUNT"
echo "Destination files (.output/server/chunks/public/_nuxt): $DEST_COUNT"
echo ""

if [ "$SOURCE_COUNT" -eq "0" ]; then
    echo "❌ No source files found! Container might need rebuild."
    exit 1
fi

echo "Step 2: Creating destination directory (as root)..."
echo "----------------------------------------"
docker exec -u root "$CONTAINER_NAME" sh -c "
    mkdir -p /app/.output/server/chunks/public && \
    chown -R 1001:1001 /app/.output/server/chunks/public && \
    chmod -R 755 /app/.output/server/chunks/public
"
echo "✓ Directory created"
echo ""

echo "Step 3: Copying files (as root)..."
echo "----------------------------------------"
if [ "$DEST_COUNT" -eq "0" ] || [ "$DEST_COUNT" -lt "$SOURCE_COUNT" ]; then
    echo "Copying _nuxt directory..."
    docker exec -u root "$CONTAINER_NAME" sh -c "
        rm -rf /app/.output/server/chunks/public/_nuxt && \
        cp -r /app/.output/public/_nuxt /app/.output/server/chunks/public/ && \
        chown -R 1001:1001 /app/.output/server/chunks/public/_nuxt && \
        chmod -R 755 /app/.output/server/chunks/public/_nuxt && \
        echo '✓ _nuxt copied'
    "
    
    # Copy builds if exists
    if docker exec "$CONTAINER_NAME" test -d "/app/.output/public/builds" 2>/dev/null; then
        echo "Copying builds directory..."
        docker exec -u root "$CONTAINER_NAME" sh -c "
            rm -rf /app/.output/server/chunks/public/builds && \
            cp -r /app/.output/public/builds /app/.output/server/chunks/public/ && \
            chown -R 1001:1001 /app/.output/server/chunks/public/builds && \
            chmod -R 755 /app/.output/server/chunks/public/builds && \
            echo '✓ builds copied'
        "
    fi
else
    echo "Files already exist, fixing permissions only..."
fi

# Fix all permissions
echo "Fixing all permissions..."
docker exec -u root "$CONTAINER_NAME" sh -c "
    chown -R 1001:1001 /app/.output/server/chunks/public && \
    chmod -R 755 /app/.output/server/chunks/public
"
echo "✓ Permissions fixed"
echo ""

echo "Step 4: Verification..."
echo "----------------------------------------"
NEW_DEST_COUNT=$(docker exec "$CONTAINER_NAME" sh -c "find /app/.output/server/chunks/public/_nuxt -type f 2>/dev/null | wc -l" || echo "0")
echo "Files in destination: $NEW_DEST_COUNT"

if [ "$NEW_DEST_COUNT" -eq "0" ]; then
    echo "❌ No files in destination!"
    exit 1
fi

# Test specific file
echo ""
echo "Testing specific file (F_WmEIm6.js)..."
if docker exec "$CONTAINER_NAME" test -f "/app/.output/server/chunks/public/_nuxt/F_WmEIm6.js" 2>/dev/null; then
    echo "✓ F_WmEIm6.js exists"
    docker exec "$CONTAINER_NAME" ls -lh /app/.output/server/chunks/public/_nuxt/F_WmEIm6.js
else
    echo "⚠ F_WmEIm6.js not found, listing available files:"
    docker exec "$CONTAINER_NAME" ls /app/.output/server/chunks/public/_nuxt | head -5
fi

# Test file readability
echo ""
echo "Testing file readability..."
if docker exec "$CONTAINER_NAME" sh -c "test -r /app/.output/server/chunks/public/_nuxt/F_WmEIm6.js" 2>/dev/null; then
    echo "✓ File is readable"
else
    echo "❌ File is not readable!"
    exit 1
fi

echo ""
echo "Step 5: Restarting container to ensure changes take effect..."
echo "----------------------------------------"
docker restart "$CONTAINER_NAME"
echo "Waiting for container to be ready..."
sleep 10

# Check if container is healthy
if docker ps | grep -q "$CONTAINER_NAME"; then
    echo "✓ Container restarted successfully"
else
    echo "❌ Container failed to start!"
    exit 1
fi

echo ""
echo "=========================================="
echo "✓ Fix completed!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Test in browser: https://app3.anakhebat.web.id"
echo "2. Check container logs: docker logs gowa-app-prod"
echo "3. Check nginx logs: sudo tail -f /var/log/nginx/app3_error.log"
echo ""

