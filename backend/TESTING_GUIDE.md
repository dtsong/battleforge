# BattleForge Backend - Complete Testing Guide

## What Was Created

A comprehensive unit test suite for the BattleForge backend API with **50+ test cases** covering all major components, edge cases, and error scenarios.

## Files Created

### Test Files (3)
1. **`internal/httpapi/router_test.go`** - Router and endpoint registration tests
2. **`internal/httpapi/showdown_handlers_test.go`** - API handler tests (20+ cases)
3. **`internal/analysis/parser_test.go`** - Showdown log parser tests (20+ cases)
4. **`internal/httpapi/test_fixtures.go`** - Sample battle logs and fixtures

### Documentation Files (4)
1. **`TESTING.md`** - Comprehensive testing guide
2. **`TEST_COVERAGE.md`** - Detailed coverage summary
3. **`Makefile`** - Easy test runner commands
4. **`run-tests.sh`** - Bash script for quick testing

### CI/CD Files (1)
1. **`.github/workflows/backend-tests.yml`** - GitHub Actions CI pipeline

## Quick Start

### Run All Tests
```bash
cd backend
make test                    # Simple test run
make test-v                  # Verbose output
make test-coverage          # With coverage %
make test-coverage-html     # Generate HTML report
```

Or using the script:
```bash
./run-tests.sh              # Run tests
./run-tests.sh coverage-html # Open coverage report
```

## Test Coverage

### Router Tests (2 functions)
- ✅ Health check endpoint works
- ✅ All 5 API routes are registered

### Handler Tests (6 functions, 20+ cases)
- ✅ Analyze endpoint with raw logs
- ✅ Analyze endpoint with replay IDs
- ✅ Analyze endpoint with usernames
- ✅ Invalid JSON parsing
- ✅ Get replay by ID
- ✅ List replays with filtering
- ✅ Pagination (limit, offset)
- ✅ TCG Live endpoint
- ✅ Proper error codes and messages

### Parser Tests (15 functions, 20+ cases)
- ✅ Valid log parsing
- ✅ Player name extraction
- ✅ Format detection
- ✅ Turn number sequencing
- ✅ Move action parsing
- ✅ Switch action parsing
- ✅ Faint/KO tracking
- ✅ Winner determination
- ✅ Statistics calculation
- ✅ Key moment detection
- ✅ UUID uniqueness
- ✅ Edge cases (empty, minimal, malformed)
- ✅ Error resilience

## Test Types

### Happy Path Tests
Tests that verify normal, expected behavior:
- Valid battle log parsing
- Correct endpoint responses
- Proper data extraction

### Error Case Tests
Tests that verify error handling:
- Empty input handling
- Invalid JSON parsing
- Missing required parameters
- Malformed data resilience

### Edge Case Tests
Tests that verify boundary conditions:
- Minimal valid logs
- Empty inputs
- Maximum parameter values
- Type variations

## Running Specific Tests

```bash
# Run all tests
go test ./...

# Run only handler tests
go test -v ./internal/httpapi

# Run only parser tests
go test -v ./internal/analysis

# Run specific test
go test -run TestParseShowdownLogBasicValid -v ./internal/analysis

# Run tests matching pattern
go test -run "TestAnalyze.*RawLog" -v ./internal/httpapi

# Run with coverage
go test -cover ./...

# Generate coverage file
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # Open in browser
go tool cover -func=coverage.out  # Show as text
```

## Test Structure Example

All tests follow Go idioms and best practices:

```go
func TestFeature(t *testing.T) {
	// Table-driven test approach
	tests := []struct {
		name          string
		input         string
		expectedCode  int
		expectedError string
	}{
		{
			name:          "valid input",
			input:         "test-data",
			expectedCode:  200,
			expectedError: "",
		},
		{
			name:          "empty input",
			input:         "",
			expectedCode:  400,
			expectedError: "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test implementation
		})
	}
}
```

## GitHub Actions CI/CD

Automatic testing on every push and pull request:

```yaml
# Triggered on: main/develop branch updates, PRs
# Tests on: Go 1.22.x and 1.25.x
# Steps:
#   1. Format check (gofmt)
#   2. Vet check (go vet)
#   3. Run all tests
#   4. Generate coverage
#   5. Upload to Codecov
#   6. Build binary artifact
```

View test results: GitHub → Actions tab → Backend Tests workflow

## Helper Tools

