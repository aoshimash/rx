# Data Model: REST API for Domain Models

**Feature**: 002-rest-api
**Date**: 2026-01-24

## Overview

This feature exposes existing domain models via REST API. The data model is defined in [001-data-structure/data-model.md](../001-data-structure/data-model.md).

This document focuses on API-specific considerations: request/response shapes, pagination, and error handling.

## Entity Reference

See [001-data-structure/data-model.md](../001-data-structure/data-model.md) for complete entity definitions:

- **Exercise** - Catalog entry for canonical exercises
- **Workout** - Completed training session with entries
- **WorkoutEntry** - Single exercise entry within a workout (nested)
- **Program** - Training program with recursive node tree
- **ProgramNode** - Node in program tree (nested)
- **TelemetryPoint** - Time-series metric data point

## API Resource Mapping

### Top-Level Resources (Independent CRUD)

| Entity | Base Path | Operations |
|--------|-----------|------------|
| Exercise | `/api/v1/exercises` | Create, Read, Update, Delete, List |
| Workout | `/api/v1/workouts` | Create, Read, Update, Delete, List |
| Program | `/api/v1/programs` | Create, Read, Update, Delete, List |
| TelemetryPoint | `/api/v1/telemetry` | Create, Read, Update, Delete, List |

### Nested Resources (Managed with Parent)

| Entity | Parent | Management |
|--------|--------|------------|
| WorkoutEntry | Workout | Created/updated via Workout.entries array |
| ProgramNode | Program | Created/updated via Program.root_nodes array |

## Request/Response Schemas

### Create Request Patterns

```yaml
# ExerciseCreate - POST /exercises
{
  "name": "Bench Press",           # required
  "description": "Barbell...",     # optional
  "aliases": ["Flat Bench"],       # optional
  "muscle_groups": ["pectoral"],   # optional
  "load_increment": 2.5            # optional
}

# WorkoutCreate - POST /workouts
{
  "timestamp": "2026-01-24T10:00:00Z",  # required
  "entries": [                           # required, min 1
    {
      "exercise_id": "uuid",
      "entry_type": "main",
      "sets": 4,
      "reps": 8,
      "load_kg": 100,
      "rpe": 8
    }
  ],
  "session_start": "...",               # optional
  "body_weight_kg": 80.5,               # optional
  "notes": "..."                        # optional
}

# ProgramCreate - POST /programs
{
  "name": "Strength Block",        # required
  "description": "...",            # optional
  "root_nodes": [                  # optional
    {
      "name": "Week 1",
      "node_type": "week",
      "order": 0,
      "children": [...]
    }
  ]
}

# TelemetryPointCreate - POST /telemetry
{
  "timestamp": "2026-01-24T10:00:00Z",  # required
  "metric_name": "daily_volume_kg",     # required
  "value": 5000,                        # required
  "unit": "kg",                         # required
  "workout_id": "uuid"                  # optional
}
```

### Response Patterns

All responses include server-generated fields:

```yaml
# Exercise response
{
  "id": "uuid",
  "name": "Bench Press",
  "description": "...",
  "aliases": [...],
  "muscle_groups": [...],
  "load_increment": 2.5,
  "created_at": "2026-01-24T10:00:00Z",
  "updated_at": "2026-01-24T10:00:00Z"
}
```

### Paginated List Response

```yaml
{
  "data": [...],                    # Array of entities
  "next_cursor": "base64string",    # Null if no more pages
  "has_more": true                  # Boolean
}
```

## Pagination Model

### Request Parameters

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `limit` | integer | 100 | 100 | Items per page |
| `after` | string | null | - | Cursor from previous response |

### Response Structure

```go
type PaginatedResponse[T any] struct {
    Data       []T     `json:"data"`
    NextCursor *string `json:"next_cursor,omitempty"`
    HasMore    bool    `json:"has_more"`
}
```

### Cursor Format

- Opaque base64-encoded string
- Contains last record ID for keyset pagination
- Do not parse or construct client-side

## Filter Parameters

### Workout Filters

| Parameter | Type | Description |
|-----------|------|-------------|
| `timestamp_from` | datetime | Include workouts at or after this time |
| `timestamp_to` | datetime | Include workouts before this time |

### TelemetryPoint Filters

| Parameter | Type | Description |
|-----------|------|-------------|
| `metric_name` | string | Filter by exact metric name |
| `timestamp_from` | datetime | Include points at or after this time |
| `timestamp_to` | datetime | Include points before this time |

## Error Model

### Error Response Structure

```yaml
{
  "code": "NOT_FOUND",
  "message": "Exercise not found",
  "details": {
    "id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### Error Codes

| Code | HTTP Status | When |
|------|-------------|------|
| `UNAUTHORIZED` | 401 | Missing/invalid authentication |
| `VALIDATION_ERROR` | 400 | Invalid request body or parameters |
| `NOT_FOUND` | 404 | Resource does not exist |
| `CONFLICT` | 409 | Cannot delete: referenced by other records |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

### Conflict Error Details

When deletion is blocked by references:

```yaml
{
  "code": "CONFLICT",
  "message": "Cannot delete exercise: referenced by workout entries",
  "details": {
    "blocking_references": [
      {"type": "workout_entry", "count": 5}
    ]
  }
}
```

## Validation Rules (API Layer)

### Input Validation

| Field | Rule |
|-------|------|
| `limit` | 1-100, default 100 |
| `timestamp` | Valid ISO 8601 datetime |
| `uuid` fields | Valid UUID v4 format |
| String lengths | Per OpenAPI schema constraints |

### Business Validation (Domain Layer)

Delegated to existing domain models:
- Exercise name uniqueness
- Workout entry constraints
- RPE 1-10 range
- etc.

## Relationship Integrity

### On Create

| Action | Validation |
|--------|------------|
| Create WorkoutEntry | `exercise_id` must exist |
| Create ProgramNode | `exercise_id` (if present) must exist |
| Create TelemetryPoint | `workout_id` (if present) must exist |

### On Delete

| Entity | Check | Behavior |
|--------|-------|----------|
| Exercise | Referenced by WorkoutEntry? | 409 Conflict if yes |
| Program | Referenced by Workout.program_node_id? | 409 Conflict if yes |
| Workout | Has entries | Cascade delete entries |
| Program | Has nodes | Cascade delete nodes |
| TelemetryPoint | None | Always allowed |

## Payload Size Limits

| Constraint | Limit | Rationale |
|------------|-------|-----------|
| Entries per Workout | 500 | Practical session limit |
| Nodes per Program | 1000 | Tree complexity limit |
| Request body size | 10MB | Server configuration |
