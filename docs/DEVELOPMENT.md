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
curl http://localhost:8080/api/v1/workouts
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
1. Define OpenAPI specification for Workout resource
2. Scaffold Go module with project structure
3. Implement in-memory store for rapid prototyping

**Deliverables:**
- `openapi/openapi.yaml` - API specification
- `internal/domain/workout.go` - Domain entity
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
3. Link Workouts to Programs

### Phase 4: MCP Server

**Goals:**
1. Implement MCP Server in `mcp/` directory
2. Expose tools for AI agents to query and record workouts
3. Distribute as Docker image or Python script

**Key Design:**
- MCP Server runs on user's local machine (not on the backend server)
- Communicates with OPTel API via HTTP
- Backend remains a pure REST API with no MCP-specific logic

## Workflow

### Domain-Driven Schema-First Development

We use a **Domain-Driven Schema-First** approach that combines the benefits of domain modeling with schema-first API development.

**Key Principles:**

1. **Domain models define business logic** - Domain entities (`internal/domain/`) contain business rules, validation, and invariants
2. **OpenAPI spec defines API contract** - OpenAPI specification (`openapi/openapi.yaml`) defines the HTTP API contract for clients
3. **Code generation from OpenAPI** - Go types and server stubs are generated from the OpenAPI spec using `oapi-codegen`
4. **Handlers bridge the gap** - HTTP handlers convert between OpenAPI types and domain models

**Workflow:**

1. **Define domain model** - Start by defining domain entities in `internal/domain/`
   - Include business rules and validation logic
   - Define repository interfaces (ports)
   - Keep domain models focused on business logic, not HTTP concerns
2. **Create OpenAPI spec** - Define API contract in `openapi/openapi.yaml`
   - Reference domain models when designing schemas
   - Define minimal schemas first (only required fields)
   - Add fields incrementally as needed
   - Document domain model references in schema descriptions
3. **Generate code** - Run `make generate`
   - This generates Go types and server stubs from the OpenAPI spec
4. **Implement handlers** - Write handler logic
   - Convert OpenAPI types to domain models
   - Use domain models for business logic
   - Convert domain models back to OpenAPI types for responses
5. **Add tests** - Write table-driven tests
   - Test domain logic separately
   - Test conversion between OpenAPI and domain types
   - Test handlers with integration tests
6. **Run linter** - Run `make lint`
7. **Test** - Run `make test`

## Incremental Development Guide

This guide shows how to build the API incrementally, starting with a minimal working specification and gradually adding features. Each step should result in a working, testable system.

### Domain Model and OpenAPI Spec Synchronization

**Key Principle:** Domain models and OpenAPI specs are maintained separately but must stay in sync.

- **Domain models first** - Define domain entities with business rules first
- **OpenAPI spec follows** - Create OpenAPI spec referencing domain models
- **Keep them synchronized** - When modifying one, update the other
- **Use conversion layer** - Handlers convert between OpenAPI types and domain models

**Synchronization Checklist:**

When modifying domain models or OpenAPI specs:
- [ ] Update domain model (`internal/domain/`)
- [ ] Update OpenAPI spec (`openapi/openapi.yaml`)
  - [ ] Required fields match
  - [ ] Validation rules match
  - [ ] Types match (int, string, etc.)
- [ ] Run code generation (`make generate`)
- [ ] Update conversion layer in handlers
- [ ] Update tests
- [ ] Run validation (`make validate-openapi`)

**Note:** The domain reference (`.claude/skills/optel-domain/reference.md`) shows the **target state** or **reference implementation**, not a requirement to implement everything upfront. Use it as inspiration, but feel free to start simpler and build up.

### Development Cycle

For each feature addition, follow this cycle:

1. **Define domain model** - Add domain entity with business rules
2. **Create/update OpenAPI spec** - Define API contract referencing domain model
3. **Generate code** - Run `make generate`
4. **Implement repository** - Create repository implementation (memory/postgres)
5. **Implement handler** - Convert between OpenAPI types and domain models
6. **Test manually** - Use `curl` or similar to verify it works
7. **Add tests** - Write table-driven tests for domain logic and handlers
8. **Validate sync** - Ensure domain model and OpenAPI spec are synchronized
9. **Commit** - Save working state before next step

### Step-by-Step: Building Workout API

#### Step 1: Minimal OpenAPI Structure

Start with the absolute minimum - just the OpenAPI info and a single endpoint:

```yaml
openapi: 3.1.0
info:
  title: OPTel Workout API
  version: 0.1.0
servers:
  - url: http://localhost:8080/api/v1
paths:
  /workouts:
    get:
      summary: List workouts
      operationId: listWorkouts
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
curl http://localhost:8080/api/v1/workouts
# Should return: []

# Commit
git commit -m "feat(api): add minimal GET /workouts endpoint"
```

#### Step 2: Add POST Endpoint with Minimal Schema

Add the ability to create workouts with only required fields:

