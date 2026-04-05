package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

// RoundLoad rounds a weight value to 0.1kg precision.
func RoundLoad(kg float64) float64 {
	return math.Round(kg*10) / 10
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

// ValidateProgramGroup validates a ProgramGroup entity.
func ValidateProgramGroup(g *ProgramGroup) error {
	if g == nil {
		return &ValidationError{
			Field:   "program_group",
			Message: "program_group cannot be nil",
		}
	}

	if err := ValidateRequiredString("name", g.Name); err != nil {
		return err
	}
	if err := ValidateStringLength("name", g.Name, 1, 200); err != nil {
		return err
	}

	if g.Order < 0 {
		return &ValidationError{
			Field:   "order",
			Message: "order must be greater than or equal to 0",
		}
	}

	if g.Notes != nil {
		if err := ValidateStringLength("notes", *g.Notes, 0, 5000); err != nil {
			return err
		}
	}

	return nil
}

// ValidateGroupDepths walks the parent chains of all groups and rejects any
// group whose nesting depth reaches MaxGroupDepth or beyond.
// Depth 0 = top-level (no parent), depth 1 = child of a top-level group, etc.
func ValidateGroupDepths(groups []ProgramGroup) error {
	byID := make(map[uuid.UUID]*ProgramGroup, len(groups))
	for i := range groups {
		byID[groups[i].ID] = &groups[i]
	}

	for i := range groups {
		depth := 0
		current := &groups[i]
		visited := make(map[uuid.UUID]bool)
		for current.ParentGroupID != nil {
			if visited[current.ID] {
				return &ValidationError{
					Field:   fmt.Sprintf("groups[%d]", i),
					Message: "circular parent reference detected",
				}
			}
			visited[current.ID] = true
			parent, ok := byID[*current.ParentGroupID]
			if !ok {
				return &ValidationError{
					Field:   fmt.Sprintf("groups[%d]", i),
					Message: fmt.Sprintf("parent_group_id %s not found", current.ParentGroupID),
				}
			}
			depth++
			current = parent
		}
		if depth >= MaxGroupDepth {
			return &ValidationError{
				Field:   fmt.Sprintf("groups[%d]", i),
				Message: fmt.Sprintf("group nesting depth %d exceeds maximum allowed depth of %d", depth, MaxGroupDepth-1),
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

	if v, ok := e.Fields["sets"]; ok {
		switch val := v.(type) {
		case float64:
			if val <= 0 {
				return &ValidationError{Field: "fields.sets", Message: "sets must be greater than 0"}
			}
		case int:
			if val <= 0 {
				return &ValidationError{Field: "fields.sets", Message: "sets must be greater than 0"}
			}
		}
	}
	if v, ok := e.Fields["reps"]; ok {
		switch val := v.(type) {
		case float64:
			if val <= 0 {
				return &ValidationError{Field: "fields.reps", Message: "reps must be greater than 0"}
			}
		case int:
			if val <= 0 {
				return &ValidationError{Field: "fields.reps", Message: "reps must be greater than 0"}
			}
		}
	}
	if v, ok := e.Fields["load_kg"]; ok {
		if val, ok := v.(float64); ok {
			if val < 0 {
				return &ValidationError{
					Field:   "fields.load_kg",
					Message: "load_kg must be greater than or equal to 0",
				}
			}
			e.Fields["load_kg"] = RoundLoad(val)
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

// allowedFieldDefTypes is the set of valid FieldDef type values.
var allowedFieldDefTypes = map[string]struct{}{
	"text":   {},
	"number": {},
	"select": {},
	"video":  {},
}

// ValidateFieldDef validates a FieldDef entry.
func ValidateFieldDef(prefix string, f *FieldDef) error {
	if f == nil {
		return &ValidationError{
			Field:   prefix,
			Message: "field_def cannot be nil",
		}
	}
	if _, ok := allowedFieldDefTypes[f.Type]; !ok {
		return &ValidationError{
			Field:   prefix + ".type",
			Message: fmt.Sprintf("type must be one of: text, number, select, video; got %q", f.Type),
		}
	}
	if f.Description != "" {
		if err := ValidateStringLength(prefix+".description", f.Description, 0, 500); err != nil {
			return err
		}
	}
	return nil
}

// ValidateFieldGroup validates a FieldGroup entity.
func ValidateFieldGroup(fg *FieldGroup) error {
	if fg == nil {
		return &ValidationError{
			Field:   "field_group",
			Message: "field_group cannot be nil",
		}
	}

	if err := ValidateRequiredString("name", fg.Name); err != nil {
		return err
	}
	if err := ValidateStringLength("name", fg.Name, 1, 200); err != nil {
		return err
	}

	if fg.Description != nil {
		if err := ValidateStringLength("description", *fg.Description, 0, 2000); err != nil {
			return err
		}
	}

	if len(fg.ProgramFields) == 0 {
		return &ValidationError{
			Field:   "program_fields",
			Message: "at least one program field is required",
		}
	}
	for i, f := range fg.ProgramFields {
		if err := ValidateFieldDef(fmt.Sprintf("program_fields[%d]", i), &f); err != nil {
			return err
		}
		if f.Type == "video" {
			return &ValidationError{
				Field:   fmt.Sprintf("program_fields[%d].type", i),
				Message: "video type is not allowed in program_fields",
			}
		}
	}

	if len(fg.LogFields) == 0 {
		return &ValidationError{
			Field:   "log_fields",
			Message: "at least one log field is required",
		}
	}
	for i, f := range fg.LogFields {
		if err := ValidateFieldDef(fmt.Sprintf("log_fields[%d]", i), &f); err != nil {
			return err
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

	if p.Notes != nil {
		if err := ValidateStringLength("notes", *p.Notes, 0, 5000); err != nil {
			return err
		}
	}

	switch p.Status {
	case ProgramStatusDraft, ProgramStatusPublished:
		// valid
	default:
		return &ValidationError{
			Field:   "status",
			Message: "status must be 'draft' or 'published'",
		}
	}

	for i := range p.Groups {
		if err := ValidateProgramGroup(&p.Groups[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("groups[%d]", i),
				Message: err.Error(),
			}
		}
	}
	if err := ValidateGroupDepths(p.Groups); err != nil {
		return err
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

	groupIDs := make(map[uuid.UUID]struct{}, len(p.Groups))
	for _, g := range p.Groups {
		groupIDs[g.ID] = struct{}{}
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
		if s.GroupID != nil {
			if _, ok := groupIDs[*s.GroupID]; !ok {
				return &ValidationError{
					Field:   fmt.Sprintf("sessions[%d].group_id", i),
					Message: fmt.Sprintf("group_id %s not found in program groups", s.GroupID),
				}
			}
		}
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

// ValidatePlanSessionEntry validates a PlanSessionEntry entity.
func ValidatePlanSessionEntry(e *PlanSessionEntry) error {
	if e == nil {
		return &ValidationError{
			Field:   "plan_session_entry",
			Message: "plan_session_entry cannot be nil",
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

	if e.Notes != nil {
		if err := ValidateStringLength("notes", *e.Notes, 0, 2000); err != nil {
			return err
		}
	}

	return nil
}

// ValidatePlanSession validates a PlanSession entity.
func ValidatePlanSession(s *PlanSession) error {
	if s == nil {
		return &ValidationError{
			Field:   "plan_session",
			Message: "plan_session cannot be nil",
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
		if err := ValidatePlanSessionEntry(&s.Entries[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("entries[%d]", i),
				Message: err.Error(),
			}
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

	if p.Name != nil {
		if err := ValidateStringLength("name", *p.Name, 1, 200); err != nil {
			return err
		}
	}

	if p.Notes != nil {
		if err := ValidateStringLength("notes", *p.Notes, 0, 5000); err != nil {
			return err
		}
	}

	if len(p.Sessions) > 200 {
		return &ValidationError{
			Field:   "sessions",
			Message: "plan cannot have more than 200 sessions",
		}
	}

	for i := range p.Sessions {
		if err := ValidatePlanSession(&p.Sessions[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("sessions[%d]", i),
				Message: err.Error(),
			}
		}
	}

	return nil
}

// ValidateLogSet validates a LogSet entity.
func ValidateLogSet(s *LogSet) error {
	if s == nil {
		return &ValidationError{
			Field:   "log_set",
			Message: "log_set cannot be nil",
		}
	}

	if s.SetNumber < 1 {
		return &ValidationError{
			Field:   "set_number",
			Message: "set_number must be greater than or equal to 1",
		}
	}

	if s.Fields == nil {
		return &ValidationError{
			Field:   "fields",
			Message: "fields is required",
		}
	}

	if s.VideoObjectKey != nil {
		if err := ValidateStringLength("video_object_key", *s.VideoObjectKey, 1, 500); err != nil {
			return err
		}
	}

	if s.Notes != nil {
		if err := ValidateStringLength("notes", *s.Notes, 0, 2000); err != nil {
			return err
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

	if v, ok := e.Fields["sets"]; ok {
		switch val := v.(type) {
		case float64:
			if val <= 0 {
				return &ValidationError{Field: "fields.sets", Message: "sets must be greater than 0"}
			}
		case int:
			if val <= 0 {
				return &ValidationError{Field: "fields.sets", Message: "sets must be greater than 0"}
			}
		}
	}
	if v, ok := e.Fields["reps"]; ok {
		switch val := v.(type) {
		case float64:
			if val <= 0 {
				return &ValidationError{Field: "fields.reps", Message: "reps must be greater than 0"}
			}
		case int:
			if val <= 0 {
				return &ValidationError{Field: "fields.reps", Message: "reps must be greater than 0"}
			}
		}
	}
	if v, ok := e.Fields["load_kg"]; ok {
		if val, ok := v.(float64); ok {
			if val < 0 {
				return &ValidationError{
					Field:   "fields.load_kg",
					Message: "load_kg must be greater than or equal to 0",
				}
			}
			e.Fields["load_kg"] = RoundLoad(val)
		}
	}
	if e.Notes != nil {
		if err := ValidateStringLength("notes", *e.Notes, 0, 2000); err != nil {
			return err
		}
	}
	if err := ValidateTimeRange("started_at", e.StartedAt, e.FinishedAt); err != nil {
		return err
	}

	for i := range e.Sets {
		if err := ValidateLogSet(&e.Sets[i]); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("sets[%d]", i),
				Message: err.Error(),
			}
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
