# Quickstart: REST API for Domain Models

**Feature**: 002-rest-api
**Date**: 2026-01-24

## Prerequisites

- Docker and Docker Compose installed
- curl or similar HTTP client for testing

## Start the API Server

```bash
# From repository root
docker compose up api
```

The API will be available at `http://localhost:8080/api/v1`.

## Authentication

All endpoints require authentication. Include an `Authorization` header:

```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/exercises
```

> Note: For MVP, any non-empty Bearer token is accepted. Production will implement proper JWT/OAuth2.

## Quick Examples

### 1. Create an Exercise

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

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Bench Press",
  "description": "Barbell bench press",
  "muscle_groups": ["pectoral", "triceps", "anterior_deltoid"],
  "created_at": "2026-01-24T10:00:00Z",
  "updated_at": "2026-01-24T10:00:00Z"
}
```

### 2. Create a Workout with Entries

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

### 3. List Workouts with Pagination

```bash
# First page
curl -H "Authorization: Bearer test-token" \
  "http://localhost:8080/api/v1/workouts?limit=10"

# Next page (using cursor from previous response)
curl -H "Authorization: Bearer test-token" \
  "http://localhost:8080/api/v1/workouts?limit=10&after=eyJpZCI6IjU1MGU4NDAwLi4uIn0="
```

### 4. Filter Workouts by Date Range

```bash
curl -H "Authorization: Bearer test-token" \
  "http://localhost:8080/api/v1/workouts?timestamp_from=2026-01-01T00:00:00Z&timestamp_to=2026-01-31T23:59:59Z"
```

### 5. Create a Program with Nodes

```bash
curl -X POST http://localhost:8080/api/v1/programs \
  -H "Authorization: Bearer test-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Strength Block",
    "description": "4-week strength focus",
    "root_nodes": [
      {
        "name": "Week 1",
        "node_type": "week",
        "order": 0,
        "children": [
          {
            "name": "Day 1 - Upper",
            "node_type": "day",
            "order": 0,
            "children": [
              {
                "name": "Bench Press 4x8",
                "node_type": "exercise",
                "order": 0,
                "exercise_id": "550e8400-e29b-41d4-a716-446655440000",
                "target_sets": 4,
                "target_reps": 8,
                "target_rpe": 8
              }
            ]
          }
        ]
      }
    ]
  }'
```

### 6. Create a Telemetry Point

```bash
curl -X POST http://localhost:8080/api/v1/telemetry \
  -H "Authorization: Bearer test-token" \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2026-01-24T23:59:59Z",
    "metric_name": "daily_volume_kg",
    "value": 5000,
    "unit": "kg"
  }'
```

## API Endpoints Summary

### Exercise
| Method | Path | Description |
|--------|------|-------------|
| GET | `/exercises` | List exercises (paginated) |
| POST | `/exercises` | Create exercise |
| GET | `/exercises/{id}` | Get exercise |
| PUT | `/exercises/{id}` | Update exercise |
| DELETE | `/exercises/{id}` | Delete exercise |

### Workout
| Method | Path | Description |
|--------|------|-------------|
| GET | `/workouts` | List workouts (paginated, filterable) |
| POST | `/workouts` | Create workout with entries |
| GET | `/workouts/{id}` | Get workout with entries |
| PUT | `/workouts/{id}` | Update workout |
| DELETE | `/workouts/{id}` | Delete workout (cascades entries) |

### Program
| Method | Path | Description |
|--------|------|-------------|
| GET | `/programs` | List programs (paginated) |
| POST | `/programs` | Create program with nodes |
| GET | `/programs/{id}` | Get program with node tree |
| PUT | `/programs/{id}` | Update program |
| DELETE | `/programs/{id}` | Delete program (cascades nodes) |

### Telemetry
| Method | Path | Description |
|--------|------|-------------|
| GET | `/telemetry` | List telemetry points (paginated, filterable) |
| POST | `/telemetry` | Create telemetry point |
| GET | `/telemetry/{id}` | Get telemetry point |
| PUT | `/telemetry/{id}` | Update telemetry point |
| DELETE | `/telemetry/{id}` | Delete telemetry point |

## Error Handling

### 401 Unauthorized
```json
{"code": "UNAUTHORIZED", "message": "Authentication required"}
```

### 400 Validation Error
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Invalid request body",
  "details": {"field": "name", "reason": "required field missing"}
}
```

### 404 Not Found
```json
{"code": "NOT_FOUND", "message": "Exercise not found"}
```

### 409 Conflict (Referential Integrity)
```json
{
  "code": "CONFLICT",
  "message": "Cannot delete exercise - referenced by workout entries",
  "details": {"blocking_references": [{"type": "workout_entry", "count": 5}]}
}
```

## Running Tests

```bash
# From api/ directory
make test
```

## Development Workflow

1. Modify OpenAPI spec in `api/openapi/openapi.yaml`
2. Regenerate code: `make generate`
3. Implement handlers in `api/internal/handler/`
4. Run tests: `make test`
5. Start server: `make run` or `docker compose up api`
