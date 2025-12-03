#!/bin/bash

# Script to validate environment variables before deployment
# Run this before deploying to production

set -e

echo "🔍 Validating environment variables..."
echo ""

# Colors for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

ERRORS=0
WARNINGS=0

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${RED}❌ ERROR: .env file not found${NC}"
    echo "   Please copy .env.example to .env and configure it"
    exit 1
fi

# Source the .env file
export $(cat .env | grep -v '^#' | xargs)

# Function to check required variable
check_required() {
    local var_name=$1
    local var_value="${!var_name}"
    
    if [ -z "$var_value" ]; then
        echo -e "${RED}❌ ERROR: $var_name is not set${NC}"
        ERRORS=$((ERRORS + 1))
        return 1
    fi
    return 0
}

# Function to check for placeholder values
check_placeholder() {
    local var_name=$1
    local var_value="${!var_name}"
    local placeholder=$2
    
    if [[ "$var_value" == *"$placeholder"* ]]; then
        echo -e "${YELLOW}⚠️  WARNING: $var_name contains placeholder value '$placeholder'${NC}"
        WARNINGS=$((WARNINGS + 1))
        return 1
    fi
    return 0
}

# Function to check minimum length
check_min_length() {
    local var_name=$1
    local var_value="${!var_name}"
    local min_length=$2
    
    if [ ${#var_value} -lt $min_length ]; then
        echo -e "${RED}❌ ERROR: $var_name must be at least $min_length characters (current: ${#var_value})${NC}"
        ERRORS=$((ERRORS + 1))
        return 1
    fi
    return 0
}

echo "Checking required variables..."
echo ""

# Check critical variables
check_required "JWT_SECRET"
check_required "DB_HOST"
check_required "DB_USER"
check_required "DB_PASSWORD"
check_required "DB_NAME"

echo ""
echo "Checking JWT_SECRET security..."
echo ""

# Check JWT_SECRET length
if check_required "JWT_SECRET"; then
    check_min_length "JWT_SECRET" 32
    check_placeholder "JWT_SECRET" "change_me"
    check_placeholder "JWT_SECRET" "secret"
fi

echo ""
echo "Checking database security..."
echo ""

# Check DB_PASSWORD
if check_required "DB_PASSWORD"; then
    check_min_length "DB_PASSWORD" 12
    check_placeholder "DB_PASSWORD" "change_me"
    check_placeholder "DB_PASSWORD" "password"
fi

echo ""
echo "Checking production settings..."
echo ""

# Check ENV setting
if [ "$ENV" = "production" ]; then
    echo -e "${GREEN}✅ ENV is set to production${NC}"
    
    # Additional checks for production
    if [ -z "$CORS_ALLOWED_ORIGINS" ]; then
        echo -e "${YELLOW}⚠️  WARNING: CORS_ALLOWED_ORIGINS not set for production${NC}"
        WARNINGS=$((WARNINGS + 1))
    elif [[ "$CORS_ALLOWED_ORIGINS" == *"localhost"* ]]; then
        echo -e "${YELLOW}⚠️  WARNING: CORS_ALLOWED_ORIGINS contains localhost in production${NC}"
        WARNINGS=$((WARNINGS + 1))
    fi
else
    echo -e "${YELLOW}ℹ️  ENV is set to: $ENV (not production)${NC}"
fi

echo ""
echo "Checking optional API keys..."
echo ""

# Check optional but recommended variables
if [ -z "$GEMINI_API_KEY" ] || [[ "$GEMINI_API_KEY" == *"your_"* ]]; then
    echo -e "${YELLOW}⚠️  WARNING: GEMINI_API_KEY not configured (AI features will be disabled)${NC}"
    WARNINGS=$((WARNINGS + 1))
fi

if [ -z "$GOOGLE_CLIENT_ID" ] || [[ "$GOOGLE_CLIENT_ID" == *"your-"* ]]; then
    echo -e "${YELLOW}ℹ️  INFO: Google OAuth not configured (optional)${NC}"
fi

echo ""
echo "================================"
echo "Validation Summary"
echo "================================"

if [ $ERRORS -gt 0 ]; then
    echo -e "${RED}❌ Found $ERRORS error(s)${NC}"
fi

if [ $WARNINGS -gt 0 ]; then
    echo -e "${YELLOW}⚠️  Found $WARNINGS warning(s)${NC}"
fi

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ All checks passed!${NC}"
    echo ""
    echo "Your environment is ready for deployment."
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  Validation passed with warnings${NC}"
    echo ""
    echo "You can proceed, but please review the warnings above."
    exit 0
else
    echo -e "${RED}❌ Validation failed${NC}"
    echo ""
    echo "Please fix the errors above before deploying."
    echo ""
    echo "Quick fixes:"
    echo "  - Generate JWT_SECRET: openssl rand -base64 32"
    echo "  - Generate DB_PASSWORD: openssl rand -base64 24"
    exit 1
fi
