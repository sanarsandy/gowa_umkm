#!/bin/bash

# ============================================================================
# Quick Setup Script untuk Static Files
# ============================================================================
# Script ini membantu setup static files directory dan permissions
# Usage: ./scripts/setup-static-files.sh [static_root]
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
STATIC_ROOT="${1:-/var/www/gowa_umkm/static}"

echo "=========================================="
echo "  Setup Static Files Directory"
echo "=========================================="
echo "Static Root: $STATIC_ROOT"
echo ""

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then 
    SUDO="sudo"
else
    SUDO=""
fi

# Create directory
echo -e "${BLUE}[1/3]${NC} Creating directory..."
$SUDO mkdir -p "$STATIC_ROOT"
echo -e "${GREEN}✓${NC} Directory created"
echo ""

# Set permissions
echo -e "${BLUE}[2/3]${NC} Setting permissions..."
$SUDO chown -R www-data:www-data "$STATIC_ROOT"
$SUDO chmod -R 755 "$STATIC_ROOT"
echo -e "${GREEN}✓${NC} Permissions set"
echo ""

# Verify
echo -e "${BLUE}[3/3]${NC} Verification..."
if [ -d "$STATIC_ROOT" ]; then
    OWNER=$(stat -c '%U:%G' "$STATIC_ROOT" 2>/dev/null || stat -f '%Su:%Sg' "$STATIC_ROOT" 2>/dev/null)
    PERMS=$(stat -c '%a' "$STATIC_ROOT" 2>/dev/null || stat -f '%A' "$STATIC_ROOT" 2>/dev/null)
    echo -e "${GREEN}✓${NC} Directory exists"
    echo -e "${GREEN}✓${NC} Owner: $OWNER"
    echo -e "${GREEN}✓${NC} Permissions: $PERMS"
else
    echo -e "${RED}✗${NC} Directory creation failed"
    exit 1
fi

echo ""
echo "=========================================="
echo -e "${GREEN}Setup completed!${NC}"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Build frontend: docker compose -f docker-compose.prod.yml build app"
echo "2. Start app: docker compose -f docker-compose.prod.yml up -d app"
echo "3. Extract static files: ./scripts/extract-static-files.sh"
echo ""


