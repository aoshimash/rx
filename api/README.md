# OPTel Workout API

REST API for managing workout data (Workouts, Exercises, Programs, and Telemetry).

## Overview

This API provides CRUD operations for the following domain models:
- **Exercise**: Exercise definitions with metadata
- **Workout**: Training sessions with performed entries
- **Program**: Structured training programs with nested nodes
- **TelemetryPoint**: Time-series telemetry data

## Quick Start

### Prerequisites

- **For Development**: [aqua](https://aquaproj.github.io/) (tool version manager)
- **For Smoke-Testing**: Docker and Docker Compose installed

### Development (Recommended)

1. Install development tools via aqua:

```bash
aqua install
```

2. Run development commands:

```bash
cd api
task generate  # Generate code from OpenAPI spec
task run       # Start server
```

The API will be available at `http://localhost:8080/api/v1`.

### Smoke-Testing with Docker Compose

For production-like local testing:

```bash
# From repository root
docker compose up -d
curl http://localhost:8080/api/v1/workouts
docker compose logs -f
docker compose down
```

### Authentication

All endpoints require authentication. Include an `Authorization` header:

```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/exercises
```

> **Note**: For MVP, any non-empty Bearer token is accepted. Production will implement proper JWT/OAuth2.

## API Endpoints

### Exercise

| Method | Path | Description |
|--------|------|-------------|
| GET | `/exercises` | List exercises (paginated) |
| POST | `/exercises` | Create exercise |
| GET | `/exercises/{id}` | Get exercise by ID |
| PUT | `/exercises/{id}` | Update exercise (full replacement) |
| DELETE | `/exercises/{id}` | Delete exercise (rejected if referenced by workouts) |

### Workout

| Method | Path | Description |
|--------|------|-------------|
| GET | `/workouts` | List workouts (paginated, filterable by date range) |
| POST | `/workouts` | Create workout with entries |
| GET | `/workouts/{id}` | Get workout with all entries |
| PUT | `/workouts/{id}` | Update workout (full replacement) |
| DELETE | `/workouts/{id}` | Delete workout (cascades to entries) |

**Query Parameters for List:**
- `limit` (1-100, default: 100): Maximum number of records
- `after`: Cursor for pagination (from previous response)
- `timestamp_from`: Filter workouts at or after this timestamp (RFC3339)
- `timestamp_to`: Filter workouts before this timestamp (RFC3339)

### Program

| Method | Path | Description |
|--------|------|-------------|
| GET | `/programs` | List programs (paginated) |
| POST | `/programs` | Create program with nested nodes |
| GET | `/programs/{id}` | Get program with complete node tree |
| PUT | `/programs/{id}` | Update program (full replacement) |
| DELETE | `/programs/{id}` | Delete program (cascades to nodes) |

### Telemetry

| Method | Path | Description |
|--------|------|-------------|
| GET | `/telemetry` | List telemetry points (paginated, filterable) |
| POST | `/telemetry` | Create telemetry point |
| GET | `/telemetry/{id}` | Get telemetry point by ID |
| PUT | `/telemetry/{id}` | Update telemetry point (full replacement) |
| DELETE | `/telemetry/{id}` | Delete telemetry point |

**Query Parameters for List:**
- `limit` (1-100, default: 100): Maximum number of records
- `after`: Cursor for pagination (from previous response)
- `metric_name`: Filter by metric name (required for filtering)
- `timestamp_from`: Filter points at or after this timestamp (RFC3339)
- `timestamp_to`: Filter points before this timestamp (RFC3339)

### Video Upload

Video upload allows attaching videos to workout entries. Files are stored in external object storage (S3, R2). The API only handles pre-signed URL generation.

| Method | Path | Description |
|--------|------|-------------|
| POST | `/videos/upload-url` | Generate pre-signed URL for uploading a video |
| POST | `/videos/download-url` | Generate pre-signed URL for downloading a video |

**Upload Flow:**
1. Request upload URL with `content_type`, `filename`, and `content_length`
2. Upload video directly to the returned `upload_url` using HTTP PUT
3. Save the returned `object_key` in `WorkoutEntry.video_object_key`

**Download Flow:**
1. Request download URL with `object_key` from `WorkoutEntry.video_object_key`
2. Download video directly from the returned `download_url` using HTTP GET

> **Note**: Video upload requires storage configuration. See `docs/DEVELOPMENT.md` for environment variables.

## Examples

### Create an Exercise

```bash
curl -X POST http://localhost:8080/api/v1/exercises \
  -H "Authorization: Bearer test-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bench Press",
    "description": "Barbell bench press",
    "muscle_groups": ["pectoral", "triceps", "anterior_deltoid"]
  }'
```

### Create a Workout with Entries

```bash
curl -X POST http://localhost:8080/api/v1/workouts \
  -H "Authorization: Bearer test-token" \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2026-01-24T10:00:00Z",
    "body_weight_kg": 80.5,
    "entries": [
      {
        "exercise_id": "550e8400-e29b-41d4-a716-446655440000",
        "entry_type": "main",
        "sets": 4,
        "reps": 8,
        "load_kg": 100,
        "rpe": 8
      }
    ]
  }'
```

### List with Pagination

```bash
# First page
curl -H "Authorization: Bearer test-token" \
  "http://localhost:8080/api/v1/workouts?limit=10"

# Next page (using cursor from previous response)
curl -H "Authorization: Bearer test-token" \
  "http://localhost:8080/api/v1/workouts?limit=10&after=<cursor>"
```

## Error Handling

All errors follow a consistent format:

```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable error message",
  "details": {
    "field": "field_name",
    "message": "Additional details"
  }
}
```

### Error Codes

- **401 Unauthorized**: Missing or invalid authentication
- **400 Validation Error**: Invalid request data
- **404 Not Found**: Resource not found
- **409 Conflict**: Referential integrity violation (e.g., deleting referenced entity)
- **500 Internal Error**: Server error

### Example Error Responses

**401 Unauthorized:**
```json
{"code": "UNAUTHORIZED", "message": "Authentication required"}
```

**400 Validation Error:**
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Invalid request body",
  "details": {"field": "name", "message": "required field cannot be empty"}
}
```

**404 Not Found:**
```json
{"code": "NOT_FOUND", "message": "Exercise not found"}
```

**409 Conflict:**
```json
{
  "code": "CONFLICT",
  "message": "Cannot delete exercise - referenced by workout entries",
  "details": {"blocking_references": [{"type": "workout_entry", "count": 5}]}
}
```

## Development

### Project Structure

```
api/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Domain models and validation
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # HTTP middleware (auth, error handling)
│   ├── repository/      # Repository interfaces
│   └── store/memory/    # In-memory repository implementations
├── openapi/             # OpenAPI specification
└── pkg/openapi/         # Generated OpenAPI code (gitignored)
```

### Taskfile Tasks

```bash
# Generate code from OpenAPI spec
task generate

