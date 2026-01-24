# Implementation Plan: Define Core Data Structures

**Branch**: `001-data-structure` | **Date**: 2026-01-24 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-data-structure/spec.md`

## Summary

Define the core domain entities for OPTel Training: **Workout**, **WorkoutEntry**, **Program**, **ProgramNode**, **Exercise**, and **TelemetryPoint**. These structures form the foundation for all training data storage and retrieval, supporting both performed workout logs and planned training programs with flexible hierarchy.

**Technical Approach**: Domain-Driven Schema-First development—define Go domain models with business rules, then generate OpenAPI spec and API types from the domain model definitions.

## Implementation Scope

This plan covers:
- ✅ Domain model struct definitions (Go types)
- ✅ Validation functions for each entity
- ✅ Domain error types
- ✅ OpenAPI schema components integration
- ✅ Unit tests for validation logic
- ❌ API handlers (future feature)
- ❌ Repository implementations (future feature)
- ❌ Database schema (future feature)

## Technical Context

**Language/Version**: Go 1.25+  
**Primary Dependencies**: chi (HTTP router), oapi-codegen (OpenAPI code generation)  
**Storage**: In-memory initially; PostgreSQL with TimescaleDB (future)  
**Testing**: Table-driven tests (standard Go testing package)  
**Target Platform**: Linux server (containerized)  
**Project Type**: Monorepo - `api/` component  
**Performance Goals**: N/A for data structure definition phase  
**Constraints**: Load precision 0.1 kg; RPE range 1-10; fatigue_level 1-5  
**Scale/Scope**: Personal training logs (single user initially)

## Dependencies

**Required:**
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/google/uuid` - UUID generation (v4)
- `github.com/oapi-codegen/runtime` - OpenAPI runtime utilities
- `github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen` - Code generator (dev tool)

**Future (not in this PR):**
- Database driver (PostgreSQL)
- Repository implementations

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| **Dumb Backend** | ✅ PASS | No health calculations; pure data storage/retrieval |
| **Domain-Driven Schema-First** | ✅ PASS | Domain models → OpenAPI spec → code generation |
| **Terminology** | ✅ PASS | Intuitive terms: Workout, Exercise, Program, etc. |
| **Data Integrity (Immutable)** | ⚠️ VIOLATION | See Complexity Tracking |
| **Clean Architecture** | ✅ PASS | Domain layer separated from infrastructure |
| **Monorepo Structure** | ✅ PASS | Changes contained within `api/` component |

## Project Structure

### Documentation (this feature)

```text
specs/001-data-structure/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (OpenAPI schemas)
│   └── openapi-entities.yaml
├── checklists/
│   └── requirements.md  # Spec quality checklist
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   └── domain/           # Domain models (NEW)
│       ├── workout.go        # Workout, WorkoutEntry
│       ├── program.go        # Program, ProgramNode
│       ├── exercise.go       # Exercise catalog
│       ├── telemetry.go      # TelemetryPoint
│       └── validation.go     # Shared validation rules
├── openapi/
│   └── openapi.yaml      # OpenAPI spec (UPDATE)
├── Makefile
├── go.mod
└── go.sum
```

