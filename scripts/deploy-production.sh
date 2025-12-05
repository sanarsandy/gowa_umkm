#!/bin/bash

# ============================================================================
# Production Deployment Script
# ============================================================================
# Script lengkap untuk deploy aplikasi ke VPS production
# Mengikuti best practices untuk serving static files
#
# Usage: ./scripts/deploy-production.sh [domain]
# Example: ./scripts/deploy-production.sh app.example.com
# ============================================================================

set -e  # Exit on error

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
DOMAIN="${1:-yourdomain.com}"
STATIC_ROOT="${STATIC_ROOT:-/var/www/gowa_umkm/static}"
NGINX_SITES_AVAILABLE="/etc/nginx/sites-available"
NGINX_SITES_ENABLED="/etc/nginx/sites-enabled"
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "=========================================="
echo "  Production Deployment Script"
echo "=========================================="
echo "Domain: $DOMAIN"
echo "Static Root: $STATIC_ROOT"
echo "Project Dir: $PROJECT_DIR"
echo ""

# Check if running as root for certain operations
if [ "$EUID" -ne 0 ]; then 
    SUDO="sudo"
else
    SUDO=""
fi

# ============================================================================
# Step 1: Pre-deployment Checks
# ============================================================================
echo -e "${BLUE}[1/10]${NC} Pre-deployment checks..."

# Check if .env file exists
if [ ! -f "$PROJECT_DIR/.env" ]; then
    echo -e "${RED}✗${NC} .env file not found"
    echo -e "${YELLOW}ℹ${NC} Please create .env file from .env.example"
    exit 1
fi
echo -e "${GREEN}✓${NC} .env file found"

# Check if docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗${NC} Docker is not installed"
    exit 1
fi
echo -e "${GREEN}✓${NC} Docker is installed"

# Check if docker compose is available
if ! docker compose version &> /dev/null; then
    echo -e "${RED}✗${NC} Docker Compose is not available"
    exit 1
fi
echo -e "${GREEN}✓${NC} Docker Compose is available"

# Check if nginx is installed
if ! command -v nginx &> /dev/null; then
    echo -e "${YELLOW}⚠${NC} Nginx is not installed. Installing..."
    $SUDO apt-get update
    $SUDO apt-get install -y nginx
fi
echo -e "${GREEN}✓${NC} Nginx is installed"
echo ""

# ============================================================================
# Step 2: Backup Database
# ============================================================================
echo -e "${BLUE}[2/10]${NC} Creating database backup..."
BACKUP_DIR="$PROJECT_DIR/backups"
mkdir -p "$BACKUP_DIR"
BACKUP_FILE="$BACKUP_DIR/backup_$(date +%Y%m%d_%H%M%S).sql"

if docker compose -f "$PROJECT_DIR/docker-compose.prod.yml" ps | grep -q "gowa-db-prod"; then
    if docker compose -f "$PROJECT_DIR/docker-compose.prod.yml" exec -T db pg_dump -U gowa_user gowa_db > "$BACKUP_FILE" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} Backup saved to: $BACKUP_FILE"
    else
        echo -e "${YELLOW}⚠${NC} Backup failed, continuing..."
    fi
else
    echo -e "${YELLOW}⚠${NC} Database container not running, skipping backup"
fi
echo ""

# ============================================================================
# Step 3: Pull Latest Code
# ============================================================================
echo -e "${BLUE}[3/10]${NC} Pulling latest code from git..."
cd "$PROJECT_DIR"
if [ -d ".git" ]; then
    if git pull origin main 2>/dev/null || git pull origin master 2>/dev/null; then
        echo -e "${GREEN}✓${NC} Code updated"
    else
        echo -e "${YELLOW}⚠${NC} Git pull failed or not a git repository, continuing..."
    fi
else
    echo -e "${YELLOW}⚠${NC} Not a git repository, skipping pull"
fi
echo ""

# ============================================================================
# Step 4: Build Containers
# ============================================================================
echo -e "${BLUE}[4/10]${NC} Building containers..."
cd "$PROJECT_DIR"
docker compose -f docker-compose.prod.yml build --no-cache
echo -e "${GREEN}✓${NC} Containers built"
echo ""

# ============================================================================
# Step 5: Start Services
# ============================================================================
echo -e "${BLUE}[5/10]${NC} Starting services..."
docker compose -f docker-compose.prod.yml up -d
echo -e "${GREEN}✓${NC} Services started"
echo ""

