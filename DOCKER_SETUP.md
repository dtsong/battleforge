# Docker Setup Guide

This guide explains the Docker configuration for running the Pokemon VGC Corner application locally with all services containerized.

## Architecture Overview

The application consists of three main services running in Docker containers:

- **PostgreSQL (port 5432)**: Database for persistent storage
- **Backend API (port 8080)**: Go REST API for battle analysis
- **Frontend (port 3000)**: Next.js web application

All services communicate through a Docker network and are orchestrated with Docker Compose.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- 2GB+ available disk space for Docker volumes and images

## Quick Start

### 1. Build and Start All Services

```bash
docker-compose up --build
```

This command will:
- Build the backend Docker image from `backend/Dockerfile`
- Build the frontend Docker image from `frontend/Dockerfile`
- Create and start PostgreSQL container
- Run database migrations automatically
- Start all services in dependency order

### 2. Wait for Services to be Healthy

Docker health checks will verify each service:
- PostgreSQL: Ready within 10-50s
- Backend: Ready within 15-30s (after database is ready)
- Frontend: Ready within 30-60s (after backend is ready)

### 3. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Database**: localhost:5432 (postgres:vgccorner_dev_password)

## Service Details

### PostgreSQL

**Container**: vgccorner-postgres
**Port**: 5432
**Credentials**:
  - User: `vgccorner`
  - Password: `vgccorner_dev_password`
  - Database: `vgccorner`

**Volumes**:
- `postgres_data`: Persistent data storage (survives container restarts)
- `backend/migrations`: SQL migration files auto-loaded on startup

**Health Check**: Verifies connectivity to the database every 10s

### Backend API

**Container**: vgccorner-backend
**Port**: 8080
**Build**: Multi-stage Go build
  - Stage 1: Compiles Go code in golang:1.25-alpine
  - Stage 2: Runs in alpine:latest (minimal image)

**Environment Variables**:
```
DB_HOST=postgres              # Docker service hostname
DB_PORT=5432                  # PostgreSQL port
DB_USER=vgccorner             # Database user
DB_PASSWORD=vgccorner_dev_password  # Database password
DB_NAME=vgccorner             # Database name
DB_SSL_MODE=disable           # SSL disabled for local dev
SERVER_PORT=8080              # API listening port
LOG_LEVEL=info                # Logging verbosity
```

**Health Check**: Checks `/health` endpoint every 30s
**Restart Policy**: Restarts unless stopped explicitly

**Dependencies**:
- Waits for PostgreSQL to be healthy before starting
- Runs database migrations on first start

### Frontend

**Container**: vgccorner-frontend
**Port**: 3000
**Build**: Multi-stage Node.js build
  - Stage 1: Dependencies installation (node:22-alpine)
  - Stage 2: Build Next.js app
  - Stage 3: Production runtime (minimal image)

**Environment Variables**:
```
NEXT_PUBLIC_API_URL=http://localhost:8080  # Backend API URL
NODE_ENV=production                         # Production mode
```

**Health Check**: Checks root URL every 30s
**Restart Policy**: Restarts unless stopped explicitly

**Dependencies**:
- Waits for backend service (not health check) before starting

## Common Commands

### View Running Containers

```bash
docker-compose ps
```

Output shows:
- Container names
- Image names
- Status (Up, Exited)
- Port mappings

### View Logs

View all services:
```bash
docker-compose logs -f
```

View specific service (last 100 lines):
```bash
docker-compose logs --tail=100 -f backend
docker-compose logs --tail=100 -f frontend
docker-compose logs --tail=100 -f postgres
```

### Stop Services

```bash
docker-compose stop
```

Services are paused but not removed. Restart with `docker-compose start`.

### Stop and Remove Everything

```bash
docker-compose down
```

Removes containers and networks but keeps volumes (database data persists).

### Remove Everything Including Data

```bash
docker-compose down -v
```

**WARNING**: This deletes the PostgreSQL volume. Database data will be lost.

### Rebuild Images

```bash
docker-compose up --build
```

