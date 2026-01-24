.PHONY: generate lint test run build clean docker-up docker-down docker-build-prod deps

# Native commands (tools must be on PATH; DevContainer provides them automatically)
# These commands work when tools are installed locally or via DevContainer

generate: deps
	oapi-codegen -generate types,chi-server -package openapi \
		-o pkg/openapi/server.gen.go openapi/openapi.yaml

lint:
	golangci-lint run ./...

test: deps
	go test -v -race ./...

deps:
	go mod download
	go mod tidy

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

clean:
	rm -rf bin/

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
