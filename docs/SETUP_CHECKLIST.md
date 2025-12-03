# Setup Checklist

Use this checklist when setting up a new project from this template.

## Initial Setup

- [ ] Generate project using `./project-template/scripts/generate-project.sh <project-name>`
- [ ] Navigate to project directory
- [ ] Run `make setup` or `./scripts/setup.sh`
- [ ] Review and edit `.env` file with your configuration

## Environment Configuration

- [ ] Generate secure `JWT_SECRET` (use `openssl rand -base64 32`)
- [ ] Update database credentials in `.env`
- [ ] Set `DB_USER`, `DB_PASSWORD`, `DB_NAME` to your values
- [ ] Configure `FRONTEND_URL` for your domain
- [ ] Set `CORS_ALLOWED_ORIGINS` with your frontend URL(s)

## Google OAuth Setup

- [ ] Create project in [Google Cloud Console](https://console.cloud.google.com/)
- [ ] Enable Google+ API
- [ ] Create OAuth 2.0 credentials
- [ ] Set authorized redirect URIs:
  - Development: `http://localhost:8080/api/auth/google/callback`
  - Production: `https://yourdomain.com/api/auth/google/callback`
- [ ] Update `.env` with `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET`
- [ ] Update `GOOGLE_REDIRECT_URL` in `.env`

## Database Setup

- [ ] Verify database container is running: `docker compose ps`
- [ ] Run migrations: `make migrate` or `cd backend && ./scripts/run-migrations.sh docker`
- [ ] Verify database connection: `make shell-db`
- [ ] Test database queries

## Development

- [ ] Start services: `make up`
- [ ] Check logs: `make logs`
- [ ] Verify frontend: http://localhost:3000
- [ ] Verify backend API: http://localhost:8080/health
- [ ] Test authentication flow
- [ ] Test Google OAuth login

## Code Customization

- [ ] Update project name in all files (if not done by generator)
- [ ] Customize frontend pages and components
- [ ] Add your API endpoints in `backend/handlers/`
- [ ] Create database models in `backend/models/`
- [ ] Add database migrations as needed
- [ ] Update README.md with project-specific information

## Production Deployment

- [ ] Set up VPS/server
- [ ] Install Docker and Docker Compose
- [ ] Clone repository to server
- [ ] Copy `.env.example` to `.env` and configure
- [ ] Set production environment variables
- [ ] Build and start: `make prod-up`
- [ ] Set up Nginx (copy `nginx.conf.example`)
- [ ] Configure SSL with Let's Encrypt
- [ ] Test production deployment
- [ ] Set up automated backups
- [ ] Configure monitoring and logging

## Security

- [ ] Review `.gitignore` to ensure sensitive files are excluded
- [ ] Verify `.env` is not committed to git
- [ ] Use strong passwords for database
- [ ] Enable HTTPS in production
- [ ] Review CORS settings
- [ ] Set up firewall rules
- [ ] Enable database SSL in production (if needed)
- [ ] Review JWT token expiration settings
- [ ] Set up rate limiting (if needed)

## Testing

- [ ] Test user registration
- [ ] Test user login
- [ ] Test Google OAuth flow
- [ ] Test protected API endpoints
- [ ] Test error handling
- [ ] Test database migrations
- [ ] Test backup and restore process

## Documentation

- [ ] Update `README.md` with project-specific information
- [ ] Document API endpoints
- [ ] Document environment variables
- [ ] Document deployment process
- [ ] Create user documentation (if needed)

## Final Checks

- [ ] All services are running: `docker compose ps`
- [ ] No errors in logs: `docker compose logs`
- [ ] Frontend is accessible
- [ ] Backend API is responding
- [ ] Authentication is working
- [ ] Database is accessible
- [ ] Backups are configured
- [ ] Monitoring is set up (if applicable)

---

**Note:** This checklist should be customized based on your specific project requirements.

