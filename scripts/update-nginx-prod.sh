#!/bin/bash
# Script to update Nginx configuration on production server

set -e

echo "🔧 Updating Nginx Configuration for app2.anakhebat.web.id"
echo "=========================================================="
echo ""

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then 
    echo "⚠️  This script needs to be run with sudo"
    echo "Usage: sudo ./scripts/update-nginx-prod.sh"
    exit 1
fi

NGINX_SITE_CONFIG="/etc/nginx/sites-available/app2.anakhebat"
NGINX_SITE_ENABLED="/etc/nginx/sites-enabled/app2.anakhebat"

# Backup current configuration
echo "📦 Backing up current configuration..."
if [ -f "$NGINX_SITE_CONFIG" ]; then
    cp "$NGINX_SITE_CONFIG" "$NGINX_SITE_CONFIG.backup.$(date +%Y%m%d_%H%M%S)"
    echo "✅ Backup created: $NGINX_SITE_CONFIG.backup.$(date +%Y%m%d_%H%M%S)"
else
    echo "⚠️  No existing configuration found at $NGINX_SITE_CONFIG"
fi

# Copy new configuration
echo ""
echo "📝 Installing new configuration..."
cp nginx.conf.prod "$NGINX_SITE_CONFIG"
echo "✅ Configuration copied to $NGINX_SITE_CONFIG"

# Create symlink if it doesn't exist
if [ ! -L "$NGINX_SITE_ENABLED" ]; then
    echo ""
    echo "🔗 Creating symlink..."
    ln -s "$NGINX_SITE_CONFIG" "$NGINX_SITE_ENABLED"
    echo "✅ Symlink created"
fi

# Test Nginx configuration
echo ""
echo "🧪 Testing Nginx configuration..."
if nginx -t; then
    echo "✅ Nginx configuration is valid"
else
    echo "❌ Nginx configuration test failed!"
    echo "Restoring backup..."
    if [ -f "$NGINX_SITE_CONFIG.backup.$(date +%Y%m%d_%H%M%S)" ]; then
        cp "$NGINX_SITE_CONFIG.backup.$(date +%Y%m%d_%H%M%S)" "$NGINX_SITE_CONFIG"
        echo "✅ Backup restored"
    fi
    exit 1
fi

# Reload Nginx
echo ""
echo "🔄 Reloading Nginx..."
if systemctl reload nginx; then
    echo "✅ Nginx reloaded successfully"
else
    echo "❌ Failed to reload Nginx"
    exit 1
fi

echo ""
echo "=========================================================="
echo "✅ Nginx configuration updated successfully!"
echo ""
echo "Key changes applied:"
echo "  ✅ Added WebSocket upgrade map"
echo "  ✅ Fixed /api/ location with WebSocket support"
echo "  ✅ Extended WebSocket timeouts to 24 hours"
echo "  ✅ Added error handling for _nuxt assets"
echo ""
echo "Next steps:"
echo "1. Test WebSocket connection: https://app2.anakhebat.web.id/"
echo "2. Check browser console for errors"
echo "3. Monitor logs: sudo tail -f /var/log/nginx/app2.anakhebat-error.log"
