#!/bin/bash
# Quick start script for Docker Compose
# Usage: ./scripts/docker-quickstart.sh [up|down|logs|restart|clean]

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}=== $1 ===${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

show_help() {
    cat << EOF
Pokemon VGC Corner - Docker Quick Start

Usage: ./scripts/docker-quickstart.sh [COMMAND]

Commands:
    up              Start all services with rebuild
    build           Build all Docker images
    down            Stop all services (keep volumes)
    restart         Restart all services
    logs            Show logs from all services
    logs-backend    Show backend logs only
    logs-frontend   Show frontend logs only
    logs-db         Show database logs only
    status          Show service status
    clean           Remove all containers and images (keeps volumes)
    clean-all       Remove everything including volumes (⚠️  loses database data)
    shell-backend   Open shell in backend container
    shell-frontend  Open shell in frontend container
    shell-db        Open PostgreSQL shell in database container
    test-api        Test API health endpoint
    test-db         Test database connectivity
    help            Show this help message

Examples:
    ./scripts/docker-quickstart.sh up        # Start all services
    ./scripts/docker-quickstart.sh logs      # View all logs
    ./scripts/docker-quickstart.sh down      # Stop services
    ./scripts/docker-quickstart.sh shell-db  # Access database

Configuration:
    Environment variables from docker-compose.yml:
    - POSTGRES_USER: vgccorner
    - DB_HOST: postgres (service name)
    - API: http://localhost:8080
    - Frontend: http://localhost:3000
EOF
}

# Main commands
cmd_up() {
    print_header "Starting Pokemon VGC Corner"
    docker-compose up --build
}

cmd_build() {
    print_header "Building Docker images"
    docker-compose build
    print_success "Build complete"
}

cmd_down() {
    print_header "Stopping services"
    docker-compose down
    print_success "Services stopped"
}

cmd_restart() {
    print_header "Restarting services"
    docker-compose restart
    print_success "Services restarted"
}

cmd_logs() {
    print_header "Showing all logs (Ctrl+C to exit)"
    docker-compose logs -f
}

cmd_logs_backend() {
    print_header "Backend logs (Ctrl+C to exit)"
    docker-compose logs -f backend
}

cmd_logs_frontend() {
    print_header "Frontend logs (Ctrl+C to exit)"
    docker-compose logs -f frontend
}

cmd_logs_db() {
    print_header "Database logs (Ctrl+C to exit)"
    docker-compose logs -f postgres
}

cmd_status() {
    print_header "Service Status"
    docker-compose ps
    echo ""
    print_header "Network Information"
    docker network inspect vgccorner_vgccorner-network 2>/dev/null | grep -A 5 "Containers" || print_warning "Network not found"
}

cmd_clean() {
    print_header "Cleaning up containers and images"
    docker-compose down
    docker image rm vgccorner_backend vgccorner_frontend 2>/dev/null || true
    print_success "Cleanup complete"
}

cmd_clean_all() {
    print_warning "This will remove ALL containers, images, AND database data!"
    read -p "Are you sure? (yes/no): " -r
    echo
    if [[ $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
        docker-compose down -v
        docker image rm vgccorner_backend vgccorner_frontend 2>/dev/null || true
        print_success "Complete cleanup done"
    else
        print_warning "Cleanup cancelled"
    fi
}

cmd_shell_backend() {
    print_header "Opening shell in backend container"
    docker-compose exec backend /bin/sh
}

cmd_shell_frontend() {
    print_header "Opening shell in frontend container"
    docker-compose exec frontend /bin/sh
}

cmd_shell_db() {
    print_header "Opening PostgreSQL shell"
    docker-compose exec postgres psql -U vgccorner -d vgccorner
}

cmd_test_api() {
    print_header "Testing API health endpoint"
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        print_success "API is responding at http://localhost:8080"
        echo ""
        curl -s http://localhost:8080/health | head -c 200
        echo ""
    else
        print_error "API is not responding"
        echo "Try: ./scripts/docker-quickstart.sh logs-backend"
    fi
}

cmd_test_db() {
    print_header "Testing database connectivity"
    if docker-compose exec postgres pg_isready -U vgccorner -d vgccorner > /dev/null 2>&1; then
        print_success "Database is ready"
        echo ""
        docker-compose exec postgres psql -U vgccorner -d vgccorner -c "SELECT version();" 2>/dev/null
    else
        print_error "Database is not responding"
        echo "Try: ./scripts/docker-quickstart.sh logs-db"
    fi
}

# Parse command
COMMAND="${1:-help}"

case "$COMMAND" in
    up)
        cmd_up
        ;;
    build)
        cmd_build
        ;;
    down)
        cmd_down
        ;;
    restart)
        cmd_restart
        ;;
    logs)
        cmd_logs
        ;;
    logs-backend)
        cmd_logs_backend
        ;;
    logs-frontend)
        cmd_logs_frontend
        ;;
    logs-db)
        cmd_logs_db
        ;;
    status)
        cmd_status
        ;;
    clean)
        cmd_clean
        ;;
    clean-all)
        cmd_clean_all
        ;;
    shell-backend)
        cmd_shell_backend
        ;;
    shell-frontend)
        cmd_shell_frontend
        ;;
    shell-db)
        cmd_shell_db
        ;;
    test-api)
        cmd_test_api
        ;;
    test-db)
        cmd_test_db
        ;;
    help)
        show_help
        ;;
    *)
        print_error "Unknown command: $COMMAND"
        echo ""
        show_help
        exit 1
        ;;
esac
