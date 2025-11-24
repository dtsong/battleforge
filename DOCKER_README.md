# Docker Configuration - Pokemon VGC Corner

This directory contains all Docker configuration files for containerizing the Pokemon VGC Corner application.

## Files Overview

### `docker-compose.yml`
Main orchestration file that defines all services:
- **postgres**: PostgreSQL 16 database
- **backend**: Go REST API service
- **frontend**: Next.js web application

All services are interconnected through a Docker network and configured with health checks and dependency management.

### `backend/Dockerfile`
Multi-stage build for the Go backend:
- **Stage 1 (builder)**: Compiles Go code with dependencies
- **Stage 2 (runtime)**: Minimal Alpine Linux runtime

Includes:
- PostgreSQL driver (`lib/pq`) support
- Health check endpoint
- Migration file copying
- Port 8080 exposure

### `frontend/Dockerfile`
Multi-stage build for the Next.js frontend:
- **Stage 1 (deps)**: Installs dependencies
- **Stage 2 (builder)**: Builds Next.js application
- **Stage 3 (runner)**: Optimized production runtime

Includes:
- Multi-package manager support (npm, yarn, pnpm)
- Production-only dependency installation
- Health check endpoint
- Port 3000 exposure

### `.dockerignore`
Excludes unnecessary files from Docker build context:
- Source control files (.git)
- IDE configurations
- Build artifacts
- Documentation
- Test files

Reduces build context size and improves build performance.

### `.env.example`
Template for environment variables used by docker-compose:
- Database credentials
- Service ports
- Logging configuration
- API URLs

**Note**: This is a template. Copy to `.env` for local development or use docker-compose environment variables.

## Quick Start

```bash
# Start all services
docker-compose up --build

# In another terminal, access the application
curl http://localhost:3000        # Frontend
curl http://localhost:8080/health # Backend health check
psql -h localhost -U vgccorner -d vgccorner  # Database
```

See `DOCKER_SETUP.md` for complete documentation and troubleshooting.

## Service Architecture

```
┌─────────────────────────────────────────────────┐
│        Docker Network: vgccorner-network         │
├─────────────────────────────────────────────────┤
│                                                 │
│  ┌──────────────┐  ┌──────────────────────┐   │
│  │  PostgreSQL  │  │  Backend (Go)        │   │
│  │  :5432       │  │  :8080               │   │
│  │  vgccorner   │──│  ▪ Health check      │   │
│  │  Database    │  │  ▪ API endpoints     │   │
│  └──────────────┘  │  ▪ DB integration    │   │
│       ▲            └──────────────────────┘   │
│       │                     ▲                  │
│       │                     │                  │
│       └─────────────────────┴──────┐           │
│                                    │           │
│                          ┌──────────┴────────┐ │
│                          │ Frontend (React)  │ │
│                          │ :3000             │ │
│                          │ ▪ Web UI          │ │
│                          │ ▪ Battle Analysis │ │
│                          └───────────────────┘ │
│                                                 │
└─────────────────────────────────────────────────┘

Host:
├─ localhost:3000   → Frontend
├─ localhost:8080   → Backend API
└─ localhost:5432   → PostgreSQL
```

## Environment Variables

### Shared (docker-compose.yml)
```bash
# PostgreSQL
POSTGRES_USER=vgccorner
POSTGRES_PASSWORD=vgccorner_dev_password
POSTGRES_DB=vgccorner
```

### Backend
```bash
DB_HOST=postgres              # Service name in Docker network
DB_PORT=5432
DB_USER=vgccorner
DB_PASSWORD=vgccorner_dev_password
DB_NAME=vgccorner
DB_SSL_MODE=disable           # Local dev only
SERVER_PORT=8080
LOG_LEVEL=info
```

