#!/bin/bash
# Quick check script - run on VPS to diagnose the issue

echo "🔍 Quick Diagnostic Check"
echo "========================="
echo ""

# 1. Check if files exist in container
echo "1. Checking if JS files exist in container..."
docker exec gowa-app-prod ls -la /app/public/_nuxt/ | grep -E "DQUWNHGA|CcVpZqDF|BHZQsgze"

echo ""
echo "2. Checking symlink..."
docker exec gowa-app-prod ls -la /app/server/chunks/

echo ""
echo "3. Testing direct access to Nuxt server..."
docker exec gowa-app-prod wget -q -O /dev/null http://localhost:3000/_nuxt/DQUWNHGA.js && echo "✅ File accessible from inside container" || echo "❌ File NOT accessible"

echo ""
echo "4. Testing from host..."
curl -I http://localhost:3002/_nuxt/DQUWNHGA.js 2>&1 | head -5

echo ""
echo "5. Checking current Nginx config..."
sudo grep -A 5 "location /_nuxt/" /etc/nginx/sites-available/app2.anakhebat || echo "⚠️ No _nuxt location block found!"

echo ""
echo "========================="
echo "Diagnostic complete"
