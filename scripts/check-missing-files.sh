#!/bin/bash

# ============================================================================
# Check Missing Files - File yang diminta browser vs yang ada di container
# ============================================================================

set -e

CONTAINER_NAME="gowa-app-prod"

echo "=========================================="
echo "  Check Missing Files"
echo "=========================================="
echo ""

# Files yang diminta browser (dari error)
MISSING_FILES=(
    "BcxIJ4tb.js"
    "Czrdq3ge.js"
    "yOviaWq-.js"
    "DlAUqK2U.js"
    "BK-TQHNH.js"
    "C6qEwHNt.js"
    "DeXWKxqL.js"
    "F_WmEIm6.js"
    "DItC9poW.js"
)

echo "Files yang diminta browser:"
for file in "${MISSING_FILES[@]}"; do
    echo "  - $file"
done
echo ""

echo "Checking if files exist in container..."
echo ""

FOUND=0
NOT_FOUND=0

for file in "${MISSING_FILES[@]}"; do
    if docker exec "$CONTAINER_NAME" test -f "/app/.output/server/chunks/public/_nuxt/$file" 2>/dev/null || \
       docker exec "$CONTAINER_NAME" test -f "/app/.output/public/_nuxt/$file" 2>/dev/null; then
        echo "✓ $file - EXISTS"
        FOUND=$((FOUND + 1))
    else
        echo "✗ $file - NOT FOUND"
        NOT_FOUND=$((NOT_FOUND + 1))
    fi
done

echo ""
echo "Summary:"
echo "  Found: $FOUND"
echo "  Not Found: $NOT_FOUND"
echo ""

if [ "$NOT_FOUND" -gt 0 ]; then
    echo "PROBLEM: HTML page references files that don't exist!"
    echo "Solution: Rebuild container to regenerate HTML with correct file references"
    echo ""
    echo "Run:"
    echo "  docker compose -f docker-compose.prod.yml build --no-cache app"
    echo "  docker compose -f docker-compose.prod.yml up -d app"
fi

