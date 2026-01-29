# Tasks: PostgreSQL導入

**Input**: Design documents from `/specs/005-postgresql-setup/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are included using Testcontainers as specified in FR-008.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Monorepo (OPTel)**: Component-specific paths:
  - API: `api/internal/`, `api/pkg/`, `api/cmd/`
  - Migrations: `api/migrations/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and dependency setup

- [x] T001 Add PostgreSQL dependencies to api/go.mod (pgx/v5, pgxpool, golang-migrate, cenkalti/backoff)
- [x] T002 Add Testcontainers dependencies to api/go.mod (testcontainers-go, postgres module)
- [x] T003 Update docker-compose.yml to use PostgreSQL 17 (change from postgres:18 to postgres:17)
- [x] T004 [P] Create api/migrations/ directory structure

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core database infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 Extend database config in api/internal/config/config.go (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE, DB_MAX_CONNS, DB_MIN_CONNS, STORAGE_TYPE)
- [x] T006 Create database connection module with connection pool in api/internal/store/postgres/db.go
- [x] T007 Implement Exponential Backoff retry logic for connection in api/internal/store/postgres/db.go
- [x] T008 Add Makefile targets for database operations in api/Makefile (migrate, migrate-down, migrate-down-1, migrate-status, migrate-create)

**Checkpoint**: Foundation ready - database connection infrastructure established

---

## Phase 3: User Story 1 - 開発環境でのデータベース接続 (Priority: P1) 🎯 MVP

**Goal**: 開発者がローカル環境でPostgreSQLに接続してワークアウトデータを永続化できる

**Independent Test**: `docker compose up -d postgres` でDBを起動し、アプリからCRUD操作が正常に動作することを確認

### Migration Files for User Story 1

- [ ] T009 [P] [US1] Create exercises table migration in api/migrations/000001_create_exercises.up.sql
- [ ] T010 [P] [US1] Create exercises rollback migration in api/migrations/000001_create_exercises.down.sql
- [ ] T011 [P] [US1] Create programs table migration in api/migrations/000002_create_programs.up.sql
- [ ] T012 [P] [US1] Create programs rollback migration in api/migrations/000002_create_programs.down.sql
- [ ] T013 [P] [US1] Create program_nodes table migration in api/migrations/000002_create_programs.up.sql (same file, after programs)
- [ ] T014 [P] [US1] Create workouts and workout_entries tables migration in api/migrations/000003_create_workouts.up.sql
- [ ] T015 [P] [US1] Create workouts rollback migration in api/migrations/000003_create_workouts.down.sql
- [ ] T016 [P] [US1] Create telemetry_points table migration in api/migrations/000004_create_telemetry.up.sql
- [ ] T017 [P] [US1] Create telemetry rollback migration in api/migrations/000004_create_telemetry.down.sql

### PostgreSQL Repository Implementations for User Story 1

- [ ] T018 [P] [US1] Implement ExerciseRepository in api/internal/store/postgres/exercise.go
- [ ] T019 [P] [US1] Implement ProgramRepository in api/internal/store/postgres/program.go
- [ ] T020 [P] [US1] Implement WorkoutRepository in api/internal/store/postgres/workout.go
- [ ] T021 [P] [US1] Implement TelemetryPointRepository in api/internal/store/postgres/telemetry.go

### Integration for User Story 1

- [ ] T022 [US1] Update main.go to initialize PostgreSQL connection pool in api/cmd/server/main.go
- [ ] T023 [US1] Add storage type switch (memory/postgres) in api/cmd/server/main.go
- [ ] T024 [US1] Update health handler to include database connectivity check in api/internal/handler/health.go

**Checkpoint**: User Story 1 complete - Application can connect to PostgreSQL and perform CRUD operations

---

## Phase 4: User Story 2 - データベーススキーマの初期化 (Priority: P2)

**Goal**: 開発者がCLIコマンド（`make migrate`）でデータベーススキーマを適用できる

**Independent Test**: 空のデータベースに対して`make migrate`を実行し、必要なテーブルが作成されることを確認

### Implementation for User Story 2

- [ ] T025 [US2] Create migration runner utility in api/internal/store/postgres/migrate.go
- [ ] T026 [US2] Add embedded migrations support using Go embed in api/internal/store/postgres/migrate.go
- [ ] T027 [US2] Implement migrate CLI command integration in api/Makefile
- [ ] T028 [US2] Add migration status check functionality in api/internal/store/postgres/migrate.go

**Checkpoint**: User Story 2 complete - Migrations can be run via `make migrate`

---