### Frontend
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
NODE_ENV=production
```

## Volumes

### `postgres_data`
Named volume for PostgreSQL data persistence:
- **Location**: Docker-managed location
- **Persistence**: Survives container restarts
- **Reset**: `docker-compose down -v` deletes volume
- **Backup**: Use `pg_dump` to export data

### `backend/migrations`
Bind mount for database migrations:
- **Location**: `./backend/migrations/` on host
- **Purpose**: Auto-loaded on PostgreSQL startup
- **Format**: `.sql` files

## Health Checks

Each service has a health check to ensure readiness:

### PostgreSQL
- **Interval**: 10s
- **Command**: `pg_isready -U vgccorner -d vgccorner`
- **Timeout**: 5s
- **Retries**: 5

### Backend
- **Interval**: 30s
- **Command**: HTTP GET `/health`
- **Timeout**: 10s
- **Retries**: 3
- **Start Period**: 15s (grace time before first check)

### Frontend
- **Interval**: 30s
- **Command**: HTTP GET `/` (root page)
- **Timeout**: 10s
- **Retries**: 3
- **Start Period**: 15s

## Dependency Management

Services start in order:
1. **PostgreSQL** (no dependencies)
2. **Backend** (waits for PostgreSQL health check)
3. **Frontend** (waits for backend service - not health check)

This ensures database is ready before API starts, and API is running before frontend connects.

## Networking

All containers are connected via `vgccorner-network` (bridge network):
- Services communicate using service names (DNS resolution)
- Ports are not exposed between containers (only via -p mappings)
- Example: Backend connects to `postgres:5432`, not `localhost:5432`

## Building

### Build All Services
```bash
docker-compose build
```

### Build Specific Service
```bash
docker-compose build backend
docker-compose build frontend
```

### Build Without Cache
```bash
docker-compose build --no-cache
```

### View Build Logs
```bash
docker-compose build --verbose backend
```

## Performance Considerations

### Image Sizes
- **Backend**: ~50MB (multi-stage, Alpine base)
- **Frontend**: ~200MB (Node.js 22, Next.js build)
- **PostgreSQL**: ~250MB (official image)
- **Total**: ~500MB (all three)

### Layer Caching
- Dependencies are installed in separate layers
- Changes to source code don't invalidate dependency cache
- Significantly speeds up iterative builds

### Build Time
- First build: 3-5 minutes (downloads base images)
- Subsequent builds: 30-60s (uses cached layers)
- Rebuild changed layer: 10-30s

## Security Notes

### Development Configuration
⚠️ **NOT FOR PRODUCTION**
- Database password is simple
- SSL is disabled
- Logging is verbose
- No authentication on API endpoints

### Production Configuration
See `DOCKER_SETUP.md` Production Deployment section for:
- Strong password generation
- SSL/TLS configuration
- External database usage
- Reverse proxy setup
- Backup strategies

## Debugging

### Access Container Shell
```bash
docker-compose exec postgres /bin/bash
docker-compose exec backend /bin/sh
docker-compose exec frontend /bin/sh
```

### View Full Logs
```bash
docker-compose logs postgres
docker-compose logs backend
docker-compose logs frontend
```

### Check Network Configuration
```bash
docker network inspect vgccorner_vgccorner-network
```

### Check Volume Mounts
```bash
docker inspect vgccorner-postgres | grep -A 20 Mounts
```

## Cleanup

### Stop Services (Keep Volumes)
```bash
docker-compose stop
```

### Remove Containers (Keep Volumes)
```bash
docker-compose down
```

### Remove Everything (Delete Volumes)
```bash
docker-compose down -v
```

### Remove Dangling Images
```bash
docker image prune -a
```

## Integration with Development

### Local Development Without Docker
```bash
# Terminal 1: Start PostgreSQL
docker-compose up postgres

# Terminal 2: Run backend locally
cd backend && go run ./cmd/vgccorner-api/main.go

# Terminal 3: Run frontend locally
cd frontend && npm run dev
```

### Hybrid Mode
```bash
# Start only database in Docker
docker-compose up postgres

# Run backend and frontend locally
# (Useful for faster iteration)
```

## References

- See `DOCKER_SETUP.md` for comprehensive setup and troubleshooting guide
- See `ARCHITECTURE.md` for system design and deployment strategy
- See `DATABASE_SCHEMA.md` for database structure details
