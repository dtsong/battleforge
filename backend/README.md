# BattleForge Backend API

Go-based HTTP API for analyzing competitive Pokémon gameplay with support for Pokémon Showdown replays.

## Project Structure

```
backend/
├── cmd/
│   └── battleforge-api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── analysis/
│   │   ├── parser.go              # Showdown log parser
│   │   └── types.go               # BattleSummary type definitions
│   ├── db/
│   │   ├── db.go                  # Database operations
│   │   └── types.go               # Database model types
│   ├── httpapi/
│   │   ├── router.go              # Chi router setup
│   │   ├── showdown_handlers.go   # Showdown analysis endpoints
│   │   └── tcglive_handlers.go    # TCG Live analysis endpoints (future)
│   └── observability/
│       └── logging.go              # Logging utilities
├── migrations/
│   └── 001_initial_schema.sql     # Database schema
├── go.mod
├── go.sum
└── openapi.yaml                   # API specification
```

## Setup

### Prerequisites

- Go 1.22+
- PostgreSQL 13+
- `go mod` for dependency management

### Installation

1. **Install dependencies:**
   ```bash
   cd backend
   go mod download
   ```

2. **Set up the database:**
   ```bash
   # Create PostgreSQL database
   createdb battleforge

   # Run migrations
   psql battleforge < migrations/001_initial_schema.sql
   ```

3. **Build the application:**
   ```bash
   go build -o battleforge-api ./cmd/battleforge-api
   ```

### Running Locally

```bash
# Set environment variables (optional)
export BATTLEFORGE_API_ADDR=:8080
export DATABASE_URL="postgres://user:password@localhost:5432/battleforge?sslmode=disable"

# Run the server
go run ./cmd/battleforge-api
```

The API will be available at `http://localhost:8080`

### Health Check

```bash
curl http://localhost:8080/healthz
# Response: ok
```

## API Endpoints

All endpoints are documented in `openapi.yaml` and follow the OpenAPI 3.0.0 specification.

### Showdown Analysis

- **POST** `/api/showdown/analyze` - Analyze a Showdown replay
  - Supports three input methods: replayId, username, or rawLog
  - Returns structured BattleSummary with detailed analytics

- **GET** `/api/showdown/replays` - List analyzed replays
  - Filter by username, format, privacy status
  - Pagination support (limit, offset)

- **GET** `/api/showdown/replays/{replayId}` - Get specific replay analysis

### TCG Live Analysis

- **POST** `/api/tcglive/analyze` - Analyze TCG Live game (planned)

## Development

### Code Style

- Follow Go conventions (go fmt, go vet)
- Use chi for HTTP routing
- Package-level documentation for public APIs
- Comments for complex logic

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestParseShowdownLog ./internal/analysis
```

### Dependencies

- **github.com/go-chi/chi/v5** - HTTP routing
- **github.com/lib/pq** - PostgreSQL driver

## API Design

The API follows REST conventions:

- **Request Format**: JSON in request body
- **Response Format**: JSON with consistent structure
- **Error Handling**: Structured error responses with codes
- **Versioning**: Implicit in URL paths (v1 in future if needed)

### Request/Response Example

```bash
# Analyze a raw Showdown log
curl -X POST http://localhost:8080/api/showdown/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "analysisType": "rawLog",
    "rawLog": "...",
    "isPrivate": true
  }'

# Response:
{
  "status": "success",
  "battleId": "550e8400-e29b-41d4-a716-446655440000",
  "data": { /* BattleSummary */ },
  "metadata": {
    "parseTimeMs": 45,
    "analysisTimeMs": 120,
    "cached": false
  }
}
```

## Database Schema

The schema includes:

- **battles**: Main battle records
- **battle_analysis**: Computed statistics per battle
- **key_moments**: Notable events in battles
- **pokemon**, **pokemon_species**: Pokémon data
- **moves**, **items**: Reference data
- **pokemon_moves**: Pokémon move mappings
- **battle_turns**, **battle_actions**: Turn-by-turn action logs

For details, see `../DATABASE_SCHEMA.md`

## Future Enhancements

- [ ] TCG Live game export parsing
- [ ] User authentication & authorization
- [ ] Caching layer (Redis)
- [ ] Metrics collection (Prometheus)
- [ ] Distributed tracing (Jaeger)
- [ ] AI-powered battle analysis
- [ ] WebSocket support for live battle analysis

## Contributing

When adding new features:

1. Update OpenAPI spec first (`openapi.yaml`)
2. Implement handlers in appropriate package
3. Add database operations if needed
4. Add tests
5. Update this README if adding new endpoints

## License

MIT
