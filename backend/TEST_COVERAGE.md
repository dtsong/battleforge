# Test Coverage Summary

## Overview

Comprehensive unit test suite for the BattleForge backend with **25+ test cases** covering all major components, edge cases, and error scenarios.

## Test Files Created

### 1. **router_test.go** (2 test functions, 7 test cases)
   - ✅ Health check endpoint responds with "ok"
   - ✅ Health check only accepts GET
   - ✅ All OpenAPI routes are registered (5 routes tested)

### 2. **showdown_handlers_test.go** (6 test functions, 20+ test cases)
   - ✅ Raw log analysis with valid/invalid inputs
   - ✅ Replay ID analysis (valid and empty ID)
   - ✅ Username analysis (missing username/format)
   - ✅ Invalid JSON request parsing
   - ✅ Get specific replay endpoint
   - ✅ List replays with filters (username, format, pagination)
   - ✅ TCG Live analysis endpoint
   - ✅ Error response codes (INVALID_REQUEST, NOT_FOUND, NOT_IMPLEMENTED)

### 3. **parser_test.go** (15 test functions, 20+ test cases)
   - ✅ Valid battle log parsing
   - ✅ Player name extraction
   - ✅ Format parsing
   - ✅ Turn sequencing (correct turn numbers)
   - ✅ Action type parsing (moves and switches)
   - ✅ Winner determination
   - ✅ Statistics calculation
   - ✅ Key moment detection (KO events)
   - ✅ Player losses tracking (fainting)
   - ✅ UUID uniqueness
   - ✅ Minimal log handling
   - ✅ Empty log handling
   - ✅ Malformed log resilience
   - ✅ Move parsing details
   - ✅ Switch action parsing

### 4. **test_fixtures.go** (6 helper functions)
   - Sample valid battle logs (full and minimal)
   - Malformed log examples
   - Edge case fixtures (empty, missing data)

## Test Statistics

| Category | Count |
|----------|-------|
| Test Functions | 23 |
| Test Cases | 50+ |
| Code Files with Tests | 2 (handlers, parser) |
| Test Fixtures | 6 |

## Coverage Areas

### HTTP API Layer
- **Endpoint Routing**: All 5 OpenAPI endpoints registered and accessible
- **Request Handling**: JSON parsing, parameter validation, query string processing
- **Error Handling**: Proper HTTP status codes and error response structures
- **Request Types**: All three analysis types (rawLog, username, replayId) tested

### Battle Log Parsing
- **Valid Input**: Complete battle logs parsed correctly
- **Data Extraction**: Player names, format, turns, actions all extracted
- **Edge Cases**: Empty logs, minimal logs, malformed data handled gracefully
- **Action Types**: Moves, switches, faints parsed and counted
- **Statistics**: Move frequency, effectiveness tracking validated

### Error Scenarios
- ✅ Empty input fields
- ✅ Missing required parameters
- ✅ Invalid JSON
- ✅ Malformed battle logs
- ✅ Unimplemented features (future endpoints)
- ✅ Invalid analysis types

## How to Run Tests

### Quick Test
```bash
cd backend
go test ./...
```

### Verbose Output
```bash
make test-v  # or: go test ./... -v
```

### With Coverage Report
```bash
make test-coverage-report  # or: go test ./... -cover
```

### HTML Coverage Report
```bash
make test-coverage-html  # Generates and opens coverage.out in browser
```

### Run Specific Test Package
```bash
make test-httpapi    # Only handler/router tests
make test-analysis   # Only parser tests
```

### Run Individual Test
```bash
go test -run TestParseShowdownLogBasicValid -v ./internal/analysis
```

## Test Design Principles

### 1. **Table-Driven Tests**
Each test uses a slice of test cases to cover multiple scenarios:
```go
tests := []struct {
    name          string
    input         string
    expectedCode  string
    expectedError string
}{
    {"valid input", "...", "success", ""},
    {"empty input", "", "error", "required"},
}
```

