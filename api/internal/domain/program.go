package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ProgramEntry represents a flat training prescription entry in a program.
// It carries a free-form metadata field for contextual grouping (e.g., week/day).
type ProgramEntry struct {
	ID        uuid.UUID       `json:"id"`
	ProgramID uuid.UUID       `json:"program_id"`
	Name      string          `json:"name"`
	Order     int             `json:"order"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	// Prescription fields
	ExerciseID         *uuid.UUID `json:"exercise_id,omitempty"`
	TargetSets         *int       `json:"target_sets,omitempty"`
	TargetReps         *int       `json:"target_reps,omitempty"`
	TargetRPE          *int       `json:"target_rpe,omitempty"`
	Percent1RM         *float64   `json:"percent_1rm,omitempty"`
	PlannedRestSeconds *int       `json:"planned_rest_seconds,omitempty"`
	MuscleGroups       []string   `json:"muscle_groups,omitempty"`
	Notes              *string    `json:"notes,omitempty"`
}

// Program represents a training program containing a flat list of entries.
type Program struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Entries     []ProgramEntry `json:"entries,omitempty"`
}
