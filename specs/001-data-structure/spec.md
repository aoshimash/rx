# Feature Specification: Define Core Data Structures

**Feature Branch**: `001-data-structure`  
**Created**: 2026-01-24  
**Status**: Draft  
**Input**: User description: "データ構造を決める"

## Clarifications

### Session 2026-01-24

- Q: What granularity should rest/interval be recorded? → A: Entry-level planned/performed rest with optional per-set overrides.
- Q: What granularity should time be recorded? → A: Session start/end plus optional per-entry start/end (no per-set timing required).
- Q: How should planned vs performed differences be preserved? → A: Link each performed entry to its planned prescription and store a plan snapshot on the performed entry.
- Q: How should load units be represented? → A: Store load in kilograms only; unit conversion is handled outside of storage.
- Q: How should exercises be identified? → A: Require an exercise catalog ID for all entries; free-text names are display-only.
- Q: What attributes should the Exercise catalog have? → A: Minimal (id, name) with extensibility for future attributes (e.g., description, aliases, muscle_groups, load_increment).
- Q: What precision should load values have? → A: Store to 1 decimal place (0.1 kg); input increment is configurable per exercise (default 2.5 kg or 0.5 kg).
- Q: Should Workout logs be immutable? → A: No. Workouts are editable (corrections can overwrite existing records). Immutability is not required for personal training logs.
- Q: Should Programs be immutable? → A: No. Programs are editable. Historical plan values are preserved via snapshots on WorkoutEntry (FR-009).
- Q: Should body condition (weight, fatigue, etc.) be recorded? → A: Yes, as session-level metadata on Workout (e.g., body_weight_kg, fatigue_level (1-5 scale), sleep_hours, notes).
- Q: What should the entry type field be called? → A: `entry_type` (not `role`, to avoid confusion with user roles).
- Q: How should Program hierarchy be structured? → A: Recursive `ProgramNode` structure with arbitrary depth and user-defined node types (e.g., cycle, week, day, block). Fixed hierarchy (Phase→Day→Block→Slot) is replaced.
- Q: How should Workout show which hierarchy level it belongs to? → A: Store both `program_node_id` (reference) and `program_context` (snapshot of hierarchy path, e.g., ["Cycle 1", "Week 3", "Day 2"]).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Record a Workout Session with Structured Entries (Priority: P1)

As an agent or client integrating with OPTel Workout, I can represent a completed training session as a workout that contains ordered performed entries so that logs remain consistent, searchable, and comparable over time.

**Why this priority**: A canonical workout record is the foundation for all storage, retrieval, and downstream analysis. Without it, every integration becomes bespoke.

**Independent Test**: Provide a sample workout with two entries (top + backoff) and confirm the record is valid and round-trips without losing ordering or key attributes.

**Acceptance Scenarios**:

1. **Given** a workout session timestamp and one or more performed entries (exercise, entry_type, sets, reps, load, RPE), **When** the client records the workout, **Then** the record is accepted and can be retrieved with the same values in the same entry order.
2. **Given** a workout entry that is associated with a planned prescription, **When** the performed values differ from the plan (e.g., load/reps/RPE/rest), **Then** both planned and performed values can be captured in structured fields (not only in freeform notes).
3. **Given** session start/end times and per-entry start/end times are available, **When** the client records the workout, **Then** those times are persisted as optional attributes; **And given** times are not available, **Then** the workout remains valid without them.

---

### User Story 2 - Attach Time-Series Telemetry to Workouts (Priority: P2)

As an agent or analysis tool, I can record and query time-series telemetry points (metrics) that optionally reference a workout so that trends and derived views can be produced without changing the underlying workout logs.

**Why this priority**: Telemetry enables longitudinal analysis and supports additional metrics without forcing changes to the workout entity.

**Independent Test**: Provide a set of telemetry points for a single day and confirm they can be represented using one telemetry-point structure, with and without linking to a workout.

**Acceptance Scenarios**:

1. **Given** a telemetry point (timestamp, metric name, numeric value, unit), **When** it is recorded, **Then** it can be retrieved and interpreted without requiring knowledge of the internal system implementation.
2. **Given** multiple telemetry points that reference the same workout, **When** they are queried, **Then** they are returned as an ordered time-series associated with that workout.

---

### User Story 3 - Represent Planned Training Programs and Link Workouts (Priority: P3)

As an agent managing a training plan, I can represent a structured program (phases → days → blocks → slots) and optionally link completed workouts to that program so that planned vs. executed training can be analyzed.

**Why this priority**: Programs provide context for workouts, but they are optional and should not block recording workouts.

**Independent Test**: Provide a minimal program with one phase/day/block/slot and confirm it is representable; then record a workout that references the program.

**Acceptance Scenarios**:

1. **Given** a minimal program tree, **When** it is recorded, **Then** the structure preserves ordering and hierarchy without ambiguity.
2. **Given** a workout that references a program, **When** the workout is retrieved, **Then** the reference to the program is available as contextual metadata (without modifying the workout log semantics).

---

### Edge Cases

