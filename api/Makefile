.PHONY: generate format lint test test-coverage test-integration test-all run build clean docker-up docker-down docker-build-prod deps check

# Native commands (tools must be on PATH; DevContainer provides them automatically)
# These commands work when tools are installed locally or via DevContainer

generate: deps
	oapi-codegen -generate types,chi-server -package openapi \
		-o pkg/openapi/server.gen.go openapi/openapi.yaml
	# Note: go mod tidy is run in deps target, but if generated code adds new dependencies,
	# you may need to run 'go mod tidy' manually after generate

format:
	@echo "Checking code formatting..."
	@FORMAT_OUTPUT=$$(gofmt -l .); \
	if [ -n "$$FORMAT_OUTPUT" ]; then \
		echo "Error: Code is not formatted. Run 'go fmt ./...' to fix."; \
		echo "$$FORMAT_OUTPUT"; \
		exit 1; \
	fi

lint: deps
	@echo "Building packages to ensure all dependencies are available..."
	@go build ./...
	golangci-lint run ./...

test: deps
	go test -v -race ./...

# Integration tests (requires MinIO via docker-compose)
# Run integration tests using build tags
# Usage: make test-integration
# Prerequisite: docker compose up -d minio
test-integration: deps
	go test -v -tags=integration ./...

# Run all tests (unit + integration)
# This runs unit tests first, then integration tests
# Useful for comprehensive local testing before pushing
test-all: deps
	go test -v -race ./...
	go test -v -tags=integration ./...

test-coverage: deps
	@set +e; \
	go test -v -race -coverprofile=coverage.out ./... 2>&1 | grep -v "no such tool \"covdata\"" || true; \
	test_exit=$$?; \
	if [ -f coverage.out ]; then \
		go tool cover -func=coverage.out || true; \
	fi; \
	if [ $$test_exit -ne 0 ]; then \
		exit $$test_exit; \
	fi

deps:
	go mod download
	go mod tidy
	# Build covdata tool for Go 1.25+ (not included as prebuilt binary)
	# This ensures the tool is available in both local and CI environments
	@go tool covdata 2>&1 | head -1 > /dev/null || go install -a cmd/covdata || true

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

clean:
	rm -rf bin/
	rm -f coverage.out

# Pre-commit check: Run format, lint, and tests
# This is used by the pre-commit hook and can be run manually
check: format lint test
	@echo "✅ All checks passed!"

# Docker Compose management (for smoke-testing, not iterative development)
# Use 'docker compose' (space) instead of 'docker-compose' (hyphen) for compatibility
DOCKER_COMPOSE = docker compose -f ../docker-compose.yml

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

# Production build (builds outside container, requires local Docker)
# Note: Run from api/ directory, so build context is current directory (.)
docker-build-prod:
	docker build -t optel-training:latest -f Dockerfile .
