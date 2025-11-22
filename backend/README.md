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

## API Specification

The complete API specification is documented in `openapi.yaml` using the **OpenAPI 3.0.0** standard. The specification includes:

- All endpoint definitions with request/response schemas
- Error codes and responses
- Example payloads for each endpoint
- Type definitions for all request/response bodies

The API follows REST conventions with consistent JSON request/response formats.

### API Endpoints

#### Health Check
- **GET** `/healthz` - API health status
  - Returns: Plain text "ok"
  - Used for load balancer health checks

#### Showdown Analysis

**POST** `/api/showdown/analyze` - Analyze a Pokémon Showdown replay
- Supports three input methods (discriminator-based):
  1. **By Replay ID**: Fetch and analyze from Showdown servers
     ```json
     {
       "analysisType": "replayId",
       "replayId": "gen9vgc2025reghbo3-2481642254",
       "isPrivate": false
     }
     ```
  2. **By Username**: Fetch recent battles for a player
     ```json
     {
       "analysisType": "username",
       "username": "Player1",
       "format": "gen9vgc2025reghbo3",
       "isPrivate": false,
       "limit": 5
     }
     ```
  3. **By Raw Log**: Analyze raw battle log (pipe-delimited format)
     ```json
     {
       "analysisType": "rawLog",
       "rawLog": "|j|☆Player1\n|j|☆Player2\n|turn|1\n...",
       "isPrivate": true
     }
     ```

- Returns: `AnalyzeShowdownResponse` containing:
  - `battleId`: UUID of stored analysis
  - `data`: Complete `BattleSummary` with battle analysis
  - `metadata`: Parse time, analysis time, cache status
  - `status`: "success" or "error"

- Status Codes:
  - `200`: Successfully analyzed
  - `400`: Invalid request (missing required fields, invalid JSON)
  - `404`: Replay/user not found
  - `500`: Parse error or internal error

**GET** `/api/showdown/replays` - List analyzed replays
- Query Parameters:
  - `username` (string): Filter by player name
  - `format` (string): Filter by battle format (e.g., "gen9vgc2025reghbo3")
  - `isPrivate` (boolean): Filter by privacy status
  - `limit` (integer): Max results (1-100, default 10)
  - `offset` (integer): Pagination offset (default 0)

- Returns: `ListReplaysResponse` with:
  - `battles`: Array of BattleSummary objects
  - `total`: Total matching replays
  - `limit`/`offset`: Pagination info

**GET** `/api/showdown/replays/{replayId}` - Get specific replay analysis
- Path Parameter: `replayId` (string) - The replay UUID or Showdown ID
- Returns: `AnalyzeShowdownResponse` with full BattleSummary

#### TCG Live Analysis

**POST** `/api/tcglive/analyze` - Analyze TCG Live game (planned)
- Status: Currently returns `501 NOT_IMPLEMENTED`
- Future: Will analyze TCG Live game exports similar to Showdown

## Development

### Code Style

- Follow Go conventions (go fmt, go vet)
- Use chi for HTTP routing
- Package-level documentation for public APIs
- Comments for complex logic

### Request Validation

All API requests go through validation in the handlers (`internal/httpapi/showdown_handlers.go`):

1. **Parse JSON**: Request body is unmarshaled into appropriate request struct
2. **Validate discriminator**: Check `analysisType` field to determine input method
3. **Validate required fields**: Ensure all fields needed for that input type are present
4. **Sanitize inputs**: Clean player names, validate format strings, etc.
5. **Business logic validation**: Check if format exists, user is known, etc.

Example validation flow for raw log analysis:
```
Request → Parse JSON → Check analysisType=rawLog →
Validate rawLog field present → Check log is not empty →
Pass to parser → Return BattleSummary or error
```

### Data Flow Summary

```
HTTP Request
    ↓
Handler Validation
    ↓
Parser (raw log → BattleSummary)
    ↓
Database Storage (optional)
    ↓
JSON Response
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestParseShowdownLog ./internal/analysis

# Use Makefile for common tasks
make test              # Run all tests
make test-coverage     # Show coverage %
make test-coverage-html # Generate HTML report
```

See `TESTING.md` and `TEST_COVERAGE.md` for comprehensive testing documentation.

### Dependencies

- **github.com/go-chi/chi/v5** - HTTP routing
- **github.com/lib/pq** - PostgreSQL driver

## API Design

The API follows REST conventions:

- **Request Format**: JSON in request body
- **Response Format**: JSON with consistent structure
- **Error Handling**: Structured error responses with codes
- **Versioning**: Implicit in URL paths (v1 in future if needed)

### Request/Response Examples

#### Example 1: Analyze by Raw Log