- Recording an entry with load \(= 0\) (e.g., bodyweight/mobility) while still being a valid performed entry.
- Recording a workout whose timestamp is in the future (should be treated as invalid).
- Recording a performed entry with RPE outside the 1–10 range (should be treated as invalid).
- Recording workouts with missing optional fields (muscle groups, notes, program reference).
- Recording a workout with:
  - session start time but no end time,
  - entry start/end times for some entries but not others,
  - overlapping or out-of-order timestamps (must remain interpretable).
- Capturing rest/interval information:
  - planned rest exists but actual rest is unknown,
  - actual rest is known but planned rest is absent,
  - rest differs between sets within the same entry.
- Recording telemetry points with:
  - duplicate timestamps,
  - missing workout linkage,
  - unknown/new metric names,
  - different units for the same metric name (must remain interpretable).
- Correcting previously recorded data by editing existing workout logs directly.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST define canonical entities for training data at minimum: **Workout**, **WorkoutEntry**, **Program**, **ProgramNode**, **TelemetryPoint**, and **Exercise**.
- **FR-001a**: The **Exercise** entity MUST include at minimum: a unique identifier and a name. The structure MUST be extensible to support future attributes (e.g., description, aliases, muscle_groups) without breaking existing references.
- **FR-002**: A **Workout** MUST include: a unique identifier and a timestamp representing when the session occurred.
- **FR-002a**: A **Workout** MAY include session-level condition metadata such as: body weight (kg), fatigue level (1-5 scale), sleep hours, and freeform condition notes. These fields are optional but structured (not embedded in general notes).
- **FR-002b**: A **Workout** MAY include a reference to a Program node (`program_node_id`) and a snapshot of the hierarchy path (`program_context`, e.g., ["Cycle 1", "Week 3", "Day 2"]) so that the training context is immediately visible in the log.
- **FR-003**: A **Workout** MUST contain one or more **WorkoutEntries** that represent what was actually performed. Each entry MUST be representable using the existing CSV row granularity (see “Existing Data Source (CSV)”).
- **FR-004**: Each **WorkoutEntry** MUST include, at minimum: an exercise catalog identifier, an `entry_type` (e.g., top/main/backoff/accessory), an effort/intensity indicator (RPE 1–10), and a performed prescription (sets, reps per set, and load in kilograms) sufficient to reproduce the CSV meaning. A display name MAY be stored for readability.
- **FR-005**: Load (weight) values MUST be stored in kilograms as the canonical unit, with precision to 1 decimal place (0.1 kg). If source data uses other units (e.g., lb), it MUST be converted to kilograms and rounded to 1 decimal place before being stored.
- **FR-006**: The data structure MUST support session start/end timestamps and optional per-entry start/end timestamps so that training time can be recorded when available (per-set timing is not required).
- **FR-007**: The data structure MUST support representing rest/interval information in a structured way, with **entry-level** planned/performed rest values and **optional per-set overrides**:
  - planned rest/interval for a prescribed exercise (Program), and
  - performed rest/interval for an executed entry (Workout).
- **FR-008**: When a Workout is associated with a Program, the data structure MUST allow linking a performed entry to its planned prescription (e.g., a Program Slot reference) and MUST allow recording planned vs performed values (e.g., sets/reps/load/RPE/rest) without relying solely on freeform notes.
- **FR-009**: To preserve historical truth when plans change, the data structure MUST allow storing a **snapshot of the planned values at execution time** on the performed entry, in addition to the reference to the program prescription.
- **FR-010**: Workout logs MAY be edited after recording. Corrections can overwrite existing records directly. (Immutability is not required for personal training logs.)
- **FR-011**: The system MUST reject workout records that violate basic integrity rules: timestamp in the future, missing entries, or entries with invalid values (e.g., sets \(\le 0\), reps \(\le 0\), load \(< 0\), or RPE outside 1–10).
- **FR-012**: A **Program** MUST be representable as a **recursive tree of ProgramNodes** with arbitrary depth and user-defined node types (e.g., "cycle", "week", "day", "block", "exercise"). Each node has at least: an identifier, a name, a `node_type`, and an order within its parent.
- **FR-012a**: Programs MAY be edited after creation. Changes to a Program do not affect plan snapshots already stored on WorkoutEntries.
- **FR-013**: A **ProgramNode** that represents an exercise prescription (leaf node) MAY include planned targets such as: exercise_id, target_sets, target_reps, target_rpe (1–10), `percent_1rm` guidance, and planned_rest, plus associated muscle_groups for context.
- **FR-014**: A **TelemetryPoint** MUST include: timestamp, metric name, numeric value, and unit. It MAY include an optional reference to a workout.
- **FR-015**: Telemetry metric names MUST be treated as opaque identifiers (no health-scoring or wellness logic). The data model MUST allow introduction of new metric names without requiring changes to existing recorded data.
- **FR-016**: The data structures MUST be documented with clear required/optional attributes and validation rules so that independent clients can create valid records without out-of-band guidance.
- **FR-017**: Changes to the data structures MUST preserve interpretability of previously recorded data (e.g., via explicit versioning and change documentation).

