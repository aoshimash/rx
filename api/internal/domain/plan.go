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

// PlanEntry represents a single exercise prescription entry within a training plan.
type PlanEntry struct {
	ID           uuid.UUID       `json:"id"`
	PlanID       uuid.UUID       `json:"plan_id"`
	Order        int             `json:"order"`
	Date         *DateOnly       `json:"date,omitempty"`
	ExerciseName string          `json:"exercise_name"`
	Sets         *int            `json:"sets,omitempty"`
	Reps         *int            `json:"reps,omitempty"`
	LoadKg       *float64        `json:"load_kg,omitempty"`
	RPE          *int            `json:"rpe,omitempty"`
	Notes        *string         `json:"notes,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}

// Plan represents a concrete training schedule, optionally derived from a Program.
type Plan struct {
	ID          uuid.UUID       `json:"id"`
	ProgramID   *uuid.UUID      `json:"program_id,omitempty"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Notes       *string         `json:"notes,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Entries     []PlanEntry     `json:"entries,omitempty"`
}
