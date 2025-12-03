-- Migration: Initial Schema
-- Description: Create initial database schema

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'parent',
    google_id VARCHAR(255) UNIQUE,
    auth_provider VARCHAR(50) DEFAULT 'email',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_auth_provider 
    ON users (email, auth_provider) 
    WHERE auth_provider IN ('email', 'both');

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_google 
    ON users (email) 
    WHERE auth_provider = 'google';

