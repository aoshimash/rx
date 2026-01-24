# Tasks: Define Core Data Structures

**Input**: Design documents from `/specs/001-data-structure/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**Implementation Scope**: Domain model struct definitions, validation functions, domain error types, OpenAPI schema integration, and unit tests. API handlers and repositories are out of scope for this PR.

**Organization**: Tasks are organized by user story to enable independent implementation and testing, though all entities are foundational and may be implemented together.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., [US1], [US2], [US3])
- Include exact file paths in descriptions

---

## Phase 1: Setup (Project Initialization)

**Purpose**: Initialize Go project structure and dependencies

- [x] T001 Create `api/internal/domain/` directory structure per plan.md
- [x] T002 [P] Update `api/go.mod` with required dependencies: `github.com/google/uuid`, `github.com/oapi-codegen/runtime`
- [x] T003 [P] Verify `api/Makefile` has `generate` target for oapi-codegen

---

## Phase 2: Foundational (Domain Errors & Validation Helpers)

**Purpose**: Core error types and shared validation helpers that all entities depend on

**⚠️ CRITICAL**: No entity implementation can begin until this phase is complete

- [x] T004 Create `api/internal/domain/errors.go` with ValidationError and DomainError types
- [x] T005 [P] Implement `RoundLoad(kg float64) float64` helper in `api/internal/domain/validation.go`
- [x] T006 [P] Implement `ValidateRPE(rpe int) error` helper in `api/internal/domain/validation.go`
- [x] T007 [P] Implement `ValidateFatigueLevel(level int) error` helper in `api/internal/domain/validation.go`
- [x] T008 [P] Implement `ValidateTimestamp(t time.Time) error` helper in `api/internal/domain/validation.go`

**Checkpoint**: Foundation ready - entity implementation can now begin

---

## Phase 3: User Story 1 - Record a Workout Session with Structured Entries (Priority: P1) 🎯 MVP

**Goal**: Define Workout, WorkoutEntry, and Exercise entities with validation to support recording completed training sessions.

**Independent Test**: Create a sample workout with two entries (top + backoff) and validate that the record is valid and round-trips without losing ordering or key attributes.

**Entities Required**: Exercise, Workout, WorkoutEntry, PlanSnapshot (embedded)

### Implementation for User Story 1

- [x] T009 [P] [US1] Create Exercise struct in `api/internal/domain/exercise.go` with all fields from data-model.md
- [x] T010 [P] [US1] Create PlanSnapshot struct in `api/internal/domain/workout.go` (embedded in WorkoutEntry)
- [x] T011 [P] [US1] Create WorkoutEntry struct in `api/internal/domain/workout.go` with all fields from data-model.md
- [x] T012 [US1] Create Workout struct in `api/internal/domain/workout.go` with all fields from data-model.md (depends on T011 for WorkoutEntry)
- [x] T013 [US1] Implement `ValidateExercise(e *Exercise) error` in `api/internal/domain/validation.go` (depends on T009)
- [x] T014 [US1] Implement `ValidateWorkoutEntry(e *WorkoutEntry) error` in `api/internal/domain/validation.go` (depends on T011, T010)
- [x] T015 [US1] Implement `ValidateWorkout(w *Workout) error` in `api/internal/domain/validation.go` (depends on T012, T014)
- [x] T016 [P] [US1] Write table-driven tests for `ValidateExercise` in `api/internal/domain/exercise_test.go` covering all validation rules
- [x] T017 [P] [US1] Write table-driven tests for `ValidateWorkoutEntry` in `api/internal/domain/workout_test.go` covering all validation rules and edge cases
- [x] T018 [US1] Write table-driven tests for `ValidateWorkout` in `api/internal/domain/workout_test.go` covering all validation rules and edge cases (depends on T015)

**Checkpoint**: User Story 1 entities are fully defined, validated, and tested. A workout with entries can be created and validated.

---

## Phase 4: User Story 2 - Attach Time-Series Telemetry to Workouts (Priority: P2)

**Goal**: Define TelemetryPoint entity with validation to support time-series metric recording.

**Independent Test**: Create a set of telemetry points for a single day and validate they can be represented with and without linking to a workout.

**Entities Required**: TelemetryPoint

### Implementation for User Story 2

- [x] T019 [P] [US2] Create TelemetryPoint struct in `api/internal/domain/telemetry.go` with all fields from data-model.md
- [x] T020 [US2] Implement `ValidateTelemetryPoint(t *TelemetryPoint) error` in `api/internal/domain/validation.go` (depends on T019)
- [x] T021 [P] [US2] Write table-driven tests for `ValidateTelemetryPoint` in `api/internal/domain/telemetry_test.go` covering all validation rules and edge cases

**Checkpoint**: User Story 2 entity is fully defined, validated, and tested. Telemetry points can be created with or without workout linkage.

---

## Phase 5: User Story 3 - Represent Planned Training Programs and Link Workouts (Priority: P3)

**Goal**: Define Program and ProgramNode entities with validation to support recursive program tree structure.

**Independent Test**: Create a minimal program with one cycle/week/day/exercise node and validate the structure preserves ordering and hierarchy.

**Entities Required**: Program, ProgramNode

### Implementation for User Story 3

- [x] T022 [P] [US3] Create ProgramNode struct in `api/internal/domain/program.go` with all fields from data-model.md (including prescription fields)
- [x] T023 [US3] Create Program struct in `api/internal/domain/program.go` with all fields from data-model.md (depends on T022 for ProgramNode)
- [x] T024 [US3] Implement `ValidateProgramNode(n *ProgramNode) error` in `api/internal/domain/validation.go` (depends on T022)
- [x] T025 [US3] Implement `ValidateProgram(p *Program) error` in `api/internal/domain/validation.go` (depends on T023, T024)
- [x] T026 [P] [US3] Write table-driven tests for `ValidateProgramNode` in `api/internal/domain/program_test.go` covering all validation rules
- [x] T027 [US3] Write table-driven tests for `ValidateProgram` in `api/internal/domain/program_test.go` covering all validation rules and recursive structure (depends on T025)

**Checkpoint**: User Story 3 entities are fully defined, validated, and tested. Programs with recursive node trees can be created and validated.

---

## Phase 6: OpenAPI Integration & Code Generation

**Purpose**: Integrate entity schemas into OpenAPI spec and verify code generation works

- [x] T028 Merge `specs/001-data-structure/contracts/openapi-entities.yaml` components into `api/openapi/openapi.yaml`
- [x] T029 [P] Verify component names match domain model names (Exercise, Workout, WorkoutEntry, etc.)
- [x] T030 [P] Ensure all required fields from data-model.md are marked as required in OpenAPI schemas
- [ ] T031 Run `make generate` in `api/` directory to generate OpenAPI types
- [ ] T032 Verify generated code compiles without errors
- [ ] T033 [P] Compare generated OpenAPI types with domain models to ensure consistency

**Checkpoint**: OpenAPI schemas are integrated and code generation works correctly

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final validation, documentation, and quality checks

- [ ] T034 [P] Run `make lint` in `api/` directory and fix all linter errors (requires Docker)
- [ ] T035 [P] Verify all validation functions have 100% test coverage (requires go test)
- [ ] T036 [P] Verify all edge cases from spec.md Edge Cases section are covered in tests
- [ ] T037 [P] Run all tests with `go test -v -race ./...` in `api/internal/domain/` (requires Docker/go test)
- [ ] T038 [P] Verify all spec.md FR-001 to FR-017 requirements are met
- [x] T039 [P] Update `api/internal/domain/README.md` (if exists) or create documentation for domain models

**Checkpoint**: All quality gates passed, ready for PR review

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all entity implementation
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories can proceed in parallel (different entities) or sequentially in priority order
  - US1 (P1) is MVP and should be completed first
- **OpenAPI Integration (Phase 6)**: Depends on all entity implementations (Phase 3-5)
- **Polish (Phase 7)**: Depends on all previous phases

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2)
  - Exercise can be implemented independently
  - WorkoutEntry depends on Exercise
  - Workout depends on WorkoutEntry
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Completely independent
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Completely independent

### Within Each User Story

- Struct definitions can be created in parallel (different files)
- Validation functions depend on struct definitions
- Tests depend on validation functions
- Entity complete before moving to next

### Parallel Opportunities

- **Phase 1**: All setup tasks marked [P] can run in parallel
- **Phase 2**: All validation helper functions marked [P] can run in parallel
- **Phase 3 (US1)**: 
  - Exercise, PlanSnapshot, WorkoutEntry structs can be created in parallel
  - Validation tests can be written in parallel
- **Phase 4 (US2)**: All tasks can run independently
- **Phase 5 (US3)**: ProgramNode and Program structs can be created sequentially, but tests can be written in parallel
- **Phase 6**: Schema verification tasks marked [P] can run in parallel
- **Phase 7**: All polish tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all struct definitions for User Story 1 together:
Task: "Create Exercise struct in api/internal/domain/exercise.go"
Task: "Create PlanSnapshot struct in api/internal/domain/workout.go"
Task: "Create WorkoutEntry struct in api/internal/domain/workout.go"

# Launch all validation tests together (after validation functions are done):
Task: "Write table-driven tests for ValidateExercise in api/internal/domain/exercise_test.go"
Task: "Write table-driven tests for ValidateWorkoutEntry in api/internal/domain/workout_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Exercise, Workout, WorkoutEntry)
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Verify entities can be created and validated correctly

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Validate (MVP!)
3. Add User Story 2 → Test independently → Validate
4. Add User Story 3 → Test independently → Validate
5. Integrate OpenAPI schemas → Verify code generation
6. Polish and quality checks → Ready for PR

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Workout entities)
   - Developer B: User Story 2 (TelemetryPoint)
   - Developer C: User Story 3 (Program entities)
3. All entities complete → OpenAPI integration
4. All polish tasks in parallel

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- All structs must match data-model.md field definitions exactly
- All validation rules from data-model.md must be implemented
- All edge cases from spec.md must be covered in tests
- Commit after each entity or logical group
- Stop at any checkpoint to validate independently
- Avoid: vague tasks, same file conflicts, missing validation rules
