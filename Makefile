.PHONY: help setup dev dev-down dev-logs dev-restart dev-build up down logs ps restart build clean migrate backup

# Default target
help:
	@echo "🚀 Gowa UMKM - Docker Commands"
	@echo ""
	@echo "📦 DEVELOPMENT MODE (Recommended for local development):"
	@echo "  make dev           - Start all services in development mode (hot-reload)"
	@echo "  make dev-down      - Stop development services"
	@echo "  make dev-logs      - View development logs"
	@echo "  make dev-restart   - Restart development services"
	@echo "  make dev-build     - Rebuild development containers"
	@echo ""
	@echo "🏭 PRODUCTION MODE:"
	@echo "  make up            - Start all services in production mode"
	@echo "  make down          - Stop production services"
	@echo "  make logs          - View production logs"
	@echo "  make restart       - Restart production services"
	@echo "  make build         - Build production containers"
	@echo ""
	@echo "🛠️  UTILITIES:"
	@echo "  make setup         - Initial project setup"
	@echo "  make clean         - Stop and remove containers, volumes"
	@echo "  make migrate       - Run database migrations"
	@echo "  make backup        - Backup database"
	@echo "  make shell-api     - Open shell in API container"
	@echo "  make shell-db      - Open psql in database"
	@echo ""
	@echo "💡 TIP: Use 'make dev' for local development (much faster!)"

# Initial setup
setup:
	@chmod +x scripts/*.sh
	@./scripts/setup.sh

# Development Mode (Hot-reload)
dev:
	@echo "🚀 Starting development mode with hot-reload..."
	docker compose -f docker-compose.dev.yml up -d

dev-down:
	docker compose -f docker-compose.dev.yml down

dev-logs:
	docker compose -f docker-compose.dev.yml logs -f

dev-restart:
	docker compose -f docker-compose.dev.yml restart

dev-build:
	docker compose -f docker-compose.dev.yml up -d --build

dev-clean:
	docker compose -f docker-compose.dev.yml down -v

# Production Mode
up:
	@echo "🏭 Starting production mode..."
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

restart:
	docker compose restart

build:
	docker compose build

clean:
	docker compose down -v
	docker system prune -f

# Database
migrate:
	cd backend && ./scripts/run-migrations.sh docker

backup:
	./scripts/backup-db.sh

# Shell access (works for both dev and prod)
shell-api:
	@if docker ps | grep -q gowa-api-dev; then \
		docker compose -f docker-compose.dev.yml exec api sh; \
	else \
		docker compose exec api sh; \
	fi

shell-app:
	@if docker ps | grep -q gowa-app-dev; then \
		docker compose -f docker-compose.dev.yml exec app sh; \
	else \
		docker compose exec app sh; \
	fi

shell-db:
	@if docker ps | grep -q gowa-db-dev; then \
		docker compose -f docker-compose.dev.yml exec db psql -U gowa_user -d gowa_db; \
	else \
		docker compose exec db psql -U gowa_user -d gowa_db; \
	fi

# Production deployment
prod-up:
	docker compose -f docker-compose.prod.yml up -d --build

prod-down:
	docker compose -f docker-compose.prod.yml down

prod-logs:
	docker compose -f docker-compose.prod.yml logs -f

prod-restart:
	docker compose -f docker-compose.prod.yml restart
