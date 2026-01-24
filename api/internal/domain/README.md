# Domain Models

This package contains the core domain entities for OPTel Training.

## Entities

### Exercise
A catalog entry representing a canonical exercise.

**Required fields**: `id`, `name`, `created_at`, `updated_at`  
**Optional fields**: `description`, `aliases`, `muscle_groups`, `load_increment`

### Workout
A completed training session containing performed entries.

**Required fields**: `id`, `timestamp`, `entries`, `created_at`, `updated_at`  
**Optional fields**: `session_start`, `session_end`, `body_weight_kg`, `fatigue_level`, `sleep_hours`, `condition_notes`, `program_node_id`, `program_context`, `notes`

### WorkoutEntry
A single performed exercise entry within a workout session.

**Required fields**: `id`, `workout_id`, `order`, `exercise_id`, `entry_type`, `sets`, `reps`, `load_kg`, `rpe`  
**Optional fields**: `display_name`, `entry_start`, `entry_end`, `planned_rest_seconds`, `performed_rest_seconds`, `per_set_rest_overrides`, `program_node_id`, `plan_snapshot`, `notes`

### PlanSnapshot
A snapshot of planned values at execution time (embedded in WorkoutEntry).

**All fields optional**: `program_node_id`, `target_sets`, `target_reps`, `target_rpe`, `target_load_kg`, `percent_1rm`, `planned_rest_seconds`

### Program
A training program containing a recursive tree of nodes.

**Required fields**: `id`, `name`, `created_at`, `updated_at`  
**Optional fields**: `description`, `root_nodes`

### ProgramNode
A node in the program tree (cycle, week, day, block, or exercise prescription).

**Required fields**: `id`, `program_id`, `name`, `node_type`, `order`  
**Optional fields**: `parent_id`, `children`, `exercise_id`, `target_sets`, `target_reps`, `target_rpe`, `percent_1rm`, `planned_rest_seconds`, `muscle_groups`, `notes`

### TelemetryPoint
A time-series metric data point.

**Required fields**: `id`, `timestamp`, `metric_name`, `value`, `unit`, `created_at`  
**Optional fields**: `workout_id`

## Validation

All entities have corresponding validation functions:

- `ValidateExercise(e *Exercise) error`
- `ValidateWorkout(w *Workout) error`
- `ValidateWorkoutEntry(e *WorkoutEntry) error`
- `ValidateProgram(p *Program) error`
- `ValidateProgramNode(n *ProgramNode) error`
- `ValidateTelemetryPoint(t *TelemetryPoint) error`

### Common Validation Helpers

- `RoundLoad(kg float64) float64` - Round to 0.1kg precision
- `ValidateRPE(rpe int) error` - Check 1-10 range
- `ValidateFatigueLevel(level int) error` - Check 1-5 range
- `ValidateTimestamp(t time.Time) error` - Check not in future
- `ValidateEntryType(entryType string) error` - Check enum values

## Error Types

- `ValidationError` - Field-level validation errors
- `DomainError` - Domain-level errors with error codes

See `errors.go` for error code constants.

## Testing

All validation functions have comprehensive table-driven tests in `*_test.go` files.

Run tests:
```bash
go test ./internal/domain/...
```

## References

- [Data Model Specification](../../../specs/001-data-structure/data-model.md)
- [Feature Specification](../../../specs/001-data-structure/spec.md)
- [Implementation Plan](../../../specs/001-data-structure/plan.md)
