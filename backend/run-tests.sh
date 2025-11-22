#!/bin/bash
# Quick test runner script for BattleForge backend
# Usage: ./run-tests.sh [command]
# Commands: test, test-v, coverage, coverage-html, clean

set -e

cd "$(dirname "$0")" || exit 1

echo "🧪 BattleForge Backend Test Runner"
echo "=================================="
echo ""

command="${1:-test}"

case "$command" in
  test)
    echo "Running all tests..."
    go test ./...
    echo "✅ Tests passed!"
    ;;
  test-v)
    echo "Running tests with verbose output..."
    go test ./... -v
    echo "✅ Tests passed!"
    ;;
  test-short)
    echo "Running tests (fast mode)..."
    go test ./... -short
    echo "✅ Tests passed!"
    ;;
  coverage)
    echo "Running tests with coverage report..."
    go test ./... -cover
    echo ""
    go test ./... -coverprofile=coverage.out
    echo "📊 Coverage report:"
    go tool cover -func=coverage.out | tail -1
    ;;
  coverage-html)
    echo "Generating HTML coverage report..."
    go test ./... -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html
    echo "📊 Coverage report generated: coverage.html"
    if command -v open &> /dev/null; then
      echo "Opening in default browser..."
      open coverage.html
    elif command -v xdg-open &> /dev/null; then
      echo "Opening in default browser..."
      xdg-open coverage.html
    else
      echo "Open coverage.html in your browser to view the report"
    fi
    ;;
  clean)
    echo "Cleaning test artifacts..."
    rm -f coverage.out coverage.html
    go clean -testcache
    echo "✅ Cleaned!"
    ;;
  help|--help|-h)
    echo "Usage: ./run-tests.sh [command]"
    echo ""
    echo "Commands:"
    echo "  test           Run all tests"
    echo "  test-v         Run tests with verbose output"
    echo "  test-short     Run tests in fast mode"
    echo "  coverage       Show coverage percentage"
    echo "  coverage-html  Generate and open HTML coverage report"
    echo "  clean          Remove test artifacts"
    echo "  help           Show this help message"
    echo ""
    echo "Examples:"
    echo "  ./run-tests.sh              # Run tests (default)"
    echo "  ./run-tests.sh test-v       # Run with verbose output"
    echo "  ./run-tests.sh coverage     # Show coverage"
    echo "  ./run-tests.sh coverage-html # Open coverage in browser"
    ;;
  *)
    echo "❌ Unknown command: $command"
    echo "Run './run-tests.sh help' for usage"
    exit 1
    ;;
esac

echo ""
echo "Done! 🎉"