# ============================================================================
# Step 6: Wait for Services to be Ready
# ============================================================================
echo -e "${BLUE}[6/10]${NC} Waiting for services to be ready..."
sleep 10

# Check API health
MAX_RETRIES=30
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -s http://localhost:8087/api/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC} API is healthy"
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo -e "${YELLOW}ℹ${NC} Waiting for API... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo -e "${RED}✗${NC} API health check failed"
    echo -e "${YELLOW}ℹ${NC} Check logs: docker compose -f docker-compose.prod.yml logs api"
fi
echo ""

# ============================================================================
# Step 7: Extract Static Files
# ============================================================================
echo -e "${BLUE}[7/10]${NC} Extracting static files from container..."
if [ -f "$PROJECT_DIR/scripts/extract-static-files.sh" ]; then
    chmod +x "$PROJECT_DIR/scripts/extract-static-files.sh"
    STATIC_ROOT="$STATIC_ROOT" "$PROJECT_DIR/scripts/extract-static-files.sh"
else
    echo -e "${YELLOW}⚠${NC} extract-static-files.sh not found, skipping..."
fi
echo ""

# ============================================================================
# Step 8: Setup Nginx Configuration
# ============================================================================
echo -e "${BLUE}[8/10]${NC} Setting up Nginx configuration..."

# Create nginx config from template
if [ -f "$PROJECT_DIR/nginx.conf.production" ]; then
    NGINX_CONFIG="$NGINX_SITES_AVAILABLE/$DOMAIN"
    
    # Replace placeholders in nginx config
    $SUDO cp "$PROJECT_DIR/nginx.conf.production" "$NGINX_CONFIG"
    $SUDO sed -i "s|{{DOMAIN}}|$DOMAIN|g" "$NGINX_CONFIG"
    $SUDO sed -i "s|{{STATIC_ROOT}}|$STATIC_ROOT|g" "$NGINX_CONFIG"
    
    # Create symlink if it doesn't exist
    if [ ! -L "$NGINX_SITES_ENABLED/$DOMAIN" ]; then
        $SUDO ln -s "$NGINX_CONFIG" "$NGINX_SITES_ENABLED/$DOMAIN"
    fi
    
    # Test nginx configuration
    if $SUDO nginx -t; then
        echo -e "${GREEN}✓${NC} Nginx configuration is valid"
    else
        echo -e "${RED}✗${NC} Nginx configuration test failed"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠${NC} nginx.conf.production not found, skipping..."
fi
echo ""

# ============================================================================
# Step 9: Setup SSL Certificate (if certbot is available)
# ============================================================================
echo -e "${BLUE}[9/10]${NC} Checking SSL certificate..."
if command -v certbot &> /dev/null; then
    if [ ! -f "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" ]; then
        echo -e "${YELLOW}ℹ${NC} SSL certificate not found"
        echo -e "${YELLOW}ℹ${NC} To obtain SSL certificate, run:"
        echo "   sudo certbot --nginx -d $DOMAIN -d www.$DOMAIN"
    else
        echo -e "${GREEN}✓${NC} SSL certificate found"
    fi
else
    echo -e "${YELLOW}⚠${NC} Certbot not installed. Install with:"
    echo "   sudo apt-get install certbot python3-certbot-nginx"
fi
echo ""

# ============================================================================
# Step 10: Reload Nginx
# ============================================================================
echo -e "${BLUE}[10/10]${NC} Reloading Nginx..."
if $SUDO systemctl reload nginx; then
    echo -e "${GREEN}✓${NC} Nginx reloaded"
else
    echo -e "${RED}✗${NC} Failed to reload Nginx"
    exit 1
fi
echo ""

# ============================================================================
# Final Status
# ============================================================================
echo "=========================================="
echo -e "${GREEN}Deployment completed!${NC}"
echo "=========================================="
echo ""
echo "Service Status:"
docker compose -f docker-compose.prod.yml ps
echo ""
echo "Next steps:"
echo "1. Verify application: https://$DOMAIN"
echo "2. Check logs: docker compose -f docker-compose.prod.yml logs -f"
echo "3. Monitor static files: ls -lh $STATIC_ROOT/_nuxt"
echo ""
echo -e "${YELLOW}Important:${NC}"
echo "- If SSL certificate is not set up, configure it with certbot"
echo "- Update DNS records to point to this server"
echo "- Verify firewall allows ports 80 and 443"
echo ""


