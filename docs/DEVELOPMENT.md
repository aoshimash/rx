# Development Guide

## Prerequisites

**Local Requirements (minimal):**
- Docker & Docker Compose
- Make
- curl (for testing API endpoints)

**Note:** All Go development tools (Go compiler, oapi-codegen, golangci-lint) are bundled in the Docker container. You don't need to install them locally.

## Project Setup

### Docker-Based Development Environment

All development tools are containerized. You only need Docker, Make, and curl installed locally.

```bash
# Clone the repository
git clone https://github.com/aoshimash/optel-workload.git
cd optel-workload

# Navigate to API directory
cd api

# Start Docker development container
make docker-up

# Generate OpenAPI code (runs inside container)
make generate

# Run linter (runs inside container)
make lint

# Run tests (runs inside container)
make test

# Start development server (runs inside container)
make run

# In another terminal, test the API
curl http://localhost:8080/api/v1/workloads
```

### Docker Commands

```bash
# Start development container
make docker-up

# Stop development container
make docker-down

# Open shell in container
make docker-shell

# All make commands (generate, lint, test, run, build) run inside the container
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

In Schema-First development, **the OpenAPI specification is the source of truth**. Data structures, API contracts, and domain models all derive from the spec, not the other way around.

**Workflow:**

1. **Edit OpenAPI spec** - Start by modifying `openapi/openapi.yaml`
   - Define minimal schemas first (only required fields)
   - Add fields incrementally as needed
   - Data structures emerge from the spec, not predefined
2. **Generate code** - Run `make generate`
   - This generates Go types and server stubs from the spec
3. **Implement handlers** - Write handler logic
   - Use generated types from `pkg/openapi/`
4. **Add tests** - Write table-driven tests
5. **Run linter** - Run `make lint`
6. **Test** - Run `make test`

## Incremental Development Guide

This guide shows how to build the API incrementally, starting with a minimal working specification and gradually adding features. Each step should result in a working, testable system.

### Important: Data Structures Are Defined During Development

**Key Principle:** In Schema-First development, data structures are **not** predefined. They emerge from the OpenAPI specification as you build it incrementally.

- **Start minimal** - Begin with the smallest possible schema that works
- **Expand gradually** - Add fields and relationships as needed
- **OpenAPI is the source of truth** - The spec defines the data structure, not the other way around
- **Domain models follow the spec** - Go domain entities are derived from OpenAPI schemas

**Note:** The domain reference (`.claude/skills/optel-domain/reference.md`) shows the **target state** or **reference implementation**, not a requirement to implement everything upfront. Use it as inspiration, but feel free to start simpler and build up.

### Development Cycle

For each feature addition, follow this cycle:

1. **Define minimal spec** - Add only what's needed for the current step
2. **Generate code** - Run `make generate`
3. **Implement minimal handler** - Get it working with basic logic
4. **Test manually** - Use `curl` or similar to verify it works
5. **Add tests** - Write table-driven tests
6. **Commit** - Save working state before next step

### Step-by-Step: Building Workload API

#### Step 1: Minimal OpenAPI Structure

Start with the absolute minimum - just the OpenAPI info and a single endpoint:

```yaml
openapi: 3.1.0
info:
  title: OPTel Workload API
  version: 0.1.0
servers:
  - url: http://localhost:8080/api/v1
paths:
  /workloads:
    get:
      summary: List workloads
      operationId: listWorkloads
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id:
                      type: string
                    timestamp:
                      type: string
                      format: date-time
```

**Action:**
```bash
# Start Docker development container (if not already running)
cd api
make docker-up

# Generate code (runs inside container)
make generate

# Create minimal handler that returns empty array
# Start server (runs inside container, in background or another terminal)
make run

# Test it works (from local machine)
curl http://localhost:8080/api/v1/workloads
# Should return: []

# Commit
git commit -m "feat(api): add minimal GET /workloads endpoint"
```

#### Step 2: Add POST Endpoint with Minimal Schema

Add the ability to create workloads with only required fields:

```yaml
paths:
  /workloads:
    get:
      # ... existing ...
    post:
      summary: Create a workload
      operationId: createWorkload
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/WorkloadCreate'
      responses:
        '201':
          description: Created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Workload'
components:
  schemas:
    WorkloadCreate:
      type: object
      required:
        - timestamp
        - duration_seconds
        - intensity_rpe
        - volume_kg
      properties:
        timestamp:
          type: string
          format: date-time
        duration_seconds:
          type: integer
          minimum: 1
        intensity_rpe:
          type: integer
          minimum: 1
          maximum: 10
        volume_kg:
          type: number
          minimum: 0
    Workload:
      allOf:
        - $ref: '#/components/schemas/WorkloadCreate'
        - type: object
          properties:
            id:
              type: string
            created_at:
              type: string
              format: date-time
