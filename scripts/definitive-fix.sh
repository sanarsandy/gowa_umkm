#!/bin/bash
# DEFINITIVE FIX - Rebuild and restart with proven configuration

set -e

echo "🔧 DEFINITIVE FIX - Rebuilding Frontend"
echo "======================================="
echo ""

cd /var/rumah_afiat/gowa_umkm

# Pull latest changes
echo "1. Pulling latest changes..."
git pull origin main

# Stop and remove old container
echo ""
echo "2. Stopping and removing old container..."
docker compose -f docker-compose.prod.yml stop app
docker compose -f docker-compose.prod.yml rm -f app

# Remove old image to force complete rebuild
echo ""
echo "3. Removing old image..."
docker rmi gowa_umkm-app || true

# Rebuild with no cache
echo ""
echo "4. Rebuilding with --no-cache..."
docker compose -f docker-compose.prod.yml build --no-cache app

# Start container
echo ""
echo "5. Starting container..."
docker compose -f docker-compose.prod.yml up -d app

# Wait for startup
echo ""
echo "6. Waiting for container to start..."
sleep 10

# Check status
echo ""
echo "7. Checking container status..."
docker ps | grep gowa-app-prod

# Check logs
echo ""
echo "8. Checking logs (last 30 lines)..."
docker logs gowa-app-prod --tail 30

echo ""
echo "======================================="
echo "✅ Rebuild complete!"
echo ""
echo "Next: Test at https://app2.anakhebat.web.id/"
