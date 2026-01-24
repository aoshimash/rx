# Development Guide

## Prerequisites

- Go 1.25+
- Docker & Docker Compose
- oapi-codegen (`go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest`)
- golangci-lint (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)

## Project Setup

```bash
# Clone the repository
git clone https://github.com/aoshimash/optel-workload.git
cd optel-workload

# Navigate to API directory
cd api

# Install dependencies
go mod download

# Generate OpenAPI code
make generate

# Run linter
make lint

# Run tests
make test

# Start development server
make run
```

## Development Phases

### Phase 1: Core API (Current)

**Goals:**
1. Define OpenAPI specification for Workload resource
2. Scaffold Go module with project structure
3. Implement in-memory store for rapid prototyping

**Deliverables:**
- `openapi/openapi.yaml` - API specification
- `internal/domain/workload.go` - Domain entity
- `internal/store/memory/` - In-memory implementation
- `internal/handler/` - HTTP handlers

### Phase 2: Persistence

**Goals:**
1. Add PostgreSQL support
2. Implement database migrations
3. Add Docker Compose for local development

### Phase 3: Program Entity

**Goals:**
1. Define OpenAPI specification for Program resource
2. Implement recursive tree structure
3. Link Workloads to Programs

### Phase 4: MCP Server

**Goals:**
1. Implement MCP Server in `mcp/` directory
2. Expose tools for AI agents to query and record workloads
3. Distribute as Docker image or Python script

**Key Design:**
- MCP Server runs on user's local machine (not on the backend server)
- Communicates with OPTel API via HTTP
- Backend remains a pure REST API with no MCP-specific logic

## Workflow

### Schema-First Development

1. **Edit OpenAPI spec** - Start by modifying `openapi/openapi.yaml`
2. **Generate code** - Run `make generate`
3. **Implement handlers** - Write handler logic
4. **Add tests** - Write table-driven tests
5. **Run linter** - Run `make lint`
6. **Test** - Run `make test`

### Making Changes

```bash
# 1. Create a feature branch
git checkout -b feature/add-workload-filtering

# 2. Make changes following Schema-First approach

# 3. Generate code
make generate

# 4. Run checks
make lint
make test

# 5. Commit with conventional commit message
git commit -m "feat(api): add workload filtering by date range"
```

## Makefile Targets

```makefile
# api/Makefile

.PHONY: generate lint test run build clean

generate:
	oapi-codegen -generate types,chi-server -package openapi \
		-o pkg/openapi/server.gen.go openapi/openapi.yaml

lint:
	golangci-lint run ./...

test:
	go test -v -race ./...

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

clean:
	rm -rf bin/
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `DATABASE_URL` | PostgreSQL connection string (Phase 2) | - |

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./internal/domain/...
```

### Writing Tests

All tests must be table-driven. See `.claude/skills/optel-go-standards/reference.md` for examples.

## Docker

### Local Development (Phase 2+)

```bash
# Start PostgreSQL
docker-compose up -d postgres

# Start API with hot reload (using air or similar)
docker-compose up api
```

### Building Production Image

```bash
docker build -t optel-workload:latest -f Dockerfile .
```

## Troubleshooting

### Common Issues

**oapi-codegen not found:**
```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

**golangci-lint errors:**
- Review `.golangci.yml` for enabled linters
- Fix issues or add exclusions if false positives

**Generated code conflicts:**
- Always run `make generate` after changing OpenAPI spec
- Commit generated code with spec changes