```

**Action:**
```bash
make generate
# Implement minimal handler with in-memory store
# Test it works
curl -X POST http://localhost:8080/api/v1/workloads \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2026-01-24T10:00:00Z",
    "duration_seconds": 3600,
    "intensity_rpe": 7,
    "volume_kg": 5000
  }'

# Verify GET returns the created workload
curl http://localhost:8080/api/v1/workloads

# Commit
git commit -m "feat(api): add POST /workloads with minimal schema"
```

#### Step 3: Add Optional Fields

Extend the schema with optional fields one at a time:

```yaml
    WorkloadCreate:
      type: object
      required:
        - timestamp
        - duration_seconds
        - intensity_rpe
        - volume_kg
      properties:
        # ... existing required fields ...
        subsystems:
          type: array
          items:
            type: string
        notes:
          type: string
    Workload:
      allOf:
        - $ref: '#/components/schemas/WorkloadCreate'
        - type: object
          properties:
            id:
              type: string
            created_at:
              type: string
              format: date-time
            # Add optional fields from WorkloadCreate
            subsystems:
              type: array
              items:
                type: string
            notes:
              type: string
```

**Action:**
```bash
make generate
# Update handler to handle optional fields
# Test with and without optional fields
curl -X POST http://localhost:8080/api/v1/workloads \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2026-01-24T10:00:00Z",
    "duration_seconds": 3600,
    "intensity_rpe": 7,
    "volume_kg": 5000,
    "subsystems": ["chest", "triceps"],
    "notes": "Bench press session"
  }'

# Commit
git commit -m "feat(api): add optional subsystems and notes fields"
```

#### Step 4: Add GET by ID Endpoint

Add the ability to retrieve a specific workload:

```yaml
paths:
  /workloads:
    # ... existing ...
  /workloads/{id}:
    get:
      summary: Get a workload by ID
      operationId: getWorkload
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Workload'
        '404':
          description: Not found
```

**Action:**
```bash
make generate
# Implement handler
# Test it works
WORKLOAD_ID=$(curl -s -X POST http://localhost:8080/api/v1/workloads \
  -H "Content-Type: application/json" \
  -d '{"timestamp":"2026-01-24T10:00:00Z","duration_seconds":3600,"intensity_rpe":7,"volume_kg":5000}' \
  | jq -r '.id')

curl http://localhost:8080/api/v1/workloads/$WORKLOAD_ID

# Commit
git commit -m "feat(api): add GET /workloads/{id} endpoint"
```

#### Step 5: Add Query Parameters for Filtering

Add filtering capabilities to the list endpoint:

```yaml
  /workloads:
    get:
      summary: List workloads
      operationId: listWorkloads
      parameters:
        - name: from
          in: query
          schema:
            type: string
            format: date-time
        - name: to
          in: query
          schema:
            type: string
            format: date-time
        - name: subsystem
          in: query
          schema:
            type: string
      responses:
        # ... existing ...
```

**Action:**
```bash
make generate
# Implement filtering logic in handler
# Test filters
curl "http://localhost:8080/api/v1/workloads?from=2026-01-01T00:00:00Z&to=2026-01-31T23:59:59Z"
curl "http://localhost:8080/api/v1/workloads?subsystem=chest"

# Commit
git commit -m "feat(api): add query parameters for filtering workloads"
```

#### Step 6: Add DELETE Endpoint

Add soft-delete capability:

```yaml
  /workloads/{id}:
    get:
      # ... existing ...
    delete:
      summary: Delete a workload
      operationId: deleteWorkload
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '204':
          description: No content
        '404':
          description: Not found
```

**Action:**
```bash
make generate
# Implement soft-delete (mark as deleted, don't remove)
# Test it works
curl -X DELETE http://localhost:8080/api/v1/workloads/$WORKLOAD_ID
curl http://localhost:8080/api/v1/workloads/$WORKLOAD_ID
# Should return 404