## Phase 5: User Story 3 - テスト用データベース環境 (Priority: P3)

**Goal**: 統合テスト用の独立したTestcontainers環境を構築

**Independent Test**: `make test`を実行し、Testcontainersで一時PostgreSQLコンテナが起動され、テスト後にクリーンアップされることを確認

### Test Infrastructure for User Story 3

- [ ] T029 [US3] Create Testcontainers helper in api/internal/store/postgres/testhelper_test.go
- [ ] T030 [US3] Implement test database setup with automatic migration in api/internal/store/postgres/testhelper_test.go

### Repository Tests for User Story 3

- [ ] T031 [P] [US3] Add ExerciseRepository integration tests in api/internal/store/postgres/exercise_test.go
- [ ] T032 [P] [US3] Add ProgramRepository integration tests in api/internal/store/postgres/program_test.go
- [ ] T033 [P] [US3] Add WorkoutRepository integration tests in api/internal/store/postgres/workout_test.go
- [ ] T034 [P] [US3] Add TelemetryPointRepository integration tests in api/internal/store/postgres/telemetry_test.go

**Checkpoint**: User Story 3 complete - Integration tests run with isolated Testcontainers PostgreSQL

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T035 [P] Add .env.example with database configuration examples
- [ ] T036 [P] Update api/README.md with PostgreSQL setup instructions
- [ ] T037 Run quickstart.md validation (verify all steps work)
- [ ] T038 Verify all acceptance scenarios from spec.md pass

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User Story 1 (P1): Required first - provides core PostgreSQL functionality
  - User Story 2 (P2): Depends on T009-T017 migration files from US1
  - User Story 3 (P3): Depends on T018-T021 repository implementations from US1
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Depends on migration files from US1 (T009-T017) but can be developed in parallel
- **User Story 3 (P3)**: Depends on repository implementations from US1 (T018-T021)

### Within Each User Story

- Migration files can be created in parallel [P]
- Repository implementations can be created in parallel [P]
- Integration tasks (T022-T024) depend on repositories being complete
- Tests depend on repositories being complete

### Parallel Opportunities

- All migration file tasks (T009-T017) can run in parallel
- All repository implementation tasks (T018-T021) can run in parallel
- All repository test tasks (T031-T034) can run in parallel
- Polish tasks (T035-T036) can run in parallel

---

## Parallel Example: User Story 1 - Migrations

```bash
# Launch all migration file tasks together:
Task: "Create exercises table migration in api/migrations/000001_create_exercises.up.sql"
Task: "Create programs table migration in api/migrations/000002_create_programs.up.sql"
Task: "Create workouts table migration in api/migrations/000003_create_workouts.up.sql"
Task: "Create telemetry table migration in api/migrations/000004_create_telemetry.up.sql"
```

## Parallel Example: User Story 1 - Repositories

```bash
# Launch all repository implementation tasks together:
Task: "Implement ExerciseRepository in api/internal/store/postgres/exercise.go"
Task: "Implement ProgramRepository in api/internal/store/postgres/program.go"
Task: "Implement WorkoutRepository in api/internal/store/postgres/workout.go"
Task: "Implement TelemetryPointRepository in api/internal/store/postgres/telemetry.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T004)
2. Complete Phase 2: Foundational (T005-T008)
3. Complete Phase 3: User Story 1 (T009-T024)
4. **STOP and VALIDATE**: Test User Story 1 independently
   - `docker compose up -d postgres`
   - `make migrate`
   - Start application with `STORAGE_TYPE=postgres`
   - Verify CRUD operations work
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Verify `make migrate` workflow
4. Add User Story 3 → Test independently → Verify Testcontainers integration
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: Migration files (T009-T017)
   - Developer B: Repository implementations (T018-T021)
3. After repositories complete:
   - Developer A: Integration tasks (T022-T024)
   - Developer B: Test infrastructure (T029-T034)

---

## Requirement Traceability

| Requirement | Tasks |
|-------------|-------|
| FR-001 (PostgreSQL connection with retry) | T005, T006, T007, T022 |
| FR-002 (Environment variable config) | T005, T035 |
| FR-003 (Health check endpoint) | T024 |
| FR-004 (Migration via CLI) | T008, T025, T026, T027, T028 |
| FR-005 (Connection pooling) | T006 |
| FR-006 (Database error logging) | T006, T007 |
| FR-007 (Docker Compose for dev) | T003 |
| FR-008 (Testcontainers for testing) | T002, T029, T030, T031-T034 |

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Existing memory store remains unchanged - PostgreSQL store is additive
- STORAGE_TYPE environment variable controls which store to use
