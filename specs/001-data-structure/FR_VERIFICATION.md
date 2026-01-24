# Functional Requirements Verification

**Date**: 2026-01-24  
**Status**: Verification Report

## FR-001: Canonical Entities ✅

**Requirement**: The system MUST define canonical entities: Workout, WorkoutEntry, Program, ProgramNode, TelemetryPoint, and Exercise.

**Verification**:
- ✅ Exercise: `api/internal/domain/exercise.go`
- ✅ Workout: `api/internal/domain/workout.go`
- ✅ WorkoutEntry: `api/internal/domain/workout.go`
- ✅ Program: `api/internal/domain/program.go`
- ✅ ProgramNode: `api/internal/domain/program.go`
- ✅ TelemetryPoint: `api/internal/domain/telemetry.go`

## FR-001a: Exercise Entity ✅

**Requirement**: Exercise MUST include unique identifier and name. Structure MUST be extensible.

**Verification**:
- ✅ `ID uuid.UUID` - unique identifier
- ✅ `Name string` - required field
- ✅ Extensible: `Description`, `Aliases`, `MuscleGroups`, `LoadIncrement` are optional pointers/slices

## FR-002: Workout Basic Fields ✅

**Requirement**: Workout MUST include unique identifier and timestamp.

**Verification**:
- ✅ `ID uuid.UUID` - unique identifier
- ✅ `Timestamp time.Time` - required field

## FR-002a: Workout Condition Metadata ✅

**Requirement**: Workout MAY include body weight, fatigue level, sleep hours, condition notes.

**Verification**:
- ✅ `BodyWeightKg *float64` - optional
- ✅ `FatigueLevel *int` - optional (1-5 scale validated)
- ✅ `SleepHours *float64` - optional (0-24 range validated)
- ✅ `ConditionNotes *string` - optional (max 2000 chars validated)

## FR-002b: Workout Program Context ✅

**Requirement**: Workout MAY include program_node_id and program_context snapshot.

**Verification**:
- ✅ `ProgramNodeID *uuid.UUID` - optional reference
- ✅ `ProgramContext []string` - optional hierarchy path snapshot

## FR-003: Workout Entries ✅

**Requirement**: Workout MUST contain one or more WorkoutEntries.

**Verification**:
- ✅ `Entries []WorkoutEntry` - required field
- ✅ Validation: `len(w.Entries) == 0` returns error

## FR-004: WorkoutEntry Required Fields ✅

**Requirement**: WorkoutEntry MUST include exercise_id, entry_type, RPE, sets, reps, load.

**Verification**:
- ✅ `ExerciseID uuid.UUID` - required
- ✅ `EntryType string` - required (enum: top/main/backoff/accessory)
- ✅ `RPE int` - required (1-10 validated)
- ✅ `Sets int` - required (> 0 validated)
- ✅ `Reps int` - required (> 0 validated)
- ✅ `LoadKg float64` - required (≥ 0 validated, rounded to 0.1kg)
- ✅ `DisplayName *string` - optional

## FR-005: Load Precision ✅

**Requirement**: Load values MUST be stored in kilograms with 0.1kg precision.

**Verification**:
- ✅ `LoadKg float64` - stored as float64
- ✅ `RoundLoad(kg float64) float64` - rounds to 0.1kg precision
- ✅ Validation applies rounding: `e.LoadKg = RoundLoad(e.LoadKg)`

## FR-006: Timing Support ✅

**Requirement**: Support session start/end and optional per-entry start/end timestamps.

**Verification**:
- ✅ `SessionStart *time.Time` - optional
- ✅ `SessionEnd *time.Time` - optional
- ✅ `EntryStart *time.Time` - optional on WorkoutEntry
- ✅ `EntryEnd *time.Time` - optional on WorkoutEntry
- ✅ Validation: `session_start ≤ session_end` if both provided

## FR-007: Rest/Interval Information ✅

**Requirement**: Support entry-level planned/performed rest with optional per-set overrides.

**Verification**:
- ✅ `PlannedRestSeconds *int` - optional entry-level
- ✅ `PerformedRestSeconds *int` - optional entry-level
- ✅ `PerSetRestOverrides []int` - optional per-set overrides

## FR-008: Program Linkage ✅

**Requirement**: Link performed entry to planned prescription and record planned vs performed values.

**Verification**:
- ✅ `ProgramNodeID *uuid.UUID` - optional on WorkoutEntry
- ✅ `PlanSnapshot *PlanSnapshot` - optional embedded struct with planned values

## FR-009: Plan Snapshot ✅

**Requirement**: Store snapshot of planned values at execution time on performed entry.

**Verification**:
- ✅ `PlanSnapshot` struct with fields:
  - `ProgramNodeID *uuid.UUID`
  - `TargetSets *int`
  - `TargetReps *int`
  - `TargetRPE *int`
  - `TargetLoadKg *float64`
  - `Percent1RM *float64`
  - `PlannedRestSeconds *int`