```yaml
paths:
  /workouts:
    get:
      # ... existing ...
    post:
      summary: Create a workout
      operationId: createWorkout
      requestBody:
        required: true
        content:
          application/json:
            schema:
                $ref: '#/components/schemas/WorkoutCreate'
      responses:
        '201':
          description: Created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Workout'
components:
  schemas:
    WorkoutCreate:
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
    Workout:
      allOf:
        - $ref: '#/components/schemas/WorkoutCreate'
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
curl -X POST http://localhost:8080/api/v1/workouts \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2026-01-24T10:00:00Z",
    "duration_seconds": 3600,
    "intensity_rpe": 7,
    "volume_kg": 5000
  }'

# Verify GET returns the created workout
curl http://localhost:8080/api/v1/workouts

# Commit
git commit -m "feat(api): add POST /workouts with minimal schema"
```

#### Step 3: Add Optional Fields

Extend the schema with optional fields one at a time:

```yaml
    WorkoutCreate:
      type: object
      required:
        - timestamp
        - duration_seconds
        - intensity_rpe
        - volume_kg
      properties:
        # ... existing required fields ...
        muscle_groups:
          type: array
          items:
            type: string
        notes:
          type: string
    Workout:
      allOf:
        - $ref: '#/components/schemas/WorkoutCreate'
        - type: object
          properties:
            id:
              type: string
            created_at:
              type: string
              format: date-time
            # Add optional fields from WorkoutCreate
            muscle_groups:
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
curl -X POST http://localhost:8080/api/v1/workouts \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2026-01-24T10:00:00Z",
    "duration_seconds": 3600,
    "intensity_rpe": 7,
    "volume_kg": 5000,
    "muscle_groups": ["chest", "triceps"],
    "notes": "Bench press session"
  }'

# Commit
git commit -m "feat(api): add optional muscle_groups and notes fields"
```

#### Step 4: Add GET by ID Endpoint

Add the ability to retrieve a specific workout:

```yaml
paths:
  /workouts:
    # ... existing ...
  /workouts/{id}:
    get:
      summary: Get a workout by ID
      operationId: getWorkout
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
                $ref: '#/components/schemas/Workout'
        '404':
          description: Not found
```

**Action:**
```bash
make generate
# Implement handler
# Test it works
WORKOUT_ID=$(curl -s -X POST http://localhost:8080/api/v1/workouts \
  -H "Content-Type: application/json" \
  -d '{"timestamp":"2026-01-24T10:00:00Z","duration_seconds":3600,"intensity_rpe":7,"volume_kg":5000}' \
  | jq -r '.id')

curl http://localhost:8080/api/v1/workouts/$WORKOUT_ID

# Commit
git commit -m "feat(api): add GET /workouts/{id} endpoint"
```

#### Step 5: Add Query Parameters for Filtering

Add filtering capabilities to the list endpoint:

```yaml
  /workouts:
    get:
      summary: List workouts
      operationId: listWorkouts
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
        - name: muscle_group
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
curl "http://localhost:8080/api/v1/workouts?from=2026-01-01T00:00:00Z&to=2026-01-31T23:59:59Z"
curl "http://localhost:8080/api/v1/workouts?muscle_group=chest"

# Commit
git commit -m "feat(api): add query parameters for filtering workouts"
```

#### Step 6: Add DELETE Endpoint

Add soft-delete capability:

```yaml
  /workouts/{id}:
    get:
      # ... existing ...
    delete:
      summary: Delete a workout
      operationId: deleteWorkout
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
curl -X DELETE http://localhost:8080/api/v1/workouts/$WORKOUT_ID
curl http://localhost:8080/api/v1/workouts/$WORKOUT_ID
# Should return 404

# Commit
git commit -m "feat(api): add DELETE /workouts/{id} endpoint"
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
4. **Validate frequently** - Run `make validate-openapi`, `make generate`, and `make test` after each change
5. **Maintain sync** - When modifying domain models, update OpenAPI spec and vice versa
6. **Document decisions** - Add comments in spec and domain models for non-obvious choices
7. **Reference domain models** - Include domain model references in OpenAPI schema descriptions

### Making Changes

```bash
# 1. Create a feature branch
git checkout -b feature/add-workout-filtering

# 2. Make changes following Domain-Driven Schema-First approach
#    - Update domain model (internal/domain/)
#    - Update OpenAPI spec (openapi/openapi.yaml)

# 3. Validate and generate code
make validate-openapi
make generate

# 4. Run checks
make lint
make test
make check-sync  # Ensure domain model and OpenAPI spec are in sync

# 5. Commit with conventional commit message
git commit -m "feat(api): add workout filtering by date range"
```

## Makefile Targets

```makefile
# api/Makefile

.PHONY: generate lint test run build clean validate-openapi check-sync

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

# Validate OpenAPI spec
validate-openapi:
	@echo "Validating OpenAPI spec..."
	@npx @apidevtools/swagger-cli validate openapi/openapi.yaml || \
		(echo "ERROR: OpenAPI spec validation failed. Install swagger-cli: npm install -g @apidevtools/swagger-cli" && exit 1)

# Check domain model and OpenAPI spec synchronization
# Note: This requires a custom script (scripts/check-sync.go) to be implemented
check-sync:
	@echo "Checking domain model and OpenAPI spec sync..."
	@go run scripts/check-sync.go || echo "WARNING: Sync check script not implemented yet"
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
curl http://localhost:8080/api/v1/workouts
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
