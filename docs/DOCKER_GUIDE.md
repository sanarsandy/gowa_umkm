# Gowa UMKM - Docker Setup Guide

## 🚀 Quick Start

### Development Mode (Recommended for Local Development)

Development mode uses hot-reload for both frontend and backend, making development much faster and lighter.

```bash
# Start all services in development mode
docker compose -f docker-compose.dev.yml up

# Start in detached mode
docker compose -f docker-compose.dev.yml up -d

# View logs
docker compose -f docker-compose.dev.yml logs -f

# Stop services
docker compose -f docker-compose.dev.yml down
```

**Features:**
- ✅ Hot-reload for frontend (Nuxt dev server)
- ✅ Hot-reload for backend (Air)
- ✅ Source code mounted as volumes
- ✅ No build step required
- ✅ Fast startup (~30 seconds)
- ✅ Changes reflected immediately

### Production Mode

Production mode builds optimized Docker images.

```bash
# Start all services in production mode
docker compose up

# Or explicitly
docker compose -f docker-compose.yml up

# Build and start
docker compose up --build

# Stop services
docker compose down
```

**Features:**
- ✅ Optimized production builds
- ✅ Multi-stage Docker builds
- ✅ Smaller image sizes
- ✅ Production-ready configuration
- ⚠️ Slower build time (~5-10 minutes)
- ⚠️ Requires rebuild for code changes

## 📦 Docker Files Overview

### Development Files
- `backend/Dockerfile.dev` - Backend development with Air hot-reload
- `frontend/Dockerfile.dev` - Frontend development with Nuxt dev server
- `docker-compose.dev.yml` - Development orchestration
- `backend/.air.toml` - Air configuration for Go hot-reload

### Production Files
- `backend/Dockerfile.prod` - Backend production build
- `frontend/Dockerfile.prod` - Frontend production build (multi-stage)
- `docker-compose.yml` - Production orchestration
- `docker-compose.prod.yml` - Production deployment configuration

## 🔧 Development Workflow

### 1. First Time Setup

```bash
# Copy environment file
cp .env.example .env

# Edit .env with your configuration
nano .env

# Start development environment
docker compose -f docker-compose.dev.yml up -d

# Check logs
docker compose -f docker-compose.dev.yml logs -f
```

### 2. Daily Development

```bash
# Start services
docker compose -f docker-compose.dev.yml up -d

# Make code changes - they will auto-reload!
# Frontend: Changes in /frontend will trigger Nuxt HMR
# Backend: Changes in /backend will trigger Air rebuild

# View specific service logs
docker compose -f docker-compose.dev.yml logs -f api
docker compose -f docker-compose.dev.yml logs -f app

# Restart a specific service
docker compose -f docker-compose.dev.yml restart api

# Stop services
docker compose -f docker-compose.dev.yml down
```

### 3. Accessing Services

**Development Mode:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Database: localhost:5432
- Redis: localhost:6379

**Container Names:**
- `gowa-db-dev` - PostgreSQL
- `gowa-redis-dev` - Redis
- `gowa-api-dev` - Backend API
- `gowa-app-dev` - Frontend App

## 🏗️ Building for Production

### Local Production Build

```bash
# Build production images
docker compose build

# Start production containers
docker compose up -d

# Check status
docker compose ps
```

### Production Deployment

```bash
# Use production compose file
docker compose -f docker-compose.prod.yml up -d

# Or with custom env file
docker compose -f docker-compose.prod.yml --env-file .env.production up -d
```

## 🐛 Troubleshooting

### Development Mode Issues

**Frontend not hot-reloading:**
```bash
# Restart frontend container
docker compose -f docker-compose.dev.yml restart app

# Check if volume is mounted correctly
docker compose -f docker-compose.dev.yml exec app ls -la /app
```

**Backend not hot-reloading:**
```bash
# Check Air logs
docker compose -f docker-compose.dev.yml logs -f api

# Restart backend container
docker compose -f docker-compose.dev.yml restart api
```

**Port already in use:**
```bash
# Check what's using the port
lsof -i :3000
lsof -i :8080

# Kill the process or change ports in docker-compose.dev.yml
```

### Clean Rebuild

```bash
# Stop all containers
docker compose -f docker-compose.dev.yml down

# Remove volumes (WARNING: This deletes database data)
docker compose -f docker-compose.dev.yml down -v

# Remove images
docker compose -f docker-compose.dev.yml down --rmi all

# Rebuild from scratch
docker compose -f docker-compose.dev.yml up --build
```

### Database Issues

```bash
# Access database
docker compose -f docker-compose.dev.yml exec db psql -U gowa_user -d gowa_db

# Run migrations manually
docker compose -f docker-compose.dev.yml exec api sh
# Inside container:
# Run your migration commands
```

## 📊 Performance Comparison

| Aspect | Development Mode | Production Mode |
|--------|-----------------|-----------------|
| **Initial Build** | ~30 seconds | ~5-10 minutes |
| **Rebuild** | ~5 seconds | ~5-10 minutes |
| **Hot Reload** | ✅ Yes | ❌ No |
| **Image Size** | Larger (~500MB) | Smaller (~100MB) |
| **Memory Usage** | Higher | Lower |
| **Best For** | Local development | Deployment |

## 🔐 Environment Variables

### Required Variables

```bash
# Database
DB_USER=gowa_user
DB_PASSWORD=gowa_password
DB_NAME=gowa_db

# JWT
JWT_SECRET=your-secret-key-here

# Frontend
FRONTEND_URL=http://localhost:3000
```

### Optional Variables

```bash
# Google OAuth
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

## 📝 Tips & Best Practices

### Development Mode
1. **Always use development mode for local development** - It's much faster
2. **Keep containers running** - No need to restart for code changes
3. **Use logs to debug** - `docker compose -f docker-compose.dev.yml logs -f`
4. **Commit often** - Changes are reflected immediately

### Production Mode
1. **Test production builds locally** before deploying
2. **Use environment-specific .env files**
3. **Always backup database** before updates
4. **Monitor container health** with `docker compose ps`

### General
1. **Don't commit .env files** - Use .env.example as template
2. **Use volumes for data persistence** - Database, Redis, WhatsApp sessions
3. **Clean up unused images** - `docker system prune`
4. **Monitor disk usage** - `docker system df`

## 🆘 Common Commands Cheatsheet

```bash
# Development
docker compose -f docker-compose.dev.yml up -d        # Start dev
docker compose -f docker-compose.dev.yml down         # Stop dev
docker compose -f docker-compose.dev.yml logs -f      # View logs
docker compose -f docker-compose.dev.yml restart api  # Restart backend
docker compose -f docker-compose.dev.yml restart app  # Restart frontend

# Production
docker compose up -d                                  # Start prod
docker compose down                                   # Stop prod
docker compose build                                  # Build images
docker compose up --build                             # Build and start

# Debugging
docker compose -f docker-compose.dev.yml exec api sh  # Shell into backend
docker compose -f docker-compose.dev.yml exec app sh  # Shell into frontend
docker compose -f docker-compose.dev.yml exec db psql -U gowa_user -d gowa_db  # DB access

# Cleanup
docker compose -f docker-compose.dev.yml down -v      # Stop and remove volumes
docker system prune -a                                # Clean all unused data
```

## 📚 Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Air (Go Hot Reload)](https://github.com/cosmtrek/air)
- [Nuxt Development](https://nuxt.com/docs/getting-started/introduction)
- [Echo Framework](https://echo.labstack.com/)
