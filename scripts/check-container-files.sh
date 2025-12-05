#!/bin/bash

# ============================================================================
# Check Container Files Structure
# ============================================================================
# Script untuk check struktur file di container dan verify apakah
# static files ada di lokasi yang benar
# ============================================================================

set -e

CONTAINER_NAME="gowa-app-prod"

echo "=========================================="
echo "  Check Container Files Structure"
echo "=========================================="
echo ""

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "❌ Container $CONTAINER_NAME is not running"
    exit 1
fi

echo "1. Checking .output/public/_nuxt/ (source location):"
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "if [ -d '/app/.output/public/_nuxt' ]; then echo '✓ Directory exists'; find /app/.output/public/_nuxt -type f | wc -l | xargs echo 'Files count:'; ls -la /app/.output/public/_nuxt | head -10; else echo '✗ Directory does not exist'; fi"
echo ""

echo "2. Checking .output/server/chunks/public/_nuxt/ (Nitro lookup location):"
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "if [ -d '/app/.output/server/chunks/public/_nuxt' ]; then echo '✓ Directory exists'; find /app/.output/server/chunks/public/_nuxt -type f | wc -l | xargs echo 'Files count:'; ls -la /app/.output/server/chunks/public/_nuxt | head -10; else echo '✗ Directory does not exist'; fi"
echo ""

echo "3. Checking .output/server/chunks/public/ (parent directory):"
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "if [ -d '/app/.output/server/chunks/public' ]; then echo '✓ Directory exists'; ls -la /app/.output/server/chunks/public/; else echo '✗ Directory does not exist'; fi"
echo ""

echo "4. Checking .output/server/chunks/ (chunks directory):"
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "if [ -d '/app/.output/server/chunks' ]; then echo '✓ Directory exists'; ls -la /app/.output/server/chunks/; else echo '✗ Directory does not exist'; fi"
echo ""

echo "5. Checking .output/public/builds/ (builds directory):"
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "if [ -d '/app/.output/public/builds' ]; then echo '✓ Directory exists'; find /app/.output/public/builds -type f | wc -l | xargs echo 'Files count:'; else echo '✗ Directory does not exist'; fi"
echo ""

echo "6. Checking .output/server/chunks/public/builds/ (Nitro lookup for builds):"
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "if [ -d '/app/.output/server/chunks/public/builds' ]; then echo '✓ Directory exists'; find /app/.output/server/chunks/public/builds -type f | wc -l | xargs echo 'Files count:'; else echo '✗ Directory does not exist'; fi"
echo ""

echo "7. Sample file check (F_WmEIm6.js):"
echo "----------------------------------------"
docker exec "$CONTAINER_NAME" sh -c "
    if [ -f '/app/.output/public/_nuxt/F_WmEIm6.js' ]; then
        echo '✓ Found in .output/public/_nuxt/'
    else
        echo '✗ Not found in .output/public/_nuxt/'
    fi
    if [ -f '/app/.output/server/chunks/public/_nuxt/F_WmEIm6.js' ]; then
        echo '✓ Found in .output/server/chunks/public/_nuxt/'
    else
        echo '✗ Not found in .output/server/chunks/public/_nuxt/'
    fi
"
echo ""

echo "=========================================="
echo "  Summary"
echo "=========================================="
echo ""

