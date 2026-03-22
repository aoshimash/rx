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

// ValidateTimeRange checks that started_at is before finished_at when both are provided.
func ValidateTimeRange(field string, startedAt, finishedAt *time.Time) error {
	if startedAt != nil && finishedAt != nil {
		if !startedAt.Before(*finishedAt) {
			return &ValidationError{
				Field:   field,
				Message: "started_at must be before finished_at",
			}
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

// ValidateProgramTemplateEntry validates a ProgramTemplateEntry entity.
func ValidateProgramTemplateEntry(e *ProgramTemplateEntry) error {
	if e == nil {
		return &ValidationError{
			Field:   "program_template_entry",
			Message: "program_template_entry cannot be nil",
		}
	}

	if err := ValidateRequiredString("exercise_name", e.ExerciseName); err != nil {
		return err
	}
	if err := ValidateStringLength("exercise_name", e.ExerciseName, 1, 200); err != nil {
		return err
	}

	if e.Order < 0 {
		return &ValidationError{
			Field:   "order",
			Message: "order must be greater than or equal to 0",
		}
	}

	if e.Sets != nil && *e.Sets <= 0 {
		return &ValidationError{Field: "sets", Message: "sets must be greater than 0"}
	}
	if e.Reps != nil && *e.Reps <= 0 {
		return &ValidationError{Field: "reps", Message: "reps must be greater than 0"}
	}
	if e.RPE != nil {
		if err := ValidateRPE(*e.RPE); err != nil {
			return err
		}
	}
	if e.Percent1RM != nil {
		if *e.Percent1RM < 0 || *e.Percent1RM > 1 {
			return &ValidationError{
				Field:   "percent_1rm",
				Message: "percent_1rm must be between 0.0 and 1.0",
			}
		}
	}
	if e.Notes != nil {
		if err := ValidateStringLength("notes", *e.Notes, 0, 2000); err != nil {
			return err
		}
	}

	return nil
}

// ValidateProgramTemplate validates a ProgramTemplate entity.
func ValidateProgramTemplate(p *ProgramTemplate) error {
	if p == nil {
		return &ValidationError{
			Field:   "program_template",
			Message: "program_template cannot be nil",
		}
	}

	if err := ValidateRequiredString("name", p.Name); err != nil {
		return err
	}
	if err := ValidateStringLength("name", p.Name, 1, 200); err != nil {
		return err
	}
	if p.Description != nil {
		if err := ValidateStringLength("description", *p.Description, 0, 2000); err != nil {
			return err
		}
	}
	if p.Notes != nil {
		if err := ValidateStringLength("notes", *p.Notes, 0, 5000); err != nil {
			return err
		}
	}
	if p.Weeks != nil {
		if err := ValidateStringLength("weeks", *p.Weeks, 0, 50); err != nil {
			return err
		}
	}
	if p.DaysPerWeek != nil {
		if err := ValidateStringLength("days_per_week", *p.DaysPerWeek, 0, 50); err != nil {
			return err
		}
	}

	if len(p.Entries) > 1000 {
		return &ValidationError{
			Field:   "entries",
			Message: "program template cannot have more than 1000 entries",
		}
	}

	for i := range p.Entries {
		if err := ValidateProgramTemplateEntry(&p.Entries[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("entries[%d]", i),
				Message: err.Error(),
			}
		}
	}

	return nil
}

// ValidateProgramSessionEntry validates a ProgramSessionEntry entity.
func ValidateProgramSessionEntry(e *ProgramSessionEntry) error {
	if e == nil {
		return &ValidationError{
			Field:   "program_session_entry",
			Message: "program_session_entry cannot be nil",
		}
	}

	if err := ValidateRequiredString("exercise_name", e.ExerciseName); err != nil {
		return err
	}
	if err := ValidateStringLength("exercise_name", e.ExerciseName, 1, 200); err != nil {
		return err
	}

	if e.Order < 0 {
		return &ValidationError{
			Field:   "order",
			Message: "order must be greater than or equal to 0",
		}
	}

	if e.Sets != nil && *e.Sets <= 0 {
		return &ValidationError{Field: "sets", Message: "sets must be greater than 0"}
	}
	if e.Reps != nil && *e.Reps <= 0 {
		return &ValidationError{Field: "reps", Message: "reps must be greater than 0"}
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

// ValidateProgramSession validates a ProgramSession entity.
func ValidateProgramSession(s *ProgramSession) error {
	if s == nil {
		return &ValidationError{
			Field:   "program_session",
			Message: "program_session cannot be nil",
		}
	}

	if err := ValidateRequiredString("session_name", s.SessionName); err != nil {
		return err
	}
	if err := ValidateStringLength("session_name", s.SessionName, 1, 200); err != nil {
		return err
	}
	if s.Order < 0 {
		return &ValidationError{
			Field:   "order",
			Message: "order must be greater than or equal to 0",
		}
	}

	if len(s.Entries) > 1000 {
		return &ValidationError{
			Field:   "entries",
			Message: "session cannot have more than 1000 entries",
		}
	}

	for i := range s.Entries {
		if err := ValidateProgramSessionEntry(&s.Entries[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("entries[%d]", i),
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

	if err := ValidateRequiredString("name", p.Name); err != nil {
		return err
	}
	if err := ValidateStringLength("name", p.Name, 1, 200); err != nil {
		return err
	}

	if p.Status != ProgramStatusActive && p.Status != ProgramStatusCompleted && p.Status != ProgramStatusPlanned {
		return &ValidationError{
			Field:   "status",
			Message: "status must be 'active', 'completed', or 'planned'",
		}
	}

	if p.Notes != nil {
		if err := ValidateStringLength("notes", *p.Notes, 0, 5000); err != nil {
			return err
		}
	}

	if len(p.Sessions) == 0 {
		return &ValidationError{
			Field:   "sessions",
			Message: "program must have at least one session",
		}
	}
	if len(p.Sessions) > 100 {
		return &ValidationError{
			Field:   "sessions",
			Message: "program cannot have more than 100 sessions",
		}
	}

	seenSessions := make(map[string]struct{}, len(p.Sessions))
	for i := range p.Sessions {
		if err := ValidateProgramSession(&p.Sessions[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("sessions[%d]", i),
				Message: err.Error(),
			}
		}
		s := &p.Sessions[i]
		if _, exists := seenSessions[s.SessionName]; exists {
			return &ValidationError{
				Field:   fmt.Sprintf("sessions[%d]", i),
				Message: fmt.Sprintf("duplicate session_name: %s", s.SessionName),
			}
		}
		seenSessions[s.SessionName] = struct{}{}
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

	if err := ValidateRequiredString("exercise_name", e.ExerciseName); err != nil {
		return err
	}
	if err := ValidateStringLength("exercise_name", e.ExerciseName, 1, 200); err != nil {
		return err
	}

	if e.Order < 0 {
		return &ValidationError{
			Field:   "order",
			Message: "order must be greater than or equal to 0",
		}
	}

	if e.Sets != nil && *e.Sets <= 0 {
		return &ValidationError{Field: "sets", Message: "sets must be greater than 0"}
	}
	if e.Reps != nil && *e.Reps <= 0 {
		return &ValidationError{Field: "reps", Message: "reps must be greater than 0"}
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
	if err := ValidateTimeRange("started_at", e.StartedAt, e.FinishedAt); err != nil {
		return err
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

	if err := ValidateTimestamp(l.PerformedAt); err != nil {
		return err
	}
	if err := ValidateTimeRange("started_at", l.StartedAt, l.FinishedAt); err != nil {
		return err
	}

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

	for i := range l.Entries {
		if err := ValidateLogEntry(&l.Entries[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("entries[%d]", i),
				Message: err.Error(),
			}
		}
	}

	if l.Notes != nil {
		if err := ValidateStringLength("notes", *l.Notes, 0, 5000); err != nil {
			return err
		}
	}

	if l.SessionName != nil {
		if err := ValidateStringLength("session_name", *l.SessionName, 0, 200); err != nil {
			return err
		}
	}

	return nil
}