### 2. **Isolation**
- No shared state between tests
- Each test is independent and can run in any order
- Tests can run in parallel with `go test -race`

### 3. **Clear Assertions**
- Error messages explain what was expected vs. actual
- Each assertion has a descriptive message

### 4. **Fixture Reuse**
- Sample data in `test_fixtures.go` reduces duplication
- Easy to add new test cases by reusing fixtures

### 5. **Edge Case Coverage**
- Valid "happy path" scenarios
- Boundary conditions (empty, minimal)
- Error conditions (malformed, missing data)
- Type variations (different input types)

## Continuous Integration

### GitHub Actions Workflow
- **File**: `.github/workflows/backend-tests.yml`
- **Triggers**: Push to main/develop, Pull Requests
- **Tests On**: Go 1.22.x and 1.25.x
- **Steps**:
  1. Format checking (gofmt)
  2. Vet checking (go vet)
  3. Run all tests
  4. Generate coverage report
  5. Upload to Codecov
  6. Build binary artifact

## Makefile Targets

```bash
make test              # Run all tests
make test-v            # Run with verbose output
make test-coverage     # Show coverage percentages
make test-coverage-html # Generate HTML report
make test-httpapi      # Test only HTTP API
make test-analysis     # Test only parser
make fmt               # Format code
make lint              # Run formatter + linter
make build             # Build binary
make run               # Run API locally
make clean             # Remove test/build artifacts
```

## Test Examples

### Handler Test (Table-Driven)
```go
func TestAnalyzeShowdownRawLog(t *testing.T) {
	tests := []struct {
		name           string
		request        AnalyzeShowdownRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "valid raw log analysis",
			request: AnalyzeShowdownRequest{
				AnalysisType: "rawLog",
				RawLog:       sampleShowdownLog(),
				IsPrivate:    false,
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "empty raw log returns error",
			request: AnalyzeShowdownRequest{
				AnalysisType: "rawLog",
				RawLog:       "",
				IsPrivate:    false,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "rawLog is required",
		},
	}
	// ... test execution
}
```

### Parser Test
```go
func TestParseShowdownLogBasicValid(t *testing.T) {
	log := sampleBattleLog()
	summary, err := ParseShowdownLog(log)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if summary == nil {
		t.Fatal("expected summary, got nil")
	}

	if summary.ID == "" {
		t.Error("expected battle ID to be set")
	}
}
```

## Future Test Expansion

Once database integration is implemented:
- [ ] Database storage and retrieval tests
- [ ] Transaction rollback scenarios
- [ ] Connection pool exhaustion handling
- [ ] Filtering and pagination edge cases

Once API integrations are implemented:
- [ ] Showdown API client tests
- [ ] Mocked HTTP responses
- [ ] Rate limiting scenarios
- [ ] Retry logic

Once frontend integration is complete:
- [ ] End-to-end integration tests
- [ ] API response validation
- [ ] Data serialization/deserialization

## Notes for Go Testing Newcomers

### Key Concepts
1. **Test files** end with `_test.go` and are in the same package
2. **Test functions** start with `Test` and take `*testing.T` parameter
3. **`t.Error()`** reports test failure but continues running
4. **`t.Fatal()`** reports failure and stops test immediately
5. **`t.Run()`** creates subtests for better organization
6. **Table-driven** tests are idiomatic Go for multiple test cases

### Running Tests
- `go test ./...` - Run all tests in current directory and subdirectories
- `go test -v` - Verbose output (show each test)
- `go test -run TestName` - Run only specific test
- `go test -short` - Run only short tests (skip long-running)
- `go test -race` - Check for race conditions

### Coverage
- `go test -cover` - Show coverage percentage
- `go test -coverprofile=coverage.out` - Generate coverage file
- `go tool cover -html=coverage.out` - View as HTML

For more info: https://golang.org/pkg/testing/
