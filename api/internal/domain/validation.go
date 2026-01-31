package domain

import (
	"fmt"
	"math"
	"time"
)

// RoundLoad rounds a weight value to 0.1kg precision.
func RoundLoad(kg float64) float64 {
	return math.Round(kg*10) / 10
}

// ValidateRPE checks if RPE is in the valid range (1-10).
func ValidateRPE(rpe int) error {
	if rpe < 1 || rpe > 10 {
		return &DomainError{
			Code:    ErrCodeInvalidRPE,
			Message: "RPE must be between 1 and 10",
			Details: map[string]interface{}{
				"value": rpe,
				"min":   1,
				"max":   10,
			},
		}
	}
	return nil
}

// ValidateFatigueLevel checks if fatigue level is in the valid range (1-5).
func ValidateFatigueLevel(level int) error {
	if level < 1 || level > 5 {
		return &DomainError{
			Code:    ErrCodeInvalidFatigueLevel,
			Message: "Fatigue level must be between 1 and 5",
			Details: map[string]interface{}{
				"value": level,
				"min":   1,
				"max":   5,
			},
		}
	}
	return nil
}

// ValidateTimestamp checks if timestamp is not in the future.
func ValidateTimestamp(t time.Time) error {
	now := time.Now()
	if t.After(now) {
		return &DomainError{
			Code:    ErrCodeInvalidTimestamp,
			Message: "Timestamp cannot be in the future",
			Details: map[string]interface{}{
				"timestamp": t,
				"now":       now,
			},
		}
	}
	return nil
}

// ValidateEntryType checks if entry_type is valid (nullable, max 50 chars).
// User-defined values are allowed.
func ValidateEntryType(entryType *string) error {
	// entry_type is nullable, so nil is valid
	if entryType == nil {
		return nil
	}

	// Validate max length (50 characters)
	if len(*entryType) > 50 {
		return &DomainError{
			Code:    ErrCodeInvalidEntryType,
			Message: "Entry type must be at most 50 characters",
			Details: map[string]interface{}{
				"value":     *entryType,
				"maxLength": 50,
			},
		}
	}
	return nil
}

// ValidateRequiredString checks if a required string field is not empty.
func ValidateRequiredString(field, value string) error {
	if value == "" {
		return &ValidationError{
			Field:   field,
			Message: "required field cannot be empty",
		}
	}
	return nil
}

// ValidateStringLength checks if a string is within the specified length range.
func ValidateStringLength(field, value string, min, max int) error {
	length := len(value)
	if length < min || length > max {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("string length must be between %d and %d characters", min, max),
		}
	}
	return nil
}

// ValidateExercise validates an Exercise entity.
func ValidateExercise(e *Exercise) error {
	if e == nil {
		return &ValidationError{
			Field:   "exercise",
			Message: "exercise cannot be nil",
		}
	}

	// Validate required fields
	if err := ValidateRequiredString("name", e.Name); err != nil {
		return err
	}

	// Validate name length
	if err := ValidateStringLength("name", e.Name, 1, 200); err != nil {
		return err
	}

	// Validate description length if provided
	if e.Description != nil {
		if err := ValidateStringLength("description", *e.Description, 0, 2000); err != nil {
			return err
		}
	}

	// Validate load_increment if provided
	if e.LoadIncrement != nil {
		if *e.LoadIncrement <= 0 {
			return &ValidationError{
				Field:   "load_increment",
				Message: "load_increment must be greater than 0",
			}
		}
	}

	return nil
}

