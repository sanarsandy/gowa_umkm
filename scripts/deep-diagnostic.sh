#!/bin/bash
# DEEP DIAGNOSTIC - Find the real root cause

echo "🔍 DEEP DIAGNOSTIC - Finding Root Cause"
echo "========================================"
echo ""

# 1. Check container status
echo "1. Container Status:"
docker ps -a | grep gowa-app-prod
echo ""

# 2. Check if Nuxt is actually running
echo "2. Testing Nuxt Server Response:"
curl -I http://localhost:3002/ 2>&1 | head -10
echo ""

# 3. Check specific JS file from INSIDE container
echo "3. Testing file access INSIDE container:"
docker exec gowa-app-prod ls -la /app/public/_nuxt/BHZQsgze.js 2>&1
docker exec gowa-app-prod ls -la /app/server/chunks/public/_nuxt/BHZQsgze.js 2>&1
echo ""

# 4. Test accessing file through Nuxt server INSIDE container
echo "4. Testing Nuxt serving file (inside container):"
docker exec gowa-app-prod wget -q -O /dev/null http://localhost:3000/_nuxt/BHZQsgze.js && echo "✅ SUCCESS" || echo "❌ FAILED"
echo ""

# 5. Test from host to container
echo "5. Testing from host to container port 3002:"
curl -I http://localhost:3002/_nuxt/BHZQsgze.js 2>&1 | head -10
echo ""

# 6. Check Nginx error logs
echo "6. Recent Nginx errors:"
sudo tail -20 /var/log/nginx/app2.anakhebat-error.log
echo ""

# 7. Check container logs for errors
echo "7. Container logs (last 30 lines):"
docker logs gowa-app-prod --tail 30
echo ""

# 8. Check Nginx config for _nuxt location
echo "8. Nginx _nuxt location config:"
sudo grep -A 20 "location /_nuxt/" /etc/nginx/sites-available/app2.anakhebat
echo ""

echo "========================================"
echo "Diagnostic complete"
