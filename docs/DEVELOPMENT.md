# Development Guide

## Prerequisites

**Local Requirements:**
- Docker & Docker Compose
- [aqua](https://aquaproj.github.io/) - Tool version manager
- curl (for testing API endpoints)

## Project Setup

#### 1. Install aqua

Install aqua following the [official installation guide](https://aquaproj.github.io/docs/tutorial-basics/installation):

```bash
# macOS (Homebrew)
brew install aquaproj/aqua/aqua

# Or using the installer script
curl -sSfL https://raw.githubusercontent.com/aquaproj/aqua-installer/v4.0.0/aqua-installer | bash
```

#### 2. Install development tools via aqua

This project uses aqua to manage tool versions, ensuring consistency between local development and CI environments.

```bash
# Install all tools defined in aqua.yaml
aqua install

# Verify tools are available
go version
oapi-codegen --version
golangci-lint --version
gh --version
```

**Note:** After installing aqua, tools are installed to `~/.local/share/aquaproj-aqua/bin` by default. Make sure this directory is in your PATH, or use `aqua exec` to run commands.

#### 3. Start PostgreSQL with Docker Compose

Start PostgreSQL and other services:

```bash
# Start PostgreSQL (and API service for smoke-testing)
docker compose up -d postgres

# Verify PostgreSQL is running
docker compose ps
```

PostgreSQL connection details:
- Host: `localhost`
- Port: `5432` (configurable via `POSTGRES_PORT` environment variable)
- User: `rx` (configurable via `POSTGRES_USER`)
- Password: `rx` (configurable via `POSTGRES_PASSWORD`)
- Database: `rx_workout` (configurable via `POSTGRES_DB`)

### 4. Development Workflow

```bash
cd api

# Install dependencies
task deps

# Generate OpenAPI code
task generate

# Run linter
task lint

# Run tests
task test

# Start development server
task run

# In another terminal, test the API
curl http://localhost:8080/api/v1/workouts
```

All `task` commands (generate, lint, test, run, build) work natively on your host machine using tools managed by aqua.

#### 4. Setup Git Hooks (lefthook)

This repository uses [lefthook](https://github.com/evilmartians/lefthook) to manage git hooks. Run once after cloning:

```bash
# Install all tools (including lefthook)
aqua install

# Register git hooks
lefthook install
```

> **Upgrading from the old `githooks/` setup?** Run `git config --unset core.hooksPath` first.

**What the hooks do:**

| Hook | Trigger | Checks |
|---|---|---|
| pre-commit | `api/**` staged | `task format`, `task lint`, `task test` |
| pre-commit | `web/**` staged | `pnpm check` |
| pre-push | `api/openapi/openapi.yaml` changed | `task generate` + uncommitted-changes check |

The pre-push hook catches stale generated code (the same check run by CI) before it reaches GitHub.

**Manual check before committing:**

```bash
cd api
task check  # Runs generate + format + lint + test
```

**Note:** The pre-commit hook is enforced for all commits, including AI agent commits. Do not use `git commit --no-verify` to skip hooks.

### Docker Compose Services

Docker Compose provides PostgreSQL for local development and the API service for smoke-testing. See the "Docker Compose" section below for details.

## API Documentation

The project includes an interactive API documentation viewer powered by Scalar API Reference. This allows developers to view and test API endpoints directly from their browser.

### Starting the Documentation Server

Start the API documentation server:

```bash
docker compose up -d api-docs
```

The documentation will be available at:

```
http://localhost:8081
```

### Using the Documentation

Once the server is running:

1. Open `http://localhost:8081` in your browser
2. Browse all available API endpoints with their request/response schemas
3. Use the "Try it" feature to send test requests to a running API server

### Testing API Endpoints

To use the "Try it" feature, ensure both the API server and documentation server are running:

```bash
# Start both services
docker compose up -d api api-docs

# The documentation UI will allow you to send requests to http://localhost:8080/api/v1
```

The "Try it" feature sends requests directly from your browser to the API server, allowing you to test endpoints without external tools like curl or Postman.

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
- Communicates with Rx API via HTTP
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
3. **Generate code** - Run `task generate`
   - This generates Go types and server stubs from the OpenAPI spec
4. **Implement handlers** - Write handler logic
   - Convert OpenAPI types to domain models
   - Use domain models for business logic
   - Convert domain models back to OpenAPI types for responses
5. **Add tests** - Write table-driven tests
   - Test domain logic separately
   - Test conversion between OpenAPI and domain types
   - Test handlers with integration tests
6. **Run linter** - Run `task lint`
7. **Test** - Run `task test`

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
- [ ] Run code generation (`task generate`)
- [ ] Update conversion layer in handlers
- [ ] Update tests
- [ ] Run validation (`task validate-openapi`)

**Note:** The domain model ([docs/DOMAIN_MODEL.md](DOMAIN_MODEL.md)) shows the **target state** or **reference implementation**, not a requirement to implement everything upfront. Use it as inspiration, but feel free to start simpler and build up.

### Development Cycle

For each feature addition, follow this cycle:

1. **Define domain model** - Add domain entity with business rules
2. **Create/update OpenAPI spec** - Define API contract referencing domain model
3. **Generate code** - Run `task generate`
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
  title: Rx API
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
# Navigate to api directory
cd api

# Generate code (runs natively on host with aqua-managed tools)
task generate

# Create minimal handler that returns empty array
# Start server (runs natively in DevContainer, in background or another terminal)
task run

# Test it works (from local machine or DevContainer terminal)
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
task generate
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
task generate
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
task generate
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
task generate
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
task generate
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
   task generate
   go build ./...
   ```

3. **Handlers work:**
   ```bash
   task run
   # Test with curl in another terminal
   ```

4. **Tests pass:**
   ```bash
   task test
   ```

5. **Linter passes:**
   ```bash
   task lint
   ```

### Best Practices

1. **One feature per commit** - Each step should be a single, focused commit
2. **Test before moving on** - Don't proceed to the next step until current one works
3. **Keep specs minimal** - Only add what's needed for current functionality
4. **Validate frequently** - Run `task validate-openapi`, `task generate`, and `task test` after each change
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
task validate-openapi
task generate

# 4. Run checks
task check  # Runs generate + format + lint + test
task check-sync  # Ensure domain model and OpenAPI spec are in sync

# 5. Commit with conventional commit message
# Pre-commit hook will automatically run format, lint, and test
git commit -m "feat(api): add workout filtering by date range"
```

## Taskfile Tasks

```yaml
# api/Taskfile.yml

version: '3'

tasks:
  generate:
    desc: Generate code from OpenAPI spec
    cmds:
      - oapi-codegen -generate types,chi-server -package openapi -o pkg/openapi/server.gen.go openapi/openapi.yaml

  lint:
    desc: Run linter
    cmds:
      - golangci-lint run ./...

  test:
    desc: Run unit tests with race detection
    cmds:
      - go test -v -race ./...

  check:
    desc: Generate code and run format, lint, and tests (pre-commit check)
    cmds:
      - task: generate
      - task: format
      - task: lint
      - task: test
      - echo "All checks passed!"

  run:
    desc: Start the development server
    cmds:
      - go run cmd/server/main.go

  build:
    desc: Build the server binary
    cmds:
      - go build -o bin/server cmd/server/main.go

  clean:
    desc: Clean build artifacts
    cmds:
      - rm -rf bin/

  # Validate OpenAPI spec
  validate-openapi:
    desc: Validate OpenAPI specification
    cmds:
      - echo "Validating OpenAPI spec..."
      - npx @apidevtools/swagger-cli validate openapi/openapi.yaml

  # Check domain model and OpenAPI spec synchronization
  # Note: This requires a custom script (scripts/check-sync.go) to be implemented
  check-sync:
    desc: Check domain model and OpenAPI spec sync
    cmds:
      - echo "Checking domain model and OpenAPI spec sync..."
      - go run scripts/check-sync.go || echo "WARNING: Sync check script not implemented yet"
```

## Versioning Strategy

This project follows [Semantic Versioning](https://semver.org/) (SemVer):

- **MAJOR** version: Breaking changes
  - Incompatible API changes (e.g., removing endpoints, changing request/response formats)
  - Breaking changes to data structures (field removals, type changes)
  - Principle removals or redefinitions
- **MINOR** version: New features (backward compatible)
  - New endpoints or optional parameters
  - Additive changes to data structures (new optional fields, new entities)
  - New principles or materially expanded guidance
- **PATCH** version: Bug fixes (backward compatible)
  - Bug fixes that don't change API contracts
  - Documentation updates
  - Validation rule clarifications
  - Non-semantic refinements

**API Versioning:**
- API versioning uses URL path pattern `/api/v1/` for the current API version
- When breaking changes are introduced, increment the API version (e.g., `/api/v2/`)
- Breaking changes require explicit migration scripts and version fields on entities

**Data Structure Versioning:**
- Initial implementation: Domain models are stable; all changes will be additive only
- Future breaking changes will require explicit migration scripts and version fields on entities

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `DATABASE_URL` | PostgreSQL connection string (Phase 2) | - |
| `AUTH_PROVIDER` | Authentication provider (stub, jwt, cognito) | `stub` |

### Video Upload Configuration

Video upload functionality requires object storage configuration. When not configured, video endpoints return 503 Service Unavailable.

| Variable | Description | Default |
|----------|-------------|---------|
| `STORAGE_PROVIDER` | Storage provider type: `s3` or `r2` | - (disabled) |
| `STORAGE_BUCKET` | Storage bucket name | - |
| `STORAGE_REGION` | AWS region (for S3) | - |
| `STORAGE_ENDPOINT` | Custom endpoint URL (for R2/MinIO) | - |
| `STORAGE_ACCESS_KEY` | Access key ID | - |
| `STORAGE_SECRET_KEY` | Secret access key | - |
| `VIDEO_MAX_SIZE_MB` | Maximum video file size in MB | `500` |
| `VIDEO_PRESIGN_UPLOAD_EXPIRE_MINUTES` | Upload URL expiration in minutes | `15` |
| `VIDEO_PRESIGN_DOWNLOAD_EXPIRE_MINUTES` | Download URL expiration in minutes | `60` |

#### Example: AWS S3

```bash
export STORAGE_PROVIDER=s3
export STORAGE_BUCKET=my-training-videos
export STORAGE_REGION=ap-northeast-1
export STORAGE_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
export STORAGE_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

#### Example: Cloudflare R2

```bash
export STORAGE_PROVIDER=r2
export STORAGE_BUCKET=my-training-videos
export STORAGE_ENDPOINT=https://ACCOUNT_ID.r2.cloudflarestorage.com
export STORAGE_ACCESS_KEY=your-access-key
export STORAGE_SECRET_KEY=your-secret-key
```

## Testing

### Running Tests

```bash
# Run all tests
task test

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./internal/domain/...
```

### Writing Tests

All tests must be table-driven. See [docs/GO_STANDARDS.md](GO_STANDARDS.md) for examples.

## Docker

### Development Environment

The development environment uses **aqua** for tool version management and **Docker Compose** for PostgreSQL and smoke-testing.

**Benefits:**
- Consistent tool versions across local development and CI
- No need to manually install Go, oapi-codegen, golangci-lint, etc.
- Tools are version-controlled in `aqua.yaml`
- PostgreSQL runs in Docker Compose for easy setup and cleanup

**Workflow:**
1. Install aqua (see Prerequisites above)
2. Run `aqua install` to install all development tools
3. Start PostgreSQL: `docker compose up -d postgres`
4. Run development commands on host: `task generate`, `task lint`, `task test`, `task run`
5. All commands use aqua-managed tools with versions matching CI

**Files:**
- `aqua.yaml` - Tool version definitions (shared with CI)
- `docker-compose.yml` - PostgreSQL and API services

### Docker Compose Services

Docker Compose provides two services:

#### PostgreSQL (Development Database)

**Configuration:**
- Service name: `postgres`
- Image: `postgres:17`
- Default port: `5432` (configurable via `POSTGRES_PORT`)
- Environment variables:
  - `POSTGRES_USER` (default: `rx`)
  - `POSTGRES_PASSWORD` (default: `rx`)
  - `POSTGRES_DB` (default: `rx_workout`)
- Volume: `postgres_data` (persistent data storage)

**Usage:**

```bash
# Start PostgreSQL
docker compose up -d postgres

# View logs
docker compose logs -f postgres

# Stop PostgreSQL
docker compose stop postgres

# Remove PostgreSQL and data (WARNING: deletes all data)
docker compose down -v postgres
```

#### API Service (Smoke-Testing)

**Configuration:**
- Service name: `api`
- Build context: `./api`
- Dockerfile: `Dockerfile` (production build)
- Default port: `8080` (configurable via `HOST_PORT` environment variable)
- Environment variables:
  - `PORT` (default: `8080`) - Container port
  - `LOG_LEVEL` (default: `info`) - Logging level

**Usage:**

```bash
# Start the API service
docker compose up -d api

# View logs
docker compose logs -f api

# Test the API
curl http://localhost:8080/api/v1/workouts

# Stop the service
docker compose down api

# Custom port (if 8080 is in use)
HOST_PORT=8081 docker compose up -d api
```

**Important Notes:**
- The API service is for **smoke-testing** the production container build
- **Do not run development commands** (generate, lint, test, etc.) inside the docker-compose container
- Development commands should be run on the host using aqua-managed tools
- The production container uses a distroless image with no shell, so interactive debugging is not possible
- Use `docker compose logs -f` to view logs

### Building Production Image

**Production images are built from the single `api/Dockerfile` (no separate development Dockerfile).**

```bash
# Build production image (from api/ directory)
docker build -t rx:latest -f Dockerfile api

# Run production container
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e LOG_LEVEL=info \
  rx:latest

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
- Self-contained: builds without requiring source code access after build stage

**Files:**
- `api/Dockerfile` - Single production Dockerfile (optimized, minimal)

### Smoke-Checking Production Image

After building the production image, verify it works correctly:

```bash
# Build the image
docker build -t rx:latest -f Dockerfile api

# Run the container
docker run -d -p 8080:8080 \
  -e PORT=8080 \
  -e LOG_LEVEL=info \
  --name rx-test \
  rx:latest

# Wait a few seconds for startup, then test
sleep 3
curl http://localhost:8080/api/v1/workouts

# Check logs
docker logs rx-test

# Clean up
docker stop rx-test
docker rm rx-test
```

## Troubleshooting

### Common Issues

**aqua tools not found:**
- Ensure aqua is installed: `aqua --version`
- Run `aqua install` to install all tools
- Add `~/.local/share/aquaproj-aqua/bin` to your PATH, or use `aqua exec <command>`
- On macOS/Linux, add to `~/.bashrc` or `~/.zshrc`: `export PATH="$HOME/.local/share/aquaproj-aqua/bin:$PATH"`

**Permission errors:**
- Ensure Docker has proper permissions
- On Linux, you may need to add your user to the docker group: `sudo usermod -aG docker $USER`

**PostgreSQL connection errors:**
- Verify PostgreSQL is running: `docker compose ps`
- Check connection details match environment variables
- View PostgreSQL logs: `docker compose logs postgres`

**oapi-codegen or golangci-lint not found:**
- These tools are managed by aqua
- Run `aqua install` to install them
- Verify installation: `aqua list`
- Ensure aqua bin directory is in PATH

**golangci-lint errors:**
- Review `.golangci.yml` for enabled linters
- Fix issues or add exclusions if false positives
- Run linter: `task lint` (uses aqua-managed golangci-lint)

**Generated code conflicts:**
- Always run `task generate` after changing OpenAPI spec
- Commit generated code with spec changes
- If generation fails, check OpenAPI spec syntax

**Port 8080 already in use:**
- Change port via `HOST_PORT` environment variable: `HOST_PORT=8081 docker compose up -d api`
- Or stop the conflicting service