Use `--no-cache` for complete rebuild:
```bash
docker-compose build --no-cache
```

## Database Management

### Access PostgreSQL Directly

```bash
docker-compose exec postgres psql -U vgccorner -d vgccorner
```

Then in psql:
```sql
\dt                    -- List all tables
\d battles             -- Describe battles table
SELECT * FROM battles; -- Query battles
\q                     -- Exit
```

### Reset Database

```bash
# Keep data (if needed)
docker-compose restart postgres

# Completely reset (delete data)
docker-compose down -v
docker-compose up
```

### Import Data

```bash
docker-compose exec postgres psql -U vgccorner -d vgccorner < backup.sql
```

### Export Database

```bash
docker-compose exec postgres pg_dump -U vgccorner vgccorner > backup.sql
```

## API Testing

### Health Check

```bash
curl http://localhost:8080/health
```

### Analyze Battle

```bash
curl -X POST http://localhost:8080/api/showdown/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "analysisType": "raw_log",
    "rawLog": "[2025-01-01 12:00:00] Battle started...",
    "metadata": {
      "format": "vgc2025",
      "timestamp": "2025-01-01T12:00:00Z"
    }
  }'
```

## Troubleshooting

### Services Won't Start

Check logs:
```bash
docker-compose logs -f
```

Common issues:
- Port already in use: Change ports in docker-compose.yml
- Insufficient disk space: Free up space or increase Docker disk limit
- Docker daemon not running: Start Docker Desktop

### Database Connection Failed

Verify PostgreSQL is healthy:
```bash
docker-compose ps postgres
```

Should show status "Up (healthy)".

If not healthy:
```bash
docker-compose logs postgres
docker-compose restart postgres
```

### Migrations Not Running

Check if migrations folder path is correct:
```bash
docker-compose exec postgres ls /docker-entrypoint-initdb.d/
```

Should list SQL files from backend/migrations.

### Backend Can't Connect to Database

Verify network connectivity:
```bash
docker-compose exec backend ping postgres
```

Check database credentials match environment variables:
```bash
docker-compose logs backend | grep -i database
```

### Frontend Can't Reach Backend API

Verify backend is running:
```bash
curl http://localhost:8080/health
```

Check NEXT_PUBLIC_API_URL in docker-compose.yml matches actual backend URL.

## Performance Optimization

### Database

For production, optimize PostgreSQL in docker-compose.yml:

```yaml
postgres:
  command:
    - "postgres"
    - "-c"
    - "shared_buffers=256MB"
    - "-c"
    - "max_connections=100"
```

### Caching

Frontend builds can be optimized by caching layers:

```bash
docker-compose build --cache frontend
```

### Resource Limits

Add resource constraints in docker-compose.yml:

```yaml
backend:
  deploy:
    resources:
      limits:
        cpus: '1'
        memory: 512M
      reservations:
        cpus: '0.5'
        memory: 256M
```

## Production Deployment

For production deployment:

1. **Use environment variables** from `.env.production` file:
   ```bash
   docker-compose --env-file .env.production up -d
   ```

2. **Change database credentials** to strong passwords

3. **Enable SSL/TLS** for database connections:
   ```
   DB_SSL_MODE=require
   ```

4. **Use external PostgreSQL** instead of container for better reliability

5. **Add reverse proxy** (nginx) for SSL/TLS termination

6. **Enable persistent backups** of database volume

7. **Remove debug logging**:
   ```
   LOG_LEVEL=warn
   ```

See ARCHITECTURE.md for more production guidelines.

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Reference](https://docs.docker.com/compose/compose-file/)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [Node.js Docker Image](https://hub.docker.com/_/node)
- [Go Docker Image](https://hub.docker.com/_/golang)

## Next Steps

1. Run `docker-compose up --build` to start all services
2. Visit http://localhost:3000 to access the frontend
3. Test the API at http://localhost:8080/api/showdown/analyze
4. Use database tools to query PostgreSQL at localhost:5432

All services are now containerized and ready for development or deployment!
