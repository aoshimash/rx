# Implementation Plan: REST API for Domain Models

**Branch**: `002-rest-api` | **Date**: 2026-01-24 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-rest-api/spec.md`

## Summary

Expose the existing domain models (Workout, WorkoutEntry, Exercise, Program, ProgramNode, TelemetryPoint) as a REST API with full CRUD operations. The API follows Domain-Driven Schema-First development: domain models already exist, OpenAPI spec defines the API contract, and handlers will bridge between generated types and domain models.

Key decisions from clarification:
- All operations require authentication (401 for unauthenticated)
- Cursor-based pagination for list endpoints (limit + cursor)
- WorkoutEntry and ProgramNode are nested resources only
- Referential integrity enforced: deletion rejected (409) if referenced

Contract strategy notes:
- The single source of truth for the API contract is `api/openapi/openapi.yaml` (used by `make generate` / oapi-codegen).
- The existing `/workouts` GET response shape in `api/openapi/openapi.yaml` is currently an array; it will be updated to the paginated response shape required by the spec (data/next_cursor/has_more).
- `specs/002-rest-api/contracts/openapi-endpoints.yaml` is a design artifact used to draft endpoint definitions and must be merged into `api/openapi/openapi.yaml` before code generation.

## Technical Context

**Language/Version**: Go 1.25+
**Primary Dependencies**: chi (HTTP router), oapi-codegen (OpenAPI code generation)
**Storage**: In-memory store initially (repository pattern enables future PostgreSQL migration)
**Testing**: Standard Go testing package with table-driven tests
**Target Platform**: Linux server (containerized)
**Project Type**: Monorepo - API component (`api/`)
**Performance Goals**: <1s response for single record, <2s for paginated list (up to 100 records)
**Constraints**: 100 records/page max, 1000 records total pagination depth
**Scale/Scope**: MVP for single-user telemetry backend
**Authentication**: Application-level middleware with pluggable providers (stub/JWT/Cognito) configured via environment variables. Same docker image works across docker compose, Kubernetes, and AWS deployments.

**Automation Note**: Speckit scripts may warn if multiple spec directories share the same numeric prefix (e.g., `specs/002-*`). Ensure only one `specs/002-*` directory exists to avoid tooling ambiguity.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| Dumb Backend | ✅ PASS | Pure CRUD operations, no health calculations or business logic |
| Domain-Driven Schema-First | ✅ PASS | Domain models exist, extending OpenAPI spec, code generation |
| Terminology | ✅ PASS | Using established terms: Workout, Exercise, Program, TelemetryPoint |
| Clean Architecture | ✅ PASS | Repository pattern, handlers bridge HTTP and domain |
| Monorepo Structure | ✅ PASS | All changes within `api/` component |

**Result**: All gates pass. Proceeding to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/002-rest-api/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (references 001-data-structure)
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── openapi-endpoints.yaml
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/          # Existing domain models (no changes needed)
│   ├── handler/         # NEW: HTTP handlers for each entity
│   │   ├── exercise.go
│   │   ├── workout.go
│   │   ├── program.go
│   │   ├── telemetry.go
│   │   └── middleware/
│   │       └── auth.go  # Authentication middleware stub
│   ├── repository/      # NEW: Repository interfaces
│   │   ├── exercise.go
│   │   ├── workout.go
│   │   ├── program.go
│   │   └── telemetry.go
│   └── store/
│       └── memory/      # NEW: In-memory implementations
│           ├── exercise.go
│           ├── workout.go
│           ├── program.go
│           └── telemetry.go
├── pkg/
│   └── openapi/         # Generated from OpenAPI spec
├── openapi/
│   └── openapi.yaml     # Extended with all CRUD endpoints
├── go.mod
└── go.sum
```

**Structure Decision**: Extending existing `api/` structure following optel-go-standards. Adding handler, repository, and store layers per Clean Architecture principle.

## Complexity Tracking

> No constitution violations. No complexity justification needed.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| (none) | - | - |
