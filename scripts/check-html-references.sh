#!/bin/bash

# ============================================================================
# Check HTML References vs Files in Container
# ============================================================================

set -e

CONTAINER_NAME="gowa-app-prod"

echo "=========================================="
echo "  Check HTML References"
echo "=========================================="
echo ""

# Get HTML from container
echo -e "${BLUE}[1/3]${NC} Getting HTML from container..."
HTML=$(curl -s http://127.0.0.1:3002/ 2>/dev/null || echo "")

if [ -z "$HTML" ]; then
    echo -e "${RED}✗${NC} Cannot get HTML from container"
    exit 1
fi

echo -e "${GREEN}✓${NC} HTML retrieved"
echo ""

# Extract file references from HTML
echo -e "${BLUE}[2/3]${NC} Extracting file references from HTML..."
echo "$HTML" | grep -oE '/_nuxt/[^"]+\.js' | sort -u > /tmp/html_files.txt
echo "$HTML" | grep -oE '/_nuxt/[^"]+\.css' | sort -u >> /tmp/html_files.txt

HTML_FILES=$(cat /tmp/html_files.txt | wc -l)
echo "Found $HTML_FILES file references in HTML"
echo ""

# Get actual files in container
echo -e "${BLUE}[3/3]${NC} Checking if referenced files exist..."
echo ""

CONTAINER_FILES=$(docker exec "$CONTAINER_NAME" find /app/.output/server/chunks/public/_nuxt -type f -name "*.js" -o -name "*.css" 2>/dev/null | xargs -n1 basename | sort -u)

MISSING_COUNT=0
FOUND_COUNT=0

while IFS= read -r ref; do
    if [ -z "$ref" ]; then
        continue
    fi
    filename=$(basename "$ref")
    
    # Check if file exists
    if echo "$CONTAINER_FILES" | grep -q "^$filename$"; then
        echo "✓ $filename - EXISTS"
        FOUND_COUNT=$((FOUND_COUNT + 1))
    else
        echo "✗ $filename - NOT FOUND"
        MISSING_COUNT=$((MISSING_COUNT + 1))
    fi
done < /tmp/html_files.txt

echo ""
echo "=========================================="
echo "  Summary"
echo "=========================================="
echo "Files referenced in HTML: $HTML_FILES"
echo "Files found in container: $FOUND_COUNT"
echo "Files missing: $MISSING_COUNT"
echo ""

if [ "$MISSING_COUNT" -gt 0 ]; then
    echo -e "${RED}PROBLEM:${NC} HTML references files that don't exist!"
    echo ""
    echo "This means:"
    echo "  1. HTML was generated with wrong file references"
    echo "  2. Or files were not copied correctly"
    echo ""
    echo "Solution:"
    echo "  1. Rebuild container: docker compose -f docker-compose.prod.yml build --no-cache app"
    echo "  2. Or check if build process is correct"
else
    echo -e "${GREEN}✓${NC} All referenced files exist in container"
    echo ""
    echo "If still getting 503 errors, the problem might be:"
    echo "  1. Browser cache (clear cache or use incognito)"
    echo "  2. Nginx proxy issue"
    echo "  3. Container cannot serve files (check logs)"
fi

# Cleanup
rm -f /tmp/html_files.txt

echo ""

