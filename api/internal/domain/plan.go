package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PlanEntry represents a single exercise prescription entry within a training plan.
type PlanEntry struct {
	ID           uuid.UUID       `json:"id"`
	PlanID       uuid.UUID       `json:"plan_id"`
	Order        int             `json:"order"`
	ExerciseName string          `json:"exercise_name"`
	Sets         *int            `json:"sets,omitempty"`
	Reps         *int            `json:"reps,omitempty"`
	LoadKg       *float64        `json:"load_kg,omitempty"`
	RPE          *int            `json:"rpe,omitempty"`
	Notes        *string         `json:"notes,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}

// Plan represents a training plan, typically created by AI.
type Plan struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Notes       *string         `json:"notes,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Entries     []PlanEntry     `json:"entries,omitempty"`
}