```bash
curl -X POST http://localhost:8080/api/showdown/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "analysisType": "rawLog",
    "rawLog": "|j|☆Player1\n|j|☆Player2\n|switch|p1a: Pikachu|pikachu|50|M|100/100\n|switch|p2a: Charizard|charizard|50|M|100/100\n|turn|1\n|move|p1|Thunderbolt\n|-damage|p2a: Charizard|70/100|[from] move: Thunderbolt [of] p1a: Pikachu\n|turn|2\n|move|p2|Flamethrower\n|-damage|p1a: Pikachu|50/100\n|win|Player1",
    "isPrivate": true
  }'
```

**Response (200 OK):**
```json
{
  "status": "success",
  "battleId": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "format": "gen9vgc2025reghbo3",
    "timestamp": "2025-11-22T10:30:00Z",
    "duration": 60,
    "player1": {
      "name": "Player1",
      "team": [/* 6 Pokémon */],
      "active": {/* Pikachu */},
      "losses": 0,
      "totalLeft": 6
    },
    "player2": {
      "name": "Player2",
      "team": [/* 6 Pokémon */],
      "active": {/* Charizard */},
      "losses": 0,
      "totalLeft": 6
    },
    "winner": "player1",
    "turns": [/* turn data */],
    "stats": {/* battle statistics */},
    "keyMoments": [/* notable moments */]
  },
  "metadata": {
    "parseTimeMs": 12,
    "analysisTimeMs": 45,
    "cached": false
  }
}
```

#### Example 2: Analyze by Showdown Replay ID

```bash
curl -X POST http://localhost:8080/api/showdown/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "analysisType": "replayId",
    "replayId": "gen9vgc2025reghbo3-2481642254",
    "isPrivate": false
  }'
```

**Response:** Same structure as above, but fetched from Showdown servers

#### Example 3: List Replays with Filters

```bash
curl "http://localhost:8080/api/showdown/replays?username=Player1&format=gen9vgc2025reghbo3&limit=5&offset=0"
```

**Response (200 OK):**
```json
{
  "status": "success",
  "battles": [
    {/* BattleSummary 1 */},
    {/* BattleSummary 2 */},
    {/* BattleSummary 3 */}
  ],
  "total": 47,
  "limit": 5,
  "offset": 0
}
```

#### Example 4: Get Specific Replay

```bash
curl "http://localhost:8080/api/showdown/replays/550e8400-e29b-41d4-a716-446655440000"
```

**Response:** Single BattleSummary object (same as analyze response)

#### Error Response Example

```bash
curl -X POST http://localhost:8080/api/showdown/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "analysisType": "invalidType"
  }'
```

**Response (400 Bad Request):**
```json
{
  "status": "error",
  "error": "Invalid request: must provide either replayId, username, or rawLog",
  "code": "INVALID_REQUEST",
  "details": "Missing required field for analysisType: invalidType"
}
```

**Common Error Codes:**
- `INVALID_REQUEST` - Missing/invalid parameters
- `NOT_FOUND` - Replay/user not found
- `PARSE_ERROR` - Failed to parse battle log
- `NOT_IMPLEMENTED` - Feature not yet available (TCG Live)
- `INTERNAL_ERROR` - Server error

## Showdown Replay Data Flow

### Input → Processing → Output

```
┌─────────────────────────────────────────────────────────────────┐
│                   Input Sources                                 │
├─────────────────┬──────────────────┬────────────────────────────┤
│ Replay ID       │ Username         │ Raw Log                    │
│ (Fetch from     │ (Query Showdown  │ (User-provided pipe-      │
│  Showdown)      │  servers)        │  delimited format)         │
└────────┬────────┴────────┬─────────┴────────────┬───────────────┘
         │                 │                      │
         └─────────────────┴──────────────────────┘
                           │
                           ▼
         ┌─────────────────────────────────┐
         │   Parser (internal/analysis)    │
         │  ParseShowdownLog()              │
         │                                 │
         │ • Extract player names          │
         │ • Parse format/rule set         │
         │ • Process turn-by-turn actions  │
         │ • Identify switches & moves     │
         │ • Calculate statistics          │
         └────────────┬────────────────────┘
                      │
                      ▼
         ┌─────────────────────────────────┐
         │   BattleSummary (types.go)      │
         │                                 │
         │ • Metadata (ID, format, time)   │
         │ • Players & teams               │
         │ • Turns & actions               │
         │ • Statistics & key moments      │
         └────────────┬────────────────────┘
                      │
         ┌────────────┴──────────────┐
         │                           │
         ▼                           ▼
    ┌──────────────┐      ┌──────────────────┐
    │ API Response │      │ Database Storage │
    │ (JSON)       │      │ (PostgreSQL)     │
    └──────────────┘      └──────────────────┘
```

### Showdown Log Format

Showdown replays use a **pipe-delimited format** where each line represents an event:

