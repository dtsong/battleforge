.PHONY: test test-v test-coverage test-coverage-html test-short clean help

help:
	@echo "BattleForge Backend - Available targets:"
	@echo "  make test              - Run all tests"
	@echo "  make test-v            - Run all tests with verbose output"
	@echo "  make test-short        - Run tests without network/caching (fast)"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make test-coverage-html - Generate and open HTML coverage report"
	@echo "  make test-httpapi      - Run only httpapi tests"
	@echo "  make test-analysis     - Run only analysis tests"
	@echo "  make build             - Build the API binary"
	@echo "  make run               - Run the API locally"
	@echo "  make fmt               - Format code"
	@echo "  make lint              - Run linter"
	@echo "  make clean             - Clean build artifacts and test files"

test:
	go test ./...

test-v:
	go test ./... -v

test-short:
	go test ./... -short

test-coverage:
	go test ./... -cover

test-coverage-report:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

test-coverage-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

test-httpapi:
	go test -v ./internal/httpapi

test-analysis:
	go test -v ./internal/analysis

test-observability:
	go test -v ./internal/observability

build:
	go build -o battleforge-api ./cmd/battleforge-api

run:
	go run ./cmd/battleforge-api

fmt:
	go fmt ./...

vet:
	go vet ./...

lint: fmt vet
	golangci-lint run ./...

clean:
	rm -f battleforge-api coverage.out coverage.html
	go clean -testcache

.DEFAULT_GOAL := help