**Structure Decision**: Domain models in `api/internal/domain/` following clean architecture. OpenAPI spec in `api/openapi/` references domain model definitions.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| **Mutable Workout logs** (vs Constitution's "Immutable Logs") | Personal training logs frequently need corrections (e.g., "entered 100kg but meant 102.5kg"). Immutability adds complexity (new record + linking) without proportional benefit for single-user personal logs. | Immutable logs with correction records rejected because: (1) adds operational complexity for simple corrections, (2) no audit/compliance requirement for personal training data, (3) user explicitly requested simple editing capability. |

## Implementation Phases

### Phase 0: Research ✅ COMPLETE
- [research.md](./research.md) - Technical decisions and patterns

### Phase 1: Design ✅ COMPLETE
- [data-model.md](./data-model.md) - Entity definitions with fields and validation
- [contracts/openapi-entities.yaml](./contracts/openapi-entities.yaml) - OpenAPI schema components
- [quickstart.md](./quickstart.md) - Quick reference for implementers

### Phase 2: Domain Models (THIS PR)
- Implement Go structs
- Implement validation functions
- Write unit tests
- Integrate OpenAPI schemas

### Phase 3: API Implementation (FUTURE)
- Implement handlers
- Implement repositories
- Add integration tests

## Implementation Order

1. **Domain Models** (`internal/domain/*.go`)
   - Define structs for all entities
   - Add field tags for JSON/validation
   
2. **Validation** (`internal/domain/validation.go`)
   - Implement validation functions per entity
   - Define domain error types
   
3. **OpenAPI Integration**
   - Merge `contracts/openapi-entities.yaml` into `api/openapi/openapi.yaml`
   - Ensure schema matches domain models
   
4. **Code Generation**
   - Run `make generate` to generate OpenAPI types
   - Verify generated code compiles
   
5. **Tests**
   - Write table-driven tests for validation
   - Test edge cases from spec.md

## Validation Functions

Each entity requires validation functions:

- `ValidateExercise(e *Exercise) error`
- `ValidateWorkout(w *Workout) error`
- `ValidateWorkoutEntry(e *WorkoutEntry) error`
- `ValidateProgram(p *Program) error`
- `ValidateProgramNode(n *ProgramNode) error`
- `ValidateTelemetryPoint(t *TelemetryPoint) error`

**Common validation helpers:**
- `RoundLoad(kg float64) float64` - Round to 0.1kg precision
- `ValidateRPE(rpe int) error` - Check 1-10 range
- `ValidateFatigueLevel(level int) error` - Check 1-5 range
- `ValidateTimestamp(t time.Time) error` - Check not in future

## Domain Errors

Define error types in `internal/domain/errors.go`:

```go
type ValidationError struct {
    Field   string
    Message string
}

type DomainError struct {
    Code    string
    Message string
    Details map[string]interface{}
}
```

**Error codes:**
- `INVALID_TIMESTAMP` - Timestamp in future
- `INVALID_RPE` - RPE outside 1-10
- `INVALID_FATIGUE_LEVEL` - Fatigue level outside 1-5
- `MISSING_REQUIRED_FIELD` - Required field is empty
- `INVALID_ENTRY_TYPE` - Entry type not in enum

## Testing Strategy

**Unit Tests** (`*_test.go` files):
- Table-driven tests for each validation function
- Test all validation rules from data-model.md
- Test edge cases from spec.md Edge Cases section

**Test Coverage Goals:**
- Validation functions: 100%
- Edge cases: All covered
- Error paths: All covered

**Example test structure:**
```go
func TestValidateWorkout(t *testing.T) {
    tests := []struct {
        name    string
        workout *Workout
        wantErr bool
        errCode string
    }{
        // Test cases from spec.md FR-011
    }
    // ...
}
```

## OpenAPI Spec Integration

**Steps:**
1. Copy `contracts/openapi-entities.yaml` components into `api/openapi/openapi.yaml`
2. Ensure component names match (Exercise, Workout, etc.)
3. Add paths section (minimal for this PR - just schema definitions)
4. Run `make generate` to verify schema is valid
5. Ensure generated types match domain models

**Note:** Full API endpoints (GET/POST/PUT/DELETE) are out of scope for this PR.

## Integration with Existing Code

**Current state:**
- `api/openapi/openapi.yaml` has basic structure with `/workouts` GET endpoint
- `cmd/server/main.go` has placeholder handler

**Integration approach:**
- Merge new entity schemas into existing `openapi.yaml`
- Keep existing `/workouts` path structure
- Add new components to `components.schemas` section
- Update Makefile `generate` target if needed

## Definition of Done

This PR is complete when:
- ✅ All 6 domain entities defined as Go structs
- ✅ All validation functions implemented and tested
- ✅ OpenAPI schemas integrated and code generation works
- ✅ All tests pass (100% validation coverage)
- ✅ No linter errors (`make lint` passes)
- ✅ All spec.md FR-001 to FR-017 requirements met