```
|j|☆Player1          # Join Player 1
|j|☆Player2          # Join Player 2
|switch|p1a: Pikachu|pikachu|50|M|25/100  # Switch to Pokémon with HP/maxHP
|move|p1|Thunderbolt # Player 1 uses Thunderbolt
|damage|p2a: Charizard|30/100|[from] move: Thunderbolt
|turn|1              # Turn number marker
|-damage|p2a: Charizard|20/100|[from] move: Thunderbolt [of] p1a: Pikachu
|faint|p2a: Charizard # Pokémon faints
|switch|p2a: Dragonite|dragonite|50|M|45/100 # Opponent switches
|turn|2
```

**Key Line Types:**
- `|j|`: Player join (player name)
- `|switch|`: Pokémon switch (species, ID, level, gender, current/max HP)
- `|move|`: Move used
- `|damage|`: Damage dealt
- `|-damage|`: Damage in turn action (shows who dealt it)
- `|faint|`: Pokémon faints/knocked out
- `|turn|`: Turn counter
- `|win|`: Battle winner

### Parser Processing

The parser (`internal/analysis/parser.go`) performs these steps:

1. **Split** log into lines
2. **Extract players** from join events
3. **Identify format** from metadata
4. **Process each turn**:
   - Record active Pokémon
   - Parse moves and switches
   - Calculate damage
   - Track stat changes
5. **Calculate statistics**:
   - Move frequency (which moves used most)
   - Type coverage (which types appear)
   - Effectiveness tracking (super effective, not very effective)
   - Player performance metrics
6. **Identify key moments**:
   - Knockouts (faints)
   - Critical moments
   - Pivotal switches
7. **Generate output** as `BattleSummary`

### BattleSummary Structure

The analysis output contains:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "format": "gen9vgc2025reghbo3",
  "timestamp": "2025-11-22T10:30:00Z",
  "duration": 720,
  "player1": {
    "name": "Player1",
    "team": [/* 6 Pokémon */],
    "active": {/* current Pokémon */},
    "losses": 2,
    "totalLeft": 4
  },
  "player2": {
    "name": "Player2",
    "team": [/* 6 Pokémon */],
    "active": {/* current Pokémon */},
    "losses": 1,
    "totalLeft": 5
  },
  "winner": "player1",
  "turns": [
    {
      "turnNumber": 1,
      "actions": [
        {
          "player": "player1",
          "actionType": "move",
          "move": {
            "id": "thunderbolt",
            "name": "Thunderbolt",
            "type": "Electric",
            "power": 90,
            "accuracy": 100
          }
        }
      ],
      "stateAfter": {
        "player1Active": {/* Pikachu */},
        "player2Active": {/* Charizard */},
        "player1Team": ["Pikachu", "Dragonite", ...],
        "player2Team": ["Charizard", "Hydreigon", ...]
      },
      "damageDealt": {"Player1": 70},
      "healingDone": {"Player2": 25}
    }
  ],
  "stats": {
    "totalTurns": 12,
    "moveFrequency": {"thunderbolt": 4, "earthquake": 3, ...},
    "typeCoverage": {"Electric": 4, "Ground": 3, ...},
    "switches": 8,
    "criticalHits": 1,
    "superEffective": 3,
    "notVeryEffective": 1,
    "avgDamagePerTurn": 58.3,
    "avgHealPerTurn": 12.5,
    "player1Stats": {
      "moveCount": 12,
      "switchCount": 4,
      "damageDealt": 700,
      "damageTaken": 450,
      "effectiveness": {"superEffective": 3, "notVeryEffective": 0, "neutral": 9}
    },
    "player2Stats": {/* similar */}
  },
  "keyMoments": [
    {
      "turnNumber": 3,
      "description": "Charizard fainted",
      "type": "kO",
      "significance": 8
    },
    {
      "turnNumber": 5,
      "description": "Player1 switched to Dragonite",
      "type": "switch",
      "significance": 6
    }
  ]
}
```

### Database Integration

When `StoreBattle()` is called in `internal/db/db.go`:

1. **Insert battle record** into `battles` table
   - ID, format, players, winner, duration
2. **Insert analysis data** into `battle_analysis` table
   - Move frequencies, statistics, type coverage
3. **Insert key moments** into `key_moments` table
   - Each moment with turn number and significance
4. **Insert turn data** into `battle_turns` table
   - State snapshots at each turn
5. **Insert actions** into `battle_actions` table
   - Individual moves and switches

This enables:
- Fast replay of specific battles
- Statistical queries across battles
- Filtering/searching by player, format, statistics
- Trend analysis over time

## Database Schema

The schema includes:

- **battles**: Main battle records with players and metadata
- **battle_analysis**: Computed statistics per battle
- **key_moments**: Notable events in battles
- **battle_turns**: Turn-by-turn state snapshots
- **battle_actions**: Individual actions (moves, switches)
- **pokemon**, **pokemon_species**: Pokémon reference data
- **moves**, **items**: Move and item reference data
- **pokemon_moves**: Pokémon move availability mappings

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