- ✅ Embedded in `WorkoutEntry` as optional field

## FR-010: Editable Workouts ✅

**Requirement**: Workout logs MAY be edited after recording.

**Verification**:
- ✅ Workout struct has `UpdatedAt time.Time` field
- ✅ No immutability constraints in validation
- ✅ Documented in plan.md as justified violation of Constitution

## FR-011: Validation Rules ✅

**Requirement**: Reject invalid records (future timestamp, missing entries, invalid values).

**Verification**:
- ✅ `ValidateTimestamp` - rejects future timestamps
- ✅ `ValidateWorkout` - requires at least one entry
- ✅ `ValidateWorkoutEntry` - validates:
  - `sets > 0`
  - `reps > 0`
  - `load_kg ≥ 0`
  - `rpe` in 1-10 range
  - `entry_type` enum validation

## FR-012: Recursive Program Structure ✅

**Requirement**: Program MUST be representable as recursive tree of ProgramNodes.

**Verification**:
- ✅ `Program.RootNodes []ProgramNode` - root nodes
- ✅ `ProgramNode.Children []ProgramNode` - recursive children
- ✅ `ProgramNode.ParentID *uuid.UUID` - parent reference (nil for root)
- ✅ `ProgramNode.NodeType string` - user-defined type
- ✅ `ProgramNode.Order int` - order among siblings
- ✅ Validation recursively validates children

## FR-012a: Editable Programs ✅

**Requirement**: Programs MAY be edited after creation.

**Verification**:
- ✅ Program struct has `UpdatedAt time.Time` field
- ✅ No immutability constraints in validation
- ✅ Plan snapshots on WorkoutEntry preserve historical values

## FR-013: ProgramNode Prescription Fields ✅

**Requirement**: ProgramNode (leaf) MAY include planned targets.

**Verification**:
- ✅ `ExerciseID *uuid.UUID` - optional
- ✅ `TargetSets *int` - optional (> 0 validated)
- ✅ `TargetReps *int` - optional (> 0 validated)
- ✅ `TargetRPE *int` - optional (1-10 validated)
- ✅ `Percent1RM *float64` - optional (0.0-1.0 validated)
- ✅ `PlannedRestSeconds *int` - optional (≥ 0 validated)
- ✅ `MuscleGroups []string` - optional

## FR-014: TelemetryPoint Required Fields ✅

**Requirement**: TelemetryPoint MUST include timestamp, metric name, value, unit. MAY include workout reference.

**Verification**:
- ✅ `Timestamp time.Time` - required
- ✅ `MetricName string` - required (1-100 chars validated)
- ✅ `Value float64` - required
- ✅ `Unit string` - required (1-50 chars validated)
- ✅ `WorkoutID *uuid.UUID` - optional

## FR-015: Opaque Metric Names ✅

**Requirement**: Telemetry metric names MUST be treated as opaque identifiers.

**Verification**:
- ✅ `MetricName string` - no validation of specific values
- ✅ No health-scoring or wellness logic in domain layer
- ✅ New metric names can be introduced without schema changes

## FR-016: Documentation ✅

**Requirement**: Data structures MUST be documented with required/optional attributes and validation rules.

**Verification**:
- ✅ `api/internal/domain/README.md` - comprehensive documentation
- ✅ `specs/001-data-structure/data-model.md` - detailed field definitions
- ✅ Validation functions documented with error messages
- ✅ Test files demonstrate validation rules

## FR-017: Data Structure Versioning ✅

**Requirement**: Changes MUST preserve interpretability of previously recorded data.

**Verification**:
- ✅ Documented in plan.md "Future Considerations" section
- ✅ Initial implementation: additive-only changes
- ✅ Future strategy: semantic versioning (MAJOR/MINOR/PATCH)
- ✅ Migration approach documented for breaking changes

---

## Edge Cases Coverage ✅

All edge cases from spec.md are covered in tests:

1. ✅ Load = 0 (bodyweight) - `workout_test.go`: "load_kg zero (bodyweight allowed)"
2. ✅ Future timestamp - `workout_test.go`: "timestamp in future"
3. ✅ Invalid RPE - `workout_test.go`: "RPE too low", "RPE too high"
4. ✅ Missing optional fields - All test files: "valid ... with minimal fields"
5. ✅ Session start without end - `workout_test.go`: "valid workout with all fields"
6. ✅ Entry timing optional - `workout_test.go`: entries with/without timing
7. ✅ Rest/interval variations - `workout_test.go`: planned/performed rest tests
8. ✅ Telemetry without workout - `telemetry_test.go`: "valid telemetry point with minimal fields"

---

## Summary

**Total Requirements**: 17 (FR-001 to FR-017)  
**Verified**: 17/17 ✅  
**Edge Cases Covered**: 8/8 ✅

**Status**: All functional requirements are met. Implementation is complete and verified.