### Makefile Commands
```bash
make test              # Run tests
make test-v            # Verbose output
make test-short        # Fast mode
make test-coverage     # Show coverage %
make test-coverage-html # Open in browser
make build             # Build binary
make run               # Run locally
make fmt               # Format code
make lint              # Lint code
make clean             # Clean artifacts
```

### Test Script
```bash
./run-tests.sh test           # Run tests
./run-tests.sh test-v         # Verbose
./run-tests.sh coverage       # Coverage
./run-tests.sh coverage-html  # HTML report
./run-tests.sh clean          # Clean artifacts
./run-tests.sh help           # Show help
```

## Test Fixtures

Sample data for testing:
- `sampleBattleLog()` - Full valid battle log
- `minimalBattleLog()` - Minimal valid structure
- `malformedBattleLog()` - Invalid data

Used across all handler and parser tests for consistency.

## Coverage Goals

### Current Coverage
- ✅ Router registration
- ✅ Handler HTTP logic
- ✅ Parser core functionality
- ✅ Error paths
- ✅ Edge cases

### Future Coverage (After DB Integration)
- [ ] Database operations
- [ ] Transaction handling
- [ ] Filtering/pagination logic
- [ ] Connection pooling

### Future Coverage (After API Integration)
- [ ] Showdown API client
- [ ] TCG Live parser
- [ ] Caching layer
- [ ] Rate limiting

## Best Practices Used

### 1. Table-Driven Tests
Multiple test cases in a single test function:
- Easy to add new cases
- Clear comparison of inputs/outputs
- DRY (Don't Repeat Yourself)

### 2. Test Isolation
Each test is independent:
- No shared state
- Can run in any order
- Can run in parallel

### 3. Descriptive Names
Test names clearly describe what they test:
- `TestParseShowdownLogBasicValid` - describes the test
- `TestAnalyzeShowdownByUsername` - clear intent

### 4. Fixture Reuse
Common test data in separate file:
- Reduces code duplication
- Easy to update test data
- One source of truth

### 5. Comprehensive Assertions
Clear error messages:
- Expected vs actual values shown
- Helps debug failures quickly

## For Go Testing Newcomers

### Key Concepts
1. **Test files**: End with `_test.go`, same package as code being tested
2. **Test functions**: Start with `Test`, take `*testing.T` parameter
3. **Table-driven**: Common Go pattern with slice of test cases
4. **Subtests**: Use `t.Run()` for better organization
5. **Helper functions**: Reusable test utilities

### Important Methods
- `t.Error()` - Fail but continue
- `t.Fatal()` - Fail and stop
- `t.Errorf()` - Formatted error
- `t.Fatalf()` - Formatted fatal
- `t.Run()` - Subtest

### Running Tests
```bash
go test ./...              # All tests
go test -v ./...           # Verbose
go test -run TestName ./...# Specific
go test -short ./...       # Skip long tests
go test -race ./...        # Check races
go test -cover ./...       # Show coverage
```

## Next Steps

1. **Review the tests**: Open `internal/httpapi/*_test.go` and `internal/analysis/*_test.go`
2. **Run the tests**: `make test` or `./run-tests.sh`
3. **Check coverage**: `make test-coverage-html`
4. **Understand patterns**: Look at table-driven test examples
5. **Add your own**: Follow the same structure for new features

## Resources

- [Go Testing Package Docs](https://golang.org/pkg/testing/)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Best Go Test Practices](https://golang.org/doc/effective_go#names)
- [Testing in Go](https://blog.golang.org/package-names)

## Troubleshooting

### Tests not found
```bash
# Make sure test files exist
ls -la internal/httpapi/*_test.go
ls -la internal/analysis/*_test.go

# Make sure they're in the right location
# (same directory as code being tested)
```

### Import errors
```bash
# Download dependencies
go mod download

# Make sure you're in backend directory
cd backend
```

### Coverage not generating
```bash
# Make sure you have permission to write files
chmod +x ./run-tests.sh

# Run with verbose to see errors
go test ./... -v -coverprofile=coverage.out
```

## Summary

You now have:
- ✅ 50+ comprehensive test cases
- ✅ Full handler endpoint coverage
- ✅ Complete parser validation
- ✅ Edge case and error handling
- ✅ Easy test runner commands
- ✅ CI/CD automation
- ✅ Detailed documentation
- ✅ Best practice examples

Tests are ready to run and can be extended as new features are added!
