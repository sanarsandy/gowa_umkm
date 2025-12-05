#!/bin/bash

# ============================================================================
# Fix Container Static Files - Immediate Fix
# ============================================================================
# Script untuk fix static files di container yang sedang running
# Copy files dari .output/public ke .output/server/chunks/public
# ============================================================================

set -e

CONTAINER_NAME="gowa-app-prod"

echo "=========================================="
echo "  Fix Container Static Files - Immediate"
echo "=========================================="
echo ""

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "❌ Container $CONTAINER_NAME is not running"
    exit 1
fi

echo "1. Checking current structure..."
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "
    echo 'Source: .output/public/_nuxt/'
    if [ -d '/app/.output/public/_nuxt' ]; then
        find /app/.output/public/_nuxt -type f | wc -l | xargs echo 'Files:'
    else
        echo '✗ Directory does not exist'
        exit 1
    fi
    echo ''
    echo 'Destination: .output/server/chunks/public/_nuxt/'
    if [ -d '/app/.output/server/chunks/public/_nuxt' ]; then
        find /app/.output/server/chunks/public/_nuxt -type f | wc -l | xargs echo 'Files:'
    else
        echo '✗ Directory does not exist'
    fi
"
echo ""

echo "2. Creating destination directory..."
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "mkdir -p /app/.output/server/chunks/public"
echo "✓ Directory created"
echo ""

echo "3. Copying _nuxt directory..."
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "
    if [ -d '/app/.output/public/_nuxt' ]; then
        cp -a /app/.output/public/_nuxt /app/.output/server/chunks/public/ && \
        echo '✓ _nuxt directory copied'
    else
        echo '✗ Source directory does not exist'
        exit 1
    fi
"
echo ""

echo "4. Copying builds directory (if exists)..."
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "
    if [ -d '/app/.output/public/builds' ]; then
        cp -a /app/.output/public/builds /app/.output/server/chunks/public/ && \
        echo '✓ builds directory copied'
    else
        echo '⚠ builds directory does not exist (optional)'
    fi
"
echo ""

echo "5. Copying other files from public root..."
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "
    find /app/.output/public -maxdepth 1 -type f -exec cp {} /app/.output/server/chunks/public/ \; 2>/dev/null || true && \
    find /app/.output/public -mindepth 1 -maxdepth 1 -type d ! -name '_nuxt' ! -name 'builds' -exec cp -a {} /app/.output/server/chunks/public/ \; 2>/dev/null || true && \
    echo '✓ Other files copied'
"
echo ""

echo "6. Setting permissions..."
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "
    chown -R 1001:1001 /app/.output/server/chunks/public && \
    echo '✓ Permissions set'
"
echo ""

echo "7. Verification..."
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "
    if [ -d '/app/.output/server/chunks/public/_nuxt' ]; then
        FILE_COUNT=\$(find /app/.output/server/chunks/public/_nuxt -type f | wc -l)
        echo \"✓ Files in destination: \$FILE_COUNT\"
        echo ''
        echo 'Sample files:'
        ls -la /app/.output/server/chunks/public/_nuxt | head -5
    else
        echo '✗ Directory does not exist'
        exit 1
    fi
"
echo ""

echo "8. Testing specific file (F_WmEIm6.js)..."
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "
    if [ -f '/app/.output/server/chunks/public/_nuxt/F_WmEIm6.js' ]; then
        echo '✓ F_WmEIm6.js found'
        ls -lh /app/.output/server/chunks/public/_nuxt/F_WmEIm6.js
    else
        echo '✗ F_WmEIm6.js not found'
        echo 'Available files:'
        ls /app/.output/server/chunks/public/_nuxt | head -10
    fi
"
echo ""

echo "=========================================="
echo "✓ Fix completed!"
echo "=========================================="
echo ""
echo "The container should now be able to serve static files."
echo "Test: curl -I https://app3.anakhebat.web.id/_nuxt/F_WmEIm6.js"
echo ""

