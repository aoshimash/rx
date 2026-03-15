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

// ValidatePlanEntry validates a PlanEntry entity.
// Metadata contents are not validated (free-form JSON).
func ValidatePlanEntry(e *PlanEntry) error {
	if e == nil {
		return &ValidationError{
			Field:   "plan_entry",
			Message: "plan_entry cannot be nil",
		}
	}

	// Validate exercise_name
	if err := ValidateRequiredString("exercise_name", e.ExerciseName); err != nil {
		return err
	}
	if err := ValidateStringLength("exercise_name", e.ExerciseName, 1, 200); err != nil {
		return err
	}

	// Validate order
	if e.Order < 0 {
		return &ValidationError{
			Field:   "order",
			Message: "order must be greater than or equal to 0",
		}
	}

	// Validate optional fields
	if e.Sets != nil && *e.Sets <= 0 {
		return &ValidationError{
			Field:   "sets",
			Message: "sets must be greater than 0",
		}
	}

	if e.Reps != nil && *e.Reps <= 0 {
		return &ValidationError{
			Field:   "reps",
			Message: "reps must be greater than 0",
		}
	}

	if e.LoadKg != nil {
		if *e.LoadKg < 0 {
			return &ValidationError{
				Field:   "load_kg",
				Message: "load_kg must be greater than or equal to 0",
			}
		}
		*e.LoadKg = RoundLoad(*e.LoadKg)
	}

	if e.RPE != nil {
		if err := ValidateRPE(*e.RPE); err != nil {
			return err
		}
	}

	if e.Notes != nil {
		if err := ValidateStringLength("notes", *e.Notes, 0, 2000); err != nil {
			return err
		}
	}

	return nil
}

// ValidatePlan validates a Plan entity.
func ValidatePlan(p *Plan) error {
	if p == nil {
		return &ValidationError{
			Field:   "plan",
			Message: "plan cannot be nil",
		}
	}

	// Validate required fields
	if err := ValidateRequiredString("name", p.Name); err != nil {
		return err
	}
	if err := ValidateStringLength("name", p.Name, 1, 200); err != nil {
		return err
	}

	// Validate description length if provided
	if p.Description != nil {
		if err := ValidateStringLength("description", *p.Description, 0, 2000); err != nil {
			return err
		}
	}

	// Validate notes length if provided
	if p.Notes != nil {
		if err := ValidateStringLength("notes", *p.Notes, 0, 5000); err != nil {
			return err
		}
	}

	// Validate entries (max 1000)
	if len(p.Entries) > 1000 {
		return &ValidationError{
			Field:   "entries",
			Message: "plan cannot have more than 1000 entries",
		}
	}

	for i := range p.Entries {
		if err := ValidatePlanEntry(&p.Entries[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("entries[%d]", i),
				Message: err.Error(),
			}
		}
	}

	return nil
}

// ValidateLogEntry validates a LogEntry entity.
func ValidateLogEntry(e *LogEntry) error {
	if e == nil {
		return &ValidationError{
			Field:   "log_entry",
			Message: "log_entry cannot be nil",
		}
	}

	// Validate exercise_name
	if err := ValidateRequiredString("exercise_name", e.ExerciseName); err != nil {
		return err
	}
	if err := ValidateStringLength("exercise_name", e.ExerciseName, 1, 200); err != nil {
		return err
	}

	// Validate order
	if e.Order < 0 {
		return &ValidationError{
			Field:   "order",
			Message: "order must be greater than or equal to 0",
		}
	}

	// Validate optional fields
	if e.Sets != nil && *e.Sets <= 0 {
		return &ValidationError{
			Field:   "sets",
			Message: "sets must be greater than 0",
		}
	}

	if e.Reps != nil && *e.Reps <= 0 {
		return &ValidationError{
			Field:   "reps",
			Message: "reps must be greater than 0",
		}
	}

	if e.LoadKg != nil {
		if *e.LoadKg < 0 {
			return &ValidationError{
				Field:   "load_kg",
				Message: "load_kg must be greater than or equal to 0",
			}
		}
		*e.LoadKg = RoundLoad(*e.LoadKg)
	}

	if e.RPE != nil {
		if err := ValidateRPE(*e.RPE); err != nil {
			return err
		}
	}

	if e.Notes != nil {
		if err := ValidateStringLength("notes", *e.Notes, 0, 2000); err != nil {
			return err
		}
	}

	if e.VideoObjectKey != nil {
		if err := ValidateStringLength("video_object_key", *e.VideoObjectKey, 1, 500); err != nil {
			return err
		}
	}

	return nil
}

// ValidateLog validates a Log entity.
func ValidateLog(l *Log) error {
	if l == nil {
		return &ValidationError{
			Field:   "log",
			Message: "log cannot be nil",
		}
	}

	// Validate performed_at (not in future)
	if err := ValidateTimestamp(l.PerformedAt); err != nil {
		return err
	}

	// Validate entries (must have at least one, max 500)
	if len(l.Entries) == 0 {
		return &ValidationError{
			Field:   "entries",
			Message: "log must have at least one entry",
		}
	}
	if len(l.Entries) > 500 {
		return &ValidationError{
			Field:   "entries",
			Message: "log cannot have more than 500 entries",
		}
	}

	// Validate each entry
	for i := range l.Entries {
		if err := ValidateLogEntry(&l.Entries[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("entries[%d]", i),
				Message: err.Error(),
			}
		}
	}

	// Validate notes length if provided
	if l.Notes != nil {
		if err := ValidateStringLength("notes", *l.Notes, 0, 5000); err != nil {
			return err
		}
	}

	return nil
}