# Run unit tests
task test

# Run integration tests (requires MinIO)
task test-integration

# Run all tests (unit + integration)
task test-all

# Run linter
task lint

# Start server
task run
```

### Development Workflow

1. Install development tools: `aqua install`
2. Modify OpenAPI spec in `api/openapi/openapi.yaml`
3. Regenerate code: `task generate`
4. Implement handlers in `api/internal/handler/`
5. Run tests: `task test`
6. Start server: `task run`

For production-like smoke-testing, use `docker compose up` from the repository root.

### Testing

```bash
# Run all unit tests
task test

# Run integration tests (requires MinIO)
task test-integration

# Run all tests (unit + integration)
task test-all

# Run specific package tests
go test ./internal/handler -v

# Run quickstart validation tests
go test ./internal/handler -run TestQuickstart -v
```

#### Integration Tests

Integration tests verify real interactions with MinIO (S3-compatible object storage) for video upload/download functionality.

**Prerequisites:**
- Start MinIO: `docker compose up -d minio` (from repository root)
- MinIO console available at http://localhost:9001 (login: minioadmin/minioadmin)

**Running Integration Tests:**
```bash
cd api
task test-integration  # Runs tests with -tags=integration
```

**Test Behavior:**
- Tests automatically skip if MinIO is unavailable (no failures)
- Tests use fixed bucket name `optel-test-videos`
- Test objects are automatically cleaned up via `t.Cleanup()`

**Environment Variables:**
| Variable | Default | Description |
|----------|---------|-------------|
| `MINIO_ENDPOINT` | `http://localhost:9000` | MinIO server endpoint |
| `MINIO_ROOT_USER` | `minioadmin` | MinIO access key |
| `MINIO_ROOT_PASSWORD` | `minioadmin` | MinIO secret key |

For more details, see `specs/004-minio-test-coverage/quickstart.md`.

## Architecture

### Design Principles

- **Schema-First**: OpenAPI spec is the single source of truth
- **Domain-Driven**: Business logic in domain layer, handlers are thin
- **Repository Pattern**: Abstract storage layer for future PostgreSQL migration
- **Dumb Backend**: No business logic for "health" - strictly stores and retrieves telemetry data

### Authentication

Authentication is handled via application-level middleware with pluggable providers:
- **Stub** (default): Development-only, accepts any non-empty Bearer token
- **JWT**: Planned for production
- **AWS Cognito**: Planned for AWS deployments

Configure via `AUTH_PROVIDER` environment variable.

### Storage

The API supports two storage backends:

- **PostgreSQL** (default): Persistent storage using PostgreSQL 17
- **Memory**: In-memory storage for quick testing (data lost on restart)

Configure via `STORAGE_TYPE` environment variable (`postgres` or `memory`).

**PostgreSQL Setup:**

1. Start PostgreSQL: `docker compose up -d postgres` (from repository root)
2. Run migrations: `cd api && make migrate`
3. Start server: `STORAGE_TYPE=postgres make run`

See `specs/005-postgresql-setup/quickstart.md` for detailed setup instructions.

## OpenAPI Specification

The complete API specification is available in `api/openapi/openapi.yaml`. This file is the single source of truth for the API contract.

## License

See [LICENSE](../LICENSE) file.
