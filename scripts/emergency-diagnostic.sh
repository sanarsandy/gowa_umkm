#!/bin/bash
# Emergency diagnostic script for 502/503 errors
# Run this on your VPS to diagnose the issue

set -e

echo "🚨 Emergency Diagnostic for 502/503 Errors"
echo "=========================================="
echo ""

# Check if containers are running
echo "1. Checking container status..."
echo "---"
docker ps -a | grep gowa

echo ""
echo "2. Checking frontend container logs (last 50 lines)..."
echo "---"
docker logs gowa-app-prod --tail 50

echo ""
echo "3. Checking if frontend port 3002 is accessible..."
echo "---"
curl -I http://localhost:3002/ 2>&1 | head -10

echo ""
echo "4. Checking if files exist in container..."
echo "---"
if docker ps | grep -q gowa-app-prod; then
    echo "Container is running, checking files..."
    docker exec gowa-app-prod ls -la /app/public/_nuxt/ 2>&1 | head -20
else
    echo "❌ Container is NOT running!"
fi

echo ""
echo "5. Checking Nginx error logs..."
echo "---"
sudo tail -30 /var/log/nginx/app2.anakhebat-error.log

echo ""
echo "=========================================="
echo "Diagnostic complete. Please share the output."
