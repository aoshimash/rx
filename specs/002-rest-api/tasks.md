# Tasks: REST API for Domain Models

**Input**: Design documents from `/specs/002-rest-api/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., [US1], [US2], [US3])
- Include exact file paths in descriptions

## Path Conventions

- **Monorepo (OPTel)**: API component paths:
  - `api/internal/handler/` - HTTP handlers
  - `api/internal/repository/` - Repository interfaces
  - `api/internal/store/memory/` - In-memory implementations
  - `api/internal/middleware/` - Middleware
  - `api/openapi/openapi.yaml` - OpenAPI specification
  - `api/pkg/openapi/` - Generated code (from oapi-codegen)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create directory structure for handlers, repositories, and stores in `api/internal/`
- [x] T002 [P] Verify chi and oapi-codegen dependencies in `api/go.mod`
- [x] T003 [P] Create `api/internal/config/config.go` for application configuration (auth provider, etc.)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Extend `api/openapi/openapi.yaml` with all CRUD endpoints from `specs/002-rest-api/contracts/openapi-endpoints.yaml` (merge paths, components, security schemes)
- [x] T005 Update `/workouts` GET response in `api/openapi/openapi.yaml` from array to paginated response shape (data/next_cursor/has_more)
- [x] T006 Run `make generate` in `api/` to generate OpenAPI server code in `api/pkg/openapi/server.gen.go`
- [x] T007 [P] Create `api/internal/middleware/auth.go` with AuthProvider interface and stub provider implementation
- [x] T008 [P] Create `api/internal/middleware/error.go` for error response formatting (Error struct with code/message/details)
- [x] T009 [P] Create `api/internal/repository/exercise.go` with ExerciseRepository interface (Create, GetByID, Update, Delete, List)
- [x] T010 [P] Create `api/internal/repository/workout.go` with WorkoutRepository interface (Create, GetByID, Update, Delete, List, ListByTimestampRange)
- [x] T011 [P] Create `api/internal/repository/program.go` with ProgramRepository interface (Create, GetByID, Update, Delete, List)
- [x] T012 [P] Create `api/internal/repository/telemetry.go` with TelemetryPointRepository interface (Create, GetByID, Update, Delete, List, ListByMetricAndTimeRange)
- [x] T013 [P] Create `api/internal/store/memory/exercise.go` implementing ExerciseRepository with in-memory map storage
- [x] T014 [P] Create `api/internal/store/memory/workout.go` implementing WorkoutRepository with in-memory map storage
- [x] T015 [P] Create `api/internal/store/memory/program.go` implementing ProgramRepository with in-memory map storage
- [x] T016 [P] Create `api/internal/store/memory/telemetry.go` implementing TelemetryPointRepository with in-memory map storage
- [x] T017 Update `api/cmd/server/main.go` to initialize chi router, register authentication middleware, and wire up generated OpenAPI handlers

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Create and Retrieve Training Data Records (Priority: P1) 🎯 MVP

**Goal**: Enable clients to create and retrieve all entity types (Exercise, Workout, Program, TelemetryPoint). This is the minimum viable product.

**Independent Test**: Create a single record of each entity type, then retrieve it by ID. All operations should work independently without other stories.

### Implementation for User Story 1

- [x] T018 [P] [US1] Create `api/internal/handler/exercise.go` with CreateExercise and GetExercise handlers (convert OpenAPI types to domain models, call repository, return responses)
- [x] T019 [P] [US1] Create `api/internal/handler/workout.go` with CreateWorkout and GetWorkout handlers (handle nested WorkoutEntry creation, convert types)
- [x] T020 [P] [US1] Create `api/internal/handler/program.go` with CreateProgram and GetProgram handlers (handle nested ProgramNode tree creation, convert types)
- [x] T021 [P] [US1] Create `api/internal/handler/telemetry.go` with CreateTelemetryPoint and GetTelemetryPoint handlers (convert types, validate workout_id reference)
- [x] T022 [US1] Register Exercise handlers in `api/cmd/server/main.go` (POST /exercises, GET /exercises/{id})
- [x] T023 [US1] Register Workout handlers in `api/cmd/server/main.go` (POST /workouts, GET /workouts/{id})
- [x] T024 [US1] Register Program handlers in `api/cmd/server/main.go` (POST /programs, GET /programs/{id})
- [x] T025 [US1] Register TelemetryPoint handlers in `api/cmd/server/main.go` (POST /telemetry, GET /telemetry/{id})
- [x] T026 [US1] Implement validation in handlers: check required fields, validate Exercise exists when creating WorkoutEntry, validate Program exists when creating ProgramNode
- [x] T027 [US1] Implement error handling in handlers: return 400 for validation errors, 404 for not found, 401 for unauthenticated (via middleware)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. Can create and retrieve all entity types.

---

## Phase 4: User Story 2 - Update Existing Training Data Records (Priority: P2)

**Goal**: Enable clients to update existing records (Exercise, Workout, Program, TelemetryPoint) using PUT semantics (full replacement).

**Independent Test**: Create a record, update one or more fields using PUT, then retrieve it to verify changes. Works independently of other stories.

### Implementation for User Story 2

- [x] T028 [P] [US2] Implement UpdateExercise handler in `api/internal/handler/exercise.go` (PUT /exercises/{id}, full replacement)
- [x] T029 [P] [US2] Implement UpdateWorkout handler in `api/internal/handler/workout.go` (PUT /workouts/{id}, replace workout including entries)
- [x] T030 [P] [US2] Implement UpdateProgram handler in `api/internal/handler/program.go` (PUT /programs/{id}, replace program including node tree)
- [x] T031 [P] [US2] Implement UpdateTelemetryPoint handler in `api/internal/handler/telemetry.go` (PUT /telemetry/{id}, full replacement)
- [x] T032 [US2] Register update endpoints in `api/cmd/server/main.go` (PUT routes for all entities)
- [x] T033 [US2] Add validation in update handlers: check record exists (404 if not), validate input data, validate references (exercise_id, workout_id, etc.)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently. Can create, retrieve, and update all entity types.

---

## Phase 5: User Story 3 - Delete Training Data Records (Priority: P2)

**Goal**: Enable clients to delete records with referential integrity checks. Parent entities cascade delete children (Workout→Entry, Program→Node). Referenced entities reject deletion (409 Conflict).

**Independent Test**: Create a record, delete it by identifier, then attempt to retrieve it to confirm deletion. Test referential integrity: try to delete Exercise referenced by WorkoutEntry (should return 409).

### Implementation for User Story 3

- [x] T034 [P] [US3] Implement DeleteExercise handler in `api/internal/handler/exercise.go` (DELETE /exercises/{id}, check for WorkoutEntry references, return 409 if referenced)
- [x] T035 [P] [US3] Implement DeleteWorkout handler in `api/internal/handler/workout.go` (DELETE /workouts/{id}, cascade delete WorkoutEntries)
- [x] T036 [P] [US3] Implement DeleteProgram handler in `api/internal/handler/program.go` (DELETE /programs/{id}, check for Workout references, cascade delete ProgramNodes, return 409 if referenced by Workout)
- [x] T037 [P] [US3] Implement DeleteTelemetryPoint handler in `api/internal/handler/telemetry.go` (DELETE /telemetry/{id}, no dependencies)
- [x] T038 [US3] Add referential integrity checks in repository implementations: WorkoutRepository.ListByExerciseID for Exercise deletion check, WorkoutRepository.ListByProgramNodeID for Program deletion check
- [x] T039 [US3] Register delete endpoints in `api/cmd/server/main.go` (DELETE routes for all entities)
- [x] T040 [US3] Implement 409 Conflict error response format in handlers with blocking_references details

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should all work independently. Can create, retrieve, update, and delete all entity types with proper referential integrity.

---

## Phase 6: User Story 4 - List and Filter Training Data Records (Priority: P3)

**Goal**: Enable clients to list all records with cursor-based pagination and filtering (Workouts by timestamp range, TelemetryPoints by metric_name and timestamp range).

**Independent Test**: Create multiple records, then list them with pagination. Test filtering: list workouts by date range, list telemetry points by metric name and date range.

### Implementation for User Story 4

- [x] T041 [P] [US4] Implement ListExercises handler in `api/internal/handler/exercise.go` (GET /exercises with limit/after pagination parameters)
- [x] T042 [P] [US4] Implement ListWorkouts handler in `api/internal/handler/workout.go` (GET /workouts with limit/after pagination and timestamp_from/timestamp_to filters)
- [x] T043 [P] [US4] Implement ListPrograms handler in `api/internal/handler/program.go` (GET /programs with limit/after pagination parameters)
- [x] T044 [P] [US4] Implement ListTelemetryPoints handler in `api/internal/handler/telemetry.go` (GET /telemetry with limit/after pagination and metric_name/timestamp_from/timestamp_to filters)
- [x] T045 [US4] Implement cursor-based pagination logic in handlers: encode last record ID as base64 cursor, calculate has_more, return PaginatedResponse structure
- [x] T046 [US4] Implement filtering logic in repository methods: WorkoutRepository.ListByTimestampRange, TelemetryPointRepository.ListByMetricAndTimeRange
- [x] T047 [US4] Add pagination parameter validation: limit 1-100 (default 100), validate cursor format
- [x] T048 [US4] Register list endpoints in `api/cmd/server/main.go` (GET routes with query parameters)

**Checkpoint**: At this point, all user stories should be complete. Can create, retrieve, update, delete, and list all entity types with pagination and filtering.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T049 [P] Add structured logging (log/slog) to all handlers in `api/internal/handler/`
- [x] T050 [P] Add request ID middleware in `api/internal/middleware/request_id.go` for tracing
- [x] T051 [P] Add CORS middleware in `api/internal/middleware/cors.go` if needed for frontend integration
- [x] T052 [P] Add request validation middleware in `api/internal/middleware/validate.go` for common validations (UUID format, etc.)
- [x] T053 [P] Update `api/internal/domain/errors.go` to include all error codes (UNAUTHORIZED, VALIDATION_ERROR, NOT_FOUND, CONFLICT, INTERNAL_ERROR)
- [x] T054 [P] Add payload size limits validation: max 500 entries per Workout, max 1000 nodes per Program, 10MB request body limit
- [x] T055 [P] Add unit tests for repository implementations in `api/internal/store/memory/*_test.go` (table-driven tests per optel-go-standards)
- [x] T056 [P] Add integration tests for handlers in `api/internal/handler/*_test.go` (test full request/response cycle)
- [x] T057 Run `make lint` in `api/` and fix any golangci-lint issues
- [x] T058 Run `make test` in `api/` and ensure all tests pass
- [x] T059 Validate quickstart.md examples: test all API endpoints manually or with integration tests
- [x] T060 Update `api/README.md` or create `api/docs/API.md` with API usage documentation
- [x] T061 [P] Add load testing to validate 100 concurrent clients requirement (SC-006): create load test script or integration test that simulates 100 concurrent clients performing CRUD operations and verify response times remain within acceptable limits
- [x] T062 [P] Add reliability validation for 95% success rate requirement (SC-008): create integration test suite that performs a large number of valid API operations and verify at least 95% complete successfully on first attempt, or document acceptance criteria for manual validation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
  - T004-T006 (OpenAPI extension and code generation) must complete before handlers can be implemented
  - T007-T016 (Repository interfaces and implementations) must complete before handlers can use them
  - T017 (Server setup) must complete to wire everything together
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can proceed sequentially in priority order (P1 → P2 → P2 → P3)
  - Or in parallel if team capacity allows (each story is independent)
- **Polish (Phase 7)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 handlers being registered but can be implemented independently
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Depends on US1/US2 for testing but implementation is independent
- **User Story 4 (P3)**: Can start after Foundational (Phase 2) - Depends on repository List methods but can be implemented independently

### Within Each User Story

- Handlers depend on repository interfaces (from Phase 2)
- Handler registration depends on handler implementation
- Validation and error handling should be added as handlers are implemented
- Story complete before moving to next priority

### Parallel Opportunities

- **Phase 1**: T002 and T003 can run in parallel
- **Phase 2**: 
  - T007-T012 (Repository interfaces) can all run in parallel
  - T013-T016 (Repository implementations) can all run in parallel
  - T004-T005 (OpenAPI updates) must complete before T006 (code generation)
- **Phase 3 (US1)**: T018-T021 (Handler implementations) can run in parallel
- **Phase 4 (US2)**: T028-T031 (Update handlers) can run in parallel
- **Phase 5 (US3)**: T034-T037 (Delete handlers) can run in parallel
- **Phase 6 (US4)**: T041-T044 (List handlers) can run in parallel
- **Phase 7**: Most tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all handler implementations for User Story 1 together:
Task: "Create api/internal/handler/exercise.go with CreateExercise and GetExercise handlers"
Task: "Create api/internal/handler/workout.go with CreateWorkout and GetWorkout handlers"
Task: "Create api/internal/handler/program.go with CreateProgram and GetProgram handlers"
Task: "Create api/internal/handler/telemetry.go with CreateTelemetryPoint and GetTelemetryPoint handlers"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Create and Retrieve)
4. **STOP and VALIDATE**: Test User Story 1 independently
   - Create Exercise, Workout, Program, TelemetryPoint
   - Retrieve each by ID
   - Verify all operations work
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add User Story 4 → Test independently → Deploy/Demo
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Exercise + Workout)
   - Developer B: User Story 1 (Program + TelemetryPoint)
   - Developer C: User Story 2 (all entities)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Domain models already exist in `api/internal/domain/` - no changes needed
- OpenAPI code generation (`make generate`) must run after OpenAPI spec changes
- Repository pattern enables future PostgreSQL migration without changing handlers
- Authentication middleware (stub provider) is implemented in Phase 2, can be upgraded to JWT/Cognito later
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
