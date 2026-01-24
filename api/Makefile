.PHONY: generate lint test run build clean docker-up docker-down docker-shell docker-build-prod

# Docker-based commands (run inside container)
DOCKER_COMPOSE = docker-compose -f ../docker-compose.yml
DOCKER_RUN = $(DOCKER_COMPOSE) exec dev

generate:
	$(DOCKER_RUN) oapi-codegen -generate types,chi-server -package openapi \
		-o pkg/openapi/server.gen.go openapi/openapi.yaml

lint:
	$(DOCKER_RUN) golangci-lint run ./...

test:
	$(DOCKER_RUN) go test -v -race ./...

run:
	$(DOCKER_RUN) go run cmd/server/main.go

build:
	$(DOCKER_RUN) go build -o bin/server cmd/server/main.go

# Docker management
docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-shell:
	$(DOCKER_RUN) /bin/sh

clean:
	$(DOCKER_RUN) rm -rf bin/

# Production build (builds outside container, requires local Docker)
docker-build-prod:
	docker build -t optel-workload:latest -f Dockerfile.prod .
