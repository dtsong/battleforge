# BattleForge Backend Tests

Comprehensive unit tests for the BattleForge backend API with coverage for all major components.

## Test Structure

### Router Tests (`internal/httpapi/router_test.go`)

Tests the HTTP router setup and route registration:

- **TestHealthCheck**: Verifies `/healthz` endpoint returns `ok` with status 200
- **TestRouterEndpointsExist**: Confirms all OpenAPI-defined routes are registered

### Handler Tests (`internal/httpapi/showdown_handlers_test.go`)

Tests all Showdown analysis API endpoints:

#### Analyze Endpoint Tests
- **TestAnalyzeShowdownRawLog**: Tests raw battle log analysis
  - Valid log returns 200 OK with BattleSummary
  - Empty log returns 400 Bad Request
  - Invalid analysis type returns 400

- **TestAnalyzeShowdownByReplayID**: Tests replay ID analysis
  - Empty ID returns 400
  - Valid ID returns 501 Not Implemented (planned feature)

- **TestAnalyzeShowdownByUsername**: Tests username-based analysis
  - Missing username returns 400
  - Missing format returns 400
  - Valid inputs return 501 Not Implemented (planned)

- **TestAnalyzeShowdownInvalidJSON**: Tests request parsing
  - Malformed JSON returns 400 with INVALID_REQUEST code

#### Get Replay Tests
- **TestGetShowdownReplay**: Tests replay retrieval endpoint
  - Returns 404 with NOT_FOUND code (DB integration pending)

#### List Replays Tests
- **TestListShowdownReplays**: Tests replay listing with filtering
  - No filters returns 200 with empty data
  - Username filter works
  - Format filter works
  - Limit and offset pagination handled correctly
  - Invalid limits default to safe values

#### TCG Live Tests
- **TestAnalyzeTCGLive**: Tests TCG Live analysis endpoint
  - Valid game export returns 501 Not Implemented
  - Empty game export returns 400

### Parser Tests (`internal/analysis/parser_test.go`)

Comprehensive tests for the Showdown battle log parser:

#### Basic Parsing
- **TestParseShowdownLogBasicValid**: Validates basic log parsing
  - Returns non-nil summary
  - Sets battle ID
  - Sets format
  - Parses player names
  - Extracts turns

#### Data Extraction
- **TestParseShowdownLogPlayerNames**: Verifies player names are correctly extracted
- **TestParseShowdownLogFormat**: Checks format string parsing
- **TestParseShowdownLogTurns**: Validates turn number sequencing
- **TestParseShowdownLogActions**: Ensures actions are parsed with correct player/type
- **TestParseShowdownLogWinner**: Verifies winner determination
- **TestParseShowdownLogStats**: Checks statistical calculations

#### Edge Cases
- **TestParseShowdownLogMinimalLog**: Handles minimal valid logs
- **TestParseShowdownLogEmptyLog**: Gracefully handles empty input
- **TestParseShowdownLogMalformedLog**: Resilient to malformed data

#### Action Type Parsing
- **TestParseShowdownLogMoveParsing**: Verifies move action parsing
  - Move ID and name are set
  - Appears in move frequency stats

- **TestParseShowdownLogSwitchParsing**: Validates switch action parsing
  - Switch target is correctly identified

- **TestParseShowdownLogKeyMoments**: Tests key moment detection
  - KO events are recorded
  - Proper significance levels

#### State Tracking
- **TestParseShowdownLogPlayerLosses**: Tracks fainting Pokémon
- **TestParseShowdownLogUUIDUniqueness**: Verifies unique battle IDs

## Test Fixtures (`internal/httpapi/test_fixtures.go`)

Helper functions providing sample data:

- `sampleShowdownLog()`: Valid complete battle log
- `minimalShowdownLog()`: Minimal valid log structure
- `malformedLog()`: Invalid log for error testing
- `emptyLog()`: Empty input
- `logWithoutPlayers()`: Missing player data
- `logWithoutTurns()`: Battle log with no turns

## Running Tests

### Run all tests
```bash
cd backend
go test ./...
```

### Run with verbose output
```bash
go test ./... -v
```

### Run specific test package
```bash
go test -v ./internal/httpapi
go test -v ./internal/analysis
```

### Run specific test
```bash
go test -v -run TestParseShowdownLogBasicValid ./internal/analysis
```

### Run with coverage report
```bash
go test ./... -cover
```

### Generate detailed coverage report
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out  # Open in browser
```

### Coverage by package
```bash
go test -cover ./internal/httpapi
go test -cover ./internal/analysis
go test -cover ./internal/observability
```

## Test Coverage Goals

### httpapi package (Router & Handlers)
- ✅ Health check endpoint
- ✅ Analyze endpoint (all input types)
- ✅ List replays endpoint
- ✅ Get specific replay endpoint
- ✅ Error responses with codes
- ✅ Request validation
- ✅ JSON marshaling/unmarshaling

### analysis package (Parser)
- ✅ Valid log parsing
- ✅ Player name extraction
- ✅ Format parsing
- ✅ Turn sequencing
- ✅ Action parsing (moves, switches, faints)
- ✅ Key moment detection
- ✅ Statistics calculation
- ✅ Edge cases (empty, minimal, malformed)
- ✅ UUID generation uniqueness

### observability package
- ✅ Logger initialization
- ✅ Log level output

## Planned Test Coverage (Future)

Once database integration is complete:

- Database storage operations
- Battle retrieval from database
- Filtering and pagination
- Transaction handling
- Connection pooling
- Error recovery

Once Showdown/TCG Live API integration is complete:

- Fetching replays from Showdown API
- Parsing exported game files
- Caching strategies
- API rate limiting handling

## Best Practices Used

1. **Table-driven tests**: Each handler test uses a slice of test cases for comprehensive coverage
2. **Test isolation**: Each test sets up its own fixtures and doesn't depend on others
3. **Meaningful assertions**: Error messages clearly indicate what was expected vs actual
4. **Edge case coverage**: Tests include empty inputs, malformed data, and boundary conditions
5. **Descriptive names**: Test function names clearly describe what they're testing
6. **Test helpers**: Reusable fixture functions reduce duplication
7. **No side effects**: Tests don't modify global state

## Adding New Tests

When adding new features:

1. Write tests first (TDD approach)
2. Place tests in `*_test.go` file in same package
3. Use table-driven pattern for multiple scenarios
4. Add fixtures to `test_fixtures.go` if reusable
5. Include both happy path and error cases
6. Run `go test -cover` to verify coverage

Example test structure:
```go
func TestNewFeature(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		expectedError bool
	}{
		{"valid input", "foo", "bar", false},
		{"empty input", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewFeature(tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
```

## Test Execution Notes

- Tests run in parallel by default with `go test ./...`
- Each package's tests are independent
- Tests take <1 second to run
- No external dependencies (all mocked/fixture-based)
- No database required for tests
- No network calls needed
