package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DateOnly is a date without time, formatted as "2006-01-02" in JSON.
type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

// PlanSessionEntry represents an exercise prescription within a plan session.
type PlanSessionEntry struct {
	ID           uuid.UUID              `json:"id"`
	SessionID    uuid.UUID              `json:"session_id"`
	Order        int                    `json:"order"`
	ExerciseName string                 `json:"exercise_name"`
	Fields       map[string]interface{} `json:"fields,omitempty"`
	Notes        *string                `json:"notes,omitempty"`
}

// PlanSession represents a single session in the user's execution plan.
type PlanSession struct {
	ID              uuid.UUID          `json:"id"`
	PlanID          uuid.UUID          `json:"plan_id"`
	FieldGroupID    *uuid.UUID         `json:"field_group_id,omitempty"`
	SessionName     string             `json:"session_name"`
	Order           int                `json:"order"`
	Date            *DateOnly          `json:"date,omitempty"`
	SourceProgramID *uuid.UUID         `json:"source_program_id,omitempty"`
	SourceSessionID *uuid.UUID         `json:"source_session_id,omitempty"`
	Entries         []PlanSessionEntry `json:"entries,omitempty"`
}

// Plan represents the user's working execution queue of upcoming sessions.
// Exactly one Plan per user. Persistent but disposable — no lifecycle constraints.
type Plan struct {
	ID        uuid.UUID     `json:"id"`
	Name      *string       `json:"name,omitempty"`
	Notes     *string       `json:"notes,omitempty"`
	Sessions  []PlanSession `json:"sessions,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
