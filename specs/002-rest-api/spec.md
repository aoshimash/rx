# Feature Specification: REST API for Domain Models

**Feature Branch**: `002-rest-api`  
**Created**: 2026-01-24  
**Status**: Draft  
**Input**: User description: "Expose the domain models (Workout, WorkoutEntry, Exercise, Program, ProgramNode, TelemetryPoint) as a REST API so that clients can create, read, update, and delete training data records."

## Clarifications

### Session 2026-01-24

- Q: APIの認証要否（アクセス制御）は？ → A: 全操作で認証必須（Create/Read/Update/Delete/Listすべて）
- Q: 参照されているエンティティの削除時の振る舞いは？ → A: 削除拒否（参照が存在する場合は409 Conflictエラーを返す）
- Q: リスト取得時のページネーション方式は？ → A: カーソルベース（limit + cursor/after パラメータ）
- Q: WorkoutEntryのAPI公開方式は？ → A: ネストリソースのみ（/workouts/{id}/entries で管理）
- Q: ProgramNodeのAPI公開方式は？ → A: ネストリソースのみ（/programs/{id}/nodes で管理）

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create and Retrieve Training Data Records (Priority: P1)

A client needs to create training data records (exercises, workouts, programs, telemetry points) and retrieve them later for analysis or display.

**Why this priority**: Core functionality - without the ability to create and read records, the API provides no value. This is the minimum viable product.

**Independent Test**: Can be fully tested by creating a single record of each entity type, then retrieving it by ID. Delivers the ability to store and access training data.

**Acceptance Scenarios**:

1. **Given** no existing exercise record, **When** a client creates an exercise with required fields (name), **Then** the system stores the exercise and returns it with a unique identifier
2. **Given** an existing exercise, **When** a client retrieves it by identifier, **Then** the system returns the complete exercise record
3. **Given** no existing workout record, **When** a client creates a workout with required fields (timestamp, entries), **Then** the system stores the workout and returns it with a unique identifier
4. **Given** an existing workout, **When** a client retrieves it by identifier, **Then** the system returns the complete workout record including all entries
5. **Given** no existing program record, **When** a client creates a program with required fields (name), **Then** the system stores the program and returns it with a unique identifier
6. **Given** an existing program, **When** a client retrieves it by identifier, **Then** the system returns the complete program record including all nodes
7. **Given** no existing telemetry point record, **When** a client creates a telemetry point with required fields (timestamp, metric_name, value, unit), **Then** the system stores the telemetry point and returns it with a unique identifier
8. **Given** an existing telemetry point, **When** a client retrieves it by identifier, **Then** the system returns the complete telemetry point record

---

### User Story 2 - Update Existing Training Data Records (Priority: P2)

A client needs to modify existing training data records to correct errors or update information.

**Why this priority**: Essential for data maintenance and accuracy. Users need to fix mistakes or update records as information becomes available.

**Independent Test**: Can be fully tested by creating a record, updating one or more fields, then retrieving it to verify changes. Delivers the ability to maintain data accuracy.

**Acceptance Scenarios**:

1. **Given** an existing exercise record, **When** a client updates its name or description, **Then** the system updates the record and returns the modified version
2. **Given** an existing workout record, **When** a client updates its notes or condition information, **Then** the system updates the record and returns the modified version
3. **Given** an existing program record, **When** a client updates its name or description, **Then** the system updates the record and returns the modified version
4. **Given** an existing telemetry point record, **When** a client updates its value, **Then** the system updates the record and returns the modified version

---

### User Story 3 - Delete Training Data Records (Priority: P2)

A client needs to remove training data records that are no longer needed or were created in error.

**Why this priority**: Important for data management and cleanup. Users need to remove incorrect or obsolete records.

**Independent Test**: Can be fully tested by creating a record, deleting it by identifier, then attempting to retrieve it to confirm deletion. Delivers the ability to manage data lifecycle.

**Acceptance Scenarios**:

1. **Given** an existing exercise record, **When** a client deletes it by identifier, **Then** the system removes the record and subsequent retrieval attempts return a not found response
2. **Given** an existing workout record, **When** a client deletes it by identifier, **Then** the system removes the record and all associated entries, and subsequent retrieval attempts return a not found response
3. **Given** an existing program record, **When** a client deletes it by identifier, **Then** the system removes the record and all associated nodes, and subsequent retrieval attempts return a not found response
4. **Given** an existing telemetry point record, **When** a client deletes it by identifier, **Then** the system removes the record and subsequent retrieval attempts return a not found response

