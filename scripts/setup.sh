#!/bin/bash

# Initial setup script for new project
# Usage: ./setup.sh

set -e

echo "=========================================="
echo "  Project Setup Script"
echo "=========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env exists
if [ -f ".env" ]; then
    echo -e "${YELLOW}⚠${NC} .env file already exists. Skipping..."
else
    echo -e "${YELLOW}[1/5]${NC} Creating .env file from template..."
    if [ -f ".env.example" ]; then
        cp .env.example .env
        echo -e "${GREEN}✓${NC} .env file created"
        echo -e "${YELLOW}⚠${NC} Please edit .env file with your configuration"
    else
        echo -e "${RED}✗${NC} .env.example not found!"
        exit 1
    fi
fi
echo ""

# Make scripts executable
echo -e "${YELLOW}[2/5]${NC} Making scripts executable..."
find . -name "*.sh" -type f -exec chmod +x {} \;
echo -e "${GREEN}✓${NC} Scripts are now executable"
echo ""

# Start Docker services
echo -e "${YELLOW}[3/5]${NC} Starting Docker services..."
if command -v docker-compose &> /dev/null || command -v docker &> /dev/null; then
    docker compose up -d
    echo -e "${GREEN}✓${NC} Docker services started"
    echo -e "${YELLOW}⚠${NC} Waiting for database to be ready..."
    sleep 5
else
    echo -e "${RED}✗${NC} Docker not found. Please install Docker first."
    exit 1
fi
echo ""

# Run migrations
echo -e "${YELLOW}[4/5]${NC} Running database migrations..."
if [ -f "backend/scripts/run-migrations.sh" ]; then
    cd backend && ./scripts/run-migrations.sh docker && cd ..
    echo -e "${GREEN}✓${NC} Migrations completed"
else
    echo -e "${YELLOW}⚠${NC} Migration script not found, skipping..."
fi
echo ""

# Install frontend dependencies (if not in Docker)
echo -e "${YELLOW}[5/5]${NC} Checking frontend dependencies..."
if [ -d "frontend/node_modules" ]; then
    echo -e "${GREEN}✓${NC} Frontend dependencies already installed"
else
    echo -e "${YELLOW}⚠${NC} Frontend dependencies will be installed by Docker"
fi
echo ""

echo "=========================================="
echo -e "${GREEN}Setup completed!${NC}"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Edit .env file with your configuration"
echo "2. Update Google OAuth credentials in .env"
echo "3. Generate a secure JWT_SECRET (openssl rand -base64 32)"
echo "4. Access frontend: http://localhost:3000"
echo "5. Access backend API: http://localhost:8080"
echo ""
echo "Useful commands:"
echo "  docker compose logs -f          # View logs"
echo "  docker compose ps               # Check service status"
echo "  docker compose down             # Stop services"
echo ""