# Commit
git commit -m "feat(api): add DELETE /workloads/{id} endpoint"
```

### Validation and Testing

After each step, verify:

1. **OpenAPI spec is valid:**
   ```bash
   # Use swagger-cli or similar
   npx @apidevtools/swagger-cli validate openapi/openapi.yaml
   ```

2. **Generated code compiles:**
   ```bash
   make generate
   go build ./...
   ```

3. **Handlers work:**
   ```bash
   make run
   # Test with curl in another terminal
   ```

4. **Tests pass:**
   ```bash
   make test
   ```

5. **Linter passes:**
   ```bash
   make lint
   ```

### Best Practices

1. **One feature per commit** - Each step should be a single, focused commit
2. **Test before moving on** - Don't proceed to the next step until current one works
3. **Keep specs minimal** - Only add what's needed for current functionality
4. **Validate frequently** - Run `make generate` and `make test` after each spec change
5. **Document decisions** - Add comments in spec for non-obvious choices

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

### Development vs Production Containers

**Important:** Development and production containers are **completely separate**.

| Aspect | Development (`Dockerfile`) | Production (`Dockerfile.prod`) |
|--------|---------------------------|-------------------------------|
| Purpose | Local development | Production deployment |
| Base Image | `golang:1.25` (Ubuntu-based) | `gcr.io/distroless/static-debian11:nonroot` |
| Size | Larger (includes Go compiler, tools) | Minimal (binary only, no shell, no package manager) |
| Tools | oapi-codegen, golangci-lint, make | None (distroless has no shell) |
| User | root (for development) | nonroot (65532:65532, provided by distroless) |
| Security | Standard (development tools) | High (minimal attack surface, no shell) |
| Volumes | Code mounted for live editing | No volumes |
| Build | Single stage | Multi-stage build (Ubuntu builder + distroless runtime) |
| Use Case | `make docker-up`, `make generate` | CI/CD, Kubernetes, production |

### Development Environment

The development environment is fully containerized. All Go tools (Go compiler, oapi-codegen, golangci-lint) run inside a Docker container.

**Benefits:**
- No local Go installation required
- Consistent development environment across machines
- Isolated from system dependencies
- Easy to reset or recreate

**Workflow:**
1. Start container: `make docker-up`
2. Run commands: `make generate`, `make lint`, `make test`, `make run`
3. All commands execute inside the container
4. Code changes are reflected immediately (volume mount)

**Files:**
- `api/Dockerfile` - Development container (includes all tools)
- `docker-compose.yml` - Development container orchestration

### Container Management

```bash
# Start development container (runs in background)
make docker-up

# Stop development container
make docker-down

# Open interactive shell in container
make docker-shell

# View container logs
docker-compose -f ../docker-compose.yml logs -f dev
```

### Local Development with PostgreSQL (Phase 2+)

When PostgreSQL is needed, it will be added to `docker-compose.yml`:

```bash
# Start PostgreSQL and development container
docker-compose -f ../docker-compose.yml up -d postgres dev
```

### Building Production Image

**Production images are built separately from development containers.**

```bash
# Build production image (from api/ directory)
docker build -t optel-workload:latest -f Dockerfile.prod .

# Run production container
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e LOG_LEVEL=info \
  optel-workload:latest

# Test production container
curl http://localhost:8080/api/v1/workloads
```

**Production Image Features:**
- Multi-stage build (Ubuntu-based builder + distroless runtime)
- Distroless base image (minimal attack surface, no shell, no package manager)
- Non-root user (nonroot:nonroot, 65532:65532) for security
- Static binary (CGO_ENABLED=0) for compatibility with distroless
- Minimal image size (only binary and CA certificates)
- No shell access (enhanced security)
- Health check support (when endpoint is implemented)

**Files:**
- `api/Dockerfile.prod` - Production container (optimized, minimal)

## Troubleshooting

### Common Issues

**Container not running:**
```bash
# Check if container is running
docker-compose -f ../docker-compose.yml ps

# Start container if not running
make docker-up

# Rebuild container if needed
docker-compose -f ../docker-compose.yml build --no-cache dev
```

**Permission errors:**
- Ensure Docker has proper permissions
- On Linux, you may need to add your user to the docker group

**oapi-codegen or golangci-lint not found:**
- These tools are installed in the Docker container
- If you see "command not found", ensure the container is running: `make docker-up`
- Rebuild the container if tools are missing: `docker-compose -f ../docker-compose.yml build dev`

**golangci-lint errors:**
- Review `.golangci.yml` for enabled linters
- Fix issues or add exclusions if false positives
- Run linter inside container: `make lint`

**Generated code conflicts:**
- Always run `make generate` after changing OpenAPI spec
- Commit generated code with spec changes
- If generation fails, check OpenAPI spec syntax

**Port 8080 already in use:**
- Change port in `docker-compose.yml` or stop the conflicting service
- Update `ports` mapping: `"8080:8080"` → `"8081:8080"`