// ValidateWorkoutEntry validates a WorkoutEntry entity.
func ValidateWorkoutEntry(e *WorkoutEntry) error {
	if e == nil {
		return &ValidationError{
			Field:   "workout_entry",
			Message: "workout_entry cannot be nil",
		}
	}

	// Validate entry_type (nullable, max 50 chars)
	if err := ValidateEntryType(e.EntryType); err != nil {
		return err
	}

	// Validate sets
	if e.Sets <= 0 {
		return &ValidationError{
			Field:   "sets",
			Message: "sets must be greater than 0",
		}
	}

	// Validate reps
	if e.Reps <= 0 {
		return &ValidationError{
			Field:   "reps",
			Message: "reps must be greater than 0",
		}
	}

	// Validate load_kg (≥ 0, rounded to 0.1kg)
	if e.LoadKg < 0 {
		return &ValidationError{
			Field:   "load_kg",
			Message: "load_kg must be greater than or equal to 0",
		}
	}
	e.LoadKg = RoundLoad(e.LoadKg)

	// Validate RPE
	if err := ValidateRPE(e.RPE); err != nil {
		return err
	}

	// Validate order
	if e.Order < 0 {
		return &ValidationError{
			Field:   "order",
			Message: "order must be greater than or equal to 0",
		}
	}

	// Validate display_name length if provided
	if e.DisplayName != nil {
		if err := ValidateStringLength("display_name", *e.DisplayName, 0, 200); err != nil {
			return err
		}
	}

	// Validate notes length if provided
	if e.Notes != nil {
		if err := ValidateStringLength("notes", *e.Notes, 0, 2000); err != nil {
			return err
		}
	}

	// Validate video_object_key length if provided
	if e.VideoObjectKey != nil {
		if err := ValidateStringLength("video_object_key", *e.VideoObjectKey, 1, 500); err != nil {
			return err
		}
	}

	// Validate rest seconds if provided
	if e.PlannedRestSeconds != nil && *e.PlannedRestSeconds < 0 {
		return &ValidationError{
			Field:   "planned_rest_seconds",
			Message: "planned_rest_seconds must be greater than or equal to 0",
		}
	}
	if e.PerformedRestSeconds != nil && *e.PerformedRestSeconds < 0 {
		return &ValidationError{
			Field:   "performed_rest_seconds",
			Message: "performed_rest_seconds must be greater than or equal to 0",
		}
	}

	return nil
}

// ValidateWorkout validates a Workout entity.
func ValidateWorkout(w *Workout) error {
	if w == nil {
		return &ValidationError{
			Field:   "workout",
			Message: "workout cannot be nil",
		}
	}

	// Validate timestamp (not in future)
	if err := ValidateTimestamp(w.Timestamp); err != nil {
		return err
	}

	// Validate entries (must have at least one, max 500 per FR-017, FR-018)
	if len(w.Entries) == 0 {
		return &ValidationError{
			Field:   "entries",
			Message: "workout must have at least one entry",
		}
	}
	if len(w.Entries) > 500 {
		return &ValidationError{
			Field:   "entries",
			Message: "workout cannot have more than 500 entries",
		}
	}

	// Validate each entry
	// Use index-based loop to allow RoundLoad to modify original entries
	for i := range w.Entries {
		if err := ValidateWorkoutEntry(&w.Entries[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("entries[%d]", i),
				Message: err.Error(),
			}
		}
	}

	// Validate session times
	if w.SessionStart != nil && w.SessionEnd != nil {
		if w.SessionStart.After(*w.SessionEnd) {
			return &ValidationError{
				Field:   "session_start",
				Message: "session_start must be before or equal to session_end",
			}
		}
	}

	// Validate fatigue_level if provided
	if w.FatigueLevel != nil {
		if err := ValidateFatigueLevel(*w.FatigueLevel); err != nil {
			return err
		}
	}

	// Validate body_weight_kg if provided
	if w.BodyWeightKg != nil {
		if *w.BodyWeightKg <= 0 {
			return &ValidationError{
				Field:   "body_weight_kg",
				Message: "body_weight_kg must be greater than 0",
			}
		}
		*w.BodyWeightKg = RoundLoad(*w.BodyWeightKg)
	}

	// Validate sleep_hours if provided
	if w.SleepHours != nil {
		if *w.SleepHours < 0 || *w.SleepHours > 24 {
			return &ValidationError{
				Field:   "sleep_hours",
				Message: "sleep_hours must be between 0 and 24",
			}
		}
	}

	// Validate string lengths
	if w.ConditionNotes != nil {
		if err := ValidateStringLength("condition_notes", *w.ConditionNotes, 0, 2000); err != nil {
			return err
		}
	}
	if w.Notes != nil {
		if err := ValidateStringLength("notes", *w.Notes, 0, 5000); err != nil {
			return err
		}
	}

	return nil
}