---

### User Story 4 - List and Filter Training Data Records (Priority: P3)

A client needs to retrieve multiple records at once, potentially filtered by criteria, to display lists or perform bulk operations.

**Why this priority**: Enhances usability by allowing clients to work with collections of records. Not essential for basic CRUD operations but important for practical use cases.

**Independent Test**: Can be fully tested by creating multiple records, then listing them with optional filters. Delivers the ability to browse and search training data.

**Acceptance Scenarios**:

1. **Given** multiple existing exercise records, **When** a client requests a list of all exercises, **Then** the system returns all exercise records
2. **Given** multiple existing workout records, **When** a client requests workouts filtered by date range, **Then** the system returns only workouts within the specified range
3. **Given** multiple existing program records, **When** a client requests a list of all programs, **Then** the system returns all program records
4. **Given** multiple existing telemetry point records, **When** a client requests telemetry points filtered by metric name and date range, **Then** the system returns only matching telemetry points

---

### Edge Cases

| シナリオ | 期待される振る舞い | 関連FR |
|----------|-------------------|--------|
| 無効または必須フィールド欠落でレコード作成 | バリデーションエラーを返す | FR-017, FR-018 |
| 存在しないレコードを取得 | 404 Not Found を返す | FR-019 |
| 存在しないレコードを更新 | 404 Not Found を返す | FR-019 |
| 存在しないレコードを削除 | 404 Not Found を返す | FR-019 |
| 存在しないExerciseを参照するWorkoutEntryを作成 | バリデーションエラーを返す | FR-017, FR-018 |
| 存在しないProgramを参照するProgramNodeを作成 | バリデーションエラーを返す | FR-017, FR-018 |
| WorkoutEntryに参照されているExerciseを削除 | 409 Conflict を返す | FR-026 |
| Workoutに参照されているProgramを削除 | 409 Conflict を返す | FR-026 |
| 同一レコードへの同時更新 | 計画フェーズで詳細化（楽観的ロック等） | - |
| 非常に大きなペイロード（数百エントリのWorkout等） | 計画フェーズで詳細化（上限設定等） | - |

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow clients to create Exercise records with required fields (name) and optional fields (description, aliases, muscle_groups, load_increment)
- **FR-002**: System MUST allow clients to retrieve Exercise records by unique identifier
- **FR-003**: System MUST allow clients to update existing Exercise records
- **FR-004**: System MUST allow clients to delete Exercise records by unique identifier
- **FR-005**: System MUST allow clients to create Workout records with required fields (timestamp, entries) and optional fields (session_start, session_end, body_weight_kg, fatigue_level, sleep_hours, condition_notes, program_node_id, program_context, notes)
- **FR-006**: System MUST allow clients to retrieve Workout records by unique identifier, including all associated WorkoutEntry records
- **FR-007**: System MUST allow clients to update existing Workout records
- **FR-008**: System MUST allow clients to delete Workout records by unique identifier, including all associated WorkoutEntry records
- **FR-009**: System MUST allow clients to create Program records with required fields (name) and optional fields (description, root_nodes)
- **FR-010**: System MUST allow clients to retrieve Program records by unique identifier, including all associated ProgramNode records in the tree structure
- **FR-011**: System MUST allow clients to update existing Program records
- **FR-012**: System MUST allow clients to delete Program records by unique identifier, including all associated ProgramNode records
- **FR-013**: System MUST allow clients to create TelemetryPoint records with required fields (timestamp, metric_name, value, unit) and optional fields (workout_id)
- **FR-014**: System MUST allow clients to retrieve TelemetryPoint records by unique identifier
- **FR-015**: System MUST allow clients to update existing TelemetryPoint records
- **FR-016**: System MUST allow clients to delete TelemetryPoint records by unique identifier
- **FR-017**: System MUST validate all input data according to domain model validation rules before storing
- **FR-018**: System MUST return appropriate error responses when validation fails, including clear error messages
- **FR-019**: System MUST return appropriate error responses when requested records do not exist (404 Not Found)
- **FR-020**: System MUST allow clients to list all records of a given entity type
- **FR-021**: System MUST allow clients to filter records by common criteria (e.g., date ranges, identifiers)
- **FR-022**: System MUST support cursor-based pagination for list endpoints using limit and cursor/after parameters, with a default and maximum limit of 100 records per request, and MUST allow clients to paginate through at least 1000 records for a given entity type
- **FR-023**: System MUST handle relationships between entities (e.g., WorkoutEntry references Exercise, Workout contains WorkoutEntry, Program contains ProgramNode)
- **FR-024**: WorkoutEntry records MUST be managed exclusively as nested resources under Workout (e.g., /workouts/{id}/entries); no independent WorkoutEntry endpoints shall be provided
- **FR-025**: ProgramNode records MUST be managed exclusively as nested resources under Program (e.g., /programs/{id}/nodes); no independent ProgramNode endpoints shall be provided
- **FR-026**: System MUST enforce referential integrity when deleting records by rejecting deletion requests for entities that are referenced by other records and returning a 409 Conflict error with a clear message indicating the blocking references; this does not prohibit deleting a parent entity together with its nested child entities (e.g., deleting a Workout deletes its nested WorkoutEntry records; deleting a Program deletes its nested ProgramNode records)
- **FR-027**: System MUST return records in a consistent format that matches the domain model structure
- **FR-028**: System MUST require authentication for all API operations (Create, Read, Update, Delete, List); unauthenticated requests MUST be rejected with a 401 Unauthorized error response regardless of whether the target resource exists (to prevent information disclosure about resource existence)
- **FR-029**: List operations for nested entities MUST be exposed only under their parent resources (e.g., WorkoutEntry under Workout, ProgramNode under Program)
- **FR-030**: The system MUST support the following minimum filters: Workouts by timestamp range; TelemetryPoints by metric_name and timestamp range

