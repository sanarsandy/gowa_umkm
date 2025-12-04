# Nuxt 503 Static Asset Fix - Applied 2025-12-04

## Problem
The production environment was experiencing 503 Service Unavailable errors for Nuxt static assets (`_nuxt/*.js` files) and WebSocket connection failures.

## Root Cause
The `Dockerfile.prod` was creating an unnecessary symlink from `/app/server/chunks/public` to `/app/public`, which doesn't align with Nuxt 3's build output structure. Nuxt 3 with the `node-server` preset automatically serves static files from `.output/public/` without needing symlinks.

## Changes Applied

### 1. Frontend Dockerfile (`frontend/Dockerfile.prod`)
- **Removed**: Unnecessary symlink creation logic
- **Added**: Build output verification steps
- **Added**: Static asset verification in final image
- **Added**: Proper file permissions for public directory

### 2. Nuxt Configuration (`frontend/nuxt.config.ts`)
- **Added**: Explicit `publicAssets` configuration for Nitro
- **Added**: WebSocket support via `experimental.websocket`
- **Enhanced**: Static asset caching with 1-year max-age

### 3. Nginx Configuration (`nginx.conf.example`)
- **Added**: WebSocket upgrade map at configuration top
- **Added**: Dedicated `/_nuxt/` location block with aggressive caching
- **Added**: Error handling and fallback for Nuxt errors
- **Enhanced**: WebSocket proxy configuration for `/api/` endpoints
- **Added**: Extended timeouts for WebSocket connections (24 hours)

## Deployment Instructions

### On Your VPS Server

1. **Pull the latest changes:**
   ```bash
   cd /path/to/gowa_umkm
   git pull origin main
   ```

2. **Run the deployment script:**
   ```bash
   ./scripts/deploy-nuxt-fix.sh
   ```

3. **Update Nginx configuration:**
   ```bash
   sudo cp nginx.conf.example /etc/nginx/sites-available/app2.anakhebat.web.id
   sudo nginx -t
   sudo systemctl reload nginx
   ```

4. **Verify the fix:**
   - Open https://app2.anakhebat.web.id/ in your browser
   - Check browser console - no 503 errors for `_nuxt/*.js` files
   - Verify WebSocket connection succeeds

## Verification Checklist

- [ ] Frontend container rebuilt successfully
- [ ] Static assets exist in `/app/public/_nuxt/` directory
- [ ] Nuxt server responds to requests
- [ ] Nginx configuration updated and reloaded
- [ ] Browser shows no 503 errors for static assets
- [ ] WebSocket connection to `/api/ws` succeeds
- [ ] Dashboard loads and functions correctly

## Rollback Procedure

If issues occur:

```bash
# Revert to previous commit
git revert HEAD

# Rebuild and restart
docker compose -f docker-compose.prod.yml build app
docker compose -f docker-compose.prod.yml up -d app

# Restore previous Nginx config
sudo cp /etc/nginx/sites-available/app2.anakhebat.web.id.backup \
     /etc/nginx/sites-available/app2.anakhebat.web.id
sudo systemctl reload nginx
```

## Technical Details

### Why the Symlink Approach Failed
Nuxt 3's Nitro server with `node-server` preset expects static files to be in `.output/public/` and serves them automatically. Creating symlinks to `/app/server/chunks/public` was unnecessary and potentially caused path resolution issues.

### How the Fix Works
1. The Dockerfile now simply copies `.output/` directory from the builder stage
2. Nuxt's Nitro server automatically serves files from `/app/public/_nuxt/`
3. Nginx proxies `/_nuxt/` requests to the Nuxt server with proper caching
4. WebSocket connections are properly upgraded via Nginx map directive

### Performance Improvements
- Static assets now cached for 1 year (immutable)
- Reduced server load with proper cache headers
- WebSocket connections more stable with extended timeouts

## Support

If you encounter any issues:
1. Check container logs: `docker logs gowa-app-prod`
2. Verify static assets: `docker exec gowa-app-prod ls -la /app/public/_nuxt/`
3. Test Nuxt server: `docker exec gowa-app-prod wget -O- http://localhost:3000/`
4. Check Nginx logs: `sudo tail -f /var/log/nginx/error.log`