// ValidateTelemetryPoint validates a TelemetryPoint entity.
func ValidateTelemetryPoint(t *TelemetryPoint) error {
	if t == nil {
		return &ValidationError{
			Field:   "telemetry_point",
			Message: "telemetry_point cannot be nil",
		}
	}

	// Validate required fields
	if err := ValidateRequiredString("metric_name", t.MetricName); err != nil {
		return err
	}

	// Validate metric_name length
	if err := ValidateStringLength("metric_name", t.MetricName, 1, 100); err != nil {
		return err
	}

	// Validate unit
	if err := ValidateRequiredString("unit", t.Unit); err != nil {
		return err
	}

	// Validate unit length
	if err := ValidateStringLength("unit", t.Unit, 1, 50); err != nil {
		return err
	}

	// Validate timestamp (not in future)
	if err := ValidateTimestamp(t.Timestamp); err != nil {
		return err
	}

	return nil
}

// ValidateProgramNode validates a ProgramNode entity.
func ValidateProgramNode(n *ProgramNode) error {
	if n == nil {
		return &ValidationError{
			Field:   "program_node",
			Message: "program_node cannot be nil",
		}
	}

	// Validate required fields
	if err := ValidateRequiredString("name", n.Name); err != nil {
		return err
	}

	// Validate name length
	if err := ValidateStringLength("name", n.Name, 1, 200); err != nil {
		return err
	}

	// Validate node_type
	if err := ValidateRequiredString("node_type", n.NodeType); err != nil {
		return err
	}

	// Validate order
	if n.Order < 0 {
		return &ValidationError{
			Field:   "order",
			Message: "order must be greater than or equal to 0",
		}
	}

	// Validate prescription fields if provided
	if n.TargetSets != nil && *n.TargetSets <= 0 {
		return &ValidationError{
			Field:   "target_sets",
			Message: "target_sets must be greater than 0",
		}
	}

	if n.TargetReps != nil && *n.TargetReps <= 0 {
		return &ValidationError{
			Field:   "target_reps",
			Message: "target_reps must be greater than 0",
		}
	}

	if n.TargetRPE != nil {
		if err := ValidateRPE(*n.TargetRPE); err != nil {
			return err
		}
	}

	if n.Percent1RM != nil {
		if *n.Percent1RM < 0 || *n.Percent1RM > 1 {
			return &ValidationError{
				Field:   "percent_1rm",
				Message: "percent_1rm must be between 0.0 and 1.0",
			}
		}
	}

	if n.PlannedRestSeconds != nil && *n.PlannedRestSeconds < 0 {
		return &ValidationError{
			Field:   "planned_rest_seconds",
			Message: "planned_rest_seconds must be greater than or equal to 0",
		}
	}

	// Validate notes length if provided
	if n.Notes != nil {
		if err := ValidateStringLength("notes", *n.Notes, 0, 2000); err != nil {
			return err
		}
	}

	// Recursively validate children
	// Use index-based loop to allow validation to modify original children
	for i := range n.Children {
		if err := ValidateProgramNode(&n.Children[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("children[%d]", i),
				Message: err.Error(),
			}
		}
	}

	return nil
}

// ValidateProgram validates a Program entity.
func ValidateProgram(p *Program) error {
	if p == nil {
		return &ValidationError{
			Field:   "program",
			Message: "program cannot be nil",
		}
	}

	// Validate required fields
	if err := ValidateRequiredString("name", p.Name); err != nil {
		return err
	}

	// Validate name length
	if err := ValidateStringLength("name", p.Name, 1, 200); err != nil {
		return err
	}

	// Validate description length if provided
	if p.Description != nil {
		if err := ValidateStringLength("description", *p.Description, 0, 2000); err != nil {
			return err
		}
	}

	// Validate root nodes recursively
	// Use index-based loop to allow validation to modify original nodes
	// Count total nodes (max 1000 per FR-017, FR-018)
	var totalNodeCount int
	var countNodes func(nodes []ProgramNode)
	countNodes = func(nodes []ProgramNode) {
		for i := range nodes {
			totalNodeCount++
			if totalNodeCount > 1000 {
				return
			}
			if len(nodes[i].Children) > 0 {
				countNodes(nodes[i].Children)
			}
		}
	}
	countNodes(p.RootNodes)
	if totalNodeCount > 1000 {
		return &ValidationError{
			Field:   "root_nodes",
			Message: "program cannot have more than 1000 nodes (entire tree)",
		}
	}

	for i := range p.RootNodes {
		if err := ValidateProgramNode(&p.RootNodes[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("root_nodes[%d]", i),
				Message: err.Error(),
			}
		}
	}

	return nil
}