### Key Entities *(include if feature involves data)*

- **Exercise**: A catalog entry representing a canonical exercise. Key attributes: unique identifier, name (required), description, aliases, muscle_groups, load_increment, timestamps. Relationships: referenced by WorkoutEntry.
- **Workout**: A completed training session containing performed entries. Key attributes: unique identifier, timestamp (required), session_start, session_end, body_weight_kg, fatigue_level, sleep_hours, condition_notes, program_node_id, program_context, notes, timestamps. Relationships: contains WorkoutEntry records (1:N), optionally references ProgramNode.
- **WorkoutEntry**: A single performed exercise entry within a workout session. Key attributes: unique identifier, workout_id (required), order (required), exercise_id (required), display_name, entry_type (required), sets (required), reps (required), load_kg (required), rpe (required), entry_start, entry_end, planned_rest_seconds, performed_rest_seconds, per_set_rest_overrides, program_node_id, plan_snapshot, notes. Relationships: belongs to Workout (N:1), references Exercise (N:1), optionally references ProgramNode.
- **Program**: A training program containing a recursive tree of nodes. Key attributes: unique identifier, name (required), description, timestamps. Relationships: contains ProgramNode records in tree structure (1:N).
- **ProgramNode**: A node in the program tree (cycle, week, day, block, or exercise prescription). Key attributes: unique identifier, program_id (required), parent_id, name (required), node_type (required), order (required), children (recursive), exercise_id, target_sets, target_reps, target_rpe, percent_1rm, planned_rest_seconds, muscle_groups, notes. Relationships: belongs to Program (N:1), optionally references parent ProgramNode (recursive), optionally references Exercise.
- **TelemetryPoint**: A time-series metric data point. Key attributes: unique identifier, timestamp (required), metric_name (required), value (required), unit (required), workout_id, created_at (required). Relationships: optionally references Workout.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A client can create a complete Workout record (including at least one WorkoutEntry) via the API in under 5 seconds from request submission to receiving the response
- **SC-002**: A client can retrieve any existing record by identifier in under 1 second from request submission to receiving the response
- **SC-003**: 100% of valid create requests result in successfully stored records that can be immediately retrieved
- **SC-004**: 100% of invalid create/update requests (missing required fields, invalid data types, constraint violations) return appropriate error responses with clear error messages
- **SC-005**: Using pagination, a client can retrieve up to 1000 records of a given entity type; each page (up to 100 records) is returned in under 2 seconds from request submission to receiving the response
- **SC-006**: The API supports at least 100 concurrent clients performing CRUD operations without degradation in response times
- **SC-007**: All entity relationships are correctly maintained (e.g., WorkoutEntry references to Exercise are valid, ProgramNode tree structures are preserved)
- **SC-008**: 95% of API operations complete successfully on first attempt when provided with valid input data
