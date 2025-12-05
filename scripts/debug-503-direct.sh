#!/bin/bash

# ============================================================================
# Debug 503 Error - Check langsung di container
# ============================================================================

set -e

CONTAINER_NAME="gowa-app-prod"

# Files yang diminta browser
FILES=(
    "DSzMySIv.js"
    "BcxIJ4tb.js"
    "Czrdq3ge.js"
    "CqpLj_Ac.js"
    "yOviaWq-.js"
    "BK-TQHNH.js"
    "DlAUqK2U.js"
    "C6qEwHNt.js"
)

echo "=========================================="
echo "  Debug 503 - Check Files in Container"
echo "=========================================="
echo ""

echo "Checking if files exist..."
echo ""

for file in "${FILES[@]}"; do
    echo "File: $file"
    
    # Check di server/chunks/public/_nuxt
    if docker exec "$CONTAINER_NAME" test -f "/app/.output/server/chunks/public/_nuxt/$file" 2>/dev/null; then
        echo "  ✓ EXISTS in server/chunks/public/_nuxt/"
        SIZE=$(docker exec "$CONTAINER_NAME" stat -c%s "/app/.output/server/chunks/public/_nuxt/$file" 2>/dev/null || echo "unknown")
        echo "  Size: $SIZE bytes"
    elif docker exec "$CONTAINER_NAME" test -f "/app/.output/public/_nuxt/$file" 2>/dev/null; then
        echo "  ✓ EXISTS in public/_nuxt/ (but not in server/chunks/public/_nuxt/)"
        echo "  ✗ PROBLEM: File not in location where Nitro expects it!"
    else
        echo "  ✗ NOT FOUND anywhere"
    fi
    
    # Test HTTP access from container
    echo -n "  HTTP test from container: "
    HTTP_CODE=$(docker exec "$CONTAINER_NAME" sh -c "wget -q --spider -S http://localhost:3000/_nuxt/$file 2>&1 | grep -i 'HTTP' | head -1 | awk '{print \$2}'" 2>/dev/null || echo "000")
    if [ "$HTTP_CODE" = "200" ]; then
        echo "200 OK"
    elif [ "$HTTP_CODE" = "404" ]; then
        echo "404 NOT FOUND"
    elif [ "$HTTP_CODE" = "500" ]; then
        echo "500 ERROR"
    else
        echo "FAILED (code: $HTTP_CODE)"
    fi
    
    # Test HTTP access from host
    echo -n "  HTTP test from host: "
    HOST_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 "http://127.0.0.1:3002/_nuxt/$file" 2>/dev/null || echo "000")
    echo "$HOST_CODE"
    
    echo ""
done

echo "=========================================="
echo "  Summary"
echo "=========================================="
echo ""
echo "If files don't exist:"
echo "  - HTML page references wrong files"
echo "  - Need to rebuild container"
echo ""
echo "If files exist but 404/500:"
echo "  - Nitro cannot find files"
echo "  - Check file location"
echo ""
echo "If files exist but 503 from host:"
echo "  - Nginx proxy issue"
echo "  - Check nginx config and logs"
echo ""

