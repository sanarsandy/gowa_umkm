#!/bin/bash

# ============================================================================
# Find Missing File - Cari file yang hilang
# ============================================================================

set -e

CONTAINER_NAME="gowa-app-prod"
MISSING_FILE="4098426c-7e8d-4aee-a1df-08d1f457e220.js"

echo "=========================================="
echo "  Find Missing File: $MISSING_FILE"
echo "=========================================="
echo ""

# Search in all possible locations
echo "Searching for file in container..."
echo ""

# Check public/_nuxt
if docker exec "$CONTAINER_NAME" test -f "/app/.output/public/_nuxt/$MISSING_FILE" 2>/dev/null; then
    echo "✓ Found in: /app/.output/public/_nuxt/$MISSING_FILE"
    docker exec "$CONTAINER_NAME" ls -lh "/app/.output/public/_nuxt/$MISSING_FILE"
fi

# Check server/chunks/public/_nuxt
if docker exec "$CONTAINER_NAME" test -f "/app/.output/server/chunks/public/_nuxt/$MISSING_FILE" 2>/dev/null; then
    echo "✓ Found in: /app/.output/server/chunks/public/_nuxt/$MISSING_FILE"
    docker exec "$CONTAINER_NAME" ls -lh "/app/.output/server/chunks/public/_nuxt/$MISSING_FILE"
fi

# Check builds directory
if docker exec "$CONTAINER_NAME" find /app/.output -name "$MISSING_FILE" 2>/dev/null | head -1; then
    FOUND_PATH=$(docker exec "$CONTAINER_NAME" find /app/.output -name "$MISSING_FILE" 2>/dev/null | head -1)
    echo "✓ Found in: $FOUND_PATH"
    docker exec "$CONTAINER_NAME" ls -lh "$FOUND_PATH"
fi

# Search everywhere
echo ""
echo "Searching everywhere in .output..."
ALL_MATCHES=$(docker exec "$CONTAINER_NAME" find /app/.output -name "*4098426c*" 2>/dev/null || echo "")

if [ -n "$ALL_MATCHES" ]; then
    echo "Found similar files:"
    echo "$ALL_MATCHES" | while read -r match; do
        echo "  - $match"
    done
else
    echo "✗ File not found anywhere in .output"
    echo ""
    echo "This file might be:"
    echo "  1. Generated dynamically by Nuxt"
    echo "  2. Missing from build"
    echo "  3. In a different location"
fi

echo ""