### Existing Data Source (CSV)

This feature must support the user’s current training record CSVs (provided here as PDFs). The purpose of capturing this information in the spec is to constrain the **data structure**, not to define an import implementation.

#### Source A: Session Log CSV (“training_menu_cycle_01”)

Observed structure:

- **Week**: e.g., “Week1 (RPE6)”
- **Day**: “Day N (YYYY/MM/DD)” and sometimes a time range (e.g., “17:00-”)
- **Rows**: each row represents a performed entry with an entry type label (e.g., Top/Main/Backoff/Accessory) and an exercise name.
  - Note: Because the system requires an exercise catalog identifier, each row’s exercise name must be mappable to a catalog entry for structured storage.

Observed columns (CSV header):

- `種目` (entry_type): e.g., トップ / メイン / バックオフ / 補助
- `RPE`: numeric 1–10 (session or entry perceived effort)
- `重量`: load (commonly in kg; sometimes notes indicate plate-only, lbs, or special cases)
- `REP`: repetitions per set
- `SET`: number of sets (count)
- `メモ`: freeform notes

Additional observed realities that influence the data model:

- Exercise naming can include variations and modifiers (e.g., “ワイドDL（ポーズ）”, “ボトムバウンドBP”).
- Some rows may encode exceptional values in notes (e.g., “55lbs”, “重りのみ”, “スキップ”).
- A day is a logical container for multiple entries; a workout session should preserve entry ordering.
- Explicit timing (session start/end, per-entry start/end) is not reliably present in the CSV today; the data structure should allow capturing it as structured fields when available, rather than embedding it in notes.
- Load units may be mixed (e.g., kg and lb); canonical storage is kilograms, so unit conversion is required before storage.

#### Source B: Program Spreadsheet PDF (“Example Program”)

The document includes both **training program prescriptions** and additional non-training information (e.g., nutrition guidance, profile). For OPTel Workout’s “dumb backend” scope:

- **In-scope for Program**: week/day structure, per-exercise prescription fields.
- **Out-of-scope as domain entities (for now)**: nutrition and profile calculations. If retained, they should be captured as generic telemetry metrics rather than first-class “health” logic.

Observed program/prescription columns (examples):

- `Week` / `Day`
- `部位` (body part grouping)
- `種目` (exercise)
- `RPE` (target)
- `％RM` (percent-of-1RM guidance)
- `Sets`, `Reps`, `TotalReps`
- `インターバル` (recommended rest interval)
- `参考重量` (reference weight)
- `実際の重量` (actual weight used)
- `体感RPE` (perceived RPE)
- `メモ`

The data structure must be able to represent the program prescription fields without requiring per-program bespoke columns.

### Key Entities *(include if feature involves data)*

- **Workout**: A completed training session, primarily a container for ordered performed entries. Key attributes include an identifier, a session timestamp, optional session metadata (body_weight_kg, fatigue_level 1-5, sleep_hours, condition_notes), optional session start/end times, optional `program_node_id` with `program_context` (hierarchy path snapshot). Editable after recording.
- **WorkoutEntry**: A performed entry inside a workout session, aligned with a single CSV row: `entry_type` (top/main/backoff/accessory), **required exercise catalog identifier** (plus optional display name), performed sets/reps/load (kg), perceived RPE, optional per-entry start/end, **entry-level planned/performed rest with optional per-set overrides**, optional link to a planned prescription (ProgramNode), optional plan snapshot (planned sets/reps/load/RPE/rest), and optional notes. Preserves ordering within the session.
- **Program**: A recursive tree of **ProgramNodes** with arbitrary depth and user-defined node types (e.g., "cycle", "week", "day", "block"). Leaf nodes represent exercise prescriptions with targets. Provides context and intent for planned training.
- **ProgramNode**: A node in the Program tree. Attributes: id, name, node_type (user-defined string), order, children (recursive). Leaf nodes (exercise prescriptions) include: exercise_id, target_sets, target_reps, target_rpe, percent_1rm, planned_rest, muscle_groups.
- **TelemetryPoint**: A time-series metric data point with timestamp, name, value, and unit; optionally linked to a workout for contextual analysis.
- **Exercise**: A catalog entry representing a canonical exercise. Required attributes: unique identifier and name. Extensible to support future attributes (e.g., description, aliases, muscle_groups, load_increment for per-exercise input precision) without breaking existing references.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A new integrator can create a valid **Workout** record (including at least one structured **WorkoutEntry**) in under 15 minutes using only the published specification.
- **SC-002**: For a sample of 30 performed entries linked to a program prescription, 100% can be represented with structured planned vs performed values (sets/reps/load/RPE/rest) without relying solely on freeform notes.
- **SC-003**: For at least one recorded workout session, session start/end and per-entry start/end timestamps can be stored and retrieved as structured fields when provided.
- **SC-004**: At least 10 distinct telemetry metric types can be recorded and retrieved using the single **TelemetryPoint** structure, without schema changes.
- **SC-005**: Users can edit a recorded workout and verify that the updated values are persisted correctly.
