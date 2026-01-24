# OpenAPI Schema and Domain Model Consistency Check

**Date**: 2026-01-24  
**Status**: Manual Verification (T033)

## Verification Method

Since Docker is not available, this is a manual code review comparing:
- Domain models in `api/internal/domain/*.go`
- OpenAPI schemas in `api/openapi/openapi.yaml`

## Exercise Entity ✅

### Domain Model
```go
type Exercise struct {
    ID            uuid.UUID `json:"id"`
    Name          string    `json:"name"`
    Description   *string   `json:"description,omitempty"`
    Aliases       []string  `json:"aliases,omitempty"`
    MuscleGroups  []string  `json:"muscle_groups,omitempty"`
    LoadIncrement *float64  `json:"load_increment,omitempty"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

### OpenAPI Schema
- ✅ Required: `id`, `name`, `created_at`, `updated_at`
- ✅ Optional: `description`, `aliases`, `muscle_groups`, `load_increment`
- ✅ Types match: UUID → string format uuid, string → string, float64 → number double
- ✅ Constraints: `name` minLength 1, maxLength 200; `load_increment` minimum 0 exclusive

**Status**: ✅ Consistent

## Workout Entity ✅

### Domain Model
```go
type Workout struct {
    ID             uuid.UUID      `json:"id"`
    Timestamp      time.Time      `json:"timestamp"`
    SessionStart   *time.Time     `json:"session_start,omitempty"`
    SessionEnd     *time.Time     `json:"session_end,omitempty"`
    BodyWeightKg   *float64       `json:"body_weight_kg,omitempty"`
    FatigueLevel   *int           `json:"fatigue_level,omitempty"`
    SleepHours     *float64       `json:"sleep_hours,omitempty"`
    ConditionNotes *string         `json:"condition_notes,omitempty"`
    ProgramNodeID  *uuid.UUID     `json:"program_node_id,omitempty"`
    ProgramContext []string       `json:"program_context,omitempty"`
    Notes          *string        `json:"notes,omitempty"`
    CreatedAt      time.Time      `json:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
    Entries        []WorkoutEntry `json:"entries"`
}
```

### OpenAPI Schema
- ✅ Required: `id`, `timestamp`, `entries`, `created_at`, `updated_at`
- ✅ Optional: All other fields match
- ✅ Types match: time.Time → date-time, float64 → number double, int → integer
- ✅ Constraints: `fatigue_level` min 1 max 5, `sleep_hours` min 0 max 24
- ✅ `entries` references `WorkoutEntry` schema

**Status**: ✅ Consistent

## WorkoutEntry Entity ✅

### Domain Model
```go
type WorkoutEntry struct {
    ID                    uuid.UUID     `json:"id"`
    WorkoutID             uuid.UUID     `json:"workout_id"`
    Order                 int           `json:"order"`
    ExerciseID            uuid.UUID     `json:"exercise_id"`
    DisplayName           *string       `json:"display_name,omitempty"`
    EntryType             string        `json:"entry_type"`
    Sets                  int           `json:"sets"`
    Reps                  int           `json:"reps"`
    LoadKg                float64       `json:"load_kg"`
    RPE                   int           `json:"rpe"`
    EntryStart            *time.Time    `json:"entry_start,omitempty"`
    EntryEnd              *time.Time    `json:"entry_end,omitempty"`
    PlannedRestSeconds    *int          `json:"planned_rest_seconds,omitempty"`
    PerformedRestSeconds  *int          `json:"performed_rest_seconds,omitempty"`
    PerSetRestOverrides   []int         `json:"per_set_rest_overrides,omitempty"`
    ProgramNodeID         *uuid.UUID    `json:"program_node_id,omitempty"`
    PlanSnapshot          *PlanSnapshot `json:"plan_snapshot,omitempty"`
    Notes                 *string       `json:"notes,omitempty"`
}
```

### OpenAPI Schema
- ✅ Required: `id`, `workout_id`, `order`, `exercise_id`, `entry_type`, `sets`, `reps`, `load_kg`, `rpe`
- ✅ Optional: All other fields match
- ✅ Types match
- ✅ Constraints: `entry_type` enum [top, main, backoff, accessory], `sets`/`reps` min 1, `load_kg` min 0, `rpe` min 1 max 10
- ✅ `plan_snapshot` references `PlanSnapshot` schema

**Status**: ✅ Consistent

## PlanSnapshot Entity ✅

### Domain Model
```go
type PlanSnapshot struct {
    ProgramNodeID      *uuid.UUID `json:"program_node_id,omitempty"`
    TargetSets         *int       `json:"target_sets,omitempty"`
    TargetReps         *int       `json:"target_reps,omitempty"`
    TargetRPE          *int       `json:"target_rpe,omitempty"`
    TargetLoadKg       *float64   `json:"target_load_kg,omitempty"`
    Percent1RM         *float64   `json:"percent_1rm,omitempty"`
    PlannedRestSeconds *int       `json:"planned_rest_seconds,omitempty"`
}
```

### OpenAPI Schema
- ✅ All fields optional (no required section)
- ✅ Types match
- ✅ Constraints: `target_sets`/`target_reps` min 1, `target_rpe` min 1 max 10, `percent_1rm` min 0 max 1

**Status**: ✅ Consistent

## Program Entity ✅

### Domain Model
```go
type Program struct {
    ID          uuid.UUID     `json:"id"`
    Name        string        `json:"name"`
    Description *string       `json:"description,omitempty"`
    CreatedAt   time.Time     `json:"created_at"`
    UpdatedAt   time.Time     `json:"updated_at"`
    RootNodes   []ProgramNode `json:"root_nodes,omitempty"`
}
```

### OpenAPI Schema
- ✅ Required: `id`, `name`, `created_at`, `updated_at`
- ✅ Optional: `description`, `root_nodes`
- ✅ Types match
- ✅ `root_nodes` references `ProgramNode` schema (recursive)

**Status**: ✅ Consistent

## ProgramNode Entity ✅

### Domain Model
```go
type ProgramNode struct {
    ID                uuid.UUID     `json:"id"`
    ProgramID         uuid.UUID     `json:"program_id"`
    ParentID          *uuid.UUID    `json:"parent_id,omitempty"`
    Name              string        `json:"name"`
    NodeType          string        `json:"node_type"`
    Order             int           `json:"order"`
    Children          []ProgramNode `json:"children,omitempty"`
    ExerciseID         *uuid.UUID `json:"exercise_id,omitempty"`
    TargetSets         *int       `json:"target_sets,omitempty"`
    TargetReps         *int       `json:"target_reps,omitempty"`
    TargetRPE          *int       `json:"target_rpe,omitempty"`
    Percent1RM         *float64   `json:"percent_1rm,omitempty"`
    PlannedRestSeconds *int       `json:"planned_rest_seconds,omitempty"`
    MuscleGroups       []string   `json:"muscle_groups,omitempty"`
    Notes              *string    `json:"notes,omitempty"`
}
```

### OpenAPI Schema
- ✅ Required: `id`, `program_id`, `name`, `node_type`, `order`
- ✅ Optional: All other fields match
- ✅ Types match
- ✅ Constraints: `name` minLength 1 maxLength 200, `node_type` minLength 1 maxLength 50
- ✅ `children` references `ProgramNode` schema (recursive)

**Status**: ✅ Consistent

## TelemetryPoint Entity ✅

### Domain Model
```go
type TelemetryPoint struct {
    ID         uuid.UUID  `json:"id"`
    Timestamp  time.Time  `json:"timestamp"`
    MetricName string     `json:"metric_name"`
    Value      float64    `json:"value"`
    Unit       string     `json:"unit"`
    WorkoutID  *uuid.UUID `json:"workout_id,omitempty"`
    CreatedAt  time.Time  `json:"created_at"`
}
```

### OpenAPI Schema
- ✅ Required: `id`, `timestamp`, `metric_name`, `value`, `unit`, `created_at`
- ✅ Optional: `workout_id`
- ✅ Types match
- ✅ Constraints: `metric_name` minLength 1 maxLength 100, `unit` minLength 1 maxLength 50

**Status**: ✅ Consistent

## Create Schemas ✅

All `*Create` schemas (ExerciseCreate, WorkoutCreate, WorkoutEntryCreate, ProgramCreate, ProgramNodeCreate, TelemetryPointCreate) have:
- ✅ Appropriate required fields (excluding auto-generated: id, created_at, updated_at)
- ✅ Optional fields match domain models
- ✅ Constraints match domain validation rules

**Status**: ✅ Consistent

## Summary

| Entity | Status | Notes |
|--------|--------|-------|
| Exercise | ✅ | All fields and constraints match |
| Workout | ✅ | All fields and constraints match |
| WorkoutEntry | ✅ | All fields and constraints match |
| PlanSnapshot | ✅ | All fields and constraints match |
| Program | ✅ | All fields and constraints match |
| ProgramNode | ✅ | All fields and constraints match |
| TelemetryPoint | ✅ | All fields and constraints match |
| Create Schemas | ✅ | All create schemas consistent |

**Overall Status**: ✅ **CONSISTENT**

All domain models and OpenAPI schemas are consistent. Field names, types, required/optional status, and constraints all match.

**Note**: Full code generation verification (T031, T032) requires Docker. Manual verification confirms consistency.
