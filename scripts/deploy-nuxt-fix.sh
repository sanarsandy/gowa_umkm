#!/bin/bash
# Deploy script for fixing Nuxt 503 static asset errors
# This script rebuilds and restarts the frontend container

set -e

echo "🔧 Deploying Nuxt Static Asset Fix"
echo "===================================="
echo ""

# Check if we're in the correct directory
if [ ! -f "docker-compose.prod.yml" ]; then
    echo "❌ Error: docker-compose.prod.yml not found"
    echo "Please run this script from the project root directory"
    exit 1
fi

# Stop the frontend container
echo "📦 Stopping frontend container..."
docker compose -f docker-compose.prod.yml stop app

# Rebuild the frontend container
echo "🔨 Rebuilding frontend container..."
docker compose -f docker-compose.prod.yml build --no-cache app

# Start the frontend container
echo "🚀 Starting frontend container..."
docker compose -f docker-compose.prod.yml up -d app

# Wait for container to be ready
echo "⏳ Waiting for container to start..."
sleep 5

# Verify the container is running
if docker ps | grep -q "gowa-app-prod"; then
    echo "✅ Container is running"
else
    echo "❌ Container failed to start"
    echo "Checking logs..."
    docker logs gowa-app-prod --tail 50
    exit 1
fi

# Verify static assets exist in the container
echo ""
echo "🔍 Verifying static assets in container..."
docker exec gowa-app-prod ls -la /app/public/_nuxt/ | head -10

# Count JS files
JS_COUNT=$(docker exec gowa-app-prod find /app/public/_nuxt/ -name "*.js" | wc -l)
echo ""
echo "📊 Found $JS_COUNT JavaScript files in _nuxt directory"

if [ "$JS_COUNT" -gt 0 ]; then
    echo "✅ Static assets verified successfully"
else
    echo "⚠️  Warning: No JavaScript files found in _nuxt directory"
    echo "This may indicate a build issue"
fi

# Test if Nuxt server can serve a static file
echo ""
echo "🧪 Testing static asset serving..."
if docker exec gowa-app-prod wget -q -O /dev/null http://localhost:3000/ 2>/dev/null; then
    echo "✅ Nuxt server is responding"
else
    echo "⚠️  Warning: Nuxt server may not be responding correctly"
fi

echo ""
echo "===================================="
echo "✅ Deployment complete!"
echo ""
echo "Next steps:"
echo "1. Update your Nginx configuration on the server:"
echo "   sudo cp nginx.conf.example /etc/nginx/sites-available/your-site"
echo "   sudo nginx -t"
echo "   sudo systemctl reload nginx"
echo ""
echo "2. Test the application in your browser:"
echo "   https://app2.anakhebat.web.id/"
echo ""
echo "3. Check browser console for any remaining errors"
echo ""
echo "4. View container logs if needed:"
echo "   docker logs -f gowa-app-prod"
